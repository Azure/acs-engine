package acsengine

import (
	"strconv"

	"github.com/Azure/acs-engine/pkg/api"
)

func setControllerManagerConfig(cs *api.ContainerService) {
	o := cs.Properties.OrchestratorProfile
	staticLinuxControllerManagerConfig := map[string]string{
		"--kubeconfig":                       "/var/lib/kubelet/kubeconfig",
		"--allocate-node-cidrs":              strconv.FormatBool(!o.IsAzureCNI()),
		"--cluster-cidr":                     o.KubernetesConfig.ClusterSubnet,
		"--cloud-provider":                   "azure",
		"--cloud-config":                     "/etc/kubernetes/azure.json",
		"--root-ca-file":                     "/etc/kubernetes/certs/ca.crt",
		"--cluster-signing-cert-file":        "/etc/kubernetes/certs/ca.crt",
		"--cluster-signing-key-file":         "/etc/kubernetes/certs/ca.key",
		"--service-account-private-key-file": "/etc/kubernetes/certs/apiserver.key",
		"--leader-elect":                     "true",
		"--v":                                "2",
		"--profiling":                        "False",
	}

	// Set --cluster-name based on appropriate DNS prefix
	if cs.Properties.MasterProfile != nil {
		staticLinuxControllerManagerConfig["--cluster-name"] = cs.Properties.MasterProfile.DNSPrefix
	} else if cs.Properties.HostedMasterProfile != nil {
		staticLinuxControllerManagerConfig["--cluster-name"] = cs.Properties.HostedMasterProfile.DNSPrefix
	}

	staticWindowsControllerManagerConfig := make(map[string]string)
	for key, val := range staticLinuxControllerManagerConfig {
		staticWindowsControllerManagerConfig[key] = val
	}
	// Windows controller-manager config overrides
	// TODO placeholder for specific config overrides for Windows clusters

	// Default controller-manager config
	defaultControllerManagerConfig := map[string]string{
		"--node-monitor-grace-period":   DefaultKubernetesCtrlMgrNodeMonitorGracePeriod,
		"--pod-eviction-timeout":        DefaultKubernetesCtrlMgrPodEvictionTimeout,
		"--route-reconciliation-period": DefaultKubernetesCtrlMgrRouteReconciliationPeriod,
	}

	// If no user-configurable controller-manager config values exists, use the defaults
	if o.KubernetesConfig.ControllerManagerConfig == nil {
		o.KubernetesConfig.ControllerManagerConfig = defaultControllerManagerConfig
	} else {
		for key, val := range defaultControllerManagerConfig {
			// If we don't have a user-configurable controller-manager config for each option
			if _, ok := o.KubernetesConfig.ControllerManagerConfig[key]; !ok {
				// then assign the default value
				o.KubernetesConfig.ControllerManagerConfig[key] = val
			}
		}
	}

	// We don't support user-configurable values for the following,
	// so any of the value assignments below will override user-provided values
	var overrideControllerManagerConfig map[string]string
	if cs.Properties.HasWindows() {
		overrideControllerManagerConfig = staticWindowsControllerManagerConfig
	} else {
		overrideControllerManagerConfig = staticLinuxControllerManagerConfig
	}
	for key, val := range overrideControllerManagerConfig {
		o.KubernetesConfig.ControllerManagerConfig[key] = val
	}

	if *o.KubernetesConfig.EnableRbac {
		o.KubernetesConfig.ControllerManagerConfig["--use-service-account-credentials"] = "true"
	}
}
