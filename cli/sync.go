package cli

import (
	"github.com/arelate/gog_media"
	"github.com/arelate/vangogh_api/cli/hours"
	"github.com/arelate/vangogh_api/cli/lines"
	"github.com/arelate/vangogh_api/cli/url_helpers"
	"github.com/arelate/vangogh_downloads"
	"github.com/arelate/vangogh_images"
	"github.com/arelate/vangogh_products"
	"github.com/arelate/vangogh_properties"
	"github.com/arelate/vangogh_urls"
	"github.com/boggydigital/gost"
	"github.com/boggydigital/nod"
	"net/url"
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
		data:             url_helpers.Flag(u, SyncOptionData),
		images:           url_helpers.Flag(u, SyncOptionImages),
		screenshots:      url_helpers.Flag(u, SyncOptionScreenshots),
		videos:           url_helpers.Flag(u, SyncOptionVideos),
		downloadsUpdates: url_helpers.Flag(u, SyncOptionDownloadsUpdates),
	}

	if url_helpers.Flag(u, "all") {
		so.data = !url_helpers.Flag(u, NegOpt(SyncOptionData))
		so.images = !url_helpers.Flag(u, NegOpt(SyncOptionImages))
		so.screenshots = !url_helpers.Flag(u, NegOpt(SyncOptionScreenshots))
		so.videos = !url_helpers.Flag(u, NegOpt(SyncOptionVideos))
		so.downloadsUpdates = !url_helpers.Flag(u, NegOpt(SyncOptionDownloadsUpdates))
	}

	return so
}

func SyncHandler(u *url.URL) error {
	mt := gog_media.Parse(url_helpers.Value(u, "media"))

	syncOpts := initSyncOptions(u)

	sha, err := hours.Atoi(url_helpers.Value(u, "since-hours-ago"))
	if err != nil {
		return err
	}

	operatingSystems := url_helpers.OperatingSystems(u)
	downloadTypes := url_helpers.DownloadTypes(u)
	langCodes := url_helpers.Values(u, "language-code")

	tempDir := url_helpers.Value(u, "temp-directory")

	updatesOnly := url_helpers.Flag(u, "updates-only")

	return Sync(
		mt,
		sha,
		syncOpts,
		operatingSystems,
		downloadTypes,
		langCodes,
		tempDir,
		updatesOnly)
}

func Sync(
	mt gog_media.Media,
	sinceHoursAgo int,
	syncOpts *syncOptions,
	operatingSystems []vangogh_downloads.OperatingSystem,
	downloadTypes []vangogh_downloads.DownloadType,
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
		paData := vangogh_products.Array()
		paData = append(paData, vangogh_products.Paged()...)
		for _, pt := range paData {
			if err := GetData(gost.NewStrSet(), nil, pt, mt, syncStart, tempDir, false, false); err != nil {
				return sa.EndWithError(err)
			}
		}

		//get main - detail data
		for _, pt := range vangogh_products.Detail() {
			denyIds := lines.Read(vangogh_urls.AbsSkiplistPath(pt))
			if err := GetData(gost.NewStrSet(), denyIds, pt, mt, syncStart, tempDir, true, true); err != nil {
				return sa.EndWithError(err)
			}
		}

		//extract data
		if err := Extract(syncStart, mt, vangogh_properties.Extracted()); err != nil {
			return sa.EndWithError(err)
		}
	}

	// get images
	if syncOpts.images {
		imageTypes := make([]vangogh_images.ImageType, 0, len(vangogh_images.All()))
		for _, it := range vangogh_images.All() {
			if !syncOpts.screenshots && it == vangogh_images.Screenshots {
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
