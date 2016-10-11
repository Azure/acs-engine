package vlabs

const (
	// DCOS is the string constant for DCOS orchestrator type and defaults to DCOS184
	DCOS = "DCOS"
	// DCOS184 is the string constant for DCOS 1.8.4 orchestrator type
	DCOS184 = "DCOS184"
	// DCOS173 is the string constant for DCOS 1.7.3 orchestrator type
	DCOS173 = "DCOS173"
	// Swarm is the string constant for the Swarm orchestrator type
	Swarm = "Swarm"
	// Kubernetes is the string constant for the Kubernetes orchestrator type
	Kubernetes = "Kubernetes"
	// MinAgentCount are the minimum number of agents
	MinAgentCount = 1
	// MaxAgentCount are the maximum number of agents
	MaxAgentCount = 100
	// MinPort specifies the minimum tcp port to open
	MinPort = 1
	// MaxPort specifies the maximum tcp port to open
	MaxPort = 65535
	// MaxDisks specifies the maximum attached disks to add to the cluster
	MaxDisks = 4
	// StorageExternal equates to VMSS where attached disks are unsupported (Default)
	StorageExternal = "External"
	// StorageVolumes equates to AS where attached disks are supported
	StorageVolumes = "Volumes"
	// StorageHAVolumes are managed disks that provide fault domain coverage for volumes.
	StorageHAVolumes = "HAVolumes"
	// OSTypeWindows specifies the Windows OS
	OSTypeWindows = "Windows"
	// OSTypeLinux specifies the Linux OS
	OSTypeLinux = "Linux"
)
