package vets

import (
	"fmt"
	"github.com/arelate/gog_atu"
	"github.com/arelate/vangogh_data"
	"github.com/boggydigital/nod"
)

func UnresolvedManualUrls(
	mt gog_atu.Media,
	operatingSystems []vangogh_data.OperatingSystem,
	downloadTypes []vangogh_data.DownloadType,
	langCodes []string,
	fix bool) error {

	cumu := nod.NewProgress("checking unresolved manual-urls...")
	defer cumu.End()

	rxa, err := vangogh_data.ConnectReduxAssets(
		vangogh_data.TitleProperty,
		vangogh_data.NativeLanguageNameProperty,
		vangogh_data.LocalManualUrlProperty)
	if err != nil {
		return cumu.EndWithError(err)
	}

	vrDetails, err := vangogh_data.NewReader(vangogh_data.Details, mt)
	if err != nil {
		return cumu.EndWithError(err)
	}

	allDetails := vrDetails.Keys()
	unresolvedIds := vangogh_data.NewIdSet()

	cumu.TotalInt(len(allDetails))
	for _, id := range allDetails {

		det, err := vrDetails.Details(id)
		if err != nil {
			cumu.Error(err)
			cumu.Increment()
			continue
		}

		downloadsList, err := vangogh_data.FromDetails(det, mt, rxa)
		if err != nil {
			cumu.Error(err)
			cumu.Increment()
			continue
		}

		downloadsList = downloadsList.Only(operatingSystems, downloadTypes, langCodes)

		for _, dl := range downloadsList {
			if _, ok := rxa.GetFirstVal(vangogh_data.LocalManualUrlProperty, dl.ManualUrl); !ok {
				unresolvedIds.Add(id)
			}
		}

		cumu.Increment()
	}

	if unresolvedIds.Len() == 0 {
		cumu.EndWithResult("all good")
	} else {

		summary, err := vangogh_data.PropertyListsFromIdSet(
			unresolvedIds,
			nil,
			[]string{vangogh_data.TitleProperty},
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
