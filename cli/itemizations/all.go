package itemizations

import (
	"github.com/arelate/gog_atu"
	"github.com/arelate/vangogh_data"
)

func All(
	idSet vangogh_data.IdSet,
	missing, updated bool,
	modifiedAfter int64,
	pt vangogh_data.ProductType,
	mt gog_atu.Media) (vangogh_data.IdSet, error) {

	for _, mainPt := range vangogh_data.MainProductTypes(pt) {
		if missing {
			missingIds, err := missingDetail(pt, mainPt, mt, modifiedAfter)
			if err != nil {
				return idSet, err
			}
			idSet.AddSet(missingIds)
		}
		if updated {
			modifiedIds, err := Modified(modifiedAfter, mainPt, mt)
			if err != nil {
				return idSet, err
			}
			idSet.AddSet(modifiedIds)
		}
	}

	return idSet, nil
}
