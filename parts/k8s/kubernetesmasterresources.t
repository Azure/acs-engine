{{if .MasterProfile.IsManagedDisks}}
    {
      "apiVersion": "[variables('apiVersionStorageManagedDisks')]",
      "location": "[variables('location')]",
      "name": "[variables('masterAvailabilitySet')]",
      "properties":
        {
            "platformFaultDomainCount": 2,
            "platformUpdateDomainCount": 3,
		        "managed" : true
        },
      "type": "Microsoft.Compute/availabilitySets"
    },
{{else if .MasterProfile.IsStorageAccount}}
    {
      "apiVersion": "[variables('apiVersionDefault')]",
      "location": "[variables('location')]",
      "name": "[variables('masterAvailabilitySet')]",
      "properties": {},
      "type": "Microsoft.Compute/availabilitySets"
    },
    {
      "apiVersion": "[variables('apiVersionStorage')]",
{{if not IsPrivateCluster}}
      "dependsOn": [
        "[concat('Microsoft.Network/publicIPAddresses/', variables('masterPublicIPAddressName'))]"
      ],
{{end}}
      "location": "[variables('location')]",
      "name": "[variables('masterStorageAccountName')]",
      "properties": {
        "accountType": "[variables('vmSizesMap')[parameters('masterVMSize')].storageAccountType]"
      },
      "type": "Microsoft.Storage/storageAccounts"
    },
{{end}}
{{if not .MasterProfile.IsCustomVNET}}
{
      "apiVersion": "[variables('apiVersionDefault')]",
      "dependsOn": [
{{if RequireRouteTable}}
        "[concat('Microsoft.Network/routeTables/', variables('routeTableName'))]"{{if not IsOpenShift}},{{end}}
{{end}}
{{if not IsOpenShift}}
        "[concat('Microsoft.Network/networkSecurityGroups/', variables('nsgName'))]"
{{end}}
      ],
      "location": "[variables('location')]",
      "name": "[variables('virtualNetworkName')]",
      "properties": {
        "addressSpace": {
          "addressPrefixes": [
            "[parameters('vnetCidr')]"
          ]
        },
        "subnets": [
          {
            "name": "[variables('subnetName')]",
            "properties": {
              "addressPrefix": "[parameters('masterSubnet')]"
{{if not IsOpenShift}}
              ,
              "networkSecurityGroup": {
                "id": "[variables('nsgID')]"
              }
{{end}}
{{if RequireRouteTable}}
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
              "destinationPortRange": {{if IsOpenShift}}"8443-8443"{{else}}"443-443"{{end}},
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
{{if RequireRouteTable}}
    {
      "apiVersion": "[variables('apiVersionDefault')]",
      "location": "[variables('location')]",
      "name": "[variables('routeTableName')]",
      "type": "Microsoft.Network/routeTables"
    },
{{end}}
{{if not IsPrivateCluster}}
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
              "frontendPort": {{if IsOpenShift}}8443{{else}}443{{end}},
              "backendPort": {{if IsOpenShift}}8443{{else}}443{{end}},
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
              "port": {{if IsOpenShift}}8443{{else}}443{{end}},
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
      "copy": {
        "count": "[sub(variables('masterCount'), variables('masterOffset'))]",
        "name": "masterLbLoopNode"
      },
      "dependsOn": [
        "[variables('masterLbID')]"
      ],
      "location": "[variables('location')]",
      "name": "[concat(variables('masterLbName'), '/', 'SSH-', variables('masterVMNamePrefix'), copyIndex(variables('masterOffset')))]",
      "properties": {
        "backendPort": 22,
        "enableFloatingIP": false,
        "frontendIPConfiguration": {
          "id": "[variables('masterLbIPConfigID')]"
        },
        "frontendPort": "[variables('sshNatPorts')[copyIndex(variables('masterOffset'))]]",
        "protocol": "tcp"
      },
      "type": "Microsoft.Network/loadBalancers/inboundNatRules"
    },
    {
      "apiVersion": "[variables('apiVersionDefault')]",
      "copy": {
        "count": "[sub(variables('masterCount'), variables('masterOffset'))]",
        "name": "nicLoopNode"
      },
      "dependsOn": [
{{if not IsOpenShift}}
{{if .MasterProfile.IsCustomVNET}}
        "[variables('nsgID')]",
{{else}}
        "[variables('vnetID')]",
{{end}}
{{else}}
        "[variables('nsgID')]",
{{if not .MasterProfile.IsCustomVNET}}
        "[variables('vnetID')]",
{{end}}
{{end}}
        "[concat(variables('masterLbID'),'/inboundNatRules/SSH-',variables('masterVMNamePrefix'),copyIndex(variables('masterOffset')))]"
{{if gt .MasterProfile.Count 1}}
        ,"[variables('masterInternalLbName')]"
{{end}}
      ],
      "location": "[variables('location')]",
      "name": "[concat(variables('masterVMNamePrefix'), 'nic-', copyIndex(variables('masterOffset')))]",
      "properties": {
        "ipConfigurations": [
          {
            "name": "ipconfig1",
            "properties": {
              "loadBalancerBackendAddressPools": [
                {
                  "id": "[concat(variables('masterLbID'), '/backendAddressPools/', variables('masterLbBackendPoolName'))]"
                }
{{if gt .MasterProfile.Count 1}}
                ,
                {
                   "id": "[concat(variables('masterInternalLbID'), '/backendAddressPools/', variables('masterLbBackendPoolName'))]"
                }
{{end}}
              ],
              "loadBalancerInboundNatRules": [
                {
                  "id": "[concat(variables('masterLbID'),'/inboundNatRules/SSH-',variables('masterVMNamePrefix'),copyIndex(variables('masterOffset')))]"
                }
              ],
              "privateIPAddress": "[variables('masterPrivateIpAddrs')[copyIndex(variables('masterOffset'))]]",
              "primary": true,
              "privateIPAllocationMethod": "Static",
              "subnet": {
                "id": "[variables('vnetSubnetID')]"
              }
            }
          }
{{if IsAzureCNI}}
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
{{if not IsAzureCNI}}
        ,
        "enableIPForwarding": true
{{end}}
{{if HasCustomNodesDNS}}
 ,"dnsSettings": {
          "dnsServers": [
              "[parameters('dnsServer')]"
          ]
      }
{{end}}
{{if or .MasterProfile.IsCustomVNET IsOpenShift}}
        ,"networkSecurityGroup": {
          "id": "[variables('nsgID')]"
        }
{{end}}
      },
      "type": "Microsoft.Network/networkInterfaces"
    },
{{else}}
      {
        "apiVersion": "[variables('apiVersionDefault')]",
        "copy": {
          "count": "[sub(variables('masterCount'), variables('masterOffset'))]",
          "name": "nicLoopNode"
        },
        "dependsOn": [
  {{if not IsOpenShift}}
  {{if .MasterProfile.IsCustomVNET}}
          "[variables('nsgID')]"
  {{else}}
          "[variables('vnetID')]"
  {{end}}
  {{else}}
          "[variables('nsgID')]"
  {{if not .MasterProfile.IsCustomVNET}}
          ,"[variables('vnetID')]"
  {{end}}
  {{end}}
  {{if gt .MasterProfile.Count 1}}
          ,"[variables('masterInternalLbName')]"
  {{end}}
        ],
        "location": "[variables('location')]",
        "name": "[concat(variables('masterVMNamePrefix'), 'nic-', copyIndex(variables('masterOffset')))]",
        "properties": {
          "ipConfigurations": [
            {
              "name": "ipconfig1",
              "properties": {
                "loadBalancerBackendAddressPools": [
  {{if gt .MasterProfile.Count 1}}
                  {
                    "id": "[concat(variables('masterInternalLbID'), '/backendAddressPools/', variables('masterLbBackendPoolName'))]"
                  }
  {{end}}
                ],
                "loadBalancerInboundNatRules": [
                ],
                "privateIPAddress": "[variables('masterPrivateIpAddrs')[copyIndex(variables('masterOffset'))]]",
                "primary": true,
                "privateIPAllocationMethod": "Static",
                "subnet": {
                  "id": "[variables('vnetSubnetID')]"
                }
              }
            }
  {{if IsAzureCNI}}
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
  {{if not IsAzureCNI}}
          ,
          "enableIPForwarding": true
  {{end}}
  {{if HasCustomNodesDNS}}
   ,"dnsSettings": {
          "dnsServers": [
              "[parameters('dnsServer')]"
          ]
      }
  {{end}}
  {{if or .MasterProfile.IsCustomVNET IsOpenShift}}
          ,"networkSecurityGroup": {
            "id": "[variables('nsgID')]"
          }
  {{end}}
        },
        "type": "Microsoft.Network/networkInterfaces"
      },
  {{if ProvisionJumpbox}}
    {
      "type": "Microsoft.Compute/virtualMachines",
      "name": "[parameters('jumpboxVMName')]",
      {{if JumpboxIsManagedDisks}}
      "apiVersion": "[variables('apiVersionStorageManagedDisks')]",
      {{else}}
      "apiVersion": "[variables('apiVersionDefault')]",
      {{end}}
      "location": "[variables('location')]",
      "properties": {
          "osProfile": {
            {{GetKubernetesJumpboxCustomData .}}
              "computerName": "[parameters('jumpboxVMName')]",
              "adminUsername": "[parameters('jumpboxUsername')]",
              "linuxConfiguration": {
                  "disablePasswordAuthentication": true,
                  "ssh": {
                      "publicKeys": [
                          {
                              "path": "[concat('/home/', parameters('jumpboxUsername'), '/.ssh/authorized_keys')]",
                              "keyData": "[parameters('jumpboxPublicKey')]"
                          }
                      ]
                  }
              }
          },
          "hardwareProfile": {
              "vmSize": "[parameters('jumpboxVMSize')]"
          },
          "storageProfile": {
              "imageReference": {
                  "publisher": "Canonical",
                  "offer": "UbuntuServer",
                  "sku": "16.04-LTS",
                  "version": "latest"
              },
            {{if JumpboxIsManagedDisks}}
              "osDisk": {
                  "createOption": "FromImage",
                  "diskSizeGB": "[parameters('jumpboxOSDiskSizeGB')]",
                  "managedDisk": {
                      "storageAccountType": "[variables('vmSizesMap')[parameters('jumpboxVMSize')].storageAccountType]"
                  }
              },
            {{else}}
              "osDisk": {
                "createOption": "fromImage",
                "vhd": {
                    "uri": "[concat(reference(concat('Microsoft.Storage/storageAccounts/',variables('jumpboxStorageAccountName')),variables('apiVersionStorage')).primaryEndpoints.blob,'vhds/',parameters('jumpboxVMName'),'jumpboxdisk.vhd')]"
                },
                "name": "[variables('jumpboxOSDiskName')]"
              },
            {{end}}
          "dataDisks": []
          },
          "networkProfile": {
              "networkInterfaces": [
                  {
                      "id": "[resourceId('Microsoft.Network/networkInterfaces', variables('jumpboxNetworkInterfaceName'))]"
                  }
              ]
          }
        },
        "dependsOn": [
            "[concat('Microsoft.Network/networkInterfaces/', variables('jumpboxNetworkInterfaceName'))]"
        ]
    },
    {{if not JumpboxIsManagedDisks}}
    {
            "type": "Microsoft.Storage/storageAccounts",
            "name": "[variables('jumpboxStorageAccountName')]",
            "apiVersion": "[variables('apiVersionStorage')]",
            "location": "[variables('location')]",
            "properties": {
                "accountType": "[variables('vmSizesMap')[parameters('jumpboxVMSize')].storageAccountType]"
            }
    },
    {{end}}
    {
      "type": "Microsoft.Network/networkSecurityGroups",
      "name": "[variables('jumpboxNetworkSecurityGroupName')]",
      "apiVersion": "[variables('apiVersionDefault')]",
      "location": "[variables('location')]",
      "properties": {
          "securityRules": [
              {
                  "name": "default-allow-ssh",
                  "properties": {
                      "priority": 1000,
                      "protocol": "Tcp",
                      "access": "Allow",
                      "direction": "Inbound",
                      "sourceAddressPrefix": "*",
                      "sourcePortRange": "*",
                      "destinationAddressPrefix": "*",
                      "destinationPortRange": "22"
                  }
              }
          ]
      }
    },
    {
      "type": "Microsoft.Network/publicIpAddresses",
      "sku": {
          "name": "Basic"
      },
      "name": "[variables('jumpboxPublicIpAddressName')]",
      "apiVersion": "[variables('apiVersionDefault')]",
      "location": "[variables('location')]",
      "properties": {
          "dnsSettings": {
            "domainNameLabel": "[variables('masterFqdnPrefix')]"
          },
          "publicIpAllocationMethod": "Dynamic"
      }
    },
    {
      "type": "Microsoft.Network/networkInterfaces",
      "name": "[variables('jumpboxNetworkInterfaceName')]",
      "apiVersion": "[variables('apiVersionDefault')]",
      "location": "[variables('location')]",
      "properties": {
          "ipConfigurations": [
              {
                  "name": "ipconfig1",
                  "properties": {
                      "subnet": {
                          "id": "[variables('vnetSubnetID')]"
                      },
                      "primary": true,
                      "privateIPAllocationMethod": "Dynamic",
                      "publicIpAddress": {
                          "id": "[resourceId('Microsoft.Network/publicIpAddresses', variables('jumpboxPublicIpAddressName'))]"
                      }
                  }
              }
          ],
          "networkSecurityGroup": {
              "id": "[resourceId('Microsoft.Network/networkSecurityGroups', variables('jumpboxNetworkSecurityGroupName'))]"
          }
      },
      "dependsOn": [
          "[concat('Microsoft.Network/publicIpAddresses/', variables('jumpboxPublicIpAddressName'))]",
          "[concat('Microsoft.Network/networkSecurityGroups/', variables('jumpboxNetworkSecurityGroupName'))]"
          {{if not .MasterProfile.IsCustomVNET}}
            ,"[variables('vnetID')]"
          {{end}}
      ]
    },
  {{end}}
{{end}}
{{if gt .MasterProfile.Count 1}}
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
              "privateIPAddress": "[variables('kubernetesAPIServerIP')]",
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
              "backendPort": {{if IsOpenShift}}8443{{else}}4443{{end}},
              "enableFloatingIP": false,
              "frontendIPConfiguration": {
                "id": "[variables('masterInternalLbIPConfigID')]"
              },
              "frontendPort": {{if IsOpenShift}}8443{{else}}443{{end}},
              "idleTimeoutInMinutes": 5,
              "protocol": "tcp",
              "probe": {
                "id": "[concat(variables('masterInternalLbID'),'/probes/tcpHTTPSProbe')]"
              }
            }
          }
        ],
        "probes": [
          {
            "name": "tcpHTTPSProbe",
            "properties": {
              "intervalInSeconds": "5",
              "numberOfProbes": "2",
              "port": {{if IsOpenShift}}8443{{else}}4443{{end}},
              "protocol": "tcp"
            }
          }
        ]
      },
      "type": "Microsoft.Network/loadBalancers"
    },
{{end}}
{{if EnableEncryptionWithExternalKms}}
     {
       "type": "Microsoft.Storage/storageAccounts",
       "name": "[variables('clusterKeyVaultName')]",
       "apiVersion": "[variables('apiVersionStorage')]",
       "location": "[variables('location')]",
       "properties": {
         "accountType": "Standard_LRS"
       }
     },
     {
       "type": "Microsoft.KeyVault/vaults",
       "name": "[variables('clusterKeyVaultName')]",
       "apiVersion": "[variables('apiVersionKeyVault')]",
       "location": "[variables('location')]",
       {{ if UseManagedIdentity}}
       "dependsOn": 
       [
          {{$max := .MasterProfile.Count}}
          {{$c := subtract $max 1}}
          {{range $i := loop 0 $max}}
            {{if (lt $i $c)}}
                "[concat('Microsoft.Compute/virtualMachines/', variables('masterVMNamePrefix'), '{{$i}}')]",
                "[concat('Microsoft.Authorization/roleAssignments/', guid(concat('Microsoft.Compute/virtualMachines/', variables('masterVMNamePrefix'), '{{$i}}', 'vmidentity')))]",
            {{else}}
                {{ if (lt $i $max)}}
                "[concat('Microsoft.Compute/virtualMachines/', variables('masterVMNamePrefix'), '{{$i}}')]",
                "[concat('Microsoft.Authorization/roleAssignments/', guid(concat('Microsoft.Compute/virtualMachines/', variables('masterVMNamePrefix'), '{{$i}}', 'vmidentity')))]"
                {{end}}
            {{end}}
          {{end}}
        ],
       {{end}}
       "properties": {
         "enabledForDeployment": "false",
         "enabledForDiskEncryption": "false",
         "enabledForTemplateDeployment": "false",
         "tenantId": "[variables('tenantID')]",
 {{if not UseManagedIdentity}}
         "accessPolicies": [
           {
             "tenantId": "[variables('tenantID')]",
             "objectId": "[parameters('servicePrincipalObjectId')]",
             "permissions": {
               "keys": ["create", "encrypt", "decrypt", "get", "list"]
             }
           }
         ],
 {{else}}
         "accessPolicies": 
         [
          {{$max := .MasterProfile.Count}}
          {{$c := subtract $max 1}}
          {{range $i := loop 0 $max}}
            {{if (lt $i $c)}}
            {
                "objectId": "[reference(concat('Microsoft.Compute/virtualMachines/', variables('masterVMNamePrefix'), '{{$i}}'), '2017-03-30', 'Full').identity.principalId]",
                "permissions": {
                "keys": [
                    "create",
                    "encrypt",
                    "decrypt",
                    "get",
                    "list"
                ]
                },
                "tenantId": "[variables('tenantID')]"
            },
            {{else}}
                {{ if (lt $i $max)}}
                {
                    "objectId": "[reference(concat('Microsoft.Compute/virtualMachines/', variables('masterVMNamePrefix'), '{{$i}}'), '2017-03-30', 'Full').identity.principalId]",
                    "permissions": {
                    "keys": [
                        "create",
                        "encrypt",
                        "decrypt",
                        "get",
                        "list"
                    ]
                    },
                    "tenantId": "[variables('tenantID')]"
                }
                {{end}}
            {{end}}
          {{end}}
         ],
 {{end}}
         "sku": {
           "name": "[parameters('clusterKeyVaultSku')]",
           "family": "A"
         }
       }
     },
 {{end}}
    {
    {{if UserAssignedIDEnabled}}
      "apiVersion": "[variables('apiVersionUserMSI')]",
    {{else}}
    {{if .MasterProfile.IsManagedDisks}}
      "apiVersion": "[variables('apiVersionStorageManagedDisks')]",
    {{else}}
      "apiVersion": "[variables('apiVersionDefault')]",
    {{end}}
    {{end}}
      "copy": {
        "count": "[sub(variables('masterCount'), variables('masterOffset'))]",
        "name": "vmLoopNode"
      },
      "dependsOn": [
        "[concat('Microsoft.Network/networkInterfaces/', variables('masterVMNamePrefix'), 'nic-', copyIndex(variables('masterOffset')))]"
        ,"[concat('Microsoft.Compute/availabilitySets/',variables('masterAvailabilitySet'))]"
{{if .MasterProfile.IsStorageAccount}}
        ,"[variables('masterStorageAccountName')]"
{{end}}
      ],
      "tags":
      {
        "creationSource" : "[concat(parameters('generatorCode'), '-', variables('masterVMNamePrefix'), copyIndex(variables('masterOffset')))]",
        "resourceNameSuffix" : "[parameters('nameSuffix')]",
        "orchestrator" : "[variables('orchestratorNameVersionTag')]",
        "acsengineVersion" : "[parameters('acsengineVersion')]",
        "poolName" : "master"
      },
      "location": "[variables('location')]",
      "name": "[concat(variables('masterVMNamePrefix'), copyIndex(variables('masterOffset')))]",
      {{if UseManagedIdentity}}
      {{if UserAssignedIDEnabled}}
      "identity": {
        "type": "userAssigned",
        "userAssignedIdentities": {
          "[resourceId('Microsoft.ManagedIdentity/userAssignedIdentities/', variables('userAssignedID'))]":{}
        }
      },
      {{else}}
      "identity": {
        "type": "systemAssigned"
      },
      {{end}}
      {{end}}
      {{if and IsOpenShift (not UseMasterCustomImage)}}
      "plan": {
        "name": "[parameters('osImageSku')]",
        "publisher": "[parameters('osImagePublisher')]",
        "product": "[parameters('osImageOffer')]"
      },
      {{end}}
      "properties": {
        "availabilitySet": {
          "id": "[resourceId('Microsoft.Compute/availabilitySets',variables('masterAvailabilitySet'))]"
        },
        "hardwareProfile": {
          "vmSize": "[parameters('masterVMSize')]"
        },
        "networkProfile": {
          "networkInterfaces": [
            {
              "id": "[resourceId('Microsoft.Network/networkInterfaces',concat(variables('masterVMNamePrefix'),'nic-', copyIndex(variables('masterOffset'))))]"
            }
          ]
        },
        "osProfile": {
          "adminUsername": "[parameters('linuxAdminUsername')]",
          "computername": "[concat(variables('masterVMNamePrefix'), copyIndex(variables('masterOffset')))]",
          {{if not IsOpenShift}}
          {{GetKubernetesMasterCustomData .}}
          {{end}}
          "linuxConfiguration": {
            "disablePasswordAuthentication": true,
            "ssh": {
              "publicKeys": [
                {
                  "keyData": "[parameters('sshRSAPublicKey')]",
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
          {{if and (not UseMasterCustomImage) (not IsOpenShift)}}
          "dataDisks": [
            {
              "createOption": "Empty"
              ,"diskSizeGB": "[parameters('etcdDiskSizeGB')]"
              ,"lun": 0
              ,"name": "[concat(variables('masterVMNamePrefix'), copyIndex(variables('masterOffset')),'-etcddisk')]"
              {{if .MasterProfile.IsStorageAccount}}
              ,"vhd": {
                "uri": "[concat(reference(concat('Microsoft.Storage/storageAccounts/',variables('masterStorageAccountName')),variables('apiVersionStorage')).primaryEndpoints.blob,'vhds/', variables('masterVMNamePrefix'),copyIndex(variables('masterOffset')),'-etcddisk.vhd')]"
              }
              {{end}}
            }
          ],
          {{end}}
          "imageReference": {
            {{if UseMasterCustomImage}}
            "id": "[resourceId(parameters('osImageResourceGroup'), 'Microsoft.Compute/images', parameters('osImageName'))]"
            {{else}}
            "offer": "[parameters('osImageOffer')]",
            "publisher": "[parameters('osImagePublisher')]",
            "sku": "[parameters('osImageSku')]",
            "version": "[parameters('osImageVersion')]"
            {{end}}
          },
          "osDisk": {
            "caching": "ReadWrite"
            ,"createOption": "FromImage"
{{if .MasterProfile.IsStorageAccount}}
            ,"name": "[concat(variables('masterVMNamePrefix'), copyIndex(variables('masterOffset')),'-osdisk')]"
            ,"vhd": {
              "uri": "[concat(reference(concat('Microsoft.Storage/storageAccounts/',variables('masterStorageAccountName')),variables('apiVersionStorage')).primaryEndpoints.blob,'vhds/',variables('masterVMNamePrefix'),copyIndex(variables('masterOffset')),'-osdisk.vhd')]"
            }
{{end}}
{{if ne .MasterProfile.OSDiskSizeGB 0}}
            ,"diskSizeGB": {{.MasterProfile.OSDiskSizeGB}}
{{end}}

          }
        }
      },
      "type": "Microsoft.Compute/virtualMachines"
    },
    {{if UseManagedIdentity}}
    {{if (not UserAssignedIDEnabled)}}
    {
      "apiVersion": "2014-10-01-preview",
      "copy": {
         "count": "[variables('masterCount')]",
         "name": "vmLoopNode"
       },
      "name": "[guid(concat('Microsoft.Compute/virtualMachines/', variables('masterVMNamePrefix'), copyIndex(),'vmidentity'))]",
      "type": "Microsoft.Authorization/roleAssignments",
      "properties": {
        "roleDefinitionId": "[variables('contributorRoleDefinitionId')]",
        "principalId": "[reference(concat('Microsoft.Compute/virtualMachines/', variables('masterVMNamePrefix'), copyIndex()), '2017-03-30', 'Full').identity.principalId]"
      }
    },
    {{end}}
     {
       "type": "Microsoft.Compute/virtualMachines/extensions",
       "name": "[concat(variables('masterVMNamePrefix'), copyIndex(), '/ManagedIdentityExtension')]",
       "copy": {
         "count": "[variables('masterCount')]",
         "name": "vmLoopNode"
       },
       "apiVersion": "2015-05-01-preview",
       "location": "[resourceGroup().location]",
       {{if (not UserAssignedIDEnabled)}}
       "dependsOn": [
         "[concat('Microsoft.Compute/virtualMachines/', variables('masterVMNamePrefix'), copyIndex())]",
         "[concat('Microsoft.Authorization/roleAssignments/', guid(concat('Microsoft.Compute/virtualMachines/', variables('masterVMNamePrefix'), copyIndex(), 'vmidentity')))]"
       ],
       {{else}}
       "dependsOn": [
        "[concat('Microsoft.Compute/virtualMachines/', variables('masterVMNamePrefix'), copyIndex())]"
       ],
       {{end}}
       "properties": {
         "publisher": "Microsoft.ManagedIdentity",
         "type": "ManagedIdentityExtensionForLinux",
         "typeHandlerVersion": "1.0",
         "autoUpgradeMinorVersion": true,
         "settings": {
           "port": 50343
         },
         "protectedSettings": {}
       }
     },
     {{end}}
    {
      "apiVersion": "[variables('apiVersionDefault')]",
      "copy": {
        "count": "[sub(variables('masterCount'), variables('masterOffset'))]",
        "name": "vmLoopNode"
      },
      "dependsOn": [
        {{if UseManagedIdentity}}
        "[concat('Microsoft.Compute/virtualMachines/', variables('masterVMNamePrefix'), copyIndex(variables('masterOffset')), '/extensions/ManagedIdentityExtension')]"
        {{else}}
        "[concat('Microsoft.Compute/virtualMachines/', variables('masterVMNamePrefix'), copyIndex(variables('masterOffset')))]"
        {{end}}
      ],
      "location": "[variables('location')]",
      "type": "Microsoft.Compute/virtualMachines/extensions",
      "name": "[concat(variables('masterVMNamePrefix'), copyIndex(variables('masterOffset')),'/cse', '-master-', copyIndex(variables('masterOffset')))]",
      "properties": {
        "publisher": "Microsoft.Azure.Extensions",
        "type": "CustomScript",
        "typeHandlerVersion": "2.0",
        "autoUpgradeMinorVersion": true,
        "settings": {},
        "protectedSettings": {
        {{if IsOpenShift}}
          "script": "{{ Base64 OpenShiftGetMasterSh }}"
        {{else}}
          "commandToExecute": "[concat('for i in $(seq 1 1200); do if [ -f /opt/azure/containers/provision.sh ]; then break; fi; if [ $i -eq 1200 ]; then exit 100; else sleep 1; fi; done; ', variables('provisionScriptParametersCommon'),' ',variables('provisionScriptParametersMaster'), ' /usr/bin/nohup /bin/bash -c \"stat /opt/azure/containers/provision.complete > /dev/null 2>&1 || /bin/bash /opt/azure/containers/provision.sh >> /var/log/azure/cluster-provision.log 2>&1\"')]"
        {{end}}
        }
      }
    }
    {{if UseAksExtension}}
    ,{
      "type": "Microsoft.Compute/virtualMachines/extensions",
      "name": "[concat(variables('masterVMNamePrefix'), copyIndex(variables('masterOffset')), '/computeAksLinuxBilling')]",
      "apiVersion": "[variables('apiVersionDefault')]",
      "copy": {
        "count": "[sub(variables('masterCount'), variables('masterOffset'))]",
        "name": "vmLoopNode"
      },
      "location": "[variables('location')]",
      "dependsOn": [
        "[concat('Microsoft.Compute/virtualMachines/', variables('masterVMNamePrefix'), copyIndex(variables('masterOffset')))]"
      ],
      "properties": {
        "publisher": "Microsoft.AKS",
        "type": "Compute.AKS-Engine.Linux.Billing",
        "typeHandlerVersion": "1.0",
        "autoUpgradeMinorVersion": true,
        "settings": {
        }
      }
    }
    {{end}}
    {{WriteLinkedTemplatesForExtensions}}
