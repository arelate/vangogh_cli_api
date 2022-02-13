package reductions

import (
	"github.com/arelate/gog_integration"
	"github.com/arelate/vangogh_local_data"
	"github.com/boggydigital/nod"
	"strconv"
	"time"
)

func Orders(modifiedAfter int64) error {

	oa := nod.NewProgress(" %s...", vangogh_local_data.GOGOrderDateProperty)
	defer oa.End()

	rxa, err := vangogh_local_data.ConnectReduxAssets(vangogh_local_data.GOGOrderDateProperty)
	if err != nil {
		return oa.EndWithError(err)
	}

	vrOrders, err := vangogh_local_data.NewReader(vangogh_local_data.Orders, gog_integration.Game)
	if err != nil {
		return oa.EndWithError(err)
	}

	gogOrderDates := make(map[string][]string, 0)

	var modifiedOrders []string
	if modifiedAfter > 0 {
		modifiedOrders = vrOrders.ModifiedAfter(modifiedAfter, false)
	} else {
		modifiedOrders = vrOrders.Keys()
	}

	oa.TotalInt(len(modifiedOrders))

	for _, orderId := range modifiedOrders {
		order, err := vrOrders.Order(orderId)
		if err != nil {
			return oa.EndWithError(err)
		}

		orderTimestamp, err := strconv.Atoi(orderId)
		if err != nil {
			return oa.EndWithError(err)
		}

		orderDate := time.Unix(int64(orderTimestamp), 0)

		for _, orderProduct := range order.Products {
			gogOrderDates[orderProduct.Id] = []string{orderDate.Format("2006-01-02 15:04:05")}
		}

		oa.Increment()
	}

	if err := rxa.BatchReplaceValues(vangogh_local_data.GOGOrderDateProperty, gogOrderDates); err != nil {
		return oa.EndWithError(err)
	}

	oa.EndWithResult("done")

	return nil
}
