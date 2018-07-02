package acsengine

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"hash/fnv"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"text/template"

	//log "github.com/sirupsen/logrus"
	"github.com/Azure/acs-engine/pkg/api"
	"github.com/Azure/acs-engine/pkg/api/common"
	"github.com/Azure/acs-engine/pkg/helpers"
	"github.com/ghodss/yaml"
	"github.com/pkg/errors"
)

var commonTemplateFiles = []string{agentOutputs, agentParams, classicParams, masterOutputs, iaasOutputs, masterParams, windowsParams}
var dcosTemplateFiles = []string{dcosBaseFile, dcosAgentResourcesVMAS, dcosAgentResourcesVMSS, dcosAgentVars, dcosMasterResources, dcosMasterVars, dcosParams, dcosWindowsAgentResourcesVMAS, dcosWindowsAgentResourcesVMSS}
var dcos2TemplateFiles = []string{dcos2BaseFile, dcosAgentResourcesVMAS, dcosAgentResourcesVMSS, dcosAgentVars, dcos2MasterResources, dcos2BootstrapResources, dcos2MasterVars, dcosParams, dcosWindowsAgentResourcesVMAS, dcosWindowsAgentResourcesVMSS, dcos2BootstrapVars, dcos2BootstrapParams}
var kubernetesTemplateFiles = []string{kubernetesBaseFile, kubernetesAgentResourcesVMAS, kubernetesAgentResourcesVMSS, kubernetesAgentVars, kubernetesMasterResources, kubernetesMasterVars, kubernetesParams, kubernetesWinAgentVars, kubernetesWinAgentVarsVMSS}
var swarmTemplateFiles = []string{swarmBaseFile, swarmParams, swarmAgentResourcesVMAS, swarmAgentVars, swarmAgentResourcesVMSS, swarmAgentResourcesClassic, swarmBaseFile, swarmMasterResources, swarmMasterVars, swarmWinAgentResourcesVMAS, swarmWinAgentResourcesVMSS}
var swarmModeTemplateFiles = []string{swarmBaseFile, swarmParams, swarmAgentResourcesVMAS, swarmAgentVars, swarmAgentResourcesVMSS, swarmAgentResourcesClassic, swarmBaseFile, swarmMasterResources, swarmMasterVars, swarmWinAgentResourcesVMAS, swarmWinAgentResourcesVMSS}
var openshiftTemplateFiles = append(
	kubernetesTemplateFiles,
	openshiftInfraResources,
	openshiftNodeScript,
	openshiftMasterScript,
	openshift39NodeScript,
	openshift39MasterScript,
)

var keyvaultSecretPathRe *regexp.Regexp

func init() {
	keyvaultSecretPathRe = regexp.MustCompile(`^(/subscriptions/\S+/resourceGroups/\S+/providers/Microsoft.KeyVault/vaults/\S+)/secrets/([^/\s]+)(/(\S+))?$`)
}

// GenerateClusterID creates a unique 8 string cluster ID
func GenerateClusterID(properties *api.Properties) string {
	uniqueNameSuffixSize := 8
	// the name suffix uniquely identifies the cluster and is generated off a hash
	// from the master dns name
	h := fnv.New64a()
	if properties.MasterProfile != nil {
		h.Write([]byte(properties.MasterProfile.DNSPrefix))
	} else if properties.HostedMasterProfile != nil {
		h.Write([]byte(properties.HostedMasterProfile.DNSPrefix))
	} else {
		h.Write([]byte(properties.AgentPoolProfiles[0].Name))
	}
	rand.Seed(int64(h.Sum64()))
	return fmt.Sprintf("%08d", rand.Uint32())[:uniqueNameSuffixSize]
}

// GenerateKubeConfig returns a JSON string representing the KubeConfig
func GenerateKubeConfig(properties *api.Properties, location string) (string, error) {
	if properties == nil {
		return "", errors.New("Properties nil in GenerateKubeConfig")
	}
	if properties.CertificateProfile == nil {
		return "", errors.New("CertificateProfile property may not be nil in GenerateKubeConfig")
	}
	b, err := Asset(kubeConfigJSON)
	if err != nil {
		return "", errors.Wrapf(err, "error reading kube config template file %s", kubeConfigJSON)
	}
	kubeconfig := string(b)
	// variable replacement
	kubeconfig = strings.Replace(kubeconfig, "{{WrapAsVerbatim \"variables('caCertificate')\"}}", base64.StdEncoding.EncodeToString([]byte(properties.CertificateProfile.CaCertificate)), -1)
	if properties.OrchestratorProfile != nil &&
		properties.OrchestratorProfile.KubernetesConfig != nil &&
		properties.OrchestratorProfile.KubernetesConfig.PrivateCluster != nil &&
		helpers.IsTrueBoolPointer(properties.OrchestratorProfile.KubernetesConfig.PrivateCluster.Enabled) {
		if properties.MasterProfile.Count > 1 {
			// more than 1 master, use the internal lb IP
			firstMasterIP := net.ParseIP(properties.MasterProfile.FirstConsecutiveStaticIP).To4()
			if firstMasterIP == nil {
				return "", errors.Errorf("MasterProfile.FirstConsecutiveStaticIP '%s' is an invalid IP address", properties.MasterProfile.FirstConsecutiveStaticIP)
			}
			lbIP := net.IP{firstMasterIP[0], firstMasterIP[1], firstMasterIP[2], firstMasterIP[3] + byte(DefaultInternalLbStaticIPOffset)}
			kubeconfig = strings.Replace(kubeconfig, "{{WrapAsVerbatim \"reference(concat('Microsoft.Network/publicIPAddresses/', variables('masterPublicIPAddressName'))).dnsSettings.fqdn\"}}", lbIP.String(), -1)
		} else {
			// Master count is 1, use the master IP
			kubeconfig = strings.Replace(kubeconfig, "{{WrapAsVerbatim \"reference(concat('Microsoft.Network/publicIPAddresses/', variables('masterPublicIPAddressName'))).dnsSettings.fqdn\"}}", properties.MasterProfile.FirstConsecutiveStaticIP, -1)
		}
	} else {
		kubeconfig = strings.Replace(kubeconfig, "{{WrapAsVerbatim \"reference(concat('Microsoft.Network/publicIPAddresses/', variables('masterPublicIPAddressName'))).dnsSettings.fqdn\"}}", FormatAzureProdFQDN(properties.MasterProfile.DNSPrefix, location), -1)
	}
	kubeconfig = strings.Replace(kubeconfig, "{{WrapAsVariable \"resourceGroup\"}}", properties.MasterProfile.DNSPrefix, -1)

	var authInfo string
	if properties.AADProfile == nil {
		authInfo = fmt.Sprintf("{\"client-certificate-data\":\"%v\",\"client-key-data\":\"%v\"}",
			base64.StdEncoding.EncodeToString([]byte(properties.CertificateProfile.KubeConfigCertificate)),
			base64.StdEncoding.EncodeToString([]byte(properties.CertificateProfile.KubeConfigPrivateKey)))
	} else {
		tenantID := properties.AADProfile.TenantID
		if len(tenantID) == 0 {
			tenantID = "common"
		}

		authInfo = fmt.Sprintf("{\"auth-provider\":{\"name\":\"azure\",\"config\":{\"environment\":\"%v\",\"tenant-id\":\"%v\",\"apiserver-id\":\"%v\",\"client-id\":\"%v\"}}}",
			getCloudTargetEnv(location),
			tenantID,
			properties.AADProfile.ServerAppID,
			properties.AADProfile.ClientAppID)
	}
	kubeconfig = strings.Replace(kubeconfig, "{{authInfo}}", authInfo, -1)

	return kubeconfig, nil
}

