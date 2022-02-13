package v1

import (
	"github.com/arelate/gog_integration"
	"github.com/arelate/vangogh_local_data"
)

func getValueReader(pt vangogh_local_data.ProductType, mt gog_integration.Media) (*vangogh_local_data.ValueReader, error) {
	ptm := productTypeMedia{productType: pt, media: mt}
	if vr, ok := valueReaders[ptm]; !ok || vr == nil {
		return vangogh_local_data.NewReader(pt, mt)
	}
	return valueReaders[ptm], nil
}
