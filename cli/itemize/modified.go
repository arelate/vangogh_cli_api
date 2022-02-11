package itemize

import (
	"github.com/arelate/gog_atu"
	"github.com/arelate/vangogh_data"
	"github.com/boggydigital/gost"
	"github.com/boggydigital/kvas"
	"github.com/boggydigital/nod"
)

func Modified(
	since int64,
	pt vangogh_data.ProductType,
	mt gog_atu.Media) (gost.StrSet, error) {

	ma := nod.Begin(" finding modified %s...", pt)
	defer ma.End()

	//licence products can only update through creation, and we've already handled
	//newly created in itemizeMissing func
	if pt == vangogh_data.LicenceProducts {
		return nil, nil
	}

	destUrl, err := vangogh_data.AbsLocalProductTypeDir(pt, mt)
	if err != nil {
		return nil, ma.EndWithError(err)
	}

	kv, err := kvas.ConnectLocal(destUrl, kvas.JsonExt)
	if err != nil {
		return nil, ma.EndWithError(err)
	}

	modSet := gost.NewStrSetWith(kv.ModifiedAfter(since, false)...)

	ma.EndWithResult(itemizationResult(modSet))

	return modSet, nil
}
