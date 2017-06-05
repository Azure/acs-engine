package kubernetesupgrade

import (
	"testing"

	"github.com/Azure/acs-engine/pkg/api"
	"github.com/Azure/acs-engine/pkg/armhelpers"
	"github.com/onsi/ginkgo/reporters"
	. "github.com/onsi/gomega"
	"github.com/satori/go.uuid"

	. "github.com/onsi/ginkgo"
)

func TestUpgradeCluster(t *testing.T) {
	RegisterFailHandler(Fail)
	junitReporter := reporters.NewJUnitReporter("junit.xml")
	RunSpecsWithDefaultAndCustomReporters(t, "Server Suite", []Reporter{junitReporter})
}

var _ = Describe("Upgrade Kubernetes cluster tests", func() {
	It("Should return error message when failing to list VMs during upgrade operation", func() {
		cs := api.ContainerService{}
		ucs := api.UpgradeContainerService{}

		uc := UpgradeCluster{}

		mockClient := armhelpers.MockACSEngineClient{}
		mockClient.FailListVirtualMachines = true
		uc.Client = &mockClient

		subID, _ := uuid.FromString("DEC923E3-1EF1-4745-9516-37906D56DEC4")

		err := uc.UpgradeCluster(subID, "TestRg", &cs, &ucs, "12345678")

		Expect(err.Error()).To(Equal("Error while querying ARM for resources: ListVirtualMachines failed"))
	})

	It("Should return error message when failing to detete VMs during upgrade operation", func() {
		cs := createContainerService("testcluster", 1, 1)

		ucs := api.UpgradeContainerService{}
		ucs.OrchestratorProfile = &api.OrchestratorProfile{}
		ucs.OrchestratorProfile.OrchestratorType = api.Kubernetes
		ucs.OrchestratorProfile.OrchestratorVersion = api.Kubernetes162

		uc := UpgradeCluster{}

		mockClient := armhelpers.MockACSEngineClient{}
		mockClient.FailDeleteVirtualMachine = true
		uc.Client = &mockClient

		subID, _ := uuid.FromString("DEC923E3-1EF1-4745-9516-37906D56DEC4")

		err := uc.UpgradeCluster(subID, "TestRg", cs, &ucs, "12345678")

		Expect(err.Error()).To(Equal("DeleteVirtualMachine failed"))
	})
})

func createContainerService(containerServiceName string, masterCount int, agentCount int) *api.ContainerService {
	cs := api.ContainerService{}
	cs.ID = uuid.NewV4().String()
	cs.Location = "eastus"
	cs.Name = containerServiceName

	cs.Properties = &api.Properties{}

	cs.Properties.MasterProfile = &api.MasterProfile{}
	cs.Properties.MasterProfile.Count = masterCount
	cs.Properties.MasterProfile.DNSPrefix = "testmaster"
	cs.Properties.MasterProfile.VMSize = "Standard_D2_v2"

	cs.Properties.AgentPoolProfiles = []*api.AgentPoolProfile{}
	agentPool := &api.AgentPoolProfile{}
	agentPool.Count = agentCount
	agentPool.Name = "agentpool1"
	agentPool.VMSize = "Standard_D2_v2"
	agentPool.OSType = "Linux"
	agentPool.AvailabilityProfile = "AvailabilitySet"
	agentPool.StorageProfile = "StorageAccount"

	cs.Properties.AgentPoolProfiles = append(cs.Properties.AgentPoolProfiles, agentPool)

	cs.Properties.LinuxProfile = &api.LinuxProfile{
		AdminUsername: "azureuser",
		SSH: struct {
			PublicKeys []struct {
				KeyData string `json:"keyData"`
			} `json:"publicKeys"`
		}{
			PublicKeys: []struct {
				KeyData string `json:"keyData"`
			}{

				{
					KeyData: "test",
				},
			},
		},
	}

	cs.Properties.ServicePrincipalProfile = &api.ServicePrincipalProfile{}
	cs.Properties.ServicePrincipalProfile.ClientID = "DEC923E3-1EF1-4745-9516-37906D56DEC4"
	cs.Properties.ServicePrincipalProfile.Secret = "DEC923E3-1EF1-4745-9516-37906D56DEC4"

	cs.Properties.OrchestratorProfile = &api.OrchestratorProfile{}
	cs.Properties.OrchestratorProfile.OrchestratorType = api.Kubernetes
	cs.Properties.OrchestratorProfile.OrchestratorVersion = api.Kubernetes153

	return &cs
}

// FailDeployTemplate         bool
// FailGetVirtualMachine      boolmake
// FailGetStorageClient       bool
// FailDeleteNetworkInterface bool
