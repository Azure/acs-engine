package cmd

import (
	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	BuildSHA  = "unset"
	BuildTime = "unset"
)

func NewVersionCmd() *cobra.Command {
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
