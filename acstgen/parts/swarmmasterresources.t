    {
      "apiVersion": "[variables('storageApiVersion')]", 
      "dependsOn": [
        "[concat('Microsoft.Network/publicIPAddresses/', variables('masterPublicIPAddressName'))]"
      ], 
      "location": "[resourceGroup().location]", 
      "name": "[variables('masterStorageAccountName')]", 
      "properties": {
        "accountType": "[variables('vmSizesMap')[variables('masterVMSize')].storageAccountType]"
      }, 
      "type": "Microsoft.Storage/storageAccounts"
    }, 
    {
      "apiVersion": "[variables('networkApiVersion')]", 
      "location": "[resourceGroup().location]", 
      "name": "[variables('virtualNetworkName')]", 
      "properties": {
        "addressSpace": {
          "addressPrefixes": [
            {{GetVNETAddressPrefixes}}
          ]
        }, 
        "subnets": [
          {{GetVNETSubnets}}
        ]
      }, 
      "type": "Microsoft.Network/virtualNetworks"
    }, 
    {
      "apiVersion": "[variables('computeApiVersion')]", 
      "location": "[resourceGroup().location]", 
      "name": "[variables('masterAvailabilitySet')]", 
      "properties": {}, 
      "type": "Microsoft.Compute/availabilitySets"
    }, 
    {
      "apiVersion": "[variables('networkApiVersion')]", 
      "location": "[resourceGroup().location]", 
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
      "apiVersion": "[variables('networkApiVersion')]", 
      "dependsOn": [
        "[concat('Microsoft.Network/publicIPAddresses/', variables('masterPublicIPAddressName'))]"
      ], 
      "location": "[resourceGroup().location]", 
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
      "apiVersion": "[variables('networkApiVersion')]", 
      "copy": {
        "count": "[variables('masterCount')]", 
        "name": "masterLbLoopNode"
      }, 
      "dependsOn": [
        "[variables('masterLbID')]"
      ], 
      "location": "[resourceGroup().location]", 
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
      "apiVersion": "[variables('networkApiVersion')]", 
      "copy": {
        "count": "[variables('masterCount')]", 
        "name": "nicLoopNode"
      }, 
      "dependsOn": [
        "[variables('masterLbID')]", 
        "[variables('vnetID')]", 
        "[concat(variables('masterLbID'),'/inboundNatRules/SSH-',variables('masterVMNamePrefix'),copyIndex())]"
      ], 
      "location": "[resourceGroup().location]", 
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
              "loadBalancerInboundNatRules": [
                {
                  "id": "[concat(variables('masterLbID'),'/inboundNatRules/SSH-',variables('masterVMNamePrefix'),copyIndex())]"
                }
              ], 
              "privateIPAddress": "[concat(split(variables('masterSubnet'),'0/24')[0], copyIndex(variables('masterFirstAddr')))]", 
              "privateIPAllocationMethod": "Static", 
              "subnet": {
                "id": "[variables('masterSubnetRef')]"
              }
            }
          }
        ]
      }, 
      "type": "Microsoft.Network/networkInterfaces"
    }, 
    {
      "apiVersion": "[variables('computeApiVersion')]", 
      "copy": {
        "count": "[variables('masterCount')]", 
        "name": "vmLoopNode"
      }, 
      "dependsOn": [
        "[concat('Microsoft.Network/networkInterfaces/', variables('masterVMNamePrefix'), 'nic-', copyIndex())]", 
        "[concat('Microsoft.Compute/availabilitySets/',variables('masterAvailabilitySet'))]", 
        "[variables('masterStorageAccountName')]"
      ], 
      "location": "[resourceGroup().location]", 
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
          "customData": "[base64({{template "swarmcustomdata.t" .}})]",
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
              "uri": "[concat(reference(concat('Microsoft.Storage/storageAccounts/', variables('masterStorageAccountName')), variables('storageApiVersion')).primaryEndpoints.blob, 'vhds/', variables('masterVMNamePrefix'), copyIndex(), '-osdisk.vhd')]"
            }
          }
        }
      }, 
      "type": "Microsoft.Compute/virtualMachines"
    },
    {
      "apiVersion": "[variables('computeApiVersion')]", 
      "copy": {
        "count": "[variables('masterCount')]", 
        "name": "vmLoopNode"
      }, 
      "dependsOn": [
          "[concat('Microsoft.Compute/virtualMachines/', concat(variables('masterVMNamePrefix'), copyIndex()))]"
      ], 
      "location": "[resourceGroup().location]", 
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