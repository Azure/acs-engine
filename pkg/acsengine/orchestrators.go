package acsengine

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Azure/acs-engine/pkg/api"
	"github.com/Azure/acs-engine/pkg/api/common"
)

type orchestratorsFunc func(*api.OrchestratorEdition) ([]*api.OrchestratorVersionProfile, error)

var funcmap map[string]orchestratorsFunc

type versionNumber struct {
	major, minor, patch int64
}

func (v *versionNumber) parse(ver string) (err error) {
	arr := strings.Split(ver, ".")
	if len(arr) != 3 {
		return fmt.Errorf("Illegal version format '%s'", ver)
	}
	if v.major, err = strconv.ParseInt(arr[0], 10, 32); err != nil {
		return
	}
	if v.minor, err = strconv.ParseInt(arr[1], 10, 32); err != nil {
		return
	}
	if v.patch, err = strconv.ParseInt(arr[2], 10, 32); err != nil {
		return
	}
	return nil
}

func (v *versionNumber) greaterThan(o *versionNumber) bool {
	// check major
	if v.major != o.major {
		return v.major > o.major
	}
	// same major; check minor
	if v.minor != o.minor {
		return v.minor > o.minor
	}
	// same minor; check patch
	return v.patch > o.patch
}

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
			arr, err := f(&api.OrchestratorEdition{})
			if err != nil {
				return nil, err
			}
			orch.Orchestrators = append(orch.Orchestrators, arr...)
		}
		return orch, nil
	}
	if orch.Orchestrators, err = funcmap[orchestrator](&api.OrchestratorEdition{OrchestratorRelease: release}); err != nil {
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
	arr, err := kubernetesInfo(&api.OrchestratorEdition{
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

func kubernetesInfo(csOrch *api.OrchestratorEdition) ([]*api.OrchestratorVersionProfile, error) {
	orchs := []*api.OrchestratorVersionProfile{}
	if len(csOrch.OrchestratorRelease) == 0 {
		// get info for all supported versions
		for rel, ver := range common.KubeReleaseToVersion {
			upgrades, err := kubernetesUpgrades(&api.OrchestratorEdition{OrchestratorRelease: rel, OrchestratorVersion: ver})
			if err != nil {
				return nil, err
			}
			orchs = append(orchs,
				&api.OrchestratorVersionProfile{
					OrchestratorType: api.Kubernetes,
					OrchestratorEdition: api.OrchestratorEdition{
						OrchestratorRelease: rel,
						OrchestratorVersion: ver,
					},
					Default:     rel == common.KubernetesDefaultRelease,
					Upgradables: upgrades,
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
				OrchestratorType: api.Kubernetes,
				OrchestratorEdition: api.OrchestratorEdition{
					OrchestratorRelease: csOrch.OrchestratorRelease,
					OrchestratorVersion: ver,
				},
				Default:     csOrch.OrchestratorRelease == common.KubernetesDefaultRelease,
				Upgradables: upgrades,
			})
	}
	return orchs, nil
}

func kubernetesUpgrades(csOrch *api.OrchestratorEdition) ([]*api.OrchestratorEdition, error) {
	ret := []*api.OrchestratorEdition{}
	var csVer versionNumber
	var err error

	if err = csVer.parse(csOrch.OrchestratorVersion); err != nil {
		return ret, err
	}
	switch csOrch.OrchestratorRelease {
	case common.KubernetesRelease1Dot5:
		// add next release
		ret = append(ret, &api.OrchestratorEdition{
			OrchestratorRelease: common.KubernetesRelease1Dot6,
			OrchestratorVersion: common.KubeReleaseToVersion[common.KubernetesRelease1Dot6],
		})
	case common.KubernetesRelease1Dot6:
		// check for patch upgrade
		if ret, err = addPatchUpgrade(ret, &csVer, csOrch.OrchestratorRelease); err != nil {
			return ret, err
		}
		// add next release
		ret = append(ret, &api.OrchestratorEdition{
			OrchestratorRelease: common.KubernetesRelease1Dot7,
			OrchestratorVersion: common.KubeReleaseToVersion[common.KubernetesRelease1Dot7],
		})
	case common.KubernetesRelease1Dot7:
		// check for patch upgrade
		if ret, err = addPatchUpgrade(ret, &csVer, csOrch.OrchestratorRelease); err != nil {
			return ret, err
		}
	}
	return ret, nil
}

func addPatchUpgrade(upgrades []*api.OrchestratorEdition, csVer *versionNumber, release string) ([]*api.OrchestratorEdition, error) {
	var pVer versionNumber
	patchVersion := common.KubeReleaseToVersion[release]
	if err := pVer.parse(patchVersion); err != nil {
		return upgrades, err
	}
	if pVer.greaterThan(csVer) {
		upgrades = append(upgrades, &api.OrchestratorEdition{OrchestratorRelease: release, OrchestratorVersion: patchVersion})
	}
	return upgrades, nil
}

func dcosInfo(csOrch *api.OrchestratorEdition) ([]*api.OrchestratorVersionProfile, error) {
	orchs := []*api.OrchestratorVersionProfile{}
	if len(csOrch.OrchestratorRelease) == 0 {
		// get info for all supported versions
		for rel, ver := range common.DCOSReleaseToVersion {
			orchs = append(orchs,
				&api.OrchestratorVersionProfile{
					OrchestratorType: api.DCOS,
					OrchestratorEdition: api.OrchestratorEdition{
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
				OrchestratorType: api.DCOS,
				OrchestratorEdition: api.OrchestratorEdition{
					OrchestratorRelease: csOrch.OrchestratorRelease,
					OrchestratorVersion: ver,
				},
				Default: csOrch.OrchestratorRelease == common.DCOSDefaultRelease,
			})
	}
	return orchs, nil
}

func swarmInfo(csOrch *api.OrchestratorEdition) ([]*api.OrchestratorVersionProfile, error) {
	return []*api.OrchestratorVersionProfile{
		{
			OrchestratorType: api.Swarm,
			OrchestratorEdition: api.OrchestratorEdition{
				OrchestratorVersion: SwarmVersion,
			},
		},
	}, nil
}

func dockerceInfo(csOrch *api.OrchestratorEdition) ([]*api.OrchestratorVersionProfile, error) {
	return []*api.OrchestratorVersionProfile{
		{
			OrchestratorType: api.SwarmMode,
			OrchestratorEdition: api.OrchestratorEdition{
				OrchestratorVersion: DockerCEVersion,
			},
		},
	}, nil
}
