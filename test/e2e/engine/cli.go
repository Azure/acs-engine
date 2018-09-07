package engine

import (
	"log"
	"os/exec"

	"github.com/Azure/acs-engine/test/e2e/kubernetes/util"
)

// Generate will run acs-engine generate on a given cluster definition
func (e *Engine) Generate() error {
	cmd := exec.Command("./bin/acs-engine", "generate", e.Config.ClusterDefinitionTemplate, "--output-directory", e.Config.GeneratedDefinitionPath)
	util.PrintCommand(cmd)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Error while trying to generate acs-engine template with cluster definition - %s: %s\n", e.Config.ClusterDefinitionTemplate, err)
		log.Printf("Command:./bin/acs-engine generate %s --output-directory %s\n", e.Config.ClusterDefinitionTemplate, e.Config.GeneratedDefinitionPath)
		log.Printf("Output:%s\n", out)
		return err
	}
	return nil
}

// Deploy will run acs-engine deploy on a given cluster definition
func (e *Engine) Deploy(location string) error {
	cmd := exec.Command("./bin/acs-engine", "deploy",
		"--location", location,
		"--api-model", e.Config.ClusterDefinitionPath,
		"--dns-prefix", e.Config.DefinitionName,
		"--output-directory", e.Config.GeneratedDefinitionPath,
		"--resource-group", e.Config.DefinitionName,
	)
	util.PrintCommand(cmd)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Error while trying to deploy acs-engine template with cluster definition - %s: %s\n", e.Config.ClusterDefinitionTemplate, err)
		log.Printf("Output:%s\n", out)
		return err
	}
	return nil
}
