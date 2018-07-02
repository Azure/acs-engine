package api

import (
	"strconv"

	"github.com/Azure/acs-engine/pkg/api/agentPoolOnlyApi/v20170831"
	"github.com/Azure/acs-engine/pkg/api/agentPoolOnlyApi/v20180331"
	"github.com/Azure/acs-engine/pkg/api/agentPoolOnlyApi/vlabs"
	"github.com/Azure/acs-engine/pkg/helpers"
)

///////////////////////////////////////////////////////////
// The converter exposes functions to convert the top level
// ContainerService resource
//
// All other functions are internal helper functions used
// for converting.
///////////////////////////////////////////////////////////

// ConvertContainerServiceToV20170831AgentPoolOnly converts an unversioned ContainerService to a v20170831 ContainerService
func ConvertContainerServiceToV20170831AgentPoolOnly(api *ContainerService) *v20170831.ManagedCluster {
	v20170831HCP := &v20170831.ManagedCluster{}
	v20170831HCP.ID = api.ID
	v20170831HCP.Location = api.Location
	v20170831HCP.Name = api.Name
	if api.Plan != nil {
		v20170831HCP.Plan = &v20170831.ResourcePurchasePlan{}
		convertResourcePurchasePlanToV20170831AgentPoolOnly(api.Plan, v20170831HCP.Plan)
	}
	v20170831HCP.Tags = map[string]string{}
	for k, v := range api.Tags {
		v20170831HCP.Tags[k] = v
	}
	v20170831HCP.Type = api.Type
	v20170831HCP.Properties = &v20170831.Properties{}
	convertPropertiesToV20170831AgentPoolOnly(api.Properties, v20170831HCP.Properties)
	return v20170831HCP
}

// ConvertContainerServiceToV20180331AgentPoolOnly converts an unversioned ContainerService to a v20180331 ContainerService
func ConvertContainerServiceToV20180331AgentPoolOnly(api *ContainerService) *v20180331.ManagedCluster {
	v20180331HCP := &v20180331.ManagedCluster{}
	v20180331HCP.ID = api.ID
	v20180331HCP.Location = api.Location
	v20180331HCP.Name = api.Name
	if api.Plan != nil {
		v20180331HCP.Plan = &v20180331.ResourcePurchasePlan{}
		convertResourcePurchasePlanToV20180331AgentPoolOnly(api.Plan, v20180331HCP.Plan)
	}
	v20180331HCP.Tags = map[string]string{}
	for k, v := range api.Tags {
		v20180331HCP.Tags[k] = v
	}
	v20180331HCP.Type = api.Type
	v20180331HCP.Properties = &v20180331.Properties{}
	convertPropertiesToV20180331AgentPoolOnly(api.Properties, v20180331HCP.Properties)
	return v20180331HCP
}

// ConvertContainerServiceToVLabsAgentPoolOnly converts an unversioned ContainerService to a vlabs ContainerService
func ConvertContainerServiceToVLabsAgentPoolOnly(api *ContainerService) *vlabs.ManagedCluster {
	vlabsHCP := &vlabs.ManagedCluster{}
	vlabsHCP.ID = api.ID
	vlabsHCP.Location = api.Location
	vlabsHCP.Name = api.Name
	if api.Plan != nil {
		vlabsHCP.Plan = &vlabs.ResourcePurchasePlan{}
		convertResourcePurchasePlanToVLabsAgentPoolOnly(api.Plan, vlabsHCP.Plan)
	}
	vlabsHCP.Tags = map[string]string{}
	for k, v := range api.Tags {
		vlabsHCP.Tags[k] = v
	}
	vlabsHCP.Type = api.Type
	vlabsHCP.Properties = &vlabs.Properties{}
	convertPropertiesToVLabsAgentPoolOnly(api.Properties, vlabsHCP.Properties)
	return vlabsHCP
}

// convertResourcePurchasePlanToVLabsAgentPoolOnly converts a vlabs ResourcePurchasePlan to an unversioned ResourcePurchasePlan
func convertResourcePurchasePlanToVLabsAgentPoolOnly(api *ResourcePurchasePlan, vlabs *vlabs.ResourcePurchasePlan) {
	vlabs.Name = api.Name
	vlabs.Product = api.Product
	vlabs.PromotionCode = api.PromotionCode
	vlabs.Publisher = api.Publisher
}

