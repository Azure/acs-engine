# Microsoft Azure Container Service Engine - Key vault referencing for k8s parameters

## Overview

ACS-Engine enables you to retrieve the following k8s deployment parameters from Microsoft Azure KeyVault:

*	apiServerCertificate
*	apiServerPrivateKey
*	caCertificate
*	clientCertificate
*	clientPrivateKey
*	kubeConfigCertificate
*	kubeConfigPrivateKey
*	servicePrincipalClientSecret

The parameters above could still be set as plain text.

To refer to a keyvault secret, the value of the parameter in the api model file should be formatted as: 

	"<PARAMETER>": "/subscriptions/<SUB_ID>/resourceGroups/<RG_NAME>/providers/Microsoft.KeyVault/vaults/<KV_NAME>/secrets/<NAME>[/<VERSION>]"
where: 
- **SUB_ID** is the subscription ID of the keyvault
- **RG_NAME** is the resource group of the keyvault
- **KV_NAME** is the name of the keyvault
- **NAME** is the name of the secret in the keyvault
- **VERSION** (optional) is the version of the secret (default: the latest version)

The example **kubernetes.json** shows you how to refer deployment parameter to a secret in a keyvault.
