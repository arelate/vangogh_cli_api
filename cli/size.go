package cli

import (
	"github.com/arelate/gog_integration"
	"github.com/arelate/vangogh_cli_api/cli/itemizations"
	"github.com/arelate/vangogh_local_data"
	"github.com/boggydigital/nod"
	"net/url"
)

func SizeHandler(u *url.URL) error {
	idSet, err := vangogh_local_data.IdSetFromUrl(u)
	if err != nil {
		return err
	}

	return Size(
		idSet,
		vangogh_local_data.MediaFromUrl(u),
		vangogh_local_data.OperatingSystemsFromUrl(u),
		vangogh_local_data.DownloadTypesFromUrl(u),
		vangogh_local_data.ValuesFromUrl(u, "language-code"),
		vangogh_local_data.FlagFromUrl(u, "missing"),
		vangogh_local_data.FlagFromUrl(u, "all"))
}

func Size(
	idSet vangogh_local_data.IdSet,
	mt gog_integration.Media,
	operatingSystems []vangogh_local_data.OperatingSystem,
	downloadTypes []vangogh_local_data.DownloadType,
	langCodes []string,
	missing bool,
	all bool) error {

	sa := nod.NewProgress("estimating downloads size...")
	defer sa.End()

	rxa, err := vangogh_local_data.ConnectReduxAssets(
		vangogh_local_data.LocalManualUrlProperty,
		vangogh_local_data.NativeLanguageNameProperty,
		vangogh_local_data.SlugProperty,
		vangogh_local_data.DownloadStatusErrorProperty)
	if err != nil {
		return sa.EndWithError(err)
	}

	if missing {
		missingIds, err := itemizations.MissingLocalDownloads(mt, rxa, operatingSystems, downloadTypes, langCodes)
		if err != nil {
			return sa.EndWithError(err)
		}

		if missingIds.Len() == 0 {
			sa.EndWithResult("no missing downloads")
			return nil
		}

		idSet.AddSet(missingIds)
	}

	if all {
		vrDetails, err := vangogh_local_data.NewReader(vangogh_local_data.Details, mt)
		if err != nil {
			return sa.EndWithError(err)
		}
		idSet.Add(vrDetails.Keys()...)
	}

	if idSet.Len() == 0 {
		sa.EndWithResult("no ids to estimate size")
		return nil
	}

	sd := &sizeDelegate{}

	sa.TotalInt(idSet.Len())

	if err := vangogh_local_data.MapDownloads(
		idSet,
		mt,
		rxa,
		operatingSystems,
		downloadTypes,
		langCodes,
		sd,
		sa); err != nil {
		return err
	}

	sa.EndWithResult("%.2fGB", sd.TotalGBsEstimate())

	return nil
}

type sizeDelegate struct {
	dlList vangogh_local_data.DownloadsList
}

func (sd *sizeDelegate) Process(_, _ string, list vangogh_local_data.DownloadsList) error {
	if sd.dlList == nil {
		sd.dlList = make(vangogh_local_data.DownloadsList, 0)
	}
	sd.dlList = append(sd.dlList, list...)
	return nil
}

func (sd *sizeDelegate) TotalGBsEstimate() float64 {
	if sd.dlList != nil {
		return sd.dlList.TotalGBsEstimate()
	}
	return 0
}
