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
})
