package api

import (
	"strings"

	"github.com/Azure/acs-engine/pkg/api/v20160330"
	"github.com/Azure/acs-engine/pkg/api/v20160930"
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

// ConvertContainerServiceToV20160930 converts an unversioned ContainerService to a v20160930 ContainerService to
func ConvertContainerServiceToV20160930(api *ContainerService) *v20160930.ContainerService {
	v20160930 := &v20160930.ContainerService{}
	v20160930.ID = api.ID
	v20160930.Location = api.Location
	v20160930.Name = api.Name
	convertResourcePurchasePlanToV20160930(&api.Plan, &v20160930.Plan)
	v20160930.Tags = map[string]string{}
	for k, v := range api.Tags {
		v20160930.Tags[k] = v
	}
	v20160930.Type = api.Type
	convertPropertiesToV20160930(&api.Properties, &v20160930.Properties)
	return v20160930
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

// convertResourcePurchasePlanToV20160930 converts a v20160930 ResourcePurchasePlan to an unversioned ResourcePurchasePlan
func convertResourcePurchasePlanToV20160930(api *ResourcePurchasePlan, v20160930 *v20160930.ResourcePurchasePlan) {
	v20160930.Name = api.Name
	v20160930.Product = api.Product
	v20160930.PromotionCode = api.PromotionCode
	v20160930.Publisher = api.Publisher
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

func convertPropertiesToV20160930(api *Properties, p *v20160930.Properties) {
	p.ProvisioningState = v20160930.ProvisioningState(api.ProvisioningState)
	convertOrchestratorProfileToV20160930(&api.OrchestratorProfile, &p.OrchestratorProfile)
	convertMasterProfileToV20160930(&api.MasterProfile, &p.MasterProfile)
	p.AgentPoolProfiles = []v20160930.AgentPoolProfile{}
	for _, apiProfile := range api.AgentPoolProfiles {
		v20160930Profile := &v20160930.AgentPoolProfile{}
		convertAgentPoolProfileToV20160930(&apiProfile, v20160930Profile)
		p.AgentPoolProfiles = append(p.AgentPoolProfiles, *v20160930Profile)
	}
	convertLinuxProfileToV20160930(&api.LinuxProfile, &p.LinuxProfile)
	convertWindowsProfileToV20160930(&api.WindowsProfile, &p.WindowsProfile)
	convertDiagnosticsProfileToV20160930(&api.DiagnosticsProfile, &p.DiagnosticsProfile)
	convertJumpboxProfileToV20160930(&api.JumpboxProfile, &p.JumpboxProfile)
	convertServicePrincipalProfileToV20160930(&api.ServicePrincipalProfile, &p.ServicePrincipalProfile)
	convertCustomProfileToV20160930(&api.CustomProfile, &p.CustomProfile)
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

func convertLinuxProfileToV20160930(api *LinuxProfile, v20160930 *v20160930.LinuxProfile) {
	v20160930.AdminUsername = api.AdminUsername
	v20160930.SSH.PublicKeys = []struct {
		KeyData string `json:"keyData"`
	}{}
	for _, d := range api.SSH.PublicKeys {
		v20160930.SSH.PublicKeys = append(v20160930.SSH.PublicKeys, d)
	}
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

func convertWindowsProfileToV20160930(api *WindowsProfile, v20160930 *v20160930.WindowsProfile) {
	v20160930.AdminUsername = api.AdminUsername
	v20160930.AdminPassword = api.AdminPassword
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

func convertOrchestratorProfileToV20160930(api *OrchestratorProfile, o *v20160930.OrchestratorProfile) {
	if strings.HasPrefix(string(api.OrchestratorType), string(v20160930.DCOS)) {
		o.OrchestratorType = v20160930.OrchestratorType(v20160930.DCOS)
	} else {
		o.OrchestratorType = v20160930.OrchestratorType(api.OrchestratorType)
	}

	o.DCOSConfig.DCOS173_BootstrapDownloadURL = api.DCOSConfig.DCOS173_BootstrapDownloadURL
	o.DCOSConfig.DCOS184_BootstrapDownloadURL = api.DCOSConfig.DCOS184_BootstrapDownloadURL
	o.DCOSConfig.DCOS187_BootstrapDownloadURL = api.DCOSConfig.DCOS187_BootstrapDownloadURL
}

func convertOrchestratorProfileToV20160330(api *OrchestratorProfile, o *v20160330.OrchestratorProfile) {
	o.OrchestratorType = v20160330.OrchestratorType(api.OrchestratorType)
	o.DCOSConfig.DCOS173_BootstrapDownloadURL = api.DCOSConfig.DCOS173_BootstrapDownloadURL
	o.DCOSConfig.DCOS184_BootstrapDownloadURL = api.DCOSConfig.DCOS184_BootstrapDownloadURL
	o.DCOSConfig.DCOS187_BootstrapDownloadURL = api.DCOSConfig.DCOS187_BootstrapDownloadURL
}

func convertOrchestratorProfileToVLabs(api *OrchestratorProfile, o *vlabs.OrchestratorProfile) {
	o.OrchestratorType = vlabs.OrchestratorType(api.OrchestratorType)
	o.DCOSConfig.DCOS173_BootstrapDownloadURL = api.DCOSConfig.DCOS173_BootstrapDownloadURL
	o.DCOSConfig.DCOS184_BootstrapDownloadURL = api.DCOSConfig.DCOS184_BootstrapDownloadURL
	o.DCOSConfig.DCOS187_BootstrapDownloadURL = api.DCOSConfig.DCOS187_BootstrapDownloadURL
}

func convertMasterProfileToV20160930(api *MasterProfile, v20160930 *v20160930.MasterProfile) {
	v20160930.Count = api.Count
	v20160930.DNSPrefix = api.DNSPrefix
	v20160930.FQDN = api.FQDN
	v20160930.SetSubnet(api.Subnet)
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
	vlabsProfile.StorageProfile = api.StorageProfile
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

func convertAgentPoolProfileToV20160930(api *AgentPoolProfile, p *v20160930.AgentPoolProfile) {
	p.Name = api.Name
	p.Count = api.Count
	p.VMSize = api.VMSize
	p.DNSPrefix = api.DNSPrefix
	p.FQDN = api.FQDN
	p.OSType = v20160930.OSType(api.OSType)
	p.SetSubnet(api.Subnet)
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

func convertDiagnosticsProfileToV20160930(api *DiagnosticsProfile, v20160930 *v20160930.DiagnosticsProfile) {
	convertVMDiagnosticsToV20160930(&api.VMDiagnostics, &v20160930.VMDiagnostics)
}

func convertVMDiagnosticsToV20160930(api *VMDiagnostics, v20160930 *v20160930.VMDiagnostics) {
	v20160930.Enabled = api.Enabled
	v20160930.StorageURL = api.StorageURL
}

func convertDiagnosticsProfileToV20160330(api *DiagnosticsProfile, v20160330 *v20160330.DiagnosticsProfile) {
	convertVMDiagnosticsToV20160330(&api.VMDiagnostics, &v20160330.VMDiagnostics)
}

func convertVMDiagnosticsToV20160330(api *VMDiagnostics, v20160330 *v20160330.VMDiagnostics) {
	v20160330.Enabled = api.Enabled
	v20160330.StorageURL = api.StorageURL
}

func convertJumpboxProfileToV20160930(api *JumpboxProfile, jb *v20160930.JumpboxProfile) {
	jb.OSType = v20160930.OSType(api.OSType)
	jb.DNSPrefix = api.DNSPrefix
	jb.FQDN = api.FQDN
}

func convertJumpboxProfileToV20160330(api *JumpboxProfile, jb *v20160330.JumpboxProfile) {
	jb.OSType = v20160330.OSType(api.OSType)
	jb.DNSPrefix = api.DNSPrefix
	jb.FQDN = api.FQDN
}

func convertServicePrincipalProfileToV20160930(api *ServicePrincipalProfile, v20160930 *v20160930.ServicePrincipalProfile) {
	v20160930.ClientID = api.ClientID
	v20160930.Secret = api.Secret
}

func convertCustomProfileToV20160930(api *CustomProfile, v20160930 *v20160930.CustomProfile) {
	v20160930.Orchestrator = api.Orchestrator
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
