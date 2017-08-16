package cmd

import (
	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	// BuildSHA holds the git commit SHA at `make build` time.
	BuildSHA = "unset"

	// BuildTag holds the `git tag` if this is a tagged build/release
	BuildTag = "unset"
)

func newVersionCmd() *cobra.Command {
	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Print the version of ACS-Engine",
		Long:  "Print the version of ACS-Engine",

		Run: func(cmd *cobra.Command, args []string) {
			log.Infof("ACS-Engine Version: %s (%s)", BuildTag, BuildSHA)
		},
	}
	return versionCmd
}
