package v1

import (
	"github.com/boggydigital/nod"
	"net/http"
)

func HandleFuncs() {
	v1PatternHandlers := map[string]func(w http.ResponseWriter, r *http.Request){
		"/v1/keys":      nod.RequestLog(GetKeys),
		"/v1/all_redux": nod.RequestLog(GetAllRedux),
		"/v1/redux":     nod.RequestLog(GetRedux),
		"/v1/data":      nod.RequestLog(GetData),
		"/v1/search":    nod.RequestLog(Search),
		"/v1/downloads": nod.RequestLog(GetDownloads),
	}

	for p, h := range v1PatternHandlers {
		http.HandleFunc(p, h)
	}
}
