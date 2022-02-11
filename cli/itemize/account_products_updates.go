package itemize

import (
	"github.com/arelate/gog_atu"
	"github.com/arelate/vangogh_data"
	"github.com/boggydigital/gost"
	"github.com/boggydigital/nod"
	"strconv"
)

func AccountProductsUpdates(mt gog_atu.Media) (gost.StrSet, error) {

	apua := nod.Begin(" finding %s updates...", vangogh_data.AccountProducts)
	defer apua.End()

	updatesSet := gost.NewStrSet()
	vrAccountPages, err := vangogh_data.NewReader(vangogh_data.AccountPage, mt)
	if err != nil {
		return updatesSet, apua.EndWithError(err)
	}

	for _, page := range vrAccountPages.Keys() {
		accountPage, err := vrAccountPages.AccountPage(page)
		if err != nil {
			return updatesSet, apua.EndWithError(err)
		}
		for _, ap := range accountPage.Products {
			if ap.Updates > 0 ||
				ap.IsNew {
				nod.Log("%s #%d Updates, isNew: %d, %v", vangogh_data.AccountProducts, ap.Id, ap.Updates, ap.IsNew)
				updatesSet.Add(strconv.Itoa(ap.Id))
			}
		}
	}

	apua.EndWithResult(itemizationResult(updatesSet))

	return updatesSet, nil
}
