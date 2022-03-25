package reductions

import (
	"github.com/arelate/gog_integration"
	"github.com/arelate/vangogh_local_data"
	"github.com/boggydigital/nod"
)

func Wishlisted(mt gog_integration.Media) error {

	wa := nod.Begin(" %s...", vangogh_local_data.WishlistedProperty)
	defer wa.End()

	vrStoreProducts, err := vangogh_local_data.NewReader(vangogh_local_data.StoreProducts, mt)
	if err != nil {
		return wa.EndWithError(err)
	}

	vrWishlisted, err := vangogh_local_data.NewReader(vangogh_local_data.WishlistProducts, mt)
	if err != nil {
		return wa.EndWithError(err)
	}

	wishlisted := map[string][]string{}

	for _, id := range vrStoreProducts.Keys() {
		wishlisted[id] = []string{"false"}
	}

	for _, id := range vrWishlisted.Keys() {
		wishlisted[id] = []string{"true"}
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
