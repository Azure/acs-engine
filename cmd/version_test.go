package cmd

import (
	"encoding/json"
	"fmt"

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

		expectedOutput, _ := json.MarshalIndent(version, "", "  ")

		Expect(output).Should(Equal(string(expectedOutput)))
	})
})
