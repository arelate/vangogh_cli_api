package vets

import (
	"fmt"
	"github.com/arelate/gog_atu"
	"github.com/arelate/vangogh_api/cli/expand"
	"github.com/arelate/vangogh_downloads"
	"github.com/arelate/vangogh_products"
	"github.com/arelate/vangogh_properties"
	"github.com/arelate/vangogh_values"
	"github.com/boggydigital/gost"
	"github.com/boggydigital/nod"
)

func UnresolvedManualUrls(
	mt gog_atu.Media,
	operatingSystems []vangogh_downloads.OperatingSystem,
	downloadTypes []vangogh_downloads.DownloadType,
	langCodes []string,
	fix bool) error {

	cumu := nod.NewProgress("checking unresolved manual-urls...")
	defer cumu.End()

	rxa, err := vangogh_properties.ConnectReduxAssets(
		vangogh_properties.TitleProperty,
		vangogh_properties.NativeLanguageNameProperty,
		vangogh_properties.LocalManualUrl)
	if err != nil {
		return cumu.EndWithError(err)
	}

	vrDetails, err := vangogh_values.NewReader(vangogh_products.Details, mt)
	if err != nil {
		return cumu.EndWithError(err)
	}

	allDetails := vrDetails.Keys()
	unresolvedIds := gost.NewStrSet()

	cumu.TotalInt(len(allDetails))
	for _, id := range allDetails {

		det, err := vrDetails.Details(id)
		if err != nil {
			cumu.Error(err)
			cumu.Increment()
			continue
		}

		downloadsList, err := vangogh_downloads.FromDetails(det, mt, rxa)
		if err != nil {
			cumu.Error(err)
			cumu.Increment()
			continue
		}

		downloadsList = downloadsList.Only(operatingSystems, downloadTypes, langCodes)

		for _, dl := range downloadsList {
			if _, ok := rxa.GetFirstVal(vangogh_properties.LocalManualUrl, dl.ManualUrl); !ok {
				unresolvedIds.Add(id)
			}
		}

		cumu.Increment()
	}

	if unresolvedIds.Len() == 0 {
		cumu.EndWithResult("all good")
	} else {

		summary, err := expand.IdsToPropertyLists(
			unresolvedIds,
			nil,
			[]string{vangogh_properties.TitleProperty},
			rxa)

		heading := fmt.Sprintf("found %d problems:", unresolvedIds.Len())
		if fix {
			heading = "found problems (run get-downloads to fix):"
		}

		if err != nil {
			return cumu.EndWithError(err)
		}
		cumu.EndWithSummary(heading, summary)
	}

	return nil
}
