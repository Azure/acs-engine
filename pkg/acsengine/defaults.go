package acsengine

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net"
	"sort"
	"strconv"
	"strings"

	"github.com/Azure/acs-engine/pkg/api"
	"github.com/Azure/acs-engine/pkg/api/common"
	"github.com/Azure/acs-engine/pkg/helpers"
	"github.com/Azure/acs-engine/pkg/openshift/certgen"
	"github.com/blang/semver"
	"github.com/pkg/errors"
)

const (
	// AzureCniPluginVerLinux specifies version of Azure CNI plugin, which has been mirrored from
	// https://github.com/Azure/azure-container-networking/releases/download/${AZURE_PLUGIN_VER}/azure-vnet-cni-linux-amd64-${AZURE_PLUGIN_VER}.tgz
	// to https://acs-mirror.azureedge.net/cni
	AzureCniPluginVerLinux = "v1.0.11"
	// AzureCniPluginVerWindows specifies version of Azure CNI plugin, which has been mirrored from
	// https://github.com/Azure/azure-container-networking/releases/download/${AZURE_PLUGIN_VER}/azure-vnet-cni-windows-amd64-${AZURE_PLUGIN_VER}.tgz
	// to https://acs-mirror.azureedge.net/cni
	AzureCniPluginVerWindows = "v1.0.11"
)

// setPropertiesDefaults for the container Properties, returns true if certs are generated
func setPropertiesDefaults(cs *api.ContainerService, isUpgrade, isScale bool) (bool, error) {
	properties := cs.Properties

	setOrchestratorDefaults(cs, isUpgrade || isScale)

	// Set master profile defaults if this cluster configuration includes master node(s)
	if cs.Properties.MasterProfile != nil {
		setMasterProfileDefaults(properties, isUpgrade)
	}
	// Set VMSS Defaults for Masters
	if cs.Properties.MasterProfile != nil && cs.Properties.MasterProfile.IsVirtualMachineScaleSets() {
		setVMSSDefaultsForMasters(properties)
	}

	setAgentProfileDefaults(properties, isUpgrade, isScale)

	setStorageDefaults(properties)
	setExtensionDefaults(properties)
	// Set VMSS Defaults for Agents
	if cs.Properties.HasVMSSAgentPool() {
		setVMSSDefaultsForAgents(properties)
	}

	// Set hosted master profile defaults if this cluster configuration has a hosted control plane
	if cs.Properties.HostedMasterProfile != nil {
		setHostedMasterProfileDefaults(properties)
	}

	certsGenerated, e := setDefaultCerts(properties)
	if e != nil {
		return false, e
	}
	return certsGenerated, nil
}

