package reductions

import (
	"github.com/arelate/gog_atu"
	"github.com/arelate/vangogh_data"
	"github.com/boggydigital/nod"
)

func Types(mt gog_atu.Media) error {

	ta := nod.Begin(" %s...", vangogh_data.TypesProperty)
	defer ta.End()

	idsTypes := make(map[string][]string)

	for _, pt := range vangogh_data.LocalProducts() {

		vr, err := vangogh_data.NewReader(pt, mt)
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

	typesEx, err := vangogh_data.ConnectReduxAssets(vangogh_data.TypesProperty)
	if err != nil {
		return ta.EndWithError(err)
	}

	if err := typesEx.BatchReplaceValues(vangogh_data.TypesProperty, idsTypes); err != nil {
		return ta.EndWithError(err)
	}

	ta.EndWithResult("done")

	return nil
}
