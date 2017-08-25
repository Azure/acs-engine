package api

import (
	"fmt"
	"strings"

	"github.com/Azure/acs-engine/pkg/api/common"
)

const yes = "yes"

type OrchestratorInfo struct {
	Orchestrator         string `json:"orchestrator"`
	Release              string `json:"release,omitempty"`
	Version              string `json:"version"`
	DockerComposeVersion string `json:"docker-compose-version,omitempty"`
	Default              string `json:"default,omitempty"`
}

type OrchestratorsInfo struct {
	Orchestrators []OrchestratorInfo `json:"orchestrators"`
}

type orchestratorsFunc func(string) ([]OrchestratorInfo, error)

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

func NewOrchestratorsInfo(orchestrator, release string) (*OrchestratorsInfo, error) {
	var err error
	if orchestrator, err = validate(orchestrator, release); err != nil {
		return nil, err
	}

	orch := &OrchestratorsInfo{}

	funcmap := map[string]orchestratorsFunc{
		Kubernetes: kubernetesInfo,
		DCOS:       dcosInfo,
		Swarm:      swarmInfo,
		SwarmMode:  dockerceInfo,
	}

	if len(orchestrator) == 0 {
		// return all orchestrators
		orch.Orchestrators = []OrchestratorInfo{}
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

func kubernetesInfo(release string) ([]OrchestratorInfo, error) {
	orchs := []OrchestratorInfo{}
	if len(release) == 0 {
		// get info for all supported versions
		for rel, ver := range common.KubeReleaseToVersion {
			var def string
			if rel == common.KubernetesDefaultRelease {
				def = yes
			}
			orchs = append(orchs,
				OrchestratorInfo{
					Orchestrator: Kubernetes,
					Release:      rel,
					Version:      ver,
					Default:      def,
				})
		}
	} else {
		// get info for the specified release
		ver, ok := common.KubeReleaseToVersion[release]
		if !ok {
			return nil, fmt.Errorf("Kubernetes release %s is not supported", release)
		}
		var def string
		if release == common.KubernetesDefaultRelease {
			def = yes
		}
		orchs = append(orchs,
			OrchestratorInfo{
				Orchestrator: Kubernetes,
				Release:      release,
				Version:      ver,
				Default:      def,
			})
	}
	return orchs, nil
}

func dcosInfo(release string) ([]OrchestratorInfo, error) {
	orchs := []OrchestratorInfo{}
	if len(release) == 0 {
		// get info for all supported versions
		for rel, ver := range common.DCOSReleaseToVersion {
			var def string
			if rel == common.DCOSDefaultRelease {
				def = yes
			}
			orchs = append(orchs,
				OrchestratorInfo{
					Orchestrator: DCOS,
					Release:      rel,
					Version:      ver,
					Default:      def,
				})
		}
	} else {
		// get info for the specified release
		ver, ok := common.DCOSReleaseToVersion[release]
		if !ok {
			return nil, fmt.Errorf("DCOS release %s is not supported", release)
		}
		var def string
		if release == common.DCOSDefaultRelease {
			def = yes
		}
		orchs = append(orchs,
			OrchestratorInfo{
				Orchestrator: DCOS,
				Release:      release,
				Version:      ver,
				Default:      def,
			})
	}
	return orchs, nil
}

func swarmInfo(release string) ([]OrchestratorInfo, error) {
	return []OrchestratorInfo{
		OrchestratorInfo{
			Orchestrator:         Swarm,
			Version:              SwarmVersion,
			DockerComposeVersion: SwarmDockerComposeVersion,
		},
	}, nil
}

func dockerceInfo(release string) ([]OrchestratorInfo, error) {
	return []OrchestratorInfo{
		OrchestratorInfo{
			Orchestrator:         SwarmMode,
			Version:              DockerCEVersion,
			DockerComposeVersion: DockerCEDockerComposeVersion,
		},
	}, nil
}
