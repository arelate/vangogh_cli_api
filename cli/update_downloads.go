package cli

import (
	"github.com/arelate/gog_media"
	"github.com/arelate/vangogh_api/cli/hours"
	"github.com/arelate/vangogh_api/cli/itemize"
	"github.com/arelate/vangogh_downloads"
	"github.com/arelate/vangogh_extracts"
	"github.com/arelate/vangogh_products"
	"github.com/arelate/vangogh_properties"
	"github.com/arelate/vangogh_urls"
	"github.com/boggydigital/nod"
	"net/url"
	"time"
)

func UpdateDownloadsHandler(u *url.URL) error {

	sha, err := hours.Atoi(vangogh_urls.UrlValue(u, "since-hours-ago"))
	if err != nil {
		return err
	}
	since := time.Now().Unix() - int64(sha*60*60)

	return UpdateDownloads(
		vangogh_urls.UrlMedia(u),
		vangogh_downloads.UrlOperatingSystems(u),
		vangogh_downloads.UrlDownloadTypes(u),
		vangogh_urls.UrlValues(u, "language-code"),
		since,
		vangogh_urls.UrlValue(u, "temp-directory"),
		vangogh_urls.UrlFlag(u, "updates-only"))
}

func UpdateDownloads(
	mt gog_media.Media,
	operatingSystems []vangogh_downloads.OperatingSystem,
	downloadTypes []vangogh_downloads.DownloadType,
	langCodes []string,
	since int64,
	tempDir string,
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
	updAccountProductIds, err := itemize.AccountProductsUpdates(mt)
	if err != nil {
		return uda.EndWithError(err)
	}

	//Additionally itemize required games for newly acquired DLCs
	requiredGamesForNewDLCs, err := itemize.RequiredAndIncluded(since)
	if err != nil {
		return uda.EndWithError(err)
	}

	updAccountProductIds.AddSet(requiredGamesForNewDLCs)

	//Additionally add modified details in case the sync was interrupted and
	//account-products doesn't have .IsNew or .Updates > 0 items
	modifiedDetails, err := itemize.Modified(since, vangogh_products.Details, mt)
	if err != nil {
		return uda.EndWithError(err)
	}

	updAccountProductIds.AddSet(modifiedDetails)

	if len(updAccountProductIds) == 0 {
		uda.EndWithResult("all downloads are up to date")
		return nil
	}

	//filter updAccountProductIds to products that have already been downloaded
	//note that this would exclude, for example, pre-order products automatic downloads
	if updatesOnly {
		exl, err := vangogh_extracts.NewList(vangogh_properties.SlugProperty)
		if err != nil {
			return uda.EndWithError(err)
		}

		for _, id := range updAccountProductIds.All() {
			ok, err := vangogh_downloads.ProductDownloaded(id, exl)
			if err != nil {
				return uda.EndWithError(err)
			}
			if !ok {
				updAccountProductIds.Hide(id)
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
		tempDir,
		false,
		true)
}
