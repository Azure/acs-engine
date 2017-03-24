package api

import (
	"github.com/Azure/acs-engine/pkg/api/v20160330"
	"github.com/Azure/acs-engine/pkg/api/v20160930"
	"github.com/Azure/acs-engine/pkg/api/vlabs"
)

///////////////////////////////////////////////////////////
// The converter exposes functions to convert the top level
// ContainerService resource
//
// All other functions are internal helper functions used
// for converting.
///////////////////////////////////////////////////////////

// ConvertV20160930ContainerService converts a v20160930 ContainerService to an unversioned ContainerService
func ConvertV20160930ContainerService(v20160930 *v20160930.ContainerService) *ContainerService {
	c := &ContainerService{}
	c.ID = v20160930.ID
	c.Location = v20160930.Location
	c.Name = v20160930.Name
	convertV20160930ResourcePurchasePlan(&v20160930.Plan, &c.Plan)
	c.Tags = map[string]string{}
	for k, v := range v20160930.Tags {
		c.Tags[k] = v
	}
	c.Type = v20160930.Type
	convertV20160930Properties(&v20160930.Properties, &c.Properties)
	return c
}

// ConvertV20160330ContainerService converts a v20160330 ContainerService to an unversioned ContainerService
func ConvertV20160330ContainerService(v20160330 *v20160330.ContainerService) *ContainerService {
	c := &ContainerService{}
	c.ID = v20160330.ID
	c.Location = v20160330.Location
	c.Name = v20160330.Name
	convertV20160330ResourcePurchasePlan(&v20160330.Plan, &c.Plan)
	c.Tags = map[string]string{}
	for k, v := range v20160330.Tags {
		c.Tags[k] = v
	}
	c.Type = v20160330.Type
	convertV20160330Properties(&v20160330.Properties, &c.Properties)
	return c
}

// ConvertVLabsContainerService converts a vlabs ContainerService to an unversioned ContainerService
func ConvertVLabsContainerService(vlabs *vlabs.ContainerService) *ContainerService {
	c := &ContainerService{}
	c.ID = vlabs.ID
	c.Location = vlabs.Location
	c.Name = vlabs.Name
	convertVLabsResourcePurchasePlan(&vlabs.Plan, &c.Plan)
	c.Tags = map[string]string{}
	for k, v := range vlabs.Tags {
		c.Tags[k] = v
	}
	c.Type = vlabs.Type
	convertVLabsProperties(&vlabs.Properties, &c.Properties)
	return c
}

// convertV20160930ResourcePurchasePlan converts a v20160930 ResourcePurchasePlan to an unversioned ResourcePurchasePlan
func convertV20160930ResourcePurchasePlan(v20160930 *v20160930.ResourcePurchasePlan, api *ResourcePurchasePlan) {
	api.Name = v20160930.Name
	api.Product = v20160930.Product
	api.PromotionCode = v20160930.PromotionCode
	api.Publisher = v20160930.Publisher
}

// convertV20160330ResourcePurchasePlan converts a v20160330 ResourcePurchasePlan to an unversioned ResourcePurchasePlan
func convertV20160330ResourcePurchasePlan(v20160330 *v20160330.ResourcePurchasePlan, api *ResourcePurchasePlan) {
	api.Name = v20160330.Name
	api.Product = v20160330.Product
	api.PromotionCode = v20160330.PromotionCode
	api.Publisher = v20160330.Publisher
}

// convertVLabsResourcePurchasePlan converts a vlabs ResourcePurchasePlan to an unversioned ResourcePurchasePlan
func convertVLabsResourcePurchasePlan(vlabs *vlabs.ResourcePurchasePlan, api *ResourcePurchasePlan) {
	api.Name = vlabs.Name
	api.Product = vlabs.Product
	api.PromotionCode = vlabs.PromotionCode
	api.Publisher = vlabs.Publisher
}

