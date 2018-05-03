package acsengine

import (
	"strconv"

	"github.com/Azure/acs-engine/pkg/api"
)

func setCloudControllerManagerConfig(cs *api.ContainerService) {
	o := cs.Properties.OrchestratorProfile
	staticLinuxCloudControllerManagerConfig := map[string]string{
		"--allocate-node-cidrs":    strconv.FormatBool(!o.IsAzureCNI()),
		"--configure-cloud-routes": strconv.FormatBool(o.RequireRouteTable()),
		"--cloud-provider":         "azure",
		"--cloud-config":           "/etc/kubernetes/azure.json",
		"--cluster-cidr":           o.KubernetesConfig.ClusterSubnet,
		"--kubeconfig":             "/var/lib/kubelet/kubeconfig",
		"--leader-elect":           "true",
		"--v":                      "2",
	}

	// Set --cluster-name based on appropriate DNS prefix
	if cs.Properties.MasterProfile != nil {
		staticLinuxCloudControllerManagerConfig["--cluster-name"] = cs.Properties.MasterProfile.DNSPrefix
	} else if cs.Properties.HostedMasterProfile != nil {
		staticLinuxCloudControllerManagerConfig["--cluster-name"] = cs.Properties.HostedMasterProfile.DNSPrefix
	}

	staticWindowsCloudControllerManagerConfig := make(map[string]string)
	for key, val := range staticLinuxCloudControllerManagerConfig {
		staticWindowsCloudControllerManagerConfig[key] = val
	}
	// Windows cloud-controller-manager config overrides
	// TODO placeholder for specific config overrides for Windows clusters

	// Default cloud-controller-manager config
	defaultCloudControllerManagerConfig := map[string]string{
		"--route-reconciliation-period": DefaultKubernetesCtrlMgrRouteReconciliationPeriod,
	}

	// If no user-configurable cloud-controller-manager config values exists, use the defaults
	if o.KubernetesConfig.CloudControllerManagerConfig == nil {
		o.KubernetesConfig.CloudControllerManagerConfig = defaultCloudControllerManagerConfig
	} else {
		for key, val := range defaultCloudControllerManagerConfig {
			// If we don't have a user-configurable cloud-controller-manager config for each option
			if _, ok := o.KubernetesConfig.CloudControllerManagerConfig[key]; !ok {
				// then assign the default value
				o.KubernetesConfig.CloudControllerManagerConfig[key] = val
			}
		}
	}

	// We don't support user-configurable values for the following,
	// so any of the value assignments below will override user-provided values
	var overrideCloudControllerManagerConfig map[string]string
	if cs.Properties.HasWindows() {
		overrideCloudControllerManagerConfig = staticWindowsCloudControllerManagerConfig
	} else {
		overrideCloudControllerManagerConfig = staticLinuxCloudControllerManagerConfig
	}
	for key, val := range overrideCloudControllerManagerConfig {
		o.KubernetesConfig.CloudControllerManagerConfig[key] = val
	}

	// TODO add RBAC support
	/*if *o.KubernetesConfig.EnableRbac {
		o.KubernetesConfig.CloudControllerManagerConfig["--use-service-account-credentials"] = "true"
	}*/
}
