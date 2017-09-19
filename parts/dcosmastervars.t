    "adminUsername": "[parameters('linuxAdminUsername')]",
    "targetEnvironment": "[parameters('targetEnvironment')]",
    "maxVMsPerPool": 100,
    "apiVersionDefault": "2016-03-30",
    "apiVersionLinkDefault": "2015-01-01",
    "singleQuote": "'",
    "doubleSingleQuote": "''",
{{if .LinuxProfile.HasSecrets}}
    "linuxProfileSecrets" :
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
      ],
{{end}}
{{if .HasWindows}}
    "windowsAdminUsername": "[parameters('windowsAdminUsername')]",
    "windowsAdminPassword": "[parameters('windowsAdminPassword')]",
    "agentWindowsBackendPort": 3389,
    "agentWindowsPublisher": "MicrosoftWindowsServer",
    "agentWindowsOffer": "WindowsServer",
    "agentWindowsSku": "2016-Datacenter-with-Containers",
    "agentWindowsVersion": "[parameters('agentWindowsVersion')]",
    "dcosWindowsBootstrapURL" : "[parameters('dcosWindowsBootstrapURL')]",
    "windowsCustomScriptSuffix": " $inputFile = '%SYSTEMDRIVE%\\AzureData\\CustomData.bin' ; $outputFile = '%SYSTEMDRIVE%\\AzureData\\dcosWindowsProvision.ps1' ; $inputStream = New-Object System.IO.FileStream $inputFile, ([IO.FileMode]::Open), ([IO.FileAccess]::Read), ([IO.FileShare]::Read) ; $sr = New-Object System.IO.StreamReader(New-Object System.IO.Compression.GZipStream($inputStream, [System.IO.Compression.CompressionMode]::Decompress)) ; $sr.ReadToEnd() | Out-File($outputFile) ; Invoke-Expression('{0} {1}' -f $outputFile, $arguments) ; ",
    "windowsMasterCustomScriptArguments": "[concat('$arguments = ', variables('singleQuote'),'-MasterCount ', variables('masterCount'), ' -firstMasterIP ', parameters('firstConsecutiveStaticIP'), variables('singleQuote'), ' ; ')]",

    "windowsMasterCustomScript": "[concat('powershell.exe -ExecutionPolicy Unrestricted -command \"', variables('windowsMasterCustomScriptArguments'), variables('windowsCustomScriptSuffix'), '\" > %SYSTEMDRIVE%\\AzureData\\dcosWindowsProvision.log 2>&1')]",
{{end}}
    "masterAvailabilitySet": "[concat(variables('orchestratorName'), '-master-availabilitySet-', variables('nameSuffix'))]", 
    "masterCount": {{.MasterProfile.Count}}, 
    "masterEndpointDNSNamePrefix": "[tolower(parameters('masterEndpointDNSNamePrefix'))]",
    "masterHttpSourceAddressPrefix": "{{.MasterProfile.HTTPSourceAddressPrefix}}",
    "masterLbBackendPoolName": "[concat(variables('orchestratorName'), '-master-pool-', variables('nameSuffix'))]",
    "masterLbID": "[resourceId('Microsoft.Network/loadBalancers',variables('masterLbName'))]",
    "masterLbIPConfigID": "[concat(variables('masterLbID'),'/frontendIPConfigurations/', variables('masterLbIPConfigName'))]",
    "masterLbIPConfigName": "[concat(variables('orchestratorName'), '-master-lbFrontEnd-', variables('nameSuffix'))]",
    "masterLbName": "[concat(variables('orchestratorName'), '-master-lb-', variables('nameSuffix'))]",
    "masterNSGID": "[resourceId('Microsoft.Network/networkSecurityGroups',variables('masterNSGName'))]",
    "masterNSGName": "[concat(variables('orchestratorName'), '-master-nsg-', variables('nameSuffix'))]",
    "masterPublicIPAddressName": "[concat(variables('orchestratorName'), '-master-ip-', variables('masterEndpointDNSNamePrefix'), '-', variables('nameSuffix'))]",
    "apiVersionStorage": "2015-06-15",

    "storageAccountBaseName": "[uniqueString(concat(variables('masterEndpointDNSNamePrefix'),variables('location'),variables('orchestratorName')))]", 
    "masterStorageAccountExhibitorName": "[concat(variables('storageAccountBaseName'), 'exhb0')]", 
    "storageAccountType": "Standard_LRS",
{{if .HasStorageAccountDisks}}
    "maxVMsPerStorageAccount": 20,
    "maxStorageAccountsPerAgent": "[div(variables('maxVMsPerPool'),variables('maxVMsPerStorageAccount'))]",
    "dataStorageAccountPrefixSeed": 97,
    "storageAccountPrefixes": [ "0", "6", "c", "i", "o", "u", "1", "7", "d", "j", "p", "v", "2", "8", "e", "k", "q", "w", "3", "9", "f", "l", "r", "x", "4", "a", "g", "m", "s", "y", "5", "b", "h", "n", "t", "z" ],
    "storageAccountPrefixesCount": "[length(variables('storageAccountPrefixes'))]",
    {{GetSizeMap}},
{{else}}
    "storageAccountPrefixes": [],
{{end}}
{{if .HasManagedDisks}}
    "apiVersionStorageManagedDisks": "2016-04-30-preview",
{{end}}
{{if .MasterProfile.IsStorageAccount}}
    "masterStorageAccountName": "[concat(variables('storageAccountBaseName'), 'mstr0')]",
{{end}}
{{if .MasterProfile.IsCustomVNET}}
    "masterVnetSubnetID": "[parameters('masterVnetSubnetID')]",
{{else}}
    "masterSubnet": "[parameters('masterSubnet')]",
    "masterSubnetName": "[concat(variables('orchestratorName'), '-masterSubnet')]",
    "vnetID": "[resourceId('Microsoft.Network/virtualNetworks',variables('virtualNetworkName'))]",
    "masterVnetSubnetID": "[concat(variables('vnetID'),'/subnets/',variables('masterSubnetName'))]",
    "virtualNetworkName": "[concat(variables('orchestratorName'), '-vnet-', variables('nameSuffix'))]",
{{end}}
    "masterFirstAddrOctets": "[split(parameters('firstConsecutiveStaticIP'),'.')]",
    "masterFirstAddrOctet4": "[variables('masterFirstAddrOctets')[3]]",
    "masterFirstAddrPrefix": "[concat(variables('masterFirstAddrOctets')[0],'.',variables('masterFirstAddrOctets')[1],'.',variables('masterFirstAddrOctets')[2],'.')]",
    "masterVMNamePrefix": "[concat(variables('orchestratorName'), '-master-', variables('nameSuffix'), '-')]",
    "masterVMNic": [
      "[concat(variables('masterVMNamePrefix'), 'nic-0')]",
      "[concat(variables('masterVMNamePrefix'), 'nic-1')]",
      "[concat(variables('masterVMNamePrefix'), 'nic-2')]",
      "[concat(variables('masterVMNamePrefix'), 'nic-3')]",
      "[concat(variables('masterVMNamePrefix'), 'nic-4')]",
      "[concat(variables('masterVMNamePrefix'), 'nic-5')]",
      "[concat(variables('masterVMNamePrefix'), 'nic-6')]"
    ],
    "masterVMSize": "[parameters('masterVMSize')]",
    "nameSuffix": "[parameters('nameSuffix')]",
    "oauthEnabled": "{{.MasterProfile.OAuthEnabled}}",
    "orchestratorName": "dcos",
    "osImageOffer": "[parameters('osImageOffer')]", 
    "osImagePublisher": "[parameters('osImagePublisher')]", 
    "osImageSKU": "[parameters('osImageSKU')]", 
    "osImageVersion": "[parameters('osImageVersion')]",
    "sshKeyPath": "[concat('/home/', variables('adminUsername'), '/.ssh/authorized_keys')]",
    "sshRSAPublicKey": "[parameters('sshRSAPublicKey')]",
    "locations": [
         "[resourceGroup().location]",
         "[parameters('location')]"
    ],
    "location": "[variables('locations')[mod(add(2,length(parameters('location'))),add(1,length(parameters('location'))))]]",
{{if IsDCOS19}}
    "masterSshInboundNatRuleIdPrefix": "[concat(variables('masterLbID'),'/inboundNatRules/SSH-',variables('masterVMNamePrefix'))]",
    "masterSshPort22InboundNatRuleIdPrefix": "[concat(variables('masterLbID'),'/inboundNatRules/SSHPort22-',variables('masterVMNamePrefix'))]",
    "masterLbInboundNatRules": [
            [
                {
                    "id": "[concat(variables('masterSshInboundNatRuleIdPrefix'),'0')]"
                },
                {
                    "id": "[concat(variables('masterSshPort22InboundNatRuleIdPrefix'),'0')]"
                }
            ],
            [
                {
                    "id": "[concat(variables('masterSshInboundNatRuleIdPrefix'),'1')]"
                }
            ],
            [
                {
                    "id": "[concat(variables('masterSshInboundNatRuleIdPrefix'),'2')]"
                }
            ],
            [
                {
                    "id": "[concat(variables('masterSshInboundNatRuleIdPrefix'),'3')]"
                }
            ],
            [
                {
                    "id": "[concat(variables('masterSshInboundNatRuleIdPrefix'),'4')]"
                }
            ]
        ],
{{end}}
    "dcosBootstrapURL": "[parameters('dcosBootstrapURL')]"

