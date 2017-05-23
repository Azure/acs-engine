package api

import "github.com/Azure/acs-engine/pkg/api/vlabs"

///////////////////////////////////////////////////////////
// The converter exposes functions to convert the top level
// UpgradeContainerService API model
//
// All other functions are internal helper functions used
// for converting.
///////////////////////////////////////////////////////////

// ConvertVLabsUpgradeContainerService converts a vlabs UpgradeContainerService to an unversioned UpgradeContainerService
func ConvertVLabsUpgradeContainerService(vlabs *vlabs.UpgradeContainerService) *UpgradeContainerService {
	ucs := &UpgradeContainerService{}
	ucs.OrchestratorProfile = &OrchestratorProfile{}
	convertVLabsOrchestratorProfile(vlabs.OrchestratorProfile, ucs.OrchestratorProfile)
	return ucs
}

func convertVLabsUpgradeOrchestratorProfile(vlabscs *vlabs.OrchestratorProfile, api *OrchestratorProfile) {
	api.OrchestratorType = OrchestratorType(vlabscs.OrchestratorType)
	if api.OrchestratorType == Kubernetes {
		switch vlabscs.OrchestratorVersion {
		case vlabs.Kubernetes162:
			api.OrchestratorVersion = Kubernetes162
		case vlabs.Kubernetes160:
			api.OrchestratorVersion = Kubernetes160
		default:
			api.OrchestratorVersion = KubernetesLatest
		}
	}
}
