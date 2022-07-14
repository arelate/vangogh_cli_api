package rest

import (
	"github.com/arelate/vangogh_local_data"
	"github.com/boggydigital/nod"
	"net/http"
)

func PutTag(
	httpClient *http.Client,
	ids map[string]bool,
	tagId string,
	w http.ResponseWriter) {

	// PUT /tag?id&tag-id

	if err := vangogh_local_data.AddTag(httpClient, ids, tagId, rxa, nil); err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
