# Microsoft Azure Container Service Engine - Kubernetes Walkthrough

* [Kubernetes Windows Walkthrough](kubernetes.windows.md) - shows how to create a Kubernetes cluster on Windows.
* [Kubernetes with GPU support Walkthrough](kubernetes.gpu.md) - shows how to create a Kubernetes cluster with GPU support.

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

## Features

|Feature|Status|API Version|Example|Description|
|---|---|---|---|---|
|Managed Identity|Alpha|`vlabs`|[kubernetes-msi.json](../examples/managed-identity/kubernetes-msi.json)|[Description](#feat-kubernetes-msi)|

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

## Troubleshooting

### Scaling up or down

Scaling your cluster up or down requires different parameters and template than the create. More details here [Scale up](../examples/scale-up/README.md)

If your cluster is not reachable, you can run the following command to check for common failures.

### Misconfigured Service Principal

If your Service Principal is misconfigured, none of the Kubernetes components will come up in a healthy manner.
You can check to see if this the problem:

```shell
ssh -i ~/.ssh/id_rsa USER@MASTERFQDN sudo journalctl -u kubelet | grep --text autorest
```

If you see output that looks like the following, then you have **not** configured the Service Principal correctly.
You may need to check to ensure the credentials were provided accurately, and that the configured Service Principal has
read and **write** permissions to the target Subscription.

`Nov 10 16:35:22 k8s-master-43D6F832-0 docker[3177]: E1110 16:35:22.840688    3201 kubelet_node_status.go:69] Unable to construct api.Node object for kubelet: failed to get external ID from cloud provider: autorest#WithErrorUnlessStatusCode: POST https://login.microsoftonline.com/72f988bf-86f1-41af-91ab-2d7cd011db47/oauth2/token?api-version=1.0 failed with 400 Bad Request: StatusCode=400`

3. [Link](serviceprincipal.md) to documentation on how to create/configure a service principal for an ACS-Engine Kubernetes cluster.

### Managed Disks

While [Managed disks](../examples/disks-managed/README.md) are supported for the node OS disks, they are currently not supported for persistent volumes. See https://github.com/kubernetes/kubernetes/pull/46360 for details.

## Known issues and mitigations

### Node "NotReady" due to lost TCP connection

Nodes might appear in the "NotReady" state for approx. 15 minutes if master stops receiving updates from agents.
This is a known upstream kubernetes [issue #41916](https://github.com/kubernetes/kubernetes/issues/41916#issuecomment-312428731). This fixing PR is currently under review.

ACS-Engine partially mitigates this issue on Linux by detecting dead TCP connections more quickly via **net.ipv4.tcp_retries2=8**.

## Learning More

Here are recommended links to learn more about Kubernetes:

1. [Kubernetes Bootcamp](https://kubernetesbootcamp.github.io/kubernetes-bootcamp/index.html) - shows you how to deploy, scale, update, and debug containerized applications.
2. [Kubernetes Userguide](http://kubernetes.io/docs/user-guide/) - provides information on running programs in an existing Kubernetes cluster.
3. [Kubernetes Examples](https://github.com/kubernetes/kubernetes/tree/master/examples) - provides a number of examples on how to run real applications with Kubernetes.
