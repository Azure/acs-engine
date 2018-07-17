# Azure Key Vault FlexVolume Add-on

[The Azure Key Vault FlexVolume](https://github.com/Azure/kubernetes-keyvault-flexvol) integrates Azure Key Vault with Kubernetes via a FlexVolume.  

With the Azure Key Vault FlexVolume, developers can access application-specific secrets, keys, and certs stored in Azure Key Vault directly from their pods.

Add this add-on to your apimodel as shown below to automatically enable Key Vault FlexVolume in your new Kubernetes cluster.

```json
{
    "apiVersion": "vlabs",
    "properties": {
      "orchestratorProfile": {
        "orchestratorType": "Kubernetes",
        "kubernetesConfig": {
          "addons": [
            {
              "name": "keyvault-flexvolume",
              "enabled" : true
            }
          ]
        }
      },
      "masterProfile": {
        "count": 1,
        "dnsPrefix": "",
        "vmSize": "Standard_DS2_v2",
      },
      "agentPoolProfiles": [
        {
          "name": "agentpool",
          "count": 3,
          "vmSize": "Standard_DS2_v2",
          "availabilityProfile": "VirtualMachineScaleSets"
        }
      ],
      "linuxProfile": {
        "adminUsername": "azureuser",
        "ssh": {
          "publicKeys": [
            {
              "keyData": ""
            }
          ]
        }
      },
      "servicePrincipalProfile": {
        "clientId": "",
        "secret": ""
      }
    }
  }

```

To validate the add-on is running as expected, run the following commands:

You should see the keyvault flexvolume installer pods running on each agent node:

```bash
kubectl get pods -n kv

kv-flexvol-installer-f7bx8   1/1       Running   0          3m
kv-flexvol-installer-rcxbl   1/1       Running   0          3m
kv-flexvol-installer-z6jm6   1/1       Running   0          3m
```

Follow the README at https://github.com/Azure/kubernetes-keyvault-flexvol for get started steps.

## Supported Orchestrators

Kubernetes