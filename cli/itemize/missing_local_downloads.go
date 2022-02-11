package itemize

import (
	"github.com/arelate/gog_atu"
	"github.com/arelate/vangogh_data"
	"github.com/boggydigital/gost"
	"github.com/boggydigital/kvas"
	"github.com/boggydigital/nod"
	"os"
	"path/filepath"
)

func MissingLocalDownloads(
	mt gog_atu.Media,
	rxa kvas.ReduxAssets,
	operatingSystems []vangogh_data.OperatingSystem,
	downloadTypes []vangogh_data.DownloadType,
	langCodes []string) (gost.StrSet, error) {
	//enumerating missing local downloads is a bit more complicated than images and videos
	//due to the fact that actual filenames are resolved when downloads are processed, so we can't compare
	//manualUrls and available files, we need to resolve manualUrls to actual local filenames first.
	//with this in mind we'll use different approach:
	//1. for all vangogh_data.Details ids:
	//2. check if there are unresolved manualUrls -> add to missingIds
	//3. check if slug dir is not present in downloads -> add to missingIds
	//4. check if any expected (resolved manualUrls) files are not present -> add to missingIds

	mlda := nod.NewProgress(" itemizing missing local downloads")
	defer mlda.End()

	if err := rxa.IsSupported(
		vangogh_data.LocalManualUrl,
		vangogh_data.DownloadStatusError); err != nil {
		return nil, mlda.EndWithError(err)
	}

	vrDetails, err := vangogh_data.NewReader(vangogh_data.Details, mt)
	if err != nil {
		return nil, mlda.EndWithError(err)
	}

	//1
	allIds := gost.NewStrSetWith(vrDetails.Keys()...)

	mlda.TotalInt(allIds.Len())

	mdd := &missingDownloadsDelegate{
		rxa: rxa}

	if err := vangogh_data.MapDownloads(
		allIds,
		mt,
		rxa,
		operatingSystems,
		downloadTypes,
		langCodes,
		mdd,
		mlda); err != nil {
		return mdd.missingIds, mlda.EndWithError(err)
	}

	return mdd.missingIds, nil
}

type missingDownloadsDelegate struct {
	rxa        kvas.ReduxAssets
	missingIds gost.StrSet
}

func (mdd *missingDownloadsDelegate) Process(id, slug string, list vangogh_data.DownloadsList) error {

	if mdd.missingIds == nil {
		mdd.missingIds = gost.NewStrSet()
	}

	//pDir = s/slug
	relDir, err := vangogh_data.RelProductDownloadsDir(slug)
	if err != nil {
		return err
	}

	expectedFiles := gost.NewStrSet()

	for _, dl := range list {

		//skip manualUrls that have produced error status codes, while they're technically missing
		//it's due to remote status for this URL, not a problem we can resolve locally
		status, ok := mdd.rxa.GetFirstVal(vangogh_data.DownloadStatusError, dl.ManualUrl)
		if ok && status == "404" {
			continue
		}

		localFilename, ok := mdd.rxa.GetFirstVal(vangogh_data.LocalManualUrl, dl.ManualUrl)
		// 2
		if !ok || localFilename == "" {
			mdd.missingIds.Add(id)
			break
		}
		//local filenames are saved as relative to root downloads folder (e.g. s/slug/local_filename)
		//so filepath.Rel would trim to local_filename (or dlc/local_filename, extra/local_filename)
		relFilename, err := filepath.Rel(relDir, localFilename)
		if err != nil {
			return err
		}

		expectedFiles.Add(relFilename)
	}

	if expectedFiles.Len() == 0 {
		return nil
	}

	// 3
	absDir, err := vangogh_data.AbsProductDownloadsDir(slug)
	if err != nil {
		return err
	}
	if _, err := os.Stat(absDir); os.IsNotExist(err) {
		mdd.missingIds.Add(id)
		return nil
	}

	presentFiles, err := vangogh_data.LocalSlugDownloads(slug)
	if err != nil {
		return nil
	}

	// 4
	missingFiles := expectedFiles.Except(presentFiles)
	if len(missingFiles) > 0 {
		mdd.missingIds.Add(id)
	}

	return nil
}
