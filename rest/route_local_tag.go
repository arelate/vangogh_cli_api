package rest

import (
	"github.com/arelate/vangogh_local_data"
	"github.com/boggydigital/nod"
	"net/http"
)

func RouteLocalTag(w http.ResponseWriter, r *http.Request) {

	ids, err := vangogh_local_data.IdSetFromUrl(r.URL)
	if err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusBadRequest)
		return
	}

	tag := vangogh_local_data.ValueFromUrl(r.URL, "tag")

	switch r.Method {
	case http.MethodPut:
		PutLocalTag(ids, tag, w)
	case http.MethodDelete:
		DeleteLocalTag(ids, tag, w)
	default:
		http.Error(w, "unexpected local-tag method", http.StatusMethodNotAllowed)
	}
}
