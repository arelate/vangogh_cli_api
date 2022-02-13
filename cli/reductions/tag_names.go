package reductions

import (
	"fmt"
	"github.com/arelate/gog_integration"
	"github.com/arelate/vangogh_local_data"
	"github.com/boggydigital/nod"
)

func TagNames(mt gog_integration.Media) error {

	tna := nod.Begin(" %s...", vangogh_local_data.TagNameProperty)
	defer tna.End()

	vrAccountPage, err := vangogh_local_data.NewReader(vangogh_local_data.AccountPage, mt)
	if err != nil {
		return tna.EndWithError(err)
	}

	const fpId = "1"
	if !vrAccountPage.Has(fpId) {
		err := fmt.Errorf("%s doesn't contain page %s", vangogh_local_data.AccountPage, fpId)
		return tna.EndWithError(err)
	}

	firstPage, err := vrAccountPage.AccountPage(fpId)
	if err != nil {
		return tna.EndWithError(err)
	}

	tagNameEx, err := vangogh_local_data.ConnectReduxAssets(vangogh_local_data.TagNameProperty)
	if err != nil {
		return tna.EndWithError(err)
	}

	tagIdNames := make(map[string][]string, 0)

	for _, tag := range firstPage.Tags {
		tagIdNames[tag.Id] = []string{tag.Name}
	}

	if err := tagNameEx.BatchReplaceValues(vangogh_local_data.TagNameProperty, tagIdNames); err != nil {
		return tna.EndWithError(err)
	}

	tna.EndWithResult("done")

	return nil
}
