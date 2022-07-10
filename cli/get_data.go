package cli

import (
	"fmt"
	"github.com/arelate/gog_integration"
	"github.com/arelate/vangogh_cli_api/cli/fetchers"
	"github.com/arelate/vangogh_cli_api/cli/itemizations"
	"github.com/arelate/vangogh_local_data"
	"github.com/boggydigital/coost"
	"github.com/boggydigital/nod"
	"net/url"
)

func GetDataHandler(u *url.URL) error {
	idSet, err := vangogh_local_data.IdSetFromUrl(u)
	if err != nil {
		return err
	}

	skipIds := vangogh_local_data.ValuesFromUrl(u, "skip-id")

	updated := vangogh_local_data.FlagFromUrl(u, "updated")
	since, err := vangogh_local_data.SinceFromUrl(u)
	if err != nil {
		return err
	}

	return GetData(
		idSet,
		skipIds,
		vangogh_local_data.ProductTypeFromUrl(u),
		vangogh_local_data.MediaFromUrl(u),
		since,
		vangogh_local_data.FlagFromUrl(u, "missing"),
		updated)
}

//GetData gets remote data from GOG.com and stores as local products (splitting as paged data if needed)
func GetData(
	idSet map[string]bool,
	skipIds []string,
	pt vangogh_local_data.ProductType,
	mt gog_integration.Media,
	since int64,
	missing bool,
	updated bool) error {

	gda := nod.NewProgress("getting %s (%s) data...", pt, mt)
	defer gda.End()

	if !vangogh_local_data.IsValidProductType(pt) {
		gda.EndWithResult("%s is not a valid product type", pt)
		return nil
	}

	if !vangogh_local_data.IsMediaSupported(pt, mt) {
		gda.EndWithResult("%s is not a supported media for %s", mt, pt)
		return nil
	}

	hc, err := coost.NewHttpClientFromFile(
		vangogh_local_data.AbsCookiePath(),
		gog_integration.GogHost)
	if err != nil {
		return gda.EndWithError(err)
	}

	if vangogh_local_data.IsProductRequiresAuth(pt) {
		li, err := gog_integration.LoggedIn(hc)
		if err != nil {
			return gda.EndWithError(err)
		}

		if !li {
			return gda.EndWithError(fmt.Errorf("user is not logged in"))
		}
	}

	if vangogh_local_data.IsPagedProduct(pt) {
		if err := fetchers.Pages(pt, mt, since, hc, gda); err != nil {
			return gda.EndWithError(err)
		}
		return split(pt, mt, since)
	}

	if vangogh_local_data.IsArrayProduct(pt) {
		ids := []string{pt.String()}
		if err := fetchers.Items(ids, pt, mt, hc); err != nil {
			return gda.EndWithError(err)
		}
		return split(pt, mt, since)
	}

	idSet, err = itemizations.All(idSet, missing, updated, since, pt, mt)
	if err != nil {
		return gda.EndWithError(err)
	}

	skipIdSet := make(map[string]bool, len(skipIds))
	for _, id := range skipIds {
		skipIdSet[id] = true
	}

	approvedIds := make([]string, 0, len(idSet))

	for id := range idSet {
		if !skipIdSet[id] {
			approvedIds = append(approvedIds, id)
		}
	}

	return fetchers.Items(approvedIds, pt, mt, hc)
}