// setOrchestratorDefaults for orchestrators
func setOrchestratorDefaults(cs *api.ContainerService, isUpdate bool) {
	a := cs.Properties

	cloudSpecConfig := cs.GetCloudSpecConfig()
	if a.OrchestratorProfile == nil {
		return
	}
	o := a.OrchestratorProfile
	o.OrchestratorVersion = common.GetValidPatchVersion(
		o.OrchestratorType,
		o.OrchestratorVersion, isUpdate, a.HasWindows())

	switch o.OrchestratorType {
	case api.Kubernetes:
		if o.KubernetesConfig == nil {
			o.KubernetesConfig = &api.KubernetesConfig{}
		}
		// For backwards compatibility with original, overloaded "NetworkPolicy" config vector
		// we translate deprecated NetworkPolicy usage to the NetworkConfig equivalent
		// and set a default network policy enforcement configuration
		switch o.KubernetesConfig.NetworkPolicy {
		case NetworkPluginAzure:
			if o.KubernetesConfig.NetworkPlugin == "" {
				o.KubernetesConfig.NetworkPlugin = NetworkPluginAzure
				o.KubernetesConfig.NetworkPolicy = DefaultNetworkPolicy
			}
		case NetworkPolicyNone:
			o.KubernetesConfig.NetworkPlugin = NetworkPluginKubenet
			o.KubernetesConfig.NetworkPolicy = DefaultNetworkPolicy
		case NetworkPolicyCalico:
			o.KubernetesConfig.NetworkPlugin = NetworkPluginKubenet
		case NetworkPolicyCilium:
			o.KubernetesConfig.NetworkPlugin = NetworkPolicyCilium
		}

		if o.KubernetesConfig.KubernetesImageBase == "" {
			o.KubernetesConfig.KubernetesImageBase = cloudSpecConfig.KubernetesSpecConfig.KubernetesImageBase
		}
		if o.KubernetesConfig.EtcdVersion == "" {
			o.KubernetesConfig.EtcdVersion = DefaultEtcdVersion
		}
		if a.HasWindows() {
			if o.KubernetesConfig.NetworkPlugin == "" {
				o.KubernetesConfig.NetworkPlugin = DefaultNetworkPluginWindows
			}
		} else {
			if o.KubernetesConfig.NetworkPlugin == "" {
				o.KubernetesConfig.NetworkPlugin = DefaultNetworkPlugin
			}
		}
		if o.KubernetesConfig.ContainerRuntime == "" {
			o.KubernetesConfig.ContainerRuntime = DefaultContainerRuntime
		}
		if o.KubernetesConfig.ClusterSubnet == "" {
			if o.IsAzureCNI() {
				// When Azure CNI is enabled, all masters, agents and pods share the same large subnet.
				// Except when master is VMSS, then masters and agents have separate subnets within the same large subnet.
				o.KubernetesConfig.ClusterSubnet = DefaultKubernetesSubnet
			} else {
				o.KubernetesConfig.ClusterSubnet = DefaultKubernetesClusterSubnet
			}
		}
		if o.KubernetesConfig.GCHighThreshold == 0 {
			o.KubernetesConfig.GCHighThreshold = DefaultKubernetesGCHighThreshold
		}
		if o.KubernetesConfig.GCLowThreshold == 0 {
			o.KubernetesConfig.GCLowThreshold = DefaultKubernetesGCLowThreshold
		}
		if o.KubernetesConfig.DNSServiceIP == "" {
			o.KubernetesConfig.DNSServiceIP = DefaultKubernetesDNSServiceIP
		}
		if o.KubernetesConfig.DockerBridgeSubnet == "" {
			o.KubernetesConfig.DockerBridgeSubnet = DefaultDockerBridgeSubnet
		}
		if o.KubernetesConfig.ServiceCIDR == "" {
			o.KubernetesConfig.ServiceCIDR = DefaultKubernetesServiceCIDR
		}

		if o.KubernetesConfig.CloudProviderBackoff == nil {
			o.KubernetesConfig.CloudProviderBackoff = helpers.PointerToBool(DefaultKubernetesCloudProviderBackoff)
		}
		// Enforce sane cloudprovider backoff defaults, if CloudProviderBackoff is true in KubernetesConfig
		if helpers.IsTrueBoolPointer(o.KubernetesConfig.CloudProviderBackoff) {
			o.KubernetesConfig.SetCloudProviderBackoffDefaults()
		}

		if o.KubernetesConfig.CloudProviderRateLimit == nil {
			o.KubernetesConfig.CloudProviderRateLimit = helpers.PointerToBool(DefaultKubernetesCloudProviderRateLimit)
		}
		// Enforce sane cloudprovider rate limit defaults, if CloudProviderRateLimit is true in KubernetesConfig
		if helpers.IsTrueBoolPointer(o.KubernetesConfig.CloudProviderRateLimit) {
			o.KubernetesConfig.SetCloudProviderRateLimitDefaults()
		}

		if o.KubernetesConfig.PrivateCluster == nil {
			o.KubernetesConfig.PrivateCluster = &api.PrivateCluster{}
		}

		if o.KubernetesConfig.PrivateCluster.Enabled == nil {
			o.KubernetesConfig.PrivateCluster.Enabled = helpers.PointerToBool(api.DefaultPrivateClusterEnabled)
		}

		if "" == a.OrchestratorProfile.KubernetesConfig.EtcdDiskSizeGB {
			switch {
			case a.TotalNodes() > 20:
				a.OrchestratorProfile.KubernetesConfig.EtcdDiskSizeGB = DefaultEtcdDiskSizeGT20Nodes
			case a.TotalNodes() > 10:
				a.OrchestratorProfile.KubernetesConfig.EtcdDiskSizeGB = DefaultEtcdDiskSizeGT10Nodes
			case a.TotalNodes() > 3:
				a.OrchestratorProfile.KubernetesConfig.EtcdDiskSizeGB = DefaultEtcdDiskSizeGT3Nodes
			default:
				a.OrchestratorProfile.KubernetesConfig.EtcdDiskSizeGB = DefaultEtcdDiskSize
			}
		}

		if helpers.IsTrueBoolPointer(o.KubernetesConfig.EnableDataEncryptionAtRest) {
			if "" == a.OrchestratorProfile.KubernetesConfig.EtcdEncryptionKey {
				a.OrchestratorProfile.KubernetesConfig.EtcdEncryptionKey = generateEtcdEncryptionKey()
			}
		}

		if a.OrchestratorProfile.KubernetesConfig.PrivateJumpboxProvision() && a.OrchestratorProfile.KubernetesConfig.PrivateCluster.JumpboxProfile.OSDiskSizeGB == 0 {
			a.OrchestratorProfile.KubernetesConfig.PrivateCluster.JumpboxProfile.OSDiskSizeGB = DefaultJumpboxDiskSize
		}

		if a.OrchestratorProfile.KubernetesConfig.PrivateJumpboxProvision() && a.OrchestratorProfile.KubernetesConfig.PrivateCluster.JumpboxProfile.Username == "" {
			a.OrchestratorProfile.KubernetesConfig.PrivateCluster.JumpboxProfile.Username = DefaultJumpboxUsername
		}

		if a.OrchestratorProfile.KubernetesConfig.PrivateJumpboxProvision() && a.OrchestratorProfile.KubernetesConfig.PrivateCluster.JumpboxProfile.StorageProfile == "" {
			a.OrchestratorProfile.KubernetesConfig.PrivateCluster.JumpboxProfile.StorageProfile = api.ManagedDisks
		}

		if !helpers.IsFalseBoolPointer(a.OrchestratorProfile.KubernetesConfig.EnableRbac) {
			if common.IsKubernetesVersionGe(o.OrchestratorVersion, "1.9.0") {
				// TODO make EnableAggregatedAPIs a pointer to bool so that a user can opt out of it
				a.OrchestratorProfile.KubernetesConfig.EnableAggregatedAPIs = true
			}
			if a.OrchestratorProfile.KubernetesConfig.EnableRbac == nil {
				a.OrchestratorProfile.KubernetesConfig.EnableRbac = helpers.PointerToBool(api.DefaultRBACEnabled)
			}
		}

		if a.OrchestratorProfile.KubernetesConfig.EnableSecureKubelet == nil {
			a.OrchestratorProfile.KubernetesConfig.EnableSecureKubelet = helpers.PointerToBool(api.DefaultSecureKubeletEnabled)
		}

		if a.OrchestratorProfile.KubernetesConfig.UseInstanceMetadata == nil {
			a.OrchestratorProfile.KubernetesConfig.UseInstanceMetadata = helpers.PointerToBool(api.DefaultUseInstanceMetadata)
		}

		if !a.HasAvailabilityZones() && a.OrchestratorProfile.KubernetesConfig.LoadBalancerSku == "" {
			a.OrchestratorProfile.KubernetesConfig.LoadBalancerSku = api.DefaultLoadBalancerSku
		}

		if common.IsKubernetesVersionGe(a.OrchestratorProfile.OrchestratorVersion, "1.11.0") && a.OrchestratorProfile.KubernetesConfig.LoadBalancerSku == "Standard" && a.OrchestratorProfile.KubernetesConfig.ExcludeMasterFromStandardLB == nil {
			a.OrchestratorProfile.KubernetesConfig.ExcludeMasterFromStandardLB = helpers.PointerToBool(api.DefaultExcludeMasterFromStandardLB)
		}

		if a.OrchestratorProfile.IsAzureCNI() {
			if a.HasWindows() {
				a.OrchestratorProfile.KubernetesConfig.AzureCNIVersion = AzureCniPluginVerWindows
			} else {
				a.OrchestratorProfile.KubernetesConfig.AzureCNIVersion = AzureCniPluginVerLinux
			}
		}

		// Configure addons
		setAddonsConfig(cs)
		// Configure kubelet
		setKubeletConfig(cs)
		// Configure controller-manager
		setControllerManagerConfig(cs)
		// Configure cloud-controller-manager
		setCloudControllerManagerConfig(cs)
		// Configure apiserver
		setAPIServerConfig(cs)
		// Configure scheduler
		setSchedulerConfig(cs)

	case api.DCOS:
		if o.DcosConfig == nil {
			o.DcosConfig = &api.DcosConfig{}
		}
		dcosSemVer, _ := semver.Make(o.OrchestratorVersion)
		dcosBootstrapSemVer, _ := semver.Make(common.DCOSVersion1Dot11Dot0)
		if !dcosSemVer.LT(dcosBootstrapSemVer) {
			if o.DcosConfig.BootstrapProfile == nil {
				o.DcosConfig.BootstrapProfile = &api.BootstrapProfile{}
			}
			if len(o.DcosConfig.BootstrapProfile.VMSize) == 0 {
				o.DcosConfig.BootstrapProfile.VMSize = "Standard_D2s_v3"
			}
		}
	case api.OpenShift:
		kc := a.OrchestratorProfile.OpenShiftConfig.KubernetesConfig
		if kc == nil {
			kc = &api.KubernetesConfig{}
		}
		if kc.ContainerRuntime == "" {
			kc.ContainerRuntime = DefaultContainerRuntime
		}
		if kc.NetworkPlugin == "" {
			kc.NetworkPlugin = DefaultNetworkPlugin
		}
	}
}

