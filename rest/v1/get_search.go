package v1

import (
	"fmt"
	"github.com/arelate/vangogh_local_data"
	"github.com/boggydigital/nod"
	"golang.org/x/exp/maps"
	"net/http"
	"strings"
)

func Search(w http.ResponseWriter, r *http.Request) {

	// GET /v1/search?text&(searchable properties)&format

	if r.Method != http.MethodGet {
		err := fmt.Errorf("unsupported method")
		http.Error(w, nod.Error(err).Error(), http.StatusMethodNotAllowed)
		return
	}

	query := make(map[string][]string)
	q := r.URL.Query()

	for _, p := range vangogh_local_data.SearchableProperties() {
		if q.Has(p) {
			val := q.Get(p)
			if val == "" {
				continue
			}
			query[p] = strings.Split(val, "+")
		}
	}

	rxa, err := vangogh_local_data.ConnectReduxAssets(vangogh_local_data.SearchableProperties()...)
	if err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
		return
	}

	found := rxa.Match(query, true)
	keys, err := vangogh_local_data.SortIds(
		maps.Keys(found),
		rxa,
		vangogh_local_data.TitleProperty,
		false)

	if err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
		return
	}

	if err := encode(keys, w, r); err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
	}
}
