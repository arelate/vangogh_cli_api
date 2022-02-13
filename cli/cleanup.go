package cli

import (
	"fmt"
	"github.com/arelate/gog_integration"
	"github.com/arelate/vangogh_local_data"
	"github.com/boggydigital/gost"
	"github.com/boggydigital/kvas"
	"github.com/boggydigital/nod"
	"math"
	"net/url"
	"os"
	"path/filepath"
)

const spaceSavingsSummary = "est. disk space savings:"

func CleanupHandler(u *url.URL) error {
	idSet, err := vangogh_local_data.IdSetFromUrl(u)
	if err != nil {
		return err
	}

	return Cleanup(
		idSet,
		vangogh_local_data.MediaFromUrl(u),
		vangogh_local_data.OperatingSystemsFromUrl(u),
		vangogh_local_data.DownloadTypesFromUrl(u),
		vangogh_local_data.ValuesFromUrl(u, "language-code"),
		vangogh_local_data.FlagFromUrl(u, "all"),
		vangogh_local_data.FlagFromUrl(u, "test"))
}

func Cleanup(
	idSet vangogh_local_data.IdSet,
	mt gog_integration.Media,
	operatingSystems []vangogh_local_data.OperatingSystem,
	downloadTypes []vangogh_local_data.DownloadType,
	langCodes []string,
	all, test bool) error {

	rxa, err := vangogh_local_data.ConnectReduxAssets(
		vangogh_local_data.SlugProperty,
		vangogh_local_data.NativeLanguageNameProperty,
		vangogh_local_data.LocalManualUrlProperty)
	if err != nil {
		return err
	}

	ca := nod.NewProgress("cleaning up...")
	defer ca.End()

	if all {
		vrDetails, err := vangogh_local_data.NewReader(vangogh_local_data.Details, mt)
		if err != nil {
			return err
		}
		idSet.Add(vrDetails.Keys()...)
	}

	cd := &cleanupDelegate{
		rxa:  rxa,
		all:  all,
		test: test,
	}

	if err := vangogh_local_data.MapDownloads(
		idSet,
		mt,
		rxa,
		operatingSystems,
		downloadTypes,
		langCodes,
		cd,
		ca); err != nil {
		return err
	}

	if cd.totalBytes > 0 {
		summary := make(map[string][]string)
		summary[spaceSavingsSummary] = []string{
			fmt.Sprintf("%.2fGB", float64(cd.totalBytes)/math.Pow(1000, 3)),
		}
		ca.EndWithSummary("cleanup summary:", summary)
	} else {
		ca.EndWithResult("done")
	}

	return nil
}

func moveToRecycleBin(fp string) error {
	rbFilepath := filepath.Join(vangogh_local_data.AbsRecycleBinDir(), fp)
	rbDir, _ := filepath.Split(rbFilepath)
	if _, err := os.Stat(rbDir); os.IsNotExist(err) {
		if err := os.MkdirAll(rbDir, 0755); err != nil {
			return err
		}
	}
	return os.Rename(fp, rbFilepath)
}

type cleanupDelegate struct {
	rxa        kvas.ReduxAssets
	all        bool
	test       bool
	totalBytes int64
}

func (cd *cleanupDelegate) Process(_ string, slug string, list vangogh_local_data.DownloadsList) error {

	csa := nod.QueueBegin(slug)
	defer csa.End()

	if err := cd.rxa.IsSupported(vangogh_local_data.LocalManualUrlProperty); err != nil {
		return csa.EndWithError(err)
	}

	//cleanup process:
	//1. enumerate all expected files for a downloadList
	//2. enumerate all files present for a slug (files present in a `downloads/slug` folder)
	//3. delete (present files).Except(expected files) and corresponding xml files

	expectedSet := gost.NewStrSet()

	//pDir = s/slug
	pDir, err := vangogh_local_data.RelProductDownloadsDir(slug)
	if err != nil {
		return csa.EndWithError(err)
	}

	for _, dl := range list {
		if localFilename, ok := cd.rxa.GetFirstVal(vangogh_local_data.LocalManualUrlProperty, dl.ManualUrl); ok {
			//local filenames are saved as relative to root downloads folder (e.g. s/slug/local_filename)
			//so filepath.Rel would trim to local_filename (or dlc/local_filename, extra/local_filename)
			relFilename, err := filepath.Rel(pDir, localFilename)
			if err != nil {
				return csa.EndWithError(err)
			}
			expectedSet.Add(relFilename)
		}
	}

	//LocalSlugDownloads returns list of files relative to s/slug product directory
	presentSet, err := vangogh_local_data.LocalSlugDownloads(slug)
	if err != nil {
		return csa.EndWithError(err)
	}

	//TODO: rework this as part of gost deprecation
	presentGost := gost.NewStrSet()
	for file := range presentSet {
		presentGost.Add(file)
	}

	unexpectedFiles := presentGost.Except(expectedSet)

	if len(unexpectedFiles) == 0 {
		if !cd.all {
			csa.EndWithResult("already clean")
			csa.Flush()
		}
		return nil
	}

	//given some unexpected files - flush message queue to output slug and put the files
	//output next in context of a slug we've queued earlier
	csa.Flush()

	for _, unexpectedFile := range unexpectedFiles {
		//restore absolute from local_filename to s/slug/local_filename
		downloadFilename := vangogh_local_data.AbsDownloadDirFromRel(filepath.Join(pDir, unexpectedFile))
		if stat, err := os.Stat(downloadFilename); err == nil {
			cd.totalBytes += stat.Size()
		} else if os.IsNotExist(err) {
			continue
		} else {
			return csa.EndWithError(err)
		}

		prefix := "DELETE"
		if cd.test {
			prefix = "TEST"
		}

		dft := nod.Begin(" %s %s", prefix, downloadFilename)
		if !cd.test {
			if err := moveToRecycleBin(downloadFilename); err != nil {
				return dft.EndWithError(err)
			}
		}
		dft.End()

		checksumFile := vangogh_local_data.AbsLocalChecksumPath(downloadFilename)
		if stat, err := os.Stat(checksumFile); err == nil {
			cd.totalBytes += stat.Size()
		} else if os.IsNotExist(err) {
			continue
		} else {
			return csa.EndWithError(err)
		}

		cft := nod.Begin(" %s %s", prefix, checksumFile)
		if !cd.test {
			if err := moveToRecycleBin(checksumFile); err != nil {
				return cft.EndWithError(err)
			}
		}
		cft.End()
	}

	return nil
}
