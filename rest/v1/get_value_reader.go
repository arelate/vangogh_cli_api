package v1

import (
	"github.com/arelate/gog_atu"
	"github.com/arelate/vangogh_products"
	"github.com/arelate/vangogh_values"
)

func getValueReader(pt vangogh_products.ProductType, mt gog_atu.Media) (*vangogh_values.ValueReader, error) {
	ptm := productTypeMedia{productType: pt, media: mt}
	if vr, ok := valueReaders[ptm]; !ok || vr == nil {
		return vangogh_values.NewReader(pt, mt)
	}
	return valueReaders[ptm], nil
}
