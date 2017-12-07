package armhelpers

import (
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/arm/authorization"
	"github.com/Azure/azure-sdk-for-go/arm/compute"
	"github.com/Azure/azure-sdk-for-go/arm/disk"
	"github.com/Azure/azure-sdk-for-go/arm/graphrbac"
	"github.com/Azure/azure-sdk-for-go/arm/resources/resources"
	"github.com/Azure/go-autorest/autorest"
	log "github.com/sirupsen/logrus"
	"k8s.io/client-go/pkg/api/v1"
)

//MockACSEngineClient is an implementation of ACSEngineClient where all requests error out
type MockACSEngineClient struct {
	FailDeployTemplate              bool
	FailEnsureResourceGroup         bool
	FailListVirtualMachines         bool
	FailListVirtualMachineScaleSets bool
	FailGetVirtualMachine           bool
	FailDeleteVirtualMachine        bool
	FailGetStorageClient            bool
	FailDeleteNetworkInterface      bool
	FailGetKubernetesClient         bool
	FailListProviders               bool
	MockKubernetesClient            *MockKubernetesClient
}

//MockStorageClient mock implementation of StorageClient
type MockStorageClient struct{}

//MockKubernetesClient mock implementation of KubernetesClient
type MockKubernetesClient struct {
	FailListPods          bool
	FailGetNode           bool
	UpdateNodeFunc        func(*v1.Node) (*v1.Node, error)
	FailUpdateNode        bool
	FailSupportEviction   bool
	FailDeletePod         bool
	FailEvictPod          bool
	FailWaitForDelete     bool
	ShouldSupportEviction bool
	PodsList              *v1.PodList
}

//ListPods returns all Pods running on the passed in node
func (mkc *MockKubernetesClient) ListPods(node *v1.Node) (*v1.PodList, error) {
	if mkc.FailListPods {
		return nil, fmt.Errorf("ListPods failed")
	}
	if mkc.PodsList != nil {
		return mkc.PodsList, nil
	}
	return &v1.PodList{}, nil
}

//GetNode returns details about node with passed in name
func (mkc *MockKubernetesClient) GetNode(name string) (*v1.Node, error) {
	if mkc.FailGetNode {
		return nil, fmt.Errorf("GetNode failed")
	}
	node := &v1.Node{}
	node.Status.Conditions = append(node.Status.Conditions, v1.NodeCondition{Type: v1.NodeReady, Status: v1.ConditionTrue})
	return node, nil
}

//UpdateNode updates the node in the api server with the passed in info
func (mkc *MockKubernetesClient) UpdateNode(node *v1.Node) (*v1.Node, error) {
	if mkc.UpdateNodeFunc != nil {
		return mkc.UpdateNodeFunc(node)
	}
	if mkc.FailUpdateNode {
		return nil, fmt.Errorf("UpdateNode failed")
	}
	return node, nil
}

//SupportEviction queries the api server to discover if it supports eviction, and returns supported type if it is supported
func (mkc *MockKubernetesClient) SupportEviction() (string, error) {
	if mkc.FailSupportEviction {
		return "", fmt.Errorf("SupportEviction failed")
	}
	if mkc.ShouldSupportEviction {
		return "version", nil
	}
	return "", nil
}

//DeletePod deletes the passed in pod
func (mkc *MockKubernetesClient) DeletePod(pod *v1.Pod) error {
	if mkc.FailDeletePod {
		return fmt.Errorf("DeletePod failed")
	}
	return nil
}

//EvictPod evicts the passed in pod using the passed in api version
func (mkc *MockKubernetesClient) EvictPod(pod *v1.Pod, policyGroupVersion string) error {
	if mkc.FailEvictPod {
		return fmt.Errorf("EvictPod failed")
	}
	return nil
}

//WaitForDelete waits until all pods are deleted. Returns all pods not deleted and an error on failure
func (mkc *MockKubernetesClient) WaitForDelete(logger *log.Entry, pods []v1.Pod, usingEviction bool) ([]v1.Pod, error) {
	if mkc.FailWaitForDelete {
		return nil, fmt.Errorf("WaitForDelete failed")
	}
	return []v1.Pod{}, nil
}

//DeleteBlob mock
func (msc *MockStorageClient) DeleteBlob(container, blob string) error {
	return nil
}

//AddAcceptLanguages mock
func (mc *MockACSEngineClient) AddAcceptLanguages(languages []string) {
	return
}

//DeployTemplate mock
func (mc *MockACSEngineClient) DeployTemplate(resourceGroup, name string, template, parameters map[string]interface{}, cancel <-chan struct{}) (*resources.DeploymentExtended, error) {
	if mc.FailDeployTemplate {
		return nil, fmt.Errorf("DeployTemplate failed")
	}

	return nil, nil
}

