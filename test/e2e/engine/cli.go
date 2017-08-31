package engine

import (
	"log"
	"os/exec"
)

// Generate will run acs-engine generate on a given cluster definition
func (e *Engine) Generate() error {
	out, err := exec.Command("./bin/acs-engine", "generate", e.ClusterDefinitionTemplate, "--output-directory", e.GeneratedDefinitionPath).CombinedOutput()
	if err != nil {
		log.Printf("Error while trying to generate acs-engine template with cluster definition - %s: %s\n", e.ClusterDefinitionTemplate, err)
		log.Printf("Command:./bin/acs-engine generate %s --output-directory %s\n", e.ClusterDefinitionTemplate, e.GeneratedDefinitionPath)
		log.Printf("Output:%s\n", out)
		return err
	}
	return nil
}
