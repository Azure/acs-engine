package cmd

import (
	"github.com/Azure/acs-engine/pkg/acsengine"

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
	kubeconfigPath      string
	rawClientID         string
	rawSubscriptionID   string
	rawAzureEnvironment string

	// parsed
	clientID         uuid.UUID
	subscriptionID   uuid.UUID
	azureEnvironment azure.Environment

	// derived
	tenantID              string
	servicePrincipalToken *adal.ServicePrincipalToken
}

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
	f.StringVar(&uc.rawAzureEnvironment, "azure-env", "AzurePublicCloud", "the target Azure cloud (default:`AzurePublicCloud`)")
	f.StringVar(&uc.authMethod, "auth-method", "client-secret", "auth method (default:`client_secret`)")
	f.StringVar(&uc.rawClientID, "client-id", "", "the client ID for the Service Principal to use for authenticating to Azure")
	f.StringVar(&uc.clientSecret, "client-secret", "", "the client secret for the Service Principal to use for authenticating to Azure")
	f.StringVar(&uc.rawSubscriptionID, "subscription-id", "", "the subscription ID where the cluster is deployed")
	f.StringVar(&uc.resourceGroupName, "resource-group", "", "the resource group where the cluster is deployed")
	f.StringVar(&uc.kubeconfigPath, "kubeconfig", "", "path to the kubeconfig file for the cluster")

	return upgradeCmd
}

func (uc *upgradeCmd) validate(cmd *cobra.Command, args []string) {
	log.Warnln("validating...")

	var err error

	if uc.azureEnvironment, err = azure.EnvironmentFromName(uc.rawAzureEnvironment); err != nil {
		log.Fatal("failed to parse --azure-env as a valid target Azure cloud environment")
	}

	if uc.rawSubscriptionID == "" {
		log.Fatal("--subscription-id must be specified")
	}

	if uc.kubeconfigPath == "" {
		log.Fatal("--kubeconfig must be specified")
	}

	if uc.resourceGroupName == "" {
		log.Fatal("--resource-group must be specified")
	}

	if uc.authMethod == "client-secret" {
		if uc.rawClientID == "" || uc.clientSecret == "" {
			log.Fatal("--client-id and --client-secret must be specified when --auth-method=\"client_secret\"")
		}
	} else {
		log.Fatal("only client secret authentication is currently supported")
	}

	if uc.clientID, err = uuid.FromString(uc.rawClientID); err != nil {
		log.Fatalf("failed to parse client id as a GUID. (client id must be specified as the application ID GUID, not the identifier_uri): %s", err.Error())
	}

	if uc.subscriptionID, err = uuid.FromString(uc.rawSubscriptionID); err != nil {
		log.Fatalf("failed to parse subscription id as a GUID: %s", err.Error())
	}

	// get the actual ServicePrincipalToken here
	// one, bounce off sub, get tenant id
	// two, create NewServicePrincipalToken()

	tenantID, err := acsengine.GetTenantID(uc.azureEnvironment, uc.subscriptionID.String())
	if err != nil {
		log.Fatalf("failed to determine tenant id based on subscription id: %s", err.Error())
	}

	oauthConfig, err := adal.NewOAuthConfig(uc.azureEnvironment.ActiveDirectoryEndpoint, tenantID)
	if err != nil {
		log.Fatalf("failed to create oauth configuration: %s", err.Error())
	}

	uc.servicePrincipalToken, err = adal.NewServicePrincipalToken(*oauthConfig, uc.clientID.String(), uc.clientSecret, uc.azureEnvironment.ResourceManagerEndpoint)
	if err != nil {
		log.Fatalf("failed to retrieve AccessToken for the Service Principal")
	}

	err = uc.servicePrincipalToken.Refresh()
	if err != nil {
		log.Fatalf("failed to refresh AccessToken: %s", err.Error())
	}
}

func (uc *upgradeCmd) run(cmd *cobra.Command, args []string) error {
	uc.validate(cmd, args)
	log.Infoln("upgrade procedure... beginning.")

	// drive the actual upgrade process using uc.ServicePrincipalToken

	log.Infoln("upgrade procedure... finished.")

	return nil
}
