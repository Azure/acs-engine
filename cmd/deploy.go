package cmd

import (
	"fmt"
	"io/ioutil"
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
	"github.com/Azure/acs-engine/pkg/acsengine/transform"
	"github.com/Azure/acs-engine/pkg/api"
	"github.com/Azure/acs-engine/pkg/armhelpers"
	"github.com/Azure/acs-engine/pkg/helpers"
	"github.com/Azure/acs-engine/pkg/i18n"
	"github.com/Azure/azure-sdk-for-go/arm/graphrbac"
	"github.com/Azure/go-autorest/autorest/to"
)

const (
	deployName             = "deploy"
	deployShortDescription = "Deploy an Azure Resource Manager template"
	deployLongDescription  = "Deploy an Azure Resource Manager template, parameters file and other assets for a cluster"

	// aadServicePrincipal is a hard-coded service principal which represents
	// Azure Active Dirctory (see az ad sp list)
	aadServicePrincipal = "00000002-0000-0000-c000-000000000000"

	// aadPermissionUserRead is the User.Read hard-coded permission on
	// aadServicePrincipal (see az ad sp list)
	aadPermissionUserRead = "311a71cc-e848-46a1-bdf8-97ff7156d8e6"
)

type deployCmd struct {
	authArgs

	apimodelPath      string
	dnsPrefix         string
	autoSuffix        bool
	outputDirectory   string // can be auto-determined from clusterDefinition
	forceOverwrite    bool
	caCertificatePath string
	caPrivateKeyPath  string
	classicMode       bool
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
			if err := dc.validate(cmd, args); err != nil {
				log.Fatalf(fmt.Sprintf("error validating deployCmd: %s", err.Error()))
			}
			if err := dc.load(cmd, args); err != nil {
				log.Fatalln("failed to load apimodel: %s", err.Error())
			}
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
	f.StringVarP(&dc.resourceGroup, "resource-group", "g", "", "resource group to deploy to (will use the DNS prefix from the apimodel if not specified)")
	f.StringVarP(&dc.location, "location", "l", "", "location to deploy to (required)")
	f.BoolVarP(&dc.forceOverwrite, "force-overwrite", "f", false, "automatically overwrite existing files in the output directory")

	addAuthFlags(&dc.authArgs, f)

	return deployCmd
}

func (dc *deployCmd) validate(cmd *cobra.Command, args []string) error {
	var err error

	dc.locale, err = i18n.LoadTranslations()
	if err != nil {
		return fmt.Errorf(fmt.Sprintf("error loading translation files: %s", err.Error()))
	}

	if dc.apimodelPath == "" {
		if len(args) == 1 {
			dc.apimodelPath = args[0]
		} else if len(args) > 1 {
			cmd.Usage()
			return fmt.Errorf(fmt.Sprintf("too many arguments were provided to 'deploy'"))
		} else {
			cmd.Usage()
			return fmt.Errorf(fmt.Sprintf("--api-model was not supplied, nor was one specified as a positional argument"))
		}
	}

	if _, err := os.Stat(dc.apimodelPath); os.IsNotExist(err) {
		return fmt.Errorf(fmt.Sprintf("specified api model does not exist (%s)", dc.apimodelPath))
	}

	if dc.location == "" {
		return fmt.Errorf(fmt.Sprintf("--location must be specified"))
	}
	dc.location = helpers.NormalizeAzureRegion(dc.location)

	return nil
}

func (dc *deployCmd) load(cmd *cobra.Command, args []string) error {
	var err error

	apiloader := &api.Apiloader{
		Translator: &i18n.Translator{
			Locale: dc.locale,
		},
	}

	// do not validate when initially loading the apimodel, validation is done later after autofilling values
	dc.containerService, dc.apiVersion, err = apiloader.LoadContainerServiceFromFile(dc.apimodelPath, false, false, nil)
	if err != nil {
		return fmt.Errorf(fmt.Sprintf("error parsing the api model: %s", err.Error()))
	}

	if dc.containerService.Location == "" {
		dc.containerService.Location = dc.location
	} else if dc.containerService.Location != dc.location {
		return fmt.Errorf(fmt.Sprintf("--location does not match api model location"))
	}

	if err = dc.authArgs.validateAuthArgs(); err != nil {
		return fmt.Errorf("%s", err)
	}

	dc.client, err = dc.authArgs.getClient()
	if err != nil {
		return fmt.Errorf("failed to get client: %s", err.Error())
	}

	if err = autofillApimodel(dc); err != nil {
		return err
	}

	_, _, err = validateApimodel(apiloader, dc.containerService, dc.apiVersion)
	if err != nil {
		return fmt.Errorf("Failed to validate the apimodel after populating values: %s", err)
	}

	dc.random = rand.New(rand.NewSource(time.Now().UnixNano()))

	return nil
}

