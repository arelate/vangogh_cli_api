package cli

import (
	"fmt"
	"github.com/arelate/gog_integration"
	"github.com/arelate/vangogh_local_data"
	"github.com/boggydigital/nod"
	"net/url"
	"sort"
	"time"
)

var filterNewProductTypes = map[vangogh_local_data.ProductType]bool{
	vangogh_local_data.Orders: true,
	//not all licence-products have associated api-products-v1/api-products-v2,
	//so in some cases we won't get a meaningful information like a title
	vangogh_local_data.LicenceProducts: true,
	//both ApiProductsVx are not interesting since they correspond to store-products or account-products
	vangogh_local_data.ApiProductsV1: true,
	vangogh_local_data.ApiProductsV2: true,
}

var filterUpdatedProductTypes = map[vangogh_local_data.ProductType]bool{
	//most of the updates are price changes for a sale, not that interesting for recurring sync
	vangogh_local_data.StoreProducts: true,
	// wishlist-products are basically store-products, so see above
	vangogh_local_data.WishlistProducts: true,
	//meaningful updates for account products come from details, not account-products
	vangogh_local_data.AccountProducts: true,
	//same as above for those product types
	vangogh_local_data.ApiProductsV1: true,
	vangogh_local_data.ApiProductsV2: true,
}

func SummaryHandler(u *url.URL) error {
	since, err := vangogh_local_data.SinceFromUrl(u)
	if err != nil {
		return err
	}

	return Summary(
		vangogh_local_data.MediaFromUrl(u),
		since)
}

func Summary(mt gog_integration.Media, since int64) error {

	sa := nod.Begin("key changes since %s:", time.Unix(since, 0).Format("01/02 03:04PM"))
	defer sa.End()

	updates := make(map[string]map[string]bool, 0)

	for _, pt := range vangogh_local_data.LocalProducts() {

		if filterNewProductTypes[pt] {
			continue
		}

		vr, err := vangogh_local_data.NewReader(pt, mt)
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

	rxa, err := vangogh_local_data.ConnectReduxAssets(vangogh_local_data.TitleProperty)
	if err != nil {
		return sa.EndWithError(err)
	}

	summary := make(map[string][]string)

	for cat, items := range updates {
		summary[cat] = make([]string, 0, len(items))
		for id := range items {
			if title, ok := rxa.GetFirstVal(vangogh_local_data.TitleProperty, id); ok {
				summary[cat] = append(summary[cat], fmt.Sprintf("%s %s", id, title))
			}
		}
	}

	sa.EndWithSummary("", summary)

	return nil
}

func humanReadable(productTypes map[vangogh_local_data.ProductType]bool) []string {
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