// convertResourcePurchasePlanToV20170831 converts a v20170831 ResourcePurchasePlan to an unversioned ResourcePurchasePlan
func convertResourcePurchasePlanToV20170831AgentPoolOnly(api *ResourcePurchasePlan, v20170831 *v20170831.ResourcePurchasePlan) {
	v20170831.Name = api.Name
	v20170831.Product = api.Product
	v20170831.PromotionCode = api.PromotionCode
	v20170831.Publisher = api.Publisher
}

func convertPropertiesToV20170831AgentPoolOnly(api *Properties, p *v20170831.Properties) {
	p.ProvisioningState = v20170831.ProvisioningState(api.ProvisioningState)
	if api.OrchestratorProfile != nil {
		if api.OrchestratorProfile.OrchestratorVersion != "" {
			p.KubernetesVersion = api.OrchestratorProfile.OrchestratorVersion
		}
	}
	if api.HostedMasterProfile != nil {
		p.DNSPrefix = api.HostedMasterProfile.DNSPrefix
		p.FQDN = api.HostedMasterProfile.FQDN
	}
	p.AgentPoolProfiles = []*v20170831.AgentPoolProfile{}
	for _, apiProfile := range api.AgentPoolProfiles {
		v20170831Profile := &v20170831.AgentPoolProfile{}
		convertAgentPoolProfileToV20170831AgentPoolOnly(apiProfile, v20170831Profile)
		p.AgentPoolProfiles = append(p.AgentPoolProfiles, v20170831Profile)
	}
	if api.LinuxProfile != nil {
		p.LinuxProfile = &v20170831.LinuxProfile{}
		convertLinuxProfileToV20170831AgentPoolOnly(api.LinuxProfile, p.LinuxProfile)
	}
	if api.WindowsProfile != nil {
		p.WindowsProfile = &v20170831.WindowsProfile{}
		convertWindowsProfileToV20170831AgentPoolOnly(api.WindowsProfile, p.WindowsProfile)
	}
	if api.ServicePrincipalProfile != nil {
		p.ServicePrincipalProfile = &v20170831.ServicePrincipalProfile{}
		convertServicePrincipalProfileToV20170831AgentPoolOnly(api.ServicePrincipalProfile, p.ServicePrincipalProfile)
	}
}

func convertLinuxProfileToV20170831AgentPoolOnly(api *LinuxProfile, obj *v20170831.LinuxProfile) {
	obj.AdminUsername = api.AdminUsername
	obj.SSH.PublicKeys = []v20170831.PublicKey{}
	for _, d := range api.SSH.PublicKeys {
		obj.SSH.PublicKeys = append(obj.SSH.PublicKeys, v20170831.PublicKey{
			KeyData: d.KeyData,
		})
	}
}

func convertWindowsProfileToV20170831AgentPoolOnly(api *WindowsProfile, v20170831Profile *v20170831.WindowsProfile) {
	v20170831Profile.AdminUsername = api.AdminUsername
	v20170831Profile.AdminPassword = api.AdminPassword
}

func convertAgentPoolProfileToV20170831AgentPoolOnly(api *AgentPoolProfile, p *v20170831.AgentPoolProfile) {
	p.Name = api.Name
	p.Count = api.Count
	p.VMSize = api.VMSize
	p.OSType = v20170831.OSType(api.OSType)
	p.SetSubnet(api.Subnet)
	p.OSDiskSizeGB = api.OSDiskSizeGB
	p.StorageProfile = api.StorageProfile
	p.VnetSubnetID = api.VnetSubnetID
}

func convertServicePrincipalProfileToV20170831AgentPoolOnly(api *ServicePrincipalProfile, v20170831 *v20170831.ServicePrincipalProfile) {
	v20170831.ClientID = api.ClientID
	v20170831.Secret = api.Secret
	// v20170831.KeyvaultSecretRef = api.KeyvaultSecretRef
}

