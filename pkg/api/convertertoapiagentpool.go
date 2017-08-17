package api

import (
	"github.com/Azure/acs-engine/pkg/api/agentPoolOnlyApi/v20170831"
	"github.com/Azure/acs-engine/pkg/api/agentPoolOnlyApi/vlabs"
)

///////////////////////////////////////////////////////////
// The converter exposes functions to convert the top level
// ContainerService resource
//
// All other functions are internal helper functions used
// for converting.
///////////////////////////////////////////////////////////

// ConvertV20170831AgentPool converts an AgentPool object into an in-memory container service
func ConvertV20170831AgentPool(v20170831 *v20170831.HostedMaster) *ContainerService {
	c := &ContainerService{}
	c.ID = v20170831.ID
	c.Location = v20170831.Location
	c.Name = v20170831.Name
	if v20170831.Plan != nil {
		c.Plan = convertv20170831AgentPoolResourcePurchasePlan(v20170831.Plan)
	}
	c.Tags = map[string]string{}
	for k, v := range v20170831.Tags {
		c.Tags[k] = v
	}
	c.Type = v20170831.Type
	c.Properties = convertV20170831AgentPoolProperties(v20170831.Properties)
	return c
}

func convertv20170831AgentPoolResourcePurchasePlan(v20170831 *v20170831.ResourcePurchasePlan) *ResourcePurchasePlan {
	return &ResourcePurchasePlan{
		Name:          v20170831.Name,
		Product:       v20170831.Product,
		PromotionCode: v20170831.PromotionCode,
		Publisher:     v20170831.Publisher,
	}
}

func convertV20170831AgentPoolProperties(obj *v20170831.Properties) *Properties {
	properties := &Properties{
		ProvisioningState: ProvisioningState(obj.ProvisioningState),
		MasterProfile:     nil,
		FQDN:              obj.FQDN,
	}

	properties.OrchestratorProfile = convertV20170831AgentPoolOrchestratorProfile(obj)
	properties.MasterProfile = nil
	properties.AgentPoolProfiles = make([]*AgentPoolProfile, len(obj.AgentPoolProfiles))
	for i := range obj.AgentPoolProfiles {
		properties.AgentPoolProfiles[i] = convertV20170831AgentPoolProfile(obj.AgentPoolProfiles[i])
	}
	if obj.LinuxProfile != nil {
		properties.LinuxProfile = convertV20170831AgentPoolLinuxProfile(obj.LinuxProfile)
	}
	if obj.WindowsProfile != nil {
		properties.WindowsProfile = convertV20170831AgentPoolWindowsProfile(obj.WindowsProfile)
	}
	//	if obj.JumpboxProfile != nil {
	//		properties.JumpboxProfile = convertV20170831JumpboxProfile(obj.JumpboxProfile)
	//	}
	if obj.ServicePrincipalProfile != nil {
		properties.ServicePrincipalProfile = convertV20170831AgentPoolServicePrincipalProfile(obj.ServicePrincipalProfile)
	}
	//if obj.NetworkProfile != nil {
	//	properties.NetworkProfile = convertV20170831NetworkProfile(obj.NetworkProfile)
	//}
	//if obj.AccessProfile != nil {
	//	properties.AccessProfile = convertV20170831AccessProfile(obj.AccessProfile)
	//}
	return properties
}

// ConvertVLabsContainerService converts a vlabs ContainerService to an unversioned ContainerService
func ConvertVLabsAgentPool(vlabs *vlabs.HostedMaster) *ContainerService {
	c := &ContainerService{}
	c.ID = vlabs.ID
	c.Location = vlabs.Location
	c.Name = vlabs.Name
	if vlabs.Plan != nil {
		c.Plan = &ResourcePurchasePlan{}
		convertVLabsAgentPoolResourcePurchasePlan(vlabs.Plan, c.Plan)
	}
	c.Tags = map[string]string{}
	for k, v := range vlabs.Tags {
		c.Tags[k] = v
	}
	c.Type = vlabs.Type
	c.Properties = &Properties{}
	convertVLabsAgentPoolProperties(vlabs.Properties, c.Properties)
	return c
}

