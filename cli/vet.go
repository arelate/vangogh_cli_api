package cli

import (
	"github.com/arelate/gog_integration"
	"github.com/arelate/vangogh_cli_api/cli/vets"
	"github.com/arelate/vangogh_local_data"
	"github.com/boggydigital/nod"
	"net/url"
)

const (
	VetOptionLocalOnlyData             = "local-only-data"
	VetOptionLocalOnlyImages           = "local-only-images"
	VetOptionRecycleBin                = "recycle-bin"
	VetOptionInvalidData               = "invalid-data"
	VetOptionUnresolvedManualUrls      = "unresolved-manual-urls"
	VetOptionInvalidResolvedManualUrls = "invalid-resolved-manual-urls"
)

type vetOptions struct {
	localOnlyData               bool
	localOnlyImages             bool
	recycleBin                  bool
	invalidData                 bool
	unresolvedManualUrls        bool
	invalidUnresolvedManualUrls bool
}

func initVetOptions(u *url.URL) *vetOptions {

	vo := &vetOptions{
		localOnlyData:               vangogh_local_data.FlagFromUrl(u, VetOptionLocalOnlyData),
		localOnlyImages:             vangogh_local_data.FlagFromUrl(u, VetOptionLocalOnlyImages),
		recycleBin:                  vangogh_local_data.FlagFromUrl(u, VetOptionRecycleBin),
		invalidData:                 vangogh_local_data.FlagFromUrl(u, VetOptionInvalidData),
		unresolvedManualUrls:        vangogh_local_data.FlagFromUrl(u, VetOptionUnresolvedManualUrls),
		invalidUnresolvedManualUrls: vangogh_local_data.FlagFromUrl(u, VetOptionInvalidResolvedManualUrls),
	}

	if vangogh_local_data.FlagFromUrl(u, "all") {
		vo.localOnlyData = !vangogh_local_data.FlagFromUrl(u, NegOpt(VetOptionLocalOnlyData))
		vo.localOnlyImages = !vangogh_local_data.FlagFromUrl(u, NegOpt(VetOptionLocalOnlyImages))
		vo.recycleBin = !vangogh_local_data.FlagFromUrl(u, NegOpt(VetOptionRecycleBin))
		vo.invalidData = !vangogh_local_data.FlagFromUrl(u, NegOpt(VetOptionInvalidData))
		vo.unresolvedManualUrls = !vangogh_local_data.FlagFromUrl(u, NegOpt(VetOptionUnresolvedManualUrls))
		vo.invalidUnresolvedManualUrls = !vangogh_local_data.FlagFromUrl(u, NegOpt(VetOptionInvalidResolvedManualUrls))
	}

	return vo
}

func VetHandler(u *url.URL) error {

	vetOpts := initVetOptions(u)

	return Vet(
		vangogh_local_data.MediaFromUrl(u),
		vetOpts,
		vangogh_local_data.OperatingSystemsFromUrl(u),
		vangogh_local_data.DownloadTypesFromUrl(u),
		vangogh_local_data.ValuesFromUrl(u, "language-code"),
		vangogh_local_data.FlagFromUrl(u, "fix"))
}

func Vet(
	mt gog_integration.Media,
	vetOpts *vetOptions,
	operatingSystems []vangogh_local_data.OperatingSystem,
	downloadTypes []vangogh_local_data.DownloadType,
	langCodes []string,
	fix bool) error {

	sda := nod.Begin("vetting local data...")
	defer sda.End()

	if vetOpts.localOnlyData {
		if err := vets.LocalOnlySplitProducts(mt, fix); err != nil {
			return sda.EndWithError(err)
		}
	}

	if vetOpts.localOnlyImages {
		if err := vets.LocalOnlyImages(fix); err != nil {
			return sda.EndWithError(err)
		}
	}

	if vetOpts.recycleBin {
		if err := vets.FilesInRecycleBin(fix); err != nil {
			return sda.EndWithError(err)
		}
	}

	if vetOpts.invalidData {
		if err := vets.InvalidLocalProductData(mt, fix); err != nil {
			return sda.EndWithError(err)
		}
	}

	if vetOpts.unresolvedManualUrls {
		if err := vets.UnresolvedManualUrls(mt, operatingSystems, downloadTypes, langCodes, fix); err != nil {
			return sda.EndWithError(err)
		}
	}

	if vetOpts.invalidUnresolvedManualUrls {
		if err := vets.InvalidResolvedManualUrls(fix); err != nil {
			return sda.EndWithError(err)
		}
	}

	//products with values different from redux
	//videos that are not linked to a product
	//logs older than 30 days
	//checksum file errors

	return nil
}