// convertResourcePurchasePlanToV20180331 converts a v20180331 ResourcePurchasePlan to an unversioned ResourcePurchasePlan
func convertResourcePurchasePlanToV20180331AgentPoolOnly(api *ResourcePurchasePlan, v20180331 *v20180331.ResourcePurchasePlan) {
	v20180331.Name = api.Name
	v20180331.Product = api.Product
	v20180331.PromotionCode = api.PromotionCode
	v20180331.Publisher = api.Publisher
}

func convertKubernetesConfigToEnableRBACV20180331AgentPoolOnly(kc *KubernetesConfig) *bool {
	if kc == nil {
		return helpers.PointerToBool(false)
	}
	// We use KubernetesConfig.EnableRbac to convert to versioned api model
	// The assumption here is KubernetesConfig.EnableSecureKubelet is set to be same
	if kc != nil && kc.EnableRbac != nil && *kc.EnableRbac {
		return helpers.PointerToBool(true)
	}
	return helpers.PointerToBool(false)
}

func convertPropertiesToV20180331AgentPoolOnly(api *Properties, p *v20180331.Properties) {
	p.ProvisioningState = v20180331.ProvisioningState(api.ProvisioningState)

	if api.OrchestratorProfile != nil {
		p.EnableRBAC = convertKubernetesConfigToEnableRBACV20180331AgentPoolOnly(api.OrchestratorProfile.KubernetesConfig)
		p.KubernetesVersion, p.NetworkProfile = convertOrchestratorProfileToV20180331AgentPoolOnly(api.OrchestratorProfile)
	}
	if api.HostedMasterProfile != nil {
		p.DNSPrefix = api.HostedMasterProfile.DNSPrefix
		p.FQDN = api.HostedMasterProfile.FQDN
	}
	p.AgentPoolProfiles = []*v20180331.AgentPoolProfile{}
	for _, apiProfile := range api.AgentPoolProfiles {
		v20180331Profile := &v20180331.AgentPoolProfile{}
		convertAgentPoolProfileToV20180331AgentPoolOnly(apiProfile, v20180331Profile)
		p.AgentPoolProfiles = append(p.AgentPoolProfiles, v20180331Profile)
	}
	if api.LinuxProfile != nil {
		p.LinuxProfile = &v20180331.LinuxProfile{}
		convertLinuxProfileToV20180331AgentPoolOnly(api.LinuxProfile, p.LinuxProfile)
	}
	if api.WindowsProfile != nil {
		p.WindowsProfile = &v20180331.WindowsProfile{}
		convertWindowsProfileToV20180331AgentPoolOnly(api.WindowsProfile, p.WindowsProfile)
	}
	if api.ServicePrincipalProfile != nil {
		p.ServicePrincipalProfile = &v20180331.ServicePrincipalProfile{}
		convertServicePrincipalProfileToV20180331AgentPoolOnly(api.ServicePrincipalProfile, p.ServicePrincipalProfile)
	}
	if api.AddonProfiles != nil {
		p.AddonProfiles = make(map[string]v20180331.AddonProfile)
		convertAddonsProfileToV20180331AgentPoolOnly(api.AddonProfiles, p.AddonProfiles)
	}
	if api.AADProfile != nil {
		p.AADProfile = &v20180331.AADProfile{}
		convertAADProfileToV20180331AgentPoolOnly(api.AADProfile, p.AADProfile)
	}
}

func convertOrchestratorProfileToV20180331AgentPoolOnly(orchestratorProfile *OrchestratorProfile) (kubernetesVersion string, networkProfile *v20180331.NetworkProfile) {
	if orchestratorProfile.OrchestratorVersion != "" {
		kubernetesVersion = orchestratorProfile.OrchestratorVersion
	}

	if orchestratorProfile.KubernetesConfig != nil {
		k := orchestratorProfile.KubernetesConfig
		if k.NetworkPlugin != "" {
			networkProfile = &v20180331.NetworkProfile{}
			networkProfile.NetworkPlugin = v20180331.NetworkPlugin(k.NetworkPlugin)
			networkProfile.NetworkPolicy = v20180331.NetworkPolicy(k.NetworkPolicy)
			if k.NetworkPlugin == string(v20180331.Kubenet) {
				networkProfile.PodCidr = k.ClusterSubnet
			}
			networkProfile.ServiceCidr = k.ServiceCIDR
			networkProfile.DNSServiceIP = k.DNSServiceIP
			networkProfile.DockerBridgeCidr = k.DockerBridgeSubnet
		} else if k.NetworkPolicy != "" {
			networkProfile = &v20180331.NetworkProfile{}
			// ACS-E uses "none" in the old un-versioned model to represent kubenet.
			if k.NetworkPolicy == "none" {
				networkProfile.NetworkPlugin = v20180331.Kubenet
				networkProfile.PodCidr = k.ClusterSubnet
			} else {
				networkProfile.NetworkPlugin = v20180331.NetworkPlugin(k.NetworkPolicy)
			}
			networkProfile.ServiceCidr = k.ServiceCIDR
			networkProfile.DNSServiceIP = k.DNSServiceIP
			networkProfile.DockerBridgeCidr = k.DockerBridgeSubnet
		}
	}

	return kubernetesVersion, networkProfile
}

