package v1

import (
	"encoding/json"
	"fmt"
	"github.com/arelate/vangogh_local_data"
	"github.com/boggydigital/nod"
	"net/http"
	"strings"
)

func GetAllRedux(w http.ResponseWriter, r *http.Request) {

	// GET /v1/all_redux?property&product-type&media

	nod.Log("GET %v", r.URL)

	if r.Method != http.MethodGet {
		err := fmt.Errorf("unsupported method")
		http.Error(w, nod.Error(err).Error(), 405)
		return
	}

	properties := strings.Split(r.URL.Query().Get("property"), ",")
	for _, prop := range properties {
		if err := rxa.IsSupported(prop); err != nil {
			http.Error(w, fmt.Sprintf("unsupported property %s", prop), 400)
			return
		}
	}

	productType := vangogh_local_data.ProductTypeFromUrl(r.URL)
	media := vangogh_local_data.MediaFromUrl(r.URL)

	values := make(map[string]map[string][]string)

	vr, err := getValueReader(productType, media)
	if err != nil {
		http.Error(w, nod.Error(err).Error(), 500)
		return
	}

	for _, id := range vr.Keys() {
		propValues := make(map[string][]string)
		for _, prop := range properties {
			propValues[prop], _ = rxa.GetAllValues(prop, id)
		}
		values[id] = propValues
	}
	if err := json.NewEncoder(w).Encode(values); err != nil {
		http.Error(w, nod.Error(err).Error(), 500)
	}
}
