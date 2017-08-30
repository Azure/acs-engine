package acsengine

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Azure/acs-engine/pkg/api"
	"github.com/Azure/acs-engine/pkg/api/common"
)

// OrchestratorInfos contains list of release info for supported orchestrators
type OrchestratorInfos struct {
	Orchestrators []*api.OrchestratorInfo `json:"orchestrators"`
}

type orchestratorsFunc func(*api.VersionInfo) ([]*api.OrchestratorInfo, error)

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

// NewOrchestratorInfos returns OrchestratorInfos object per (optionally) specified orchestrator and release
func NewOrchestratorInfos(orchestrator, release string) (*OrchestratorInfos, error) {
	var err error
	if orchestrator, err = validate(orchestrator, release); err != nil {
		return nil, err
	}
	orch := &OrchestratorInfos{}

	if len(orchestrator) == 0 {
		// return all orchestrators
		orch.Orchestrators = []*api.OrchestratorInfo{}
		for _, f := range funcmap {
			arr, err := f(&api.VersionInfo{})
			if err != nil {
				return nil, err
			}
			orch.Orchestrators = append(orch.Orchestrators, arr...)
		}
		return orch, nil
	}
	if orch.Orchestrators, err = funcmap[orchestrator](&api.VersionInfo{Release: release}); err != nil {
		return nil, err
	}
	return orch, nil
}

// GetOrchestratorUpgradeInfo returns orchestrator info for upgradable container service
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
	arr, err := kubernetesInfo(&api.VersionInfo{
		Release: cs.Properties.OrchestratorProfile.OrchestratorRelease,
		Version: cs.Properties.OrchestratorProfile.OrchestratorVersion})
	if err != nil {
		return nil, err
	}
	// has to be exactly one element per specified orchestrator/release
	if len(arr) != 1 {
		return nil, fmt.Errorf("Umbiguous Orchestrator Releases")
	}
	return arr[0], nil
}

func kubernetesInfo(csInfo *api.VersionInfo) ([]*api.OrchestratorInfo, error) {
	orchs := []*api.OrchestratorInfo{}
	if len(csInfo.Release) == 0 {
		// get info for all supported versions
		for rel, ver := range common.KubeReleaseToVersion {
			upgrades, err := kubernetesUpgrades(&api.VersionInfo{Release: rel, Version: ver})
			if err != nil {
				return nil, err
			}
			orchs = append(orchs,
				&api.OrchestratorInfo{
					Orchestrator: api.Kubernetes,
					VersionInfo: api.VersionInfo{
						Release: rel,
						Version: ver,
					},
					Default:    rel == common.KubernetesDefaultRelease,
					Upgradable: upgrades,
				})
		}
	} else {
		// get info for the specified release
		ver, ok := common.KubeReleaseToVersion[csInfo.Release]
		if !ok {
			return nil, fmt.Errorf("Kubernetes release %s is not supported", csInfo.Release)
		}
		// set defaulr version if empty
		if len(csInfo.Version) == 0 {
			csInfo.Version = ver
		}
		upgrades, err := kubernetesUpgrades(csInfo)
		if err != nil {
			return nil, err
		}
		orchs = append(orchs,
			&api.OrchestratorInfo{
				Orchestrator: api.Kubernetes,
				VersionInfo: api.VersionInfo{
					Release: csInfo.Release,
					Version: ver,
				},
				Default:    csInfo.Release == common.KubernetesDefaultRelease,
				Upgradable: upgrades,
			})
	}
	return orchs, nil
}

func kubernetesUpgrades(csInfo *api.VersionInfo) ([]*api.VersionInfo, error) {
	ret := []*api.VersionInfo{}
	var csVer, pVer versionNumber

	if err := csVer.parse(csInfo.Version); err != nil {
		return ret, err
	}
	switch csInfo.Release {
	case common.KubernetesRelease1Dot5:
		// add next release
		ret = append(ret, &api.VersionInfo{Release: common.KubernetesRelease1Dot6, Version: common.KubeReleaseToVersion[common.KubernetesRelease1Dot6]})
	case common.KubernetesRelease1Dot6:
		// check for patch upgrade
		patchVersion := common.KubeReleaseToVersion[common.KubernetesRelease1Dot6]
		if err := pVer.parse(patchVersion); err != nil {
			return ret, err
		}
		if pVer.greaterThan(&csVer) {
			ret = append(ret, &api.VersionInfo{Release: common.KubernetesRelease1Dot6, Version: patchVersion})
		}
		// add next release
		ret = append(ret, &api.VersionInfo{Release: common.KubernetesRelease1Dot7, Version: common.KubeReleaseToVersion[common.KubernetesRelease1Dot7]})
	case common.KubernetesRelease1Dot7:
		// check for patch upgrade
		patchVersion := common.KubeReleaseToVersion[common.KubernetesRelease1Dot7]
		if err := pVer.parse(patchVersion); err != nil {
			return ret, err
		}
		if pVer.greaterThan(&csVer) {
			ret = append(ret, &api.VersionInfo{Release: common.KubernetesRelease1Dot7, Version: patchVersion})
		}
	}
	return ret, nil
}

func dcosInfo(csInfo *api.VersionInfo) ([]*api.OrchestratorInfo, error) {
	orchs := []*api.OrchestratorInfo{}
	if len(csInfo.Release) == 0 {
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
		ver, ok := common.DCOSReleaseToVersion[csInfo.Release]
		if !ok {
			return nil, fmt.Errorf("DCOS release %s is not supported", csInfo.Release)
		}
		orchs = append(orchs,
			&api.OrchestratorInfo{
				Orchestrator: api.DCOS,
				VersionInfo: api.VersionInfo{
					Release: csInfo.Release,
					Version: ver,
				},
				Default: csInfo.Release == common.DCOSDefaultRelease,
			})
	}
	return orchs, nil
}

func swarmInfo(csInfo *api.VersionInfo) ([]*api.OrchestratorInfo, error) {
	return []*api.OrchestratorInfo{
		{
			Orchestrator: api.Swarm,
			VersionInfo: api.VersionInfo{
				Version: SwarmVersion,
			},
			DockerComposeVersion: SwarmDockerComposeVersion,
		},
	}, nil
}

func dockerceInfo(csInfo *api.VersionInfo) ([]*api.OrchestratorInfo, error) {
	return []*api.OrchestratorInfo{
		{
			Orchestrator: api.SwarmMode,
			VersionInfo: api.VersionInfo{
				Version: DockerCEVersion,
			},
			DockerComposeVersion: DockerCEDockerComposeVersion,
		},
	}, nil
}
