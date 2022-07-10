package rest

import (
	"encoding/gob"
	"github.com/arelate/gog_integration"
	"github.com/arelate/steam_integration"
	"github.com/arelate/vangogh_local_data"
	"github.com/boggydigital/kvas"
)

var rxa kvas.ReduxAssets

func Init() error {

	//GOG.com types
	gob.Register(gog_integration.AccountPage{})
	gob.Register(gog_integration.AccountProduct{})
	gob.Register(gog_integration.ApiProductV1{})
	gob.Register(gog_integration.ApiProductV2{})
	gob.Register(gog_integration.Details{})
	gob.Register(gog_integration.Licences{})
	gob.Register(gog_integration.OrderPage{})
	gob.Register(gog_integration.Order{})
	gob.Register(gog_integration.StorePage{})
	gob.Register(gog_integration.StoreProduct{})
	gob.Register(gog_integration.WishlistPage{})
	//Steam types
	gob.Register(steam_integration.AppList{})
	gob.Register(steam_integration.GetNewsForAppResponse{})
	gob.Register(steam_integration.AppReviews{})

	var err error
	properties := vangogh_local_data.ReduxProperties()
	//used by get_downloads
	properties = append(properties, vangogh_local_data.NativeLanguageNameProperty)
	rxa, err = vangogh_local_data.ConnectReduxAssets(properties...)
	return err
}