// formatAzureProdFQDNs constructs all possible Azure prod fqdn
func formatAzureProdFQDNs(fqdnPrefix string) []string {
	var fqdns []string
	for _, location := range AzureLocations {
		fqdns = append(fqdns, FormatAzureProdFQDN(fqdnPrefix, location))
	}
	return fqdns
}

// FormatAzureProdFQDN constructs an Azure prod fqdn
func FormatAzureProdFQDN(fqdnPrefix string, location string) string {
	var FQDNFormat string
	switch getCloudTargetEnv(location) {
	case azureChinaCloud:
		FQDNFormat = AzureChinaCloudSpec.EndpointConfig.ResourceManagerVMDNSSuffix
	case azureGermanCloud:
		FQDNFormat = AzureGermanCloudSpec.EndpointConfig.ResourceManagerVMDNSSuffix
	case azureUSGovernmentCloud:
		FQDNFormat = AzureUSGovernmentCloud.EndpointConfig.ResourceManagerVMDNSSuffix
	default:
		FQDNFormat = AzureCloudSpec.EndpointConfig.ResourceManagerVMDNSSuffix
	}
	return fmt.Sprintf("%s.%s."+FQDNFormat, fqdnPrefix, location)
}

//getCloudSpecConfig returns the kubenernetes container images url configurations based on the deploy target environment
//for example: if the target is the public azure, then the default container image url should be k8s-gcrio.azureedge.net/...
//if the target is azure china, then the default container image should be mirror.azure.cn:5000/google_container/...
func getCloudSpecConfig(location string) AzureEnvironmentSpecConfig {
	switch getCloudTargetEnv(location) {
	case azureChinaCloud:
		return AzureChinaCloudSpec
	case azureGermanCloud:
		return AzureGermanCloudSpec
	case azureUSGovernmentCloud:
		return AzureUSGovernmentCloud
	default:
		return AzureCloudSpec
	}
}

// validateDistro checks if the requested orchestrator type is supported on the requested Linux distro.
func validateDistro(cs *api.ContainerService) bool {
	// Check Master distro
	if cs.Properties.MasterProfile != nil && cs.Properties.MasterProfile.Distro == api.RHEL &&
		(cs.Properties.OrchestratorProfile.OrchestratorType != api.SwarmMode && cs.Properties.OrchestratorProfile.OrchestratorType != api.OpenShift) {
		log.Fatalf("Orchestrator type %s not suported on RHEL Master", cs.Properties.OrchestratorProfile.OrchestratorType)
		return false
	}
	// Check Agent distros
	for _, agentProfile := range cs.Properties.AgentPoolProfiles {
		if agentProfile.Distro == api.RHEL &&
			(cs.Properties.OrchestratorProfile.OrchestratorType != api.SwarmMode && cs.Properties.OrchestratorProfile.OrchestratorType != api.OpenShift) {
			log.Fatalf("Orchestrator type %s not suported on RHEL Agent", cs.Properties.OrchestratorProfile.OrchestratorType)
			return false
		}
	}
	return true
}

// getCloudTargetEnv determines and returns whether the region is a sovereign cloud which
// have their own data compliance regulations (China/Germany/USGov) or standard
//  Azure public cloud
func getCloudTargetEnv(location string) string {
	loc := strings.ToLower(strings.Join(strings.Fields(location), ""))
	switch {
	case loc == "chinaeast" || loc == "chinanorth" || loc == "chinaeast2" || loc == "chinanorth2":
		return azureChinaCloud
	case loc == "germanynortheast" || loc == "germanycentral":
		return azureGermanCloud
	case strings.HasPrefix(loc, "usgov") || strings.HasPrefix(loc, "usdod"):
		return azureUSGovernmentCloud
	default:
		return azurePublicCloud
	}
}

func getOpenshiftMasterShAsset(version string) string {
	switch version {
	case common.OpenShiftVersion3Dot9Dot0:
		return openshift39MasterScript
	case common.OpenShiftVersionUnstable:
		return openshiftMasterScript
	default:
		panic(fmt.Sprintf("BUG: invalid OpenShift version %s", version))
	}
}

func getOpenshiftNodeShAsset(version string) string {
	switch version {
	case common.OpenShiftVersion3Dot9Dot0:
		return openshift39NodeScript
	case common.OpenShiftVersionUnstable:
		return openshiftNodeScript
	default:
		panic(fmt.Sprintf("BUG: invalid OpenShift version %s", version))
	}
}

func generateIPList(count int, firstAddr string) []string {
	ipaddr := net.ParseIP(firstAddr).To4()
	if ipaddr == nil {
		panic(fmt.Sprintf("IPAddr '%s' is an invalid IP address", firstAddr))
	}
	ret := make([]string, count)
	for i := 0; i < count; i++ {
		ret[i] = fmt.Sprintf("%d.%d.%d.%d", ipaddr[0], ipaddr[1], ipaddr[2], ipaddr[3]+byte(i))
	}
	return ret
}

func addValue(m paramsMap, k string, v interface{}) {
	m[k] = paramsMap{
		"value": v,
	}
}

func addKeyvaultReference(m paramsMap, k string, vaultID, secretName, secretVersion string) {
	m[k] = paramsMap{
		"reference": &KeyVaultRef{
			KeyVault: KeyVaultID{
				ID: vaultID,
			},
			SecretName:    secretName,
			SecretVersion: secretVersion,
		},
	}
}

func addSecret(m paramsMap, k string, v interface{}, encode bool) {
	str, ok := v.(string)
	if !ok {
		addValue(m, k, v)
		return
	}
	parts := keyvaultSecretPathRe.FindStringSubmatch(str)
	if parts == nil || len(parts) != 5 {
		if encode {
			addValue(m, k, base64.StdEncoding.EncodeToString([]byte(str)))
		} else {
			addValue(m, k, str)
		}
		return
	}
	addKeyvaultReference(m, k, parts[1], parts[2], parts[4])
}

