package api

import (
	"strconv"
	"testing"

	"github.com/Azure/acs-engine/pkg/api/agentPoolOnlyApi/v20180331"
	"github.com/Azure/acs-engine/pkg/helpers"
)

func TestConvertOrchestratorProfileToV20180331AgentPoolOnly(t *testing.T) {
	orchestratorVersion := "1.7.9"
	networkPolicy := "azure"
	serviceCIDR := "10.0.0.0/8"
	dnsServiceIP := "10.0.0.10"
	dockerBridgeSubnet := "172.17.0.1/16"

	kubernetesConfig := &KubernetesConfig{
		NetworkPolicy:      networkPolicy,
		ServiceCIDR:        serviceCIDR,
		DNSServiceIP:       dnsServiceIP,
		DockerBridgeSubnet: dockerBridgeSubnet,
	}
	api := &OrchestratorProfile{
		OrchestratorVersion: orchestratorVersion,
		KubernetesConfig:    kubernetesConfig,
	}

	var version string
	p := &v20180331.NetworkProfile{}
	convertOrchestratorProfileToV20180331AgentPoolOnly(api, &version, p)

	if version != orchestratorVersion {
		t.Error("error in orchestrator profile orchestratorVersion conversion")
	}

	if string(p.NetworkPlugin) != networkPolicy {
		t.Error("error in orchestrator profile networkPlugin conversion")
	}

	if p.ServiceCidr != serviceCIDR {
		t.Error("error in orchestrator profile serviceCidr conversion")
	}

	if p.DNSServiceIP != dnsServiceIP {
		t.Error("error in orchestrator profile dnsServiceIP conversion")
	}

	if p.DockerBridgeCidr != dockerBridgeSubnet {
		t.Error("error in orchestrator profile dockerBridgeCidr conversion")
	}
}

func TestConvertAgentPoolProfileToV20180331AgentPoolOnly(t *testing.T) {
	maxPods := 25

	kubernetesConfig := &KubernetesConfig{
		KubeletConfig: map[string]string{"--max-pods": strconv.Itoa(maxPods)},
	}
	api := &AgentPoolProfile{
		KubernetesConfig: kubernetesConfig,
	}

	p := &v20180331.AgentPoolProfile{}
	convertAgentPoolProfileToV20180331AgentPoolOnly(api, p)

	if p.MaxPods != maxPods {
		t.Error("error in agent pool profile max pods conversion")
	}
}

func TestConvertToV20180331AddonProfile(t *testing.T) {
	addonName := "AddonFoo"
	api := map[string]AddonProfile{
		addonName: {
			Enabled: true,
			Config: map[string]string{
				"opt1": "value1",
			},
		},
	}

	p := make(map[string]v20180331.AddonProfile)
	convertAddonsProfileToV20180331AgentPoolOnly(api, p)

	if len(p) != 1 {
		t.Error("there has to be one addon")
	}
	if _, ok := p[addonName]; !ok {
		t.Error("addon is not found")
	}
	if p[addonName].Enabled != true {
		t.Error("addon should be enabled")
	}
	v, ok := p[addonName].Config["opt1"]
	if !ok {
		t.Error("Addon config opt1 is not found")
	}
	if v != "value1" {
		t.Error("addon config value does not match")
	}
}

func TestConvertKubernetesConfigToEnableRBACV20180331AgentPoolOnly(t *testing.T) {
	var kc *KubernetesConfig
	kc = nil
	enableRBAC := convertKubernetesConfigToEnableRBACV20180331AgentPoolOnly(kc)
	if enableRBAC == nil {
		t.Error("EnableRBAC expected not to be nil")
	}
	if *enableRBAC {
		t.Error("EnableRBAC expected to be false")
	}

	kc = &KubernetesConfig{
		EnableRbac:          nil,
		EnableSecureKubelet: helpers.PointerToBool(true),
	}
	enableRBAC = convertKubernetesConfigToEnableRBACV20180331AgentPoolOnly(kc)
	if enableRBAC == nil {
		t.Error("EnableRBAC expected not to be nil")
	}
	if *enableRBAC {
		t.Error("EnableRBAC expected to be false")
	}

	kc = &KubernetesConfig{
		EnableRbac:          helpers.PointerToBool(false),
		EnableSecureKubelet: helpers.PointerToBool(true),
	}
	enableRBAC = convertKubernetesConfigToEnableRBACV20180331AgentPoolOnly(kc)
	if enableRBAC == nil {
		t.Error("EnableRBAC expected not to be nil")
	}
	if *enableRBAC {
		t.Error("EnableRBAC expected to be false")
	}

	kc = &KubernetesConfig{
		EnableRbac:          helpers.PointerToBool(false),
		EnableSecureKubelet: helpers.PointerToBool(false),
	}
	enableRBAC = convertKubernetesConfigToEnableRBACV20180331AgentPoolOnly(kc)
	if enableRBAC == nil {
		t.Error("EnableRBAC expected not to be nil")
	}
	if *enableRBAC {
		t.Error("EnableRBAC expected to be false")
	}

	kc = &KubernetesConfig{
		EnableRbac:          helpers.PointerToBool(true),
		EnableSecureKubelet: helpers.PointerToBool(true),
	}
	enableRBAC = convertKubernetesConfigToEnableRBACV20180331AgentPoolOnly(kc)
	if enableRBAC == nil {
		t.Error("EnableRBAC expected not to be nil")
	}
	if !*enableRBAC {
		t.Error("EnableRBAC expected to be true")
	}

	kc = &KubernetesConfig{
		EnableRbac:          helpers.PointerToBool(true),
		EnableSecureKubelet: helpers.PointerToBool(false),
	}
	enableRBAC = convertKubernetesConfigToEnableRBACV20180331AgentPoolOnly(kc)
	if enableRBAC == nil {
		t.Error("EnableRBAC expected not to be nil")
	}
	if !*enableRBAC {
		t.Error("EnableRBAC expected to be true")
	}

	kc = &KubernetesConfig{
		EnableRbac:          helpers.PointerToBool(true),
		EnableSecureKubelet: nil,
	}
	enableRBAC = convertKubernetesConfigToEnableRBACV20180331AgentPoolOnly(kc)
	if enableRBAC == nil {
		t.Error("EnableRBAC expected not to be nil")
	}
	if !*enableRBAC {
		t.Error("EnableRBAC expected to be true")
	}

}
