package operations

import (
	"testing"

	"github.com/Azure/acs-engine/pkg/api"
	"github.com/Azure/acs-engine/pkg/armhelpers"
	. "github.com/onsi/gomega"
	"github.com/satori/go.uuid"
)

func Test_UpgradeCluster_FailListVirtualMachines(t *testing.T) {
	RegisterTestingT(t)

	cs := api.ContainerService{}
	ucs := api.UpgradeContainerService{}

	uc := UpgradeCluster{}

	mockClient := armhelpers.MockACSEngineClient{}
	mockClient.FailListVirtualMachines = true
	uc.Client = &mockClient

	subID, _ := uuid.FromString("DEC923E3-1EF1-4745-9516-37906D56DEC4")

	err := uc.UpgradeCluster(subID, "TestRg", &cs, &ucs, "12345678")

	Expect(err.Error()).To(Equal("Error while querying ARM for resources: ListVirtualMachines failed"))
}
