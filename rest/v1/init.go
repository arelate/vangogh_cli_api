package v1

import (
	"github.com/arelate/gog_atu"
	"github.com/arelate/vangogh_data"
	"github.com/boggydigital/kvas"
)

type productTypeMedia struct {
	productType vangogh_data.ProductType
	media       gog_atu.Media
}

type productTypeMediaSort struct {
	productTypeMedia
	sort string
	desc bool
}

var rxa kvas.ReduxAssets
var valueReaders map[productTypeMedia]*vangogh_data.ValueReader
var sortedIds map[productTypeMediaSort][]string
var defaultSort = vangogh_data.TitleProperty

func Init() error {
	var err error

	rxa, err = vangogh_data.ConnectReduxAssets(vangogh_data.Extracted()...)
	if err != nil {
		return err
	}

	valueReaders = make(map[productTypeMedia]*vangogh_data.ValueReader)
	mt := gog_atu.Game
	for _, pt := range vangogh_data.LocalProducts() {
		ptm := productTypeMedia{productType: pt, media: mt}
		valueReaders[ptm], err = vangogh_data.NewReader(pt, mt)
		if err != nil {
			return err
		}
	}

	//TODO: consider priming that with default sort for a type
	sortedIds = make(map[productTypeMediaSort][]string)

	return nil
}
