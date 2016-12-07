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
    "kubectlVersion": "[parameters('kubectlVersion')]",
    "servicePrincipalClientId": "[parameters('servicePrincipalClientId')]",
    "servicePrincipalClientSecret": "[parameters('servicePrincipalClientSecret')]",
    "username": "[parameters('linuxAdminUsername')]",
    "masterFqdnPrefix": "[parameters('masterEndpointDNSNamePrefix')]",
    "masterPrivateIp": "[parameters('firstConsecutiveStaticIP')]",
    "masterVMSize": "[parameters('masterVMSize')]",
    "sshPublicKeyData": "[parameters('sshRSAPublicKey')]",
    "masterCount": {{.MasterProfile.Count}},   
    "apiVersionDefault": "2016-03-30",
    "apiVersionStorage": "2015-06-15",
    "location": "[resourceGroup().location]", 
    "masterAvailabilitySet": "master-availabilityset",
    "storageAccountBaseName": "[uniqueString(concat(variables('masterFqdnPrefix'),resourceGroup().location, variables('orchestratorName')))]",
    "masterStorageAccountName": "[concat(variables('storageAccountBaseName'), 'mstr0')]",
    "nameSuffix": "[parameters('nameSuffix')]", 
    "orchestratorName": "k8s",  
    "osImageOffer": "UbuntuServer", 
    "osImagePublisher": "Canonical", 
    "osImageSKU": "16.04.0-LTS", 
    "osImageVersion": "16.04.201606270",
    "resourceGroup": "[resourceGroup().name]", 
    "routeTableName": "[concat(variables('masterVMNamePrefix'),'routetable')]",
    "routeTableID": "[resourceId('Microsoft.Network/routeTables', variables('routeTableName'))]",
    "sshNatPorts": [22,2201,2202,2203,2204],
    "sshKeyPath": "[concat('/home/',variables('username'),'/.ssh/authorized_keys')]", 
    "storageAccountBaseName": "[uniqueString(concat(variables('masterFqdnPrefix'),resourceGroup().location))]", 
    "storageAccountPrefixes": [ "0", "6", "c", "i", "o", "u", "1", "7", "d", "j", "p", "v", "2", "8", "e", "k", "q", "w", "3", "9", "f", "l", "r", "x", "4", "a", "g", "m", "s", "y", "5", "b", "h", "n", "t", "z" ], 
    "storageAccountPrefixesCount": "[length(variables('storageAccountPrefixes'))]",
    "vmsPerStorageAccount": 20,
{{if AnyAgentHasDisks}}
    "dataStorageAccountPrefixSeed": 97,
{{end}}
{{if .MasterProfile.IsCustomVNET}}
    "vnetSubnetID": "[parameters('masterVnetSubnetID')]",
    "subnetName": "[parameters('masterVnetSubnetID')]",
    "vnetParts": "[split(parameters('masterVnetSubnetID'),'/subnets/')]",
    "virtualNetworkName": "[variables('vnetParts')[0]]",
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
    "nsgName": "[concat(variables('masterVMNamePrefix'), 'nsg')]",
    "nsgID": "[resourceId('Microsoft.Network/networkSecurityGroups',variables('nsgName'))]",
    "primaryAvailablitySetName": "[concat('{{ (index .AgentPoolProfiles 0).Name }}-availabilitySet-',variables('nameSuffix'))]",
    "masterPublicIPAddressName": "[concat(variables('orchestratorName'), '-master-ip-', variables('masterFqdnPrefix'), '-', variables('nameSuffix'))]",
    "masterLbID": "[resourceId('Microsoft.Network/loadBalancers',variables('masterLbName'))]", 
    "masterLbIPConfigID": "[concat(variables('masterLbID'),'/frontendIPConfigurations/', variables('masterLbIPConfigName'))]", 
    "masterLbIPConfigName": "[concat(variables('orchestratorName'), '-master-lbFrontEnd-', variables('nameSuffix'))]",
    "masterLbName": "[concat(variables('orchestratorName'), '-master-lb-', variables('nameSuffix'))]",
    "masterLbBackendPoolName": "[concat(variables('orchestratorName'), '-master-pool-', variables('nameSuffix'))]",
    "masterFirstAddrComment": "these MasterFirstAddrComment are used to place multiple masters consecutively in the address space",
    "masterFirstAddrOctets": "[split(parameters('firstConsecutiveStaticIP'),'.')]",
    "masterFirstAddrOctet4": "[variables('masterFirstAddrOctets')[3]]",
    "masterFirstAddrPrefix": "[concat(variables('masterFirstAddrOctets')[0],'.',variables('masterFirstAddrOctets')[1],'.',variables('masterFirstAddrOctets')[2],'.')]",
    "masterVMNamePrefix": "[concat(variables('orchestratorName'), '-master-', variables('nameSuffix'), '-')]",
    "subscriptionId": "[subscription().subscriptionId]",
    "tenantId": "[subscription().tenantId]"
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


    
 
