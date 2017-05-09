package cmd

import (
	"os"
	"path"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/Azure/acs-engine/pkg/acsengine"
	"github.com/Azure/acs-engine/pkg/api"
)

const (
	generateName             = "generate"
	generateShortDescription = "generate ARM assets"
	generateLongDescription  = "generates ARM template, parameters for a container orchestrator cluster (also includes PKI assets and a `kubeconfig` for Kubernetes)"
)

type generateCmd struct {
	clusterDefinitionPath string
	outputDirectory       string // will be auto-determined from clusterDefinition
	caCertificatePath     string
	caPrivateKeyPath      string
	classicMode           bool
	noPrettyPrint         bool
	parametersOnly        bool

	// Parsed from inputs
	containerService *api.ContainerService
	apiVersion       string
}

func NewGenerateCmd() *cobra.Command {
	gc := generateCmd{}

	generateCmd := &cobra.Command{
		Use:   generateName,
		Short: generateShortDescription,
		Long:  generateLongDescription,
		RunE: func(cmd *cobra.Command, args []string) error {
			return gc.run(cmd, args)
		},
	}

	f := generateCmd.Flags()
	f.StringVar(&gc.clusterDefinitionPath, "cluster-definition", "", "")
	f.StringVar(&gc.outputDirectory, "output-directory", "", "output directory (derived from FQDN if absent)")
	f.StringVar(&gc.caCertificatePath, "ca-certificate-path", "", "path to the CA certificate to use for Kubernetes PKI assets")
	f.StringVar(&gc.caPrivateKeyPath, "ca-private-key-path", "", "path to the CA private key to use for Kubernetes PKI assets")
	f.BoolVar(&gc.classicMode, "classic-mode", false, "enable classic parameters and outputs")
	f.BoolVar(&gc.noPrettyPrint, "no-pretty-print", false, "skip pretty printing the output")
	f.BoolVar(&gc.parametersOnly, "parameters-only", false, "only output parameters files")

	return generateCmd
}

func (gc *generateCmd) validate(cmd *cobra.Command, args []string) {
	log.Warnln("Validating...")

	var caCertificateBytes []byte
	var caKeyBytes []byte

	if gc.clusterDefinitionPath == "" {
		if len(args) > 0 {
			gc.clusterDefinitionPath = args[0]
		} else if len(args) > 1 {
			log.Fatalln("too many arguments were provided to 'generate'")
		} else {
			log.Fatalln("--cluster-definition was not supplied, nor was one specified as a positional argument")
		}
	}

	if _, err := os.Stat(gc.clusterDefinitionPath); os.IsNotExist(err) {
		log.Fatalf("specified cluster definition does not exist (%s)", gc.clusterDefinitionPath)
	}

	containerService, apiVersion, err := api.LoadContainerServiceFromFile(gc.clusterDefinitionPath)
	if err != nil {
		log.Fatalf("error loading the container service model from the cluster definition: %s", err.Error())
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
}

func (gc *generateCmd) run(cmd *cobra.Command, args []string) error {
	gc.validate(cmd, args)
	log.Infoln("Generating...")

	templateGenerator, err := acsengine.InitializeTemplateGenerator(gc.classicMode)
	if err != nil {
		log.Fatalln("failed to initialize template generator: %s", err.Error())
	}

	certsGenerated := false
	template, parameters, certsGenerated, err := templateGenerator.GenerateTemplate(gc.containerService)
	if err != nil {
		log.Fatalf("error generating template %s: %s", gc.clusterDefinitionPath, err.Error())
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

	return nil
}
