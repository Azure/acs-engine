package acsengine

import (
	"github.com/Azure/acs-engine/pkg/api"
)

// Default debug Config
var defaultDebugConfig = map[string]string{
	"waitForNodes": "false",
}

func setKubernetesDebugConfig(cs *api.ContainerService) {
	o := cs.Properties.OrchestratorProfile

	if o.KubernetesConfig.Debug == nil {
		o.KubernetesConfig.Debug = defaultDebugConfig
	} else {
		for key, val := range defaultDebugConfig {
			// If we don't have a user-configurable debug config for each option
			if _, ok := o.KubernetesConfig.Debug[key]; !ok {
				// then assign the default value
				o.KubernetesConfig.Debug[key] = val
			}
		}
	}
}
