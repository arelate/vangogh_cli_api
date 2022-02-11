package extract

import (
	"github.com/arelate/gog_atu"
	"github.com/arelate/vangogh_data"
	"github.com/boggydigital/nod"
	"strconv"
	"time"
)

func Orders(modifiedAfter int64) error {

	oa := nod.NewProgress(" %s...", vangogh_data.GOGOrderDate)
	defer oa.End()

	rxa, err := vangogh_data.ConnectReduxAssets(vangogh_data.GOGOrderDate)
	if err != nil {
		return oa.EndWithError(err)
	}

	vrOrders, err := vangogh_data.NewReader(vangogh_data.Orders, gog_atu.Game)
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

	if err := rxa.BatchReplaceValues(vangogh_data.GOGOrderDate, gogOrderDates); err != nil {
		return oa.EndWithError(err)
	}

	oa.EndWithResult("done")

	return nil
}
