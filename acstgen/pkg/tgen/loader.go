package tgen

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/Azure/acs-labs/acstgen/pkg/api"
	"github.com/Azure/acs-labs/acstgen/pkg/api/v20160330"
	"github.com/Azure/acs-labs/acstgen/pkg/api/vlabs"
)

// LoadContainerServiceFromFile loads an ACS Cluster API Model from a JSON file
func LoadContainerServiceFromFile(jsonFile string) (*api.ContainerService, error) {
	contents, e := ioutil.ReadFile(jsonFile)
	if e != nil {
		return nil, fmt.Errorf("error reading file %s: %s", jsonFile, e.Error())
	}
	return LoadContainerService(contents)
}

// LoadContainerService loads an ACS Cluster API Model from a JSON file
func LoadContainerService(contents []byte) (*api.ContainerService, error) {
	m := &api.TypeMeta{}
	if err := json.Unmarshal(contents, &m); err != nil {
		return nil, err
	}

	switch m.APIVersion {
	case v20160330.APIVersion:
		containerService := &v20160330.ContainerService{}
		if e := json.Unmarshal(contents, &containerService); e != nil {
			return nil, e
		}

		if e := containerService.Properties.Validate(); e != nil {
			return nil, e
		}
		return api.ConvertV20160330ContainerService(containerService), nil

	case vlabs.APIVersion:
		containerService := &vlabs.ContainerService{}
		if e := json.Unmarshal(contents, &containerService); e != nil {
			return nil, e
		}

		if e := containerService.Properties.Validate(); e != nil {
			return nil, e
		}
		return api.ConvertVLabsContainerService(containerService), nil

	default:
		return nil, fmt.Errorf("unrecognized APIVersion '%s'", m.APIVersion)
	}
}
