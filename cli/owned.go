package cli

import (
	"fmt"
	"github.com/arelate/gog_atu"
	"github.com/arelate/vangogh_api/cli/url_helpers"
	"github.com/arelate/vangogh_products"
	"github.com/arelate/vangogh_properties"
	"github.com/arelate/vangogh_values"
	"github.com/boggydigital/gost"
	"github.com/boggydigital/nod"
	"net/url"
)

const (
	ownedSection    = "owned"
	notOwnedSection = "not owned"
)

func OwnedHandler(u *url.URL) error {
	idSet, err := url_helpers.IdSet(u)
	if err != nil {
		return err
	}

	return Owned(idSet)
}

func Owned(idSet gost.StrSet) error {

	oa := nod.Begin("checking ownership...")
	defer oa.End()

	ownedSet := gost.NewStrSet()
	propSet := gost.NewStrSetWith(
		vangogh_properties.TitleProperty,
		vangogh_properties.SlugProperty,
		vangogh_properties.IncludesGamesProperty)

	rxa, err := vangogh_properties.ConnectReduxAssets(propSet.All()...)
	if err != nil {
		return err
	}

	vrLicenceProducts, err := vangogh_values.NewReader(vangogh_products.LicenceProducts, gog_atu.Game)
	if err != nil {
		return err
	}

	for _, id := range idSet.All() {

		if vrLicenceProducts.Has(id) {
			ownedSet.Add(id)
			continue
		}

		includesGames, ok := rxa.GetAllUnchangedValues(vangogh_properties.IncludesGamesProperty, id)
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
	for id := range ownedSet {
		if title, ok := rxa.GetFirstVal(vangogh_properties.TitleProperty, id); ok {
			ownSummary[ownedSection] = append(ownSummary[ownedSection], fmt.Sprintf("%s %s", id, title))
		}
	}

	notOwned := idSet.Except(ownedSet)

	ownSummary[notOwnedSection] = make([]string, 0, len(notOwned))
	for _, id := range notOwned {
		if title, ok := rxa.GetFirstVal(vangogh_properties.TitleProperty, id); ok {
			ownSummary[notOwnedSection] = append(ownSummary[notOwnedSection], fmt.Sprintf("%s %s", id, title))
		}
	}

	oa.EndWithSummary("ownership results:", ownSummary)

	return nil
}
