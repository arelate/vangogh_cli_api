package cli

import (
	"fmt"
	"github.com/arelate/gog_atu"
	"github.com/arelate/vangogh_api/cli/fetchers"
	"github.com/arelate/vangogh_api/cli/itemizations"
	"github.com/arelate/vangogh_data"
	"github.com/boggydigital/coost"
	"github.com/boggydigital/gost"
	"github.com/boggydigital/nod"
	"net/url"
	"path/filepath"
	"time"
)

func GetDataHandler(u *url.URL) error {
	idSet, err := vangogh_data.IdSetFromUrl(u)
	if err != nil {
		return err
	}

	skipIds := vangogh_data.ValuesFromUrl(u, "skip-id")

	updated := vangogh_data.FlagFromUrl(u, "updated")
	since := time.Now().Unix()
	if updated {
		since = time.Now().Add(-time.Hour * 24).Unix()
	}

	return GetData(
		idSet,
		skipIds,
		vangogh_data.ProductTypeFromUrl(u),
		vangogh_data.MediaFromUrl(u),
		since,
		vangogh_data.ValueFromUrl(u, "temp-directory"),
		vangogh_data.FlagFromUrl(u, "missing"),
		updated)
}

//GetData gets remote data from GOG.com and stores as local products (splitting as paged data if needed)
func GetData(
	idSet vangogh_data.IdSet,
	skipIds []string,
	pt vangogh_data.ProductType,
	mt gog_atu.Media,
	since int64,
	tempDir string,
	missing bool,
	updated bool) error {

	gda := nod.NewProgress("getting %s (%s) data...", pt, mt)
	defer gda.End()

	if !vangogh_data.IsValidProductType(pt) {
		gda.EndWithResult("%s is not a valid product type", pt)
		return nil
	}

	if !vangogh_data.IsMediaSupported(pt, mt) {
		gda.EndWithResult("%s is not a supported media for %s", mt, pt)
		return nil
	}

	hc, err := coost.NewHttpClientFromFile(
		filepath.Join(tempDir, cookiesFilename), gog_atu.GogHost)
	if err != nil {
		return gda.EndWithError(err)
	}

	if vangogh_data.IsProductRequiresAuth(pt) {
		li, err := gog_atu.LoggedIn(hc)
		if err != nil {
			return gda.EndWithError(err)
		}

		if !li {
			return gda.EndWithError(fmt.Errorf("user is not logged in"))
		}
	}

	if vangogh_data.IsPagedProduct(pt) {
		if err := fetchers.Pages(pt, mt, hc, gda); err != nil {
			return gda.EndWithError(err)
		}
		return split(pt, mt, since)
	}

	if vangogh_data.IsArrayProduct(pt) {
		// using "licences" as id, since that's how we store that data in kvas
		ids := []string{vangogh_data.Licences.String()}
		if err := fetchers.Items(ids, pt, mt, hc); err != nil {
			return gda.EndWithError(err)
		}
		return split(pt, mt, since)
	}

	idSet, err = itemizations.All(idSet, missing, updated, since, pt, mt)
	if err != nil {
		return gda.EndWithError(err)
	}

	approvedIds := idSet.Except(gost.NewStrSetWith(skipIds...))

	return fetchers.Items(approvedIds, pt, mt, hc)
}
