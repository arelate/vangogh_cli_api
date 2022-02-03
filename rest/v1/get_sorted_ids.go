package v1

import (
	"github.com/arelate/gog_atu"
	"github.com/arelate/vangogh_products"
	"github.com/arelate/vangogh_sets"
)

func getSortedIds(pt vangogh_products.ProductType, mt gog_atu.Media, sort string, desc bool) ([]string, error) {

	ptms := productTypeMediaSort{
		productTypeMedia: productTypeMedia{productType: pt, media: mt},
		sort:             sort,
		desc:             desc,
	}

	if sids, ok := sortedIds[ptms]; ok {
		return sids, nil
	}

	if vr, err := getValueReader(pt, mt); err != nil {
		return nil, err
	} else {
		idSet := vangogh_sets.IdSetWith(vr.All()...)
		sortedIds[ptms] = idSet.Sort(exl, sort, desc)
	}

	return sortedIds[ptms], nil
}
