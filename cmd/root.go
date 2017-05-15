package cmd

import (
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
)

const (
	rootName             = "acs-engine"
	rootShortDescription = "ACS-Engine deploys and manages container orchestrators in Azure"
	rootLongDescription  = "ACS-Engine deploys and manages Kubernetes, Swarm Mode, and DC/OS clusters in Azure"
)

var (
	debug bool
)

// NewRootCmd returns the root command for ACS-Engine.
func NewRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   rootName,
		Short: rootShortDescription,
		Long:  rootLongDescription,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if debug {
				log.SetLevel(log.DebugLevel)
			}
		},
	}

	p := rootCmd.PersistentFlags()
	p.BoolVar(&debug, "debug", false, "enable verbose debug logs")

	rootCmd.AddCommand(newVersionCmd())
	rootCmd.AddCommand(newGenerateCmd())

	if val := os.Getenv("ACSENGINE_EXPERIMENTAL_FEATURES"); val == "1" {
		rootCmd.AddCommand(newUpgradeCmd())
	}

	return rootCmd
}
