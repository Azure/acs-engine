package vlabs

const (
	// APIVersion is the version of this API
	APIVersion = "vlabs"
)

const (
	// Windows OSType
	Windows OSType = "Windows"
	// Linux OSType
	Linux OSType = "Linux"
)

// validation values
const (
	// MinAgentCount are the minimum number of agents
	MinAgentCount = 1
	// MaxAgentCount are the maximum number of agents
	MaxAgentCount = 100
	// MinDiskSizeGB specifies the minimum attached disk size
	MinDiskSizeGB = 1
	// MaxDiskSizeGB specifies the maximum attached disk size
	MaxDiskSizeGB = 1023
)

// storage profiles
const (
	// StorageAccount means that the nodes use raw storage accounts for their os and attached volumes
	StorageAccount = "StorageAccount"
	// ManagedDisks means that the nodes use managed disks for their os and attached volumes
	ManagedDisks = "ManagedDisks"
)

// distro values
const (
	// Supported distros by OpenShift
	OpenShift39RHEL Distro = "openshift39_rhel"
	OpenShiftCentOS Distro = "openshift39_centos"
)

// agent pool profiles
const (
	// AgentPoolProfileRoleEmpty is the empty role
	AgentPoolProfileRoleEmpty AgentPoolProfileRole = ""
	// AgentPoolProfileRoleInfra is the infra role
	AgentPoolProfileRoleInfra AgentPoolProfileRole = "infra"
)
