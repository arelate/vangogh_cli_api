package reductions

import (
	"github.com/arelate/gog_integration"
	"github.com/arelate/vangogh_local_data"
	"github.com/boggydigital/nod"
)

func Types(mt gog_integration.Media) error {

	ta := nod.Begin(" %s...", vangogh_local_data.TypesProperty)
	defer ta.End()

	idsTypes := make(map[string][]string)

	for _, pt := range vangogh_local_data.LocalProducts() {

		vr, err := vangogh_local_data.NewReader(pt, mt)
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

	typesEx, err := vangogh_local_data.ConnectReduxAssets(vangogh_local_data.TypesProperty)
	if err != nil {
		return ta.EndWithError(err)
	}

	if err := typesEx.BatchReplaceValues(vangogh_local_data.TypesProperty, idsTypes); err != nil {
		return ta.EndWithError(err)
	}

	ta.EndWithResult("done")

	return nil
}
