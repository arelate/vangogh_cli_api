package cli

import (
	"github.com/arelate/vangogh_api/cli/expand"
	"github.com/arelate/vangogh_api/cli/url_helpers"
	"github.com/arelate/vangogh_properties"
	"github.com/arelate/vangogh_urls"
	"github.com/boggydigital/gost"
	"github.com/boggydigital/nod"
	"net/url"
)

func InfoHandler(u *url.URL) error {
	idSet, err := url_helpers.IdSet(u)
	if err != nil {
		return err
	}

	return Info(
		idSet,
		vangogh_urls.UrlFlag(u, "all-text"),
		vangogh_urls.UrlFlag(u, "images"),
		vangogh_urls.UrlFlag(u, "video-id"))
}

func Info(idSet gost.StrSet, allText, images, videoId bool) error {

	ia := nod.Begin("information:")
	defer ia.End()

	propSet := gost.NewStrSetWith(vangogh_properties.TypesProperty)

	propSet.Add(vangogh_properties.Text()...)
	if allText {
		propSet.Add(vangogh_properties.AllText()...)
	}
	if images {
		propSet.Add(vangogh_properties.ImageId()...)
	}
	if videoId {
		propSet.Add(vangogh_properties.VideoId()...)
	}

	rxa, err := vangogh_properties.ConnectReduxAssets(propSet.All()...)
	if err != nil {
		return ia.EndWithError(err)
	}

	itp, err := expand.IdsToPropertyLists(
		idSet,
		nil,
		propSet.All(),
		rxa)

	if err != nil {
		return ia.EndWithError(err)
	}

	ia.EndWithSummary("", itp)

	return nil
}
