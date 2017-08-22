# Microsoft Azure Container Service Engine - Extensions

Extensions in acs-engine provide an easy way for acs-engine users to add pre-packaged functionality into their cluster.  For example, an extension could configure a monitoring solution on an ACS cluster.  The user would not need to know the details of how to install the monitoring solution.  Rather, the user would simply add the extension into the extensionsProfile section of the template.

# extensionsProfile

The extensionsProfile contains the extensions that the cluster will install. The following illustrates a template with a hello-world-k8s extension.

``` javascript
{ 
  ...
  "extensionsProfile": [
    {
        "name": "hello-world-k8s",
        "version": "v1",
        "extensionParameters": "parameters",
        "rootURL": "http://mytestlocation.com/hello-world-k8s/"
    }
  ]
}
```

|Name|Required|Description|
|---|---|---|
|name|yes|the name of the extension.  This has to exactly match the name of a folder under the extensions folder|
|version|yes|the version of the extension.  This has to exactly match the name of the folder under the extension name folder|
|extensionParameters|optional|extension parameters may be required by extensions.  The format of the parameters is also extension dependant. If the index in the vm pool is needed add EXTENSION_LOOP_INDEX at the location you wan the index and it will be replaced with the string representation of the index(zero based)|
|rootURL|optional|url to the root location of extensions.  The rootURL must have an extensions child folder that follows the extensions convention.  The rootURL is mainly used for testing purposes.|

# rootURL
You normally would not provide a rootURL.  The extensions are normally loaded from the extensions folder in GitHub.  However, you may specify the rootURL when testing a new extension.  The rootURL must adhere to the extensions conventions.  For example, in order to use an Azure Storage account to test an extension named extension-one, you would do the following:
- Create a storage account.  For the purposes of this example, we will call it 'mystorageaccount'
- Create a blob container called 'extensions'
- Under 'extensions', create a folder called 'extension-one'
- Under 'extension-one', create a folder called 'v1'
- Under 'v1', upload your files (see Required Extension Files) 
- Set the rootURL to: 'https://mystorageaccount.blob.core.windows.net/'

# masterProfile
Extensions, in the current implementation run a script on a master node. The extensions array in the masterProfile define that the master pool will have the script run on a single node on it.

``` javascript
{
  "masterProfile": {
      "count": 3,
      "dnsPrefix": "dnsprefix",
      "vmSize": "Standard_D2_v2",
      "osType": "Linux",
      "firstConsecutiveStaticIP": "10.240.255.5",
      "extensions": [
        { 
          "name": "hello-world-k8s"
        }
     ]
  },
  "extensionsProfile": [
    {
        "name": "hello-world-k8s",
        "version": "v1",
        "extensionParameters": "parameters"
    }
  ]
}
```

|Name|Required|Description|
|---|---|---|
|name|yes|The name of the extension. This must match the name in the extensionsProfile| 

# Required Extension Files

In order to install an extension, there are four required files - supported-orchestrators.json, template.json, template-link.json and EXTENSION-NAME.sh. Following is a description of each file.

|File Name|Description|
|-----------------------------|---|
|supported-orchestrators.json |Defines what orchestrators are supported by the extension (Swarm, Dcos, or Kubernetes)|
|template.json               |The ARM template used to deploy the extension|
|template-link.json          |The ARM template snippet which will be injected into azuredeploy.json to call template.json|
|EXTENSION-NAME.sh           |The script file that will execute on the VM itself via Custom Script Extension to perform installation of the extension|

# Creating supported-orchestrators.json

The supported-orchestrators.json file is a simple one line file that contains the list of supported orchestrators for which the extension can be installed into.

``` javascript
["Kubernetes"]
```

# Creating extension template.json

The template.json file is a linked template that will be called by the main cluster deployment template and must adhere to all the rules of a normal ARM template. All the necessary parameters needed from the azuredeploy.json file must be passed into this template and defined appropriately.

Additional variables can be defined for use in creating additional resources. Additional resources can also be created.  The key resource for installing the extension is the custom script extension. 

Modify the commandToExecute entry with the necessary command and paramters to install the desired extension. Replace EXTENSION-NAME with the name of the extension. The resource name of the custom script extension has to have the same name as the other custom script on the box as we aren't allowed to have two, this is also why we use a linked deployment so we can have the same resource twice and just make this one depend on the other so that it always runs after the provision extension is done.

The following is an example of the template.json file.

