package itemize

import (
	"github.com/arelate/vangogh_data"
	"github.com/boggydigital/gost"
	"github.com/boggydigital/kvas"
	"github.com/boggydigital/nod"
)

type imageExtractsGetter struct {
	imageType   vangogh_data.ImageType
	reduxAssets kvas.ReduxAssets
}

func NewImageExtractsGetter(
	it vangogh_data.ImageType,
	rxa kvas.ReduxAssets) *imageExtractsGetter {
	return &imageExtractsGetter{
		imageType:   it,
		reduxAssets: rxa,
	}
}

func (ieg *imageExtractsGetter) GetImageIds(id string) ([]string, bool) {
	return ieg.reduxAssets.GetAllValues(vangogh_data.PropertyFromImageType(ieg.imageType), id)
}

func MissingLocalImages(
	it vangogh_data.ImageType,
	rxa kvas.ReduxAssets,
	localImageIds gost.StrSet) (gost.StrSet, error) {

	all := rxa.Keys(vangogh_data.PropertyFromImageType(it))

	if localImageIds == nil {
		var err error
		if localImageIds, err = vangogh_data.LocalImageIds(); err != nil {
			return nil, err
		}
	}

	ieg := NewImageExtractsGetter(it, rxa)

	mlia := nod.NewProgress(" itemizing local images (%s)...", it)
	defer mlia.EndWithResult("done")

	return missingLocalFiles(all, localImageIds, ieg.GetImageIds, nil, mlia)
}
