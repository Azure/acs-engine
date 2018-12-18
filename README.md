# Pardon our Dust!

This codebase has been deprecated in favor of aks-engine, the natural evolution from acs-engine:

https://github.com/Azure/aks-engine

All future development and maintenance will occur there as an outcome of this deprecation. We're sorry for any inconvenience!

We've moved the Kubernetes code over 100% as-is (with the exception of the boilerplate renaming overhead that accompanies such a move); we're confident this housekeeping manouver will more effectively track the close affinity between the AKS managed service and the "build and manage your own configurable Kubernetes" stories that folks use this tool for.

See you at https://github.com/Azure/aks-engine!

The historical documentation remains below.

# Microsoft Azure Container Service Engine - Builds Docker Enabled Clusters

[![Coverage Status](https://codecov.io/gh/Azure/acs-engine/branch/master/graph/badge.svg)](https://codecov.io/gh/Azure/acs-engine)
[![CircleCI](https://circleci.com/gh/Azure/acs-engine/tree/master.svg?style=svg)](https://circleci.com/gh/Azure/acs-engine/tree/master)
[![GoDoc](https://godoc.org/github.com/Azure/acs-engine?status.svg)](https://godoc.org/github.com/Azure/acs-engine)

## Overview

The Azure Container Service Engine (`acs-engine`) generates ARM (Azure Resource Manager) templates for Docker enabled clusters on Microsoft Azure with your choice of DC/OS, Kubernetes, OpenShift, Swarm Mode, or Swarm orchestrators. The input to the tool is a cluster definition. The cluster definition (or apimodel) is very similar to (in many cases the same as) the ARM template syntax used to deploy a Microsoft Azure Container Service cluster.

The cluster definition file enables you to customize your Docker enabled cluster in many ways including:

* Choice of DC/OS, Kubernetes, OpenShift, Swarm Mode, or Swarm orchestrators
* Multiple agent pools where each agent pool can specify:
  * Standard or premium VM Sizes, including GPU optimized VM sizes
  * Node count
  * Virtual Machine ScaleSets or Availability Sets
  * Storage Account Disks or Managed Disks
  * OS and distro
* Custom VNET
* Extensions

More info, including a thorough walkthrough is [here](docs/acsengine.md).

## User guides

These guides show how to create your first deployment for each orchestrator:

* [DC/OS Walkthrough](docs/dcos.md) - shows how to create a DC/OS cluster on Azure
* [Kubernetes Walkthrough](docs/kubernetes.md) - shows how to create a Linux or Windows Kubernetes cluster on Azure
* [OpenShift Walkthrough](docs/openshift.md) - shows how to create an OpenShift cluster on Azure
* [Swarm Mode Walkthrough](docs/swarmmode.md) - shows how to create a [Docker Swarm Mode](https://docs.docker.com/engine/swarm/) cluster on Azure
* [Standalone Swarm Walkthrough](docs/swarm.md) - shows how to create a [Docker Standalone Swarm](https://docs.docker.com/swarm/) cluster on Azure

These guides cover more advanced features to try out after you have built your first cluster:

* [Cluster Definition](docs/clusterdefinition.md) - describes the components of the cluster definition file
* [Custom VNET](examples/vnet) - shows how to use a custom VNET
* [Attached Disks](examples/disks-storageaccount) - shows how to attach up to 4 disks per node
* [Managed Disks](examples/disks-managed) - shows how to use managed disks
* [Large Clusters](examples/largeclusters) - shows how to create cluster sizes of up to 1200 nodes

## Usage

### Generate Templates

Usage is best demonstrated with an example:

```sh
$ vim examples/kubernetes.json

# insert your preferred, unique DNS prefix
# insert your SSH public key

$ ./acs-engine generate examples/kubernetes.json
```

This produces a new directory inside `_output/` that contains an ARM template for deploying Kubernetes into Azure. (In the case of Kubernetes, some additional needed assets are generated and placed in the output directory.)

## Code of conduct

This project has adopted the [Microsoft Open Source Code of Conduct](https://opensource.microsoft.com/codeofconduct/). For more information see the [Code of Conduct FAQ](https://opensource.microsoft.com/codeofconduct/faq) or contact [opencode@microsoft.com](mailto:opencode@microsoft.com) with any additional questions or comments.
