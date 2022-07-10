package rest

import (
	"github.com/boggydigital/middleware"
	"github.com/boggydigital/nod"
	"net/http"
)

func HandleFuncs() {
	patternHandlers := map[string]http.Handler{
		"/all_redux": middleware.GetMethodOnly(nod.RequestLog(http.HandlerFunc(GetAllRedux))),
		"/data":      middleware.GetMethodOnly(nod.RequestLog(http.HandlerFunc(GetData))),
		"/has_data":  middleware.GetMethodOnly(nod.RequestLog(http.HandlerFunc(GetHasData))),
		"/digest":    IfReduxModifiedSince(middleware.GetMethodOnly(nod.RequestLog(http.HandlerFunc(GetDigest)))),
		"/downloads": middleware.GetMethodOnly(nod.RequestLog(http.HandlerFunc(GetDownloads))),
		"/keys":      middleware.GetMethodOnly(nod.RequestLog(http.HandlerFunc(GetKeys))),
		"/redux":     IfReduxModifiedSince(middleware.GetMethodOnly(nod.RequestLog(http.HandlerFunc(GetRedux)))),
		"/has_redux": IfReduxModifiedSince(middleware.GetMethodOnly(nod.RequestLog(http.HandlerFunc(GetHasRedux)))),
		"/search":    IfReduxModifiedSince(middleware.GetMethodOnly(nod.RequestLog(http.HandlerFunc(Search)))),
		"/updates":   middleware.GetMethodOnly(nod.RequestLog(http.HandlerFunc(GetUpdates))),
	}

	for p, h := range patternHandlers {
		http.HandleFunc(p, h.ServeHTTP)
	}
}
