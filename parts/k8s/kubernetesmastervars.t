    "etcdDiskSizeGB": "[parameters('etcdDiskSizeGB')]",
    "etcdDownloadURLBase": "[parameters('etcdDownloadURLBase')]",
    "etcdVersion": "[parameters('etcdVersion')]",
    "maxVMsPerPool": 100,
{{ if not IsOpenShift }}
    "apiServerCertificate": "[parameters('apiServerCertificate')]",
{{ end }}
{{ if IsOpenShift }}
    "routerNSGName": "[concat(variables('orchestratorName'), '-router-', variables('nameSuffix'), '-nsg')]",
    "routerNSGID": "[resourceId('Microsoft.Network/networkSecurityGroups', variables('routerNSGName'))]",
    "routerIPName": "[concat(variables('orchestratorName'), '-router-ip-', variables('masterFqdnPrefix'), '-', variables('nameSuffix'))]",
    "routerLBName": "[concat(variables('orchestratorName'), '-router-lb-', variables('nameSuffix'))]",
    "routerLBID": "[resourceId('Microsoft.Network/loadBalancers', variables('routerLBName'))]",
{{ end }}
{{ if not IsHostedMaster }}
{{ if not IsOpenShift }}
    "apiServerPrivateKey": "[parameters('apiServerPrivateKey')]",
    "etcdServerCertificate": "[parameters('etcdServerCertificate')]",
    "etcdServerPrivateKey": "[parameters('etcdServerPrivateKey')]",
    "etcdClientPrivateKey": "[parameters('etcdClientPrivateKey')]",
    "etcdClientCertificate": "[parameters('etcdClientCertificate')]",
    {{if eq .MasterProfile.Count 1}}
    "etcdPeerPrivateKeys": [
        "[parameters('etcdPeerPrivateKey0')]"
    ],
    "etcdPeerCertificates": [
        "[parameters('etcdPeerCertificate0')]"
    ],
    {{end}}
    {{if eq .MasterProfile.Count 3}}
    "etcdPeerPrivateKeys": [
        "[parameters('etcdPeerPrivateKey0')]",
        "[parameters('etcdPeerPrivateKey1')]",
        "[parameters('etcdPeerPrivateKey2')]"
    ],
    "etcdPeerCertificates": [
        "[parameters('etcdPeerCertificate0')]",
        "[parameters('etcdPeerCertificate1')]",
        "[parameters('etcdPeerCertificate2')]"
    ],
    {{end}}
    {{if eq .MasterProfile.Count 5}}
    "etcdPeerPrivateKeys": [
        "[parameters('etcdPeerPrivateKey0')]",
        "[parameters('etcdPeerPrivateKey1')]",
        "[parameters('etcdPeerPrivateKey2')]",
        "[parameters('etcdPeerPrivateKey3')]",
        "[parameters('etcdPeerPrivateKey4')]"
    ],
    "etcdPeerCertificates": [
        "[parameters('etcdPeerCertificate0')]",
        "[parameters('etcdPeerCertificate1')]",
        "[parameters('etcdPeerCertificate2')]",
        "[parameters('etcdPeerCertificate3')]",
        "[parameters('etcdPeerCertificate4')]"
    ],
    {{end}}
    "etcdPeerCertFilepath":[
        "/etc/kubernetes/certs/etcdpeer0.crt",
        "/etc/kubernetes/certs/etcdpeer1.crt",
        "/etc/kubernetes/certs/etcdpeer2.crt",
        "/etc/kubernetes/certs/etcdpeer3.crt",
        "/etc/kubernetes/certs/etcdpeer4.crt"
    ],
    "etcdPeerKeyFilepath":[
        "/etc/kubernetes/certs/etcdpeer0.key",
        "/etc/kubernetes/certs/etcdpeer1.key",
        "/etc/kubernetes/certs/etcdpeer2.key",
        "/etc/kubernetes/certs/etcdpeer3.key",
        "/etc/kubernetes/certs/etcdpeer4.key"
    ],
    "etcdCaFilepath": "/etc/kubernetes/certs/ca.crt",
    "etcdClientCertFilepath": "/etc/kubernetes/certs/etcdclient.crt",
    "etcdClientKeyFilepath": "/etc/kubernetes/certs/etcdclient.key",
    "etcdServerCertFilepath": "/etc/kubernetes/certs/etcdserver.crt",
    "etcdServerKeyFilepath": "/etc/kubernetes/certs/etcdserver.key",
{{end}}
{{end}}
{{ if not IsOpenShift }}
    "caCertificate": "[parameters('caCertificate')]",
    "caPrivateKey": "[parameters('caPrivateKey')]",
    "clientCertificate": "[parameters('clientCertificate')]",
    "clientPrivateKey": "[parameters('clientPrivateKey')]",
    "kubeConfigCertificate": "[parameters('kubeConfigCertificate')]",
    "kubeConfigPrivateKey": "[parameters('kubeConfigPrivateKey')]",
{{end}}
    "kubernetesHyperkubeSpec": "[parameters('kubernetesHyperkubeSpec')]",
    "kubernetesCcmImageSpec": "[parameters('kubernetesCcmImageSpec')]",
    "kubernetesAddonManagerSpec": "[parameters('kubernetesAddonManagerSpec')]",
    "kubernetesAddonResizerSpec": "[parameters('kubernetesAddonResizerSpec')]",
    "kubernetesDashboardSpec": "[parameters('kubernetesDashboardSpec')]",
    "kubernetesDashboardCPURequests": "[parameters('kubernetesDashboardCPURequests')]",
    "kubernetesDashboardMemoryRequests": "[parameters('kubernetesDashboardMemoryRequests')]",
    "kubernetesDashboardCPULimit": "[parameters('kubernetesDashboardCPULimit')]",
    "kubernetesDashboardMemoryLimit": "[parameters('kubernetesDashboardMemoryLimit')]",
    "kubernetesExecHealthzSpec": "[parameters('kubernetesExecHealthzSpec')]",
    "kubernetesDNSSidecarSpec": "[parameters('kubernetesDNSSidecarSpec')]",
    "kubernetesHeapsterSpec": "[parameters('kubernetesHeapsterSpec')]",
    "kubernetesMetricsServerSpec": "[parameters('kubernetesMetricsServerSpec')]",
    "kubernetesNVIDIADevicePluginSpec": "[parameters('kubernetesNVIDIADevicePluginSpec')]",
    "kubernetesNVIDIADevicePluginCPURequests": "[parameters('kubernetesNVIDIADevicePluginCPURequests')]",
    "kubernetesNVIDIADevicePluginMemoryRequests": "[parameters('kubernetesNVIDIADevicePluginMemoryRequests')]",
    "kubernetesNVIDIADevicePluginCPULimit": "[parameters('kubernetesNVIDIADevicePluginCPULimit')]",
    "kubernetesNVIDIADevicePluginMemoryLimit": "[parameters('kubernetesNVIDIADevicePluginMemoryLimit')]",
    "kubernetesTillerSpec": "[parameters('kubernetesTillerSpec')]",
    "kubernetesTillerCPURequests": "[parameters('kubernetesTillerCPURequests')]",
    "kubernetesTillerMemoryRequests": "[parameters('kubernetesTillerMemoryRequests')]",
    "kubernetesTillerCPULimit": "[parameters('kubernetesTillerCPULimit')]",
    "kubernetesTillerMemoryLimit": "[parameters('kubernetesTillerMemoryLimit')]",
    "kubernetesTillerMaxHistory": "[parameters('kubernetesTillerMaxHistory')]",
    "kubernetesAADPodIdentityEnabled": "[parameters('kubernetesAADPodIdentityEnabled')]",
    "kubernetesACIConnectorEnabled": "[parameters('kubernetesACIConnectorEnabled')]",
    "kubernetesACIConnectorSpec": "[parameters('kubernetesACIConnectorSpec')]",
    "kubernetesACIConnectorNodeName": "[parameters('kubernetesACIConnectorNodeName')]",
    "kubernetesACIConnectorOS": "[parameters('kubernetesACIConnectorOS')]",
    "kubernetesACIConnectorTaint": "[parameters('kubernetesACIConnectorTaint')]",
    "kubernetesACIConnectorRegion": "[parameters('kubernetesACIConnectorRegion')]",
    "kubernetesACIConnectorCPURequests": "[parameters('kubernetesACIConnectorCPURequests')]",
    "kubernetesACIConnectorMemoryRequests": "[parameters('kubernetesACIConnectorMemoryRequests')]",
    "kubernetesACIConnectorCPULimit": "[parameters('kubernetesACIConnectorCPULimit')]",
    "kubernetesACIConnectorMemoryLimit": "[parameters('kubernetesACIConnectorMemoryLimit')]",
    "kubernetesClusterAutoscalerSpec": "[parameters('kubernetesClusterAutoscalerSpec')]",
    "kubernetesClusterAutoscalerAzureCloud": "[parameters('kubernetesClusterAutoscalerAzureCloud')]",
    "kubernetesClusterAutoscalerCPULimit": "[parameters('kubernetesClusterAutoscalerCPULimit')]",
    "kubernetesClusterAutoscalerMemoryLimit": "[parameters('kubernetesClusterAutoscalerMemoryLimit')]",
    "kubernetesClusterAutoscalerCPURequests": "[parameters('kubernetesClusterAutoscalerCPURequests')]",
    "kubernetesClusterAutoscalerMemoryRequests": "[parameters('kubernetesClusterAutoscalerMemoryRequests')]",
    "kubernetesClusterAutoscalerMinNodes": "[parameters('kubernetesClusterAutoscalerMinNodes')]",
    "kubernetesClusterAutoscalerMaxNodes": "[parameters('kubernetesClusterAutoscalerMaxNodes')]",
    "kubernetesClusterAutoscalerEnabled": "[parameters('kubernetesClusterAutoscalerEnabled')]",
    "kubernetesClusterAutoscalerUseManagedIdentity": "[parameters('kubernetesClusterAutoscalerUseManagedIdentity')]",
    "kubernetesKeyVaultFlexVolumeInstallerCPURequests": "[parameters('kubernetesKeyVaultFlexVolumeInstallerCPURequests')]",
    "kubernetesKeyVaultFlexVolumeInstallerMemoryRequests": "[parameters('kubernetesKeyVaultFlexVolumeInstallerMemoryRequests')]",
    "kubernetesKeyVaultFlexVolumeInstallerCPULimit": "[parameters('kubernetesKeyVaultFlexVolumeInstallerCPULimit')]",
    "kubernetesKeyVaultFlexVolumeInstallerMemoryLimit": "[parameters('kubernetesKeyVaultFlexVolumeInstallerMemoryLimit')]",
    "kubernetesReschedulerSpec": "[parameters('kubernetesReschedulerSpec')]",
    "kubernetesReschedulerCPURequests": "[parameters('kubernetesReschedulerCPURequests')]",
    "kubernetesReschedulerMemoryRequests": "[parameters('kubernetesReschedulerMemoryRequests')]",
    "kubernetesReschedulerCPULimit": "[parameters('kubernetesReschedulerCPULimit')]",
    "kubernetesReschedulerMemoryLimit": "[parameters('kubernetesReschedulerMemoryLimit')]",
    "kubernetesPodInfraContainerSpec": "[parameters('kubernetesPodInfraContainerSpec')]",
    "cloudProviderBackoff": "[parameters('cloudProviderBackoff')]",
    "cloudProviderBackoffRetries": "[parameters('cloudProviderBackoffRetries')]",
    "cloudProviderBackoffExponent": "[parameters('cloudProviderBackoffExponent')]",
    "cloudProviderBackoffDuration": "[parameters('cloudProviderBackoffDuration')]",
    "cloudProviderBackoffJitter": "[parameters('cloudProviderBackoffJitter')]",
    "cloudProviderRatelimit": "[parameters('cloudProviderRatelimit')]",
    "cloudProviderRatelimitQPS": "[parameters('cloudProviderRatelimitQPS')]",
    "cloudProviderRatelimitBucket": "[parameters('cloudProviderRatelimitBucket')]",
    "useManagedIdentityExtension": "{{ UseManagedIdentity }}",
    "useInstanceMetadata": "{{ UseInstanceMetadata }}",
    "kubernetesKubeDNSSpec": "[parameters('kubernetesKubeDNSSpec')]",
    "kubernetesDNSMasqSpec": "[parameters('kubernetesDNSMasqSpec')]",
    "networkPlugin": "[parameters('networkPlugin')]",
    "networkPolicy": "[parameters('networkPolicy')]",
    "containerRuntime": "[parameters('containerRuntime')]",
    "containerdDownloadURLBase": "[parameters('containerdDownloadURLBase')]",
    "cniPluginsURL":"[parameters('cniPluginsURL')]",
    "vnetCniLinuxPluginsURL":"[parameters('vnetCniLinuxPluginsURL')]",
    "vnetCniWindowsPluginsURL":"[parameters('vnetCniWindowsPluginsURL')]",
    "kubernetesNonMasqueradeCidr": "[parameters('kubernetesNonMasqueradeCidr')]",
    "kubernetesKubeletClusterDomain": "[parameters('kubernetesKubeletClusterDomain')]",
    "maxPods": "[parameters('maxPods')]",
    "vnetCidr": "[parameters('vnetCidr')]",
    "gcHighThreshold":"[parameters('gcHighThreshold')]",
    "gcLowThreshold":"[parameters('gcLowThreshold')]",
    "omsAgentVersion": "[parameters('omsAgentVersion')]",
    "omsAgentDockerProviderVersion": "[parameters('omsAgentDockerProviderVersion')]",
    "omsAgentImage": "[parameters('omsAgentImage')]",
    "omsAgentWorkspaceGuid": "[parameters('omsAgentWorkspaceGuid')]",
    "omsAgentWorkspaceKey": "[parameters('omsAgentWorkspaceKey')]",
    "kubernetesOMSAgentCPURequests": "[parameters('kubernetesOMSAgentCPURequests')]",
    "kubernetesOMSAgentMemoryRequests": "[parameters('kubernetesOMSAgentMemoryRequests')]",
    "kubernetesOMSAgentCPULimit": "[parameters('kubernetesOMSAgentCPULimit')]",
    "kubernetesOMSAgentMemoryLimit": "[parameters('kubernetesOMSAgentMemoryLimit')]",
{{if EnableDataEncryptionAtRest}}
    "etcdEncryptionKey": "[parameters('etcdEncryptionKey')]",
{{end}}
{{ if UseManagedIdentity }}
    "servicePrincipalClientId": "msi",
    "servicePrincipalClientSecret": "msi",
{{ else }}
    "servicePrincipalClientId": "[parameters('servicePrincipalClientId')]",
    "servicePrincipalClientSecret": "[parameters('servicePrincipalClientSecret')]",
{{ end }}
    "username": "[parameters('linuxAdminUsername')]",
    "masterFqdnPrefix": "[tolower(parameters('masterEndpointDNSNamePrefix'))]",
{{ if not IsHostedMaster }}
    "masterPrivateIp": "[parameters('firstConsecutiveStaticIP')]",
    "masterVMSize": "[parameters('masterVMSize')]",
{{end}}
    "sshPublicKeyData": "[parameters('sshRSAPublicKey')]",
{{if .HasAadProfile}}
    "aadTenantId": "[parameters('aadTenantId')]",
    "aadAdminGroupId": "[parameters('aadAdminGroupId')]",
{{end}}
{{if not IsHostedMaster}}
  {{if GetClassicMode}}
    "masterCount": "[parameters('masterCount')]",
  {{else}}
    "masterCount": {{.MasterProfile.Count}},
  {{end}}
    "masterOffset": "[parameters('masterOffset')]",
{{end}}
    "apiVersionDefault": "2016-03-30",
    "apiVersionAcceleratedNetworking": "2018-04-01",
    "apiVersionLinkDefault": "2015-01-01",
    "locations": [
         "[resourceGroup().location]",
         "[parameters('location')]"
    ],
    "location": "[variables('locations')[mod(add(2,length(parameters('location'))),add(1,length(parameters('location'))))]]",
    "masterAvailabilitySet": "[concat('master-availabilityset-', variables('nameSuffix'))]",
    "nameSuffix": "[parameters('nameSuffix')]",
    "orchestratorName": "[parameters('orchestratorName')]",
    "generatorCode": "[parameters('generatorCode')]",
    "fqdnEndpointSuffix":"[parameters('fqdnEndpointSuffix')]",
    "osImageOffer": "[parameters('osImageOffer')]",
    "osImagePublisher": "[parameters('osImagePublisher')]",
    "osImageSKU": "[parameters('osImageSKU')]",
    "osImageVersion": "[parameters('osImageVersion')]",
    "osImageName": "[parameters('osImageName')]",
    "osImageResourceGroup": "[parameters('osImageResourceGroup')]",
    "resourceGroup": "[resourceGroup().name]",
    "truncatedResourceGroup": "[take(replace(replace(resourceGroup().name, '(', '-'), ')', '-'), 63)]",
    "labelResourceGroup": "[if(or(or(endsWith(variables('truncatedResourceGroup'), '-'), endsWith(variables('truncatedResourceGroup'), '_')), endsWith(variables('truncatedResourceGroup'), '.')), concat(take(variables('truncatedResourceGroup'), 62), 'z'), variables('truncatedResourceGroup'))]",
{{if not IsHostedMaster}}
    "routeTableName": "[concat(variables('masterVMNamePrefix'),'routetable')]",
{{else}}
    "routeTableName": "[concat(variables('agentNamePrefix'), 'routetable')]",
{{end}}
    "routeTableID": "[resourceId('Microsoft.Network/routeTables', variables('routeTableName'))]",
    "sshNatPorts": [22,2201,2202,2203,2204],
    "sshKeyPath": "[concat('/home/',variables('username'),'/.ssh/authorized_keys')]",

