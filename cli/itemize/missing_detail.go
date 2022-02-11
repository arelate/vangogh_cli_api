package itemize

import (
	"fmt"
	"github.com/arelate/gog_atu"
	"github.com/arelate/vangogh_data"
	"github.com/boggydigital/gost"
	"github.com/boggydigital/kvas"
	"github.com/boggydigital/nod"
)

func itemizationResult(set gost.StrSet) string {
	if set.Len() == 0 {
		return "found nothing"
	} else {
		return fmt.Sprintf("found %d", set.Len())
	}
}

func missingDetail(
	detailPt, mainPt vangogh_data.ProductType,
	mt gog_atu.Media,
	since int64) (gost.StrSet, error) {

	//api-products-v2 provides
	//includes-games, is-included-by-games,
	//requires-games, is-required-by-games
	if mainPt == vangogh_data.ApiProductsV2 &&
		detailPt == vangogh_data.ApiProductsV2 {
		lg, err := linkedGames(since)
		if err != nil {
			return lg, err
		}
		return lg, nil
	}

	//licences give a signal when DLC has been purchased, this would add
	//required (base) game details to the updates
	if mainPt == vangogh_data.LicenceProducts &&
		detailPt == vangogh_data.Details {
		rg, err := RequiredAndIncluded(since)
		if err != nil {
			return rg, err
		}
		return rg, nil
	}

	mda := nod.Begin(" finding missing %s for %s...", detailPt, mainPt)
	defer mda.End()

	missingIdSet := gost.NewStrSet()

	mainDestUrl, err := vangogh_data.AbsLocalProductTypeDir(mainPt, mt)
	if err != nil {
		return missingIdSet, mda.EndWithError(err)
	}

	detailDestUrl, err := vangogh_data.AbsLocalProductTypeDir(detailPt, mt)
	if err != nil {
		return missingIdSet, mda.EndWithError(err)
	}

	kvMain, err := kvas.ConnectLocal(mainDestUrl, kvas.JsonExt)
	if err != nil {
		return missingIdSet, mda.EndWithError(err)
	}

	kvDetail, err := kvas.ConnectLocal(detailDestUrl, kvas.JsonExt)
	if err != nil {
		return missingIdSet, mda.EndWithError(err)
	}

	for _, id := range kvMain.Keys() {
		if !kvDetail.Has(id) {
			nod.Log("adding missing %s: #%s", detailPt, id)
			missingIdSet.Add(id)
		}
	}

	mda.EndWithResult(itemizationResult(missingIdSet))

	if mainPt == vangogh_data.AccountProducts &&
		detailPt == vangogh_data.Details {
		updatedAccountProducts, err := AccountProductsUpdates(mt)
		if err != nil {
			return missingIdSet, err
		}
		for uapId := range updatedAccountProducts {
			missingIdSet.Add(uapId)
		}
	}

	return missingIdSet, nil
}
