    {
      "apiVersion": "[variables('apiVersionStorageManagedDisks')]",
      "location": "[variables('location')]",
      "name": "[variables('bootstrapAvailabilitySet')]",
      "properties": {
        "platformFaultDomainCount": "2",
        "platformUpdateDomainCount": "3",
        "managed": "true"
      },
      "type": "Microsoft.Compute/availabilitySets"
    },
    {
      "apiVersion": "[variables('apiVersionDefault')]",
      "location": "[variables('location')]",
      "name": "[variables('bootstrapPublicIPAddressName')]",
      "properties": {
        "dnsSettings": {
          "domainNameLabel": "[variables('bootstrapEndpointDNSNamePrefix')]"
        },
        "publicIPAllocationMethod": "Dynamic"
      },
      "type": "Microsoft.Network/publicIPAddresses"
    },
    {
      "apiVersion": "[variables('apiVersionDefault')]",
      "dependsOn": [
        "[concat('Microsoft.Network/publicIPAddresses/', variables('bootstrapPublicIPAddressName'))]"
      ],
      "location": "[variables('location')]",
      "name": "[variables('bootstrapLbName')]",
      "properties": {
        "backendAddressPools": [
          {
            "name": "[variables('bootstrapLbBackendPoolName')]"
          }
        ],
        "frontendIPConfigurations": [
          {
            "name": "[variables('bootstrapLbIPConfigName')]",
            "properties": {
              "publicIPAddress": {
                "id": "[resourceId('Microsoft.Network/publicIPAddresses',variables('bootstrapPublicIPAddressName'))]"
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
        "count": "[variables('bootstrapCount')]",
        "name": "bootstrapLbLoopNode"
      },
      "dependsOn": [
        "[variables('bootstrapLbID')]"
      ],
      "location": "[variables('location')]",
      "name": "[concat(variables('bootstrapLbName'), '/', 'SSH-', variables('bootstrapVMNamePrefix'), copyIndex())]",
      "properties": {
        "backendPort": 22,
        "enableFloatingIP": false,
        "frontendIPConfiguration": {
          "id": "[variables('bootstrapLbIPConfigID')]"
        },
        "frontendPort": "[copyIndex(22)]",
        "protocol": "tcp"
      },
      "type": "Microsoft.Network/loadBalancers/inboundNatRules"
    },
    {
      "apiVersion": "[variables('apiVersionDefault')]",
      "copy": {
        "count": "[variables('bootstrapCount')]",
        "name": "bootstrapLbLoopNode"
      },
      "dependsOn": [
        "[variables('bootstrapLbID')]"
      ],
      "location": "[variables('location')]",
      "name": "[concat(variables('bootstrapLbName'), '/', 'bootstrapService-', variables('bootstrapVMNamePrefix'), copyIndex())]",
      "properties": {
        "backendPort": 8086,
        "enableFloatingIP": false,
        "frontendIPConfiguration": {
          "id": "[variables('bootstrapLbIPConfigID')]"
        },
        "frontendPort": "[copyIndex(8086)]",
        "protocol": "tcp"
      },
      "type": "Microsoft.Network/loadBalancers/inboundNatRules"
    },
    {
      "apiVersion": "[variables('apiVersionDefault')]",
      "location": "[variables('location')]",
      "name": "[variables('bootstrapNSGName')]",
      "properties": {
        "securityRules": [
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
            },
            {
                "properties": {
                    "priority": 201,
                    "access": "Allow",
                    "direction": "Inbound",
                    "destinationPortRange": "8086",
                    "sourcePortRange": "*",
                    "destinationAddressPrefix": "*",
                    "protocol": "Tcp",
                    "description": "Allow bootstrap service",
                    "sourceAddressPrefix": "*"
                },
                "name": "Port8086"
            }
        ]
      },
      "type": "Microsoft.Network/networkSecurityGroups"
    },
    {
      "apiVersion": "[variables('apiVersionDefault')]",
      "copy": {
        "count": "[variables('bootstrapCount')]",
        "name": "nicLoopNode"
      },
      "dependsOn": [
        "[variables('bootstrapNSGID')]",
{{if not .MasterProfile.IsCustomVNET}}
        "[variables('vnetID')]",
{{end}}
        "[variables('bootstrapLbID')]",
        "[concat(variables('bootstrapLbID'),'/inboundNatRules/SSH-',variables('bootstrapVMNamePrefix'),copyIndex())]",
        "[concat(variables('bootstrapLbID'),'/inboundNatRules/bootstrapService-',variables('bootstrapVMNamePrefix'),copyIndex())]"
      ],
      "location": "[variables('location')]",
      "name": "[concat(variables('bootstrapVMNamePrefix'), 'nic-', copyIndex())]",
      "properties": {
        "ipConfigurations": [
          {
            "name": "ipConfigNode",
            "properties": {
              "loadBalancerBackendAddressPools": [
                {
                  "id": "[concat(variables('bootstrapLbID'), '/backendAddressPools/', variables('bootstrapLbBackendPoolName'))]"
                }
              ],
              "loadBalancerInboundNatRules": [
                {
                  "id": "[concat(variables('bootstrapLbID'),'/inboundNatRules/SSH-',variables('bootstrapVMNamePrefix'),copyIndex())]"
                },
                {
                  "id": "[concat(variables('bootstrapLbID'),'/inboundNatRules/bootstrapService-',variables('bootstrapVMNamePrefix'),copyIndex())]"
                }
              ],
              "privateIPAddress": "[concat(variables('bootstrapFirstAddrPrefix'), copyIndex(int(variables('bootstrapFirstAddrOctet4'))))]",
              "privateIPAllocationMethod": "Static",
              "subnet": {
                "id": "[variables('masterVnetSubnetID')]"
              }
            }
          }
        ]
        ,"networkSecurityGroup": {
          "id": "[variables('bootstrapNSGID')]"
        }
      },
      "type": "Microsoft.Network/networkInterfaces"
    },
    {
      "apiVersion": "[variables('apiVersionStorageManagedDisks')]",
      "copy": {
        "count": "[variables('bootstrapCount')]",
        "name": "vmLoopNode"
      },
      "dependsOn": [
        "[concat('Microsoft.Network/networkInterfaces/', variables('bootstrapVMNamePrefix'), 'nic-', copyIndex())]",
        "[concat('Microsoft.Compute/availabilitySets/',variables('bootstrapAvailabilitySet'))]",
{{if .MasterProfile.IsStorageAccount}}
        "[variables('masterStorageAccountName')]",
{{end}}
        "[variables('masterStorageAccountExhibitorName')]"
      ],
      "tags":
      {
        "creationSource" : "[concat('acsengine-', variables('bootstrapVMNamePrefix'), copyIndex())]"
      },
      "location": "[variables('location')]",
      "name": "[concat(variables('bootstrapVMNamePrefix'), copyIndex())]",
      "properties": {
        "availabilitySet": {
          "id": "[resourceId('Microsoft.Compute/availabilitySets',variables('bootstrapAvailabilitySet'))]"
        },
        "hardwareProfile": {
          "vmSize": "[variables('bootstrapVMSize')]"
        },
        "networkProfile": {
          "networkInterfaces": [
            {
              "id": "[resourceId('Microsoft.Network/networkInterfaces',concat(variables('bootstrapVMNamePrefix'), 'nic-', copyIndex()))]"
            }
          ]
        },
        "osProfile": {
          "adminUsername": "[variables('adminUsername')]",
          "computername": "[concat(variables('bootstrapVMNamePrefix'), copyIndex())]",
          {{GetDCOSBootstrapCustomData}}
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
            "caching": "ReadWrite"
            ,"createOption": "FromImage"
{{if .MasterProfile.IsStorageAccount}}
            ,"name": "[concat(variables('bootstrapVMNamePrefix'), copyIndex(),'-osdisk')]"
            ,"vhd": {
              "uri": "[concat(reference(concat('Microsoft.Storage/storageAccounts/',variables('masterStorageAccountName')),variables('apiVersionStorage')).primaryEndpoints.blob,'vhds/',variables('bootstrapVMNamePrefix'),copyIndex(),'-osdisk.vhd')]"
            }
{{end}}
{{if ne .OrchestratorProfile.DcosConfig.BootstrapProfile.OSDiskSizeGB 0}}
            ,"diskSizeGB": "60"
{{end}}
          }
        }
      },
      "type": "Microsoft.Compute/virtualMachines"
    },
    {
      "apiVersion": "[variables('apiVersionDefault')]",
      "dependsOn": [
        "[concat('Microsoft.Compute/virtualMachines/', variables('bootstrapVMNamePrefix'), sub(variables('bootstrapCount'), 1))]"
      ],
      "location": "[variables('location')]",
      "name": "[concat(variables('bootstrapVMNamePrefix'), sub(variables('bootstrapCount'), 1), '/bootstrapready')]",
      "properties": {
        "autoUpgradeMinorVersion": true,
        "publisher": "Microsoft.OSTCExtensions",
        "settings": {
          "commandToExecute": "[concat('/bin/bash -c \"until curl -f http://', variables('bootstrapFirstConsecutiveStaticIP'), ':8086/dcos_install.sh > /dev/null; do echo waiting for bootstrap node; sleep 15; done; echo bootstrap node up\"')]"
        },
        "type": "CustomScriptForLinux",
        "typeHandlerVersion": "1.4"
      },
      "type": "Microsoft.Compute/virtualMachines/extensions"
    }{{WriteLinkedTemplatesForExtensions}}