{{if .HasStorageAccountDisks}}
    "apiVersionStorage": "2015-06-15",
    "maxVMsPerStorageAccount": 20,
    "maxStorageAccountsPerAgent": "[div(variables('maxVMsPerPool'),variables('maxVMsPerStorageAccount'))]",
    "dataStorageAccountPrefixSeed": 97,
    "storageAccountPrefixes": [ "0", "6", "c", "i", "o", "u", "1", "7", "d", "j", "p", "v", "2", "8", "e", "k", "q", "w", "3", "9", "f", "l", "r", "x", "4", "a", "g", "m", "s", "y", "5", "b", "h", "n", "t", "z" ],
    "storageAccountPrefixesCount": "[length(variables('storageAccountPrefixes'))]",
    "vmsPerStorageAccount": 20,
    "storageAccountBaseName": "[uniqueString(concat(variables('masterFqdnPrefix'),variables('location')))]",
    {{GetSizeMap}},
{{else}}
    "storageAccountPrefixes": [],
    "storageAccountBaseName": "",
{{end}}
{{if .HasManagedDisks}}
    "apiVersionStorageManagedDisks": "2016-04-30-preview",
{{end}}
{{if .HasVirtualMachineScaleSets}}
    "apiVersionVirtualMachineScaleSets": "2017-12-01",
{{end}}
{{if not IsHostedMaster}}
  {{if .MasterProfile.IsStorageAccount}}
    "masterStorageAccountName": "[concat(variables('storageAccountBaseName'), 'mstr0')]",
  {{end}}
{{end}}
    "provisionScript": "{{GetKubernetesB64Provision}}",
    "provisionSource": "{{GetKubernetesB64ProvisionSource}}",
    "mountetcdScript": "{{GetKubernetesB64Mountetcd}}",
    "customSearchDomainsScript": "{{GetKubernetesB64CustomSearchDomainsScript}}",
{{if not IsOpenShift}}
    {{if not IsHostedMaster}}
    "provisionScriptParametersMaster": "[concat('MASTER_VM_NAME=',variables('masterVMNames')[variables('masterOffset')],' ETCD_PEER_URL=',variables('masterEtcdPeerURLs')[variables('masterOffset')],' ETCD_CLIENT_URL=',variables('masterEtcdClientURLs')[variables('masterOffset')],' MASTER_NODE=true CLUSTER_AUTOSCALER_ADDON=',variables('kubernetesClusterAutoscalerEnabled'),' ACI_CONNECTOR_ADDON=',variables('kubernetesACIConnectorEnabled'),' APISERVER_PRIVATE_KEY=',variables('apiServerPrivateKey'),' CA_CERTIFICATE=',variables('caCertificate'),' CA_PRIVATE_KEY=',variables('caPrivateKey'),' MASTER_FQDN=',variables('masterFqdnPrefix'),' KUBECONFIG_CERTIFICATE=',variables('kubeConfigCertificate'),' KUBECONFIG_KEY=',variables('kubeConfigPrivateKey'),' ETCD_SERVER_CERTIFICATE=',variables('etcdServerCertificate'),' ETCD_CLIENT_CERTIFICATE=',variables('etcdClientCertificate'),' ETCD_SERVER_PRIVATE_KEY=',variables('etcdServerPrivateKey'),' ETCD_CLIENT_PRIVATE_KEY=',variables('etcdClientPrivateKey'),' ETCD_PEER_CERTIFICATES=',string(variables('etcdPeerCertificates')),' ETCD_PEER_PRIVATE_KEYS=',string(variables('etcdPeerPrivateKeys')))]",
        {{if EnableEncryptionWithExternalKms}}
            {{ if not UseManagedIdentity}}
    "servicePrincipalObjectId": "[parameters('servicePrincipalObjectId')]",
            {{end}}
    "provisionScriptParametersCommon": "[concat('ADMINUSER=',variables('username'),' DOCKER_ENGINE_VERSION=',variables('dockerEngineVersion'),' DOCKER_REPO=',variables('dockerEngineDownloadRepo'),' TENANT_ID=',variables('tenantID'),' HYPERKUBE_URL=',variables('kubernetesHyperkubeSpec'),' APISERVER_PUBLIC_KEY=',variables('apiserverCertificate'),' SUBSCRIPTION_ID=',variables('subscriptionId'),' RESOURCE_GROUP=',variables('resourceGroup'),' LOCATION=',variables('location'),' VM_TYPE=',variables('vmType'),' SUBNET=',variables('subnetName'),' NETWORK_SECURITY_GROUP=',variables('nsgName'),' VIRTUAL_NETWORK=',variables('virtualNetworkName'),' VIRTUAL_NETWORK_RESOURCE_GROUP=',variables('virtualNetworkResourceGroupName'),' ROUTE_TABLE=',variables('routeTableName'),' PRIMARY_AVAILABILITY_SET=',variables('primaryAvailabilitySetName'),' PRIMARY_SCALE_SET=',variables('primaryScaleSetName'),' SERVICE_PRINCIPAL_CLIENT_ID=',variables('servicePrincipalClientId'),' SERVICE_PRINCIPAL_CLIENT_SECRET=',variables('singleQuote'),variables('servicePrincipalClientSecret'),variables('singleQuote'),' KUBELET_PRIVATE_KEY=',variables('clientPrivateKey'),' TARGET_ENVIRONMENT=',variables('targetEnvironment'),' NETWORK_PLUGIN=',variables('networkPlugin'),' VNET_CNI_PLUGINS_URL=',variables('vnetCniLinuxPluginsURL'),' CNI_PLUGINS_URL=',variables('cniPluginsURL'),' CLOUDPROVIDER_BACKOFF=',variables('cloudProviderBackoff'),' CLOUDPROVIDER_BACKOFF_RETRIES=',variables('cloudProviderBackoffRetries'),' CLOUDPROVIDER_BACKOFF_EXPONENT=',variables('cloudProviderBackoffExponent'),' CLOUDPROVIDER_BACKOFF_DURATION=',variables('cloudProviderBackoffDuration'),' CLOUDPROVIDER_BACKOFF_JITTER=',variables('cloudProviderBackoffJitter'),' CLOUDPROVIDER_RATELIMIT=',variables('cloudProviderRatelimit'),' CLOUDPROVIDER_RATELIMIT_QPS=',variables('cloudProviderRatelimitQPS'),' CLOUDPROVIDER_RATELIMIT_BUCKET=',variables('cloudProviderRatelimitBucket'),' USE_MANAGED_IDENTITY_EXTENSION=',variables('useManagedIdentityExtension'),' USE_INSTANCE_METADATA=',variables('useInstanceMetadata'),' CONTAINER_RUNTIME=',variables('containerRuntime'),' CONTAINERD_DOWNLOAD_URL_BASE=',variables('containerdDownloadURLBase'),' POD_INFRA_CONTAINER_SPEC=',variables('kubernetesPodInfraContainerSpec'),' KUBECONFIG_SERVER=',variables('kubeconfigServer'),' KMS_PROVIDER_VAULT_NAME=',variables('clusterKeyVaultName'), ' EnableEncryptionWithExternalKms=true')]",
        {{else}}
    "provisionScriptParametersCommon": "[concat('ADMINUSER=',variables('username'),' DOCKER_ENGINE_VERSION=',variables('dockerEngineVersion'),' DOCKER_REPO=',variables('dockerEngineDownloadRepo'),' TENANT_ID=',variables('tenantID'),' HYPERKUBE_URL=',variables('kubernetesHyperkubeSpec'),' APISERVER_PUBLIC_KEY=',variables('apiserverCertificate'),' SUBSCRIPTION_ID=',variables('subscriptionId'),' RESOURCE_GROUP=',variables('resourceGroup'),' LOCATION=',variables('location'),' VM_TYPE=',variables('vmType'),' SUBNET=',variables('subnetName'),' NETWORK_SECURITY_GROUP=',variables('nsgName'),' VIRTUAL_NETWORK=',variables('virtualNetworkName'),' VIRTUAL_NETWORK_RESOURCE_GROUP=',variables('virtualNetworkResourceGroupName'),' ROUTE_TABLE=',variables('routeTableName'),' PRIMARY_AVAILABILITY_SET=',variables('primaryAvailabilitySetName'),' PRIMARY_SCALE_SET=',variables('primaryScaleSetName'),' SERVICE_PRINCIPAL_CLIENT_ID=',variables('servicePrincipalClientId'),' SERVICE_PRINCIPAL_CLIENT_SECRET=',variables('singleQuote'),variables('servicePrincipalClientSecret'),variables('singleQuote'),' KUBELET_PRIVATE_KEY=',variables('clientPrivateKey'),' TARGET_ENVIRONMENT=',variables('targetEnvironment'),' NETWORK_PLUGIN=',variables('networkPlugin'),' VNET_CNI_PLUGINS_URL=',variables('vnetCniLinuxPluginsURL'),' CNI_PLUGINS_URL=',variables('cniPluginsURL'),' CLOUDPROVIDER_BACKOFF=',variables('cloudProviderBackoff'),' CLOUDPROVIDER_BACKOFF_RETRIES=',variables('cloudProviderBackoffRetries'),' CLOUDPROVIDER_BACKOFF_EXPONENT=',variables('cloudProviderBackoffExponent'),' CLOUDPROVIDER_BACKOFF_DURATION=',variables('cloudProviderBackoffDuration'),' CLOUDPROVIDER_BACKOFF_JITTER=',variables('cloudProviderBackoffJitter'),' CLOUDPROVIDER_RATELIMIT=',variables('cloudProviderRatelimit'),' CLOUDPROVIDER_RATELIMIT_QPS=',variables('cloudProviderRatelimitQPS'),' CLOUDPROVIDER_RATELIMIT_BUCKET=',variables('cloudProviderRatelimitBucket'),' USE_MANAGED_IDENTITY_EXTENSION=',variables('useManagedIdentityExtension'),' USE_INSTANCE_METADATA=',variables('useInstanceMetadata'),' CONTAINER_RUNTIME=',variables('containerRuntime'),' CONTAINERD_DOWNLOAD_URL_BASE=',variables('containerdDownloadURLBase'),' POD_INFRA_CONTAINER_SPEC=',variables('kubernetesPodInfraContainerSpec'),' KUBECONFIG_SERVER=',variables('kubeconfigServer'))]",
        {{end}}
    {{else}}
        {{if EnableEncryptionWithExternalKms}}
            {{ if not UseManagedIdentity}}
    "servicePrincipalObjectId": "[parameters('servicePrincipalObjectId')]",
            {{end}}
    "provisionScriptParametersCommon": "[concat('ADMINUSER=',variables('username'),' DOCKER_ENGINE_VERSION=',variables('dockerEngineVersion'),' DOCKER_REPO=',variables('dockerEngineDownloadRepo'),' TENANT_ID=',variables('tenantID'),' HYPERKUBE_URL=',variables('kubernetesHyperkubeSpec'),' APISERVER_PUBLIC_KEY=',variables('apiserverCertificate'),' SUBSCRIPTION_ID=',variables('subscriptionId'),' RESOURCE_GROUP=',variables('resourceGroup'),' LOCATION=',variables('location'),' VM_TYPE=',variables('vmType'),' SUBNET=',variables('subnetName'),' NETWORK_SECURITY_GROUP=',variables('nsgName'),' VIRTUAL_NETWORK=',variables('virtualNetworkName'),' VIRTUAL_NETWORK_RESOURCE_GROUP=',variables('virtualNetworkResourceGroupName'),' ROUTE_TABLE=',variables('routeTableName'),' PRIMARY_AVAILABILITY_SET=',variables('primaryAvailabilitySetName'),' PRIMARY_SCALE_SET=',variables('primaryScaleSetName'),' SERVICE_PRINCIPAL_CLIENT_ID=',variables('servicePrincipalClientId'),' SERVICE_PRINCIPAL_CLIENT_SECRET=',variables('singleQuote'),variables('servicePrincipalClientSecret'),variables('singleQuote'),' KUBELET_PRIVATE_KEY=',variables('clientPrivateKey'),' TARGET_ENVIRONMENT=',variables('targetEnvironment'),' NETWORK_PLUGIN=',variables('networkPlugin'),' VNET_CNI_PLUGINS_URL=',variables('vnetCniLinuxPluginsURL'),' CNI_PLUGINS_URL=',variables('cniPluginsURL'),' CLOUDPROVIDER_BACKOFF=',variables('cloudProviderBackoff'),' CLOUDPROVIDER_BACKOFF_RETRIES=',variables('cloudProviderBackoffRetries'),' CLOUDPROVIDER_BACKOFF_EXPONENT=',variables('cloudProviderBackoffExponent'),' CLOUDPROVIDER_BACKOFF_DURATION=',variables('cloudProviderBackoffDuration'),' CLOUDPROVIDER_BACKOFF_JITTER=',variables('cloudProviderBackoffJitter'),' CLOUDPROVIDER_RATELIMIT=',variables('cloudProviderRatelimit'),' CLOUDPROVIDER_RATELIMIT_QPS=',variables('cloudProviderRatelimitQPS'),' CLOUDPROVIDER_RATELIMIT_BUCKET=',variables('cloudProviderRatelimitBucket'),' USE_MANAGED_IDENTITY_EXTENSION=',variables('useManagedIdentityExtension'),' USE_INSTANCE_METADATA=',variables('useInstanceMetadata'),' CONTAINER_RUNTIME=',variables('containerRuntime'),' CONTAINERD_DOWNLOAD_URL_BASE=',variables('containerdDownloadURLBase'),' POD_INFRA_CONTAINER_SPEC=',variables('kubernetesPodInfraContainerSpec'),' KMS_PROVIDER_VAULT_NAME=',variables('clusterKeyVaultName'), ' EnableEncryptionWithExternalKms=true')]",
        {{else}}
    "provisionScriptParametersCommon": "[concat('ADMINUSER=',variables('username'),' DOCKER_ENGINE_VERSION=',variables('dockerEngineVersion'),' DOCKER_REPO=',variables('dockerEngineDownloadRepo'),' TENANT_ID=',variables('tenantID'),' HYPERKUBE_URL=',variables('kubernetesHyperkubeSpec'),' APISERVER_PUBLIC_KEY=',variables('apiserverCertificate'),' SUBSCRIPTION_ID=',variables('subscriptionId'),' RESOURCE_GROUP=',variables('resourceGroup'),' LOCATION=',variables('location'),' VM_TYPE=',variables('vmType'),' SUBNET=',variables('subnetName'),' NETWORK_SECURITY_GROUP=',variables('nsgName'),' VIRTUAL_NETWORK=',variables('virtualNetworkName'),' VIRTUAL_NETWORK_RESOURCE_GROUP=',variables('virtualNetworkResourceGroupName'),' ROUTE_TABLE=',variables('routeTableName'),' PRIMARY_AVAILABILITY_SET=',variables('primaryAvailabilitySetName'),' PRIMARY_SCALE_SET=',variables('primaryScaleSetName'),' SERVICE_PRINCIPAL_CLIENT_ID=',variables('servicePrincipalClientId'),' SERVICE_PRINCIPAL_CLIENT_SECRET=',variables('singleQuote'),variables('servicePrincipalClientSecret'),variables('singleQuote'),' KUBELET_PRIVATE_KEY=',variables('clientPrivateKey'),' TARGET_ENVIRONMENT=',variables('targetEnvironment'),' NETWORK_PLUGIN=',variables('networkPlugin'),' VNET_CNI_PLUGINS_URL=',variables('vnetCniLinuxPluginsURL'),' CNI_PLUGINS_URL=',variables('cniPluginsURL'),' CLOUDPROVIDER_BACKOFF=',variables('cloudProviderBackoff'),' CLOUDPROVIDER_BACKOFF_RETRIES=',variables('cloudProviderBackoffRetries'),' CLOUDPROVIDER_BACKOFF_EXPONENT=',variables('cloudProviderBackoffExponent'),' CLOUDPROVIDER_BACKOFF_DURATION=',variables('cloudProviderBackoffDuration'),' CLOUDPROVIDER_BACKOFF_JITTER=',variables('cloudProviderBackoffJitter'),' CLOUDPROVIDER_RATELIMIT=',variables('cloudProviderRatelimit'),' CLOUDPROVIDER_RATELIMIT_QPS=',variables('cloudProviderRatelimitQPS'),' CLOUDPROVIDER_RATELIMIT_BUCKET=',variables('cloudProviderRatelimitBucket'),' USE_MANAGED_IDENTITY_EXTENSION=',variables('useManagedIdentityExtension'),' USE_INSTANCE_METADATA=',variables('useInstanceMetadata'),' CONTAINER_RUNTIME=',variables('containerRuntime'),' CONTAINERD_DOWNLOAD_URL_BASE=',variables('containerdDownloadURLBase'),' POD_INFRA_CONTAINER_SPEC=',variables('kubernetesPodInfraContainerSpec'))]",
        {{end}}
    {{end}}
{{end}}
    "generateProxyCertsScript": "{{GetKubernetesB64GenerateProxyCerts}}",
    "acsengineVersion": "[parameters('acsengineVersion')]",
    "orchestratorNameVersionTag": "{{.OrchestratorProfile.OrchestratorType}}:{{.OrchestratorProfile.OrchestratorVersion}}",

