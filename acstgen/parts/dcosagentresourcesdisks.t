    {
      "apiVersion": "[variables('apiVersionDefault')]", 
      "location": "[resourceGroup().location]", 
      "name": "[variables('{{.Name}}NSGName')]", 
      "properties": {
        "securityRules": [
            {{GetSecurityRules .Ports}}
        ]
      }, 
      "type": "Microsoft.Network/networkSecurityGroups"
    },
    {
      "apiVersion": "[variables('apiVersionDefault')]", 
      "copy": {
        "count": "[variables('{{.Name}}Count')]", 
        "name": "loop"
      }, 
      "dependsOn": [
{{if .IsCustomVNET}}
      "[concat('Microsoft.Network/networkSecurityGroups/', variables('{{.Name}}NSGName'))]" 
{{else}}
      "[variables('vnetID')]"
{{end}}
{{if IsPublic .Ports}}
	  ,"[variables('{{.Name}}LbID')]"
{{end}}
      ], 
      "location": "[resourceGroup().location]", 
      "name": "[concat(variables('{{.Name}}VMNamePrefix'), 'nic-', copyIndex())]", 
      "properties": {
{{if .IsCustomVNET}}                  
	    "networkSecurityGroup": {
		  "id": "[resourceId('Microsoft.Network/networkSecurityGroups/', variables('{{.Name}}NSGName'))]"
	    },
{{end}}
        "ipConfigurations": [
          {
            "name": "ipConfigNode", 
            "properties": {
{{if IsPublic .Ports}}
              "loadBalancerBackendAddressPools": [
		        {
		      	  "id": "[concat('/subscriptions/', subscription().subscriptionId,'/resourceGroups/', resourceGroup().name, '/providers/Microsoft.Network/loadBalancers/', variables('{{.Name}}LbName'), '/backendAddressPools/',variables('{{.Name}}LbBackendPoolName'))]"
		        }
		      ], 
{{end}}  
              "privateIPAllocationMethod": "Dynamic", 
              "subnet": {
                "id": "[variables('{{.Name}}VnetSubnetID')]"
             }
            }
          }
        ]
      }, 
      "type": "Microsoft.Network/networkInterfaces"
    },
    {
      "apiVersion": "[variables('apiVersionStorage')]", 
      "copy": {
        "count": "[variables('{{.Name}}StorageAccountsCount')]", 
        "name": "loop"
      }, 
      "dependsOn": [
        "[concat('Microsoft.Network/publicIPAddresses/', variables('masterPublicIPAddressName'))]"
      ], 
      "location": "[variables('storageLocation')]", 
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
      "dependsOn": [
        "[concat('Microsoft.Network/publicIPAddresses/', variables('masterPublicIPAddressName'))]"
      ], 
      "location": "[variables('storageLocation')]", 
      "name": "[concat(variables('storageAccountPrefixes')[mod(add(copyIndex(variables('dataStorageAccountPrefixSeed')),variables('{{.Name}}StorageAccountOffset')),variables('storageAccountPrefixesCount'))],variables('storageAccountPrefixes')[div(add(copyIndex(variables('dataStorageAccountPrefixSeed')),variables('{{.Name}}StorageAccountOffset')),variables('storageAccountPrefixesCount'))],variables('{{.Name}}DataAccountName'))]", 
      "properties": {
        "accountType": "[variables('vmSizesMap')[variables('{{.Name}}VMSize')].storageAccountType]"
      }, 
      "type": "Microsoft.Storage/storageAccounts"
    }, 
{{end}}
    {
      "apiVersion": "[variables('apiVersionDefault')]", 
      "location": "[resourceGroup().location]", 
      "name": "[variables('{{.Name}}AvailabilitySet')]", 
      "properties": {}, 
      "type": "Microsoft.Compute/availabilitySets"
    },
{{if IsPublic .Ports}}
    {
      "apiVersion": "[variables('apiVersionDefault')]", 
      "location": "[resourceGroup().location]", 
      "name": "[variables('{{.Name}}IPAddressName')]", 
      "properties": {
        "dnsSettings": {
          "domainNameLabel": "[variables('{{.Name}}EndpointDNSNamePrefix')]"
        }, 
        "publicIPAllocationMethod": "Dynamic"
      }, 
      "type": "Microsoft.Network/publicIPAddresses"
    },
    {
      "apiVersion": "[variables('apiVersionDefault')]", 
      "dependsOn": [
        "[concat('Microsoft.Network/publicIPAddresses/', variables('{{.Name}}IPAddressName'))]"
      ], 
      "location": "[resourceGroup().location]", 
      "name": "[variables('{{.Name}}LbName')]", 
      "properties": {
        "backendAddressPools": [
          {
            "name": "[variables('{{.Name}}LbBackendPoolName')]"
          }
        ], 
        "frontendIPConfigurations": [
          {
            "name": "[variables('{{.Name}}LbIPConfigName')]", 
            "properties": {
              "publicIPAddress": {
                "id": "[resourceId('Microsoft.Network/publicIPAddresses',variables('{{.Name}}IPAddressName'))]"
              }
            }
          }
        ], 
        "inboundNatRules": [], 
        "loadBalancingRules": [
            {{(GetLBRules .Name .Ports)}}
        ], 
        "probes": [
            {{(GetProbes .Ports)}}
        ]
      }, 
      "type": "Microsoft.Network/loadBalancers"
    },
{{end}}
    {
      "apiVersion": "[variables('apiVersionDefault')]", 
      "copy": {
        "count": "[variables('{{.Name}}Count')]", 
        "name": "vmLoopNode"
      }, 
      "dependsOn": [
        "[concat('Microsoft.Storage/storageAccounts/',variables('storageAccountPrefixes')[mod(add(div(copyIndex(),variables('maxVMsPerStorageAccount')),variables('{{.Name}}StorageAccountOffset')),variables('storageAccountPrefixesCount'))],variables('storageAccountPrefixes')[div(add(div(copyIndex(),variables('maxVMsPerStorageAccount')),variables('{{.Name}}StorageAccountOffset')),variables('storageAccountPrefixesCount'))],variables('{{.Name}}AccountName'))]",
{{if .HasDisks}}
        "[concat('Microsoft.Storage/storageAccounts/',variables('storageAccountPrefixes')[mod(add(add(div(copyIndex(),variables('maxVMsPerStorageAccount')),variables('{{.Name}}StorageAccountOffset')),variables('dataStorageAccountPrefixSeed')),variables('storageAccountPrefixesCount'))],variables('storageAccountPrefixes')[div(add(add(div(copyIndex(),variables('maxVMsPerStorageAccount')),variables('{{.Name}}StorageAccountOffset')),variables('dataStorageAccountPrefixSeed')),variables('storageAccountPrefixesCount'))],variables('{{.Name}}DataAccountName'))]",
{{end}}
        "[concat('Microsoft.Network/networkInterfaces/', variables('{{.Name}}VMNamePrefix'), 'nic-', copyIndex())]", 
        "[concat('Microsoft.Compute/availabilitySets/', variables('{{.Name}}AvailabilitySet'))]"
      ], 
      "location": "[resourceGroup().location]", 
      "name": "[concat(variables('{{.Name}}VMNamePrefix'), copyIndex())]", 
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
              "id": "[resourceId('Microsoft.Network/networkInterfaces',concat(variables('{{.Name}}VMNamePrefix'), 'nic-', copyIndex()))]"
            }
          ]
        }, 
        "osProfile": {
          "adminUsername": "[variables('adminUsername')]", 
          "computername": "[concat(variables('{{.Name}}VMNamePrefix'), copyIndex())]", 
          "customData": "[base64(concat({{if IsDCOS173}}{{template "dcoscustomdata173.t" dict "DCOSCustomDataPublicIPStr" GetDCOSCustomDataPublicIPStr "DCOSGUID" GetDCOSGUID "RolesString" (GetAgentRolesFileContents .Ports)}}{{else if IsDCOS184}}{{template "dcoscustomdata184.t" dict "DCOSCustomDataPublicIPStr" GetDCOSCustomDataPublicIPStr "DCOSGUID" GetDCOSGUID "RolesString" (GetAgentRolesFileContents .Ports)}}{{end}}))]", 
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
            "caching": "ReadOnly", 
            "createOption": "FromImage", 
            "name": "[concat(variables('{{.Name}}VMNamePrefix'), copyIndex(),'-osdisk')]", 
            "vhd": {
              "uri": "[concat(reference(concat('Microsoft.Storage/storageAccounts/',variables('storageAccountPrefixes')[mod(add(div(copyIndex(),variables('maxVMsPerStorageAccount')),variables('{{.Name}}StorageAccountOffset')),variables('storageAccountPrefixesCount'))],variables('storageAccountPrefixes')[div(add(div(copyIndex(),variables('maxVMsPerStorageAccount')),variables('{{.Name}}StorageAccountOffset')),variables('storageAccountPrefixesCount'))],variables('{{.Name}}AccountName')),variables('apiVersionStorage')).primaryEndpoints.blob,'osdisk/', variables('{{.Name}}VMNamePrefix'), copyIndex(), '-osdisk.vhd')]"
            }
          }
        }
      }, 
      "type": "Microsoft.Compute/virtualMachines"
    }
