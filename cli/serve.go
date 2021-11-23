package cli

import (
	"fmt"
	"github.com/arelate/vangogh_api/cli/url_helpers"
	"github.com/arelate/vangogh_api/rest/v1"
	"github.com/boggydigital/nod"
	"net/http"
	"net/url"
	"strconv"
)

func ServeHandler(u *url.URL) error {
	portStr := url_helpers.Value(u, "port")
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return err
	}

	return Serve(port)
}

func Serve(port int) error {
	sa := nod.Begin("serving at port %d...", port)
	defer sa.End()

	// API Version 1

	if err := v1.Init(); err != nil {
		return err
	}

	v1PatternHandlers := map[string]func(w http.ResponseWriter, r *http.Request){
		"/v1/indexes":  v1.GetIndexes,
		"/v1/extracts": v1.GetExtracts,
		"/v1/data":     v1.GetData,
		"/v1/images":   v1.GetImages,
		"/v1/videos":   v1.GetVideos,
	}

	for p, h := range v1PatternHandlers {
		http.HandleFunc(p, h)
	}

	return http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
