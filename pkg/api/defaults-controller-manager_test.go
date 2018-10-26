package api

import (
	"testing"

	"github.com/Azure/acs-engine/pkg/helpers"
)

func TestControllerManagerConfigEnableRbac(t *testing.T) {
	// Test EnableRbac = true
	cs := CreateMockContainerService("testcluster", defaultTestClusterVer, 3, 2, false)
	cs.Properties.OrchestratorProfile.KubernetesConfig.EnableRbac = helpers.PointerToBool(true)
	cs.setControllerManagerConfig()
	cm := cs.Properties.OrchestratorProfile.KubernetesConfig.ControllerManagerConfig
	if cm["--use-service-account-credentials"] != "true" {
		t.Fatalf("got unexpected '--use-service-account-credentials' Controller Manager config value for EnableRbac=true: %s",
			cm["--use-service-account-credentials"])
	}

	// Test default
	cs = CreateMockContainerService("testcluster", defaultTestClusterVer, 3, 2, false)
	cs.Properties.OrchestratorProfile.KubernetesConfig.EnableRbac = helpers.PointerToBool(false)
	cs.setControllerManagerConfig()
	cm = cs.Properties.OrchestratorProfile.KubernetesConfig.ControllerManagerConfig
	if cm["--use-service-account-credentials"] != DefaultKubernetesCtrlMgrUseSvcAccountCreds {
		t.Fatalf("got unexpected '--use-service-account-credentials' Controller Manager config value for EnableRbac=false: %s",
			cm["--use-service-account-credentials"])
	}

}
func TestControllerManagerConfigEnableProfiling(t *testing.T) {
	// Test
	// "controllerManagerConfig": {
	// 	"--profiling": "true"
	// },
	cs := CreateMockContainerService("testcluster", defaultTestClusterVer, 3, 2, false)
	cs.Properties.OrchestratorProfile.KubernetesConfig.ControllerManagerConfig = map[string]string{
		"--profiling": "true",
	}
	cs.setControllerManagerConfig()
	cm := cs.Properties.OrchestratorProfile.KubernetesConfig.ControllerManagerConfig
	if cm["--profiling"] != "true" {
		t.Fatalf("got unexpected '--profiling' Controller Manager config value for \"--profiling\": \"true\": %s",
			cm["--profiling"])
	}

	// Test default
	cs = CreateMockContainerService("testcluster", defaultTestClusterVer, 3, 2, false)
	cs.setControllerManagerConfig()
	cm = cs.Properties.OrchestratorProfile.KubernetesConfig.ControllerManagerConfig
	if cm["--profiling"] != DefaultKubernetesCtrMgrEnableProfiling {
		t.Fatalf("got unexpected default value for '--profiling' Controller Manager config: %s",
			cm["--profiling"])
	}
}
