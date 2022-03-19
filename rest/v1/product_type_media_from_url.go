package v1

import (
	"fmt"
	"github.com/arelate/gog_integration"
	"github.com/arelate/vangogh_local_data"
	"net/url"
)

func productTypeMediaFromUrl(u *url.URL) (vangogh_local_data.ProductType, gog_integration.Media, error) {
	q := u.Query()

	productType := q.Get("product-type")
	pt := vangogh_local_data.ParseProductType(productType)
	if pt == vangogh_local_data.UnknownProductType {
		return pt, gog_integration.Unknown, fmt.Errorf("unknown product-type %s", productType)
	}

	media := q.Get("media")
	mt := gog_integration.ParseMedia(media)
	if mt == gog_integration.Unknown {
		return pt, mt, fmt.Errorf("unknown media %s", media)
	}

	return pt, mt, nil
}