//EnsureResourceGroup mock
func (mc *MockACSEngineClient) EnsureResourceGroup(resourceGroup, location string, managedBy *string) (*resources.Group, error) {
	if mc.FailEnsureResourceGroup {
		return nil, fmt.Errorf("EnsureResourceGroup failed")
	}

	return nil, nil
}

//ListVirtualMachines mock
func (mc *MockACSEngineClient) ListVirtualMachines(resourceGroup string) (compute.VirtualMachineListResult, error) {
	if mc.FailListVirtualMachines {
		return compute.VirtualMachineListResult{}, fmt.Errorf("ListVirtualMachines failed")
	}

	vm1Name := "k8s-agentpool1-12345678-0"

	creationSourceString := "creationSource"
	orchestratorString := "orchestrator"
	resourceNameSuffixString := "resourceNameSuffix"
	poolnameString := "poolName"

	creationSource := "acsengine-k8s-agentpool1-12345678-0"
	orchestrator := "Kubernetes:1.5.8"
	resourceNameSuffix := "12345678"
	poolname := "agentpool1"

	tags := map[string]*string{
		creationSourceString:     &creationSource,
		orchestratorString:       &orchestrator,
		resourceNameSuffixString: &resourceNameSuffix,
		poolnameString:           &poolname,
	}

	vm1 := compute.VirtualMachine{
		Name: &vm1Name,
		Tags: &tags,
		VirtualMachineProperties: &compute.VirtualMachineProperties{
			StorageProfile: &compute.StorageProfile{
				OsDisk: &compute.OSDisk{
					Vhd: &compute.VirtualHardDisk{
						URI: &validOSDiskResourceName},
				},
			},
			NetworkProfile: &compute.NetworkProfile{
				NetworkInterfaces: &[]compute.NetworkInterfaceReference{
					{
						ID: &validNicResourceName,
					},
				},
			},
		},
	}

	vmr := compute.VirtualMachineListResult{}
	vmr.Value = &[]compute.VirtualMachine{vm1}

	return vmr, nil
}

//ListVirtualMachineScaleSets mock
func (mc *MockACSEngineClient) ListVirtualMachineScaleSets(resourceGroup string) (compute.VirtualMachineScaleSetListResult, error) {
	if mc.FailListVirtualMachineScaleSets {
		return compute.VirtualMachineScaleSetListResult{}, fmt.Errorf("ListVirtualMachines failed")
	}

	return compute.VirtualMachineScaleSetListResult{}, nil
}

//GetVirtualMachine mock
func (mc *MockACSEngineClient) GetVirtualMachine(resourceGroup, name string) (compute.VirtualMachine, error) {
	if mc.FailGetVirtualMachine {
		return compute.VirtualMachine{}, fmt.Errorf("GetVirtualMachine failed")
	}

	vm1Name := "k8s-agentpool1-12345678-0"

	creationSourceString := "creationSource"
	orchestratorString := "orchestrator"
	resourceNameSuffixString := "resourceNameSuffix"
	poolnameString := "poolName"

	creationSource := "acsengine-k8s-agentpool1-12345678-0"
	orchestrator := "Kubernetes:1.5.8"
	resourceNameSuffix := "12345678"
	poolname := "agentpool1"

	tags := map[string]*string{
		creationSourceString:     &creationSource,
		orchestratorString:       &orchestrator,
		resourceNameSuffixString: &resourceNameSuffix,
		poolnameString:           &poolname,
	}

	return compute.VirtualMachine{
		Name: &vm1Name,
		Tags: &tags,
		VirtualMachineProperties: &compute.VirtualMachineProperties{
			StorageProfile: &compute.StorageProfile{
				OsDisk: &compute.OSDisk{
					Vhd: &compute.VirtualHardDisk{
						URI: &validOSDiskResourceName},
				},
			},
			NetworkProfile: &compute.NetworkProfile{
				NetworkInterfaces: &[]compute.NetworkInterfaceReference{
					{
						ID: &validNicResourceName,
					},
				},
			},
		},
	}, nil
}

//DeleteVirtualMachine mock
func (mc *MockACSEngineClient) DeleteVirtualMachine(resourceGroup, name string, cancel <-chan struct{}) (<-chan compute.OperationStatusResponse, <-chan error) {
	if mc.FailDeleteVirtualMachine {
		errChan := make(chan error)
		respChan := make(chan compute.OperationStatusResponse)
		go func() {
			defer func() {
				close(errChan)
			}()
			defer func() {
				close(respChan)
			}()
			errChan <- fmt.Errorf("DeleteVirtualMachine failed")
		}()
		return respChan, errChan
	}

	errChan := make(chan error)
	respChan := make(chan compute.OperationStatusResponse)
	go func() {
		defer func() {
			close(errChan)
		}()
		defer func() {
			close(respChan)
		}()
		errChan <- nil
		respChan <- compute.OperationStatusResponse{}
	}()
	return respChan, errChan
}

