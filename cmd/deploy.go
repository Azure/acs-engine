package cmd

import (
	"fmt"
	"math/rand"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/leonelquinteros/gotext"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"encoding/json"

	"github.com/Azure/acs-engine/pkg/acsengine"
	"github.com/Azure/acs-engine/pkg/api"
	"github.com/Azure/acs-engine/pkg/armhelpers"
	"github.com/Azure/acs-engine/pkg/i18n"
)

const (
	deployName             = "deploy"
	deployShortDescription = "deploy an Azure Resource Manager template"
	deployLongDescription  = "deploys an Azure Resource Manager template, parameters file and other assets for a cluster"
)

type deployCmd struct {
	authArgs

	apimodelPath      string
	dnsPrefix         string
	autoSuffix        bool
	outputDirectory   string // can be auto-determined from clusterDefinition
	caCertificatePath string
	caPrivateKeyPath  string
	classicMode       bool
	noPrettyPrint     bool
	parametersOnly    bool

	// derived
	containerService *api.ContainerService
	apiVersion       string
	locale           *gotext.Locale

	client        armhelpers.ACSEngineClient
	resourceGroup string
	random        *rand.Rand
	location      string
}

func newDeployCmd() *cobra.Command {
	dc := deployCmd{}

	deployCmd := &cobra.Command{
		Use:   deployName,
		Short: deployShortDescription,
		Long:  deployLongDescription,
		RunE: func(cmd *cobra.Command, args []string) error {
			dc.validate(cmd, args)
			return dc.run()
		},
	}

	f := deployCmd.Flags()
	f.StringVar(&dc.apimodelPath, "api-model", "", "path to the apimodel file")
	f.StringVar(&dc.dnsPrefix, "dns-prefix", "", "dns prefix (unique name for the cluster)")
	f.BoolVar(&dc.autoSuffix, "auto-suffix", false, "automatically append a compressed timestamp to the dnsPrefix to ensure unique cluster name automatically")
	f.StringVar(&dc.outputDirectory, "output-directory", "", "output directory (derived from FQDN if absent)")
	f.StringVar(&dc.caCertificatePath, "ca-certificate-path", "", "path to the CA certificate to use for Kubernetes PKI assets")
	f.StringVar(&dc.caPrivateKeyPath, "ca-private-key-path", "", "path to the CA private key to use for Kubernetes PKI assets")
	f.StringVar(&dc.resourceGroup, "resource-group", "", "resource group to deploy to")
	f.StringVar(&dc.location, "location", "", "location to deploy to")

	addAuthFlags(&dc.authArgs, f)

	return deployCmd
}

func (dc *deployCmd) validate(cmd *cobra.Command, args []string) {
	var err error

	dc.locale, err = i18n.LoadTranslations()
	if err != nil {
		log.Fatalf("error loading translation files: %s", err.Error())
	}

	if dc.apimodelPath == "" {
		if len(args) > 0 {
			dc.apimodelPath = args[0]
		} else if len(args) > 1 {
			cmd.Usage()
			log.Fatalln("too many arguments were provided to 'deploy'")
		} else {
			cmd.Usage()
			log.Fatalln("--api-model was not supplied, nor was one specified as a positional argument")
		}
	}

	if _, err := os.Stat(dc.apimodelPath); os.IsNotExist(err) {
		log.Fatalf("specified api model does not exist (%s)", dc.apimodelPath)
	}

	apiloader := &api.Apiloader{
		Translator: &i18n.Translator{
			Locale: dc.locale,
		},
	}
	// skip validating the model fields for now
	dc.containerService, dc.apiVersion, err = apiloader.LoadContainerServiceFromFile(dc.apimodelPath, false, nil)
	if err != nil {
		log.Fatalf("error parsing the api model: %s", err.Error())
	}

	if dc.location == "" {
		log.Fatalf("--location must be specified")
	}

	dc.client, err = dc.authArgs.getClient()
	if err != nil {
		log.Fatalf("failed to get client") // TODO: cleanup
	}

	// autofillApimodel calls log.Fatal() directly and does not return errors
	autofillApimodel(dc)

	_, _, err = revalidateApimodel(apiloader, dc.containerService, dc.apiVersion)
	if err != nil {
		log.Fatalf("Failed to validate the apimodel after populating values: %s", err)
	}

	dc.random = rand.New(rand.NewSource(time.Now().UnixNano()))
}

