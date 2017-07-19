package agentpool

import (
	"github.com/Azure/acs-engine/pkg/api/kubernetesagentpool"
	"text/template"
	//"github.com/go-playground/locales/cs"
)

func getTemplateFuncMap(agentPool *kubernetesagentpool.AgentPool) map[string]interface{} {
	return template.FuncMap{
		//"IsKubernetesVersionGe": func(version string) bool {
		//	targetVersion := api.OrchestratorVersion(version)
		//	targetVersionOrdinal := VersionOrdinal(targetVersion)
		//	orchestratorVersionOrdinal := VersionOrdinal(cs.Properties.OrchestratorProfile.OrchestratorVersion)
		//	return cs.Properties.OrchestratorProfile.OrchestratorType == api.Kubernetes &&
		//		orchestratorVersionOrdinal >= targetVersionOrdinal
		//},
		//"GetKubernetesLabels": func(profile *api.AgentPoolProfile) string {
		//	var buf bytes.Buffer
		//	buf.WriteString(fmt.Sprintf("role=agent,agentpool=%s", profile.Name))
		//	for k, v := range profile.CustomNodeLabels {
		//		buf.WriteString(fmt.Sprintf(",%s=%s", k, v))
		//	}
		//
		//	return buf.String()
		//},
		//"RequiresFakeAgentOutput": func() bool {
		//	return cs.Properties.OrchestratorProfile.OrchestratorType == api.Kubernetes
		//},
		//"IsSwarmMode": func() bool {
		//	return cs.Properties.OrchestratorProfile.IsSwarmMode()
		//},
		//"IsKubernetes": func() bool {
		//	return cs.Properties.OrchestratorProfile.IsKubernetes()
		//},
		//"IsPublic": func(ports []int) bool {
		//	return len(ports) > 0
		//},
		//"IsVNETIntegrated": func() bool {
		//	return cs.Properties.OrchestratorProfile.IsVNETIntegrated()
		//},
		//"GetVNETSubnetDependencies": func() string {
		//	return getVNETSubnetDependencies(cs.Properties)
		//},
		//"GetLBRules": func(name string, ports []int) string {
		//	return getLBRules(name, ports)
		//},
		//"GetProbes": func(ports []int) string {
		//	return getProbes(ports)
		//},
		//"GetSecurityRules": func(ports []int) string {
		//	return getSecurityRules(ports)
		//},
		//"GetUniqueNameSuffix": func() string {
		//	return GenerateClusterID(cs.Properties)
		//},
		//"GetVNETAddressPrefixes": func() string {
		//	return getVNETAddressPrefixes(cs.Properties)
		//},
		//"GetVNETSubnets": func(addNSG bool) string {
		//	return getVNETSubnets(cs.Properties, addNSG)
		//},
		//"GetDataDisks": func(profile *api.AgentPoolProfile) string {
		//	return getDataDisks(profile)
		//},
		//"GetDCOSMasterCustomData": func() string {
		//	masterProvisionScript := getDCOSMasterProvisionScript()
		//	masterAttributeContents := getDCOSMasterCustomNodeLabels()
		//	str := getSingleLineDCOSCustomData(cs.Properties.OrchestratorProfile.OrchestratorType, cs.Properties.OrchestratorProfile.OrchestratorVersion, cs.Properties.MasterProfile.Count, masterProvisionScript, masterAttributeContents)
		//
		//	return fmt.Sprintf("\"customData\": \"[base64(concat('#cloud-config\\n\\n', '%s'))]\",", str)
		//},
		//"GetDCOSAgentCustomData": func(profile *api.AgentPoolProfile) string {
		//	agentProvisionScript := getDCOSAgentProvisionScript(profile)
		//	attributeContents := getDCOSAgentCustomNodeLabels(profile)
		//	str := getSingleLineDCOSCustomData(cs.Properties.OrchestratorProfile.OrchestratorType, cs.Properties.OrchestratorProfile.OrchestratorVersion, cs.Properties.MasterProfile.Count, agentProvisionScript, attributeContents)
		//
		//	return fmt.Sprintf("\"customData\": \"[base64(concat('#cloud-config\\n\\n', '%s'))]\",", str)
		//},
		//"GetMasterAllowedSizes": func() string {
		//	if t.ClassicMode {
		//		return GetClassicAllowedSizes()
		//	} else if cs.Properties.OrchestratorProfile.OrchestratorType == api.DCOS {
		//		return GetDCOSMasterAllowedSizes()
		//	}
		//	return GetMasterAgentAllowedSizes()
		//},
		//"GetAgentAllowedSizes": func() string {
		//	if t.ClassicMode {
		//		return GetClassicAllowedSizes()
		//	} else if cs.Properties.OrchestratorProfile.OrchestratorType == api.Kubernetes {
		//		return GetKubernetesAgentAllowedSizes()
		//	}
		//	return GetMasterAgentAllowedSizes()
		//},
		//"GetSizeMap": func() string {
		//	if t.ClassicMode {
		//		return GetClassicSizeMap()
		//	}
		//	return GetSizeMap()
		//},
		//"GetClassicMode": func() bool {
		//	return t.ClassicMode
		//},
		//"Base64": func(s string) string {
		//	return base64.StdEncoding.EncodeToString([]byte(s))
		//},
		//"GetDefaultInternalLbStaticIPOffset": func() int {
		//	return DefaultInternalLbStaticIPOffset
		//},
		//"GetKubernetesMasterCustomScript": func() string {
		//	return getBase64CustomScript(kubernetesMasterCustomScript)
		//},
		//"GetKubernetesMasterCustomData": func(profile *api.Properties) string {
		//	str, e := t.getSingleLineForTemplate(kubernetesMasterCustomDataYaml, cs, profile)
		//	if e != nil {
		//		return ""
		//	}
		//
		//	for placeholder, filename := range kubernetesManifestYamls {
		//		manifestTextContents := getBase64CustomScript(filename)
		//		str = strings.Replace(str, placeholder, manifestTextContents, -1)
		//	}
		//
		//	// add artifacts and addons
		//	for placeholder, filename := range kubernetesAritfacts {
		//		addonTextContents := getBase64CustomScript(filename)
		//		str = strings.Replace(str, placeholder, addonTextContents, -1)
		//	}
		//
		//	var addonYamls map[string]string
		//	if profile.OrchestratorProfile.OrchestratorVersion == api.Kubernetes153 ||
		//		profile.OrchestratorProfile.OrchestratorVersion == api.Kubernetes157 {
		//		addonYamls = kubernetesAddonYamls15
		//	} else {
		//		addonYamls = kubernetesAddonYamls
		//	}
		//	for placeholder, filename := range addonYamls {
		//		addonTextContents := getBase64CustomScript(filename)
		//		str = strings.Replace(str, placeholder, addonTextContents, -1)
		//	}
		//
		//	// add calico manifests
		//	if profile.OrchestratorProfile.KubernetesConfig.NetworkPolicy == "calico" {
		//		for placeholder, filename := range calicoAddonYamls {
		//			addonTextContents := getBase64CustomScript(filename)
		//			str = strings.Replace(str, placeholder, addonTextContents, -1)
		//		}
		//	}
		//
		//	// return the custom data
		//	return fmt.Sprintf("\"customData\": \"[base64(concat('%s'))]\",", str)
		//},
		//"GetKubernetesAgentCustomData": func(profile *api.AgentPoolProfile) string {
		//	str, e := t.getSingleLineForTemplate(kubernetesAgentCustomDataYaml, cs, profile)
		//	if e != nil {
		//		return ""
		//	}
		//
		//	// add artifacts
		//	for placeholder, filename := range kubernetesAritfacts {
		//		addonTextContents := getBase64CustomScript(filename)
		//		str = strings.Replace(str, placeholder, addonTextContents, -1)
		//	}
		//
		//	return fmt.Sprintf("\"customData\": \"[base64(concat('%s'))]\",", str)
		//},
		//"GetKubernetesB64Provision": func() string {
		//	return getBase64CustomScript(kubernetesMasterCustomScript)
		//},
		//"GetMasterSwarmCustomData": func() string {
		//	files := []string{swarmProvision}
		//	str := buildYamlFileWithWriteFiles(files)
		//	str = escapeSingleLine(str)
		//	return fmt.Sprintf("\"customData\": \"[base64('%s')]\",", str)
		//},
		//"GetAgentSwarmCustomData": func() string {
		//	files := []string{swarmProvision}
		//	str := buildYamlFileWithWriteFiles(files)
		//	str = escapeSingleLine(str)
		//	return fmt.Sprintf("\"customData\": \"[base64(concat('%s',variables('agentRunCmdFile'),variables('agentRunCmd')))]\",", str)
		//},
		//"GetLocation": func() string {
		//	return cs.Location
		//},
		//"GetWinAgentSwarmCustomData": func() string {
		//	str := getBase64CustomScript(swarmWindowsProvision)
		//	return fmt.Sprintf("\"customData\": \"%s\"", str)
		//},
		//"GetWinAgentSwarmModeCustomData": func() string {
		//	str := getBase64CustomScript(swarmModeWindowsProvision)
		//	return fmt.Sprintf("\"customData\": \"%s\"", str)
		//},
		//"GetKubernetesWindowsAgentCustomData": func(profile *api.AgentPoolProfile) string {
		//	str, e := t.getSingleLineForTemplate(kubernetesWindowsAgentCustomDataPS1, cs, profile)
		//	if e != nil {
		//		return ""
		//	}
		//	return fmt.Sprintf("\"customData\": \"[base64(concat('%s'))]\",", str)
		//},
		//"GetMasterSwarmModeCustomData": func() string {
		//	files := []string{swarmModeProvision}
		//	str := buildYamlFileWithWriteFiles(files)
		//	str = escapeSingleLine(str)
		//	return fmt.Sprintf("\"customData\": \"[base64('%s')]\",", str)
		//},
		//"GetAgentSwarmModeCustomData": func() string {
		//	files := []string{swarmModeProvision}
		//	str := buildYamlFileWithWriteFiles(files)
		//	str = escapeSingleLine(str)
		//	return fmt.Sprintf("\"customData\": \"[base64(concat('%s',variables('agentRunCmdFile'),variables('agentRunCmd')))]\",", str)
		//},
		//"GetKubernetesSubnets": func() string {
		//	return getKubernetesSubnets(cs.Properties)
		//},
		//"GetKubernetesPodStartIndex": func() string {
		//	return fmt.Sprintf("%d", getKubernetesPodStartIndex(cs.Properties))
		//},
		//"WrapAsVariable": func(s string) string {
		//	return fmt.Sprintf("',variables('%s'),'", s)
		//},
		//"WrapAsVerbatim": func(s string) string {
		//	return fmt.Sprintf("',%s,'", s)
		//},
		//"AnyAgentUsesAvailablilitySets": func() bool {
		//	for _, agentProfile := range cs.Properties.AgentPoolProfiles {
		//		if agentProfile.IsAvailabilitySets() {
		//			return true
		//		}
		//	}
		//	return false
		//},
		//"HasLinuxAgents": func() bool {
		//	for _, agentProfile := range cs.Properties.AgentPoolProfiles {
		//		if agentProfile.IsLinux() {
		//			return true
		//		}
		//	}
		//	return false
		//},
		//"HasLinuxSecrets": func() bool {
		//	return cs.Properties.LinuxProfile.HasSecrets()
		//},
		//"HasWindowsSecrets": func() bool {
		//	return cs.Properties.WindowsProfile.HasSecrets()
		//},
		//"PopulateClassicModeDefaultValue": func(attr string) string {
		//	var val string
		//	if !t.ClassicMode {
		//		val = ""
		//	} else {
		//		kubernetesVersion := cs.Properties.OrchestratorProfile.OrchestratorVersion
		//		cloudSpecConfig := GetCloudSpecConfig(cs.Location)
		//		switch attr {
		//		case "kubernetesHyperkubeSpec":
		//			val = cs.Properties.OrchestratorProfile.KubernetesConfig.KubernetesImageBase + KubeImages[kubernetesVersion]["hyperkube"]
		//		case "kubernetesAddonManagerSpec":
		//			val = cloudSpecConfig.KubernetesSpecConfig.KubernetesImageBase + KubeImages[kubernetesVersion]["addonmanager"]
		//		case "kubernetesAddonResizerSpec":
		//			val = cloudSpecConfig.KubernetesSpecConfig.KubernetesImageBase + KubeImages[kubernetesVersion]["addonresizer"]
		//		case "kubernetesDashboardSpec":
		//			val = cloudSpecConfig.KubernetesSpecConfig.KubernetesImageBase + KubeImages[kubernetesVersion]["dashboard"]
		//		case "kubernetesDNSMasqSpec":
		//			val = cloudSpecConfig.KubernetesSpecConfig.KubernetesImageBase + KubeImages[kubernetesVersion]["dnsmasq"]
		//		case "kubernetesExecHealthzSpec":
		//			val = cloudSpecConfig.KubernetesSpecConfig.KubernetesImageBase + KubeImages[kubernetesVersion]["exechealthz"]
		//		case "kubernetesHeapsterSpec":
		//			val = cloudSpecConfig.KubernetesSpecConfig.KubernetesImageBase + KubeImages[kubernetesVersion]["heapster"]
		//		case "kubernetesKubeDNSSpec":
		//			val = cloudSpecConfig.KubernetesSpecConfig.KubernetesImageBase + KubeImages[kubernetesVersion]["dns"]
		//		case "kubernetesPodInfraContainerSpec":
		//			val = cloudSpecConfig.KubernetesSpecConfig.KubernetesImageBase + KubeImages[kubernetesVersion]["pause"]
		//		case "kubernetesNodeStatusUpdateFrequency":
		//			val = KubeImages[kubernetesVersion]["nodestatusfreq"]
		//		case "kubernetesCtrlMgrNodeMonitorGracePeriod":
		//			val = KubeImages[kubernetesVersion]["nodegraceperiod"]
		//		case "kubernetesCtrlMgrPodEvictionTimeout":
		//			val = KubeImages[kubernetesVersion]["podeviction"]
		//		case "kubernetesCtrlMgrRouteReconciliationPeriod":
		//			val = KubeImages[kubernetesVersion]["routeperiod"]
		//		case "cloudProviderBackoff":
		//			val = KubeImages[kubernetesVersion]["backoff"]
		//		case "cloudProviderBackoffRetries":
		//			val = KubeImages[kubernetesVersion]["backoffretries"]
		//		case "cloudProviderBackoffExponent":
		//			val = KubeImages[kubernetesVersion]["backoffexponent"]
		//		case "cloudProviderBackoffDuration":
		//			val = KubeImages[kubernetesVersion]["backoffduration"]
		//		case "cloudProviderBackoffJitter":
		//			val = KubeImages[kubernetesVersion]["backoffjitter"]
		//		case "cloudProviderRatelimit":
		//			val = KubeImages[kubernetesVersion]["ratelimit"]
		//		case "cloudProviderRatelimitQPS":
		//			val = KubeImages[kubernetesVersion]["ratelimitqps"]
		//		case "cloudProviderRatelimitBucket":
		//			val = KubeImages[kubernetesVersion]["ratelimitbucket"]
		//		case "kubeBinariesSASURL":
		//			val = cloudSpecConfig.KubernetesSpecConfig.KubeBinariesSASURLBase + KubeImages[kubernetesVersion]["windowszip"]
		//		case "kubeClusterCidr":
		//			val = "10.244.0.0/16"
		//		case "kubeBinariesVersion":
		//			val = string(api.KubernetesLatest)
		//		case "caPrivateKey":
		//			// The base64 encoded "NotAvailable"
		//			val = "Tm90QXZhaWxhYmxlCg=="
		//		case "dockerBridgeCidr":
		//			val = DefaultDockerBridgeSubnet
		//		default:
		//			val = ""
		//		}
		//	}
		//	return fmt.Sprintf("\"defaultValue\": \"%s\",", val)
		//},
		//// inspired by http://stackoverflow.com/questions/18276173/calling-a-template-with-several-pipeline-parameters/18276968#18276968
		//"dict": func(values ...interface{}) (map[string]interface{}, error) {
		//	if len(values)%2 != 0 {
		//		return nil, errors.New("invalid dict call")
		//	}
		//	dict := make(map[string]interface{}, len(values)/2)
		//	for i := 0; i < len(values); i += 2 {
		//		key, ok := values[i].(string)
		//		if !ok {
		//			return nil, errors.New("dict keys must be strings")
		//		}
		//		dict[key] = values[i+1]
		//	}
		//	return dict, nil
		//},
		//"loop": func(min, max int) []int {
		//	var s []int
		//	for i := min; i <= max; i++ {
		//		s = append(s, i)
		//	}
		//	return s
		//},
	}
}
