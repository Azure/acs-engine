package agentpool

import (
	"fmt"
	"github.com/Azure/acs-engine/pkg/acsengine"
	"github.com/Azure/acs-engine/pkg/api/kubernetesagentpool"
)

// getParameters is an unexported function that will create the parameters for the azuredeploy.parameters.json file
// This was intentionally pulled out of acsenging so we can have a unique and decoupled grouping of parameters for agent pools only.
func getParameters(agentPool *kubernetesagentpool.AgentPool) (map[string]interface{}, error) {
	properties := agentPool.Properties
	location := agentPool.Location
	parametersMap := map[string]interface{}{}
	acsengine.AddValue(parametersMap, "location", location)
	//acsengine.AddValue(parametersMap, "targetEnvironment", acsengine.GetCloudTargetEnv(location))
	acsengine.AddValue(parametersMap, "linuxAdminUsername", properties.LinuxProfile.AdminUsername)
	acsengine.AddValue(parametersMap, "sshRSAPublicKey", properties.LinuxProfile.SSH.PublicKeys[0].KeyData)
	for i, s := range properties.LinuxProfile.Secrets {
		acsengine.AddValue(parametersMap, fmt.Sprintf("linuxKeyVaultID%d", i), s.SourceVault.ID)
		for j, c := range s.VaultCertificates {
			acsengine.AddValue(parametersMap, fmt.Sprintf("linuxKeyVaultID%dCertificateURL%d", i, j), c.CertificateURL)
		}
	}
	//var KubernetesVersion api.OrchestratorVersion
	//KubernetesVersion = api.OrchestratorVersion(properties.KubernetesVersion)
	cloudSpecConfig := acsengine.GetCloudSpecConfig(location)

	// Agentpool parameters
	for _, agentProfile := range properties.AgentPoolProfiles {
		acsengine.AddValue(parametersMap, fmt.Sprintf("%sCount", agentProfile.Name), agentProfile.Count)
		acsengine.AddValue(parametersMap, fmt.Sprintf("%sVMSize", agentProfile.Name), agentProfile.VMSize)
	}

	// Jumpbox
	acsengine.AddValue(parametersMap, "jumpboxVmSize", properties.JumpBoxProfile.VMSize)
	acsengine.AddValue(parametersMap, "jumpboxCount", properties.JumpBoxProfile.Count)
	acsengine.AddValue(parametersMap, "jumpboxEndpointDNSNamePrefix", properties.DnsPrefix)

	// Certificate information
	acsengine.AddSecret(parametersMap, "apiServerCertificate", properties.CertificateProfile.APIServerCertificate, true)
	acsengine.AddSecret(parametersMap, "apiServerPrivateKey", properties.CertificateProfile.APIServerPrivateKey, true)
	acsengine.AddSecret(parametersMap, "caCertificate", properties.CertificateProfile.CaCertificate, true)
	acsengine.AddSecret(parametersMap, "caPrivateKey", properties.CertificateProfile.CaPrivateKey, true)
	acsengine.AddSecret(parametersMap, "clientCertificate", properties.CertificateProfile.ClientCertificate, true)
	acsengine.AddSecret(parametersMap, "clientPrivateKey", properties.CertificateProfile.ClientPrivateKey, true)
	acsengine.AddSecret(parametersMap, "kubeConfigCertificate", properties.CertificateProfile.KubeConfigCertificate, true)
	acsengine.AddSecret(parametersMap, "kubeConfigPrivateKey", properties.CertificateProfile.KubeConfigPrivateKey, true)

	acsengine.AddValue(parametersMap, "dockerEngineDownloadRepo", cloudSpecConfig.DockerSpecConfig.DockerEngineRepo)
	////acsengine.AddValue(parametersMap, "kubernetesHyperkubeSpec", properties.OrchestratorProfile.KubernetesConfig.KubernetesImageBase+acsengine.KubeImages[KubernetesVersion]["hyperkube"])
	//acsengine.AddValue(parametersMap, "kubernetesAddonManagerSpec", cloudSpecConfig.KubernetesSpecConfig.KubernetesImageBase+acsengine.KubeImages[KubernetesVersion]["addonmanager"])
	//acsengine.AddValue(parametersMap, "kubernetesAddonResizerSpec", cloudSpecConfig.KubernetesSpecConfig.KubernetesImageBase+acsengine.KubeImages[KubernetesVersion]["addonresizer"])
	//acsengine.AddValue(parametersMap, "kubernetesDashboardSpec", cloudSpecConfig.KubernetesSpecConfig.KubernetesImageBase+acsengine.KubeImages[KubernetesVersion]["dashboard"])
	//acsengine.AddValue(parametersMap, "kubernetesDNSMasqSpec", cloudSpecConfig.KubernetesSpecConfig.KubernetesImageBase+acsengine.KubeImages[KubernetesVersion]["dnsmasq"])
	//acsengine.AddValue(parametersMap, "kubernetesExecHealthzSpec", cloudSpecConfig.KubernetesSpecConfig.KubernetesImageBase+acsengine.KubeImages[KubernetesVersion]["exechealthz"])
	//acsengine.AddValue(parametersMap, "kubernetesHeapsterSpec", cloudSpecConfig.KubernetesSpecConfig.KubernetesImageBase+acsengine.KubeImages[KubernetesVersion]["heapster"])
	//acsengine.AddValue(parametersMap, "kubernetesKubeDNSSpec", cloudSpecConfig.KubernetesSpecConfig.KubernetesImageBase+acsengine.KubeImages[KubernetesVersion]["dns"])
	//acsengine.AddValue(parametersMap, "kubernetesPodInfraContainerSpec", cloudSpecConfig.KubernetesSpecConfig.KubernetesImageBase+acsengine.KubeImages[KubernetesVersion]["pause"])
	//acsengine.AddValue(parametersMap, "kubernetesNodeStatusUpdateFrequency", acsengine.KubeImages[KubernetesVersion]["nodestatusfreq"])
	//acsengine.AddValue(parametersMap, "kubernetesCtrlMgrNodeMonitorGracePeriod", acsengine.KubeImages[KubernetesVersion]["nodegraceperiod"])
	//acsengine.AddValue(parametersMap, "kubernetesCtrlMgrPodEvictionTimeout", acsengine.KubeImages[KubernetesVersion]["podeviction"])
	//acsengine.AddValue(parametersMap, "kubernetesCtrlMgrRouteReconciliationPeriod", acsengine.KubeImages[KubernetesVersion]["routeperiod"])
	//acsengine.AddValue(parametersMap, "cloudProviderBackoff", acsengine.KubeImages[KubernetesVersion]["backoff"])
	//acsengine.AddValue(parametersMap, "cloudProviderBackoffRetries", acsengine.KubeImages[KubernetesVersion]["backoffretries"])
	//acsengine.AddValue(parametersMap, "cloudProviderBackoffExponent", acsengine.KubeImages[KubernetesVersion]["backoffexponent"])
	//acsengine.AddValue(parametersMap, "cloudProviderBackoffDuration", acsengine.KubeImages[KubernetesVersion]["backoffduration"])
	//acsengine.AddValue(parametersMap, "cloudProviderBackoffJitter", acsengine.KubeImages[KubernetesVersion]["backoffjitter"])
	//acsengine.AddValue(parametersMap, "cloudProviderRatelimit", acsengine.KubeImages[KubernetesVersion]["ratelimit"])
	//acsengine.AddValue(parametersMap, "cloudProviderRatelimitQPS", acsengine.KubeImages[KubernetesVersion]["ratelimitqps"])
	//acsengine.AddValue(parametersMap, "cloudProviderRatelimitBucket", acsengine.KubeImages[KubernetesVersion]["ratelimitbucket"])
	acsengine.AddValue(parametersMap, "jumpboxSubnet", properties.NetworkProfile.ServiceCIDR)
	////acsengine.AddValue(parametersMap, "dockerBridgeCidr", properties.OrchestratorProfile.KubernetesConfig.DockerBridgeSubnet)
	////acsengine.AddValue(parametersMap, "networkPolicy", properties.OrchestratorProfile.KubernetesConfig.NetworkPolicy)
	acsengine.AddValue(parametersMap, "servicePrincipalClientId", properties.ServicePrincipalProfile.ClientID)
	acsengine.AddSecret(parametersMap, "servicePrincipalClientSecret", properties.ServicePrincipalProfile.Secret, false)

	return parametersMap, nil
}
