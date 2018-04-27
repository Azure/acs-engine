{{if .IsStorageAccount}}
  {
    "apiVersion": "[variables('apiVersionStorage')]",
    "copy": {
      "count": "[variables('{{.Name}}Config').StorageAccountsCount]",
      "name": "loop"
    },
    {{if not IsHostedMaster}}
      {{if not IsPrivateCluster}}
        "dependsOn": [
          "[concat('Microsoft.Network/publicIPAddresses/', variables('masterPublicIPAddressName'))]"
        ],
      {{end}}
    {{end}}
    "location": "[variables('location')]",
    "name": "[concat(variables('storageAccountPrefixes')[mod(add(copyIndex(),variables('{{.Name}}Config').StorageAccountOffset),variables('storageAccountPrefixesCount'))],variables('storageAccountPrefixes')[div(add(copyIndex(),variables('{{.Name}}Config').StorageAccountOffset),variables('storageAccountPrefixesCount'))],variables('{{.Name}}Config').AccountName)]",
    "properties": {
      "accountType": "[variables('vmSizesMap')[variables('{{.Name}}Config').VMSize].storageAccountType]"
    },
    "type": "Microsoft.Storage/storageAccounts"
  },
  {{if .HasDisks}}
  {
    "apiVersion": "[variables('apiVersionStorage')]",
    "copy": {
      "count": "[variables('{{.Name}}Config').StorageAccountsCount]",
      "name": "datadiskLoop"
    },
    {{if not IsHostedMaster}}
      {{if not IsPrivateCluster}}
        "dependsOn": [
          "[concat('Microsoft.Network/publicIPAddresses/', variables('masterPublicIPAddressName'))]"
        ],
      {{end}}
    {{end}}
    "location": "[variables('location')]",
    "name": "[concat(variables('storageAccountPrefixes')[mod(add(copyIndex(variables('dataStorageAccountPrefixSeed')),variables('{{.Name}}Config').StorageAccountOffset),variables('storageAccountPrefixesCount'))],variables('storageAccountPrefixes')[div(add(copyIndex(variables('dataStorageAccountPrefixSeed')),variables('{{.Name}}Config').StorageAccountOffset),variables('storageAccountPrefixesCount'))],variables('{{.Name}}Config').DataAccountName)]",
    "properties": {
      "accountType": "[variables('vmSizesMap')[variables('{{.Name}}Config').VMSize].storageAccountType]"
    },
    "type": "Microsoft.Storage/storageAccounts"
  },
  {{end}}
{{end}}
{{if UseManagedIdentity}}
  {
    "apiVersion": "2014-10-01-preview",
    "name": "[guid(concat('Microsoft.Compute/virtualMachineScaleSets/', variables('{{.Name}}Config').VMNamePrefix, 'vmidentity'))]",
    "type": "Microsoft.Authorization/roleAssignments",
    "properties": {
      "roleDefinitionId": "[variables('readerRoleDefinitionId')]",
      "principalId": "[reference(concat('Microsoft.Compute/virtualMachineScaleSets/', variables('{{.Name}}Config').VMNamePrefix), '2017-03-30', 'Full').identity.principalId]"
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
    {{if .IsStorageAccount}}
        ,"[concat('Microsoft.Storage/storageAccounts/',variables('storageAccountPrefixes')[mod(add(div(0,variables('maxVMsPerStorageAccount')),variables('{{.Name}}Config').StorageAccountOffset),variables('storageAccountPrefixesCount'))],variables('storageAccountPrefixes')[div(add(div(0,variables('maxVMsPerStorageAccount')),variables('{{.Name}}Config').StorageAccountOffset),variables('storageAccountPrefixesCount'))],variables('{{.Name}}Config').AccountName)]"
    {{if .HasDisks}}
        ,"[concat('Microsoft.Storage/storageAccounts/',variables('storageAccountPrefixes')[mod(add(add(div(0,variables('maxVMsPerStorageAccount')),variables('{{.Name}}Config').StorageAccountOffset),variables('dataStorageAccountPrefixSeed')),variables('storageAccountPrefixesCount'))],variables('storageAccountPrefixes')[div(add(add(div(0,variables('maxVMsPerStorageAccount')),variables('{{.Name}}Config').StorageAccountOffset),variables('dataStorageAccountPrefixSeed')),variables('storageAccountPrefixesCount'))],variables('{{.Name}}Config').DataAccountName)]"
    {{end}}
    {{end}}
    ],
    "tags":
    {
      "creationSource" : "[concat(variables('generatorCode'), '-', variables('{{.Name}}Config').VMNamePrefix)]",
      "resourceNameSuffix" : "[variables('nameSuffix')]",
      "orchestrator" : "[variables('orchestratorNameVersionTag')]",
      "poolName" : "{{.Name}}"
    },
    "location": "[variables('location')]",
    "name": "[variables('{{.Name}}Config').VMNamePrefix]",
    {{if UseManagedIdentity}}
    "identity": {
      "type": "systemAssigned"
    },
    {{end}}
    "sku": {
      "tier": "Standard",
      "capacity": "[variables('{{.Name}}Config').Count]",
      "name": "[variables('{{.Name}}Config').VMSize]"
    },
    "properties": {
      "overprovision": false,
      "upgradePolicy": {
        "mode": "Manual"
      },
      "virtualMachineProfile": {
        "networkProfile": {
          "networkInterfaceConfigurations": [
            {
              "name": "[variables('{{.Name}}Config').VMNamePrefix]",
              "properties": {
                "primary": true,
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
                {{if not IsAzureCNI}}
                ,"enableIPForwarding": true
                {{end}}
              }
            }
          ]
        },
        "osProfile": {
          "adminUsername": "[variables('username')]",
          "computerNamePrefix": "[variables('{{.Name}}Config').VMNamePrefix]",
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
            "id": "[resourceId(variables('{{.Name}}Config').osImageResourceGroup, 'Microsoft.Compute/images', variables('{{.Name}}Config').osImageName)]"
            {{else}}
            "offer": "[variables('{{.Name}}Config').osImageOffer]",
            "publisher": "[variables('{{.Name}}Config').osImagePublisher]",
            "sku": "[variables('{{.Name}}Config').osImageSKU]",
            "version": "[variables('{{.Name}}Config').osImageVersion]"
            {{end}}
          },
          "osDisk": {
            "createOption": "FromImage",
            "caching": "ReadWrite"
          {{if .IsStorageAccount}}
            ,"name": "[concat(variables('{{.Name}}Config').VMNamePrefix,'-osdisk')]"
            ,"vhdContainers": [
              "[concat(reference(concat('Microsoft.Storage/storageAccounts/',variables('storageAccountPrefixes')[mod(add(0,variables('{{.Name}}Config').StorageAccountOffset),variables('storageAccountPrefixesCount'))],variables('storageAccountPrefixes')[div(add(0,variables('{{.Name}}Config').StorageAccountOffset),variables('storageAccountPrefixesCount'))],variables('{{.Name}}Config').AccountName),variables('apiVersionStorage')).primaryEndpoints.blob,'osdisk')]"
            ]
          {{end}}
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