package cli

import (
	"github.com/arelate/gog_atu"
	"github.com/arelate/vangogh_api/cli/vets"
	"github.com/arelate/vangogh_data"
	"github.com/boggydigital/nod"
	"net/url"
)

const (
	VetOptionLocalOnlyData        = "local-only-data"
	VetOptionRecycleBin           = "recycle-bin"
	VetOptionInvalidData          = "invalid-data"
	VetOptionUnresolvedManualUrls = "unresolved-manual-urls"
)

type vetOptions struct {
	localOnlyData        bool
	recycleBin           bool
	invalidData          bool
	unresolvedManualUrls bool
}

func initVetOptions(u *url.URL) *vetOptions {

	vo := &vetOptions{
		localOnlyData:        vangogh_data.FlagFromUrl(u, VetOptionLocalOnlyData),
		recycleBin:           vangogh_data.FlagFromUrl(u, VetOptionRecycleBin),
		invalidData:          vangogh_data.FlagFromUrl(u, VetOptionInvalidData),
		unresolvedManualUrls: vangogh_data.FlagFromUrl(u, VetOptionUnresolvedManualUrls),
	}

	if vangogh_data.FlagFromUrl(u, "all") {
		vo.localOnlyData = !vangogh_data.FlagFromUrl(u, NegOpt(VetOptionLocalOnlyData))
		vo.recycleBin = !vangogh_data.FlagFromUrl(u, NegOpt(VetOptionRecycleBin))
		vo.invalidData = !vangogh_data.FlagFromUrl(u, NegOpt(VetOptionInvalidData))
		vo.unresolvedManualUrls = !vangogh_data.FlagFromUrl(u, NegOpt(VetOptionUnresolvedManualUrls))
	}

	return vo
}

func VetHandler(u *url.URL) error {

	vetOpts := initVetOptions(u)

	return Vet(
		vangogh_data.MediaFromUrl(u),
		vetOpts,
		vangogh_data.OperatingSystemsFromUrl(u),
		vangogh_data.DownloadTypesFromUrl(u),
		vangogh_data.ValuesFromUrl(u, "language-code"),
		vangogh_data.FlagFromUrl(u, "fix"))
}

func Vet(
	mt gog_atu.Media,
	vetOpts *vetOptions,
	operatingSystems []vangogh_data.OperatingSystem,
	downloadTypes []vangogh_data.DownloadType,
	langCodes []string,
	fix bool) error {

	sda := nod.Begin("vetting local data...")
	defer sda.End()

	if vetOpts.localOnlyData {
		if err := vets.LocalOnlySplitProducts(mt, fix); err != nil {
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

	//products with values different from extracts
	//images that are not linked to a product
	//videos that are not linked to a product
	//logs older than 30 days
	//checksum file errors

	return nil
}
