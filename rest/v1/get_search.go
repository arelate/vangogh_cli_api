package v1

import (
	"github.com/arelate/vangogh_local_data"
	"github.com/boggydigital/nod"
	"golang.org/x/exp/maps"
	"net/http"
	"strings"
)

func Search(w http.ResponseWriter, r *http.Request) {

	// GET /v1/search?text&(searchable properties)&sort&desc&format

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

	sort := q.Get("sort")
	if sort == "" {
		sort = vangogh_local_data.TitleProperty
	}
	desc := q.Get("desc") == "true"

	properties := []string{sort}
	for p := range query {
		properties = append(properties, p)
	}

	detailedProperties := vangogh_local_data.DetailAllAggregateProperties(properties...)

	if err := RefreshReduxAssets(maps.Keys(detailedProperties)...); err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
		return
	}

	found := rxa.Match(query, true)
	keys, err := vangogh_local_data.SortIds(
		maps.Keys(found),
		rxa,
		sort,
		desc)

	if err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
		return
	}

	if err := encode(keys, w, r); err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
	}
}
