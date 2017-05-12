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
	generateName             = "generate"
	generateShortDescription = "Generate an Azure Resource Manager template"
	generateLongDescription  = "Generates an Azure Resource Manager template, parameters file and other assets for a cluster"
)

type generateCmd struct {
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
	client        armhelpers.UberClient
	deploy        bool
	resourceGroup string
	random        *rand.Rand
	location      string
}

func NewGenerateCmd() *cobra.Command {
	gc := generateCmd{}

	generateCmd := &cobra.Command{
		Use:   generateName,
		Short: generateShortDescription,
		Long:  generateLongDescription,
		RunE: func(cmd *cobra.Command, args []string) error {
			gc.validate(cmd, args)
			return gc.run()
		},
	}

	f := generateCmd.Flags()
	f.StringVar(&gc.apimodelPath, "api-model", "", "")
	f.StringVar(&gc.outputDirectory, "output-directory", "", "output directory (derived from FQDN if absent)")
	f.StringVar(&gc.caCertificatePath, "ca-certificate-path", "", "path to the CA certificate to use for Kubernetes PKI assets")
	f.StringVar(&gc.caPrivateKeyPath, "ca-private-key-path", "", "path to the CA private key to use for Kubernetes PKI assets")
	f.BoolVar(&gc.classicMode, "classic-mode", false, "enable classic parameters and outputs")
	f.BoolVar(&gc.noPrettyPrint, "no-pretty-print", false, "skip pretty printing the output")
	f.BoolVar(&gc.parametersOnly, "parameters-only", false, "only output parameters files")
	f.BoolVar(&gc.deploy, "deploy", false, "deploy as well")
	f.StringVar(&gc.resourceGroup, "resource-group", "", "resource group to deploy to")
	f.StringVar(&gc.location, "location", "", "location to deploy to")

	addAuthFlags(&gc.authArgs, f)

	return generateCmd
}

func (gc *generateCmd) validate(cmd *cobra.Command, args []string) {
	var caCertificateBytes []byte
	var caKeyBytes []byte

	if gc.apimodelPath == "" {
		if len(args) > 0 {
			gc.apimodelPath = args[0]
		} else if len(args) > 1 {
			cmd.Usage()
			log.Fatalln("too many arguments were provided to 'generate'")
		} else {
			cmd.Usage()
			log.Fatalln("--api-model was not supplied, nor was one specified as a positional argument")
		}
	}

	if _, err := os.Stat(gc.apimodelPath); os.IsNotExist(err) {
		log.Fatalf("specified api model does not exist (%s)", gc.apimodelPath)
	}

	containerService, apiVersion, err := api.LoadContainerServiceFromFile(gc.apimodelPath)
	if err != nil {
		log.Fatalf("error parsing the api model: %s", err.Error())
	}

	if gc.outputDirectory == "" {
		gc.outputDirectory = path.Join("_output", containerService.Properties.MasterProfile.DNSPrefix)
	}

	if len(caKeyBytes) != 0 {
		// the caKey is not in the api model, and should be stored separately from the model
		// we put these in the model after model is deserialized
		containerService.Properties.CertificateProfile.CaCertificate = string(caCertificateBytes)
		containerService.Properties.CertificateProfile.SetCAPrivateKey(string(caKeyBytes))
	}

	gc.containerService = containerService
	gc.apiVersion = apiVersion

	if gc.deploy {
		gc.client, err = gc.authArgs.getClient()
		if err != nil {
			log.Fatalf("failed to get client") // TODO: cleanup
		}

		if gc.resourceGroup == "" {
			cmd.Usage()
			log.Fatal("--resource-group is required when deploying")
		}
	}

	gc.random = rand.New(rand.NewSource(time.Now().UnixNano()))
}

func (gc *generateCmd) run() error {
	log.Infoln("Generating...")

	templateGenerator, err := acsengine.InitializeTemplateGenerator(gc.classicMode)
	if err != nil {
		log.Fatalln("failed to initialize template generator: %s", err.Error())
	}

	certsGenerated := false
	template, parameters, certsGenerated, err := templateGenerator.GenerateTemplate(gc.containerService)
	if err != nil {
		log.Fatalf("error generating template %s: %s", gc.apimodelPath, err.Error())
		os.Exit(1)
	}

	if !gc.noPrettyPrint {
		if template, err = acsengine.PrettyPrintArmTemplate(template); err != nil {
			log.Fatalf("error pretty printing template: %s \n", err.Error())
		}
		if parameters, err = acsengine.PrettyPrintJSON(parameters); err != nil {
			log.Fatalf("error pretty printing template parameters: %s \n", err.Error())
		}
	}

	if err = acsengine.WriteArtifacts(gc.containerService, gc.apiVersion, template, parameters, gc.outputDirectory, certsGenerated, gc.parametersOnly); err != nil {
		log.Fatalf("error writing artifacts: %s \n", err.Error())
	}

	if gc.deploy {
		templateJSON := make(map[string]interface{})
		parametersJSON := make(map[string]interface{})

		err = json.Unmarshal([]byte(template), &templateJSON)
		if err != nil {
			log.Fatalln(err)
		}

		deploymentSuffix := gc.random.Int31()

		// TODO(colemick): precreate resource group based on location

		err = json.Unmarshal([]byte(parameters), &parametersJSON)
		if err != nil {
			log.Fatalln(err)
		}
		_, err := armhelpers.DeployTemplate(
			gc.client.TemplateDeployer(),
			gc.resourceGroup,
			fmt.Sprintf("%s-%d", gc.resourceGroup, deploymentSuffix),
			templateJSON,
			parametersJSON)
		if err != nil {
			log.Fatalln(err)
		}
	}

	return nil
}
