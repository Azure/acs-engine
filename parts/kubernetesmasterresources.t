    {
      "apiVersion": "[variables('apiVersionDefault')]",
      "location": "[variables('location')]",
      "name": "[variables('masterAvailabilitySet')]",
      "properties": {},
      "type": "Microsoft.Compute/availabilitySets"
    },
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
      "dependsOn": [
        "[concat('Microsoft.Network/networkSecurityGroups/', variables('nsgName'))]"
{{if not IsVNETIntegrated}}
        ,
        "[concat('Microsoft.Network/routeTables/', variables('routeTableName'))]"
{{end}}
      ],
      "location": "[variables('location')]",
      "name": "[variables('virtualNetworkName')]",
      "properties": {
        "addressSpace": {
          "addressPrefixes": [
            "[variables('vnetCidr')]"
          ]
        },
        "subnets": [
          {
            "name": "[variables('subnetName')]",
            "properties": {
              "addressPrefix": "[variables('subnet')]",
              "networkSecurityGroup": {
                "id": "[variables('nsgID')]"
              }
{{if not IsVNETIntegrated}}
              ,
              "routeTable": {
                "id": "[variables('routeTableID')]"
              }
{{end}}
            }
          }
        ]
      },
      "type": "Microsoft.Network/virtualNetworks"
    },
{{end}}
    {
      "apiVersion": "[variables('apiVersionDefault')]",
      "location": "[variables('location')]",
      "name": "[variables('nsgName')]",
      "properties": {
        "securityRules": [
{{if .HasWindows}}
          {
            "name": "allow_rdp", 
            "properties": {
              "access": "Allow", 
              "description": "Allow RDP traffic to master", 
              "destinationAddressPrefix": "*", 
              "destinationPortRange": "3389-3389", 
              "direction": "Inbound", 
              "priority": 102, 
              "protocol": "Tcp", 
              "sourceAddressPrefix": "*", 
              "sourcePortRange": "*"
            }
          },
{{end}}       
          {
            "name": "allow_ssh",
            "properties": {
              "access": "Allow",
              "description": "Allow SSH traffic to master",
              "destinationAddressPrefix": "*",
              "destinationPortRange": "22-22",
              "direction": "Inbound",
              "priority": 101,
              "protocol": "Tcp",
              "sourceAddressPrefix": "*",
              "sourcePortRange": "*"
            }
          },
          {
            "name": "allow_kube_tls",
            "properties": {
              "access": "Allow",
              "description": "Allow kube-apiserver (tls) traffic to master",
              "destinationAddressPrefix": "*",
              "destinationPortRange": "443-443",
              "direction": "Inbound",
              "priority": 100,
              "protocol": "Tcp",
              "sourceAddressPrefix": "*",
              "sourcePortRange": "*"
            }
          }
        ]
      },
      "type": "Microsoft.Network/networkSecurityGroups"
    },
{{if not IsVNETIntegrated}}
    {
      "apiVersion": "[variables('apiVersionDefault')]",
      "location": "[variables('location')]",
      "name": "[variables('routeTableName')]",
      "type": "Microsoft.Network/routeTables"
    },
{{end}}
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
        ],
        "loadBalancingRules": [
         {
            "name": "LBRuleHTTPS",
            "properties": {
              "frontendIPConfiguration": {
                "id": "[variables('masterLbIPConfigID')]"
              },
              "backendAddressPool": {
                "id": "[concat(variables('masterLbID'), '/backendAddressPools/', variables('masterLbBackendPoolName'))]"
              },
              "protocol": "tcp",
              "frontendPort": 443,
              "backendPort": 443,
              "enableFloatingIP": false,
              "idleTimeoutInMinutes": 5,
              "loadDistribution": "Default",
              "probe": {
                "id": "[concat(variables('masterLbID'),'/probes/tcpHTTPSProbe')]"
              }
            }
          }
        ],
        "probes": [
          {
            "name": "tcpHTTPSProbe",
            "properties": {
              "protocol": "tcp",
              "port": 443,
              "intervalInSeconds": "5",
              "numberOfProbes": "2"
            }
          }
        ]
      },
      "type": "Microsoft.Network/loadBalancers"
    },
    {
      "apiVersion": "[variables('apiVersionDefault')]",
      "dependsOn": [
{{if .MasterProfile.IsCustomVNET}}
        "[variables('nsgID')]"
{{else}}
        "[variables('vnetID')]"
{{end}}
      ],
      "location": "[variables('location')]",
      "name": "[variables('masterInternalLbName')]",
      "properties": {
        "backendAddressPools": [
          {
            "name": "[variables('masterLbBackendPoolName')]"
          }
        ],
        "frontendIPConfigurations": [
          {
            "name": "[variables('masterInternalLbIPConfigName')]",
            "properties": {
              "privateIPAddress": "[variables('masterInternalLbIp')]",
              "privateIPAllocationMethod": "Static",
              "subnet": {
                "id": "[variables('vnetSubnetID')]"
              }
            }
          }
        ],
        "loadBalancingRules": [
          {
            "name": "InternalLBRuleHTTPS",
            "properties": {
              "backendAddressPool": {
                "id": "[concat(variables('masterInternalLbID'), '/backendAddressPools/', variables('masterLbBackendPoolName'))]"
              },
              "backendPort": 4443,
              "enableFloatingIP": false,
              "frontendIPConfiguration": {
                "id": "[variables('masterInternalLbIPConfigID')]"
              },
              "frontendPort": 443,
              "idleTimeoutInMinutes": 5,
              "protocol": "tcp"
            }
          }
        ],
        "probes": [
          {
            "name": "tcpHTTPSProbe",
            "properties": {
              "intervalInSeconds": "5",
              "numberOfProbes": "2",
              "port": 4443,
              "protocol": "tcp"
            }
          }
        ]
      },
      "type": "Microsoft.Network/loadBalancers"
    },
    {
      "apiVersion": "[variables('apiVersionDefault')]",
      "location": "[variables('location')]",
      "name": "[variables('masterPublicIPAddressName')]",
      "properties": {
        "dnsSettings": {
          "domainNameLabel": "[variables('masterFqdnPrefix')]"
        },
        "publicIPAllocationMethod": "Dynamic"
      },
      "type": "Microsoft.Network/publicIPAddresses"
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
        "frontendPort": "[variables('sshNatPorts')[copyIndex()]]",
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
{{if .MasterProfile.IsCustomVNET}}
        "[variables('nsgID')]",
{{else}}
        "[variables('vnetID')]",
{{end}}
        "[concat(variables('masterLbID'),'/inboundNatRules/SSH-',variables('masterVMNamePrefix'),copyIndex())]",
        "[variables('masterInternalLbName')]"
      ],
      "location": "[variables('location')]",
      "name": "[concat(variables('masterVMNamePrefix'), 'nic-', copyIndex())]",
      "properties": {
        "ipConfigurations": [
          {
            "name": "ipconfig1",
            "properties": {
              "loadBalancerBackendAddressPools": [
                {
                  "id": "[concat(variables('masterLbID'), '/backendAddressPools/', variables('masterLbBackendPoolName'))]"
                },
                {
                   "id": "[concat(variables('masterInternalLbID'), '/backendAddressPools/', variables('masterLbBackendPoolName'))]"
                }
              ],
              "loadBalancerInboundNatRules": [
                {
                  "id": "[concat(variables('masterLbID'),'/inboundNatRules/SSH-',variables('masterVMNamePrefix'),copyIndex())]"
                }
              ],
              "privateIPAddress": "[concat(variables('masterFirstAddrPrefix'), copyIndex(int(variables('masterFirstAddrOctet4'))))]",
              "primary": true,
              "privateIPAllocationMethod": "Static",
              "subnet": {
                "id": "[variables('vnetSubnetID')]"
              }
            }
          }
{{if IsVNETIntegrated}}
          {{range $seq := loop 2 .MasterProfile.IPAddressCount}}
          ,
          {
            "name": "ipconfig{{$seq}}",
            "properties": {
              "primary": false,
              "privateIPAllocationMethod": "Dynamic",
              "subnet": {
                "id": "[variables('vnetSubnetID')]"
              }
            }
          }
          {{end}}
{{end}}
        ]
{{if not IsVNETIntegrated}}
        ,
        "enableIPForwarding": true
{{end}}
{{if .MasterProfile.IsCustomVNET}}
        ,"networkSecurityGroup": {
          "id": "[variables('nsgID')]"
        }
{{end}}
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
        "creationSource" : "[concat('acsengine-', variables('masterVMNamePrefix'), copyIndex())]",
        "resourceNameSuffix" : "[variables('nameSuffix')]"
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
              "id": "[resourceId('Microsoft.Network/networkInterfaces',concat(variables('masterVMNamePrefix'),'nic-', copyIndex()))]"
            }
          ]
        },
        "osProfile": {
          "adminUsername": "[variables('username')]",
          "computername": "[concat(variables('masterVMNamePrefix'), copyIndex())]",
          {{GetKubernetesMasterCustomData}}
          "linuxConfiguration": {
            "disablePasswordAuthentication": "true",
            "ssh": {
              "publicKeys": [
                {
                  "keyData": "[variables('sshPublicKeyData')]",
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
          "dataDisks": [
            {
              "createOption": "Empty",
              "diskSizeGB": "128",
              "lun": 0,
              "name": "[concat(variables('masterVMNamePrefix'), copyIndex(),'-etcddisk')]",
              "vhd": {
                "uri": "[concat(reference(concat('Microsoft.Storage/storageAccounts/',variables('masterStorageAccountName')),variables('apiVersionStorage')).primaryEndpoints.blob,'vhds/', variables('masterVMNamePrefix'),copyIndex(),'-etcddisk.vhd')]"
              }
            }
          ],
          "imageReference": {
            "offer": "[variables('osImageOffer')]",
            "publisher": "[variables('osImagePublisher')]",
            "sku": "[variables('osImageSku')]",
            "version": "[variables('osImageVersion')]"
          },
          "osDisk": {
            "caching": "ReadWrite",
            "createOption": "FromImage",
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
      "copy": {
        "count": "[variables('masterCount')]",
        "name": "vmLoopNode"
      },
      "dependsOn": [
        "[concat('Microsoft.Compute/virtualMachines/', variables('masterVMNamePrefix'), copyIndex())]"
      ],
      "location": "[variables('location')]",
      "type": "Microsoft.Compute/virtualMachines/extensions",
      "name": "[concat(variables('masterVMNamePrefix'), copyIndex(),'/cse', copyIndex())]",
      "properties": {
        "publisher": "Microsoft.Azure.Extensions",
        "type": "CustomScript",
        "typeHandlerVersion": "2.0",
        "autoUpgradeMinorVersion": true,
        "settings": {},
        "protectedSettings": {
          "commandToExecute": "[concat('/usr/bin/nohup /bin/bash -c \"/bin/bash /opt/azure/containers/provision.sh ',variables('tenantID'),' ',variables('subscriptionId'),' ',variables('resourceGroup'),' ',variables('location'),' ',variables('subnetName'),' ',variables('nsgName'),' ',variables('virtualNetworkName'),' ',variables('routeTableName'),' ',variables('primaryAvailablitySetName'),' ',variables('servicePrincipalClientId'),' ',variables('servicePrincipalClientSecret'),' ',variables('clientPrivateKey'),' ',variables('targetEnvironment'),' ',variables('networkPolicy'),' ',variables('apiServerPrivateKey'),' ',variables('caCertificate'),' ',variables('masterFqdnPrefix'),' ',variables('kubeConfigCertificate'),' ',variables('kubeConfigPrivateKey'),' ',variables('username'),' >> /var/log/azure/cluster-provision.log 2>&1\"')]"
        }
      }
    }
