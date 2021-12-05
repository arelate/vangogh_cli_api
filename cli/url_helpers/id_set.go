package url_helpers

import (
	"github.com/arelate/vangogh_api/cli/lines"
	"github.com/arelate/vangogh_extracts"
	"github.com/arelate/vangogh_properties"
	"github.com/arelate/vangogh_urls"
	"github.com/boggydigital/gost"
	"net/url"
)

func SlugIds(exl *vangogh_extracts.ExtractsList, slugs []string) (slugId gost.StrSet, err error) {

	if exl == nil && len(slugs) > 0 {
		exl, err = vangogh_extracts.NewList(vangogh_properties.SlugProperty)
		if err != nil {
			return nil, err
		}
	}

	if exl != nil {
		if err := exl.AssertSupport(vangogh_properties.SlugProperty); err != nil {
			return nil, err
		}
	}

	idSet := gost.NewStrSet()
	for _, slug := range slugs {
		if slug != "" && exl != nil {
			idSet.Add(exl.Search(map[string][]string{vangogh_properties.SlugProperty: {slug}}, true)...)
		}
	}

	return idSet, nil
}

func IdSet(u *url.URL) (idSet gost.StrSet, err error) {

	idSet = gost.NewStrSetWith(vangogh_urls.UrlValues(u, "id")...)

	if vangogh_urls.UrlFlag(u, "read-ids") {
		pipedIds, err := lines.ReadPipedIds()
		if err != nil {
			return idSet, err
		}
		idSet.AddSet(pipedIds)
	}

	slugs := vangogh_urls.UrlValues(u, "slug")

	slugIds, err := SlugIds(nil, slugs)
	if err != nil {
		return idSet, err
	}
	idSet.AddSet(slugIds)

	return idSet, err
}
