# Microsoft Azure Container Service Engine - Custom VNET

## Overview

These examples show you how to build a customized Docker enabled cluster on Microsoft Azure where you can provide your own VNET.

**Note**: Kubernetes must have the custom VNET deployed in the same resource group as Kubernetes.

To try:

1. first deploy a custom vnet.  An example of an arm template that does this is under directory vnetarmtemplate.
2. next configure the example templates and deploy according to the examples:
 1. **dcos.json** - deploying and using [DC/OS](../../docs/dcos.md)
 1. **kubernetes.json** - deploying and using [Kubernetes](../../docs/kubernetes.md)
 1. **swarmmodevnet.json** - deploying and using [Swarm Mode](../../docs/swarmmode.md)

