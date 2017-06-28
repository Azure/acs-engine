package cmd

import (
	"fmt"
	"math/rand"
	"os"
	"path"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"

	"encoding/json"

	"github.com/Azure/acs-engine/pkg/acsengine"
	"github.com/Azure/acs-engine/pkg/api"
	"github.com/Azure/acs-engine/pkg/armhelpers"
)

const (
	deployName             = "deploy"
	deployShortDescription = "deploy an Azure Resource Manager template"
	deployLongDescription  = "deploys an Azure Resource Manager template, parameters file and other assets for a cluster"
)

type deployCmd struct {
	authArgs

	apimodelPath      string
	outputDirectory   string // can be auto-determined from clusterDefinition
	caCertificatePath string
	caPrivateKeyPath  string
	classicMode       bool
	noPrettyPrint     bool
	parametersOnly    bool

	// derived
	containerService *api.ContainerService
	apiVersion       string

	// experimental
	client        armhelpers.ACSEngineClient
	deploy        bool
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
	f.StringVar(&dc.apimodelPath, "api-model", "", "")
	f.StringVar(&dc.outputDirectory, "output-directory", "", "output directory (derived from FQDN if absent)")
	f.StringVar(&dc.caCertificatePath, "ca-certificate-path", "", "path to the CA certificate to use for Kubernetes PKI assets")
	f.StringVar(&dc.caPrivateKeyPath, "ca-private-key-path", "", "path to the CA private key to use for Kubernetes PKI assets")
	f.BoolVar(&dc.deploy, "deploy", false, "deploy as well")
	f.StringVar(&dc.resourceGroup, "resource-group", "", "resource group to deploy to")
	f.StringVar(&dc.location, "location", "", "location to deploy to")

	addAuthFlags(&dc.authArgs, f)

	return deployCmd
}

func (dc *deployCmd) validate(cmd *cobra.Command, args []string) {
	var caCertificateBytes []byte
	var caKeyBytes []byte
	var err error

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

	dc.containerService, dc.apiVersion, err = api.LoadContainerServiceFromFile(dc.apimodelPath)
	if err != nil {
		log.Fatalf("error parsing the api model: %s", err.Error())
	}

	if dc.outputDirectory == "" {
		dc.outputDirectory = path.Join("_output", dc.containerService.Properties.MasterProfile.DNSPrefix)
	}

	if dc.resourceGroup == "" {
		dnsPrefix := dc.containerService.Properties.MasterProfile.DNSPrefix
		log.Warnf("--resource-group was not specified. Using the DNS prefix from the apimodel as the resource group name: %s")
		dc.resourceGroup = dnsPrefix
		if dc.location == "" {
			// TODO: move this so we only require location for a non-pre-existing RG?
			log.Fatal("--resource-group was not specified. --location must be specified in case the resource group needs creation.")
		}
	}

	if len(caKeyBytes) != 0 {
		// the caKey is not in the api model, and should be stored separately from the model
		// we put these in the model after model is deserialized
		dc.containerService.Properties.CertificateProfile.CaCertificate = string(caCertificateBytes)
		dc.containerService.Properties.CertificateProfile.CaPrivateKey = string(caKeyBytes)
	}

	dc.client, err = dc.authArgs.getClient()
	if err != nil {
		log.Fatalf("failed to get client") // TODO: cleanup
	}

	dc.random = rand.New(rand.NewSource(time.Now().UnixNano()))
}

func (dc *deployCmd) run() error {
	templateGenerator, err := acsengine.InitializeTemplateGenerator(dc.classicMode)
	if err != nil {
		log.Fatalln("failed to initialize template generator: %s", err.Error())
	}

	certsgenerated := false
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

	if err = acsengine.WriteArtifacts(dc.containerService, dc.apiVersion, template, parametersFile, dc.outputDirectory, certsgenerated, dc.parametersOnly); err != nil {
		log.Fatalf("error writing artifacts: %s \n", err.Error())
	}

	log.Infoln("deploying...")

	_, err = dc.client.EnsureResourceGroup(dc.resourceGroup, dc.location)
	if err != nil {
		log.Fatalln(err)
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