``` javascript
{
   "$schema": "http://schema.management.azure.com/schemas/2015-01-01/deploymentTemplate.json#",
   "contentVersion": "1.0.0.0",
   "parameters": {  
		"apiVersionStorage": {
			"type": "string",
			"minLength": 1,
			"metadata": {
				"description": "Storage API Version"
			}
		},
		"apiVersionDefault": {
			"type": "string",
			"minLength": 1,
			"metadata": {
				"description": "Compute API Version"
			}
		},
		"username": {
			"type": "string",
			"minLength": 1,
			"metadata": {
				"description": "Username for OS"
			}
		},
		"storageAccountBaseName": {
			"type": "string",
			"minLength": 1,
			"metadata": {
				"description": "Base Name of Storage Account"
			}
		},
		"extensionParameters": {
			"type": "securestring",
			"minLength": 1,
			"metadata": {
				"description": "Custom Parameter for Extension"
			}
		}
   },
   "variables": {  
		"singleQuote": "'",
		"sampleStorageAccountName": "[concat(uniqueString(concat(parameters('storageAccountBaseName'), 'sample')), 'aa')]"
		"initScriptUrl": "https://raw.githubusercontent.com/Azure/acs-engine/master/extensions/EXTENSION-NAME/v1/EXTENSION-NAME.sh"
   },
   "resources": [  
	{
      "apiVersion": "[parameters('apiVersionStorage')]", 
      "dependsOn": [], 
      "location": "[resourceGroup().location]", 
      "name": "[variables('sampleStorageAccountName')]", 
      "properties": {
        "accountType": "Standard_LRS"
      }, 
      "type": "Microsoft.Storage/storageAccounts"	
	}, {
      "apiVersion": "[parameters('apiVersionDefault')]",
      "dependsOn": [],
      "location": "[resourceGroup().location]",
      "type": "Microsoft.Compute/virtualMachines/extensions",
	  "name": "CustomExtension",
      "properties": {
        "publisher": "Microsoft.OSTCExtensions",
        "type": "CustomScriptForLinux",
        "typeHandlerVersion": "1.5",
        "autoUpgradeMinorVersion": true,
        "settings": {
			"fileUris": [ 
			   "[variables('initScriptUrl')]" 
			 ]
		},
        "protectedSettings": {
			"commandToExecute": "[concat('/bin/bash -c \"/bin/bash ./EXTENSION-NAME.sh ', variables('singleQuote'), parameters('extensionParameters'), variables('singleQuote'), ' ', variables('singleQuote'), parameters('sampleStorageAccountName'), variables('singleQuote'), ' >> /var/log/azure/sysdig-provision.log 2>&1 &\" &')]"
        }
      }
    }
	],
   "outputs": {  }
}
```
 
# Creating extension template-link.json

When acs-engine generates the azuredeploy.json file, this JSON snippet will be injected. This code calls the linked template (template.json) defined above.

Any parameters from the main azuredeploy.json file that is needed by template.json must be passed in via the parameters section. The parameter, "extensionParameters" is an optional parameter that is passed in directly by the user in the **extensionsProfile** section as defined in an earlier section. This special parameter can be used to pass in information such as an activation key or access code (as an example). If the extension does not need this capability, this optional parameter can be deleted.

Before this resource is created, all the dependencies must be satisfied first as defined by "dependsOn". The default dependency is that the entire cluster is fully provisioned before the script extension executes. This can be changed to meet your needs.

Replace "**EXTENSION-NAME**" with the name of the extension.

``` javascript
{
    "name": "EXTENSION-NAME",
    "type": "Microsoft.Resources/deployments",
    "apiVersion": "[variables('apiVersionLinkDefault')]",
    "dependsOn": [
        "vmLoopNode"
    ],
    "properties": {
        "mode": "Incremental",
        "templateLink": {
            "uri": "https://raw.githubusercontent.com/Azure/acs-engine/master/extensions/EXTENSION-NAME/v1/template.json",
            "contentVersion": "1.0.0.0"
        },
        "parameters": {
            "apiVersionStorage": {
                "value": "[variables('apiVersionStorage')]"
            },
            "apiVersionDefault": {
                "value": "[variables('apiVersionDefault')]"
            },
            "username": {
                "value": "[variables('username')]"
            },
            "storageAccountBaseName": {
                "value": "[variables('storageAccountBaseName')]"
            },
            "extensionParameters": {
                "value": "EXTENSION_PARAMETERS_REPLACE"
            }
        }
    }
}
```

# Creating extension script file

The script file will get executed on the VM to install the extension. Following is an example of a script.sh file.

``` bash
#!/bin/bash

# Add comments to explain the components of the script for easier troubleshooting
# Include echo statements so comments are written to output file
# Include necessary error checking

# Local variables

VARIABLE1=$1
VARIABLE2=$2
VARIABLE3=$3

echo $(date) " - Starting Script"

# Step 1 - example of creating config file
echo $(date) " - Creating sample.yaml file using local variables"

cat > sample.yaml <<EOF
line 1 $VARIABLE1
line 2 $VARIABLE2
line 3 $VARIABLE3
EOF

# Step 2 - example of downloading file
echo $(date) " - Downloading file"

curl -sSL http://example.com/sample/install

# Step 3 - example of executing other commands
echo $(date) " - Executing command"

install sample.yaml

echo $(date) " - Script complete"
```

# Current list of extensions
- [hello-world-dcos] (../extensions/hello-world-dcos/README.md)
- [hello-world-k8s] (../extensions/hello-world-k8s/README.md)