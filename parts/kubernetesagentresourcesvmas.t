    {
      "apiVersion": "[variables('apiVersionDefault')]",
      "copy": {
        "count": "[sub(variables('{{.Name}}Count'), variables('{{.Name}}Offset'))]",
        "name": "loop"
      },
      "dependsOn": [
{{if .IsCustomVNET}}
      "[variables('nsgID')]"
{{else}}
      "[variables('vnetID')]"
{{end}}
      ],
      "location": "[variables('location')]",
      "name": "[concat(variables('{{.Name}}VMNamePrefix'), 'nic-', copyIndex(variables('{{.Name}}Offset')))]",
      "properties": {
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
              "privateIPAllocationMethod": "Dynamic",
              "subnet": {
                "id": "[variables('{{$.Name}}VnetSubnetID')]"
              }
            }
          }
          {{if lt $seq $.IPAddressCount}},{{end}}
          {{end}}
        ]
{{if not IsVNETIntegrated}}
        ,
        "enableIPForwarding": true
{{end}}
      },
      "type": "Microsoft.Network/networkInterfaces"
    },
{{if .IsManagedDisks}} 
   {
      "location": "[variables('location')]",
      "name": "[variables('{{.Name}}AvailabilitySet')]",
      "apiVersion": "[variables('apiVersionStorageManagedDisks')]",
      "properties":
        {
            "platformFaultDomainCount": "2",
            "platformUpdateDomainCount": "3",
		"managed" : "true"
        },
  
      "type": "Microsoft.Compute/availabilitySets"
    },
{{else if .IsStorageAccount}}
    {
      "apiVersion": "[variables('apiVersionStorage')]",
      "copy": {
        "count": "[variables('{{.Name}}StorageAccountsCount')]",
        "name": "loop"
      },
      {{if not IsHostedMaster}}
      "dependsOn": [
        "[concat('Microsoft.Network/publicIPAddresses/', variables('masterPublicIPAddressName'))]"
      ],
      {{end}}
      "location": "[variables('location')]",
      "name": "[concat(variables('storageAccountPrefixes')[mod(add(copyIndex(),variables('{{.Name}}StorageAccountOffset')),variables('storageAccountPrefixesCount'))],variables('storageAccountPrefixes')[div(add(copyIndex(),variables('{{.Name}}StorageAccountOffset')),variables('storageAccountPrefixesCount'))],variables('{{.Name}}AccountName'))]",
      "properties": {
        "accountType": "[variables('vmSizesMap')[variables('{{.Name}}VMSize')].storageAccountType]"
      },
      "type": "Microsoft.Storage/storageAccounts"
    },
    {{if .HasDisks}}
    {
      "apiVersion": "[variables('apiVersionStorage')]",
      "copy": {
        "count": "[variables('{{.Name}}StorageAccountsCount')]",
        "name": "datadiskLoop"
      },
      {{if not IsHostedMaster}}
      "dependsOn": [
        "[concat('Microsoft.Network/publicIPAddresses/', variables('masterPublicIPAddressName'))]"
      ],
      {{end}}
      "location": "[variables('location')]",
      "name": "[concat(variables('storageAccountPrefixes')[mod(add(copyIndex(variables('dataStorageAccountPrefixSeed')),variables('{{.Name}}StorageAccountOffset')),variables('storageAccountPrefixesCount'))],variables('storageAccountPrefixes')[div(add(copyIndex(variables('dataStorageAccountPrefixSeed')),variables('{{.Name}}StorageAccountOffset')),variables('storageAccountPrefixesCount'))],variables('{{.Name}}DataAccountName'))]",
      "properties": {
        "accountType": "[variables('vmSizesMap')[variables('{{.Name}}VMSize')].storageAccountType]"
      },
      "type": "Microsoft.Storage/storageAccounts"
    },
    {{end}}
    {
      "location": "[variables('location')]",
      "name": "[variables('{{.Name}}AvailabilitySet')]",
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
        "count": "[sub(variables('{{.Name}}Count'), variables('{{.Name}}Offset'))]",
        "name": "vmLoopNode"
      },
      "dependsOn": [
{{if .IsStorageAccount}}
        "[concat('Microsoft.Storage/storageAccounts/',variables('storageAccountPrefixes')[mod(add(div(copyIndex(variables('{{.Name}}Offset')),variables('maxVMsPerStorageAccount')),variables('{{.Name}}StorageAccountOffset')),variables('storageAccountPrefixesCount'))],variables('storageAccountPrefixes')[div(add(div(copyIndex(variables('{{.Name}}Offset')),variables('maxVMsPerStorageAccount')),variables('{{.Name}}StorageAccountOffset')),variables('storageAccountPrefixesCount'))],variables('{{.Name}}AccountName'))]",

  {{if .HasDisks}}
        "[concat('Microsoft.Storage/storageAccounts/',variables('storageAccountPrefixes')[mod(add(add(div(copyIndex(variables('{{.Name}}Offset')),variables('maxVMsPerStorageAccount')),variables('{{.Name}}StorageAccountOffset')),variables('dataStorageAccountPrefixSeed')),variables('storageAccountPrefixesCount'))],variables('storageAccountPrefixes')[div(add(add(div(copyIndex(variables('{{.Name}}Offset')),variables('maxVMsPerStorageAccount')),variables('{{.Name}}StorageAccountOffset')),variables('dataStorageAccountPrefixSeed')),variables('storageAccountPrefixesCount'))],variables('{{.Name}}DataAccountName'))]",
  {{end}}
{{end}}
        "[concat('Microsoft.Network/networkInterfaces/', variables('{{.Name}}VMNamePrefix'), 'nic-', copyIndex(variables('{{.Name}}Offset')))]",
        "[concat('Microsoft.Compute/availabilitySets/', variables('{{.Name}}AvailabilitySet'))]"
      ],
      "tags":
      {
        "creationSource" : "[concat('acsengine-', variables('{{.Name}}VMNamePrefix'), copyIndex(variables('{{.Name}}Offset')))]",
        "resourceNameSuffix" : "[variables('nameSuffix')]",
        "orchestrator" : "[variables('orchestratorNameVersionTag')]",
        "poolName" : "{{.Name}}"
      },
      "location": "[variables('location')]",
      "name": "[concat(variables('{{.Name}}VMNamePrefix'), copyIndex(variables('{{.Name}}Offset')))]",
      {{if UseManagedIdentity}}
      "identity": {
        "type": "systemAssigned"
      },
      {{end}}
      "properties": {
        "availabilitySet": {
          "id": "[resourceId('Microsoft.Compute/availabilitySets',variables('{{.Name}}AvailabilitySet'))]"
        },
        "hardwareProfile": {
          "vmSize": "[variables('{{.Name}}VMSize')]"
        },
        "networkProfile": {
          "networkInterfaces": [
            {
              "id": "[resourceId('Microsoft.Network/networkInterfaces',concat(variables('{{.Name}}VMNamePrefix'), 'nic-', copyIndex(variables('{{.Name}}Offset'))))]"
            }
          ]
        },
        "osProfile": {
          "adminUsername": "[variables('username')]",
          "computername": "[concat(variables('{{.Name}}VMNamePrefix'), copyIndex(variables('{{.Name}}Offset')))]",
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
            "offer": "[variables('osImageOffer')]",
            "publisher": "[variables('osImagePublisher')]",
            "sku": "[variables('osImageSKU')]",
            "version": "[variables('osImageVersion')]"
          },
          "osDisk": {
            "createOption": "FromImage"
            ,"caching": "ReadWrite"
          {{if .IsStorageAccount}}
            ,"name": "[concat(variables('{{.Name}}VMNamePrefix'), copyIndex(variables('{{.Name}}Offset')),'-osdisk')]"
            ,"vhd": {
              "uri": "[concat(reference(concat('Microsoft.Storage/storageAccounts/',variables('storageAccountPrefixes')[mod(add(div(copyIndex(variables('{{.Name}}Offset')),variables('maxVMsPerStorageAccount')),variables('{{.Name}}StorageAccountOffset')),variables('storageAccountPrefixesCount'))],variables('storageAccountPrefixes')[div(add(div(copyIndex(variables('{{.Name}}Offset')),variables('maxVMsPerStorageAccount')),variables('{{.Name}}StorageAccountOffset')),variables('storageAccountPrefixesCount'))],variables('{{.Name}}AccountName')),variables('apiVersionStorage')).primaryEndpoints.blob,'osdisk/', variables('{{.Name}}VMNamePrefix'), copyIndex(variables('{{.Name}}Offset')), '-osdisk.vhd')]"
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
         "count": "[variables('{{.Name}}Count')]",
         "name": "vmLoopNode"
       },
      "name": "[guid(concat('Microsoft.Compute/virtualMachines/', variables('{{.Name}}VMNamePrefix'), copyIndex(), 'vmidentity'))]",
      "type": "Microsoft.Authorization/roleAssignments",
      "properties": {
        "roleDefinitionId": "[variables('readerRoleDefinitionId')]",
        "principalId": "[reference(concat('Microsoft.Compute/virtualMachines/', variables('{{.Name}}VMNamePrefix'), copyIndex()), '2017-03-30', 'Full').identity.principalId]"
      }
    },
     {
       "type": "Microsoft.Compute/virtualMachines/extensions",
       "name": "[concat(variables('{{.Name}}VMNamePrefix'), copyIndex(), '/ManagedIdentityExtension')]",
       "copy": {
         "count": "[variables('{{.Name}}Count')]",
         "name": "vmLoopNode"
       },
       "apiVersion": "2015-05-01-preview",
       "location": "[resourceGroup().location]",
       "dependsOn": [
         "[concat('Microsoft.Compute/virtualMachines/', variables('{{.Name}}VMNamePrefix'), copyIndex())]",
         "[concat('Microsoft.Authorization/roleAssignments/', guid(concat('Microsoft.Compute/virtualMachines/', variables('{{.Name}}VMNamePrefix'), copyIndex(), 'vmidentity')))]"
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
        "count": "[sub(variables('{{.Name}}Count'), variables('{{.Name}}Offset'))]",
        "name": "vmLoopNode"
      },
      "dependsOn": [
        {{if UseManagedIdentity}}
        "[concat('Microsoft.Compute/virtualMachines/', variables('{{.Name}}VMNamePrefix'), copyIndex(), '/extensions/ManagedIdentityExtension')]"
        {{else}}
        "[concat('Microsoft.Compute/virtualMachines/', variables('{{.Name}}VMNamePrefix'), copyIndex(variables('{{.Name}}Offset')))]"
        {{end}}
      ],
      "location": "[variables('location')]",
      "type": "Microsoft.Compute/virtualMachines/extensions",
      "name": "[concat(variables('{{.Name}}VMNamePrefix'), copyIndex(variables('{{.Name}}Offset')),'/cse', copyIndex(variables('{{.Name}}Offset')))]",
      "properties": {
        "publisher": "Microsoft.Azure.Extensions",
        "type": "CustomScript",
        "typeHandlerVersion": "2.0",
        "autoUpgradeMinorVersion": true,
        "settings": {},
        "protectedSettings": {
          "commandToExecute": "[concat('/usr/bin/nohup /bin/bash -c \"/bin/bash /opt/azure/containers/provision.sh ',variables('tenantID'),' ',variables('subscriptionId'),' ',variables('resourceGroup'),' ',variables('location'),' ',variables('subnetName'),' ',variables('nsgName'),' ',variables('virtualNetworkName'),' ',variables('routeTableName'),' ',variables('primaryAvailablitySetName'),' ',variables('servicePrincipalClientId'),' ',variables('servicePrincipalClientSecret'),' ',variables('clientPrivateKey'),' ',variables('targetEnvironment'),' ',variables('networkPolicy'),' ',variables('cloudProviderBackoff'),' ',variables('cloudProviderBackoffRetries'),' ',variables('cloudProviderBackoffExponent'),' ',variables('cloudProviderBackoffDuration'),' ',variables('cloudProviderBackoffJitter'),' ',variables('cloudProviderRatelimit'),' ',variables('cloudProviderRatelimitQPS'),' ',variables('cloudProviderRatelimitBucket'),' ', variables('useManagedIdentityExtension'),' ',variables('useInstanceMetadata'),' >> /var/log/azure/cluster-provision.log 2>&1 &\" &')]"
        }
      }
    }