// getStorageAccountType returns the support managed disk storage tier for a give VM size
func getStorageAccountType(sizeName string) (string, error) {
	spl := strings.Split(sizeName, "_")
	if len(spl) < 2 {
		return "", errors.Errorf("Invalid sizeName: %s", sizeName)
	}
	capability := spl[1]
	if strings.Contains(strings.ToLower(capability), "s") {
		return "Premium_LRS", nil
	}
	return "Standard_LRS", nil
}

func makeMasterExtensionScriptCommands(cs *api.ContainerService) string {
	copyIndex := "',copyIndex(),'"
	if cs.Properties.OrchestratorProfile.IsKubernetes() || cs.Properties.OrchestratorProfile.IsOpenShift() {
		copyIndex = "',copyIndex(variables('masterOffset')),'"
	}
	return makeExtensionScriptCommands(cs.Properties.MasterProfile.PreprovisionExtension,
		cs.Properties.ExtensionProfiles, copyIndex)
}

func makeAgentExtensionScriptCommands(cs *api.ContainerService, profile *api.AgentPoolProfile) string {
	copyIndex := "',copyIndex(),'"
	if profile.IsAvailabilitySets() {
		copyIndex = fmt.Sprintf("',copyIndex(variables('%sOffset')),'", profile.Name)
	}
	if profile.OSType == api.Windows {
		return makeWindowsExtensionScriptCommands(profile.PreprovisionExtension,
			cs.Properties.ExtensionProfiles, copyIndex)
	}
	return makeExtensionScriptCommands(profile.PreprovisionExtension,
		cs.Properties.ExtensionProfiles, copyIndex)
}

func makeExtensionScriptCommands(extension *api.Extension, extensionProfiles []*api.ExtensionProfile, copyIndex string) string {
	var extensionProfile *api.ExtensionProfile
	for _, eP := range extensionProfiles {
		if strings.EqualFold(eP.Name, extension.Name) {
			extensionProfile = eP
			break
		}
	}

	if extensionProfile == nil {
		panic(fmt.Sprintf("%s extension referenced was not found in the extension profile", extension.Name))
	}

	extensionsParameterReference := fmt.Sprintf("parameters('%sParameters')", extensionProfile.Name)
	scriptURL := getExtensionURL(extensionProfile.RootURL, extensionProfile.Name, extensionProfile.Version, extensionProfile.Script, extensionProfile.URLQuery)
	scriptFilePath := fmt.Sprintf("/opt/azure/containers/extensions/%s/%s", extensionProfile.Name, extensionProfile.Script)
	return fmt.Sprintf("- sudo /usr/bin/curl --retry 5 --retry-delay 10 --retry-max-time 30 -o %s --create-dirs \"%s\" \n- sudo /bin/chmod 744 %s \n- sudo %s ',%s,' > /var/log/%s-output.log",
		scriptFilePath, scriptURL, scriptFilePath, scriptFilePath, extensionsParameterReference, extensionProfile.Name)
}

func makeWindowsExtensionScriptCommands(extension *api.Extension, extensionProfiles []*api.ExtensionProfile, copyIndex string) string {
	var extensionProfile *api.ExtensionProfile
	for _, eP := range extensionProfiles {
		if strings.EqualFold(eP.Name, extension.Name) {
			extensionProfile = eP
			break
		}
	}

	if extensionProfile == nil {
		panic(fmt.Sprintf("%s extension referenced was not found in the extension profile", extension.Name))
	}

	scriptURL := getExtensionURL(extensionProfile.RootURL, extensionProfile.Name, extensionProfile.Version, extensionProfile.Script, extensionProfile.URLQuery)
	scriptFileDir := fmt.Sprintf("$env:SystemDrive:/AzureData/extensions/%s", extensionProfile.Name)
	scriptFilePath := fmt.Sprintf("%s/%s", scriptFileDir, extensionProfile.Script)
	return fmt.Sprintf("New-Item -ItemType Directory -Force -Path \"%s\" ; Invoke-WebRequest -Uri \"%s\" -OutFile \"%s\" ; powershell \"%s %s\"\n", scriptFileDir, scriptURL, scriptFilePath, scriptFilePath, "$preprovisionExtensionParams")
}

func getDCOSWindowsAgentPreprovisionParameters(cs *api.ContainerService, profile *api.AgentPoolProfile) string {
	extension := profile.PreprovisionExtension

	var extensionProfile *api.ExtensionProfile

	for _, eP := range cs.Properties.ExtensionProfiles {
		if strings.EqualFold(eP.Name, extension.Name) {
			extensionProfile = eP
			break
		}
	}

	parms := extensionProfile.ExtensionParameters
	return parms
}

func getDCOSDefaultBootstrapInstallerURL(profile *api.OrchestratorProfile) string {
	if profile.OrchestratorType == api.DCOS {
		switch profile.OrchestratorVersion {
		case common.DCOSVersion1Dot11Dot2:
			return "https://dcos-mirror.azureedge.net/dcos-1-11-2/dcos_generate_config.sh"
		case common.DCOSVersion1Dot11Dot0:
			return "https://dcos-mirror.azureedge.net/dcos-1-11-0/dcos_generate_config.sh"
		}
	}
	return ""
}

func getDCOSDefaultWindowsBootstrapInstallerURL(profile *api.OrchestratorProfile) string {
	if profile.OrchestratorType == api.DCOS {
		switch profile.OrchestratorVersion {
		case common.DCOSVersion1Dot11Dot2:
			return "https://dcos-mirror.azureedge.net/dcos-windows/1-11-2"
		case common.DCOSVersion1Dot11Dot0:
			return "https://dcos-mirror.azureedge.net/dcos-windows/1-11-0"
		}
	}
	return ""
}

func getDCOSDefaultProviderPackageGUID(orchestratorType string, orchestratorVersion string, masterCount int) string {
	if orchestratorType == api.DCOS {
		switch orchestratorVersion {
		case common.DCOSVersion1Dot10Dot0:
			switch masterCount {
			case 1:
				return "c4ec6210f396b8e435177b82e3280a2cef0ce721"
			case 3:
				return "08197947cb57d479eddb077a429fa15c139d7d20"
			case 5:
				return "f286ad9d3641da5abb622e4a8781f73ecd8492fa"
			}
		case common.DCOSVersion1Dot9Dot0:
			switch masterCount {
			case 1:
				return "bcc883b7a3191412cf41824bdee06c1142187a0b"
			case 3:
				return "dcff7e24c0c1827bebeb7f1a806f558054481b33"
			case 5:
				return "b41bfa84137a6374b2ff5eb1655364d7302bd257"
			}
		case common.DCOSVersion1Dot9Dot8:
			switch masterCount {
			case 1:
				return "e8b0e3fc4a16394dc6dd5b19fc54bf1543bff429"
			case 3:
				return "2d36c3f570d9dd7d187c699f9a322ed9d95e7dfa"
			case 5:
				return "c03c9587f88929f310b80af4f448b7b51654f1c8"
			}
		case common.DCOSVersion1Dot8Dot8:
			switch masterCount {
			case 1:
				return "441385ce2f5942df7e29075c12fb38fa5e92cbba"
			case 3:
				return "b1cd359287504efb780257bd12cc3a63704e42d4"
			case 5:
				return "d9b61156dfcc9383e014851529738aa550ef57d9"
			}
		}
	}
	return ""
}

