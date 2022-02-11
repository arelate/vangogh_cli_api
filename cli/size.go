package cli

import (
	"github.com/arelate/gog_atu"
	"github.com/arelate/vangogh_api/cli/itemize"
	"github.com/arelate/vangogh_api/cli/url_helpers"
	"github.com/arelate/vangogh_data"
	"github.com/boggydigital/gost"
	"github.com/boggydigital/nod"
	"net/url"
)

func SizeHandler(u *url.URL) error {
	idSet, err := url_helpers.IdSet(u)
	if err != nil {
		return err
	}

	return Size(
		idSet,
		vangogh_data.MediaFromUrl(u),
		vangogh_data.OperatingSystemsFromUrl(u),
		vangogh_data.DownloadTypesFromUrl(u),
		vangogh_data.ValuesFromUrl(u, "language-code"),
		vangogh_data.FlagFromUrl(u, "missing"),
		vangogh_data.FlagFromUrl(u, "all"))
}

func Size(
	idSet gost.StrSet,
	mt gog_atu.Media,
	operatingSystems []vangogh_data.OperatingSystem,
	downloadTypes []vangogh_data.DownloadType,
	langCodes []string,
	missing bool,
	all bool) error {

	sa := nod.NewProgress("estimating downloads size...")
	defer sa.End()

	rxa, err := vangogh_data.ConnectReduxAssets(
		vangogh_data.LocalManualUrl,
		vangogh_data.NativeLanguageNameProperty,
		vangogh_data.SlugProperty,
		vangogh_data.DownloadStatusError)
	if err != nil {
		return sa.EndWithError(err)
	}

	if missing {
		missingIds, err := itemize.MissingLocalDownloads(mt, rxa, operatingSystems, downloadTypes, langCodes)
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
		vrDetails, err := vangogh_data.NewReader(vangogh_data.Details, mt)
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

	if err := vangogh_data.MapDownloads(
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
	dlList vangogh_data.DownloadsList
}

func (sd *sizeDelegate) Process(_, _ string, list vangogh_data.DownloadsList) error {
	if sd.dlList == nil {
		sd.dlList = make(vangogh_data.DownloadsList, 0)
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
