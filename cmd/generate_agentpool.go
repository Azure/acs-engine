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

type GenerateOptions struct {
	ApiModelPath string
	AgentPool    *kubernetesagentpool.AgentPool
	ApiVersion   string
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
			return genOptions.Run()
		},
	}

	f := genAgentpoolCmd.Flags()
	f.StringVar(&genOptions.ApiModelPath, "api-model", "", "Define the API model to use")
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
	gc.AgentPool, gc.ApiVersion, err = kubernetesagentpool.LoadAgentPoolFromFile(gc.ApiModelPath)
	if err != nil {
		return fmt.Errorf("error parsing the api model: %v", err)
	}


	return nil
}

func (gc *GenerateOptions) Validate(cmd *cobra.Command, args []string) {
	// todo validate
}

func (gc *GenerateOptions) Run() error {
	//log.Infoln("Generating assets...")

<<<<<<< HEAD
	fmt.Println(gc.ContainerService)
=======
	interpolator := agentpool.NewAgentPoolInterpolator(gc.AgentPool, "kubernetes/agentpool")
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
=======
	iw := interpolatorwriter.NewInterpolatorWriter("./_output", "azuredeploy.json", "azuredeploy.params.json", interpolator)
=======
	iw := interpolatorwriter.NewInterpolatorWriter(fmt.Sprintf("./_output/%s", gc.AgentPool.Name), "azuredeploy.json", "azuredeploy.parameterss.json", interpolator)
>>>>>>> Adding work
	err = iw.Write()
	if err != nil {
		return fmt.Errorf("Unable to write template: %v", err)
	}
>>>>>>> Clean. Simple. Go.
	return nil

}
