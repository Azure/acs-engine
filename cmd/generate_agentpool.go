package cmd

import (
	"fmt"
	"github.com/Azure/acs-engine/pkg/api/kubernetesagentpool"
	"github.com/Azure/acs-engine/pkg/interpolator/agentpool"
	"github.com/Azure/acs-engine/pkg/interpolatorwriter"
	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
)

// GenerateOptions defines the options a user can define to work with the Agent Pool API
type GenerateOptions struct {
	APIModelPath string // The path on the local filesystem where the input object is
	agentPool    *kubernetesagentpool.AgentPool
	apiVersion   string
}

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
			err = genOptions.Validate(cmd, args)
			if err != nil {
				log.Fatal(err)
			}
			err = genOptions.Run()
			if err != nil {
				log.Fatal(err)
			}
			return nil
		},
	}

	f := genAgentpoolCmd.Flags()
	f.StringVar(&genOptions.APIModelPath, "api-model", "", "Define the API model to use")
	return genAgentpoolCmd
}

// Init will initialize the GenerateOptions struct, and calculate runtime configuration
func (gc *GenerateOptions) Init(cmd *cobra.Command, args []string) error {

	if gc.APIModelPath == "" {
		if len(args) > 0 {
			gc.APIModelPath = args[0]
		} else if len(args) > 1 {
			cmd.Usage()
			return fmt.Errorf("too many arguments were provided to 'generate'")
		} else {
			cmd.Usage()
			return fmt.Errorf("--api-model was not supplied, nor was one specified as a positional argument")
		}
	}

	var err error
	gc.agentPool, gc.apiVersion, err = kubernetesagentpool.LoadAgentPoolFromFile(gc.APIModelPath)
	if err != nil {
		return fmt.Errorf("error parsing the api model: %v", err)
	}

	return nil
}

// Validate will validate that the input object is sane and valid
func (gc *GenerateOptions) Validate(cmd *cobra.Command, args []string) error {
	return gc.agentPool.Validate()
}

// Run will run the GenerateOptions struct. This will interpolate the ARM template, and atomically commit to disk
func (gc *GenerateOptions) Run() error {
	interpolator := agentpool.NewAgentPoolInterpolator(gc.agentPool, "kubernetes/agentpool")
	err := interpolator.Interpolate()
	if err != nil {
		return fmt.Errorf("Major error on interpolate: %v", err)
	}

	iw := interpolatorwriter.NewInterpolatorWriter(fmt.Sprintf("./_output/%s", gc.agentPool.Name), "azuredeploy.json", "azuredeploy.parameters.json", interpolator)
	err = iw.Write()
	if err != nil {
		return fmt.Errorf("Unable to write template: %v", err)
	}
	return nil
}
