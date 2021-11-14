package cli

import (
	"github.com/arelate/gog_media"
	"github.com/arelate/vangogh_api/cli/itemize"
	"github.com/arelate/vangogh_api/cli/url_helpers"
	"github.com/arelate/vangogh_downloads"
	"github.com/arelate/vangogh_extracts"
	"github.com/arelate/vangogh_products"
	"github.com/arelate/vangogh_properties"
	"github.com/arelate/vangogh_values"
	"github.com/boggydigital/gost"
	"github.com/boggydigital/nod"
	"net/url"
)

func SizeHandler(u *url.URL) error {
	idSet, err := url_helpers.IdSet(u)
	if err != nil {
		return err
	}

	mt := gog_media.Parse(url_helpers.Value(u, "media"))

	operatingSystems := url_helpers.OperatingSystems(u)
	downloadTypes := url_helpers.DownloadTypes(u)
	langCodes := url_helpers.Values(u, "language-code")

	missing := url_helpers.Flag(u, "missing")
	all := url_helpers.Flag(u, "all")

	return Size(idSet, mt, operatingSystems, downloadTypes, langCodes, missing, all)
}

func Size(
	idSet gost.StrSet,
	mt gog_media.Media,
	operatingSystems []vangogh_downloads.OperatingSystem,
	downloadTypes []vangogh_downloads.DownloadType,
	langCodes []string,
	missing bool,
	all bool) error {

	sa := nod.NewProgress("estimating downloads size...")
	defer sa.End()

	exl, err := vangogh_extracts.NewList(
		vangogh_properties.LocalManualUrl,
		vangogh_properties.NativeLanguageNameProperty,
		vangogh_properties.SlugProperty,
		vangogh_properties.DownloadStatusError)
	if err != nil {
		return sa.EndWithError(err)
	}

	if missing {
		missingIds, err := itemize.MissingLocalDownloads(mt, exl, operatingSystems, downloadTypes, langCodes)
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
		vrDetails, err := vangogh_values.NewReader(vangogh_products.Details, mt)
		if err != nil {
			return sa.EndWithError(err)
		}
		idSet.Add(vrDetails.All()...)
	}

	if idSet.Len() == 0 {
		sa.EndWithResult("no ids to estimate size")
		return nil
	}

	sd := &sizeDelegate{}

	sa.TotalInt(idSet.Len())

	if err := vangogh_downloads.Map(
		idSet,
		mt,
		exl,
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
	dlList vangogh_downloads.DownloadsList
}

func (sd *sizeDelegate) Process(_, _ string, list vangogh_downloads.DownloadsList) error {
	if sd.dlList == nil {
		sd.dlList = make(vangogh_downloads.DownloadsList, 0)
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
