package v1

import (
	"github.com/arelate/gog_integration"
	"github.com/arelate/vangogh_local_data"
)

func getSortedIds(pt vangogh_local_data.ProductType, mt gog_integration.Media, sort string, desc bool) ([]string, error) {

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
		idSet := vangogh_local_data.IdSetFromSlice(vr.Keys()...)
		sortedIds[ptms] = idSet.Sort(rxa, sort, desc)
	}

	return sortedIds[ptms], nil
}
