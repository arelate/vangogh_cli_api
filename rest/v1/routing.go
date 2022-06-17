package v1

import (
	"github.com/boggydigital/middleware"
	"github.com/boggydigital/nod"
	"net/http"
)

func HandleFuncs() {
	v1PatternHandlers := map[string]http.Handler{
		"/v1/all_redux": middleware.GetMethodOnly(nod.RequestLog(http.HandlerFunc(GetAllRedux))),
		"/v1/data":      middleware.GetMethodOnly(nod.RequestLog(http.HandlerFunc(GetData))),
		"/v1/has_data":  middleware.GetMethodOnly(nod.RequestLog(http.HandlerFunc(GetHasData))),
		"/v1/digest":    IfReduxModifiedSince(middleware.GetMethodOnly(nod.RequestLog(http.HandlerFunc(GetDigest)))),
		"/v1/downloads": middleware.GetMethodOnly(nod.RequestLog(http.HandlerFunc(GetDownloads))),
		"/v1/keys":      middleware.GetMethodOnly(nod.RequestLog(http.HandlerFunc(GetKeys))),
		"/v1/redux":     IfReduxModifiedSince(middleware.GetMethodOnly(nod.RequestLog(http.HandlerFunc(GetRedux)))),
		"/v1/has_redux": IfReduxModifiedSince(middleware.GetMethodOnly(nod.RequestLog(http.HandlerFunc(GetHasRedux)))),
		"/v1/search":    IfReduxModifiedSince(middleware.GetMethodOnly(nod.RequestLog(http.HandlerFunc(Search)))),
		"/v1/updates":   middleware.GetMethodOnly(nod.RequestLog(http.HandlerFunc(GetUpdates))),
	}

	for p, h := range v1PatternHandlers {
		http.HandleFunc(p, h.ServeHTTP)
	}
}
