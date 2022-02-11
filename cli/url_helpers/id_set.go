package url_helpers

import (
	"github.com/arelate/vangogh_data"
	"github.com/boggydigital/gost"
	"github.com/boggydigital/kvas"
	"net/url"
)

func SlugIds(rxa kvas.ReduxAssets, slugs []string) (slugId gost.StrSet, err error) {

	if rxa == nil && len(slugs) > 0 {
		rxa, err = vangogh_data.ConnectReduxAssets(vangogh_data.SlugProperty)
		if err != nil {
			return nil, err
		}
	}

	if rxa != nil {
		if err := rxa.IsSupported(vangogh_data.SlugProperty); err != nil {
			return nil, err
		}
	}

	idSet := gost.NewStrSet()
	for _, slug := range slugs {
		if slug != "" && rxa != nil {
			idSet.AddSet(rxa.Match(map[string][]string{vangogh_data.SlugProperty: {slug}}, true))
		}
	}

	return idSet, nil
}

func IdSet(u *url.URL) (idSet gost.StrSet, err error) {

	idSet = gost.NewStrSetWith(vangogh_data.ValuesFromUrl(u, "id")...)

	slugs := vangogh_data.ValuesFromUrl(u, "slug")

	slugIds, err := SlugIds(nil, slugs)
	if err != nil {
		return idSet, err
	}
	idSet.AddSet(slugIds)

	return idSet, err
}
