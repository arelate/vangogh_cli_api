package cli

import (
	"fmt"
	"github.com/arelate/gog_integration"
	v1 "github.com/arelate/vangogh_cli_api/rest/v1"
	"github.com/arelate/vangogh_local_data"
	"github.com/boggydigital/nod"
	"net/url"
	"sort"
	"time"
)

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

	updates, err := v1.Updates(mt, since)
	if err != nil {
		return sa.EndWithError(err)
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
