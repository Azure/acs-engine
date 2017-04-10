package api

import (
	"github.com/Azure/acs-engine/pkg/api/v20160330"
	"github.com/Azure/acs-engine/pkg/api/v20160930"
	"github.com/Azure/acs-engine/pkg/api/v20170131"
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
	if v20160930.Plan != nil {
		c.Plan = &ResourcePurchasePlan{}
		convertV20160930ResourcePurchasePlan(v20160930.Plan, c.Plan)
	}
	c.Tags = map[string]string{}
	for k, v := range v20160930.Tags {
		c.Tags[k] = v
	}
	c.Type = v20160930.Type
	c.Properties = &Properties{}
	convertV20160930Properties(v20160930.Properties, c.Properties)
	return c
}

// ConvertV20160330ContainerService converts a v20160330 ContainerService to an unversioned ContainerService
func ConvertV20160330ContainerService(v20160330 *v20160330.ContainerService) *ContainerService {
	c := &ContainerService{}
	c.ID = v20160330.ID
	c.Location = v20160330.Location
	c.Name = v20160330.Name
	if v20160330.Plan != nil {
		c.Plan = &ResourcePurchasePlan{}
		convertV20160330ResourcePurchasePlan(v20160330.Plan, c.Plan)
	}
	c.Tags = map[string]string{}
	for k, v := range v20160330.Tags {
		c.Tags[k] = v
	}
	c.Type = v20160330.Type
	c.Properties = &Properties{}
	convertV20160330Properties(v20160330.Properties, c.Properties)
	return c
}

// ConvertV20170131ContainerService converts a v20170131 ContainerService to an unversioned ContainerService
func ConvertV20170131ContainerService(v20170131 *v20170131.ContainerService) *ContainerService {
	c := &ContainerService{}
	c.ID = v20170131.ID
	c.Location = v20170131.Location
	c.Name = v20170131.Name
	if v20170131.Plan != nil {
		c.Plan = &ResourcePurchasePlan{}
		convertV20170131ResourcePurchasePlan(v20170131.Plan, c.Plan)
	}
	c.Tags = map[string]string{}
	for k, v := range v20170131.Tags {
		c.Tags[k] = v
	}
	c.Type = v20170131.Type
	c.Properties = &Properties{}
	convertV20170131Properties(v20170131.Properties, c.Properties)
	return c
}