func getDCOSDefaultRepositoryURL(orchestratorType string, orchestratorVersion string) string {
	if orchestratorType == api.DCOS {
		switch orchestratorVersion {
		case common.DCOSVersion1Dot10Dot0:
			return "https://dcosio.azureedge.net/dcos/stable/1.10.0"
		case common.DCOSVersion1Dot9Dot8:
			return "https://dcosio.azureedge.net/dcos/stable/1.9.8"
		default:
			return "https://dcosio.azureedge.net/dcos/stable"
		}
	}
	return ""
}

func isNSeriesSKU(profile *api.AgentPoolProfile) bool {
	return strings.Contains(profile.VMSize, "Standard_N")
}

func isCustomVNET(a []*api.AgentPoolProfile) bool {
	if a != nil {
		for _, agentPoolProfile := range a {
			if !agentPoolProfile.IsCustomVNET() {
				return false
			}
		}
		return true
	}
	return false
}

func getGPUDriversInstallScript(profile *api.AgentPoolProfile) string {

	// latest version of the drivers. Later this parameter could be bubbled up so that users can choose specific driver versions.
	dv := "396.26"
	dest := "/usr/local/nvidia"
	nvidiaDockerVersion := "2.0.3"
	dockerVersion := "1.13.1-1"
	nvidiaContainerRuntimeVersion := "2.0.0"
	/*
		First we remove the nouveau drivers, which are the open source drivers for NVIDIA cards. Nouveau is installed on NV Series VMs by default.
		We also installed needed dependencies.
	*/
	installScript := fmt.Sprintf(`- rmmod nouveau
- sh -c "echo \"blacklist nouveau\" >> /etc/modprobe.d/blacklist.conf"
- update-initramfs -u
- mkdir -p %s
- cd %s`, dest, dest)

	/*
		Installing nvidia-docker, setting nvidia runtime as default and restarting docker daemon
	*/
	installScript += fmt.Sprintf(`
- retrycmd_if_failure_no_stats 180 1 5 curl -fsSL https://nvidia.github.io/nvidia-docker/gpgkey > /tmp/aptnvidia.gpg
- cat /tmp/aptnvidia.gpg | apt-key add -
- retrycmd_if_failure_no_stats 180 1 5 curl -fsSL https://nvidia.github.io/nvidia-docker/ubuntu16.04/amd64/nvidia-docker.list > /tmp/nvidia-docker.list
- cat /tmp/nvidia-docker.list | sudo tee /etc/apt/sources.list.d/nvidia-docker.list
- apt_get_update
- retrycmd_if_failure 5 5 300 apt-get install -y linux-headers-$(uname -r) gcc make
- retrycmd_if_failure 5 5 300 apt-get -o Dpkg::Options::="--force-confold" install -y nvidia-docker2=%s+docker%s nvidia-container-runtime=%s+docker%s
- sudo pkill -SIGHUP dockerd
- mkdir -p %s
- cd %s`, nvidiaDockerVersion, dockerVersion, nvidiaContainerRuntimeVersion, dockerVersion, dest, dest)

	/*
		Download the .run file from NVIDIA.
		Nvidia libraries are always install in /usr/lib/x86_64-linux-gnu, and there is no option in the run file to change this.
		Instead we use Overlayfs to move the newly installed libraries under /usr/local/nvidia/lib64
	*/
	installScript += fmt.Sprintf(`
- retrycmd_if_failure 5 10 60 curl -fLS https://us.download.nvidia.com/tesla/%s/NVIDIA-Linux-x86_64-%s.run -o nvidia-drivers-%s
- mkdir -p lib64 overlay-workdir
- mount -t overlay -o lowerdir=/usr/lib/x86_64-linux-gnu,upperdir=lib64,workdir=overlay-workdir none /usr/lib/x86_64-linux-gnu`, dv, dv, dv)

	/*
		Install the drivers and update /etc/ld.so.conf.d/nvidia.conf which will make the libraries discoverable through $LD_LIBRARY_PATH.
		Run nvidia-smi to test the installation, unmount overlayfs and restard kubelet (GPUs are only discovered when kubelet starts)
	*/
	installScript += fmt.Sprintf(`
- sh nvidia-drivers-%s --silent --accept-license --no-drm --utility-prefix="%s" --opengl-prefix="%s"
- echo "%s" > /etc/ld.so.conf.d/nvidia.conf
- sudo ldconfig
- umount /usr/lib/x86_64-linux-gnu
- nvidia-modprobe -u -c0
- %s/bin/nvidia-smi
- sudo ldconfig
- retrycmd_if_failure 5 10 60 systemctl restart kubelet`, dv, dest, dest, fmt.Sprintf("%s/lib64", dest), dest)

	/* If a new GPU sku becomes available, add a key to this map, but only provide an installation script if you have a confirmation
	   that we have an agreement with NVIDIA for this specific gpu. Otherwise use the warning message.
	*/
	dm := map[string]string{
		// K80
		"Standard_NC6":   installScript,
		"Standard_NC12":  installScript,
		"Standard_NC24":  installScript,
		"Standard_NC24r": installScript,
		// M60
		"Standard_NV6":   installScript,
		"Standard_NV12":  installScript,
		"Standard_NV24":  installScript,
		"Standard_NV24r": installScript,
		// P40
		"Standard_ND6s":   installScript,
		"Standard_ND12s":  installScript,
		"Standard_ND24s":  installScript,
		"Standard_ND24rs": installScript,
		// P100
		"Standard_NC6s_v2":   installScript,
		"Standard_NC12s_v2":  installScript,
		"Standard_NC24s_v2":  installScript,
		"Standard_NC24rs_v2": installScript,
		// V100
		"Standard_NC6s_v3":   installScript,
		"Standard_NC12s_v3":  installScript,
		"Standard_NC24s_v3":  installScript,
		"Standard_NC24rs_v3": installScript,
	}
	if _, ok := dm[profile.VMSize]; ok {
		return dm[profile.VMSize]
	}

	// The VM is not part of the GPU skus, no extra steps.
	return ""
}

func getDCOSCustomDataPublicIPStr(orchestratorType string, masterCount int) string {
	if orchestratorType == api.DCOS {
		var buf bytes.Buffer
		for i := 0; i < masterCount; i++ {
			buf.WriteString(fmt.Sprintf("reference(variables('masterVMNic')[%d]).ipConfigurations[0].properties.privateIPAddress,", i))
			if i < (masterCount - 1) {
				buf.WriteString(`'\\\", \\\"', `)
			}
		}
		return buf.String()
	}
	return ""
}

