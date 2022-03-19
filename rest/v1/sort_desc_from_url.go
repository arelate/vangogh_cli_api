package v1

import (
	"github.com/arelate/vangogh_local_data"
	"net/url"
)

func sortDescFromUrl(u *url.URL) (string, bool) {
	q := u.Query()
	sort := q.Get("sort")
	if sort == "" {
		sort = vangogh_local_data.TitleProperty
	}
	desc := false
	if q.Get("desc") == "true" {
		desc = true
	}

	return sort, desc
}