func convertV20160930Properties(v20160930 *v20160930.Properties, api *Properties) {
	api.ProvisioningState = ProvisioningState(v20160930.ProvisioningState)
	convertV20160930OrchestratorProfile(&v20160930.OrchestratorProfile, &api.OrchestratorProfile)
	convertV20160930MasterProfile(&v20160930.MasterProfile, &api.MasterProfile)
	api.AgentPoolProfiles = []AgentPoolProfile{}
	for _, p := range v20160930.AgentPoolProfiles {
		apiProfile := &AgentPoolProfile{}
		convertV20160930AgentPoolProfile(&p, apiProfile)
		api.AgentPoolProfiles = append(api.AgentPoolProfiles, *apiProfile)
	}
	convertV20160930LinuxProfile(&v20160930.LinuxProfile, &api.LinuxProfile)
	convertV20160930WindowsProfile(&v20160930.WindowsProfile, &api.WindowsProfile)
	convertV20160930DiagnosticsProfile(&v20160930.DiagnosticsProfile, &api.DiagnosticsProfile)
	convertV20160930JumpboxProfile(&v20160930.JumpboxProfile, &api.JumpboxProfile)
	convertV20160930ServicePrincipalProfile(&v20160930.ServicePrincipalProfile, &api.ServicePrincipalProfile)
	convertV20160930CustomProfile(&v20160930.CustomProfile, &api.CustomProfile)
}

func convertV20160330Properties(v20160330 *v20160330.Properties, api *Properties) {
	api.ProvisioningState = ProvisioningState(v20160330.ProvisioningState)
	convertV20160330OrchestratorProfile(&v20160330.OrchestratorProfile, &api.OrchestratorProfile)
	convertV20160330MasterProfile(&v20160330.MasterProfile, &api.MasterProfile)
	api.AgentPoolProfiles = []AgentPoolProfile{}
	for _, p := range v20160330.AgentPoolProfiles {
		apiProfile := &AgentPoolProfile{}
		convertV20160330AgentPoolProfile(&p, apiProfile)
		api.AgentPoolProfiles = append(api.AgentPoolProfiles, *apiProfile)
	}
	convertV20160330LinuxProfile(&v20160330.LinuxProfile, &api.LinuxProfile)
	convertV20160330WindowsProfile(&v20160330.WindowsProfile, &api.WindowsProfile)
	convertV20160330DiagnosticsProfile(&v20160330.DiagnosticsProfile, &api.DiagnosticsProfile)
	convertV20160330JumpboxProfile(&v20160330.JumpboxProfile, &api.JumpboxProfile)
}

func convertVLabsProperties(vlabs *vlabs.Properties, api *Properties) {
	api.ProvisioningState = ProvisioningState(vlabs.ProvisioningState)
	convertVLabsOrchestratorProfile(&vlabs.OrchestratorProfile, &api.OrchestratorProfile)
	convertVLabsMasterProfile(&vlabs.MasterProfile, &api.MasterProfile)
	api.AgentPoolProfiles = []AgentPoolProfile{}
	for _, p := range vlabs.AgentPoolProfiles {
		apiProfile := &AgentPoolProfile{}
		convertVLabsAgentPoolProfile(&p, apiProfile)
		api.AgentPoolProfiles = append(api.AgentPoolProfiles, *apiProfile)
	}
	convertVLabsLinuxProfile(&vlabs.LinuxProfile, &api.LinuxProfile)
	convertVLabsWindowsProfile(&vlabs.WindowsProfile, &api.WindowsProfile)
	convertVLabsServicePrincipalProfile(&vlabs.ServicePrincipalProfile, &api.ServicePrincipalProfile)
	convertVLabsCertificateProfile(&vlabs.CertificateProfile, &api.CertificateProfile)
}

func convertV20160930LinuxProfile(v20160930 *v20160930.LinuxProfile, api *LinuxProfile) {
	api.AdminUsername = v20160930.AdminUsername
	api.SSH.PublicKeys = []struct {
		KeyData string `json:"keyData"`
	}{}
	for _, d := range v20160930.SSH.PublicKeys {
		api.SSH.PublicKeys = append(api.SSH.PublicKeys, d)
	}
}

func convertV20160330LinuxProfile(v20160330 *v20160330.LinuxProfile, api *LinuxProfile) {
	api.AdminUsername = v20160330.AdminUsername
	api.SSH.PublicKeys = []struct {
		KeyData string `json:"keyData"`
	}{}
	for _, d := range v20160330.SSH.PublicKeys {
		api.SSH.PublicKeys = append(api.SSH.PublicKeys, d)
	}
}

func convertVLabsLinuxProfile(vlabs *vlabs.LinuxProfile, api *LinuxProfile) {
	api.AdminUsername = vlabs.AdminUsername
	api.SSH.PublicKeys = []struct {
		KeyData string `json:"keyData"`
	}{}
	for _, d := range vlabs.SSH.PublicKeys {
		api.SSH.PublicKeys = append(api.SSH.PublicKeys, d)
	}
	api.Secrets = []KeyVaultSecrets{}
	for _, s := range vlabs.Secrets {
		secret := &KeyVaultSecrets{}
		convertVLabsKeyVaultSecrets(&s, secret)
		api.Secrets = append(api.Secrets, *secret)
	}
}

