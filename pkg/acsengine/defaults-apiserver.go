package acsengine

import (
	"strconv"

	"github.com/Azure/acs-engine/pkg/api"
	"github.com/Azure/acs-engine/pkg/api/common"
	"github.com/Azure/acs-engine/pkg/helpers"
)

func setAPIServerConfig(cs *api.ContainerService) {
	o := cs.Properties.OrchestratorProfile
	staticAPIServerConfig := map[string]string{
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

	// Default apiserver config
	defaultAPIServerConfig := map[string]string{
		"--audit-log-maxage":    "30",
		"--audit-log-maxbackup": "10",
		"--audit-log-maxsize":   "100",
	}

	// Data Encryption at REST configuration conditions
	if helpers.IsTrueBoolPointer(o.KubernetesConfig.EnableDataEncryptionAtRest) || helpers.IsTrueBoolPointer(o.KubernetesConfig.EnableEncryptionWithExternalKms) {
		staticAPIServerConfig["--experimental-encryption-provider-config"] = "/etc/kubernetes/encryption-config.yaml"
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
		staticAPIServerConfig["--cloud-provider"] = "azure"
		staticAPIServerConfig["--cloud-config"] = "/etc/kubernetes/azure.json"
	}

	// AAD configuration
	if cs.Properties.HasAadProfile() {
		defaultAPIServerConfig["--oidc-username-claim"] = "oid"
		defaultAPIServerConfig["--oidc-groups-claim"] = "groups"
		defaultAPIServerConfig["--oidc-client-id"] = "spn:" + cs.Properties.AADProfile.ServerAppID
		issuerHost := "sts.windows.net"
		if getCloudTargetEnv(cs.Location) == "AzureChinaCloud" {
			issuerHost = "sts.chinacloudapi.cn"
		}
		defaultAPIServerConfig["--oidc-issuer-url"] = "https://" + issuerHost + "/" + cs.Properties.AADProfile.TenantID + "/"
	}

	// Audit Policy configuration
	if common.IsKubernetesVersionGe(o.OrchestratorVersion, "1.8.0") {
		defaultAPIServerConfig["--audit-policy-file"] = "/etc/kubernetes/addons/audit-policy.yaml"
	}

	// RBAC configuration
	if helpers.IsTrueBoolPointer(o.KubernetesConfig.EnableRbac) {
		if common.IsKubernetesVersionGe(o.OrchestratorVersion, "1.7.0") {
			defaultAPIServerConfig["--authorization-mode"] = "Node,RBAC"
		} else {
			defaultAPIServerConfig["--authorization-mode"] = "RBAC"
		}
	}

	// Set default admission controllers
	admissionControlKey, admissionControlValues := getDefaultAdmissionControls(cs)
	defaultAPIServerConfig[admissionControlKey] = admissionControlValues

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
	for key, val := range staticAPIServerConfig {
		o.KubernetesConfig.APIServerConfig[key] = val
	}

	// Remove flags for secure communication to kubelet, if configured
	if !helpers.IsTrueBoolPointer(o.KubernetesConfig.EnableSecureKubelet) {
		for _, key := range []string{"--kubelet-client-certificate", "--kubelet-client-key"} {
			delete(o.KubernetesConfig.APIServerConfig, key)
		}
	}

	// Enforce flags removal that don't work with specific versions, to accommodate upgrade
	// Remove flags that are not compatible with 1.10
	if common.IsKubernetesVersionGe(o.OrchestratorVersion, "1.10.0") {
		for _, key := range []string{"--admission-control"} {
			delete(o.KubernetesConfig.APIServerConfig, key)
		}
	}
}

func getDefaultAdmissionControls(cs *api.ContainerService) (string, string) {
	o := cs.Properties.OrchestratorProfile
	admissionControlKey := "--enable-admission-plugins"
	var admissionControlValues string

	// --admission-control was used in v1.9 and earlier and was deprecated in 1.10
	if !common.IsKubernetesVersionGe(o.OrchestratorVersion, "1.10.0") {
		admissionControlKey = "--admission-control"
	}

	// Add new version case when applying admission controllers only available in that version or later
	switch {
	case common.IsKubernetesVersionGe(o.OrchestratorVersion, "1.9.0"):
		admissionControlValues = "NamespaceLifecycle,LimitRanger,ServiceAccount,DefaultStorageClass,DefaultTolerationSeconds,MutatingAdmissionWebhook,ValidatingAdmissionWebhook,ResourceQuota,AlwaysPullImages,ExtendedResourceToleration"
	default:
		admissionControlValues = "NamespaceLifecycle,LimitRanger,ServiceAccount,DefaultStorageClass,ResourceQuota,AlwaysPullImages"
	}

	// Pod Security Policy configuration
	if helpers.IsTrueBoolPointer(o.KubernetesConfig.EnablePodSecurityPolicy) {
		admissionControlValues += ",PodSecurityPolicy"
	}

	return admissionControlKey, admissionControlValues
}
