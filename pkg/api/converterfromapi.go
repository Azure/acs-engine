package api

import (
	"regexp"
	"strings"

	"github.com/Azure/acs-engine/pkg/api/v20160330"
	"github.com/Azure/acs-engine/pkg/api/v20160930"
	"github.com/Azure/acs-engine/pkg/api/v20170131"
	"github.com/Azure/acs-engine/pkg/api/v20170701"
	"github.com/Azure/acs-engine/pkg/api/vlabs"
)

///////////////////////////////////////////////////////////
// The converter exposes functions to convert the top level
// ContainerService resource
//
// All other functions are internal helper functions used
// for converting.
///////////////////////////////////////////////////////////

// ConvertContainerServiceToV20160930 converts an unversioned ContainerService to a v20160930 ContainerService
func ConvertContainerServiceToV20160930(api *ContainerService) *v20160930.ContainerService {
	v20160930CS := &v20160930.ContainerService{}
	v20160930CS.ID = api.ID
	v20160930CS.Location = api.Location
	v20160930CS.Name = api.Name
	if api.Plan != nil {
		v20160930CS.Plan = &v20160930.ResourcePurchasePlan{}
		convertResourcePurchasePlanToV20160930(api.Plan, v20160930CS.Plan)
	}
	v20160930CS.Tags = map[string]string{}
	for k, v := range api.Tags {
		v20160930CS.Tags[k] = v
	}
	v20160930CS.Type = api.Type
	v20160930CS.Properties = &v20160930.Properties{}
	convertPropertiesToV20160930(api.Properties, v20160930CS.Properties)
	return v20160930CS
}

// ConvertContainerServiceToV20160330 converts an unversioned ContainerService to a v20160330 ContainerService
func ConvertContainerServiceToV20160330(api *ContainerService) *v20160330.ContainerService {
	v20160330CS := &v20160330.ContainerService{}
	v20160330CS.ID = api.ID
	v20160330CS.Location = api.Location
	v20160330CS.Name = api.Name
	if api.Plan != nil {
		v20160330CS.Plan = &v20160330.ResourcePurchasePlan{}
		convertResourcePurchasePlanToV20160330(api.Plan, v20160330CS.Plan)
	}
	v20160330CS.Tags = map[string]string{}
	for k, v := range api.Tags {
		v20160330CS.Tags[k] = v
	}
	v20160330CS.Type = api.Type
	v20160330CS.Properties = &v20160330.Properties{}
	convertPropertiesToV20160330(api.Properties, v20160330CS.Properties)
	return v20160330CS
}

// ConvertContainerServiceToV20170131 converts an unversioned ContainerService to a v20170131 ContainerService
func ConvertContainerServiceToV20170131(api *ContainerService) *v20170131.ContainerService {
	v20170131CS := &v20170131.ContainerService{}
	v20170131CS.ID = api.ID
	v20170131CS.Location = api.Location
	v20170131CS.Name = api.Name
	if api.Plan != nil {
		v20170131CS.Plan = &v20170131.ResourcePurchasePlan{}
		convertResourcePurchasePlanToV20170131(api.Plan, v20170131CS.Plan)
	}
	v20170131CS.Tags = map[string]string{}
	for k, v := range api.Tags {
		v20170131CS.Tags[k] = v
	}
	v20170131CS.Type = api.Type
	v20170131CS.Properties = &v20170131.Properties{}
	convertPropertiesToV20170131(api.Properties, v20170131CS.Properties)
	return v20170131CS
}

// ConvertContainerServiceToV20170701 converts an unversioned ContainerService to a v20170701 ContainerService
func ConvertContainerServiceToV20170701(api *ContainerService) *v20170701.ContainerService {
	v20170701CS := &v20170701.ContainerService{}
	v20170701CS.ID = api.ID
	v20170701CS.Location = api.Location
	v20170701CS.Name = api.Name
	if api.Plan != nil {
		v20170701CS.Plan = &v20170701.ResourcePurchasePlan{}
		convertResourcePurchasePlanToV20170701(api.Plan, v20170701CS.Plan)
	}
	v20170701CS.Tags = map[string]string{}
	for k, v := range api.Tags {
		v20170701CS.Tags[k] = v
	}
	v20170701CS.Type = api.Type
	v20170701CS.Properties = &v20170701.Properties{}
	convertPropertiesToV20170701(api.Properties, v20170701CS.Properties)
	return v20170701CS
}

