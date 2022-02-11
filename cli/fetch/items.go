package fetch

import (
	"fmt"
	"github.com/arelate/gog_atu"
	"github.com/arelate/vangogh_data"
	"github.com/boggydigital/dolo"
	"github.com/boggydigital/nod"
	"net/http"
	"net/url"
)

//Items fetches all individual data items (details, api-products-v1/v2) using provided ids
func Items(
	ids []string,
	pt vangogh_data.ProductType,
	mt gog_atu.Media,
	httpClient *http.Client) error {

	ia := nod.NewProgress(" fetching %s (%s)...", pt, mt)
	defer ia.End()

	if !vangogh_data.IsGetItemsSupported(pt) {
		return ia.EndWithError(fmt.Errorf("getting %s is not supported", pt))
	}

	ia.TotalInt(len(ids))

	sourceUrl, err := vangogh_data.RemoteProductsUrl(pt)
	if err != nil {
		return ia.EndWithError(err)
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
