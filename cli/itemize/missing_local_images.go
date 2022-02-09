package itemize

import (
	"github.com/arelate/vangogh_images"
	"github.com/arelate/vangogh_properties"
	"github.com/arelate/vangogh_urls"
	"github.com/boggydigital/gost"
	"github.com/boggydigital/kvas"
	"github.com/boggydigital/nod"
)

type imageExtractsGetter struct {
	imageType   vangogh_images.ImageType
	reduxAssets kvas.ReduxAssets
}

func NewImageExtractsGetter(
	it vangogh_images.ImageType,
	rxa kvas.ReduxAssets) *imageExtractsGetter {
	return &imageExtractsGetter{
		imageType:   it,
		reduxAssets: rxa,
	}
}

func (ieg *imageExtractsGetter) GetImageIds(id string) ([]string, bool) {
	return ieg.reduxAssets.GetAllValues(vangogh_properties.FromImageType(ieg.imageType), id)
}

func MissingLocalImages(
	it vangogh_images.ImageType,
	rxa kvas.ReduxAssets,
	localImageIds gost.StrSet) (gost.StrSet, error) {

	all := rxa.Keys(vangogh_properties.FromImageType(it))

	if localImageIds == nil {
		var err error
		if localImageIds, err = vangogh_urls.LocalImageIds(); err != nil {
			return nil, err
		}
	}

	ieg := NewImageExtractsGetter(it, rxa)

	mlia := nod.NewProgress(" itemizing local images (%s)...", it)
	defer mlia.EndWithResult("done")

	return missingLocalFiles(all, localImageIds, ieg.GetImageIds, nil, mlia)
}