func autofillApimodel(dc *deployCmd) error {
	var err error

	if dc.containerService.Properties.LinuxProfile != nil {
		if dc.containerService.Properties.LinuxProfile.AdminUsername == "" {
			log.Warnf("apimodel: no linuxProfile.adminUsername was specified. Will use 'azureuser'.")
			dc.containerService.Properties.LinuxProfile.AdminUsername = "azureuser"
		}
	}

	if dc.dnsPrefix != "" && dc.containerService.Properties.MasterProfile.DNSPrefix != "" {
		return fmt.Errorf("invalid configuration: the apimodel masterProfile.dnsPrefix and --dns-prefix were both specified")
	}
	if dc.containerService.Properties.MasterProfile.DNSPrefix == "" {
		if dc.dnsPrefix == "" {
			return fmt.Errorf("apimodel: missing masterProfile.dnsPrefix and --dns-prefix was not specified")
		}
		log.Warnf("apimodel: missing masterProfile.dnsPrefix will use %q", dc.dnsPrefix)
		dc.containerService.Properties.MasterProfile.DNSPrefix = dc.dnsPrefix
	}

	if dc.autoSuffix {
		suffix := strconv.FormatInt(time.Now().Unix(), 16)
		dc.containerService.Properties.MasterProfile.DNSPrefix += "-" + suffix
	}

	if dc.outputDirectory == "" {
		dc.outputDirectory = path.Join("_output", dc.containerService.Properties.MasterProfile.DNSPrefix)
	}

	if _, err := os.Stat(dc.outputDirectory); !dc.forceOverwrite && err == nil {
		return fmt.Errorf("Output directory already exists and forceOverwrite flag is not set: %s", dc.outputDirectory)
	}

	if dc.resourceGroup == "" {
		dnsPrefix := dc.containerService.Properties.MasterProfile.DNSPrefix
		log.Warnf("--resource-group was not specified. Using the DNS prefix from the apimodel as the resource group name: %s", dnsPrefix)
		dc.resourceGroup = dnsPrefix
		if dc.location == "" {
			return fmt.Errorf("--resource-group was not specified. --location must be specified in case the resource group needs creation")
		}
	}

	if dc.containerService.Properties.LinuxProfile != nil && (dc.containerService.Properties.LinuxProfile.SSH.PublicKeys == nil ||
		len(dc.containerService.Properties.LinuxProfile.SSH.PublicKeys) == 0 ||
		dc.containerService.Properties.LinuxProfile.SSH.PublicKeys[0].KeyData == "") {
		translator := &i18n.Translator{
			Locale: dc.locale,
		}
		_, publicKey, err := acsengine.CreateSaveSSH(dc.containerService.Properties.LinuxProfile.AdminUsername, dc.outputDirectory, translator)
		if err != nil {
			return fmt.Errorf("Failed to generate SSH Key: %s", err.Error())
		}

		dc.containerService.Properties.LinuxProfile.SSH.PublicKeys = []api.PublicKey{{KeyData: publicKey}}
	}

	_, err = dc.client.EnsureResourceGroup(dc.resourceGroup, dc.location, nil)
	if err != nil {
		return err
	}

	useManagedIdentity := dc.containerService.Properties.OrchestratorProfile.KubernetesConfig != nil &&
		dc.containerService.Properties.OrchestratorProfile.KubernetesConfig.UseManagedIdentity

	if !useManagedIdentity {
		spp := dc.containerService.Properties.ServicePrincipalProfile
		if spp != nil && spp.ClientID == "" && spp.Secret == "" && spp.KeyvaultSecretRef == nil && (dc.ClientID.String() == "" || dc.ClientID.String() == "00000000-0000-0000-0000-000000000000") && dc.ClientSecret == "" {
			log.Warnln("apimodel: ServicePrincipalProfile was missing or empty, creating application...")

			// TODO: consider caching the creds here so they persist between subsequent runs of 'deploy'
			appName := dc.containerService.Properties.MasterProfile.DNSPrefix
			appURL := fmt.Sprintf("https://%s/", appName)
			var replyURLs *[]string
			var requiredResourceAccess *[]graphrbac.RequiredResourceAccess
			if dc.containerService.Properties.OrchestratorProfile.OrchestratorType == api.OpenShift {
				appName = fmt.Sprintf("%s.%s.cloudapp.azure.com", appName, dc.containerService.Properties.AzProfile.Location)
				appURL = fmt.Sprintf("https://%s:8443/", appName)
				replyURLs = to.StringSlicePtr([]string{fmt.Sprintf("https://%s:8443/oauth2callback/Azure%%20AD", appName)})
				requiredResourceAccess = &[]graphrbac.RequiredResourceAccess{
					{
						ResourceAppID: to.StringPtr(aadServicePrincipal),
						ResourceAccess: &[]graphrbac.ResourceAccess{
							{
								ID:   to.StringPtr(aadPermissionUserRead),
								Type: to.StringPtr("Scope"),
							},
						},
					},
				}
			}
			applicationID, servicePrincipalObjectID, secret, err := dc.client.CreateApp(appName, appURL, replyURLs, requiredResourceAccess)
			if err != nil {
				return fmt.Errorf("apimodel invalid: ServicePrincipalProfile was empty, and we failed to create valid credentials: %q", err)
			}
			log.Warnf("created application with applicationID (%s) and servicePrincipalObjectID (%s).", applicationID, servicePrincipalObjectID)

			log.Warnln("apimodel: ServicePrincipalProfile was empty, assigning role to application...")

			err = dc.client.CreateRoleAssignmentSimple(dc.resourceGroup, servicePrincipalObjectID)
			if err != nil {
				return fmt.Errorf("apimodel: could not create or assign ServicePrincipal: %q", err)

			}

			dc.containerService.Properties.ServicePrincipalProfile = &api.ServicePrincipalProfile{
				ClientID: applicationID,
				Secret:   secret,
				ObjectID: servicePrincipalObjectID,
			}
		} else if (dc.containerService.Properties.ServicePrincipalProfile == nil || ((dc.containerService.Properties.ServicePrincipalProfile.ClientID == "" || dc.containerService.Properties.ServicePrincipalProfile.ClientID == "00000000-0000-0000-0000-000000000000") && dc.containerService.Properties.ServicePrincipalProfile.Secret == "")) && dc.ClientID.String() != "" && dc.ClientSecret != "" {
			dc.containerService.Properties.ServicePrincipalProfile = &api.ServicePrincipalProfile{
				ClientID: dc.ClientID.String(),
				Secret:   dc.ClientSecret,
			}
		}
	}
	return nil
}

