package v1

import (
	"fmt"
	"github.com/arelate/vangogh_local_data"
	"github.com/boggydigital/nod"
	"net/http"
)

func GetImages(w http.ResponseWriter, r *http.Request) {

	// GET /v1/images?id

	nod.Log("GET %v", r.URL)

	if r.Method != http.MethodGet {
		err := fmt.Errorf("unsupported method")
		http.Error(w, nod.Error(err).Error(), 405)
		return
	}

	q := r.URL.Query()
	imageId := q.Get("id")
	if imageId == "" {
		err := fmt.Errorf("empty image id")
		http.Error(w, nod.Error(err).Error(), 400)
		return
	}
	if localImagePath := vangogh_local_data.AbsLocalImagePath(imageId); localImagePath != "" {
		http.ServeFile(w, r, localImagePath)
	} else {
		_ = nod.Error(fmt.Errorf("local image path for image id %s is empty", imageId))
		http.NotFound(w, r)
	}
}
