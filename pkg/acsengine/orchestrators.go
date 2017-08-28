package acsengine

import (
	"fmt"
	"strings"

	"github.com/Azure/acs-engine/pkg/api"
	"github.com/Azure/acs-engine/pkg/api/common"
)

type OrchestratorsInfo struct {
	Orchestrators []*api.OrchestratorInfo `json:"orchestrators"`
}

type orchestratorsFunc func(string) ([]*api.OrchestratorInfo, error)

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

func NewOrchestratorsInfo(orchestrator, release string) (*OrchestratorsInfo, error) {
	var err error
	if orchestrator, err = validate(orchestrator, release); err != nil {
		return nil, err
	}
	orch := &OrchestratorsInfo{}

	if len(orchestrator) == 0 {
		// return all orchestrators
		orch.Orchestrators = []*api.OrchestratorInfo{}
		for _, f := range funcmap {
			arr, err := f(release)
			if err != nil {
				return nil, err
			}
			orch.Orchestrators = append(orch.Orchestrators, arr...)
		}
		return orch, nil
	}
	if orch.Orchestrators, err = funcmap[orchestrator](release); err != nil {
		return nil, err
	}
	return orch, nil
}

func GetOrchestratorUpgradeInfo(cs *api.ContainerService) (*api.OrchestratorInfo, error) {
	if cs == nil || cs.Properties == nil || cs.Properties.OrchestratorProfile == nil {
		return nil, fmt.Errorf("Incomplete ContainerService")
	}
	if len(cs.Properties.OrchestratorProfile.OrchestratorRelease) == 0 {
		return nil, fmt.Errorf("Missing Orchestrator Release")
	}
	if cs.Properties.OrchestratorProfile.OrchestratorType != api.Kubernetes {
		return nil, fmt.Errorf("Upgrade operation is not supported for '%s'", cs.Properties.OrchestratorProfile.OrchestratorType)
	}
	orchestrator, err := validate(cs.Properties.OrchestratorProfile.OrchestratorType, cs.Properties.OrchestratorProfile.OrchestratorRelease)
	if err != nil {
		return nil, err
	}
	arr, err := funcmap[orchestrator](cs.Properties.OrchestratorProfile.OrchestratorRelease)
	if err != nil {
		return nil, err
	}
	// has to be exactly one element per specified orchestrator/release
	if len(arr) != 1 {
		return nil, fmt.Errorf("Umbiguous Orchestrator Releases")
	}
	return arr[0], nil
}

func kubernetesInfo(release string) ([]*api.OrchestratorInfo, error) {
	orchs := []*api.OrchestratorInfo{}
	if len(release) == 0 {
		// get info for all supported versions
		for rel, ver := range common.KubeReleaseToVersion {
			orchs = append(orchs,
				&api.OrchestratorInfo{
					Orchestrator: api.Kubernetes,
					VersionInfo: api.VersionInfo{
						Release: rel,
						Version: ver,
					},
					Default:    rel == common.KubernetesDefaultRelease,
					Upgradable: kubernetesUpgrades(rel),
				})
		}
	} else {
		// get info for the specified release
		ver, ok := common.KubeReleaseToVersion[release]
		if !ok {
			return nil, fmt.Errorf("Kubernetes release %s is not supported", release)
		}
		orchs = append(orchs,
			&api.OrchestratorInfo{
				Orchestrator: api.Kubernetes,
				VersionInfo: api.VersionInfo{
					Release: release,
					Version: ver,
				},
				Default:    release == common.KubernetesDefaultRelease,
				Upgradable: kubernetesUpgrades(release),
			})
	}
	return orchs, nil
}

func kubernetesUpgrades(release string) []*api.VersionInfo {
	ret := []*api.VersionInfo{}
	switch release {
	case common.KubernetesRelease1Dot5:
		ret = append(ret, &api.VersionInfo{Release: common.KubernetesRelease1Dot6, Version: common.KubeReleaseToVersion[common.KubernetesRelease1Dot6]})
	case common.KubernetesRelease1Dot6:
		ret = append(ret, &api.VersionInfo{Release: common.KubernetesRelease1Dot6, Version: common.KubeReleaseToVersion[common.KubernetesRelease1Dot6]})
		ret = append(ret, &api.VersionInfo{Release: common.KubernetesRelease1Dot7, Version: common.KubeReleaseToVersion[common.KubernetesRelease1Dot7]})
	case common.KubernetesRelease1Dot7:
		ret = append(ret, &api.VersionInfo{Release: common.KubernetesRelease1Dot7, Version: common.KubeReleaseToVersion[common.KubernetesRelease1Dot7]})
	}
	return ret
}

func dcosInfo(release string) ([]*api.OrchestratorInfo, error) {
	orchs := []*api.OrchestratorInfo{}
	if len(release) == 0 {
		// get info for all supported versions
		for rel, ver := range common.DCOSReleaseToVersion {
			orchs = append(orchs,
				&api.OrchestratorInfo{
					Orchestrator: api.DCOS,
					VersionInfo: api.VersionInfo{
						Release: rel,
						Version: ver,
					},
					Default: rel == common.DCOSDefaultRelease,
				})
		}
	} else {
		// get info for the specified release
		ver, ok := common.DCOSReleaseToVersion[release]
		if !ok {
			return nil, fmt.Errorf("DCOS release %s is not supported", release)
		}
		orchs = append(orchs,
			&api.OrchestratorInfo{
				Orchestrator: api.DCOS,
				VersionInfo: api.VersionInfo{
					Release: release,
					Version: ver,
				},
				Default: release == common.DCOSDefaultRelease,
			})
	}
	return orchs, nil
}

func swarmInfo(release string) ([]*api.OrchestratorInfo, error) {
	return []*api.OrchestratorInfo{
		&api.OrchestratorInfo{
			Orchestrator: api.Swarm,
			VersionInfo: api.VersionInfo{
				Version: SwarmVersion,
			},
			DockerComposeVersion: SwarmDockerComposeVersion,
		},
	}, nil
}

func dockerceInfo(release string) ([]*api.OrchestratorInfo, error) {
	return []*api.OrchestratorInfo{
		&api.OrchestratorInfo{
			Orchestrator: api.SwarmMode,
			VersionInfo: api.VersionInfo{
				Version: DockerCEVersion,
			},
			DockerComposeVersion: DockerCEDockerComposeVersion,
		},
	}, nil
}
