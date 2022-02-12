package fetchers

import (
	"encoding/json"
	"fmt"
	"github.com/arelate/gog_atu"
	"github.com/arelate/vangogh_data"
	"github.com/boggydigital/dolo"
	"github.com/boggydigital/kvas"
	"github.com/boggydigital/nod"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

//TODO: move this to kvas
type kvasIndexSetter struct {
	keyValues kvas.KeyValues
	ids       []string
}

func (kis *kvasIndexSetter) Len() int {
	return len(kis.ids)
}

func (kis *kvasIndexSetter) Exists(int) bool {
	//kvas performs hash computation to track modified files,
	//so we want all set attempts to go through (we need to
	//read src to compute that hash)
	return false
}

func (kis *kvasIndexSetter) Set(index int, src io.ReadCloser, results chan *dolo.IndexResult, errors chan *dolo.IndexError) {

	defer src.Close()

	if index < 0 || index >= len(kis.ids) {
		errors <- dolo.NewIndexError(index, fmt.Errorf("id index out of bounds"))
	}

	if err := kis.keyValues.Set(kis.ids[index], src); err != nil {
		errors <- dolo.NewIndexError(index, err)
	}

	results <- dolo.NewIndexResult(index, true)
}

func (kis *kvasIndexSetter) Get(key string) (io.ReadCloser, error) {
	return kis.keyValues.Get(key)
}

func NewKvasIndexSetter(pt vangogh_data.ProductType, mt gog_atu.Media, ids []string) (*kvasIndexSetter, error) {

	localDir, err := vangogh_data.AbsLocalProductTypeDir(pt, mt)
	if err != nil {
		return nil, err
	}

	valueSet, err := kvas.ConnectLocal(localDir, kvas.JsonExt)
	if err != nil {
		return nil, err
	}

	return &kvasIndexSetter{
		keyValues: valueSet,
		ids:       ids,
	}, nil
}

//Pages fetches all paged product type pages concurrently (using dolo.GetSet).
//To do that it downloads the first page, decodes that to get TotalPages,
//then constructs a slice of URLs and page ids to download all the remaining
//pages from 2nd to TotalPages using kvas index setter.
func Pages(pt vangogh_data.ProductType, mt gog_atu.Media, httpClient *http.Client, tpw nod.TotalProgressWriter) error {

	gfp := nod.Begin(" getting the first %s (%s)...", pt, mt)
	defer gfp.End()

	remoteUrl, err := vangogh_data.RemoteProductsUrl(pt)
	if err != nil {
		return err
	}

	//we need to handle the first page of the paged product type get-data request
	//separately from the rest, because:
	//1) initially we don't know how many pages paged data source would have (at least one is guaranteed)
	//2) first page contains the amount of pages data source has, so upon saving that we'll read it back
	//3) after figuring out how many pages data source has, we can construct URLs and page ids to get
	//the remaining set using dolo.GetSet and kvasIndexSetter.
	//Here is how we put that plan in motion:

	firstPage := "1"

	//construct a list of a single first page URL and page id "1"
	urls, ids := make([]*url.URL, 1), make([]string, 1)
	urls[0], ids[0] = remoteUrl(firstPage, mt), firstPage

	//initiate kvasIndexSetter using single page id "1"
	kis, err := NewKvasIndexSetter(pt, mt, ids)
	if err != nil {
		return err
	}

	dc := dolo.NewClient(httpClient, dolo.Defaults())

	//get the first page payload and set it in kvas
	if err := dc.GetSet(urls, kis, tpw); err != nil {
		return err
	}

	//get downloaded first page from kvas...
	fpReadCloser, err := kis.Get(firstPage)
	defer fpReadCloser.Close()
	if err != nil {
		return err
	}

	//...and decode it using minimal data structure to get total pages amount
	var page gog_atu.TotalPages
	if err = json.NewDecoder(fpReadCloser).Decode(&page); err != nil {
		return tpw.EndWithError(err)
	}

	gfp.EndWithResult("done")

	//now that we know how many pages we have in total - reinitialize URLs and ids to
	//that number minus the first page we already got
	urls, ids = make([]*url.URL, page.TotalPages-1), make([]string, page.TotalPages-1)

	for i := 2; i <= page.TotalPages; i++ {
		page := strconv.Itoa(i)
		//i-2 = first page (index: 0) will be 2
		urls[i-2] = remoteUrl(page, mt)
		ids[i-2] = page
	}

	kis, err = NewKvasIndexSetter(pt, mt, ids)
	if err != nil {
		return err
	}

	if err := dc.GetSet(urls, kis, tpw); err != nil {
		return err
	}

	tpw.EndWithResult("done")

	return nil
}
