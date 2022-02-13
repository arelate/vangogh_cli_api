package cli

import (
	"fmt"
	"github.com/arelate/vangogh_local_data"
	"github.com/boggydigital/gost"
	"github.com/boggydigital/nod"
	"net/url"
)

func DigestHandler(u *url.URL) error {
	return Digest(
		vangogh_local_data.PropertyFromUrl(u))
}

func Digest(property string) error {

	da := nod.Begin("digesting...")
	defer da.End()

	rxa, err := vangogh_local_data.ConnectReduxAssets(property)
	if err != nil {
		return err
	}

	distValues := make(map[string]int, 0)

	for _, id := range rxa.Keys(property) {
		values, ok := rxa.GetAllValues(property, id)
		if !ok || len(values) == 0 {
			continue
		}

		for _, val := range values {
			if val == "" {
				continue
			}
			distValues[val] = distValues[val] + 1
		}
	}

	keys := make([]string, 0, len(distValues))
	for key := range distValues {
		keys = append(keys, key)
	}

	_, sorted := gost.NewIntSortedStrSetWith(distValues, vangogh_local_data.DefaultDesc)

	summary := make(map[string][]string)
	summary[""] = make([]string, 0, len(sorted))
	for _, key := range sorted {
		summary[""] = append(summary[""], fmt.Sprintf("%s: %d items", key, distValues[key]))
	}

	da.EndWithSummary("digested properties:", summary)

	return nil
}