//GetStorageClient mock
func (mc *MockACSEngineClient) GetStorageClient(resourceGroup, accountName string) (ACSStorageClient, error) {
	if mc.FailGetStorageClient {
		return nil, fmt.Errorf("GetStorageClient failed")
	}

	return &MockStorageClient{}, nil
}

//DeleteNetworkInterface mock
func (mc *MockACSEngineClient) DeleteNetworkInterface(resourceGroup, nicName string, cancel <-chan struct{}) (<-chan autorest.Response, <-chan error) {
	if mc.FailDeleteNetworkInterface {
		errChan := make(chan error)
		respChan := make(chan autorest.Response)
		go func() {
			defer func() {
				close(errChan)
			}()
			defer func() {
				close(respChan)
			}()
			errChan <- fmt.Errorf("DeleteNetworkInterface failed")
		}()
		return respChan, errChan
	}

	errChan := make(chan error)
	respChan := make(chan autorest.Response)
	go func() {
		defer func() {
			close(errChan)
		}()
		defer func() {
			close(respChan)
		}()
		errChan <- nil
		respChan <- autorest.Response{}
	}()
	return respChan, errChan
}

var validOSDiskResourceName = "https://00k71r4u927seqiagnt0.blob.core.windows.net/osdisk/k8s-agentpool1-12345678-0-osdisk.vhd"
var validNicResourceName = "/subscriptions/DEC923E3-1EF1-4745-9516-37906D56DEC4/resourceGroups/acsK8sTest/providers/Microsoft.Network/networkInterfaces/k8s-agent-12345678-nic-0"

// Active Directory
// Mocks

// Graph Mocks

// CreateGraphApplication creates an application via the graphrbac client
func (mc *MockACSEngineClient) CreateGraphApplication(applicationCreateParameters graphrbac.ApplicationCreateParameters) (graphrbac.Application, error) {
	return graphrbac.Application{}, nil
}

// CreateGraphPrincipal creates a service principal via the graphrbac client
func (mc *MockACSEngineClient) CreateGraphPrincipal(servicePrincipalCreateParameters graphrbac.ServicePrincipalCreateParameters) (graphrbac.ServicePrincipal, error) {
	return graphrbac.ServicePrincipal{}, nil
}

// CreateApp is a simpler method for creating an application
func (mc *MockACSEngineClient) CreateApp(applicationName, applicationURL string) (applicationID, servicePrincipalObjectID, secret string, err error) {
	return "app-id", "client-id", "client-secret", nil
}

// RBAC Mocks

// CreateRoleAssignment creates a role assignment via the authorization client
func (mc *MockACSEngineClient) CreateRoleAssignment(scope string, roleAssignmentName string, parameters authorization.RoleAssignmentCreateParameters) (authorization.RoleAssignment, error) {
	return authorization.RoleAssignment{}, nil
}

// CreateRoleAssignmentSimple is a wrapper around RoleAssignmentsClient.Create
func (mc *MockACSEngineClient) CreateRoleAssignmentSimple(applicationID, roleID string) error {
	return nil
}

// DeleteManagedDisk is a wrapper around disksClient.Delete
func (mc *MockACSEngineClient) DeleteManagedDisk(resourceGroupName string, diskName string, cancel <-chan struct{}) (<-chan disk.OperationStatusResponse, <-chan error) {
	errChan := make(chan error)
	respChan := make(chan disk.OperationStatusResponse)
	go func() {
		defer func() {
			close(errChan)
		}()
		defer func() {
			close(respChan)
		}()
		errChan <- nil
		respChan <- disk.OperationStatusResponse{}
	}()
	return respChan, errChan
}

// ListManagedDisksByResourceGroup is a wrapper around disksClient.ListManagedDisksByResourceGroup
func (mc *MockACSEngineClient) ListManagedDisksByResourceGroup(resourceGroupName string) (result disk.ListType, err error) {
	return disk.ListType{}, nil
}

//GetKubernetesClient mock
func (mc *MockACSEngineClient) GetKubernetesClient(masterURL, kubeConfig string, interval, timeout time.Duration) (KubernetesClient, error) {
	if mc.FailGetKubernetesClient {
		return nil, fmt.Errorf("GetKubernetesClient failed")
	}

	if mc.MockKubernetesClient == nil {
		mc.MockKubernetesClient = &MockKubernetesClient{}
	}
	return mc.MockKubernetesClient, nil
}

// ListProviders mock
func (mc *MockACSEngineClient) ListProviders() (resources.ProviderListResult, error) {
	if mc.FailListProviders {
		return resources.ProviderListResult{}, fmt.Errorf("ListProviders failed")
	}

	return resources.ProviderListResult{}, nil
}
