package reductions

import (
	"github.com/arelate/vangogh_local_data"
	"github.com/boggydigital/nod"
)

var cascadingProperties = []string{
	vangogh_local_data.GOGOrderDateProperty,
	vangogh_local_data.GOGReleaseDateProperty,
	vangogh_local_data.SteamAppIdProperty,
}

//Cascade is a method to assign reductions to products that don't have them,
//and can get them through parent products. Current implementation is a
//template for additional properties and currently only cascades
//GOGOrderDateProperty from store-products (that are referenced in orders)
//to account-products that are linked as store-product.IncludesGames.
func Cascade() error {

	ca := nod.NewProgress("cascading supported properties...")
	defer ca.End()

	rxa, err := vangogh_local_data.ConnectReduxAssets(vangogh_local_data.ReduxProperties()...)
	if err != nil {
		return ca.EndWithError(err)
	}

	if err := rxa.IsSupported(vangogh_local_data.IncludesGamesProperty); err != nil {
		return ca.EndWithError(err)
	}

	ids := rxa.Keys(vangogh_local_data.IncludesGamesProperty)

	ca.TotalInt(len(ids))

	for _, id := range ids {
		includesIds, ok := rxa.GetAllUnchangedValues(vangogh_local_data.IncludesGamesProperty, id)
		if !ok {
			ca.Increment()
			continue
		}
		for _, prop := range cascadingProperties {
			mainValues, ok := rxa.GetAllUnchangedValues(prop, id)
			if !ok {
				continue
			}
			for _, includesId := range includesIds {
				if _, ok := rxa.GetAllUnchangedValues(prop, includesId); !ok {
					if err := rxa.ReplaceValues(prop, includesId, mainValues...); err != nil {
						return ca.EndWithError(err)
					}
				}
			}
		}
		ca.Increment()
	}

	return nil
}
