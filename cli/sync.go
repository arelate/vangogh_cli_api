package cli

import (
	"github.com/arelate/gog_integration"
	"github.com/arelate/vangogh_cli_api/cli/reductions"
	"github.com/arelate/vangogh_local_data"
	"github.com/boggydigital/nod"
	"github.com/boggydigital/wits"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	SyncOptionData             = "data"
	SyncOptionItems            = "items"
	SyncOptionImages           = "images"
	SyncOptionScreenshots      = "screenshots"
	SyncOptionVideos           = "videos"
	SyncOptionThumbnails       = "thumbnails"
	SyncOptionDownloadsUpdates = "downloads-Updates"
	negativePrefix             = "no-"
)

type syncOptions struct {
	data             bool
	items            bool
	images           bool
	screenshots      bool
	videos           bool
	thumbnails       bool
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
		thumbnails:       vangogh_local_data.FlagFromUrl(u, SyncOptionThumbnails),
		downloadsUpdates: vangogh_local_data.FlagFromUrl(u, SyncOptionDownloadsUpdates),
	}

	if vangogh_local_data.FlagFromUrl(u, "all") {
		so.data = !vangogh_local_data.FlagFromUrl(u, NegOpt(SyncOptionData))
		so.items = !vangogh_local_data.FlagFromUrl(u, NegOpt(SyncOptionItems))
		so.images = !vangogh_local_data.FlagFromUrl(u, NegOpt(SyncOptionImages))
		so.screenshots = !vangogh_local_data.FlagFromUrl(u, NegOpt(SyncOptionScreenshots))
		so.videos = !vangogh_local_data.FlagFromUrl(u, NegOpt(SyncOptionVideos))
		so.thumbnails = !vangogh_local_data.FlagFromUrl(u, NegOpt(SyncOptionThumbnails))
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
		vangogh_local_data.ValuesFromUrl(u, "language-code"))
}

func Sync(
	mt gog_integration.Media,
	since int64,
	syncOpts *syncOptions,
	operatingSystems []vangogh_local_data.OperatingSystem,
	downloadTypes []vangogh_local_data.DownloadType,
	langCodes []string) error {

	sa := nod.Begin("syncing source data...")
	defer sa.End()

	syncEventsRxa, err := vangogh_local_data.ConnectReduxAssets(vangogh_local_data.SyncEventsProperty)
	if err != nil {
		return sa.EndWithError(err)
	}

	syncStart := since
	if syncStart == 0 {
		syncStart = time.Now().Unix()
	}

	if err := syncEventsRxa.AddVal(
		vangogh_local_data.SyncEventsProperty,
		vangogh_local_data.SyncStartKey,
		strconv.Itoa(int(syncStart))); err != nil {
		return sa.EndWithError(err)
	}

	if syncOpts.data {
		//get array and paged data
		paData := append(vangogh_local_data.ArrayProducts(),
			vangogh_local_data.PagedProducts()...)

		for _, pt := range paData {
			if err := GetData(map[string]bool{}, nil, pt, mt, since, false, false); err != nil {
				return sa.EndWithError(err)
			}
		}

		//get GOG.com main - detail data
		if err := getDetailData(vangogh_local_data.GOGDetailProducts(), mt, since); err != nil {
			return sa.EndWithError(err)
		}

		//reduce Steam AppId
		if err := Reduce(mt, since, []string{vangogh_local_data.SteamAppIdProperty}, false); err != nil {
			return sa.EndWithError(err)
		}

		//get Steam main - detail data
		//this needs to happen after reduce, since Steam AppId - GOG.com ProductId
		//connection is established at reduce. And the earlier data set cannot be retrieved post reduce,
		//since SteamAppList is fetched with initial data
		if err := getDetailData(vangogh_local_data.SteamDetailProducts(), mt, since); err != nil {
			return sa.EndWithError(err)
		}

		// finally, reduce all properties
		if err := Reduce(mt, since, vangogh_local_data.ReduxProperties(), false); err != nil {
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

		pr := nod.Begin("reducing post image download...")
		if err := reductions.DehydratedImages(); err != nil {
			return pr.EndWithError(err)
		}
		pr.EndWithResult("done")
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

	// get videos
	if syncOpts.videos {
		if err := GetVideos(map[string]bool{}, true); err != nil {
			return sa.EndWithError(err)
		}
	}

	// get thumbnails
	if syncOpts.thumbnails {
		if err := GetThumbnails(map[string]bool{}, true); err != nil {
			return sa.EndWithError(err)
		}
	}

	sa.EndWithResult("done")

	if err := syncEventsRxa.AddVal(
		vangogh_local_data.SyncEventsProperty,
		vangogh_local_data.SyncCompleteKey,
		strconv.Itoa(int(time.Now().Unix()))); err != nil {
		return sa.EndWithError(err)
	}

	// summarize sync updates
	if err := Summarize(mt, syncStart); err != nil {
		return sa.EndWithError(err)
	}

	// print new, updated
	return Summary(mt)
}

func getDetailData(pts []vangogh_local_data.ProductType, mt gog_integration.Media, since int64) error {
	for _, pt := range pts {

		var skipList wits.KeyValues

		if _, err := os.Stat(vangogh_local_data.AbsSkipListPath()); err == nil {
			slFile, err := os.Open(vangogh_local_data.AbsSkipListPath())
			if err != nil {
				slFile.Close()
				return err
			}

			skipList, err = wits.ReadKeyValues(slFile)
			slFile.Close()
			if err != nil {
				return err
			}
		}

		skipIds := skipList[pt.String()]
		if err := GetData(map[string]bool{}, skipIds, pt, mt, since, true, true); err != nil {
			return err
		}
	}

	return nil
}
