{{if UseManagedIdentity}}
  {
    "apiVersion": "2014-10-01-preview",
    "name": "[guid(concat('Microsoft.Compute/virtualMachineScaleSets/', variables('{{.Name}}VMNamePrefix'), 'vmidentity'))]",
    "type": "Microsoft.Authorization/roleAssignments",
    "properties": {
      "roleDefinitionId": "[variables('readerRoleDefinitionId')]",
      "principalId": "[reference(concat('Microsoft.Compute/virtualMachineScaleSets/', variables('{{.Name}}VMNamePrefix')), '2017-03-30', 'Full').identity.principalId]"
    }
  },
{{end}}
  {
    "apiVersion": "[variables('apiVersionVirtualMachineScaleSets')]",
    "dependsOn": [
    {{if .IsCustomVNET}}
      "[variables('nsgID')]"
    {{else}}
      "[variables('vnetID')]"
    {{end}}
    ],
    "tags":
    {
      "creationSource" : "[concat(variables('generatorCode'), '-', variables('{{.Name}}VMNamePrefix'))]",
      "resourceNameSuffix" : "[variables('nameSuffix')]",
      "orchestrator" : "[variables('orchestratorNameVersionTag')]",
      "poolName" : "{{.Name}}"
    },
    "location": "[variables('location')]",
    "name": "[variables('{{.Name}}VMNamePrefix')]",
    {{if UseManagedIdentity}}
    "identity": {
      "type": "systemAssigned"
    },
    {{end}}
    "sku": {
      "tier": "Standard",
      "capacity": "[variables('{{.Name}}Count')]",
      "name": "[variables('{{.Name}}VMSize')]"
    },
    "properties": {
      "overprovision": false,
      "upgradePolicy": {
        "mode": "Manual"
      },
      "virtualMachineProfile": {
        {{if .IsLowPriorityScaleSet}}
        "priority": "[variables('{{.Name}}ScaleSetPriority')]",
        "evictionPolicy": "[variables('{{.Name}}ScaleSetEvictionPolicy')]",
        {{end}}
        "networkProfile": {
          "networkInterfaceConfigurations": [
            {
              "name": "[variables('{{.Name}}VMNamePrefix')]",
              "properties": {
                "primary": true,
                "enableAcceleratedNetworking" : "{{.AcceleratedNetworkingEnabled}}",
                {{if .IsCustomVNET}}
                "networkSecurityGroup": {
                  "id": "[variables('nsgID')]"
                },
                {{end}}
                "ipConfigurations": [
                  {{range $seq := loop 1 .IPAddressCount}}
                  {
                    "name": "ipconfig{{$seq}}",
                    "properties": {
                      {{if eq $seq 1}}
                      "primary": true,
                      {{end}}
                      "subnet": {
                        "id": "[variables('{{$.Name}}VnetSubnetID')]"
                      }
                    }
                  }
                  {{if lt $seq $.IPAddressCount}},{{end}}
                  {{end}}
                ]
{{if HasCustomNodesDNS}}
                 ,"dnsSettings": {
                    "dnsServers": [
                        "[variables('dnsServer')]"
                    ]
                }
{{end}}
                {{if not IsAzureCNI}}
                ,"enableIPForwarding": true
                {{end}}
              }
            }
          ]
        },
        "osProfile": {
          "adminUsername": "[variables('username')]",
          "computerNamePrefix": "[variables('{{.Name}}VMNamePrefix')]",
          {{GetKubernetesAgentCustomData .}}
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
            {{if HasLinuxSecrets}}
              ,
              "secrets": "[variables('linuxProfileSecrets')]"
            {{end}}
        },
        "storageProfile": {
          {{if not (UseAgentCustomImage .)}}
            {{GetDataDisks .}}
          {{end}}
          "imageReference": {
            {{if UseAgentCustomImage .}}
            "id": "[resourceId(variables('{{.Name}}osImageResourceGroup'), 'Microsoft.Compute/images', variables('{{.Name}}osImageName'))]"
            {{else}}
            "offer": "[variables('{{.Name}}osImageOffer')]",
            "publisher": "[variables('{{.Name}}osImagePublisher')]",
            "sku": "[variables('{{.Name}}osImageSKU')]",
            "version": "[variables('{{.Name}}osImageVersion')]"
            {{end}}
          },
          "osDisk": {
            "createOption": "FromImage",
            "caching": "ReadWrite"
          {{if ne .OSDiskSizeGB 0}}
            ,"diskSizeGB": {{.OSDiskSizeGB}}
          {{end}}
          }
        },
        "extensionProfile": {
          "extensions": [
            {
              "name": "vmssCSE",
              "properties": {
                "publisher": "Microsoft.Azure.Extensions",
                "type": "CustomScript",
                "typeHandlerVersion": "2.0",
                "autoUpgradeMinorVersion": true,
                "settings": {},
                "protectedSettings": {
                  "commandToExecute": "[concat(variables('provisionScriptParametersCommon'),' /usr/bin/nohup /bin/bash -c \"/bin/bash /opt/azure/containers/provision.sh >> /var/log/azure/cluster-provision.log 2>&1\"')]"
                }
              }
            },
            {
              "name": "[concat(variables('{{.Name}}VMNamePrefix'), '-computeAksLinuxBilling')]",
              "location": "[variables('location')]",
              "properties": {
                "publisher": "Microsoft.AKS",
                "type": "Compute.AKS-Engine.Linux.Billing",
                "typeHandlerVersion": "1.0",
                "autoUpgradeMinorVersion": true,
                "settings": {}
              }
            }
            {{if UseManagedIdentity}}
            ,{
              "name": "managedIdentityExtension",
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
            }
            {{end}}
          ]
        }
      }
    },
    "type": "Microsoft.Compute/virtualMachineScaleSets"
  }
