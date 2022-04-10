package cli

import (
	"github.com/arelate/gog_integration"
	"github.com/arelate/vangogh_local_data"
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

	vrStoreProducts, err := vangogh_local_data.NewReader(vangogh_local_data.StoreProducts, mt)
	if err != nil {
		return gia.EndWithError(err)
	}
	vrApiProductsV2, err := vangogh_local_data.NewReader(vangogh_local_data.ApiProductsV2, mt)
	if err != nil {
		return gia.EndWithError(err)
	}
	vrApiProductsV1, err := vangogh_local_data.NewReader(vangogh_local_data.ApiProductsV1, mt)
	if err != nil {
		return gia.EndWithError(err)
	}

	all := vrStoreProducts.ModifiedAfter(since, false)

	gia.TotalInt(len(all))

	for _, id := range all {

		var items []string
		var title string

		if apv2, err := vrApiProductsV2.ApiProductV2(id); err != nil {
			gia.Error(err)
			gia.Increment()
			continue
		} else if apv2 != nil {
			items = apv2.GetDescriptionItems()
			title = apv2.GetTitle()
		} else {
			if apv1, err := vrApiProductsV1.ApiProductV1(id); err != nil {
				gia.Error(err)
				gia.Increment()
				continue
			} else {
				items = apv1.GetDescriptionItems()
				title = apv1.GetTitle()
			}
		}

		if len(items) < 1 {
			gia.Increment()
			continue
		}

		dia := nod.NewProgress("%s %s", id, title)
		dia.TotalInt(len(items))

		//for _, itemUrl := range items {
		//	u, err := url.Parse(itemUrl)
		//	if err != nil {
		//		dia.Error(err)
		//		dia.Increment()
		//		continue
		//	}
		//
		//	dia.Increment()
		//}

		dia.EndWithResult("done")
		gia.Increment()
	}

	gia.EndWithResult("done")

	return nil
}
