# Microsoft Azure Container Service Engine - Network Policy

## Overview

> The Azure VNET network policy for Linux containers is currently available as **public preview** for Kubernetes. During public preview, Kubernetes pods do not have outbound Internet access from their private IP addresses. See [Issue #561](https://github.com/Azure/acs-engine/issues/561). Exposed Kubernetes deployments behind a VIP work as expected. Support for Windows containers and other orchestrators is coming soon.

These cluster definition examples demonstrate how to create customized Docker enabled clusters with a network policy provider.

1. **kubernetes-azure.json** - deploying and using [Kubernetes](../../docs/kubernetes.md) with Azure VNET network policy.
2. **kubernetes-azure-hybrid.json** - deploying and using [Kubernetes](../../docs/kubernetes.md) with both Linux and Windows agent pools and Azure VNET network policy.
3. **kubernetes-calico.json** - deploying and using [Kubernetes](../../docs/kubernetes.md) with Calico network policy.
