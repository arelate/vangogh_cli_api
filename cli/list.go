package cli

import (
	"fmt"
	"github.com/arelate/gog_integration"
	"github.com/arelate/vangogh_local_data"
	"github.com/boggydigital/nod"
	"golang.org/x/exp/maps"
	"net/url"
	"time"
)

func ListHandler(u *url.URL) error {
	idSet, err := vangogh_local_data.IdSetFromUrl(u)
	if err != nil {
		return err
	}

	since, err := vangogh_local_data.SinceFromUrl(u)
	if err != nil {
		return err
	}

	return List(
		idSet,
		since,
		vangogh_local_data.ProductTypeFromUrl(u),
		vangogh_local_data.MediaFromUrl(u),
		vangogh_local_data.PropertiesFromUrl(u))
}

//List prints products of a certain type and media.
//Can be filtered to products that were created or modified since a certain time.
//Provided properties will be printed for each product (if supported) in addition to default ID, Title.
func List(
	idSet map[string]bool,
	modifiedSince int64,
	pt vangogh_local_data.ProductType,
	mt gog_integration.Media,
	properties []string) error {

	la := nod.Begin("listing %s...", pt)
	defer la.End()

	if !vangogh_local_data.IsValidProductType(pt) {
		return la.EndWithError(fmt.Errorf("can't list invalid product type %s", pt))
	}
	if !gog_integration.IsValidMedia(mt) {
		return la.EndWithError(fmt.Errorf("can't list invalid media %s", mt))
	}

	propSet := make(map[string]bool)
	for _, p := range properties {
		propSet[p] = true
	}

	//if no properties have been provided - print ID, Title
	if len(propSet) == 0 {
		propSet[vangogh_local_data.IdProperty] = true
		propSet[vangogh_local_data.TitleProperty] = true
	}

	//if Title property has not been provided - add it.
	//we'll always print the title.
	//same goes for sort-by property
	propSet[vangogh_local_data.TitleProperty] = true

	//rules for collecting IDs to print:
	//1. start with user provided IDs
	//2. if createdAfter has been provided - add products created since that time
	//3. if modifiedAfter has been provided - add products modified (not by creation!) since that time
	//4. if no IDs have been collected and the request have not provided createdAfter or modifiedAfter:
	// add all product IDs

	vr, err := vangogh_local_data.NewReader(pt, mt)
	if err != nil {
		return la.EndWithError(err)
	}

	if modifiedSince > 0 {
		for _, mid := range vr.ModifiedAfter(modifiedSince, false) {
			idSet[mid] = true
		}
		if len(idSet) == 0 {
			la.EndWithResult("no new or updated %s (%s) since %v\n", pt, mt, time.Unix(modifiedSince, 0).Format(time.Kitchen))
		}
	}

	if len(idSet) == 0 &&
		modifiedSince == 0 {
		for _, id := range vr.Keys() {
			idSet[id] = true
		}
	}

	itp, err := vangogh_local_data.PropertyListsFromIdSet(
		idSet,
		nil,
		vangogh_local_data.SupportedPropertiesOnly(pt, maps.Keys(propSet)),
		nil)

	if err != nil {
		return la.EndWithError(err)
	}

	la.EndWithSummary("", itp)

	return nil
}