func getDCOSMasterCustomNodeLabels() string {
	// return empty string for DCOS since no attribtutes needed on master
	return ""
}

func getDCOSAgentCustomNodeLabels(profile *api.AgentPoolProfile) string {
	var buf bytes.Buffer
	var attrstring string
	buf.WriteString("")
	// always write MESOS_ATTRIBUTES because
	// the provision script will add FD/UD attributes
	// at node provisioning time
	if len(profile.OSType) > 0 {
		attrstring = fmt.Sprintf("MESOS_ATTRIBUTES=\"os:%s", profile.OSType)
	} else {
		attrstring = fmt.Sprintf("MESOS_ATTRIBUTES=\"os:%s", api.Linux)
	}

	if len(profile.Ports) > 0 {
		attrstring += ";public_ip:yes"
	}

	buf.WriteString(attrstring)
	if len(profile.CustomNodeLabels) > 0 {
		for k, v := range profile.CustomNodeLabels {
			buf.WriteString(fmt.Sprintf(";%s:%s", k, v))
		}
	}
	buf.WriteString("\"")
	return buf.String()
}

func getDCOSWindowsAgentCustomAttributes(profile *api.AgentPoolProfile) string {
	var buf bytes.Buffer
	var attrstring string
	buf.WriteString("")
	if len(profile.OSType) > 0 {
		attrstring = fmt.Sprintf("os:%s", profile.OSType)
	} else {
		attrstring = fmt.Sprintf("os:windows")
	}
	if len(profile.Ports) > 0 {
		attrstring += ";public_ip:yes"
	}
	buf.WriteString(attrstring)
	if len(profile.CustomNodeLabels) > 0 {
		for k, v := range profile.CustomNodeLabels {
			buf.WriteString(fmt.Sprintf(";%s:%s", k, v))
		}
	}
	return buf.String()
}

func getVNETAddressPrefixes(properties *api.Properties) string {
	visitedSubnets := make(map[string]bool)
	var buf bytes.Buffer
	buf.WriteString(`"[variables('masterSubnet')]"`)
	visitedSubnets[properties.MasterProfile.Subnet] = true
	for _, profile := range properties.AgentPoolProfiles {
		if _, ok := visitedSubnets[profile.Subnet]; !ok {
			buf.WriteString(fmt.Sprintf(",\n            \"[variables('%sSubnet')]\"", profile.Name))
		}
	}
	return buf.String()
}

func getVNETSubnetDependencies(properties *api.Properties) string {
	agentString := `        "[concat('Microsoft.Network/networkSecurityGroups/', variables('%sNSGName'))]"`
	var buf bytes.Buffer
	for index, agentProfile := range properties.AgentPoolProfiles {
		if index > 0 {
			buf.WriteString(",\n")
		}
		buf.WriteString(fmt.Sprintf(agentString, agentProfile.Name))
	}
	return buf.String()
}

func getVNETSubnets(properties *api.Properties, addNSG bool) string {
	masterString := `{
            "name": "[variables('masterSubnetName')]",
            "properties": {
              "addressPrefix": "[variables('masterSubnet')]"
            }
          }`
	agentString := `          {
            "name": "[variables('%sSubnetName')]",
            "properties": {
              "addressPrefix": "[variables('%sSubnet')]"
            }
          }`
	agentStringNSG := `          {
            "name": "[variables('%sSubnetName')]",
            "properties": {
              "addressPrefix": "[variables('%sSubnet')]",
              "networkSecurityGroup": {
                "id": "[resourceId('Microsoft.Network/networkSecurityGroups', variables('%sNSGName'))]"
              }
            }
          }`
	var buf bytes.Buffer
	buf.WriteString(masterString)
	for _, agentProfile := range properties.AgentPoolProfiles {
		buf.WriteString(",\n")
		if addNSG {
			buf.WriteString(fmt.Sprintf(agentStringNSG, agentProfile.Name, agentProfile.Name, agentProfile.Name))
		} else {
			buf.WriteString(fmt.Sprintf(agentString, agentProfile.Name, agentProfile.Name))
		}

	}
	return buf.String()
}

func getLBRule(name string, port int) string {
	return fmt.Sprintf(`	          {
            "name": "LBRule%d",
            "properties": {
              "backendAddressPool": {
                "id": "[concat(variables('%sLbID'), '/backendAddressPools/', variables('%sLbBackendPoolName'))]"
              },
              "backendPort": %d,
              "enableFloatingIP": false,
              "frontendIPConfiguration": {
                "id": "[variables('%sLbIPConfigID')]"
              },
              "frontendPort": %d,
              "idleTimeoutInMinutes": 5,
              "loadDistribution": "Default",
              "probe": {
                "id": "[concat(variables('%sLbID'),'/probes/tcp%dProbe')]"
              },
              "protocol": "tcp"
            }
          }`, port, name, name, port, name, port, name, port)
}

func getLBRules(name string, ports []int) string {
	var buf bytes.Buffer
	for index, port := range ports {
		if index > 0 {
			buf.WriteString(",\n")
		}
		buf.WriteString(getLBRule(name, port))
	}
	return buf.String()
}

func getProbe(port int) string {
	return fmt.Sprintf(`          {
            "name": "tcp%dProbe",
            "properties": {
              "intervalInSeconds": "5",
              "numberOfProbes": "2",
              "port": %d,
              "protocol": "tcp"
            }
          }`, port, port)
}

func getProbes(ports []int) string {
	var buf bytes.Buffer
	for index, port := range ports {
		if index > 0 {
			buf.WriteString(",\n")
		}
		buf.WriteString(getProbe(port))
	}
	return buf.String()
}

func getSecurityRule(port int, portIndex int) string {
	// BaseLBPriority specifies the base lb priority.
	BaseLBPriority := 200
	return fmt.Sprintf(`          {
            "name": "Allow_%d",
            "properties": {
              "access": "Allow",
              "description": "Allow traffic from the Internet to port %d",
              "destinationAddressPrefix": "*",
              "destinationPortRange": "%d",
              "direction": "Inbound",
              "priority": %d,
              "protocol": "*",
              "sourceAddressPrefix": "Internet",
              "sourcePortRange": "*"
            }
          }`, port, port, port, BaseLBPriority+portIndex)
}

