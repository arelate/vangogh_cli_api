package v1

import (
	"fmt"
	"github.com/arelate/vangogh_local_data"
	"github.com/boggydigital/kvas"
	"github.com/boggydigital/nod"
	"golang.org/x/exp/maps"
	"net/http"
	"sort"
)

func PropertyValuesCounts(rxa kvas.ReduxAssets, property string) map[string]int {
	distValues := make(map[string]int)

	for _, id := range rxa.Keys(property) {
		values, ok := rxa.GetAllValues(property, id)
		if !ok || len(values) == 0 {
			continue
		}

		for _, val := range values {
			if val == "" {
				continue
			}
			distValues[val] = distValues[val] + 1
		}
	}

	return distValues
}

func GetDigest(w http.ResponseWriter, r *http.Request) {

	// GET /v1/digest?property&format

	if r.Method != http.MethodGet {
		err := fmt.Errorf("unsupported method")
		http.Error(w, nod.Error(err).Error(), http.StatusMethodNotAllowed)
		return
	}

	properties := vangogh_local_data.PropertiesFromUrl(r.URL)

	var err error
	if rxa, err = rxa.RefreshReduxAssets(); err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
		return
	}

	if err := rxa.IsSupported(properties...); err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
		return
	}

	propertyValues := make(map[string][]string)

	for _, p := range properties {
		pvc := PropertyValuesCounts(rxa, p)
		values := maps.Keys(pvc)
		sort.Strings(values)
		propertyValues[p] = values
	}

	if err := encode(propertyValues, w, r); err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
	}
}
