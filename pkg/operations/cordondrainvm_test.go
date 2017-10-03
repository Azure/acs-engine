package operations

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	log "github.com/sirupsen/logrus"
)

var _ = Describe("Safely Drain node operation tests", func() {
	It("Should return error messages for invalid kube config", func() {
		err := SafelyDrainNode(log.NewEntry(log.New()), "http://bad.com/", "bad", "node")
		Expect(err).Should(HaveOccurred())
	})
})
