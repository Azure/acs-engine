    {
      "apiVersion": "[variables('apiVersionNetwork')]",
      "copy": {
        "count": "[sub(variables('{{.Name}}Count'), variables('{{.Name}}Variables').Offset)]",
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
{{if .IsCustomVNET}}
      "[concat(variables('masterVMNamePrefix'), 'nic-0')]",
{{else}}
      "[variables('vnetID')]",
{{end}}
{{if eq .Role "infra"}}
      "[variables('routerLBName')]",
      "[variables('routerNSGID')]"
{{else}}
      "[variables('nsgID')]"
{{end}}
{{end}}
      ],
      "location": "[variables('location')]",
      "name": "[concat(variables('{{.Name}}Variables').VMNamePrefix, 'nic-', copyIndex(variables('{{.Name}}Variables').Offset))]",
      "properties": {
        "enableAcceleratedNetworking" : {{.AcceleratedNetworkingEnabled}},
{{if not IsOpenShift}}
{{if .IsCustomVNET}}
        "networkSecurityGroup": {
          "id": "[variables('nsgID')]"
        },
{{end}}
{{else}}
        "networkSecurityGroup": {
          {{if eq .Role "infra"}}
          "id": "[variables('routerNSGID')]"
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
                {{if eq $.Name "system"}}
                "privateIPAddress": "[concat(variables('aciSystemNodeAddrPrefix'), copyIndex(add(50, int(variables('aciPrimaryIPOctet4')))))]",
                {{else if eq $.Name "agentpool1"}}
                "privateIPAddress": "[concat(variables('aciSystemNodeAddrPrefix'), copyIndex(add(100, int(variables('aciPrimaryIPOctet4')))))]",
                {{else}}
                "privateIPAddress": "[concat(variables('aciCustomerNodeAddrPrefix'), copyIndex(mul(25, sub(int(variables('{{$.Name}}Index')), 2))), '.', variables('aciPrimaryIPOctet4'))]",
                {{end}}
              {{else if eq $.Name "system"}}
              "privateIPAddress": "[concat(variables('aciSystemPodAddrPrefix'), copyIndex(add(50, int(variables('aciPrimaryIPOctet4')))), '.', add(sub({{$seq}}, 1), int(variables('aciPrimaryIPOctet4'))))]",
              {{else if eq $.Name "agentpool1"}}
              "privateIPAddress": "[concat(variables('aciSystemPodAddrPrefix'), copyIndex(add(100, int(variables('aciPrimaryIPOctet4')))), '.', add(sub({{$seq}}, 1), int(variables('aciPrimaryIPOctet4'))))]",
              {{else}}
              "privateIPAddress": "[concat(variables('aciCustomerPodAddrPrefix'), copyIndex(mul(25, sub(int(variables('{{$.Name}}Index')), 2))), '.', add(sub({{$seq}}, 1), int(variables('aciPrimaryIPOctet4'))))]",
              {{end}}
              "privateIPAllocationMethod": "Static",
              "subnet": {
                "id": "[variables('{{$.Name}}Variables').VnetSubnetID]"
              }
{{if eq $.Role "infra"}}
              ,
              "loadBalancerBackendAddressPools": [
                {
                    "id": "[concat(resourceId('Microsoft.Network/loadBalancers', variables('routerLBName')), '/backendAddressPools/backend')]"
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
      "name": "[variables('{{.Name}}Variables').AvailabilitySet]",
      "apiVersion": "[variables('apiVersionStorageManagedDisks')]",
      "properties":
        {
            "platformFaultDomainCount": 2,
            "platformUpdateDomainCount": 3
            "managed" : "true"
        },
      "sku": {
        "name": "Aligned"
      },
      "type": "Microsoft.Compute/availabilitySets"
    },
{{else if .IsStorageAccount}}
    {
      "apiVersion": "[variables('apiVersionStorage')]",
      "copy": {
        "count": "[variables('{{.Name}}Variables').StorageAccountsCount]",
        "name": "loop"
      },
      {{if not IsHostedMaster}}
        {{if not IsPrivateCluster}}
          "dependsOn": [
            "[concat('Microsoft.Network/publicIPAddresses/', variables('masterPublicIPAddressName'))]"
          ],
        {{end}}
      {{end}}
      "kind": "Storage",
      "location": "[variables('location')]",
      "name": "[concat(variables('storageAccountPrefixes')[mod(add(copyIndex(),variables('{{.Name}}Variables').StorageAccountOffset),variables('storageAccountPrefixesCount'))],variables('storageAccountPrefixes')[div(add(copyIndex(),variables('{{.Name}}Variables').StorageAccountOffset),variables('storageAccountPrefixesCount'))],variables('{{.Name}}AccountName'))]",
      "properties": {
        "encryption": {
          "keySource": "Microsoft.Storage",
          "services": {
            "blob": {
              "enabled": true
            },
            "file": {
              "enabled": true
            }
          }
        },
        "supportsHttpsTrafficOnly": true
      },
      "sku": {
        "name": "[variables('vmSizesMap')[variables('{{.Name}}Variables').VMSize].storageAccountType]"
      },
      "type": "Microsoft.Storage/storageAccounts"
    },
    {{if .HasDisks}}
    {
      "apiVersion": "[variables('apiVersionStorage')]",
      "copy": {
        "count": "[variables('{{.Name}}Variables').StorageAccountsCount]",
        "name": "datadiskLoop"
      },
      {{if not IsHostedMaster}}
        {{if not IsPrivateCluster}}
          "dependsOn": [
            "[concat('Microsoft.Network/publicIPAddresses/', variables('masterPublicIPAddressName'))]"
          ],
        {{end}}
      {{end}}
      "kind": "Storage",
      "location": "[variables('location')]",
      "name": "[concat(variables('storageAccountPrefixes')[mod(add(copyIndex(variables('dataStorageAccountPrefixSeed')),variables('{{.Name}}Variables').StorageAccountOffset),variables('storageAccountPrefixesCount'))],variables('storageAccountPrefixes')[div(add(copyIndex(variables('dataStorageAccountPrefixSeed')),variables('{{.Name}}Variables').StorageAccountOffset),variables('storageAccountPrefixesCount'))],variables('{{.Name}}DataAccountName'))]",
      "properties": {
        "encryption": {
          "keySource": "Microsoft.Storage",
          "services": {
            "blob": {
              "enabled": true
            },
            "file": {
              "enabled": true
            }
          }
        },
        "supportsHttpsTrafficOnly": true
      },
      "sku": {
        "name": "[variables('vmSizesMap')[variables('{{.Name}}Variables').VMSize].storageAccountType]"
      },
      "type": "Microsoft.Storage/storageAccounts"
    },
    {{end}}
    {
      "location": "[variables('location')]",
      "name": "[variables('{{.Name}}Variables').AvailabilitySet]",
      "apiVersion": "[variables('apiVersionCompute')]",
      "properties": {},
      "type": "Microsoft.Compute/availabilitySets"
    },
{{end}}
  {
      "apiVersion": "[variables('apiVersionCompute')]",
      "copy": {
        "count": "[sub(variables('{{.Name}}Count'), variables('{{.Name}}Variables').Offset)]",
        "name": "vmLoopNode"
      },
      "dependsOn": [
{{if .IsStorageAccount}}
        "[concat('Microsoft.Storage/storageAccounts/',variables('storageAccountPrefixes')[mod(add(div(copyIndex(variables('{{.Name}}Variables').Offset),variables('maxVMsPerStorageAccount')),variables('{{.Name}}Variables').StorageAccountOffset),variables('storageAccountPrefixesCount'))],variables('storageAccountPrefixes')[div(add(div(copyIndex(variables('{{.Name}}Variables').Offset),variables('maxVMsPerStorageAccount')),variables('{{.Name}}Variables').StorageAccountOffset),variables('storageAccountPrefixesCount'))],variables('{{.Name}}AccountName'))]",

  {{if .HasDisks}}
        "[concat('Microsoft.Storage/storageAccounts/',variables('storageAccountPrefixes')[mod(add(add(div(copyIndex(variables('{{.Name}}Variables').Offset),variables('maxVMsPerStorageAccount')),variables('{{.Name}}Variables').StorageAccountOffset),variables('dataStorageAccountPrefixSeed')),variables('storageAccountPrefixesCount'))],variables('storageAccountPrefixes')[div(add(add(div(copyIndex(variables('{{.Name}}Variables').Offset),variables('maxVMsPerStorageAccount')),variables('{{.Name}}Variables').StorageAccountOffset),variables('dataStorageAccountPrefixSeed')),variables('storageAccountPrefixesCount'))],variables('{{.Name}}DataAccountName'))]",
  {{end}}
{{end}}
        "[concat('Microsoft.Network/networkInterfaces/', variables('{{.Name}}Variables').VMNamePrefix, 'nic-', copyIndex(variables('{{.Name}}Variables').Offset))]",
        "[concat('Microsoft.Compute/availabilitySets/', variables('{{.Name}}Variables').AvailabilitySet)]"
      ],
      "tags":
      {
        "creationSource" : "[concat(parameters('generatorCode'), '-', variables('{{.Name}}Variables').VMNamePrefix, copyIndex(variables('{{.Name}}Variables').Offset))]",
        "resourceNameSuffix" : "[parameters('nameSuffix')]",
        "orchestrator" : "[variables('orchestratorNameVersionTag')]",
        "acsengineVersion" : "[parameters('acsengineVersion')]",
        "poolName" : "{{.Name}}"
      },
      "location": "[variables('location')]",
      "name": "[concat(variables('{{.Name}}Variables').VMNamePrefix, copyIndex(variables('{{.Name}}Variables').Offset))]",
      {{if UseManagedIdentity}}
      {{if UserAssignedIDEnabled}}
      "identity": {
        "type": "userAssigned",
        "userAssignedIdentities": {
          "[variables('userAssignedIDReference')]":{}
        }
      },
      {{else}}
      "identity": {
        "type": "systemAssigned"
      },
      {{end}}
      {{end}}
      {{if and IsOpenShift (not (UseAgentCustomImage .))}}
      "plan": {
        "name": "[variables('{{.Name}}osImageSKU')]",
        "publisher": "[variables('{{.Name}}osImagePublisher')]",
        "product": "[variables('{{.Name}}osImageOffer')]"
      },
      {{end}}
      "properties": {
        "availabilitySet": {
          "id": "[resourceId('Microsoft.Compute/availabilitySets',variables('{{.Name}}Variables').AvailabilitySet)]"
        },
        "hardwareProfile": {
          "vmSize": "[variables('{{.Name}}Variables').VMSize]"
        },
        "networkProfile": {
          "networkInterfaces": [
            {
              "id": "[resourceId('Microsoft.Network/networkInterfaces',concat(variables('{{.Name}}Variables').VMNamePrefix, 'nic-', copyIndex(variables('{{.Name}}Variables').Offset)))]"
            }
          ]
        },
        "osProfile": {
          "adminUsername": "[parameters('linuxAdminUsername')]",
          "computername": "[concat(variables('{{.Name}}Variables').VMNamePrefix, copyIndex(variables('{{.Name}}Variables').Offset))]",
          {{if not IsOpenShift}}
          {{GetKubernetesAgentCustomData .}}
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
            "id": "[resourceId(variables('{{.Name}}Variables').osImageResourceGroup, 'Microsoft.Compute/images', variables('{{.Name}}Variables').osImageName)]"
            {{else}}
            "offer": "[variables('{{.Name}}Variables').osImageOffer]",
            "publisher": "[variables('{{.Name}}Variables').osImagePublisher]",
            "sku": "[variables('{{.Name}}Variables').osImageSKU]",
            "version": "[variables('{{.Name}}Variables').osImageVersion]"
            {{end}}
          },
          "osDisk": {
            "createOption": "FromImage"
            ,"caching": "ReadWrite"
          {{if .IsStorageAccount}}
            ,"name": "[concat(variables('{{.Name}}Variables').VMNamePrefix, copyIndex(variables('{{.Name}}Variables').Offset),'-osdisk')]"
            ,"vhd": {
              "uri": "[concat(reference(concat('Microsoft.Storage/storageAccounts/',variables('storageAccountPrefixes')[mod(add(div(copyIndex(variables('{{.Name}}Variables').Offset),variables('maxVMsPerStorageAccount')),variables('{{.Name}}Variables').StorageAccountOffset),variables('storageAccountPrefixesCount'))],variables('storageAccountPrefixes')[div(add(div(copyIndex(variables('{{.Name}}Variables').Offset),variables('maxVMsPerStorageAccount')),variables('{{.Name}}Variables').StorageAccountOffset),variables('storageAccountPrefixesCount'))],variables('{{.Name}}AccountName')),variables('apiVersionStorage')).primaryEndpoints.blob,'osdisk/', variables('{{.Name}}Variables').VMNamePrefix, copyIndex(variables('{{.Name}}Variables').Offset), '-osdisk.vhd')]"
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
    {{if (not UserAssignedIDEnabled)}}
    {
      "apiVersion": "[variables('apiVersionAuthorization')]",
      "copy": {
         "count": "[sub(variables('{{.Name}}Count'), variables('{{.Name}}Variables').Offset)]",
         "name": "vmLoopNode"
       },
      "name": "[guid(concat('Microsoft.Compute/virtualMachines/', variables('{{.Name}}Variables').VMNamePrefix, copyIndex(variables('{{.Name}}Variables').Offset), 'vmidentity'))]",
      "type": "Microsoft.Authorization/roleAssignments",
      "properties": {
        "roleDefinitionId": "[variables('readerRoleDefinitionId')]",
        "principalId": "[reference(concat('Microsoft.Compute/virtualMachines/', variables('{{.Name}}Variables').VMNamePrefix, copyIndex(variables('{{.Name}}Variables').Offset)), '2017-03-30', 'Full').identity.principalId]"
      }
    },
    {{end}}
     {
       "type": "Microsoft.Compute/virtualMachines/extensions",
       "name": "[concat(variables('{{.Name}}Variables').VMNamePrefix, copyIndex(variables('{{.Name}}Variables').Offset), '/ManagedIdentityExtension')]",
       "copy": {
         "count": "[sub(variables('{{.Name}}Count'), variables('{{.Name}}Variables').Offset)]",
         "name": "vmLoopNode"
       },
       "apiVersion": "[variables('apiVersionCompute')]",
       "location": "[resourceGroup().location]",
       {{if UserAssignedIDEnabled}}
       "dependsOn": [
         "[concat('Microsoft.Compute/virtualMachines/', variables('{{.Name}}Variables').VMNamePrefix, copyIndex(variables('{{.Name}}Variables').Offset))]",
        "[concat('Microsoft.Authorization/roleAssignments/',guid(concat(variables('userAssignedID'), 'roleAssignment', resourceGroup().id)))]"
       ],
       {{else}}
       "dependsOn": [
         "[concat('Microsoft.Compute/virtualMachines/', variables('{{.Name}}Variables').VMNamePrefix, copyIndex(variables('{{.Name}}Variables').Offset))]",
         "[concat('Microsoft.Authorization/roleAssignments/', guid(concat('Microsoft.Compute/virtualMachines/', variables('{{.Name}}Variables').VMNamePrefix, copyIndex(variables('{{.Name}}Variables').Offset), 'vmidentity')))]"
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
      "apiVersion": "[variables('apiVersionCompute')]",
      "copy": {
        "count": "[sub(variables('{{.Name}}Count'), variables('{{.Name}}Variables').Offset)]",
        "name": "vmLoopNode"
      },
      "dependsOn": [
        {{if UseManagedIdentity}}
        "[concat('Microsoft.Compute/virtualMachines/', variables('{{.Name}}Variables').VMNamePrefix, copyIndex(variables('{{.Name}}Variables').Offset), '/extensions/ManagedIdentityExtension')]"
        {{else}}
        "[concat('Microsoft.Compute/virtualMachines/', variables('{{.Name}}Variables').VMNamePrefix, copyIndex(variables('{{.Name}}Variables').Offset))]"
        {{end}}
      ],
      "location": "[variables('location')]",
      "type": "Microsoft.Compute/virtualMachines/extensions",
      "name": "[concat(variables('{{.Name}}Variables').VMNamePrefix, copyIndex(variables('{{.Name}}Variables').Offset),'/cse', '-agent-', copyIndex(variables('{{.Name}}Variables').Offset))]",
      "properties": {
        "publisher": "Microsoft.Azure.Extensions",
        "type": "CustomScript",
        "typeHandlerVersion": "2.0",
        "autoUpgradeMinorVersion": true,
        "settings": {},
        "protectedSettings": {
          {{if eq $.Name "system"}}
          "commandToExecute": "[concat('retrycmd_if_failure() { r=$1; w=$2; t=$3; shift && shift && shift; for i in $(seq 1 $r); do timeout $t ${@}; [ $? -eq 0  ] && break || if [ $i -eq $r ]; then return 1; else sleep $w; fi; done };{{if not (IsFeatureEnabled "BlockOutboundInternet")}} ERR_OUTBOUND_CONN_FAIL=50; retrycmd_if_failure 50 1 3 nc -vz {{if IsMooncake}}gcr.azk8s.cn 80{{else}}k8s.gcr.io 443 && retrycmd_if_failure 50 1 3 nc -vz gcr.io 443 && retrycmd_if_failure 50 1 3 nc -vz docker.io 443{{end}} || exit $ERR_OUTBOUND_CONN_FAIL;{{end}} for i in $(seq 1 1200); do if [ -f /opt/azure/containers/provision.sh ]; then break; fi; if [ $i -eq 1200 ]; then exit 100; else sleep 1; fi; done; ', variables('provisionScriptParametersCommon'),' VNET_CNI_PLUGINS_URL=', parameters('vnetCniLinuxPluginsURL'),' GPU_NODE={{IsNSeriesSKU .}} /usr/bin/nohup /bin/bash -c \"/bin/bash /opt/azure/containers/provision.sh >> /var/log/azure/cluster-provision.log 2>&1{{if IsFeatureEnabled "CSERunInBackground" }} &{{end}}\"')]"
          {{else if eq $.Name "agentpool1"}}
          "commandToExecute": "[concat('retrycmd_if_failure() { r=$1; w=$2; t=$3; shift && shift && shift; for i in $(seq 1 $r); do timeout $t ${@}; [ $? -eq 0  ] && break || if [ $i -eq $r ]; then return 1; else sleep $w; fi; done };{{if not (IsFeatureEnabled "BlockOutboundInternet")}} ERR_OUTBOUND_CONN_FAIL=50; retrycmd_if_failure 50 1 3 nc -vz {{if IsMooncake}}gcr.azk8s.cn 80{{else}}k8s.gcr.io 443 && retrycmd_if_failure 50 1 3 nc -vz gcr.io 443 && retrycmd_if_failure 50 1 3 nc -vz docker.io 443{{end}} || exit $ERR_OUTBOUND_CONN_FAIL;{{end}} for i in $(seq 1 1200); do if [ -f /opt/azure/containers/provision.sh ]; then break; fi; if [ $i -eq 1200 ]; then exit 100; else sleep 1; fi; done; ', variables('provisionScriptParametersCommon'),' VNET_CNI_PLUGINS_URL=', parameters('vnetCniLinuxPluginsURL'),' GPU_NODE={{IsNSeriesSKU .}} /usr/bin/nohup /bin/bash -c \"/bin/bash /opt/azure/containers/provision.sh >> /var/log/azure/cluster-provision.log 2>&1{{if IsFeatureEnabled "CSERunInBackground" }} &{{end}}\"')]"
          {{else}}
          "commandToExecute": "[concat('retrycmd_if_failure() { r=$1; w=$2; t=$3; shift && shift && shift; for i in $(seq 1 $r); do timeout $t ${@}; [ $? -eq 0  ] && break || if [ $i -eq $r ]; then return 1; else sleep $w; fi; done };{{if not (IsFeatureEnabled "BlockOutboundInternet")}} ERR_OUTBOUND_CONN_FAIL=50; retrycmd_if_failure 50 1 3 nc -vz {{if IsMooncake}}gcr.azk8s.cn 80{{else}}k8s.gcr.io 443 && retrycmd_if_failure 50 1 3 nc -vz gcr.io 443 && retrycmd_if_failure 50 1 3 nc -vz docker.io 443{{end}} || exit $ERR_OUTBOUND_CONN_FAIL;{{end}} for i in $(seq 1 1200); do if [ -f /opt/azure/containers/provision.sh ]; then break; fi; if [ $i -eq 1200 ]; then exit 100; else sleep 1; fi; done; ', variables('provisionScriptParametersCommon'),' VNET_CNI_PLUGINS_URL=', parameters('vnetCniMultitenancyLinuxPluginsURL'),' GPU_NODE={{IsNSeriesSKU .}} /usr/bin/nohup /bin/bash -c \"/bin/bash /opt/azure/containers/provision.sh >> /var/log/azure/cluster-provision.log 2>&1{{if IsFeatureEnabled "CSERunInBackground" }} &{{end}}\"')]"
          {{end}}
        }
      }
    }
    {{if UseAksExtension}}
    ,{
      "type": "Microsoft.Compute/virtualMachines/extensions",
      "name": "[concat(variables('{{.Name}}Variables').VMNamePrefix, copyIndex(variables('{{.Name}}Variables').Offset), '/computeAksLinuxBilling')]",
      "apiVersion": "[variables('apiVersionCompute')]",
      "copy": {
        "count": "[sub(variables('{{.Name}}Count'), variables('{{.Name}}Variables').Offset)]",
        "name": "vmLoopNode"
      },
      "location": "[variables('location')]",
      "dependsOn": [
        "[concat('Microsoft.Compute/virtualMachines/', variables('{{.Name}}Variables').VMNamePrefix, copyIndex(variables('{{.Name}}Variables').Offset))]"
      ],
      "properties": {
        "publisher": "Microsoft.AKS",
        "type": {{if IsHostedMaster}}"Compute.AKS.Linux.Billing"{{else}}"Compute.AKS-Engine.Linux.Billing"{{end}},
        "typeHandlerVersion": "1.0",
        "autoUpgradeMinorVersion": true,
        "settings": {
        }
      }
    }
    {{end}}
    
