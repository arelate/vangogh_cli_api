package cli

import (
	"github.com/arelate/vangogh_local_data"
	"github.com/boggydigital/issa"
	"github.com/boggydigital/nod"
	"image"
	"net/url"
	"os"
)

const vangoghSamplingRate = 24

func DehydrateImagesHandler(u *url.URL) error {
	return DehydrateImages()
}

func DehydrateImages() error {
	dia := nod.NewProgress("dehydrating images...")
	defer dia.End()

	rxa, err := vangogh_local_data.ConnectReduxAssets(
		vangogh_local_data.ImageProperty,
		vangogh_local_data.DehydratedImageProperty)
	if err != nil {
		return dia.EndWithError(err)
	}

	dehydratedImages := make(map[string][]string)

	for _, id := range rxa.Keys(vangogh_local_data.ImageProperty) {
		if dip, ok := rxa.GetFirstVal(vangogh_local_data.DehydratedImageProperty, id); !ok || dip == "" {
			dehydratedImages[id] = nil
		}
	}

	dia.TotalInt(len(dehydratedImages))

	for id := range dehydratedImages {

		if imageId, ok := rxa.GetFirstVal(vangogh_local_data.ImageProperty, id); ok {
			absLocalImagePath := vangogh_local_data.AbsLocalImagePath(imageId)
			if fi, err := os.Open(absLocalImagePath); err == nil {
				if jpegImage, _, err := image.Decode(fi); err == nil {
					gifImage := issa.GIFImage(jpegImage, issa.StdPalette(), vangoghSamplingRate)

					if dhi, err := issa.Dehydrate(gifImage); err == nil {
						dehydratedImages[id] = []string{dhi}
					} else {
						dia.Error(err)
					}
				} else {
					dia.Error(err)
				}
			} else if !os.IsNotExist(err) {
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
