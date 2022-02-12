package itemizations

import (
	"github.com/arelate/vangogh_data"
	"github.com/boggydigital/gost"
	"github.com/boggydigital/kvas"
	"github.com/boggydigital/nod"
)

type imageReduxGetter struct {
	imageType   vangogh_data.ImageType
	reduxAssets kvas.ReduxAssets
}

func NewImageReduxGetter(
	it vangogh_data.ImageType,
	rxa kvas.ReduxAssets) *imageReduxGetter {
	return &imageReduxGetter{
		imageType:   it,
		reduxAssets: rxa,
	}
}

func (ieg *imageReduxGetter) GetImageIds(id string) ([]string, bool) {
	return ieg.reduxAssets.GetAllValues(vangogh_data.PropertyFromImageType(ieg.imageType), id)
}

func MissingLocalImages(
	it vangogh_data.ImageType,
	rxa kvas.ReduxAssets,
	localImageIds gost.StrSet) (vangogh_data.IdSet, error) {

	all := rxa.Keys(vangogh_data.PropertyFromImageType(it))

	if localImageIds == nil {
		var err error
		if localImageIds, err = vangogh_data.LocalImageIds(); err != nil {
			return vangogh_data.NewIdSet(), err
		}
	}

	ieg := NewImageReduxGetter(it, rxa)

	mlia := nod.NewProgress(" itemizing local images (%s)...", it)
	defer mlia.EndWithResult("done")

	return missingLocalFiles(all, localImageIds, ieg.GetImageIds, nil, mlia)
}
