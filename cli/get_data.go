package cli

import (
	"fmt"
	"github.com/arelate/gog_auth"
	"github.com/arelate/gog_media"
	"github.com/arelate/vangogh_api/cli/fetch"
	"github.com/arelate/vangogh_api/cli/itemize"
	"github.com/arelate/vangogh_api/cli/split"
	"github.com/arelate/vangogh_api/cli/url_helpers"
	"github.com/arelate/vangogh_products"
	"github.com/arelate/vangogh_urls"
	"github.com/boggydigital/coost"
	"github.com/boggydigital/gost"
	"github.com/boggydigital/nod"
	"net/url"
	"time"
)

func GetDataHandler(u *url.URL) error {
	idSet, err := url_helpers.IdSet(u)
	if err != nil {
		return err
	}

	skipIds := vangogh_urls.UrlValues(u, "skip-id")

	updated := vangogh_urls.UrlFlag(u, "updated")
	since := time.Now().Unix()
	if updated {
		since = time.Now().Add(-time.Hour * 24).Unix()
	}

	return GetData(
		idSet,
		skipIds,
		vangogh_urls.UrlProductType(u),
		vangogh_urls.UrlMedia(u),
		since,
		vangogh_urls.UrlValue(u, "temp-directory"),
		vangogh_urls.UrlFlag(u, "missing"),
		updated)
}

//GetData gets remote data from GOG.com and stores as local products (splitting as paged data if needed)
func GetData(
	idSet gost.StrSet,
	skipIds []string,
	pt vangogh_products.ProductType,
	mt gog_media.Media,
	since int64,
	tempDir string,
	missing bool,
	updated bool) error {

	gda := nod.NewProgress("getting %s (%s) data...", pt, mt)
	defer gda.End()

	if !vangogh_products.Valid(pt) {
		gda.EndWithResult("%s is not a valid product type", pt)
		return nil
	}

	if !vangogh_products.SupportsMedia(pt, mt) {
		gda.EndWithResult("%s is not a supported media for %s", mt, pt)
		return nil
	}

	cj, err := coost.NewJar(gogHosts, tempDir)
	if err != nil {
		return gda.EndWithError(err)
	}

	hc := cj.NewHttpClient()

	if vangogh_products.RequiresAuth(pt) {
		li, err := gog_auth.LoggedIn(hc)
		if err != nil {
			return gda.EndWithError(err)
		}

		if !li {
			return gda.EndWithError(fmt.Errorf("user is not logged in"))
		}
	}

	if vangogh_products.IsPaged(pt) {
		if err := fetch.Pages(pt, mt, hc, gda); err != nil {
			return gda.EndWithError(err)
		}
		return split.Pages(pt, mt, since)
	}

	if vangogh_products.IsArray(pt) {
		// using "licences" as id, since that's how we store that data in kvas
		ids := []string{vangogh_products.Licences.String()}
		if err := fetch.Items(ids, pt, mt, hc); err != nil {
			return gda.EndWithError(err)
		}
		return split.Pages(pt, mt, since)
	}

	idSet, err = itemize.All(idSet, missing, updated, since, pt, mt)
	if err != nil {
		return gda.EndWithError(err)
	}

	approvedIds := idSet.Except(gost.NewStrSetWith(skipIds...))

	return fetch.Items(approvedIds, pt, mt, hc)
}
