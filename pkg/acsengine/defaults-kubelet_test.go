package acsengine

import (
	"strconv"
	"testing"

	"github.com/Azure/acs-engine/pkg/helpers"
)

func TestKubeletConfigDefaults(t *testing.T) {
	cs := createContainerService("testcluster", "1.8.6", 3, 2)
	setKubeletConfig(cs)
	k := cs.Properties.OrchestratorProfile.KubernetesConfig.KubeletConfig
	// TODO test all default config values
	for key, val := range map[string]string{"--azure-container-registry-config": "/etc/kubernetes/azure.json"} {
		if k[key] != val {
			t.Fatalf("got unexpected kubelet config value for %s: %s, expected %s",
				key, k[key], val)
		}
	}

	cs = createContainerService("testcluster", "1.8.6", 3, 2)
	// TODO test all default overrides
	overrideVal := "/etc/override"
	cs.Properties.OrchestratorProfile.KubernetesConfig.KubeletConfig = map[string]string{
		"--azure-container-registry-config": overrideVal,
	}
	setKubeletConfig(cs)
	k = cs.Properties.OrchestratorProfile.KubernetesConfig.KubeletConfig
	for key, val := range map[string]string{"--azure-container-registry-config": overrideVal} {
		if k[key] != val {
			t.Fatalf("got unexpected kubelet config value for %s: %s, expected %s",
				key, k[key], val)
		}
	}
}

func TestKubeletConfigUseCloudControllerManager(t *testing.T) {
	// Test UseCloudControllerManager = true
	cs := createContainerService("testcluster", defaultTestClusterVer, 3, 2)
	cs.Properties.OrchestratorProfile.KubernetesConfig.UseCloudControllerManager = helpers.PointerToBool(true)
	setKubeletConfig(cs)
	k := cs.Properties.OrchestratorProfile.KubernetesConfig.KubeletConfig
	if k["--cloud-provider"] != "external" {
		t.Fatalf("got unexpected '--cloud-provider' kubelet config value for UseCloudControllerManager=true: %s",
			k["--cloud-provider"])
	}

	// Test UseCloudControllerManager = false
	cs = createContainerService("testcluster", defaultTestClusterVer, 3, 2)
	cs.Properties.OrchestratorProfile.KubernetesConfig.UseCloudControllerManager = helpers.PointerToBool(false)
	setKubeletConfig(cs)
	k = cs.Properties.OrchestratorProfile.KubernetesConfig.KubeletConfig
	if k["--cloud-provider"] != "azure" {
		t.Fatalf("got unexpected '--cloud-provider' kubelet config value for UseCloudControllerManager=false: %s",
			k["--cloud-provider"])
	}

}

func TestKubeletConfigCloudConfig(t *testing.T) {
	// Test default value and custom value for --cloud-config
	cs := createContainerService("testcluster", defaultTestClusterVer, 3, 2)
	setKubeletConfig(cs)
	k := cs.Properties.OrchestratorProfile.KubernetesConfig.KubeletConfig
	if k["--cloud-config"] != "/etc/kubernetes/azure.json" {
		t.Fatalf("got unexpected '--cloud-config' kubelet config default value: %s",
			k["--cloud-config"])
	}

	cs = createContainerService("testcluster", defaultTestClusterVer, 3, 2)
	cs.Properties.OrchestratorProfile.KubernetesConfig.KubeletConfig["--cloud-config"] = "custom.json"
	setKubeletConfig(cs)
	k = cs.Properties.OrchestratorProfile.KubernetesConfig.KubeletConfig
	if k["--cloud-config"] != "custom.json" {
		t.Fatalf("got unexpected '--cloud-config' kubelet config default value: %s",
			k["--cloud-config"])
	}
}

func TestKubeletConfigAzureContainerRegistryCofig(t *testing.T) {
	// Test default value and custom value for --azure-container-registry-config
	cs := createContainerService("testcluster", defaultTestClusterVer, 3, 2)
	setKubeletConfig(cs)
	k := cs.Properties.OrchestratorProfile.KubernetesConfig.KubeletConfig
	if k["--azure-container-registry-config"] != "/etc/kubernetes/azure.json" {
		t.Fatalf("got unexpected '--azure-container-registry-config' kubelet config default value: %s",
			k["--azure-container-registry-config"])
	}

	cs = createContainerService("testcluster", defaultTestClusterVer, 3, 2)
	cs.Properties.OrchestratorProfile.KubernetesConfig.KubeletConfig["--azure-container-registry-config"] = "custom.json"
	setKubeletConfig(cs)
	k = cs.Properties.OrchestratorProfile.KubernetesConfig.KubeletConfig
	if k["--azure-container-registry-config"] != "custom.json" {
		t.Fatalf("got unexpected '--azure-container-registry-config' kubelet config default value: %s",
			k["--azure-container-registry-config"])
	}
}

