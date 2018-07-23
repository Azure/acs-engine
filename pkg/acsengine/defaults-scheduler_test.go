package acsengine

import (
	"testing"
)

func TestSchedulerDefaultConfig(t *testing.T) {
	cs := CreateMockContainerService("testcluster", "1.9.6", 3, 2, false)
	setSchedulerConfig(cs)
	s := cs.Properties.OrchestratorProfile.KubernetesConfig.SchedulerConfig
	for key, val := range staticSchedulerConfig {
		if val != s[key] {
			t.Fatalf("got unexpected kube-scheduler static config value for %s. Expected %s, got %s",
				key, val, s[key])
		}
	}
	for key, val := range defaultSchedulerConfig {
		if val != s[key] {
			t.Fatalf("got unexpected kube-scheduler default config value for %s. Expected %s, got %s",
				key, val, s[key])
		}
	}
}

func TestSchedulerUserConfig(t *testing.T) {
	cs := CreateMockContainerService("testcluster", "1.9.6", 3, 2, false)
	assignmentMap := map[string]string{
		"--scheduler-name": "my-custom-name",
		"--feature-gates":  "APIListChunking=true,APIResponseCompression=true,Accelerators=true,AdvancedAuditing=true",
	}
	cs.Properties.OrchestratorProfile.KubernetesConfig.SchedulerConfig = assignmentMap
	setSchedulerConfig(cs)
	for key, val := range assignmentMap {
		if val != cs.Properties.OrchestratorProfile.KubernetesConfig.SchedulerConfig[key] {
			t.Fatalf("got unexpected kube-scheduler config value for %s. Expected %s, got %s",
				key, val, cs.Properties.OrchestratorProfile.KubernetesConfig.SchedulerConfig[key])
		}
	}
}

func TestSchedulerStaticConfig(t *testing.T) {
	cs := CreateMockContainerService("testcluster", "1.9.6", 3, 2, false)
	cs.Properties.OrchestratorProfile.KubernetesConfig.SchedulerConfig = map[string]string{
		"--kubeconfig":   "user-override",
		"--leader-elect": "user-override",
		"--profiling":    "user-override",
	}
	setSchedulerConfig(cs)
	for key, val := range staticSchedulerConfig {
		if val != cs.Properties.OrchestratorProfile.KubernetesConfig.SchedulerConfig[key] {
			t.Fatalf("kube-scheduler static config did not override user values for %s. Expected %s, got %s",
				key, val, cs.Properties.OrchestratorProfile.KubernetesConfig.SchedulerConfig)
		}
	}
}
