package cli

import (
	"fmt"
	"github.com/arelate/vangogh_api/cli/itemize"
	"github.com/arelate/vangogh_api/cli/url_helpers"
	"github.com/arelate/vangogh_extracts"
	"github.com/arelate/vangogh_images"
	"github.com/arelate/vangogh_properties"
	"github.com/arelate/vangogh_urls"
	"github.com/boggydigital/dolo"
	"github.com/boggydigital/gost"
	"github.com/boggydigital/nod"
	"net/url"
	"path/filepath"
)

func GetImagesHandler(u *url.URL) error {
	idSet, err := url_helpers.IdSet(u)
	if err != nil {
		return err
	}

	imageTypes := url_helpers.Values(u, "image-type")
	its := make([]vangogh_images.ImageType, 0, len(imageTypes))
	for _, imageType := range imageTypes {
		its = append(its, vangogh_images.Parse(imageType))
	}
	missing := url_helpers.Flag(u, "missing")

	return GetImages(idSet, its, missing)
}

//GetImages fetches remote images for a given type (box-art, screenshots, background, etc.).
//If requested it can check locally present files and download all missing (used in data files,
//but not present locally) images for a given type.
func GetImages(
	idSet gost.StrSet,
	its []vangogh_images.ImageType,
	missing bool) error {

	gia := nod.NewProgress("getting images...")
	defer gia.End()

	for _, it := range its {
		if !vangogh_images.Valid(it) {
			return gia.EndWithError(fmt.Errorf("invalid image type %s", it))
		}
	}

	propSet := gost.NewStrSetWith(
		vangogh_properties.TitleProperty)

	for _, it := range its {
		propSet.Add(vangogh_properties.FromImageType(it))
	}

	exl, err := vangogh_extracts.NewList(propSet.All()...)
	if err != nil {
		return gia.EndWithError(err)
	}

	//for every product we'll collect image types missing for id and download only those
	idMissingTypes := map[string][]vangogh_images.ImageType{}

	if missing {
		localImageSet, err := vangogh_urls.LocalImageIds()
		if err != nil {
			return gia.EndWithError(err)
		}
		//to track image types missing for each product we do the following:
		//1. for every image type requested to be downloaded - get product ids that don't have this image type locally
		//2. for every product id we get this way - add this image type to idMissingTypes[id]
		for _, it := range its {
			//1
			missingImageIds, err := itemize.MissingLocalImages(it, exl, localImageSet)
			if err != nil {
				return gia.EndWithError(err)
			}

			//2
			for _, id := range missingImageIds.All() {
				if idMissingTypes[id] == nil {
					idMissingTypes[id] = make([]vangogh_images.ImageType, 0)
				}
				idMissingTypes[id] = append(idMissingTypes[id], it)
			}
		}
	} else {
		for _, id := range idSet.All() {
			idMissingTypes[id] = its
		}
	}

	if len(idMissingTypes) == 0 {
		if missing {
			gia.EndWithResult("all images are available locally")
		} else {
			gia.EndWithResult("need at least one product id to get images")
		}
		return nil
	}

	gia.TotalInt(len(idMissingTypes))

	for id, missingIts := range idMissingTypes {

		//for every product collect all image URLs and all corresponding local filenames
		//to pass to dolo.GetSet, that'll concurrently download all required product images

		title, ok := exl.Get(vangogh_properties.TitleProperty, id)
		if !ok {
			title = id
		}

		mita := nod.NewProgress("%s %s", id, title)

		//sensible assumption - we'll have at least as many URLs and filenames
		//as types of images we're missing
		urls := make([]*url.URL, 0, len(missingIts))
		filenames := make([]string, 0, len(missingIts))

		for _, it := range missingIts {

			images, ok := exl.GetAll(vangogh_properties.FromImageType(it), id)
			if !ok || len(images) == 0 {
				continue
			}

			srcUrls, err := vangogh_urls.PropImageUrls(images, it)
			if err != nil {
				return mita.EndWithError(err)
			}

			urls = append(urls, srcUrls...)

			for _, srcUrl := range srcUrls {
				dstDir := vangogh_urls.ImageDir(srcUrl.Path)
				filenames = append(filenames, filepath.Join(dstDir, srcUrl.Path))
			}
		}

		imagesIndexSetter := dolo.NewFileIndexSetter(filenames)

		//using http.DefaultClient as no image types require authentication
		//(this might change in the future)
		if err := dolo.DefaultClient.GetSet(urls, imagesIndexSetter, mita); err != nil {
			return mita.EndWithError(err)
		}

		mita.EndWithResult("done")
		gia.Increment()
	}

	gia.EndWithResult("done")

	return nil
}
