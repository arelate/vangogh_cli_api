package cli

import (
	"github.com/arelate/gog_media"
	"github.com/arelate/gog_urls"
	"github.com/arelate/vangogh_api/cli/remove"
	"github.com/arelate/vangogh_products"
	"github.com/arelate/vangogh_urls"
	"github.com/arelate/vangogh_values"
	"github.com/boggydigital/coost"
	"github.com/boggydigital/nod"
	"net/http"
	"net/url"
	"path/filepath"
)

func WishlistHandler(u *url.URL) error {
	return Wishlist(
		vangogh_urls.UrlMedia(u),
		vangogh_urls.UrlValues(u, "add"),
		vangogh_urls.UrlValues(u, "remove"),
		vangogh_urls.UrlValue(u, "temp-directory"))
}

func Wishlist(mt gog_media.Media, addProductIds, removeProductIds []string, tempDir string) error {

	wa := nod.Begin("performing requested wishlist operations...")
	defer wa.End()

	hc, err := coost.NewHttpClientFromFile(
		filepath.Join(tempDir, cookiesFilename), gog_urls.GogHost)
	if err != nil {
		return wa.EndWithError(err)
	}

	vrStoreProducts, err := vangogh_values.NewReader(vangogh_products.StoreProducts, mt)
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
	vrStoreProducts *vangogh_values.ValueReader,
	mt gog_media.Media) error {

	waa := nod.NewProgress(" adding product(s) to local wishlist...")
	defer waa.End()

	waa.TotalInt(len(ids))

	for _, id := range ids {
		if err := vrStoreProducts.CopyToType(id, vangogh_products.WishlistProducts, mt); err != nil {
			return waa.EndWithError(err)
		}
		waa.Increment()
	}

	waa.EndWithResult("done")

	return remoteWishlistCommand(
		ids,
		gog_urls.AddToWishlist,
		httpClient,
		vrStoreProducts)
}

func wishlistRemove(
	ids []string,
	httpClient *http.Client,
	vrStoreProducts *vangogh_values.ValueReader,
	mt gog_media.Media) error {

	wra := nod.NewProgress(" removing product(s) from local wishlist...")
	defer wra.End()

	if err := remove.Data(ids, vangogh_products.WishlistProducts, mt); err != nil {
		return wra.EndWithError(err)
	}

	wra.EndWithResult("done")

	return remoteWishlistCommand(
		ids,
		gog_urls.RemoveFromWishlist,
		httpClient,
		vrStoreProducts)
}

func remoteWishlistCommand(
	ids []string,
	wishlistUrl func(string) *url.URL,
	httpClient *http.Client,
	vrStoreProducts *vangogh_values.ValueReader) error {

	rwca := nod.NewProgress(" syncing to remote wishlist...")
	defer rwca.End()

	rwca.TotalInt(len(ids))

	for _, id := range ids {
		if !vrStoreProducts.Contains(id) {
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
