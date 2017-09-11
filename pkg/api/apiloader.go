package api

import (
	"encoding/json"
	"io/ioutil"

	"github.com/Azure/acs-engine/pkg/api/common"
	"github.com/Azure/acs-engine/pkg/api/v20170930"

	"github.com/Azure/acs-engine/pkg/api/agentPoolOnlyApi/v20170831"
	apvlabs "github.com/Azure/acs-engine/pkg/api/agentPoolOnlyApi/vlabs"
	"github.com/Azure/acs-engine/pkg/api/v20160330"
	"github.com/Azure/acs-engine/pkg/api/v20160930"
	"github.com/Azure/acs-engine/pkg/api/v20170131"
	"github.com/Azure/acs-engine/pkg/api/v20170701"
	"github.com/Azure/acs-engine/pkg/api/vlabs"
	"github.com/Azure/acs-engine/pkg/i18n"
	log "github.com/sirupsen/logrus"
)

// Apiloader represents the object that loads api model
type Apiloader struct {
	Translator *i18n.Translator
}

// LoadContainerServiceFromFile loads an ACS Cluster API Model from a JSON file
func (a *Apiloader) LoadContainerServiceFromFile(jsonFile string, validate bool, existingContainerService *ContainerService) (*ContainerService, string, error) {
	contents, e := ioutil.ReadFile(jsonFile)
	if e != nil {
		return nil, "", a.Translator.Errorf("error reading file %s: %s", jsonFile, e.Error())
	}
	return a.DeserializeContainerService(contents, validate, existingContainerService)
}

// DeserializeContainerService loads an ACS Cluster API Model, validates it, and returns the unversioned representation
func (a *Apiloader) DeserializeContainerService(contents []byte, validate bool, existingContainerService *ContainerService) (*ContainerService, string, error) {
	m := &TypeMeta{}
	if err := json.Unmarshal(contents, &m); err != nil {
		return nil, "", err
	}

	version := m.APIVersion
	service, err := a.LoadContainerService(contents, version, validate, existingContainerService)
	if service == nil || err != nil {
		log.Infof("Error returned by LoadContainerService: %+v. Attempting to load container service using LoadContainerServiceForAgentPoolOnlyCluster", err)
		service, err = a.LoadContainerServiceForAgentPoolOnlyCluster(contents, version, validate)
	}

	return service, version, err
}

