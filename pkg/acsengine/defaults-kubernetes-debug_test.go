package acsengine

import (
	"testing"
)

var userDebugConfig = map[string]string{
	"waitForNodes": "true",
}

func TestKubernetesDebugDefaults(t *testing.T) {
	cs := createContainerService("testcluster", "1.8.6", 3, 2)
	setKubernetesDebugConfig(cs)
	d := cs.Properties.OrchestratorProfile.KubernetesConfig.Debug
	for key, val := range defaultDebugConfig {
		if d[key] != val {
			t.Fatalf("got unexpected kubernetes debug config value for %s: %s, expected %s",
				key, d[key], val)
		}
	}
}

func TestKubernetesDebug(t *testing.T) {
	cs := createContainerService("testcluster", "1.8.6", 3, 2)
	cs.Properties.OrchestratorProfile.KubernetesConfig.Debug = userDebugConfig
	setKubernetesDebugConfig(cs)
	d := cs.Properties.OrchestratorProfile.KubernetesConfig.Debug
	for key, val := range userDebugConfig {
		if d[key] != val {
			t.Fatalf("got unexpected kubernetes debug config value for %s: %s, expected %s",
				key, d[key], val)
		}
	}
}
