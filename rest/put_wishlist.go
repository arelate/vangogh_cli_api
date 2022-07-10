package rest

import (
	"net/http"
)

func PutWishlist(ids map[string]bool, w http.ResponseWriter) {

	// PUT /wishlist?id

	//if err := cli.Wishlist(gog_integration.Game, maps.Keys(ids), nil); err != nil {
	//	http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
	//	return
	//}

	w.WriteHeader(http.StatusOK)
}
