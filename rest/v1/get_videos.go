package v1

import (
	"github.com/arelate/vangogh_urls"
	"io"
	"net/http"
)

func GetVideos(w http.ResponseWriter, r *http.Request) {

	// GET /v1/videos?id

	q := r.URL.Query()
	videoId := q.Get("id")
	if videoId == "" {
		w.WriteHeader(400)
		_, _ = io.WriteString(w, "empty video-id")
		return
	}
	if localVideoPath := vangogh_urls.LocalVideoPath(videoId); localVideoPath != "" {
		http.ServeFile(w, r, localVideoPath)
	} else {
		w.WriteHeader(404)
	}
}
