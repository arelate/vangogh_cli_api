package cli

import (
	"github.com/arelate/vangogh_api/cli/itemize"
	"github.com/arelate/vangogh_api/cli/url_helpers"
	"github.com/arelate/vangogh_data"
	"github.com/boggydigital/dolo"
	"github.com/boggydigital/gost"
	"github.com/boggydigital/nod"
	"github.com/boggydigital/yt_urls"
	"net/http"
	"net/url"
)

const (
	missingStr = "missing"
)

func GetVideosHandler(u *url.URL) error {
	idSet, err := url_helpers.IdSet(u)
	if err != nil {
		return err
	}

	return GetVideos(
		idSet,
		vangogh_data.FlagFromUrl(u, "missing"))
}

func GetVideos(idSet gost.StrSet, missing bool) error {

	gva := nod.NewProgress("getting videos...")
	defer gva.End()

	rxa, err := vangogh_data.ConnectReduxAssets(
		vangogh_data.TitleProperty,
		vangogh_data.SlugProperty,
		vangogh_data.VideoIdProperty,
		vangogh_data.MissingVideoUrlProperty)

	if err != nil {
		return gva.EndWithError(err)
	}

	if missing {
		missingIds, err := itemize.MissingLocalVideos(rxa)
		if err != nil {
			return gva.EndWithError(err)
		}
		idSet.AddSet(missingIds)
	}

	if idSet.Len() == 0 {
		if missing {
			gva.EndWithResult("all videos are available locally")
		} else {
			gva.EndWithResult("no ids to get videos for")
		}
		return nil
	}

	gva.TotalInt(idSet.Len())

	for _, id := range idSet.All() {
		videoIds, ok := rxa.GetAllUnchangedValues(vangogh_data.VideoIdProperty, id)
		if !ok || len(videoIds) == 0 {
			gva.Increment()
			continue
		}

		title, _ := rxa.GetFirstVal(vangogh_data.TitleProperty, id)

		va := nod.Begin("%s %s", id, title)

		dl := dolo.DefaultClient

		for _, videoId := range videoIds {

			vp, err := yt_urls.GetVideoPage(http.DefaultClient, videoId)
			if err != nil {
				va.Error(err)
				if addErr := rxa.AddVal(vangogh_data.MissingVideoUrlProperty, videoId, err.Error()); addErr != nil {
					return addErr
				}
				continue
			}

			vfa := nod.NewProgress(" %s", vp.Title())

			vidUrls := vp.StreamingFormats()

			if len(vidUrls) == 0 {
				if err := rxa.AddVal(vangogh_data.MissingVideoUrlProperty, videoId, missingStr); err != nil {
					return vfa.EndWithError(err)
				}
			}

			for _, vidUrl := range vidUrls {

				if vidUrl.Url == "" {
					if err := rxa.AddVal(vangogh_data.MissingVideoUrlProperty, videoId, missingStr); err != nil {
						return vfa.EndWithError(err)
					}
					continue
				}

				dir := vangogh_data.AbsDirByVideoId(videoId)

				u, err := url.Parse(vidUrl.Url)
				if err != nil {
					return vfa.EndWithError(err)
				}

				//get-videos is not using dolo.GetSetMany unlike get-images, and is downloading
				//videos sequentially for two main reasons:
				//1) each video has a list of bitrate-sorted URLs, and we're attempting to download "the best" quality
				//moving to the next available on failure
				//2) currently dolo.GetSetMany doesn't support nod progress reporting on each individual concurrent
				//download (ok, well, StdOutPresenter doesn't, nod likely does) and for video files this would mean
				//long pauses as we download individual files
				if err = dl.Download(u, vfa, dir, videoId+yt_urls.DefaultExt); err != nil {
					vfa.Error(err)
					continue
				}

				//yt_urls.StreamingUrls returns bitrate-sorted video urls,
				//so we can stop, if we've successfully got the best available one
				break
			}
		}

		va.End()
		gva.Increment()
	}

	//gva.EndWithResult("done")

	return nil
}