func convertLinuxProfileToV20180331AgentPoolOnly(api *LinuxProfile, obj *v20180331.LinuxProfile) {
	obj.AdminUsername = api.AdminUsername
	obj.SSH.PublicKeys = []v20180331.PublicKey{}
	for _, d := range api.SSH.PublicKeys {
		obj.SSH.PublicKeys = append(obj.SSH.PublicKeys, v20180331.PublicKey{
			KeyData: d.KeyData,
		})
	}
}

func convertWindowsProfileToV20180331AgentPoolOnly(api *WindowsProfile, v20180331Profile *v20180331.WindowsProfile) {
	v20180331Profile.AdminUsername = api.AdminUsername
	v20180331Profile.AdminPassword = api.AdminPassword
}

func convertAgentPoolProfileToV20180331AgentPoolOnly(api *AgentPoolProfile, p *v20180331.AgentPoolProfile) {
	p.Name = api.Name
	p.Count = api.Count
	p.VMSize = api.VMSize
	p.OSType = v20180331.OSType(api.OSType)
	p.SetSubnet(api.Subnet)
	p.OSDiskSizeGB = api.OSDiskSizeGB
	p.StorageProfile = api.StorageProfile
	p.VnetSubnetID = api.VnetSubnetID
	if api.KubernetesConfig != nil && api.KubernetesConfig.KubeletConfig != nil {
		if maxPods, ok := api.KubernetesConfig.KubeletConfig["--max-pods"]; ok {
			agentPoolMaxPods, _ := strconv.Atoi(maxPods)
			p.MaxPods = &agentPoolMaxPods
		}
	}
}

func convertServicePrincipalProfileToV20180331AgentPoolOnly(api *ServicePrincipalProfile, v20180331 *v20180331.ServicePrincipalProfile) {
	v20180331.ClientID = api.ClientID
	v20180331.Secret = api.Secret
	// v20180331.KeyvaultSecretRef = api.KeyvaultSecretRef
}

func convertAddonsProfileToV20180331AgentPoolOnly(api map[string]AddonProfile, p map[string]v20180331.AddonProfile) {
	if api == nil {
		return
	}

	for k, v := range api {
		p[k] = v20180331.AddonProfile{
			Enabled: v.Enabled,
			Config:  v.Config,
		}
	}
}

func convertAADProfileToV20180331AgentPoolOnly(api *AADProfile, v20180331 *v20180331.AADProfile) {
	v20180331.ClientAppID = api.ClientAppID
	v20180331.ServerAppID = api.ServerAppID
	v20180331.TenantID = api.TenantID
	if api.Authenticator == Webhook {
		v20180331.ServerAppSecret = api.ServerAppSecret
	}
}

