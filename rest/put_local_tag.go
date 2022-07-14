package rest

import (
	"github.com/arelate/vangogh_local_data"
	"github.com/boggydigital/nod"
	"golang.org/x/exp/maps"
	"net/http"
)

func PutLocalTag(
	ids map[string]bool,
	tag string,
	w http.ResponseWriter) {

	// PUT /local-tag?id&tag

	if err := vangogh_local_data.AddLocalTag(maps.Keys(ids), tag, nil); err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
