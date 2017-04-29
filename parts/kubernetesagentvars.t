    "{{.Name}}StorageAccountOffset": "[mul(variables('maxStorageAccountsPerAgent'),variables('{{.Name}}Index'))]",
    "{{.Name}}Count": "[parameters('{{.Name}}Count')]",
    "{{.Name}}Offset": "[parameters('{{.Name}}Offset')]",
    "{{.Name}}AvailabilitySet": "[concat('{{.Name}}-availabilitySet-', variables('nameSuffix'))]",
    "{{.Name}}StorageAccountsCount": "[add(div(variables('{{.Name}}Count'), variables('maxVMsPerStorageAccount')), mod(add(mod(variables('{{.Name}}Count'), variables('maxVMsPerStorageAccount')),2), add(mod(variables('{{.Name}}Count'), variables('maxVMsPerStorageAccount')),1)))]",
{{if .IsWindows}}
    "winResourceNamePrefix" : "[substring(variables('nameSuffix'), 0, 5)]",
    "{{.Name}}VMNamePrefix": "[concat(variables('winResourceNamePrefix'), 'acs', add(900,variables('{{.Name}}Index')))]",
{{else}}
    "{{.Name}}VMNamePrefix": "[concat(variables('orchestratorName'), '-{{.Name}}-', variables('nameSuffix'), '-')]", 
{{end}}
    "{{.Name}}VMSize": "[parameters('{{.Name}}VMSize')]",
{{if .IsCustomVNET}}
    "{{.Name}}VnetSubnetID": "[parameters('{{.Name}}VnetSubnetID')]",
    "{{.Name}}SubnetName": "[parameters('{{.Name}}VnetSubnetID')]",
    "{{.Name}}VnetParts": "[split(parameters('{{.Name}}VnetSubnetID'),'/subnets/')]",
{{else}}
    "{{.Name}}VnetSubnetID": "[variables('vnetSubnetID')]",
    "{{.Name}}SubnetName": "[variables('subnetName')]",
{{end}}