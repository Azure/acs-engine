package engine

import (
	"log"
	"os/exec"
)

// Generate will run acs-engine generate on a given cluster definition
func (e *Engine) Generate() error {
	cmd := exec.Command("acs-engine", "generate", e.ClusterDefinitionTemplate, "--output-directory", e.GeneratedDefinitionPath)
	err := cmd.Start()
	if err != nil {
		log.Printf("Error while trying to start generate:%s\n", err)
		return err
	}
	err = cmd.Wait()
	if err != nil {
		log.Printf("Error while trying to generate acs-engine template with cluster definition - %s: %s", e.GeneratedDefinitionPath, err)
		return err
	}
	return nil
}
