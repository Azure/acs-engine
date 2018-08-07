package acsengine

import (
	"strconv"
	"strings"

	"github.com/Azure/acs-engine/pkg/api"
	"github.com/Azure/acs-engine/pkg/api/common"
	"github.com/Azure/acs-engine/pkg/helpers"
)

func setKubeletConfig(cs *api.ContainerService) {
	o := cs.Properties.OrchestratorProfile
	cloudSpecConfig := getCloudSpecConfig(cs.Location)
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
	staticWindowsKubeletConfig["--azure-container-registry-config"] = "c:\\k\\azure.json"
	staticWindowsKubeletConfig["--pod-infra-container-image"] = "kubletwin/pause"
	staticWindowsKubeletConfig["--kubeconfig"] = "c:\\k\\config"
	staticWindowsKubeletConfig["--cloud-config"] = "c:\\k\\azure.json"
	staticWindowsKubeletConfig["--cgroups-per-qos"] = "false"
	staticWindowsKubeletConfig["--enforce-node-allocatable"] = "\"\""

	// Default Kubelet config
	defaultKubeletConfig := map[string]string{
		"--cluster-domain":                  "cluster.local",
		"--network-plugin":                  "cni",
		"--pod-infra-container-image":       cloudSpecConfig.KubernetesSpecConfig.KubernetesImageBase + KubeConfigs[o.OrchestratorVersion]["pause"],
		"--max-pods":                        strconv.Itoa(DefaultKubernetesMaxPods),
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
		"--pod-max-pids":                    strconv.Itoa(DefaultKubeletPodMaxPIDs),
		"--image-pull-progress-deadline":    "30m",
	}

	// Apply Azure CNI-specific --max-pods value
	if o.KubernetesConfig.NetworkPlugin == NetworkPluginAzure {
		defaultKubeletConfig["--max-pods"] = strconv.Itoa(DefaultKubernetesMaxPodsVNETIntegrated)
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
		if o.KubernetesConfig.NetworkPolicy != NetworkPolicyCalico {
			o.KubernetesConfig.KubeletConfig["--network-plugin"] = NetworkPluginKubenet
		}
	}

	// We don't support user-configurable values for the following,
	// so any of the value assignments below will override user-provided values
	for key, val := range staticLinuxKubeletConfig {
		o.KubernetesConfig.KubeletConfig[key] = val
	}

	// Get rid of values not supported in v1.5 clusters
	if !common.IsKubernetesVersionGe(o.OrchestratorVersion, "1.6.0") {
		for _, key := range []string{"--non-masquerade-cidr", "--cgroups-per-qos", "--enforce-node-allocatable"} {
			delete(o.KubernetesConfig.KubeletConfig, key)
		}
	}

	// Get rid of values not supported in v1.10 clusters
	if !common.IsKubernetesVersionGe(o.OrchestratorVersion, "1.10.0") {
		for _, key := range []string{"--pod-max-pids"} {
			delete(o.KubernetesConfig.KubeletConfig, key)
		}
	}

	// Get rid of values not supported in v1.12 and up
	if common.IsKubernetesVersionGe(o.OrchestratorVersion, "1.12.0-alpha.1") {
		for _, key := range []string{"--cadvisor-port"} {
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
			if profile.OSType == "Windows" {
				for key, val := range staticWindowsKubeletConfig {
					profile.KubernetesConfig.KubeletConfig[key] = val
				}
			}
		}
		setMissingKubeletValues(profile.KubernetesConfig, o.KubernetesConfig.KubeletConfig)

		// For N Series (GPU) VMs
		if strings.Contains(profile.VMSize, "Standard_N") {
			if !cs.Properties.IsNVIDIADevicePluginEnabled() && !common.IsKubernetesVersionGe(o.OrchestratorVersion, "1.11.0") {
				// enabling accelerators for Kubernetes >= 1.6 to <= 1.9
				addDefaultFeatureGates(profile.KubernetesConfig.KubeletConfig, o.OrchestratorVersion, "1.6.0", "Accelerators=true")
			}
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
