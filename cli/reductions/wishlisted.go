package reductions

import (
	"github.com/arelate/gog_integration"
	"github.com/arelate/vangogh_local_data"
	"github.com/boggydigital/nod"
	"strconv"
)

func Wishlisted(mt gog_integration.Media) error {

	wa := nod.Begin(" %s...", vangogh_local_data.WishlistedProperty)
	defer wa.End()

	vrStoreProducts, err := vangogh_local_data.NewReader(vangogh_local_data.StoreProducts, mt)
	if err != nil {
		return wa.EndWithError(err)
	}

	//using WishlistPage and not WishlistProduct for the remote source of truth
	vrWishlistPages, err := vangogh_local_data.NewReader(vangogh_local_data.WishlistPage, mt)
	if err != nil {
		return wa.EndWithError(err)
	}

	wishlisted := map[string][]string{}

	//important to set all to false as a starting point to overwrite status
	//for product no longer wishlisted at the remote source of truth
	for _, id := range vrStoreProducts.Keys() {
		wishlisted[id] = []string{"false"}
	}

	for _, page := range vrWishlistPages.Keys() {
		page, err := vrWishlistPages.ProductsGetter(page)
		if err != nil {
			wa.EndWithError(err)
		}
		for _, prod := range page.GetProducts() {
			id := strconv.Itoa(prod.GetId())
			wishlisted[id] = []string{"true"}
		}
	}

	wishlistedRdx, err := vangogh_local_data.ConnectReduxAssets(vangogh_local_data.WishlistedProperty)
	if err != nil {
		return wa.EndWithError(err)
	}

	if err := wishlistedRdx.BatchReplaceValues(vangogh_local_data.WishlistedProperty, wishlisted); err != nil {
		return wa.EndWithError(err)
	}

	wa.EndWithResult("done")

	return nil
}
