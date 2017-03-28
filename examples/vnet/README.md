# Microsoft Azure Container Service Engine - Custom VNET

## Overview

These examples show you how to build a customized Docker enabled cluster on Microsoft Azure where you can provide your own VNET.

By providing your own VNET, you can change the VNET address space from the default of 10.0.0.0/8 to a smaller subnet, like
10.1.0.0/16. This may be useful if you have existing VNETs that you would like to peer with, or plan to have them in the future.
Kubernetes will still use 10.0.0.0/8 internally within the cluster, but all the node ports will be within the modified VNET
address space.

Modifying the address space to a completely different network is not yet supported. See
[issue #180](https://github.com/Azure/acs-engine/issues/180).

**Note**: Kubernetes must have the custom VNET deployed in the same resource group as Kubernetes.

To try: 

1. first deploy a custom vnet.  An example of an arm template that does this is under directory vnetarmtemplate.
2. next configure the example templates and deploy according to the examples:
 1. **dcos.json** - deploying and using [DC/OS](../../docs/dcos.md)
 2. **kubernetes.json** - deploying and using [Kubernetes](../../docs/kubernetes.md)
 3. **swarm.json** - deploying and using [Swarm](../../docs/swarm.md)
 4. **swarmmodevnet.json** - deploying and using [Swarm Mode](../../docs/swarmmode.md)

