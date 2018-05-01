{{ range $index, $agent := .AgentPoolProfiles }}
"{{.Name}}Config": {
    "Index": {{$index}},
    {{if .IsStorageAccount}}
        {{if .HasDisks}}
        "DataAccountName": "[concat(variables('storageAccountBaseName'), 'data{{$index}}')]",
        {{end}}
        "AccountName": "[concat(variables('storageAccountBaseName'), 'agnt{{$index}}')]",
    {{end}}
    {{if .IsStorageAccount}}
        "StorageAccountOffset": "[mul(variables('maxStorageAccountsPerAgent'),variables('{{.Name}}Config').Index)]",
        "StorageAccountsCount": "[add(div(variables('{{.Name}}Config').Count, variables('maxVMsPerStorageAccount')), mod(add(mod(variables('{{.Name}}Config').Count, variables('maxVMsPerStorageAccount')),2), add(mod(variables('{{.Name}}Config').Count, variables('maxVMsPerStorageAccount')),1)))]",
    {{end}}
        "Count": "[parameters('{{.Name}}Count')]",
    {{if .IsAvailabilitySets}}
        "Offset": "[parameters('{{.Name}}Offset')]",
        "AvailabilitySet": "[concat('{{.Name}}-availabilitySet-', variables('nameSuffix'))]",
    {{end}}
    {{if .IsWindows}}
        "winResourceNamePrefix" : "[substring(variables('nameSuffix'), 0, 5)]",
        "VMNamePrefix": "[concat(variables('winResourceNamePrefix'), variables('orchestratorName'), add(900, {{ $index }}))]",
    {{else}}
    {{if .IsAvailabilitySets}}
        "VMNamePrefix": "[concat(variables('orchestratorName'), '-{{.Name}}-', variables('nameSuffix'), '-')]",
    {{else}}
        "VMNamePrefix": "[concat(variables('orchestratorName'), '-{{.Name}}-', variables('nameSuffix'), '-vmss')]",
    {{end}}
    {{end}}
        "VMSize": "[parameters('{{.Name}}VMSize')]",
    {{if .IsCustomVNET}}
        "VnetSubnetID": "[parameters('{{.Name}}VnetSubnetID')]",
        "SubnetName": "[parameters('{{.Name}}VnetSubnetID')]",
        "VnetParts": "[split(parameters('{{.Name}}VnetSubnetID'),'/subnets/')]",
    {{else}}
        "VnetSubnetID": "[variables('vnetSubnetID')]",
        "SubnetName": "[variables('subnetName')]",
    {{end}}
        "osImageOffer": "[parameters('{{.Name}}osImageOffer')]",
        "osImageSKU": "[parameters('{{.Name}}osImageSKU')]",
        "osImagePublisher": "[parameters('{{.Name}}osImagePublisher')]",
        "osImageVersion": "[parameters('{{.Name}}osImageVersion')]",
        "osImageName": "[parameters('{{.Name}}osImageName')]",
        "osImageResourceGroup": "[parameters('{{.Name}}osImageResourceGroup')]"
},
{{end}}