package cmd

import (
	"os"

	"github.com/Azure/acs-engine/pkg/armhelpers"

	"github.com/Azure/go-autorest/autorest/azure"
	log "github.com/Sirupsen/logrus"
	"github.com/satori/go.uuid"
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
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
		rootCmd.AddCommand(newDeployCmd())
	}

	return rootCmd
}

type authArgs struct {
	RawAzureEnvironment string
	rawSubscriptionID   string
	SubscriptionID      uuid.UUID
	AuthMethod          string
	rawClientID         string

	ClientID        uuid.UUID
	ClientSecret    string
	CertificatePath string
	PrivateKeyPath  string
}

func addAuthFlags(authArgs *authArgs, f *flag.FlagSet) {
	f.StringVar(&authArgs.RawAzureEnvironment, "azure-env", "AzurePublicCloud", "the target Azure cloud")
	f.StringVar(&authArgs.rawSubscriptionID, "subscription-id", "", "azure subscription id")
	f.StringVar(&authArgs.AuthMethod, "auth-method", "device", "auth method (default:`device`, `client_secret`, `client_certificate`)")
	f.StringVar(&authArgs.rawClientID, "client-id", "", "client id (used with --auth-method=[client_secret|client_certificate])")
	f.StringVar(&authArgs.ClientSecret, "client-secret", "", "client secret (used with --auth-mode=client_secret)")
	f.StringVar(&authArgs.CertificatePath, "certificate-path", "", "path to client certificate (used with --auth-method=client_certificate)")
	f.StringVar(&authArgs.PrivateKeyPath, "private-key-path", "", "path to private key (used with --auth-method=client_certificate)")
}

func (authArgs *authArgs) getClient() (*armhelpers.AzureClient, error) {
	authArgs.ClientID, _ = uuid.FromString(authArgs.rawClientID)
	authArgs.SubscriptionID, _ = uuid.FromString(authArgs.rawSubscriptionID)

	if authArgs.AuthMethod == "client_secret" {
		if authArgs.ClientID.String() == "" || authArgs.ClientSecret == "" {
			log.Fatal(`--client-id and --client-secret must be specified when --auth-method="client_secret"`)
		}
		// try parse the UUID
	} else if authArgs.AuthMethod == "client_certificate" {
		if authArgs.ClientID.String() == "" || authArgs.CertificatePath == "" || authArgs.PrivateKeyPath == "" {
			log.Fatal(`--client-id and --certificate-path, and --private-key-path must be specified when --auth-method="client_certificate"`)
		}
	}

	if authArgs.SubscriptionID.String() == "" {
		log.Fatal("--subscription-id is required (and must be a valid UUID)")
	}

	env, err := azure.EnvironmentFromName(authArgs.RawAzureEnvironment)
	if err != nil {
		log.Fatal("failed to parse --azure-env as a valid target Azure cloud environment")
	}

	switch authArgs.AuthMethod {
	case "device":
		return armhelpers.NewAzureClientWithDeviceAuth(env, authArgs.SubscriptionID.String())
	case "client_secret":
		return armhelpers.NewAzureClientWithClientSecret(env, authArgs.SubscriptionID.String(), authArgs.ClientID.String(), authArgs.ClientSecret)
	case "client_certificate":
		return armhelpers.NewAzureClientWithClientCertificate(env, authArgs.SubscriptionID.String(), authArgs.ClientID.String(), authArgs.CertificatePath, authArgs.PrivateKeyPath)
	default:
		log.Fatalf("--auth-method: ERROR: method unsupported. method=%q.", authArgs.AuthMethod)
	}

	return nil, nil // unreachable
}
