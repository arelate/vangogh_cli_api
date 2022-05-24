package cli

import (
	"github.com/arelate/vangogh_local_data"
	"github.com/boggydigital/nod"
	"net/url"
)

func LocalTagHandler(u *url.URL) error {
	idSet, err := vangogh_local_data.IdSetFromUrl(u)
	if err != nil {
		return err
	}

	return LocalTag(
		idSet,
		vangogh_local_data.ValueFromUrl(u, "operation"),
		vangogh_local_data.ValueFromUrl(u, "tag-name"))
}

func LocalTag(idSet map[string]bool, operation string, tagName string) error {

	lta := nod.NewProgress("%s local tag %s...", operation, tagName)
	defer lta.End()

	rxa, err := vangogh_local_data.ConnectReduxAssets(vangogh_local_data.LocalTagsProperty)
	if err != nil {
		return lta.EndWithError(err)
	}

	lta.TotalInt(len(idSet))

	for id := range idSet {
		switch operation {
		case "add":
			if err := rxa.AddVal(vangogh_local_data.LocalTagsProperty, id, tagName); err != nil {
				return err
			}
		case "remove":
			if err := rxa.CutVal(vangogh_local_data.LocalTagsProperty, id, tagName); err != nil {
				return err
			}
		}

		lta.Increment()
	}

	lta.EndWithResult("done")

	return nil
}
