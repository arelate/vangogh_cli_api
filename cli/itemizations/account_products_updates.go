package itemizations

import (
	"github.com/arelate/gog_integration"
	"github.com/arelate/vangogh_local_data"
	"github.com/boggydigital/nod"
	"strconv"
)

func AccountProductsUpdates(mt gog_integration.Media) (*vangogh_local_data.IdSet, error) {

	apua := nod.Begin(" finding %s updates...", vangogh_local_data.AccountProducts)
	defer apua.End()

	updatesSet := vangogh_local_data.NewIdSet()
	vrAccountPages, err := vangogh_local_data.NewReader(vangogh_local_data.AccountPage, mt)
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
				nod.Log("%s #%d Updates, isNew: %d, %v", vangogh_local_data.AccountProducts, ap.Id, ap.Updates, ap.IsNew)
				updatesSet.Add(strconv.Itoa(ap.Id))
			}
		}
	}

	apua.EndWithResult(itemizationResult(updatesSet))

	return updatesSet, nil
}
