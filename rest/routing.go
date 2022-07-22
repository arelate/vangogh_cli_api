package rest

import (
	"github.com/boggydigital/middleware"
	"github.com/boggydigital/nod"
	"net/http"
)

func HandleFuncs() {
	patternHandlers := map[string]http.Handler{
		"/data":      middleware.GetMethodOnly(nod.RequestLog(http.HandlerFunc(GetData))),
		"/digest":    IfReduxModifiedSince(middleware.GetMethodOnly(nod.RequestLog(http.HandlerFunc(GetDigest)))),
		"/downloads": middleware.GetMethodOnly(nod.RequestLog(http.HandlerFunc(GetDownloads))),
		"/has_data":  middleware.GetMethodOnly(nod.RequestLog(http.HandlerFunc(GetHasData))),
		"/has_redux": IfReduxModifiedSince(middleware.GetMethodOnly(nod.RequestLog(http.HandlerFunc(GetHasRedux)))),
		"/local_tag": middleware.PatchMethodOnly(nod.RequestLog(http.HandlerFunc(PatchLocalTag))),
		"/redux":     IfReduxModifiedSince(middleware.GetMethodOnly(nod.RequestLog(http.HandlerFunc(GetRedux)))),
		"/search":    IfReduxModifiedSince(middleware.GetMethodOnly(nod.RequestLog(http.HandlerFunc(Search)))),
		"/tag":       middleware.PatchMethodOnly(nod.RequestLog(http.HandlerFunc(PatchTag))),
		"/updates":   middleware.GetMethodOnly(nod.RequestLog(http.HandlerFunc(GetUpdates))),
		"/wishlist":  nod.RequestLog(http.HandlerFunc(RouteWishlist)),
	}

	for p, h := range patternHandlers {
		http.HandleFunc(p, h.ServeHTTP)
	}
}
