package extract

import (
	"github.com/arelate/gog_atu"
	"github.com/arelate/vangogh_products"
	"github.com/arelate/vangogh_properties"
	"github.com/arelate/vangogh_values"
	"github.com/boggydigital/nod"
)

func Types(mt gog_atu.Media) error {

	ta := nod.Begin(" %s...", vangogh_properties.TypesProperty)
	defer ta.End()

	idsTypes := make(map[string][]string)

	for _, pt := range vangogh_products.Local() {

		vr, err := vangogh_values.NewReader(pt, mt)
		if err != nil {
			return ta.EndWithError(err)
		}

		for _, id := range vr.Keys() {

			if idsTypes[id] == nil {
				idsTypes[id] = make([]string, 0)
			}

			idsTypes[id] = append(idsTypes[id], pt.String())
		}
	}

	typesEx, err := vangogh_properties.ConnectReduxAssets(vangogh_properties.TypesProperty)
	if err != nil {
		return ta.EndWithError(err)
	}

	if err := typesEx.BatchReplaceValues(vangogh_properties.TypesProperty, idsTypes); err != nil {
		return ta.EndWithError(err)
	}

	ta.EndWithResult("done")

	return nil
}
