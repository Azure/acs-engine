package api

import (
	"github.com/Azure/acs-engine/pkg/api/common"
	"github.com/Azure/acs-engine/pkg/i18n"
	"github.com/leonelquinteros/gotext"

	"path"
	"testing"
)

func TestLoadContainerServiceFromFile(t *testing.T) {
	existingContainerService := &ContainerService{Name: "test",
		Properties: &Properties{OrchestratorProfile: &OrchestratorProfile{OrchestratorType: Kubernetes, OrchestratorVersion: "1.6.9"}}}

	locale := gotext.NewLocale(path.Join("..", "..", "translations"), "en_US")
	i18n.Initialize(locale)
	apiloader := &Apiloader{
		Translator: &i18n.Translator{
			Locale: locale,
		},
	}

	containerService, _, err := apiloader.LoadContainerServiceFromFile("../acsengine/testdata/v20170701/kubernetes.json", true, false, existingContainerService)
	if err != nil {
		t.Error(err.Error())
	}
	if containerService.Properties.OrchestratorProfile.OrchestratorVersion != "1.8.12" {
		t.Errorf("Failed to set orcherstator version when it is set in the json, expected 1.8.12 but got %s", containerService.Properties.OrchestratorProfile.OrchestratorVersion)
	}

	containerService, _, err = apiloader.LoadContainerServiceFromFile("../acsengine/testdata/v20170701/kubernetes-default-version.json", true, false, existingContainerService)
	if err != nil {
		t.Error(err.Error())
	}
	if containerService.Properties.OrchestratorProfile.OrchestratorVersion != "1.6.9" {
		t.Errorf("Failed to set orcherstator version when it is not set in the json, got %s", containerService.Properties.OrchestratorProfile.OrchestratorVersion)
	}

	containerService, _, err = apiloader.LoadContainerServiceFromFile("../acsengine/testdata/v20170131/kubernetes.json", true, false, existingContainerService)
	if err != nil {
		t.Error(err.Error())
	}
	if containerService.Properties.OrchestratorProfile.OrchestratorVersion != "1.6.9" {
		t.Errorf("Failed to set orcherstator version when it is not set in the json, got %s", containerService.Properties.OrchestratorProfile.OrchestratorVersion)
	}

	containerService, _, err = apiloader.LoadContainerServiceFromFile("../acsengine/testdata/v20160930/kubernetes.json", true, false, existingContainerService)
	if err != nil {
		t.Error(err.Error())
	}
	if containerService.Properties.OrchestratorProfile.OrchestratorVersion != "1.6.9" {
		t.Errorf("Failed to set orcherstator version when it is not set in the json, got %s", containerService.Properties.OrchestratorProfile.OrchestratorVersion)
	}

	containerService, _, err = apiloader.LoadContainerServiceFromFile("../acsengine/testdata/v20170701/kubernetes-default-version.json", true, false, nil)
	if err != nil {
		t.Error(err.Error())
	}
	if containerService.Properties.OrchestratorProfile.OrchestratorVersion != common.GetDefaultKubernetesVersion() {
		t.Errorf("Failed to set orcherstator version when it is not set in the json API v20170701, got %s but expected %s", containerService.Properties.OrchestratorProfile.OrchestratorVersion, common.GetDefaultKubernetesVersion())
	}

	containerService, _, err = apiloader.LoadContainerServiceFromFile("../acsengine/testdata/v20170701/kubernetes-win-default-version.json", true, false, nil)
	if err != nil {
		t.Error(err.Error())
	}
	if containerService.Properties.OrchestratorProfile.OrchestratorVersion != common.GetDefaultKubernetesVersionWindows() {
		t.Errorf("Failed to set orcherstator version to windows default when it is not set in the json API v20170701, got %s but expected %s", containerService.Properties.OrchestratorProfile.OrchestratorVersion, common.GetDefaultKubernetesVersionWindows())
	}

	containerService, _, err = apiloader.LoadContainerServiceFromFile("../acsengine/testdata/v20170131/kubernetes.json", true, false, nil)
	if err != nil {
		t.Error(err.Error())
	}
	if containerService.Properties.OrchestratorProfile.OrchestratorVersion != common.GetDefaultKubernetesVersion() {
		t.Errorf("Failed to set orcherstator version when it is not set in the json API v20170131, got %s but expected %s", containerService.Properties.OrchestratorProfile.OrchestratorVersion, common.GetDefaultKubernetesVersion())
	}

	containerService, _, err = apiloader.LoadContainerServiceFromFile("../acsengine/testdata/v20170131/kubernetes-win.json", true, false, nil)
	if err != nil {
		t.Error(err.Error())
	}
	if containerService.Properties.OrchestratorProfile.OrchestratorVersion != common.GetDefaultKubernetesVersionWindows() {
		t.Errorf("Failed to set orcherstator version to windows default when it is not set in the json API v20170131, got %s but expected %s", containerService.Properties.OrchestratorProfile.OrchestratorVersion, common.GetDefaultKubernetesVersionWindows())
	}
}
