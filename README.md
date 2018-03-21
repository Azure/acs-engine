# Microsoft Azure Container Service Engine - Builds Docker Enabled Clusters
[![Coverage Status](https://coveralls.io/repos/github/Azure/acs-engine/badge.svg?branch=master)](https://coveralls.io/github/Azure/acs-engine?branch=master)
[![CircleCI](https://circleci.com/gh/Azure/acs-engine/tree/master.svg?style=svg)](https://circleci.com/gh/Azure/acs-engine/tree/master)

## Overview

The Azure Container Service Engine (`acs-engine`) generates ARM (Azure Resource Manager) templates for Docker enabled clusters on Microsoft Azure with your choice of DC/OS, Kubernetes, Swarm Mode, or Swarm orchestrators. The input to the tool is a cluster definition. The cluster definition is very similar to (in many cases the same as) the ARM template syntax used to deploy a Microsoft Azure Container Service cluster.

The cluster definition file enables the following customizations to your Docker enabled cluster:
* choice of DC/OS, Kubernetes, Swarm Mode, or Swarm orchestrators
* multiple agent pools where each agent pool can specify:
 * standard or premium VM Sizes,
 * node count,
 * Virtual Machine ScaleSets or Availability Sets,
 * Storage Account Disks or Managed Disks (under private preview)
* Docker cluster sizes of 1200
* Custom VNET

## User guides

* [ACS Engine](docs/acsengine.md) - shows you how to build and use the ACS engine to generate custom Docker enabled container clusters
* [Cluster Definition](docs/clusterdefinition.md) - describes the components of the cluster definition file
* [DC/OS Walkthrough](docs/dcos.md) - shows how to create a DC/OS enabled Docker cluster on Azure
* [Kubernetes Walkthrough](docs/kubernetes.md) - shows how to create a Kubernetes enabled Docker cluster on Azure
* [Swarm Walkthrough](docs/swarm.md) - shows how to create a Swarm enabled Docker cluster on Azure
* [Swarm Mode Walkthrough](docs/swarmmode.md) - shows how to create a Swarm Mode cluster on Azure
* [Custom VNET](examples/vnet) - shows how to use a custom VNET
* [Attached Disks](examples/disks-storageaccount) - shows how to attach up to 4 disks per node
* [Managed Disks](examples/disks-managed) - shows how to use managed disks
* [Large Clusters](examples/largeclusters) - shows how to create cluster sizes of up to 1200 nodes

## Contributing

Please follow these instructions before submitting a PR:

1. Execute `make ci` to run the checkin validation tests.

2. Manually test deployments if you are making modifications to the templates.
   For example, if you have to change the expected resulting templates then you
   should deploy the relevant example cluster definitions to ensure you're not
   introducing any sort of regression.

## Usage

### Generate Templates

Usage is best demonstrated with an example:

```shell
$ vim examples/classic/kubernetes.classic.json

# insert your preferred, unique DNS prefix
# insert your SSH public key

$ ./acs-engine generate examples/classic/kubernetes.classic.json
```

This produces a new directory inside `_output/` that contains an ARM template
for deploying Kubernetes into Azure. (In the case of Kubernetes, some additional
needed assets are generated and placed in the output directory.)

## Code of conduct

This project has adopted the [Microsoft Open Source Code of Conduct](https://opensource.microsoft.com/codeofconduct/). For more information see the [Code of Conduct FAQ](https://opensource.microsoft.com/codeofconduct/faq) or contact [opencode@microsoft.com](mailto:opencode@microsoft.com) with any additional questions or comments.
