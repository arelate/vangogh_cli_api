package v1

import (
	"github.com/arelate/gog_atu"
	"github.com/arelate/vangogh_products"
	"github.com/arelate/vangogh_properties"
	"github.com/arelate/vangogh_values"
	"github.com/boggydigital/kvas"
)

type productTypeMedia struct {
	productType vangogh_products.ProductType
	media       gog_atu.Media
}

type productTypeMediaSort struct {
	productTypeMedia
	sort string
	desc bool
}

var rxa kvas.ReduxAssets
var valueReaders map[productTypeMedia]*vangogh_values.ValueReader
var sortedIds map[productTypeMediaSort][]string
var defaultSort = vangogh_properties.TitleProperty

func Init() error {
	var err error

	rxa, err = vangogh_properties.ConnectReduxAssets(vangogh_properties.Extracted()...)
	if err != nil {
		return err
	}

	valueReaders = make(map[productTypeMedia]*vangogh_values.ValueReader)
	mt := gog_atu.Game
	for _, pt := range vangogh_products.Local() {
		ptm := productTypeMedia{productType: pt, media: mt}
		valueReaders[ptm], err = vangogh_values.NewReader(pt, mt)
		if err != nil {
			return err
		}
	}

	//TODO: consider priming that with default sort for a type
	sortedIds = make(map[productTypeMediaSort][]string)

	return nil
}
