package api

import (
	"fmt"
	"strings"

	"github.com/Azure/acs-engine/pkg/api/common"
	"github.com/Azure/acs-engine/pkg/api/v20170930"
	"github.com/Azure/acs-engine/pkg/api/vlabs"
	"github.com/Masterminds/semver"
)

type orchestratorsFunc func(*OrchestratorProfile) ([]*OrchestratorVersionProfile, error)

var funcmap map[string]orchestratorsFunc
var versionsMap map[string][]string

func init() {
	funcmap = map[string]orchestratorsFunc{
		Kubernetes: kubernetesInfo,
		DCOS:       dcosInfo,
		Swarm:      swarmInfo,
		SwarmMode:  dockerceInfo,
		OpenShift:  openShiftInfo,
	}
	versionsMap = map[string][]string{
		Kubernetes: common.GetAllSupportedKubernetesVersions(),
		DCOS:       common.GetAllSupportedDCOSVersions(),
		Swarm:      common.GetAllSupportedSwarmVersions(),
		SwarmMode:  common.GetAllSupportedDockerCEVersions(),
		OpenShift:  common.GetAllSupportedOpenShiftVersions(),
	}
}

func validate(orchestrator, version string) (string, error) {
	switch {
	case strings.EqualFold(orchestrator, Kubernetes):
		return Kubernetes, nil
	case strings.EqualFold(orchestrator, DCOS):
		return DCOS, nil
	case strings.EqualFold(orchestrator, Swarm):
		return Swarm, nil
	case strings.EqualFold(orchestrator, SwarmMode):
		return SwarmMode, nil
	case strings.EqualFold(orchestrator, OpenShift):
		return OpenShift, nil
	case orchestrator == "":
		if version != "" {
			return "", fmt.Errorf("Must specify orchestrator for version '%s'", version)
		}
	default:
		return "", fmt.Errorf("Unsupported orchestrator '%s'", orchestrator)
	}
	return "", nil
}

func isVersionSupported(csOrch *OrchestratorProfile) bool {
	supported := false
	for _, version := range versionsMap[csOrch.OrchestratorType] {

		if version == csOrch.OrchestratorVersion {
			supported = true
			break
		}
	}
	return supported
}

// GetOrchestratorVersionProfileListVLabs returns vlabs OrchestratorVersionProfileList object per (optionally) specified orchestrator and version
func GetOrchestratorVersionProfileListVLabs(orchestrator, version string) (*vlabs.OrchestratorVersionProfileList, error) {
	apiOrchs, err := getOrchestratorVersionProfileList(orchestrator, version)
	if err != nil {
		return nil, err
	}
	orchList := &vlabs.OrchestratorVersionProfileList{}
	orchList.Orchestrators = []*vlabs.OrchestratorVersionProfile{}
	for _, orch := range apiOrchs {
		orchList.Orchestrators = append(orchList.Orchestrators, ConvertOrchestratorVersionProfileToVLabs(orch))
	}
	return orchList, nil
}

// GetOrchestratorVersionProfileListV20170930 returns v20170930 OrchestratorVersionProfileList object per (optionally) specified orchestrator and version
func GetOrchestratorVersionProfileListV20170930(orchestrator, version string) (*v20170930.OrchestratorVersionProfileList, error) {
	apiOrchs, err := getOrchestratorVersionProfileList(orchestrator, version)
	if err != nil {
		return nil, err
	}
	orchList := &v20170930.OrchestratorVersionProfileList{}
	for _, orch := range apiOrchs {
		orchList.Properties.Orchestrators = append(orchList.Properties.Orchestrators, ConvertOrchestratorVersionProfileToV20170930(orch))
	}
	return orchList, nil
}

func getOrchestratorVersionProfileList(orchestrator, version string) ([]*OrchestratorVersionProfile, error) {
	var err error
	if orchestrator, err = validate(orchestrator, version); err != nil {
		return nil, err
	}
	orchs := []*OrchestratorVersionProfile{}
	if len(orchestrator) == 0 {
		// return all orchestrators
		for _, f := range funcmap {
			arr, err := f(&OrchestratorProfile{})
			if err != nil {
				return nil, err
			}
			orchs = append(orchs, arr...)
		}
	} else {
		if orchs, err = funcmap[orchestrator](&OrchestratorProfile{OrchestratorType: orchestrator, OrchestratorVersion: version}); err != nil {
			return nil, err
		}
	}
	return orchs, nil
}

