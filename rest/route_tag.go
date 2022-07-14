package rest

import (
	"github.com/arelate/gog_integration"
	"github.com/arelate/vangogh_local_data"
	"github.com/boggydigital/coost"
	"github.com/boggydigital/nod"
	"net/http"
)

func RouteTag(w http.ResponseWriter, r *http.Request) {

	ids, err := vangogh_local_data.IdSetFromUrl(r.URL)
	if err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusBadRequest)
		return
	}

	tagId := vangogh_local_data.ValueFromUrl(r.URL, "tag-id")

	hc, err := coost.NewHttpClientFromFile(vangogh_local_data.AbsCookiePath(), gog_integration.GogHost)
	if err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
		return
	}

	switch r.Method {
	case http.MethodPut:
		PutTag(hc, ids, tagId, w)
	case http.MethodDelete:
		DeleteTag(hc, ids, tagId, w)
	default:
		http.Error(w, "unexpected tag method", http.StatusMethodNotAllowed)
	}
}
