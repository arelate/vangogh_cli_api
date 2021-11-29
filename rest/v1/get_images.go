package v1

import (
	"github.com/arelate/vangogh_urls"
	"net/http"
)

func GetImages(w http.ResponseWriter, r *http.Request) {

	// GET /v1/images?id

	if r.Method != http.MethodGet {
		http.Error(w, "unsupported method", 405)
		return
	}

	q := r.URL.Query()
	imageId := q.Get("id")
	if imageId == "" {
		http.Error(w, "empty image-id", 400)
		return
	}
	if localImagePath := vangogh_urls.LocalImagePath(imageId); localImagePath != "" {
		http.ServeFile(w, r, localImagePath)
	} else {
		http.NotFound(w, r)
	}
}
