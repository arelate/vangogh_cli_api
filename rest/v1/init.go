package v1

import (
	"github.com/arelate/vangogh_local_data"
	"github.com/boggydigital/kvas"
)

var rxa kvas.ReduxAssets

func Init() error {
	var err error
	properties := vangogh_local_data.ReduxProperties()
	//used by get_downloads
	properties = append(properties, vangogh_local_data.NativeLanguageNameProperty)
	rxa, err = vangogh_local_data.ConnectReduxAssets(properties...)
	return err
}
