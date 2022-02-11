package cli

import (
	"github.com/arelate/vangogh_api/cli/expand"
	"github.com/arelate/vangogh_data"
	"github.com/boggydigital/gost"
	"github.com/boggydigital/nod"
	"net/url"
)

func SearchHandler(u *url.URL) error {
	query := make(map[string][]string)

	for _, prop := range vangogh_data.Searchable() {
		if values := vangogh_data.ValuesFromUrl(u, prop); len(values) > 0 {
			query[prop] = values
		}
	}

	return Search(query)
}

func Search(query map[string][]string) error {

	sa := nod.Begin("searching...")
	defer sa.End()

	//prepare a list of all properties to load extracts for and
	//always start with a `title` property since it is printed for all matched item
	//(even if the match is for another property)
	propSet := gost.NewStrSetWith(vangogh_data.TitleProperty)
	for qp := range query {
		propSet.Add(qp)
	}

	rxa, err := vangogh_data.ConnectReduxAssets(propSet.All()...)
	if err != nil {
		return sa.EndWithError(err)
	}

	results := rxa.Match(query, true)

	//expand query properties for use in printInfo filter
	//since it's not aware of collapsed/expanded properties concept
	propertyFilter := make(map[string][]string, 0)

	//TODO: This needs to be restored
	//for prop, terms := range query {
	//	if vangogh_data.IsCollapsed(prop) {
	//		for _, ep := range vangogh_data.Expand(prop) {
	//			propertyFilter[ep] = terms
	//		}
	//	} else {
	//		propertyFilter[prop] = terms
	//	}
	//}

	if len(results) == 0 {
		sa.EndWithResult("no products found")
		return nil
	}

	itp, err := expand.IdsToPropertyLists(
		results,
		propertyFilter,
		//similarly for propertyFilter (see comment above) - expand all properties to display
		//TODO: restore this as well
		//vangogh_data.ExpandAll(propSet.All()),
		propSet.All(),
		rxa)

	if err != nil {
		return sa.EndWithError(err)
	}

	sa.EndWithSummary("found products:", itp)

	return nil
}