func validateApimodel(apiloader *api.Apiloader, containerService *api.ContainerService, apiVersion string) (*api.ContainerService, string, error) {
	// This isn't terribly elegant, but it's the easiest way to go for now w/o duplicating a bunch of code
	rawVersionedAPIModel, err := apiloader.SerializeContainerService(containerService, apiVersion)
	if err != nil {
		return nil, "", err
	}
	return apiloader.DeserializeContainerService(rawVersionedAPIModel, true, false, nil)
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

	template, parameters, certsgenerated, err := templateGenerator.GenerateTemplate(dc.containerService, acsengine.DefaultGeneratorCode, false, BuildTag)
	if err != nil {
		log.Fatalf("error generating template %s: %s", dc.apimodelPath, err.Error())
		os.Exit(1)
	}

	if template, err = transform.PrettyPrintArmTemplate(template); err != nil {
		log.Fatalf("error pretty printing template: %s \n", err.Error())
	}
	var parametersFile string
	if parametersFile, err = transform.BuildAzureParametersFile(parameters); err != nil {
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

	if res, err := dc.client.DeployTemplate(
		dc.resourceGroup,
		fmt.Sprintf("%s-%d", dc.resourceGroup, deploymentSuffix),
		templateJSON,
		parametersJSON,
		nil); err != nil {
		if res != nil && res.Response.Response != nil && res.Body != nil {
			defer res.Body.Close()
			body, _ := ioutil.ReadAll(res.Body)
			log.Errorf(string(body))
		}
		log.Fatalln(err)
	}

	return nil
}