// LoadContainerService loads an ACS Cluster API Model, validates it, and returns the unversioned representation
func (a *Apiloader) LoadContainerService(
	contents []byte,
	version string,
	validate bool,
	existingContainerService *ContainerService) (*ContainerService, error) {
	switch version {
	case v20160930.APIVersion:
		containerService := &v20160930.ContainerService{}
		if e := json.Unmarshal(contents, &containerService); e != nil {
			return nil, e
		}
		if existingContainerService != nil {
			vecs := ConvertContainerServiceToV20160930(existingContainerService)
			if e := containerService.Merge(vecs); e != nil {
				return nil, e
			}
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
		if existingContainerService != nil {
			vecs := ConvertContainerServiceToV20160330(existingContainerService)
			if e := containerService.Merge(vecs); e != nil {
				return nil, e
			}
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
		if existingContainerService != nil {
			vecs := ConvertContainerServiceToV20170131(existingContainerService)
			if e := containerService.Merge(vecs); e != nil {
				return nil, e
			}
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
		if existingContainerService != nil {
			vecs := ConvertContainerServiceToV20170701(existingContainerService)
			if e := containerService.Merge(vecs); e != nil {
				return nil, e
			}
		}
		setContainerServiceDefaultsv20170701(containerService)
		if e := containerService.Properties.Validate(); validate && e != nil {
			return nil, e
		}
		return ConvertV20170701ContainerService(containerService), nil

	case vlabs.APIVersion:
		containerService := &vlabs.ContainerService{}
		if e := json.Unmarshal(contents, &containerService); e != nil {
			return nil, e
		}
		if existingContainerService != nil {
			vecs := ConvertContainerServiceToVLabs(existingContainerService)
			if e := containerService.Merge(vecs); e != nil {
				return nil, e
			}
		}
		setContainerServiceDefaultsvlabs(containerService)
		if e := containerService.Properties.Validate(); validate && e != nil {
			return nil, e
		}
		return ConvertVLabsContainerService(containerService), nil

	default:
		return nil, a.Translator.Errorf("unrecognized APIVersion '%s'", version)
	}
}

// LoadContainerServiceForAgentPoolOnlyCluster loads an ACS Cluster API Model, validates it, and returns the unversioned representation
func (a *Apiloader) LoadContainerServiceForAgentPoolOnlyCluster(contents []byte, version string, validate bool) (*ContainerService, error) {
	switch version {
	case v20170831.APIVersion:
		managedCluster := &v20170831.ManagedCluster{}
		if e := json.Unmarshal(contents, &managedCluster); e != nil {
			return nil, e
		}
		setManagedClusterDefaultsv20170831(managedCluster)
		if e := managedCluster.Properties.Validate(); validate && e != nil {
			return nil, e
		}
		return ConvertV20170831AgentPoolOnly(managedCluster), nil
	case apvlabs.APIVersion:
		managedCluster := &apvlabs.ManagedCluster{}
		if e := json.Unmarshal(contents, &managedCluster); e != nil {
			return nil, e
		}
		setManagedClusterDefaultsvlabs(managedCluster)
		if e := managedCluster.Properties.Validate(); validate && e != nil {
			return nil, e
		}
		return ConvertVLabsAgentPoolOnly(managedCluster), nil
	default:
		return nil, a.Translator.Errorf("unrecognized APIVersion in LoadContainerServiceForAgentPoolOnlyCluster '%s'", version)
	}
}

// UpdateContainerServiceForUpgrade pre-validates upgrade operation and updates container service
func (a *Apiloader) UpdateContainerServiceForUpgrade(
	contents []byte,
	version string,
	cs *ContainerService,
	allowCurrentVersionUpgrade bool) error {
	unverOrch := &OrchestratorProfile{}

	switch version {
	case v20170930.APIVersion:
		up := &v20170930.OrchestratorProfile{}
		if e := json.Unmarshal(contents, up); e != nil {
			return a.Translator.Errorf(e.Error())
		}
		if e := up.ValidateForUpgrade(); e != nil {
			return a.Translator.Errorf(e.Error())
		}
		convertV20170930OrchestratorProfile(up, unverOrch)

	case vlabs.APIVersion:
		up := &vlabs.OrchestratorProfile{}
		if e := json.Unmarshal(contents, up); e != nil {
			return a.Translator.Errorf(e.Error())
		}
		if e := up.ValidateForUpgrade(); e != nil {
			return a.Translator.Errorf(e.Error())
		}
		convertVLabsOrchestratorProfile(up, unverOrch)

	default:
		return a.Translator.Errorf("unrecognized APIVersion in UpdateContainerServiceForUpgrade '%s'", version)
	}

	// get available upgrades for container service
	orchestratorInfo, e := GetOrchestratorVersionProfile(cs.Properties.OrchestratorProfile)
	if e != nil {
		return e
	}

	// add current version if upgrade has failed
	if allowCurrentVersionUpgrade {
		release := cs.Properties.OrchestratorProfile.OrchestratorRelease
		orchestratorInfo.Upgrades = append(orchestratorInfo.Upgrades, &OrchestratorProfile{
			OrchestratorRelease: release,
			OrchestratorVersion: common.KubeReleaseToVersion[release]})
	}
	// validate desired upgrade version and set goal state
	for _, up := range orchestratorInfo.Upgrades {
		if up.OrchestratorRelease == unverOrch.OrchestratorRelease {
			cs.Properties.OrchestratorProfile.OrchestratorRelease = up.OrchestratorRelease
			cs.Properties.OrchestratorProfile.OrchestratorVersion = up.OrchestratorVersion
			return nil
		}
	}
	return a.Translator.Errorf("Kubernetes %s cannot be upgraded to %s",
		cs.Properties.OrchestratorProfile.OrchestratorRelease, unverOrch.OrchestratorRelease)
}

// SerializeContainerService takes an unversioned container service and returns the bytes
func (a *Apiloader) SerializeContainerService(containerService *ContainerService, version string) ([]byte, error) {
	if containerService.Properties != nil && containerService.Properties.HostedMasterProfile != nil {
		b, err := a.serializeHostedContainerService(containerService, version)
		if err == nil && b != nil {
			return b, nil
		}
	}

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
		return nil, a.Translator.Errorf("invalid version %s for conversion back from unversioned object", version)
	}
}

func (a *Apiloader) serializeHostedContainerService(containerService *ContainerService, version string) ([]byte, error) {
	switch version {
	case v20170831.APIVersion:
		v20170831ContainerService := ConvertContainerServiceToV20170831AgentPoolOnly(containerService)
		armContainerService := &V20170831ARMManagedContainerService{}
		armContainerService.ManagedCluster = v20170831ContainerService
		armContainerService.APIVersion = version
		b, err := json.MarshalIndent(armContainerService, "", "  ")
		if err != nil {
			return nil, err
		}
		return b, nil
	default:
		return nil, a.Translator.Errorf("invalid version %s for conversion back from unversioned object", version)
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

// Sets default container service property values for any appropriate zero values
func setContainerServiceDefaultsv20170701(c *v20170701.ContainerService) {
	if c.Properties.OrchestratorProfile != nil {
		c.Properties.OrchestratorProfile.OrchestratorVersion = ""
	}
}

// Sets default container service property values for any appropriate zero values
func setContainerServiceDefaultsvlabs(c *vlabs.ContainerService) {
	if c.Properties.OrchestratorProfile != nil {
		c.Properties.OrchestratorProfile.OrchestratorVersion = ""
	}
}

// Sets default HostedMaster property values for any appropriate zero values
func setManagedClusterDefaultsv20170831(hm *v20170831.ManagedCluster) {
	hm.Properties.KubernetesVersion = ""
}

// Sets default HostedMaster property values for any appropriate zero values
func setManagedClusterDefaultsvlabs(hm *apvlabs.ManagedCluster) {
	hm.Properties.KubernetesVersion = ""
}
