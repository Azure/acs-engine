    {
      "apiVersion": "[variables('apiVersionDefault')]", 
      "copy": {
        "count": "[sub(variables('{{.Name}}Count'), variables('{{.Name}}Offset'))]", 
        "name": "loop"
      }, 
      "dependsOn": [
{{if .IsCustomVNET}}
      "[variables('nsgID')]" 
{{else}}
      "[variables('vnetID')]"
{{end}}
      ], 
      "location": "[variables('location')]", 
      "name": "[concat(variables('{{.Name}}VMNamePrefix'), 'nicp-', copyIndex(variables('{{.Name}}Offset')))]", 
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
              "privateIPAllocationMethod": "Dynamic",
              "subnet": {
                "id": "[variables('{{.Name}}VnetSubnetID')]"
             }
            }
          }
        ],
        "enableIPForwarding": true
      }, 
      "type": "Microsoft.Network/networkInterfaces"
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
        "count": "[sub(variables('{{.Name}}Count'), variables('{{.Name}}Offset'))]", 
        "name": "vmLoopNode"
      }, 
      "dependsOn": [
        "[concat('Microsoft.Storage/storageAccounts/',variables('storageAccountPrefixes')[mod(add(div(copyIndex(variables('{{.Name}}Offset')),variables('maxVMsPerStorageAccount')),variables('{{.Name}}StorageAccountOffset')),variables('storageAccountPrefixesCount'))],variables('storageAccountPrefixes')[div(add(div(copyIndex(variables('{{.Name}}Offset')),variables('maxVMsPerStorageAccount')),variables('{{.Name}}StorageAccountOffset')),variables('storageAccountPrefixesCount'))],variables('{{.Name}}AccountName'))]",
{{if .HasDisks}}
        "[concat('Microsoft.Storage/storageAccounts/',variables('storageAccountPrefixes')[mod(add(add(div(copyIndex(variables('{{.Name}}Offset')),variables('maxVMsPerStorageAccount')),variables('{{.Name}}StorageAccountOffset')),variables('dataStorageAccountPrefixSeed')),variables('storageAccountPrefixesCount'))],variables('storageAccountPrefixes')[div(add(add(div(copyIndex(variables('{{.Name}}Offset')),variables('maxVMsPerStorageAccount')),variables('{{.Name}}StorageAccountOffset')),variables('dataStorageAccountPrefixSeed')),variables('storageAccountPrefixesCount'))],variables('{{.Name}}DataAccountName'))]",
{{end}}
        "[concat('Microsoft.Network/networkInterfaces/', variables('{{.Name}}VMNamePrefix'), 'nicp-', copyIndex(variables('{{.Name}}Offset')))]",
        "[concat('Microsoft.Compute/availabilitySets/', variables('{{.Name}}AvailabilitySet'))]"
      ], 
      "tags":
      {
        "creationSource" : "[concat('acsengine-', variables('{{.Name}}VMNamePrefix'), copyIndex(variables('{{.Name}}Offset')))]",
        "resourceNameSuffix" : "[variables('winResourceNamePrefix')]"
      },
      "location": "[variables('location')]", 
      "name": "[concat(variables('{{.Name}}VMNamePrefix'), copyIndex(variables('{{.Name}}Offset')))]", 
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
              "id": "[resourceId('Microsoft.Network/networkInterfaces',concat(variables('{{.Name}}VMNamePrefix'), 'nicp-', copyIndex(variables('{{.Name}}Offset'))))]"
            }
          ]
        }, 
        "osProfile": {
          "computername": "[concat(variables('{{.Name}}VMNamePrefix'), copyIndex(variables('{{.Name}}Offset')))]",
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
            "version": "[variables('agentWindowsVersion')]"
          }, 
          "osDisk": {
            "caching": "ReadWrite", 
            "createOption": "FromImage", 
            "name": "[concat(variables('{{.Name}}VMNamePrefix'), copyIndex(variables('{{.Name}}Offset')),'-osdisk')]", 
            "vhd": {
              "uri": "[concat(reference(concat('Microsoft.Storage/storageAccounts/',variables('storageAccountPrefixes')[mod(add(div(copyIndex(variables('{{.Name}}Offset')),variables('maxVMsPerStorageAccount')),variables('{{.Name}}StorageAccountOffset')),variables('storageAccountPrefixesCount'))],variables('storageAccountPrefixes')[div(add(div(copyIndex(variables('{{.Name}}Offset')),variables('maxVMsPerStorageAccount')),variables('{{.Name}}StorageAccountOffset')),variables('storageAccountPrefixesCount'))],variables('{{.Name}}AccountName')),variables('apiVersionStorage')).primaryEndpoints.blob,'osdisk/', variables('{{.Name}}VMNamePrefix'), copyIndex(variables('{{.Name}}Offset')), '-osdisk.vhd')]"
            }
          }
        }
      }, 
      "type": "Microsoft.Compute/virtualMachines"
    }, 
    {
      "apiVersion": "[variables('apiVersionDefault')]", 
      "copy": {
        "count": "[sub(variables('{{.Name}}Count'), variables('{{.Name}}Offset'))]", 
        "name": "vmLoopNode"
      }, 
      "dependsOn": [
        "[concat('Microsoft.Compute/virtualMachines/', variables('{{.Name}}VMNamePrefix'), copyIndex(variables('{{.Name}}Offset')))]"
      ], 
      "location": "[variables('location')]",
      "name": "[concat(variables('{{.Name}}VMNamePrefix'), copyIndex(variables('{{.Name}}Offset')), '/cse')]",
      "properties": {
        "publisher": "Microsoft.Compute",
        "type": "CustomScriptExtension",
        "typeHandlerVersion": "1.8",
        "autoUpgradeMinorVersion": true,
        "protectedSettings": {
          "commandToExecute": "[concat('start /b powershell.exe -ExecutionPolicy Unrestricted -command \"', '$arguments = ', variables('singleQuote'),'-MasterIP ',variables('masterInternalLbIp'),' -KubeDnsServiceIp ',variables('kubeDnsServiceIp'),' -MasterFQDNPrefix ',variables('masterFqdnPrefix'),' -Location ',variables('location'),' -AgentKey ',variables('clientPrivateKey'),' -AzureHostname ',variables('{{.Name}}VMNamePrefix'),copyIndex(variables('{{.Name}}Offset')),' -AADClientId ',variables('servicePrincipalClientId'),' -AADClientSecret ',variables('servicePrincipalClientSecret'),variables('singleQuote'), ' ; ', variables('windowsCustomScriptSuffix'), '\" > %SYSTEMDRIVE%\\AzureData\\CustomDataSetupScript.log 2>&1')]"
        }
      }, 
      "type": "Microsoft.Compute/virtualMachines/extensions"
    }
    
