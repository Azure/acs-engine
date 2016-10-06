    {
      "apiVersion": "[variables('apiVersionDefault')]", 
      "copy": {
        "count": "[variables('{{.Name}}Count')]", 
        "name": "loop"
      }, 
      "dependsOn": [
{{if .IsCustomVNET}}
      "[concat('Microsoft.Network/networkSecurityGroups/', variables('masterNSGName'))]" 
{{else}}
      "[variables('vnetID')]"
{{end}}
      ], 
      "location": "[variables('location')]", 
      "name": "[concat(variables('{{.Name}}VMNamePrefix'), 'nic-', copyIndex())]", 
      "properties": {
{{if .IsCustomVNET}}                  
	    "networkSecurityGroup": {
		  "id": "[resourceId('Microsoft.Network/networkSecurityGroups/', variables('masterNSGName'))]"
	    },
{{end}}
        "ipConfigurations": [
          {
            "name": "ipconfig1", 
            "properties": {
              "privateIPAllocationMethod": "Dynamic", 
              "subnet": {
                "id": "[variables('subnetRef')]"
             }
            }
          }
        ],
        "enableIPForwarding": true
      }, 
      "type": "Microsoft.Network/networkInterfaces"
    },
    {
      "apiVersion": "[variables('apiVersionStorage')]", 
      "copy": {
        "count": "[variables('{{.Name}}StorageAccountsCount')]", 
        "name": "loop"
      }, 
      "dependsOn": [
        "[concat('Microsoft.Network/publicIPAddresses/', variables('masterPublicIPAddressName'))]"
      ], 
      "location": "[variables('location')]",
      "name": "[concat(variables('storageAccountPrefixes')[mod(add(copyIndex(),variables('{{.Name}}StorageAccountOffset')),variables('storageAccountPrefixesCount'))],variables('storageAccountPrefixes')[div(add(copyIndex(),variables('{{.Name}}StorageAccountOffset')),variables('storageAccountPrefixesCount'))],variables('{{.Name}}AccountName'))]",
      "properties": {
        "accountType": "[variables('vmSizesMap')[variables('{{.Name}}VMSize')].storageAccountType]"
      }, 
      "type": "Microsoft.Storage/storageAccounts"
    },
    {
      "apiVersion": "[variables('apiVersionStorage')]", 
      "copy": {
        "count": "[variables('{{.Name}}StorageAccountsCount')]", 
        "name": "datadiskLoop"
      }, 
      "dependsOn": [
        "[concat('Microsoft.Network/publicIPAddresses/', variables('masterPublicIPAddressName'))]"
      ], 
      "location": "[variables('location')]", 
      "name": "[concat(variables('storageAccountPrefixes')[mod(add(copyIndex(variables('dataStorageAccountPrefixSeed')),variables('{{.Name}}StorageAccountOffset')),variables('storageAccountPrefixesCount'))],variables('storageAccountPrefixes')[div(add(copyIndex(variables('dataStorageAccountPrefixSeed')),variables('{{.Name}}StorageAccountOffset')),variables('storageAccountPrefixesCount'))],variables('{{.Name}}DataAccountName'))]", 
      "properties": {
        "accountType": "[variables('vmSizesMap')[variables('{{.Name}}VMSize')].storageAccountType]"
      }, 
      "type": "Microsoft.Storage/storageAccounts"
    },
    {
      "apiVersion": "[variables('apiVersionDefault')]", 
      "location": "[variables('location')]",  
      "name": "[variables('{{.Name}}AvailabilitySet')]", 
      "properties": {}, 
      "type": "Microsoft.Compute/availabilitySets"
    }