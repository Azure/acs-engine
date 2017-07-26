package cmd

import (
	"fmt"
<<<<<<< HEAD
	"github.com/Azure/acs-engine/pkg/api"
=======
	"github.com/Azure/acs-engine/pkg/api/kubernetesagentpool"
	"github.com/Azure/acs-engine/pkg/interpolator/agentpool"
<<<<<<< HEAD
>>>>>>> Refactor into agentpool instead of container service
=======
	"github.com/Azure/acs-engine/pkg/interpolatorwriter"
>>>>>>> Clean. Simple. Go.
	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
)

// GenerateOptions defines the options a user can define to work with the Agent Pool API
type GenerateOptions struct {
	ApiModelPath string // The path on the local filesystem where the input object is
	agentPool    *kubernetesagentpool.AgentPool
	apiVersion   string
}



//    LoadHostedControlPlane is a method in ACS Engine
//      api.LoadHostedControlPlane is a helper function in ACS Engine to convert versioned API model to unversioned internal
//      representation.
//      rawBody: is the versioned HostedControl API model received from the caller: github.com\Azure\acs-engine\pkg\agentPoolOnlyApi\v20170831\types.go
//    unversionedHostedControlPlane: is the unversioned internal model saveb  by ACS RP and used to generate ARM template
//unversionedHostedControlPlane, err := api.LoadHostedControlPlane(rawBody, apiVersion) // apiVersion = "2017-08-31"

// InitializeTemplateGenerator and GenerateAgentPoolTemplate are methods in ACS Engine to generate agent pool template
//templateGenerator, e := acsengine.InitializeTemplateGenerator(false /*classicMode*/)
//templateJSON, parametersJSON, _, e = templateGenerator.GenerateAgentPoolTemplate(unversionedHostedControlPlane)

// NewGenerateAgentpoolCmd will create a new Agent Pool cobra command
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
			return genOptions.Run()
		},
	}

	f := genAgentpoolCmd.Flags()
	f.StringVar(&genOptions.ApiModelPath, "api-model", "", "Define the API model to use")
	return genAgentpoolCmd
}

// Init will initialize the GenerateOptions struct, and calculate runtime configuration
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
	gc.agentPool, gc.apiVersion, err = kubernetesagentpool.LoadAgentPoolFromFile(gc.ApiModelPath)
	if err != nil {
		return fmt.Errorf("error parsing the api model: %v", err)
	}


	return nil
}

// Validate will validate that the input object is sane and valid
func (gc *GenerateOptions) Validate(cmd *cobra.Command, args []string) {
	// TODO (@kris-nova) We need to figure out what we want to validate on and code it here
}

// Run will run the GenerateOptions struct. This will interpolate the ARM template, and atomically commit to disk
func (gc *GenerateOptions) Run() error {
<<<<<<< HEAD
	//log.Infoln("Generating assets...")

<<<<<<< HEAD
	fmt.Println(gc.ContainerService)
=======
	interpolator := agentpool.NewAgentPoolInterpolator(gc.AgentPool, "kubernetes/agentpool")
=======
	interpolator := agentpool.NewAgentPoolInterpolator(gc.agentPool, "kubernetes/agentpool")
>>>>>>> Docs, docs, docs
	err := interpolator.Interpolate()
	if err != nil {
		return fmt.Errorf("Major error on interpolate: %v", err)
	}
>>>>>>> Refactor into agentpool instead of container service

<<<<<<< HEAD
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

<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
=======
	iw := interpolatorwriter.NewInterpolatorWriter("./_output", "azuredeploy.json", "azuredeploy.params.json", interpolator)
=======
	iw := interpolatorwriter.NewInterpolatorWriter(fmt.Sprintf("./_output/%s", gc.AgentPool.Name), "azuredeploy.json", "azuredeploy.parameterss.json", interpolator)
>>>>>>> Adding work
=======
	iw := interpolatorwriter.NewInterpolatorWriter(fmt.Sprintf("./_output/%s", gc.agentPool.Name), "azuredeploy.json", "azuredeploy.parameterss.json", interpolator)
>>>>>>> Docs, docs, docs
=======
	iw := interpolatorwriter.NewInterpolatorWriter(fmt.Sprintf("./_output/%s", gc.agentPool.Name), "azuredeploy.json", "azuredeploy.parameters.json", interpolator)
>>>>>>> Validation is now passing - just dialing in the arm template
	err = iw.Write()
	if err != nil {
		return fmt.Errorf("Unable to write template: %v", err)
	}
>>>>>>> Clean. Simple. Go.
	return nil
}
