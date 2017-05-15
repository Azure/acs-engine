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
    {
      "apiVersion": "[variables('apiVersionStorage')]",
      "dependsOn": [
        "[concat('Microsoft.Network/publicIPAddresses/', variables('masterPublicIPAddressName'))]"
      ],
      "location": "[variables('location')]",
      "name": "[variables('masterStorageAccountExhibitorName')]",
      "properties": {
        "accountType": "Standard_LRS"
      },
      "type": "Microsoft.Storage/storageAccounts"
    },
{{if not .MasterProfile.IsCustomVNET}}
    {
      "apiVersion": "[variables('apiVersionDefault')]",
      "dependsOn": [
          {{GetVNETSubnetDependencies}}
      ],
      "location": "[variables('location')]",
      "name": "[variables('virtualNetworkName')]",
      "properties": {
        "addressSpace": {
          "addressPrefixes": [
            {{GetVNETAddressPrefixes}}
          ]
        },
        "subnets": [
          {{GetVNETSubnets true}}
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
{{if IsDCOS190}}    
    {
      "apiVersion": "[variables('apiVersionDefault')]",
      "dependsOn": [
        "[variables('masterLbID')]"
      ],
      "location": "[resourceGroup().location]",

      "name": "[concat(variables('masterLbName'), '/', 'SSHPort22-', variables('masterVMNamePrefix'), '0')]",
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
{{end}}
    {
      "apiVersion": "[variables('apiVersionDefault')]",
      "location": "[variables('location')]",
      "name": "[variables('masterNSGName')]",
      "properties": {
        "securityRules": [
{{if IsDCOS190}} 
            {
                "properties": {
                    "priority": 201,
                    "access": "Allow",
                    "direction": "Inbound",
                    "destinationPortRange": "2222",
                    "sourcePortRange": "*",
                    "destinationAddressPrefix": "*",
                    "protocol": "Tcp",
                    "description": "Allow SSH",
                    "sourceAddressPrefix": "*"
                },
                "name": "sshPort22"
            },
{{end}}
            {
                "properties": {
                    "priority": 200,
                    "access": "Allow",
                    "direction": "Inbound",
                    "destinationPortRange": "22",
                    "sourcePortRange": "*",
                    "destinationAddressPrefix": "*",
                    "protocol": "Tcp",
                    "description": "Allow SSH",
                    "sourceAddressPrefix": "*"
                },
                "name": "ssh"
            }
        ]
      },
      "type": "Microsoft.Network/networkSecurityGroups"
    },
    {
      "apiVersion": "[variables('apiVersionDefault')]",
      "copy": {
        "count": "[variables('masterCount')]",
        "name": "nicLoopNode"
      },
      "dependsOn": [
        "[variables('masterNSGID')]",
{{if not .MasterProfile.IsCustomVNET}}        
        "[variables('vnetID')]",
{{end}}
        "[variables('masterLbID')]",
        "[concat(variables('masterLbID'),'/inboundNatRules/SSH-',variables('masterVMNamePrefix'),copyIndex())]"
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
{{if IsDCOS190}}
              "loadBalancerInboundNatRules": "[variables('masterLbInboundNatRules')[copyIndex()]]",
{{else}}
              "loadBalancerInboundNatRules": [
                {
                  "id": "[concat(variables('masterLbID'),'/inboundNatRules/SSH-',variables('masterVMNamePrefix'),copyIndex())]"
                }
              ],
{{end}}
              "privateIPAddress": "[concat(variables('masterFirstAddrPrefix'), copyIndex(int(variables('masterFirstAddrOctet4'))))]",
              "privateIPAllocationMethod": "Static",
              "subnet": {
                "id": "[variables('masterVnetSubnetID')]"
              }
            }
          }
        ]
        ,"networkSecurityGroup": {
          "id": "[variables('masterNSGID')]"
        }
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
        "[variables('masterStorageAccountName')]",
        "[variables('masterStorageAccountExhibitorName')]"
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
          {{GetDCOSMasterCustomData}}
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
{{if ne .MasterProfile.OSDiskSizeGB 0}}
            "diskSizeGB": {{.MasterProfile.OSDiskSizeGB}},
{{end}}
            "name": "[concat(variables('masterVMNamePrefix'), copyIndex(),'-osdisk')]",
            "vhd": {
              "uri": "[concat(reference(concat('Microsoft.Storage/storageAccounts/',variables('masterStorageAccountName')),variables('apiVersionStorage')).primaryEndpoints.blob,'vhds/',variables('masterVMNamePrefix'),copyIndex(),'-osdisk.vhd')]"
            }
          }
        }
      },
      "type": "Microsoft.Compute/virtualMachines"
    },
    {
      "apiVersion": "[variables('apiVersionDefault')]",
      "dependsOn": [
        "[concat('Microsoft.Compute/virtualMachines/', variables('masterVMNamePrefix'), sub(variables('masterCount'), 1))]"
      ],
      "location": "[variables('location')]",
      "name": "[concat(variables('masterVMNamePrefix'), sub(variables('masterCount'), 1), '/waitforleader')]",
      "properties": {
        "autoUpgradeMinorVersion": true,
        "publisher": "Microsoft.OSTCExtensions",
        "settings": {
          "commandToExecute": "sh -c 'until ping -c1 leader.mesos;do echo waiting for leader.mesos;sleep 15;done;echo leader.mesos up'"
        },
        "type": "CustomScriptForLinux",
        "typeHandlerVersion": "1.4"
      },
      "type": "Microsoft.Compute/virtualMachines/extensions"
    }
