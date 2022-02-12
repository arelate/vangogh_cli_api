package vets

import (
	"fmt"
	"github.com/arelate/gog_atu"
	"github.com/arelate/vangogh_data"
	"github.com/boggydigital/gost"
	"strconv"
)

func findLocalOnlySplitProducts(pagedPt vangogh_data.ProductType, mt gog_atu.Media) (vangogh_data.IdSet, error) {
	emptyIdSet := vangogh_data.NewIdSet()

	if !vangogh_data.IsPagedProduct(pagedPt) {
		return emptyIdSet, fmt.Errorf("%s is not a paged type", pagedPt)
	}

	pagedIds := gost.NewStrSet()

	vrPaged, err := vangogh_data.NewReader(pagedPt, mt)
	if err != nil {
		return emptyIdSet, err
	}

	for _, id := range vrPaged.Keys() {
		productGetter, err := vrPaged.ProductsGetter(id)
		if err != nil {
			return emptyIdSet, err
		}
		for _, idGetter := range productGetter.GetProducts() {
			pagedIds.Add(strconv.Itoa(idGetter.GetId()))
		}
	}

	splitPt := vangogh_data.SplitProductType(pagedPt)
	vrSplit, err := vangogh_data.NewReader(splitPt, mt)
	if err != nil {
		return emptyIdSet, err
	}

	splitIdSet := gost.NewStrSetWith(vrSplit.Keys()...)

	return vangogh_data.IdSetFromSlice(splitIdSet.Except(pagedIds)...), nil
}
