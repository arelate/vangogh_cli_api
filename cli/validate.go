package cli

import (
	"crypto/md5"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/arelate/gog_atu"
	"github.com/arelate/vangogh_api/cli/url_helpers"
	"github.com/arelate/vangogh_api/cli/validation"
	"github.com/arelate/vangogh_downloads"
	"github.com/arelate/vangogh_products"
	"github.com/arelate/vangogh_properties"
	"github.com/arelate/vangogh_urls"
	"github.com/arelate/vangogh_values"
	"github.com/boggydigital/dolo"
	"github.com/boggydigital/gost"
	"github.com/boggydigital/kvas"
	"github.com/boggydigital/nod"
	"net/url"
	"os"
)

var (
	ErrUnresolvedManualUrl    = errors.New("unresolved manual-url")
	ErrMissingDownload        = errors.New("not downloaded")
	ErrMissingChecksum        = errors.New("missing checksum")
	ErrValidationNotSupported = errors.New("validation not supported")
	ErrValidationFailed       = errors.New("failed validation")
)

func ValidateHandler(u *url.URL) error {
	idSet, err := url_helpers.IdSet(u)
	if err != nil {
		return err
	}

	return Validate(
		idSet,
		vangogh_urls.UrlMedia(u),
		vangogh_downloads.UrlOperatingSystems(u),
		vangogh_downloads.UrlDownloadTypes(u),
		vangogh_urls.UrlValues(u, "language-code"),
		vangogh_urls.UrlFlag(u, "all"))
}

func Validate(
	idSet gost.StrSet,
	mt gog_atu.Media,
	operatingSystems []vangogh_downloads.OperatingSystem,
	downloadTypes []vangogh_downloads.DownloadType,
	langCodes []string,
	all bool) error {

	va := nod.NewProgress("validating...")
	defer va.End()

	rxa, err := vangogh_properties.ConnectReduxAssets(
		vangogh_properties.SlugProperty,
		vangogh_properties.NativeLanguageNameProperty,
		vangogh_properties.LocalManualUrl)
	if err != nil {
		return err
	}

	if all {
		vrDetails, err := vangogh_values.NewReader(vangogh_products.Details, mt)
		if err != nil {
			return err
		}
		idSet.Add(vrDetails.Keys()...)
	}

	vd := &validateDelegate{rxa: rxa}

	if err := vangogh_downloads.Map(
		idSet,
		mt,
		rxa,
		operatingSystems,
		downloadTypes,
		langCodes,
		vd,
		va); err != nil {
		return err
	}

	summary := map[string][]string{}
	tp := fmt.Sprintf("%d product(s) successfully validated", len(vd.validated))
	summary[tp] = []string{}
	maybeAddTopic(summary, "%d product(s) have unresolved manual-url (not downloaded)", vd.unresolvedManualUrl)
	maybeAddTopic(summary, "%d product(s) missing downloads", vd.missingDownloads)
	maybeAddTopic(summary, "%d product(s) without checksum", vd.missingChecksum)
	maybeAddTopic(summary, "%d product(s) failed validation", vd.failed)
	if len(vd.slugLastError) > 0 {
		tp = fmt.Sprintf("%d product(s) validation caused an error", len(vd.slugLastError))
		summary[tp] = make([]string, 0, len(vd.slugLastError))
		for slug, err := range vd.slugLastError {
			summary[tp] = append(summary[tp], fmt.Sprintf(" %s: %s", slug, err))
		}
	}

	va.EndWithSummary("", summary)

	return nil
}

func maybeAddTopic(summary map[string][]string, tmpl string, col map[string]bool) {
	if len(col) > 0 {
		tp := fmt.Sprintf(tmpl, len(col))
		summary[tp] = make([]string, 0, len(col))
		for it := range col {
			summary[tp] = append(summary[tp], it)
		}
	}
}

