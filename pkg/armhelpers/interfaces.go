package armhelpers

import (
	"time"

	"github.com/Azure/azure-sdk-for-go/arm/authorization"
	"github.com/Azure/azure-sdk-for-go/arm/compute"
	"github.com/Azure/azure-sdk-for-go/arm/disk"
	"github.com/Azure/azure-sdk-for-go/arm/graphrbac"
	"github.com/Azure/azure-sdk-for-go/arm/resources/resources"
	"github.com/Azure/go-autorest/autorest"
	log "github.com/sirupsen/logrus"
	"k8s.io/api/core/v1"
)

// ACSEngineClient is the interface used to talk to an Azure environment.
// This interface exposes just the subset of Azure APIs and clients needed for
// ACS-Engine.
type ACSEngineClient interface {

	//AddAcceptLanguages sets the list of languages to accept on this request
	AddAcceptLanguages(languages []string)
	//
	// RESOURCES

	// DeployTemplate can deploy a template into Azure ARM
	DeployTemplate(resourceGroup, name string, template, parameters map[string]interface{}, cancel <-chan struct{}) (*resources.DeploymentExtended, error)

	// EnsureResourceGroup ensures the specified resource group exists in the specified location
	EnsureResourceGroup(resourceGroup, location string, managedBy *string) (*resources.Group, error)

	//
	// COMPUTE

	// List lists VM resources
	ListVirtualMachines(resourceGroup string) (compute.VirtualMachineListResult, error)

	// GetVirtualMachine retrieves the specified virtual machine.
	GetVirtualMachine(resourceGroup, name string) (compute.VirtualMachine, error)

	// DeleteVirtualMachine deletes the specified virtual machine.
	DeleteVirtualMachine(resourceGroup, name string, cancel <-chan struct{}) (<-chan compute.OperationStatusResponse, <-chan error)

	// ListVirtualMachineScaleSets lists the vmss resources in the resource group
	ListVirtualMachineScaleSets(resourceGroup string) (compute.VirtualMachineScaleSetListResult, error)

	// ListVirtualMachineScaleSetVMs lists the virtual machines contained in a vmss
	ListVirtualMachineScaleSetVMs(resourceGroup, virtualMachineScaleSet string) (compute.VirtualMachineScaleSetVMListResult, error)

	// DeleteVirtualMachineScaleSetVM deletes a VM in a VMSS
	DeleteVirtualMachineScaleSetVM(resourceGroup, virtualMachineScaleSet, instanceID string, cancel <-chan struct{}) (<-chan compute.OperationStatusResponse, <-chan error)

	// SetVirtualMachineScaleSetCapacity sets the VMSS capacity
	SetVirtualMachineScaleSetCapacity(resourceGroup, virtualMachineScaleSet string, sku compute.Sku, location string, cancel <-chan struct{}) (<-chan compute.VirtualMachineScaleSet, <-chan error)

	//
	// STORAGE

	// GetStorageClient uses SRP to retrieve keys, and then an authenticated client for talking to the specified storage
	// account.
	GetStorageClient(resourceGroup, accountName string) (ACSStorageClient, error)

	//
	// NETWORK

	// DeleteNetworkInterface deletes the specified network interface.
	DeleteNetworkInterface(resourceGroup, nicName string, cancel <-chan struct{}) (<-chan autorest.Response, <-chan error)

	//
	// GRAPH

	// CreateGraphAppliction creates an application via the graphrbac client
	CreateGraphApplication(applicationCreateParameters graphrbac.ApplicationCreateParameters) (graphrbac.Application, error)

	// CreateGraphPrincipal creates a service principal via the graphrbac client
	CreateGraphPrincipal(servicePrincipalCreateParameters graphrbac.ServicePrincipalCreateParameters) (graphrbac.ServicePrincipal, error)
	CreateApp(applicationName, applicationURL string, replyURLs *[]string, requiredResourceAccess *[]graphrbac.RequiredResourceAccess) (applicationID, servicePrincipalObjectID, secret string, err error)

	// RBAC
	CreateRoleAssignment(scope string, roleAssignmentName string, parameters authorization.RoleAssignmentCreateParameters) (authorization.RoleAssignment, error)
	CreateRoleAssignmentSimple(applicationID, roleID string) error
	DeleteRoleAssignmentByID(roleAssignmentNameID string) (authorization.RoleAssignment, error)
	ListRoleAssignmentsForPrincipal(scope string, principalID string) (authorization.RoleAssignmentListResult, error)

	// MANAGED DISKS
	DeleteManagedDisk(resourceGroupName string, diskName string, cancel <-chan struct{}) (<-chan disk.OperationStatusResponse, <-chan error)
	ListManagedDisksByResourceGroup(resourceGroupName string) (result disk.ListType, err error)

	GetKubernetesClient(masterURL, kubeConfig string, interval, timeout time.Duration) (KubernetesClient, error)

	ListProviders() (resources.ProviderListResult, error)

	// DEPLOYMENTS

	// ListDeploymentOperations gets all deployments operations for a deployment.
	ListDeploymentOperations(resourceGroupName string, deploymentName string, top *int32) (result resources.DeploymentOperationsListResult, err error)

	// ListDeploymentOperationsNextResults retrieves the next set of results, if any.
	ListDeploymentOperationsNextResults(lastResults resources.DeploymentOperationsListResult) (result resources.DeploymentOperationsListResult, err error)
}

// ACSStorageClient interface models the azure storage client
type ACSStorageClient interface {
	// DeleteBlob deletes the specified blob in the specified container.
	DeleteBlob(container, blob string) error
}

// KubernetesClient interface models client for interacting with kubernetes api server
type KubernetesClient interface {
	//ListPods returns all Pods running on the passed in node
	ListPods(node *v1.Node) (*v1.PodList, error)
	//GetNode returns details about node with passed in name
	GetNode(name string) (*v1.Node, error)
	//UpdateNode updates the node in the api server with the passed in info
	UpdateNode(node *v1.Node) (*v1.Node, error)
	//DeleteNode deregisters node in the api server
	DeleteNode(name string) error
	//SupportEviction queries the api server to discover if it supports eviction, and returns supported type if it is supported
	SupportEviction() (string, error)
	//DeletePod deletes the passed in pod
	DeletePod(pod *v1.Pod) error
	//EvictPod evicts the passed in pod using the passed in api version
	EvictPod(pod *v1.Pod, policyGroupVersion string) error
	//WaitForDelete waits until all pods are deleted. Returns all pods not deleted and an error on failure
	WaitForDelete(logger *log.Entry, pods []v1.Pod, usingEviction bool) ([]v1.Pod, error)
}
