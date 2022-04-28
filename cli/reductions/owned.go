package reductions

import (
	"github.com/arelate/gog_integration"
	"github.com/arelate/vangogh_local_data"
	"github.com/boggydigital/kvas"
	"github.com/boggydigital/nod"
)

func CheckOwnership(idSet map[string]bool, rxa kvas.ReduxAssets) (map[string]bool, error) {

	ownedSet := make(map[string]bool)

	if err := rxa.IsSupported(vangogh_local_data.SlugProperty, vangogh_local_data.IncludesGamesProperty); err != nil {
		return ownedSet, err
	}

	vrLicenceProducts, err := vangogh_local_data.NewReader(vangogh_local_data.LicenceProducts, gog_integration.Game)
	if err != nil {
		return ownedSet, err
	}

	for id := range idSet {

		if vrLicenceProducts.Has(id) {
			ownedSet[id] = true
			continue
		}

		includesGames, ok := rxa.GetAllUnchangedValues(vangogh_local_data.IncludesGamesProperty, id)
		if !ok || len(includesGames) == 0 {
			continue
		}

		ownAllIncludedGames := true
		for _, igId := range includesGames {
			ownAllIncludedGames = ownAllIncludedGames && vrLicenceProducts.Has(igId)
			if !ownAllIncludedGames {
				break
			}
		}

		if ownAllIncludedGames {
			ownedSet[id] = true
		}
	}

	return ownedSet, nil
}

func Owned(mt gog_integration.Media) error {

	oa := nod.Begin(" %s...", vangogh_local_data.OwnedProperty)
	defer oa.End()

	rxa, err := vangogh_local_data.ConnectReduxAssets(
		vangogh_local_data.TitleProperty,
		vangogh_local_data.OwnedProperty,
		vangogh_local_data.SlugProperty,
		vangogh_local_data.IncludesGamesProperty)
	if err != nil {
		return oa.EndWithError(err)
	}

	//vrStoreProducts, err := vangogh_local_data.NewReader(vangogh_local_data.StoreProducts, mt)
	//if err != nil {
	//	return oa.EndWithError(err)
	//}

	idSet := make(map[string]bool)
	for _, id := range rxa.Keys(vangogh_local_data.TitleProperty) {
		idSet[id] = true
	}

	owned, err := CheckOwnership(idSet, rxa)
	if err != nil {
		return oa.EndWithError(err)
	}

	ownedRdx := make(map[string][]string)

	for _, id := range rxa.Keys(vangogh_local_data.TitleProperty) {
		if _, ok := owned[id]; ok {
			ownedRdx[id] = []string{"true"}
		} else {
			ownedRdx[id] = []string{"false"}
		}
	}

	if err := rxa.BatchReplaceValues(vangogh_local_data.OwnedProperty, ownedRdx); err != nil {
		return oa.EndWithError(err)
	}

	oa.EndWithResult("done")

	return nil
}
