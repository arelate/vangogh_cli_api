package itemizations

import (
	"github.com/arelate/gog_integration"
	"github.com/arelate/vangogh_local_data"
	"github.com/boggydigital/nod"
)

func linkedGames(modifiedAfter int64) (vangogh_local_data.IdSet, error) {

	lga := nod.Begin(" finding missing linked %s...", vangogh_local_data.ApiProductsV2)
	defer lga.End()

	missingSet := vangogh_local_data.NewIdSet()

	//currently, api-products-v2 support only gog_integration.Game, and since this method is exclusively
	//using api-products-v2 we're fine specifying media directly and not taking as a parameter
	vrApv2, err := vangogh_local_data.NewReader(vangogh_local_data.ApiProductsV2, gog_integration.Game)
	if err != nil {
		return missingSet, lga.EndWithError(err)
	}

	modifiedApv2 := vrApv2.ModifiedAfter(modifiedAfter, false)
	if len(modifiedApv2) > 0 {
		nod.Log("modified %s: %v", vangogh_local_data.ApiProductsV2, modifiedApv2)
	}

	for _, id := range modifiedApv2 {

		// have to use product reader and not extracts here, since extracts wouldn't be ready
		// while we're still getting data. Attempting to minimize the impact by only querying
		// new or updated api-product-v2 items since start to the sync
		apv2, err := vrApv2.ApiProductV2(id)

		if err != nil {
			return missingSet, lga.EndWithError(err)
		}

		gig := apv2.GetIncludesGames()
		if len(gig) > 0 {
			nod.Log("%s #%s includes-games: %v", vangogh_local_data.ApiProductsV2, id, gig)
		}
		lgs := gig

		giiig := apv2.GetIsIncludedInGames()
		if len(giiig) > 0 {
			nod.Log("%s #%s is-included-in-games: %v", vangogh_local_data.ApiProductsV2, id, giiig)
		}
		lgs = append(lgs, giiig...)

		grg := apv2.GetRequiresGames()
		if len(grg) > 0 {
			nod.Log("%s #%s requires-games: %v", vangogh_local_data.ApiProductsV2, id, grg)
		}
		lgs = append(lgs, grg...)

		girbg := apv2.GetIsRequiredByGames()
		if len(girbg) > 0 {
			nod.Log("%s #%s is-required-by-games: %v", vangogh_local_data.ApiProductsV2, id, girbg)
		}
		lgs = append(lgs, girbg...)

		for _, lid := range lgs {
			if !vrApv2.Has(lid) {
				missingSet.Add(lid)
			}
		}
	}

	lga.EndWithResult(itemizationResult(missingSet))

	return missingSet, nil
}
