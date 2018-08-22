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
	// OpenShift is the string constant for the OpenShift orchestrator type
	OpenShift string = "OpenShift"
)

// the OSTypes supported by vlabs
const (
	Windows OSType = "Windows"
	Linux   OSType = "Linux"
)

// the LinuxDistros supported by vlabs
const (
	Ubuntu Distro = "ubuntu"
	RHEL   Distro = "rhel"
	CoreOS Distro = "coreos"
	AKS    Distro = "aks"
	// Supported distros by OpenShift
	OpenShift39RHEL Distro = "openshift39_rhel"
	OpenShiftCentOS Distro = "openshift39_centos"
)

const (
	// SwarmVersion is the Swarm orchestrator version
	SwarmVersion = "swarm:1.1.0"
	// SwarmDockerComposeVersion is the Docker Compose version
	SwarmDockerComposeVersion = "1.6.2"
	// DockerCEVersion is the DockerCE orchestrator version
	DockerCEVersion = "17.03.*"
	// DockerCEDockerComposeVersion is the Docker Compose version
	DockerCEDockerComposeVersion = "1.14.0"
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
	// DefaultFirstConsecutiveKubernetesStaticIP specifies the static IP address on Kubernetes master 0
	DefaultFirstConsecutiveKubernetesStaticIP = "10.240.255.5"
	// DefaultFirstConsecutiveKubernetesStaticIPVMSS specifies the static IP address on Kubernetes master 0 of VMSS
	DefaultFirstConsecutiveKubernetesStaticIPVMSS = "10.240.0.4"
	// DefaultKubernetesFirstConsecutiveStaticIPOffset specifies the IP address offset of master 0
	// when VNET integration is enabled.
	DefaultKubernetesFirstConsecutiveStaticIPOffset = 5
	// DefaultKubernetesFirstConsecutiveStaticIPOffsetVMSS specifies the IP address offset of master 0 in VMSS
	// when VNET integration is enabled.
	DefaultKubernetesFirstConsecutiveStaticIPOffsetVMSS = 4
	// VirtualMachineScaleSets means that the vms are in a virtual machine scaleset
	VirtualMachineScaleSets = "VirtualMachineScaleSets"
	// ScaleSetPriorityRegular is the default ScaleSet Priority
	ScaleSetPriorityRegular = "Regular"
	// ScaleSetPriorityLow means the ScaleSet will use Low-priority VMs
	ScaleSetPriorityLow = "Low"
	// ScaleSetEvictionPolicyDelete is the default Eviction Policy for Low-priority VM ScaleSets
	ScaleSetEvictionPolicyDelete = "Delete"
	// ScaleSetEvictionPolicyDeallocate means a Low-priority VM ScaleSet will deallocate, rather than delete, VMs.
	ScaleSetEvictionPolicyDeallocate = "Deallocate"
)

// storage profiles
const (
	// StorageAccount means that the nodes use raw storage accounts for their os and attached volumes
	StorageAccount = "StorageAccount"
	// ManagedDisks means that the nodes use managed disks for their os and attached volumes
	ManagedDisks = "ManagedDisks"
)

// To identify programmatically generated public agent pools
const publicAgentPoolSuffix = "-public"

