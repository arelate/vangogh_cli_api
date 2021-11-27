package cli

import (
	"fmt"
	"github.com/arelate/gog_media"
	"github.com/arelate/vangogh_api/cli/hours"
	"github.com/arelate/vangogh_api/cli/url_helpers"
	"github.com/arelate/vangogh_extracts"
	"github.com/arelate/vangogh_products"
	"github.com/arelate/vangogh_properties"
	"github.com/arelate/vangogh_values"
	"github.com/boggydigital/nod"
	"net/url"
	"sort"
	"time"
)

var filterNewProductTypes = map[vangogh_products.ProductType]bool{
	vangogh_products.Orders: true,
	//not all licence-products have associated api-products-v1/api-products-v2,
	//so in some cases we won't get a meaningful information like a title
	vangogh_products.LicenceProducts: true,
	//both ApiProductsVx are not interesting since they correspond to store-products or account-products
	vangogh_products.ApiProductsV1: true,
	vangogh_products.ApiProductsV2: true,
}

var filterUpdatedProductTypes = map[vangogh_products.ProductType]bool{
	//most of the updates are price changes for a sale, not that interesting for recurring sync
	vangogh_products.StoreProducts: true,
	// wishlist-products are basically store-products, so see above
	vangogh_products.WishlistProducts: true,
	//meaningful updates for account products come from details, not account-products
	vangogh_products.AccountProducts: true,
	//same as above for those product types
	vangogh_products.ApiProductsV1: true,
	vangogh_products.ApiProductsV2: true,
}

func SummaryHandler(u *url.URL) error {
	sha, err := hours.Atoi(url_helpers.Value(u, "since-hours-ago"))
	if err != nil {
		return err
	}
	since := time.Now().Unix() - int64(sha*60*60)

	mt := gog_media.Parse(url_helpers.Value(u, "media"))

	return Summary(mt, since)
}

func Summary(mt gog_media.Media, since int64) error {

	sa := nod.Begin("key changes since %s:", time.Unix(since, 0).Format("01/02 03:04PM"))
	defer sa.End()

	updates := make(map[string]map[string]bool, 0)

	for _, pt := range vangogh_products.Local() {

		if filterNewProductTypes[pt] {
			continue
		}

		vr, err := vangogh_values.NewReader(pt, mt)
		if err != nil {
			return sa.EndWithError(err)
		}

		categorize(vr.CreatedAfter(since),
			fmt.Sprintf("new in %s", pt.HumanReadableString()),
			updates)

		if filterUpdatedProductTypes[pt] {
			continue
		}

		categorize(vr.ModifiedAfter(since, true),
			fmt.Sprintf("updated in %s", pt.HumanReadableString()),
			updates)
	}

	if len(updates) == 0 {
		sa.EndWithResult("no new or updated products")
		return nil
	}

	exl, err := vangogh_extracts.NewList(vangogh_properties.TitleProperty)
	if err != nil {
		return sa.EndWithError(err)
	}

	summary := make(map[string][]string)

	for cat, items := range updates {
		summary[cat] = make([]string, 0, len(items))
		for id := range items {
			if title, ok := exl.Get(vangogh_properties.TitleProperty, id); ok {
				summary[cat] = append(summary[cat], fmt.Sprintf("%s %s", id, title))
			}
		}
	}

	sa.EndWithSummary("", summary)

	return nil
}

func humanReadable(productTypes map[vangogh_products.ProductType]bool) []string {
	hrStrings := make(map[string]bool, 0)
	for key, ok := range productTypes {
		if !ok {
			continue
		}
		hrStrings[key.HumanReadableString()] = true
	}

	keys := make([]string, 0, len(hrStrings))
	for key := range hrStrings {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	return keys
}

func categorize(ids []string, cat string, updates map[string]map[string]bool) {
	for _, id := range ids {
		if updates[cat] == nil {
			updates[cat] = make(map[string]bool, 0)
		}
		updates[cat][id] = true
	}
}
