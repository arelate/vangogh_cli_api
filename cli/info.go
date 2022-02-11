package cli

import (
	"github.com/arelate/vangogh_api/cli/expand"
	"github.com/arelate/vangogh_api/cli/url_helpers"
	"github.com/arelate/vangogh_data"
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
		vangogh_data.FlagFromUrl(u, "all-text"),
		vangogh_data.FlagFromUrl(u, "images"),
		vangogh_data.FlagFromUrl(u, "video-id"))
}

func Info(idSet gost.StrSet, allText, images, videoId bool) error {

	ia := nod.Begin("information:")
	defer ia.End()

	propSet := gost.NewStrSetWith(vangogh_data.TypesProperty)

	propSet.Add(vangogh_data.Text()...)
	if allText {
		propSet.Add(vangogh_data.AllText()...)
	}
	if images {
		propSet.Add(vangogh_data.ImageId()...)
	}
	if videoId {
		propSet.Add(vangogh_data.VideoId()...)
	}

	rxa, err := vangogh_data.ConnectReduxAssets(propSet.All()...)
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
