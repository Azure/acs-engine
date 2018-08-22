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
	// CNIPluginVer specifies the version of CNI implementation
	// https://github.com/containernetworking/plugins
	CNIPluginVer = "v0.7.1"
)

var (
	//DefaultKubernetesSpecConfig is the default Docker image source of Kubernetes
	DefaultKubernetesSpecConfig = KubernetesSpecConfig{
		KubernetesImageBase:              "k8s.gcr.io/",
		TillerImageBase:                  "gcr.io/kubernetes-helm/",
		ACIConnectorImageBase:            "microsoft/",
		NVIDIAImageBase:                  "nvidia/",
		AzureCNIImageBase:                "containernetworking/",
		EtcdDownloadURLBase:              "https://acs-mirror.azureedge.net/github-coreos",
		KubeBinariesSASURLBase:           "https://acs-mirror.azureedge.net/wink8s/",
		WindowsPackageSASURLBase:         "https://acs-mirror.azureedge.net/wink8s/",
		WindowsTelemetryGUID:             "fb801154-36b9-41bc-89c2-f4d4f05472b0",
		CNIPluginsDownloadURL:            "https://acs-mirror.azureedge.net/cni/cni-plugins-amd64-" + CNIPluginVer + ".tgz",
		VnetCNILinuxPluginsDownloadURL:   "https://acs-mirror.azureedge.net/cni/azure-vnet-cni-linux-amd64-" + AzureCniPluginVerLinux + ".tgz",
		VnetCNIWindowsPluginsDownloadURL: "https://acs-mirror.azureedge.net/cni/azure-vnet-cni-windows-amd64-" + AzureCniPluginVerWindows + ".zip",
		ContainerdDownloadURLBase:        "https://storage.googleapis.com/cri-containerd-release/",
	}

	//DefaultDCOSSpecConfig is the default DC/OS binary download URL.
	DefaultDCOSSpecConfig = DCOSSpecConfig{
		DCOS188BootstrapDownloadURL:     fmt.Sprintf(AzureEdgeDCOSBootstrapDownloadURL, "stable", "5df43052907c021eeb5de145419a3da1898c58a5"),
		DCOS190BootstrapDownloadURL:     fmt.Sprintf(AzureEdgeDCOSBootstrapDownloadURL, "stable", "58fd0833ce81b6244fc73bf65b5deb43217b0bd7"),
		DCOS198BootstrapDownloadURL:     fmt.Sprintf(AzureEdgeDCOSBootstrapDownloadURL, "stable/1.9.8", "f4ae0d20665fc68ee25282d6f78681b2773c6e10"),
		DCOS110BootstrapDownloadURL:     fmt.Sprintf(AzureEdgeDCOSBootstrapDownloadURL, "stable/1.10.0", "4d92536e7381176206e71ee15b5ffe454439920c"),
		DCOS111BootstrapDownloadURL:     fmt.Sprintf(AzureEdgeDCOSBootstrapDownloadURL, "stable/1.11.0", "a0654657903fb68dff60f6e522a7f241c1bfbf0f"),
		DCOSWindowsBootstrapDownloadURL: "http://dcos-win.westus.cloudapp.azure.com/dcos-windows/stable/",
		DcosRepositoryURL:               "https://dcosio.azureedge.net/dcos/stable/1.11.0",
		DcosClusterPackageListID:        "248a66388bba1adbcb14a52fd3b7b424ab06fa76",
	}

	//DefaultDockerSpecConfig is the default Docker engine repo.
	DefaultDockerSpecConfig = DockerSpecConfig{
		DockerEngineRepo:         "https://aptdocker.azureedge.net/repo",
		DockerComposeDownloadURL: "https://github.com/docker/compose/releases/download",
	}

	//DefaultUbuntuImageConfig is the default Linux distribution.
	DefaultUbuntuImageConfig = AzureOSImageConfig{
		ImageOffer:     "UbuntuServer",
		ImageSku:       "16.04-LTS",
		ImagePublisher: "Canonical",
		ImageVersion:   "16.04.201808140",
	}

	//DefaultRHELOSImageConfig is the RHEL Linux distribution.
	DefaultRHELOSImageConfig = AzureOSImageConfig{
		ImageOffer:     "RHEL",
		ImageSku:       "7.3",
		ImagePublisher: "RedHat",
		ImageVersion:   "latest",
	}

	//DefaultCoreOSImageConfig is the CoreOS Linux distribution.
	DefaultCoreOSImageConfig = AzureOSImageConfig{
		ImageOffer:     "CoreOS",
		ImageSku:       "Stable",
		ImagePublisher: "CoreOS",
		ImageVersion:   "latest",
	}

	// DefaultAKSOSImageConfig is the AKS image based on Ubuntu 16.04.
	DefaultAKSOSImageConfig = AzureOSImageConfig{
		ImageOffer:     "aks",
		ImageSku:       "aksbase",
		ImagePublisher: "microsoft-aks",
		ImageVersion:   "0.13.0",
	}

	//DefaultOpenShift39RHELImageConfig is the OpenShift on RHEL distribution.
	DefaultOpenShift39RHELImageConfig = AzureOSImageConfig{
		ImageOffer:     "acsengine-preview",
		ImageSku:       "rhel74",
		ImagePublisher: "redhat",
		ImageVersion:   "latest",
	}

	//DefaultOpenShift39CentOSImageConfig is the OpenShift on CentOS distribution.
	DefaultOpenShift39CentOSImageConfig = AzureOSImageConfig{
		ImageOffer:     "origin-acsengine-preview",
		ImageSku:       "centos7",
		ImagePublisher: "redhat",
		ImageVersion:   "latest",
	}

	//AzureCloudSpec is the default configurations for global azure.
	AzureCloudSpec = AzureEnvironmentSpecConfig{
		CloudName: azurePublicCloud,
		//DockerSpecConfig specify the docker engine download repo
		DockerSpecConfig: DefaultDockerSpecConfig,
		//KubernetesSpecConfig is the default kubernetes container image url.
		KubernetesSpecConfig: DefaultKubernetesSpecConfig,
		DCOSSpecConfig:       DefaultDCOSSpecConfig,

		EndpointConfig: AzureEndpointConfig{
			ResourceManagerVMDNSSuffix: "cloudapp.azure.com",
		},

		OSImageConfig: map[api.Distro]AzureOSImageConfig{
			api.Ubuntu: DefaultUbuntuImageConfig,
			api.RHEL:   DefaultRHELOSImageConfig,
			api.CoreOS: DefaultCoreOSImageConfig,
			api.AKS:    DefaultAKSOSImageConfig,
			// Image config supported for OpenShift
			api.OpenShift39RHEL: DefaultOpenShift39RHELImageConfig,
			api.OpenShiftCentOS: DefaultOpenShift39CentOSImageConfig,
		},
	}

	//AzureGermanCloudSpec is the German cloud config.
	AzureGermanCloudSpec = AzureEnvironmentSpecConfig{
		CloudName:            azureGermanCloud,
		DockerSpecConfig:     DefaultDockerSpecConfig,
		KubernetesSpecConfig: DefaultKubernetesSpecConfig,
		DCOSSpecConfig:       DefaultDCOSSpecConfig,
		EndpointConfig: AzureEndpointConfig{
			ResourceManagerVMDNSSuffix: "cloudapp.microsoftazure.de",
		},
		OSImageConfig: map[api.Distro]AzureOSImageConfig{
			api.Ubuntu: {
				ImageOffer:     "UbuntuServer",
				ImageSku:       "16.04-LTS",
				ImagePublisher: "Canonical",
				ImageVersion:   "16.04.201801050",
			},
			api.RHEL:   DefaultRHELOSImageConfig,
			api.CoreOS: DefaultCoreOSImageConfig,
		},
	}

	//AzureUSGovernmentCloud is the US government config.
	AzureUSGovernmentCloud = AzureEnvironmentSpecConfig{
		CloudName:            azureUSGovernmentCloud,
		DockerSpecConfig:     DefaultDockerSpecConfig,
		KubernetesSpecConfig: DefaultKubernetesSpecConfig,
		DCOSSpecConfig:       DefaultDCOSSpecConfig,
		EndpointConfig: AzureEndpointConfig{
			ResourceManagerVMDNSSuffix: "cloudapp.usgovcloudapi.net",
		},
		OSImageConfig: map[api.Distro]AzureOSImageConfig{
			api.Ubuntu: {
				ImageOffer:     "UbuntuServer",
				ImageSku:       "16.04-LTS",
				ImagePublisher: "Canonical",
				ImageVersion:   "latest",
			},
			api.RHEL:   DefaultRHELOSImageConfig,
			api.CoreOS: DefaultCoreOSImageConfig,
		},
	}

	//AzureChinaCloudSpec is the configurations for Azure China (Mooncake)
	AzureChinaCloudSpec = AzureEnvironmentSpecConfig{
		CloudName: azureChinaCloud,
		//DockerSpecConfig specify the docker engine download repo
		DockerSpecConfig: DockerSpecConfig{
			DockerEngineRepo:         "https://mirror.azure.cn/docker-engine/apt/repo/",
			DockerComposeDownloadURL: "https://mirror.azure.cn/docker-toolbox/linux/compose",
		},
		//KubernetesSpecConfig - Due to Chinese firewall issue, the default containers from google is blocked, use the Chinese local mirror instead
		KubernetesSpecConfig: KubernetesSpecConfig{
			KubernetesImageBase:              "gcr.akscn.io/google_containers/",
			TillerImageBase:                  "gcr.akscn.io/kubernetes-helm/",
			ACIConnectorImageBase:            DefaultKubernetesSpecConfig.ACIConnectorImageBase,
			EtcdDownloadURLBase:              DefaultKubernetesSpecConfig.EtcdDownloadURLBase,
			KubeBinariesSASURLBase:           DefaultKubernetesSpecConfig.KubeBinariesSASURLBase,
			WindowsPackageSASURLBase:         DefaultKubernetesSpecConfig.WindowsPackageSASURLBase,
			WindowsTelemetryGUID:             DefaultKubernetesSpecConfig.WindowsTelemetryGUID,
			CNIPluginsDownloadURL:            DefaultKubernetesSpecConfig.CNIPluginsDownloadURL,
			VnetCNILinuxPluginsDownloadURL:   DefaultKubernetesSpecConfig.VnetCNILinuxPluginsDownloadURL,
			VnetCNIWindowsPluginsDownloadURL: DefaultKubernetesSpecConfig.VnetCNIWindowsPluginsDownloadURL,
			ContainerdDownloadURLBase:        "https://mirror.azure.cn/kubernetes/containerd/",
		},
		DCOSSpecConfig: DCOSSpecConfig{
			DCOS188BootstrapDownloadURL:     fmt.Sprintf(AzureChinaCloudDCOSBootstrapDownloadURL, "5df43052907c021eeb5de145419a3da1898c58a5"),
			DCOSWindowsBootstrapDownloadURL: "https://dcosdevstorage.blob.core.windows.net/dcos-windows",
			DCOS190BootstrapDownloadURL:     fmt.Sprintf(AzureChinaCloudDCOSBootstrapDownloadURL, "58fd0833ce81b6244fc73bf65b5deb43217b0bd7"),
			DCOS198BootstrapDownloadURL:     fmt.Sprintf(AzureChinaCloudDCOSBootstrapDownloadURL, "f4ae0d20665fc68ee25282d6f78681b2773c6e10"),
		},

		EndpointConfig: AzureEndpointConfig{
			ResourceManagerVMDNSSuffix: "cloudapp.chinacloudapi.cn",
		},
		OSImageConfig: map[api.Distro]AzureOSImageConfig{
			api.Ubuntu: {
				ImageOffer:     "UbuntuServer",
				ImageSku:       "16.04-LTS",
				ImagePublisher: "Canonical",
				ImageVersion:   "latest",
			},
			api.RHEL:   DefaultRHELOSImageConfig,
			api.CoreOS: DefaultCoreOSImageConfig,
		},
	}
)

