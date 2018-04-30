package acsengine

import (
	"strconv"
	"strings"

	"github.com/Azure/acs-engine/pkg/api/common"
)

var k8sComponentVersions = map[string]map[string]string{
	"1.11": {
		"dockerEngine":    "1.13.*",
		"dashboard":       "kubernetes-dashboard-amd64:v1.8.3",
		"exechealthz":     "exechealthz-amd64:1.2",
		"addon-resizer":   "addon-resizer:1.8.1",
		"heapster":        "heapster-amd64:v1.5.1",
		"metrics-server":  "metrics-server-amd64:v0.2.1",
		"kube-dns":        "k8s-dns-kube-dns-amd64:1.14.8",
		"addon-manager":   "kube-addon-manager-amd64:v8.6",
		"dnsmasq":         "k8s-dns-dnsmasq-nanny-amd64:1.14.8",
		"pause":           "pause-amd64:3.1",
		"tiller":          "tiller:v2.8.1",
		"rescheduler":     "rescheduler:v0.3.1",
		"aci-connector":   "virtual-kubelet:latest",
		"nodestatusfreq":  DefaultKubernetesNodeStatusUpdateFrequency,
		"nodegraceperiod": DefaultKubernetesCtrlMgrNodeMonitorGracePeriod,
		"podeviction":     DefaultKubernetesCtrlMgrPodEvictionTimeout,
		"routeperiod":     DefaultKubernetesCtrlMgrRouteReconciliationPeriod,
		"backoffretries":  strconv.Itoa(DefaultKubernetesCloudProviderBackoffRetries),
		"backoffjitter":   strconv.FormatFloat(DefaultKubernetesCloudProviderBackoffJitter, 'f', -1, 64),
		"backoffduration": strconv.Itoa(DefaultKubernetesCloudProviderBackoffDuration),
		"backoffexponent": strconv.FormatFloat(DefaultKubernetesCloudProviderBackoffExponent, 'f', -1, 64),
		"ratelimitqps":    strconv.FormatFloat(DefaultKubernetesCloudProviderRateLimitQPS, 'f', -1, 64),
		"ratelimitbucket": strconv.Itoa(DefaultKubernetesCloudProviderRateLimitBucket),
		"gchighthreshold": strconv.Itoa(DefaultKubernetesGCHighThreshold),
		"gclowthreshold":  strconv.Itoa(DefaultKubernetesGCLowThreshold),
	},
	"1.10": {
		"dockerEngine":    "1.13.*",
		"dashboard":       "kubernetes-dashboard-amd64:v1.8.3",
		"exechealthz":     "exechealthz-amd64:1.2",
		"addon-resizer":   "addon-resizer:1.8.1",
		"heapster":        "heapster-amd64:v1.5.1",
		"metrics-server":  "metrics-server-amd64:v0.2.1",
		"kube-dns":        "k8s-dns-kube-dns-amd64:1.14.8",
		"addon-manager":   "kube-addon-manager-amd64:v8.6",
		"dnsmasq":         "k8s-dns-dnsmasq-nanny-amd64:1.14.8",
		"pause":           "pause-amd64:3.1",
		"tiller":          "tiller:v2.8.1",
		"rescheduler":     "rescheduler:v0.3.1",
		"aci-connector":   "virtual-kubelet:latest",
		"nodestatusfreq":  DefaultKubernetesNodeStatusUpdateFrequency,
		"nodegraceperiod": DefaultKubernetesCtrlMgrNodeMonitorGracePeriod,
		"podeviction":     DefaultKubernetesCtrlMgrPodEvictionTimeout,
		"routeperiod":     DefaultKubernetesCtrlMgrRouteReconciliationPeriod,
		"backoffretries":  strconv.Itoa(DefaultKubernetesCloudProviderBackoffRetries),
		"backoffjitter":   strconv.FormatFloat(DefaultKubernetesCloudProviderBackoffJitter, 'f', -1, 64),
		"backoffduration": strconv.Itoa(DefaultKubernetesCloudProviderBackoffDuration),
		"backoffexponent": strconv.FormatFloat(DefaultKubernetesCloudProviderBackoffExponent, 'f', -1, 64),
		"ratelimitqps":    strconv.FormatFloat(DefaultKubernetesCloudProviderRateLimitQPS, 'f', -1, 64),
		"ratelimitbucket": strconv.Itoa(DefaultKubernetesCloudProviderRateLimitBucket),
		"gchighthreshold": strconv.Itoa(DefaultKubernetesGCHighThreshold),
		"gclowthreshold":  strconv.Itoa(DefaultKubernetesGCLowThreshold),
	},
	"1.9": {
		"dockerEngine":    "1.13.*",
		"dashboard":       "kubernetes-dashboard-amd64:v1.8.3",
		"exechealthz":     "exechealthz-amd64:1.2",
		"addon-resizer":   "addon-resizer:1.8.1",
		"heapster":        "heapster-amd64:v1.5.1",
		"metrics-server":  "metrics-server-amd64:v0.2.1",
		"kube-dns":        "k8s-dns-kube-dns-amd64:1.14.8",
		"addon-manager":   "kube-addon-manager-amd64:v8.6",
		"dnsmasq":         "k8s-dns-dnsmasq-nanny-amd64:1.14.8",
		"pause":           "pause-amd64:3.1",
		"tiller":          "tiller:v2.8.1",
		"rescheduler":     "rescheduler:v0.3.1",
		"aci-connector":   "virtual-kubelet:latest",
		"nodestatusfreq":  DefaultKubernetesNodeStatusUpdateFrequency,
		"nodegraceperiod": DefaultKubernetesCtrlMgrNodeMonitorGracePeriod,
		"podeviction":     DefaultKubernetesCtrlMgrPodEvictionTimeout,
		"routeperiod":     DefaultKubernetesCtrlMgrRouteReconciliationPeriod,
		"backoffretries":  strconv.Itoa(DefaultKubernetesCloudProviderBackoffRetries),
		"backoffjitter":   strconv.FormatFloat(DefaultKubernetesCloudProviderBackoffJitter, 'f', -1, 64),
		"backoffduration": strconv.Itoa(DefaultKubernetesCloudProviderBackoffDuration),
		"backoffexponent": strconv.FormatFloat(DefaultKubernetesCloudProviderBackoffExponent, 'f', -1, 64),
		"ratelimitqps":    strconv.FormatFloat(DefaultKubernetesCloudProviderRateLimitQPS, 'f', -1, 64),
		"ratelimitbucket": strconv.Itoa(DefaultKubernetesCloudProviderRateLimitBucket),
		"gchighthreshold": strconv.Itoa(DefaultKubernetesGCHighThreshold),
		"gclowthreshold":  strconv.Itoa(DefaultKubernetesGCLowThreshold),
	},
	"1.8": {
		"dockerEngine":    "1.13.*",
		"dashboard":       "kubernetes-dashboard-amd64:v1.8.3",
		"exechealthz":     "exechealthz-amd64:1.2",
		"addon-resizer":   "addon-resizer:1.7",
		"heapster":        "heapster-amd64:v1.5.1",
		"metrics-server":  "metrics-server-amd64:v0.2.1",
		"kube-dns":        "k8s-dns-kube-dns-amd64:1.14.8",
		"addon-manager":   "kube-addon-manager-amd64:v8.6",
		"dnsmasq":         "k8s-dns-dnsmasq-nanny-amd64:1.14.8",
		"pause":           "pause-amd64:3.1",
		"tiller":          "tiller:v2.8.1",
		"rescheduler":     "rescheduler:v0.3.1",
		"aci-connector":   "virtual-kubelet:latest",
		"nodestatusfreq":  DefaultKubernetesNodeStatusUpdateFrequency,
		"nodegraceperiod": DefaultKubernetesCtrlMgrNodeMonitorGracePeriod,
		"podeviction":     DefaultKubernetesCtrlMgrPodEvictionTimeout,
		"routeperiod":     DefaultKubernetesCtrlMgrRouteReconciliationPeriod,
		"backoffretries":  strconv.Itoa(DefaultKubernetesCloudProviderBackoffRetries),
		"backoffjitter":   strconv.FormatFloat(DefaultKubernetesCloudProviderBackoffJitter, 'f', -1, 64),
		"backoffduration": strconv.Itoa(DefaultKubernetesCloudProviderBackoffDuration),
		"backoffexponent": strconv.FormatFloat(DefaultKubernetesCloudProviderBackoffExponent, 'f', -1, 64),
		"ratelimitqps":    strconv.FormatFloat(DefaultKubernetesCloudProviderRateLimitQPS, 'f', -1, 64),
		"ratelimitbucket": strconv.Itoa(DefaultKubernetesCloudProviderRateLimitBucket),
		"gchighthreshold": strconv.Itoa(DefaultKubernetesGCHighThreshold),
		"gclowthreshold":  strconv.Itoa(DefaultKubernetesGCLowThreshold),
	},
	"1.7": {
		"dockerEngine":    "1.13.*",
		"dashboard":       "kubernetes-dashboard-amd64:v1.6.3",
		"exechealthz":     "exechealthz-amd64:1.2",
		"addon-resizer":   "addon-resizer:1.7",
		"heapster":        "heapster-amd64:v1.5.1",
		"metrics-server":  "metrics-server-amd64:v0.2.1",
		"kube-dns":        "k8s-dns-kube-dns-amd64:1.14.5",
		"addon-manager":   "kube-addon-manager-amd64:v8.6",
		"dnsmasq":         "k8s-dns-dnsmasq-nanny-amd64:1.14.5",
		"pause":           "pause-amd64:3.1",
		"tiller":          "tiller:v2.8.1",
		"rescheduler":     "rescheduler:v0.3.1",
		"aci-connector":   "virtual-kubelet:latest",
		"nodestatusfreq":  DefaultKubernetesNodeStatusUpdateFrequency,
		"nodegraceperiod": DefaultKubernetesCtrlMgrNodeMonitorGracePeriod,
		"podeviction":     DefaultKubernetesCtrlMgrPodEvictionTimeout,
		"routeperiod":     DefaultKubernetesCtrlMgrRouteReconciliationPeriod,
		"backoffretries":  strconv.Itoa(DefaultKubernetesCloudProviderBackoffRetries),
		"backoffjitter":   strconv.FormatFloat(DefaultKubernetesCloudProviderBackoffJitter, 'f', -1, 64),
		"backoffduration": strconv.Itoa(DefaultKubernetesCloudProviderBackoffDuration),
		"backoffexponent": strconv.FormatFloat(DefaultKubernetesCloudProviderBackoffExponent, 'f', -1, 64),
		"ratelimitqps":    strconv.FormatFloat(DefaultKubernetesCloudProviderRateLimitQPS, 'f', -1, 64),
		"ratelimitbucket": strconv.Itoa(DefaultKubernetesCloudProviderRateLimitBucket),
		"gchighthreshold": strconv.Itoa(DefaultKubernetesGCHighThreshold),
		"gclowthreshold":  strconv.Itoa(DefaultKubernetesGCLowThreshold),
	},
	"1.6": {
		"dockerEngine":    "1.12.*",
		"dashboard":       "kubernetes-dashboard-amd64:v1.6.3",
		"exechealthz":     "exechealthz-amd64:1.2",
		"addon-resizer":   "addon-resizer:1.7",
		"heapster":        "heapster-amd64:v1.3.0",
		"metrics-server":  "metrics-server-amd64:v0.2.1",
		"kube-dns":        "k8s-dns-kube-dns-amd64:1.14.5",
		"addon-manager":   "kube-addon-manager-amd64:v6.5",
		"dnsmasq":         "k8s-dns-dnsmasq-nanny-amd64:1.14.5",
		"pause":           "pause-amd64:3.0",
		"tiller":          "tiller:v2.8.1",
		"rescheduler":     "rescheduler:v0.3.1",
		"aci-connector":   "virtual-kubelet:latest",
		"nodestatusfreq":  DefaultKubernetesNodeStatusUpdateFrequency,
		"nodegraceperiod": DefaultKubernetesCtrlMgrNodeMonitorGracePeriod,
		"podeviction":     DefaultKubernetesCtrlMgrPodEvictionTimeout,
		"routeperiod":     DefaultKubernetesCtrlMgrRouteReconciliationPeriod,
		"backoffretries":  strconv.Itoa(DefaultKubernetesCloudProviderBackoffRetries),
		"backoffjitter":   strconv.FormatFloat(DefaultKubernetesCloudProviderBackoffJitter, 'f', -1, 64),
		"backoffduration": strconv.Itoa(DefaultKubernetesCloudProviderBackoffDuration),
		"backoffexponent": strconv.FormatFloat(DefaultKubernetesCloudProviderBackoffExponent, 'f', -1, 64),
		"ratelimitqps":    strconv.FormatFloat(DefaultKubernetesCloudProviderRateLimitQPS, 'f', -1, 64),
		"ratelimitbucket": strconv.Itoa(DefaultKubernetesCloudProviderRateLimitBucket),
		"gchighthreshold": strconv.Itoa(DefaultKubernetesGCHighThreshold),
		"gclowthreshold":  strconv.Itoa(DefaultKubernetesGCLowThreshold),
	},
}

