package rest

import (
	"github.com/arelate/gog_integration"
	"github.com/arelate/vangogh_cli_api/cli"
	"github.com/boggydigital/nod"
	"golang.org/x/exp/maps"
	"net/http"
)

func DeleteWishlist(ids map[string]bool, w http.ResponseWriter) {

	// DELETE /wishlist?id

	if err := cli.Wishlist(gog_integration.Game, nil, maps.Keys(ids)); err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