// ConvertVLabsContainerService converts a vlabs ContainerService to an unversioned ContainerService
func ConvertVLabsContainerService(vlabs *vlabs.ContainerService) *ContainerService {
	c := &ContainerService{}
	c.ID = vlabs.ID
	c.Location = vlabs.Location
	c.Name = vlabs.Name
	if vlabs.Plan != nil {
		c.Plan = &ResourcePurchasePlan{}
		convertVLabsResourcePurchasePlan(vlabs.Plan, c.Plan)
	}
	c.Tags = map[string]string{}
	for k, v := range vlabs.Tags {
		c.Tags[k] = v
	}
	c.Type = vlabs.Type
	c.Properties = &Properties{}
	convertVLabsProperties(vlabs.Properties, c.Properties)
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

// convertV20170131ResourcePurchasePlan converts a v20170131 ResourcePurchasePlan to an unversioned ResourcePurchasePlan
func convertV20170131ResourcePurchasePlan(v20170131 *v20170131.ResourcePurchasePlan, api *ResourcePurchasePlan) {
	api.Name = v20170131.Name
	api.Product = v20170131.Product
	api.PromotionCode = v20170131.PromotionCode
	api.Publisher = v20170131.Publisher
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
	if v20160930.OrchestratorProfile != nil {
		api.OrchestratorProfile = &OrchestratorProfile{}
		convertV20160930OrchestratorProfile(v20160930.OrchestratorProfile, api.OrchestratorProfile)
	}
	if v20160930.MasterProfile != nil {
		api.MasterProfile = &MasterProfile{}
		convertV20160930MasterProfile(v20160930.MasterProfile, api.MasterProfile)
	}
	api.AgentPoolProfiles = []AgentPoolProfile{}
	for _, p := range v20160930.AgentPoolProfiles {
		apiProfile := &AgentPoolProfile{}
		convertV20160930AgentPoolProfile(&p, apiProfile)
		api.AgentPoolProfiles = append(api.AgentPoolProfiles, *apiProfile)
	}
	if v20160930.LinuxProfile != nil {
		api.LinuxProfile = &LinuxProfile{}
		convertV20160930LinuxProfile(v20160930.LinuxProfile, api.LinuxProfile)
	}
	if v20160930.WindowsProfile != nil {
		api.WindowsProfile = &WindowsProfile{}
		convertV20160930WindowsProfile(v20160930.WindowsProfile, api.WindowsProfile)
	}
	if v20160930.DiagnosticsProfile != nil {
		api.DiagnosticsProfile = &DiagnosticsProfile{}
		convertV20160930DiagnosticsProfile(v20160930.DiagnosticsProfile, api.DiagnosticsProfile)
	}
	if v20160930.JumpboxProfile != nil {
		api.JumpboxProfile = &JumpboxProfile{}
		convertV20160930JumpboxProfile(v20160930.JumpboxProfile, api.JumpboxProfile)
	}
	if v20160930.ServicePrincipalProfile != nil {
		api.ServicePrincipalProfile = &ServicePrincipalProfile{}
		convertV20160930ServicePrincipalProfile(v20160930.ServicePrincipalProfile, api.ServicePrincipalProfile)
	}
	if v20160930.CustomProfile != nil {
		api.CustomProfile = &CustomProfile{}
		convertV20160930CustomProfile(v20160930.CustomProfile, api.CustomProfile)
	}
}

func convertV20160330Properties(v20160330 *v20160330.Properties, api *Properties) {
	api.ProvisioningState = ProvisioningState(v20160330.ProvisioningState)
	if v20160330.OrchestratorProfile != nil {
		api.OrchestratorProfile = &OrchestratorProfile{}
		convertV20160330OrchestratorProfile(v20160330.OrchestratorProfile, api.OrchestratorProfile)
	}
	if v20160330.MasterProfile != nil {
		api.MasterProfile = &MasterProfile{}
		convertV20160330MasterProfile(v20160330.MasterProfile, api.MasterProfile)
	}
	api.AgentPoolProfiles = []AgentPoolProfile{}
	for _, p := range v20160330.AgentPoolProfiles {
		apiProfile := &AgentPoolProfile{}
		convertV20160330AgentPoolProfile(&p, apiProfile)
		api.AgentPoolProfiles = append(api.AgentPoolProfiles, *apiProfile)

	}
	if v20160330.LinuxProfile != nil {
		api.LinuxProfile = &LinuxProfile{}
		convertV20160330LinuxProfile(v20160330.LinuxProfile, api.LinuxProfile)
	}
	if v20160330.WindowsProfile != nil {
		api.WindowsProfile = &WindowsProfile{}
		convertV20160330WindowsProfile(v20160330.WindowsProfile, api.WindowsProfile)
	}
	if v20160330.DiagnosticsProfile != nil {
		api.DiagnosticsProfile = &DiagnosticsProfile{}
		convertV20160330DiagnosticsProfile(v20160330.DiagnosticsProfile, api.DiagnosticsProfile)
	}
	if v20160330.JumpboxProfile != nil {
		api.JumpboxProfile = &JumpboxProfile{}
		convertV20160330JumpboxProfile(v20160330.JumpboxProfile, api.JumpboxProfile)
	}
}

func convertV20170131Properties(v20170131 *v20170131.Properties, api *Properties) {
	api.ProvisioningState = ProvisioningState(v20170131.ProvisioningState)
	if v20170131.OrchestratorProfile != nil {
		api.OrchestratorProfile = &OrchestratorProfile{}
		convertV20170131OrchestratorProfile(v20170131.OrchestratorProfile, api.OrchestratorProfile)
	}
	if v20170131.MasterProfile != nil {
		api.MasterProfile = &MasterProfile{}
		convertV20170131MasterProfile(v20170131.MasterProfile, api.MasterProfile)
	}
	api.AgentPoolProfiles = []AgentPoolProfile{}
	for _, p := range v20170131.AgentPoolProfiles {
		apiProfile := &AgentPoolProfile{}
		convertV20170131AgentPoolProfile(&p, apiProfile)
		api.AgentPoolProfiles = append(api.AgentPoolProfiles, *apiProfile)
	}
	if v20170131.LinuxProfile != nil {
		api.LinuxProfile = &LinuxProfile{}
		convertV20170131LinuxProfile(v20170131.LinuxProfile, api.LinuxProfile)
	}
	if v20170131.WindowsProfile != nil {
		api.WindowsProfile = &WindowsProfile{}
		convertV20170131WindowsProfile(v20170131.WindowsProfile, api.WindowsProfile)
	}
	if v20170131.DiagnosticsProfile != nil {
		api.DiagnosticsProfile = &DiagnosticsProfile{}
		convertV20170131DiagnosticsProfile(v20170131.DiagnosticsProfile, api.DiagnosticsProfile)
	}
	if v20170131.JumpboxProfile != nil {
		api.JumpboxProfile = &JumpboxProfile{}
		convertV20170131JumpboxProfile(v20170131.JumpboxProfile, api.JumpboxProfile)
	}
	if v20170131.ServicePrincipalProfile != nil {
		api.ServicePrincipalProfile = &ServicePrincipalProfile{}
		convertV20170131ServicePrincipalProfile(v20170131.ServicePrincipalProfile, api.ServicePrincipalProfile)
	}
	if v20170131.CustomProfile != nil {
		api.CustomProfile = &CustomProfile{}
		convertV20170131CustomProfile(v20170131.CustomProfile, api.CustomProfile)
	}
}

func convertVLabsProperties(vlabs *vlabs.Properties, api *Properties) {
	api.ProvisioningState = ProvisioningState(vlabs.ProvisioningState)
	if vlabs.OrchestratorProfile != nil {
		api.OrchestratorProfile = &OrchestratorProfile{}
		convertVLabsOrchestratorProfile(vlabs.OrchestratorProfile, api.OrchestratorProfile)
	}
	if vlabs.MasterProfile != nil {
		api.MasterProfile = &MasterProfile{}
		convertVLabsMasterProfile(vlabs.MasterProfile, api.MasterProfile)
	}
	api.AgentPoolProfiles = []AgentPoolProfile{}
	for _, p := range vlabs.AgentPoolProfiles {
		apiProfile := &AgentPoolProfile{}
		convertVLabsAgentPoolProfile(&p, apiProfile)
		api.AgentPoolProfiles = append(api.AgentPoolProfiles, *apiProfile)
	}
	if vlabs.LinuxProfile != nil {
		api.LinuxProfile = &LinuxProfile{}
		convertVLabsLinuxProfile(vlabs.LinuxProfile, api.LinuxProfile)
	}
	if vlabs.WindowsProfile != nil {
		api.WindowsProfile = &WindowsProfile{}
		convertVLabsWindowsProfile(vlabs.WindowsProfile, api.WindowsProfile)
	}
	if vlabs.ServicePrincipalProfile != nil {
		api.ServicePrincipalProfile = &ServicePrincipalProfile{}
		convertVLabsServicePrincipalProfile(vlabs.ServicePrincipalProfile, api.ServicePrincipalProfile)
	}
	if vlabs.CertificateProfile != nil {
		api.CertificateProfile = &CertificateProfile{}
		convertVLabsCertificateProfile(vlabs.CertificateProfile, api.CertificateProfile)
	}
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

func convertV20170131LinuxProfile(v20170131 *v20170131.LinuxProfile, api *LinuxProfile) {
	api.AdminUsername = v20170131.AdminUsername
	api.SSH.PublicKeys = []struct {
		KeyData string `json:"keyData"`
	}{}
	for _, d := range v20170131.SSH.PublicKeys {
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

func convertV20170131WindowsProfile(v20170131 *v20170131.WindowsProfile, api *WindowsProfile) {
	api.AdminUsername = v20170131.AdminUsername
	api.AdminPassword = v20170131.AdminPassword
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

func convertV20170131OrchestratorProfile(v20170131 *v20170131.OrchestratorProfile, api *OrchestratorProfile) {
	api.OrchestratorType = OrchestratorType(v20170131.OrchestratorType)
}

func convertVLabsOrchestratorProfile(vlabs *vlabs.OrchestratorProfile, api *OrchestratorProfile) {
	api.OrchestratorType = OrchestratorType(vlabs.OrchestratorType)
	if api.OrchestratorType == Kubernetes && vlabs.KubernetesConfig != nil {
		api.KubernetesConfig = &KubernetesConfig{}
		convertVLabsKubernetesConfig(vlabs.KubernetesConfig, api.KubernetesConfig)
	}
}

func convertVLabsKubernetesConfig(vlabs *vlabs.KubernetesConfig, api *KubernetesConfig) {
	api.KubernetesImageBase = vlabs.KubernetesImageBase
	api.NetworkPolicy = vlabs.NetworkPolicy
	api.ClusterCIDR = vlabs.ClusterCidr
	api.DnsServiceIP = vlabs.DnsServiceIP
	api.ServiceCIDR = vlabs.ServiceCidr
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
	// Set default VMSize
	api.VMSize = "Standard_D2_v2"
}

func convertV20170131MasterProfile(v20170131 *v20170131.MasterProfile, api *MasterProfile) {
	api.Count = v20170131.Count
	api.DNSPrefix = v20170131.DNSPrefix
	api.FQDN = v20170131.FQDN
	api.Subnet = v20170131.GetSubnet()
	// Set default VMSize
	// TODO: Use azureconst.go to set
	api.VMSize = "Standard_D2_v2"
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
	if api.DNSPrefix != "" {
		// Set default Ports when DNSPrefix specified
		api.Ports = []int{80, 443, 8080}
	}
	api.FQDN = v20160330.FQDN
	api.OSType = OSType(v20160330.OSType)
	api.Subnet = v20160330.GetSubnet()
}

func convertV20170131AgentPoolProfile(v20170131 *v20170131.AgentPoolProfile, api *AgentPoolProfile) {
	api.Name = v20170131.Name
	api.Count = v20170131.Count
	api.VMSize = v20170131.VMSize
	api.DNSPrefix = v20170131.DNSPrefix
	if api.DNSPrefix != "" {
		// Set default Ports when DNSPrefix specified
		api.Ports = []int{80, 443, 8080}
	}
	api.FQDN = v20170131.FQDN
	api.OSType = OSType(v20170131.OSType)
	api.Subnet = v20170131.GetSubnet()
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
	api.SourceVault = &KeyVaultID{ID: vlabs.SourceVault.ID}
	api.VaultCertificates = []KeyVaultCertificate{}
	for _, c := range vlabs.VaultCertificates {
		cert := KeyVaultCertificate{}
		cert.CertificateStore = c.CertificateStore
		cert.CertificateURL = c.CertificateURL
		api.VaultCertificates = append(api.VaultCertificates, cert)
	}
}

func convertV20160930DiagnosticsProfile(v20160930 *v20160930.DiagnosticsProfile, api *DiagnosticsProfile) {
	if v20160930.VMDiagnostics != nil {
		api.VMDiagnostics = &VMDiagnostics{}
		convertV20160930VMDiagnostics(v20160930.VMDiagnostics, api.VMDiagnostics)
	}
}

func convertV20160930VMDiagnostics(v20160930 *v20160930.VMDiagnostics, api *VMDiagnostics) {
	api.Enabled = v20160930.Enabled
	api.StorageURL = v20160930.StorageURL
}

func convertV20160330DiagnosticsProfile(v20160330 *v20160330.DiagnosticsProfile, api *DiagnosticsProfile) {
	if v20160330.VMDiagnostics != nil {
		api.VMDiagnostics = &VMDiagnostics{}
		convertV20160330VMDiagnostics(v20160330.VMDiagnostics, api.VMDiagnostics)
	}
}

func convertV20160330VMDiagnostics(v20160330 *v20160330.VMDiagnostics, api *VMDiagnostics) {
	api.Enabled = v20160330.Enabled
	api.StorageURL = v20160330.StorageURL
}

func convertV20170131DiagnosticsProfile(v20170131 *v20170131.DiagnosticsProfile, api *DiagnosticsProfile) {
	if v20170131.VMDiagnostics != nil {
		api.VMDiagnostics = &VMDiagnostics{}
		convertV20170131VMDiagnostics(v20170131.VMDiagnostics, api.VMDiagnostics)
	}
}

func convertV20170131VMDiagnostics(v20170131 *v20170131.VMDiagnostics, api *VMDiagnostics) {
	api.Enabled = v20170131.Enabled
	api.StorageURL = v20170131.StorageURL
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

func convertV20170131JumpboxProfile(v20170131 *v20170131.JumpboxProfile, api *JumpboxProfile) {
	api.OSType = OSType(v20170131.OSType)
	api.DNSPrefix = v20170131.DNSPrefix
	api.FQDN = v20170131.FQDN
}

func convertV20160930ServicePrincipalProfile(v20160930 *v20160930.ServicePrincipalProfile, api *ServicePrincipalProfile) {
	api.ClientID = v20160930.ClientID
	api.Secret = v20160930.Secret
}

func convertV20170131ServicePrincipalProfile(v20170131 *v20170131.ServicePrincipalProfile, api *ServicePrincipalProfile) {
	api.ClientID = v20170131.ClientID
	api.Secret = v20170131.Secret
}

func convertVLabsServicePrincipalProfile(vlabs *vlabs.ServicePrincipalProfile, api *ServicePrincipalProfile) {
	api.ClientID = vlabs.ClientID
	api.Secret = vlabs.Secret
}

func convertV20160930CustomProfile(v20160930 *v20160930.CustomProfile, api *CustomProfile) {
	api.Orchestrator = v20160930.Orchestrator
}

func convertV20170131CustomProfile(v20170131 *v20170131.CustomProfile, api *CustomProfile) {
	api.Orchestrator = v20170131.Orchestrator
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