func setExtensionDefaults(a *api.Properties) {
	if a.ExtensionProfiles == nil {
		return
	}
	for _, extension := range a.ExtensionProfiles {
		if extension.RootURL == "" {
			extension.RootURL = DefaultExtensionsRootURL
		}
	}
}

func setMasterProfileDefaults(a *api.Properties, isUpgrade bool) {
	if a.MasterProfile.Distro == "" {
		if a.OrchestratorProfile.IsKubernetes() {
			a.MasterProfile.Distro = api.AKS
		} else if !a.OrchestratorProfile.IsOpenShift() {
			a.MasterProfile.Distro = api.Ubuntu
		}
	}
	// set default to VMAS for now
	if len(a.MasterProfile.AvailabilityProfile) == 0 {
		a.MasterProfile.AvailabilityProfile = api.AvailabilitySet
	}

	if !a.MasterProfile.IsCustomVNET() {
		if a.OrchestratorProfile.OrchestratorType == api.Kubernetes {
			if a.OrchestratorProfile.IsAzureCNI() {
				// When VNET integration is enabled, all masters, agents and pods share the same large subnet.
				a.MasterProfile.Subnet = a.OrchestratorProfile.KubernetesConfig.ClusterSubnet
				// FirstConsecutiveStaticIP is not reset if it is upgrade and some value already exists
				if !isUpgrade || len(a.MasterProfile.FirstConsecutiveStaticIP) == 0 {
					if a.MasterProfile.IsVirtualMachineScaleSets() {
						a.MasterProfile.FirstConsecutiveStaticIP = api.DefaultFirstConsecutiveKubernetesStaticIPVMSS
						a.MasterProfile.Subnet = DefaultKubernetesMasterSubnet
						a.MasterProfile.AgentSubnet = DefaultKubernetesAgentSubnetVMSS
					} else {
						a.MasterProfile.FirstConsecutiveStaticIP = a.MasterProfile.GetFirstConsecutiveStaticIPAddress(a.MasterProfile.Subnet)
					}
				}
			} else {
				a.MasterProfile.Subnet = DefaultKubernetesMasterSubnet
				// FirstConsecutiveStaticIP is not reset if it is upgrade and some value already exists
				if !isUpgrade || len(a.MasterProfile.FirstConsecutiveStaticIP) == 0 {
					if a.MasterProfile.IsVirtualMachineScaleSets() {
						a.MasterProfile.FirstConsecutiveStaticIP = api.DefaultFirstConsecutiveKubernetesStaticIPVMSS
						a.MasterProfile.AgentSubnet = DefaultKubernetesAgentSubnetVMSS
					} else {
						a.MasterProfile.FirstConsecutiveStaticIP = api.DefaultFirstConsecutiveKubernetesStaticIP
					}
				}
			}
		} else if a.OrchestratorProfile.OrchestratorType == api.OpenShift {
			a.MasterProfile.Subnet = DefaultOpenShiftMasterSubnet
			if !isUpgrade || len(a.MasterProfile.FirstConsecutiveStaticIP) == 0 {
				a.MasterProfile.FirstConsecutiveStaticIP = DefaultOpenShiftFirstConsecutiveStaticIP
			}
		} else if a.OrchestratorProfile.OrchestratorType == api.DCOS {
			a.MasterProfile.Subnet = DefaultDCOSMasterSubnet
			// FirstConsecutiveStaticIP is not reset if it is upgrade and some value already exists
			if !isUpgrade || len(a.MasterProfile.FirstConsecutiveStaticIP) == 0 {
				a.MasterProfile.FirstConsecutiveStaticIP = DefaultDCOSFirstConsecutiveStaticIP
			}
			if a.OrchestratorProfile.DcosConfig != nil && a.OrchestratorProfile.DcosConfig.BootstrapProfile != nil {
				if !isUpgrade || len(a.OrchestratorProfile.DcosConfig.BootstrapProfile.StaticIP) == 0 {
					a.OrchestratorProfile.DcosConfig.BootstrapProfile.StaticIP = DefaultDCOSBootstrapStaticIP
				}
			}
		} else if a.HasWindows() {
			a.MasterProfile.Subnet = DefaultSwarmWindowsMasterSubnet
			// FirstConsecutiveStaticIP is not reset if it is upgrade and some value already exists
			if !isUpgrade || len(a.MasterProfile.FirstConsecutiveStaticIP) == 0 {
				a.MasterProfile.FirstConsecutiveStaticIP = DefaultSwarmWindowsFirstConsecutiveStaticIP
			}
		} else {
			a.MasterProfile.Subnet = DefaultMasterSubnet
			// FirstConsecutiveStaticIP is not reset if it is upgrade and some value already exists
			if !isUpgrade || len(a.MasterProfile.FirstConsecutiveStaticIP) == 0 {
				a.MasterProfile.FirstConsecutiveStaticIP = DefaultFirstConsecutiveStaticIP
			}
		}
	}

	if a.MasterProfile.IsCustomVNET() && a.MasterProfile.IsVirtualMachineScaleSets() {
		if a.OrchestratorProfile.OrchestratorType == api.Kubernetes {
			a.MasterProfile.FirstConsecutiveStaticIP = a.MasterProfile.GetFirstConsecutiveStaticIPAddress(a.MasterProfile.VnetCidr)
		}
	}
	// Set the default number of IP addresses allocated for masters.
	if a.MasterProfile.IPAddressCount == 0 {
		// Allocate one IP address for the node.
		a.MasterProfile.IPAddressCount = 1

		// Allocate IP addresses for pods if VNET integration is enabled.
		if a.OrchestratorProfile.IsAzureCNI() {
			if a.OrchestratorProfile.OrchestratorType == api.Kubernetes {
				masterMaxPods, _ := strconv.Atoi(a.MasterProfile.KubernetesConfig.KubeletConfig["--max-pods"])
				a.MasterProfile.IPAddressCount += masterMaxPods
			}
		}
	}

	if a.MasterProfile.HTTPSourceAddressPrefix == "" {
		a.MasterProfile.HTTPSourceAddressPrefix = "*"
	}
}