// GetOrchestratorVersionProfile returns orchestrator info for upgradable container service
func GetOrchestratorVersionProfile(orch *OrchestratorProfile) (*OrchestratorVersionProfile, error) {
	if orch.OrchestratorVersion == "" {
		return nil, fmt.Errorf("Missing Orchestrator Version")
	}
	switch orch.OrchestratorType {
	case Kubernetes, DCOS:
		arr, err := funcmap[orch.OrchestratorType](orch)
		if err != nil {
			return nil, err
		}
		// has to be exactly one element per specified orchestrator/version
		if len(arr) != 1 {
			return nil, fmt.Errorf("Umbiguous Orchestrator Versions")
		}
		return arr[0], nil
	default:
		return nil, fmt.Errorf("Upgrade operation is not supported for '%s'", orch.OrchestratorType)
	}
}

func kubernetesInfo(csOrch *OrchestratorProfile) ([]*OrchestratorVersionProfile, error) {
	orchs := []*OrchestratorVersionProfile{}
	if csOrch.OrchestratorVersion == "" {
		// get info for all supported versions
		for _, ver := range common.GetAllSupportedKubernetesVersions() {
			upgrades, err := kubernetesUpgrades(&OrchestratorProfile{OrchestratorVersion: ver})
			if err != nil {
				return nil, err
			}
			orchs = append(orchs,
				&OrchestratorVersionProfile{
					OrchestratorProfile: OrchestratorProfile{
						OrchestratorType:    Kubernetes,
						OrchestratorVersion: ver,
					},
					Default:  ver == common.GetDefaultKubernetesVersion(),
					Upgrades: upgrades,
				})
		}
	} else {
		if !isVersionSupported(csOrch) {
			return nil, fmt.Errorf("Kubernetes version %s is not supported", csOrch.OrchestratorVersion)
		}

		upgrades, err := kubernetesUpgrades(csOrch)
		if err != nil {
			return nil, err
		}
		orchs = append(orchs,
			&OrchestratorVersionProfile{
				OrchestratorProfile: OrchestratorProfile{
					OrchestratorType:    Kubernetes,
					OrchestratorVersion: csOrch.OrchestratorVersion,
				},
				Default:  csOrch.OrchestratorVersion == common.GetDefaultKubernetesVersion(),
				Upgrades: upgrades,
			})
	}
	return orchs, nil
}

func kubernetesUpgrades(csOrch *OrchestratorProfile) ([]*OrchestratorProfile, error) {
	ret := []*OrchestratorProfile{}

	currentVer, err := semver.NewVersion(csOrch.OrchestratorVersion)
	if err != nil {
		return nil, err
	}
	nextNextMinorVersion := currentVer.IncMinor().IncMinor()
	upgradeableVersions := common.GetVersionsBetween(common.GetAllSupportedKubernetesVersions(), csOrch.OrchestratorVersion, nextNextMinorVersion.String(), false, true)
	for _, ver := range upgradeableVersions {
		ret = append(ret, &OrchestratorProfile{
			OrchestratorType:    Kubernetes,
			OrchestratorVersion: ver,
		})
	}
	return ret, nil
}

func dcosInfo(csOrch *OrchestratorProfile) ([]*OrchestratorVersionProfile, error) {
	orchs := []*OrchestratorVersionProfile{}
	if csOrch.OrchestratorVersion == "" {
		// get info for all supported versions
		for _, ver := range common.AllDCOSSupportedVersions {
			upgrades, err := dcosUpgrades(&OrchestratorProfile{OrchestratorVersion: ver})
			if err != nil {
				return nil, err
			}
			orchs = append(orchs,
				&OrchestratorVersionProfile{
					OrchestratorProfile: OrchestratorProfile{
						OrchestratorType:    DCOS,
						OrchestratorVersion: ver,
					},
					Default:  ver == common.DCOSDefaultVersion,
					Upgrades: upgrades,
				})
		}
	} else {
		if !isVersionSupported(csOrch) {
			return nil, fmt.Errorf("DCOS version %s is not supported", csOrch.OrchestratorVersion)
		}

		// get info for the specified version
		upgrades, err := dcosUpgrades(csOrch)
		if err != nil {
			return nil, err
		}
		orchs = append(orchs,
			&OrchestratorVersionProfile{
				OrchestratorProfile: OrchestratorProfile{
					OrchestratorType:    DCOS,
					OrchestratorVersion: csOrch.OrchestratorVersion,
				},
				Default:  csOrch.OrchestratorVersion == common.DCOSDefaultVersion,
				Upgrades: upgrades,
			})
	}
	return orchs, nil
}

