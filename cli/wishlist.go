package cli

import (
	"github.com/arelate/gog_integration"
	"github.com/arelate/vangogh_local_data"
	"github.com/boggydigital/coost"
	"github.com/boggydigital/kvas"
	"github.com/boggydigital/nod"
	"net/http"
	"net/url"
)

func WishlistHandler(u *url.URL) error {
	return Wishlist(
		vangogh_local_data.MediaFromUrl(u),
		vangogh_local_data.ValuesFromUrl(u, "add"),
		vangogh_local_data.ValuesFromUrl(u, "remove"))
}

func Wishlist(mt gog_integration.Media, addProductIds, removeProductIds []string) error {

	wa := nod.Begin("performing requested wishlist operations...")
	defer wa.End()

	hc, err := coost.NewHttpClientFromFile(vangogh_local_data.AbsCookiePath(), gog_integration.GogHost)
	if err != nil {
		return wa.EndWithError(err)
	}

	vrStoreProducts, err := vangogh_local_data.NewReader(vangogh_local_data.StoreProducts, mt)
	if err != nil {
		return wa.EndWithError(err)
	}

	rxa, err := vangogh_local_data.ConnectReduxAssets(vangogh_local_data.WishlistedProperty)
	if err != nil {
		return wa.EndWithError(err)
	}

	if len(addProductIds) > 0 {
		if err := wishlistAdd(addProductIds, hc, vrStoreProducts, rxa, mt); err != nil {
			return wa.EndWithError(err)
		}
	}

	if len(removeProductIds) > 0 {
		if err := wishlistRemove(removeProductIds, hc, rxa, mt); err != nil {
			return wa.EndWithError(err)
		}
	}

	wa.EndWithResult("done")

	return nil
}

func wishlistAdd(
	ids []string,
	httpClient *http.Client,
	vrStoreProducts *vangogh_local_data.ValueReader,
	rxa kvas.ReduxAssets,
	mt gog_integration.Media) error {

	waa := nod.NewProgress(" adding product(s) to local wishlist...")
	defer waa.End()

	waa.TotalInt(len(ids))

	for _, id := range ids {

		if err := vrStoreProducts.CopyToType(id, vangogh_local_data.WishlistProducts, mt); err != nil {
			return waa.EndWithError(err)
		}

		if !rxa.HasVal(vangogh_local_data.WishlistedProperty, id, "true") {
			if err := rxa.AddVal(vangogh_local_data.WishlistedProperty, id, "true"); err != nil {
				return waa.EndWithError(err)
			}
		}

		waa.Increment()
	}

	waa.EndWithResult("done")

	return gog_integration.AddToWishlist(httpClient, ids...)
}

func wishlistRemove(
	ids []string,
	httpClient *http.Client,
	rxa kvas.ReduxAssets,
	mt gog_integration.Media) error {

	wra := nod.NewProgress(" removing product(s) from local wishlist...")
	defer wra.End()

	idSet := make(map[string]bool)
	for _, id := range ids {
		idSet[id] = true

		if err := rxa.CutVal(vangogh_local_data.WishlistedProperty, id, "true"); err != nil {
			return wra.EndWithError(err)
		}
	}

	if err := vangogh_local_data.Cut(idSet, vangogh_local_data.WishlistProducts, mt); err != nil {
		return wra.EndWithError(err)
	}

	wra.EndWithResult("done")

	return gog_integration.RemoveFromWishlist(httpClient, ids...)
}
