package fetchers

import (
	"fmt"
	"github.com/arelate/gog_integration"
	"github.com/arelate/vangogh_local_data"
	"github.com/boggydigital/dolo"
	"github.com/boggydigital/nod"
	"net/http"
	"net/url"
)

//Items fetches all individual data items (details, api-products-v1/v2) using provided ids
func Items(
	ids []string,
	pt vangogh_local_data.ProductType,
	mt gog_integration.Media,
	httpClient *http.Client) error {

	ia := nod.NewProgress(" fetching %s (%s)...", pt, mt)
	defer ia.End()

	if !vangogh_local_data.IsGetItemsSupported(pt) {
		return ia.EndWithError(fmt.Errorf("getting %s is not supported", pt))
	}

	ia.TotalInt(len(ids))

	sourceUrl, err := vangogh_local_data.RemoteProductsUrl(pt)
	if err != nil {
		return ia.EndWithError(err)
	}

	// intercept steam-app-news sourceUrl and set to an id-transforming provider
	if sourceUrl == nil && pt == vangogh_local_data.SteamAppNews {

		rxa, err := vangogh_local_data.ConnectReduxAssets(vangogh_local_data.SteamAppIdProperty)
		if err != nil {
			return err
		}

		sanup := &SteamAppNewsUrlProvider{rxa: rxa}
		sourceUrl = sanup.DefaultSourceUrl
	}

	//since we know how many ids need to be fetched, allocate URLs and idStrs to that number
	urls, idStr := make([]*url.URL, len(ids)), make([]string, len(ids))

	for i := 0; i < len(ids); i++ {
		urls[i], idStr[i] = sourceUrl(ids[i], mt), ids[i]
	}

	kis, err := NewKvasIndexSetter(pt, mt, idStr)
	if err != nil {
		return ia.EndWithError(err)
	}

	dc := dolo.NewClient(httpClient, dolo.Defaults())

	if err := dc.GetSet(urls, kis, ia); err != nil {
		return ia.EndWithError(err)
	}

	ia.EndWithResult("done")

	return nil
}
