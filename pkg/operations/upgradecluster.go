package operations

import (
	"github.com/Azure/acs-engine/pkg/api/vlabs"
	"github.com/satori/uuid"
)

// ClusterTopology contains resources of the cluster the upgrade operation
// is targeting
type ClusterTopology struct {
	*vlabs.ContainerService
}

// UpgradeCluster upgrades a cluster with Orchestrator version X
// (or X.X or X.X.X) to version y (or Y.Y or X.X.X). RIght now
// upgrades are supported for Kubernetes cluster only.
type UpgradeCluster struct{}

// UpgradeCluster runs the workflow to upgrade a Kubernetes cluster.
// UpgradeContainerService contains target state of the cluster that
// the operation will drive towards.
func (uc *UpgradeCluster) UpgradeCluster(subscription uuid.UUID, rg string,
	cs *vlabs.ContainerService, ucs *vlabs.UpgradeContainerService) {
	// 	ListClusterResources(subscription uuid.UUID, rg string) error
}

// UpgradeWorkFlow outlines various individual high level steps
// that need to be run (one or more times) in the upgrade workflow.
type UpgradeWorkFlow interface {
	ClusterPreflightCheck()

	// upgrade masters
	// upgrade agent nodes
	RunUpgrade() error

	Validate() error
}

// UpgradeNode drives work flow of deleting and replacing a master or agent node to a
// specified target version of Kubernetes
type UpgradeNode interface {
	// DeleteNode takes state/resources of the master/agent node from ListNodeResources
	// backs up/preserves state as needed by a specific version of Kubernetes and then deletes
	// the node
	DeleteNode() error

	// CreateNode creates a new master/agent node with the targeted version of Kubernetes
	CreateNode() error

	// Validate will verify the that master/agent node has been upgraded as expected.
	Validate() error
}
