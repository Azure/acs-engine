  {
      "apiVersion": "[variables('apiVersionVirtualMachineScaleSets')]",
      "dependsOn": [
        "[variables('vnetID')]"
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
        "overprovision": true,
        "upgradePolicy": {
          "mode": "Automatic"
        },
        "virtualMachineProfile": {
          "networkProfile": {
            "networkInterfaceConfigurations": [
              {
                "name": "[variables('{{.Name}}VMNamePrefix')]",
                "properties": {
                  "primary": true,
                  "enableIPForwarding": true,
                  "ipConfigurations": [
                    {
                      "name": "[variables('{{.Name}}VMNamePrefix')]",
                      "properties": {
                        "subnet": {
                          "id": "[variables('{{$.Name}}VnetSubnetID')]"
                        }
                      }
                    }
                  ]
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
            {{GetDataDisks .}}
            "imageReference": {
              "offer": "[variables('{{.Name}}osImageOffer')]",
              "publisher": "[variables('{{.Name}}osImagePublisher')]",
              "sku": "[variables('{{.Name}}osImageSKU')]",
              "version": "[variables('{{.Name}}osImageVersion')]"
            },
            "osDisk": {
              "createOption": "FromImage",
              "caching": "ReadWrite"
            {{if ne .OSDiskSizeGB 0}}
              ,"diskSizeGB": {{.OSDiskSizeGB}}
            {{end}}
            }
          }
        }
      },
      "type": "Microsoft.Compute/virtualMachineScaleSets"
    },
    {
    "apiVersion": "[variables('apiVersionVirtualMachineScaleSets')]",
    "dependsOn": [
      "[concat('Microsoft.Compute/virtualMachineScaleSets/', variables('{{.Name}}VMNamePrefix'))]"
    ],
    "location": "[variables('location')]",
    "type": "Microsoft.Compute/virtualMachineScaleSets/extensions",
    "name": "[concat(variables('{{.Name}}VMNamePrefix'),'/cse')]",
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
  }