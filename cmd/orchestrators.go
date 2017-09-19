package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/Azure/acs-engine/pkg/api"
	"github.com/spf13/cobra"
)

const (
	cmdName             = "orchestrators"
	cmdShortDescription = "provide info about supported orchestrators"
	cmdLongDescription  = "provide info about versions of supported orchestrators"
)

type orchestratorsCmd struct {
	// user input
	orchestrator string
	release      string
}

func newOrchestratorsCmd() *cobra.Command {
	oc := orchestratorsCmd{}

	command := &cobra.Command{
		Use:   cmdName,
		Short: cmdShortDescription,
		Long:  cmdLongDescription,
		RunE: func(cmd *cobra.Command, args []string) error {
			return oc.run(cmd, args)
		},
	}

	f := command.Flags()
	f.StringVar(&oc.orchestrator, "orchestrator", "", "orchestrator name (optional) ")
	f.StringVar(&oc.release, "release", "", "orchestrator release (optional)")

	return command
}

func (oc *orchestratorsCmd) run(cmd *cobra.Command, args []string) error {
	orchs, err := api.GetOrchestratorVersionProfileList(oc.orchestrator, oc.release)
	if err != nil {
		return err
	}
	data, err := json.MarshalIndent(orchs, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(data))
	return nil
}
