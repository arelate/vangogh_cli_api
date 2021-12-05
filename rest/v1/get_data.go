package v1

import (
	"encoding/json"
	"fmt"
	"github.com/boggydigital/nod"
	"net/http"
	"strings"
)

func GetData(w http.ResponseWriter, r *http.Request) {

	// GET /v1/data?product-type&media&id

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

	ids := strings.Split(r.URL.Query().Get("id"), ",")

	values := make(map[string]interface{}, len(ids))

	if vr, err := getValueReader(pt, mt); err == nil {

		//var err error
		for i := 0; i < len(ids); i++ {
			if values[ids[i]], err = vr.ReadValue(ids[i]); err != nil {
				http.Error(w, nod.Error(err).Error(), 500)
				return
			}
		}

	} else {
		http.Error(w, nod.Error(err).Error(), 500)
		return
	}

	if err := json.NewEncoder(w).Encode(values); err != nil {
		http.Error(w, nod.Error(err).Error(), 500)
	}
}
