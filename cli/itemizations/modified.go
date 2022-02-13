package itemizations

import (
	"github.com/arelate/gog_integration"
	"github.com/arelate/vangogh_local_data"
	"github.com/boggydigital/kvas"
	"github.com/boggydigital/nod"
)

func Modified(
	since int64,
	pt vangogh_local_data.ProductType,
	mt gog_integration.Media) (vangogh_local_data.IdSet, error) {

	ma := nod.Begin(" finding modified %s...", pt)
	defer ma.End()

	modSet := vangogh_local_data.NewIdSet()

	//licence products can only update through creation, and we've already handled
	//newly created in itemizeMissing func
	if pt == vangogh_local_data.LicenceProducts {
		return modSet, nil
	}

	destUrl, err := vangogh_local_data.AbsLocalProductTypeDir(pt, mt)
	if err != nil {
		return modSet, ma.EndWithError(err)
	}

	kv, err := kvas.ConnectLocal(destUrl, kvas.JsonExt)
	if err != nil {
		return modSet, ma.EndWithError(err)
	}

	modSet.Add(kv.ModifiedAfter(since, false)...)

	ma.EndWithResult(itemizationResult(modSet))

	return modSet, nil
}
