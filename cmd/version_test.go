package cmd

import (
	"fmt"

	"github.com/Azure/acs-engine/pkg/helpers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("the version command", func() {
	It("should print a humanized version of ACS-Engine", func() {
		output := getVersion("human")

		expectedOutput := fmt.Sprintf("Version: %s\nGitCommit: %s\nGitTreeState: %s",
			BuildTag,
			BuildSHA,
			GitTreeState)

		Expect(output).Should(Equal(expectedOutput))
	})

	It("should print a json version of ACS-Engine", func() {
		output := getVersion("json")

		expectedOutput, _ := helpers.JSONMarshalIndent(version, "", "  ", false)

		Expect(output).Should(Equal(string(expectedOutput)))
	})
})
