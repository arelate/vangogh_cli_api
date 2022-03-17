package vets

import (
	"fmt"
	"github.com/arelate/gog_integration"
	"github.com/arelate/vangogh_local_data"
	"github.com/boggydigital/gost"
	"strconv"
)

func findLocalOnlySplitProducts(pagedPt vangogh_local_data.ProductType, mt gog_integration.Media) (*vangogh_local_data.IdSet, error) {
	emptyIdSet := vangogh_local_data.NewIdSet()

	if !vangogh_local_data.IsPagedProduct(pagedPt) {
		return emptyIdSet, fmt.Errorf("%s is not a paged type", pagedPt)
	}

	pagedIds := gost.NewStrSet()

	vrPaged, err := vangogh_local_data.NewReader(pagedPt, mt)
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

	splitPt := vangogh_local_data.SplitProductType(pagedPt)
	vrSplit, err := vangogh_local_data.NewReader(splitPt, mt)
	if err != nil {
		return emptyIdSet, err
	}

	splitIdSet := gost.NewStrSetWith(vrSplit.Keys()...)

	return vangogh_local_data.IdSetFromSlice(splitIdSet.Except(pagedIds)...), nil
}