// KubeConfigs represents Docker images used for Kubernetes components based on Kubernetes versions (major.minor.patch)
var KubeConfigs = getKubeConfigs()

func getKubeConfigs() map[string]map[string]string {
	ret := make(map[string]map[string]string)
	for _, version := range common.GetAllSupportedKubernetesVersions() {
		ret[version] = getK8sVersionComponents(version, getVersionOverrides(version))
	}
	return ret
}

func getVersionOverrides(v string) map[string]string {
	switch v {
	case "1.8.11":
		return map[string]string{"kube-dns": "k8s-dns-kube-dns-amd64:1.14.9"}
	case "1.8.9":
		return map[string]string{"windowszip": "v1.8.9-2int.zip"}
	case "1.8.6":
		return map[string]string{"windowszip": "v1.8.6-2int.zip"}
	case "1.8.2":
		return map[string]string{"windowszip": "v1.8.2-2int.zip"}
	case "1.8.1":
		return map[string]string{"windowszip": "v1.8.1-2int.zip"}
	case "1.8.0":
		return map[string]string{"windowszip": "v1.8.0-2int.zip"}
	case "1.7.16":
		return map[string]string{"windowszip": "v1.7.16-1int.zip"}
	case "1.7.15":
		return map[string]string{"windowszip": "v1.7.15-1int.zip"}
	case "1.7.14":
		return map[string]string{"windowszip": "v1.7.14-1int.zip"}
	case "1.7.13":
		return map[string]string{"windowszip": "v1.7.13-1int.zip"}
	case "1.7.12":
		return map[string]string{"windowszip": "v1.7.12-2int.zip"}
	case "1.7.10":
		return map[string]string{"windowszip": "v1.7.10-1int.zip"}
	case "1.7.9":
		return map[string]string{"windowszip": "v1.7.9-2int.zip"}
	case "1.7.7":
		return map[string]string{"windowszip": "v1.7.7-2int.zip"}
	case "1.7.5":
		return map[string]string{"windowszip": "v1.7.5-4int.zip"}
	case "1.7.4":
		return map[string]string{"windowszip": "v1.7.4-2int.zip"}
	case "1.7.2":
		return map[string]string{"windowszip": "v1.7.2-1int.zip"}
	default:
		return nil
	}
}

