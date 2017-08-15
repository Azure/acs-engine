package cmd

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("The info command", func() {
	It("should fail on unsupported orchestrator", func() {
		infoCmd := &infoCmd{
			orchestrator: "unsupported",
		}

		err := infoCmd.run(nil, nil)
		Expect(err).NotTo(BeNil())
		Expect(err.Error()).To(Equal("Unsupported orchestrator 'unsupported'"))
	})

	It("should fail on unprovided orchestrator", func() {
		infoCmd := &infoCmd{
			release: "1.1",
		}

		err := infoCmd.run(nil, nil)
		Expect(err).NotTo(BeNil())
		Expect(err.Error()).To(Equal("Must specify orchestrator for release '1.1'"))
	})

	It("should fail on unsupported release", func() {
		infoCmd := &infoCmd{
			orchestrator: "kubernetes",
			release:      "1.1",
		}

		err := infoCmd.run(nil, nil)
		Expect(err).NotTo(BeNil())
		Expect(err.Error()).To(Equal("Kubernetes release 1.1 is not supported"))
	})

	It("should succeed", func() {
		infoCmd := &infoCmd{
			orchestrator: "kubernetes",
			release:      "1.7",
		}

		err := infoCmd.run(nil, nil)
		Expect(err).To(BeNil())
	})
})