// convertVLabsResourcePurchasePlan converts a vlabs ResourcePurchasePlan to an unversioned ResourcePurchasePlan
func convertVLabsAgentPoolResourcePurchasePlan(vlabs *vlabs.ResourcePurchasePlan, api *ResourcePurchasePlan) {
	api.Name = vlabs.Name
	api.Product = vlabs.Product
	api.PromotionCode = vlabs.PromotionCode
	api.Publisher = vlabs.Publisher
}

func convertVLabsAgentPoolProperties(vlabs *vlabs.Properties, api *Properties) {
	api.ProvisioningState = ProvisioningState(vlabs.ProvisioningState)
	api.OrchestratorProfile = convertVLabsAgentPoolOrchestratorProfile(vlabs)
	api.MasterProfile = nil

	api.AgentPoolProfiles = []*AgentPoolProfile{}
	for _, p := range vlabs.AgentPoolProfiles {
		apiProfile := &AgentPoolProfile{}
		convertVLabsAgentPoolAgentPoolProfile(p, apiProfile)
		// by default vlabs will use managed disks for all orchestrators but kubernetes as it has encryption at rest.
		if !api.OrchestratorProfile.IsKubernetes() {
			// by default vlabs will use managed disks for all orchestrators but kubernetes as it has encryption at rest.
			if len(p.StorageProfile) == 0 {
				apiProfile.StorageProfile = ManagedDisks
			}
		}
		api.AgentPoolProfiles = append(api.AgentPoolProfiles, apiProfile)
	}
	if vlabs.LinuxProfile != nil {
		api.LinuxProfile = &LinuxProfile{}
		convertVLabsAgentPoolLinuxProfile(vlabs.LinuxProfile, api.LinuxProfile)
	}
	if vlabs.WindowsProfile != nil {
		api.WindowsProfile = &WindowsProfile{}
		convertVLabsAgentPoolWindowsProfile(vlabs.WindowsProfile, api.WindowsProfile)
	}
	if vlabs.ServicePrincipalProfile != nil {
		api.ServicePrincipalProfile = &ServicePrincipalProfile{}
		convertVLabsAgentPoolServicePrincipalProfile(vlabs.ServicePrincipalProfile, api.ServicePrincipalProfile)
	}
	if vlabs.CertificateProfile != nil {
		api.CertificateProfile = &CertificateProfile{}
		convertVLabsAgentPoolCertificateProfile(vlabs.CertificateProfile, api.CertificateProfile)
	}
}

func convertVLabsAgentPoolLinuxProfile(vlabs *vlabs.LinuxProfile, api *LinuxProfile) {
	api.AdminUsername = vlabs.AdminUsername
	api.SSH.PublicKeys = []PublicKey{}
	for _, d := range vlabs.SSH.PublicKeys {
		api.SSH.PublicKeys = append(api.SSH.PublicKeys,
			PublicKey{KeyData: d.KeyData})
	}
	// api.Secrets = []KeyVaultSecrets{}
	// for _, s := range vlabs.Secrets {
	// 	secret := &KeyVaultSecrets{}
	// 	convertVLabsKeyVaultSecrets(&s, secret)
	// 	api.Secrets = append(api.Secrets, *secret)
	// }
}

func convertV20170831AgentPoolLinuxProfile(obj *v20170831.LinuxProfile) *LinuxProfile {
	api := &LinuxProfile{
		AdminUsername: obj.AdminUsername,
	}
	api.SSH.PublicKeys = []PublicKey{}
	for _, d := range obj.SSH.PublicKeys {
		api.SSH.PublicKeys = append(api.SSH.PublicKeys, PublicKey{KeyData: d.KeyData})
	}
	return api
}

