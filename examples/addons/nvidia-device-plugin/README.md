# NVIDIA Device Plugin

This is the NVIDIA Device Plugin add-on for Kubernetes (v1.8+). This add-on will be automatically be enabled if you are using a Kubernetes cluster (v1.8+) with an N-series agent pool (which contains an NVIDIA GPU). You can use this add-on to your json file as shown below to enable or disable NVIDIA Device Plugin explicitly.

```
{
  "apiVersion": "vlabs",
  "properties": {
    "orchestratorProfile": {
      "orchestratorType": "Kubernetes",
      "orchestratorRelease": "1.9",
      "kubernetesConfig": {
        "addons": [
          {
            "name": "nvidia-device-plugin",
            "enabled": true
          }
        ]
      }
    },
    "masterProfile": {
      "count": 1,
      "dnsPrefix": "",
      "vmSize": "Standard_DS2_v2"
    },
    "agentPoolProfiles": [
      {
        "name": "agentpool",
        "count": 3,
        "vmSize": "Standard_NC6",
        "availabilityProfile": "AvailabilitySet"
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

You can validate that the add-on is running as expected with the following commands:

Make sure to create resource group:

```
az group create \
    --name "[resource group name]" \
    --location "[location]"
```

You should see NVIDIA Device Plugin daemonset as running after running:

```
$ kubectl get pods -n kube-system
```

Follow the README at https://github.com/NVIDIA/k8s-device-plugin for examples.

# Supported Orchestrators

Kubernetes
