package v1

import (
	"encoding/json"
	"github.com/arelate/vangogh_local_data"
	"github.com/boggydigital/nod"
	"net/http"
)

func Search(w http.ResponseWriter, r *http.Request) {

	// GET /v1/search?text&(searchable properties)

	query := make(map[string][]string)
	q := r.URL.Query()

	for _, p := range vangogh_local_data.SearchableProperties() {
		if q.Has(p) {
			query[p] = []string{q.Get(p)}
		}
	}

	found := rxa.Match(query, true)
	keys := make([]string, 0, len(found))

	if err := json.NewEncoder(w).Encode(keys); err != nil {
		http.Error(w, nod.Error(err).Error(), 500)
	}
}