func autofillApimodel(dc *deployCmd) {
	var err error

	if dc.containerService.Properties.LinuxProfile.AdminUsername == "" {
		log.Warnf("apimodel: no linuxProfile.adminUsername was specified. Will use 'azureuser'.")
		dc.containerService.Properties.LinuxProfile.AdminUsername = "azureuser"
	}

	if dc.dnsPrefix != "" && dc.containerService.Properties.MasterProfile.DNSPrefix != "" {
		log.Fatalf("invalid configuration: the apimodel masterProfile.dnsPrefix and --dns-prefix were both specified")
	}
	if dc.containerService.Properties.MasterProfile.DNSPrefix == "" {
		if dc.dnsPrefix == "" {
			log.Fatalf("apimodel: missing masterProfile.dnsPrefix and --dns-prefix was not specified")
		}

		dnsPrefix := dc.dnsPrefix
		if dc.autoSuffix {
			suffix := strconv.FormatInt(time.Now().Unix(), 16)
			dnsPrefix = dnsPrefix + "-" + suffix
		}

		log.Warnf("apimodel: missing masterProfile.dnsPrefix will use %q", dnsPrefix)
		dc.containerService.Properties.MasterProfile.DNSPrefix = dnsPrefix
	}

	if dc.outputDirectory == "" {
		dc.outputDirectory = path.Join("_output", dc.containerService.Properties.MasterProfile.DNSPrefix)
	}

	if dc.resourceGroup == "" {
		dnsPrefix := dc.containerService.Properties.MasterProfile.DNSPrefix
		log.Warnf("--resource-group was not specified. Using the DNS prefix from the apimodel as the resource group name: %s", dnsPrefix)
		dc.resourceGroup = dnsPrefix
		if dc.location == "" {
			log.Fatal("--resource-group was not specified. --location must be specified in case the resource group needs creation.")
		}
	}

	if dc.containerService.Properties.LinuxProfile.SSH.PublicKeys == nil ||
		len(dc.containerService.Properties.LinuxProfile.SSH.PublicKeys) == 0 ||
		dc.containerService.Properties.LinuxProfile.SSH.PublicKeys[0].KeyData == "" {
		creator := &acsengine.SSHCreator{
			Translator: &i18n.Translator{
				Locale: dc.locale,
			},
		}
		_, publicKey, err := creator.CreateSaveSSH(dc.containerService.Properties.LinuxProfile.AdminUsername, dc.outputDirectory)
		if err != nil {
			log.Fatal("Failed to generate SSH Key")
		}

		dc.containerService.Properties.LinuxProfile.SSH.PublicKeys = []api.PublicKey{{KeyData: publicKey}}
	}

	_, err = dc.client.EnsureResourceGroup(dc.resourceGroup, dc.location)
	if err != nil {
		log.Fatalln(err)
	}

	useManagedIdentity := dc.containerService.Properties.OrchestratorProfile.KubernetesConfig != nil &&
		dc.containerService.Properties.OrchestratorProfile.KubernetesConfig.UseManagedIdentity

	if !useManagedIdentity {
		spp := dc.containerService.Properties.ServicePrincipalProfile
		if spp != nil && spp.ClientID == "" && spp.Secret == "" && spp.KeyvaultSecretRef == nil {
			log.Warnln("apimodel: ServicePrincipalProfile was missing or empty, creating application...")

			// TODO: consider caching the creds here so they persist between subsequent runs of 'deploy'
			appName := dc.containerService.Properties.MasterProfile.DNSPrefix
			appURL := fmt.Sprintf("https://%s/", appName)
			applicationID, servicePrincipalObjectID, secret, err := dc.client.CreateApp(appName, appURL)
			if err != nil {
				log.Fatalf("apimodel invalid: ServicePrincipalProfile was empty, and we failed to create valid credentials: %q", err)
			}
			log.Warnf("created application with applicationID (%s) and servicePrincipalObjectID (%s).", applicationID, servicePrincipalObjectID)

			log.Warnln("apimodel: ServicePrincipalProfile was empty, assigning role to application...")
			for {
				err = dc.client.CreateRoleAssignmentSimple(dc.resourceGroup, servicePrincipalObjectID)
				if err != nil {
					log.Debugf("Failed to create role assignment (will retry): %q", err)
					time.Sleep(3 * time.Second)
					continue
				}
				break
			}

			dc.containerService.Properties.ServicePrincipalProfile = &api.ServicePrincipalProfile{
				ClientID: applicationID,
				Secret:   secret,
			}
		}
	}
}

func revalidateApimodel(apiloader *api.Apiloader, containerService *api.ContainerService, apiVersion string) (*api.ContainerService, string, error) {
	// This isn't terribly elegant, but it's the easiest way to go for now w/o duplicating a bunch of code
	rawVersionedAPIModel, err := apiloader.SerializeContainerService(containerService, apiVersion)
	if err != nil {
		return nil, "", err
	}
	return apiloader.DeserializeContainerService(rawVersionedAPIModel, true, nil)
}

func (dc *deployCmd) run() error {
	ctx := acsengine.Context{
		Translator: &i18n.Translator{
			Locale: dc.locale,
		},
	}

	templateGenerator, err := acsengine.InitializeTemplateGenerator(ctx, dc.classicMode)
	if err != nil {
		log.Fatalln("failed to initialize template generator: %s", err.Error())
	}

	template, parameters, certsgenerated, err := templateGenerator.GenerateTemplate(dc.containerService)
	if err != nil {
		log.Fatalf("error generating template %s: %s", dc.apimodelPath, err.Error())
		os.Exit(1)
	}

	if template, err = acsengine.PrettyPrintArmTemplate(template); err != nil {
		log.Fatalf("error pretty printing template: %s \n", err.Error())
	}
	var parametersFile string
	if parametersFile, err = acsengine.BuildAzureParametersFile(parameters); err != nil {
		log.Fatalf("error pretty printing template parameters: %s \n", err.Error())
	}

	writer := &acsengine.ArtifactWriter{
		Translator: &i18n.Translator{
			Locale: dc.locale,
		},
	}
	if err = writer.WriteTLSArtifacts(dc.containerService, dc.apiVersion, template, parametersFile, dc.outputDirectory, certsgenerated, dc.parametersOnly); err != nil {
		log.Fatalf("error writing artifacts: %s \n", err.Error())
	}

	templateJSON := make(map[string]interface{})
	parametersJSON := make(map[string]interface{})

	err = json.Unmarshal([]byte(template), &templateJSON)
	if err != nil {
		log.Fatalln(err)
	}

	err = json.Unmarshal([]byte(parameters), &parametersJSON)
	if err != nil {
		log.Fatalln(err)
	}

	deploymentSuffix := dc.random.Int31()

	_, err = dc.client.DeployTemplate(
		dc.resourceGroup,
		fmt.Sprintf("%s-%d", dc.resourceGroup, deploymentSuffix),
		templateJSON,
		parametersJSON,
		nil)
	if err != nil {
		log.Fatalln(err)
	}

	return nil
}
