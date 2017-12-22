package acsengine

import (
	"strconv"

	"github.com/Azure/acs-engine/pkg/api/common"
)

// KubeConfigs represents Docker images used for Kubernetes components based on Kubernetes versions (major.minor.patch)
var KubeConfigs = map[string]map[string]string{
	common.KubernetesVersion1Dot9Dot0: {
		"hyperkube": "hyperkube-amd64:v1.9.0",
		"ccm":       "cloud-controller-manager-amd64:v1.9.0",
		"dockerEngineVersion":       "1.12.*",
		DefaultDashboardAddonName:   "kubernetes-dashboard-amd64:v1.9.0", // TODO fix missing gcrio.azureedge.net/google_containers/kubernetes-dashboard-amd64:v1.9.0
		"exechealthz":               "exechealthz-amd64:1.2",
		"addonresizer":              "addon-resizer:1.7",
		"heapster":                  "heapster-amd64:v1.4.2",
		"dns":                       "k8s-dns-kube-dns-amd64:1.14.5",
		"addonmanager":              "kube-addon-manager-amd64:v6.4-beta.2",
		"dnsmasq":                   "k8s-dns-dnsmasq-nanny-amd64:1.14.5",
		"pause":                     "pause-amd64:3.0",
		DefaultTillerAddonName:      DefaultTillerImage,
		DefaultReschedulerAddonName: DefaultReschedulerImage,
		"windowszip":                "v1.9.0-1int.zip",
		"nodestatusfreq":            DefaultKubernetesNodeStatusUpdateFrequency,
		"nodegraceperiod":           DefaultKubernetesCtrlMgrNodeMonitorGracePeriod,
		"podeviction":               DefaultKubernetesCtrlMgrPodEvictionTimeout,
		"routeperiod":               DefaultKubernetesCtrlMgrRouteReconciliationPeriod,
		"backoffretries":            strconv.Itoa(DefaultKubernetesCloudProviderBackoffRetries),
		"backoffjitter":             strconv.FormatFloat(DefaultKubernetesCloudProviderBackoffJitter, 'f', -1, 64),
		"backoffduration":           strconv.Itoa(DefaultKubernetesCloudProviderBackoffDuration),
		"backoffexponent":           strconv.FormatFloat(DefaultKubernetesCloudProviderBackoffExponent, 'f', -1, 64),
		"ratelimitqps":              strconv.FormatFloat(DefaultKubernetesCloudProviderRateLimitQPS, 'f', -1, 64),
		"ratelimitbucket":           strconv.Itoa(DefaultKubernetesCloudProviderRateLimitBucket),
		"gchighthreshold":           strconv.Itoa(DefaultKubernetesGCHighThreshold),
		"gclowthreshold":            strconv.Itoa(DefaultKubernetesGCLowThreshold),
	},
	common.KubernetesVersion1Dot8Dot4: {
		"hyperkube": "hyperkube-amd64:v1.8.4",
		"ccm":       "cloud-controller-manager-amd64:v1.8.4",
		"dockerEngineVersion":       "1.12.*",
		DefaultDashboardAddonName:   "kubernetes-dashboard-amd64:v1.8.0",
		"exechealthz":               "exechealthz-amd64:1.2",
		"addonresizer":              "addon-resizer:1.7",
		"heapster":                  "heapster-amd64:v1.4.2",
		"dns":                       "k8s-dns-kube-dns-amd64:1.14.5",
		"addonmanager":              "kube-addon-manager-amd64:v6.4-beta.2",
		"dnsmasq":                   "k8s-dns-dnsmasq-nanny-amd64:1.14.5",
		"pause":                     "pause-amd64:3.0",
		DefaultTillerAddonName:      DefaultTillerImage,
		DefaultReschedulerAddonName: DefaultReschedulerImage,
		"windowszip":                "v1.8.4-1int.zip",
		"nodestatusfreq":            DefaultKubernetesNodeStatusUpdateFrequency,
		"nodegraceperiod":           DefaultKubernetesCtrlMgrNodeMonitorGracePeriod,
		"podeviction":               DefaultKubernetesCtrlMgrPodEvictionTimeout,
		"routeperiod":               DefaultKubernetesCtrlMgrRouteReconciliationPeriod,
		"backoffretries":            strconv.Itoa(DefaultKubernetesCloudProviderBackoffRetries),
		"backoffjitter":             strconv.FormatFloat(DefaultKubernetesCloudProviderBackoffJitter, 'f', -1, 64),
		"backoffduration":           strconv.Itoa(DefaultKubernetesCloudProviderBackoffDuration),
		"backoffexponent":           strconv.FormatFloat(DefaultKubernetesCloudProviderBackoffExponent, 'f', -1, 64),
		"ratelimitqps":              strconv.FormatFloat(DefaultKubernetesCloudProviderRateLimitQPS, 'f', -1, 64),
		"ratelimitbucket":           strconv.Itoa(DefaultKubernetesCloudProviderRateLimitBucket),
		"gchighthreshold":           strconv.Itoa(DefaultKubernetesGCHighThreshold),
		"gclowthreshold":            strconv.Itoa(DefaultKubernetesGCLowThreshold),
	},
	common.KubernetesVersion1Dot8Dot2: {
		"hyperkube": "hyperkube-amd64:v1.8.2",
		"ccm":       "cloud-controller-manager-amd64:v1.8.2",
		"dockerEngineVersion":        "1.12.*",
		DefaultDashboardAddonName:    "kubernetes-dashboard-amd64:v1.7.1",
		"exechealthz":                "exechealthz-amd64:1.2",
		"addonresizer":               "addon-resizer:1.7",
		"heapster":                   "heapster-amd64:v1.4.2",
		"dns":                        "k8s-dns-kube-dns-amd64:1.14.5",
		"addonmanager":               "kube-addon-manager-amd64:v6.4-beta.2",
		"dnsmasq":                    "k8s-dns-dnsmasq-nanny-amd64:1.14.5",
		"pause":                      "pause-amd64:3.0",
		DefaultTillerAddonName:       DefaultTillerImage,
		DefaultACIConnectorAddonName: DefaultACIConnectorImage,
		DefaultReschedulerAddonName:  DefaultReschedulerImage,
		"windowszip":                 "v1.8.2-2int.zip",
		"nodestatusfreq":             DefaultKubernetesNodeStatusUpdateFrequency,
		"nodegraceperiod":            DefaultKubernetesCtrlMgrNodeMonitorGracePeriod,
		"podeviction":                DefaultKubernetesCtrlMgrPodEvictionTimeout,
		"routeperiod":                DefaultKubernetesCtrlMgrRouteReconciliationPeriod,
		"backoffretries":             strconv.Itoa(DefaultKubernetesCloudProviderBackoffRetries),
		"backoffjitter":              strconv.FormatFloat(DefaultKubernetesCloudProviderBackoffJitter, 'f', -1, 64),
		"backoffduration":            strconv.Itoa(DefaultKubernetesCloudProviderBackoffDuration),
		"backoffexponent":            strconv.FormatFloat(DefaultKubernetesCloudProviderBackoffExponent, 'f', -1, 64),
		"ratelimitqps":               strconv.FormatFloat(DefaultKubernetesCloudProviderRateLimitQPS, 'f', -1, 64),
		"ratelimitbucket":            strconv.Itoa(DefaultKubernetesCloudProviderRateLimitBucket),
		"gchighthreshold":            strconv.Itoa(DefaultKubernetesGCHighThreshold),
		"gclowthreshold":             strconv.Itoa(DefaultKubernetesGCLowThreshold),
	},
	common.KubernetesVersion1Dot8Dot1: {
		"hyperkube": "hyperkube-amd64:v1.8.1",
		"ccm":       "cloud-controller-manager-amd64:v1.8.1",
		"dockerEngineVersion":        "1.12.*",
		DefaultDashboardAddonName:    "kubernetes-dashboard-amd64:v1.7.1",
		"exechealthz":                "exechealthz-amd64:1.2",
		"addonresizer":               "addon-resizer:1.7",
		"heapster":                   "heapster-amd64:v1.4.2",
		"dns":                        "k8s-dns-kube-dns-amd64:1.14.5",
		"addonmanager":               "kube-addon-manager-amd64:v6.4-beta.2",
		"dnsmasq":                    "k8s-dns-dnsmasq-nanny-amd64:1.14.5",
		"pause":                      "pause-amd64:3.0",
		DefaultTillerAddonName:       DefaultTillerImage,
		DefaultACIConnectorAddonName: DefaultACIConnectorImage,
		DefaultReschedulerAddonName:  DefaultReschedulerImage,
		"windowszip":                 "v1.8.1-2int.zip",
		"nodestatusfreq":             DefaultKubernetesNodeStatusUpdateFrequency,
		"nodegraceperiod":            DefaultKubernetesCtrlMgrNodeMonitorGracePeriod,
		"podeviction":                DefaultKubernetesCtrlMgrPodEvictionTimeout,
		"routeperiod":                DefaultKubernetesCtrlMgrRouteReconciliationPeriod,
		"backoffretries":             strconv.Itoa(DefaultKubernetesCloudProviderBackoffRetries),
		"backoffjitter":              strconv.FormatFloat(DefaultKubernetesCloudProviderBackoffJitter, 'f', -1, 64),
		"backoffduration":            strconv.Itoa(DefaultKubernetesCloudProviderBackoffDuration),
		"backoffexponent":            strconv.FormatFloat(DefaultKubernetesCloudProviderBackoffExponent, 'f', -1, 64),
		"ratelimitqps":               strconv.FormatFloat(DefaultKubernetesCloudProviderRateLimitQPS, 'f', -1, 64),
		"ratelimitbucket":            strconv.Itoa(DefaultKubernetesCloudProviderRateLimitBucket),
		"gchighthreshold":            strconv.Itoa(DefaultKubernetesGCHighThreshold),
		"gclowthreshold":             strconv.Itoa(DefaultKubernetesGCLowThreshold),
	},
	common.KubernetesVersion1Dot8Dot0: {
		"hyperkube": "hyperkube-amd64:v1.8.0",
		"ccm":       "cloud-controller-manager-amd64:v1.8.0",
		"dockerEngineVersion":        "1.12.*",
		DefaultDashboardAddonName:    "kubernetes-dashboard-amd64:v1.7.1",
		"exechealthz":                "exechealthz-amd64:1.2",
		"addonresizer":               "addon-resizer:1.7",
		"heapster":                   "heapster-amd64:v1.4.2",
		"dns":                        "k8s-dns-kube-dns-amd64:1.14.5",
		"addonmanager":               "kube-addon-manager-amd64:v6.4-beta.2",
		"dnsmasq":                    "k8s-dns-dnsmasq-nanny-amd64:1.14.5",
		"pause":                      "pause-amd64:3.0",
		DefaultTillerAddonName:       DefaultTillerImage,
		DefaultACIConnectorAddonName: DefaultACIConnectorImage,
		DefaultReschedulerAddonName:  DefaultReschedulerImage,
		"windowszip":                 "v1.8.0-2int.zip",
		"nodestatusfreq":             DefaultKubernetesNodeStatusUpdateFrequency,
		"nodegraceperiod":            DefaultKubernetesCtrlMgrNodeMonitorGracePeriod,
		"podeviction":                DefaultKubernetesCtrlMgrPodEvictionTimeout,
		"routeperiod":                DefaultKubernetesCtrlMgrRouteReconciliationPeriod,
		"backoffretries":             strconv.Itoa(DefaultKubernetesCloudProviderBackoffRetries),
		"backoffjitter":              strconv.FormatFloat(DefaultKubernetesCloudProviderBackoffJitter, 'f', -1, 64),
		"backoffduration":            strconv.Itoa(DefaultKubernetesCloudProviderBackoffDuration),
		"backoffexponent":            strconv.FormatFloat(DefaultKubernetesCloudProviderBackoffExponent, 'f', -1, 64),
		"ratelimitqps":               strconv.FormatFloat(DefaultKubernetesCloudProviderRateLimitQPS, 'f', -1, 64),
		"ratelimitbucket":            strconv.Itoa(DefaultKubernetesCloudProviderRateLimitBucket),
		"gchighthreshold":            strconv.Itoa(DefaultKubernetesGCHighThreshold),
		"gclowthreshold":             strconv.Itoa(DefaultKubernetesGCLowThreshold),
	},
	common.KubernetesVersion1Dot7Dot10: {
		"hyperkube":                  "hyperkube-amd64:v1.7.10",
		"dockerEngineVersion":        "1.12.*",
		DefaultDashboardAddonName:    "kubernetes-dashboard-amd64:v1.6.3",
		"exechealthz":                "exechealthz-amd64:1.2",
		"addonresizer":               "addon-resizer:1.7",
		"heapster":                   "heapster-amd64:v1.4.2",
		"dns":                        "k8s-dns-kube-dns-amd64:1.14.5",
		"addonmanager":               "kube-addon-manager-amd64:v6.4-beta.2",
		"dnsmasq":                    "k8s-dns-dnsmasq-nanny-amd64:1.14.5",
		"pause":                      "pause-amd64:3.0",
		DefaultTillerAddonName:       DefaultTillerImage,
		DefaultACIConnectorAddonName: DefaultACIConnectorImage,
		DefaultReschedulerAddonName:  DefaultReschedulerImage,
		"windowszip":                 "v1.7.10-1int.zip",
		"nodestatusfreq":             DefaultKubernetesNodeStatusUpdateFrequency,
		"nodegraceperiod":            DefaultKubernetesCtrlMgrNodeMonitorGracePeriod,
		"podeviction":                DefaultKubernetesCtrlMgrPodEvictionTimeout,
		"routeperiod":                DefaultKubernetesCtrlMgrRouteReconciliationPeriod,
		"backoffretries":             strconv.Itoa(DefaultKubernetesCloudProviderBackoffRetries),
		"backoffjitter":              strconv.FormatFloat(DefaultKubernetesCloudProviderBackoffJitter, 'f', -1, 64),
		"backoffduration":            strconv.Itoa(DefaultKubernetesCloudProviderBackoffDuration),
		"backoffexponent":            strconv.FormatFloat(DefaultKubernetesCloudProviderBackoffExponent, 'f', -1, 64),
		"ratelimitqps":               strconv.FormatFloat(DefaultKubernetesCloudProviderRateLimitQPS, 'f', -1, 64),
		"ratelimitbucket":            strconv.Itoa(DefaultKubernetesCloudProviderRateLimitBucket),
		"gchighthreshold":            strconv.Itoa(DefaultKubernetesGCHighThreshold),
		"gclowthreshold":             strconv.Itoa(DefaultKubernetesGCLowThreshold),
	},
	common.KubernetesVersion1Dot7Dot9: {
		"hyperkube":                  "hyperkube-amd64:v1.7.9",
		"dockerEngineVersion":        "1.12.*",
		DefaultDashboardAddonName:    "kubernetes-dashboard-amd64:v1.6.3",
		"exechealthz":                "exechealthz-amd64:1.2",
		"addonresizer":               "addon-resizer:1.7",
		"heapster":                   "heapster-amd64:v1.4.2",
		"dns":                        "k8s-dns-kube-dns-amd64:1.14.5",
		"addonmanager":               "kube-addon-manager-amd64:v6.4-beta.2",
		"dnsmasq":                    "k8s-dns-dnsmasq-nanny-amd64:1.14.5",
		"pause":                      "pause-amd64:3.0",
		DefaultTillerAddonName:       DefaultTillerImage,
		DefaultACIConnectorAddonName: DefaultACIConnectorImage,
		DefaultReschedulerAddonName:  DefaultReschedulerImage,
		"windowszip":                 "v1.7.9-2int.zip",
		"nodestatusfreq":             DefaultKubernetesNodeStatusUpdateFrequency,
		"nodegraceperiod":            DefaultKubernetesCtrlMgrNodeMonitorGracePeriod,
		"podeviction":                DefaultKubernetesCtrlMgrPodEvictionTimeout,
		"routeperiod":                DefaultKubernetesCtrlMgrRouteReconciliationPeriod,
		"backoffretries":             strconv.Itoa(DefaultKubernetesCloudProviderBackoffRetries),
		"backoffjitter":              strconv.FormatFloat(DefaultKubernetesCloudProviderBackoffJitter, 'f', -1, 64),
		"backoffduration":            strconv.Itoa(DefaultKubernetesCloudProviderBackoffDuration),
		"backoffexponent":            strconv.FormatFloat(DefaultKubernetesCloudProviderBackoffExponent, 'f', -1, 64),
		"ratelimitqps":               strconv.FormatFloat(DefaultKubernetesCloudProviderRateLimitQPS, 'f', -1, 64),
		"ratelimitbucket":            strconv.Itoa(DefaultKubernetesCloudProviderRateLimitBucket),
		"gchighthreshold":            strconv.Itoa(DefaultKubernetesGCHighThreshold),
		"gclowthreshold":             strconv.Itoa(DefaultKubernetesGCLowThreshold),
	},
	common.KubernetesVersion1Dot7Dot7: {
		"hyperkube":                  "hyperkube-amd64:v1.7.7",
		"dockerEngineVersion":        "1.12.*",
		DefaultDashboardAddonName:    "kubernetes-dashboard-amd64:v1.6.3",
		"exechealthz":                "exechealthz-amd64:1.2",
		"addonresizer":               "addon-resizer:1.7",
		"heapster":                   "heapster-amd64:v1.4.2",
		"dns":                        "k8s-dns-kube-dns-amd64:1.14.5",
		"addonmanager":               "kube-addon-manager-amd64:v6.4-beta.2",
		"dnsmasq":                    "k8s-dns-dnsmasq-nanny-amd64:1.14.5",
		"pause":                      "pause-amd64:3.0",
		DefaultTillerAddonName:       DefaultTillerImage,
		DefaultACIConnectorAddonName: DefaultACIConnectorImage,
		DefaultReschedulerAddonName:  DefaultReschedulerImage,
		"windowszip":                 "v1.7.7-2int.zip",
		"nodestatusfreq":             DefaultKubernetesNodeStatusUpdateFrequency,
		"nodegraceperiod":            DefaultKubernetesCtrlMgrNodeMonitorGracePeriod,
		"podeviction":                DefaultKubernetesCtrlMgrPodEvictionTimeout,
		"routeperiod":                DefaultKubernetesCtrlMgrRouteReconciliationPeriod,
		"backoffretries":             strconv.Itoa(DefaultKubernetesCloudProviderBackoffRetries),
		"backoffjitter":              strconv.FormatFloat(DefaultKubernetesCloudProviderBackoffJitter, 'f', -1, 64),
		"backoffduration":            strconv.Itoa(DefaultKubernetesCloudProviderBackoffDuration),
		"backoffexponent":            strconv.FormatFloat(DefaultKubernetesCloudProviderBackoffExponent, 'f', -1, 64),
		"ratelimitqps":               strconv.FormatFloat(DefaultKubernetesCloudProviderRateLimitQPS, 'f', -1, 64),
		"ratelimitbucket":            strconv.Itoa(DefaultKubernetesCloudProviderRateLimitBucket),
		"gchighthreshold":            strconv.Itoa(DefaultKubernetesGCHighThreshold),
		"gclowthreshold":             strconv.Itoa(DefaultKubernetesGCLowThreshold),
	},
	common.KubernetesVersion1Dot7Dot5: {
		"hyperkube":                  "hyperkube-amd64:v1.7.5",
		"dockerEngineVersion":        "1.12.*",
		DefaultDashboardAddonName:    "kubernetes-dashboard-amd64:v1.6.3",
		"exechealthz":                "exechealthz-amd64:1.2",
		"addonresizer":               "addon-resizer:1.7",
		"heapster":                   "heapster-amd64:v1.4.2",
		"dns":                        "k8s-dns-kube-dns-amd64:1.14.4",
		"addonmanager":               "kube-addon-manager-amd64:v6.4-beta.2",
		"dnsmasq":                    "k8s-dns-dnsmasq-nanny-amd64:1.14.4",
		"pause":                      "pause-amd64:3.0",
		DefaultTillerAddonName:       DefaultTillerImage,
		DefaultACIConnectorAddonName: DefaultACIConnectorImage,
		DefaultReschedulerAddonName:  DefaultReschedulerImage,
		"windowszip":                 "v1.7.5-4int.zip",
		"nodestatusfreq":             DefaultKubernetesNodeStatusUpdateFrequency,
		"nodegraceperiod":            DefaultKubernetesCtrlMgrNodeMonitorGracePeriod,
		"podeviction":                DefaultKubernetesCtrlMgrPodEvictionTimeout,
		"routeperiod":                DefaultKubernetesCtrlMgrRouteReconciliationPeriod,
		"backoffretries":             strconv.Itoa(DefaultKubernetesCloudProviderBackoffRetries),
		"backoffjitter":              strconv.FormatFloat(DefaultKubernetesCloudProviderBackoffJitter, 'f', -1, 64),
		"backoffduration":            strconv.Itoa(DefaultKubernetesCloudProviderBackoffDuration),
		"backoffexponent":            strconv.FormatFloat(DefaultKubernetesCloudProviderBackoffExponent, 'f', -1, 64),
		"ratelimitqps":               strconv.FormatFloat(DefaultKubernetesCloudProviderRateLimitQPS, 'f', -1, 64),
		"ratelimitbucket":            strconv.Itoa(DefaultKubernetesCloudProviderRateLimitBucket),
		"gchighthreshold":            strconv.Itoa(DefaultKubernetesGCHighThreshold),
		"gclowthreshold":             strconv.Itoa(DefaultKubernetesGCLowThreshold),
	},
	common.KubernetesVersion1Dot7Dot4: {
		"hyperkube":                  "hyperkube-amd64:v1.7.4",
		"dockerEngineVersion":        "1.12.*",
		DefaultDashboardAddonName:    "kubernetes-dashboard-amd64:v1.6.3",
		"exechealthz":                "exechealthz-amd64:1.2",
		"addonresizer":               "addon-resizer:1.7",
		"heapster":                   "heapster-amd64:v1.4.1",
		"dns":                        "k8s-dns-kube-dns-amd64:1.14.4",
		"addonmanager":               "kube-addon-manager-amd64:v6.4-beta.2",
		"dnsmasq":                    "k8s-dns-dnsmasq-nanny-amd64:1.14.4",
		"pause":                      "pause-amd64:3.0",
		DefaultTillerAddonName:       DefaultTillerImage,
		DefaultACIConnectorAddonName: DefaultACIConnectorImage,
		DefaultReschedulerAddonName:  DefaultReschedulerImage,
		"windowszip":                 "v1.7.4-2int.zip",
		"nodestatusfreq":             DefaultKubernetesNodeStatusUpdateFrequency,
		"nodegraceperiod":            DefaultKubernetesCtrlMgrNodeMonitorGracePeriod,
		"podeviction":                DefaultKubernetesCtrlMgrPodEvictionTimeout,
		"routeperiod":                DefaultKubernetesCtrlMgrRouteReconciliationPeriod,
		"backoffretries":             strconv.Itoa(DefaultKubernetesCloudProviderBackoffRetries),
		"backoffjitter":              strconv.FormatFloat(DefaultKubernetesCloudProviderBackoffJitter, 'f', -1, 64),
		"backoffduration":            strconv.Itoa(DefaultKubernetesCloudProviderBackoffDuration),
		"backoffexponent":            strconv.FormatFloat(DefaultKubernetesCloudProviderBackoffExponent, 'f', -1, 64),
		"ratelimitqps":               strconv.FormatFloat(DefaultKubernetesCloudProviderRateLimitQPS, 'f', -1, 64),
		"ratelimitbucket":            strconv.Itoa(DefaultKubernetesCloudProviderRateLimitBucket),
		"gchighthreshold":            strconv.Itoa(DefaultKubernetesGCHighThreshold),
		"gclowthreshold":             strconv.Itoa(DefaultKubernetesGCLowThreshold),
	},
	common.KubernetesVersion1Dot7Dot2: {
		"hyperkube":                  "hyperkube-amd64:v1.7.2",
		"dockerEngineVersion":        "1.12.*",
		DefaultDashboardAddonName:    "kubernetes-dashboard-amd64:v1.6.3",
		"exechealthz":                "exechealthz-amd64:1.2",
		"addonresizer":               "addon-resizer:1.7",
		"heapster":                   "heapster-amd64:v1.4.1",
		"dns":                        "k8s-dns-kube-dns-amd64:1.14.4",
		"addonmanager":               "kube-addon-manager-amd64:v6.4-beta.2",
		"dnsmasq":                    "k8s-dns-dnsmasq-nanny-amd64:1.14.4",
		"pause":                      "pause-amd64:3.0",
		DefaultTillerAddonName:       DefaultTillerImage,
		DefaultACIConnectorAddonName: DefaultACIConnectorImage,
		DefaultReschedulerAddonName:  DefaultReschedulerImage,
		"windowszip":                 "v1.7.2-1int.zip",
		"nodestatusfreq":             DefaultKubernetesNodeStatusUpdateFrequency,
		"nodegraceperiod":            DefaultKubernetesCtrlMgrNodeMonitorGracePeriod,
		"podeviction":                DefaultKubernetesCtrlMgrPodEvictionTimeout,
		"routeperiod":                DefaultKubernetesCtrlMgrRouteReconciliationPeriod,
		"backoffretries":             strconv.Itoa(DefaultKubernetesCloudProviderBackoffRetries),
		"backoffjitter":              strconv.FormatFloat(DefaultKubernetesCloudProviderBackoffJitter, 'f', -1, 64),
		"backoffduration":            strconv.Itoa(DefaultKubernetesCloudProviderBackoffDuration),
		"backoffexponent":            strconv.FormatFloat(DefaultKubernetesCloudProviderBackoffExponent, 'f', -1, 64),
		"ratelimitqps":               strconv.FormatFloat(DefaultKubernetesCloudProviderRateLimitQPS, 'f', -1, 64),
		"ratelimitbucket":            strconv.Itoa(DefaultKubernetesCloudProviderRateLimitBucket),
		"gchighthreshold":            strconv.Itoa(DefaultKubernetesGCHighThreshold),
		"gclowthreshold":             strconv.Itoa(DefaultKubernetesGCLowThreshold),
	},
	common.KubernetesVersion1Dot7Dot1: {
		"hyperkube":                  "hyperkube-amd64:v1.7.1",
		"dockerEngineVersion":        "1.12.*",
		DefaultDashboardAddonName:    "kubernetes-dashboard-amd64:v1.6.1",
		"exechealthz":                "exechealthz-amd64:1.2",
		"addonresizer":               "addon-resizer:1.7",
		"heapster":                   "heapster:v1.4.0",
		"dns":                        "k8s-dns-kube-dns-amd64:1.14.4",
		"addonmanager":               "kube-addon-manager-amd64:v6.4-beta.2",
		"dnsmasq":                    "k8s-dns-dnsmasq-nanny-amd64:1.14.4",
		"pause":                      "pause-amd64:3.0",
		DefaultTillerAddonName:       DefaultTillerImage,
		DefaultACIConnectorAddonName: DefaultACIConnectorImage,
		DefaultReschedulerAddonName:  DefaultReschedulerImage,
		"nodestatusfreq":             DefaultKubernetesNodeStatusUpdateFrequency,
		"nodegraceperiod":            DefaultKubernetesCtrlMgrNodeMonitorGracePeriod,
		"podeviction":                DefaultKubernetesCtrlMgrPodEvictionTimeout,
		"routeperiod":                DefaultKubernetesCtrlMgrRouteReconciliationPeriod,
		"backoffretries":             strconv.Itoa(DefaultKubernetesCloudProviderBackoffRetries),
		"backoffjitter":              strconv.FormatFloat(DefaultKubernetesCloudProviderBackoffJitter, 'f', -1, 64),
		"backoffduration":            strconv.Itoa(DefaultKubernetesCloudProviderBackoffDuration),
		"backoffexponent":            strconv.FormatFloat(DefaultKubernetesCloudProviderBackoffExponent, 'f', -1, 64),
		"ratelimitqps":               strconv.FormatFloat(DefaultKubernetesCloudProviderRateLimitQPS, 'f', -1, 64),
		"ratelimitbucket":            strconv.Itoa(DefaultKubernetesCloudProviderRateLimitBucket),
		"gchighthreshold":            strconv.Itoa(DefaultKubernetesGCHighThreshold),
		"gclowthreshold":             strconv.Itoa(DefaultKubernetesGCLowThreshold),
	},
	common.KubernetesVersion1Dot7Dot0: {
		"hyperkube":                  "hyperkube-amd64:v1.7.0",
		"dockerEngineVersion":        "1.12.*",
		DefaultDashboardAddonName:    "kubernetes-dashboard-amd64:v1.6.1",
		"exechealthz":                "exechealthz-amd64:1.2",
		"addonresizer":               "addon-resizer:1.7",
		"heapster":                   "heapster:v1.4.0",
		"dns":                        "k8s-dns-kube-dns-amd64:1.14.4",
		"addonmanager":               "kube-addon-manager-amd64:v6.4-beta.2",
		"dnsmasq":                    "k8s-dns-dnsmasq-nanny-amd64:1.14.4",
		"pause":                      "pause-amd64:3.0",
		DefaultTillerAddonName:       DefaultTillerImage,
		DefaultACIConnectorAddonName: DefaultACIConnectorImage,
		DefaultReschedulerAddonName:  DefaultReschedulerImage,
		"nodestatusfreq":             DefaultKubernetesNodeStatusUpdateFrequency,
		"nodegraceperiod":            DefaultKubernetesCtrlMgrNodeMonitorGracePeriod,
		"podeviction":                DefaultKubernetesCtrlMgrPodEvictionTimeout,
		"routeperiod":                DefaultKubernetesCtrlMgrRouteReconciliationPeriod,
		"backoffretries":             strconv.Itoa(DefaultKubernetesCloudProviderBackoffRetries),
		"backoffjitter":              strconv.FormatFloat(DefaultKubernetesCloudProviderBackoffJitter, 'f', -1, 64),
		"backoffduration":            strconv.Itoa(DefaultKubernetesCloudProviderBackoffDuration),
		"backoffexponent":            strconv.FormatFloat(DefaultKubernetesCloudProviderBackoffExponent, 'f', -1, 64),
		"ratelimitqps":               strconv.FormatFloat(DefaultKubernetesCloudProviderRateLimitQPS, 'f', -1, 64),
		"ratelimitbucket":            strconv.Itoa(DefaultKubernetesCloudProviderRateLimitBucket),
		"gchighthreshold":            strconv.Itoa(DefaultKubernetesGCHighThreshold),
		"gclowthreshold":             strconv.Itoa(DefaultKubernetesGCLowThreshold),
	},
	common.KubernetesVersion1Dot6Dot13: {
		"hyperkube":                  "hyperkube-amd64:v1.6.13",
		"dockerEngineVersion":        "1.12.*",
		DefaultDashboardAddonName:    "kubernetes-dashboard-amd64:v1.6.3",
		"exechealthz":                "exechealthz-amd64:1.2",
		"addonresizer":               "addon-resizer:1.7",
		"heapster":                   "heapster-amd64:v1.3.0",
		"dns":                        "k8s-dns-kube-dns-amd64:1.14.5",
		"addonmanager":               "kube-addon-manager-amd64:v6.4-beta.2",
		"dnsmasq":                    "k8s-dns-dnsmasq-nanny-amd64:1.14.5",
		"pause":                      "pause-amd64:3.0",
		DefaultTillerAddonName:       DefaultTillerImage,
		DefaultACIConnectorAddonName: DefaultACIConnectorImage,
		DefaultReschedulerAddonName:  DefaultReschedulerImage,
		"nodestatusfreq":             DefaultKubernetesNodeStatusUpdateFrequency,
		"nodegraceperiod":            DefaultKubernetesCtrlMgrNodeMonitorGracePeriod,
		"podeviction":                DefaultKubernetesCtrlMgrPodEvictionTimeout,
		"routeperiod":                DefaultKubernetesCtrlMgrRouteReconciliationPeriod,
		"backoffretries":             strconv.Itoa(DefaultKubernetesCloudProviderBackoffRetries),
		"backoffjitter":              strconv.FormatFloat(DefaultKubernetesCloudProviderBackoffJitter, 'f', -1, 64),
		"backoffduration":            strconv.Itoa(DefaultKubernetesCloudProviderBackoffDuration),
		"backoffexponent":            strconv.FormatFloat(DefaultKubernetesCloudProviderBackoffExponent, 'f', -1, 64),
		"ratelimitqps":               strconv.FormatFloat(DefaultKubernetesCloudProviderRateLimitQPS, 'f', -1, 64),
		"ratelimitbucket":            strconv.Itoa(DefaultKubernetesCloudProviderRateLimitBucket),
		"gchighthreshold":            strconv.Itoa(DefaultKubernetesGCHighThreshold),
		"gclowthreshold":             strconv.Itoa(DefaultKubernetesGCLowThreshold),
	},
	common.KubernetesVersion1Dot6Dot12: {
		"hyperkube":                  "hyperkube-amd64:v1.6.12",
		"dockerEngineVersion":        "1.12.*",
		DefaultDashboardAddonName:    "kubernetes-dashboard-amd64:v1.6.3",
		"exechealthz":                "exechealthz-amd64:1.2",
		"addonresizer":               "addon-resizer:1.7",
		"heapster":                   "heapster-amd64:v1.3.0",
		"dns":                        "k8s-dns-kube-dns-amd64:1.14.5",
		"addonmanager":               "kube-addon-manager-amd64:v6.4-beta.2",
		"dnsmasq":                    "k8s-dns-dnsmasq-nanny-amd64:1.14.5",
		"pause":                      "pause-amd64:3.0",
		DefaultTillerAddonName:       DefaultTillerImage,
		DefaultACIConnectorAddonName: DefaultACIConnectorImage,
		DefaultReschedulerAddonName:  DefaultReschedulerImage,
		"nodestatusfreq":             DefaultKubernetesNodeStatusUpdateFrequency,
		"nodegraceperiod":            DefaultKubernetesCtrlMgrNodeMonitorGracePeriod,
		"podeviction":                DefaultKubernetesCtrlMgrPodEvictionTimeout,
		"routeperiod":                DefaultKubernetesCtrlMgrRouteReconciliationPeriod,
		"backoffretries":             strconv.Itoa(DefaultKubernetesCloudProviderBackoffRetries),
		"backoffjitter":              strconv.FormatFloat(DefaultKubernetesCloudProviderBackoffJitter, 'f', -1, 64),
		"backoffduration":            strconv.Itoa(DefaultKubernetesCloudProviderBackoffDuration),
		"backoffexponent":            strconv.FormatFloat(DefaultKubernetesCloudProviderBackoffExponent, 'f', -1, 64),
		"ratelimitqps":               strconv.FormatFloat(DefaultKubernetesCloudProviderRateLimitQPS, 'f', -1, 64),
		"ratelimitbucket":            strconv.Itoa(DefaultKubernetesCloudProviderRateLimitBucket),
		"gchighthreshold":            strconv.Itoa(DefaultKubernetesGCHighThreshold),
		"gclowthreshold":             strconv.Itoa(DefaultKubernetesGCLowThreshold),
	},
	common.KubernetesVersion1Dot6Dot11: {
		"hyperkube":                  "hyperkube-amd64:v1.6.11",
		"dockerEngineVersion":        "1.12.*",
		DefaultDashboardAddonName:    "kubernetes-dashboard-amd64:v1.6.3",
		"exechealthz":                "exechealthz-amd64:1.2",
		"addonresizer":               "addon-resizer:1.7",
		"heapster":                   "heapster-amd64:v1.3.0",
		"dns":                        "k8s-dns-kube-dns-amd64:1.14.5",
		"addonmanager":               "kube-addon-manager-amd64:v6.4-beta.2",
		"dnsmasq":                    "k8s-dns-dnsmasq-nanny-amd64:1.14.5",
		"pause":                      "pause-amd64:3.0",
		DefaultTillerAddonName:       DefaultTillerImage,
		DefaultACIConnectorAddonName: DefaultACIConnectorImage,
		DefaultReschedulerAddonName:  DefaultReschedulerImage,
		"nodestatusfreq":             DefaultKubernetesNodeStatusUpdateFrequency,
		"nodegraceperiod":            DefaultKubernetesCtrlMgrNodeMonitorGracePeriod,
		"podeviction":                DefaultKubernetesCtrlMgrPodEvictionTimeout,
		"routeperiod":                DefaultKubernetesCtrlMgrRouteReconciliationPeriod,
		"backoffretries":             strconv.Itoa(DefaultKubernetesCloudProviderBackoffRetries),
		"backoffjitter":              strconv.FormatFloat(DefaultKubernetesCloudProviderBackoffJitter, 'f', -1, 64),
		"backoffduration":            strconv.Itoa(DefaultKubernetesCloudProviderBackoffDuration),
		"backoffexponent":            strconv.FormatFloat(DefaultKubernetesCloudProviderBackoffExponent, 'f', -1, 64),
		"ratelimitqps":               strconv.FormatFloat(DefaultKubernetesCloudProviderRateLimitQPS, 'f', -1, 64),
		"ratelimitbucket":            strconv.Itoa(DefaultKubernetesCloudProviderRateLimitBucket),
		"gchighthreshold":            strconv.Itoa(DefaultKubernetesGCHighThreshold),
		"gclowthreshold":             strconv.Itoa(DefaultKubernetesGCLowThreshold),
	},
	common.KubernetesVersion1Dot6Dot9: {
		"hyperkube":                  "hyperkube-amd64:v1.6.9",
		"dockerEngineVersion":        "1.12.*",
		DefaultDashboardAddonName:    "kubernetes-dashboard-amd64:v1.6.3",
		"exechealthz":                "exechealthz-amd64:1.2",
		"addonresizer":               "addon-resizer:1.7",
		"heapster":                   "heapster-amd64:v1.3.0",
		"dns":                        "k8s-dns-kube-dns-amd64:1.14.4",
		"addonmanager":               "kube-addon-manager-amd64:v6.4-beta.2",
		"dnsmasq":                    "k8s-dns-dnsmasq-nanny-amd64:1.14.4",
		"pause":                      "pause-amd64:3.0",
		DefaultTillerAddonName:       DefaultTillerImage,
		DefaultACIConnectorAddonName: DefaultACIConnectorImage,
		DefaultReschedulerAddonName:  DefaultReschedulerImage,
		"nodestatusfreq":             DefaultKubernetesNodeStatusUpdateFrequency,
		"nodegraceperiod":            DefaultKubernetesCtrlMgrNodeMonitorGracePeriod,
		"podeviction":                DefaultKubernetesCtrlMgrPodEvictionTimeout,
		"routeperiod":                DefaultKubernetesCtrlMgrRouteReconciliationPeriod,
		"backoffretries":             strconv.Itoa(DefaultKubernetesCloudProviderBackoffRetries),
		"backoffjitter":              strconv.FormatFloat(DefaultKubernetesCloudProviderBackoffJitter, 'f', -1, 64),
		"backoffduration":            strconv.Itoa(DefaultKubernetesCloudProviderBackoffDuration),
		"backoffexponent":            strconv.FormatFloat(DefaultKubernetesCloudProviderBackoffExponent, 'f', -1, 64),
		"ratelimitqps":               strconv.FormatFloat(DefaultKubernetesCloudProviderRateLimitQPS, 'f', -1, 64),
		"ratelimitbucket":            strconv.Itoa(DefaultKubernetesCloudProviderRateLimitBucket),
		"gchighthreshold":            strconv.Itoa(DefaultKubernetesGCHighThreshold),
		"gclowthreshold":             strconv.Itoa(DefaultKubernetesGCLowThreshold),
	},
	common.KubernetesVersion1Dot6Dot6: {
		"hyperkube":                  "hyperkube-amd64:v1.6.6",
		"dockerEngineVersion":        "1.12.*",
		DefaultDashboardAddonName:    "kubernetes-dashboard-amd64:v1.6.3",
		"exechealthz":                "exechealthz-amd64:1.2",
		"addonresizer":               "addon-resizer:1.7",
		"heapster":                   "heapster-amd64:v1.3.0",
		"dns":                        "k8s-dns-kube-dns-amd64:1.14.4",
		"addonmanager":               "kube-addon-manager-amd64:v6.4-beta.2",
		"dnsmasq":                    "k8s-dns-dnsmasq-nanny-amd64:1.14.4",
		"pause":                      "pause-amd64:3.0",
		DefaultTillerAddonName:       DefaultTillerImage,
		DefaultACIConnectorAddonName: DefaultACIConnectorImage,
		DefaultReschedulerAddonName:  DefaultReschedulerImage,
		"nodestatusfreq":             DefaultKubernetesNodeStatusUpdateFrequency,
		"nodegraceperiod":            DefaultKubernetesCtrlMgrNodeMonitorGracePeriod,
		"podeviction":                DefaultKubernetesCtrlMgrPodEvictionTimeout,
		"routeperiod":                DefaultKubernetesCtrlMgrRouteReconciliationPeriod,
		"backoffretries":             strconv.Itoa(DefaultKubernetesCloudProviderBackoffRetries),
		"backoffjitter":              strconv.FormatFloat(DefaultKubernetesCloudProviderBackoffJitter, 'f', -1, 64),
		"backoffduration":            strconv.Itoa(DefaultKubernetesCloudProviderBackoffDuration),
		"backoffexponent":            strconv.FormatFloat(DefaultKubernetesCloudProviderBackoffExponent, 'f', -1, 64),
		"ratelimitqps":               strconv.FormatFloat(DefaultKubernetesCloudProviderRateLimitQPS, 'f', -1, 64),
		"ratelimitbucket":            strconv.Itoa(DefaultKubernetesCloudProviderRateLimitBucket),
		"gchighthreshold":            strconv.Itoa(DefaultKubernetesGCHighThreshold),
		"gclowthreshold":             strconv.Itoa(DefaultKubernetesGCLowThreshold),
	},
	common.KubernetesVersion1Dot5Dot8: {
		"hyperkube":                  "hyperkube-amd64:v1.5.8",
		"dockerEngineVersion":        "1.12.*",
		DefaultDashboardAddonName:    "kubernetes-dashboard-amd64:v1.5.1",
		"exechealthz":                "exechealthz-amd64:1.2",
		"addonresizer":               "addon-resizer:1.6",
		"heapster":                   "heapster:v1.2.0",
		"dns":                        "kubedns-amd64:1.7",
		"addonmanager":               "kube-addon-manager-amd64:v6.4-beta.2",
		"dnsmasq":                    "kube-dnsmasq-amd64:1.3",
		"pause":                      "pause-amd64:3.0",
		DefaultTillerAddonName:       "tiller:v2.5.1",
		DefaultACIConnectorAddonName: "virtual-kubelet:latest",
		"nodestatusfreq":             DefaultKubernetesNodeStatusUpdateFrequency,
		"nodegraceperiod":            DefaultKubernetesCtrlMgrNodeMonitorGracePeriod,
		"podeviction":                DefaultKubernetesCtrlMgrPodEvictionTimeout,
		"routeperiod":                DefaultKubernetesCtrlMgrRouteReconciliationPeriod,
		"gchighthreshold":            strconv.Itoa(DefaultKubernetesGCHighThreshold),
		"gclowthreshold":             strconv.Itoa(DefaultKubernetesGCLowThreshold),
	},
	common.KubernetesVersion1Dot5Dot7: {
		"hyperkube":                  "hyperkube-amd64:v1.5.7",
		"dockerEngineVersion":        "1.12.*",
		DefaultDashboardAddonName:    "kubernetes-dashboard-amd64:v1.5.1",
		"exechealthz":                "exechealthz-amd64:1.2",
		"addonresizer":               "addon-resizer:1.6",
		"heapster":                   "heapster:v1.2.0",
		"dns":                        "kubedns-amd64:1.7",
		"addonmanager":               "kube-addon-manager-amd64:v6.4-beta.2",
		"dnsmasq":                    "kube-dnsmasq-amd64:1.3",
		"pause":                      "pause-amd64:3.0",
		DefaultTillerAddonName:       "tiller:v2.5.1",
		DefaultACIConnectorAddonName: "virtual-kubelet:latest",
		"nodestatusfreq":             DefaultKubernetesNodeStatusUpdateFrequency,
		"nodegraceperiod":            DefaultKubernetesCtrlMgrNodeMonitorGracePeriod,
		"podeviction":                DefaultKubernetesCtrlMgrPodEvictionTimeout,
		"routeperiod":                DefaultKubernetesCtrlMgrRouteReconciliationPeriod,
		"gchighthreshold":            strconv.Itoa(DefaultKubernetesGCHighThreshold),
		"gclowthreshold":             strconv.Itoa(DefaultKubernetesGCLowThreshold),
	},
}
