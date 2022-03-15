package cli

import (
	"fmt"
	"github.com/arelate/gog_integration"
	"github.com/arelate/vangogh_local_data"
	"github.com/boggydigital/gost"
	"github.com/boggydigital/nod"
	"net/url"
)

const (
	ownedSection    = "owned"
	notOwnedSection = "not owned"
)

func OwnedHandler(u *url.URL) error {
	idSet, err := vangogh_local_data.IdSetFromUrl(u)
	if err != nil {
		return err
	}

	return Owned(idSet)
}

func Owned(idSet vangogh_local_data.IdSet) error {

	oa := nod.Begin("checking ownership...")
	defer oa.End()

	ownedSet := vangogh_local_data.NewIdSet()
	propSet := gost.NewStrSetWith(
		vangogh_local_data.TitleProperty,
		vangogh_local_data.SlugProperty,
		vangogh_local_data.IncludesGamesProperty)

	rxa, err := vangogh_local_data.ConnectReduxAssets(propSet.All()...)
	if err != nil {
		return err
	}

	vrLicenceProducts, err := vangogh_local_data.NewReader(vangogh_local_data.LicenceProducts, gog_integration.Game)
	if err != nil {
		return err
	}

	for _, id := range idSet.All() {

		if vrLicenceProducts.Has(id) {
			ownedSet.Add(id)
			continue
		}

		includesGames, ok := rxa.GetAllUnchangedValues(vangogh_local_data.IncludesGamesProperty, id)
		if !ok || len(includesGames) == 0 {
			continue
		}

		ownAllIncludedGames := true
		for _, igId := range includesGames {
			ownAllIncludedGames = ownAllIncludedGames && vrLicenceProducts.Has(igId)
			if !ownAllIncludedGames {
				break
			}
		}

		if ownAllIncludedGames {
			ownedSet.Add(id)
		}
	}

	ownSummary := make(map[string][]string)
	ownSummary[ownedSection] = make([]string, 0, ownedSet.Len())
	for _, id := range ownedSet.All() {
		if title, ok := rxa.GetFirstVal(vangogh_local_data.TitleProperty, id); ok {
			ownSummary[ownedSection] = append(ownSummary[ownedSection], fmt.Sprintf("%s %s", id, title))
		}
	}

	notOwned := idSet.Except(ownedSet)

	ownSummary[notOwnedSection] = make([]string, 0, len(notOwned))
	for _, id := range notOwned {
		if title, ok := rxa.GetFirstVal(vangogh_local_data.TitleProperty, id); ok {
			ownSummary[notOwnedSection] = append(ownSummary[notOwnedSection], fmt.Sprintf("%s %s", id, title))
		}
	}

	oa.EndWithSummary("ownership results:", ownSummary)

	return nil
}
