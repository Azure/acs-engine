package cmd

import (
	"fmt"

	"github.com/Azure/acs-engine/pkg/helpers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	flag "github.com/spf13/pflag"
)

var _ = Describe("the version command", func() {
	It("should create a version command", func() {
		output := newVersionCmd()

		flag := &flag.Flag{
			Name:      "output",
			Shorthand: "o",
			Usage:     "Output format to use: [human json]",
			DefValue:  "human",
		}
		flag.Value.Set("human")

		Expect(output.Use).Should(Equal(versionName))
		Expect(output.Short).Should(Equal(versionShortDescription))
		Expect(output.Long).Should(Equal(versionLongDescription))
		Expect(output.Flag("output")).Should(Equal(flag))
	})

	It("should print a json version of ACS-Engine", func() {
		output := getVersion("json")

		expectedOutput, _ := helpers.JSONMarshalIndent(version, "", "  ", false)

		Expect(output).Should(Equal(string(expectedOutput)))
	})
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