func validateManualUrl(
	slug string,
	dl *vangogh_downloads.Download,
	rxa kvas.ReduxAssets) error {

	if err := rxa.IsSupported(vangogh_properties.LocalManualUrl); err != nil {
		return err
	}

	mua := nod.NewProgress(" %s...", dl.String())
	defer mua.End()

	//local filenames are saved as relative to root downloads folder (e.g. s/slug/local_filename)
	localFile, ok := rxa.GetFirstVal(vangogh_properties.LocalManualUrl, dl.ManualUrl)
	if !ok {
		mua.EndWithResult(ErrUnresolvedManualUrl.Error())
		return ErrUnresolvedManualUrl
	}

	//absolute path (given a downloads/ root) for a s/slug/local_filename,
	//e.g. downloads/s/slug/local_filename
	absLocalFile := vangogh_urls.DownloadDirRelToAbs(localFile)
	if !vangogh_urls.CanValidate(absLocalFile) {
		mua.EndWithResult(ErrValidationNotSupported.Error())
		return ErrValidationNotSupported
	}

	if _, err := os.Stat(absLocalFile); os.IsNotExist(err) {
		mua.EndWithResult(ErrMissingDownload.Error())
		return ErrMissingDownload
	}

	absChecksumFile := vangogh_urls.LocalChecksumPath(absLocalFile)

	if _, err := os.Stat(absChecksumFile); os.IsNotExist(err) {
		mua.EndWithResult(ErrMissingChecksum.Error())
		return ErrMissingChecksum
	}

	chkFile, err := os.Open(absChecksumFile)
	if err != nil {
		return mua.EndWithError(err)
	}
	defer chkFile.Close()

	var chkData validation.File
	if err := xml.NewDecoder(chkFile).Decode(&chkData); err != nil {
		return mua.EndWithError(err)
	}

	sourceFile, err := os.Open(absLocalFile)
	if err != nil {
		return mua.EndWithError(err)
	}
	defer sourceFile.Close()

	h := md5.New()

	stat, err := sourceFile.Stat()
	if err != nil {
		return mua.EndWithError(err)
	}

	mua.Total(uint64(stat.Size()))
	if err != nil {
		return mua.EndWithError(err)
	}

	if err := dolo.CopyWithProgress(h, sourceFile, mua); err != nil {
		return mua.EndWithError(err)
	}

	sourceFileMD5 := fmt.Sprintf("%x", h.Sum(nil))

	if chkData.MD5 != sourceFileMD5 {
		mua.EndWithResult("error")
		return ErrValidationFailed
	} else {
		mua.EndWithResult("valid")
	}

	return nil
}

type validateDelegate struct {
	rxa                 kvas.ReduxAssets
	validated           map[string]bool
	unresolvedManualUrl map[string]bool
	missingDownloads    map[string]bool
	missingChecksum     map[string]bool
	failed              map[string]bool
	slugLastError       map[string]string
}

func (vd *validateDelegate) Process(_, slug string, list vangogh_downloads.DownloadsList) error {

	sva := nod.Begin(slug)
	defer sva.End()

	if vd.validated == nil {
		vd.validated = make(map[string]bool)
	}
	if vd.unresolvedManualUrl == nil {
		vd.unresolvedManualUrl = make(map[string]bool)
	}
	if vd.missingDownloads == nil {
		vd.missingDownloads = make(map[string]bool)
	}
	if vd.missingChecksum == nil {
		vd.missingChecksum = make(map[string]bool)
	}
	if vd.failed == nil {
		vd.failed = make(map[string]bool)
	}
	if vd.slugLastError == nil {
		vd.slugLastError = make(map[string]string)
	}

	hasValidationTargets := false

	for _, dl := range list {
		if err := validateManualUrl(slug, &dl, vd.rxa); errors.Is(err, ErrValidationNotSupported) {
			continue
		} else if errors.Is(err, ErrMissingChecksum) {
			vd.missingChecksum[slug] = true
		} else if errors.Is(err, ErrUnresolvedManualUrl) {
			vd.unresolvedManualUrl[slug] = true
		} else if errors.Is(err, ErrMissingDownload) {
			vd.missingDownloads[slug] = true
		} else if errors.Is(err, ErrValidationFailed) {
			vd.failed[slug] = true
		} else if err != nil {
			vd.slugLastError[slug] = err.Error()
			continue
		}
		// don't attempt to assess success for files that don't support validation
		hasValidationTargets = true
	}

	if hasValidationTargets &&
		!vd.missingChecksum[slug] &&
		!vd.unresolvedManualUrl[slug] &&
		!vd.missingDownloads[slug] &&
		!vd.failed[slug] &&
		vd.slugLastError[slug] == "" {
		vd.validated[slug] = true
	}

	return nil
}
