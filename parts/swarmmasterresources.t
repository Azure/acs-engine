    {
      "apiVersion": "[variables('apiVersionStorage')]",
      "dependsOn": [
        "[concat('Microsoft.Network/publicIPAddresses/', variables('masterPublicIPAddressName'))]"
      ],
      "location": "[variables('location')]",
      "name": "[variables('masterStorageAccountName')]",
      "properties": {
        "accountType": "[variables('vmSizesMap')[variables('masterVMSize')].storageAccountType]"
      },
      "type": "Microsoft.Storage/storageAccounts"
    },
{{if not .MasterProfile.IsCustomVNET}}
    {
      "apiVersion": "[variables('apiVersionDefault')]",
      "location": "[variables('location')]",
      "name": "[variables('virtualNetworkName')]",
      "properties": {
        "addressSpace": {
          "addressPrefixes": [
            {{GetVNETAddressPrefixes}}
          ]
        },
        "subnets": [
          {{GetVNETSubnets false}}
        ]
      },
      "type": "Microsoft.Network/virtualNetworks"
    },
{{end}}
    {
      "apiVersion": "[variables('apiVersionDefault')]",
      "location": "[variables('location')]",
      "name": "[variables('masterAvailabilitySet')]",
      "properties": {},
      "type": "Microsoft.Compute/availabilitySets"
    },
    {
      "apiVersion": "[variables('apiVersionDefault')]",
      "location": "[variables('location')]",
      "name": "[variables('masterPublicIPAddressName')]",
      "properties": {
        "dnsSettings": {
          "domainNameLabel": "[variables('masterEndpointDNSNamePrefix')]"
        },
        "publicIPAllocationMethod": "Dynamic"
      },
      "type": "Microsoft.Network/publicIPAddresses"
    },
    {
      "apiVersion": "[variables('apiVersionDefault')]",
      "dependsOn": [
        "[concat('Microsoft.Network/publicIPAddresses/', variables('masterPublicIPAddressName'))]"
      ],
      "location": "[variables('location')]",
      "name": "[variables('masterLbName')]",
      "properties": {
        "backendAddressPools": [
          {
            "name": "[variables('masterLbBackendPoolName')]"
          }
        ],
        "frontendIPConfigurations": [
          {
            "name": "[variables('masterLbIPConfigName')]",
            "properties": {
              "publicIPAddress": {
                "id": "[resourceId('Microsoft.Network/publicIPAddresses',variables('masterPublicIPAddressName'))]"
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
        "count": "[variables('masterCount')]",
        "name": "masterLbLoopNode"
      },
      "dependsOn": [
        "[variables('masterLbID')]"
      ],
      "location": "[variables('location')]",
      "name": "[concat(variables('masterLbName'), '/', 'SSH-', variables('masterVMNamePrefix'), copyIndex())]",
      "properties": {
        "backendPort": 22,
        "enableFloatingIP": false,
        "frontendIPConfiguration": {
          "id": "[variables('masterLbIPConfigID')]"
        },
        "frontendPort": "[copyIndex(2200)]",
        "protocol": "tcp"
      },
      "type": "Microsoft.Network/loadBalancers/inboundNatRules"
    },
    {
      "apiVersion": "[variables('apiVersionDefault')]",
      "dependsOn": [
        "[variables('masterLbID')]"
      ],
      "location": "[variables('location')]",
      "name": "[concat(variables('masterSshPort22InboundNatRuleNamePrefix'), '0')]",
      "properties": {
        "backendPort": 2222,
        "enableFloatingIP": false,
        "frontendIPConfiguration": {
          "id": "[variables('masterLbIPConfigID')]"
        },
        "frontendPort": "22",
        "protocol": "tcp"
      },
      "type": "Microsoft.Network/loadBalancers/inboundNatRules"
    },
    {
      "apiVersion": "[variables('apiVersionDefault')]",
      "copy": {
        "count": "[variables('masterCount')]",
        "name": "nicLoopNode"
      },
      "dependsOn": [
{{if not .MasterProfile.IsCustomVNET}}
        "[variables('vnetID')]",
{{end}}
        "[variables('masterLbID')]",
        "[concat(variables('masterSshPort22InboundNatRuleIdPrefix'),'0')]",
        "[concat(variables('masterSshInboundNatRuleIdPrefix'),copyIndex())]"        
      ],
      "location": "[variables('location')]",
      "name": "[concat(variables('masterVMNamePrefix'), 'nic-', copyIndex())]",
      "properties": {
        "ipConfigurations": [
          {
            "name": "ipConfigNode",
            "properties": {
              "loadBalancerBackendAddressPools": [
                {
                  "id": "[concat(variables('masterLbID'), '/backendAddressPools/', variables('masterLbBackendPoolName'))]"
                }
              ],
              "loadBalancerInboundNatRules": "[variables('masterLbInboundNatRules')[copyIndex()]]",
              "privateIPAddress": "[concat(variables('masterFirstAddrPrefix'), copyIndex(int(variables('masterFirstAddrOctet4'))))]",
              "privateIPAllocationMethod": "Static",
              "subnet": {
                "id": "[variables('masterVnetSubnetID')]"
              }
            }
          }
        ]
      },
      "type": "Microsoft.Network/networkInterfaces"
    },
    {
      "apiVersion": "[variables('apiVersionDefault')]",
      "copy": {
        "count": "[variables('masterCount')]",
        "name": "vmLoopNode"
      },
      "dependsOn": [
        "[concat('Microsoft.Network/networkInterfaces/', variables('masterVMNamePrefix'), 'nic-', copyIndex())]",
        "[concat('Microsoft.Compute/availabilitySets/',variables('masterAvailabilitySet'))]",
        "[variables('masterStorageAccountName')]"
      ],
      "tags":
      {
        "creationSource" : "[concat('acsengine-', variables('masterVMNamePrefix'), copyIndex())]"
      },
      "location": "[variables('location')]",
      "name": "[concat(variables('masterVMNamePrefix'), copyIndex())]",
      "properties": {
        "availabilitySet": {
          "id": "[resourceId('Microsoft.Compute/availabilitySets',variables('masterAvailabilitySet'))]"
        },
        "hardwareProfile": {
          "vmSize": "[variables('masterVMSize')]"
        },
        "networkProfile": {
          "networkInterfaces": [
            {
              "id": "[resourceId('Microsoft.Network/networkInterfaces',concat(variables('masterVMNamePrefix'), 'nic-', copyIndex()))]"
            }
          ]
        },
        "osProfile": {
          "adminUsername": "[variables('adminUsername')]",
          "computername": "[concat(variables('masterVMNamePrefix'), copyIndex())]",
          {{if .OrchestratorProfile.IsSwarmMode}}
            {{GetMasterSwarmModeCustomData}}
          {{else}}
            {{GetMasterSwarmCustomData}}
          {{end}}
          "linuxConfiguration": {
            "disablePasswordAuthentication": "true",
            "ssh": {
                "publicKeys": [
                    {
                        "keyData": "[variables('sshRSAPublicKey')]",
                        "path": "[variables('sshKeyPath')]"
                    }
                ]
            }
          }
          {{if .LinuxProfile.HasSecrets}}
          ,
          "secrets": "[variables('linuxProfileSecrets')]"
          {{end}}
        },
        "storageProfile": {
          "imageReference": {
            "offer": "[variables('osImageOffer')]",
            "publisher": "[variables('osImagePublisher')]",
            "sku": "[variables('osImageSKU')]",
            "version": "[variables('osImageVersion')]"
          },
          "osDisk": {
            "caching": "ReadWrite",
            "createOption": "FromImage",
            "name": "[concat(variables('masterVMNamePrefix'), copyIndex(),'-osdisk')]",
            "vhd": {
              "uri": "[concat(reference(concat('Microsoft.Storage/storageAccounts/', variables('masterStorageAccountName')), variables('apiVersionStorage')).primaryEndpoints.blob, 'vhds/', variables('masterVMNamePrefix'), copyIndex(), '-osdisk.vhd')]"
            }
          }
        }
      },
      "type": "Microsoft.Compute/virtualMachines"
    },
    {
      "apiVersion": "[variables('apiVersionDefault')]",
      "copy": {
        "count": "[variables('masterCount')]",
        "name": "vmLoopNode"
      },
      "dependsOn": [
          "[concat('Microsoft.Compute/virtualMachines/', concat(variables('masterVMNamePrefix'), copyIndex()))]"
      ],
      "location": "[variables('location')]",
      "name": "[concat(variables('masterVMNamePrefix'), copyIndex(), '/configuremaster')]",
      "properties": {
        "publisher": "Microsoft.OSTCExtensions",
        "settings": {
          "commandToExecute": "[variables('masterCustomScript')]",
          "fileUris": []
        },
        "type": "CustomScriptForLinux",
        "typeHandlerVersion": "1.4"
      },
      "type": "Microsoft.Compute/virtualMachines/extensions"
    }