// setPropertiesDefaults for the container Properties, returns true if certs are generated
func setPropertiesDefaults(cs *api.ContainerService, isUpgrade, isScale bool) (bool, error) {
	properties := cs.Properties

	setOrchestratorDefaults(cs, isUpgrade || isScale)

	// Set master profile defaults if this cluster configuration includes master node(s)
	if cs.Properties.MasterProfile != nil {
		setMasterProfileDefaults(properties, isUpgrade)
	}

	setAgentProfileDefaults(properties, isUpgrade, isScale)

	setStorageDefaults(properties)
	setExtensionDefaults(properties)
	setVMSSDefaults(properties)

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
	location := cs.Location
	a := cs.Properties

	cloudSpecConfig := getCloudSpecConfig(location)
	if a.OrchestratorProfile == nil {
		return
	}
	o := a.OrchestratorProfile
	o.OrchestratorVersion = common.GetValidPatchVersion(
		o.OrchestratorType,
		o.OrchestratorVersion, isUpdate, a.HasWindows())

	switch o.OrchestratorType {
	case api.Kubernetes:
		k8sVersion := o.OrchestratorVersion

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
		// Enforce sane cloudprovider backoff defaults, if CloudProviderBackoff is true in KubernetesConfig
		o.KubernetesConfig.CloudProviderBackoff = true
		if o.KubernetesConfig.CloudProviderBackoffDuration == 0 {
			o.KubernetesConfig.CloudProviderBackoffDuration = DefaultKubernetesCloudProviderBackoffDuration
		}
		if o.KubernetesConfig.CloudProviderBackoffExponent == 0 {
			o.KubernetesConfig.CloudProviderBackoffExponent = DefaultKubernetesCloudProviderBackoffExponent
		}
		if o.KubernetesConfig.CloudProviderBackoffJitter == 0 {
			o.KubernetesConfig.CloudProviderBackoffJitter = DefaultKubernetesCloudProviderBackoffJitter
		}
		if o.KubernetesConfig.CloudProviderBackoffRetries == 0 {
			o.KubernetesConfig.CloudProviderBackoffRetries = DefaultKubernetesCloudProviderBackoffRetries
		}
		k8sSemVer, _ := semver.Make(k8sVersion)
		minVersion, _ := semver.Make("1.6.6")
		// Enforce sane cloudprovider rate limit defaults, if CloudProviderRateLimit is true in KubernetesConfig
		// For k8s version greater or equal to 1.6.6, we will set the default CloudProviderRate* settings
		o.KubernetesConfig.CloudProviderRateLimit = true
		if o.KubernetesConfig.CloudProviderRateLimit && k8sSemVer.GTE(minVersion) {
			if o.KubernetesConfig.CloudProviderRateLimitQPS == 0 {
				o.KubernetesConfig.CloudProviderRateLimitQPS = DefaultKubernetesCloudProviderRateLimitQPS
			}
			if o.KubernetesConfig.CloudProviderRateLimitBucket == 0 {
				o.KubernetesConfig.CloudProviderRateLimitBucket = DefaultKubernetesCloudProviderRateLimitBucket
			}
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

		if a.OrchestratorProfile.KubernetesConfig.LoadBalancerSku == "" {
			a.OrchestratorProfile.KubernetesConfig.LoadBalancerSku = api.DefaultLoadBalancerSku
		}

		if common.IsKubernetesVersionGe(a.OrchestratorProfile.OrchestratorVersion, "1.11.0") && a.OrchestratorProfile.KubernetesConfig.LoadBalancerSku == "Standard" {
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
	// don't default Distro for OpenShift
	if !a.OrchestratorProfile.IsOpenShift() {
		if a.MasterProfile.Distro == "" {
			a.MasterProfile.Distro = api.AKS
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

// setVMSSDefaults
func setVMSSDefaults(a *api.Properties) {
	for _, profile := range a.AgentPoolProfiles {
		if profile.AvailabilityProfile == api.VirtualMachineScaleSets {
			if profile.Count > 100 {
				profile.SinglePlacementGroup = helpers.PointerToBool(false)
			}
			if profile.SinglePlacementGroup == nil {
				profile.SinglePlacementGroup = helpers.PointerToBool(api.DefaultSinglePlacementGroup)
			}
			if profile.SinglePlacementGroup == helpers.PointerToBool(false) {
				profile.StorageProfile = api.ManagedDisks
			}
			if profile.HasAvailabilityZones() {
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

		// don't default Distro for OpenShift
		if !a.OrchestratorProfile.IsOpenShift() {
			if profile.Distro == "" {
				if profile.OSDiskSizeGB != 0 && profile.OSDiskSizeGB < api.VHDDiskSizeAKS {
					profile.Distro = api.Ubuntu
				} else {
					profile.Distro = api.AKS
				}
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

func setDefaultCerts(a *api.Properties) (bool, error) {
	if a.MasterProfile != nil && a.OrchestratorProfile.OrchestratorType == api.OpenShift {
		return certgen.OpenShiftSetDefaultCerts(a, DefaultOpenshiftOrchestratorName, GenerateClusterID(a))
	}

	if a.MasterProfile == nil || a.OrchestratorProfile.OrchestratorType != api.Kubernetes {
		return false, nil
	}

	provided := certsAlreadyPresent(a.CertificateProfile, a.MasterProfile.Count)

	if areAllTrue(provided) {
		return false, nil
	}

	masterExtraFQDNs := append(formatAzureProdFQDNs(a.MasterProfile.DNSPrefix), a.MasterProfile.SubjectAltNames...)
	firstMasterIP := net.ParseIP(a.MasterProfile.FirstConsecutiveStaticIP).To4()

	if firstMasterIP == nil {
		return false, errors.Errorf("MasterProfile.FirstConsecutiveStaticIP '%s' is an invalid IP address", a.MasterProfile.FirstConsecutiveStaticIP)
	}

	ips := []net.IP{firstMasterIP}
	// Add the Internal Loadbalancer IP which is always at at a known offset from the firstMasterIP
	ips = append(ips, net.IP{firstMasterIP[0], firstMasterIP[1], firstMasterIP[2], firstMasterIP[3] + byte(DefaultInternalLbStaticIPOffset)})
	// Include the Internal load balancer as well

	if a.MasterProfile.IsVirtualMachineScaleSets() {
		// Include the Internal load balancer as well
		for i := 1; i < a.MasterProfile.Count; i++ {
			offset := i * a.MasterProfile.IPAddressCount
			ip := net.IP{firstMasterIP[0], firstMasterIP[1], firstMasterIP[2], firstMasterIP[3] + byte(offset)}
			ips = append(ips, ip)
		}
	} else {
		for i := 1; i < a.MasterProfile.Count; i++ {
			ip := net.IP{firstMasterIP[0], firstMasterIP[1], firstMasterIP[2], firstMasterIP[3] + byte(i)}
			ips = append(ips, ip)
		}
	}
	if a.CertificateProfile == nil {
		a.CertificateProfile = &api.CertificateProfile{}
	}

	// use the specified Certificate Authority pair, or generate a new pair
	var caPair *PkiKeyCertPair
	if provided["ca"] {
		caPair = &PkiKeyCertPair{CertificatePem: a.CertificateProfile.CaCertificate, PrivateKeyPem: a.CertificateProfile.CaPrivateKey}
	} else {
		caCertificate, caPrivateKey, err := createCertificate("ca", nil, nil, false, false, nil, nil, nil)
		if err != nil {
			return false, err
		}
		caPair = &PkiKeyCertPair{CertificatePem: string(certificateToPem(caCertificate.Raw)), PrivateKeyPem: string(privateKeyToPem(caPrivateKey))}
		a.CertificateProfile.CaCertificate = caPair.CertificatePem
		a.CertificateProfile.CaPrivateKey = caPair.PrivateKeyPem
	}

	cidrFirstIP, err := common.CidrStringFirstIP(a.OrchestratorProfile.KubernetesConfig.ServiceCIDR)
	if err != nil {
		return false, err
	}
	ips = append(ips, cidrFirstIP)

	apiServerPair, clientPair, kubeConfigPair, etcdServerPair, etcdClientPair, etcdPeerPairs, err := CreatePki(masterExtraFQDNs, ips, DefaultKubernetesClusterDomain, caPair, a.MasterProfile.Count)
	if err != nil {
		return false, err
	}

	// If no Certificate Authority pair or no cert/key pair was provided, use generated cert/key pairs signed by provided Certificate Authority pair
	if !provided["apiserver"] || !provided["ca"] {
		a.CertificateProfile.APIServerCertificate = apiServerPair.CertificatePem
		a.CertificateProfile.APIServerPrivateKey = apiServerPair.PrivateKeyPem
	}
	if !provided["client"] || !provided["ca"] {
		a.CertificateProfile.ClientCertificate = clientPair.CertificatePem
		a.CertificateProfile.ClientPrivateKey = clientPair.PrivateKeyPem
	}
	if !provided["kubeconfig"] || !provided["ca"] {
		a.CertificateProfile.KubeConfigCertificate = kubeConfigPair.CertificatePem
		a.CertificateProfile.KubeConfigPrivateKey = kubeConfigPair.PrivateKeyPem
	}
	if !provided["etcd"] || !provided["ca"] {
		a.CertificateProfile.EtcdServerCertificate = etcdServerPair.CertificatePem
		a.CertificateProfile.EtcdServerPrivateKey = etcdServerPair.PrivateKeyPem
		a.CertificateProfile.EtcdClientCertificate = etcdClientPair.CertificatePem
		a.CertificateProfile.EtcdClientPrivateKey = etcdClientPair.PrivateKeyPem
		a.CertificateProfile.EtcdPeerCertificates = make([]string, a.MasterProfile.Count)
		a.CertificateProfile.EtcdPeerPrivateKeys = make([]string, a.MasterProfile.Count)
		for i, v := range etcdPeerPairs {
			a.CertificateProfile.EtcdPeerCertificates[i] = v.CertificatePem
			a.CertificateProfile.EtcdPeerPrivateKeys[i] = v.PrivateKeyPem
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
