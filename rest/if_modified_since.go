package rest

import (
	"github.com/arelate/gog_integration"
	"github.com/arelate/vangogh_local_data"
	"net/http"
	"time"
)

const (
	lastModifiedHeader    = "Last-Modified"
	ifModifiedSinceHeader = "If-Modified-Since"
)

func isNotModified(w http.ResponseWriter, r *http.Request, since int64) bool {
	utcSince := time.Unix(since, 0).UTC()
	w.Header().Set(lastModifiedHeader, utcSince.Format(time.RFC1123))
	if imsh := r.Header.Get(ifModifiedSinceHeader); imsh != "" {
		if ims, err := time.Parse(time.RFC1123, imsh); err == nil {
			return utcSince.Unix() <= ims.UTC().Unix()
		}
	}
	return false
}

func IfReduxModifiedSince(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// redux assets mod time is used to:
		// 1) set Last-Modified header
		// 2) check if content was modified since client cache
		if ramt, err := rxa.ReduxAssetsModTime(); err == nil {
			if isNotModified(w, r, ramt) {
				w.WriteHeader(http.StatusNotModified)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}

func IfDataModifiedSince(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// data assets mod time is used to:
		// 1) set Last-Modified header
		// 2) check if content was modified since client cache
		//pt := vangogh_local_data.ProductTypeFromUrl(r.URL)
		pts := vangogh_local_data.ValuesFromUrl(r.URL, "product-type")
		mt := vangogh_local_data.MediaFromUrl(r.URL)
		notModified := true
		for _, pt := range pts {
			productType := vangogh_local_data.ParseProductType(pt)
			if vr, err := vangogh_local_data.NewReader(productType, mt); err == nil {
				if imt, err := vr.IndexCurrentModTime(); err == nil {
					notModified = notModified && isNotModified(w, r, imt)
				}
			}
		}
		if notModified {
			w.WriteHeader(http.StatusNotModified)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func IfDetailsModifiedSince(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// data assets mod time is used to:
		// 1) set Last-Modified header
		// 2) check if content was modified since client cache
		//pt := vangogh_local_data.ProductTypeFromUrl(r.URL)
		pt := vangogh_local_data.Details
		mt := gog_integration.Game
		if vr, err := vangogh_local_data.NewReader(pt, mt); err == nil {
			if imt, err := vr.IndexCurrentModTime(); err == nil {
				if isNotModified(w, r, imt) {
					w.WriteHeader(http.StatusNotModified)
					return
				}
			}
		}
		next.ServeHTTP(w, r)
	})
}