func getDataDisks(a *api.AgentPoolProfile) string {
	if !a.HasDisks() {
		return ""
	}
	var buf bytes.Buffer
	buf.WriteString("\"dataDisks\": [\n")
	dataDisks := `            {
              "createOption": "Empty",
              "diskSizeGB": "%d",
              "lun": %d,
              "name": "[concat(variables('%sVMNamePrefix'), copyIndex(),'-datadisk%d')]",
              "vhd": {
                "uri": "[concat('http://',variables('storageAccountPrefixes')[mod(add(add(div(copyIndex(),variables('maxVMsPerStorageAccount')),variables('%sStorageAccountOffset')),variables('dataStorageAccountPrefixSeed')),variables('storageAccountPrefixesCount'))],variables('storageAccountPrefixes')[div(add(add(div(copyIndex(),variables('maxVMsPerStorageAccount')),variables('%sStorageAccountOffset')),variables('dataStorageAccountPrefixSeed')),variables('storageAccountPrefixesCount'))],variables('%sDataAccountName'),'.blob.core.windows.net/vhds/',variables('%sVMNamePrefix'),copyIndex(), '--datadisk%d.vhd')]"
              }
            }`
	managedDataDisks := `            {
              "diskSizeGB": "%d",
              "lun": %d,
              "createOption": "Empty"
            }`
	for i, diskSize := range a.DiskSizesGB {
		if i > 0 {
			buf.WriteString(",\n")
		}
		if a.StorageProfile == api.StorageAccount {
			buf.WriteString(fmt.Sprintf(dataDisks, diskSize, i, a.Name, i, a.Name, a.Name, a.Name, a.Name, i))
		} else if a.StorageProfile == api.ManagedDisks {
			buf.WriteString(fmt.Sprintf(managedDataDisks, diskSize, i))
		}
	}
	buf.WriteString("\n          ],")
	return buf.String()
}

func getSecurityRules(ports []int) string {
	var buf bytes.Buffer
	for index, port := range ports {
		if index > 0 {
			buf.WriteString(",\n")
		}
		buf.WriteString(getSecurityRule(port, index))
	}
	return buf.String()
}

// getSingleLineForTemplate returns the file as a single line for embedding in an arm template
func (t *TemplateGenerator) getSingleLineForTemplate(textFilename string, cs *api.ContainerService, profile interface{}) (string, error) {
	b, err := Asset(textFilename)
	if err != nil {
		return "", t.Translator.Errorf("yaml file %s does not exist", textFilename)
	}

	// use go templates to process the text filename
	templ := template.New("customdata template").Funcs(t.getTemplateFuncMap(cs))
	if _, err = templ.New(textFilename).Parse(string(b)); err != nil {
		return "", t.Translator.Errorf("error parsing file %s: %v", textFilename, err)
	}

	var buffer bytes.Buffer
	if err = templ.ExecuteTemplate(&buffer, textFilename, profile); err != nil {
		return "", t.Translator.Errorf("error executing template for file %s: %v", textFilename, err)
	}
	expandedTemplate := buffer.String()

	textStr := escapeSingleLine(string(expandedTemplate))

	return textStr, nil
}

func escapeSingleLine(escapedStr string) string {
	// template.JSEscapeString leaves undesirable chars that don't work with pretty print
	escapedStr = strings.Replace(escapedStr, "\\", "\\\\", -1)
	escapedStr = strings.Replace(escapedStr, "\r\n", "\\n", -1)
	escapedStr = strings.Replace(escapedStr, "\n", "\\n", -1)
	escapedStr = strings.Replace(escapedStr, "\"", "\\\"", -1)
	return escapedStr
}

// getBase64CustomScript will return a base64 of the CSE
func getBase64CustomScript(csFilename string) string {
	b, err := Asset(csFilename)
	if err != nil {
		// this should never happen and this is a bug
		panic(fmt.Sprintf("BUG: %s", err.Error()))
	}
	// translate the parameters
	csStr := string(b)
	csStr = strings.Replace(csStr, "\r\n", "\n", -1)
	return getBase64CustomScriptFromStr(csStr)
}

// getBase64CustomScript will return a base64 of the CSE
func getBase64CustomScriptFromStr(str string) string {
	var gzipB bytes.Buffer
	w := gzip.NewWriter(&gzipB)
	w.Write([]byte(str))
	w.Close()
	return base64.StdEncoding.EncodeToString(gzipB.Bytes())
}

func getDCOSProvisionScript(script string) string {
	// add the provision script
	bp, err := Asset(script)
	if err != nil {
		panic(fmt.Sprintf("BUG: %s", err.Error()))
	}

	provisionScript := string(bp)
	if strings.Contains(provisionScript, "'") {
		panic(fmt.Sprintf("BUG: %s may not contain character '", script))
	}

	return strings.Replace(strings.Replace(provisionScript, "\r\n", "\n", -1), "\n", "\n\n    ", -1)
}

func getDCOSAgentProvisionScript(profile *api.AgentPoolProfile, orchProfile *api.OrchestratorProfile, bootstrapIP string) string {
	// add the provision script
	scriptname := dcos2Provision
	if orchProfile.DcosConfig == nil || orchProfile.DcosConfig.BootstrapProfile == nil {
		if profile.OSType == api.Windows {
			scriptname = dcosWindowsProvision
		} else {
			scriptname = dcosProvision
		}
	}

	bp, err := Asset(scriptname)
	if err != nil {
		panic(fmt.Sprintf("BUG: %s", err.Error()))
	}

	provisionScript := string(bp)
	if strings.Contains(provisionScript, "'") {
		panic(fmt.Sprintf("BUG: %s may not contain character '", dcosProvision))
	}

	// the embedded roleFileContents
	var roleFileContents string
	if len(profile.Ports) > 0 {
		// public agents
		roleFileContents = "touch /etc/mesosphere/roles/slave_public"
	} else {
		roleFileContents = "touch /etc/mesosphere/roles/slave"
	}
	provisionScript = strings.Replace(provisionScript, "ROLESFILECONTENTS", roleFileContents, -1)
	provisionScript = strings.Replace(provisionScript, "BOOTSTRAP_IP", bootstrapIP, -1)

	var b bytes.Buffer
	b.WriteString(provisionScript)
	b.WriteString("\n")

	if len(orchProfile.DcosConfig.Registry) == 0 {
		b.WriteString("rm /etc/docker.tar.gz\n")
	}

	return strings.Replace(strings.Replace(b.String(), "\r\n", "\n", -1), "\n", "\n\n    ", -1)
}

func getDCOSMasterProvisionScript(orchProfile *api.OrchestratorProfile, bootstrapIP string) string {
	scriptname := dcos2Provision
	if orchProfile.DcosConfig == nil || orchProfile.DcosConfig.BootstrapProfile == nil {
		scriptname = dcosProvision
	}

	// add the provision script
	bp, err := Asset(scriptname)
	if err != nil {
		panic(fmt.Sprintf("BUG: %s", err.Error()))
	}

	provisionScript := string(bp)
	if strings.Contains(provisionScript, "'") {
		panic(fmt.Sprintf("BUG: %s may not contain character '", scriptname))
	}

	// the embedded roleFileContents
	roleFileContents := `touch /etc/mesosphere/roles/master
touch /etc/mesosphere/roles/azure_master`
	provisionScript = strings.Replace(provisionScript, "ROLESFILECONTENTS", roleFileContents, -1)
	provisionScript = strings.Replace(provisionScript, "BOOTSTRAP_IP", bootstrapIP, -1)

	var b bytes.Buffer
	b.WriteString(provisionScript)
	b.WriteString("\n")

	return strings.Replace(strings.Replace(b.String(), "\r\n", "\n", -1), "\n", "\n\n    ", -1)
}

