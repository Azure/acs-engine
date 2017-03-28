# Microsoft Azure Container Service Engine - Custom VNET

## Overview

These examples show you how to build a customized Docker enabled cluster on Microsoft Azure where you can provide your own VNET.

By providing your own VNET, you can change the VNET address space from the default of 10.0.0.0/8 to a smaller subnet, like
10.1.0.0/16. This may be useful if you have existing VNETs that you would like to peer with, or plan to have them in the future.
For Kubernetes, a change of the address space only affects the IP addresses assigned to the master and agent nodes. The service
network of 10.0.0.0/16 is still implemented by the kube-proxy in each host, and the pod network stays as 10.244.0.0/16. This
means that access to services from another VNET will have to be done via node port through the host addresses, or via an
internal load balancer (which Kubernetes can create automatically with the 
[appropriate service annotation](https://kubernetes.io/docs/concepts/services-networking/service/#internal-load-balancer)).

Future versions will make the service and pod networks configurable. See
[pull #546](https://github.com/Azure/acs-engine/pull/546).

To try: 

1. first deploy a custom vnet.  An example of an arm template that does this is under directory vnetarmtemplate.
2. next configure the example templates and deploy according to the examples:
 1. **dcos.json** - deploying and using [DC/OS](../../docs/dcos.md)
 2. **kubernetes.json** - deploying and using [Kubernetes](../../docs/kubernetes.md)
 3. **swarm.json** - deploying and using [Swarm](../../docs/swarm.md)
 4. **swarmmodevnet.json** - deploying and using [Swarm Mode](../../docs/swarmmode.md)

