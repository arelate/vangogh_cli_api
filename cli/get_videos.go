package cli

import (
	"github.com/arelate/vangogh_cli_api/cli/itemizations"
	"github.com/arelate/vangogh_local_data"
	"github.com/boggydigital/nod"
	"github.com/boggydigital/yet/yet"
	"github.com/boggydigital/yt_urls"
	"net/http"
	"net/url"
	"os/exec"
	"path/filepath"
)

const (
	missingStr = "missing"
)

func GetVideosHandler(u *url.URL) error {
	idSet, err := vangogh_local_data.IdSetFromUrl(u)
	if err != nil {
		return err
	}

	ffmpegCmd := vangogh_local_data.ValueFromUrl(u, "ffmpeg-cmd")

	if ffmpegCmd == "" {
		if path, err := exec.LookPath("ffmpeg"); err == nil {
			ffmpegCmd = path
		}
	}

	return GetVideos(
		idSet,
		ffmpegCmd,
		vangogh_local_data.FlagFromUrl(u, "missing"))
}

func GetVideos(idSet map[string]bool, ffmpegCmd string, missing bool) error {

	gva := nod.NewProgress("getting videos...")
	defer gva.End()

	rxa, err := vangogh_local_data.ConnectReduxAssets(
		vangogh_local_data.TitleProperty,
		vangogh_local_data.SlugProperty,
		vangogh_local_data.VideoIdProperty,
		vangogh_local_data.MissingVideoUrlProperty)

	if err != nil {
		return gva.EndWithError(err)
	}

	if missing {
		missingIds, err := itemizations.MissingLocalVideos(rxa)
		if err != nil {
			return gva.EndWithError(err)
		}
		for id := range missingIds {
			idSet[id] = true
		}
	}

	if len(idSet) == 0 {
		if missing {
			gva.EndWithResult("all videos are available locally")
		} else {
			gva.EndWithResult("no ids to get videos for")
		}
		return nil
	}

	gva.TotalInt(len(idSet))

	for id := range idSet {
		videoIds, ok := rxa.GetAllUnchangedValues(vangogh_local_data.VideoIdProperty, id)
		if !ok || len(videoIds) == 0 {
			gva.Increment()
			continue
		}

		title, _ := rxa.GetFirstVal(vangogh_local_data.TitleProperty, id)

		va := nod.Begin("%s %s", id, title)

		for _, videoId := range videoIds {
			if err := yet.DownloadVideos(http.DefaultClient, vgFnDelegate, ffmpegCmd, videoId); err != nil {
				if vErr := rxa.AddVal(vangogh_local_data.MissingVideoUrlProperty, videoId, err.Error()); vErr != nil {
					return err
				}
				va.Error(err)
				continue
			}
		}

		va.End()
		gva.Increment()
	}

	return nil
}

func vgFnDelegate(videoId string, videoPage *yt_urls.InitialPlayerResponse) string {
	return filepath.Join(
		vangogh_local_data.AbsDirByVideoId(videoId),
		videoId+yt_urls.DefaultExt)
}
