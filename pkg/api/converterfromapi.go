package api

import (
	"github.com/Azure/acs-engine/pkg/api/v20160330"
	"github.com/Azure/acs-engine/pkg/api/vlabs"
)

///////////////////////////////////////////////////////////
// The converter exposes functions to convert the 2 top
// level resources:
// 1. Subscription
// 2. ContainerService
//
// All other functions are internal helper functions used
// for converting.
///////////////////////////////////////////////////////////

// ConvertSubscriptionToV20160330 converts a v20160330 Subscription to an unversioned Subscription
func ConvertSubscriptionToV20160330(api *Subscription) *v20160330.Subscription {
	s := &v20160330.Subscription{}
	s.ID = api.ID
	s.State = v20160330.SubscriptionState(api.State)
	return s
}

// ConvertSubscriptionToVLabs converts a vlabs Subscription to an unversioned Subscription
func ConvertSubscriptionToVLabs(api *Subscription) *vlabs.Subscription {
	s := &vlabs.Subscription{}
	s.ID = api.ID
	s.State = vlabs.SubscriptionState(api.State)
	return s
}

// ConvertContainerServiceToV20160330 converts a v20160330 ContainerService to an unversioned ContainerService
func ConvertContainerServiceToV20160330(api *ContainerService) *v20160330.ContainerService {
	v20160330 := &v20160330.ContainerService{}
	v20160330.ID = api.ID
	v20160330.Location = api.Location
	v20160330.Name = api.Name
	convertResourcePurchasePlanToV20160330(&api.Plan, &v20160330.Plan)
	v20160330.Tags = map[string]string{}
	for k, v := range api.Tags {
		v20160330.Tags[k] = v
	}
	v20160330.Type = api.Type
	convertPropertiesToV20160330(&api.Properties, &v20160330.Properties)
	return v20160330
}

// ConvertContainerServiceToVLabs converts a vlabs ContainerService to an unversioned ContainerService
func ConvertContainerServiceToVLabs(api *ContainerService) *vlabs.ContainerService {
	vlabs := &vlabs.ContainerService{}
	vlabs.ID = api.ID
	vlabs.Location = api.Location
	vlabs.Name = api.Name
	convertResourcePurchasePlanToVLabs(&api.Plan, &vlabs.Plan)
	vlabs.Tags = map[string]string{}
	for k, v := range api.Tags {
		vlabs.Tags[k] = v
	}
	vlabs.Type = api.Type
	convertPropertiesToVLabs(&api.Properties, &vlabs.Properties)
	return vlabs
}

// convertResourcePurchasePlanToV20160330 converts a v20160330 ResourcePurchasePlan to an unversioned ResourcePurchasePlan
func convertResourcePurchasePlanToV20160330(api *ResourcePurchasePlan, v20160330 *v20160330.ResourcePurchasePlan) {
	v20160330.Name = api.Name
	v20160330.Product = api.Product
	v20160330.PromotionCode = api.PromotionCode
	v20160330.Publisher = api.Publisher
}

// convertResourcePurchasePlanToVLabs converts a vlabs ResourcePurchasePlan to an unversioned ResourcePurchasePlan
func convertResourcePurchasePlanToVLabs(api *ResourcePurchasePlan, vlabs *vlabs.ResourcePurchasePlan) {
	vlabs.Name = api.Name
	vlabs.Product = api.Product
	vlabs.PromotionCode = api.PromotionCode
	vlabs.Publisher = api.Publisher
}

func convertPropertiesToV20160330(api *Properties, p *v20160330.Properties) {
	p.ProvisioningState = v20160330.ProvisioningState(api.ProvisioningState)
	convertOrchestratorProfileToV20160330(&api.OrchestratorProfile, &p.OrchestratorProfile)
	convertMasterProfileToV20160330(&api.MasterProfile, &p.MasterProfile)
	p.AgentPoolProfiles = []v20160330.AgentPoolProfile{}
	for _, apiProfile := range api.AgentPoolProfiles {
		v20160330Profile := &v20160330.AgentPoolProfile{}
		convertAgentPoolProfileToV20160330(&apiProfile, v20160330Profile)
		p.AgentPoolProfiles = append(p.AgentPoolProfiles, *v20160330Profile)
	}
	convertLinuxProfileToV20160330(&api.LinuxProfile, &p.LinuxProfile)
	convertWindowsProfileToV20160330(&api.WindowsProfile, &p.WindowsProfile)
	convertDiagnosticsProfileToV20160330(&api.DiagnosticsProfile, &p.DiagnosticsProfile)
	convertJumpboxProfileToV20160330(&api.JumpboxProfile, &p.JumpboxProfile)
}

