package v1

import (
	"fmt"
	"github.com/arelate/vangogh_local_data"
	"github.com/boggydigital/nod"
	"net/http"
	"strings"
)

func GetData(w http.ResponseWriter, r *http.Request) {

	// GET /v1/data?product-type&media&id&format

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

	ids := strings.Split(r.URL.Query().Get("id"), ",")

	values := make(map[string]interface{}, len(ids))

	if vr, err := vangogh_local_data.NewReader(pt, mt); err == nil {

		//var err error
		for i := 0; i < len(ids); i++ {
			if values[ids[i]], err = vr.ReadValue(ids[i]); err != nil {
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
