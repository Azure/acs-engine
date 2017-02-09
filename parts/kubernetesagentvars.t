    "{{.Name}}StorageAccountOffset": "[mul(variables('maxStorageAccountsPerAgent'),variables('{{.Name}}Index'))]",
    "{{.Name}}Count": "[parameters('{{.Name}}Count')]",
    "{{.Name}}AvailabilitySet": "[concat('{{.Name}}-availabilitySet-', variables('nameSuffix'))]",
    "{{.Name}}StorageAccountsCount": "[variables('maxStorageAccountsPerAgent')]",
    "{{.Name}}VMNamePrefix": "[concat(variables('orchestratorName'), '-{{.Name}}-', variables('nameSuffix'), '-')]", 
    "{{.Name}}VMSize": "[parameters('{{.Name}}VMSize')]",
{{if .IsWindows}}
    "{{.Name}}IPAddressName": "[concat(variables('orchestratorName'), '-agent-ip-', variables('{{.Name}}Index'), '-', variables('nameSuffix'))]",
    "{{.Name}}LbBackendPoolName": "[concat(variables('orchestratorName'), '-{{.Name}}-', variables('nameSuffix'))]", 
    "{{.Name}}LbID": "[resourceId('Microsoft.Network/loadBalancers',variables('{{.Name}}LbName'))]", 
    "{{.Name}}LbIPConfigID": "[concat(variables('{{.Name}}LbID'),'/frontendIPConfigurations/', variables('{{.Name}}LbIPConfigName'))]", 
    "{{.Name}}LbIPConfigName": "[concat(variables('orchestratorName'), '-{{.Name}}-', variables('nameSuffix'))]", 
    "{{.Name}}LbName": "[concat(variables('orchestratorName'), '-{{.Name}}-', variables('nameSuffix'))]",
    "{{.Name}}WindowsRDPNatRangeStart": 3389,
    "{{.Name}}WindowsRDPEndRangeStop": "[add(variables('{{.Name}}WindowsRDPNatRangeStart'), add(variables('{{.Name}}Count'),variables('{{.Name}}Count')))]",
{{end}}
{{if .IsCustomVNET}}
    "{{.Name}}VnetSubnetID": "[parameters('{{.Name}}VnetSubnetID')]",
    "{{.Name}}SubnetName": "[parameters('{{.Name}}VnetSubnetID')]",
    "{{.Name}}VnetParts": "[split(parameters('{{.Name}}VnetSubnetID'),'/subnets/')]",
{{else}}
    "{{.Name}}VnetSubnetID": "[variables('vnetSubnetID')]",
    "{{.Name}}SubnetName": "[variables('subnetName')]",
{{end}}