const (
	// DefaultTillerAddonEnabled determines the acs-engine provided default for enabling tiller addon
	DefaultTillerAddonEnabled = true
	// DefaultAADPodIdentityAddonEnabled determines the acs-engine provided default for enabling aad-pod-identity addon
	DefaultAADPodIdentityAddonEnabled = false
	// DefaultACIConnectorAddonEnabled determines the acs-engine provided default for enabling aci connector addon
	DefaultACIConnectorAddonEnabled = false
	// DefaultClusterAutoscalerAddonEnabled determines the acs-engine provided default for enabling cluster autoscaler addon
	DefaultClusterAutoscalerAddonEnabled = false
	// DefaultBlobfuseFlexVolumeAddonEnabled determines the acs-engine provided default for enabling blobfuse flexvolume addon
	DefaultBlobfuseFlexVolumeAddonEnabled = true
	// DefaultSMBFlexVolumeAddonEnabled determines the acs-engine provided default for enabling smb flexvolume addon
	DefaultSMBFlexVolumeAddonEnabled = true
	// DefaultKeyVaultFlexVolumeAddonEnabled determines the acs-engine provided default for enabling key vault flexvolume addon
	DefaultKeyVaultFlexVolumeAddonEnabled = true
	// DefaultDashboardAddonEnabled determines the acs-engine provided default for enabling kubernetes-dashboard addon
	DefaultDashboardAddonEnabled = true
	// DefaultReschedulerAddonEnabled determines the acs-engine provided default for enabling kubernetes-rescheduler addon
	DefaultReschedulerAddonEnabled = false
	// DefaultRBACEnabled determines the acs-engine provided default for enabling kubernetes RBAC
	DefaultRBACEnabled = true
	// DefaultUseInstanceMetadata determines the acs-engine provided default for enabling Azure cloudprovider instance metadata service
	DefaultUseInstanceMetadata = true
	// DefaultLoadBalancerSku determines the acs-engine provided default for enabling Azure cloudprovider load balancer SKU
	DefaultLoadBalancerSku = "Basic"
	// DefaultExcludeMasterFromStandardLB determines the acs-engine provided default for excluding master nodes from standard load balancer.
	DefaultExcludeMasterFromStandardLB = true
	// DefaultSecureKubeletEnabled determines the acs-engine provided default for securing kubelet communications
	DefaultSecureKubeletEnabled = true
	// DefaultMetricsServerAddonEnabled determines the acs-engine provided default for enabling kubernetes metrics-server addon
	DefaultMetricsServerAddonEnabled = false
	// DefaultNVIDIADevicePluginAddonEnabled determines the acs-engine provided default for enabling NVIDIA Device Plugin
	DefaultNVIDIADevicePluginAddonEnabled = false
	// DefaultContainerMonitoringAddonEnabled determines the acs-engine provided default for enabling kubernetes container monitoring addon
	DefaultContainerMonitoringAddonEnabled = false
	// DefaultAzureCNINetworkMonitoringAddonEnabled Azure CNI networkmonitor addon default
	DefaultAzureCNINetworkMonitoringAddonEnabled = false
	// DefaultTillerAddonName is the name of the tiller addon deployment
	DefaultTillerAddonName = "tiller"
	// DefaultAADPodIdentityAddonName is the name of the aad-pod-identity addon deployment
	DefaultAADPodIdentityAddonName = "aad-pod-identity"
	// DefaultACIConnectorAddonName is the name of the aci-connector addon deployment
	DefaultACIConnectorAddonName = "aci-connector"
	// DefaultClusterAutoscalerAddonName is the name of the cluster autoscaler addon deployment
	DefaultClusterAutoscalerAddonName = "cluster-autoscaler"
	// DefaultBlobfuseFlexVolumeAddonName is the name of the blobfuse flexvolume addon
	DefaultBlobfuseFlexVolumeAddonName = "blobfuse-flexvolume"
	// DefaultSMBFlexVolumeAddonName is the name of the smb flexvolume addon
	DefaultSMBFlexVolumeAddonName = "smb-flexvolume"
	// DefaultKeyVaultFlexVolumeAddonName is the name of the key vault flexvolume addon deployment
	DefaultKeyVaultFlexVolumeAddonName = "keyvault-flexvolume"
	// DefaultDashboardAddonName is the name of the kubernetes-dashboard addon deployment
	DefaultDashboardAddonName = "kubernetes-dashboard"
	// DefaultReschedulerAddonName is the name of the rescheduler addon deployment
	DefaultReschedulerAddonName = "rescheduler"
	// DefaultMetricsServerAddonName is the name of the kubernetes metrics server addon deployment
	DefaultMetricsServerAddonName = "metrics-server"
	// NVIDIADevicePluginAddonName is the name of the NVIDIA device plugin addon deployment
	NVIDIADevicePluginAddonName = "nvidia-device-plugin"
	// ContainerMonitoringAddonName is the name of the kubernetes Container Monitoring addon deployment
	ContainerMonitoringAddonName = "container-monitoring"
	// DefaultPrivateClusterEnabled determines the acs-engine provided default for enabling kubernetes Private Cluster
	DefaultPrivateClusterEnabled = false
	// NetworkPolicyAzure is the string expression for Azure CNI network policy manager
	NetworkPolicyAzure = "azure"
	// NetworkPolicyNone is the string expression for the deprecated NetworkPolicy usage pattern "none"
	NetworkPolicyNone = "none"
	// NetworkPluginKubenet is the string expression for the kubenet NetworkPlugin config
	NetworkPluginKubenet = "kubenet"
	// NetworkPluginAzure is thee string expression for Azure CNI plugin.
	NetworkPluginAzure = "azure"
	// DefaultSinglePlacementGroup determines the acs-engine provided default for supporting large VMSS
	// (true = single placement group 0-100 VMs, false = multiple placement group 0-1000 VMs)
	DefaultSinglePlacementGroup = true
)

const (
	// AgentPoolProfileRoleEmpty is the empty role.  Deprecated; only used in
	// acs-engine.
	AgentPoolProfileRoleEmpty AgentPoolProfileRole = ""
	// AgentPoolProfileRoleCompute is the compute role
	AgentPoolProfileRoleCompute AgentPoolProfileRole = "compute"
	// AgentPoolProfileRoleInfra is the infra role
	AgentPoolProfileRoleInfra AgentPoolProfileRole = "infra"
	// AgentPoolProfileRoleMaster is the master role
	AgentPoolProfileRoleMaster AgentPoolProfileRole = "master"
)

const (
	// VHDDiskSizeAKS maps to the OSDiskSizeGB for AKS VHD image
	VHDDiskSizeAKS = 100
)
