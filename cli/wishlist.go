package cli

import (
	"github.com/arelate/gog_integration"
	"github.com/arelate/vangogh_local_data"
	"github.com/boggydigital/coost"
	"github.com/boggydigital/nod"
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

	if len(addProductIds) > 0 {
		if processedIds, err := wishlistAdd(addProductIds, mt); err == nil {
			if err := gog_integration.AddToWishlist(hc, processedIds...); err != nil {
				return wa.EndWithError(err)
			}
		} else {
			return wa.EndWithError(err)

		}
	}

	if len(removeProductIds) > 0 {
		if processedIds, err := wishlistRemove(removeProductIds, mt); err == nil {
			if err := gog_integration.RemoveFromWishlist(hc, processedIds...); err != nil {
				return wa.EndWithError(err)
			}
		} else {
			return wa.EndWithError(err)
		}
	}

	wa.EndWithResult("done")

	return nil
}

func wishlistAdd(
	ids []string,
	mt gog_integration.Media) ([]string, error) {

	waa := nod.NewProgress(" adding product(s) to local wishlist...")
	defer waa.End()

	pids, err := vangogh_local_data.AddToLocalWishlist(ids, mt, waa)
	if err != nil {
		waa.EndWithError(err)
	} else {
		waa.EndWithResult("done")
	}

	return pids, err
}

func wishlistRemove(
	ids []string,
	mt gog_integration.Media) ([]string, error) {

	wra := nod.NewProgress(" removing product(s) from local wishlist...")
	defer wra.End()

	pids, err := vangogh_local_data.RemoveFromLocalWishlist(ids, mt, wra)
	if err != nil {
		wra.EndWithError(err)
	} else {
		wra.EndWithResult("done")
	}

	return pids, err
}
