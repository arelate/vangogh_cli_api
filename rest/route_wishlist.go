package rest

import (
	"github.com/arelate/vangogh_local_data"
	"github.com/boggydigital/nod"
	"net/http"
)

func RouteWishlist(w http.ResponseWriter, r *http.Request) {

	ids, err := vangogh_local_data.IdSetFromUrl(r.URL)
	if err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodPut:
		PutWishlist(ids, w)
	case http.MethodDelete:
		DeleteWishlist(ids, w)
	default:
		http.Error(w, "unexpected wishlist method", http.StatusMethodNotAllowed)
	}
}
