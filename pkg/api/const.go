package api

// the orchestrators supported by vlabs
const (
	// Mesos is the string constant for MESOS orchestrator type
	Mesos string = "Mesos"
	// DCOS is the string constant for DCOS orchestrator type and defaults to DCOS188
	DCOS string = "DCOS"
	// Swarm is the string constant for the Swarm orchestrator type
	Swarm string = "Swarm"
	// Kubernetes is the string constant for the Kubernetes orchestrator type
	Kubernetes string = "Kubernetes"
	// SwarmMode is the string constant for the Swarm Mode orchestrator type
	SwarmMode string = "SwarmMode"
)

// the OSTypes supported by vlabs
const (
	Windows OSType = "Windows"
	Linux   OSType = "Linux"
)

// validation values
const (
	// MinAgentCount are the minimum number of agents per agent pool
	MinAgentCount = 1
	// MaxAgentCount are the maximum number of agents per agent pool
	MaxAgentCount = 100
	// MinPort specifies the minimum tcp port to open
	MinPort = 1
	// MaxPort specifies the maximum tcp port to open
	MaxPort = 65535
	// MaxDisks specifies the maximum attached disks to add to the cluster
	MaxDisks = 4
)

// Availability profiles
const (
	// AvailabilitySet means that the vms are in an availability set
	AvailabilitySet = "AvailabilitySet"
	// VirtualMachineScaleSets means that the vms are in a virtual machine scaleset
	VirtualMachineScaleSets = "VirtualMachineScaleSets"
)

// storage profiles
const (
	// StorageAccount means that the nodes use raw storage accounts for their os and attached volumes
	StorageAccount = "StorageAccount"
	// ManagedDisks means that the nodes use managed disks for their os and attached volumes
	ManagedDisks = "ManagedDisks"
)

const (
	KubernetesVersionHint17      string = "1.7"
	KubernetesVersionHint16      string = "1.6"
	KubernetesVersionHint15      string = "1.5"
	KubernetesDefaultVersionHint string = KubernetesVersionHint16
)

const (
	DCOSVersionHint19      string = "1.9"
	DCOSVersionHint18      string = "1.8"
	DCOSVersionHint17      string = "1.7"
	DCOSDefaultVersionHint string = DCOSVersionHint19
)

// DCOSHintToVersion is the hint to actual version map
var DCOSHintToVersion = map[string]string{
	DCOSVersionHint19: "1.9.0",
	DCOSVersionHint18: "1.8.8",
	DCOSVersionHint17: "1.7.3",
}

// To identify programmatically generated public agent pools
const publicAgentPoolSuffix = "-public"

// KubeImages is the map from version hint to corresponding artifacts
var KubeImages = map[string]map[string]string{
	KubernetesVersionHint17: {
		"version":      "1.7.1",
		"hyperkube":    "hyperkube-amd64:v1.7.1",
		"dashboard":    "kubernetes-dashboard-amd64:v1.6.1",
		"exechealthz":  "exechealthz-amd64:1.2",
		"addonresizer": "addon-resizer:2.0",
		"heapster":     "heapster:v1.4.0",
		"dns":          "k8s-dns-kube-dns-amd64:1.14.4",
		"addonmanager": "kube-addon-manager-amd64:v6.4-beta.2",
		"dnsmasq":      "k8s-dns-dnsmasq-amd64:1.14.4",
		"pause":        "pause-amd64:3.0",
		"windowszip":   "v1.7.1intwinnat.zip",
	},
	KubernetesVersionHint16: {
		"version":      "1.6.6",
		"hyperkube":    "hyperkube-amd64:v1.6.6",
		"dashboard":    "kubernetes-dashboard-amd64:v1.6.1",
		"exechealthz":  "exechealthz-amd64:1.2",
		"addonresizer": "addon-resizer:1.7",
		"heapster":     "heapster:v1.3.0",
		"dns":          "k8s-dns-kube-dns-amd64:1.14.4",
		"addonmanager": "kube-addon-manager-amd64:v6.4-beta.2",
		"dnsmasq":      "k8s-dns-dnsmasq-amd64:1.13.0",
		"pause":        "pause-amd64:3.0",
		"windowszip":   "v1.6.6intwinnat.zip",
	},
	KubernetesVersionHint15: {
		"version":      "1.5.7",
		"hyperkube":    "hyperkube-amd64:v1.5.7",
		"dashboard":    "kubernetes-dashboard-amd64:v1.5.1",
		"exechealthz":  "exechealthz-amd64:1.2",
		"addonresizer": "addon-resizer:1.6",
		"heapster":     "heapster:v1.2.0",
		"dns":          "kubedns-amd64:1.7",
		"addonmanager": "kube-addon-manager-amd64:v6.2",
		"dnsmasq":      "kube-dnsmasq-amd64:1.3",
		"pause":        "pause-amd64:3.0",
		"windowszip":   "v1.5.7intwinnat.zip",
	},
}
