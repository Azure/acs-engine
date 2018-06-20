    {
      "type": "Microsoft.Network/networkSecurityGroups",
      "apiVersion": "[variables('apiVersionDefault')]",
      "location": "[variables('location')]",
      "name": "[variables('routerNSGName')]",
      "properties": {
        "securityRules": [
          {
            "name": "allow_http",
            "properties": {
              "access": "Allow",
              "description": "Allow http traffic to infra nodes",
              "destinationAddressPrefix": "*",
              "destinationPortRange": "80",
              "direction": "Inbound",
              "priority": 110,
              "protocol": "Tcp",
              "sourceAddressPrefix": "*",
              "sourcePortRange": "*"
            }
          },
          {
            "name": "allow_https",
            "properties": {
              "access": "Allow",
              "description": "Allow https traffic to infra nodes",
              "destinationAddressPrefix": "*",
              "destinationPortRange": "443",
              "direction": "Inbound",
              "priority": 111,
              "protocol": "Tcp",
              "sourceAddressPrefix": "*",
              "sourcePortRange": "*"
            }
          }
        ]
      }
    },
    {
        "name": "[variables('routerIPName')]",
        "type": "Microsoft.Network/publicIPAddresses",
        "apiVersion": "2017-08-01",
        "location": "[variables('location')]",
        "properties": {
            "publicIPAllocationMethod": "Static",
            "dnsSettings": {
              "domainNameLabel": "[concat(variables('masterFqdnPrefix'), '-router')]"
            }
        },
        "sku": {
            "name": "Basic"
        }
    },
    {
        "name": "[variables('routerLBName')]",
        "type": "Microsoft.Network/loadBalancers",
        "apiVersion": "2017-10-01",
        "location": "[variables('location')]",
        "dependsOn": [
            "[concat('Microsoft.Network/publicIPAddresses/', variables('routerIPName'))]"
        ],
        "properties": {
            "frontendIPConfigurations": [
                {
                    "name": "frontend",
                    "properties": {
                        "privateIPAllocationMethod": "Dynamic",
                        "publicIPAddress": {
                            "id": "[resourceId('Microsoft.Network/publicIPAddresses', variables('routerIPName'))]"
                        }
                    }
                }
            ],
            "backendAddressPools": [
                {
                    "name": "backend"
                }
            ],
            "loadBalancingRules": [
                {
                    "name": "port-80",
                    "properties": {
                        "frontendIPConfiguration": {
                            "id": "[concat(variables('routerLBID'), '/frontendIPConfigurations/frontend')]"
                        },
                        "frontendPort": 80,
                        "backendPort": 80,
                        "enableFloatingIP": false,
                        "idleTimeoutInMinutes": 4,
                        "protocol": "Tcp",
                        "loadDistribution": "Default",
                        "backendAddressPool": {
                            "id": "[concat(variables('routerLBID'), '/backendAddressPools/backend')]"
                        },
                        "probe": {
                            "id": "[concat(variables('routerLBID'), '/probes/port-80')]"
                        }
                    }
                },
                {
                    "name": "port-443",
                    "properties": {
                        "frontendIPConfiguration": {
                            "id": "[concat(variables('routerLBID'), '/frontendIPConfigurations/frontend')]"
                        },
                        "frontendPort": 443,
                        "backendPort": 443,
                        "enableFloatingIP": false,
                        "idleTimeoutInMinutes": 4,
                        "protocol": "Tcp",
                        "loadDistribution": "Default",
                        "backendAddressPool": {
                            "id": "[concat(variables('routerLBID'), '/backendAddressPools/backend')]"
                        },
                        "probe": {
                            "id": "[concat(variables('routerLBID'), '/probes/port-443')]"
                        }
                    }
                }
            ],
            "probes": [
                {
                    "name": "port-80",
                    "properties": {
                        "protocol": "Tcp",
                        "port": 80,
                        "intervalInSeconds": 5,
                        "numberOfProbes": 2
                    }
                },
                {
                    "name": "port-443",
                    "properties": {
                        "protocol": "Tcp",
                        "port": 443,
                        "intervalInSeconds": 5,
                        "numberOfProbes": 2
                    }
                }
            ],
            "inboundNatRules": [],
            "outboundNatRules": [],
            "inboundNatPools": []
        },
        "sku": {
            "name": "Basic"
        }
    },
    {
      "type": "Microsoft.Storage/storageAccounts",
      "apiVersion": "[variables('apiVersionStorage')]",
      "name": "[concat(variables('storageAccountBaseName'), 'registry')]",
      "location": "[variables('location')]",
      "properties": {
        "accountType": "Standard_LRS"
      }
    },