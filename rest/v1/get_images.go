package v1

import (
	"github.com/arelate/vangogh_urls"
	"io"
	"net/http"
)

func GetImages(w http.ResponseWriter, r *http.Request) {

	// GET /v1/images?id

	q := r.URL.Query()
	imageId := q.Get("id")
	if imageId == "" {
		w.WriteHeader(400)
		_, _ = io.WriteString(w, "empty image-id")
		return
	}
	if localImagePath := vangogh_urls.LocalImagePath(imageId); localImagePath != "" {
		http.ServeFile(w, r, localImagePath)
	} else {
		w.WriteHeader(404)
	}
}
