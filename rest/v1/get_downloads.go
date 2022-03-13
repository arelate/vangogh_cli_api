package v1

import (
	"encoding/json"
	"fmt"
	"github.com/arelate/gog_integration"
	"github.com/arelate/vangogh_local_data"
	"github.com/boggydigital/nod"
	"net/http"
)

func GetDownloads(w http.ResponseWriter, r *http.Request) {

	// GET /v1/downloads?id&os&lang

	if r.Method != http.MethodGet {
		err := fmt.Errorf("unsupported method")
		http.Error(w, nod.Error(err).Error(), 405)
		return
	}

	id := r.URL.Query().Get("id")
	mt := gog_integration.Game

	vrDetails, err := getValueReader(vangogh_local_data.Details, mt)
	if err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
		return
	}

	det, err := vrDetails.Details(id)
	if err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
		return
	}

	dl, err := vangogh_local_data.FromDetails(det, mt, rxa)
	if err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
		return
	}

	// TODO: Only os, lang

	if err := json.NewEncoder(w).Encode(dl); err != nil {
		http.Error(w, nod.Error(err).Error(), 500)
	}
}
