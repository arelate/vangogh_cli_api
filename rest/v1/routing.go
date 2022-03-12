package v1

import "net/http"

func HandleFuncs() {
	v1PatternHandlers := map[string]func(w http.ResponseWriter, r *http.Request){
		"/v1/keys":      GetKeys,
		"/v1/all_redux": GetAllRedux,
		"/v1/redux":     GetRedux,
		"/v1/data":      GetData,
		"/v1/images":    GetImages,
		"/v1/videos":    GetVideos,
		"/v1/search":    Search,
	}

	for p, h := range v1PatternHandlers {
		http.HandleFunc(p, h)
	}
}
