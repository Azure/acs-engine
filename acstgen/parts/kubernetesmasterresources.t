    {
      "apiVersion": "[variables('computeApiVersion')]", 
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
        "[concat('Microsoft.Network/networkSecurityGroups/', variables('nsgName'))]", 
        "[concat('Microsoft.Network/routeTables/', variables('routeTableName'))]"
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
              "addressPrefix": "[variables('subnetCidr')]", 
              "networkSecurityGroup": {
                "id": "[resourceId('Microsoft.Network/networkSecurityGroups', variables('nsgName'))]"
              }, 
              "routeTable": {
                "id": "[resourceId('Microsoft.Network/routeTables', variables('routeTableName'))]"
              }
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
    {
      "apiVersion": "[variables('apiVersionDefault')]", 
      "location": "[variables('location')]", 
      "name": "[variables('routeTableName')]", 
      "type": "Microsoft.Network/routeTables"
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
      "dependsOn": [
        "[concat('Microsoft.Network/publicIPAddresses/', variables('masterPublicIPAddressName')]", 
        "[concat('Microsoft.Network/virtualNetworks/', variables('virtualNetworkName'))]"
      ], 
      "location": "[variables('location')]", 
      "name": "[concat(variables('masterFqdnPrefix'), '-nic-master')]", 
      "properties": {
        "ipConfigurations": [
          {
            "name": "ipconfig1", 
            "properties": {
              "privateIPAddress": "[variables('masterPrivateIp')]", 
              "privateIPAllocationMethod": "Static", 
              "publicIPAddress": {
                "id": "[resourceId('Microsoft.Network/publicIpAddresses', concat(variables('masterFqdnPrefix'), '-pip-master'))]"
              }, 
              "subnet": {
                "id": "[variables('subnetRef')]"
              }
            }
          }
        ],
        "enableIPForwarding": true
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
        "[concat('Microsoft.Storage/storageAccounts/', variables('masterStorageAccountName'))]", 
        "[concat('Microsoft.Network/networkInterfaces/', variables('masterFqdnPrefix'), '-nic-master')]"
      ], 
      "location": "[variables('location')]", 
      "name": "[concat(variables('masterVMNamePrefix'), copyIndex())]", 
      "properties": {
        "availabilitySet": {
          "id": "[resourceId('Microsoft.Compute/availabilitySets',variables('masterAvailabilitySet'))]"
        }, 
        "hardwareProfile": {
          "vmSize": "[variables('masterSize')]"
        }, 
        "networkProfile": {
          "networkInterfaces": [
            {
              "id": "[resourceId('Microsoft.Network/networkInterfaces',concat(variables('masterFqdnPrefix'),'-nic-master'))]"
            }
          ]
        }, 
        "osProfile": {
          "adminUsername": "[variables('username')]", 
          "computername": "[concat(variables('vmNamePrefix'), 'master')]", 
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
        }, 
        "storageProfile": {
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
    }