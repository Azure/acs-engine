    "adminUsername": "[parameters('linuxAdminUsername')]",
    "maxVMsPerPool": 100,
    "maxVMsPerStorageAccount": 20,
    "maxStorageAccountsPerAgent": "[div(variables('maxVMsPerPool'),variables('maxVMsPerStorageAccount'))]",
    "dataStorageAccountPrefixSeed": 97, 
    "apiVersionDefault": "2016-03-30", 
    "apiVersionStorage": "2015-06-15", 
    "masterAvailabilitySet": "[concat(variables('orchestratorName'), '-master-availabilitySet-', variables('nameSuffix'))]", 
    "masterCount": {{.MasterProfile.Count}}, 
    "masterEndpointDNSNamePrefix": "[tolower(parameters('masterEndpointDNSNamePrefix'))]",
    "masterLbBackendPoolName": "[concat(variables('orchestratorName'), '-master-pool-', variables('nameSuffix'))]", 
    "masterLbID": "[resourceId('Microsoft.Network/loadBalancers',variables('masterLbName'))]", 
    "masterLbIPConfigID": "[concat(variables('masterLbID'),'/frontendIPConfigurations/', variables('masterLbIPConfigName'))]", 
    "masterLbIPConfigName": "[concat(variables('orchestratorName'), '-master-lbFrontEnd-', variables('nameSuffix'))]", 
    "masterLbName": "[concat(variables('orchestratorName'), '-master-lb-', variables('nameSuffix'))]", 
    "masterNSGID": "[resourceId('Microsoft.Network/networkSecurityGroups',variables('masterNSGName'))]", 
    "masterNSGName": "[concat(variables('orchestratorName'), '-master-nsg-', variables('nameSuffix'))]", 
    "masterPublicIPAddressName": "[concat(variables('orchestratorName'), '-master-ip-', variables('masterEndpointDNSNamePrefix'), '-', variables('nameSuffix'))]", 
    "masterStorageAccountExhibitorName": "[concat(variables('storageAccountBaseName'), 'exhb0')]", 
    "masterStorageAccountName": "[concat(variables('storageAccountBaseName'), 'mstr0')]",
{{if .MasterProfile.IsCustomVNET}}
    "masterVnetSubnetID": "[parameters('masterVnetSubnetID')]",
{{else}}
    "masterSubnet": "[parameters('masterSubnet')]",
    "masterSubnetName": "[concat(variables('orchestratorName'), '-masterSubnet')]",
    "vnetID": "[resourceId('Microsoft.Network/virtualNetworks',variables('virtualNetworkName'))]",
    "masterVnetSubnetID": "[concat(variables('vnetID'),'/subnets/',variables('masterSubnetName'))]",
    "virtualNetworkName": "[concat(variables('orchestratorName'), '-vnet-', variables('nameSuffix'))]", 
{{end}}
    "masterFirstAddrOctets": "[split(parameters('firstConsecutiveStaticIP'),'.')]",
    "masterFirstAddrOctet4": "[variables('masterFirstAddrOctets')[3]]",
    "masterFirstAddrPrefix": "[concat(variables('masterFirstAddrOctets')[0],'.',variables('masterFirstAddrOctets')[1],'.',variables('masterFirstAddrOctets')[2],'.')]",
    "masterVMNamePrefix": "[concat(variables('orchestratorName'), '-master-', variables('nameSuffix'), '-')]", 
    "masterVMNic": [
      "[concat(variables('masterVMNamePrefix'), 'nic-0')]", 
      "[concat(variables('masterVMNamePrefix'), 'nic-1')]", 
      "[concat(variables('masterVMNamePrefix'), 'nic-2')]", 
      "[concat(variables('masterVMNamePrefix'), 'nic-3')]", 
      "[concat(variables('masterVMNamePrefix'), 'nic-4')]", 
      "[concat(variables('masterVMNamePrefix'), 'nic-5')]", 
      "[concat(variables('masterVMNamePrefix'), 'nic-6')]"
    ], 
    "masterVMSize": "[parameters('masterVMSize')]", 
    "nameSuffix": "{{GetUniqueNameSuffix}}", 
    "oauthEnabled": "false", 
    "orchestratorName": "dcos", 
    "osImageOffer": "UbuntuServer", 
    "osImagePublisher": "Canonical", 
    "osImageSKU": "16.04.0-LTS", 
    "osImageVersion": "16.04.201606270", 
    "sshKeyPath": "[concat('/home/', variables('adminUsername'), '/.ssh/authorized_keys')]", 
    "sshRSAPublicKey": "[parameters('sshRSAPublicKey')]", 
    "storageAccountBaseName": "[uniqueString(concat(variables('masterEndpointDNSNamePrefix'),resourceGroup().location, variables('orchestratorName')))]", 
    "storageAccountPrefixes": [
      "0", 
      "6", 
      "c", 
      "i", 
      "o", 
      "u", 
      "1", 
      "7", 
      "d", 
      "j", 
      "p", 
      "v", 
      "2", 
      "8", 
      "e", 
      "k", 
      "q", 
      "w", 
      "3", 
      "9", 
      "f", 
      "l", 
      "r", 
      "x", 
      "4", 
      "a", 
      "g", 
      "m", 
      "s", 
      "y", 
      "5", 
      "b", 
      "h", 
      "n", 
      "t", 
      "z"
    ], 
    "storageAccountPrefixesCount": "[length(variables('storageAccountPrefixes'))]", 
    "storageAccountType": "Standard_LRS", 
    "storageLocation": "[resourceGroup().location]", 
    "vmSizesMap": {
      "Standard_A0": {
        "storageAccountType": "Standard_LRS"
      }, 
      "Standard_A1": {
        "storageAccountType": "Standard_LRS"
      }, 
      "Standard_A10": {
        "storageAccountType": "Standard_LRS"
      }, 
      "Standard_A11": {
        "storageAccountType": "Standard_LRS"
      }, 
      "Standard_A2": {
        "storageAccountType": "Standard_LRS"
      }, 
      "Standard_A3": {
        "storageAccountType": "Standard_LRS"
      }, 
      "Standard_A4": {
        "storageAccountType": "Standard_LRS"
      }, 
      "Standard_A5": {
        "storageAccountType": "Standard_LRS"
      }, 
      "Standard_A6": {
        "storageAccountType": "Standard_LRS"
      }, 
      "Standard_A7": {
        "storageAccountType": "Standard_LRS"
      }, 
      "Standard_A8": {
        "storageAccountType": "Standard_LRS"
      }, 
      "Standard_A9": {
        "storageAccountType": "Standard_LRS"
      }, 
      "Standard_D1": {
        "storageAccountType": "Standard_LRS"
      }, 
      "Standard_D11": {
        "storageAccountType": "Standard_LRS"
      }, 
      "Standard_D11_v2": {
        "storageAccountType": "Standard_LRS"
      }, 
      "Standard_D12": {
        "storageAccountType": "Standard_LRS"
      }, 
      "Standard_D12_v2": {
        "storageAccountType": "Standard_LRS"
      }, 
      "Standard_D13": {
        "storageAccountType": "Standard_LRS"
      }, 
      "Standard_D13_v2": {
        "storageAccountType": "Standard_LRS"
      }, 
      "Standard_D14": {
        "storageAccountType": "Standard_LRS"
      }, 
      "Standard_D14_v2": {
        "storageAccountType": "Standard_LRS"
      }, 
      "Standard_D1_v2": {
        "storageAccountType": "Standard_LRS"
      }, 
      "Standard_D2": {
        "storageAccountType": "Standard_LRS"
      }, 
      "Standard_D2_v2": {
        "storageAccountType": "Standard_LRS"
      }, 
      "Standard_D3": {
        "storageAccountType": "Standard_LRS"
      }, 
      "Standard_D3_v2": {
        "storageAccountType": "Standard_LRS"
      }, 
      "Standard_D4": {
        "storageAccountType": "Standard_LRS"
      }, 
      "Standard_D4_v2": {
        "storageAccountType": "Standard_LRS"
      }, 
      "Standard_D5_v2": {
        "storageAccountType": "Standard_LRS"
      }, 
      "Standard_DS1": {
        "storageAccountType": "Premium_LRS"
      }, 
      "Standard_DS11": {
        "storageAccountType": "Premium_LRS"
      }, 
      "Standard_DS12": {
        "storageAccountType": "Premium_LRS"
      }, 
      "Standard_DS13": {
        "storageAccountType": "Premium_LRS"
      }, 
      "Standard_DS14": {
        "storageAccountType": "Premium_LRS"
      }, 
      "Standard_DS2": {
        "storageAccountType": "Premium_LRS"
      }, 
      "Standard_DS3": {
        "storageAccountType": "Premium_LRS"
      }, 
      "Standard_DS4": {
        "storageAccountType": "Premium_LRS"
      }, 
      "Standard_G1": {
        "storageAccountType": "Standard_LRS"
      }, 
      "Standard_G2": {
        "storageAccountType": "Standard_LRS"
      }, 
      "Standard_G3": {
        "storageAccountType": "Standard_LRS"
      }, 
      "Standard_G4": {
        "storageAccountType": "Standard_LRS"
      }, 
      "Standard_G5": {
        "storageAccountType": "Standard_LRS"
      }, 
      "Standard_GS1": {
        "storageAccountType": "Premium_LRS"
      }, 
      "Standard_GS2": {
        "storageAccountType": "Premium_LRS"
      }, 
      "Standard_GS3": {
        "storageAccountType": "Premium_LRS"
      }, 
      "Standard_GS4": {
        "storageAccountType": "Premium_LRS"
      }, 
      "Standard_GS5": {
        "storageAccountType": "Premium_LRS"
      }
    }
