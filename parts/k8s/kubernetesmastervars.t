    "maxVMsPerPool": 100,
{{ if IsOpenShift }}
    "routerNSGName": "[concat(parameters('orchestratorName'), '-router-', parameters('nameSuffix'), '-nsg')]",
    "routerNSGID": "[resourceId('Microsoft.Network/networkSecurityGroups', variables('routerNSGName'))]",
    "routerIPName": "[concat(parameters('orchestratorName'), '-router-ip-', variables('masterFqdnPrefix'), '-', parameters('nameSuffix'))]",
    "routerLBName": "[concat(parameters('orchestratorName'), '-router-lb-', parameters('nameSuffix'))]",
    "routerLBID": "[resourceId('Microsoft.Network/loadBalancers', variables('routerLBName'))]",
{{ end }}
{{ if not IsHostedMaster }}
    {{ if not IsOpenShift }}
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
    "useManagedIdentityExtension": "{{ UseManagedIdentity }}",
    "userAssignedID": "{{UserAssignedID}}",
    "userAssignedClientID": "{{UserAssignedClientID}}",
    "useInstanceMetadata": "{{ UseInstanceMetadata }}",
    "loadBalancerSku": "{{ LoadBalancerSku }}",
    "excludeMasterFromStandardLB": "{{ ExcludeMasterFromStandardLB }}",
{{ if UseManagedIdentity }}
    "servicePrincipalClientId": "msi",
    "servicePrincipalClientSecret": "msi",
{{ else }}
    "servicePrincipalClientId": "[parameters('servicePrincipalClientId')]",
    "servicePrincipalClientSecret": "[parameters('servicePrincipalClientSecret')]",
{{ end }}
    "masterFqdnPrefix": "[tolower(parameters('masterEndpointDNSNamePrefix'))]",
{{if not IsHostedMaster}}
    "masterCount": {{.MasterProfile.Count}},
    "masterOffset": "[parameters('masterOffset')]",
{{end}}
    "apiVersionDefault": "2016-03-30",
    "apiVersionLinkDefault": "2015-01-01",
    "locations": [
         "[resourceGroup().location]",
         "[parameters('location')]"
    ],
    "location": "[variables('locations')[mod(add(2,length(parameters('location'))),add(1,length(parameters('location'))))]]",
    "masterAvailabilitySet": "[concat('master-availabilityset-', parameters('nameSuffix'))]",
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
    "sshKeyPath": "[concat('/home/',parameters('linuxAdminUsername'),'/.ssh/authorized_keys')]",

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
{{if UserAssignedIDEnabled}}
    "apiVersionUserMSI": "2018-06-01",
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
    "healthMonitorScript": "{{GetKubernetesB64HealthMonitorScript}}",
    "provisionInstalls": "{{GetKubernetesB64Installs}}",
    "provisionConfigs": "{{GetKubernetesB64Configs}}",
    "mountetcdScript": "{{GetKubernetesB64Mountetcd}}",
    "customSearchDomainsScript": "{{GetKubernetesB64CustomSearchDomainsScript}}",
    "sshdConfig": "{{GetB64sshdConfig}}",
{{if not IsOpenShift}}
    "provisionScriptParametersCommon": "[concat('ADMINUSER=',parameters('linuxAdminUsername'),' ETCD_DOWNLOAD_URL=',parameters('etcdDownloadURLBase'),' ETCD_VERSION=',parameters('etcdVersion'),' DOCKER_ENGINE_VERSION=',parameters('dockerEngineVersion'),' DOCKER_REPO=',parameters('dockerEngineDownloadRepo'),' TENANT_ID=',variables('tenantID'),' KUBERNETES_VERSION={{.OrchestratorProfile.OrchestratorVersion}} HYPERKUBE_URL=',parameters('kubernetesHyperkubeSpec'),' APISERVER_PUBLIC_KEY=',parameters('apiserverCertificate'),' SUBSCRIPTION_ID=',variables('subscriptionId'),' RESOURCE_GROUP=',variables('resourceGroup'),' LOCATION=',variables('location'),' VM_TYPE=',variables('vmType'),' SUBNET=',variables('subnetName'),' NETWORK_SECURITY_GROUP=',variables('nsgName'),' VIRTUAL_NETWORK=',variables('virtualNetworkName'),' VIRTUAL_NETWORK_RESOURCE_GROUP=',variables('virtualNetworkResourceGroupName'),' ROUTE_TABLE=',variables('routeTableName'),' PRIMARY_AVAILABILITY_SET=',variables('primaryAvailabilitySetName'),' PRIMARY_SCALE_SET=',variables('primaryScaleSetName'),' SERVICE_PRINCIPAL_CLIENT_ID=',variables('servicePrincipalClientId'),' SERVICE_PRINCIPAL_CLIENT_SECRET=',variables('singleQuote'),variables('servicePrincipalClientSecret'),variables('singleQuote'),' KUBELET_PRIVATE_KEY=',parameters('clientPrivateKey'),' TARGET_ENVIRONMENT=',parameters('targetEnvironment'),' NETWORK_PLUGIN=',parameters('networkPlugin'),' VNET_CNI_PLUGINS_URL=',parameters('vnetCniLinuxPluginsURL'),' CNI_PLUGINS_URL=',parameters('cniPluginsURL'),' CLOUDPROVIDER_BACKOFF=',toLower(string(parameters('cloudproviderConfig').cloudProviderBackoff)),' CLOUDPROVIDER_BACKOFF_RETRIES=',parameters('cloudproviderConfig').cloudProviderBackoffRetries,' CLOUDPROVIDER_BACKOFF_EXPONENT=',parameters('cloudproviderConfig').cloudProviderBackoffExponent,' CLOUDPROVIDER_BACKOFF_DURATION=',parameters('cloudproviderConfig').cloudProviderBackoffDuration,' CLOUDPROVIDER_BACKOFF_JITTER=',parameters('cloudproviderConfig').cloudProviderBackoffJitter,' CLOUDPROVIDER_RATELIMIT=',toLower(string(parameters('cloudproviderConfig').cloudProviderRatelimit)),' CLOUDPROVIDER_RATELIMIT_QPS=',parameters('cloudproviderConfig').cloudProviderRatelimitQPS,' CLOUDPROVIDER_RATELIMIT_BUCKET=',parameters('cloudproviderConfig').cloudProviderRatelimitBucket,' USE_MANAGED_IDENTITY_EXTENSION=',variables('useManagedIdentityExtension'),' USER_ASSIGNED_IDENTITY_ID=',variables('userAssignedClientID'),' USE_INSTANCE_METADATA=',variables('useInstanceMetadata'),' LOAD_BALANCER_SKU=',variables('loadBalancerSku'),' EXCLUDE_MASTER_FROM_STANDARD_LB=',variables('excludeMasterFromStandardLB'),' CONTAINER_RUNTIME=',parameters('containerRuntime'),' CONTAINERD_DOWNLOAD_URL_BASE=',parameters('containerdDownloadURLBase'),' POD_INFRA_CONTAINER_SPEC=',parameters('kubernetesPodInfraContainerSpec'),' KMS_PROVIDER_VAULT_NAME=',variables('clusterKeyVaultName'))]",
    {{if not IsHostedMaster}}
    "provisionScriptParametersMaster": "[concat('MASTER_VM_NAME=',variables('masterVMNames')[variables('masterOffset')],' ETCD_PEER_URL=',variables('masterEtcdPeerURLs')[variables('masterOffset')],' ETCD_CLIENT_URL=',variables('masterEtcdClientURLs')[variables('masterOffset')],' MASTER_NODE=true CLUSTER_AUTOSCALER_ADDON=',parameters('kubernetesClusterAutoscalerEnabled'),' ACI_CONNECTOR_ADDON=',parameters('kubernetesACIConnectorEnabled'),' APISERVER_PRIVATE_KEY=',parameters('apiServerPrivateKey'),' CA_CERTIFICATE=',parameters('caCertificate'),' CA_PRIVATE_KEY=',parameters('caPrivateKey'),' MASTER_FQDN=',variables('masterFqdnPrefix'),' KUBECONFIG_CERTIFICATE=',parameters('kubeConfigCertificate'),' KUBECONFIG_KEY=',parameters('kubeConfigPrivateKey'),' ETCD_SERVER_CERTIFICATE=',parameters('etcdServerCertificate'),' ETCD_CLIENT_CERTIFICATE=',parameters('etcdClientCertificate'),' ETCD_SERVER_PRIVATE_KEY=',parameters('etcdServerPrivateKey'),' ETCD_CLIENT_PRIVATE_KEY=',parameters('etcdClientPrivateKey'),' ETCD_PEER_CERTIFICATES=',string(variables('etcdPeerCertificates')),' ETCD_PEER_PRIVATE_KEYS=',string(variables('etcdPeerPrivateKeys')),' ENABLE_AGGREGATED_APIS=',string(parameters('enableAggregatedAPIs')),' KUBECONFIG_SERVER=',variables('kubeconfigServer'))]",
    {{end}}
{{end}}
    "generateProxyCertsScript": "{{GetKubernetesB64GenerateProxyCerts}}",
    "orchestratorNameVersionTag": "{{.OrchestratorProfile.OrchestratorType}}:{{.OrchestratorProfile.OrchestratorVersion}}",

