package url_helpers

import (
	"github.com/arelate/vangogh_properties"
	"github.com/arelate/vangogh_urls"
	"github.com/boggydigital/gost"
	"github.com/boggydigital/kvas"
	"net/url"
)

func SlugIds(rxa kvas.ReduxAssets, slugs []string) (slugId gost.StrSet, err error) {

	if rxa == nil && len(slugs) > 0 {
		rxa, err = vangogh_properties.ConnectReduxAssets(vangogh_properties.SlugProperty)
		if err != nil {
			return nil, err
		}
	}

	if rxa != nil {
		if err := rxa.IsSupported(vangogh_properties.SlugProperty); err != nil {
			return nil, err
		}
	}

	idSet := gost.NewStrSet()
	for _, slug := range slugs {
		if slug != "" && rxa != nil {
			idSet.AddSet(rxa.Match(map[string][]string{vangogh_properties.SlugProperty: {slug}}, true))
		}
	}

	return idSet, nil
}

func IdSet(u *url.URL) (idSet gost.StrSet, err error) {

	idSet = gost.NewStrSetWith(vangogh_urls.UrlValues(u, "id")...)

	slugs := vangogh_urls.UrlValues(u, "slug")

	slugIds, err := SlugIds(nil, slugs)
	if err != nil {
		return idSet, err
	}
	idSet.AddSet(slugIds)

	return idSet, err
}
