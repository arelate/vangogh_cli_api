package cli

import (
	"fmt"
	"github.com/arelate/gog_media"
	"github.com/arelate/vangogh_api/cli/expand"
	"github.com/arelate/vangogh_api/cli/hours"
	"github.com/arelate/vangogh_api/cli/url_helpers"
	"github.com/arelate/vangogh_products"
	"github.com/arelate/vangogh_properties"
	"github.com/arelate/vangogh_values"
	"github.com/boggydigital/gost"
	"github.com/boggydigital/nod"
	"net/url"
	"time"
)

func ListHandler(u *url.URL) error {
	idSet, err := url_helpers.IdSet(u)
	if err != nil {
		return err
	}

	sha, err := hours.Atoi(url_helpers.Value(u, "since-hours-ago"))
	if err != nil {
		return err
	}

	var since int64 = 0
	if sha > 0 {
		since = time.Now().Add(-time.Hour * time.Duration(sha)).Unix()
	}

	pt := vangogh_products.Parse(url_helpers.Value(u, "product-type"))
	mt := gog_media.Parse(url_helpers.Value(u, "media"))

	properties := url_helpers.Values(u, "property")

	return List(idSet, since, pt, mt, properties)
}

//List prints products of a certain type and media.
//Can be filtered to products that were created or modified since a certain time.
//Provided properties will be printed for each product (if supported) in addition to default ID, Title.
func List(
	idSet gost.StrSet,
	modifiedSince int64,
	pt vangogh_products.ProductType,
	mt gog_media.Media,
	properties []string) error {

	la := nod.Begin("listing %s...", pt)
	defer la.End()

	if !vangogh_products.Valid(pt) {
		return la.EndWithError(fmt.Errorf("can't list invalid product type %s", pt))
	}
	if !gog_media.Valid(mt) {
		return la.EndWithError(fmt.Errorf("can't list invalid media %s", mt))
	}

	propSet := gost.NewStrSetWith(properties...)

	//if no properties have been provided - print ID, Title
	if propSet.Len() == 0 {
		propSet.Add(
			vangogh_properties.IdProperty,
			vangogh_properties.TitleProperty)
	}

	//if Title property has not been provided - add it.
	//we'll always print the title.
	//same goes for sort-by property
	propSet.Add(vangogh_properties.TitleProperty)

	//rules for collecting IDs to print:
	//1. start with user provided IDs
	//2. if createdAfter has been provided - add products created since that time
	//3. if modifiedAfter has been provided - add products modified (not by creation!) since that time
	//4. if no IDs have been collected and the request have not provided createdAfter or modifiedAfter:
	// add all product IDs

	vr, err := vangogh_values.NewReader(pt, mt)
	if err != nil {
		return la.EndWithError(err)
	}

	if modifiedSince > 0 {
		idSet.Add(vr.ModifiedAfter(modifiedSince, false)...)
		if idSet.Len() == 0 {
			la.EndWithResult("no new or updated %s (%s) since %v\n", pt, mt, time.Unix(modifiedSince, 0).Format(time.Kitchen))
		}
	}

	if idSet.Len() == 0 &&
		modifiedSince == 0 {
		idSet.Add(vr.All()...)
	}

	itp, err := expand.IdsToPropertyLists(
		idSet.All(),
		nil,
		vangogh_properties.Supported(pt, properties),
		nil)

	if err != nil {
		return la.EndWithError(err)
	}

	la.EndWithSummary("", itp)

	return nil
}
