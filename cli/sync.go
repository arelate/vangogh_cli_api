package cli

import (
	"github.com/arelate/gog_integration"
	"github.com/arelate/vangogh_local_data"
	"github.com/boggydigital/nod"
	"github.com/boggydigital/wits"
	"net/url"
	"os"
	"strings"
)

const (
	SyncOptionData             = "data"
	SyncOptionItems            = "items"
	SyncOptionImages           = "images"
	SyncOptionScreenshots      = "screenshots"
	SyncOptionVideos           = "videos"
	SyncOptionDownloadsUpdates = "downloads-Updates"
	negativePrefix             = "no-"
)

type syncOptions struct {
	data             bool
	items            bool
	images           bool
	screenshots      bool
	videos           bool
	downloadsUpdates bool
}

func NegOpt(option string) string {
	if !strings.HasPrefix(option, negativePrefix) {
		return negativePrefix + option
	}
	return option
}

func initSyncOptions(u *url.URL) *syncOptions {

	so := &syncOptions{
		data:             vangogh_local_data.FlagFromUrl(u, SyncOptionData),
		items:            vangogh_local_data.FlagFromUrl(u, SyncOptionItems),
		images:           vangogh_local_data.FlagFromUrl(u, SyncOptionImages),
		screenshots:      vangogh_local_data.FlagFromUrl(u, SyncOptionScreenshots),
		videos:           vangogh_local_data.FlagFromUrl(u, SyncOptionVideos),
		downloadsUpdates: vangogh_local_data.FlagFromUrl(u, SyncOptionDownloadsUpdates),
	}

	if vangogh_local_data.FlagFromUrl(u, "all") {
		so.data = !vangogh_local_data.FlagFromUrl(u, NegOpt(SyncOptionData))
		so.items = !vangogh_local_data.FlagFromUrl(u, NegOpt(SyncOptionItems))
		so.images = !vangogh_local_data.FlagFromUrl(u, NegOpt(SyncOptionImages))
		so.screenshots = !vangogh_local_data.FlagFromUrl(u, NegOpt(SyncOptionScreenshots))
		so.videos = !vangogh_local_data.FlagFromUrl(u, NegOpt(SyncOptionVideos))
		so.downloadsUpdates = !vangogh_local_data.FlagFromUrl(u, NegOpt(SyncOptionDownloadsUpdates))
	}

	return so
}

func SyncHandler(u *url.URL) error {
	syncOpts := initSyncOptions(u)

	since, err := vangogh_local_data.SinceFromUrl(u)
	if err != nil {
		return err
	}

	return Sync(
		vangogh_local_data.MediaFromUrl(u),
		since,
		syncOpts,
		vangogh_local_data.OperatingSystemsFromUrl(u),
		vangogh_local_data.DownloadTypesFromUrl(u),
		vangogh_local_data.ValuesFromUrl(u, "language-code"),
		vangogh_local_data.FlagFromUrl(u, "fast"))
}

func Sync(
	mt gog_integration.Media,
	since int64,
	syncOpts *syncOptions,
	operatingSystems []vangogh_local_data.OperatingSystem,
	downloadTypes []vangogh_local_data.DownloadType,
	langCodes []string,
	fast bool) error {

	sa := nod.Begin("syncing source data...")
	defer sa.End()

	if syncOpts.data {
		//get array and paged data
		var paData []vangogh_local_data.ProductType
		if fast {
			paData = vangogh_local_data.FastSyncProducts()
		} else {
			paData = append(
				vangogh_local_data.ArrayProducts(),
				vangogh_local_data.PagedProducts()...)
		}
		for _, pt := range paData {
			if err := GetData(map[string]bool{}, nil, pt, mt, since, false, false); err != nil {
				return sa.EndWithError(err)
			}
		}

		//get main - detail data
		for _, pt := range vangogh_local_data.DetailProducts() {

			var skipList wits.KeyValues

			if _, err := os.Stat(vangogh_local_data.AbsSkipListPath()); err == nil {
				slFile, err := os.Open(vangogh_local_data.AbsSkipListPath())
				if err != nil {
					slFile.Close()
					return sa.EndWithError(err)
				}

				skipList, err = wits.ReadKeyValues(slFile)
				slFile.Close()
				if err != nil {
					return sa.EndWithError(err)
				}
			}

			skipIds := skipList[pt.String()]
			if len(skipIds) > 0 {
				sa.Log("skipping %s ids: %v", pt, skipIds)
			} else {
				sa.Log("no skip list for %s", pt)
			}

			if err := GetData(map[string]bool{}, skipIds, pt, mt, since, true, true); err != nil {
				return sa.EndWithError(err)
			}
		}

		//reduce data
		if err := Reduce(mt, since, vangogh_local_data.ReduxProperties()); err != nil {
			return sa.EndWithError(err)
		}
	}

	// get items (embedded into descriptions)
	if syncOpts.items {
		if err := GetItems(map[string]bool{}, mt, since); err != nil {
			return sa.EndWithError(err)
		}
	}

	// get images
	if syncOpts.images {
		imageTypes := make([]vangogh_local_data.ImageType, 0, len(vangogh_local_data.AllImageTypes()))
		for _, it := range vangogh_local_data.AllImageTypes() {
			if !syncOpts.screenshots && it == vangogh_local_data.Screenshots {
				continue
			}
			imageTypes = append(imageTypes, it)
		}
		if err := GetImages(map[string]bool{}, imageTypes, true); err != nil {
			return sa.EndWithError(err)
		}
	}

	// get videos, unless fast sync was requested
	if syncOpts.videos && !fast {
		if err := GetVideos(map[string]bool{}, true); err != nil {
			return sa.EndWithError(err)
		}
	}

	// get downloads Updates
	if syncOpts.downloadsUpdates {
		if err := UpdateDownloads(
			mt,
			operatingSystems,
			downloadTypes,
			langCodes,
			since,
			false); err != nil {
			return sa.EndWithError(err)
		}
	}

	sa.EndWithResult("done")

	// print new or updated
	return Summary(mt, since)
}