{{if IsAzureCNI}}
    "allocateNodeCidrs": false,
    "AzureCNINetworkMonitorImageURL": "[parameters('AzureCNINetworkMonitorImageURL')]",
{{else}}
    "allocateNodeCidrs": true,
{{end}}
{{if HasCustomSearchDomain}}
    "searchDomainName": "[parameters('searchDomainName')]",
    "searchDomainRealmUser": "[parameters('searchDomainRealmUser')]",
    "searchDomainRealmPassword": "[parameters('searchDomainRealmPassword')]",
{{end}}
{{if HasCustomNodesDNS}}
    "dnsServer": "[parameters('dnsServer')]",
{{end}}
{{if not IsHostedMaster}}
  {{if .MasterProfile.IsCustomVNET}}
    "vnetSubnetID": "[parameters('masterVnetSubnetID')]",
    "subnetNameResourceSegmentIndex": 10,
    "subnetName": "[split(parameters('masterVnetSubnetID'), '/')[variables('subnetNameResourceSegmentIndex')]]",
    "vnetNameResourceSegmentIndex": 8,
    "virtualNetworkName": "[split(parameters('masterVnetSubnetID'), '/')[variables('vnetNameResourceSegmentIndex')]]",
    "vnetResourceGroupNameResourceSegmentIndex": 4,
    "virtualNetworkResourceGroupName": "[split(parameters('masterVnetSubnetID'), '/')[variables('vnetResourceGroupNameResourceSegmentIndex')]]",
  {{else}}
    "subnet": "[parameters('masterSubnet')]",
    "subnetName": "[concat(variables('orchestratorName'), '-subnet')]",
    "vnetID": "[resourceId('Microsoft.Network/virtualNetworks',variables('virtualNetworkName'))]",
    "vnetSubnetID": "[concat(variables('vnetID'),'/subnets/',variables('subnetName'))]",
    "virtualNetworkName": "[concat(variables('orchestratorName'), '-vnet-', variables('nameSuffix'))]",
    "virtualNetworkResourceGroupName": "''",
  {{end}}
{{else}}
    "subnet": "[parameters('masterSubnet')]",
  {{if IsCustomVNET}}
    "vnetSubnetID": "[parameters('{{ (index .AgentPoolProfiles 0).Name }}VnetSubnetID')]",
    "subnetNameResourceSegmentIndex": 10,
    "subnetName": "[split(variables('vnetSubnetID'), '/')[variables('subnetNameResourceSegmentIndex')]]",
    "vnetNameResourceSegmentIndex": 8,
    "virtualNetworkName": "[split(variables('vnetSubnetID'), '/')[variables('vnetNameResourceSegmentIndex')]]",
    "vnetResourceGroupNameResourceSegmentIndex": 4,
    "virtualNetworkResourceGroupName": "[split(variables('vnetSubnetID'), '/')[variables('vnetResourceGroupNameResourceSegmentIndex')]]",
  {{else}}
    "subnetName": "[concat(variables('orchestratorName'), '-subnet')]",
    "vnetID": "[resourceId('Microsoft.Network/virtualNetworks',variables('virtualNetworkName'))]",
    "vnetSubnetID": "[concat(variables('vnetID'),'/subnets/',variables('subnetName'))]",
    "virtualNetworkName": "[concat(variables('orchestratorName'), '-vnet-', variables('nameSuffix'))]",
    "virtualNetworkResourceGroupName": "",
  {{end}}
{{end}}
    "vnetCidr": "[parameters('vnetCidr')]",
    "kubeDNSServiceIP": "[parameters('kubeDNSServiceIP')]",
    "kubeServiceCidr": "[parameters('kubeServiceCidr')]",
    "kubeClusterCidr": "[parameters('kubeClusterCidr')]",
    "dockerBridgeCidr": "[parameters('dockerBridgeCidr')]",
{{if IsKubernetesVersionGe "1.6.0"}}
    "registerWithTaints": "node-role.kubernetes.io/master=true:NoSchedule",
{{else}}
    {{if HasLinuxAgents}}
    "registerSchedulable": "false",
    {{else}}
    "registerSchedulable": "true",
    {{end}}
{{end}}
{{if not IsHostedMaster }}
    "nsgName": "[concat(variables('masterVMNamePrefix'), 'nsg')]",
{{else}}
    "nsgName": "[concat(variables('agentNamePrefix'), 'nsg')]",
{{end}}
    "nsgID": "[resourceId('Microsoft.Network/networkSecurityGroups',variables('nsgName'))]",
{{if not AnyAgentUsesVirtualMachineScaleSets}}
    "primaryAvailabilitySetName": "[concat('{{ (index .AgentPoolProfiles 0).Name }}-availabilitySet-',variables('nameSuffix'))]",
    "primaryScaleSetName": "''",
    "vmType": "standard",
{{else}}
    "primaryScaleSetName": "[concat(variables('orchestratorName'), '-{{ (index .AgentPoolProfiles 0).Name }}-',variables('nameSuffix'), '-vmss')]",
    "primaryAvailabilitySetName": "''",
    "vmType": "vmss",
{{end}}
{{if not IsHostedMaster }}
    {{if IsPrivateCluster}}
        "kubeconfigServer": "[concat('https://', variables('kubernetesAPIServerIP'), ':443')]",
        {{if ProvisionJumpbox}}
            "jumpboxVMName": "[parameters('jumpboxVMName')]",
            "jumpboxVMSize": "[parameters('jumpboxVMSize')]",
            "jumpboxOSDiskSizeGB": "[parameters('jumpboxOSDiskSizeGB')]",
            "jumpboxOSDiskName": "[concat(variables('jumpboxVMName'), '-osdisk')]",
            "jumpboxUsername": "[parameters('jumpboxUsername')]",
            "jumpboxPublicKey": "[parameters('jumpboxPublicKey')]",
            "jumpboxStorageProfile": "[parameters('jumpboxStorageProfile')]",
            "jumpboxPublicIpAddressName": "[concat(variables('jumpboxVMName'), '-ip')]",
            "jumpboxNetworkInterfaceName": "[concat(variables('jumpboxVMName'), '-nic')]",
            "jumpboxNetworkSecurityGroupName": "[concat(variables('jumpboxVMName'), '-nsg')]",
            "kubeconfig": "{{GetKubeConfig}}",
            {{if not JumpboxIsManagedDisks}}
            "jumpboxStorageAccountName": "[concat(variables('storageAccountBaseName'), 'jb')]",
            {{end}}
            {{if not .HasStorageAccountDisks}}
            {{GetSizeMap}},
            {{end}}
        {{end}}
    {{else}}
        "masterPublicIPAddressName": "[concat(variables('orchestratorName'), '-master-ip-', variables('masterFqdnPrefix'), '-', variables('nameSuffix'))]",
        "masterLbID": "[resourceId('Microsoft.Network/loadBalancers',variables('masterLbName'))]",
        "masterLbIPConfigID": "[concat(variables('masterLbID'),'/frontendIPConfigurations/', variables('masterLbIPConfigName'))]",
        "masterLbIPConfigName": "[concat(variables('orchestratorName'), '-master-lbFrontEnd-', variables('nameSuffix'))]",
        "masterLbName": "[concat(variables('orchestratorName'), '-master-lb-', variables('nameSuffix'))]",
        "kubeconfigServer": "[concat('https://', variables('masterFqdnPrefix'), '.', variables('location'), '.', variables('fqdnEndpointSuffix'))]",
    {{end}}
{{if gt .MasterProfile.Count 1}}
    "masterInternalLbName": "[concat(variables('orchestratorName'), '-master-internal-lb-', variables('nameSuffix'))]",
    "masterInternalLbID": "[resourceId('Microsoft.Network/loadBalancers',variables('masterInternalLbName'))]",
    "masterInternalLbIPConfigName": "[concat(variables('orchestratorName'), '-master-internal-lbFrontEnd-', variables('nameSuffix'))]",
    "masterInternalLbIPConfigID": "[concat(variables('masterInternalLbID'),'/frontendIPConfigurations/', variables('masterInternalLbIPConfigName'))]",
    "masterInternalLbIPOffset": {{GetDefaultInternalLbStaticIPOffset}},
    "kubernetesAPIServerIP": "[concat(variables('masterFirstAddrPrefix'), add(variables('masterInternalLbIPOffset'), int(variables('masterFirstAddrOctet4'))))]",
{{else}}
    "kubernetesAPIServerIP": "[parameters('firstConsecutiveStaticIP')]",
{{end}}
    "masterLbBackendPoolName": "[concat(variables('orchestratorName'), '-master-pool-', variables('nameSuffix'))]",
    "masterFirstAddrComment": "these MasterFirstAddrComment are used to place multiple masters consecutively in the address space",
    "masterFirstAddrOctets": "[split(parameters('firstConsecutiveStaticIP'),'.')]",
    "masterFirstAddrOctet4": "[variables('masterFirstAddrOctets')[3]]",
    "masterFirstAddrPrefix": "[concat(variables('masterFirstAddrOctets')[0],'.',variables('masterFirstAddrOctets')[1],'.',variables('masterFirstAddrOctets')[2],'.')]",
    "masterVMNamePrefix": "[concat(variables('orchestratorName'), '-master-', variables('nameSuffix'), '-')]",
    "masterVMNames": [
      "[concat(variables('masterVMNamePrefix'), '0')]",
      "[concat(variables('masterVMNamePrefix'), '1')]",
      "[concat(variables('masterVMNamePrefix'), '2')]",
      "[concat(variables('masterVMNamePrefix'), '3')]",
      "[concat(variables('masterVMNamePrefix'), '4')]"
    ],
    "masterPrivateIpAddrs": [
      "[concat(variables('masterFirstAddrPrefix'), add(0, int(variables('masterFirstAddrOctet4'))))]",
      "[concat(variables('masterFirstAddrPrefix'), add(1, int(variables('masterFirstAddrOctet4'))))]",
      "[concat(variables('masterFirstAddrPrefix'), add(2, int(variables('masterFirstAddrOctet4'))))]",
      "[concat(variables('masterFirstAddrPrefix'), add(3, int(variables('masterFirstAddrOctet4'))))]",
      "[concat(variables('masterFirstAddrPrefix'), add(4, int(variables('masterFirstAddrOctet4'))))]"
    ],
    "masterEtcdServerPort": {{GetMasterEtcdServerPort}},
    "masterEtcdClientPort": {{GetMasterEtcdClientPort}},
    "masterEtcdPeerURLs":[
      "[concat('https://', variables('masterPrivateIpAddrs')[0], ':', variables('masterEtcdServerPort'))]",
      "[concat('https://', variables('masterPrivateIpAddrs')[1], ':', variables('masterEtcdServerPort'))]",
      "[concat('https://', variables('masterPrivateIpAddrs')[2], ':', variables('masterEtcdServerPort'))]",
      "[concat('https://', variables('masterPrivateIpAddrs')[3], ':', variables('masterEtcdServerPort'))]",
      "[concat('https://', variables('masterPrivateIpAddrs')[4], ':', variables('masterEtcdServerPort'))]"
    ],
    "masterEtcdClientURLs":[
      "[concat('https://', variables('masterPrivateIpAddrs')[0], ':', variables('masterEtcdClientPort'))]",
      "[concat('https://', variables('masterPrivateIpAddrs')[1], ':', variables('masterEtcdClientPort'))]",
      "[concat('https://', variables('masterPrivateIpAddrs')[2], ':', variables('masterEtcdClientPort'))]",
      "[concat('https://', variables('masterPrivateIpAddrs')[3], ':', variables('masterEtcdClientPort'))]",
      "[concat('https://', variables('masterPrivateIpAddrs')[4], ':', variables('masterEtcdClientPort'))]"
    ],
    "masterEtcdClusterStates": [
      "[concat(variables('masterVMNames')[0], '=', variables('masterEtcdPeerURLs')[0])]",
      "[concat(variables('masterVMNames')[0], '=', variables('masterEtcdPeerURLs')[0], ',', variables('masterVMNames')[1], '=', variables('masterEtcdPeerURLs')[1], ',', variables('masterVMNames')[2], '=', variables('masterEtcdPeerURLs')[2])]",
      "[concat(variables('masterVMNames')[0], '=', variables('masterEtcdPeerURLs')[0], ',', variables('masterVMNames')[1], '=', variables('masterEtcdPeerURLs')[1], ',', variables('masterVMNames')[2], '=', variables('masterEtcdPeerURLs')[2], ',', variables('masterVMNames')[3], '=', variables('masterEtcdPeerURLs')[3], ',', variables('masterVMNames')[4], '=', variables('masterEtcdPeerURLs')[4])]"
    ],
{{else}}
    "kubernetesAPIServerIP": "[parameters('kubernetesEndpoint')]",
    "agentNamePrefix": "[concat(variables('orchestratorName'), '-agentpool-', variables('nameSuffix'), '-')]",
{{end}}
    "subscriptionId": "[subscription().subscriptionId]",
    "contributorRoleDefinitionId": "[concat('/subscriptions/', subscription().subscriptionId, '/providers/Microsoft.Authorization/roleDefinitions/', 'b24988ac-6180-42a0-ab88-20f7382dd24c')]",
    "readerRoleDefinitionId": "[concat('/subscriptions/', subscription().subscriptionId, '/providers/Microsoft.Authorization/roleDefinitions/', 'acdd72a7-3385-48ef-bd42-f606fba81ae7')]",
    "scope": "[resourceGroup().id]",
    "tenantId": "[subscription().tenantId]",
    "singleQuote": "'",
    "targetEnvironment": "[parameters('targetEnvironment')]"
{{if not IsOpenShift}}
    , "dockerEngineDownloadRepo": "[parameters('dockerEngineDownloadRepo')]"
    , "dockerEngineVersion": "[parameters('dockerEngineVersion')]"
{{end}}
{{if .LinuxProfile.HasSecrets}}
    , "linuxProfileSecrets" :
      [
          {{range  $vIndex, $vault := .LinuxProfile.Secrets}}
            {{if $vIndex}} , {{end}}
              {
                "sourceVault":{
                  "id":"[parameters('linuxKeyVaultID{{$vIndex}}')]"
                },
                "vaultCertificates":[
                {{range $cIndex, $cert := $vault.VaultCertificates}}
                  {{if $cIndex}} , {{end}}
                  {
                    "certificateUrl" :"[parameters('linuxKeyVaultID{{$vIndex}}CertificateURL{{$cIndex}}')]"
                  }
                {{end}}
                ]
              }
        {{end}}
      ]
{{end}}
{{if .HasWindows}}
    ,"windowsAdminUsername": "[parameters('windowsAdminUsername')]",
    "windowsAdminPassword": "[parameters('windowsAdminPassword')]",
    "kubeBinariesSASURL": "[parameters('kubeBinariesSASURL')]",
    "windowsPackageSASURLBase": "[parameters('windowsPackageSASURLBase')]",
    "kubeBinariesVersion": "[parameters('kubeBinariesVersion')]",
    "windowsTelemetryGUID": "[parameters('windowsTelemetryGUID')]",
    "agentWindowsPublisher": "[parameters('agentWindowsPublisher')]",
    "agentWindowsOffer": "[parameters('agentWindowsOffer')]",
    "agentWindowsSku": "[parameters('agentWindowsSku')]",
    "agentWindowsVersion": "[parameters('agentWindowsVersion')]",
    "windowsCustomScriptSuffix": " $inputFile = '%SYSTEMDRIVE%\\AzureData\\CustomData.bin' ; $outputFile = '%SYSTEMDRIVE%\\AzureData\\CustomDataSetupScript.ps1' ; Copy-Item $inputFile $outputFile ; Invoke-Expression('{0} {1}' -f $outputFile, $arguments) ; "
{{end}}
{{if EnableEncryptionWithExternalKms}}
     ,"apiVersionKeyVault": "2016-10-01",
     {{if not .HasStorageAccountDisks}}
     "apiVersionStorage": "2015-06-15",
     {{end}}
     "clusterKeyVaultName": "[take(concat('kv', tolower(uniqueString(concat(variables('masterFqdnPrefix'),variables('location'),variables('nameSuffix'))))), 22)]",
     "clusterKeyVaultSku" : "[parameters('clusterKeyVaultSku')]"
{{end}}
