package cli

import (
	"github.com/arelate/vangogh_data"
	"github.com/boggydigital/gost"
	"github.com/boggydigital/nod"
	"net/url"
)

func InfoHandler(u *url.URL) error {
	idSet, err := vangogh_data.IdSetFromUrl(u)
	if err != nil {
		return err
	}

	return Info(
		idSet,
		vangogh_data.FlagFromUrl(u, "all-text"),
		vangogh_data.FlagFromUrl(u, "images"),
		vangogh_data.FlagFromUrl(u, "video-id"))
}

func Info(idSet vangogh_data.IdSet, allText, images, videoId bool) error {

	ia := nod.Begin("information:")
	defer ia.End()

	propSet := gost.NewStrSetWith(vangogh_data.TypesProperty)

	propSet.Add(vangogh_data.TextProperties()...)
	if allText {
		propSet.Add(vangogh_data.AllTextProperties()...)
	}
	if images {
		propSet.Add(vangogh_data.ImageIdProperties()...)
	}
	if videoId {
		propSet.Add(vangogh_data.VideoIdProperties()...)
	}

	rxa, err := vangogh_data.ConnectReduxAssets(propSet.All()...)
	if err != nil {
		return ia.EndWithError(err)
	}

	itp, err := vangogh_data.PropertyListsFromIdSet(
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
