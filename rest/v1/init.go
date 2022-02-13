package v1

import (
	"github.com/arelate/gog_integration"
	"github.com/arelate/vangogh_local_data"
	"github.com/boggydigital/kvas"
)

type productTypeMedia struct {
	productType vangogh_local_data.ProductType
	media       gog_integration.Media
}

type productTypeMediaSort struct {
	productTypeMedia
	sort string
	desc bool
}

var rxa kvas.ReduxAssets
var valueReaders map[productTypeMedia]*vangogh_local_data.ValueReader
var sortedIds map[productTypeMediaSort][]string
var defaultSort = vangogh_local_data.TitleProperty

func Init() error {
	var err error

	rxa, err = vangogh_local_data.ConnectReduxAssets(vangogh_local_data.ExtractedProperties()...)
	if err != nil {
		return err
	}

	valueReaders = make(map[productTypeMedia]*vangogh_local_data.ValueReader)
	mt := gog_integration.Game
	for _, pt := range vangogh_local_data.LocalProducts() {
		ptm := productTypeMedia{productType: pt, media: mt}
		valueReaders[ptm], err = vangogh_local_data.NewReader(pt, mt)
		if err != nil {
			return err
		}
	}

	//TODO: consider priming that with default sort for a type
	sortedIds = make(map[productTypeMediaSort][]string)

	return nil
}
