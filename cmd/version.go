package cmd

import (
	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	// BuildSHA holds the git commit SHA at `make build` time.
	BuildSHA = "unset"

	// BuildTime holds the `date` at `make build` time.
	BuildTime = "unset"
)

func newVersionCmd() *cobra.Command {
	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Print the version of ACS-Engine",
		Long:  "Print the version of ACS-Engine",

		Run: func(cmd *cobra.Command, args []string) {
			log.Infof("ACS-Engine Version: %s (%s)", BuildSHA, BuildTime)
		},
	}
	return versionCmd
}
