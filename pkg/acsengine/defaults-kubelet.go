package acsengine

import (
	"strconv"

	"github.com/Azure/acs-engine/pkg/api"
	"github.com/Azure/acs-engine/pkg/api/common"
	"github.com/Azure/acs-engine/pkg/helpers"
)

func setKubeletConfig(cs *api.ContainerService) {
	o := cs.Properties.OrchestratorProfile
	cloudSpecConfig := GetCloudSpecConfig(cs.Location)
	staticLinuxKubeletConfig := map[string]string{
		"--address":                     "0.0.0.0",
		"--allow-privileged":            "true",
		"--anonymous-auth":              "false",
		"--authorization-mode":          "Webhook",
		"--client-ca-file":              "/etc/kubernetes/certs/ca.crt",
		"--pod-manifest-path":           "/etc/kubernetes/manifests",
		"--cluster-dns":                 o.KubernetesConfig.DNSServiceIP,
		"--cgroups-per-qos":             "true",
		"--enforce-node-allocatable":    "pods",
		"--kubeconfig":                  "/var/lib/kubelet/kubeconfig",
		"--keep-terminated-pod-volumes": "false",
	}

	staticWindowsKubeletConfig := make(map[string]string)
	for key, val := range staticLinuxKubeletConfig {
		staticWindowsKubeletConfig[key] = val
	}

	// Default Kubelet config
	defaultKubeletConfig := map[string]string{
		"--cluster-domain":                  "cluster.local",
		"--network-plugin":                  "cni",
		"--pod-infra-container-image":       cloudSpecConfig.KubernetesSpecConfig.KubernetesImageBase + KubeConfigs[o.OrchestratorVersion]["pause"],
		"--max-pods":                        strconv.Itoa(DefaultKubernetesMaxPodsVNETIntegrated),
		"--eviction-hard":                   DefaultKubernetesHardEvictionThreshold,
		"--node-status-update-frequency":    KubeConfigs[o.OrchestratorVersion]["nodestatusfreq"],
		"--image-gc-high-threshold":         strconv.Itoa(DefaultKubernetesGCHighThreshold),
		"--image-gc-low-threshold":          strconv.Itoa(DefaultKubernetesGCLowThreshold),
		"--non-masquerade-cidr":             o.KubernetesConfig.ClusterSubnet,
		"--cloud-provider":                  "azure",
		"--cloud-config":                    "/etc/kubernetes/azure.json",
		"--azure-container-registry-config": "/etc/kubernetes/azure.json",
		"--event-qps":                       DefaultKubeletEventQPS,
		"--cadvisor-port":                   DefaultKubeletCadvisorPort,
	}

	// If no user-configurable kubelet config values exists, use the defaults
	setMissingKubeletValues(o.KubernetesConfig, defaultKubeletConfig)
	addDefaultFeatureGates(o.KubernetesConfig.KubeletConfig, o.OrchestratorVersion, "", "")

	// Override default cloud-provider?
	if helpers.IsTrueBoolPointer(o.KubernetesConfig.UseCloudControllerManager) {
		staticLinuxKubeletConfig["--cloud-provider"] = "external"
	}

	// Override default --network-plugin?
	if o.KubernetesConfig.NetworkPlugin == NetworkPluginKubenet {
		o.KubernetesConfig.KubeletConfig["--network-plugin"] = NetworkPluginKubenet
		o.KubernetesConfig.KubeletConfig["--max-pods"] = strconv.Itoa(DefaultKubernetesMaxPods)
	}

	// We don't support user-configurable values for the following,
	// so any of the value assignments below will override user-provided values
	var overrideKubeletConfig map[string]string
	if cs.Properties.HasWindows() {
		overrideKubeletConfig = staticWindowsKubeletConfig
	} else {
		overrideKubeletConfig = staticLinuxKubeletConfig
	}
	for key, val := range overrideKubeletConfig {
		o.KubernetesConfig.KubeletConfig[key] = val
	}

	// Get rid of values not supported in v1.5 clusters
	if !common.IsKubernetesVersionGe(o.OrchestratorVersion, "1.6.0") {
		for _, key := range []string{"--non-masquerade-cidr", "--cgroups-per-qos", "--enforce-node-allocatable"} {
			delete(o.KubernetesConfig.KubeletConfig, key)
		}
	}

	// Remove secure kubelet flags, if configured
	if !helpers.IsTrueBoolPointer(o.KubernetesConfig.EnableSecureKubelet) {
		for _, key := range []string{"--anonymous-auth", "--client-ca-file"} {
			delete(o.KubernetesConfig.KubeletConfig, key)
		}
	}

	// Master-specific kubelet config changes go here
	if cs.Properties.MasterProfile != nil {
		if cs.Properties.MasterProfile.KubernetesConfig == nil {
			cs.Properties.MasterProfile.KubernetesConfig = &api.KubernetesConfig{}
			cs.Properties.MasterProfile.KubernetesConfig.KubeletConfig = copyMap(cs.Properties.MasterProfile.KubernetesConfig.KubeletConfig)
		}
		setMissingKubeletValues(cs.Properties.MasterProfile.KubernetesConfig, o.KubernetesConfig.KubeletConfig)
		addDefaultFeatureGates(cs.Properties.MasterProfile.KubernetesConfig.KubeletConfig, o.OrchestratorVersion, "", "")

	}
	// Agent-specific kubelet config changes go here
	for _, profile := range cs.Properties.AgentPoolProfiles {
		if profile.KubernetesConfig == nil {
			profile.KubernetesConfig = &api.KubernetesConfig{}
			profile.KubernetesConfig.KubeletConfig = copyMap(profile.KubernetesConfig.KubeletConfig)
		}
		setMissingKubeletValues(profile.KubernetesConfig, o.KubernetesConfig.KubeletConfig)
		if !common.IsKubernetesVersionGe(o.OrchestratorVersion, "1.11.0") {
			addDefaultFeatureGates(profile.KubernetesConfig.KubeletConfig, o.OrchestratorVersion, "1.6.0", "Accelerators=true")
		}
	}
}

func setMissingKubeletValues(p *api.KubernetesConfig, d map[string]string) {
	if p.KubeletConfig == nil {
		p.KubeletConfig = d
	} else {
		for key, val := range d {
			// If we don't have a user-configurable value for each option
			if _, ok := p.KubeletConfig[key]; !ok {
				// then assign the default value
				p.KubeletConfig[key] = val
			}
		}
	}
}
func copyMap(input map[string]string) map[string]string {
	copy := map[string]string{}
	for key, value := range input {
		copy[key] = value
	}
	return copy
}
