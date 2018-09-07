package kubernetesupgrade

import (
	"os"
	"testing"

	"github.com/Azure/acs-engine/pkg/acsengine"
	"github.com/Azure/acs-engine/pkg/api"
	"github.com/Azure/acs-engine/pkg/armhelpers"
	"github.com/Azure/acs-engine/pkg/i18n"
	. "github.com/Azure/acs-engine/pkg/test"
	. "github.com/onsi/gomega"

	"fmt"

	. "github.com/onsi/ginkgo"
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

	It("Should return error message when failing to list VMs during upgrade operation", func() {
		cs := acsengine.CreateMockContainerService("testcluster", "1.6.9", 1, 1, false)

		cs.Properties.OrchestratorProfile.OrchestratorVersion = "1.7.14"

		uc := UpgradeCluster{
			Translator: &i18n.Translator{},
			Logger:     log.NewEntry(log.New()),
		}

		mockClient := armhelpers.MockACSEngineClient{}
		mockClient.FailListVirtualMachines = true
		uc.Client = &mockClient

		subID, _ := uuid.FromString("DEC923E3-1EF1-4745-9516-37906D56DEC4")

		err := uc.UpgradeCluster(subID, "kubeConfig", "TestRg", cs, "12345678", []string{"agentpool1"}, TestACSEngineVersion)
		Expect(err).NotTo(BeNil())
		Expect(err.Error()).To(Equal("Error while querying ARM for resources: ListVirtualMachines failed"))

		// Clean up
		os.RemoveAll("./translations")
	})

	It("Should return error message when failing to delete VMs during upgrade operation", func() {
		cs := acsengine.CreateMockContainerService("testcluster", "1.6.9", 1, 1, false)

		cs.Properties.OrchestratorProfile.OrchestratorVersion = "1.7.16"
		uc := UpgradeCluster{
			Translator: &i18n.Translator{},
			Logger:     log.NewEntry(log.New()),
		}

		mockClient := armhelpers.MockACSEngineClient{}
		mockClient.FailDeleteVirtualMachine = true
		uc.Client = &mockClient

		subID, _ := uuid.FromString("DEC923E3-1EF1-4745-9516-37906D56DEC4")

		err := uc.UpgradeCluster(subID, "kubeConfig", "TestRg", cs, "12345678", []string{"agentpool1"}, TestACSEngineVersion)
		Expect(err).NotTo(BeNil())
		Expect(err.Error()).To(Equal("DeleteVirtualMachine failed"))
	})

	It("Should return error message when failing to deploy template during upgrade operation", func() {
		cs := acsengine.CreateMockContainerService("testcluster", "1.7.16", 1, 1, false)
		cs.Properties.OrchestratorProfile.OrchestratorVersion = "1.7.16"
		uc := UpgradeCluster{
			Translator: &i18n.Translator{},
			Logger:     log.NewEntry(log.New()),
		}

		mockClient := armhelpers.MockACSEngineClient{}
		mockClient.FailDeployTemplate = true
		uc.Client = &mockClient

		subID, _ := uuid.FromString("DEC923E3-1EF1-4745-9516-37906D56DEC4")

		err := uc.UpgradeCluster(subID, "kubeConfig", "TestRg", cs, "12345678", []string{"agentpool1"}, TestACSEngineVersion)
		Expect(err).NotTo(BeNil())
		Expect(err.Error()).To(Equal("DeployTemplate failed"))
	})

	It("Should return error message when failing to get a virtual machine during upgrade operation", func() {
		cs := acsengine.CreateMockContainerService("testcluster", "1.6.9", 1, 6, false)
		cs.Properties.OrchestratorProfile.OrchestratorVersion = "1.7.16"
		uc := UpgradeCluster{
			Translator: &i18n.Translator{},
			Logger:     log.NewEntry(log.New()),
		}

		mockClient := armhelpers.MockACSEngineClient{}
		mockClient.FailGetVirtualMachine = true
		uc.Client = &mockClient

		subID, _ := uuid.FromString("DEC923E3-1EF1-4745-9516-37906D56DEC4")

		err := uc.UpgradeCluster(subID, "kubeConfig", "TestRg", cs, "12345678", []string{"agentpool1"}, TestACSEngineVersion)
		Expect(err).NotTo(BeNil())
		Expect(err.Error()).To(Equal("GetVirtualMachine failed"))
	})

	It("Should return error message when failing to get storage client during upgrade operation", func() {
		cs := acsengine.CreateMockContainerService("testcluster", "1.6.9", 5, 1, false)
		cs.Properties.OrchestratorProfile.OrchestratorVersion = "1.7.16"
		uc := UpgradeCluster{
			Translator: &i18n.Translator{},
			Logger:     log.NewEntry(log.New()),
		}

		mockClient := armhelpers.MockACSEngineClient{}
		mockClient.FailGetStorageClient = true
		uc.Client = &mockClient

		subID, _ := uuid.FromString("DEC923E3-1EF1-4745-9516-37906D56DEC4")

		err := uc.UpgradeCluster(subID, "kubeConfig", "TestRg", cs, "12345678", []string{"agentpool1"}, TestACSEngineVersion)
		Expect(err).NotTo(BeNil())
		Expect(err.Error()).To(Equal("GetStorageClient failed"))
	})

	It("Should return error message when failing to delete network interface during upgrade operation", func() {
		cs := acsengine.CreateMockContainerService("testcluster", "1.6.9", 3, 2, false)
		cs.Properties.OrchestratorProfile.OrchestratorVersion = "1.7.16"
		uc := UpgradeCluster{
			Translator: &i18n.Translator{},
			Logger:     log.NewEntry(log.New()),
		}

		mockClient := armhelpers.MockACSEngineClient{}
		mockClient.FailDeleteNetworkInterface = true
		uc.Client = &mockClient

		subID, _ := uuid.FromString("DEC923E3-1EF1-4745-9516-37906D56DEC4")

		err := uc.UpgradeCluster(subID, "kubeConfig", "TestRg", cs, "12345678", []string{"agentpool1"}, TestACSEngineVersion)
		Expect(err).NotTo(BeNil())
		Expect(err.Error()).To(Equal("DeleteNetworkInterface failed"))
	})

	It("Should return error message when failing on ClusterPreflightCheck operation", func() {
		cs := acsengine.CreateMockContainerService("testcluster", "1.6.9", 3, 3, false)
		cs.Properties.OrchestratorProfile.OrchestratorVersion = "1.8.15"
		uc := UpgradeCluster{
			Translator: &i18n.Translator{},
			Logger:     log.NewEntry(log.New()),
		}

		mockClient := armhelpers.MockACSEngineClient{}
		uc.Client = &mockClient

		subID, _ := uuid.FromString("DEC923E3-1EF1-4745-9516-37906D56DEC4")

		err := uc.UpgradeCluster(subID, "kubeConfig", "TestRg", cs, "12345678", []string{"agentpool1"}, TestACSEngineVersion)
		Expect(err).NotTo(BeNil())
		fmt.Print("GOT :   ", err.Error())
		Expect(err.Error()).To(ContainSubstring("Error while querying ARM for resources: Kubernetes:1.6.9 cannot be upgraded to 1.8.15"))
	})

	It("Should return error message when failing to delete role assignment during upgrade operation", func() {
		cs := acsengine.CreateMockContainerService("testcluster", "1.6.9", 3, 2, false)
		cs.Properties.OrchestratorProfile.OrchestratorVersion = "1.7.16"
		cs.Properties.OrchestratorProfile.KubernetesConfig = &api.KubernetesConfig{}
		cs.Properties.OrchestratorProfile.KubernetesConfig.UseManagedIdentity = true
		uc := UpgradeCluster{
			Translator: &i18n.Translator{},
			Logger:     log.NewEntry(log.New()),
		}

		mockClient := armhelpers.MockACSEngineClient{}
		mockClient.FailDeleteRoleAssignment = true
		mockClient.ShouldSupportVMIdentity = true
		uc.Client = &mockClient

		subID, _ := uuid.FromString("DEC923E3-1EF1-4745-9516-37906D56DEC4")

		err := uc.UpgradeCluster(subID, "kubeConfig", "TestRg", cs, "12345678", []string{"agentpool1"}, TestACSEngineVersion)
		Expect(err).NotTo(BeNil())
		Expect(err.Error()).To(Equal("DeleteRoleAssignmentByID failed"))
	})

	It("Should not fail if no managed identity is returned by azure during upgrade operation", func() {
		cs := acsengine.CreateMockContainerService("testcluster", "1.6.9", 3, 2, false)
		cs.Properties.OrchestratorProfile.OrchestratorVersion = "1.7.16"
		cs.Properties.OrchestratorProfile.KubernetesConfig = &api.KubernetesConfig{}
		cs.Properties.OrchestratorProfile.KubernetesConfig.UseManagedIdentity = true
		uc := UpgradeCluster{
			Translator: &i18n.Translator{},
			Logger:     log.NewEntry(log.New()),
		}

		mockClient := armhelpers.MockACSEngineClient{}
		uc.Client = &mockClient

		subID, _ := uuid.FromString("DEC923E3-1EF1-4745-9516-37906D56DEC4")

		err := uc.UpgradeCluster(subID, "kubeConfig", "TestRg", cs, "12345678", []string{"agentpool1"}, TestACSEngineVersion)
		Expect(err).To(BeNil())
	})
})
