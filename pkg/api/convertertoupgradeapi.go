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
func ConvertVLabsUpgradeContainerService(vlabUCS *vlabs.UpgradeContainerService) *UpgradeContainerService {
	ucs := &UpgradeContainerService{}
	convertVLabsOrchestratorProfile((*vlabs.OrchestratorProfile)(vlabUCS), (*OrchestratorProfile)(ucs))
	return ucs
}
