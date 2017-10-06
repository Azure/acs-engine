package kubernetesupgrade

// UpgradeWorkFlow outlines various individual high level steps
// that need to be run (one or more times) in the upgrade workflow.
type UpgradeWorkFlow interface {
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
	DeleteNode(*string) error

	// CreateNode creates a new master/agent node with the targeted version of Kubernetes
	CreateNode(string, int) error

	// Validate will verify the that master/agent node has been upgraded as expected.
	Validate() error
}