func getK8sVersionComponents(version string, overrides map[string]string) map[string]string {
	s := strings.Split(version, ".")
	majorMinor := strings.Join(s[:2], ".")
	var ret map[string]string
	switch majorMinor {
	case "1.11":
		ret = map[string]string{
			"hyperkube":                   "hyperkube-amd64:v" + version,
			"ccm":                         "cloud-controller-manager-amd64:v" + version,
			"windowszip":                  "v" + version + "-1int.zip",
			"dockerEngineVersion":         k8sComponentVersions["1.11"]["dockerEngine"],
			DefaultDashboardAddonName:     k8sComponentVersions["1.11"]["dashboard"],
			"exechealthz":                 k8sComponentVersions["1.11"]["exechealthz"],
			"addonresizer":                k8sComponentVersions["1.11"]["addon-resizer"],
			"heapster":                    k8sComponentVersions["1.11"]["heapster"],
			DefaultMetricsServerAddonName: k8sComponentVersions["1.11"]["metrics-server"],
			"dns":                        k8sComponentVersions["1.11"]["kube-dns"],
			"addonmanager":               k8sComponentVersions["1.11"]["addon-manager"],
			"dnsmasq":                    k8sComponentVersions["1.11"]["dnsmasq"],
			"pause":                      k8sComponentVersions["1.11"]["pause"],
			DefaultTillerAddonName:       k8sComponentVersions["1.11"]["tiller"],
			DefaultReschedulerAddonName:  k8sComponentVersions["1.11"]["rescheduler"],
			DefaultACIConnectorAddonName: k8sComponentVersions["1.11"]["aci-connector"],
			"nodestatusfreq":             k8sComponentVersions["1.11"]["nodestatusfreq"],
			"nodegraceperiod":            k8sComponentVersions["1.11"]["nodegraceperiod"],
			"podeviction":                k8sComponentVersions["1.11"]["podeviction"],
			"routeperiod":                k8sComponentVersions["1.11"]["routeperiod"],
			"backoffretries":             k8sComponentVersions["1.11"]["backoffretries"],
			"backoffjitter":              k8sComponentVersions["1.11"]["backoffjitter"],
			"backoffduration":            k8sComponentVersions["1.11"]["backoffduration"],
			"backoffexponent":            k8sComponentVersions["1.11"]["backoffexponent"],
			"ratelimitqps":               k8sComponentVersions["1.11"]["ratelimitqps"],
			"ratelimitbucket":            k8sComponentVersions["1.11"]["ratelimitbucket"],
			"gchighthreshold":            k8sComponentVersions["1.11"]["gchighthreshold"],
			"gclowthreshold":             k8sComponentVersions["1.11"]["gclowthreshold"],
		}
	case "1.10":
		ret = map[string]string{
			"hyperkube":                   "hyperkube-amd64:v" + version,
			"ccm":                         "cloud-controller-manager-amd64:v" + version,
			"windowszip":                  "v" + version + "-1int.zip",
			"dockerEngineVersion":         k8sComponentVersions["1.10"]["dockerEngine"],
			DefaultDashboardAddonName:     k8sComponentVersions["1.10"]["dashboard"],
			"exechealthz":                 k8sComponentVersions["1.10"]["exechealthz"],
			"addonresizer":                k8sComponentVersions["1.10"]["addon-resizer"],
			"heapster":                    k8sComponentVersions["1.10"]["heapster"],
			DefaultMetricsServerAddonName: k8sComponentVersions["1.10"]["metrics-server"],
			"dns":                        k8sComponentVersions["1.10"]["kube-dns"],
			"addonmanager":               k8sComponentVersions["1.10"]["addon-manager"],
			"dnsmasq":                    k8sComponentVersions["1.10"]["dnsmasq"],
			"pause":                      k8sComponentVersions["1.10"]["pause"],
			DefaultTillerAddonName:       k8sComponentVersions["1.10"]["tiller"],
			DefaultReschedulerAddonName:  k8sComponentVersions["1.10"]["rescheduler"],
			DefaultACIConnectorAddonName: k8sComponentVersions["1.10"]["aci-connector"],
			"nodestatusfreq":             k8sComponentVersions["1.10"]["nodestatusfreq"],
			"nodegraceperiod":            k8sComponentVersions["1.10"]["nodegraceperiod"],
			"podeviction":                k8sComponentVersions["1.10"]["podeviction"],
			"routeperiod":                k8sComponentVersions["1.10"]["routeperiod"],
			"backoffretries":             k8sComponentVersions["1.10"]["backoffretries"],
			"backoffjitter":              k8sComponentVersions["1.10"]["backoffjitter"],
			"backoffduration":            k8sComponentVersions["1.10"]["backoffduration"],
			"backoffexponent":            k8sComponentVersions["1.10"]["backoffexponent"],
			"ratelimitqps":               k8sComponentVersions["1.10"]["ratelimitqps"],
			"ratelimitbucket":            k8sComponentVersions["1.10"]["ratelimitbucket"],
			"gchighthreshold":            k8sComponentVersions["1.10"]["gchighthreshold"],
			"gclowthreshold":             k8sComponentVersions["1.10"]["gclowthreshold"],
		}
	case "1.9":
		ret = map[string]string{
			"hyperkube":                   "hyperkube-amd64:v" + version,
			"ccm":                         "cloud-controller-manager-amd64:v" + version,
			"windowszip":                  "v" + version + "-1int.zip",
			"dockerEngineVersion":         k8sComponentVersions["1.9"]["dockerEngine"],
			DefaultDashboardAddonName:     k8sComponentVersions["1.9"]["dashboard"],
			"exechealthz":                 k8sComponentVersions["1.9"]["exechealthz"],
			"addonresizer":                k8sComponentVersions["1.9"]["addon-resizer"],
			"heapster":                    k8sComponentVersions["1.9"]["heapster"],
			DefaultMetricsServerAddonName: k8sComponentVersions["1.9"]["metrics-server"],
			"dns":                        k8sComponentVersions["1.9"]["kube-dns"],
			"addonmanager":               k8sComponentVersions["1.9"]["addon-manager"],
			"dnsmasq":                    k8sComponentVersions["1.9"]["dnsmasq"],
			"pause":                      k8sComponentVersions["1.9"]["pause"],
			DefaultTillerAddonName:       k8sComponentVersions["1.9"]["tiller"],
			DefaultReschedulerAddonName:  k8sComponentVersions["1.9"]["rescheduler"],
			DefaultACIConnectorAddonName: k8sComponentVersions["1.9"]["aci-connector"],
			"nodestatusfreq":             k8sComponentVersions["1.9"]["nodestatusfreq"],
			"nodegraceperiod":            k8sComponentVersions["1.9"]["nodegraceperiod"],
			"podeviction":                k8sComponentVersions["1.9"]["podeviction"],
			"routeperiod":                k8sComponentVersions["1.9"]["routeperiod"],
			"backoffretries":             k8sComponentVersions["1.9"]["backoffretries"],
			"backoffjitter":              k8sComponentVersions["1.9"]["backoffjitter"],
			"backoffduration":            k8sComponentVersions["1.9"]["backoffduration"],
			"backoffexponent":            k8sComponentVersions["1.9"]["backoffexponent"],
			"ratelimitqps":               k8sComponentVersions["1.9"]["ratelimitqps"],
			"ratelimitbucket":            k8sComponentVersions["1.9"]["ratelimitbucket"],
			"gchighthreshold":            k8sComponentVersions["1.9"]["gchighthreshold"],
			"gclowthreshold":             k8sComponentVersions["1.9"]["gclowthreshold"],
		}
	case "1.8":
		ret = map[string]string{
			"hyperkube":                   "hyperkube-amd64:v" + version,
			"ccm":                         "cloud-controller-manager-amd64:v" + version,
			"windowszip":                  "v" + version + "-1int.zip",
			"dockerEngineVersion":         k8sComponentVersions["1.8"]["dockerEngine"],
			DefaultDashboardAddonName:     k8sComponentVersions["1.8"]["dashboard"],
			"exechealthz":                 k8sComponentVersions["1.8"]["exechealthz"],
			"addonresizer":                k8sComponentVersions["1.8"]["addon-resizer"],
			"heapster":                    k8sComponentVersions["1.8"]["heapster"],
			DefaultMetricsServerAddonName: k8sComponentVersions["1.8"]["metrics-server"],
			"dns":                        k8sComponentVersions["1.8"]["kube-dns"],
			"addonmanager":               k8sComponentVersions["1.8"]["addon-manager"],
			"dnsmasq":                    k8sComponentVersions["1.8"]["dnsmasq"],
			"pause":                      k8sComponentVersions["1.8"]["pause"],
			DefaultTillerAddonName:       k8sComponentVersions["1.8"]["tiller"],
			DefaultReschedulerAddonName:  k8sComponentVersions["1.8"]["rescheduler"],
			DefaultACIConnectorAddonName: k8sComponentVersions["1.8"]["aci-connector"],
			"nodestatusfreq":             k8sComponentVersions["1.8"]["nodestatusfreq"],
			"nodegraceperiod":            k8sComponentVersions["1.8"]["nodegraceperiod"],
			"podeviction":                k8sComponentVersions["1.8"]["podeviction"],
			"routeperiod":                k8sComponentVersions["1.8"]["routeperiod"],
			"backoffretries":             k8sComponentVersions["1.8"]["backoffretries"],
			"backoffjitter":              k8sComponentVersions["1.8"]["backoffjitter"],
			"backoffduration":            k8sComponentVersions["1.8"]["backoffduration"],
			"backoffexponent":            k8sComponentVersions["1.8"]["backoffexponent"],
			"ratelimitqps":               k8sComponentVersions["1.8"]["ratelimitqps"],
			"ratelimitbucket":            k8sComponentVersions["1.8"]["ratelimitbucket"],
			"gchighthreshold":            k8sComponentVersions["1.8"]["gchighthreshold"],
			"gclowthreshold":             k8sComponentVersions["1.8"]["gclowthreshold"],
		}
	case "1.7":
		ret = map[string]string{
			"hyperkube":                   "hyperkube-amd64:v" + version,
			"dockerEngineVersion":         k8sComponentVersions["1.7"]["dockerEngine"],
			DefaultDashboardAddonName:     k8sComponentVersions["1.7"]["dashboard"],
			"exechealthz":                 k8sComponentVersions["1.7"]["exechealthz"],
			"addonresizer":                k8sComponentVersions["1.7"]["addon-resizer"],
			"heapster":                    k8sComponentVersions["1.7"]["heapster"],
			DefaultMetricsServerAddonName: k8sComponentVersions["1.7"]["metrics-server"],
			"dns":                        k8sComponentVersions["1.7"]["kube-dns"],
			"addonmanager":               k8sComponentVersions["1.7"]["addon-manager"],
			"dnsmasq":                    k8sComponentVersions["1.7"]["dnsmasq"],
			"pause":                      k8sComponentVersions["1.7"]["pause"],
			DefaultTillerAddonName:       k8sComponentVersions["1.7"]["tiller"],
			DefaultReschedulerAddonName:  k8sComponentVersions["1.7"]["rescheduler"],
			DefaultACIConnectorAddonName: k8sComponentVersions["1.7"]["aci-connector"],
			"nodestatusfreq":             k8sComponentVersions["1.7"]["nodestatusfreq"],
			"nodegraceperiod":            k8sComponentVersions["1.7"]["nodegraceperiod"],
			"podeviction":                k8sComponentVersions["1.7"]["podeviction"],
			"routeperiod":                k8sComponentVersions["1.7"]["routeperiod"],
			"backoffretries":             k8sComponentVersions["1.7"]["backoffretries"],
			"backoffjitter":              k8sComponentVersions["1.7"]["backoffjitter"],
			"backoffduration":            k8sComponentVersions["1.7"]["backoffduration"],
			"backoffexponent":            k8sComponentVersions["1.7"]["backoffexponent"],
			"ratelimitqps":               k8sComponentVersions["1.7"]["ratelimitqps"],
			"ratelimitbucket":            k8sComponentVersions["1.7"]["ratelimitbucket"],
			"gchighthreshold":            k8sComponentVersions["1.7"]["gchighthreshold"],
			"gclowthreshold":             k8sComponentVersions["1.7"]["gclowthreshold"],
		}
	case "1.6":
		ret = map[string]string{
			"hyperkube":                   "hyperkube-amd64:v" + version,
			"dockerEngineVersion":         k8sComponentVersions["1.6"]["dockerEngine"],
			DefaultDashboardAddonName:     k8sComponentVersions["1.6"]["dashboard"],
			"exechealthz":                 k8sComponentVersions["1.6"]["exechealthz"],
			"addonresizer":                k8sComponentVersions["1.6"]["addon-resizer"],
			"heapster":                    k8sComponentVersions["1.6"]["heapster"],
			DefaultMetricsServerAddonName: k8sComponentVersions["1.6"]["metrics-server"],
			"dns":                        k8sComponentVersions["1.6"]["kube-dns"],
			"addonmanager":               k8sComponentVersions["1.6"]["addon-manager"],
			"dnsmasq":                    k8sComponentVersions["1.6"]["dnsmasq"],
			"pause":                      k8sComponentVersions["1.6"]["pause"],
			DefaultTillerAddonName:       k8sComponentVersions["1.6"]["tiller"],
			DefaultReschedulerAddonName:  k8sComponentVersions["1.6"]["rescheduler"],
			DefaultACIConnectorAddonName: k8sComponentVersions["1.6"]["aci-connector"],
			"nodestatusfreq":             k8sComponentVersions["1.6"]["nodestatusfreq"],
			"nodegraceperiod":            k8sComponentVersions["1.6"]["nodegraceperiod"],
			"podeviction":                k8sComponentVersions["1.6"]["podeviction"],
			"routeperiod":                k8sComponentVersions["1.6"]["routeperiod"],
			"backoffretries":             k8sComponentVersions["1.6"]["backoffretries"],
			"backoffjitter":              k8sComponentVersions["1.6"]["backoffjitter"],
			"backoffduration":            k8sComponentVersions["1.6"]["backoffduration"],
			"backoffexponent":            k8sComponentVersions["1.6"]["backoffexponent"],
			"ratelimitqps":               k8sComponentVersions["1.6"]["ratelimitqps"],
			"ratelimitbucket":            k8sComponentVersions["1.6"]["ratelimitbucket"],
			"gchighthreshold":            k8sComponentVersions["1.6"]["gchighthreshold"],
			"gclowthreshold":             k8sComponentVersions["1.6"]["gclowthreshold"],
		}

	default:
		ret = nil
	}
	for k, v := range overrides {
		ret[k] = v
	}
	return ret
}
