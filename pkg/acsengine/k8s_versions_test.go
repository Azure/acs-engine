package acsengine

import (
	"testing"
)

func TestGetK8sVersionComponents(t *testing.T) {

	oneDotElevenDotZero := getK8sVersionComponents("1.11.0-alpha.1", nil)
	if oneDotElevenDotZero == nil {
		t.Fatalf("getK8sVersionComponents() should not return nil for valid version")
	}
	expected := map[string]string{
		"hyperkube":                   "hyperkube-amd64:v1.11.0-alpha.1",
		"ccm":                         "cloud-controller-manager-amd64:v1.11.0-alpha.1",
		"windowszip":                  "v1.11.0-alpha.1-1int.zip",
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
		ContainerMonitoringAddonName: k8sComponentVersions["1.11"][ContainerMonitoringAddonName],
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

	for k, v := range oneDotElevenDotZero {
		if expected[k] != v {
			t.Fatalf("getK8sVersionComponents() returned an unexpected map[string]string value for k8s 1.11.0-alpha.1: %s = %s", k, oneDotElevenDotZero[k])
		}
	}

	oneDotNineDotThree := getK8sVersionComponents("1.9.3", nil)
	if oneDotNineDotThree == nil {
		t.Fatalf("getK8sVersionComponents() should not return nil for valid version")
	}
	expected = map[string]string{
		"hyperkube":                   "hyperkube-amd64:v1.9.3",
		"ccm":                         "cloud-controller-manager-amd64:v1.9.3",
		"windowszip":                  "v1.9.3-1int.zip",
		"dockerEngineVersion":         k8sComponentVersions["1.9"]["dockerEngine"],
		DefaultDashboardAddonName:     k8sComponentVersions["1.9"]["dashboard"],
		"exechealthz":                 k8sComponentVersions["1.9"]["exechealthz"],
		"addonresizer":                k8sComponentVersions["1.9"]["addon-resizer"],
		"heapster":                    k8sComponentVersions["1.9"]["heapster"],
		DefaultMetricsServerAddonName: k8sComponentVersions["1.9"]["metrics-server"],
		"dns":                             k8sComponentVersions["1.9"]["kube-dns"],
		"addonmanager":                    k8sComponentVersions["1.9"]["addon-manager"],
		"dnsmasq":                         k8sComponentVersions["1.9"]["dnsmasq"],
		"pause":                           k8sComponentVersions["1.9"]["pause"],
		DefaultTillerAddonName:            k8sComponentVersions["1.9"]["tiller"],
		DefaultReschedulerAddonName:       k8sComponentVersions["1.9"]["rescheduler"],
		DefaultACIConnectorAddonName:      k8sComponentVersions["1.9"]["aci-connector"],
		ContainerMonitoringAddonName:      k8sComponentVersions["1.11"][ContainerMonitoringAddonName],
		DefaultClusterAutoscalerAddonName: k8sComponentVersions["1.9"]["cluster-autoscaler"],
		"nodestatusfreq":                  k8sComponentVersions["1.9"]["nodestatusfreq"],
		"nodegraceperiod":                 k8sComponentVersions["1.9"]["nodegraceperiod"],
		"podeviction":                     k8sComponentVersions["1.9"]["podeviction"],
		"routeperiod":                     k8sComponentVersions["1.9"]["routeperiod"],
		"backoffretries":                  k8sComponentVersions["1.9"]["backoffretries"],
		"backoffjitter":                   k8sComponentVersions["1.9"]["backoffjitter"],
		"backoffduration":                 k8sComponentVersions["1.9"]["backoffduration"],
		"backoffexponent":                 k8sComponentVersions["1.9"]["backoffexponent"],
		"ratelimitqps":                    k8sComponentVersions["1.9"]["ratelimitqps"],
		"ratelimitbucket":                 k8sComponentVersions["1.9"]["ratelimitbucket"],
		"gchighthreshold":                 k8sComponentVersions["1.9"]["gchighthreshold"],
		"gclowthreshold":                  k8sComponentVersions["1.9"]["gclowthreshold"],
	}

	for k, v := range oneDotNineDotThree {
		if expected[k] != v {
			t.Fatalf("getK8sVersionComponents() returned an unexpected map[string]string value for k8s 1.9.3: %s = %s", k, oneDotNineDotThree[k])
		}
	}

	oneDotEightDotEight := getK8sVersionComponents("1.8.8", nil)
	if oneDotEightDotEight == nil {
		t.Fatalf("getK8sVersionComponents() should not return nil for valid version")
	}
	expected = map[string]string{
		"hyperkube":                   "hyperkube-amd64:v1.8.8",
		"ccm":                         "cloud-controller-manager-amd64:v1.8.8",
		"windowszip":                  "v1.8.8-1int.zip",
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
		ContainerMonitoringAddonName: k8sComponentVersions["1.11"][ContainerMonitoringAddonName],
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
	for k, v := range oneDotEightDotEight {
		if expected[k] != v {
			t.Fatalf("getK8sVersionComponents() returned an unexpected map[string]string value for k8s 1.8.8: %s = %s", k, oneDotNineDotThree[k])
		}
	}

	oneDotSevenDotZero := getK8sVersionComponents("1.7.13", nil)
	if oneDotSevenDotZero == nil {
		t.Fatalf("getK8sVersionComponents() should not return nil for valid version")
	}
	expected = map[string]string{
		"hyperkube":                   "hyperkube-amd64:v1.7.13",
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
		ContainerMonitoringAddonName: k8sComponentVersions["1.11"][ContainerMonitoringAddonName],
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
	for k, v := range oneDotSevenDotZero {
		if expected[k] != v {
			t.Fatalf("getK8sVersionComponents() returned an unexpected map[string]string value for k8s 1.7.0: %s = %s", k, oneDotSevenDotZero[k])
		}
	}

	override := getK8sVersionComponents("1.9.3", map[string]string{"windowszip": "v1.9.3-2int.zip", "dockerEngineVersion": "1.99.*"})
	if override == nil {
		t.Fatalf("getK8sVersionComponents() should not return nil for valid version")
	}
	expected = map[string]string{
		"hyperkube":                   "hyperkube-amd64:v1.9.3",
		"ccm":                         "cloud-controller-manager-amd64:v1.9.3",
		"windowszip":                  "v1.9.3-2int.zip",
		"dockerEngineVersion":         "1.99.*",
		DefaultDashboardAddonName:     k8sComponentVersions["1.9"]["dashboard"],
		"exechealthz":                 k8sComponentVersions["1.9"]["exechealthz"],
		"addonresizer":                k8sComponentVersions["1.9"]["addon-resizer"],
		"heapster":                    k8sComponentVersions["1.9"]["heapster"],
		DefaultMetricsServerAddonName: k8sComponentVersions["1.9"]["metrics-server"],
		"dns":                             k8sComponentVersions["1.9"]["kube-dns"],
		"addonmanager":                    k8sComponentVersions["1.9"]["addon-manager"],
		"dnsmasq":                         k8sComponentVersions["1.9"]["dnsmasq"],
		"pause":                           k8sComponentVersions["1.9"]["pause"],
		DefaultTillerAddonName:            k8sComponentVersions["1.9"]["tiller"],
		DefaultReschedulerAddonName:       k8sComponentVersions["1.9"]["rescheduler"],
		DefaultACIConnectorAddonName:      k8sComponentVersions["1.9"]["aci-connector"],
		ContainerMonitoringAddonName:      k8sComponentVersions["1.11"][ContainerMonitoringAddonName],
		DefaultClusterAutoscalerAddonName: k8sComponentVersions["1.9"]["cluster-autoscaler"],
		"nodestatusfreq":                  k8sComponentVersions["1.9"]["nodestatusfreq"],
		"nodegraceperiod":                 k8sComponentVersions["1.9"]["nodegraceperiod"],
		"podeviction":                     k8sComponentVersions["1.9"]["podeviction"],
		"routeperiod":                     k8sComponentVersions["1.9"]["routeperiod"],
		"backoffretries":                  k8sComponentVersions["1.9"]["backoffretries"],
		"backoffjitter":                   k8sComponentVersions["1.9"]["backoffjitter"],
		"backoffduration":                 k8sComponentVersions["1.9"]["backoffduration"],
		"backoffexponent":                 k8sComponentVersions["1.9"]["backoffexponent"],
		"ratelimitqps":                    k8sComponentVersions["1.9"]["ratelimitqps"],
		"ratelimitbucket":                 k8sComponentVersions["1.9"]["ratelimitbucket"],
		"gchighthreshold":                 k8sComponentVersions["1.9"]["gchighthreshold"],
		"gclowthreshold":                  k8sComponentVersions["1.9"]["gclowthreshold"],
	}
	for k, v := range override {
		if expected[k] != v {
			t.Fatalf("getK8sVersionComponents() returned an unexpected map[string]string value for k8s 1.9.3 w/ overrides: %s = %s", k, override[k])
		}
	}

	unknown := getK8sVersionComponents("1.0.0", nil)
	if unknown != nil {
		t.Fatalf("getK8sVersionComponents() should return nil for unknown k8s version")
	}
}