func convertV20160930WindowsProfile(v20160930 *v20160930.WindowsProfile, api *WindowsProfile) {
	api.AdminUsername = v20160930.AdminUsername
	api.AdminPassword = v20160930.AdminPassword
}

func convertV20160330WindowsProfile(v20160330 *v20160330.WindowsProfile, api *WindowsProfile) {
	api.AdminUsername = v20160330.AdminUsername
	api.AdminPassword = v20160330.AdminPassword
}

func convertVLabsWindowsProfile(vlabs *vlabs.WindowsProfile, api *WindowsProfile) {
	api.AdminUsername = vlabs.AdminUsername
	api.AdminPassword = vlabs.AdminPassword
	api.Secrets = []KeyVaultSecrets{}
	for _, s := range vlabs.Secrets {
		secret := &KeyVaultSecrets{}
		convertVLabsKeyVaultSecrets(&s, secret)
		api.Secrets = append(api.Secrets, *secret)
	}
}

func convertV20160930OrchestratorProfile(v20160930 *v20160930.OrchestratorProfile, api *OrchestratorProfile) {
	api.OrchestratorType = OrchestratorType(v20160930.OrchestratorType)
}

func convertV20160330OrchestratorProfile(v20160330 *v20160330.OrchestratorProfile, api *OrchestratorProfile) {
	api.OrchestratorType = OrchestratorType(v20160330.OrchestratorType)
}

func convertVLabsOrchestratorProfile(vlabs *vlabs.OrchestratorProfile, api *OrchestratorProfile) {
	api.OrchestratorType = OrchestratorType(vlabs.OrchestratorType)
	if api.OrchestratorType == Kubernetes {
		convertVLabsKubernetesConfig(&vlabs.KubernetesConfig, &api.KubernetesConfig)
	}
}

func convertVLabsKubernetesConfig(vlabs *vlabs.KubernetesConfig, api *KubernetesConfig) {
	api.KubernetesImageBase = vlabs.KubernetesImageBase
	api.NetworkPolicy = vlabs.NetworkPolicy
}

func convertV20160930MasterProfile(v20160930 *v20160930.MasterProfile, api *MasterProfile) {
	api.Count = v20160930.Count
	api.DNSPrefix = v20160930.DNSPrefix
	api.FQDN = v20160930.FQDN
	api.Subnet = v20160930.GetSubnet()
	// Set default VMSize
	api.VMSize = "Standard_D2_v2"
}

func convertV20160330MasterProfile(v20160330 *v20160330.MasterProfile, api *MasterProfile) {
	api.Count = v20160330.Count
	api.DNSPrefix = v20160330.DNSPrefix
	api.FQDN = v20160330.FQDN
	api.Subnet = v20160330.GetSubnet()
}

func convertVLabsMasterProfile(vlabs *vlabs.MasterProfile, api *MasterProfile) {
	api.Count = vlabs.Count
	api.DNSPrefix = vlabs.DNSPrefix
	api.VMSize = vlabs.VMSize
	api.VnetSubnetID = vlabs.VnetSubnetID
	api.FirstConsecutiveStaticIP = vlabs.FirstConsecutiveStaticIP
	api.Subnet = vlabs.GetSubnet()
	api.IPAddressCount = vlabs.IPAddressCount
	api.FQDN = vlabs.FQDN
	api.StorageProfile = vlabs.StorageProfile
}

func convertV20160930AgentPoolProfile(v20160930 *v20160930.AgentPoolProfile, api *AgentPoolProfile) {
	api.Name = v20160930.Name
	api.Count = v20160930.Count
	api.VMSize = v20160930.VMSize
	api.DNSPrefix = v20160930.DNSPrefix
	if api.DNSPrefix != "" {
		// Set default Ports when DNSPrefix specified
		api.Ports = []int{80, 443, 8080}
	}
	api.FQDN = v20160930.FQDN
	api.OSType = OSType(v20160930.OSType)
	api.Subnet = v20160930.GetSubnet()
}

func convertV20160330AgentPoolProfile(v20160330 *v20160330.AgentPoolProfile, api *AgentPoolProfile) {
	api.Name = v20160330.Name
	api.Count = v20160330.Count
	api.VMSize = v20160330.VMSize
	api.DNSPrefix = v20160330.DNSPrefix
	api.FQDN = v20160330.FQDN
	api.OSType = OSType(v20160330.OSType)
	api.Subnet = v20160330.GetSubnet()
}

