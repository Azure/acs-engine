{{if .IsStorageAccount}}    
    {
      "apiVersion": "[variables('apiVersionStorage')]", 
      "copy": {
        "count": "[variables('{{.Name}}StorageAccountsCount')]", 
        "name": "vmLoopNode"
      }, 
      "dependsOn": [
        "[concat('Microsoft.Network/publicIPAddresses/', variables('masterPublicIPAddressName'))]"
      ], 
      "location": "[resourceGroup().location]", 
      "name": "[concat(variables('storageAccountPrefixes')[mod(add(copyIndex(),variables('{{.Name}}StorageAccountOffset')),variables('storageAccountPrefixesCount'))],variables('storageAccountPrefixes')[div(add(copyIndex(),variables('{{.Name}}StorageAccountOffset')),variables('storageAccountPrefixesCount'))],variables('{{.Name}}AccountName'))]", 
      "properties": {
        "accountType": "[variables('vmSizesMap')[variables('{{.Name}}VMSize')].storageAccountType]"
      }, 
      "type": "Microsoft.Storage/storageAccounts"
    },
{{end}}
{{if IsPublic .Ports}}
    {
      "apiVersion": "[variables('apiVersionDefault')]", 
      "location": "[resourceGroup().location]", 
      "name": "[variables('{{.Name}}IPAddressName')]", 
      "properties": {
        "dnsSettings": {
          "domainNameLabel": "[variables('{{.Name}}EndpointDNSNamePrefix')]"
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
            "name": "[variables('{{.Name}}LbBackendPoolName')]"
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
        ], 
        "inboundNatRules": [], 
        "loadBalancingRules": [
          {{(GetLBRules .Name .Ports)}}
        ], 
        "probes": [
          {{(GetProbes .Ports)}}
        ],
        "inboundNatPools": [
          {
            "name": "[concat('RDP-', variables('{{.Name}}VMNamePrefix'))]",
            "properties": {
              "frontendIPConfiguration": {
                "id": "[variables('{{.Name}}LbIPConfigID')]"
              },
              "protocol": "tcp",
              "frontendPortRangeStart": "[variables('{{.Name}}WindowsRDPNatRangeStart')]",
              "frontendPortRangeEnd": "[variables('{{.Name}}WindowsRDPEndRangeStop')]",
              "backendPort": "[variables('agentWindowsBackendPort')]"
            }
          }
        ]
      }, 
      "type": "Microsoft.Network/loadBalancers"
    }, 
{{end}}
    {
{{if .IsManagedDisks}}
      "apiVersion": "[variables('apiManagedDisksVersion')]",
{{else}} 
      "apiVersion": "[variables('apiVersionDefault')]",
{{end}} 
      "dependsOn": [
        "[concat('Microsoft.Network/publicIPAddresses/', variables('masterPublicIPAddressName'))]"
{{if .IsStorageAccount}}
        ,"[concat('Microsoft.Storage/storageAccounts/', variables('storageAccountPrefixes')[mod(add(0,variables('{{.Name}}StorageAccountOffset')),variables('storageAccountPrefixesCount'))],variables('storageAccountPrefixes')[div(add(0,variables('{{.Name}}StorageAccountOffset')),variables('storageAccountPrefixesCount'))],variables('{{.Name}}AccountName'))]", 
        "[concat('Microsoft.Storage/storageAccounts/', variables('storageAccountPrefixes')[mod(add(1,variables('{{.Name}}StorageAccountOffset')),variables('storageAccountPrefixesCount'))],variables('storageAccountPrefixes')[div(add(1,variables('{{.Name}}StorageAccountOffset')),variables('storageAccountPrefixesCount'))],variables('{{.Name}}AccountName'))]", 
		    "[concat('Microsoft.Storage/storageAccounts/', variables('storageAccountPrefixes')[mod(add(2,variables('{{.Name}}StorageAccountOffset')),variables('storageAccountPrefixesCount'))],variables('storageAccountPrefixes')[div(add(2,variables('{{.Name}}StorageAccountOffset')),variables('storageAccountPrefixesCount'))],variables('{{.Name}}AccountName'))]", 
        "[concat('Microsoft.Storage/storageAccounts/', variables('storageAccountPrefixes')[mod(add(3,variables('{{.Name}}StorageAccountOffset')),variables('storageAccountPrefixesCount'))],variables('storageAccountPrefixes')[div(add(3,variables('{{.Name}}StorageAccountOffset')),variables('storageAccountPrefixesCount'))],variables('{{.Name}}AccountName'))]", 
        "[concat('Microsoft.Storage/storageAccounts/', variables('storageAccountPrefixes')[mod(add(4,variables('{{.Name}}StorageAccountOffset')),variables('storageAccountPrefixesCount'))],variables('storageAccountPrefixes')[div(add(4,variables('{{.Name}}StorageAccountOffset')),variables('storageAccountPrefixesCount'))],variables('{{.Name}}AccountName'))]"
{{end}}
{{if not .IsCustomVNET}}
      ,"[variables('vnetID')]"
{{end}}
{{if IsPublic .Ports}} 
       ,"[variables('{{.Name}}LbID')]"
{{end}} 
      ],
      "tags":
      {
        "creationSource" : "[concat('acsengine-', variables('{{.Name}}VMNamePrefix'), '-vmss')]"
      },
      "location": "[resourceGroup().location]", 
      "name": "[concat(variables('{{.Name}}VMNamePrefix'), '-vmss')]", 
      "properties": {
        "upgradePolicy": {
          "mode": "Automatic"
        }, 
        "virtualMachineProfile": {
          "networkProfile": {
            "networkInterfaceConfigurations": [
              {
                "name": "nic", 
                "properties": {
                  "ipConfigurations": [
                    {
                      "name": "nicipconfig", 
                      "properties": {
{{if IsPublic .Ports}}
                        "loadBalancerBackendAddressPools": [
                          {
                            "id": "[concat(variables('{{.Name}}LbID'), '/backendAddressPools/', variables('{{.Name}}LbBackendPoolName'))]"
                          }
                        ],
                        "loadBalancerInboundNatPools": [
                          {
                            "id": "[concat(variables('{{.Name}}LbID'), '/inboundNatPools/', 'RDP-', variables('{{.Name}}VMNamePrefix'))]"
                          }
                        ],
{{end}}
                        "subnet": {
                          "id": "[variables('{{.Name}}VnetSubnetID')]"
                        }
                      }
                    }
                  ], 
                  "primary": "true"
                }
              }
            ]
          }, 
          "osProfile": {
            "computerNamePrefix": "[concat(substring(variables('nameSuffix'), 0, 5), 'acs')]",
            "adminUsername": "[variables('windowsAdminUsername')]",
            "adminPassword": "[variables('windowsAdminPassword')]",
            "customData": "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"
          }, 
          "storageProfile": {
            "imageReference": {
              "publisher": "[variables('agentWindowsPublisher')]",
              "offer": "[variables('agentWindowsOffer')]",
              "sku": "[variables('agentWindowsSku')]",
              "version": "latest"
            }, 
            "osDisk": {
              "caching": "ReadWrite", 
              "createOption": "FromImage"
{{if .IsStorageAccount}} 
              ,"name": "vmssosdisk", 
              "vhdContainers": [
                "[concat(reference(concat('Microsoft.Storage/storageAccounts/',variables('storageAccountPrefixes')[mod(add(0,variables('{{.Name}}StorageAccountOffset')),variables('storageAccountPrefixesCount'))],variables('storageAccountPrefixes')[div(add(0,variables('{{.Name}}StorageAccountOffset')),variables('storageAccountPrefixesCount'))],variables('{{.Name}}AccountName')), variables('apiVersionStorage') ).primaryEndpoints.blob, 'osdisk')]", 
                "[concat(reference(concat('Microsoft.Storage/storageAccounts/',variables('storageAccountPrefixes')[mod(add(1,variables('{{.Name}}StorageAccountOffset')),variables('storageAccountPrefixesCount'))],variables('storageAccountPrefixes')[div(add(1,variables('{{.Name}}StorageAccountOffset')),variables('storageAccountPrefixesCount'))],variables('{{.Name}}AccountName')), variables('apiVersionStorage')).primaryEndpoints.blob, 'osdisk')]", 
                "[concat(reference(concat('Microsoft.Storage/storageAccounts/',variables('storageAccountPrefixes')[mod(add(2,variables('{{.Name}}StorageAccountOffset')),variables('storageAccountPrefixesCount'))],variables('storageAccountPrefixes')[div(add(2,variables('{{.Name}}StorageAccountOffset')),variables('storageAccountPrefixesCount'))],variables('{{.Name}}AccountName')), variables('apiVersionStorage')).primaryEndpoints.blob, 'osdisk')]", 
                "[concat(reference(concat('Microsoft.Storage/storageAccounts/',variables('storageAccountPrefixes')[mod(add(3,variables('{{.Name}}StorageAccountOffset')),variables('storageAccountPrefixesCount'))],variables('storageAccountPrefixes')[div(add(3,variables('{{.Name}}StorageAccountOffset')),variables('storageAccountPrefixesCount'))],variables('{{.Name}}AccountName')), variables('apiVersionStorage')).primaryEndpoints.blob, 'osdisk')]", 
                "[concat(reference(concat('Microsoft.Storage/storageAccounts/',variables('storageAccountPrefixes')[mod(add(4,variables('{{.Name}}StorageAccountOffset')),variables('storageAccountPrefixesCount'))],variables('storageAccountPrefixes')[div(add(4,variables('{{.Name}}StorageAccountOffset')),variables('storageAccountPrefixesCount'))],variables('{{.Name}}AccountName')), variables('apiVersionStorage')).primaryEndpoints.blob, 'osdisk')]"
              ]
{{end}}
            }
          },
          "extensionProfile": {
            "extensions": [
              {
                "name": "vmssCustomScriptExtension",
                "properties": {
                  "publisher": "Microsoft.Compute",
                  "type": "CustomScriptExtension",
                  "typeHandlerVersion": "1.8",
                  "autoUpgradeMinorVersion": true,
                  "settings": {
                    "commandToExecute": "[variables('windowsCustomScript')]"
                  }
                }
              }
            ]
          }
        }
      }, 
      "sku": {
        "capacity": "[variables('{{.Name}}Count')]", 
        "name": "[variables('{{.Name}}VMSize')]", 
        "tier": "Standard"
      }, 
      "type": "Microsoft.Compute/virtualMachineScaleSets"
    }
