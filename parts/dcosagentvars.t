    "{{.Name}}Count": "[parameters('{{.Name}}Count')]",
    "{{.Name}}NSGID": "[resourceId('Microsoft.Network/networkSecurityGroups',variables('{{.Name}}NSGName'))]", 
    "{{.Name}}NSGName": "[concat(variables('orchestratorName'), '-{{.Name}}-nsg-', variables('nameSuffix'))]", 
{{if .IsWindows}}

    "winResourceNamePrefix" : "[substring(variables('nameSuffix'), 0, 5)]",
    "{{.Name}}VMNamePrefix": "[concat(variables('winResourceNamePrefix'), 'acs', add(900,variables('{{.Name}}Index')))]",
    {{if IsPublic .Ports}}
        "{{.Name}}windowsAgentCustomAttributes": "[concat(' -customAttrs ', variables('escapedSingleQuote'), '{{GetDCOSWindowsAgentCustomNodeAttributes . }}', variables('escapedSingleQuote') )]",
        "{{.Name}}windowsAgentCustomScriptArguments": "[concat('$arguments = ', variables('singleQuote'), '-subnet ', variables('{{.Name}}Subnet'), ' -MasterCount ', variables('masterCount'), ' -firstMasterIP ', parameters('firstConsecutiveStaticIP'), ' -bootstrapUri ', '\"', variables('dcosWindowsBootstrapURL'), '\"', ' -isAgent $true -isPublic $true ',  variables('{{.Name}}windowsAgentCustomAttributes'), variables('singleQuote'), ' ; ')]",
    {{else}}
        "{{.Name}}windowsAgentCustomScriptArguments": "[concat('$arguments = ', variables('singleQuote'),'-subnet ', variables('{{.Name}}Subnet'), ' -MasterCount ', variables('masterCount'), ' -firstMasterIP ', parameters('firstConsecutiveStaticIP'), ' -bootstrapUri ', '\"', variables('dcosWindowsBootstrapURL'), '\"', ' -isAgent $true',  variables('singleQuote'),' ; ')]",
    {{end}}
   
    "{{.Name}}windowsAgentCustomScript": "[concat('powershell.exe -ExecutionPolicy Unrestricted -command \"', variables('{{.Name}}windowsAgentCustomScriptArguments'), variables('windowsCustomScriptSuffix'), '\" > %SYSTEMDRIVE%\\AzureData\\dcosWindowsProvision.log 2>&1')]",

{{else}}
    "{{.Name}}VMNamePrefix": "[concat(variables('orchestratorName'), '-{{.Name}}-', variables('nameSuffix'), '-')]",
{{end}}

    "{{.Name}}VMSize": "[parameters('{{.Name}}VMSize')]",
    "{{.Name}}VMSizeTier": "[split(parameters('{{.Name}}VMSize'),'_')[0]]",
{{if .IsAvailabilitySets}}
    {{if .IsStorageAccount}}
    "{{.Name}}StorageAccountsCount": "[add(div(variables('{{.Name}}Count'), variables('maxVMsPerStorageAccount')), mod(add(mod(variables('{{.Name}}Count'), variables('maxVMsPerStorageAccount')),2), add(mod(variables('{{.Name}}Count'), variables('maxVMsPerStorageAccount')),1)))]",
    "{{.Name}}StorageAccountOffset": "[mul(variables('maxStorageAccountsPerAgent'),variables('{{.Name}}Index'))]",
    {{end}}
    "{{.Name}}AvailabilitySet": "[concat('{{.Name}}-availabilitySet-', variables('nameSuffix'))]",
    "{{.Name}}Offset": "[parameters('{{.Name}}Offset')]",
{{else}}
    {{if .IsStorageAccount}}
    "{{.Name}}StorageAccountsCount": "[variables('maxStorageAccountsPerAgent')]",
    {{end}}
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
     {{if .IsWindows}}
        "{{.Name}}WindowsRDPNatRangeStart": 3389,
        "{{.Name}}WindowsRDPEndRangeStop": "[add(variables('{{.Name}}WindowsRDPNatRangeStart'), add(variables('{{.Name}}Count'),variables('{{.Name}}Count')))]",

    {{end}}
{{end}}
