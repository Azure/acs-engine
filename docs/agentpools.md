# Deploying Agent Pools Only

If you chose to deploy agent pools for Kubernetes in Azure you can use the `acs-engine generate agentpool` command.

The command's only required input paramter is a path to an `kubernetesagentpool` JSON file to read from

An example is below:

```json
{
  "apiVersion": "v20170727",
  "location": "eastus",
  "name": "my-k8s-agentpools",
  "properties": {
    "kubernetesEndpoint": "api.mydomain.io",
    "kubernetesVersion": "1.7.2",
    "dnsPrefix": "my-k8s-agentpool",
    "agentPoolProfiles": [
      {
        "name": "agentpool1",
        "count": 3,
        "vmSize": "Standard_D2_v2",
        "availabilityProfile": "AvailabilitySet"
      }
    ],
    "linuxProfile": {
      "adminUsername": "azureuser",
      "ssh": {
        "publicKeys": [
          {
            "keyData": "ssh-rsa ..."
          }
        ]
      }
    },
    "JumpboxProfile": {
      "publicIpAddressId": "",
      "vmSize": "Standard_D2_v2",
      "count": 1
    },
    "servicePrincipalProfile": {
      "servicePrincipalClientID": "Service Principal App ID",
      "servicePrincipalClientSecret": "Service Principal Password"
    },
    "networkProfile": {
      "podCidr": "10.0.0.0/24",
      "serviceCIDR": "10.0.0.0/16",
      "vnetSubnetId": "",
      "kubeDnsServiceIp": ""
    },
    "certificateProfile": {
      "caCertificate": "-----BEGIN CERTIFICATE-----/n ...",
      "caPrivateKey": "-----BEGIN CERTIFICATE-----/n ...",
      "apiServerCertificate": "-----BEGIN CERTIFICATE-----/n ...",
      "apiServerPrivateKey": "-----BEGIN CERTIFICATE-----/n ...",
      "clientCertificate": "-----BEGIN CERTIFICATE-----/n ...",
      "clientPrivateKey": "-----BEGIN CERTIFICATE-----/n ...",
      "kubeConfigCertificate": "-----BEGIN CERTIFICATE-----/n ...",
      "kubeConfigPrivateKey": "-----BEGIN CERTIFICATE-----/n ..."
    }
  }
}
```

After the ARM template has been created it will be written to disk in the `_output/$clustername` directory relative to the path you ran the command from.

The ARM template can now be deployed using any avenue the user chooses.