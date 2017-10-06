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
      "location": "[variables('location')]", 
      "name": "[concat(variables('storageAccountPrefixes')[mod(copyIndex(),variables('storageAccountPrefixesCount'))],variables('storageAccountPrefixes')[div(copyIndex(),variables('storageAccountPrefixesCount'))],variables('storageAccountBaseClassicName'),copyIndex(1))]", 
      "properties": {
        "accountType": "[variables('vmSizesMap')[variables('{{.Name}}VMSize')].storageAccountType]"
      }, 
      "type": "Microsoft.Storage/storageAccounts"
    },
{{end}}
{{if IsPublic .Ports}}
    {
      "apiVersion": "[variables('apiVersionDefault')]", 
      "location": "[variables('location')]", 
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
      "location": "[variables('location')]", 
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
        ]
      }, 
      "type": "Microsoft.Network/loadBalancers"
    }, 
{{end}}
    {
{{if .IsManagedDisks}}
      "apiVersion": "[variables('apiVersionStorageManagedDisks')]",
{{else}} 
      "apiVersion": "[variables('apiVersionDefault')]",
{{end}}  
      "dependsOn": [
        "[concat('Microsoft.Network/publicIPAddresses/', variables('masterPublicIPAddressName'))]"
{{if .IsStorageAccount}}
        ,"[concat('Microsoft.Storage/storageAccounts/', variables('storageAccountPrefixes')[mod(0,variables('storageAccountPrefixesCount'))],variables('storageAccountPrefixes')[div(0,variables('storageAccountPrefixesCount'))],variables('storageAccountBaseClassicName'),1)]",
        "[concat('Microsoft.Storage/storageAccounts/', variables('storageAccountPrefixes')[mod(1,variables('storageAccountPrefixesCount'))],variables('storageAccountPrefixes')[div(1,variables('storageAccountPrefixesCount'))],variables('storageAccountBaseClassicName'),2)]",
        "[concat('Microsoft.Storage/storageAccounts/',variables('storageAccountPrefixes')[mod(2,variables('storageAccountPrefixesCount'))],variables('storageAccountPrefixes')[div(2,variables('storageAccountPrefixesCount'))],variables('storageAccountBaseClassicName'),3)]",
        "[concat('Microsoft.Storage/storageAccounts/', variables('storageAccountPrefixes')[mod(3,variables('storageAccountPrefixesCount'))],variables('storageAccountPrefixes')[div(3,variables('storageAccountPrefixesCount'))],variables('storageAccountBaseClassicName'),4)]",
        "[concat('Microsoft.Storage/storageAccounts/', variables('storageAccountPrefixes')[mod(4,variables('storageAccountPrefixesCount'))],variables('storageAccountPrefixes')[div(4,variables('storageAccountPrefixesCount'))],variables('storageAccountBaseClassicName'),5)]"
 
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
      "location": "[variables('location')]", 
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
            "adminUsername": "[variables('adminUsername')]", 
            "computerNamePrefix": "[variables('{{.Name}}VMNamePrefix')]", 
{{if IsSwarmMode}}
  {{if not .IsRHEL}}
            {{GetAgentSwarmModeCustomData .}} 
  {{end}}
{{else}}
            {{GetAgentSwarmCustomData .}} 
{{end}}
            "linuxConfiguration": {
              "disablePasswordAuthentication": "true", 
              "ssh": {
                "publicKeys": [
                  {
                    "keyData": "[parameters('sshRSAPublicKey')]", 
                    "path": "[variables('sshKeyPath')]"
                  }
                ]
              }
            }
          }, 
          "storageProfile": {
            "imageReference": {
              "offer": "[variables('osImageOffer')]", 
              "publisher": "[variables('osImagePublisher')]", 
              "sku": "[variables('osImageSKU')]", 
              "version": "latest"
            },
            {{GetDataDisks .}}
            "osDisk": {
              "caching": "ReadWrite"
              ,"createOption": "FromImage"
{{if .IsStorageAccount}}
              ,"name": "vmssosdisk"
              ,"vhdContainers": [
               "[concat(reference(concat('Microsoft.Storage/storageAccounts/',variables('storageAccountPrefixes')[mod(0,variables('storageAccountPrefixesCount'))],variables('storageAccountPrefixes')[div(0,variables('storageAccountPrefixesCount'))],variables('storageAccountBaseClassicName'),1), variables('apiVersionStorage') ).primaryEndpoints.blob, 'osdisk')]",
               "[concat(reference(concat('Microsoft.Storage/storageAccounts/',variables('storageAccountPrefixes')[mod(1,variables('storageAccountPrefixesCount'))],variables('storageAccountPrefixes')[div(1,variables('storageAccountPrefixesCount'))],variables('storageAccountBaseClassicName'),2), variables('apiVersionStorage')).primaryEndpoints.blob, 'osdisk')]",
               "[concat(reference(concat('Microsoft.Storage/storageAccounts/',variables('storageAccountPrefixes')[mod(2,variables('storageAccountPrefixesCount'))],variables('storageAccountPrefixes')[div(2,variables('storageAccountPrefixesCount'))],variables('storageAccountBaseClassicName'),3), variables('apiVersionStorage')).primaryEndpoints.blob, 'osdisk')]",
               "[concat(reference(concat('Microsoft.Storage/storageAccounts/',variables('storageAccountPrefixes')[mod(3,variables('storageAccountPrefixesCount'))],variables('storageAccountPrefixes')[div(3,variables('storageAccountPrefixesCount'))],variables('storageAccountBaseClassicName'),4), variables('apiVersionStorage')).primaryEndpoints.blob, 'osdisk')]",
               "[concat(reference(concat('Microsoft.Storage/storageAccounts/',variables('storageAccountPrefixes')[mod(4,variables('storageAccountPrefixesCount'))],variables('storageAccountPrefixes')[div(4,variables('storageAccountPrefixesCount'))],variables('storageAccountBaseClassicName'),5), variables('apiVersionStorage')).primaryEndpoints.blob, 'osdisk')]"
              ]
{{end}}
{{if ne .OSDiskSizeGB 0}}
            ,"diskSizeGB": {{.OSDiskSizeGB}}
{{end}}
            }
          }
{{if .IsRHEL}}
          ,"extensionProfile": {
            "extensions": [
              {
                "name": "configure{{.Name}}",
                "properties": {
                  "publisher": "Microsoft.Azure.Extensions",
                  "settings": {
                    "commandToExecute": "[variables('agentCustomScript')]",
                    "fileUris": [
                      "[concat('{{ GetConfigurationScriptRootURL }}', variables('configureClusterScriptFile'))]"
                    ]
                  },
                  "type": "CustomScript",
                  "typeHandlerVersion": "2.0"
                }
              }
            ]
          }
{{end}}
        }
      }, 
      "sku": {
        "capacity": "[variables('{{.Name}}Count')]", 
        "name": "[variables('{{.Name}}VMSize')]", 
        "tier": "[variables('{{.Name}}VMSizeTier')]"
      }, 
      "type": "Microsoft.Compute/virtualMachineScaleSets"
    }
