package acsengine

import (
	"path"
	"testing"

	"github.com/Azure/acs-engine/pkg/api"
	"github.com/Azure/acs-engine/pkg/i18n"
	"github.com/leonelquinteros/gotext"
)

func TestAssignParameters(t *testing.T) {
	// Initialize locale for translation
	locale := gotext.NewLocale(path.Join("..", "..", "translations"), "en_US")
	i18n.Initialize(locale)

	apiloader := &api.Apiloader{
		Translator: &i18n.Translator{
			Locale: locale,
		},
	}
	// iterate the test data directory
	apiModelTestFiles := &[]APIModelTestFile{}
	if e := IterateTestFilesDirectory(TestDataDir, apiModelTestFiles); e != nil {
		t.Error(e.Error())
		return
	}

	for _, tuple := range *apiModelTestFiles {
		containerService, _, err := apiloader.LoadContainerServiceFromFile(tuple.APIModelFilename, true, false, nil)
		if err != nil {
			t.Errorf("Loading file %s got error: %s", tuple.APIModelFilename, err.Error())
			continue
		}

		containerService.Location = "eastus"
		containerService.SetPropertiesDefaults(false, false)
		parametersMap, err := getParameters(containerService, DefaultGeneratorCode, "testversion")
		if err != nil {
			t.Errorf("should not get error when populating parameters")
		}
		for k, v := range parametersMap {
			switch val := v.(paramsMap)["value"].(type) {
			case *bool:
				t.Errorf("got a pointer to bool in paramsMap value, this is dangerous!: %s: %v", k, val)
			}
		}
	}
}
