package cli

import (
	"bytes"
	"encoding/json"
	"github.com/arelate/gog_atu"
	"github.com/arelate/vangogh_data"
	"github.com/boggydigital/kvas"
	"github.com/boggydigital/nod"
	"sort"
	"strconv"
)

func split(sourcePt vangogh_data.ProductType, mt gog_atu.Media, timestamp int64) error {

	spa := nod.NewProgress(" splitting %s (%s)...", sourcePt, mt)
	defer spa.End()

	vrPaged, err := vangogh_data.NewReader(sourcePt, mt)
	if err != nil {
		return err
	}

	modifiedIds := vrPaged.ModifiedAfter(timestamp, false)
	if len(modifiedIds) == 0 {
		spa.EndWithResult("unchanged")
		return nil
	}

	// split operates on pages and ids are expected to be numerical...
	intIds := make([]int, 0, len(modifiedIds))
	for _, id := range modifiedIds {
		inv, err := strconv.Atoi(id)
		if err == nil {
			intIds = append(intIds, inv)
		}
	}

	// ...however if some were not - just use the original modifiedIds set
	if len(intIds) == len(modifiedIds) {
		sort.Ints(intIds)
		modifiedIds = make([]string, 0, len(intIds))
		for _, id := range intIds {
			modifiedIds = append(modifiedIds, strconv.Itoa(id))
		}
	}

	spa.TotalInt(len(modifiedIds))

	for _, id := range modifiedIds {

		splitPt := vangogh_data.SplitProductType(sourcePt)

		productsGetter, err := vrPaged.ProductsGetter(id)

		if err != nil {
			return spa.EndWithError(err)
		}

		detailDstUrl, err := vangogh_data.AbsLocalProductTypeDir(splitPt, mt)
		if err != nil {
			return spa.EndWithError(err)
		}

		kvDetail, err := kvas.ConnectLocal(detailDstUrl, kvas.JsonExt)
		if err != nil {
			return spa.EndWithError(err)
		}

		products := productsGetter.GetProducts()

		if sourcePt == vangogh_data.Licences {
			spa.TotalInt(len(products))
		}

		for _, product := range products {
			buf := new(bytes.Buffer)
			if err := json.NewEncoder(buf).Encode(product); err != nil {
				return spa.EndWithError(err)
			}
			if err := kvDetail.Set(strconv.Itoa(product.GetId()), buf); err != nil {
				return spa.EndWithError(err)
			}
			if sourcePt == vangogh_data.Licences {
				spa.Increment()
			}
		}

		if sourcePt != vangogh_data.Licences {
			spa.Increment()
		}
	}

	spa.EndWithResult("done")

	return nil
}
