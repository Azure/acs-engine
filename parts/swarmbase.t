{
  "$schema": "https://schema.management.azure.com/schemas/2015-01-01/deploymentTemplate.json#",
  "contentVersion": "1.0.0.0",
  "parameters": {
    {{range .AgentPoolProfiles}}{{template "agentparams.t" .}},{{end}}
    {{if .HasWindows}}
      {{template "windowsparams.t"}},
    {{end}}
    {{template "masterparams.t" .}}
  },
  "variables": {
    {{range $index, $agent := .AgentPoolProfiles}}
        {{template "swarmagentvars.t" .}}
        {{if and .HasDisks .IsStorageAccount}}
          "{{.Name}}DataAccountName": "[concat(variables('storageAccountBaseName'), 'data{{$index}}')]",
        {{end}}
        "{{.Name}}Index": {{$index}},
        "{{.Name}}AccountName": "[concat(variables('storageAccountBaseName'), 'agnt{{$index}}')]",
    {{end}}

    {{if .HasWindows}}
      "windowsAdminUsername": "[parameters('windowsAdminUsername')]",
      "windowsAdminPassword": "[parameters('windowsAdminPassword')]",
      "agentWindowsPublisher": "MicrosoftWindowsServer",
      "agentWindowsOffer": "WindowsServer",
      "agentWindowsSku": "2016-Datacenter-with-Containers",
      "agentWindowsVersion": "latest",
      "singleQuote": "'",
      "windowsCustomScriptArguments": "[concat('$arguments = ', variables('singleQuote'),'-SwarmMasterIP ', variables('masterFirstAddrPrefix'), variables('masterFirstAddrOctet4'), variables('singleQuote'), ' ; ')]",
      "windowsCustomScriptSuffix": " $inputFile = '%SYSTEMDRIVE%\\AzureData\\CustomData.bin' ; $outputFile = '%SYSTEMDRIVE%\\AzureData\\CustomDataSetupScript.ps1' ; $inputStream = New-Object System.IO.FileStream $inputFile, ([IO.FileMode]::Open), ([IO.FileAccess]::Read), ([IO.FileShare]::Read) ; $sr = New-Object System.IO.StreamReader(New-Object System.IO.Compression.GZipStream($inputStream, [System.IO.Compression.CompressionMode]::Decompress)) ; $sr.ReadToEnd() | Out-File($outputFile) ; Invoke-Expression('{0} {1}' -f $outputFile, $arguments) ; ",
      "windowsCustomScript": "[concat('powershell.exe -ExecutionPolicy Unrestricted -command \"', variables('windowsCustomScriptArguments'), variables('windowsCustomScriptSuffix'), '\" > %SYSTEMDRIVE%\\AzureData\\CustomDataSetupScript.log 2>&1')]",
      "agentWindowsBackendPort": 3389,
    {{end}}

    {{template "swarmmastervars.t" .}},
    
    {{GetSizeMap}}
  },
  "resources": [
    {{range .AgentPoolProfiles}}
      {{if .IsWindows}}
        {{if .IsAvailabilitySets}}
          {{template "swarmwinagentresourcesvmas.t" .}},
        {{else}}
          {{template "swarmwinagentresourcesvmss.t" .}},
        {{end}}
      {{else}}
        {{if .IsAvailabilitySets}}
          {{template "swarmagentresourcesvmas.t" .}},
        {{else}}
          {{template "swarmagentresourcesvmss.t" .}},
        {{end}}
      {{end}}      
    {{end}}
    {{template "swarmmasterresources.t" .}}
  ],
  "outputs": {
    {{range .AgentPoolProfiles}}{{template "agentoutputs.t" .}}
    {{end}}{{template "masteroutputs.t" .}}
  }
}
