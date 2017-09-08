package acsengine

import (
	"fmt"
	"strings"

	"github.com/Azure/acs-engine/pkg/api"
	"github.com/Azure/acs-engine/pkg/api/common"
	"github.com/Masterminds/semver"
)

type orchestratorsFunc func(*api.OrchestratorProfile) ([]*api.OrchestratorVersionProfile, error)

var funcmap map[string]orchestratorsFunc

func init() {
	funcmap = map[string]orchestratorsFunc{
		api.Kubernetes: kubernetesInfo,
		api.DCOS:       dcosInfo,
		api.Swarm:      swarmInfo,
		api.SwarmMode:  dockerceInfo,
	}
}

func validate(orchestrator, release string) (string, error) {
	switch {
	case strings.EqualFold(orchestrator, api.Kubernetes):
		return api.Kubernetes, nil
	case strings.EqualFold(orchestrator, api.DCOS):
		return api.DCOS, nil
	case strings.EqualFold(orchestrator, api.Swarm):
		return api.Swarm, nil
	case strings.EqualFold(orchestrator, api.SwarmMode):
		return api.SwarmMode, nil
	case len(orchestrator) == 0:
		if len(release) > 0 {
			return "", fmt.Errorf("Must specify orchestrator for release '%s'", release)
		}
	default:
		return "", fmt.Errorf("Unsupported orchestrator '%s'", orchestrator)
	}
	return "", nil
}

// GetOrchestratorVersionProfileList returns OrchestratorVersionProfileList object per (optionally) specified orchestrator and release
func GetOrchestratorVersionProfileList(orchestrator, release string) (*api.OrchestratorVersionProfileList, error) {
	var err error
	if orchestrator, err = validate(orchestrator, release); err != nil {
		return nil, err
	}
	orch := &api.OrchestratorVersionProfileList{}

	if len(orchestrator) == 0 {
		// return all orchestrators
		orch.Orchestrators = []*api.OrchestratorVersionProfile{}
		for _, f := range funcmap {
			arr, err := f(&api.OrchestratorProfile{})
			if err != nil {
				return nil, err
			}
			orch.Orchestrators = append(orch.Orchestrators, arr...)
		}
		return orch, nil
	}
	if orch.Orchestrators, err = funcmap[orchestrator](&api.OrchestratorProfile{OrchestratorRelease: release}); err != nil {
		return nil, err
	}
	return orch, nil
}

// GetOrchestratorVersionProfile returns orchestrator info for upgradable container service
func GetOrchestratorVersionProfile(cs *api.ContainerService) (*api.OrchestratorVersionProfile, error) {
	if cs == nil || cs.Properties == nil || cs.Properties.OrchestratorProfile == nil {
		return nil, fmt.Errorf("Incomplete ContainerService")
	}
	if len(cs.Properties.OrchestratorProfile.OrchestratorRelease) == 0 {
		return nil, fmt.Errorf("Missing Orchestrator Release")
	}
	if cs.Properties.OrchestratorProfile.OrchestratorType != api.Kubernetes {
		return nil, fmt.Errorf("Upgrade operation is not supported for '%s'", cs.Properties.OrchestratorProfile.OrchestratorType)
	}
	arr, err := kubernetesInfo(&api.OrchestratorProfile{
		OrchestratorRelease: cs.Properties.OrchestratorProfile.OrchestratorRelease,
		OrchestratorVersion: cs.Properties.OrchestratorProfile.OrchestratorVersion})
	if err != nil {
		return nil, err
	}
	// has to be exactly one element per specified orchestrator/release
	if len(arr) != 1 {
		return nil, fmt.Errorf("Umbiguous Orchestrator Releases")
	}
	return arr[0], nil
}

func kubernetesInfo(csOrch *api.OrchestratorProfile) ([]*api.OrchestratorVersionProfile, error) {
	orchs := []*api.OrchestratorVersionProfile{}
	if len(csOrch.OrchestratorRelease) == 0 {
		// get info for all supported versions
		for rel, ver := range common.KubeReleaseToVersion {
			upgrades, err := kubernetesUpgrades(&api.OrchestratorProfile{OrchestratorRelease: rel, OrchestratorVersion: ver})
			if err != nil {
				return nil, err
			}
			orchs = append(orchs,
				&api.OrchestratorVersionProfile{
					OrchestratorProfile: api.OrchestratorProfile{
						OrchestratorType:    api.Kubernetes,
						OrchestratorRelease: rel,
						OrchestratorVersion: ver,
					},
					Default:  rel == common.KubernetesDefaultRelease,
					Upgrades: upgrades,
				})
		}
	} else {
		// get info for the specified release
		ver, ok := common.KubeReleaseToVersion[csOrch.OrchestratorRelease]
		if !ok {
			return nil, fmt.Errorf("Kubernetes release %s is not supported", csOrch.OrchestratorRelease)
		}
		// set default version if empty
		if len(csOrch.OrchestratorVersion) == 0 {
			csOrch.OrchestratorVersion = ver
		}
		upgrades, err := kubernetesUpgrades(csOrch)
		if err != nil {
			return nil, err
		}
		orchs = append(orchs,
			&api.OrchestratorVersionProfile{
				OrchestratorProfile: api.OrchestratorProfile{
					OrchestratorType:    api.Kubernetes,
					OrchestratorRelease: csOrch.OrchestratorRelease,
					OrchestratorVersion: ver,
				},
				Default:  csOrch.OrchestratorRelease == common.KubernetesDefaultRelease,
				Upgrades: upgrades,
			})
	}
	return orchs, nil
}

