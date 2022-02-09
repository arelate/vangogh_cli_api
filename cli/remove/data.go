package remove

import (
	"github.com/arelate/gog_atu"
	"github.com/arelate/vangogh_products"
	"github.com/arelate/vangogh_urls"
	"github.com/boggydigital/kvas"
)

func Data(ids []string, pt vangogh_products.ProductType, mt gog_atu.Media) error {
	ptDir, err := vangogh_urls.AbsLocalProductsDir(pt, mt)
	if err != nil {
		return err
	}
	kvPt, err := kvas.ConnectLocal(ptDir, kvas.JsonExt)
	if err != nil {
		return err
	}

	for _, id := range ids {
		//log.Printf("remove %s (%s) id %s", pt, mt, id)
		if _, err := kvPt.Cut(id); err != nil {
			return err
		}
	}

	return nil
}