func TestKubeletConfigNetworkPolicy(t *testing.T) {
	// Test NetworkPolicy = none
	cs := createContainerService("testcluster", defaultTestClusterVer, 3, 2)
	cs.Properties.OrchestratorProfile.KubernetesConfig.NetworkPolicy = NetworkPolicyNone
	setKubeletConfig(cs)
	k := cs.Properties.OrchestratorProfile.KubernetesConfig.KubeletConfig
	if k["--network-plugin"] != NetworkPluginKubenet {
		t.Fatalf("got unexpected '--network-plugin' kubelet config value for NetworkPolicy=none: %s",
			k["--network-plugin"])
	}

	// Test NetworkPolicy = azure
	cs = createContainerService("testcluster", defaultTestClusterVer, 3, 2)
	cs.Properties.OrchestratorProfile.KubernetesConfig.NetworkPolicy = "azure"
	setKubeletConfig(cs)
	k = cs.Properties.OrchestratorProfile.KubernetesConfig.KubeletConfig
	if k["--network-plugin"] != "cni" {
		t.Fatalf("got unexpected '--network-plugin' kubelet config value for NetworkPolicy=azure: %s",
			k["--network-plugin"])
	}

}

func TestKubeletConfigEnableSecureKubelet(t *testing.T) {
	// Test EnableSecureKubelet = true
	cs := createContainerService("testcluster", defaultTestClusterVer, 3, 2)
	cs.Properties.OrchestratorProfile.KubernetesConfig.EnableSecureKubelet = helpers.PointerToBool(true)
	setKubeletConfig(cs)
	k := cs.Properties.OrchestratorProfile.KubernetesConfig.KubeletConfig
	if k["--anonymous-auth"] != "false" {
		t.Fatalf("got unexpected '--anonymous-auth' kubelet config value for EnableSecureKubelet=true: %s",
			k["--anonymous-auth"])
	}
	if k["--authorization-mode"] != "Webhook" {
		t.Fatalf("got unexpected '--authorization-mode' kubelet config value for EnableSecureKubelet=true: %s",
			k["--authorization-mode"])
	}
	if k["--client-ca-file"] != "/etc/kubernetes/certs/ca.crt" {
		t.Fatalf("got unexpected '--client-ca-file' kubelet config value for EnableSecureKubelet=true: %s",
			k["--client-ca-file"])
	}

	// Test EnableSecureKubelet = false
	cs = createContainerService("testcluster", defaultTestClusterVer, 3, 2)
	cs.Properties.OrchestratorProfile.KubernetesConfig.EnableSecureKubelet = helpers.PointerToBool(false)
	setKubeletConfig(cs)
	k = cs.Properties.OrchestratorProfile.KubernetesConfig.KubeletConfig
	for _, key := range []string{"--anonymous-auth", "--client-ca-file"} {
		if _, ok := k[key]; ok {
			t.Fatalf("got unexpected '%s' kubelet config value for EnableSecureKubelet=false: %s",
				key, k[key])
		}
	}

}

func TestKubeletMaxPods(t *testing.T) {
	cs := createContainerService("testcluster", defaultTestClusterVer, 3, 2)
	cs.Properties.OrchestratorProfile.KubernetesConfig.NetworkPolicy = NetworkPolicyAzure
	setKubeletConfig(cs)
	k := cs.Properties.OrchestratorProfile.KubernetesConfig.KubeletConfig
	if k["--max-pods"] != strconv.Itoa(DefaultKubernetesMaxPodsVNETIntegrated) {
		t.Fatalf("got unexpected '--max-pods' kubelet config value for NetworkPolicy=%s: %s",
			NetworkPolicyAzure, k["--max-pods"])
	}

	cs = createContainerService("testcluster", defaultTestClusterVer, 3, 2)
	cs.Properties.OrchestratorProfile.KubernetesConfig.NetworkPolicy = NetworkPolicyNone
	setKubeletConfig(cs)
	k = cs.Properties.OrchestratorProfile.KubernetesConfig.KubeletConfig
	if k["--max-pods"] != strconv.Itoa(DefaultKubernetesMaxPods) {
		t.Fatalf("got unexpected '--max-pods' kubelet config value for NetworkPolicy=%s: %s",
			NetworkPolicyNone, k["--max-pods"])
	}
}
