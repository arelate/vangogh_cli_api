package v1

import (
	"github.com/arelate/gog_atu"
	"github.com/arelate/vangogh_data"
)

func getValueReader(pt vangogh_data.ProductType, mt gog_atu.Media) (*vangogh_data.ValueReader, error) {
	ptm := productTypeMedia{productType: pt, media: mt}
	if vr, ok := valueReaders[ptm]; !ok || vr == nil {
		return vangogh_data.NewReader(pt, mt)
	}
	return valueReaders[ptm], nil
}