func convertV20170831AgentPoolWindowsProfile(obj *v20170831.WindowsProfile) *WindowsProfile {
	return &WindowsProfile{
		AdminUsername: obj.AdminUsername,
		AdminPassword: obj.AdminPassword,
	}
}

func convertVLabsAgentPoolWindowsProfile(vlabs *vlabs.WindowsProfile, api *WindowsProfile) {
	api.AdminUsername = vlabs.AdminUsername
	api.AdminPassword = vlabs.AdminPassword
	// api.Secrets = []KeyVaultSecrets{}
	// for _, s := range vlabs.Secrets {
	// 	secret := &KeyVaultSecrets{}
	// 	convertVLabsKeyVaultSecrets(&s, secret)
	// 	api.Secrets = append(api.Secrets, *secret)
	// }
}

func convertV20170831AgentPoolOrchestratorProfile(obj *v20170831.Properties) *OrchestratorProfile {
	orchestratorProfile := &OrchestratorProfile{
		OrchestratorType:    Kubernetes,
		OrchestratorVersion: KubernetesRelease1Dot7,
		OrchestratorRelease: KubernetesReleaseToVersion[KubernetesRelease1Dot7],
	}

	return orchestratorProfile
}

func convertVLabsAgentPoolOrchestratorProfile(obj *vlabs.Properties) *OrchestratorProfile {
	orchestratorProfile := &OrchestratorProfile{
		OrchestratorType:    Kubernetes,
		OrchestratorVersion: KubernetesRelease1Dot7,
		OrchestratorRelease: KubernetesReleaseToVersion[KubernetesRelease1Dot7],
	}

	return orchestratorProfile
}

func convertV20170831AgentPoolProfile(v20170831 *v20170831.AgentPoolProfile) *AgentPoolProfile {
	api := &AgentPoolProfile{}
	api.Name = v20170831.Name
	api.Count = v20170831.Count
	api.VMSize = v20170831.VMSize
	api.OSDiskSizeGB = v20170831.OSDiskSizeGB
	api.OSType = OSType(v20170831.OSType)
	api.StorageProfile = v20170831.StorageProfile
	api.VnetSubnetID = v20170831.VnetSubnetID
	api.Subnet = v20170831.GetSubnet()
	return api
}

func convertVLabsAgentPoolAgentPoolProfile(vlabs *vlabs.AgentPoolProfile, api *AgentPoolProfile) {
	api.Name = vlabs.Name
	api.Count = vlabs.Count
	api.VMSize = vlabs.VMSize
	api.OSDiskSizeGB = vlabs.OSDiskSizeGB
	api.OSType = OSType(vlabs.OSType)
	api.StorageProfile = vlabs.StorageProfile
	api.VnetSubnetID = vlabs.VnetSubnetID
	api.Subnet = vlabs.GetSubnet()
}

func convertVLabsAgentPoolServicePrincipalProfile(vlabs *vlabs.ServicePrincipalProfile, api *ServicePrincipalProfile) {
	api.ClientID = vlabs.ClientID
	api.Secret = vlabs.Secret
	// api.KeyvaultSecretRef = vlabs.KeyvaultSecretRef
}

func convertV20170831AgentPoolServicePrincipalProfile(obj *v20170831.ServicePrincipalProfile) *ServicePrincipalProfile {
	return &ServicePrincipalProfile{
		ClientID: obj.ClientID,
		Secret:   obj.Secret,
	}
}

func convertVLabsAgentPoolCertificateProfile(vlabs *vlabs.CertificateProfile, api *CertificateProfile) {
	api.CaCertificate = vlabs.CaCertificate
	api.CaPrivateKey = vlabs.CaPrivateKey
	api.APIServerCertificate = vlabs.APIServerCertificate
	api.APIServerPrivateKey = vlabs.APIServerPrivateKey
	api.ClientCertificate = vlabs.ClientCertificate
	api.ClientPrivateKey = vlabs.ClientPrivateKey
	api.KubeConfigCertificate = vlabs.KubeConfigCertificate
	api.KubeConfigPrivateKey = vlabs.KubeConfigPrivateKey
}

