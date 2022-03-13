package v1

import (
	"encoding/json"
	"fmt"
	"github.com/boggydigital/nod"
	"net/http"
	"strings"
)

func GetRedux(w http.ResponseWriter, r *http.Request) {

	// GET /v1/redux?property&id

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

	ids := strings.Split(r.URL.Query().Get("id"), ",")

	values := make(map[string]map[string][]string, len(ids))
	for _, id := range ids {
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
