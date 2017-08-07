package agentpool

import (
	"fmt"

	"github.com/Azure/acs-engine/pkg/acsengine"
	"github.com/Azure/acs-engine/pkg/api/kubernetesagentpool"
)

const (

	// KubernetesImagebase is the base image path to pull the Kubernetes containers from
	KubernetesImagebase = "gcrio.azureedge.net/google_containers/"
)

// getParameters is an unexported function that will create the parameters for the azuredeploy.parameters.json file
// This was intentionally pulled out of acsengine so we can have a unique and decoupled grouping of parameters for agent pools only.
func getParameters(agentPool *kubernetesagentpool.AgentPool) (map[string]interface{}, error) {
	properties := agentPool.Properties
	location := agentPool.Location
	parametersMap := map[string]interface{}{}
	acsengine.AddValue(parametersMap, "location", location)
	acsengine.AddValue(parametersMap, "linuxAdminUsername", properties.LinuxProfile.AdminUsername)
	acsengine.AddValue(parametersMap, "sshRSAPublicKey", properties.LinuxProfile.SSH.PublicKeys[0].KeyData)
	for i, s := range properties.LinuxProfile.Secrets {
		acsengine.AddValue(parametersMap, fmt.Sprintf("linuxKeyVaultID%d", i), s.SourceVault.ID)
		for j, c := range s.VaultCertificates {
			acsengine.AddValue(parametersMap, fmt.Sprintf("linuxKeyVaultID%dCertificateURL%d", i, j), c.CertificateURL)
		}
	}
	KubernetesVersion := agentPool.Properties.KubernetesVersion

	cloudSpecConfig := acsengine.GetCloudSpecConfig(location)

	// Agentpool parameters
	for _, agentProfile := range properties.AgentPoolProfiles {
		acsengine.AddValue(parametersMap, fmt.Sprintf("%sCount", agentProfile.Name), agentProfile.Count)
		acsengine.AddValue(parametersMap, fmt.Sprintf("%sVMSize", agentProfile.Name), agentProfile.VMSize)
	}

	// Jumpbox
	acsengine.AddValue(parametersMap, "jumpboxVmSize", properties.JumpBoxProfile.VMSize)
	acsengine.AddValue(parametersMap, "jumpboxCount", properties.JumpBoxProfile.Count)
	acsengine.AddValue(parametersMap, "jumpboxEndpointDNSNamePrefix", properties.DNSPrefix)
	acsengine.AddValue(parametersMap, "jumpboxInternalAddress", properties.JumpBoxProfile.InternalAddress)

	// Certificate information
	acsengine.AddSecret(parametersMap, "apiServerCertificate", properties.CertificateProfile.APIServerCertificate, true)
	acsengine.AddSecret(parametersMap, "apiServerPrivateKey", properties.CertificateProfile.APIServerPrivateKey, true)
	acsengine.AddSecret(parametersMap, "caCertificate", properties.CertificateProfile.CaCertificate, true)
	acsengine.AddSecret(parametersMap, "caPrivateKey", properties.CertificateProfile.CaPrivateKey, true)
	acsengine.AddSecret(parametersMap, "clientCertificate", properties.CertificateProfile.ClientCertificate, true)
	acsengine.AddSecret(parametersMap, "clientPrivateKey", properties.CertificateProfile.ClientPrivateKey, true)
	acsengine.AddSecret(parametersMap, "kubeConfigCertificate", properties.CertificateProfile.KubeConfigCertificate, true)
	acsengine.AddSecret(parametersMap, "kubeConfigPrivateKey", properties.CertificateProfile.KubeConfigPrivateKey, true)

	// Kubernetes
	acsengine.AddValue(parametersMap, "dockerEngineDownloadRepo", cloudSpecConfig.DockerSpecConfig.DockerEngineRepo)
	acsengine.AddValue(parametersMap, "kubernetesHyperkubeSpec", KubernetesImagebase+acsengine.KubeImages[KubernetesVersion]["hyperkube"])
	acsengine.AddValue(parametersMap, "kubernetesPodInfraContainerSpec", cloudSpecConfig.KubernetesSpecConfig.KubernetesImageBase+acsengine.KubeImages[KubernetesVersion]["pause"])
	acsengine.AddValue(parametersMap, "kubernetesNodeStatusUpdateFrequency", acsengine.KubeImages[KubernetesVersion]["nodestatusfreq"])
	acsengine.AddValue(parametersMap, "kubernetesCtrlMgrNodeMonitorGracePeriod", acsengine.KubeImages[KubernetesVersion]["nodegraceperiod"])
	acsengine.AddValue(parametersMap, "kubernetesCtrlMgrPodEvictionTimeout", acsengine.KubeImages[KubernetesVersion]["podeviction"])
	acsengine.AddValue(parametersMap, "kubernetesCtrlMgrRouteReconciliationPeriod", acsengine.KubeImages[KubernetesVersion]["routeperiod"])
	acsengine.AddValue(parametersMap, "jumpboxSubnet", properties.NetworkProfile.AgentCIDR)
	acsengine.AddValue(parametersMap, "servicePrincipalClientId", properties.ServicePrincipalProfile.ClientID)
	acsengine.AddSecret(parametersMap, "servicePrincipalClientSecret", properties.ServicePrincipalProfile.Secret, false)
	acsengine.AddValue(parametersMap, "kubernetesApiServer", properties.KubernetesEndpoint)

	return parametersMap, nil
}
