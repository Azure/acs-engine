package kubernetesupgrade

import (
	"os"
	"testing"

	"fmt"

	"github.com/Azure/acs-engine/pkg/api"
	"github.com/Azure/acs-engine/pkg/armhelpers"
	"github.com/Azure/acs-engine/pkg/i18n"
	. "github.com/Azure/acs-engine/pkg/test"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
)

const TestACSEngineVersion = "1.0.0"

func TestUpgradeCluster(t *testing.T) {
	RunSpecsWithReporters(t, "kubernetesupgrade", "Server Suite")
}

var _ = Describe("Upgrade Kubernetes cluster tests", func() {
	AfterEach(func() {
		// delete temp template directory
		os.RemoveAll("_output")
	})

	It("Should succeed when cluster VMs are missing expected tags during upgrade operation", func() {
		mockK8sVersion, upgradeK8sVersion := "1.8.15", "1.9.11"
		cs := api.CreateMockContainerService("testcluster", mockK8sVersion, 1, 6, false)
		cs.Properties.OrchestratorProfile.OrchestratorVersion = upgradeK8sVersion
		uc := UpgradeCluster{
			Translator: &i18n.Translator{},
			Logger:     log.NewEntry(log.New()),
		}

		c := armhelpers.MockKubernetesClient{KubernetesVersion: mockK8sVersion}
		mockClient := armhelpers.MockACSEngineClient{MockKubernetesClient: &c}
		mockClient.FailListVirtualMachinesTags = true
		uc.Client = &mockClient

		subID, _ := uuid.FromString("DEC923E3-1EF1-4745-9516-37906D56DEC4")

		err := uc.UpgradeCluster(subID, &mockClient, "kubeConfig", "TestRg", cs, "12345678", []string{"agentpool1"}, TestACSEngineVersion)
		Expect(err).To(BeNil())
		Expect(uc.ClusterTopology.AgentPools).NotTo(BeEmpty())

		// Clean up
		os.RemoveAll("./translations")
	})

	It("Should return error message when failing to list VMs during upgrade operation", func() {
		mockK8sVersion, upgradeK8sVersion := "1.8.15", "1.9.11"
		cs := api.CreateMockContainerService("testcluster", mockK8sVersion, 1, 6, false)
		cs.Properties.OrchestratorProfile.OrchestratorVersion = upgradeK8sVersion
		uc := UpgradeCluster{
			Translator: &i18n.Translator{},
			Logger:     log.NewEntry(log.New()),
		}

		c := armhelpers.MockKubernetesClient{KubernetesVersion: mockK8sVersion}
		mockClient := armhelpers.MockACSEngineClient{MockKubernetesClient: &c}
		mockClient.FailListVirtualMachines = true
		uc.Client = &mockClient

		subID, _ := uuid.FromString("DEC923E3-1EF1-4745-9516-37906D56DEC4")

		err := uc.UpgradeCluster(subID, nil, "kubeConfig", "TestRg", cs, "12345678", []string{"agentpool1"}, TestACSEngineVersion)
		Expect(err).NotTo(BeNil())
		Expect(err.Error()).To(Equal("Error while querying ARM for resources: ListVirtualMachines failed"))

		// Clean up
		os.RemoveAll("./translations")
	})

	It("Should return error message when failing to delete VMs during upgrade operation", func() {
		mockK8sVersion, upgradeK8sVersion := "1.8.15", "1.9.11"
		cs := api.CreateMockContainerService("testcluster", mockK8sVersion, 1, 1, false)
		cs.Properties.OrchestratorProfile.OrchestratorVersion = upgradeK8sVersion
		uc := UpgradeCluster{
			Translator: &i18n.Translator{},
			Logger:     log.NewEntry(log.New()),
		}

		c := armhelpers.MockKubernetesClient{KubernetesVersion: mockK8sVersion}
		mockClient := armhelpers.MockACSEngineClient{MockKubernetesClient: &c}
		mockClient.FailDeleteVirtualMachine = true
		uc.Client = &mockClient

		subID, _ := uuid.FromString("DEC923E3-1EF1-4745-9516-37906D56DEC4")

		err := uc.UpgradeCluster(subID, nil, "kubeConfig", "TestRg", cs, "12345678", []string{"agentpool1"}, TestACSEngineVersion)
		Expect(err).NotTo(BeNil())
		Expect(err.Error()).To(Equal("DeleteVirtualMachine failed"))
	})

	It("Should return error message when failing to deploy template during upgrade operation", func() {
		mockK8sVersion, upgradeK8sVersion := "1.8.15", "1.9.11"
		cs := api.CreateMockContainerService("testcluster", mockK8sVersion, 1, 1, false)
		cs.Properties.OrchestratorProfile.OrchestratorVersion = upgradeK8sVersion
		uc := UpgradeCluster{
			Translator: &i18n.Translator{},
			Logger:     log.NewEntry(log.New()),
		}

		c := armhelpers.MockKubernetesClient{KubernetesVersion: mockK8sVersion}
		mockClient := armhelpers.MockACSEngineClient{MockKubernetesClient: &c}
		mockClient.FailDeployTemplate = true
		uc.Client = &mockClient

		subID, _ := uuid.FromString("DEC923E3-1EF1-4745-9516-37906D56DEC4")

		err := uc.UpgradeCluster(subID, nil, "kubeConfig", "TestRg", cs, "12345678", []string{"agentpool1"}, TestACSEngineVersion)
		Expect(err).NotTo(BeNil())
		Expect(err.Error()).To(Equal("DeployTemplate failed"))
	})

	It("Should return error message when failing to get a virtual machine during upgrade operation", func() {
		mockK8sVersion, upgradeK8sVersion := "1.8.15", "1.9.11"
		cs := api.CreateMockContainerService("testcluster", mockK8sVersion, 1, 6, false)
		cs.Properties.OrchestratorProfile.OrchestratorVersion = upgradeK8sVersion
		uc := UpgradeCluster{
			Translator: &i18n.Translator{},
			Logger:     log.NewEntry(log.New()),
		}

		c := armhelpers.MockKubernetesClient{KubernetesVersion: mockK8sVersion}
		mockClient := armhelpers.MockACSEngineClient{MockKubernetesClient: &c}
		mockClient.FailGetVirtualMachine = true
		uc.Client = &mockClient

		subID, _ := uuid.FromString("DEC923E3-1EF1-4745-9516-37906D56DEC4")

		err := uc.UpgradeCluster(subID, nil, "kubeConfig", "TestRg", cs, "12345678", []string{"agentpool1"}, TestACSEngineVersion)
		Expect(err).NotTo(BeNil())
		Expect(err.Error()).To(Equal("GetVirtualMachine failed"))
	})

	It("Should return error message when failing to get storage client during upgrade operation", func() {
		mockK8sVersion, upgradeK8sVersion := "1.8.15", "1.9.11"
		cs := api.CreateMockContainerService("testcluster", mockK8sVersion, 5, 1, false)
		cs.Properties.OrchestratorProfile.OrchestratorVersion = upgradeK8sVersion
		uc := UpgradeCluster{
			Translator: &i18n.Translator{},
			Logger:     log.NewEntry(log.New()),
		}

		c := armhelpers.MockKubernetesClient{KubernetesVersion: mockK8sVersion}
		mockClient := armhelpers.MockACSEngineClient{MockKubernetesClient: &c}
		mockClient.FailGetStorageClient = true
		uc.Client = &mockClient

		subID, _ := uuid.FromString("DEC923E3-1EF1-4745-9516-37906D56DEC4")

		err := uc.UpgradeCluster(subID, nil, "kubeConfig", "TestRg", cs, "12345678", []string{"agentpool1"}, TestACSEngineVersion)
		Expect(err).NotTo(BeNil())
		Expect(err.Error()).To(Equal("GetStorageClient failed"))
	})

	It("Should return error message when failing to delete network interface during upgrade operation", func() {
		mockK8sVersion, upgradeK8sVersion := "1.8.15", "1.9.11"
		cs := api.CreateMockContainerService("testcluster", mockK8sVersion, 3, 2, false)
		cs.Properties.OrchestratorProfile.OrchestratorVersion = upgradeK8sVersion
		uc := UpgradeCluster{
			Translator: &i18n.Translator{},
			Logger:     log.NewEntry(log.New()),
		}

		c := armhelpers.MockKubernetesClient{KubernetesVersion: mockK8sVersion}
		mockClient := armhelpers.MockACSEngineClient{MockKubernetesClient: &c}
		mockClient.FailDeleteNetworkInterface = true
		uc.Client = &mockClient

		subID, _ := uuid.FromString("DEC923E3-1EF1-4745-9516-37906D56DEC4")

		err := uc.UpgradeCluster(subID, nil, "kubeConfig", "TestRg", cs, "12345678", []string{"agentpool1"}, TestACSEngineVersion)
		Expect(err).NotTo(BeNil())
		Expect(err.Error()).To(Equal("DeleteNetworkInterface failed"))
	})

	It("Should return error message when failing on ClusterPreflightCheck operation", func() {
		mockK8sVersion, upgradeK8sVersion := "1.6.9", "1.9.11"
		cs := api.CreateMockContainerService("testcluster", mockK8sVersion, 3, 2, false)
		cs.Properties.OrchestratorProfile.OrchestratorVersion = upgradeK8sVersion
		uc := UpgradeCluster{
			Translator: &i18n.Translator{},
			Logger:     log.NewEntry(log.New()),
		}

		c := armhelpers.MockKubernetesClient{KubernetesVersion: mockK8sVersion}
		mockClient := armhelpers.MockACSEngineClient{MockKubernetesClient: &c}
		uc.Client = &mockClient

		subID, _ := uuid.FromString("DEC923E3-1EF1-4745-9516-37906D56DEC4")

		err := uc.UpgradeCluster(subID, nil, "kubeConfig", "TestRg", cs, "12345678", []string{"agentpool1"}, TestACSEngineVersion)
		Expect(err).NotTo(BeNil())
		fmt.Print("GOT :   ", err.Error())
		errStr := fmt.Sprintf("Error while querying ARM for resources: Kubernetes:%s cannot be upgraded to %s",
			mockK8sVersion, upgradeK8sVersion)
		Expect(err.Error()).To(ContainSubstring(errStr))
	})

	It("Should return error message when failing to delete role assignment during upgrade operation", func() {
		mockK8sVersion, upgradeK8sVersion := "1.8.15", "1.9.11"
		cs := api.CreateMockContainerService("testcluster", mockK8sVersion, 3, 2, false)
		cs.Properties.OrchestratorProfile.OrchestratorVersion = upgradeK8sVersion
		cs.Properties.OrchestratorProfile.KubernetesConfig = &api.KubernetesConfig{}
		cs.Properties.OrchestratorProfile.KubernetesConfig.UseManagedIdentity = true
		uc := UpgradeCluster{
			Translator: &i18n.Translator{},
			Logger:     log.NewEntry(log.New()),
		}

		c := armhelpers.MockKubernetesClient{KubernetesVersion: mockK8sVersion}
		mockClient := armhelpers.MockACSEngineClient{MockKubernetesClient: &c}
		mockClient.FailDeleteRoleAssignment = true
		mockClient.ShouldSupportVMIdentity = true
		uc.Client = &mockClient

		subID, _ := uuid.FromString("DEC923E3-1EF1-4745-9516-37906D56DEC4")

		err := uc.UpgradeCluster(subID, nil, "kubeConfig", "TestRg", cs, "12345678", []string{"agentpool1"}, TestACSEngineVersion)
		Expect(err).NotTo(BeNil())
		Expect(err.Error()).To(Equal("DeleteRoleAssignmentByID failed"))
	})

	It("Should not fail if no managed identity is returned by azure during upgrade operation", func() {
		mockK8sVersion, upgradeK8sVersion := "1.8.15", "1.9.11"
		cs := api.CreateMockContainerService("testcluster", mockK8sVersion, 3, 2, false)
		cs.Properties.OrchestratorProfile.OrchestratorVersion = upgradeK8sVersion
		cs.Properties.OrchestratorProfile.KubernetesConfig = &api.KubernetesConfig{}
		cs.Properties.OrchestratorProfile.KubernetesConfig.UseManagedIdentity = true
		uc := UpgradeCluster{
			Translator: &i18n.Translator{},
			Logger:     log.NewEntry(log.New()),
		}

		c := armhelpers.MockKubernetesClient{KubernetesVersion: mockK8sVersion}
		mockClient := armhelpers.MockACSEngineClient{MockKubernetesClient: &c}
		uc.Client = &mockClient

		subID, _ := uuid.FromString("DEC923E3-1EF1-4745-9516-37906D56DEC4")

		err := uc.UpgradeCluster(subID, nil, "kubeConfig", "TestRg", cs, "12345678", []string{"agentpool1"}, TestACSEngineVersion)
		Expect(err).To(BeNil())
	})
})
