package cmd

import (
	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	// BuildSHA holds the git commit SHA at `make build` time.
	BuildSHA = "unset"
)

func newVersionCmd() *cobra.Command {
	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Print the version of ACS-Engine",
		Long:  "Print the version of ACS-Engine",

		Run: func(cmd *cobra.Command, args []string) {
			log.Infof("ACS-Engine Version: %s", BuildSHA)
		},
	}
	return versionCmd
}
