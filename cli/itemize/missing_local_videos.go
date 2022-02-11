package itemize

import (
	"github.com/arelate/vangogh_data"
	"github.com/boggydigital/gost"
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
	return vpg.reduxAssets.GetAllUnchangedValues(vangogh_data.VideoIdProperty, id)
}

func (vpg *videoPropertiesGetter) IsMissingVideo(videoId string) bool {
	return vpg.reduxAssets.HasKey(vangogh_data.MissingVideoUrlProperty, videoId)
}

func MissingLocalVideos(rxa kvas.ReduxAssets) (gost.StrSet, error) {
	all := rxa.Keys(vangogh_data.VideoIdProperty)

	localVideoSet, err := vangogh_data.LocalVideoIds()
	if err != nil {
		return nil, err
	}

	vpg := NewVideoPropertiesGetter(rxa)

	mlva := nod.NewProgress(" itemizing local videos...")
	defer mlva.EndWithResult("done")

	return missingLocalFiles(all, localVideoSet, vpg.GetVideoIds, vpg.IsMissingVideo, mlva)
}
