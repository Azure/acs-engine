package api

import (
	"github.com/Azure/acs-engine/pkg/api/v20170930"
	"github.com/Azure/acs-engine/pkg/api/vlabs"
)

///////////////////////////////////////////////////////////
// The converter exposes functions to convert the top level
// UpgradeContainerService API model
//
// All other functions are internal helper functions used
// for converting.
///////////////////////////////////////////////////////////

// ConvertVLabsUpgradeContainerService converts a vlabs UpgradeContainerService to an unversioned UpgradeContainerService
func ConvertVLabsUpgradeContainerService(vlabUCS *vlabs.UpgradeContainerService) *UpgradeContainerService {
	ucs := &UpgradeContainerService{}
	convertVLabsOrchestratorProfile((*vlabs.OrchestratorProfile)(vlabUCS), (*OrchestratorProfile)(ucs))
	return ucs
}

// ConvertV20170930UpgradeContainerService converts a v20170930 UpgradeContainerService to an unversioned UpgradeContainerService
func ConvertV20170930UpgradeContainerService(vUCS *v20170930.UpgradeContainerService) *UpgradeContainerService {
	ucs := &UpgradeContainerService{}
	switch vUCS.OrchestratorType {
	case v20170930.Kubernetes:
		ucs.OrchestratorType = Kubernetes
	case v20170930.DCOS:
		ucs.OrchestratorType = DCOS
	case v20170930.Swarm:
		ucs.OrchestratorType = Swarm
	case v20170930.DockerCE:
		ucs.OrchestratorType = SwarmMode
	}
	ucs.OrchestratorRelease = vUCS.OrchestratorRelease
	ucs.OrchestratorVersion = vUCS.OrchestratorVersion
	return ucs
}