func convertVLabsAgentPoolProfile(vlabs *vlabs.AgentPoolProfile, api *AgentPoolProfile) {
	api.Name = vlabs.Name
	api.Count = vlabs.Count
	api.VMSize = vlabs.VMSize
	api.DNSPrefix = vlabs.DNSPrefix
	api.OSType = OSType(vlabs.OSType)
	api.Ports = []int{}
	api.Ports = append(api.Ports, vlabs.Ports...)
	api.AvailabilityProfile = vlabs.AvailabilityProfile
	api.StorageProfile = vlabs.StorageProfile
	api.DiskSizesGB = []int{}
	api.DiskSizesGB = append(api.DiskSizesGB, vlabs.DiskSizesGB...)
	api.VnetSubnetID = vlabs.VnetSubnetID
	api.Subnet = vlabs.GetSubnet()
	api.IPAddressCount = vlabs.IPAddressCount
	api.FQDN = vlabs.FQDN
	api.CustomNodeLabels = map[string]string{}
	for k, v := range vlabs.CustomNodeLabels {
		api.CustomNodeLabels[k] = v
	}
}

func convertVLabsKeyVaultSecrets(vlabs *vlabs.KeyVaultSecrets, api *KeyVaultSecrets) {
	api.SourceVault = KeyVaultID{ID: vlabs.SourceVault.ID}
	api.VaultCertificates = []KeyVaultCertificate{}
	for _, c := range vlabs.VaultCertificates {
		cert := KeyVaultCertificate{}
		cert.CertificateStore = c.CertificateStore
		cert.CertificateURL = c.CertificateURL
		api.VaultCertificates = append(api.VaultCertificates, cert)
	}
}

func convertV20160930DiagnosticsProfile(v20160930 *v20160930.DiagnosticsProfile, api *DiagnosticsProfile) {
	convertV20160930VMDiagnostics(&v20160930.VMDiagnostics, &api.VMDiagnostics)
}

func convertV20160930VMDiagnostics(v20160930 *v20160930.VMDiagnostics, api *VMDiagnostics) {
	api.Enabled = v20160930.Enabled
	api.StorageURL = v20160930.StorageURL
}

func convertV20160330DiagnosticsProfile(v20160330 *v20160330.DiagnosticsProfile, api *DiagnosticsProfile) {
	convertV20160330VMDiagnostics(&v20160330.VMDiagnostics, &api.VMDiagnostics)
}

func convertV20160330VMDiagnostics(v20160330 *v20160330.VMDiagnostics, api *VMDiagnostics) {
	api.Enabled = v20160330.Enabled
	api.StorageURL = v20160330.StorageURL
}

func convertV20160930JumpboxProfile(v20160930 *v20160930.JumpboxProfile, api *JumpboxProfile) {
	api.OSType = OSType(v20160930.OSType)
	api.DNSPrefix = v20160930.DNSPrefix
	api.FQDN = v20160930.FQDN
}

func convertV20160330JumpboxProfile(v20160330 *v20160330.JumpboxProfile, api *JumpboxProfile) {
	api.OSType = OSType(v20160330.OSType)
	api.DNSPrefix = v20160330.DNSPrefix
	api.FQDN = v20160330.FQDN
}

func convertV20160930ServicePrincipalProfile(v20160930 *v20160930.ServicePrincipalProfile, api *ServicePrincipalProfile) {
	api.ClientID = v20160930.ClientID
	api.Secret = v20160930.Secret
}

func convertVLabsServicePrincipalProfile(vlabs *vlabs.ServicePrincipalProfile, api *ServicePrincipalProfile) {
	api.ClientID = vlabs.ClientID
	api.Secret = vlabs.Secret
}

func convertV20160930CustomProfile(v20160930 *v20160930.CustomProfile, api *CustomProfile) {
	api.Orchestrator = v20160930.Orchestrator
}

func convertVLabsCertificateProfile(vlabs *vlabs.CertificateProfile, api *CertificateProfile) {
	api.CaCertificate = vlabs.CaCertificate
	api.APIServerCertificate = vlabs.APIServerCertificate
	api.APIServerPrivateKey = vlabs.APIServerPrivateKey
	api.ClientCertificate = vlabs.ClientCertificate
	api.ClientPrivateKey = vlabs.ClientPrivateKey
	api.KubeConfigCertificate = vlabs.KubeConfigCertificate
	api.KubeConfigPrivateKey = vlabs.KubeConfigPrivateKey
	api.SetCAPrivateKey(vlabs.GetCAPrivateKey())
}
