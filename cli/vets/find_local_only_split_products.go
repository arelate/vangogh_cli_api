package vets

import (
	"fmt"
	"github.com/arelate/gog_atu"
	"github.com/arelate/vangogh_products"
	"github.com/arelate/vangogh_values"
	"github.com/boggydigital/gost"
	"strconv"
)

func findLocalOnlySplitProducts(pagedPt vangogh_products.ProductType, mt gog_atu.Media) (gost.StrSet, error) {
	if !vangogh_products.IsPaged(pagedPt) {
		return nil, fmt.Errorf("%s is not a paged type", pagedPt)
	}

	pagedIds := gost.NewStrSet()

	vrPaged, err := vangogh_values.NewReader(pagedPt, mt)
	if err != nil {
		return nil, err
	}

	for _, id := range vrPaged.Keys() {
		productGetter, err := vrPaged.ProductsGetter(id)
		if err != nil {
			return nil, err
		}
		for _, idGetter := range productGetter.GetProducts() {
			pagedIds.Add(strconv.Itoa(idGetter.GetId()))
		}
	}

	splitPt := vangogh_products.SplitType(pagedPt)
	vrSplit, err := vangogh_values.NewReader(splitPt, mt)
	if err != nil {
		return nil, err
	}

	splitIdSet := gost.NewStrSetWith(vrSplit.Keys()...)

	return gost.NewStrSetWith(splitIdSet.Except(pagedIds)...), nil
}
