package v1

import (
	"encoding/json"
	"fmt"
	"github.com/arelate/vangogh_local_data"
	"github.com/boggydigital/nod"
	"net/http"
)

func GetIndexesList(w http.ResponseWriter, r *http.Request) {

	// GET /v1/indexes-list?product-type&media&sort&desc

	nod.Log("GET %v", r.URL)

	if r.Method != http.MethodGet {
		err := fmt.Errorf("unsupported method")
		http.Error(w, nod.Error(err).Error(), 405)
		return
	}

	pt, mt, err := getProductTypeMedia(r.URL)
	if err != nil {
		http.Error(w, nod.Error(err).Error(), 400)
		return
	}

	sort, desc := getSortDesc(r.URL)
	if !vangogh_local_data.IsValidProperty(sort) {
		err := fmt.Errorf("invalid sort property %s", sort)
		http.Error(w, nod.Error(err).Error(), 400)
		return
	}

	if sids, err := getSortedIds(pt, mt, sort, desc); err != nil {
		http.Error(w, nod.Error(err).Error(), 500)
	} else {
		if err := json.NewEncoder(w).Encode(sids); err != nil {
			http.Error(w, nod.Error(err).Error(), 500)
		}
	}
}
