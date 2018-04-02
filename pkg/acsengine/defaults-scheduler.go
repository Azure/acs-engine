package acsengine

import (
	"github.com/Azure/acs-engine/pkg/api"
)

// staticLinuxSchedulerConfig is not user-overridable
var staticLinuxSchedulerConfig = map[string]string{
	"--kubeconfig":   "/var/lib/kubelet/kubeconfig",
	"--leader-elect": "true",
	"--profiling":    "false",
}

// defaultSchedulerConfig provides targeted defaults, but is user-overridable
var defaultSchedulerConfig = map[string]string{
	"--v": "2",
}

func setSchedulerConfig(cs *api.ContainerService) {
	o := cs.Properties.OrchestratorProfile
	staticWindowsSchedulerConfig := make(map[string]string)
	for key, val := range staticLinuxSchedulerConfig {
		staticWindowsSchedulerConfig[key] = val
	}
	// Windows scheduler config overrides
	// TODO placeholder for specific config overrides for Windows clusters

	// If no user-configurable scheduler config values exists, use the defaults
	if o.KubernetesConfig.SchedulerConfig == nil {
		o.KubernetesConfig.SchedulerConfig = defaultSchedulerConfig
	} else {
		for key, val := range defaultSchedulerConfig {
			// If we don't have a user-configurable scheduler config for each option
			if _, ok := o.KubernetesConfig.SchedulerConfig[key]; !ok {
				// then assign the default value
				o.KubernetesConfig.SchedulerConfig[key] = val
			}
		}
	}

	// We don't support user-configurable values for the following,
	// so any of the value assignments below will override user-provided values
	var overrideSchedulerConfig map[string]string
	if cs.Properties.HasWindows() {
		overrideSchedulerConfig = staticWindowsSchedulerConfig
	} else {
		overrideSchedulerConfig = staticLinuxSchedulerConfig
	}
	for key, val := range overrideSchedulerConfig {
		o.KubernetesConfig.SchedulerConfig[key] = val
	}
}
