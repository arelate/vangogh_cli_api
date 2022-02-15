package cli

import (
	"encoding/json"
	"fmt"
	"github.com/arelate/gog_integration"
	"github.com/arelate/vangogh_cli_api/cli/dirs"
	"github.com/arelate/vangogh_local_data"
	"github.com/boggydigital/coost"
	"github.com/boggydigital/kvas"
	"github.com/boggydigital/nod"
	"net/url"
	"path/filepath"
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

func Tag(idSet vangogh_local_data.IdSet, operation, tagName string) error {

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
		tagId, err = tagIdByName(tagName, rxa)
		if err != nil {
			return err
		}
	}

	switch operation {
	case createOp:
		return createTag(tagName, rxa)
	case deleteOp:
		return deleteTag(tagName, tagId, rxa)
	case addOp:
		return addTag(idSet, tagName, tagId, rxa)
	case removeOp:
		return removeTag(idSet, tagName, tagId, rxa)
	default:
		return ta.EndWithError(fmt.Errorf("unknown tag operation %s", operation))
	}
}

func postResp(url *url.URL, respVal interface{}) error {
	hc, err := coost.NewHttpClientFromFile(
		filepath.Join(dirs.GetTempDir(), cookiesFilename), gog_integration.GogHost)
	if err != nil {
		return err
	}

	resp, err := hc.Post(url.String(), "", nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return fmt.Errorf("unexpected status: %s", resp.Status)
	}

	return json.NewDecoder(resp.Body).Decode(&respVal)
}

func tagIdByName(tagName string, rxa kvas.ReduxAssets) (string, error) {
	if err := rxa.IsSupported(vangogh_local_data.TagNameProperty); err != nil {
		return "", err
	}

	tagIds := rxa.Match(map[string][]string{vangogh_local_data.TagNameProperty: {tagName}}, true)
	if len(tagIds) == 0 {
		return "", fmt.Errorf("unknown tag-name %s", tagName)
	}
	if len(tagIds) > 1 {
		return "", fmt.Errorf("ambiguous tag-name %s, matching tag-ids: %v",
			tagName,
			tagIds)
	}
	tagId := ""
	for ti := range tagIds {
		tagId = ti
	}
	return tagId, nil
}

func createTag(tagName string, rxa kvas.ReduxAssets) error {

	cta := nod.Begin(" creating tag %s...", tagName)
	defer cta.End()

	if err := rxa.IsSupported(vangogh_local_data.TagNameProperty); err != nil {
		return cta.EndWithError(err)
	}

	createTagUrl := gog_integration.CreateTagUrl(tagName)
	var ctResp gog_integration.CreateTagResp
	if err := postResp(createTagUrl, &ctResp); err != nil {
		return cta.EndWithError(err)
	}
	if ctResp.Id == "" {
		return cta.EndWithError(fmt.Errorf("invalid create tag response"))
	}

	if err := rxa.AddVal(vangogh_local_data.TagNameProperty, ctResp.Id, tagName); err != nil {
		return cta.EndWithError(err)
	}

	cta.EndWithResult("done")

	return nil
}

func deleteTag(tagName, tagId string, rxa kvas.ReduxAssets) error {

	dta := nod.Begin(" deleting tag %s...", tagName)
	defer dta.End()

	if err := rxa.IsSupported(vangogh_local_data.TagNameProperty); err != nil {
		return dta.EndWithError(err)
	}

	deleteTagUrl := gog_integration.DeleteTagUrl(tagId)
	var dtResp gog_integration.DeleteTagResp
	if err := postResp(deleteTagUrl, &dtResp); err != nil {
		return dta.EndWithError(err)
	}
	if dtResp.Status != "deleted" {
		return dta.EndWithError(fmt.Errorf("invalid delete tag response"))
	}

	if err := rxa.CutVal(vangogh_local_data.TagNameProperty, tagId, tagName); err != nil {
		return dta.EndWithError(err)
	}

	dta.EndWithResult("done")

	return nil
}

func addTag(idSet vangogh_local_data.IdSet, tagName, tagId string, rxa kvas.ReduxAssets) error {

	ata := nod.NewProgress(" adding tag %s to item(s)...", tagName)
	defer ata.End()

	if err := rxa.IsSupported(vangogh_local_data.TagNameProperty, vangogh_local_data.TitleProperty); err != nil {
		return ata.EndWithError(err)
	}

	ata.TotalInt(idSet.Len())

	for _, id := range idSet.All() {
		addTagUrl := gog_integration.AddTagUrl(id, tagId)
		var artResp gog_integration.AddRemoveTagResp
		if err := postResp(addTagUrl, &artResp); err != nil {
			return ata.EndWithError(err)
		}
		if !artResp.Success {
			return ata.EndWithError(fmt.Errorf("failed to add tag %s", tagName))
		}

		if err := rxa.AddVal(vangogh_local_data.TagIdProperty, id, tagId); err != nil {
			return ata.EndWithError(err)
		}

		ata.Increment()
	}

	ata.EndWithResult("done")

	return nil
}

func removeTag(idSet vangogh_local_data.IdSet, tagName, tagId string, rxa kvas.ReduxAssets) error {

	rta := nod.NewProgress(" removing tag %s from item(s)...", tagName)
	defer rta.End()

	if err := rxa.IsSupported(vangogh_local_data.TagNameProperty, vangogh_local_data.TitleProperty); err != nil {
		return rta.EndWithError(err)
	}

	rta.TotalInt(idSet.Len())

	for _, id := range idSet.All() {
		removeTagUrl := gog_integration.RemoveTagUrl(id, tagId)
		var artResp gog_integration.AddRemoveTagResp
		if err := postResp(removeTagUrl, &artResp); err != nil {
			return rta.EndWithError(err)
		}
		if !artResp.Success {
			return rta.EndWithError(fmt.Errorf("failed to remove tag %s", tagName))
		}

		if err := rxa.CutVal(vangogh_local_data.TagIdProperty, id, tagId); err != nil {
			return rta.EndWithError(err)
		}

		rta.Increment()
	}

	rta.EndWithResult("done")

	return nil
}
