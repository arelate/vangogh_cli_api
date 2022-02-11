package cli

import (
	"github.com/arelate/gog_atu"
	"github.com/arelate/vangogh_api/cli/hours"
	"github.com/arelate/vangogh_data"
	"github.com/boggydigital/gost"
	"github.com/boggydigital/nod"
	"github.com/boggydigital/wits"
	"net/url"
	"os"
	"strings"
	"time"
)

const (
	SyncOptionData             = "data"
	SyncOptionImages           = "images"
	SyncOptionScreenshots      = "screenshots"
	SyncOptionVideos           = "videos"
	SyncOptionDownloadsUpdates = "downloads-updates"
	negativePrefix             = "no-"
)

type syncOptions struct {
	data             bool
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
		data:             vangogh_data.FlagFromUrl(u, SyncOptionData),
		images:           vangogh_data.FlagFromUrl(u, SyncOptionImages),
		screenshots:      vangogh_data.FlagFromUrl(u, SyncOptionScreenshots),
		videos:           vangogh_data.FlagFromUrl(u, SyncOptionVideos),
		downloadsUpdates: vangogh_data.FlagFromUrl(u, SyncOptionDownloadsUpdates),
	}

	if vangogh_data.FlagFromUrl(u, "all") {
		so.data = !vangogh_data.FlagFromUrl(u, NegOpt(SyncOptionData))
		so.images = !vangogh_data.FlagFromUrl(u, NegOpt(SyncOptionImages))
		so.screenshots = !vangogh_data.FlagFromUrl(u, NegOpt(SyncOptionScreenshots))
		so.videos = !vangogh_data.FlagFromUrl(u, NegOpt(SyncOptionVideos))
		so.downloadsUpdates = !vangogh_data.FlagFromUrl(u, NegOpt(SyncOptionDownloadsUpdates))
	}

	return so
}

func SyncHandler(u *url.URL) error {
	syncOpts := initSyncOptions(u)

	sha, err := hours.Atoi(vangogh_data.ValueFromUrl(u, "since-hours-ago"))
	if err != nil {
		return err
	}

	return Sync(
		vangogh_data.MediaFromUrl(u),
		sha,
		syncOpts,
		vangogh_data.OperatingSystemsFromUrl(u),
		vangogh_data.DownloadTypesFromUrl(u),
		vangogh_data.ValuesFromUrl(u, "language-code"),
		vangogh_data.ValueFromUrl(u, "temp-directory"),
		vangogh_data.FlagFromUrl(u, "updates-only"))
}

func Sync(
	mt gog_atu.Media,
	sinceHoursAgo int,
	syncOpts *syncOptions,
	operatingSystems []vangogh_data.OperatingSystem,
	downloadTypes []vangogh_data.DownloadType,
	langCodes []string,
	tempDir string,
	updatesOnly bool) error {

	var syncStart int64
	if sinceHoursAgo > 0 {
		syncStart = time.Now().Unix() - int64(sinceHoursAgo*60*60)
	} else {
		syncStart = time.Now().Unix()
	}

	sa := nod.Begin("syncing source data...")
	defer sa.End()

	if syncOpts.data {
		//get array and paged data
		paData := vangogh_data.ArrayProducts()
		paData = append(paData, vangogh_data.PagedProducts()...)
		for _, pt := range paData {
			if err := GetData(gost.NewStrSet(), nil, pt, mt, syncStart, tempDir, false, false); err != nil {
				return sa.EndWithError(err)
			}
		}

		//get main - detail data
		for _, pt := range vangogh_data.DetailProducts() {

			var skipList wits.KeyValues

			if _, err := os.Stat(vangogh_data.AbsSkipListPath()); err == nil {
				slFile, err := os.Open(vangogh_data.AbsSkipListPath())
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

			if err := GetData(gost.NewStrSet(), skipIds, pt, mt, syncStart, tempDir, true, true); err != nil {
				return sa.EndWithError(err)
			}
		}

		//extract data
		if err := Extract(syncStart, mt, vangogh_data.Extracted()); err != nil {
			return sa.EndWithError(err)
		}
	}

	// get images
	if syncOpts.images {
		imageTypes := make([]vangogh_data.ImageType, 0, len(vangogh_data.AllImageTypes()))
		for _, it := range vangogh_data.AllImageTypes() {
			if !syncOpts.screenshots && it == vangogh_data.Screenshots {
				continue
			}
			imageTypes = append(imageTypes, it)
		}
		if err := GetImages(gost.NewStrSet(), imageTypes, true); err != nil {
			return sa.EndWithError(err)
		}
	}

	// get videos
	if syncOpts.videos {
		if err := GetVideos(gost.NewStrSet(), true); err != nil {
			return sa.EndWithError(err)
		}
	}

	// get downloads updates
	if syncOpts.downloadsUpdates {
		if err := UpdateDownloads(
			mt,
			operatingSystems,
			downloadTypes,
			langCodes,
			syncStart,
			tempDir,
			updatesOnly); err != nil {
			return sa.EndWithError(err)
		}
	}

	sa.EndWithResult("done")

	// print new or updated
	return Summary(mt, syncStart)
}
