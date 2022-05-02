package v1

import (
	"fmt"
	"github.com/arelate/vangogh_local_data"
	"github.com/boggydigital/nod"
	"net/http"
	"strconv"
)

func GetKeys(w http.ResponseWriter, r *http.Request) {

	// GET /v1/keys?product-type&media&sort&desc&format

	if r.Method != http.MethodGet {
		err := fmt.Errorf("unsupported method")
		http.Error(w, nod.Error(err).Error(), http.StatusMethodNotAllowed)
		return
	}

	pt := vangogh_local_data.ProductTypeFromUrl(r.URL)
	mt := vangogh_local_data.MediaFromUrl(r.URL)
	count, err := strconv.Atoi(vangogh_local_data.ValueFromUrl(r.URL, "count"))
	if err != nil {
		count = -1
	}

	sort, desc := sortDescFromUrl(r.URL)
	if !vangogh_local_data.IsValidProperty(sort) {
		err := fmt.Errorf("invalid sort property %s", sort)
		http.Error(w, nod.Error(err).Error(), http.StatusBadRequest)
		return
	}

	vr, err := vangogh_local_data.NewReader(pt, mt)
	if err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
		return
	}

	if err := RefreshReduxAssets(sort, vangogh_local_data.TitleProperty); err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
		return
	}

	sortedIds, err := vangogh_local_data.SortIds(vr.Keys(), rxa, sort, desc)
	if err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
		return
	}

	if count > 0 && len(sortedIds) >= count {
		sortedIds = sortedIds[:count]
	}

	if err := encode(sortedIds, w, r); err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
	}
}
