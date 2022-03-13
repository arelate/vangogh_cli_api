package v1

import (
	"encoding/json"
	"fmt"
	"github.com/arelate/vangogh_local_data"
	"github.com/boggydigital/nod"
	"net/http"
)

func GetKeys(w http.ResponseWriter, r *http.Request) {

	// GET /v1/keys?product-type&media&sort&desc

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

	vr, err := getValueReader(pt, mt)
	if err != nil {
		http.Error(w, nod.Error(err).Error(), 500)
	}

	idSet := vangogh_local_data.IdSetFromSlice(vr.Keys()...)
	sortedIds := idSet.Sort(rxa, sort, desc)

	if err := json.NewEncoder(w).Encode(sortedIds); err != nil {
		http.Error(w, nod.Error(err).Error(), 500)
	}
}