// setVMSSDefaultsForMasters
func setVMSSDefaultsForMasters(a *api.Properties) {
	if a.MasterProfile.SinglePlacementGroup == nil {
		a.MasterProfile.SinglePlacementGroup = helpers.PointerToBool(api.DefaultSinglePlacementGroup)
	}
	if a.MasterProfile.HasAvailabilityZones() && (a.OrchestratorProfile.KubernetesConfig != nil && a.OrchestratorProfile.KubernetesConfig.LoadBalancerSku == "") {
		a.OrchestratorProfile.KubernetesConfig.LoadBalancerSku = "Standard"
		a.OrchestratorProfile.KubernetesConfig.ExcludeMasterFromStandardLB = helpers.PointerToBool(api.DefaultExcludeMasterFromStandardLB)
	}
}

// setVMSSDefaultsForAgents
func setVMSSDefaultsForAgents(a *api.Properties) {
	for _, profile := range a.AgentPoolProfiles {
		if profile.AvailabilityProfile == api.VirtualMachineScaleSets {
			if profile.Count > 100 {
				profile.SinglePlacementGroup = helpers.PointerToBool(false)
			}
			if profile.SinglePlacementGroup == nil {
				profile.SinglePlacementGroup = helpers.PointerToBool(api.DefaultSinglePlacementGroup)
			}
			if profile.HasAvailabilityZones() && (a.OrchestratorProfile.KubernetesConfig != nil && a.OrchestratorProfile.KubernetesConfig.LoadBalancerSku == "") {
				a.OrchestratorProfile.KubernetesConfig.LoadBalancerSku = "Standard"
				a.OrchestratorProfile.KubernetesConfig.ExcludeMasterFromStandardLB = helpers.PointerToBool(api.DefaultExcludeMasterFromStandardLB)
			}
		}

	}
}

