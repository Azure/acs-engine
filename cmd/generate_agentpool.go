package cmd

import (
	"fmt"
	"github.com/Azure/acs-engine/pkg/api"
	"github.com/Azure/acs-engine/pkg/interpolator/agentpool"
	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
)

type GenerateOptions struct {
	ApiModelPath string
	//OutputDirectory   string
	//CaCertificatePath string
	//CaPrivateKeyPath  string
	//ClassicMode       bool
	//NoPrettyPrint     bool
	//ParametersOnly    bool

	ContainerService *api.ContainerService
	ApiVersion       string
}

func NewGenerateAgentpoolCmd() *cobra.Command {
	genOptions := GenerateOptions{}

	genAgentpoolCmd := &cobra.Command{
		Use:   "agentpool",
		Short: "Create agent pools for existing Kubernetes control plane infrastructure",
		Long:  "Create agent pools for existing Kubernetes control plane infrastructure",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := genOptions.Init(cmd, args)
			if err != nil {
				log.Fatal(err)
			}
			genOptions.Validate(cmd, args)
			err = genOptions.Run()
			if err != nil {
				log.Fatal(err)
			}
			return nil
		},
	}

	f := genAgentpoolCmd.Flags()
	f.StringVar(&genOptions.ApiModelPath, "api-model", "", "Define the API model to use")
	//f.StringVar(&gc.outputDirectory, "output-directory", "", "output directory (derived from FQDN if absent)")
	//f.StringVar(&gc.caCertificatePath, "ca-certificate-path", "", "path to the CA certificate to use for Kubernetes PKI assets")
	//f.StringVar(&gc.caPrivateKeyPath, "ca-private-key-path", "", "path to the CA private key to use for Kubernetes PKI assets")
	//f.BoolVar(&gc.classicMode, "classic-mode", false, "enable classic parameters and outputs")
	//f.BoolVar(&gc.noPrettyPrint, "no-pretty-print", false, "skip pretty printing the output")
	//f.BoolVar(&gc.parametersOnly, "parameters-only", false, "only output parameters files")
	return genAgentpoolCmd
}

func (gc *GenerateOptions) Init(cmd *cobra.Command, args []string) error {

	if gc.ApiModelPath == "" {
		if len(args) > 0 {
			gc.ApiModelPath = args[0]
		} else if len(args) > 1 {
			cmd.Usage()
			return fmt.Errorf("too many arguments were provided to 'generate'")
		} else {
			cmd.Usage()
			return fmt.Errorf("--api-model was not supplied, nor was one specified as a positional argument")
		}
	}

	var err error
	gc.ContainerService, gc.ApiVersion, err = api.LoadContainerServiceFromFile(gc.ApiModelPath)
	if err != nil {
		return fmt.Errorf("error parsing the api model: %v", err)
	}

	return nil
}

func (gc *GenerateOptions) Validate(cmd *cobra.Command, args []string) {
	//var caCertificateBytes []byte
	//var caKeyBytes []byte
	//var err error
	//

	//
	//if _, err := os.Stat(gc.apimodelPath); os.IsNotExist(err) {
	//	log.Fatalf("specified api model does not exist (%s)", gc.apimodelPath)
	//}
	//
	//gc.containerService, gc.apiVersion, err = api.LoadContainerServiceFromFile(gc.apimodelPath)
	//if err != nil {
	//	log.Fatalf("error parsing the api model: %s", err.Error())
	//}
	//
	//// ------------------------------------------------------------------------------------------
	////
	//// TODO (@kris-nova) Here we need to actually validate the API model.
	//// TODO (@kris-nova) Let's code this after we know what the API is going to look like and
	//// TODO (@kris-nova) how it's supposed to behave
	//if gc.apiVersion == kubernetesagentpool.APIVersion {
	//	log.Infof("Bypassing validation for API: [%s]", kubernetesagentpool.APIVersion)
	//	return
	//}
	////
	//// ------------------------------------------------------------------------------------------
	//
	//if gc.outputDirectory == "" {
	//	gc.outputDirectory = path.Join("_output", gc.containerService.Properties.MasterProfile.DNSPrefix)
	//}
	//
	//// consume gc.caCertificatePath and gc.caPrivateKeyPath
	//
	//if (gc.caCertificatePath != "" && gc.caPrivateKeyPath == "") || (gc.caCertificatePath == "" && gc.caPrivateKeyPath != "") {
	//	log.Fatal("--ca-certificate-path and --ca-private-key-path must be specified together")
	//}
	//if gc.caCertificatePath != "" {
	//	if caCertificateBytes, err = ioutil.ReadFile(gc.caCertificatePath); err != nil {
	//		log.Fatal("failed to read CA certificate file:", err)
	//	}
	//	if caKeyBytes, err = ioutil.ReadFile(gc.caPrivateKeyPath); err != nil {
	//		log.Fatal("failed to read CA private key file:", err)
	//	}
	//
	//	prop := gc.containerService.Properties
	//	if prop.CertificateProfile == nil {
	//		prop.CertificateProfile = &api.CertificateProfile{}
	//	}
	//	prop.CertificateProfile.CaCertificate = string(caCertificateBytes)
	//	prop.CertificateProfile.CaPrivateKey = string(caKeyBytes)
	//}
}

func (gc *GenerateOptions) Run() error {
	//log.Infoln("Generating assets...")

	interpolator := agentpool.NewAgentPoolInterpolator(gc.ContainerService)
	err := interpolator.Interpolate()
	if err != nil {
		return fmt.Errorf("Major error on interpolate: %v", err)
	}

	//templateGenerator, err := acsengine.InitializeTemplateGenerator(gc.classicMode)
	//if err != nil {
	//	log.Fatalln("failed to initialize template generator: %s", err.Error())
	//}
	//
	//certsGenerated := false
	//template, parameters, certsGenerated, err := templateGenerator.GenerateTemplate(gc.containerService)
	//if err != nil {
	//	log.Fatalf("error generating template %s: %s", gc.apimodelPath, err.Error())
	//	os.Exit(1)
	//}
	//
	//if !gc.noPrettyPrint {
	//	if template, err = acsengine.PrettyPrintArmTemplate(template); err != nil {
	//		log.Fatalf("error pretty printing template: %s \n", err.Error())
	//	}
	//	if parameters, err = acsengine.BuildAzureParametersFile(parameters); err != nil {
	//		log.Fatalf("error pretty printing template parameters: %s \n", err.Error())
	//	}
	//}
	//
	//if err = acsengine.WriteArtifacts(gc.containerService, gc.apiVersion, template, parameters, gc.outputDirectory, certsGenerated, gc.parametersOnly); err != nil {
	//	log.Fatalf("error writing artifacts: %s \n", err.Error())
	//}

	fmt.Println(interpolator.GetTemplate())

	return nil
}
