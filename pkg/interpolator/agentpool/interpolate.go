package agentpool

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Azure/acs-engine/pkg/acsengine"
	"github.com/Azure/acs-engine/pkg/api"
	"github.com/Azure/acs-engine/pkg/interpolator"
	"github.com/prometheus/common/log"
	txttemplate "text/template"
)

type Interpolator struct {
	containerService *api.ContainerService
	interpolated     bool
	template         []byte
	parameters       []byte
}

func NewAgentPoolInterpolator(containerService *api.ContainerService) interpolator.Interpolator {
	return &Interpolator{
		containerService: containerService,
	}
}

func (i *Interpolator) Interpolate() error {

	// Init template
	templ := txttemplate.New("agentpool").Funcs(acsengine.GetTemplateFuncMap(i.containerService))

	// Load files
	files, err := acsengine.AssetDir("kubernetes/agentpool")
	if err != nil {
		return fmt.Errorf("Unable to parse asset dir [kubernetes/agentpool]: %v", err)
	}

	// Parse files
	for _, file := range files {
		bytes, err := acsengine.Asset(file)
		if err != nil {
			return fmt.Errorf("Error reading file %s, Error: %s", file, err.Error())
		}
		if _, err = templ.New(file).Parse(string(bytes)); err != nil {
			return fmt.Errorf("Unable to parse template: %v", err)
		}
	}

	// Watch for panics
	defer func() {
		if r := recover(); r != nil {
			//err = fmt.Errorf("%v", r)
			log.Fatal(err)
			// invalidate the template and the parameters
			//templateRaw = ""
			//parametersRaw = ""
		}
	}()

	var b bytes.Buffer
	if err = templ.ExecuteTemplate(&b, "kubernetesbase.t", i.containerService.Properties); err != nil {
		return fmt.Errorf("Unable to execute template: %v", err)
	}

	var parametersMap map[string]interface{}
	parametersMap, err = getParameters(i.containerService)
	if err != nil {
		return fmt.Errorf("Unable to get parameteres: %v", err)
	}
	var parameterBytes []byte
	parameterBytes, err = json.Marshal(parametersMap)
	if err != nil {
		return fmt.Errorf("Unable to marshal parameters map: %v", err)
	}

	// Cache the template
	i.template = b.Bytes()
	i.parameters = parameterBytes
	i.interpolated = true
	return nil
}

