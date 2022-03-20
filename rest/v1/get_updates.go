package v1

import (
	"fmt"
	"github.com/arelate/gog_integration"
	"github.com/arelate/vangogh_local_data"
	"github.com/boggydigital/nod"
	"net/http"
)

var filterNewProductTypes = map[vangogh_local_data.ProductType]bool{
	vangogh_local_data.Orders: true,
	//not all licence-products have associated api-products-v1/api-products-v2,
	//so in some cases we won't get a meaningful information like a title
	vangogh_local_data.LicenceProducts: true,
	//both ApiProductsVx are not interesting since they correspond to store-products or account-products
	vangogh_local_data.ApiProductsV1: true,
	vangogh_local_data.ApiProductsV2: true,
}

var filterUpdatedProductTypes = map[vangogh_local_data.ProductType]bool{
	//most of the Updates are price changes for a sale, not that interesting for recurring sync
	vangogh_local_data.StoreProducts: true,
	// wishlist-products are basically store-products, so see above
	vangogh_local_data.WishlistProducts: true,
	//meaningful Updates for account products come from details, not account-products
	vangogh_local_data.AccountProducts: true,
	//same as above for those product types
	vangogh_local_data.ApiProductsV1: true,
	vangogh_local_data.ApiProductsV2: true,
}

func GetUpdates(w http.ResponseWriter, r *http.Request) {

	// GET /v1/updates?media&since&format

	if r.Method != http.MethodGet {
		err := fmt.Errorf("unsupported method")
		http.Error(w, nod.Error(err).Error(), http.StatusMethodNotAllowed)
		return
	}

	_, mt, err := productTypeMediaFromUrl(r.URL)
	if err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusMethodNotAllowed)
		return
	}

	since, err := vangogh_local_data.SinceFromUrl(r.URL)
	if err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusMethodNotAllowed)
		return
	}

	updates, err := Updates(mt, since)
	if err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusMethodNotAllowed)
		return
	}

	if err := encode(updates, w, r); err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusMethodNotAllowed)
		return
	}
}

func Updates(mt gog_integration.Media, since int64) (map[string]map[string]bool, error) {
	updates := make(map[string]map[string]bool, 0)

	for _, pt := range vangogh_local_data.LocalProducts() {

		if filterNewProductTypes[pt] {
			continue
		}

		vr, err := vangogh_local_data.NewReader(pt, mt)
		if err != nil {
			return updates, err
		}

		categorize(vr.CreatedAfter(since),
			fmt.Sprintf("new in %s", pt.HumanReadableString()),
			updates)

		if filterUpdatedProductTypes[pt] {
			continue
		}

		categorize(vr.ModifiedAfter(since, true),
			fmt.Sprintf("updated in %s", pt.HumanReadableString()),
			updates)
	}

	return updates, nil
}

func categorize(ids []string, cat string, updates map[string]map[string]bool) {
	for _, id := range ids {
		if updates[cat] == nil {
			updates[cat] = make(map[string]bool, 0)
		}
		updates[cat][id] = true
	}
}
