package cmd

import (
	"encoding/json"
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
)

var (
	// BuildTag holds the `git tag` if this is a tagged build/release
	BuildTag = "canary"

	// BuildSHA holds the git commit SHA at `make build` time.
	BuildSHA = "unset"

	// GitTreeState is the state of the git tree, either clean or dirty
	GitTreeState = "unset"

	outputFormatOptions = []string{"human", "json"}
	outputFormat        string
	version             versionInfo
)

type versionInfo struct {
	GitTag       string
	GitCommit    string
	GitTreeState string
}

func init() {
	version = versionInfo{
		GitTag:       BuildTag,
		GitCommit:    BuildSHA,
		GitTreeState: GitTreeState,
	}
}

func getHumanVersion() string {
	r := fmt.Sprintf("Version: %s\nGitCommit: %s\nGitTreeState: %s",
		version.GitTag,
		version.GitCommit,
		version.GitTreeState)

	return r
}

func getJSONVersion() string {
	jsonVersion, _ := json.MarshalIndent(version, "", "  ")
	return string(jsonVersion)
}

func getVersion(outputType string) string {
	var output string

	if outputType == "human" {
		output = getHumanVersion()
	} else if outputType == "json" {
		output = getJSONVersion()
	} else {
		log.Fatalf("unsupported output format: %s\n", outputFormat)
	}

	return output
}

func newVersionCmd() *cobra.Command {
	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Print the version of ACS-Engine",
		Long:  "Print the version of ACS-Engine",

		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(getVersion(outputFormat))
		},
	}

	versionCmdDescription := fmt.Sprintf("Output format to use: %s", outputFormatOptions)

	versionCmd.Flags().StringVarP(&outputFormat, "output", "o", "human", versionCmdDescription)

	return versionCmd
}
