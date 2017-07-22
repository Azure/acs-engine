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
