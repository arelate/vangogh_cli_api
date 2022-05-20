package reductions

import (
	"github.com/arelate/gog_integration"
	"github.com/arelate/steam_integration"
	"github.com/arelate/vangogh_local_data"
	"github.com/boggydigital/nod"
	"strconv"
)

func SteamAppId(since int64) error {

	saia := nod.Begin(" %s...", vangogh_local_data.SteamAppId)
	defer saia.End()

	rxa, err := vangogh_local_data.ConnectReduxAssets(
		vangogh_local_data.TitleProperty,
		vangogh_local_data.SteamAppId)
	if err != nil {
		return saia.EndWithError(err)
	}

	vrSteamAppList, err := vangogh_local_data.NewReader(vangogh_local_data.SteamAppList, gog_integration.Game)
	if err != nil {
		return saia.EndWithError(err)
	}

	if vrSteamAppList.IsModifiedAfter(vangogh_local_data.SteamAppList.String(), since) {
		saia.EndWithResult("unchanged")
		return nil
	}

	galr, err := vrSteamAppList.SteamGetAppListResponse()
	if err != nil {
		return saia.EndWithError(err)
	}

	appMap := GetAppListResponseToMap(galr)
	gogSteamAppId := make(map[string][]string)

	for _, id := range rxa.Keys(vangogh_local_data.TitleProperty) {
		title, ok := rxa.GetFirstVal(vangogh_local_data.TitleProperty, id)
		if !ok {
			continue
		}

		if appId, ok := appMap[title]; ok {
			gogSteamAppId[id] = []string{strconv.Itoa(int(appId))}
		}

	}

	if err := rxa.BatchReplaceValues(vangogh_local_data.SteamAppId, gogSteamAppId); err != nil {
		return saia.EndWithError(err)
	}

	saia.EndWithResult("done")

	return nil
}

func GetAppListResponseToMap(galr *steam_integration.GetAppListResponse) map[string]uint32 {
	appsMap := make(map[string]uint32, len(galr.AppList.Apps))

	for _, app := range galr.AppList.Apps {
		if app.Name == "" {
			continue
		}
		appsMap[app.Name] = app.AppId
	}

	return appsMap
}
