package v1

import (
	"encoding/json"
	"fmt"
	"github.com/arelate/vangogh_properties"
	"net/http"
)

func GetIndexes(w http.ResponseWriter, r *http.Request) {

	// GET /v1/indexes?product-type&media&sort&desc

	if r.Method != http.MethodGet {
		http.Error(w, "unsupported method", 405)
		return
	}

	pt, mt, err := getProductTypeMedia(r.URL)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	sort, desc := getSortDesc(r.URL)
	if !vangogh_properties.IsValid(sort) {
		http.Error(w, fmt.Sprintf("invalid sort property %s", sort), 400)
		return
	}

	if sids, err := getSortedIds(pt, mt, sort, desc); err != nil {
		http.Error(w, err.Error(), 500)
	} else {
		if err := json.NewEncoder(w).Encode(sids); err != nil {
			http.Error(w, err.Error(), 500)
		}
	}
}