func dcosUpgrades(csOrch *OrchestratorProfile) ([]*OrchestratorProfile, error) {
	ret := []*OrchestratorProfile{}

	switch csOrch.OrchestratorVersion {
	case common.DCOSVersion1Dot11Dot0:
		ret = append(ret, &OrchestratorProfile{
			OrchestratorType:    DCOS,
			OrchestratorVersion: common.DCOSVersion1Dot11Dot2,
		})
	}
	return ret, nil
}

func swarmInfo(csOrch *OrchestratorProfile) ([]*OrchestratorVersionProfile, error) {
	if csOrch.OrchestratorVersion == "" {
		return []*OrchestratorVersionProfile{
			{
				OrchestratorProfile: OrchestratorProfile{
					OrchestratorType:    Swarm,
					OrchestratorVersion: SwarmVersion,
				},
			},
		}, nil
	}

	if !isVersionSupported(csOrch) {
		return nil, fmt.Errorf("Swarm version %s is not supported", csOrch.OrchestratorVersion)
	}
	return []*OrchestratorVersionProfile{
		{
			OrchestratorProfile: OrchestratorProfile{
				OrchestratorType:    Swarm,
				OrchestratorVersion: csOrch.OrchestratorVersion,
			},
		},
	}, nil
}

func dockerceInfo(csOrch *OrchestratorProfile) ([]*OrchestratorVersionProfile, error) {

	if csOrch.OrchestratorVersion == "" {
		return []*OrchestratorVersionProfile{
			{
				OrchestratorProfile: OrchestratorProfile{
					OrchestratorType:    SwarmMode,
					OrchestratorVersion: DockerCEVersion,
				},
			},
		}, nil
	}

	if !isVersionSupported(csOrch) {
		return nil, fmt.Errorf("Docker CE version %s is not supported", csOrch.OrchestratorVersion)
	}
	return []*OrchestratorVersionProfile{
		{
			OrchestratorProfile: OrchestratorProfile{
				OrchestratorType:    SwarmMode,
				OrchestratorVersion: csOrch.OrchestratorVersion,
			},
		},
	}, nil
}

func openShiftInfo(csOrch *OrchestratorProfile) ([]*OrchestratorVersionProfile, error) {
	orchs := []*OrchestratorVersionProfile{}
	if csOrch.OrchestratorVersion == "" {
		// get info for all supported versions
		for _, ver := range common.GetAllSupportedOpenShiftVersions() {
			if ver == common.OpenShiftVersionUnstable {
				continue
			}
			// TODO: populate OrchestratorVersionProfile.Upgrades
			orchs = append(orchs,
				&OrchestratorVersionProfile{
					OrchestratorProfile: OrchestratorProfile{
						OrchestratorType:    OpenShift,
						OrchestratorVersion: ver,
					},
					Default: ver == common.OpenShiftDefaultVersion,
				})
		}
	} else {
		if !isVersionSupported(csOrch) {
			return nil, fmt.Errorf("OpenShift version %s is not supported", csOrch.OrchestratorVersion)
		}

		// TODO: populate OrchestratorVersionProfile.Upgrades
		orchs = append(orchs,
			&OrchestratorVersionProfile{
				OrchestratorProfile: OrchestratorProfile{
					OrchestratorType:    OpenShift,
					OrchestratorVersion: csOrch.OrchestratorVersion,
				},
				Default: csOrch.OrchestratorVersion == common.OpenShiftDefaultVersion,
			})
	}
	return orchs, nil
}
