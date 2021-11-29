package v1

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func GetExtracts(w http.ResponseWriter, r *http.Request) {

	// GET /v1/extracts?product-type&media&property&sort&desc&from&to

	if r.Method != http.MethodGet {
		http.Error(w, "unsupported method", 405)
		return
	}

	pt, mt, err := getProductTypeMedia(r.URL)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	properties := strings.Split(r.URL.Query().Get("property"), ",")
	for _, prop := range properties {
		if err := exl.AssertSupport(prop); err != nil {
			http.Error(w, fmt.Sprintf("unsupported property %s", prop), 400)
			return
		}
	}

	sort, desc := getSortDesc(r.URL)

	from, to, err := getFromTo(r.URL)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	if sids, err := getSortedIds(pt, mt, sort, desc); err != nil {
		http.Error(w, err.Error(), 500)
		return
	} else {
		values := make(map[string]map[string][]string, to-from+1)
		for i := from; i <= to; i++ {
			propValues := make(map[string][]string)
			for _, prop := range properties {
				propValues[prop], _ = exl.GetAll(prop, sids[i])
			}
			values[sids[i]] = propValues
		}
		if err := json.NewEncoder(w).Encode(values); err != nil {
			http.Error(w, err.Error(), 500)
		}
	}
}
