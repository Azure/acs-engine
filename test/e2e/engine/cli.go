package engine

import (
	"log"
	"os/exec"
)

// Generate will run acs-engine generate on a given cluster definition
func (e *Engine) Generate() error {
	out, err := exec.Command("./bin/acs-engine", "generate", e.Config.ClusterDefinitionTemplate, "--output-directory", e.Config.GeneratedDefinitionPath).CombinedOutput()
	if err != nil {
		log.Printf("Error while trying to generate acs-engine template with cluster definition - %s: %s\n", e.Config.ClusterDefinitionTemplate, err)
		log.Printf("Command:./bin/acs-engine generate %s --output-directory %s\n", e.Config.ClusterDefinitionTemplate, e.Config.GeneratedDefinitionPath)
		log.Printf("Output:%s\n", out)
		return err
	}
	return nil
}
