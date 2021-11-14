package cli

import (
	"fmt"
	"github.com/arelate/gog_auth"
	"github.com/arelate/gog_media"
	"github.com/arelate/vangogh_api/cli/fetch"
	"github.com/arelate/vangogh_api/cli/http_client"
	"github.com/arelate/vangogh_api/cli/itemize"
	"github.com/arelate/vangogh_api/cli/lines"
	"github.com/arelate/vangogh_api/cli/split"
	"github.com/arelate/vangogh_api/cli/url_helpers"
	"github.com/arelate/vangogh_products"
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

	pt := vangogh_products.Parse(url_helpers.Value(u, "product-type"))
	mt := gog_media.Parse(url_helpers.Value(u, "media"))

	denyIdsFile := url_helpers.Value(u, "deny-ids-file")
	denyIds := lines.Read(denyIdsFile)

	updated := url_helpers.Flag(u, "updated")
	since := time.Now().Unix()
	if updated {
		since = time.Now().Add(-time.Hour * 24).Unix()
	}
	missing := url_helpers.Flag(u, "missing")

	return GetData(idSet, denyIds, pt, mt, since, missing, updated)
}

//GetData gets remote data from GOG.com and stores as local products (splitting as paged data if needed)
func GetData(
	idSet gost.StrSet,
	denyIds []string,
	pt vangogh_products.ProductType,
	mt gog_media.Media,
	since int64,
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

	httpClient, err := http_client.Default()
	if err != nil {
		return gda.EndWithError(err)
	}

	if vangogh_products.RequiresAuth(pt) {
		li, err := gog_auth.LoggedIn(httpClient)
		if err != nil {
			return gda.EndWithError(err)
		}

		if !li {
			return gda.EndWithError(fmt.Errorf("user is not logged in"))
		}
	}

	if vangogh_products.IsPaged(pt) {
		if err := fetch.Pages(pt, mt, httpClient, gda); err != nil {
			return gda.EndWithError(err)
		}
		return split.Pages(pt, mt, since)
	}

	if vangogh_products.IsArray(pt) {
		// using "licences" as id, since that's how we store that data in kvas
		ids := []string{vangogh_products.Licences.String()}
		if err := fetch.Items(ids, pt, mt, httpClient); err != nil {
			return gda.EndWithError(err)
		}
		return split.Pages(pt, mt, since)
	}

	idSet, err = itemize.All(idSet, missing, updated, since, pt, mt)
	if err != nil {
		return gda.EndWithError(err)
	}

	approvedIds := idSet.Except(gost.NewStrSetWith(denyIds...))

	return fetch.Items(approvedIds, pt, mt, httpClient)
}
