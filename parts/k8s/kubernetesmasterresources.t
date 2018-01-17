{{if .MasterProfile.IsManagedDisks}} 
    {
      "apiVersion": "[variables('apiVersionStorageManagedDisks')]",
      "location": "[variables('location')]",
      "name": "[variables('masterAvailabilitySet')]",
      "properties":
        {
            "platformFaultDomainCount": "2",
            "platformUpdateDomainCount": "3",
		        "managed" : "true"
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
{{end}}
{{if not .MasterProfile.IsCustomVNET}}
    {
      "apiVersion": "[variables('apiVersionDefault')]",
      "dependsOn": [
        "[concat('Microsoft.Network/networkSecurityGroups/', variables('nsgName'))]"
{{if not IsAzureCNI}}
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
{{if not IsAzureCNI}}
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
{{if not IsAzureCNI}}
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
{{end}}
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
{{if .MasterProfile.IsCustomVNET}}
        "[variables('nsgID')]",
{{else}}
        "[variables('vnetID')]",
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
          {{range $seq := loop 1 .MasterProfile.IPAddressCount}}
          ,
          {
            "name": "[concat('ipconfig', add({{$seq}}, 1))]",
            "properties": {
              "privateIPAddress": "[variables('masterSecondaryAddrs')[add(mul(copyIndex(variables('masterOffset')), variables('ipAddressCount')), sub({{$seq}}, 1))]]",
              "primary": false,
              "privateIPAllocationMethod": "Static",
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
{{if .MasterProfile.IsCustomVNET}}
        ,"networkSecurityGroup": {
          "id": "[variables('nsgID')]"
        }
{{end}}
      },
      "type": "Microsoft.Network/networkInterfaces"
    },
    {
    {{if .MasterProfile.IsManagedDisks}}
      "apiVersion": "[variables('apiVersionStorageManagedDisks')]",
    {{else}}
      "apiVersion": "[variables('apiVersionDefault')]",
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
        "creationSource" : "[concat(variables('generatorCode'), '-', variables('masterVMNamePrefix'), copyIndex(variables('masterOffset')))]",
        "resourceNameSuffix" : "[variables('nameSuffix')]",
        "orchestrator" : "[variables('orchestratorNameVersionTag')]",
        "poolName" : "master"
      },
      "location": "[variables('location')]",
      "name": "[concat(variables('masterVMNamePrefix'), copyIndex(variables('masterOffset')))]",
      {{if UseManagedIdentity}}
      "identity": {
        "type": "systemAssigned"
      },
      {{end}}
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
              "id": "[resourceId('Microsoft.Network/networkInterfaces',concat(variables('masterVMNamePrefix'),'nic-', copyIndex(variables('masterOffset'))))]"
            }
          ]
        },
        "osProfile": {
          "adminUsername": "[variables('username')]",
          "computername": "[concat(variables('masterVMNamePrefix'), copyIndex(variables('masterOffset')))]",
          {{GetKubernetesMasterCustomData .}}
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
              "createOption": "Empty"
              ,"diskSizeGB": "[variables('etcdDiskSizeGB')]"
              ,"lun": 0
              ,"name": "[concat(variables('masterVMNamePrefix'), copyIndex(variables('masterOffset')),'-etcddisk')]"
          {{if .MasterProfile.IsStorageAccount}}
              ,"vhd": {
                "uri": "[concat(reference(concat('Microsoft.Storage/storageAccounts/',variables('masterStorageAccountName')),variables('apiVersionStorage')).primaryEndpoints.blob,'vhds/', variables('masterVMNamePrefix'),copyIndex(variables('masterOffset')),'-etcddisk.vhd')]"
              }
          {{end}}
            }
          ],
          "imageReference": {
            "offer": "[variables('osImageOffer')]",
            "publisher": "[variables('osImagePublisher')]",
            "sku": "[variables('osImageSku')]",
            "version": "[variables('osImageVersion')]"
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
     {
       "type": "Microsoft.Compute/virtualMachines/extensions",
       "name": "[concat(variables('masterVMNamePrefix'), copyIndex(), '/ManagedIdentityExtension')]",
       "copy": {
         "count": "[variables('masterCount')]",
         "name": "vmLoopNode"
       },
       "apiVersion": "2015-05-01-preview",
       "location": "[resourceGroup().location]",
       "dependsOn": [
         "[concat('Microsoft.Compute/virtualMachines/', variables('masterVMNamePrefix'), copyIndex())]",
         "[concat('Microsoft.Authorization/roleAssignments/', guid(concat('Microsoft.Compute/virtualMachines/', variables('masterVMNamePrefix'), copyIndex(), 'vmidentity')))]"
       ],
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
      "name": "[concat(variables('masterVMNamePrefix'), copyIndex(variables('masterOffset')),'/cse', copyIndex(variables('masterOffset')))]",
      "properties": {
        "publisher": "Microsoft.Azure.Extensions",
        "type": "CustomScript",
        "typeHandlerVersion": "2.0",
        "autoUpgradeMinorVersion": true,
        "settings": {},
        "protectedSettings": {
          "commandToExecute": "[concat(variables('provisionScriptParametersCommon'),' ',variables('provisionScriptParametersMaster'), ' MASTER_INDEX=',copyIndex(variables('masterOffset')),' /usr/bin/nohup /bin/bash -c \"stat /opt/azure/containers/provision.complete || /bin/bash /opt/azure/containers/provision.sh >> /var/log/azure/cluster-provision.log 2>&1\"')]"
        }
      }
    }{{WriteLinkedTemplatesForExtensions}}
