package api

import "github.com/Azure/acs-engine/pkg/api/vlabs"

///////////////////////////////////////////////////////////
// The converter exposes functions to convert the top level
// ContainerService resource
//
// All other functions are internal helper functions used
// for converting.
///////////////////////////////////////////////////////////

// ConvertUpgradeContainerServiceToVLabs converts an unversioned ContainerService to a vlabs ContainerService
func ConvertUpgradeContainerServiceToVLabs(api *UpgradeContainerService) *vlabs.UpgradeContainerService {
	vlabsUCS := &vlabs.UpgradeContainerService{}
	vlabsUCS.OrchestratorProfile = &vlabs.OrchestratorProfile{}
	convertUpgradeOrchestratorProfileToVLabs(api.OrchestratorProfile, vlabsUCS.OrchestratorProfile)
	return vlabsUCS
}

func convertUpgradeOrchestratorProfileToVLabs(api *OrchestratorProfile, o *vlabs.OrchestratorProfile) {
	o.OrchestratorType = vlabs.OrchestratorType(api.OrchestratorType)

	if api.OrchestratorVersion != "" {
		o.OrchestratorVersion = vlabs.OrchestratorVersion(api.OrchestratorVersion)
	}
}
