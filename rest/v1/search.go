package v1

import (
	"encoding/json"
	"github.com/arelate/vangogh_local_data"
	"github.com/boggydigital/nod"
	"net/http"
	"strings"
)

func Search(w http.ResponseWriter, r *http.Request) {

	// GET /v1/search?text&(searchable properties)

	query := make(map[string][]string)
	q := r.URL.Query()

	for _, p := range vangogh_local_data.SearchableProperties() {
		if q.Has(p) {
			val := q.Get(p)
			if val == "" {
				continue
			}
			query[p] = strings.Split(val, " ")
		}
	}

	found := rxa.Match(query, true)
	keys := make([]string, 0, len(found))

	for id := range found {
		keys = append(keys, id)
	}

	if err := json.NewEncoder(w).Encode(keys); err != nil {
		http.Error(w, nod.Error(err).Error(), 500)
	}
}