func convertPropertiesToVLabs(api *Properties, vlabsProps *vlabs.Properties) {
	vlabsProps.ProvisioningState = vlabs.ProvisioningState(api.ProvisioningState)
	convertOrchestratorProfileToVLabs(&api.OrchestratorProfile, &vlabsProps.OrchestratorProfile)
	convertMasterProfileToVLabs(&api.MasterProfile, &vlabsProps.MasterProfile)
	vlabsProps.AgentPoolProfiles = []vlabs.AgentPoolProfile{}
	for _, apiProfile := range api.AgentPoolProfiles {
		vlabsProfile := &vlabs.AgentPoolProfile{}
		convertAgentPoolProfileToVLabs(&apiProfile, vlabsProfile)
		vlabsProps.AgentPoolProfiles = append(vlabsProps.AgentPoolProfiles, *vlabsProfile)
	}
	convertLinuxProfileToVLabs(&api.LinuxProfile, &vlabsProps.LinuxProfile)
	convertWindowsProfileToVLabs(&api.WindowsProfile, &vlabsProps.WindowsProfile)
	convertServicePrincipalProfileToVLabs(&api.ServicePrincipalProfile, &vlabsProps.ServicePrincipalProfile)
	convertCertificateProfileToVLabs(&api.CertificateProfile, &vlabsProps.CertificateProfile)
}

func convertLinuxProfileToV20160330(api *LinuxProfile, v20160330 *v20160330.LinuxProfile) {
	v20160330.AdminUsername = api.AdminUsername
	v20160330.SSH.PublicKeys = []struct {
		KeyData string `json:"keyData"`
	}{}
	for _, d := range api.SSH.PublicKeys {
		v20160330.SSH.PublicKeys = append(v20160330.SSH.PublicKeys, d)
	}
}

func convertLinuxProfileToVLabs(api *LinuxProfile, vlabsProfile *vlabs.LinuxProfile) {
	vlabsProfile.AdminUsername = api.AdminUsername
	vlabsProfile.SSH.PublicKeys = []struct {
		KeyData string `json:"keyData"`
	}{}
	for _, d := range api.SSH.PublicKeys {
		vlabsProfile.SSH.PublicKeys = append(vlabsProfile.SSH.PublicKeys, d)
	}
	vlabsProfile.Secrets = []vlabs.KeyVaultSecrets{}
	for _, s := range api.Secrets {
		secret := &vlabs.KeyVaultSecrets{}
		convertKeyVaultSecretsToVlabs(&s, secret)
		vlabsProfile.Secrets = append(vlabsProfile.Secrets, *secret)
	}
}

func convertWindowsProfileToV20160330(api *WindowsProfile, v20160330 *v20160330.WindowsProfile) {
	v20160330.AdminUsername = api.AdminUsername
	v20160330.AdminPassword = api.AdminPassword
}

func convertWindowsProfileToVLabs(api *WindowsProfile, vlabsProfile *vlabs.WindowsProfile) {
	vlabsProfile.AdminUsername = api.AdminUsername
	vlabsProfile.AdminPassword = api.AdminPassword
	vlabsProfile.Secrets = []vlabs.KeyVaultSecrets{}
	for _, s := range api.Secrets {
		secret := &vlabs.KeyVaultSecrets{}
		convertKeyVaultSecretsToVlabs(&s, secret)
		vlabsProfile.Secrets = append(vlabsProfile.Secrets, *secret)
	}
}

func convertOrchestratorProfileToV20160330(api *OrchestratorProfile, o *v20160330.OrchestratorProfile) {
	o.OrchestratorType = v20160330.OrchestratorType(api.OrchestratorType)
}

func convertOrchestratorProfileToVLabs(api *OrchestratorProfile, o *vlabs.OrchestratorProfile) {
	o.OrchestratorType = vlabs.OrchestratorType(api.OrchestratorType)
}

