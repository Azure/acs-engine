    {
      "apiVersion": "[variables('apiVersionDefault')]", 
      "copy": {
        "count": "[variables('{{.Name}}Count')]", 
        "name": "loop"
      }, 
      "dependsOn": [
      "[concat(variables('{{.Name}}LbID'), '/inboundNatRules/RDP-', variables('{{.Name}}VMNamePrefix'), copyIndex())]"  
{{if .IsCustomVNET}}
      ,"[variables('nsgID')]" 
{{else}}
      ,"[variables('vnetID')]"
{{end}}
      ], 
      "location": "[variables('location')]", 
      "name": "[concat(variables('{{.Name}}VMNamePrefix'), 'nicp-', copyIndex())]", 
      "properties": {
{{if .IsCustomVNET}}                  
	    "networkSecurityGroup": {
		    "id": "[variables('nsgID')]"
	    },
{{end}}
        "ipConfigurations": [
          {
            "name": "ipconfig1", 
            "properties": {
              "privateIPAddress": "[concat('10.240.245.', copyindex(5))]",
              "privateIPAllocationMethod": "Static", 
              "subnet": {
                "id": "[variables('{{.Name}}VnetSubnetID')]"
              },
              "loadBalancerBackendAddressPools": [
                {
                  "id": "[concat(variables('{{.Name}}LbID'), '/backendAddressPools/pool-',variables('{{.Name}}LbName'))]"
                }
              ], 
              "loadBalancerInboundNatRules": [
                {
                  "id": "[concat(variables('{{.Name}}LbID'),'/inboundNatRules/RDP-',variables('{{.Name}}VMNamePrefix'),copyIndex())]"
                }
              ]
            }
          }
        ],
        "enableIPForwarding": true
      }, 
      "type": "Microsoft.Network/networkInterfaces"
    },
    {
      "apiVersion": "[variables('apiVersionDefault')]", 
      "location": "[resourceGroup().location]", 
      "name": "[variables('{{.Name}}IPAddressName')]", 
      "properties": {
        "dnsSettings": {
          "domainNameLabel": "[concat('rdp',uniqueString(variables('masterFqdnPrefix')))]"
        }, 
        "publicIPAllocationMethod": "Dynamic"
      }, 
      "type": "Microsoft.Network/publicIPAddresses"
    },
    {
      "apiVersion": "[variables('apiVersionDefault')]", 
      "dependsOn": [
        "[concat('Microsoft.Network/publicIPAddresses/', variables('{{.Name}}IPAddressName'))]"
      ], 
      "location": "[resourceGroup().location]", 
      "name": "[variables('{{.Name}}LbName')]", 
      "properties": {
        "backendAddressPools": [
          {
            "name": "[concat('pool-',variables('{{.Name}}LbName'))]"
          }
        ], 
        "frontendIPConfigurations": [
          {
            "name": "[variables('{{.Name}}LbIPConfigName')]", 
            "properties": {
              "publicIPAddress": {
                "id": "[resourceId('Microsoft.Network/publicIPAddresses',variables('{{.Name}}IPAddressName'))]"
              }
            }
          }
        ]
      }, 
      "type": "Microsoft.Network/loadBalancers"
    },
    {
      "apiVersion": "[variables('apiVersionDefault')]", 
      "copy": {
        "count": "[variables('{{.Name}}Count')]",
        "name": "loop"
      }, 
      "dependsOn": [
        "[variables('{{.Name}}LbID')]"
      ], 
      "location": "[resourceGroup().location]", 
      "name": "[concat(variables('{{.Name}}LbName'), '/', 'RDP-', variables('{{.Name}}VMNamePrefix'), copyIndex())]", 
      "properties": {
        "backendPort": "[variables('agentWindowsBackendPort')]", 
        "enableFloatingIP": false, 
        "frontendIPConfiguration": {
          "id": "[variables('{{.Name}}LbIPConfigID')]"
        }, 
        "frontendPort": "[copyIndex(variables('agentWindowsBackendPort'))]", 
        "protocol": "tcp"
      }, 
      "type": "Microsoft.Network/loadBalancers/inboundNatRules"
    },
    {
      "apiVersion": "[variables('apiVersionStorage')]", 
      "copy": {
        "count": "[variables('{{.Name}}StorageAccountsCount')]", 
        "name": "loop"
      }, 
      "dependsOn": [
        "[concat('Microsoft.Network/publicIPAddresses/', variables('masterPublicIPAddressName'))]"
      ], 
      "location": "[variables('location')]",
      "name": "[concat(variables('storageAccountPrefixes')[mod(add(copyIndex(),variables('{{.Name}}StorageAccountOffset')),variables('storageAccountPrefixesCount'))],variables('storageAccountPrefixes')[div(add(copyIndex(),variables('{{.Name}}StorageAccountOffset')),variables('storageAccountPrefixesCount'))],variables('{{.Name}}AccountName'))]",
      "properties": {
        "accountType": "[variables('vmSizesMap')[variables('{{.Name}}VMSize')].storageAccountType]"
      }, 
      "type": "Microsoft.Storage/storageAccounts"
    },
{{if .HasDisks}}
    {
      "apiVersion": "[variables('apiVersionStorage')]", 
      "copy": {
        "count": "[variables('{{.Name}}StorageAccountsCount')]", 
        "name": "datadiskLoop"
      }, 
      "dependsOn": [
        "[concat('Microsoft.Network/publicIPAddresses/', variables('masterPublicIPAddressName'))]"
      ], 
      "location": "[variables('location')]", 
      "name": "[concat(variables('storageAccountPrefixes')[mod(add(copyIndex(variables('dataStorageAccountPrefixSeed')),variables('{{.Name}}StorageAccountOffset')),variables('storageAccountPrefixesCount'))],variables('storageAccountPrefixes')[div(add(copyIndex(variables('dataStorageAccountPrefixSeed')),variables('{{.Name}}StorageAccountOffset')),variables('storageAccountPrefixesCount'))],variables('{{.Name}}DataAccountName'))]", 
      "properties": {
        "accountType": "[variables('vmSizesMap')[variables('{{.Name}}VMSize')].storageAccountType]"
      }, 
      "type": "Microsoft.Storage/storageAccounts"
    },
{{end}}
    {
      "apiVersion": "[variables('apiVersionDefault')]", 
      "location": "[variables('location')]",  
      "name": "[variables('{{.Name}}AvailabilitySet')]", 
      "properties": {}, 
      "type": "Microsoft.Compute/availabilitySets"
    },
    {
      "apiVersion": "[variables('apiVersionDefault')]", 
      "copy": {
        "count": "[variables('{{.Name}}Count')]", 
        "name": "vmLoopNode"
      }, 
      "dependsOn": [
        "[concat('Microsoft.Storage/storageAccounts/',variables('storageAccountPrefixes')[mod(add(div(copyIndex(),variables('maxVMsPerStorageAccount')),variables('{{.Name}}StorageAccountOffset')),variables('storageAccountPrefixesCount'))],variables('storageAccountPrefixes')[div(add(div(copyIndex(),variables('maxVMsPerStorageAccount')),variables('{{.Name}}StorageAccountOffset')),variables('storageAccountPrefixesCount'))],variables('{{.Name}}AccountName'))]",
{{if .HasDisks}}
        "[concat('Microsoft.Storage/storageAccounts/',variables('storageAccountPrefixes')[mod(add(add(div(copyIndex(),variables('maxVMsPerStorageAccount')),variables('{{.Name}}StorageAccountOffset')),variables('dataStorageAccountPrefixSeed')),variables('storageAccountPrefixesCount'))],variables('storageAccountPrefixes')[div(add(add(div(copyIndex(),variables('maxVMsPerStorageAccount')),variables('{{.Name}}StorageAccountOffset')),variables('dataStorageAccountPrefixSeed')),variables('storageAccountPrefixesCount'))],variables('{{.Name}}DataAccountName'))]",
{{end}}
        "[concat('Microsoft.Network/networkInterfaces/', variables('{{.Name}}VMNamePrefix'), 'nicp-', copyIndex())]",
        "[concat('Microsoft.Compute/availabilitySets/', variables('{{.Name}}AvailabilitySet'))]"
      ], 
      "tags":
      {
        "creationSource" : "[concat('acsengine-', variables('{{.Name}}VMNamePrefix'), copyIndex())]"
      },
      "location": "[variables('location')]", 
      "name": "[concat(variables('{{.Name}}VMNamePrefix'), copyIndex())]", 
      "properties": {
        "availabilitySet": {
          "id": "[resourceId('Microsoft.Compute/availabilitySets',variables('{{.Name}}AvailabilitySet'))]"
        }, 
        "hardwareProfile": {
          "vmSize": "[variables('{{.Name}}VMSize')]"
        }, 
        "networkProfile": {
          "networkInterfaces": [
            {
              "properties": {
                                "primary": true
                            },
              "id": "[resourceId('Microsoft.Network/networkInterfaces',concat(variables('{{.Name}}VMNamePrefix'), 'nicp-', copyIndex()))]"
            }
          ]
        }, 
        "osProfile": {
          "computername": "[concat(substring(variables('nameSuffix'), 0, 5), 'acs', copyIndex(), add(900,variables('{{.Name}}Index')))]",
          {{GetKubernetesWindowsAgentCustomData}}
          "adminUsername": "[variables('windowsAdminUsername')]",
          "adminPassword": "[variables('windowsAdminPassword')]"
        }, 
        "storageProfile": {
          {{GetDataDisks .}}
          "imageReference": {
            "publisher": "[variables('agentWindowsPublisher')]",
            "offer": "[variables('agentWindowsOffer')]",
            "sku": "[variables('agentWindowsSku')]",
            "version": "latest"
          }, 
          "osDisk": {
            "caching": "ReadWrite", 
            "createOption": "FromImage", 
            "name": "[concat(variables('{{.Name}}VMNamePrefix'), copyIndex(),'-osdisk')]", 
            "vhd": {
              "uri": "[concat(reference(concat('Microsoft.Storage/storageAccounts/',variables('storageAccountPrefixes')[mod(add(div(copyIndex(),variables('maxVMsPerStorageAccount')),variables('{{.Name}}StorageAccountOffset')),variables('storageAccountPrefixesCount'))],variables('storageAccountPrefixes')[div(add(div(copyIndex(),variables('maxVMsPerStorageAccount')),variables('{{.Name}}StorageAccountOffset')),variables('storageAccountPrefixesCount'))],variables('{{.Name}}AccountName')),variables('apiVersionStorage')).primaryEndpoints.blob,'osdisk/', variables('{{.Name}}VMNamePrefix'), copyIndex(), '-osdisk.vhd')]"
            }
          }
        }
      }, 
      "type": "Microsoft.Compute/virtualMachines"
    }, 
    {
      "apiVersion": "[variables('apiVersionDefault')]", 
      "copy": {
        "count": "[variables('{{.Name}}Count')]", 
        "name": "vmLoopNode"
      }, 
      "dependsOn": [
        "[concat('Microsoft.Compute/virtualMachines/', variables('{{.Name}}VMNamePrefix'), copyIndex())]"
      ], 
      "location": "[resourceGroup().location]", 
      "name": "[concat(variables('{{.Name}}VMNamePrefix'), copyIndex(), '/cse')]", 
      "properties": {
        "publisher": "Microsoft.Compute",
        "type": "CustomScriptExtension",
        "typeHandlerVersion": "1.8",
        "autoUpgradeMinorVersion": true,
        "protectedSettings": {
          "commandToExecute": "[concat('powershell.exe -ExecutionPolicy Unrestricted -command \"', '$arguments = ', variables('singleQuote'),'-KubeDnsServiceIp ',variables('kubeDnsServiceIp'),' -MasterFQDNPrefix ',variables('masterFqdnPrefix'),' -Location ',variables('location'),' -AgentKey ',variables('clientPrivateKey'),' -AzureHostname ',variables('{{.Name}}VMNamePrefix'),copyIndex(),variables('singleQuote'), ' ; ', variables('windowsCustomScriptSuffix'), '\" > %SYSTEMDRIVE%\\AzureData\\CustomDataSetupScript.log 2>&1')]"
        }
      }, 
      "type": "Microsoft.Compute/virtualMachines/extensions"
    }
    