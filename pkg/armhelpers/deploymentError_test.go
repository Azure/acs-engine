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
		deplErr, ok := err.(*DeploymentError)
		Expect(ok).To(BeTrue())
		Expect(deplErr.TopError).NotTo(BeNil())
		Expect(deplErr.TopError.Error()).To(Equal("DeployTemplate failed"))
		Expect(deplErr.ProvisioningState).To(Equal(""))
		Expect(deplErr.StatusCode).To(Equal(0))
		Expect(len(deplErr.OperationsLists)).To(Equal(0))
	})

	It("Should return QuotaExceeded error code, specified in details", func() {
		mockClient := &MockACSEngineClient{}
		mockClient.FailDeployTemplateQuota = true
		logger := log.NewEntry(log.New())

		err := DeployTemplateSync(mockClient, logger, "rg1", "agentvm", map[string]interface{}{}, map[string]interface{}{})
		Expect(err).NotTo(BeNil())
		deplErr, ok := err.(*DeploymentError)
		Expect(ok).To(BeTrue())
		Expect(deplErr.TopError).NotTo(BeNil())
		Expect(deplErr.ProvisioningState).To(Equal(""))
		Expect(deplErr.StatusCode).To(Equal(400))
		Expect(string(deplErr.Response)).To(ContainSubstring("\"code\":\"QuotaExceeded\""))
		Expect(len(deplErr.OperationsLists)).To(Equal(0))
	})

	It("Should return Conflict error code, specified in details", func() {
		mockClient := &MockACSEngineClient{}
		mockClient.FailDeployTemplateConflict = true
		logger := log.NewEntry(log.New())

		err := DeployTemplateSync(mockClient, logger, "rg1", "agentvm", map[string]interface{}{}, map[string]interface{}{})
		Expect(err).NotTo(BeNil())
		deplErr, ok := err.(*DeploymentError)
		Expect(ok).To(BeTrue())
		Expect(deplErr.TopError).NotTo(BeNil())
		Expect(deplErr.ProvisioningState).To(Equal(""))
		Expect(deplErr.StatusCode).To(Equal(200))
		Expect(string(deplErr.Response)).To(ContainSubstring("\"code\":\"Conflict\""))
		Expect(len(deplErr.OperationsLists)).To(Equal(0))
	})
})
