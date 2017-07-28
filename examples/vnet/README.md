# Microsoft Azure Container Service Engine - Custom VNET

## Overview

These examples show you how to build a customized Docker enabled cluster on Microsoft Azure where you can provide your own VNET.

To try: 

1. first deploy a custom vnet.  An example of an arm template that does this is under directory vnetarmtemplate.
2. next configure the example templates and deploy according to the examples:
 1. **dcos.json** - deploying and using [DC/OS](../../docs/dcos.md)
 2. **kubernetes.json** - deploying and using [Kubernetes](../../docs/kubernetes.md)
 3. **swarm.json** - deploying and using [Swarm](../../docs/swarm.md)
 4. **swarmmodevnet.json** - deploying and using [Swarm Mode](../../docs/swarmmode.md)

