package v1

import (
	"fmt"
	"github.com/arelate/vangogh_local_data"
	"github.com/boggydigital/nod"
	"net/http"
	"strings"
)

func GetRedux(w http.ResponseWriter, r *http.Request) {

	// GET /v1/redux?property&id&format

	if r.Method != http.MethodGet {
		err := fmt.Errorf("unsupported method")
		http.Error(w, nod.Error(err).Error(), http.StatusMethodNotAllowed)
		return
	}

	properties := strings.Split(r.URL.Query().Get("property"), ",")
	rxa, err := vangogh_local_data.ConnectReduxAssets(properties...)
	if err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
		return
	}

	ids := strings.Split(r.URL.Query().Get("id"), ",")

	values := make(map[string]map[string][]string, len(ids))
	for _, id := range ids {
		propValues := make(map[string][]string)
		for _, prop := range properties {
			propValues[prop], _ = rxa.GetAllValues(prop, id)
		}
		values[id] = propValues
	}

	if err := encode(values, w, r); err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
	}
}
