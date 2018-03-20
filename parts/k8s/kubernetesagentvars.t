{{if .IsStorageAccount}}
    "{{.Name}}StorageAccountOffset": "[mul(variables('maxStorageAccountsPerAgent'),variables('{{.Name}}Index'))]",
    "{{.Name}}StorageAccountsCount": "[add(div(variables('{{.Name}}Count'), variables('maxVMsPerStorageAccount')), mod(add(mod(variables('{{.Name}}Count'), variables('maxVMsPerStorageAccount')),2), add(mod(variables('{{.Name}}Count'), variables('maxVMsPerStorageAccount')),1)))]",
{{end}}
    "{{.Name}}Count": "[parameters('{{.Name}}Count')]",
    "{{.Name}}Offset": "[parameters('{{.Name}}Offset')]",
    "{{.Name}}AvailabilitySet": "[concat(parameters('clusterName'), '-{{.Name}}-availabilitySet')]",
{{if .IsWindows}}
    "winResourceNamePrefix" : "[substring(variables('nameSuffix'), 0, 5)]",
    "{{.Name}}VMNamePrefix": "[concat(variables('winResourceNamePrefix'), parameters('clusterName'), add(900,variables('{{.Name}}Index')))]",
{{else}}
    "{{.Name}}VMNamePrefix": "[concat(parameters('clusterName'), '-{{.Name}}-')]",
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
    "{{.Name}}osImageOffer": "[parameters('{{.Name}}osImageOffer')]",
    "{{.Name}}osImageSKU": "[parameters('{{.Name}}osImageSKU')]",
    "{{.Name}}osImagePublisher": "[parameters('{{.Name}}osImagePublisher')]",
    "{{.Name}}osImageVersion": "[parameters('{{.Name}}osImageVersion')]",

