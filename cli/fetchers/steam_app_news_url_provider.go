package fetchers

import (
	"github.com/arelate/gog_integration"
	"github.com/arelate/steam_integration"
	"github.com/arelate/vangogh_local_data"
	"github.com/boggydigital/kvas"
	"net/url"
	"strconv"
)

type SteamAppNewsUrlProvider struct {
	rxa kvas.ReduxAssets
}

func (sanup *SteamAppNewsUrlProvider) DefaultSourceUrl(gogId string, _ gog_integration.Media) *url.URL {

	if appIdStr, ok := sanup.rxa.GetFirstVal(vangogh_local_data.SteamAppIdProperty, gogId); ok {
		if appId, err := strconv.Atoi(appIdStr); err == nil {
			return steam_integration.NewsForApp(uint32(appId))
		}
	}

	return nil
}
