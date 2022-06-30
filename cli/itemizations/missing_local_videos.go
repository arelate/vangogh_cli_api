package itemizations

import (
	"github.com/arelate/vangogh_local_data"
	"github.com/boggydigital/kvas"
	"github.com/boggydigital/nod"
)

type videoPropertiesGetter struct {
	reduxAssets kvas.ReduxAssets
}

func NewVideoPropertiesGetter(rxa kvas.ReduxAssets) *videoPropertiesGetter {
	return &videoPropertiesGetter{
		reduxAssets: rxa,
	}
}

func (vpg *videoPropertiesGetter) GetVideoIds(id string) ([]string, bool) {
	return vpg.reduxAssets.GetAllUnchangedValues(vangogh_local_data.VideoIdProperty, id)
}

func (vpg *videoPropertiesGetter) IsMissingVideo(videoId string) bool {
	return vpg.reduxAssets.HasKey(vangogh_local_data.MissingVideoUrlProperty, videoId)
}

func missingLocalVideoRelatedFiles(
	rxa kvas.ReduxAssets,
	localVideoIdsDelegate func() (map[string]bool, error),
	media string) (map[string]bool, error) {
	all := rxa.Keys(vangogh_local_data.VideoIdProperty)

	localSet, err := localVideoIdsDelegate()
	if err != nil {
		return map[string]bool{}, err
	}

	vpg := NewVideoPropertiesGetter(rxa)

	mlma := nod.NewProgress(" itemizing local %s...", media)
	defer mlma.EndWithResult("done")

	return missingLocalFiles(all, localSet, vpg.GetVideoIds, vpg.IsMissingVideo, mlma)
}

func MissingLocalVideos(rxa kvas.ReduxAssets) (map[string]bool, error) {
	return missingLocalVideoRelatedFiles(rxa, vangogh_local_data.LocalVideoIds, "videos")
}

func MissingLocalThumbnails(rxa kvas.ReduxAssets) (map[string]bool, error) {
	return missingLocalVideoRelatedFiles(rxa, vangogh_local_data.LocalVideoThumbnailIds, "thumbnails")
}
