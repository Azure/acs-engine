# Features

|Feature|Status|API Version|Example|Description|
|---|---|---|---|---|
|Managed Disks|Beta|`vlabs`|[kubernetes-vmas.json](../../examples/disks-managed/kubernetes-vmss.json)|[Description](#feat-managed-disks)|
|Calico Network Policy|Alpha|`vlabs`|[kubernetes-calico.json](../../examples/networkpolicy/kubernetes-calico.json)|[Description](#feat-calico)|
|Custom VNET|Beta|`vlabs`|[kubernetesvnet-azure-cni.json](../../examples/vnet/kubernetesvnet-azure-cni.json)|[Description](#feat-custom-vnet)|

<a name="feat-kubernetes-msi"></a>

## Managed Identity

Enabling Managed Identity configures acs-engine to include and use MSI identities for all interactions with the Azure Resource Manager (ARM) API.

Instead of using a static servic principal written to `/etc/kubernetes/azure.json`, Kubernetes will use a dynamic, time-limited token fetched from the MSI extension running on master and agent nodes. This support is currently alpha and requires Kubernetes v1.9.1 or newer.

Enable Managed Identity by adding `useManagedIdentity` in `kubernetesConfig`.

```json
"kubernetesConfig": {
  "useManagedIdentity": true,
  "customHyperkubeImage": "docker.io/colemickens/hyperkube-amd64:3b15e8a446fa09d68a2056e2a5e650c90ae849ed"
}
```

<a name="feat-managed-disks"></a>

## Optional: Disable Kubernetes Role-Based Access Control (RBAC)

By default, the cluster will be provisioned with [Role-Based Access Control](https://kubernetes.io/docs/admin/authorization/rbac/) enabled. Disable RBAC by adding `enableRbac` in `kubernetesConfig` in the api model:

```console
      "kubernetesConfig": {
        "enableRbac": false
      }
```

See [cluster definition](https://github.com/Azure/acs-engine/blob/master/docs/clusterdefinition.md#kubernetesconfig) for further detail.

## Managed Disks

[Managed disks](../examples/disks-managed/README.md) are supported for both node OS disks and Kubernetes persistent volumes.

Related [upstream PR](https://github.com/kubernetes/kubernetes/pull/46360) for details.

### Using Kubernetes Persistent Volumes

By default, each ACS-Engine cluster is bootstrapped with several StorageClass resources. This bootstrapping is handled by the addon-manager pod that creates resources defined under /etc/kubernetes/addons directory on master VMs.

#### Non-managed Disks

The default storage class has been set via the Kubernetes admission controller `DefaultStorageClass`.

The default storage class will be used if persistent volume resources don't specify a storage class as part of the resource definition.

The default storage class uses non-managed blob storage and will provision the blob within an existing storage account present in the resource group or provision a new storage account.

Non-managed persistent volume types are available on all VM sizes.

#### Managed Disks

As part of cluster bootstrapping, two storage classes will be created to provide access to create Kubernetes persistent volumes using Azure managed disks.

Nodes will be labelled as follows if they support managed disks:

```
storageprofile=managed
storagetier=<Standard_LRS|Premium_LRS>
```

They are managed-premium and managed-standard and map to Standard_LRS and Premium_LRS managed disk types respectively.

In order to use these storage classes the following conditions must be met.

* The cluster must be running Kubernetes release 1.7 or greater. Refer to this [example](../examples/kubernetes-releases/kubernetes1.7.json) for how to provision a Kubernetes cluster of a specific version.
* The node must support managed disks. See this [example](../examples/disks-managed/kubernetes-vmas.json) to provision nodes with managed disks. You can also confirm if a node has managed disks using kubectl.

```console
kubectl get nodes -l storageprofile=managed
NAME                    STATUS    AGE       VERSION
k8s-agent1-23731866-0   Ready     24m       v1.7.2
```

* The VM size must support the type of managed disk type requested. For example, Premium VM sizes with managed OS disks support both managed-standard and managed-premium storage classes whereas Standard VM sizes with managed OS disks only support managed-standard storage class.

* If you have mixed node cluster (both non-managed and managed disk types). You must use [affinity or nodeSelectors](https://kubernetes.io/docs/concepts/configuration/assign-pod-node/) on your resource definitions in order to ensure that workloads are scheduled to VMs that support the underlying disk requirements.

For example
```
spec:
  affinity:
    nodeAffinity:
      requiredDuringSchedulingIgnoredDuringExecution:
        nodeSelectorTerms:
        - matchExpressions:
          - key: storageprofile
            operator: In
            values:
            - managed
```

## Using Azure integrated networking (CNI)

Kubernetes clusters are configured by default to use the [Azure CNI plugin](https://github.com/Azure/azure-container-networking) which provides an Azure native networking experience. Pods will receive IP addresses directly from the vnet subnet on which they're hosted. If the api model doesn't specify explicitly, acs-engine will automatically provide the following `networkPolicy` configuration in `kubernetesConfig`:

```
      "kubernetesConfig": {
        "networkPolicy": "azure"
      }
```

### Additional Azure integrated networking configuration

In addition you can modify the following settings to change the networking behavior when using Azure integrated networking:

IP addresses are pre-allocated in the subnet. Using ipAddressCount you can specify how many you would like to pre-allocate. This number needs to account for number of pods you would like to run on that subnet.

```
    "masterProfile": {
      "ipAddressCount": 200
    },
```

Currently, the IP addresses that are pre-allocated aren't allowed by the default natter for Internet bound traffic. In order to work around this limitation we allow the user to specify the vnetCidr (eg. 10.0.0.0/8) to be EXCLUDED from the default masquerade rule that is applied. The result is that traffic destined for anything within that block will NOT be natted on the outbound VM interface. This field has been called vnetCidr but may be wider than the vnet cidr block if you would like POD IPs to be routable across vnets using vnet-peering or express-route.
```
    "masterProfile": {
      "vnetCidr": "10.0.0.0/8",
    },
```

When using Azure integrated networking the maxPods setting will be set to 30 by default. This number can be changed keeping in mind that there is a limit of 4,000 IPs per vnet.

```
      "kubernetesConfig": {
        "maxPods": 50
      }
```

<a name="feat-calico"></a>

## Network Policy Enforcement with Calico

Using the default configuration, Kubernetes allows communication between all
Pods within a cluster. To ensure that Pods can only be accessed by authorized
Pods, a policy enforcement is needed. To enable policy enforcement using Calico refer to the [cluster definition](https://github.com/Azure/acs-engine/blob/master/docs/clusterdefinition.md#kubernetesconfig) document under networkPolicy. There is also a reference cluster definition available [here](https://github.com/Azure/acs-engine/blob/master/examples/networkpolicy/kubernetes-calico.json).

This will deploy a Calico node controller to every instance of the cluster
using a Kubernetes DaemonSet. After a successful deployment you should be able
to see these Pods running in your cluster:

```
kubectl get pods --namespace kube-system -l k8s-app=calico-node -o wide
NAME                READY     STATUS    RESTARTS   AGE       IP             NODE
calico-node-034zh   2/2       Running   0          2h        10.240.255.5   k8s-master-30179930-0
calico-node-qmr7n   2/2       Running   0          2h        10.240.0.4     k8s-agentpool1-30179930-1
calico-node-z3p02   2/2       Running   0          2h        10.240.0.5     k8s-agentpool1-30179930-0
```

Per default Calico still allows all communication within the cluster. Using Kubernetes' NetworkPolicy API, you can define stricter policies. Good resources to get information about that are:

* [NetworkPolicy User Guide](https://kubernetes.io/docs/user-guide/networkpolicies/)
* [NetworkPolicy Example Walkthrough](https://kubernetes.io/docs/getting-started-guides/network-policy/walkthrough/)
* [Calico Kubernetes](https://github.com/Azure/acs-engine/blob/master/examples/networkpolicy)

<a name="feat-custom-vnet"></a>

## Custom VNET

ACS Engine supports deploying into an existing VNET. Operators must specify the ARM path/id of Subnets for the `masterProfile` and  any `agentPoolProfiles`, as well as the first IP address to use for IP static IP allocation in `firstConsecutiveStaticIP`. Additionally, to prevent source address NAT'ing within the VNET, we assign to the `vnetCidr` property in `masterProfile` the CIDR block that represents the usable address space in the existing VNET.

Depending upon the size of the VNET address space, during deployment, it is possible to experience IP address assignment collision between the required Kubernetes static IPs (one each per master and one for the API server load balancer, if more than one masters) and Azure CNI-assigned dynamic IPs (one for each NIC on the agent nodes). In practice, the larger the VNET the less likely this is to happen; some detail, and then a guideline.

First, the detail:

* Azure CNI assigns dynamic IP addresses from the "beginning" of the subnet IP address space (specifically, it looks for available addresses starting at ".4" ["10.0.0.4" in a "10.0.0.0/24" network])
* acs-engine will require a range of up to 16 unused IP addresses in multi-master scenarios (1 per master for up to 5 masters, and then the next 10 IP addresses immediately following the "last" master for headroom reservation, and finally 1 more for the load balancer immediately adjacent to the afore-described _n_ masters+10 sequence) to successfully scaffold the network stack for your cluster

A guideline that will remove the danger of IP address allocation collision during deployment:

* If possible, assign to the `firstConsecutiveStaticIP` configuration property an IP address that is near the "end" of the available IP address space in the desired  subnet.
  * For example, if the desired subnet is a `/24`, choose the "239" address in that network space

In larger subnets (e.g., `/16`) it's not as practically useful to push static IP assignment to the very "end" of large subnet, but as long as it's not in the "first" `/24` (for example) your deployment will be resilient to this edge case behavior.

Before provisioning, modify the `masterProfile` and `agentPoolProfiles` to match the above requirements, with the below being a representative example:

```json
"masterProfile": {
  ...
  "vnetSubnetId": "/subscriptions/SUB_ID/resourceGroups/RG_NAME/providers/Microsoft.Network/virtualNetworks/VNET_NAME/subnets/MASTER_SUBNET_NAME",
  "firstConsecutiveStaticIP": "10.239.255.239",
  "vnetCidr": "10.239.0.0/16",
  ...
},
...
"agentPoolProfiles": [
  {
    ...
    "name": "agentpri",
    "vnetSubnetId": "/subscriptions/SUB_ID/resourceGroups/RG_NAME/providers/Microsoft.Network/virtualNetworks/VNET_NAME/subnets/AGENT_SUBNET_NAME",
    ...
  },
```

### Kubenet Networking Custom VNET

If you're not using Azure CNI (e.g., `"networkPolicy": "none"` in the `kubernetesConfig` api model configuration object): After a custom VNET-configured cluster finishes provisioning, fetch the id of the Route Table resource from `Microsoft.Network` provider in your new cluster's Resource Group.

The route table resource id is of the format: `/subscriptions/SUBSCRIPTIONID/resourceGroups/RESOURCEGROUPNAME/providers/Microsoft.Network/routeTables/ROUTETABLENAME`

Existing subnets will need to use the Kubernetes-based Route Table so that machines can route to Kubernetes-based workloads.

Update properties of all subnets in the existing VNET he route table resource by appending the following to subnet properties:

```json
"routeTable": {
        "id": "/subscriptions/<SubscriptionId>/resourceGroups/<ResourceGroupName>/providers/Microsoft.Network/routeTables/k8s-master-<SOMEID>-routetable>"
      }
```

E.g.:
```json
"subnets": [
    {
      "name": "subnetname",
      "id": "/subscriptions/<SubscriptionId>/resourceGroups/<ResourceGroupName>/providers/Microsoft.Network/virtualNetworks/<VirtualNetworkName>/subnets/<SubnetName>",
      "properties": {
        "provisioningState": "Succeeded",
        "addressPrefix": "10.240.0.0/16",
        "routeTable": {
          "id": "/subscriptions/<SubscriptionId>/resourceGroups/<ResourceGroupName>/providers/Microsoft.Network/routeTables/k8s-master-<SOMEID>-routetable"
        }
      ...
      }
      ...
    }
]
```
