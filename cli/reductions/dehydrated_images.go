package reductions

import (
	"github.com/arelate/vangogh_local_data"
	"github.com/boggydigital/issa"
	"github.com/boggydigital/nod"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"
)

func DehydratedImages() error {

	dia := nod.NewProgress(" %s...", vangogh_local_data.DehydratedImageProperty)
	defer dia.End()

	rxa, err := vangogh_local_data.ConnectReduxAssets(
		vangogh_local_data.ImageProperty,
		vangogh_local_data.DehydratedImageProperty)
	if err != nil {
		return dia.EndWithError(err)
	}

	dehydratedImages := make(map[string][]string)

	for _, id := range rxa.Keys(vangogh_local_data.ImageProperty) {
		if !rxa.HasKey(vangogh_local_data.DehydratedImageProperty, id) {
			dehydratedImages[id] = nil
		}
	}

	dia.TotalInt(len(dehydratedImages))

	for id := range dehydratedImages {

		if imageId, ok := rxa.GetFirstVal(vangogh_local_data.ImageProperty, id); ok {
			absLocalImagePath := vangogh_local_data.AbsLocalImagePath(imageId)
			if fi, err := os.Open(absLocalImagePath); err == nil {
				if jpegImage, _, err := image.Decode(fi); err == nil {
					gifImage := issa.GIFImage(jpegImage, issa.StdPalette(), issa.DefaultSampling*2)

					if dhi, err := issa.Dehydrate(gifImage); err == nil {
						dehydratedImages[id] = []string{dhi}
					} else {
						dia.Error(err)
					}
				} else {
					dia.Error(err)
				}
			} else {
				dia.Error(err)
			}

		}
		dia.Increment()
	}

	if err := rxa.BatchReplaceValues(vangogh_local_data.DehydratedImageProperty, dehydratedImages); err != nil {
		dia.EndWithError(err)
		return err
	}

	dia.EndWithResult("done")

	return nil
}
