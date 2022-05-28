package v1

import (
	"github.com/arelate/gog_integration"
	"github.com/arelate/vangogh_local_data"
	"github.com/boggydigital/nod"
	"golang.org/x/exp/maps"
	"net/http"
)

func GetUpdates(w http.ResponseWriter, r *http.Request) {

	// GET /v1/updates?media&since-hours-ago&format

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

	if err := RefreshReduxAssets(vangogh_local_data.TitleProperty); err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
		return
	}

	sortedUpdates := make(map[string][]string)

	for section, ids := range updates {
		sortedIds, err := vangogh_local_data.SortIds(
			maps.Keys(ids),
			rxa,
			vangogh_local_data.DefaultSort,
			vangogh_local_data.DefaultDesc)
		if err != nil {
			http.Error(w, nod.Error(err).Error(), http.StatusMethodNotAllowed)
			return
		}
		sortedUpdates[section] = sortedIds
	}

	if err := encode(sortedUpdates, w, r); err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusMethodNotAllowed)
		return
	}
}