func setAgentProfileDefaults(a *api.Properties, isUpgrade, isScale bool) {
	// configure the subnets if not in custom VNET
	if a.MasterProfile != nil && !a.MasterProfile.IsCustomVNET() {
		subnetCounter := 0
		for _, profile := range a.AgentPoolProfiles {
			if a.OrchestratorProfile.OrchestratorType == api.Kubernetes ||
				a.OrchestratorProfile.OrchestratorType == api.OpenShift {
				if !a.MasterProfile.IsVirtualMachineScaleSets() {
					profile.Subnet = a.MasterProfile.Subnet
				}
			} else {
				profile.Subnet = fmt.Sprintf(DefaultAgentSubnetTemplate, subnetCounter)
			}

			subnetCounter++
		}
	}

	for _, profile := range a.AgentPoolProfiles {
		// set default OSType to Linux
		if profile.OSType == "" {
			profile.OSType = api.Linux
		}

		// Accelerated Networking is supported on most general purpose and compute-optimized instance sizes with 2 or more vCPUs.
		// These supported series are: D/DSv2 and F/Fs // All the others are not supported
		// On instances that support hyperthreading, Accelerated Networking is supported on VM instances with 4 or more vCPUs.
		// Supported series are: D/DSv3, E/ESv3, Fsv2, and Ms/Mms.
		if profile.AcceleratedNetworkingEnabled == nil {
			profile.AcceleratedNetworkingEnabled = helpers.PointerToBool(!isUpgrade && !isScale && helpers.AcceleratedNetworkingSupported(profile.VMSize))
		}

		if profile.AcceleratedNetworkingEnabledWindows == nil {
			profile.AcceleratedNetworkingEnabledWindows = helpers.PointerToBool(api.DefaultAcceleratedNetworkingWindowsEnabled)
		}

		if profile.Distro == "" {
			if a.OrchestratorProfile.IsKubernetes() {
				if profile.OSDiskSizeGB != 0 && profile.OSDiskSizeGB < api.VHDDiskSizeAKS {
					profile.Distro = api.Ubuntu
				} else {
					profile.Distro = api.AKS
				}
			} else if !a.OrchestratorProfile.IsOpenShift() {
				profile.Distro = api.Ubuntu
			}
		}

		// Set the default number of IP addresses allocated for agents.
		if profile.IPAddressCount == 0 {
			// Allocate one IP address for the node.
			profile.IPAddressCount = 1

			// Allocate IP addresses for pods if VNET integration is enabled.
			if a.OrchestratorProfile.IsAzureCNI() {
				agentPoolMaxPods, _ := strconv.Atoi(profile.KubernetesConfig.KubeletConfig["--max-pods"])
				profile.IPAddressCount += agentPoolMaxPods
			}
		}
	}
}

