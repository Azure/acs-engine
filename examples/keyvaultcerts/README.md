# Microsoft Azure Container Service Engine - Key vault certificate deployment

## Overview

ACS-Engine enables you to create customized Docker enabled cluster on Microsoft Azure with certs installed from key vault during deployment.

The examples show you how to configure installing a cert from keyvault. These certs are assumed to be in the secrets portion of your keyvault:

1. **dcos.json** - deploying and using [DC/OS](../../docs/dcos.md)
1. **kubernetes.json** - deploying and using [Kubernetes](../../docs/kubernetes.md)
1. **swarmmode.json** - deploying and using [Swarm Mode](../../docs/swarmmode.md)

On windows machines certificates will be installed under the machine in the specified store.
On linux machines the certificates will be installed in the folder /var/lib/waagent/. There will be two files
1. {thumbprint}.prv - this will be the private key pem formatted
2. {thumbprint}.crt - this will be the full cert chain pem formatted
