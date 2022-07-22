package cli

import (
	"fmt"
	"github.com/arelate/gog_integration"
	"github.com/arelate/vangogh_local_data"
	"github.com/boggydigital/nod"
	"golang.org/x/exp/maps"
	"net/url"
	"sort"
)

func SummaryHandler(u *url.URL) error {
	return Summary(vangogh_local_data.MediaFromUrl(u))
}

func Summarize(mt gog_integration.Media, since int64) error {

	updates, err := vangogh_local_data.Updates(mt, since)
	if err != nil {
		return err
	}

	if len(updates) == 0 {
		return nil
	}

	rxa, err := vangogh_local_data.ConnectReduxAssets(
		vangogh_local_data.LastSyncUpdatesProperty,
		vangogh_local_data.TitleProperty)
	if err != nil {
		return err
	}

	summary := make(map[string][]string)

	//set new values for each section
	for section, ids := range updates {
		sortedIds, err := vangogh_local_data.SortIds(
			maps.Keys(ids),
			rxa,
			vangogh_local_data.DefaultSort,
			vangogh_local_data.DefaultDesc)
		if err != nil {
			return err
		}
		summary[section] = sortedIds
	}

	//clean sections filled earlier that don't exist anymore
	for _, section := range rxa.Keys(vangogh_local_data.LastSyncUpdatesProperty) {
		if _, ok := updates[section]; ok {
			continue
		}
		summary[section] = nil
	}

	return rxa.BatchReplaceValues(vangogh_local_data.LastSyncUpdatesProperty, summary)

}

func Summary(mt gog_integration.Media) error {

	sa := nod.Begin("last sync summary:")
	defer sa.End()

	rxa, err := vangogh_local_data.ConnectReduxAssets(
		vangogh_local_data.LastSyncUpdatesProperty,
		vangogh_local_data.TitleProperty)
	if err != nil {
		return sa.EndWithError(err)
	}

	summary := make(map[string][]string)

	for _, section := range rxa.Keys(vangogh_local_data.LastSyncUpdatesProperty) {
		ids, _ := rxa.GetAllUnchangedValues(vangogh_local_data.LastSyncUpdatesProperty, section)
		for _, id := range ids {
			if title, ok := rxa.GetFirstVal(vangogh_local_data.TitleProperty, id); ok {
				summary[section] = append(summary[section], fmt.Sprintf("%s %s", id, title))
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
