package itemizations

import (
	"github.com/arelate/vangogh_data"
	"github.com/boggydigital/gost"
	"github.com/boggydigital/nod"
)

func missingLocalFiles(
	all []string,
	localSet gost.StrSet,
	getById func(id string) ([]string, bool),
	exclude func(id string) bool,
	tpw nod.TotalProgressWriter) (vangogh_data.IdSet, error) {

	idSet := vangogh_data.NewIdSet()
	var err error

	if tpw != nil {
		tpw.TotalInt(len(all))
	}

	for _, id := range all {
		items, ok := getById(id)
		if !ok || len(items) == 0 {
			if tpw != nil {
				tpw.Increment()
			}
			continue
		}

		for _, item := range items {
			if exclude != nil && exclude(item) {
				continue
			}
			if !localSet.Has(item) {
				idSet.Add(id)
				break
			}
		}

		if tpw != nil {
			tpw.Increment()
		}
	}

	return idSet, err
}