// func convertV20170831AgentPoolKubernetesConfig(vlabs *vlabs.KubernetesConfig, api *KubernetesConfig) {
// 	api.KubernetesImageBase = vlabs.KubernetesImageBase
// 	api.ClusterSubnet = vlabs.ClusterSubnet
// 	api.NetworkPolicy = vlabs.NetworkPolicy
// 	api.DockerBridgeSubnet = vlabs.DockerBridgeSubnet
// 	api.NodeStatusUpdateFrequency = vlabs.NodeStatusUpdateFrequency
// 	api.CtrlMgrNodeMonitorGracePeriod = vlabs.CtrlMgrNodeMonitorGracePeriod
// 	api.CtrlMgrPodEvictionTimeout = vlabs.CtrlMgrPodEvictionTimeout
// 	api.CtrlMgrRouteReconciliationPeriod = vlabs.CtrlMgrRouteReconciliationPeriod
// 	api.CloudProviderBackoff = vlabs.CloudProviderBackoff
// 	api.CloudProviderBackoffDuration = vlabs.CloudProviderBackoffDuration
// 	api.CloudProviderBackoffExponent = vlabs.CloudProviderBackoffExponent
// 	api.CloudProviderBackoffJitter = vlabs.CloudProviderBackoffJitter
// 	api.CloudProviderBackoffRetries = vlabs.CloudProviderBackoffRetries
// 	api.CloudProviderRateLimit = vlabs.CloudProviderRateLimit
// 	api.CloudProviderRateLimitBucket = vlabs.CloudProviderRateLimitBucket
// 	api.CloudProviderRateLimitQPS = vlabs.CloudProviderRateLimitQPS
// 	api.UseManagedIdentity = vlabs.UseManagedIdentity
// 	api.CustomHyperkubeImage = vlabs.CustomHyperkubeImage
// 	api.UseInstanceMetadata = vlabs.UseInstanceMetadata
// 	api.EnableRbac = vlabs.EnableRbac
// }

// func convertVLabsAgentPoolKubernetesConfig(vlabs *vlabs.KubernetesConfig, api *KubernetesConfig) {
// 	api.KubernetesImageBase = vlabs.KubernetesImageBase
// 	api.ClusterSubnet = vlabs.ClusterSubnet
// 	api.NetworkPolicy = vlabs.NetworkPolicy
// 	api.DockerBridgeSubnet = vlabs.DockerBridgeSubnet
// 	api.NodeStatusUpdateFrequency = vlabs.NodeStatusUpdateFrequency
// 	api.CtrlMgrNodeMonitorGracePeriod = vlabs.CtrlMgrNodeMonitorGracePeriod
// 	api.CtrlMgrPodEvictionTimeout = vlabs.CtrlMgrPodEvictionTimeout
// 	api.CtrlMgrRouteReconciliationPeriod = vlabs.CtrlMgrRouteReconciliationPeriod
// 	api.CloudProviderBackoff = vlabs.CloudProviderBackoff
// 	api.CloudProviderBackoffDuration = vlabs.CloudProviderBackoffDuration
// 	api.CloudProviderBackoffExponent = vlabs.CloudProviderBackoffExponent
// 	api.CloudProviderBackoffJitter = vlabs.CloudProviderBackoffJitter
// 	api.CloudProviderBackoffRetries = vlabs.CloudProviderBackoffRetries
// 	api.CloudProviderRateLimit = vlabs.CloudProviderRateLimit
// 	api.CloudProviderRateLimitBucket = vlabs.CloudProviderRateLimitBucket
// 	api.CloudProviderRateLimitQPS = vlabs.CloudProviderRateLimitQPS
// 	api.UseManagedIdentity = vlabs.UseManagedIdentity
// 	api.CustomHyperkubeImage = vlabs.CustomHyperkubeImage
// 	api.UseInstanceMetadata = vlabs.UseInstanceMetadata
// 	api.EnableRbac = vlabs.EnableRbac
// }
