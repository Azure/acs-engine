# ACI Connector Add-on


This is the ACI Connector add-on.  Add this add-on to your json file as shown below to automatically enable ACI Connector in your new Kubernetes cluster.

```
{
    "apiVersion": "vlabs",
    "properties": {
      "orchestratorProfile": {
        "orchestratorType": "Kubernetes",
        "kubernetesConfig": {
          "addons": [
            {
              "name": "aci-connector",
              "enabled" : true,
              "config": {
                  "clientId": "",
                  "clientKey": "",
                  "tenantId": "",
                  "subscriptionId": "",
                  "resourceGroup": "",
                  "region": "",
                  "nodeName": "",
                  "os": "",
                  "taint": ""
              },
              "containers": [
                {
                  "name": "aci-connector",
                  "cpuRequests": "50m",
                  "memoryRequests": "150Mi",
                  "cpuLimits": "50m",
                  "memoryLimits": "150Mi"
                }
              ]
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

You should see ACI Connector as running after running:
```
$ kubectl get pods -n kube-system
```

You should see ACI Connector node after running:
```
$ kubectl get nodes
```

Follow the README at https://github.com/virtual-kubelet/virtual-kubelet for examples.

# Configuration
|Name|Required|Description|Default Value|
|---|---|---|---|
|clientId|yes|your client id||
|clientKey|yes|your client key||
|tenantId|yes|your tenant id||
|resourceGroup|yes|your resource group||
|region|no|Azure region|"westus"|
|nodeName|no|node name|"aci-connector"|
|os|no|operating system (Linux/Windows)|"Linux"|
|taint|no|apply taint to node, making scheduling explicit|"azure.com/aci"|
|name|no|container name|"aci-connector"|
|image|no|image|"microsoft/virtual-kubelet:latest"|
|cpuRequests|no|cpu requests for the container|"50m"|
|memoryRequests|no|memory requests for the container|"150Mi"|
|cpuLimits|no|cpu limits for the container|"50m"|
|memoryLimits|no|memory limits for the container|"150Mi"|


# Supported Orchestrators
Kubernetes