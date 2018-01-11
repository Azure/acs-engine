package acsengine

import (
	"bytes"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/Azure/acs-engine/pkg/api"
	"github.com/Azure/acs-engine/pkg/helpers"
)

func setKubeletConfig(cs *api.ContainerService) {
	o := cs.Properties.OrchestratorProfile
	cloudSpecConfig := GetCloudSpecConfig(cs.Location)
	staticLinuxKubeletConfig := map[string]string{
		"--address":                         "0.0.0.0",
		"--allow-privileged":                "true",
		"--anonymous-auth":                  "false",
		"--authorization-mode":              "Webhook",
		"--client-ca-file":                  "/etc/kubernetes/certs/ca.crt",
		"--pod-manifest-path":               "/etc/kubernetes/manifests",
		"--cluster-domain":                  "cluster.local",
		"--cluster-dns":                     DefaultKubernetesDNSServiceIP,
		"--cgroups-per-qos":                 "false",
		"--enforce-node-allocatable":        "",
		"--kubeconfig":                      "/var/lib/kubelet/kubeconfig",
		"--azure-container-registry-config": "/etc/kubernetes/azure.json",
		"--read-only-port":                  "0",
		"--keep-terminated-pod-volumes":     "false",
	}

	staticWindowsKubeletConfig := make(map[string]string)
	for key, val := range staticLinuxKubeletConfig {
		staticWindowsKubeletConfig[key] = val
	}
	// Windows kubelet config overrides
	staticWindowsKubeletConfig["--network-plugin"] = NetworkPluginKubenet

	// Default Kubelet config
	defaultKubeletConfig := map[string]string{
		"--network-plugin":               "cni",
		"--pod-infra-container-image":    cloudSpecConfig.KubernetesSpecConfig.KubernetesImageBase + KubeConfigs[o.OrchestratorVersion]["pause"],
		"--max-pods":                     strconv.Itoa(DefaultKubernetesKubeletMaxPods),
		"--eviction-hard":                DefaultKubernetesHardEvictionThreshold,
		"--node-status-update-frequency": KubeConfigs[o.OrchestratorVersion]["nodestatusfreq"],
		"--image-gc-high-threshold":      strconv.Itoa(DefaultKubernetesGCHighThreshold),
		"--image-gc-low-threshold":       strconv.Itoa(DefaultKubernetesGCLowThreshold),
		"--non-masquerade-cidr":          DefaultNonMasqueradeCidr,
		"--cloud-provider":               "azure",
		"--cloud-config":                 "/etc/kubernetes/azure.json",
	}

	// If no user-configurable kubelet config values exists, use the defaults
	setMissingKubeletValues(o.KubernetesConfig, defaultKubeletConfig)
	addDefaultFeatureGates(o.KubernetesConfig.KubeletConfig, "", "")

	// Override default cloud-provider?
	if helpers.IsTrueBoolPointer(o.KubernetesConfig.UseCloudControllerManager) {
		staticLinuxKubeletConfig["--cloud-provider"] = "external"
	}

	// Override default --network-plugin?
	if o.KubernetesConfig.NetworkPolicy == NetworkPolicyNone {
		o.KubernetesConfig.KubeletConfig["--network-plugin"] = NetworkPluginKubenet
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
	if !isKubernetesVersionGe(o.OrchestratorVersion, "1.6.0") {
		for _, key := range []string{"--non-masquerade-cidr", "--cgroups-per-qos", "--enforce-node-allocatable"} {
			delete(o.KubernetesConfig.KubeletConfig, key)
		}
	}

	// Master-specific kubelet config changes go here
	if cs.Properties.MasterProfile != nil {
		if cs.Properties.MasterProfile.KubernetesConfig == nil {
			cs.Properties.MasterProfile.KubernetesConfig = &api.KubernetesConfig{}
		}
		setMissingKubeletValues(cs.Properties.MasterProfile.KubernetesConfig, o.KubernetesConfig.KubeletConfig)
		addDefaultFeatureGates(cs.Properties.MasterProfile.KubernetesConfig.KubeletConfig, "", "")

	}
	// Agent-specific kubelet config changes go here
	for _, profile := range cs.Properties.AgentPoolProfiles {
		if profile.KubernetesConfig == nil {
			profile.KubernetesConfig = &api.KubernetesConfig{}
		}
		setMissingKubeletValues(profile.KubernetesConfig, o.KubernetesConfig.KubeletConfig)
		addDefaultFeatureGates(profile.KubernetesConfig.KubeletConfig, o.OrchestratorVersion, "Accelerators=true")
	}
}

// combine user-provided --feature-gates vals with defaults
// a minimum k8s version may be declared as required for defaults assignment
func addDefaultFeatureGates(m map[string]string, minVersion string, defaults string) {
	if minVersion != "" {
		if isKubernetesVersionGe(minVersion, "1.6.0") {
			m["--feature-gates"] = combineValues(m["--feature-gates"], defaults)
		} else {
			m["--feature-gates"] = combineValues(m["--feature-gates"], "")
		}
	} else {
		m["--feature-gates"] = combineValues(m["--feature-gates"], defaults)
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

func combineValues(inputs ...string) string {
	var valueMap map[string]string
	valueMap = make(map[string]string)
	for _, input := range inputs {
		applyValueStringToMap(valueMap, input)
	}
	return mapToString(valueMap)
}

func applyValueStringToMap(valueMap map[string]string, input string) {
	values := strings.Split(input, ",")
	for index := 0; index < len(values); index++ {
		// trim spaces (e.g. if the input was "foo=true, bar=true" - we want to drop the space after the comma)
		value := strings.Trim(values[index], " ")
		valueParts := strings.Split(value, "=")
		if len(valueParts) == 2 {
			valueMap[valueParts[0]] = valueParts[1]
		}
	}
}

func mapToString(valueMap map[string]string) string {
	// Order by key for consistency
	keys := []string{}
	for key := range valueMap {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	var buf bytes.Buffer
	for _, key := range keys {
		buf.WriteString(fmt.Sprintf("%s=%s,", key, valueMap[key]))
	}
	return strings.TrimSuffix(buf.String(), ",")
}
