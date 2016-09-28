    "{{.Name}}StorageAccountOffset": "[mul(variables('agentStorageAccountsCount'),variables('{{.Name}}Index'))]",
    "{{.Name}}Count": "[parameters('{{.Name}}Count')]",  
    "{{.Name}}NSGID": "[resourceId('Microsoft.Network/networkSecurityGroups',variables('{{.Name}}NSGName'))]", 
    "{{.Name}}NSGName": "[concat(variables('orchestratorName'), '-{{.Name}}-nsg-', variables('nameSuffix'))]", 
    "{{.Name}}VMNamePrefix": "[concat(variables('orchestratorName'), '-{{.Name}}-', variables('nameSuffix'))]", 
    "{{.Name}}VMSize": "[parameters('{{.Name}}VMSize')]",
{{if .IsCustomVNET}}
    "{{.Name}}VnetSubnetID": "[parameters('{{.Name}}VnetSubnetID')]",
{{else}}
    "{{.Name}}Subnet": "[parameters('{{.Name}}Subnet')]",
    "{{.Name}}SubnetName": "[concat(variables('orchestratorName'), '-{{.Name}}Subnet')]",
    "{{.Name}}VnetSubnetID": "[concat(variables('vnetID'),'/subnets/',variables('{{.Name}}SubnetName'))]",
{{end}}
{{if IsPublic .Ports}}
    "{{.Name}}EndpointDNSNamePrefix": "[tolower(parameters('{{.Name}}EndpointDNSNamePrefix'))]",
    "{{.Name}}IPAddressName": "[concat(variables('orchestratorName'), '-agent-ip-', variables('{{.Name}}EndpointDNSNamePrefix'), '-', variables('nameSuffix'))]",
    "{{.Name}}LbBackendPoolName": "[concat(variables('orchestratorName'), '-{{.Name}}-', variables('nameSuffix'))]", 
    "{{.Name}}LbID": "[resourceId('Microsoft.Network/loadBalancers',variables('{{.Name}}LbName'))]", 
    "{{.Name}}LbIPConfigID": "[concat(variables('{{.Name}}LbID'),'/frontendIPConfigurations/', variables('{{.Name}}LbIPConfigName'))]", 
    "{{.Name}}LbIPConfigName": "[concat(variables('orchestratorName'), '-{{.Name}}-', variables('nameSuffix'))]", 
    "{{.Name}}LbName": "[concat(variables('orchestratorName'), '-{{.Name}}-', variables('nameSuffix'))]",
{{end}}