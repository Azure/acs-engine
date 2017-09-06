package api

import (
	"github.com/Azure/acs-engine/pkg/api/v20170930"
	"github.com/Azure/acs-engine/pkg/api/vlabs"
)

///////////////////////////////////////////////////////////
// The converter exposes functions to convert the top level
// ContainerService resource
//
// All other functions are internal helper functions used
// for converting.
///////////////////////////////////////////////////////////

// ConvertUpgradeContainerServiceToV20170930 converts an unversioned UpgradeContainerService to a v20170930 UpgradeContainerService
func ConvertUpgradeContainerServiceToV20170930(api *UpgradeContainerService) *v20170930.UpgradeContainerService {
	vProfile := &v20170930.UpgradeContainerService{}
	switch api.OrchestratorType {
	case Kubernetes:
		vProfile.OrchestratorType = v20170930.Kubernetes
	case DCOS:
		vProfile.OrchestratorType = v20170930.DCOS
	case Swarm:
		vProfile.OrchestratorType = v20170930.Swarm
	case SwarmMode:
		vProfile.OrchestratorType = v20170930.DockerCE
	}
	vProfile.OrchestratorVersion = api.OrchestratorVersion
	vProfile.OrchestratorRelease = api.OrchestratorRelease
	return vProfile
}

// ConvertUpgradeContainerServiceToVLabs converts an unversioned UpgradeContainerService to a vlabs UpgradeContainerService
func ConvertUpgradeContainerServiceToVLabs(ucs *UpgradeContainerService) *vlabs.UpgradeContainerService {
	vlabsProfile := &vlabs.UpgradeContainerService{}
	switch ucs.OrchestratorType {
	case Kubernetes:
		vlabsProfile.OrchestratorType = vlabs.Kubernetes
	case DCOS:
		vlabsProfile.OrchestratorType = vlabs.DCOS
	case Swarm:
		vlabsProfile.OrchestratorType = vlabs.Swarm
	case SwarmMode:
		vlabsProfile.OrchestratorType = vlabs.SwarmMode
	}
	vlabsProfile.OrchestratorVersion = ucs.OrchestratorVersion
	vlabsProfile.OrchestratorRelease = ucs.OrchestratorRelease
	return vlabsProfile
}