{{if IsAzureCNI}}
    "allocateNodeCidrs": false,
{{else}}
    "allocateNodeCidrs": true,
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
    "subnetName": "[concat(parameters('orchestratorName'), '-subnet')]",
    "vnetID": "[resourceId('Microsoft.Network/virtualNetworks',variables('virtualNetworkName'))]",
    "vnetSubnetID": "[concat(variables('vnetID'),'/subnets/',variables('subnetName'))]",
    "virtualNetworkName": "[concat(parameters('orchestratorName'), '-vnet-', parameters('nameSuffix'))]",
    "virtualNetworkResourceGroupName": "''",
  {{end}}
{{else}}
  {{if IsCustomVNET}}
    "vnetSubnetID": "[parameters('{{ (index .AgentPoolProfiles 0).Name }}VnetSubnetID')]",
    "subnetNameResourceSegmentIndex": 10,
    "subnetName": "[split(variables('vnetSubnetID'), '/')[variables('subnetNameResourceSegmentIndex')]]",
    "vnetNameResourceSegmentIndex": 8,
    "virtualNetworkName": "[split(variables('vnetSubnetID'), '/')[variables('vnetNameResourceSegmentIndex')]]",
    "vnetResourceGroupNameResourceSegmentIndex": 4,
    "virtualNetworkResourceGroupName": "[split(variables('vnetSubnetID'), '/')[variables('vnetResourceGroupNameResourceSegmentIndex')]]",
  {{else}}
    "subnetName": "[concat(parameters('orchestratorName'), '-subnet')]",
    "vnetID": "[resourceId('Microsoft.Network/virtualNetworks',variables('virtualNetworkName'))]",
    "vnetSubnetID": "[concat(variables('vnetID'),'/subnets/',variables('subnetName'))]",
    "virtualNetworkName": "[concat(parameters('orchestratorName'), '-vnet-', parameters('nameSuffix'))]",
    "virtualNetworkResourceGroupName": "",
  {{end}}
{{end}}
{{if not IsHostedMaster }}
    "nsgName": "[concat(variables('masterVMNamePrefix'), 'nsg')]",
{{else}}
    "nsgName": "[concat(variables('agentNamePrefix'), 'nsg')]",
{{end}}
    "nsgID": "[resourceId('Microsoft.Network/networkSecurityGroups',variables('nsgName'))]",
{{if not AnyAgentUsesVirtualMachineScaleSets}}
    "primaryAvailabilitySetName": "[concat('{{ (index .AgentPoolProfiles 0).Name }}-availabilitySet-',parameters('nameSuffix'))]",
    "primaryScaleSetName": "",
    "vmType": "standard",
{{else}}
    "primaryScaleSetName": "[concat(parameters('orchestratorName'), '-{{ (index .AgentPoolProfiles 0).Name }}-',parameters('nameSuffix'), '-vmss')]",
    "primaryAvailabilitySetName": "",
    "vmType": "vmss",
{{end}}
{{if not IsHostedMaster }}
    {{if IsPrivateCluster}}
        "kubeconfigServer": "[concat('https://', variables('kubernetesAPIServerIP'), ':443')]",
        {{if ProvisionJumpbox}}
            "jumpboxOSDiskName": "[concat(parameters('jumpboxVMName'), '-osdisk')]",
            "jumpboxPublicIpAddressName": "[concat(parameters('jumpboxVMName'), '-ip')]",
            "jumpboxNetworkInterfaceName": "[concat(parameters('jumpboxVMName'), '-nic')]",
            "jumpboxNetworkSecurityGroupName": "[concat(parameters('jumpboxVMName'), '-nsg')]",
            "kubeconfig": "{{GetKubeConfig}}",
            {{if not JumpboxIsManagedDisks}}
            "jumpboxStorageAccountName": "[concat(variables('storageAccountBaseName'), 'jb')]",
            {{end}}
            {{if not .HasStorageAccountDisks}}
            {{GetSizeMap}},
            {{end}}
        {{end}}
    {{else}}
        "masterPublicIPAddressName": "[concat(parameters('orchestratorName'), '-master-ip-', variables('masterFqdnPrefix'), '-', parameters('nameSuffix'))]",
        "masterLbID": "[resourceId('Microsoft.Network/loadBalancers',variables('masterLbName'))]",
        "masterLbIPConfigID": "[concat(variables('masterLbID'),'/frontendIPConfigurations/', variables('masterLbIPConfigName'))]",
        "masterLbIPConfigName": "[concat(parameters('orchestratorName'), '-master-lbFrontEnd-', parameters('nameSuffix'))]",
        "masterLbName": "[concat(parameters('orchestratorName'), '-master-lb-', parameters('nameSuffix'))]",
        "kubeconfigServer": "[concat('https://', variables('masterFqdnPrefix'), '.', variables('location'), '.', parameters('fqdnEndpointSuffix'))]",
    {{end}}
    {{if gt .MasterProfile.Count 1}}
    "masterInternalLbName": "[concat(parameters('orchestratorName'), '-master-internal-lb-', parameters('nameSuffix'))]",
    "masterInternalLbID": "[resourceId('Microsoft.Network/loadBalancers',variables('masterInternalLbName'))]",
    "masterInternalLbIPConfigName": "[concat(parameters('orchestratorName'), '-master-internal-lbFrontEnd-', parameters('nameSuffix'))]",
    "masterInternalLbIPConfigID": "[concat(variables('masterInternalLbID'),'/frontendIPConfigurations/', variables('masterInternalLbIPConfigName'))]",
    "masterInternalLbIPOffset": {{GetDefaultInternalLbStaticIPOffset}},
    "kubernetesAPIServerIP": "[concat(variables('masterFirstAddrPrefix'), add(variables('masterInternalLbIPOffset'), int(variables('masterFirstAddrOctet4'))))]",
    {{else}}
    "kubernetesAPIServerIP": "[parameters('firstConsecutiveStaticIP')]",
    {{end}}
    "masterLbBackendPoolName": "[concat(parameters('orchestratorName'), '-master-pool-', parameters('nameSuffix'))]",
    "masterFirstAddrComment": "these MasterFirstAddrComment are used to place multiple masters consecutively in the address space",
    "masterFirstAddrOctets": "[split(parameters('firstConsecutiveStaticIP'),'.')]",
    "masterFirstAddrOctet4": "[variables('masterFirstAddrOctets')[3]]",
    "masterFirstAddrPrefix": "[concat(variables('masterFirstAddrOctets')[0],'.',variables('masterFirstAddrOctets')[1],'.',variables('masterFirstAddrOctets')[2],'.')]",
    "masterVMNamePrefix": "[concat(parameters('orchestratorName'), '-master-', parameters('nameSuffix'), '-')]",
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
    "agentNamePrefix": "[concat(parameters('orchestratorName'), '-agentpool-', parameters('nameSuffix'), '-')]",
{{end}}
    "subscriptionId": "[subscription().subscriptionId]",
    "contributorRoleDefinitionId": "[concat('/subscriptions/', subscription().subscriptionId, '/providers/Microsoft.Authorization/roleDefinitions/', 'b24988ac-6180-42a0-ab88-20f7382dd24c')]",
    "readerRoleDefinitionId": "[concat('/subscriptions/', subscription().subscriptionId, '/providers/Microsoft.Authorization/roleDefinitions/', 'acdd72a7-3385-48ef-bd42-f606fba81ae7')]",
    "scope": "[resourceGroup().id]",
    "tenantId": "[subscription().tenantId]",
    "singleQuote": "'"
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
    , "windowsCustomScriptSuffix": " $inputFile = '%SYSTEMDRIVE%\\AzureData\\CustomData.bin' ; $outputFile = '%SYSTEMDRIVE%\\AzureData\\CustomDataSetupScript.ps1' ; Copy-Item $inputFile $outputFile ; Invoke-Expression('{0} {1}' -f $outputFile, $arguments) ; "
{{end}}
{{if EnableEncryptionWithExternalKms}}
     ,"apiVersionKeyVault": "2016-10-01",
     {{if not .HasStorageAccountDisks}}
     "apiVersionStorage": "2015-06-15",
     {{end}}
     "clusterKeyVaultName": "[take(concat('kv', tolower(uniqueString(concat(variables('masterFqdnPrefix'),variables('location'),parameters('nameSuffix'))))), 22)]"
{{else}}
    ,"clusterKeyVaultName": ""
{{end}}