func convertMasterProfileToV20160330(api *MasterProfile, v20160330 *v20160330.MasterProfile) {
	v20160330.Count = api.Count
	v20160330.DNSPrefix = api.DNSPrefix
	v20160330.FQDN = api.FQDN
	v20160330.SetSubnet(api.Subnet)
}

func convertMasterProfileToVLabs(api *MasterProfile, vlabsProfile *vlabs.MasterProfile) {
	vlabsProfile.Count = api.Count
	vlabsProfile.DNSPrefix = api.DNSPrefix
	vlabsProfile.VMSize = api.VMSize
	vlabsProfile.VnetSubnetID = api.VnetSubnetID
	vlabsProfile.FirstConsecutiveStaticIP = api.FirstConsecutiveStaticIP
	vlabsProfile.SetSubnet(api.Subnet)
	vlabsProfile.FQDN = api.FQDN
}

func convertKeyVaultSecretsToVlabs(api *KeyVaultSecrets, vlabsSecrets *vlabs.KeyVaultSecrets) {
	vlabsSecrets.SourceVault = vlabs.KeyVaultID{ID: api.SourceVault.ID}
	vlabsSecrets.VaultCertificates = []vlabs.KeyVaultCertificate{}
	for _, c := range api.VaultCertificates {
		cert := vlabs.KeyVaultCertificate{}
		cert.CertificateStore = c.CertificateStore
		cert.CertificateURL = c.CertificateURL
		vlabsSecrets.VaultCertificates = append(vlabsSecrets.VaultCertificates, cert)
	}
}

func convertAgentPoolProfileToV20160330(api *AgentPoolProfile, p *v20160330.AgentPoolProfile) {
	p.Name = api.Name
	p.Count = api.Count
	p.VMSize = api.VMSize
	p.DNSPrefix = api.DNSPrefix
	p.FQDN = api.FQDN
	p.OSType = v20160330.OSType(api.OSType)
	p.SetSubnet(api.Subnet)
}

func convertAgentPoolProfileToVLabs(api *AgentPoolProfile, p *vlabs.AgentPoolProfile) {
	p.Name = api.Name
	p.Count = api.Count
	p.VMSize = api.VMSize
	p.DNSPrefix = api.DNSPrefix
	p.OSType = vlabs.OSType(api.OSType)
	p.Ports = []int{}
	p.Ports = append(p.Ports, api.Ports...)
	p.AvailabilityProfile = api.AvailabilityProfile
	p.StorageProfile = api.StorageProfile
	p.DiskSizesGB = []int{}
	p.DiskSizesGB = append(p.DiskSizesGB, api.DiskSizesGB...)
	p.VnetSubnetID = api.VnetSubnetID
	p.SetSubnet(api.Subnet)
	p.FQDN = api.FQDN
}

func convertDiagnosticsProfileToV20160330(api *DiagnosticsProfile, v20160330 *v20160330.DiagnosticsProfile) {
	convertVMDiagnosticsToV20160330(&api.VMDiagnostics, &v20160330.VMDiagnostics)
}

func convertVMDiagnosticsToV20160330(api *VMDiagnostics, v20160330 *v20160330.VMDiagnostics) {
	v20160330.Enabled = api.Enabled
	v20160330.StorageURL = api.StorageURL
}

func convertJumpboxProfileToV20160330(api *JumpboxProfile, jb *v20160330.JumpboxProfile) {
	jb.OSType = v20160330.OSType(api.OSType)
	jb.DNSPrefix = api.DNSPrefix
	jb.FQDN = api.FQDN
}

func convertServicePrincipalProfileToVLabs(api *ServicePrincipalProfile, vlabs *vlabs.ServicePrincipalProfile) {
	vlabs.ClientID = api.ClientID
	vlabs.Secret = api.Secret
}

func convertCertificateProfileToVLabs(api *CertificateProfile, vlabs *vlabs.CertificateProfile) {
	vlabs.CaCertificate = api.CaCertificate
	vlabs.APIServerCertificate = api.APIServerCertificate
	vlabs.APIServerPrivateKey = api.APIServerPrivateKey
	vlabs.ClientCertificate = api.ClientCertificate
	vlabs.ClientPrivateKey = api.ClientPrivateKey
	vlabs.KubeConfigCertificate = api.KubeConfigCertificate
	vlabs.KubeConfigPrivateKey = api.KubeConfigPrivateKey
	vlabs.SetCAPrivateKey(api.GetCAPrivateKey())
}
