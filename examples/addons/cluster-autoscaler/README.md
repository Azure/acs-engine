# Cluster Autoscaler (VMSS) Add-on

Cluster Autoscaler is a tool that automatically adjusts the size of the Kubernetes cluster when:

* there are pods that failed to run in the cluster due to insufficient resources.
* some nodes in the cluster are so underutilized, for an extended period of time, that they can be deleted and their pods will be easily placed on some other, existing nodes.

This is the Kubernetes Cluster Autoscaler add-on for Virtual Machine Scale Sets. Add this add-on to your json file as shown below to automatically enable cluster autoscaler in your new Kubernetes cluster.

To use this add-on, make sure your cluster's Kubernetes version is 1.10 or above, and agent pool `availabilityProfile` is set to `VirtualMachineScaleSets`. This will automatically enable first agent pool to autoscale from 1 to 5 nodes by default. You can override these settings in `config` section of the `cluster-autoscaler` add-on.

```
{
    "apiVersion": "vlabs",
    "properties": {
      "orchestratorProfile": {
        "orchestratorType": "Kubernetes",
        "orchestratorRelease": "1.10",
        "kubernetesConfig": {
          "addons": [
            {
              "name": "cluster-autoscaler",
              "enabled" : true,
              "config": {
                "minNodes": "1",
                "maxNodes": "5"
              }
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
          "count": 1,
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

You should see cluster autoscaler as running after running:

```
$ kubectl get pods -n kube-system
```

Follow the README at https://github.com/kubernetes/autoscaler/tree/master/cluster-autoscaler for examples.

# Configuration

| Name           | Required | Description                       | Default Value                                              |
| -------------- | -------- | --------------------------------- | ---------------------------------------------------------- |
| minNodes       | no       | minimum node count                |                                                            |
| maxNodes       | no       | maximum node count                |                                                            |
| name           | no       | container name                    | "cluster-autoscaler"                                       |
| image          | no       | image                             | "gcrio.azureedge.net/google-containers/cluster-autoscaler" |
| cpuRequests    | no       | cpu requests for the container    | "100m"                                                     |
| memoryRequests | no       | memory requests for the container | "300Mi"                                                    |
| cpuLimits      | no       | cpu limits for the container      | "100m"                                                     |
| memoryLimits   | no       | memory limits for the container   | "300Mi"                                                    |

# Supported Orchestrators

Kubernetes