func getDCOSCustomDataTemplate(orchestratorType, orchestratorVersion string) string {
	switch orchestratorType {
	case api.DCOS:
		switch orchestratorVersion {
		case common.DCOSVersion1Dot8Dot8:
			return dcosCustomData188
		case common.DCOSVersion1Dot9Dot0:
			return dcosCustomData190
		case common.DCOSVersion1Dot9Dot8:
			return dcosCustomData198
		case common.DCOSVersion1Dot10Dot0:
			return dcosCustomData110
		case common.DCOSVersion1Dot11Dot0:
			return dcos2CustomData1110
		case common.DCOSVersion1Dot11Dot2:
			return dcos2CustomData1112
		}
	default:
		// it is a bug to get here
		panic(fmt.Sprintf("BUG: invalid orchestrator %s", orchestratorType))
	}
	return ""
}

// getSingleLineForTemplate returns the file as a single line for embedding in an arm template
func getSingleLineDCOSCustomData(orchestratorType, yamlFilename string, masterCount int, replaceMap map[string]string) string {
	b, err := Asset(yamlFilename)
	if err != nil {
		panic(fmt.Sprintf("BUG getting yaml custom data file: %s", err.Error()))
	}
	yamlStr := string(b)
	for k, v := range replaceMap {
		yamlStr = strings.Replace(yamlStr, k, v, -1)
	}

	// convert to json
	jsonBytes, err4 := yaml.YAMLToJSON([]byte(yamlStr))
	if err4 != nil {
		panic(fmt.Sprintf("BUG: %s", err4.Error()))
	}
	yamlStr = string(jsonBytes)

	// convert to one line
	yamlStr = strings.Replace(yamlStr, "\\", "\\\\", -1)
	yamlStr = strings.Replace(yamlStr, "\r\n", "\\n", -1)
	yamlStr = strings.Replace(yamlStr, "\n", "\\n", -1)
	yamlStr = strings.Replace(yamlStr, "\"", "\\\"", -1)

	// variable replacement
	rVariable, e1 := regexp.Compile("{{{([^}]*)}}}")
	if e1 != nil {
		panic(fmt.Sprintf("BUG: %s", e1.Error()))
	}
	yamlStr = rVariable.ReplaceAllString(yamlStr, "',variables('$1'),'")

	// replace the internal values
	publicIPStr := getDCOSCustomDataPublicIPStr(orchestratorType, masterCount)
	yamlStr = strings.Replace(yamlStr, "DCOSCUSTOMDATAPUBLICIPSTR", publicIPStr, -1)

	return yamlStr
}

func buildYamlFileWithWriteFiles(files []string) string {
	clusterYamlFile := `#cloud-config

write_files:
%s
`
	writeFileBlock := ` -  encoding: gzip
    content: !!binary |
        %s
    path: /opt/azure/containers/%s
    permissions: "0744"
`

	filelines := ""
	for _, file := range files {
		b64GzipString := getBase64CustomScript(file)
		fileNoPath := strings.TrimPrefix(file, "swarm/")
		filelines = filelines + fmt.Sprintf(writeFileBlock, b64GzipString, fileNoPath)
	}
	return fmt.Sprintf(clusterYamlFile, filelines)
}

// Identifies Master distro to use for master parameters
func getMasterDistro(m *api.MasterProfile) api.Distro {
	// Use Ubuntu distro if MasterProfile is not defined (e.g. agents-only)
	if m == nil {
		return api.Ubuntu
	}

	// MasterProfile.Distro configured by defaults#setMasterNetworkDefaults
	return m.Distro
}

func getKubernetesSubnets(properties *api.Properties) string {
	subnetString := `{
            "name": "podCIDR%d",
            "properties": {
              "addressPrefix": "10.244.%d.0/24",
              "networkSecurityGroup": {
                "id": "[variables('nsgID')]"
              },
              "routeTable": {
                "id": "[variables('routeTableID')]"
              }
            }
          }`
	var buf bytes.Buffer

	cidrIndex := getKubernetesPodStartIndex(properties)
	for _, agentProfile := range properties.AgentPoolProfiles {
		if agentProfile.OSType == api.Windows {
			for i := 0; i < agentProfile.Count; i++ {
				buf.WriteString(",\n")
				buf.WriteString(fmt.Sprintf(subnetString, cidrIndex, cidrIndex))
				cidrIndex++
			}
		}
	}
	return buf.String()
}

func getKubernetesPodStartIndex(properties *api.Properties) int {
	nodeCount := 0
	nodeCount += properties.MasterProfile.Count
	for _, agentProfile := range properties.AgentPoolProfiles {
		if agentProfile.OSType != api.Windows {
			nodeCount += agentProfile.Count
		}
	}

	return nodeCount + 1
}

// getLinkedTemplatesForExtensions returns the
// Microsoft.Resources/deployments for each extension
//func getLinkedTemplatesForExtensions(properties api.Properties) string {
func getLinkedTemplatesForExtensions(properties *api.Properties) string {
	var result string

	extensions := properties.ExtensionProfiles
	masterProfileExtensions := properties.MasterProfile.Extensions
	orchestratorType := properties.OrchestratorProfile.OrchestratorType

	for err, extensionProfile := range extensions {
		_ = err

		masterOptedForExtension, singleOrAll := validateProfileOptedForExtension(extensionProfile.Name, masterProfileExtensions)
		if masterOptedForExtension {
			result += ","
			dta, e := getMasterLinkedTemplateText(properties.MasterProfile, orchestratorType, extensionProfile, singleOrAll)
			if e != nil {
				fmt.Println(e.Error())
				return ""
			}
			result += dta
		}

		for _, agentPoolProfile := range properties.AgentPoolProfiles {
			poolProfileExtensions := agentPoolProfile.Extensions
			poolOptedForExtension, singleOrAll := validateProfileOptedForExtension(extensionProfile.Name, poolProfileExtensions)
			if poolOptedForExtension {
				result += ","
				dta, e := getAgentPoolLinkedTemplateText(agentPoolProfile, orchestratorType, extensionProfile, singleOrAll)
				if e != nil {
					fmt.Println(e.Error())
					return ""
				}
				result += dta
			}

		}
	}

	return result
}

