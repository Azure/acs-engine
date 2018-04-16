package api

import (
	"strconv"
	"testing"

	"github.com/Azure/acs-engine/pkg/api/agentPoolOnlyApi/v20180331"
	"github.com/Azure/acs-engine/pkg/helpers"
)

func TestConvertV20180331AgentPoolOnlyOrchestratorProfile(t *testing.T) {
	kubernetesVersion := "1.7.9"
	networkPlugin := v20180331.Azure
	serviceCIDR := "10.0.0.0/8"
	dnsServiceIP := "10.0.0.10"
	dockerBridgeSubnet := "172.17.0.1/16"

	// all networkProfile fields are defined
	p := &v20180331.NetworkProfile{
		NetworkPlugin:    networkPlugin,
		ServiceCidr:      serviceCIDR,
		DNSServiceIP:     dnsServiceIP,
		DockerBridgeCidr: dockerBridgeSubnet,
	}

	api := convertV20180331AgentPoolOnlyOrchestratorProfile(kubernetesVersion, p, nil)

	if api.OrchestratorVersion != kubernetesVersion {
		t.Error("error in orchestrator profile kubernetesVersion conversion")
	}

	if api.KubernetesConfig.NetworkPolicy != string(networkPlugin) {
		t.Error("error in orchestrator profile networkPlugin conversion")
	}

	if api.KubernetesConfig.ServiceCIDR != string(serviceCIDR) {
		t.Error("error in orchestrator profile networkPlugin conversion")
	}

	if api.KubernetesConfig.DNSServiceIP != string(dnsServiceIP) {
		t.Error("error in orchestrator profile networkPlugin conversion")
	}

	if api.KubernetesConfig.DockerBridgeSubnet != string(dockerBridgeSubnet) {
		t.Error("error in orchestrator profile networkPlugin conversion")
	}

	// no networkProfile is defined
	p = nil

	api = convertV20180331AgentPoolOnlyOrchestratorProfile(kubernetesVersion, p, nil)

	if api.OrchestratorVersion != kubernetesVersion {
		t.Error("error in orchestrator profile kubernetesVersion conversion")
	}

	if api.KubernetesConfig.NetworkPolicy != "none" {
		t.Error("error in orchestrator profile networkPlugin conversion")
	}

	if api.KubernetesConfig.ServiceCIDR != "10.0.0.0/16" {
		t.Error("error in orchestrator profile networkPlugin conversion")
	}

	if api.KubernetesConfig.DNSServiceIP != "10.0.0.10" {
		t.Error("error in orchestrator profile networkPlugin conversion")
	}

	if api.KubernetesConfig.DockerBridgeSubnet != "172.17.0.1/16" {
		t.Error("error in orchestrator profile networkPlugin conversion")
	}

	// only networkProfile NetworkPlugin fields is defined
	p = &v20180331.NetworkProfile{
		NetworkPlugin: networkPlugin,
	}

	api = convertV20180331AgentPoolOnlyOrchestratorProfile(kubernetesVersion, p, nil)

	if api.OrchestratorVersion != kubernetesVersion {
		t.Error("error in orchestrator profile kubernetesVersion conversion")
	}

	if api.KubernetesConfig.NetworkPolicy != string(networkPlugin) {
		t.Error("error in orchestrator profile networkPlugin conversion")
	}

	if api.KubernetesConfig.ServiceCIDR != "10.0.0.0/16" {
		t.Error("error in orchestrator profile networkPlugin conversion")
	}

	if api.KubernetesConfig.DNSServiceIP != "10.0.0.10" {
		t.Error("error in orchestrator profile networkPlugin conversion")
	}

	if api.KubernetesConfig.DockerBridgeSubnet != "172.17.0.1/16" {
		t.Error("error in orchestrator profile networkPlugin conversion")
	}
}

func TestConvertV20180331AgentPoolOnlyAgentPoolProfile(t *testing.T) {
	maxPods := 25
	availabilityProfile := "availabilityProfile"

	p := &v20180331.AgentPoolProfile{
		MaxPods: maxPods,
	}
	api := convertV20180331AgentPoolOnlyAgentPoolProfile(p, availabilityProfile)

	if api.KubernetesConfig.KubeletConfig["--max-pods"] != strconv.Itoa(maxPods) {
		t.Error("error in agent pool profile max pods conversion")
	}
}

func TestConvertFromV20180331AddonProfile(t *testing.T) {
	addonName := "AddonFoo"
	p := map[string]v20180331.AddonProfile{
		addonName: {
			Enabled: true,
			Config: map[string]string{
				"opt1": "value1",
			},
		},
	}
	api := convertV20180331AgentPoolOnlyAddonProfiles(p)

	if len(api) != 1 {
		t.Error("there has to be one addon")
	}
	if _, ok := api[addonName]; !ok {
		t.Error("addon is not found")
	}
	if api[addonName].Enabled != true {
		t.Error("addon should be enabled")
	}
	v, ok := api[addonName].Config["opt1"]
	if !ok {
		t.Error("Addon config opt1 is not found")
	}
	if v != "value1" {
		t.Error("addon config value does not match")
	}
}