// setStorageDefaults for agents
func setStorageDefaults(a *api.Properties) {
	if a.MasterProfile != nil && len(a.MasterProfile.StorageProfile) == 0 {
		if a.OrchestratorProfile.OrchestratorType == api.Kubernetes {
			a.MasterProfile.StorageProfile = api.ManagedDisks
		} else {
			a.MasterProfile.StorageProfile = api.StorageAccount
		}
	}
	for _, profile := range a.AgentPoolProfiles {
		if len(profile.StorageProfile) == 0 {
			if a.OrchestratorProfile.OrchestratorType == api.Kubernetes {
				profile.StorageProfile = api.ManagedDisks
			} else {
				profile.StorageProfile = api.StorageAccount
			}
		}
		if len(profile.AvailabilityProfile) == 0 {
			profile.AvailabilityProfile = api.VirtualMachineScaleSets
			// VMSS is not supported for k8s below 1.10.2
			if a.OrchestratorProfile.OrchestratorType == api.Kubernetes && !common.IsKubernetesVersionGe(a.OrchestratorProfile.OrchestratorVersion, "1.10.2") {
				profile.AvailabilityProfile = api.AvailabilitySet
			}
		}
		if len(profile.ScaleSetEvictionPolicy) == 0 && profile.ScaleSetPriority == api.ScaleSetPriorityLow {
			profile.ScaleSetEvictionPolicy = api.ScaleSetEvictionPolicyDelete
		}
	}
}