func kubernetesUpgrades(csOrch *api.OrchestratorProfile) ([]*api.UpgradeContainerService, error) {
	ret := []*api.UpgradeContainerService{}
	var err error

	switch csOrch.OrchestratorRelease {
	case common.KubernetesRelease1Dot5:
		// add next release
		ret = append(ret, &api.UpgradeContainerService{
			OrchestratorType:    api.Kubernetes,
			OrchestratorRelease: common.KubernetesRelease1Dot6,
			OrchestratorVersion: common.KubeReleaseToVersion[common.KubernetesRelease1Dot6],
		})
	case common.KubernetesRelease1Dot6:
		// check for patch upgrade
		if ret, err = addPatchUpgrade(ret, csOrch.OrchestratorRelease, csOrch.OrchestratorVersion); err != nil {
			return ret, err
		}
		// add next release
		ret = append(ret, &api.UpgradeContainerService{
			OrchestratorType:    api.Kubernetes,
			OrchestratorRelease: common.KubernetesRelease1Dot7,
			OrchestratorVersion: common.KubeReleaseToVersion[common.KubernetesRelease1Dot7],
		})
	case common.KubernetesRelease1Dot7:
		// check for patch upgrade
		if ret, err = addPatchUpgrade(ret, csOrch.OrchestratorRelease, csOrch.OrchestratorVersion); err != nil {
			return ret, err
		}
	}
	return ret, nil
}

func addPatchUpgrade(upgrades []*api.UpgradeContainerService, release, version string) ([]*api.UpgradeContainerService, error) {
	patchVersion, ok := common.KubeReleaseToVersion[release]
	if !ok {
		return upgrades, fmt.Errorf("Kubernetes release %s is not supported", release)
	}
	pVer, err := semver.NewVersion(patchVersion)
	if err != nil {
		return upgrades, err
	}
	constraint, err := semver.NewConstraint(">" + version)
	if err != nil {
		return upgrades, err
	}
	if constraint.Check(pVer) {
		upgrades = append(upgrades, &api.UpgradeContainerService{OrchestratorRelease: release, OrchestratorVersion: patchVersion})
	}
	return upgrades, nil
}

func dcosInfo(csOrch *api.OrchestratorProfile) ([]*api.OrchestratorVersionProfile, error) {
	orchs := []*api.OrchestratorVersionProfile{}
	if len(csOrch.OrchestratorRelease) == 0 {
		// get info for all supported versions
		for rel, ver := range common.DCOSReleaseToVersion {
			orchs = append(orchs,
				&api.OrchestratorVersionProfile{
					OrchestratorProfile: api.OrchestratorProfile{
						OrchestratorType:    api.DCOS,
						OrchestratorRelease: rel,
						OrchestratorVersion: ver,
					},
					Default: rel == common.DCOSDefaultRelease,
				})
		}
	} else {
		// get info for the specified release
		ver, ok := common.DCOSReleaseToVersion[csOrch.OrchestratorRelease]
		if !ok {
			return nil, fmt.Errorf("DCOS release %s is not supported", csOrch.OrchestratorRelease)
		}
		orchs = append(orchs,
			&api.OrchestratorVersionProfile{
				OrchestratorProfile: api.OrchestratorProfile{
					OrchestratorType:    api.DCOS,
					OrchestratorRelease: csOrch.OrchestratorRelease,
					OrchestratorVersion: ver,
				},
				Default: csOrch.OrchestratorRelease == common.DCOSDefaultRelease,
			})
	}
	return orchs, nil
}

func swarmInfo(csOrch *api.OrchestratorProfile) ([]*api.OrchestratorVersionProfile, error) {
	return []*api.OrchestratorVersionProfile{
		{
			OrchestratorProfile: api.OrchestratorProfile{
				OrchestratorType:    api.Swarm,
				OrchestratorVersion: SwarmVersion,
			},
		},
	}, nil
}

func dockerceInfo(csOrch *api.OrchestratorProfile) ([]*api.OrchestratorVersionProfile, error) {
	return []*api.OrchestratorVersionProfile{
		{
			OrchestratorProfile: api.OrchestratorProfile{
				OrchestratorType:    api.SwarmMode,
				OrchestratorVersion: DockerCEVersion,
			},
		},
	}, nil
}
