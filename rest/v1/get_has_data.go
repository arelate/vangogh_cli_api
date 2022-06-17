package v1

import (
	"github.com/arelate/vangogh_local_data"
	"github.com/boggydigital/nod"
	"net/http"
)

func GetHasData(w http.ResponseWriter, r *http.Request) {

	// GET /v1/has_data?product-type&media&id&format

	pt := vangogh_local_data.ProductTypeFromUrl(r.URL)
	mt := vangogh_local_data.MediaFromUrl(r.URL)
	ids, err := vangogh_local_data.IdSetFromUrl(r.URL)
	if err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusBadRequest)
		return
	}

	values := make(map[string]string, len(ids))

	vr, err := vangogh_local_data.NewReader(pt, mt)

	if err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
		return
	}

	for id := range ids {
		if ok := vr.Has(id); ok {
			values[id] = "true"
		} else {
			values[id] = "false"
		}
	}

	if err := encode(values, w, r); err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
	}
}
