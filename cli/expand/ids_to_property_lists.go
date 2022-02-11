package expand

import (
	"fmt"
	"github.com/arelate/vangogh_data"
	"github.com/boggydigital/gost"
	"github.com/boggydigital/kvas"
	"sort"
	"strings"
)

const (
	DefaultSort = vangogh_data.TitleProperty
	DefaultDesc = false
)

func IdsToPropertyLists(
	//ids []string,
	ids map[string]bool,
	propertyFilter map[string][]string,
	properties []string,
	rxa kvas.ReduxAssets) (map[string][]string, error) {

	propSet := gost.NewStrSetWith(properties...)
	propSet.Add(vangogh_data.TitleProperty)

	if rxa == nil {
		var err error
		rxa, err = vangogh_data.ConnectReduxAssets(propSet.All()...)
		if err != nil {
			return nil, err
		}
	}

	itps := make(map[string][]string)

	for id := range ids {
		itp, err := item(id, propertyFilter, propSet.All(), rxa)
		if err != nil {
			return itps, err
		}
		for idTitle, props := range itp {
			itps[idTitle] = props
		}
	}

	return itps, nil
}

func item(
	id string,
	propertyFilter map[string][]string,
	properties []string,
	rxa kvas.ReduxAssets) (map[string][]string, error) {

	if err := rxa.IsSupported(properties...); err != nil {
		return nil, err
	}

	title, ok := rxa.GetFirstVal(vangogh_data.TitleProperty, id)
	if !ok {
		return nil, nil
	}

	itp := make(map[string][]string)
	idTitle := fmt.Sprintf("%s %s", id, title)
	itp[idTitle] = make([]string, 0)

	sort.Strings(properties)

	for _, prop := range properties {
		if prop == vangogh_data.IdProperty ||
			prop == vangogh_data.TitleProperty {
			continue
		}
		values, ok := rxa.GetAllValues(prop, id)
		if !ok || len(values) == 0 {
			continue
		}
		filterValues := propertyFilter[prop]

		if len(values) > 1 && vangogh_data.IsPropertiesJoinPreferred(prop) {
			joinedValue := strings.Join(values, ",")
			if shouldSkip(joinedValue, filterValues) {
				continue
			}
			itp[idTitle] = append(itp[idTitle], fmt.Sprintf("%s:%s", prop, joinedValue))
			continue
		}

		for _, val := range values {
			if shouldSkip(val, filterValues) {
				continue
			}
			itp[idTitle] = append(itp[idTitle], fmt.Sprintf("%s:%s", prop, val))
		}
	}

	return itp, nil
}

func shouldSkip(value string, filterValues []string) bool {
	value = strings.ToLower(value)
	for _, fv := range filterValues {
		if strings.Contains(value, fv) {
			return false
		}
	}
	// this makes sure we don't filter values if there is no filter
	return len(filterValues) > 0
}
