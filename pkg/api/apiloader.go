package api

import (
	"encoding/json"
	"fmt"
	"github.com/Azure/acs-engine/pkg/api/v20160330"
	"github.com/Azure/acs-engine/pkg/api/v20160930"
	"github.com/Azure/acs-engine/pkg/api/v20170131"
	"github.com/Azure/acs-engine/pkg/api/v20170701"
	"github.com/Azure/acs-engine/pkg/api/vlabs"
	"io/ioutil"
)

// LoadContainerServiceFromFile loads an ACS Cluster API Model from a JSON file
func LoadContainerServiceFromFile(jsonFile string, validate bool) (*ContainerService, string, error) {
	contents, e := ioutil.ReadFile(jsonFile)
	if e != nil {
		return nil, "", fmt.Errorf("error reading file %s: %s", jsonFile, e.Error())
	}
	return DeserializeContainerService(contents, validate)
}

// DeserializeContainerService loads an ACS Cluster API Model, validates it, and returns the unversioned representation
func DeserializeContainerService(contents []byte, validate bool) (*ContainerService, string, error) {
	m := &TypeMeta{}
	if err := json.Unmarshal(contents, &m); err != nil {
		return nil, "", err
	}
	version := m.APIVersion
	service, err := LoadContainerService(contents, version, validate)

	return service, version, err
}

// LoadContainerService loads an ACS Cluster API Model, validates it, and returns the unversioned representation
func LoadContainerService(contents []byte, version string, validate bool) (*ContainerService, error) {
	switch version {
	case v20160930.APIVersion:
		containerService := &v20160930.ContainerService{}
		if e := json.Unmarshal(contents, &containerService); e != nil {
			return nil, e
		}
		setContainerServiceDefaultsv20160930(containerService)
		if e := containerService.Properties.Validate(); validate && e != nil {
			return nil, e
		}
		return ConvertV20160930ContainerService(containerService), nil

	case v20160330.APIVersion:
		containerService := &v20160330.ContainerService{}
		if e := json.Unmarshal(contents, &containerService); e != nil {
			return nil, e
		}
		setContainerServiceDefaultsv20160330(containerService)
		if e := containerService.Properties.Validate(); validate && e != nil {
			return nil, e
		}
		return ConvertV20160330ContainerService(containerService), nil

	case v20170131.APIVersion:
		containerService := &v20170131.ContainerService{}
		if e := json.Unmarshal(contents, &containerService); e != nil {
			return nil, e
		}
		setContainerServiceDefaultsv20170131(containerService)
		if e := containerService.Properties.Validate(); validate && e != nil {
			return nil, e
		}
		return ConvertV20170131ContainerService(containerService), nil

	case v20170701.APIVersion:
		containerService := &v20170701.ContainerService{}
		if e := json.Unmarshal(contents, &containerService); e != nil {
			return nil, e
		}
		if e := containerService.Properties.Validate(); validate && e != nil {
			return nil, e
		}
		return ConvertV20170701ContainerService(containerService), nil

	case vlabs.APIVersion:
		containerService := &vlabs.ContainerService{}
		if e := json.Unmarshal(contents, &containerService); e != nil {
			return nil, e
		}
		if e := containerService.Properties.Validate(); validate && e != nil {
			return nil, e
		}
		return ConvertVLabsContainerService(containerService), nil
	default:
		return nil, fmt.Errorf("unrecognized APIVersion '%s'", version)
	}
}

// SerializeContainerService takes an unversioned container service and returns the bytes
func SerializeContainerService(containerService *ContainerService, version string) ([]byte, error) {
	switch version {
	case v20160930.APIVersion:
		v20160930ContainerService := ConvertContainerServiceToV20160930(containerService)
		armContainerService := &V20160930ARMContainerService{}
		armContainerService.ContainerService = v20160930ContainerService
		armContainerService.APIVersion = version
		b, err := json.MarshalIndent(armContainerService, "", "  ")
		if err != nil {
			return nil, err
		}
		return b, nil

	case v20160330.APIVersion:
		v20160330ContainerService := ConvertContainerServiceToV20160330(containerService)
		armContainerService := &V20160330ARMContainerService{}
		armContainerService.ContainerService = v20160330ContainerService
		armContainerService.APIVersion = version
		b, err := json.MarshalIndent(armContainerService, "", "  ")
		if err != nil {
			return nil, err
		}
		return b, nil

	case v20170131.APIVersion:
		v20170131ContainerService := ConvertContainerServiceToV20170131(containerService)
		armContainerService := &V20170131ARMContainerService{}
		armContainerService.ContainerService = v20170131ContainerService
		armContainerService.APIVersion = version
		b, err := json.MarshalIndent(armContainerService, "", "  ")
		if err != nil {
			return nil, err
		}
		return b, nil

	case v20170701.APIVersion:
		v20170701ContainerService := ConvertContainerServiceToV20170701(containerService)
		armContainerService := &V20170701ARMContainerService{}
		armContainerService.ContainerService = v20170701ContainerService
		armContainerService.APIVersion = version
		b, err := json.MarshalIndent(armContainerService, "", "  ")
		if err != nil {
			return nil, err
		}
		return b, nil

	case vlabs.APIVersion:
		vlabsContainerService := ConvertContainerServiceToVLabs(containerService)
		armContainerService := &VlabsARMContainerService{}
		armContainerService.ContainerService = vlabsContainerService
		armContainerService.APIVersion = version
		b, err := json.MarshalIndent(armContainerService, "", "  ")
		if err != nil {
			return nil, err
		}
		return b, nil

	default:
		return nil, fmt.Errorf("invalid version %s for conversion back from unversioned object", version)
	}
}

// Sets default container service property values for any appropriate zero values
func setContainerServiceDefaultsv20160930(c *v20160930.ContainerService) {
	if c.Properties.OrchestratorProfile == nil {
		c.Properties.OrchestratorProfile = &v20160930.OrchestratorProfile{
			OrchestratorType: v20160930.DCOS,
		}
	}
}

// Sets default container service property values for any appropriate zero values
func setContainerServiceDefaultsv20160330(c *v20160330.ContainerService) {
	if c.Properties.OrchestratorProfile == nil {
		c.Properties.OrchestratorProfile = &v20160330.OrchestratorProfile{
			OrchestratorType: v20160330.DCOS,
		}
	}
}

// Sets default container service property values for any appropriate zero values
func setContainerServiceDefaultsv20170131(c *v20170131.ContainerService) {
	if c.Properties.OrchestratorProfile == nil {
		c.Properties.OrchestratorProfile = &v20170131.OrchestratorProfile{
			OrchestratorType: v20170131.DCOS,
		}
	}
}
