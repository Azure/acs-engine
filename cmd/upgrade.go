package cmd

import (
	"os"
	"path"

	"github.com/Azure/acs-engine/pkg/api"

	"github.com/Azure/acs-engine/pkg/operations"
	armhelpers "github.com/Azure/acs-engine/pkg/operations/armhelpers"
	"github.com/Azure/go-autorest/autorest/adal"
	"github.com/Azure/go-autorest/autorest/azure"
	log "github.com/Sirupsen/logrus"
	"github.com/satori/go.uuid"
	"github.com/spf13/cobra"
)

const (
	upgradeName             = "upgrade"
	upgradeShortDescription = "upgrades an existing Kubernetes cluster"
	upgradeLongDescription  = "upgrades an existing Kubernetes cluster, first replacing masters, then nodes"
)

type upgradeCmd struct {
	// user input
	authMethod          string
	clientSecret        string
	resourceGroupName   string
	deploymentDirectory string
	rawClientID         string
	rawSubscriptionID   string
	rawAzureEnvironment string
	upgradeModelFile    string

	// parsed
	clientID         uuid.UUID
	subscriptionID   uuid.UUID
	azureEnvironment azure.Environment

	// derived
	tenantID              string
	servicePrincipalToken *adal.ServicePrincipalToken
	containerService      *api.ContainerService
	apiVersion            string

	upgradeContainerService *api.UpgradeContainerService
	upgradeAPIVersion       string
}

// NewUpgradeCmd run a command to upgrade a Kubernetes cluster
func NewUpgradeCmd() *cobra.Command {
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
	// TODO: list supported cloud envs
	f.StringVar(&uc.rawAzureEnvironment, "azure-env", "AzurePublicCloud", "the target Azure cloud")
	f.StringVar(&uc.authMethod, "auth-method", "client-secret", "auth method")
	f.StringVar(&uc.rawClientID, "client-id", "", "the client ID for the Service Principal to use for authenticating to Azure")
	f.StringVar(&uc.clientSecret, "client-secret", "", "the client secret for the Service Principal to use for authenticating to Azure")
	f.StringVar(&uc.rawSubscriptionID, "subscription-id", "", "the subscription ID where the cluster is deployed")
	f.StringVar(&uc.resourceGroupName, "resource-group", "", "the resource group where the cluster is deployed")
	f.StringVar(&uc.deploymentDirectory, "deployment-dir", "", "the location of the output from `generate`")
	f.StringVar(&uc.upgradeModelFile, "upgrademodel-file", "", "file path to upgrade API model")

	return upgradeCmd
}

func (uc *upgradeCmd) validate(cmd *cobra.Command, args []string) {
	log.Infoln("validating...")

	var err error

	if uc.azureEnvironment, err = azure.EnvironmentFromName(uc.rawAzureEnvironment); err != nil {
		log.Fatal("failed to parse --azure-env as a valid target Azure cloud environment")
	}

	if uc.rawSubscriptionID == "" {
		cmd.Usage()
		log.Fatal("--subscription-id must be specified")
	}

	if uc.resourceGroupName == "" {
		cmd.Usage()
		log.Fatal("--resource-group must be specified")
	}

	// TODO(colemick): add in the cmd annotation to help enable autocompletion
	if uc.upgradeModelFile == "" {
		cmd.Usage()
		log.Fatal("--upgrademodel-file must be specified")
	}

	if uc.authMethod == "client-secret" {
		if uc.rawClientID == "" || uc.clientSecret == "" {
			cmd.Usage()
			log.Fatal("--client-id and --client-secret must be specified when --auth-method=\"client_secret\"")
		}
	} else {
		cmd.Usage()
		log.Fatal("only client secret authentication is currently supported")
	}

	if uc.clientID, err = uuid.FromString(uc.rawClientID); err != nil {
		log.Fatalf("failed to parse client id as a GUID. (client id must be specified as the application ID GUID, not the identifier_uri): %s", err.Error())
	}

	if uc.subscriptionID, err = uuid.FromString(uc.rawSubscriptionID); err != nil {
		log.Fatalf("failed to parse subscription id as a GUID: %s", err.Error())
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
}

func (uc *upgradeCmd) run(cmd *cobra.Command, args []string) error {
	uc.validate(cmd, args)

	client, err := armhelpers.NewAzureClientWithClientSecret(uc.azureEnvironment, uc.subscriptionID.String(), uc.clientID.String(), uc.clientSecret)
	if err != nil {
		log.Fatalln("Failed to retrive access token for Azure:", err)
	}

	upgradeCluster := operations.UpgradeCluster{
		AzureClient: client,
	}

	upgradeCluster.UpgradeCluster(uc.resourceGroupName, uc.containerService, uc.upgradeContainerService)

	return nil
}
