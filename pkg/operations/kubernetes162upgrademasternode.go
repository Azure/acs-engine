package operations

// Compiler to verify QueueMessageProcessor implements OperationsProcessor
var _ UpgradeNode = &UpgradeMasterNode{}

// UpgradeMasterNode upgrades a Kubernetes 1.5.3 master node to 1.6.2
type UpgradeMasterNode struct {
}

// DeleteNode takes state/resources of the master/agent node from ListNodeResources
// backs up/preserves state as needed by a specific version of Kubernetes and then deletes
// the node
func (kmn *UpgradeMasterNode) DeleteNode() error {
	return nil
}

// CreateNode creates a new master/agent node with the targeted version of Kubernetes
func (kmn *UpgradeMasterNode) CreateNode() error {
	return nil
}

// Validate will verify the that master/agent node has been upgraded as expected.
func (kmn *UpgradeMasterNode) Validate() error {
	return nil
}
