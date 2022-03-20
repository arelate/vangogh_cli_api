package v1

import (
	"fmt"
	"github.com/arelate/gog_integration"
	"github.com/arelate/vangogh_local_data"
	"github.com/boggydigital/nod"
	"net/http"
)

func GetUpdates(w http.ResponseWriter, r *http.Request) {

	// GET /v1/updates?media&since-hours-ago&format

	if r.Method != http.MethodGet {
		err := fmt.Errorf("unsupported method")
		http.Error(w, nod.Error(err).Error(), http.StatusMethodNotAllowed)
		return
	}

	mt := vangogh_local_data.MediaFromUrl(r.URL)
	if mt == gog_integration.Unknown {
		mt = gog_integration.Game
	}

	since, err := vangogh_local_data.SinceFromUrl(r.URL)
	if err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusMethodNotAllowed)
		return
	}

	updates, err := vangogh_local_data.Updates(mt, since)
	if err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusMethodNotAllowed)
		return
	}

	if err := encode(updates, w, r); err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusMethodNotAllowed)
		return
	}
}
