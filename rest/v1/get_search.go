package v1

import (
	"fmt"
	"github.com/arelate/vangogh_local_data"
	"github.com/boggydigital/nod"
	"golang.org/x/exp/maps"
	"net/http"
	"strings"
	"time"
)

func Search(w http.ResponseWriter, r *http.Request) {

	// GET /v1/search?text&(searchable properties)&sort&desc&format

	if r.Method != http.MethodGet {
		err := fmt.Errorf("unsupported method")
		http.Error(w, nod.Error(err).Error(), http.StatusMethodNotAllowed)
		return
	}

	// redux assets mod time is used to:
	// 1) set Last-Modified header
	// 2) check if content was modified since client cache
	if ramt, err := rxa.ReduxAssetsModTime(); err == nil {
		w.Header().Set("Last-Modified", time.Unix(ramt, 0).Format(time.RFC1123))
		if imsh := r.Header.Get("If-Modified-Since"); imsh != "" {
			if ims, err := time.Parse(time.RFC1123, imsh); err == nil {
				if ramt <= ims.Unix() {
					w.WriteHeader(http.StatusNotModified)
					return
				}
			}
		}
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

	var err error
	if rxa, err = rxa.RefreshReduxAssets(); err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
		return
	}

	if err := rxa.IsSupported(maps.Keys(detailedProperties)...); err != nil {
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