// ConvertContainerServiceToVLabs converts an unversioned ContainerService to a vlabs ContainerService
func ConvertContainerServiceToVLabs(api *ContainerService) *vlabs.ContainerService {
	vlabsCS := &vlabs.ContainerService{}
	vlabsCS.ID = api.ID
	vlabsCS.Location = api.Location
	vlabsCS.Name = api.Name
	if api.Plan != nil {
		vlabsCS.Plan = &vlabs.ResourcePurchasePlan{}
		convertResourcePurchasePlanToVLabs(api.Plan, vlabsCS.Plan)
	}
	vlabsCS.Tags = map[string]string{}
	for k, v := range api.Tags {
		vlabsCS.Tags[k] = v
	}
	vlabsCS.Type = api.Type
	vlabsCS.Properties = &vlabs.Properties{}
	convertPropertiesToVLabs(api.Properties, vlabsCS.Properties)
	return vlabsCS
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

// convertResourcePurchasePlanToV20170131 converts an unversioned ResourcePurchasePlan to a v20170131 ResourcePurchasePlan
func convertResourcePurchasePlanToV20170131(api *ResourcePurchasePlan, v20170131 *v20170131.ResourcePurchasePlan) {
	v20170131.Name = api.Name
	v20170131.Product = api.Product
	v20170131.PromotionCode = api.PromotionCode
	v20170131.Publisher = api.Publisher
}

// convertResourcePurchasePlanToV20170701 converts a v20170701 ResourcePurchasePlan to an unversioned ResourcePurchasePlan
func convertResourcePurchasePlanToV20170701(api *ResourcePurchasePlan, v20170701 *v20170701.ResourcePurchasePlan) {
	v20170701.Name = api.Name
	v20170701.Product = api.Product
	v20170701.PromotionCode = api.PromotionCode
	v20170701.Publisher = api.Publisher
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
	if api.OrchestratorProfile != nil {
		p.OrchestratorProfile = &v20160930.OrchestratorProfile{}
		convertOrchestratorProfileToV20160930(api.OrchestratorProfile, p.OrchestratorProfile)
	}
	if api.MasterProfile != nil {
		p.MasterProfile = &v20160930.MasterProfile{}
		convertMasterProfileToV20160930(api.MasterProfile, p.MasterProfile)
	}
	p.AgentPoolProfiles = []*v20160930.AgentPoolProfile{}
	// DCOS conversion logic
	if api.OrchestratorProfile.IsDCOS() && len(api.AgentPoolProfiles) == 2 {
		var privIndex, pubIndex int
		for i, apiProfile := range api.AgentPoolProfiles {
			// We added a pool with a "-public" suffix when converting to API model;
			// we don't want to include that when converting back to a version-specific model
			matched, err := regexp.MatchString(publicAgentPoolSuffix+"$", apiProfile.Name)
			if !matched && err == nil {
				v20160930Profile := &v20160930.AgentPoolProfile{}
				convertAgentPoolProfileToV20160930(apiProfile, v20160930Profile)
				p.AgentPoolProfiles = append(p.AgentPoolProfiles, v20160930Profile)
				privIndex = i
			} else {
				pubIndex = i
			}
		}
		// Assign DNS Prefix to private agent pool from public agent pool
		p.AgentPoolProfiles[privIndex].DNSPrefix = api.AgentPoolProfiles[pubIndex].DNSPrefix
		p.AgentPoolProfiles[privIndex].FQDN = api.AgentPoolProfiles[pubIndex].FQDN
	} else {
		for _, apiProfile := range api.AgentPoolProfiles {
			v20160930Profile := &v20160930.AgentPoolProfile{}
			convertAgentPoolProfileToV20160930(apiProfile, v20160930Profile)
			p.AgentPoolProfiles = append(p.AgentPoolProfiles, v20160930Profile)
		}
	}
	if api.LinuxProfile != nil {
		p.LinuxProfile = &v20160930.LinuxProfile{}
		convertLinuxProfileToV20160930(api.LinuxProfile, p.LinuxProfile)
	}
	if api.WindowsProfile != nil {
		p.WindowsProfile = &v20160930.WindowsProfile{}
		convertWindowsProfileToV20160930(api.WindowsProfile, p.WindowsProfile)
	}
	if api.DiagnosticsProfile != nil {
		p.DiagnosticsProfile = &v20160930.DiagnosticsProfile{}
		convertDiagnosticsProfileToV20160930(api.DiagnosticsProfile, p.DiagnosticsProfile)
	}
	if api.JumpboxProfile != nil {
		p.JumpboxProfile = &v20160930.JumpboxProfile{}
		convertJumpboxProfileToV20160930(api.JumpboxProfile, p.JumpboxProfile)
	}
	if api.ServicePrincipalProfile != nil {
		p.ServicePrincipalProfile = &v20160930.ServicePrincipalProfile{}
		convertServicePrincipalProfileToV20160930(api.ServicePrincipalProfile, p.ServicePrincipalProfile)
	}
	if api.CustomProfile != nil {
		p.CustomProfile = &v20160930.CustomProfile{}
		convertCustomProfileToV20160930(api.CustomProfile, p.CustomProfile)
	}
}

func convertPropertiesToV20160330(api *Properties, p *v20160330.Properties) {
	p.ProvisioningState = v20160330.ProvisioningState(api.ProvisioningState)
	if api.OrchestratorProfile != nil {
		p.OrchestratorProfile = &v20160330.OrchestratorProfile{}
		convertOrchestratorProfileToV20160330(api.OrchestratorProfile, p.OrchestratorProfile)
	}
	if api.MasterProfile != nil {
		p.MasterProfile = &v20160330.MasterProfile{}
		convertMasterProfileToV20160330(api.MasterProfile, p.MasterProfile)
	}
	p.AgentPoolProfiles = []*v20160330.AgentPoolProfile{}
	// DCOS conversion logic
	if api.OrchestratorProfile.IsDCOS() && len(api.AgentPoolProfiles) == 2 {
		var privIndex, pubIndex int
		for i, apiProfile := range api.AgentPoolProfiles {
			// We added a pool with a "-public" suffix when converting to API model;
			// we don't want to include that when converting back to a version-specific model
			matched, err := regexp.MatchString(publicAgentPoolSuffix+"$", apiProfile.Name)
			if !matched && err == nil {
				v20160330Profile := &v20160330.AgentPoolProfile{}
				convertAgentPoolProfileToV20160330(apiProfile, v20160330Profile)
				p.AgentPoolProfiles = append(p.AgentPoolProfiles, v20160330Profile)
				privIndex = i
			} else {
				pubIndex = i
			}
		}
		// Assign DNS Prefix to private agent pool from public agent pool
		p.AgentPoolProfiles[privIndex].DNSPrefix = api.AgentPoolProfiles[pubIndex].DNSPrefix
		p.AgentPoolProfiles[privIndex].FQDN = api.AgentPoolProfiles[pubIndex].FQDN
	} else {
		for _, apiProfile := range api.AgentPoolProfiles {
			v20160330Profile := &v20160330.AgentPoolProfile{}
			convertAgentPoolProfileToV20160330(apiProfile, v20160330Profile)
			p.AgentPoolProfiles = append(p.AgentPoolProfiles, v20160330Profile)
		}
	}
	if api.LinuxProfile != nil {
		p.LinuxProfile = &v20160330.LinuxProfile{}
		convertLinuxProfileToV20160330(api.LinuxProfile, p.LinuxProfile)
	}
	if api.WindowsProfile != nil {
		p.WindowsProfile = &v20160330.WindowsProfile{}
		convertWindowsProfileToV20160330(api.WindowsProfile, p.WindowsProfile)
	}
	if api.DiagnosticsProfile != nil {
		p.DiagnosticsProfile = &v20160330.DiagnosticsProfile{}
		convertDiagnosticsProfileToV20160330(api.DiagnosticsProfile, p.DiagnosticsProfile)
	}
	if api.JumpboxProfile != nil {
		p.JumpboxProfile = &v20160330.JumpboxProfile{}
		convertJumpboxProfileToV20160330(api.JumpboxProfile, p.JumpboxProfile)
	}
}

func convertPropertiesToV20170131(api *Properties, p *v20170131.Properties) {
	p.ProvisioningState = v20170131.ProvisioningState(api.ProvisioningState)
	if api.OrchestratorProfile != nil {
		p.OrchestratorProfile = &v20170131.OrchestratorProfile{}
		convertOrchestratorProfileToV20170131(api.OrchestratorProfile, p.OrchestratorProfile)
	}
	if api.MasterProfile != nil {
		p.MasterProfile = &v20170131.MasterProfile{}
		convertMasterProfileToV20170131(api.MasterProfile, p.MasterProfile)
	}
	p.AgentPoolProfiles = []*v20170131.AgentPoolProfile{}
	// DCOS conversion logic
	if api.OrchestratorProfile.IsDCOS() && len(api.AgentPoolProfiles) == 2 {
		var privIndex, pubIndex int
		for i, apiProfile := range api.AgentPoolProfiles {
			// We added a pool with a "-public" suffix when converting to API model;
			// we don't want to include that when converting back to a version-specific model
			matched, err := regexp.MatchString(publicAgentPoolSuffix+"$", apiProfile.Name)
			if !matched && err == nil {
				v20170131Profile := &v20170131.AgentPoolProfile{}
				convertAgentPoolProfileToV20170131(apiProfile, v20170131Profile)
				p.AgentPoolProfiles = append(p.AgentPoolProfiles, v20170131Profile)
				privIndex = i
			} else {
				pubIndex = i
			}
		}
		// Assign DNS Prefix to private agent pool from public agent pool
		p.AgentPoolProfiles[privIndex].DNSPrefix = api.AgentPoolProfiles[pubIndex].DNSPrefix
		p.AgentPoolProfiles[privIndex].FQDN = api.AgentPoolProfiles[pubIndex].FQDN
	} else {
		for _, apiProfile := range api.AgentPoolProfiles {
			v20170131Profile := &v20170131.AgentPoolProfile{}
			convertAgentPoolProfileToV20170131(apiProfile, v20170131Profile)
			p.AgentPoolProfiles = append(p.AgentPoolProfiles, v20170131Profile)
		}
	}
	if api.LinuxProfile != nil {
		p.LinuxProfile = &v20170131.LinuxProfile{}
		convertLinuxProfileToV20170131(api.LinuxProfile, p.LinuxProfile)
	}
	if api.WindowsProfile != nil {
		p.WindowsProfile = &v20170131.WindowsProfile{}
		convertWindowsProfileToV20170131(api.WindowsProfile, p.WindowsProfile)
	}
	if api.DiagnosticsProfile != nil {
		p.DiagnosticsProfile = &v20170131.DiagnosticsProfile{}
		convertDiagnosticsProfileToV20170131(api.DiagnosticsProfile, p.DiagnosticsProfile)
	}
	if api.JumpboxProfile != nil {
		p.JumpboxProfile = &v20170131.JumpboxProfile{}
		convertJumpboxProfileToV20170131(api.JumpboxProfile, p.JumpboxProfile)
	}
	if api.ServicePrincipalProfile != nil {
		p.ServicePrincipalProfile = &v20170131.ServicePrincipalProfile{}
		convertServicePrincipalProfileToV20170131(api.ServicePrincipalProfile, p.ServicePrincipalProfile)
	}
	if api.CustomProfile != nil {
		p.CustomProfile = &v20170131.CustomProfile{}
		convertCustomProfileToV20170131(api.CustomProfile, p.CustomProfile)
	}
}

func convertPropertiesToV20170701(api *Properties, p *v20170701.Properties) {
	p.ProvisioningState = v20170701.ProvisioningState(api.ProvisioningState)
	if api.OrchestratorProfile != nil {
		p.OrchestratorProfile = &v20170701.OrchestratorProfile{}
		convertOrchestratorProfileToV20170701(api.OrchestratorProfile, p.OrchestratorProfile)
	}
	if api.MasterProfile != nil {
		p.MasterProfile = &v20170701.MasterProfile{}
		convertMasterProfileToV20170701(api.MasterProfile, p.MasterProfile)
	}
	p.AgentPoolProfiles = []*v20170701.AgentPoolProfile{}
	for _, apiProfile := range api.AgentPoolProfiles {
		v20170701Profile := &v20170701.AgentPoolProfile{}
		convertAgentPoolProfileToV20170701(apiProfile, v20170701Profile)
		p.AgentPoolProfiles = append(p.AgentPoolProfiles, v20170701Profile)
	}
	if api.LinuxProfile != nil {
		p.LinuxProfile = &v20170701.LinuxProfile{}
		convertLinuxProfileToV20170701(api.LinuxProfile, p.LinuxProfile)
	}
	if api.WindowsProfile != nil {
		p.WindowsProfile = &v20170701.WindowsProfile{}
		convertWindowsProfileToV20170701(api.WindowsProfile, p.WindowsProfile)
	}
	if api.ServicePrincipalProfile != nil {
		p.ServicePrincipalProfile = &v20170701.ServicePrincipalProfile{}
		convertServicePrincipalProfileToV20170701(api.ServicePrincipalProfile, p.ServicePrincipalProfile)
	}
	if api.CustomProfile != nil {
		p.CustomProfile = &v20170701.CustomProfile{}
		convertCustomProfileToV20170701(api.CustomProfile, p.CustomProfile)
	}
}

func convertPropertiesToVLabs(api *Properties, vlabsProps *vlabs.Properties) {
	vlabsProps.ProvisioningState = vlabs.ProvisioningState(api.ProvisioningState)
	if api.OrchestratorProfile != nil {
		vlabsProps.OrchestratorProfile = &vlabs.OrchestratorProfile{}
		convertOrchestratorProfileToVLabs(api.OrchestratorProfile, vlabsProps.OrchestratorProfile)
	}
	if api.MasterProfile != nil {
		vlabsProps.MasterProfile = &vlabs.MasterProfile{}
		convertMasterProfileToVLabs(api.MasterProfile, vlabsProps.MasterProfile)
	}
	vlabsProps.AgentPoolProfiles = []*vlabs.AgentPoolProfile{}
	for _, apiProfile := range api.AgentPoolProfiles {
		vlabsProfile := &vlabs.AgentPoolProfile{}
		convertAgentPoolProfileToVLabs(apiProfile, vlabsProfile)
		vlabsProps.AgentPoolProfiles = append(vlabsProps.AgentPoolProfiles, vlabsProfile)
	}
	if api.LinuxProfile != nil {
		vlabsProps.LinuxProfile = &vlabs.LinuxProfile{}
		convertLinuxProfileToVLabs(api.LinuxProfile, vlabsProps.LinuxProfile)
	}
	vlabsProps.ExtensionsProfile = []vlabs.ExtensionProfile{}
	for _, extensionProfile := range api.ExtensionsProfile {
		vlabsExtensionProfile := &vlabs.ExtensionProfile{}
		convertExtensionProfileToVLabs(&extensionProfile, vlabsExtensionProfile)
		vlabsProps.ExtensionsProfile = append(vlabsProps.ExtensionsProfile, *vlabsExtensionProfile)
	}
	if api.WindowsProfile != nil {
		vlabsProps.WindowsProfile = &vlabs.WindowsProfile{}
		convertWindowsProfileToVLabs(api.WindowsProfile, vlabsProps.WindowsProfile)
	}
	if api.ServicePrincipalProfile != nil {
		vlabsProps.ServicePrincipalProfile = &vlabs.ServicePrincipalProfile{}
		convertServicePrincipalProfileToVLabs(api.ServicePrincipalProfile, vlabsProps.ServicePrincipalProfile)
	}
	if api.CertificateProfile != nil {
		vlabsProps.CertificateProfile = &vlabs.CertificateProfile{}
		convertCertificateProfileToVLabs(api.CertificateProfile, vlabsProps.CertificateProfile)
	}
}

func convertLinuxProfileToV20160930(api *LinuxProfile, obj *v20160930.LinuxProfile) {
	obj.AdminUsername = api.AdminUsername
	obj.SSH.PublicKeys = []v20160930.PublicKey{}
	for _, d := range api.SSH.PublicKeys {
		obj.SSH.PublicKeys = append(obj.SSH.PublicKeys, v20160930.PublicKey{
			KeyData: d.KeyData,
		})
	}
}

func convertLinuxProfileToV20160330(api *LinuxProfile, obj *v20160330.LinuxProfile) {
	obj.AdminUsername = api.AdminUsername
	obj.SSH.PublicKeys = []v20160330.PublicKey{}
	for _, d := range api.SSH.PublicKeys {
		obj.SSH.PublicKeys = append(obj.SSH.PublicKeys, v20160330.PublicKey{
			KeyData: d.KeyData,
		})
	}
}

func convertLinuxProfileToV20170131(api *LinuxProfile, obj *v20170131.LinuxProfile) {
	obj.AdminUsername = api.AdminUsername
	obj.SSH.PublicKeys = []v20170131.PublicKey{}
	for _, d := range api.SSH.PublicKeys {
		obj.SSH.PublicKeys = append(obj.SSH.PublicKeys, v20170131.PublicKey{KeyData: d.KeyData})
	}
}

func convertExtensionProfileToVLabs(api *ExtensionProfile, vlabs *vlabs.ExtensionProfile) {
	vlabs.Name = api.Name
	vlabs.Version = api.Version
	vlabs.ExtensionParameters = api.ExtensionParameters
	vlabs.RootURL = api.RootURL
}

func convertExtensionToVLabs(api *Extension, vlabs *vlabs.Extension) {
	vlabs.Name = api.Name
	vlabs.SingleOrAll = api.SingleOrAll
	vlabs.Template = api.Template
}

func convertLinuxProfileToV20170701(api *LinuxProfile, obj *v20170701.LinuxProfile) {
	obj.AdminUsername = api.AdminUsername
	obj.SSH.PublicKeys = []v20170701.PublicKey{}
	for _, d := range api.SSH.PublicKeys {
		obj.SSH.PublicKeys = append(obj.SSH.PublicKeys, v20170701.PublicKey{
			KeyData: d.KeyData,
		})
	}
}

func convertLinuxProfileToVLabs(obj *LinuxProfile, vlabsProfile *vlabs.LinuxProfile) {
	vlabsProfile.AdminUsername = obj.AdminUsername
	vlabsProfile.SSH.PublicKeys = []vlabs.PublicKey{}
	for _, d := range obj.SSH.PublicKeys {
		vlabsProfile.SSH.PublicKeys = append(vlabsProfile.SSH.PublicKeys,
			vlabs.PublicKey{KeyData: d.KeyData})
	}
	vlabsProfile.Secrets = []vlabs.KeyVaultSecrets{}
	for _, s := range obj.Secrets {
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

func convertWindowsProfileToV20170131(api *WindowsProfile, v20170131 *v20170131.WindowsProfile) {
	v20170131.AdminUsername = api.AdminUsername
	v20170131.AdminPassword = api.AdminPassword
}

func convertWindowsProfileToV20170701(api *WindowsProfile, v20170701Profile *v20170701.WindowsProfile) {
	v20170701Profile.AdminUsername = api.AdminUsername
	v20170701Profile.AdminPassword = api.AdminPassword
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
	if strings.HasPrefix(api.OrchestratorType, v20160930.DCOS) {
		o.OrchestratorType = v20160930.DCOS
	} else {
		o.OrchestratorType = api.OrchestratorType
	}
}

func convertOrchestratorProfileToV20160330(api *OrchestratorProfile, o *v20160330.OrchestratorProfile) {
	if strings.HasPrefix(api.OrchestratorType, v20160330.DCOS) {
		o.OrchestratorType = v20160330.DCOS
	} else {
		o.OrchestratorType = api.OrchestratorType
	}
}

func convertOrchestratorProfileToV20170131(api *OrchestratorProfile, o *v20170131.OrchestratorProfile) {
	if strings.HasPrefix(api.OrchestratorType, v20170131.DCOS) {
		o.OrchestratorType = v20170131.DCOS
	} else {
		o.OrchestratorType = api.OrchestratorType
	}
}

func convertOrchestratorProfileToV20170701(api *OrchestratorProfile, o *v20170701.OrchestratorProfile) {
	if api.OrchestratorType == SwarmMode {
		o.OrchestratorType = v20170701.DockerCE
	} else {
		o.OrchestratorType = api.OrchestratorType
	}

	if api.OrchestratorRelease != "" {
		o.OrchestratorRelease = api.OrchestratorRelease
	}

	if api.OrchestratorVersion != "" {
		o.OrchestratorVersion = api.OrchestratorVersion
	}
}

func convertOrchestratorProfileToVLabs(api *OrchestratorProfile, o *vlabs.OrchestratorProfile) {
	o.OrchestratorType = api.OrchestratorType

	if api.OrchestratorRelease != "" {
		o.OrchestratorRelease = api.OrchestratorRelease
	}

	if api.OrchestratorVersion != "" {
		o.OrchestratorVersion = api.OrchestratorVersion
	}

	if api.KubernetesConfig != nil {
		o.KubernetesConfig = &vlabs.KubernetesConfig{}
		convertKubernetesConfigToVLabs(api.KubernetesConfig, o.KubernetesConfig)
	}
}

func convertKubernetesConfigToVLabs(api *KubernetesConfig, vlabs *vlabs.KubernetesConfig) {
	vlabs.KubernetesImageBase = api.KubernetesImageBase
	vlabs.ClusterSubnet = api.ClusterSubnet
	vlabs.NetworkPolicy = api.NetworkPolicy
	vlabs.DockerBridgeSubnet = api.DockerBridgeSubnet
	vlabs.NodeStatusUpdateFrequency = api.NodeStatusUpdateFrequency
	vlabs.CtrlMgrNodeMonitorGracePeriod = api.CtrlMgrNodeMonitorGracePeriod
	vlabs.CtrlMgrPodEvictionTimeout = api.CtrlMgrPodEvictionTimeout
	vlabs.CtrlMgrRouteReconciliationPeriod = api.CtrlMgrRouteReconciliationPeriod
	vlabs.CloudProviderBackoff = api.CloudProviderBackoff
	vlabs.CloudProviderBackoffDuration = api.CloudProviderBackoffDuration
	vlabs.CloudProviderBackoffExponent = api.CloudProviderBackoffExponent
	vlabs.CloudProviderBackoffJitter = api.CloudProviderBackoffJitter
	vlabs.CloudProviderBackoffRetries = api.CloudProviderBackoffRetries
	vlabs.CloudProviderRateLimit = api.CloudProviderRateLimit
	vlabs.CloudProviderRateLimitBucket = api.CloudProviderRateLimitBucket
	vlabs.CloudProviderRateLimitQPS = api.CloudProviderRateLimitQPS
	vlabs.UseManagedIdentity = api.UseManagedIdentity
	vlabs.CustomHyperkubeImage = api.CustomHyperkubeImage
	vlabs.UseInstanceMetadata = api.UseInstanceMetadata
	vlabs.EnableRbac = api.EnableRbac
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

func convertMasterProfileToV20170131(api *MasterProfile, v20170131 *v20170131.MasterProfile) {
	v20170131.Count = api.Count
	v20170131.DNSPrefix = api.DNSPrefix
	v20170131.FQDN = api.FQDN
	v20170131.SetSubnet(api.Subnet)
}

func convertMasterProfileToV20170701(api *MasterProfile, v20170701Profile *v20170701.MasterProfile) {
	v20170701Profile.Count = api.Count
	v20170701Profile.DNSPrefix = api.DNSPrefix
	v20170701Profile.FQDN = api.FQDN
	v20170701Profile.SetSubnet(api.Subnet)
	v20170701Profile.VMSize = api.VMSize
	v20170701Profile.OSDiskSizeGB = api.OSDiskSizeGB
	v20170701Profile.VnetSubnetID = api.VnetSubnetID
	v20170701Profile.FirstConsecutiveStaticIP = api.FirstConsecutiveStaticIP
	v20170701Profile.StorageProfile = api.StorageProfile
}

func convertMasterProfileToVLabs(api *MasterProfile, vlabsProfile *vlabs.MasterProfile) {
	vlabsProfile.Count = api.Count
	vlabsProfile.DNSPrefix = api.DNSPrefix
	vlabsProfile.VMSize = api.VMSize
	vlabsProfile.OSDiskSizeGB = api.OSDiskSizeGB
	vlabsProfile.VnetSubnetID = api.VnetSubnetID
	vlabsProfile.FirstConsecutiveStaticIP = api.FirstConsecutiveStaticIP
	vlabsProfile.SetSubnet(api.Subnet)
	vlabsProfile.FQDN = api.FQDN
	vlabsProfile.StorageProfile = api.StorageProfile
	vlabsProfile.Extensions = []vlabs.Extension{}
	for _, extension := range api.Extensions {
		vlabsExtension := &vlabs.Extension{}
		convertExtensionToVLabs(&extension, vlabsExtension)
		vlabsProfile.Extensions = append(vlabsProfile.Extensions, *vlabsExtension)
	}
}

func convertKeyVaultSecretsToVlabs(api *KeyVaultSecrets, vlabsSecrets *vlabs.KeyVaultSecrets) {
	vlabsSecrets.SourceVault = &vlabs.KeyVaultID{ID: api.SourceVault.ID}
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

func convertAgentPoolProfileToV20170131(api *AgentPoolProfile, p *v20170131.AgentPoolProfile) {
	p.Name = api.Name
	p.Count = api.Count
	p.VMSize = api.VMSize
	p.DNSPrefix = api.DNSPrefix
	p.FQDN = api.FQDN
	p.OSType = v20170131.OSType(api.OSType)
	p.SetSubnet(api.Subnet)
}

func convertAgentPoolProfileToV20170701(api *AgentPoolProfile, p *v20170701.AgentPoolProfile) {
	p.Name = api.Name
	p.Count = api.Count
	p.VMSize = api.VMSize
	p.DNSPrefix = api.DNSPrefix
	p.FQDN = api.FQDN
	p.OSType = v20170701.OSType(api.OSType)
	p.SetSubnet(api.Subnet)
	p.OSDiskSizeGB = api.OSDiskSizeGB
	p.Ports = []int{}
	p.Ports = append(p.Ports, api.Ports...)
	p.StorageProfile = api.StorageProfile
	p.VnetSubnetID = api.VnetSubnetID
}

func convertAgentPoolProfileToVLabs(api *AgentPoolProfile, p *vlabs.AgentPoolProfile) {
	p.Name = api.Name
	p.Count = api.Count
	p.VMSize = api.VMSize
	p.OSDiskSizeGB = api.OSDiskSizeGB
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
	p.CustomNodeLabels = map[string]string{}
	for k, v := range api.CustomNodeLabels {
		p.CustomNodeLabels[k] = v
	}

	p.Extensions = []vlabs.Extension{}
	for _, extension := range api.Extensions {
		vlabsExtension := &vlabs.Extension{}
		convertExtensionToVLabs(&extension, vlabsExtension)
		p.Extensions = append(p.Extensions, *vlabsExtension)
	}
}

func convertDiagnosticsProfileToV20160930(api *DiagnosticsProfile, dp *v20160930.DiagnosticsProfile) {
	if api.VMDiagnostics != nil {
		dp.VMDiagnostics = &v20160930.VMDiagnostics{}
		convertVMDiagnosticsToV20160930(api.VMDiagnostics, dp.VMDiagnostics)
	}
}

func convertVMDiagnosticsToV20160930(api *VMDiagnostics, v20160930 *v20160930.VMDiagnostics) {
	v20160930.Enabled = api.Enabled
	v20160930.StorageURL = api.StorageURL
}

func convertDiagnosticsProfileToV20160330(api *DiagnosticsProfile, dp *v20160330.DiagnosticsProfile) {
	if api.VMDiagnostics != nil {
		dp.VMDiagnostics = &v20160330.VMDiagnostics{}
		convertVMDiagnosticsToV20160330(api.VMDiagnostics, dp.VMDiagnostics)
	}
}

func convertVMDiagnosticsToV20160330(api *VMDiagnostics, v20160330 *v20160330.VMDiagnostics) {
	v20160330.Enabled = api.Enabled
	v20160330.StorageURL = api.StorageURL
}

func convertDiagnosticsProfileToV20170131(api *DiagnosticsProfile, dp *v20170131.DiagnosticsProfile) {
	if api.VMDiagnostics != nil {
		dp.VMDiagnostics = &v20170131.VMDiagnostics{}
		convertVMDiagnosticsToV20170131(api.VMDiagnostics, dp.VMDiagnostics)
	}
}

func convertVMDiagnosticsToV20170131(api *VMDiagnostics, v20170131 *v20170131.VMDiagnostics) {
	v20170131.Enabled = api.Enabled
	v20170131.StorageURL = api.StorageURL
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

func convertJumpboxProfileToV20170131(api *JumpboxProfile, jb *v20170131.JumpboxProfile) {
	jb.OSType = v20170131.OSType(api.OSType)
	jb.DNSPrefix = api.DNSPrefix
	jb.FQDN = api.FQDN
}

func convertServicePrincipalProfileToV20160930(api *ServicePrincipalProfile, v20160930 *v20160930.ServicePrincipalProfile) {
	v20160930.ClientID = api.ClientID
	v20160930.Secret = api.Secret
}

func convertServicePrincipalProfileToV20170131(api *ServicePrincipalProfile, v20170131 *v20170131.ServicePrincipalProfile) {
	v20170131.ClientID = api.ClientID
	v20170131.Secret = api.Secret
}

func convertCustomProfileToV20160930(api *CustomProfile, v20160930 *v20160930.CustomProfile) {
	v20160930.Orchestrator = api.Orchestrator
}

func convertCustomProfileToV20170131(api *CustomProfile, v20170131 *v20170131.CustomProfile) {
	v20170131.Orchestrator = api.Orchestrator
}

func convertCustomProfileToV20170701(api *CustomProfile, v20170701 *v20170701.CustomProfile) {
	v20170701.Orchestrator = api.Orchestrator
}

func convertServicePrincipalProfileToV20170701(api *ServicePrincipalProfile, v20170701 *v20170701.ServicePrincipalProfile) {
	v20170701.ClientID = api.ClientID
	v20170701.Secret = api.Secret
	v20170701.KeyvaultSecretRef = api.KeyvaultSecretRef
}

func convertServicePrincipalProfileToVLabs(api *ServicePrincipalProfile, vlabs *vlabs.ServicePrincipalProfile) {
	vlabs.ClientID = api.ClientID
	vlabs.Secret = api.Secret
	vlabs.KeyvaultSecretRef = api.KeyvaultSecretRef
}

func convertCertificateProfileToVLabs(api *CertificateProfile, vlabs *vlabs.CertificateProfile) {
	vlabs.CaCertificate = api.CaCertificate
	vlabs.CaPrivateKey = api.CaPrivateKey
	vlabs.APIServerCertificate = api.APIServerCertificate
	vlabs.APIServerPrivateKey = api.APIServerPrivateKey
	vlabs.ClientCertificate = api.ClientCertificate
	vlabs.ClientPrivateKey = api.ClientPrivateKey
	vlabs.KubeConfigCertificate = api.KubeConfigCertificate
	vlabs.KubeConfigPrivateKey = api.KubeConfigPrivateKey
}
