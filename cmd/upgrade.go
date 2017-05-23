package cmd

import (
	"os"
	"path"

	"github.com/Azure/acs-engine/pkg/api"
	"github.com/Azure/acs-engine/pkg/armhelpers"
	"github.com/Azure/acs-engine/pkg/operations"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
)

const (
	upgradeName             = "upgrade"
	upgradeShortDescription = "upgrades an existing Kubernetes cluster"
	upgradeLongDescription  = "upgrades an existing Kubernetes cluster, first replacing masters, then nodes"
)

type upgradeCmd struct {
	authArgs

	// user input
	resourceGroupName   string
	deploymentDirectory string
	upgradeModelFile    string
	containerService    *api.ContainerService
	apiVersion          string

	// derived
	upgradeContainerService *api.UpgradeContainerService
	upgradeAPIVersion       string
	client                  armhelpers.ACSEngineClient
}

// NewUpgradeCmd run a command to upgrade a Kubernetes cluster
func newUpgradeCmd() *cobra.Command {
	uc := upgradeCmd{}

	upgradeCmd := &cobra.Command{
		Use:   upgradeName,
		Short: upgradeShortDescription,
		Long:  upgradeLongDescription,
		RunE: func(cmd *cobra.Command, args []string) error {
			return uc.run(cmd, args)
		},
	}

	f := upgradeCmd.Flags()
	f.StringVar(&uc.resourceGroupName, "resource-group", "", "the resource group where the cluster is deployed")
	f.StringVar(&uc.deploymentDirectory, "deployment-dir", "", "the location of the output from `generate`")
	f.StringVar(&uc.upgradeModelFile, "upgrademodel-file", "", "file path to upgrade API model")
	addAuthFlags(&uc.authArgs, f)

	return upgradeCmd
}

func (uc *upgradeCmd) validate(cmd *cobra.Command, args []string) {
	log.Infoln("validating...")

	var err error

	if uc.resourceGroupName == "" {
		cmd.Usage()
		log.Fatal("--resource-group must be specified")
	}

	// TODO(colemick): add in the cmd annotation to help enable autocompletion
	if uc.upgradeModelFile == "" {
		cmd.Usage()
		log.Fatal("--upgrademodel-file must be specified")
	}

	if uc.client, err = uc.authArgs.getClient(); err != nil {
		log.Error("Failed to get client:", err)
	}

	if uc.deploymentDirectory == "" {
		cmd.Usage()
		log.Fatal("--deployment-dir must be specified")
	}

	// load apimodel from the deployment directory
	apiModelPath := path.Join(uc.deploymentDirectory, "apimodel.json")

	if _, err := os.Stat(apiModelPath); os.IsNotExist(err) {
		log.Fatalf("specified api model does not exist (%s)", apiModelPath)
	}

	uc.containerService, uc.apiVersion, err = api.LoadContainerServiceFromFile(apiModelPath)
	if err != nil {
		log.Fatalf("error parsing the api model: %s", err.Error())
	}

	if _, err := os.Stat(uc.upgradeModelFile); os.IsNotExist(err) {
		log.Fatalf("specified upgrade model file does not exist (%s)", uc.upgradeModelFile)
	}

	uc.upgradeContainerService, uc.upgradeAPIVersion, err = api.LoadUpgradeContainerServiceFromFile(uc.upgradeModelFile)
	if err != nil {
		log.Fatalf("error parsing the upgrade api model: %s", err.Error())
	}

	uc.client, err = uc.authArgs.getClient()
	if err != nil {
		log.Fatalf("failed to get client") // TODO: cleanup
	}

	// TODO: Validate that downgrade is not allowed
	// TODO: Validate noop case and return early
}

func (uc *upgradeCmd) run(cmd *cobra.Command, args []string) error {
	uc.validate(cmd, args)

	upgradeCluster := operations.UpgradeCluster{
		Client: uc.client,
	}

	if err := upgradeCluster.UpgradeCluster(uc.authArgs.SubscriptionID, uc.resourceGroupName,
		uc.containerService, uc.upgradeContainerService); err != nil {
		log.Fatalf("Error upgrading cluster: %s \n", err.Error())
	}

	return nil
}
