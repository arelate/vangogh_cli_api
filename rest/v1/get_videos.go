package v1

import (
	"github.com/arelate/vangogh_urls"
	"net/http"
)

func GetVideos(w http.ResponseWriter, r *http.Request) {

	// GET /v1/videos?id

	if r.Method != http.MethodGet {
		http.Error(w, "unsupported method", 405)
		return
	}

	q := r.URL.Query()
	videoId := q.Get("id")
	if videoId == "" {
		http.Error(w, "empty video-id", 400)
		return
	}
	if localVideoPath := vangogh_urls.LocalVideoPath(videoId); localVideoPath != "" {
		http.ServeFile(w, r, localVideoPath)
	} else {
		http.NotFound(w, r)
	}
}