func getMasterLinkedTemplateText(masterProfile *api.MasterProfile, orchestratorType string, extensionProfile *api.ExtensionProfile, singleOrAll string) (string, error) {
	extTargetVMNamePrefix := "variables('masterVMNamePrefix')"

	loopCount := "[variables('masterCount')]"
	loopOffset := ""
	if orchestratorType == api.Kubernetes || orchestratorType == api.OpenShift {
		// Due to upgrade k8s sometimes needs to install just some of the nodes.
		loopCount = "[sub(variables('masterCount'), variables('masterOffset'))]"
		loopOffset = "variables('masterOffset')"
	}

	if strings.EqualFold(singleOrAll, "single") {
		loopCount = "1"
	}
	return internalGetPoolLinkedTemplateText(extTargetVMNamePrefix, orchestratorType, loopCount,
		loopOffset, extensionProfile)
}

func getAgentPoolLinkedTemplateText(agentPoolProfile *api.AgentPoolProfile, orchestratorType string, extensionProfile *api.ExtensionProfile, singleOrAll string) (string, error) {
	extTargetVMNamePrefix := fmt.Sprintf("variables('%sVMNamePrefix')", agentPoolProfile.Name)
	loopCount := fmt.Sprintf("[variables('%sCount'))]", agentPoolProfile.Name)
	loopOffset := ""

	// Availability sets can have an offset since we don't redeploy vms.
	// So we don't want to rerun these extensions in scale up scenarios.
	if agentPoolProfile.IsAvailabilitySets() {
		loopCount = fmt.Sprintf("[sub(variables('%sCount'), variables('%sOffset'))]",
			agentPoolProfile.Name, agentPoolProfile.Name)
		loopOffset = fmt.Sprintf("variables('%sOffset')", agentPoolProfile.Name)
	}

	if strings.EqualFold(singleOrAll, "single") {
		loopCount = "1"
	}

	return internalGetPoolLinkedTemplateText(extTargetVMNamePrefix, orchestratorType, loopCount,
		loopOffset, extensionProfile)
}

func internalGetPoolLinkedTemplateText(extTargetVMNamePrefix, orchestratorType, loopCount, loopOffset string, extensionProfile *api.ExtensionProfile) (string, error) {
	dta, e := getLinkedTemplateTextForURL(extensionProfile.RootURL, orchestratorType, extensionProfile.Name, extensionProfile.Version, extensionProfile.URLQuery)
	if e != nil {
		return "", e
	}
	if strings.Contains(extTargetVMNamePrefix, "master") {
		dta = strings.Replace(dta, "EXTENSION_TARGET_VM_TYPE", "master", -1)
	} else {
		dta = strings.Replace(dta, "EXTENSION_TARGET_VM_TYPE", "agent", -1)
	}
	extensionsParameterReference := fmt.Sprintf("[parameters('%sParameters')]", extensionProfile.Name)
	dta = strings.Replace(dta, "EXTENSION_PARAMETERS_REPLACE", extensionsParameterReference, -1)
	dta = strings.Replace(dta, "EXTENSION_URL_REPLACE", extensionProfile.RootURL, -1)
	dta = strings.Replace(dta, "EXTENSION_TARGET_VM_NAME_PREFIX", extTargetVMNamePrefix, -1)
	if _, err := strconv.Atoi(loopCount); err == nil {
		dta = strings.Replace(dta, "\"EXTENSION_LOOP_COUNT\"", loopCount, -1)
	} else {
		dta = strings.Replace(dta, "EXTENSION_LOOP_COUNT", loopCount, -1)
	}

	dta = strings.Replace(dta, "EXTENSION_LOOP_OFFSET", loopOffset, -1)
	return dta, nil
}

func validateProfileOptedForExtension(extensionName string, profileExtensions []api.Extension) (bool, string) {
	for _, extension := range profileExtensions {
		if extensionName == extension.Name {
			return true, extension.SingleOrAll
		}
	}
	return false, ""
}

// getLinkedTemplateTextForURL returns the string data from
// template-link.json in the following directory:
// extensionsRootURL/extensions/extensionName/version
// It returns an error if the extension cannot be found
// or loaded.  getLinkedTemplateTextForURL provides the ability
// to pass a root extensions url for testing
func getLinkedTemplateTextForURL(rootURL, orchestrator, extensionName, version, query string) (string, error) {
	supportsExtension, err := orchestratorSupportsExtension(rootURL, orchestrator, extensionName, version, query)
	if !supportsExtension {
		return "", errors.Wrap(err, "Extension not supported for orchestrator")
	}

	templateLinkBytes, err := getExtensionResource(rootURL, extensionName, version, "template-link.json", query)
	if err != nil {
		return "", err
	}

	return string(templateLinkBytes), nil
}

func orchestratorSupportsExtension(rootURL, orchestrator, extensionName, version, query string) (bool, error) {
	orchestratorBytes, err := getExtensionResource(rootURL, extensionName, version, "supported-orchestrators.json", query)
	if err != nil {
		return false, err
	}

	var supportedOrchestrators []string
	err = json.Unmarshal(orchestratorBytes, &supportedOrchestrators)
	if err != nil {
		return false, errors.Errorf("Unable to parse supported-orchestrators.json for Extension %s Version %s", extensionName, version)
	}

	if !stringInSlice(orchestrator, supportedOrchestrators) {
		return false, errors.Errorf("Orchestrator: %s not in list of supported orchestrators for Extension: %s Version %s", orchestrator, extensionName, version)
	}

	return true, nil
}

func getExtensionResource(rootURL, extensionName, version, fileName, query string) ([]byte, error) {
	requestURL := getExtensionURL(rootURL, extensionName, version, fileName, query)

	res, err := http.Get(requestURL)
	if err != nil {
		return nil, errors.Wrapf(err, "Unable to GET extension resource for extension: %s with version %s with filename %s at URL: %s", extensionName, version, fileName, requestURL)
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, errors.Errorf("Unable to GET extension resource for extension: %s with version %s with filename %s at URL: %s StatusCode: %s: Status: %s", extensionName, version, fileName, requestURL, strconv.Itoa(res.StatusCode), res.Status)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrapf(err, "Unable to GET extension resource for extension: %s with version %s  with filename %s at URL: %s", extensionName, version, fileName, requestURL)
	}

	return body, nil
}

func getExtensionURL(rootURL, extensionName, version, fileName, query string) string {
	extensionsDir := "extensions"
	url := rootURL + extensionsDir + "/" + extensionName + "/" + version + "/" + fileName
	if query != "" {
		url += "?" + query
	}
	return url
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func getSwarmVersions(orchestratorVersion, dockerComposeVersion string) string {
	return fmt.Sprintf("\"orchestratorVersion\": \"%s\",\n\"dockerComposeVersion\": \"%s\",\n", orchestratorVersion, dockerComposeVersion)
}

func getAddonByName(addons []api.KubernetesAddon, name string) api.KubernetesAddon {
	for i := range addons {
		if addons[i].Name == name {
			return addons[i]
		}
	}
	return api.KubernetesAddon{}
}
