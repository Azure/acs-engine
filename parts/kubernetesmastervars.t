    "maxVMsPerPool": 100,
    "maxVMsPerStorageAccount": 20,
    "maxStorageAccountsPerAgent": "[div(variables('maxVMsPerPool'),variables('maxVMsPerStorageAccount'))]",
    "apiServerCertificate": "[parameters('apiServerCertificate')]",
    "apiServerPrivateKey": "[parameters('apiServerPrivateKey')]",
    "caCertificate": "[parameters('caCertificate')]",
    "clientCertificate": "[parameters('clientCertificate')]",
    "clientPrivateKey": "[parameters('clientPrivateKey')]",
    "kubeConfigCertificate": "[parameters('kubeConfigCertificate')]",
    "kubeConfigPrivateKey": "[parameters('kubeConfigPrivateKey')]",
    "kubernetesHyperkubeSpec": "[parameters('kubernetesHyperkubeSpec')]",
    "kubernetesAddonManagerSpec": "[parameters('kubernetesAddonManagerSpec')]",
    "kubernetesAddonResizerSpec": "[parameters('kubernetesAddonResizerSpec')]",
    "kubernetesDashboardSpec": "[parameters('kubernetesDashboardSpec')]",
    "kubernetesExecHealthzSpec": "[parameters('kubernetesExecHealthzSpec')]",
    "kubernetesHeapsterSpec": "[parameters('kubernetesHeapsterSpec')]",
    "kubernetesPodInfraContainerSpec": "[parameters('kubernetesPodInfraContainerSpec')]",
    "kubernetesKubeDNSSpec": "[parameters('kubernetesKubeDNSSpec')]",
    "kubernetesDNSMasqSpec": "[parameters('kubernetesDNSMasqSpec')]",
    "networkPolicy": "[parameters('networkPolicy')]",
    "servicePrincipalClientId": "[parameters('servicePrincipalClientId')]",
    "servicePrincipalClientSecret": "[parameters('servicePrincipalClientSecret')]",
    "username": "[parameters('linuxAdminUsername')]",
    "masterFqdnPrefix": "[tolower(parameters('masterEndpointDNSNamePrefix'))]",
    "masterPrivateIp": "[parameters('firstConsecutiveStaticIP')]",
    "masterVMSize": "[parameters('masterVMSize')]",
    "sshPublicKeyData": "[parameters('sshRSAPublicKey')]",
{{if  GetClassicMode}}
    "masterCount": "[parameters('masterCount')]",
{{else}}
    "masterCount": {{.MasterProfile.Count}}, 
{{end}}    
    "apiVersionDefault": "2016-03-30",
    "apiVersionStorage": "2015-06-15",
{{if .HasManagedDisks}}
    "apiVersionStorageManagedDisks": "2016-04-30-preview",
{{end}}
    "location": "[resourceGroup().location]", 
    "locations": [
         "[resourceGroup().location]",
         "[parameters('location')]"
    ],
    "location": "[variables('locations')[mod(add(2,length(parameters('location'))),add(1,length(parameters('location'))))]]",
    "masterAvailabilitySet": "[concat('master-availabilityset-', variables('nameSuffix'))]",
    "storageAccountBaseName": "[uniqueString(concat(variables('masterFqdnPrefix'),variables('location'),variables('orchestratorName')))]",
    "masterStorageAccountName": "[concat(variables('storageAccountBaseName'), 'mstr0')]",
    "nameSuffix": "[parameters('nameSuffix')]", 
    "orchestratorName": "k8s",  
    "osImageOffer": "UbuntuServer", 
    "osImagePublisher": "Canonical", 
    "osImageSKU": "16.04-LTS", 
    "osImageVersion": "16.04.201703070",
    "resourceGroup": "[resourceGroup().name]", 
    "routeTableName": "[concat(variables('masterVMNamePrefix'),'routetable')]",
    "routeTableID": "[resourceId('Microsoft.Network/routeTables', variables('routeTableName'))]",
    "sshNatPorts": [22,2201,2202,2203,2204],
    "sshKeyPath": "[concat('/home/',variables('username'),'/.ssh/authorized_keys')]", 
    "storageAccountBaseName": "[uniqueString(concat(variables('masterFqdnPrefix'),variables('location')))]", 
    "storageAccountPrefixes": [ "0", "6", "c", "i", "o", "u", "1", "7", "d", "j", "p", "v", "2", "8", "e", "k", "q", "w", "3", "9", "f", "l", "r", "x", "4", "a", "g", "m", "s", "y", "5", "b", "h", "n", "t", "z" ], 
    "storageAccountPrefixesCount": "[length(variables('storageAccountPrefixes'))]",
    "vmsPerStorageAccount": 20,
    "provisionScript": "{{GetKubernetesB64Provision}}",
{{if AnyAgentHasDisks}}
    "dataStorageAccountPrefixSeed": 97,
{{end}}
{{if IsVNETIntegrated}}
    "allocateNodeCidrs": false,
{{else}}
    "allocateNodeCidrs": true,
{{end}}
{{if .MasterProfile.IsCustomVNET}}
    "vnetSubnetID": "[parameters('masterVnetSubnetID')]",
    "subnetNameResourceSegmentIndex": 10,
    "subnetName": "[split(parameters('masterVnetSubnetID'), '/')[variables('subnetNameResourceSegmentIndex')]]",
    "vnetNameResourceSegmentIndex": 8,
    "virtualNetworkName": "[split(parameters('masterVnetSubnetID'), '/')[variables('vnetNameResourceSegmentIndex')]]",
{{else}}
    "subnet": "[parameters('masterSubnet')]",
    "subnetName": "[concat(variables('orchestratorName'), '-subnet')]",
    "vnetID": "[resourceId('Microsoft.Network/virtualNetworks',variables('virtualNetworkName'))]",
    "vnetSubnetID": "[concat(variables('vnetID'),'/subnets/',variables('subnetName'))]",
    "virtualNetworkName": "[concat(variables('orchestratorName'), '-vnet-', variables('nameSuffix'))]",
    "vnetCidr": "10.0.0.0/8",
{{end}}
    "kubeDnsServiceIp": "10.0.0.10", 
    "kubeServiceCidr": "10.0.0.0/16",
{{if HasLinuxAgents}}
    "registerSchedulable": "false",
{{else}}
    "registerSchedulable": "true",
{{end}}
    "nsgName": "[concat(variables('masterVMNamePrefix'), 'nsg')]",
    "nsgID": "[resourceId('Microsoft.Network/networkSecurityGroups',variables('nsgName'))]",
    "primaryAvailablitySetName": "[concat('{{ (index .AgentPoolProfiles 0).Name }}-availabilitySet-',variables('nameSuffix'))]",
    "masterPublicIPAddressName": "[concat(variables('orchestratorName'), '-master-ip-', variables('masterFqdnPrefix'), '-', variables('nameSuffix'))]",
    "masterLbID": "[resourceId('Microsoft.Network/loadBalancers',variables('masterLbName'))]", 
    "masterLbIPConfigID": "[concat(variables('masterLbID'),'/frontendIPConfigurations/', variables('masterLbIPConfigName'))]", 
    "masterLbIPConfigName": "[concat(variables('orchestratorName'), '-master-lbFrontEnd-', variables('nameSuffix'))]",
    "masterLbName": "[concat(variables('orchestratorName'), '-master-lb-', variables('nameSuffix'))]",
    "masterInternalLbName": "[concat(variables('orchestratorName'), '-master-internal-lb-', variables('nameSuffix'))]",
    "masterInternalLbID": "[resourceId('Microsoft.Network/loadBalancers',variables('masterInternalLbName'))]",
    "masterInternalLbIPConfigName": "[concat(variables('orchestratorName'), '-master-internal-lbFrontEnd-', variables('nameSuffix'))]",
    "masterInternalLbIPConfigID": "[concat(variables('masterInternalLbID'),'/frontendIPConfigurations/', variables('masterInternalLbIPConfigName'))]",
    "masterInternalLbIPOffset": {{GetDefaultInternalLbStaticIPOffset}},
    "masterInternalLbIp": "[concat(variables('masterFirstAddrPrefix'), add(variables('masterInternalLbIPOffset'), int(variables('masterFirstAddrOctet4'))))]",
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
    "masterEtcdServerPort": 2380,
    "masterEtcdClientPort": 2379,
    "masterEtcdPeerURLs":[
      "[concat('http://', variables('masterPrivateIpAddrs')[0], ':', variables('masterEtcdServerPort'))]",
      "[concat('http://', variables('masterPrivateIpAddrs')[1], ':', variables('masterEtcdServerPort'))]",
      "[concat('http://', variables('masterPrivateIpAddrs')[2], ':', variables('masterEtcdServerPort'))]",
      "[concat('http://', variables('masterPrivateIpAddrs')[3], ':', variables('masterEtcdServerPort'))]",
      "[concat('http://', variables('masterPrivateIpAddrs')[4], ':', variables('masterEtcdServerPort'))]"
    ],
    "masterEtcdClientURLs":[
      "[concat('http://', variables('masterPrivateIpAddrs')[0], ':', variables('masterEtcdClientPort'))]",
      "[concat('http://', variables('masterPrivateIpAddrs')[1], ':', variables('masterEtcdClientPort'))]",
      "[concat('http://', variables('masterPrivateIpAddrs')[2], ':', variables('masterEtcdClientPort'))]",
      "[concat('http://', variables('masterPrivateIpAddrs')[3], ':', variables('masterEtcdClientPort'))]",
      "[concat('http://', variables('masterPrivateIpAddrs')[4], ':', variables('masterEtcdClientPort'))]"    
    ],
    "masterEtcdClusterStates": [
      "[concat(variables('masterVMNames')[0], '=', variables('masterEtcdPeerURLs')[0])]",
      "[concat(variables('masterVMNames')[0], '=', variables('masterEtcdPeerURLs')[0], ',', variables('masterVMNames')[1], '=', variables('masterEtcdPeerURLs')[1], ',', variables('masterVMNames')[2], '=', variables('masterEtcdPeerURLs')[2])]",
      "[concat(variables('masterVMNames')[0], '=', variables('masterEtcdPeerURLs')[0], ',', variables('masterVMNames')[1], '=', variables('masterEtcdPeerURLs')[1], ',', variables('masterVMNames')[2], '=', variables('masterEtcdPeerURLs')[2], ',', variables('masterVMNames')[3], '=', variables('masterEtcdPeerURLs')[3], ',', variables('masterVMNames')[4], '=', variables('masterEtcdPeerURLs')[4])]"
    ],
    "subscriptionId": "[subscription().subscriptionId]",
    "tenantId": "[subscription().tenantId]",
    "targetEnvironment": "[parameters('targetEnvironment')]",
    "dockerEngineDownloadRepo": "[parameters('dockerEngineDownloadRepo')]",
    "dockerEngineVersion": "1.12.*"
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
    "kubeBinariesVersion": "[parameters('kubeBinariesVersion')]",
    "agentWindowsPublisher": "MicrosoftWindowsServer",
    "agentWindowsOffer": "WindowsServer",
    "agentWindowsSku": "2016-Datacenter-with-Containers",
    "agentWindowsVersion": "2016.0.20170127",
    "singleQuote": "'",
    "windowsCustomScriptSuffix": " $inputFile = '%SYSTEMDRIVE%\\AzureData\\CustomData.bin' ; $outputFile = '%SYSTEMDRIVE%\\AzureData\\CustomDataSetupScript.ps1' ; Copy-Item $inputFile $outputFile ; Invoke-Expression('{0} {1}' -f $outputFile, $arguments) ; "
{{end}}

    
 
