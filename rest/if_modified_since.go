package rest

import (
	"net/http"
	"time"
)

const (
	lastModifiedHeader    = "Last-Modified"
	ifModifiedSinceHeader = "If-Modified-Since"
)

func IfReduxModifiedSince(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// redux assets mod time is used to:
		// 1) set Last-Modified header
		// 2) check if content was modified since client cache
		if ramt, err := rxa.ReduxAssetsModTime(); err == nil {
			utcRamt := time.Unix(ramt, 0).UTC()
			w.Header().Set(lastModifiedHeader, utcRamt.Format(time.RFC1123))
			if imsh := r.Header.Get(ifModifiedSinceHeader); imsh != "" {
				if ims, err := time.Parse(time.RFC1123, imsh); err == nil {
					if utcRamt.Unix() <= ims.UTC().Unix() {
						w.WriteHeader(http.StatusNotModified)
						return
					}
				}
			}
		}
		next.ServeHTTP(w, r)
	})
}
