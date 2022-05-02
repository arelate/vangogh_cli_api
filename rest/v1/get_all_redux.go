package v1

import (
	"fmt"
	"github.com/arelate/vangogh_local_data"
	"github.com/boggydigital/nod"
	"net/http"
)

func GetAllRedux(w http.ResponseWriter, r *http.Request) {

	// GET /v1/all_redux?property&product-type&media&format

	if r.Method != http.MethodGet {
		err := fmt.Errorf("unsupported method")
		http.Error(w, nod.Error(err).Error(), http.StatusMethodNotAllowed)
		return
	}

	properties := vangogh_local_data.PropertiesFromUrl(r.URL)
	pt := vangogh_local_data.ProductTypeFromUrl(r.URL)
	mt := vangogh_local_data.MediaFromUrl(r.URL)

	if err := RefreshReduxAssets(properties...); err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
		return
	}

	values := make(map[string]map[string][]string)

	vr, err := vangogh_local_data.NewReader(pt, mt)
	if err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
		return
	}

	for _, id := range vr.Keys() {
		propValues := make(map[string][]string)
		for _, prop := range properties {
			propValues[prop], _ = rxa.GetAllValues(prop, id)
		}
		values[id] = propValues
	}

	if err := encode(values, w, r); err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
		return
	}
}
