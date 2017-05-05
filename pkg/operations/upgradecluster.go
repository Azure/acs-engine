package operations

import "github.com/Azure/acs-engine/pkg/api/vlabs"

// UpgradeCluster upgrades a cluster with Orchestrator version X
// (or X.X or X.X.X) to version y (or Y.Y or X.X.X). RIght now
// upgrades are supported for Kubernetes cluster only.
type UpgradeCluster struct{}

// UpgradeCluster runs the workflow to upgrade a Kubernetes
// cluster
func (uc *UpgradeCluster) UpgradeCluster(rg string, cs *vlabs.ContainerService) {
	// a.	Create API Model
	// b.	Resource group name
	// c.	Input Upgrade API (for now it will only container target orchestrator version) 1.6.2
	// d.	Subscription Id

}
