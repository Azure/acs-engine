package api

import (
	"encoding/json"
	"io/ioutil"

	"github.com/Azure/acs-engine/pkg/api/vlabs"
	"github.com/Azure/acs-engine/pkg/i18n"
)

// UpgradeApiloader represents the object that loads api model
type UpgradeApiloader struct {
	Translator *i18n.Translator
}

// LoadUpgradeContainerServiceFromFile loads an ACS Cluster API Model from a JSON file
func (ua *UpgradeApiloader) LoadUpgradeContainerServiceFromFile(jsonFile string) (*UpgradeContainerService, string, error) {
	contents, e := ioutil.ReadFile(jsonFile)
	if e != nil {
		return nil, "", ua.Translator.Errorf("error reading file %s: %s", jsonFile, e.Error())
	}
	return ua.DeserializeUpgradeContainerService(contents)
}

// DeserializeUpgradeContainerService loads an ACS Cluster API Model, validates it, and returns the unversioned representation
func (ua *UpgradeApiloader) DeserializeUpgradeContainerService(contents []byte) (*UpgradeContainerService, string, error) {
	m := &TypeMeta{}
	if err := json.Unmarshal(contents, &m); err != nil {
		return nil, "", err
	}
	version := m.APIVersion
	upgradecontainerservice, err := ua.LoadUpgradeContainerService(contents, version)

	return upgradecontainerservice, version, err
}

// LoadUpgradeContainerService loads an ACS Cluster API Model, validates it, and returns the unversioned representation
func (ua *UpgradeApiloader) LoadUpgradeContainerService(contents []byte, version string) (*UpgradeContainerService, error) {
	switch version {
	case vlabs.APIVersion:
		upgradecontainerService := &vlabs.UpgradeContainerService{}
		if e := json.Unmarshal(contents, &upgradecontainerService); e != nil {
			return nil, e
		}
		if e := upgradecontainerService.Validate(); e != nil {
			return nil, e
		}
		return ConvertVLabsUpgradeContainerService(upgradecontainerService), nil

	default:
		return nil, ua.Translator.Errorf("unrecognized APIVersion '%s'", version)
	}
}

// SerializeUpgradeContainerService takes an unversioned container service and returns the bytes
func (ua *UpgradeApiloader) SerializeUpgradeContainerService(upgradeContainerService *UpgradeContainerService, version string) ([]byte, error) {
	switch version {
	case vlabs.APIVersion:
		vlabsUpgradeContainerService := ConvertUpgradeContainerServiceToVLabs(upgradeContainerService)
		armUpgradeContainerService := &VlabsUpgradeContainerService{}
		armUpgradeContainerService.UpgradeContainerService = vlabsUpgradeContainerService
		armUpgradeContainerService.APIVersion = version
		b, err := json.MarshalIndent(armUpgradeContainerService, "", "  ")
		if err != nil {
			return nil, err
		}
		return b, nil

	default:
		return nil, ua.Translator.Errorf("invalid version %s for conversion back from unversioned object", version)
	}
}
