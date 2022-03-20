package v1

import (
	"fmt"
	"github.com/arelate/vangogh_local_data"
	"github.com/boggydigital/nod"
	"net/http"
)

func GetData(w http.ResponseWriter, r *http.Request) {

	// GET /v1/data?product-type&media&id&format

	if r.Method != http.MethodGet {
		err := fmt.Errorf("unsupported method")
		http.Error(w, nod.Error(err).Error(), http.StatusMethodNotAllowed)
		return
	}

	pt := vangogh_local_data.ProductTypeFromUrl(r.URL)
	mt := vangogh_local_data.MediaFromUrl(r.URL)
	ids, err := vangogh_local_data.IdSetFromUrl(r.URL)
	if err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusMethodNotAllowed)
		return
	}

	values := make(map[string]interface{}, len(ids))

	if vr, err := vangogh_local_data.NewReader(pt, mt); err == nil {

		for id := range ids {
			if values[id], err = vr.ReadValue(id); err != nil {
				http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
				return
			}
		}

	} else {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
		return
	}

	if err := encode(values, w, r); err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
	}
}
