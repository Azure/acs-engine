package api

import (
	"github.com/Azure/acs-labs/acstgen/pkg/api/v20160330"
	"github.com/Azure/acs-labs/acstgen/pkg/api/vlabs"
)

///////////////////////////////////////////////////////////
// The converter exposes functions to convert the 3 top
// level resources:
// 1. Subscription
// 2. ResourcePurchasePlan
// 3. ContainerService
//
// All other functions are internal helper functions used
// for converting.
///////////////////////////////////////////////////////////

// ConvertV20160330Subscription converts a v20160330 Subscription to an unversioned Subscription
func ConvertV20160330Subscription(v20160330 *v20160330.Subscription, api *Subscription) {

}

// ConvertVLabsSubscription converts a vlabs Subscription to an unversioned Subscription
func ConvertVLabsSubscription(vlabs *vlabs.Subscription, api *Subscription) {

}

// ConvertV20160330ResourcePurchasePlan converts a v20160330 ResourcePurchasePlan to an unversioned ResourcePurchasePlan
func ConvertV20160330ResourcePurchasePlan(v20160330 *v20160330.ResourcePurchasePlan, api *ResourcePurchasePlan) {

}

// ConvertVLabsResourcePurchasePlan converts a vlabs ResourcePurchasePlan to an unversioned ResourcePurchasePlan
func ConvertVLabsResourcePurchasePlan(vlabs *vlabs.ResourcePurchasePlan, api *ResourcePurchasePlan) {

}

// ConvertV20160330ContainerService converts a v20160330 ContainerService to an unversioned ContainerService
func ConvertV20160330ContainerService(v20160330 *v20160330.ContainerService, api *ContainerService) {

}

// ConvertVLabsContainerService converts a vlabs ContainerService to an unversioned ContainerService
func ConvertVLabsContainerService(v20160330 *vlabs.ContainerService, api *ContainerService) {

}

func convertV20160330Properties(v20160330 *v20160330.Properties, api *Properties) {

}

func convertVLabsProperties(vlabs *vlabs.Properties, api *Properties) {

}

func convertV20160330LinuxProfile(v20160330 *v20160330.LinuxProfile, api *LinuxProfile) {

}

func convertVLabsLinuxProfile(vlabs *vlabs.LinuxProfile, api *LinuxProfile) {

}

func convertV20160330WindowsProfile(v20160330 *v20160330.LinuxProfile, api *LinuxProfile) {

}

func convertVLabsWindowsProfile(vlabs *vlabs.LinuxProfile, api *LinuxProfile) {

}

func convertV20160330OrchestratorProfile(v20160330 *v20160330.OrchestratorProfile, api *OrchestratorProfile) {

}

func convertVLabsOrchestratorProfile(vlabs *vlabs.OrchestratorProfile, api *OrchestratorProfile) {

}

func convertV20160330MasterProfile(v20160330 *v20160330.MasterProfile, api *MasterProfile) {

}

func convertVLabsMasterProfile(vlabs *vlabs.MasterProfile, api *MasterProfile) {

}

func convertV20160330AgentPoolProfile(v20160330 *v20160330.AgentPoolProfile, api *AgentPoolProfile) {

}

func convertVLabsAgentPoolProfile(vlabs *vlabs.AgentPoolProfile, api *AgentPoolProfile) {

}

func convertV20160330DiagnosticsProfile(v20160330 *v20160330.DiagnosticsProfile, api *DiagnosticsProfile) {

}

func convertV20160330VMDiagnostics(v20160330 *v20160330.VMDiagnostics, api *VMDiagnostics) {

}

func convertVLabsServicePrincipalProfile(vlabs *vlabs.ServicePrincipalProfile, api *ServicePrincipalProfile) {

}

func convertVLabsCertificateProfile(vlabs *vlabs.CertificateProfile, api *CertificateProfile) {

}