func TestConvertV20170831AgentPoolOnlyOrchestratorProfile_KubernetesConfig(t *testing.T) {
	op := convertV20170831AgentPoolOnlyOrchestratorProfile("1.8.9")
	if op == nil {
		t.Error("OrchestratorProfile expected not to be nil")
	}

	if op.KubernetesConfig == nil {
		t.Error("OrchestratorProfile.KubernetesConfig expected not to be nil")
	}

	if op.KubernetesConfig.EnableRbac == nil || *op.KubernetesConfig.EnableRbac == true {
		t.Error("OrchestratorProfile.KubernetesConfig.EnableRbac expected to be *false")
	}

	if op.KubernetesConfig.EnableSecureKubelet == nil || *op.KubernetesConfig.EnableSecureKubelet == true {
		t.Error("OrchestratorProfile.KubernetesConfig.EnableSecureKubelet expected to be *false")
	}

}

func TestConvertV20180331AgentPoolOnlyKubernetesConfig(t *testing.T) {
	var kc *KubernetesConfig
	kc = convertV20180331AgentPoolOnlyKubernetesConfig(helpers.PointerToBool(true))
	if kc == nil {
		t.Error("kubernetesConfig expected not to be nil")
	}

	if kc.EnableRbac == nil {
		t.Error("EnableRbac expected not to be nil")
	}

	if *kc.EnableRbac != true {
		t.Error("EnableRbac expected to be true")
	}

	if kc.EnableSecureKubelet == nil {
		t.Error("EnableSecureKubelet expected not to be nil")
	}

	if *kc.EnableSecureKubelet != true {
		t.Error("EnableSecureKubelet expected to be true")
	}

	if *kc.EnableSecureKubelet != *kc.EnableRbac {
		t.Error("EnableSecureKubelet and EnableRbac expected to be same")
	}

	kc = convertV20180331AgentPoolOnlyKubernetesConfig(helpers.PointerToBool(false))
	if kc == nil {
		t.Error("kubernetesConfig expected not to be nil")
	}

	if kc.EnableRbac == nil {
		t.Error("EnableRbac expected not to be nil")
	}

	if *kc.EnableRbac != false {
		t.Error("EnableRbac expected to be false")
	}

	if kc.EnableSecureKubelet == nil {
		t.Error("EnableSecureKubelet expected not to be nil")
	}

	if *kc.EnableSecureKubelet != false {
		t.Error("EnableSecureKubelet expected to be false")
	}

	if *kc.EnableSecureKubelet != *kc.EnableRbac {
		t.Error("EnableSecureKubelet and EnableRbac expected to be same")
	}

	kc = convertV20180331AgentPoolOnlyKubernetesConfig(nil)
	if kc == nil {
		t.Error("kubernetesConfig expected not to be nil")
	}

	if kc.EnableRbac == nil {
		t.Error("EnableRbac expected not to be nil")
	}

	if *kc.EnableRbac != true {
		t.Error("EnableRbac expected to be true")
	}

	if kc.EnableSecureKubelet == nil {
		t.Error("EnableSecureKubelet expected not to be nil")
	}

	if *kc.EnableSecureKubelet != true {
		t.Error("EnableSecureKubelet expected to be true")
	}

	if *kc.EnableSecureKubelet != *kc.EnableRbac {
		t.Error("EnableSecureKubelet and EnableRbac expected to be same")
	}

}

func TestIfMasterProfileIsMissingThenApiModelIsAgentPoolOnly(t *testing.T) {
	json := `
	{
		"apiVersion": "vlabs",
		"properties": {
			"dnsPrefix": "dp",
			"fqdn": "fqdn",
			"agentPoolProfiles": [],
			"servicePrincipalProfile": {}
		}
	}
	`
	isAgentPool := isAgentPoolOnlyClusterJSON([]byte(json))
	if !isAgentPool {
		t.Error("Expected JSON without masterProfile to be interpreted as agent pool, but it was not")
	}
}

func TestIfMasterProfileIsPresentThenApiModelIsFullCluster(t *testing.T) {
	json := `
	{
		"apiVersion": "vlabs",
		"properties": {
			"orchestratorProfile": {},
			"masterProfile": {},
			"agentPoolProfiles": [],
			"servicePrincipalProfile": {}
		}
	}
	`
	isAgentPool := isAgentPoolOnlyClusterJSON([]byte(json))
	if isAgentPool {
		t.Error("Expected JSON with masterProfile not to be interpreted as agent pool, but it was")
	}
}
