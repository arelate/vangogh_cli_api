package cli

import (
	"github.com/arelate/gog_integration"
	"github.com/arelate/vangogh_local_data"
	"github.com/boggydigital/issa"
	"github.com/boggydigital/nod"
	"image"
	_ "image/jpeg"
	"net/url"
	"os"
	"strings"
)

func DehydrateImagesHandler(u *url.URL) error {
	return DehydrateImages(
		vangogh_local_data.MediaFromUrl(u))
}

func DehydrateImages(mt gog_integration.Media) error {

	dia := nod.NewProgress("dehydrating images...")
	defer dia.End()

	rxa, err := vangogh_local_data.ConnectReduxAssets(vangogh_local_data.ImageProperty)
	if err != nil {
		return dia.EndWithError(err)
	}

	vrDehydratedImages, err := vangogh_local_data.NewReader(vangogh_local_data.DehydratedImages, mt)
	if err != nil {
		return dia.EndWithError(err)
	}

	missingImages := make(map[string]interface{})

	for _, id := range rxa.Keys(vangogh_local_data.ImageProperty) {
		if imageId, ok := rxa.GetFirstVal(vangogh_local_data.ImageProperty, id); ok {
			if !vrDehydratedImages.Has(imageId) {
				missingImages[imageId] = nil
			}
		}
	}

	dia.TotalInt(len(missingImages))

	for imageId := range missingImages {

		absLocalImagePath := vangogh_local_data.AbsLocalImagePath(imageId)
		if _, err := os.Stat(absLocalImagePath); err != nil {
			dia.Error(err)
			dia.Increment()
			continue
		}

		if fi, err := os.Open(absLocalImagePath); err != nil {
			dia.Error(err)
			dia.Increment()
			continue
		} else {
			jpegImage, _, err := image.Decode(fi)
			if err != nil {
				dia.Error(err)
				dia.Increment()
				continue
			}

			gifImage := issa.GIFImage(jpegImage, issa.StdPalette(), issa.DefaultSampling)

			dhi, err := issa.Dehydrate(gifImage)
			if err != nil {
				dia.Error(err)
				dia.Increment()
				continue
			}

			if err := vrDehydratedImages.Set(imageId, strings.NewReader(dhi)); err != nil {
				dia.Error(err)
				dia.Increment()
				continue
			}

			dia.Increment()
		}

	}

	dia.EndWithResult("done")

	return nil
}