func getParameters(containerService *api.ContainerService) (map[string]interface{}, error) {
	properties := containerService.Properties
	location := containerService.Location
	parametersMap := map[string]interface{}{}

	// Master Parameters
	acsengine.AddValue(parametersMap, "location", location)
	acsengine.AddValue(parametersMap, "targetEnvironment", acsengine.GetCloudTargetEnv(location))
	acsengine.AddValue(parametersMap, "linuxAdminUsername", properties.LinuxProfile.AdminUsername)
	acsengine.AddValue(parametersMap, "masterEndpointDNSNamePrefix", properties.MasterProfile.DNSPrefix)
	if properties.MasterProfile.IsCustomVNET() {
		acsengine.AddValue(parametersMap, "masterVnetSubnetID", properties.MasterProfile.VnetSubnetID)
	} else {
		acsengine.AddValue(parametersMap, "masterSubnet", properties.MasterProfile.Subnet)
	}
	acsengine.AddValue(parametersMap, "firstConsecutiveStaticIP", properties.MasterProfile.FirstConsecutiveStaticIP)
	acsengine.AddValue(parametersMap, "masterVMSize", properties.MasterProfile.VMSize)
	acsengine.AddValue(parametersMap, "masterCount", properties.MasterProfile.Count)
	acsengine.AddValue(parametersMap, "sshRSAPublicKey", properties.LinuxProfile.SSH.PublicKeys[0].KeyData)
	for i, s := range properties.LinuxProfile.Secrets {
		acsengine.AddValue(parametersMap, fmt.Sprintf("linuxKeyVaultID%d", i), s.SourceVault.ID)
		for j, c := range s.VaultCertificates {
			acsengine.AddValue(parametersMap, fmt.Sprintf("linuxKeyVaultID%dCertificateURL%d", i, j), c.CertificateURL)
		}
	}
	KubernetesVersion := properties.OrchestratorProfile.OrchestratorVersion
	cloudSpecConfig := acsengine.GetCloudSpecConfig(location)
	acsengine.AddSecret(parametersMap, "apiServerCertificate", properties.CertificateProfile.APIServerCertificate, true)
	acsengine.AddSecret(parametersMap, "apiServerPrivateKey", properties.CertificateProfile.APIServerPrivateKey, true)
	acsengine.AddSecret(parametersMap, "caCertificate", properties.CertificateProfile.CaCertificate, true)
	acsengine.AddSecret(parametersMap, "caPrivateKey", properties.CertificateProfile.CaPrivateKey, true)
	acsengine.AddSecret(parametersMap, "clientCertificate", properties.CertificateProfile.ClientCertificate, true)
	acsengine.AddSecret(parametersMap, "clientPrivateKey", properties.CertificateProfile.ClientPrivateKey, true)
	acsengine.AddSecret(parametersMap, "kubeConfigCertificate", properties.CertificateProfile.KubeConfigCertificate, true)
	acsengine.AddSecret(parametersMap, "kubeConfigPrivateKey", properties.CertificateProfile.KubeConfigPrivateKey, true)
	acsengine.AddValue(parametersMap, "dockerEngineDownloadRepo", cloudSpecConfig.DockerSpecConfig.DockerEngineRepo)
	acsengine.AddValue(parametersMap, "kubernetesHyperkubeSpec", properties.OrchestratorProfile.KubernetesConfig.KubernetesImageBase+acsengine.KubeImages[KubernetesVersion]["hyperkube"])
	acsengine.AddValue(parametersMap, "kubernetesAddonManagerSpec", cloudSpecConfig.KubernetesSpecConfig.KubernetesImageBase+acsengine.KubeImages[KubernetesVersion]["addonmanager"])
	acsengine.AddValue(parametersMap, "kubernetesAddonResizerSpec", cloudSpecConfig.KubernetesSpecConfig.KubernetesImageBase+acsengine.KubeImages[KubernetesVersion]["addonresizer"])
	acsengine.AddValue(parametersMap, "kubernetesDashboardSpec", cloudSpecConfig.KubernetesSpecConfig.KubernetesImageBase+acsengine.KubeImages[KubernetesVersion]["dashboard"])
	acsengine.AddValue(parametersMap, "kubernetesDNSMasqSpec", cloudSpecConfig.KubernetesSpecConfig.KubernetesImageBase+acsengine.KubeImages[KubernetesVersion]["dnsmasq"])
	acsengine.AddValue(parametersMap, "kubernetesExecHealthzSpec", cloudSpecConfig.KubernetesSpecConfig.KubernetesImageBase+acsengine.KubeImages[KubernetesVersion]["exechealthz"])
	acsengine.AddValue(parametersMap, "kubernetesHeapsterSpec", cloudSpecConfig.KubernetesSpecConfig.KubernetesImageBase+acsengine.KubeImages[KubernetesVersion]["heapster"])
	acsengine.AddValue(parametersMap, "kubernetesKubeDNSSpec", cloudSpecConfig.KubernetesSpecConfig.KubernetesImageBase+acsengine.KubeImages[KubernetesVersion]["dns"])
	acsengine.AddValue(parametersMap, "kubernetesPodInfraContainerSpec", cloudSpecConfig.KubernetesSpecConfig.KubernetesImageBase+acsengine.KubeImages[KubernetesVersion]["pause"])
	acsengine.AddValue(parametersMap, "kubernetesNodeStatusUpdateFrequency", acsengine.KubeImages[KubernetesVersion]["nodestatusfreq"])
	acsengine.AddValue(parametersMap, "kubernetesCtrlMgrNodeMonitorGracePeriod", acsengine.KubeImages[KubernetesVersion]["nodegraceperiod"])
	acsengine.AddValue(parametersMap, "kubernetesCtrlMgrPodEvictionTimeout", acsengine.KubeImages[KubernetesVersion]["podeviction"])
	acsengine.AddValue(parametersMap, "kubernetesCtrlMgrRouteReconciliationPeriod", acsengine.KubeImages[KubernetesVersion]["routeperiod"])
	acsengine.AddValue(parametersMap, "cloudProviderBackoff", acsengine.KubeImages[KubernetesVersion]["backoff"])
	acsengine.AddValue(parametersMap, "cloudProviderBackoffRetries", acsengine.KubeImages[KubernetesVersion]["backoffretries"])
	acsengine.AddValue(parametersMap, "cloudProviderBackoffExponent", acsengine.KubeImages[KubernetesVersion]["backoffexponent"])
	acsengine.AddValue(parametersMap, "cloudProviderBackoffDuration", acsengine.KubeImages[KubernetesVersion]["backoffduration"])
	acsengine.AddValue(parametersMap, "cloudProviderBackoffJitter", acsengine.KubeImages[KubernetesVersion]["backoffjitter"])
	acsengine.AddValue(parametersMap, "cloudProviderRatelimit", acsengine.KubeImages[KubernetesVersion]["ratelimit"])
	acsengine.AddValue(parametersMap, "cloudProviderRatelimitQPS", acsengine.KubeImages[KubernetesVersion]["ratelimitqps"])
	acsengine.AddValue(parametersMap, "cloudProviderRatelimitBucket", acsengine.KubeImages[KubernetesVersion]["ratelimitbucket"])
	acsengine.AddValue(parametersMap, "kubeClusterCidr", properties.OrchestratorProfile.KubernetesConfig.ClusterSubnet)
	acsengine.AddValue(parametersMap, "dockerBridgeCidr", properties.OrchestratorProfile.KubernetesConfig.DockerBridgeSubnet)
	acsengine.AddValue(parametersMap, "networkPolicy", properties.OrchestratorProfile.KubernetesConfig.NetworkPolicy)
	acsengine.AddValue(parametersMap, "servicePrincipalClientId", properties.ServicePrincipalProfile.ClientID)
	acsengine.AddSecret(parametersMap, "servicePrincipalClientSecret", properties.ServicePrincipalProfile.Secret, false)

	// Agent parameters
	for _, agentProfile := range properties.AgentPoolProfiles {
		acsengine.AddValue(parametersMap, fmt.Sprintf("%sCount", agentProfile.Name), agentProfile.Count)
		acsengine.AddValue(parametersMap, fmt.Sprintf("%sVMSize", agentProfile.Name), agentProfile.VMSize)
		if agentProfile.IsCustomVNET() {
			acsengine.AddValue(parametersMap, fmt.Sprintf("%sVnetSubnetID", agentProfile.Name), agentProfile.VnetSubnetID)
		} else {
			acsengine.AddValue(parametersMap, fmt.Sprintf("%sSubnet", agentProfile.Name), agentProfile.Subnet)
		}
		if len(agentProfile.Ports) > 0 {
			acsengine.AddValue(parametersMap, fmt.Sprintf("%sEndpointDNSNamePrefix", agentProfile.Name), agentProfile.DNSPrefix)
		}
	}

	// Windows parameters
	if properties.HasWindows() {
		acsengine.AddValue(parametersMap, "windowsAdminUsername", properties.WindowsProfile.AdminUsername)
		acsengine.AddSecret(parametersMap, "windowsAdminPassword", properties.WindowsProfile.AdminPassword, false)
		if properties.OrchestratorProfile.OrchestratorType == api.Kubernetes {
			KubernetesVersion := properties.OrchestratorProfile.OrchestratorVersion
			acsengine.AddValue(parametersMap, "kubeBinariesSASURL", cloudSpecConfig.KubernetesSpecConfig.KubeBinariesSASURLBase+acsengine.KubeImages[KubernetesVersion]["windowszip"])
			acsengine.AddValue(parametersMap, "kubeBinariesVersion", KubernetesVersion)
		}
		for i, s := range properties.WindowsProfile.Secrets {
			acsengine.AddValue(parametersMap, fmt.Sprintf("windowsKeyVaultID%d", i), s.SourceVault.ID)
			for j, c := range s.VaultCertificates {
				acsengine.AddValue(parametersMap, fmt.Sprintf("windowsKeyVaultID%dCertificateURL%d", i, j), c.CertificateURL)
				acsengine.AddValue(parametersMap, fmt.Sprintf("windowsKeyVaultID%dCertificateStore%d", i, j), c.CertificateStore)
			}
		}
	}
	return parametersMap, nil
}

func (i *Interpolator) GetTemplate() ([]byte, error) {
	if i.interpolated == false {
		return []byte(""), fmt.Errorf("Unable to get template before calling Interpolate()")
	}
	return []byte(""), nil
}

func (i *Interpolator) GetParameters() ([]byte, error) {
	if i.interpolated == false {
		return []byte(""), fmt.Errorf("Unable to get template before calling Interpolate()")
	}
	return []byte(""), nil
}
