package v20160330

// v20160330 supports orchestrators Mesos, Swarm, DCOS
const (
	Mesos OrchestratorType = "Mesos"
	Swarm OrchestratorType = "Swarm"
	DCOS  OrchestratorType = "DCOS"
)

// ACS supports orchestrators Mesos, Swarm, or DCOS.
const (
	Windows OSType = "Windows"
	Linux   OSType = "Linux"
)

// subscription states
const (
	// Registered means the subscription is entitled to use the namespace
	Registered SubscriptionState = iota
	// Unregistered means the subscription is not entitled to use the namespace
	Unregistered
	// Suspended means the subscription has been suspended from the system
	Suspended
	// Deleted means the subscription has been deleted
	Deleted
	// Warned means the subscription has been warned
	Warned
)

// storage profiles
const (
	// StorageExternal equates to VMSS where attached disks are unsupported (Default)
	StorageExternal = "External"
	// StorageVolumes equates to AS where attached disks are supported
	StorageVolumes = "Volumes"
	// StorageHAVolumes are managed disks that provide fault domain coverage for volumes.
	StorageHAVolumes = "HAVolumes"
)
