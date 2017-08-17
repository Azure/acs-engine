package api

import "github.com/Azure/acs-engine/pkg/api/agentPoolOnlyApi/v20170831"

///////////////////////////////////////////////////////////
// The converter exposes functions to convert the top level
// ContainerService resource
//
// All other functions are internal helper functions used
// for converting.
///////////////////////////////////////////////////////////

// ConvertContainerServiceToV20170831 converts an unversioned ContainerService to a v20170831 ContainerService
func ConvertContainerServiceToV20170831AgentPoolOnly(api *ContainerService) *v20170831.HostedMaster {
	v20170831HCP := &v20170831.HostedMaster{}
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
		p.KubernetesVersion = api.OrchestratorProfile.OrchestratorVersion
		p.KubernetesRelease = api.OrchestratorProfile.OrchestratorRelease
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
