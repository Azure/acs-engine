package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/Azure/acs-engine/pkg/api/v20160330"
	"github.com/Azure/acs-engine/pkg/api/vlabs"
)

// LoadContainerServiceFromFile loads an ACS Cluster API Model from a JSON file
func LoadContainerServiceFromFile(jsonFile string) (*ContainerService, error) {
	contents, e := ioutil.ReadFile(jsonFile)
	if e != nil {
		return nil, fmt.Errorf("error reading file %s: %s", jsonFile, e.Error())
	}
	return LoadContainerService(contents)
}

// LoadContainerService loads an ACS Cluster API Model, validates it, and returns the unversioned representation
func LoadContainerService(contents []byte) (*ContainerService, error) {
	m := &TypeMeta{}
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
		return ConvertV20160330ContainerService(containerService), nil

	case vlabs.APIVersion:
		containerService := &vlabs.ContainerService{}
		if e := json.Unmarshal(contents, &containerService); e != nil {
			return nil, e
		}

		if e := containerService.Properties.Validate(); e != nil {
			return nil, e
		}
		return ConvertVLabsContainerService(containerService), nil

	default:
		return nil, fmt.Errorf("unrecognized APIVersion '%s'", m.APIVersion)
	}
}

// LoadContainerServiceFromAPI load an ACS Cluster API Model, validate it, and return the unversioned representation
func LoadContainerServiceFromAPI(api interface{}) (*ContainerService, error) {
	m, found := api.(*TypeMeta)
	if !found {
		return nil, fmt.Errorf("no APIVersion field")
	}
	switch m.APIVersion {
	case v20160330.APIVersion:
		if e := api.(*v20160330.ContainerService).Properties.Validate(); e != nil {
			return nil, e
		}
		return ConvertV20160330ContainerService(api.(*v20160330.ContainerService)), nil

	case vlabs.APIVersion:
		if e := api.(*vlabs.ContainerService).Properties.Validate(); e != nil {
			return nil, e
		}
		return ConvertVLabsContainerService(api.(*vlabs.ContainerService)), nil

	default:
		return nil, fmt.Errorf("unrecognized APIVersion '%s'", m.APIVersion)
	}
}
