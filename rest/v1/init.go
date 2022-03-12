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

var rxa kvas.ReduxAssets
var valueReaders map[productTypeMedia]*vangogh_local_data.ValueReader

func Init() error {
	var err error

	rxa, err = vangogh_local_data.ConnectReduxAssets(vangogh_local_data.ReduxProperties()...)
	if err != nil {
		return err
	}

	valueReaders = make(map[productTypeMedia]*vangogh_local_data.ValueReader)

	return nil
}
