package cli

import (
	"fmt"
	"github.com/arelate/gog_integration"
	"github.com/arelate/vangogh_local_data"
	"github.com/boggydigital/coost"
	"github.com/boggydigital/nod"
	"net/url"
	"strings"
)

const (
	createOp = "create"
	deleteOp = "delete"
	addOp    = "add"
	removeOp = "remove"
)

func TagHandler(u *url.URL) error {
	idSet, err := vangogh_local_data.IdSetFromUrl(u)
	if err != nil {
		return err
	}

	return Tag(
		idSet,
		vangogh_local_data.ValueFromUrl(u, "operation"),
		vangogh_local_data.ValueFromUrl(u, "tag-name"))
}

func Tag(idSet map[string]bool, operation, tagName string) error {

	ta := nod.Begin("performing requested tag operation...")
	defer ta.End()

	//matching default GOG.com capitalization for tags
	tagName = strings.ToUpper(tagName)

	rxa, err := vangogh_local_data.ConnectReduxAssets(
		vangogh_local_data.TagNameProperty,
		vangogh_local_data.TagIdProperty,
		vangogh_local_data.TitleProperty,
	)
	if err != nil {
		return err
	}

	tagId := ""
	if operation != createOp {
		tagId, err = vangogh_local_data.TagIdByName(tagName, rxa)
		if err != nil {
			return err
		}
	}

	hc, err := coost.NewHttpClientFromFile(vangogh_local_data.AbsCookiePath(), gog_integration.GogHost)
	if err != nil {
		return ta.EndWithError(err)
	}

	toa := nod.NewProgress(" %s tag %s...", operation, tagName)
	defer toa.End()

	switch operation {
	case createOp:
		return vangogh_local_data.CreateTag(hc, tagName, rxa)
	case deleteOp:
		return vangogh_local_data.DeleteTag(hc, tagName, tagId, rxa)
	case addOp:
		return vangogh_local_data.AddTag(hc, idSet, tagId, rxa, toa)
	case removeOp:
		return vangogh_local_data.RemoveTag(hc, idSet, tagId, rxa, toa)
	default:
		return ta.EndWithError(fmt.Errorf("unknown tag operation %s", operation))
	}
}
