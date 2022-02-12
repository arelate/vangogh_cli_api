package cli

import (
	"github.com/arelate/gog_atu"
	"github.com/arelate/vangogh_data"
	"github.com/boggydigital/coost"
	"github.com/boggydigital/nod"
	"net/http"
	"net/url"
	"path/filepath"
)

func WishlistHandler(u *url.URL) error {
	return Wishlist(
		vangogh_data.MediaFromUrl(u),
		vangogh_data.ValuesFromUrl(u, "add"),
		vangogh_data.ValuesFromUrl(u, "remove"),
		vangogh_data.ValueFromUrl(u, "temp-directory"))
}

func Wishlist(mt gog_atu.Media, addProductIds, removeProductIds []string, tempDir string) error {

	wa := nod.Begin("performing requested wishlist operations...")
	defer wa.End()

	hc, err := coost.NewHttpClientFromFile(
		filepath.Join(tempDir, cookiesFilename), gog_atu.GogHost)
	if err != nil {
		return wa.EndWithError(err)
	}

	vrStoreProducts, err := vangogh_data.NewReader(vangogh_data.StoreProducts, mt)
	if err != nil {
		return wa.EndWithError(err)
	}

	if len(addProductIds) > 0 {
		if err := wishlistAdd(addProductIds, hc, vrStoreProducts, mt); err != nil {
			return wa.EndWithError(err)
		}
	}

	if len(removeProductIds) > 0 {
		if err := wishlistRemove(removeProductIds, hc, vrStoreProducts, mt); err != nil {
			return wa.EndWithError(err)
		}
	}

	wa.EndWithResult("done")

	return nil
}

func wishlistAdd(
	ids []string,
	httpClient *http.Client,
	vrStoreProducts *vangogh_data.ValueReader,
	mt gog_atu.Media) error {

	waa := nod.NewProgress(" adding product(s) to local wishlist...")
	defer waa.End()

	waa.TotalInt(len(ids))

	for _, id := range ids {
		if err := vrStoreProducts.CopyToType(id, vangogh_data.WishlistProducts, mt); err != nil {
			return waa.EndWithError(err)
		}
		waa.Increment()
	}

	waa.EndWithResult("done")

	return remoteWishlistCommand(
		ids,
		gog_atu.AddToWishlistUrl,
		httpClient,
		vrStoreProducts)
}

func wishlistRemove(
	ids []string,
	httpClient *http.Client,
	vrStoreProducts *vangogh_data.ValueReader,
	mt gog_atu.Media) error {

	wra := nod.NewProgress(" removing product(s) from local wishlist...")
	defer wra.End()

	if err := vangogh_data.Cut(ids, vangogh_data.WishlistProducts, mt); err != nil {
		return wra.EndWithError(err)
	}

	wra.EndWithResult("done")

	return remoteWishlistCommand(
		ids,
		gog_atu.RemoveFromWishlistUrl,
		httpClient,
		vrStoreProducts)
}

func remoteWishlistCommand(
	ids []string,
	wishlistUrl func(string) *url.URL,
	httpClient *http.Client,
	vrStoreProducts *vangogh_data.ValueReader) error {

	rwca := nod.NewProgress(" syncing to remote wishlist...")
	defer rwca.End()

	rwca.TotalInt(len(ids))

	for _, id := range ids {
		if !vrStoreProducts.Has(id) {
			continue
		}
		wUrl := wishlistUrl(id)
		resp, err := httpClient.Get(wUrl.String())
		if err != nil {
			return rwca.EndWithError(err)
		}

		if err := resp.Body.Close(); err != nil {
			return rwca.EndWithError(err)
		}

		rwca.Increment()
	}

	rwca.EndWithResult("done")

	return nil
}