func setHostedMasterProfileDefaults(a *api.Properties) {
	a.HostedMasterProfile.Subnet = DefaultKubernetesMasterSubnet
}

func setDefaultCerts(p *api.Properties) (bool, error) {
	if p.MasterProfile != nil && p.OrchestratorProfile.OrchestratorType == api.OpenShift {
		return certgen.OpenShiftSetDefaultCerts(p, api.DefaultOpenshiftOrchestratorName, p.GetClusterID())
	}

	if p.MasterProfile == nil || p.OrchestratorProfile.OrchestratorType != api.Kubernetes {
		return false, nil
	}

	provided := certsAlreadyPresent(p.CertificateProfile, p.MasterProfile.Count)

	if areAllTrue(provided) {
		return false, nil
	}

	masterExtraFQDNs := append(formatAzureProdFQDNs(p.MasterProfile.DNSPrefix), p.MasterProfile.SubjectAltNames...)
	firstMasterIP := net.ParseIP(p.MasterProfile.FirstConsecutiveStaticIP).To4()

	if firstMasterIP == nil {
		return false, errors.Errorf("MasterProfile.FirstConsecutiveStaticIP '%s' is an invalid IP address", p.MasterProfile.FirstConsecutiveStaticIP)
	}

	ips := []net.IP{firstMasterIP}
	// Add the Internal Loadbalancer IP which is always at at p known offset from the firstMasterIP
	ips = append(ips, net.IP{firstMasterIP[0], firstMasterIP[1], firstMasterIP[2], firstMasterIP[3] + byte(DefaultInternalLbStaticIPOffset)})
	// Include the Internal load balancer as well

	if p.MasterProfile.IsVirtualMachineScaleSets() {
		// Include the Internal load balancer as well
		for i := 1; i < p.MasterProfile.Count; i++ {
			offset := i * p.MasterProfile.IPAddressCount
			ip := net.IP{firstMasterIP[0], firstMasterIP[1], firstMasterIP[2], firstMasterIP[3] + byte(offset)}
			ips = append(ips, ip)
		}
	} else {
		for i := 1; i < p.MasterProfile.Count; i++ {
			ip := net.IP{firstMasterIP[0], firstMasterIP[1], firstMasterIP[2], firstMasterIP[3] + byte(i)}
			ips = append(ips, ip)
		}
	}
	if p.CertificateProfile == nil {
		p.CertificateProfile = &api.CertificateProfile{}
	}

	// use the specified Certificate Authority pair, or generate p new pair
	var caPair *helpers.PkiKeyCertPair
	if provided["ca"] {
		caPair = &helpers.PkiKeyCertPair{CertificatePem: p.CertificateProfile.CaCertificate, PrivateKeyPem: p.CertificateProfile.CaPrivateKey}
	} else {
		var err error
		caPair, err = helpers.CreatePkiKeyCertPair("ca")
		if err != nil {
			return false, err
		}
		p.CertificateProfile.CaCertificate = caPair.CertificatePem
		p.CertificateProfile.CaPrivateKey = caPair.PrivateKeyPem
	}

	cidrFirstIP, err := common.CidrStringFirstIP(p.OrchestratorProfile.KubernetesConfig.ServiceCIDR)
	if err != nil {
		return false, err
	}
	ips = append(ips, cidrFirstIP)

	apiServerPair, clientPair, kubeConfigPair, etcdServerPair, etcdClientPair, etcdPeerPairs, err := helpers.CreatePki(masterExtraFQDNs, ips, DefaultKubernetesClusterDomain, caPair, p.MasterProfile.Count)
	if err != nil {
		return false, err
	}

	// If no Certificate Authority pair or no cert/key pair was provided, use generated cert/key pairs signed by provided Certificate Authority pair
	if !provided["apiserver"] || !provided["ca"] {
		p.CertificateProfile.APIServerCertificate = apiServerPair.CertificatePem
		p.CertificateProfile.APIServerPrivateKey = apiServerPair.PrivateKeyPem
	}
	if !provided["client"] || !provided["ca"] {
		p.CertificateProfile.ClientCertificate = clientPair.CertificatePem
		p.CertificateProfile.ClientPrivateKey = clientPair.PrivateKeyPem
	}
	if !provided["kubeconfig"] || !provided["ca"] {
		p.CertificateProfile.KubeConfigCertificate = kubeConfigPair.CertificatePem
		p.CertificateProfile.KubeConfigPrivateKey = kubeConfigPair.PrivateKeyPem
	}
	if !provided["etcd"] || !provided["ca"] {
		p.CertificateProfile.EtcdServerCertificate = etcdServerPair.CertificatePem
		p.CertificateProfile.EtcdServerPrivateKey = etcdServerPair.PrivateKeyPem
		p.CertificateProfile.EtcdClientCertificate = etcdClientPair.CertificatePem
		p.CertificateProfile.EtcdClientPrivateKey = etcdClientPair.PrivateKeyPem
		p.CertificateProfile.EtcdPeerCertificates = make([]string, p.MasterProfile.Count)
		p.CertificateProfile.EtcdPeerPrivateKeys = make([]string, p.MasterProfile.Count)
		for i, v := range etcdPeerPairs {
			p.CertificateProfile.EtcdPeerCertificates[i] = v.CertificatePem
			p.CertificateProfile.EtcdPeerPrivateKeys[i] = v.PrivateKeyPem
		}
	}

	return true, nil
}

