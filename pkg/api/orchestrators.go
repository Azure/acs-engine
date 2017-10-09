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

func init() {
	funcmap = map[string]orchestratorsFunc{
		Kubernetes: kubernetesInfo,
		DCOS:       dcosInfo,
		Swarm:      swarmInfo,
		SwarmMode:  dockerceInfo,
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
	case orchestrator == "":
		if version != "" {
			return "", fmt.Errorf("Must specify orchestrator for version '%s'", version)
		}
	default:
		return "", fmt.Errorf("Unsupported orchestrator '%s'", orchestrator)
	}
	return "", nil
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
	orchList.Orchestrators = []*v20170930.OrchestratorVersionProfile{}
	for _, orch := range apiOrchs {
		orchList.Orchestrators = append(orchList.Orchestrators, ConvertOrchestratorVersionProfileToV20170930(orch))
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
		if orchs, err = funcmap[orchestrator](&OrchestratorProfile{OrchestratorVersion: version}); err != nil {
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
	if orch.OrchestratorType != Kubernetes {
		return nil, fmt.Errorf("Upgrade operation is not supported for '%s'", orch.OrchestratorType)
	}
	arr, err := kubernetesInfo(orch)
	if err != nil {
		return nil, err
	}
	// has to be exactly one element per specified orchestrator/version
	if len(arr) != 1 {
		return nil, fmt.Errorf("Umbiguous Orchestrator Versions")
	}
	return arr[0], nil
}

func getUpgradesV20170930(orch *OrchestratorVersionProfile, allowCurrentVersionUpgrade bool) []*v20170930.OrchestratorProfile {
	var upgrades []*v20170930.OrchestratorProfile
	if orch.Upgrades != nil {
		upgrades = make([]*v20170930.OrchestratorProfile, len(orch.Upgrades))
		for i, h := range orch.Upgrades {
			upgrades[i] = &v20170930.OrchestratorProfile{
				OrchestratorType:    orch.OrchestratorType,
				OrchestratorVersion: h.OrchestratorVersion,
			}
		}
	}
	// add current version if upgrade has failed
	if allowCurrentVersionUpgrade {
		upgrades = append(upgrades, &v20170930.OrchestratorProfile{
			OrchestratorType:    orch.OrchestratorType,
			OrchestratorVersion: orch.OrchestratorVersion,
		})
	}
	return upgrades
}

func kubernetesInfo(csOrch *OrchestratorProfile) ([]*OrchestratorVersionProfile, error) {
	orchs := []*OrchestratorVersionProfile{}
	if csOrch.OrchestratorVersion == "" {
		// get info for all supported versions
		for _, ver := range common.AllKubernetesSupportedVersions {
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
					Default:  ver == common.KubernetesDefaultVersion,
					Upgrades: upgrades,
				})
		}
	} else {
		ver, _ := semver.NewVersion(csOrch.OrchestratorVersion)
		cons, _ := semver.NewConstraint("<1.5.0")
		if cons.Check(ver) {
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
				Default:  csOrch.OrchestratorVersion == common.KubernetesDefaultVersion,
				Upgrades: upgrades,
			})
	}
	return orchs, nil
}

func kubernetesUpgrades(csOrch *OrchestratorProfile) ([]*OrchestratorProfile, error) {
	ret := []*OrchestratorProfile{}
	var err error

	switch {
	case strings.HasPrefix(csOrch.OrchestratorVersion, "1.5"):
		// add next version
		ret = append(ret, &OrchestratorProfile{
			OrchestratorType:    Kubernetes,
			OrchestratorVersion: common.KubernetesVersion1Dot6Dot11,
		})
	case strings.HasPrefix(csOrch.OrchestratorVersion, "1.6"):
		// check for patch upgrade
		if ret, err = addPatchUpgrade(ret, csOrch.OrchestratorVersion, common.KubernetesVersion1Dot6Dot11); err != nil {
			return ret, err
		}
		// add next version
		ret = append(ret, &OrchestratorProfile{
			OrchestratorType:    Kubernetes,
			OrchestratorVersion: common.KubernetesVersion1Dot7Dot7,
		})
	case strings.HasPrefix(csOrch.OrchestratorVersion, "1.7"):
		// check for patch upgrade
		if ret, err = addPatchUpgrade(ret, csOrch.OrchestratorVersion, common.KubernetesVersion1Dot7Dot7); err != nil {
			return ret, err
		}
	}
	return ret, nil
}

func addPatchUpgrade(upgrades []*OrchestratorProfile, version, targetVersion string) ([]*OrchestratorProfile, error) {
	pVer, err := semver.NewVersion(targetVersion)
	if err != nil {
		return upgrades, err
	}
	constraint, err := semver.NewConstraint(">" + version)
	if err != nil {
		return upgrades, err
	}
	if constraint.Check(pVer) {
		upgrades = append(upgrades, &OrchestratorProfile{OrchestratorVersion: targetVersion})
	}
	return upgrades, nil
}

func dcosInfo(csOrch *OrchestratorProfile) ([]*OrchestratorVersionProfile, error) {
	orchs := []*OrchestratorVersionProfile{}
	if csOrch.OrchestratorVersion == "" {
		// get info for all supported versions
		for _, ver := range common.AllDCOSSupportedVersions {
			orchs = append(orchs,
				&OrchestratorVersionProfile{
					OrchestratorProfile: OrchestratorProfile{
						OrchestratorType:    DCOS,
						OrchestratorVersion: ver,
					},
					Default: ver == common.DCOSDefaultVersion,
				})
		}
	} else {
		// get info for the specified version
		orchs = append(orchs,
			&OrchestratorVersionProfile{
				OrchestratorProfile: OrchestratorProfile{
					OrchestratorType:    DCOS,
					OrchestratorVersion: csOrch.OrchestratorVersion,
				},
				Default: csOrch.OrchestratorVersion == common.DCOSDefaultVersion,
			})
	}
	return orchs, nil
}

func swarmInfo(csOrch *OrchestratorProfile) ([]*OrchestratorVersionProfile, error) {
	return []*OrchestratorVersionProfile{
		{
			OrchestratorProfile: OrchestratorProfile{
				OrchestratorType:    Swarm,
				OrchestratorVersion: SwarmVersion,
			},
		},
	}, nil
}

func dockerceInfo(csOrch *OrchestratorProfile) ([]*OrchestratorVersionProfile, error) {
	return []*OrchestratorVersionProfile{
		{
			OrchestratorProfile: OrchestratorProfile{
				OrchestratorType:    SwarmMode,
				OrchestratorVersion: DockerCEVersion,
			},
		},
	}, nil
}
