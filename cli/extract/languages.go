package extract

import (
	"github.com/arelate/gog_atu"
	"github.com/arelate/vangogh_data"
	"github.com/boggydigital/gost"
	"github.com/boggydigital/kvas"
	"github.com/boggydigital/nod"
)

func GetLanguageCodes(rxa kvas.ReduxAssets) (gost.StrSet, error) {

	lca := nod.Begin(" %s...", vangogh_data.LanguageCodeProperty)
	defer lca.EndWithResult("done")

	langCodeSet := gost.NewStrSet()

	if err := rxa.IsSupported(vangogh_data.LanguageCodeProperty); err != nil {
		return langCodeSet, lca.EndWithError(err)
	}

	//digest distinct languages codes
	for _, id := range rxa.Keys(vangogh_data.LanguageCodeProperty) {
		idCodes, ok := rxa.GetAllUnchangedValues(vangogh_data.LanguageCodeProperty, id)
		if !ok {
			continue
		}
		for _, code := range idCodes {
			langCodeSet.Add(code)
		}
	}

	return langCodeSet, nil
}

func getMissingLanguageNames(
	langCodeSet gost.StrSet,
	rxa kvas.ReduxAssets,
	property string) (gost.StrSet, error) {
	missingLangs := gost.NewStrSetWith(langCodeSet.All()...)

	// TODO: write a comment explaining all or nothing approach
	//map all language codes to names and hide existing
	for _, lc := range missingLangs.All() {
		if _, ok := rxa.GetFirstVal(property, lc); ok {
			missingLangs.Hide(lc)
		}
	}

	return missingLangs, nil
}

func updateLanguageNames(languages map[string]string, missingNames gost.StrSet, names map[string][]string) {
	for langCode, langName := range languages {
		if missingNames.Has(langCode) {
			names[langCode] = []string{langName}
			missingNames.Hide(langCode)
		}
	}
}

func LanguageNames(langCodeSet gost.StrSet) error {
	property := vangogh_data.LanguageNameProperty

	lna := nod.Begin(" %s...", property)
	defer lna.EndWithResult("done")

	langNamesEx, err := vangogh_data.ConnectReduxAssets(property)
	if err != nil {
		return lna.EndWithError(err)
	}

	missingLangs, err := getMissingLanguageNames(langCodeSet, langNamesEx, property)
	if err != nil {
		return lna.EndWithError(err)
	}

	if missingLangs.Len() == 0 {
		return nil
	}

	missingLangs = gost.NewStrSetWith(langCodeSet.All()...)
	names := make(map[string][]string, 0)

	//iterate through api-products-v1 until we fill all native names
	vrApiProductsV2, err := vangogh_data.NewReader(vangogh_data.ApiProductsV2, gog_atu.Game)
	if err != nil {
		return lna.EndWithError(err)
	}

	for _, id := range vrApiProductsV2.Keys() {
		apv2, err := vrApiProductsV2.ApiProductV2(id)
		if err != nil {
			return lna.EndWithError(err)
		}

		updateLanguageNames(apv2.GetLanguages(), missingLangs, names)

		if missingLangs.Len() == 0 {
			break
		}
	}

	if err := langNamesEx.BatchReplaceValues(property, names); err != nil {
		return lna.EndWithError(err)
	}

	return nil
}

func NativeLanguageNames(langCodeSet gost.StrSet) error {
	property := vangogh_data.NativeLanguageNameProperty

	nlna := nod.Begin(" %s...", property)
	defer nlna.End()

	langNamesEx, err := vangogh_data.ConnectReduxAssets(property)
	if err != nil {
		return nlna.EndWithError(err)
	}

	missingNativeLangs, err := getMissingLanguageNames(langCodeSet, langNamesEx, property)
	if err != nil {
		return nlna.EndWithError(err)
	}

	if missingNativeLangs.Len() == 0 {
		nlna.EndWithResult("done")
		return nil
	}

	vrApiProductsV1, err := vangogh_data.NewReader(vangogh_data.ApiProductsV1, gog_atu.Game)
	if err != nil {
		return nlna.EndWithError(err)
	}

	missingNativeLangs = gost.NewStrSetWith(langCodeSet.All()...)
	nativeNames := make(map[string][]string, 0)

	for _, id := range vrApiProductsV1.Keys() {
		apv1, err := vrApiProductsV1.ApiProductV1(id)
		if err != nil {
			return nlna.EndWithError(err)
		}

		updateLanguageNames(apv1.GetNativeLanguages(), missingNativeLangs, nativeNames)

		if missingNativeLangs.Len() == 0 {
			break
		}
	}

	if err := langNamesEx.BatchReplaceValues(property, nativeNames); err != nil {
		return nlna.EndWithError(err)
	}

	nlna.EndWithResult("done")

	return nil
}
