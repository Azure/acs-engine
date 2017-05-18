package operations

// Compiler to verify QueueMessageProcessor implements OperationsProcessor
var _ UpgradeNode = &UpgradeAgentNode{}

// UpgradeAgentNode upgrades a Kubernetes 1.5.3 agent node to 1.6.2
type UpgradeAgentNode struct {
}

// DeleteNode takes state/resources of the master/agent node from ListNodeResources
// backs up/preserves state as needed by a specific version of Kubernetes and then deletes
// the node
func (kmn *UpgradeAgentNode) DeleteNode(*string) error {
	return nil
}

// CreateNode creates a new master/agent node with the targeted version of Kubernetes
func (kmn *UpgradeAgentNode) CreateNode(int) error {
	return nil
}

// Validate will verify the that master/agent node has been upgraded as expected.
func (kmn *UpgradeAgentNode) Validate() error {
	return nil
}
