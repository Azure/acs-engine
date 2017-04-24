    "{{.Name}}Count": "[parameters('{{.Name}}Count')]",
    "{{.Name}}NSGID": "[resourceId('Microsoft.Network/networkSecurityGroups',variables('{{.Name}}NSGName'))]", 
    "{{.Name}}NSGName": "[concat(variables('orchestratorName'), '-{{.Name}}-nsg-', variables('nameSuffix'))]", 
    "{{.Name}}VMNamePrefix": "[concat(variables('orchestratorName'), '-{{.Name}}-', variables('nameSuffix'))]", 
    "{{.Name}}VMSize": "[parameters('{{.Name}}VMSize')]",
    "{{.Name}}VMSizeTier": "[split(parameters('{{.Name}}VMSize'),'_')[0]]",
{{if .IsAvailabilitySets}}
    {{if .IsStorageAccount}}
    "{{.Name}}StorageAccountsCount": "[add(div(variables('{{.Name}}Count'), variables('maxVMsPerStorageAccount')), mod(add(mod(variables('{{.Name}}Count'), variables('maxVMsPerStorageAccount')),2), add(mod(variables('{{.Name}}Count'), variables('maxVMsPerStorageAccount')),1)))]",
    {{end}}
    "{{.Name}}AvailabilitySet": "[concat('{{.Name}}-availabilitySet-', variables('nameSuffix'))]",
    "{{.Name}}Offset": "[parameters('{{.Name}}Offset')]",
{{else}}
    "{{.Name}}StorageAccountsCount": "[variables('maxStorageAccountsPerAgent')]",
{{end}}
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