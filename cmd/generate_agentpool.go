package cmd

import (
	"fmt"
	"github.com/Azure/acs-engine/pkg/api/kubernetesagentpool"
	"github.com/Azure/acs-engine/pkg/interpolator/agentpool"
	"github.com/Azure/acs-engine/pkg/interpolatorwriter"
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
			err = genOptions.Run()
			if err != nil {
				log.Fatal(err)
			}
			return nil
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

	interpolator := agentpool.NewAgentPoolInterpolator(gc.AgentPool, "kubernetes/agentpool")
	err := interpolator.Interpolate()
	if err != nil {
		return fmt.Errorf("Major error on interpolate: %v", err)
	}

	iw := interpolatorwriter.NewInterpolatorWriter(fmt.Sprintf("./_output/%s", gc.AgentPool.Name), "azuredeploy.json", "azuredeploy.parameterss.json", interpolator)
	err = iw.Write()
	if err != nil {
		return fmt.Errorf("Unable to write template: %v", err)
	}
	return nil

}
