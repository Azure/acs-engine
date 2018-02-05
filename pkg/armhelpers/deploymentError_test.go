package armhelpers

import (
	"testing"

	. "github.com/Azure/acs-engine/pkg/test"
	. "github.com/onsi/gomega"

	. "github.com/onsi/ginkgo"
	log "github.com/sirupsen/logrus"
)

func TestUpgradeCluster(t *testing.T) {
	RunSpecsWithReporters(t, "templatedeployment", "Server Suite")
}

var _ = Describe("Template deployment tests", func() {

	It("Should return InternalOperationError error code", func() {
		mockClient := &MockACSEngineClient{}
		mockClient.FailDeployTemplate = true
		logger := log.NewEntry(log.New())

		err := DeployTemplateSync(mockClient, logger, "rg1", "agentvm", map[string]interface{}{}, map[string]interface{}{})
		Expect(err).NotTo(BeNil())
		armerr, ok := err.(*ArmError)
		Expect(ok).To(BeTrue())
		Expect(armerr.Code).To(Equal("DeploymentFailed"))
	})

	It("Should return QuotaExceeded error code, specified in details", func() {
		mockClient := &MockACSEngineClient{}
		mockClient.FailDeployTemplateQuota = true
		logger := log.NewEntry(log.New())

		err := DeployTemplateSync(mockClient, logger, "rg1", "agentvm", map[string]interface{}{}, map[string]interface{}{})
		Expect(err).NotTo(BeNil())
		armerr, ok := err.(*ArmError)
		Expect(ok).To(BeTrue())
		Expect(armerr.Code).To(Equal("InvalidTemplateDeployment"))
		Expect(len(armerr.Details)).To(Equal(1))
		Expect(armerr.Details[0].Code).To(Equal("QuotaExceeded"))
	})

	It("Should return Conflict error code, specified in details", func() {
		mockClient := &MockACSEngineClient{}
		mockClient.FailDeployTemplateConflict = true
		logger := log.NewEntry(log.New())

		err := DeployTemplateSync(mockClient, logger, "rg1", "agentvm", map[string]interface{}{}, map[string]interface{}{})
		Expect(err).NotTo(BeNil())
		armerr, ok := err.(*ArmError)
		Expect(ok).To(BeTrue())
		Expect(armerr.Code).To(Equal("DeploymentFailed"))
		Expect(len(armerr.Details)).To(Equal(1))
		Expect(armerr.Details[0].Code).To(Equal("Conflict"))
	})
})
