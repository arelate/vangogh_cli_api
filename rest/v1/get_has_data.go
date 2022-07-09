package v1

import (
	"github.com/arelate/vangogh_local_data"
	"github.com/boggydigital/nod"
	"net/http"
)

func GetHasData(w http.ResponseWriter, r *http.Request) {

	// GET /v1/has_data?product-type&media&id&format

	pts := vangogh_local_data.ValuesFromUrl(r.URL, "product-type")
	mt := vangogh_local_data.MediaFromUrl(r.URL)
	ids, err := vangogh_local_data.IdSetFromUrl(r.URL)

	if err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusBadRequest)
		return
	}

	values := make(map[string]map[string]string, len(pts))

	for _, pt := range pts {

		values[pt] = make(map[string]string, len(ids))

		productType := vangogh_local_data.ParseProductType(pt)

		vr, err := vangogh_local_data.NewReader(productType, mt)

		if err != nil {
			http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
			return
		}

		for id := range ids {
			if ok := vr.Has(id); ok {
				values[pt][id] = "true"
			} else {
				values[pt][id] = "false"
			}
		}

	}

	if err := encode(values, w, r); err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
	}
}
