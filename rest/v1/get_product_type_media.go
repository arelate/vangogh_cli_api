package v1

import (
	"fmt"
	"github.com/arelate/gog_atu"
	"github.com/arelate/vangogh_products"
	"net/url"
)

func getProductTypeMedia(u *url.URL) (vangogh_products.ProductType, gog_atu.Media, error) {
	q := u.Query()

	productType := q.Get("product-type")
	pt := vangogh_products.Parse(productType)
	if pt == vangogh_products.Unknown {
		return pt, gog_atu.Unknown, fmt.Errorf("unknown product-type %s", productType)
	}

	media := q.Get("media")
	mt := gog_atu.ParseMedia(media)
	if mt == gog_atu.Unknown {
		return pt, mt, fmt.Errorf("unknown media %s", media)
	}

	return pt, mt, nil
}
