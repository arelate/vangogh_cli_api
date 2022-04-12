package cli

import (
	"github.com/arelate/gog_integration"
	"github.com/arelate/vangogh_local_data"
	"github.com/boggydigital/dolo"
	"github.com/boggydigital/nod"
	"net/url"
)

func GetItemsHandler(u *url.URL) error {
	since, err := vangogh_local_data.SinceFromUrl(u)
	if err != nil {
		return nil
	}

	return GetItems(
		vangogh_local_data.MediaFromUrl(u),
		since)
}

func GetItems(mt gog_integration.Media, since int64) error {

	gia := nod.NewProgress("getting description items...")
	defer gia.End()

	rxa, err := vangogh_local_data.ConnectReduxAssets(
		vangogh_local_data.TitleProperty,
		vangogh_local_data.DescriptionProperty)
	if err != nil {
		return err
	}

	vrStoreProducts, err := vangogh_local_data.NewReader(vangogh_local_data.StoreProducts, mt)
	if err != nil {
		return gia.EndWithError(err)
	}

	dl := dolo.DefaultClient

	all := vrStoreProducts.ModifiedAfter(since, false)

	gia.TotalInt(len(all))

	for _, id := range all {

		title, ok := rxa.GetFirstVal(vangogh_local_data.TitleProperty, id)
		if !ok {
			gia.Log("%s has no title", id)
			continue
		}
		description, ok := rxa.GetFirstVal(vangogh_local_data.DescriptionProperty, id)
		if !ok {
			continue
		}

		items := vangogh_local_data.ExtractDescItems(description)

		if len(items) < 1 {
			gia.Increment()
			continue
		}

		dia := nod.NewProgress(" %s %s", id, title)
		dia.TotalInt(len(items))

		urls := make([]*url.URL, 0, len(items))
		filenames := make([]string, 0, len(items))

		for _, itemUrl := range items {
			if u, err := url.Parse(itemUrl); err == nil {
				urls = append(urls, u)
				filenames = append(filenames, vangogh_local_data.AbsItemPath(u.Path))
			}
		}

		if err := dl.GetSet(urls, dolo.NewFileIndexSetter(filenames), dia); err != nil {
			gia.Error(err)
			continue
		}

		dia.EndWithResult("done")
		gia.Increment()
	}

	gia.EndWithResult("done")

	return nil
}
