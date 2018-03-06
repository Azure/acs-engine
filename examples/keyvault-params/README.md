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
* etcdServerCertificate
* etcdServerPrivateKey
* etcdClientCertificate
* etcdClientPrivateKey
* etcdPeerCertificates
* etcdPeerPrivateKeys
*	servicePrincipal secret (a special case)

The parameters above could still be set as plain text.

To refer to a keyvault secret, the value of the parameter in the api model file should be formatted as:

	"<PARAMETER>": "/subscriptions/<SUB_ID>/resourceGroups/<RG_NAME>/providers/Microsoft.KeyVault/vaults/<KV_NAME>/secrets/<NAME>[/<VERSION>]"
where:
- **SUB_ID** is the subscription ID of the keyvault
- **RG_NAME** is the resource group of the keyvault
- **KV_NAME** is the name of the keyvault
- **NAME** is the name of the secret in the keyvault
- **VERSION** (optional) is the version of the secret (default: the latest version)

The treatment of servicePrincipal secret is slightly different. As a plain text, the secret is set in *servicePrincipalProfile.secret*. Alternatively, as a keyvault reference, the secret is set in *servicePrincipalProfile.keyvaultSecretRef*, which separates kevault ID, secret name, and optionally the version.
The keyvault ID format is *"/subscriptions/<SUB_ID>/resourceGroups/<RG_NAME>/providers/Microsoft.KeyVault/vaults/<KV_NAME>*

The example **kubernetes.json** shows you how to refer deployment parameter to a secret in a keyvault.

**Important** The secrets in the KeyVault for the Certificates and Private Keys must be Base64 encoded, and all on a single line -- this means you can't use the `--encoding base64` option of the Azure CLI. Instead you should use the `base64` command:

```sh
  # On OSX base64 will not wrap by default
  az keyvault secret set --vault-name KV_NAME --name NAME --value "$(cat ca.crt | base64 --break=0)"

  # On Linux it will wrap at 76 chars by default
  az keyvault secret set --vault-name KV_NAME --name NAME --value "$(cat ca.crt | base64 --wrap=0)"
```

## KeyVault Configuration
To enable Azure Resource Manager to retrieve the secrets from the KeyVault, template deployment must be enabled on the KeyVault:
```sh
az keyvault update -g $RG_NAME -n $KV_NAME --enabled-for-template-deployment
```
