package rest

import (
	"net/http"
)

func DeleteWishlist(ids map[string]bool, w http.ResponseWriter) {

	// DELETE /wishlist?id

	//if err := cli.Wishlist(gog_integration.Game, nil, maps.Keys(ids)); err != nil {
	//	http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
	//	return
	//}

	w.WriteHeader(http.StatusOK)
}
