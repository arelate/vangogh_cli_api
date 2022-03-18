package cli

import (
	"github.com/arelate/gog_integration"
	"github.com/arelate/vangogh_cli_api/cli/itemizations"
	"github.com/arelate/vangogh_local_data"
	"github.com/boggydigital/nod"
	"net/url"
)

func UpdateDownloadsHandler(u *url.URL) error {

	since, err := vangogh_local_data.SinceFromUrl(u)
	if err != nil {
		return err
	}

	return UpdateDownloads(
		vangogh_local_data.MediaFromUrl(u),
		vangogh_local_data.OperatingSystemsFromUrl(u),
		vangogh_local_data.DownloadTypesFromUrl(u),
		vangogh_local_data.ValuesFromUrl(u, "language-code"),
		since,
		vangogh_local_data.FlagFromUrl(u, "updates-only"))
}

func UpdateDownloads(
	mt gog_integration.Media,
	operatingSystems []vangogh_local_data.OperatingSystem,
	downloadTypes []vangogh_local_data.DownloadType,
	langCodes []string,
	since int64,
	updatesOnly bool) error {

	uda := nod.Begin("itemizing updated downloads...")
	defer uda.End()

	//Here is a set of items we'll consider as updated for updating downloads:
	//1) account-products updates, all products that have .IsNew or .Updates > 0 -
	// basically items that GOG.com marked as new/updated
	//2) required games for newly acquired license-products -
	// making sure we update downloads for base product, when purchasing a DLC separately
	//3) modified details (since certain time) -
	// this accounts for interrupted sync, when we already processed account-products
	// updates (so .IsNew or .Updates > 0 won't be true anymore) and have updated
	// details as a result. This is somewhat excessive for general case, however would
	// allow us to capture all updated account-products at a price of some extra checks
	updAccountProductIds, err := itemizations.AccountProductsUpdates(mt)
	if err != nil {
		return uda.EndWithError(err)
	}

	//Additionally itemize required games for newly acquired DLCs
	requiredGamesForNewDLCs, err := itemizations.RequiredAndIncluded(since)
	if err != nil {
		return uda.EndWithError(err)
	}

	for rg := range requiredGamesForNewDLCs {
		updAccountProductIds[rg] = true
	}

	//Additionally add modified details in case the sync was interrupted and
	//account-products doesn't have .IsNew or .Updates > 0 items
	modifiedDetails, err := itemizations.Modified(since, vangogh_local_data.Details, mt)
	if err != nil {
		return uda.EndWithError(err)
	}

	for md := range modifiedDetails {
		updAccountProductIds[md] = true
	}

	if len(updAccountProductIds) == 0 {
		uda.EndWithResult("all downloads are up to date")
		return nil
	}

	//filter updAccountProductIds to products that have already been downloaded
	//note that this would exclude, for example, pre-order products automatic downloads
	if updatesOnly {
		rxa, err := vangogh_local_data.ConnectReduxAssets(vangogh_local_data.SlugProperty)
		if err != nil {
			return uda.EndWithError(err)
		}

		for id := range updAccountProductIds {
			ok, err := vangogh_local_data.IsProductDownloaded(id, rxa)
			if err != nil {
				return uda.EndWithError(err)
			}
			if !ok {
				delete(updAccountProductIds, id)
			}
		}
	}

	uda.EndWithResult("done")

	return GetDownloads(
		updAccountProductIds,
		mt,
		operatingSystems,
		downloadTypes,
		langCodes,
		false,
		true)
}
