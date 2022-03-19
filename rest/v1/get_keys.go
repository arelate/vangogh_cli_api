package v1

import (
	"fmt"
	"github.com/arelate/vangogh_local_data"
	"github.com/boggydigital/nod"
	"net/http"
)

func GetKeys(w http.ResponseWriter, r *http.Request) {

	// GET /v1/keys?product-type&media&sort&desc&format

	if r.Method != http.MethodGet {
		err := fmt.Errorf("unsupported method")
		http.Error(w, nod.Error(err).Error(), http.StatusMethodNotAllowed)
		return
	}

	pt, mt, err := productTypeMediaFromUrl(r.URL)
	if err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusBadRequest)
		return
	}

	sort, desc := sortDescFromUrl(r.URL)
	if !vangogh_local_data.IsValidProperty(sort) {
		err := fmt.Errorf("invalid sort property %s", sort)
		http.Error(w, nod.Error(err).Error(), http.StatusBadRequest)
		return
	}

	vr, err := vangogh_local_data.NewReader(pt, mt)
	if err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
		return
	}

	rxa, err := vangogh_local_data.ConnectReduxAssets(sort, vangogh_local_data.TitleProperty)
	if err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
		return
	}

	sortedIds, err := vangogh_local_data.SortIds(vr.Keys(), rxa, sort, desc)
	if err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
		return
	}

	if err := encode(sortedIds, w, r); err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
	}
}
