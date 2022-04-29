package v1

import (
	"github.com/boggydigital/nod"
	"net/http"
)

func HandleFuncs() {
	v1PatternHandlers := map[string]http.Handler{
		"/v1/all_redux": nod.RequestLog(http.HandlerFunc(GetAllRedux)),
		"/v1/data":      nod.RequestLog(http.HandlerFunc(GetData)),
		"/v1/digest":    nod.RequestLog(http.HandlerFunc(GetDigest)),
		"/v1/downloads": nod.RequestLog(http.HandlerFunc(GetDownloads)),
		"/v1/keys":      nod.RequestLog(http.HandlerFunc(GetKeys)),
		"/v1/redux":     nod.RequestLog(http.HandlerFunc(GetRedux)),
		"/v1/search":    nod.RequestLog(http.HandlerFunc(Search)),
		"/v1/updates":   nod.RequestLog(http.HandlerFunc(GetUpdates)),
	}

	for p, h := range v1PatternHandlers {
		http.HandleFunc(p, h.ServeHTTP)
	}
}
