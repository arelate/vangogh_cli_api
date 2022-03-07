package cli

import (
	"fmt"
	"github.com/arelate/vangogh_cli_api/rest/v1"
	"github.com/arelate/vangogh_local_data"
	"github.com/boggydigital/nod"
	"net/http"
	"net/url"
	"strconv"
)

func ServeHandler(u *url.URL) error {
	portStr := vangogh_local_data.ValueFromUrl(u, "port")
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return err
	}

	return Serve(
		port,
		vangogh_local_data.FlagFromUrl(u, "stderr"))
}

func Serve(port int, stderr bool) error {

	if stderr {
		nod.EnableStdErrLogger()
		nod.DisableOutput(nod.StdOut)
	}

	sa := nod.Begin("serving at port %d...", port)
	defer sa.End()

	// API Version 1

	if err := v1.Init(); err != nil {
		return err
	}

	v1PatternHandlers := map[string]func(w http.ResponseWriter, r *http.Request){
		"/v1/keys":      v1.GetKeys,
		"/v1/all_redux": v1.GetAllRedux,
		"/v1/redux":     v1.GetRedux,
		"/v1/data":      v1.GetData,
		"/v1/images":    v1.GetImages,
		"/v1/videos":    v1.GetVideos,
	}

	for p, h := range v1PatternHandlers {
		http.HandleFunc(p, h)
	}

	return http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
