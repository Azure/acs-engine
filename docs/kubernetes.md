# Microsoft Azure Container Service Engine - Kubernetes Walkthrough

* [Kubernetes Windows Walkthrough](kubernetes/windows.md) - shows how to create a Kubernetes cluster on Windows.
* [Kubernetes with GPU support Walkthrough](kubernetes/gpu.md) - shows how to create a Kubernetes cluster with GPU support.

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
