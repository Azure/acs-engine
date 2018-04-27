    {
      "apiVersion": "[variables('apiVersionDefault')]",
      "copy": {
        "count": "[sub(variables('{{.Name}}Config').Count, variables('{{.Name}}Config').Offset)]",
        "name": "loop"
      },
      "dependsOn": [
{{if not IsOpenShift}}
{{if .IsCustomVNET}}
      "[variables('nsgID')]"
{{else}}
      "[variables('vnetID')]"
{{end}}
{{else}}
{{if not .IsCustomVNET}}
      "[variables('vnetID')]",
{{end}}
{{if eq .Role "infra"}}
      "[resourceId('Microsoft.Network/networkSecurityGroups', 'router-nsg')]"
{{else}}
      "[variables('nsgID')]"
{{end}}
{{end}}
      ],
      "location": "[variables('location')]",
      "name": "[concat(variables('{{.Name}}Config').VMNamePrefix, 'nic-', copyIndex(variables('{{.Name}}Config').Offset))]",
      "properties": {
{{if not IsOpenShift}}
{{if .IsCustomVNET}}
        "networkSecurityGroup": {
          "id": "[variables('nsgID')]"
        },
{{end}}
{{else}}
        "networkSecurityGroup": {
          {{if eq .Role "infra"}}
          "id": "[resourceId('Microsoft.Network/networkSecurityGroups', 'router-nsg')]"
          {{else}}
          "id": "[variables('nsgID')]"
          {{end}}
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
              "privateIPAllocationMethod": "Dynamic",
              "subnet": {
                "id": "[variables('{{$.Name}}Config').VnetSubnetID]"
              }
{{if eq $.Role "infra"}}
              ,
              "loadBalancerBackendAddressPools": [
                {
                    "id": "[concat(resourceId('Microsoft.Network/loadBalancers', 'router-lb'), '/backendAddressPools/backend')]"
                }
              ]
{{end}}
            }
          }
          {{if lt $seq $.IPAddressCount}},{{end}}
          {{end}}
        ]
{{if not IsAzureCNI}}
        ,
        "enableIPForwarding": true
{{end}}
      },
      "type": "Microsoft.Network/networkInterfaces"
    },
{{if .IsManagedDisks}}
   {
      "location": "[variables('location')]",
      "name": "[variables('{{.Name}}Config').AvailabilitySet]",
      "apiVersion": "[variables('apiVersionStorageManagedDisks')]",
      "properties":
        {
            "platformFaultDomainCount": 2,
            "platformUpdateDomainCount": 3,
		"managed" : "true"
        },

      "type": "Microsoft.Compute/availabilitySets"
    },
{{else if .IsStorageAccount}}
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
    {
      "location": "[variables('location')]",
      "name": "[variables('{{.Name}}Config').AvailabilitySet]",
      "apiVersion": "[variables('apiVersionDefault')]",
      "properties": {},
      "type": "Microsoft.Compute/availabilitySets"
    },
{{end}}
  {
    {{if .IsManagedDisks}}
      "apiVersion": "[variables('apiVersionStorageManagedDisks')]",
    {{else}}
      "apiVersion": "[variables('apiVersionDefault')]",
    {{end}}
      "copy": {
        "count": "[sub(variables('{{.Name}}Config').Count, variables('{{.Name}}Config').Offset)]",
        "name": "vmLoopNode"
      },
      "dependsOn": [
{{if .IsStorageAccount}}
        "[concat('Microsoft.Storage/storageAccounts/',variables('storageAccountPrefixes')[mod(add(div(copyIndex(variables('{{.Name}}Config').Offset),variables('maxVMsPerStorageAccount')),variables('{{.Name}}Config').StorageAccountOffset),variables('storageAccountPrefixesCount'))],variables('storageAccountPrefixes')[div(add(div(copyIndex(variables('{{.Name}}Config').Offset),variables('maxVMsPerStorageAccount')),variables('{{.Name}}Config').StorageAccountOffset),variables('storageAccountPrefixesCount'))],variables('{{.Name}}Config').AccountName)]",

  {{if .HasDisks}}
        "[concat('Microsoft.Storage/storageAccounts/',variables('storageAccountPrefixes')[mod(add(add(div(copyIndex(variables('{{.Name}}Config').Offset),variables('maxVMsPerStorageAccount')),variables('{{.Name}}Config').StorageAccountOffset),variables('dataStorageAccountPrefixSeed')),variables('storageAccountPrefixesCount'))],variables('storageAccountPrefixes')[div(add(add(div(copyIndex(variables('{{.Name}}Config').Offset),variables('maxVMsPerStorageAccount')),variables('{{.Name}}Config').StorageAccountOffset),variables('dataStorageAccountPrefixSeed')),variables('storageAccountPrefixesCount'))],variables('{{.Name}}Config').DataAccountName)]",
  {{end}}
{{end}}
        "[concat('Microsoft.Network/networkInterfaces/', variables('{{.Name}}Config').VMNamePrefix, 'nic-', copyIndex(variables('{{.Name}}Config').Offset))]",
        "[concat('Microsoft.Compute/availabilitySets/', variables('{{.Name}}Config').AvailabilitySet)]"
      ],
      "tags":
      {
        "creationSource" : "[concat(variables('generatorCode'), '-', variables('{{.Name}}Config').VMNamePrefix, copyIndex(variables('{{.Name}}Config').Offset))]",
        "resourceNameSuffix" : "[variables('nameSuffix')]",
        "orchestrator" : "[variables('orchestratorNameVersionTag')]",
        "poolName" : "{{.Name}}"
      },
      "location": "[variables('location')]",
      "name": "[concat(variables('{{.Name}}Config').VMNamePrefix, copyIndex(variables('{{.Name}}Config').Offset))]",
      {{if UseManagedIdentity}}
      "identity": {
        "type": "systemAssigned"
      },
      {{end}}
      {{if and IsOpenShift (not (UseAgentCustomImage .))}}
      "plan": {
        "name": "[variables('{{.Name}}Config').osImageSKU]",
        "publisher": "[variables('{{.Name}}Config').osImagePublisher]",
        "product": "[variables('{{.Name}}Config').osImageOffer]"
      },
      {{end}}
      "properties": {
        "availabilitySet": {
          "id": "[resourceId('Microsoft.Compute/availabilitySets',variables('{{.Name}}Config').AvailabilitySet)]"
        },
        "hardwareProfile": {
          "vmSize": "[variables('{{.Name}}Config').VMSize]"
        },
        "networkProfile": {
          "networkInterfaces": [
            {
              "id": "[resourceId('Microsoft.Network/networkInterfaces',concat(variables('{{.Name}}Config').VMNamePrefix, 'nic-', copyIndex(variables('{{.Name}}Config').Offset)))]"
            }
          ]
        },
        "osProfile": {
          "adminUsername": "[variables('username')]",
          "computername": "[concat(variables('{{.Name}}Config').VMNamePrefix, copyIndex(variables('{{.Name}}Config').Offset))]",
          {{if not IsOpenShift}}
          {{GetKubernetesAgentCustomData .}}
          {{end}}
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
            "createOption": "FromImage"
            ,"caching": "ReadWrite"
          {{if .IsStorageAccount}}
            ,"name": "[concat(variables('{{.Name}}Config').VMNamePrefix, copyIndex(variables('{{.Name}}Config').Offset),'-osdisk')]"
            ,"vhd": {
              "uri": "[concat(reference(concat('Microsoft.Storage/storageAccounts/',variables('storageAccountPrefixes')[mod(add(div(copyIndex(variables('{{.Name}}Config').Offset),variables('maxVMsPerStorageAccount')),variables('{{.Name}}Config').StorageAccountOffset),variables('storageAccountPrefixesCount'))],variables('storageAccountPrefixes')[div(add(div(copyIndex(variables('{{.Name}}Config').Offset),variables('maxVMsPerStorageAccount')),variables('{{.Name}}Config').StorageAccountOffset),variables('storageAccountPrefixesCount'))],variables('{{.Name}}Config').AccountName),variables('apiVersionStorage')).primaryEndpoints.blob,'osdisk/', variables('{{.Name}}Config').VMNamePrefix, copyIndex(variables('{{.Name}}Config').Offset), '-osdisk.vhd')]"
            }
          {{end}}
          {{if ne .OSDiskSizeGB 0}}
            ,"diskSizeGB": {{.OSDiskSizeGB}}
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
         "count": "[variables('{{.Name}}Config').Count]",
         "name": "vmLoopNode"
       },
      "name": "[guid(concat('Microsoft.Compute/virtualMachines/', variables('{{.Name}}Config').VMNamePrefix, copyIndex(), 'vmidentity'))]",
      "type": "Microsoft.Authorization/roleAssignments",
      "properties": {
        "roleDefinitionId": "[variables('readerRoleDefinitionId')]",
        "principalId": "[reference(concat('Microsoft.Compute/virtualMachines/', variables('{{.Name}}Config').VMNamePrefix, copyIndex()), '2017-03-30', 'Full').identity.principalId]"
      }
    },
     {
       "type": "Microsoft.Compute/virtualMachines/extensions",
       "name": "[concat(variables('{{.Name}}Config').VMNamePrefix, copyIndex(), '/ManagedIdentityExtension')]",
       "copy": {
         "count": "[variables('{{.Name}}Config').Count]",
         "name": "vmLoopNode"
       },
       "apiVersion": "2015-05-01-preview",
       "location": "[resourceGroup().location]",
       "dependsOn": [
         "[concat('Microsoft.Compute/virtualMachines/', variables('{{.Name}}Config').VMNamePrefix, copyIndex())]",
         "[concat('Microsoft.Authorization/roleAssignments/', guid(concat('Microsoft.Compute/virtualMachines/', variables('{{.Name}}Config').VMNamePrefix, copyIndex(), 'vmidentity')))]"
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
        "count": "[sub(variables('{{.Name}}Config').Count, variables('{{.Name}}Config').Offset)]",
        "name": "vmLoopNode"
      },
      "dependsOn": [
        {{if UseManagedIdentity}}
        "[concat('Microsoft.Compute/virtualMachines/', variables('{{.Name}}Config').VMNamePrefix, copyIndex(), '/extensions/ManagedIdentityExtension')]"
        {{else}}
        "[concat('Microsoft.Compute/virtualMachines/', variables('{{.Name}}Config').VMNamePrefix, copyIndex(variables('{{.Name}}Config').Offset))]"
        {{end}}
      ],
      "location": "[variables('location')]",
      "type": "Microsoft.Compute/virtualMachines/extensions",
      "name": "[concat(variables('{{.Name}}Config').VMNamePrefix, copyIndex(variables('{{.Name}}Config').Offset),'/cse', '-agent-', copyIndex(variables('{{.Name}}Config').Offset))]",
      "properties": {
        "publisher": "Microsoft.Azure.Extensions",
        "type": "CustomScript",
        "typeHandlerVersion": "2.0",
        "autoUpgradeMinorVersion": true,
        "settings": {},
        "protectedSettings": {
        {{if IsOpenShift }}
          "script": "{{ Base64 (OpenShiftGetNodeSh .) }}"
        {{else}}
          "commandToExecute": "[concat(variables('provisionScriptParametersCommon'),' /usr/bin/nohup /bin/bash -c \"/bin/bash /opt/azure/containers/provision.sh >> /var/log/azure/cluster-provision.log 2>&1\"')]"
        {{end}}
        }
      }
    }
