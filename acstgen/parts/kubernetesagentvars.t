    "{{.Name}}StorageAccountOffset": "[mul(variables('maxStorageAccountsPerAgent'),variables('{{.Name}}Index'))]",
    "{{.Name}}Count": "[parameters('{{.Name}}Count')]",
{{if .IsStateful}}
    "{{.Name}}AvailabilitySet": "[concat('{{.Name}}-availabilitySet-', variables('nameSuffix'))]",
    "{{.Name}}StorageAccountsCount": "[add(div(variables('{{.Name}}Count'), variables('maxVMsPerStorageAccount')), mod(add(mod(variables('{{.Name}}Count'), variables('maxVMsPerStorageAccount')),2), add(mod(variables('{{.Name}}Count'), variables('maxVMsPerStorageAccount')),1)))]",
{{else}}
    "{{.Name}}StorageAccountsCount": "[variables('maxStorageAccountsPerAgent')]",
{{end}}  
    "{{.Name}}VMNamePrefix": "[concat(variables('orchestratorName'), '-{{.Name}}-', variables('nameSuffix'), '-')]", 
    "{{.Name}}VMSize": "[parameters('{{.Name}}VMSize')]",
