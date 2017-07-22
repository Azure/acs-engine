## Features

|Feature|Status|API Version|Example|Description|
|---|---|---|---|---|
|Managed Identity|Alpha|`vlabs`|[kubernetes-msi.json](../../examples/managed-identity/kubernetes-msi.json)|[Description](#feat-kubernetes-msi)|

<a name="feat-kubernetes-msi"></a>
### Managed Identity

Enabling Managed Identity configures acs-engine to include and use MSI identities for all interactions with the Azure Resource Manager (ARM) API.

Instead of using a static servic principal written to `/etc/kubernetes/azure.json`, Kubernetes will use a dynamic, time-limited token fetched from the MSI extension running on master and agent nodes. This support is currently alpha and requires Kubernetes v1.7.2 or newer.

Enable Managed Identity by adding `useManagedIdentity` in `kubernetesConfig`.

```json
"kubernetesConfig": {
  "useManagedIdentity": true,
  "customHyperkubeImage": "docker.io/colemickens/hyperkube-amd64:3b15e8a446fa09d68a2056e2a5e650c90ae849ed"
}
```

## Optional: Enable network policy enforcement using Calico

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
* [Calico Kubernetes](http://docs.projectcalico.org/v2.0/getting-started/kubernetes/)

### Managed Disks

While [Managed disks](../examples/disks-managed/README.md) are supported for the node OS disks, they are currently not supported for persistent volumes. See https://github.com/kubernetes/kubernetes/pull/46360 for details.