func areAllTrue(m map[string]bool) bool {
	for _, v := range m {
		if !v {
			return false
		}
	}
	return true
}

// certsAlreadyPresent already present returns a map where each key is a type of cert and each value is true if that cert/key pair is user-provided
func certsAlreadyPresent(c *api.CertificateProfile, m int) map[string]bool {
	g := map[string]bool{
		"ca":         false,
		"apiserver":  false,
		"kubeconfig": false,
		"client":     false,
		"etcd":       false,
	}
	if c != nil {
		etcdPeer := true
		if len(c.EtcdPeerCertificates) != m || len(c.EtcdPeerPrivateKeys) != m {
			etcdPeer = false
		} else {
			for i, p := range c.EtcdPeerCertificates {
				if !(len(p) > 0) || !(len(c.EtcdPeerPrivateKeys[i]) > 0) {
					etcdPeer = false
				}
			}
		}
		g["ca"] = len(c.CaCertificate) > 0 && len(c.CaPrivateKey) > 0
		g["apiserver"] = len(c.APIServerCertificate) > 0 && len(c.APIServerPrivateKey) > 0
		g["kubeconfig"] = len(c.KubeConfigCertificate) > 0 && len(c.KubeConfigPrivateKey) > 0
		g["client"] = len(c.ClientCertificate) > 0 && len(c.ClientPrivateKey) > 0
		g["etcd"] = etcdPeer && len(c.EtcdClientCertificate) > 0 && len(c.EtcdClientPrivateKey) > 0 && len(c.EtcdServerCertificate) > 0 && len(c.EtcdServerPrivateKey) > 0
	}
	return g
}

// combine user-provided --feature-gates vals with defaults
// a minimum k8s version may be declared as required for defaults assignment
func addDefaultFeatureGates(m map[string]string, version string, minVersion string, defaults string) {
	if minVersion != "" {
		if common.IsKubernetesVersionGe(version, minVersion) {
			m["--feature-gates"] = combineValues(m["--feature-gates"], defaults)
		} else {
			m["--feature-gates"] = combineValues(m["--feature-gates"], "")
		}
	} else {
		m["--feature-gates"] = combineValues(m["--feature-gates"], defaults)
	}
}

func combineValues(inputs ...string) string {
	valueMap := make(map[string]string)
	for _, input := range inputs {
		applyValueStringToMap(valueMap, input)
	}
	return mapToString(valueMap)
}

func applyValueStringToMap(valueMap map[string]string, input string) {
	values := strings.Split(input, ",")
	for index := 0; index < len(values); index++ {
		// trim spaces (e.g. if the input was "foo=true, bar=true" - we want to drop the space after the comma)
		value := strings.Trim(values[index], " ")
		valueParts := strings.Split(value, "=")
		if len(valueParts) == 2 {
			valueMap[valueParts[0]] = valueParts[1]
		}
	}
}

func mapToString(valueMap map[string]string) string {
	// Order by key for consistency
	keys := []string{}
	for key := range valueMap {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	var buf bytes.Buffer
	for _, key := range keys {
		buf.WriteString(fmt.Sprintf("%s=%s,", key, valueMap[key]))
	}
	return strings.TrimSuffix(buf.String(), ",")
}

func generateEtcdEncryptionKey() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)
}
