package rest

import (
	"github.com/arelate/gog_integration"
	"github.com/arelate/vangogh_cli_api/cli"
	"github.com/boggydigital/nod"
	"golang.org/x/exp/maps"
	"net/http"
)

func PutWishlist(ids map[string]bool, w http.ResponseWriter) {

	// PUT /wishlist?id

	if err := cli.Wishlist(gog_integration.Game, maps.Keys(ids), nil); err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