func convertPropertiesToVLabsAgentPoolOnly(api *Properties, p *vlabs.Properties) {
	p.ProvisioningState = vlabs.ProvisioningState(api.ProvisioningState)

	if api.OrchestratorProfile != nil {
		p.OrchestratorProfile = &vlabs.OrchestratorProfile{}
		convertOrchestratorProfileToVLabsAgentPoolOnly(api.OrchestratorProfile, p.OrchestratorProfile)
	}
	if api.HostedMasterProfile != nil {
		p.DNSPrefix = api.HostedMasterProfile.DNSPrefix
		p.FQDN = api.HostedMasterProfile.FQDN
	}
	p.AgentPoolProfiles = []*vlabs.AgentPoolProfile{}
	for _, apiProfile := range api.AgentPoolProfiles {
		vlabsProfile := &vlabs.AgentPoolProfile{}
		convertAgentPoolProfileToVLabsAgentPoolOnly(apiProfile, vlabsProfile)
		p.AgentPoolProfiles = append(p.AgentPoolProfiles, vlabsProfile)
	}
	if api.LinuxProfile != nil {
		p.LinuxProfile = &vlabs.LinuxProfile{}
		convertLinuxProfileToVLabsAgentPoolOnly(api.LinuxProfile, p.LinuxProfile)
	}
	if api.ServicePrincipalProfile != nil {
		p.ServicePrincipalProfile = &vlabs.ServicePrincipalProfile{}
		convertServicePrincipalProfileToVLabsAgentPoolOnly(api.ServicePrincipalProfile, p.ServicePrincipalProfile)
	}
	if api.CertificateProfile != nil {
		p.CertificateProfile = &vlabs.CertificateProfile{}
		convertCertificateProfileToVLabsAgentPoolOnly(api.CertificateProfile, p.CertificateProfile)
	}
	if api.AzProfile != nil {
		p.AzProfile = &vlabs.AzProfile{}
		convertAzProfileToVLabsAgentPoolOnly(api.AzProfile, p.AzProfile)
	}
}

func convertOrchestratorProfileToVLabsAgentPoolOnly(api *OrchestratorProfile, p *vlabs.OrchestratorProfile) {
	p.OrchestratorType = api.OrchestratorType
	p.OrchestratorVersion = api.OrchestratorVersion
	if api.OpenShiftConfig != nil {
		p.OpenShiftConfig = &vlabs.OpenShiftConfig{}
		convertOpenShiftConfigToVlabsAgentPoolOnly(api.OpenShiftConfig, p.OpenShiftConfig)
	}
}

func convertAgentPoolProfileToVLabsAgentPoolOnly(api *AgentPoolProfile, p *vlabs.AgentPoolProfile) {
	p.Name = api.Name
	p.Count = api.Count
	p.VMSize = api.VMSize
	p.OSType = vlabs.OSType(api.OSType)
	p.SetSubnet(api.Subnet)
	p.OSDiskSizeGB = api.OSDiskSizeGB
	p.StorageProfile = api.StorageProfile
	p.VnetSubnetID = api.VnetSubnetID
	p.AvailabilityProfile = api.AvailabilityProfile
}

func convertLinuxProfileToVLabsAgentPoolOnly(api *LinuxProfile, p *vlabs.LinuxProfile) {
	p.AdminUsername = api.AdminUsername
	p.SSH.PublicKeys = []vlabs.PublicKey{}
	for _, d := range api.SSH.PublicKeys {
		p.SSH.PublicKeys = append(p.SSH.PublicKeys, vlabs.PublicKey{
			KeyData: d.KeyData,
		})
	}
}

func convertServicePrincipalProfileToVLabsAgentPoolOnly(api *ServicePrincipalProfile, p *vlabs.ServicePrincipalProfile) {
	p.ClientID = api.ClientID
	p.Secret = api.Secret
}

func convertOpenShiftConfigToVlabsAgentPoolOnly(api *OpenShiftConfig, p *vlabs.OpenShiftConfig) {
	p.ConfigBundles = api.ConfigBundles
}

func convertCertificateProfileToVLabsAgentPoolOnly(api *CertificateProfile, p *vlabs.CertificateProfile) {
	p.CaCertificate = api.CaCertificate
	p.CaPrivateKey = api.CaPrivateKey
	p.APIServerCertificate = api.APIServerCertificate
	p.APIServerPrivateKey = api.APIServerPrivateKey
	p.ClientCertificate = api.ClientCertificate
	p.ClientPrivateKey = api.ClientPrivateKey
	p.KubeConfigCertificate = api.KubeConfigCertificate
	p.KubeConfigPrivateKey = api.KubeConfigPrivateKey
}

func convertAzProfileToVLabsAgentPoolOnly(api *AzProfile, p *vlabs.AzProfile) {
	p.TenantID = api.TenantID
	p.SubscriptionID = api.SubscriptionID
	p.ResourceGroup = api.ResourceGroup
	p.Location = api.Location
}
