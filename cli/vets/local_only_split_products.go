package vets

import (
	"fmt"
	"github.com/arelate/gog_atu"
	"github.com/arelate/vangogh_api/cli/expand"
	"github.com/arelate/vangogh_api/cli/remove"
	"github.com/arelate/vangogh_products"
	"github.com/arelate/vangogh_properties"
	"github.com/boggydigital/nod"
)

func LocalOnlySplitProducts(mt gog_atu.Media, fix bool) error {

	sloa := nod.Begin("checking for local only split products...")
	defer sloa.End()

	rxa, err := vangogh_properties.ConnectReduxAssets(vangogh_properties.TitleProperty)
	if err != nil {
		return sloa.EndWithError(err)
	}

	for _, pagedPt := range vangogh_products.Paged() {

		splitPt := vangogh_products.SplitType(pagedPt)

		pa := nod.Begin(" checking %s not present in %s...", splitPt, pagedPt)

		localOnlyProducts, err := findLocalOnlySplitProducts(pagedPt, mt)
		if err != nil {
			return pa.EndWithError(err)
		}

		if localOnlyProducts.Len() > 0 {

			summary, err := expand.IdsToPropertyLists(
				localOnlyProducts,
				nil,
				[]string{vangogh_properties.TitleProperty},
				rxa)

			if err != nil {
				_ = pa.EndWithError(err)
				continue
			}

			pa.EndWithSummary(fmt.Sprintf("found %d:", localOnlyProducts.Len()), summary)

			if fix {
				fa := nod.Begin(" removing local only %s...", splitPt)
				if err := remove.Data(localOnlyProducts.All(), splitPt, mt); err != nil {
					return fa.EndWithError(err)
				}
				fa.EndWithResult("done")
			}
		} else {
			pa.EndWithResult("none found")
		}
	}

	sloa.EndWithResult("done")

	return nil
}
