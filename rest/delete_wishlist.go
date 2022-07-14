package rest

import (
	"github.com/arelate/gog_integration"
	"github.com/arelate/vangogh_local_data"
	"github.com/boggydigital/nod"
	"golang.org/x/exp/maps"
	"net/http"
)

func DeleteWishlist(
	httpClient *http.Client,
	ids map[string]bool,
	mt gog_integration.Media,
	w http.ResponseWriter) {

	// DELETE /wishlist?id

	if pids, err := vangogh_local_data.RemoveFromLocalWishlist(maps.Keys(ids), mt, nil); err == nil {
		if err := gog_integration.RemoveFromWishlist(httpClient, pids...); err != nil {
			http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
			return
		}
	} else {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
