package api

import (
	"fmt"
	"strings"

	"github.com/Azure/acs-engine/pkg/api/common"
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

func validate(orchestrator, release string) (string, error) {
	switch {
	case strings.EqualFold(orchestrator, Kubernetes):
		return Kubernetes, nil
	case strings.EqualFold(orchestrator, DCOS):
		return DCOS, nil
	case strings.EqualFold(orchestrator, Swarm):
		return Swarm, nil
	case strings.EqualFold(orchestrator, SwarmMode):
		return SwarmMode, nil
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
func GetOrchestratorVersionProfileList(orchestrator, release string) (*OrchestratorVersionProfileList, error) {
	var err error
	if orchestrator, err = validate(orchestrator, release); err != nil {
		return nil, err
	}
	orch := &OrchestratorVersionProfileList{}

	if len(orchestrator) == 0 {
		// return all orchestrators
		orch.Orchestrators = []*OrchestratorVersionProfile{}
		for _, f := range funcmap {
			arr, err := f(&OrchestratorProfile{})
			if err != nil {
				return nil, err
			}
			orch.Orchestrators = append(orch.Orchestrators, arr...)
		}
		return orch, nil
	}
	if orch.Orchestrators, err = funcmap[orchestrator](&OrchestratorProfile{OrchestratorRelease: release}); err != nil {
		return nil, err
	}
	return orch, nil
}

// GetOrchestratorVersionProfile returns orchestrator info for upgradable container service
func GetOrchestratorVersionProfile(orch *OrchestratorProfile) (*OrchestratorVersionProfile, error) {
	if len(orch.OrchestratorRelease) == 0 {
		return nil, fmt.Errorf("Missing Orchestrator Release")
	}
	if orch.OrchestratorType != Kubernetes {
		return nil, fmt.Errorf("Upgrade operation is not supported for '%s'", orch.OrchestratorType)
	}
	arr, err := kubernetesInfo(orch)
	if err != nil {
		return nil, err
	}
	// has to be exactly one element per specified orchestrator/release
	if len(arr) != 1 {
		return nil, fmt.Errorf("Umbiguous Orchestrator Releases")
	}
	return arr[0], nil
}

func kubernetesInfo(csOrch *OrchestratorProfile) ([]*OrchestratorVersionProfile, error) {
	orchs := []*OrchestratorVersionProfile{}
	if len(csOrch.OrchestratorRelease) == 0 {
		// get info for all supported versions
		for rel, ver := range common.KubeReleaseToVersion {
			upgrades, err := kubernetesUpgrades(&OrchestratorProfile{OrchestratorRelease: rel, OrchestratorVersion: ver})
			if err != nil {
				return nil, err
			}
			orchs = append(orchs,
				&OrchestratorVersionProfile{
					OrchestratorProfile: OrchestratorProfile{
						OrchestratorType:    Kubernetes,
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
			&OrchestratorVersionProfile{
				OrchestratorProfile: OrchestratorProfile{
					OrchestratorType:    Kubernetes,
					OrchestratorRelease: csOrch.OrchestratorRelease,
					OrchestratorVersion: ver,
				},
				Default:  csOrch.OrchestratorRelease == common.KubernetesDefaultRelease,
				Upgrades: upgrades,
			})
	}
	return orchs, nil
}

func kubernetesUpgrades(csOrch *OrchestratorProfile) ([]*OrchestratorProfile, error) {
	ret := []*OrchestratorProfile{}
	var err error

	switch csOrch.OrchestratorRelease {
	case common.KubernetesRelease1Dot5:
		// add next release
		ret = append(ret, &OrchestratorProfile{
			OrchestratorType:    Kubernetes,
			OrchestratorRelease: common.KubernetesRelease1Dot6,
			OrchestratorVersion: common.KubeReleaseToVersion[common.KubernetesRelease1Dot6],
		})
	case common.KubernetesRelease1Dot6:
		// check for patch upgrade
		if ret, err = addPatchUpgrade(ret, csOrch.OrchestratorRelease, csOrch.OrchestratorVersion); err != nil {
			return ret, err
		}
		// add next release
		ret = append(ret, &OrchestratorProfile{
			OrchestratorType:    Kubernetes,
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

func addPatchUpgrade(upgrades []*OrchestratorProfile, release, version string) ([]*OrchestratorProfile, error) {
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
		upgrades = append(upgrades, &OrchestratorProfile{OrchestratorRelease: release, OrchestratorVersion: patchVersion})
	}
	return upgrades, nil
}

func dcosInfo(csOrch *OrchestratorProfile) ([]*OrchestratorVersionProfile, error) {
	orchs := []*OrchestratorVersionProfile{}
	if len(csOrch.OrchestratorRelease) == 0 {
		// get info for all supported versions
		for rel, ver := range common.DCOSReleaseToVersion {
			orchs = append(orchs,
				&OrchestratorVersionProfile{
					OrchestratorProfile: OrchestratorProfile{
						OrchestratorType:    DCOS,
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
			&OrchestratorVersionProfile{
				OrchestratorProfile: OrchestratorProfile{
					OrchestratorType:    DCOS,
					OrchestratorRelease: csOrch.OrchestratorRelease,
					OrchestratorVersion: ver,
				},
				Default: csOrch.OrchestratorRelease == common.DCOSDefaultRelease,
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
