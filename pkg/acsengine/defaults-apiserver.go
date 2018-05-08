package acsengine

import (
	"strconv"

	"github.com/Azure/acs-engine/pkg/api"
	"github.com/Azure/acs-engine/pkg/api/common"
	"github.com/Azure/acs-engine/pkg/helpers"
)

func setAPIServerConfig(cs *api.ContainerService) {
	o := cs.Properties.OrchestratorProfile
	staticLinuxAPIServerConfig := map[string]string{
		"--bind-address":               "0.0.0.0",
		"--advertise-address":          "<kubernetesAPIServerIP>",
		"--allow-privileged":           "true",
		"--anonymous-auth":             "false",
		"--audit-log-path":             "/var/log/audit.log",
		"--insecure-port":              "8080",
		"--secure-port":                "443",
		"--service-account-lookup":     "true",
		"--etcd-cafile":                "/etc/kubernetes/certs/ca.crt",
		"--etcd-certfile":              "/etc/kubernetes/certs/etcdclient.crt",
		"--etcd-keyfile":               "/etc/kubernetes/certs/etcdclient.key",
		"--etcd-servers":               "https://127.0.0.1:" + strconv.Itoa(DefaultMasterEtcdClientPort),
		"--tls-cert-file":              "/etc/kubernetes/certs/apiserver.crt",
		"--tls-private-key-file":       "/etc/kubernetes/certs/apiserver.key",
		"--client-ca-file":             "/etc/kubernetes/certs/ca.crt",
		"--profiling":                  "false",
		"--repair-malformed-updates":   "false",
		"--service-account-key-file":   "/etc/kubernetes/certs/apiserver.key",
		"--kubelet-client-certificate": "/etc/kubernetes/certs/client.crt",
		"--kubelet-client-key":         "/etc/kubernetes/certs/client.key",
		"--service-cluster-ip-range":   o.KubernetesConfig.ServiceCIDR,
		"--storage-backend":            o.GetAPIServerEtcdAPIVersion(),
		"--v":                          "4",
	}

	// Windows apiserver config overrides
	// TODO placeholder for specific config overrides for Windows clusters
	staticWindowsAPIServerConfig := make(map[string]string)
	for key, val := range staticLinuxAPIServerConfig {
		staticWindowsAPIServerConfig[key] = val
	}

	// Default apiserver config
	defaultAPIServerConfig := map[string]string{
		"--admission-control":   "NamespaceLifecycle,LimitRanger,ServiceAccount,DefaultStorageClass,ResourceQuota,DenyEscalatingExec,AlwaysPullImages",
		"--audit-log-maxage":    "30",
		"--audit-log-maxbackup": "10",
		"--audit-log-maxsize":   "100",
	}

	// Data Encryption at REST configuration
	if helpers.IsTrueBoolPointer(o.KubernetesConfig.EnableDataEncryptionAtRest) {
		staticLinuxAPIServerConfig["--experimental-encryption-provider-config"] = "/etc/kubernetes/encryption-config.yaml"
	}

	// Data Encryption at REST with external KMS configuration
	if helpers.IsTrueBoolPointer(o.KubernetesConfig.EnableEncryptionWithExternalKms) {
		staticLinuxAPIServerConfig["--experimental-encryption-provider-config"] = "/etc/kubernetes/encryption-config.yaml"
	}

	// Aggregated API configuration
	if o.KubernetesConfig.EnableAggregatedAPIs || common.IsKubernetesVersionGe(o.OrchestratorVersion, "1.9.0") {
		defaultAPIServerConfig["--requestheader-client-ca-file"] = "/etc/kubernetes/certs/proxy-ca.crt"
		defaultAPIServerConfig["--proxy-client-cert-file"] = "/etc/kubernetes/certs/proxy.crt"
		defaultAPIServerConfig["--proxy-client-key-file"] = "/etc/kubernetes/certs/proxy.key"
		defaultAPIServerConfig["--requestheader-allowed-names"] = ""
		defaultAPIServerConfig["--requestheader-extra-headers-prefix"] = "X-Remote-Extra-"
		defaultAPIServerConfig["--requestheader-group-headers"] = "X-Remote-Group"
		defaultAPIServerConfig["--requestheader-username-headers"] = "X-Remote-User"
	}

	// Enable cloudprovider if we're not using cloud controller manager
	if !helpers.IsTrueBoolPointer(o.KubernetesConfig.UseCloudControllerManager) {
		staticLinuxAPIServerConfig["--cloud-provider"] = "azure"
		staticLinuxAPIServerConfig["--cloud-config"] = "/etc/kubernetes/azure.json"
	}

	// AAD configuration
	if cs.Properties.HasAadProfile() {
		defaultAPIServerConfig["--oidc-username-claim"] = "oid"
		defaultAPIServerConfig["--oidc-groups-claim"] = "groups"
		defaultAPIServerConfig["--oidc-client-id"] = "spn:" + cs.Properties.AADProfile.ServerAppID
		issuerHost := "sts.windows.net"
		if GetCloudTargetEnv(cs.Location) == "AzureChinaCloud" {
			issuerHost = "sts.chinacloudapi.cn"
		}
		defaultAPIServerConfig["--oidc-issuer-url"] = "https://" + issuerHost + "/" + cs.Properties.AADProfile.TenantID + "/"
	}

	// Audit Policy configuration
	if common.IsKubernetesVersionGe(o.OrchestratorVersion, "1.8.0") {
		staticLinuxAPIServerConfig["--audit-policy-file"] = "/etc/kubernetes/manifests/audit-policy.yaml"
	}

	// RBAC configuration
	if helpers.IsTrueBoolPointer(o.KubernetesConfig.EnableRbac) {
		if common.IsKubernetesVersionGe(o.OrchestratorVersion, "1.7.0") {
			defaultAPIServerConfig["--authorization-mode"] = "Node,RBAC"
		} else {
			defaultAPIServerConfig["--authorization-mode"] = "RBAC"
		}
	}

	// Pod Security Policy configuration
	if helpers.IsTrueBoolPointer(o.KubernetesConfig.EnablePodSecurityPolicy) {
		defaultAPIServerConfig["--admission-control"] = defaultAPIServerConfig["--admission-control"] + ",PodSecurityPolicy"
	}

	// If no user-configurable apiserver config values exists, use the defaults
	if o.KubernetesConfig.APIServerConfig == nil {
		o.KubernetesConfig.APIServerConfig = defaultAPIServerConfig
	} else {
		for key, val := range defaultAPIServerConfig {
			// If we don't have a user-configurable apiserver config for each option
			if _, ok := o.KubernetesConfig.APIServerConfig[key]; !ok {
				// then assign the default value
				o.KubernetesConfig.APIServerConfig[key] = val
			}
		}
	}

	// We don't support user-configurable values for the following,
	// so any of the value assignments below will override user-provided values
	var overrideAPIServerConfig map[string]string
	if cs.Properties.HasWindows() {
		overrideAPIServerConfig = staticWindowsAPIServerConfig
	} else {
		overrideAPIServerConfig = staticLinuxAPIServerConfig
	}
	for key, val := range overrideAPIServerConfig {
		o.KubernetesConfig.APIServerConfig[key] = val
	}

	// Remove flags for secure communication to kubelet, if configured
	if !helpers.IsTrueBoolPointer(o.KubernetesConfig.EnableSecureKubelet) {
		for _, key := range []string{"--kubelet-client-certificate", "--kubelet-client-key"} {
			delete(o.KubernetesConfig.APIServerConfig, key)
		}
	}
}
