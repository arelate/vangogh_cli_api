package cli

import (
	"github.com/arelate/gog_media"
	"github.com/arelate/vangogh_api/cli/checks"
	"github.com/arelate/vangogh_api/cli/url_helpers"
	"github.com/arelate/vangogh_downloads"
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
		localOnlyData:        url_helpers.Flag(u, VetOptionLocalOnlyData),
		recycleBin:           url_helpers.Flag(u, VetOptionRecycleBin),
		invalidData:          url_helpers.Flag(u, VetOptionInvalidData),
		unresolvedManualUrls: url_helpers.Flag(u, VetOptionUnresolvedManualUrls),
	}

	if url_helpers.Flag(u, "all") {
		vo.localOnlyData = !url_helpers.Flag(u, NegOpt(VetOptionLocalOnlyData))
		vo.recycleBin = !url_helpers.Flag(u, NegOpt(VetOptionRecycleBin))
		vo.invalidData = !url_helpers.Flag(u, NegOpt(VetOptionInvalidData))
		vo.unresolvedManualUrls = !url_helpers.Flag(u, NegOpt(VetOptionUnresolvedManualUrls))
	}

	return vo
}

func VetHandler(u *url.URL) error {
	mt := gog_media.Parse(url_helpers.Value(u, "media"))

	operatingSystems := url_helpers.OperatingSystems(u)
	downloadTypes := url_helpers.DownloadTypes(u)
	langCodes := url_helpers.Values(u, "language-code")

	vetOpts := initVetOptions(u)

	fix := url_helpers.Flag(u, "fix")

	return Vet(mt, vetOpts, operatingSystems, downloadTypes, langCodes, fix)
}

func Vet(
	mt gog_media.Media,
	vetOpts *vetOptions,
	operatingSystems []vangogh_downloads.OperatingSystem,
	downloadTypes []vangogh_downloads.DownloadType,
	langCodes []string,
	fix bool) error {

	sda := nod.Begin("vetting local data...")
	defer sda.End()

	if vetOpts.localOnlyData {
		if err := checks.LocalOnlySplitProducts(mt, fix); err != nil {
			return sda.EndWithError(err)
		}
	}

	if vetOpts.recycleBin {
		if err := checks.FilesInRecycleBin(fix); err != nil {
			return sda.EndWithError(err)
		}
	}

	if vetOpts.invalidData {
		if err := checks.InvalidLocalProductData(mt, fix); err != nil {
			return sda.EndWithError(err)
		}
	}

	if vetOpts.unresolvedManualUrls {
		if err := checks.UnresolvedManualUrls(mt, operatingSystems, downloadTypes, langCodes, fix); err != nil {
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
