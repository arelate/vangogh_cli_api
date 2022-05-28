package v1

import (
	"github.com/arelate/gog_integration"
	"github.com/arelate/vangogh_local_data"
	"github.com/boggydigital/nod"
	"net/http"
)

func GetDownloads(w http.ResponseWriter, r *http.Request) {

	// GET /v1/downloads?id&operating-system&language-code&format&media

	id := vangogh_local_data.ValueFromUrl(r.URL, "id")
	mt := vangogh_local_data.MediaFromUrl(r.URL)
	if mt == gog_integration.Unknown {
		mt = gog_integration.Game
	}

	vrDetails, err := vangogh_local_data.NewReader(vangogh_local_data.Details, mt)
	if err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
		return
	}

	det, err := vrDetails.Details(id)
	if err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
		return
	}

	dl := make(vangogh_local_data.DownloadsList, 0)

	if err := RefreshReduxAssets(vangogh_local_data.NativeLanguageNameProperty); err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
		return
	}

	if det != nil {
		dl, err = vangogh_local_data.FromDetails(det, mt, rxa)
		if err != nil {
			http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
			return
		}
	}

	os := vangogh_local_data.OperatingSystemsFromUrl(r.URL)
	lang := vangogh_local_data.ValuesFromUrl(r.URL, "language-code")

	dl = dl.Only(os, []vangogh_local_data.DownloadType{vangogh_local_data.AnyDownloadType}, lang)

	if err := encode(dl, w, r); err != nil {
		http.Error(w, nod.Error(err).Error(), http.StatusInternalServerError)
	}
}
