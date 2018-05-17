# Microsoft Azure Container Service Engine - Builds Docker Enabled Clusters
[![Coverage Status](https://codecov.io/gh/Azure/acs-engine/branch/master/graph/badge.svg)](https://codecov.io/gh/Azure/acs-engine)
[![CircleCI](https://circleci.com/gh/Azure/acs-engine/tree/master.svg?style=svg)](https://circleci.com/gh/Azure/acs-engine/tree/master)

## Overview

The Azure Container Service Engine (`acs-engine`) generates ARM (Azure Resource Manager) templates for Docker enabled clusters on Microsoft Azure with your choice of DC/OS, Kubernetes, OpenShift, Swarm Mode, or Swarm orchestrators. The input to the tool is a cluster definition. The cluster definition (or apimodel) is very similar to (in many cases the same as) the ARM template syntax used to deploy a Microsoft Azure Container Service cluster.

The cluster definition file enables you to customize your Docker enabled cluster in many ways including:
* choice of DC/OS, Kubernetes, OpenShift, Swarm Mode, or Swarm orchestrators
* multiple agent pools where each agent pool can specify:
   * standard or premium VM Sizes, including GPU optimized VM sizes
   * node count
   * Virtual Machine ScaleSets or Availability Sets
   * Storage Account Disks or Managed Disks
   * OS and distro
* Custom VNET
* Extensions

## User guides

* [ACS Engine](docs/acsengine.md) - shows you how to build and use the ACS engine to generate custom Docker enabled container clusters
* [Cluster Definition](docs/clusterdefinition.md) - describes the components of the cluster definition file
* [DC/OS Walkthrough](docs/dcos.md) - shows how to create a DC/OS enabled Docker cluster on Azure
* [Kubernetes Walkthrough](docs/kubernetes.md) - shows how to create a Kubernetes enabled Docker cluster on Azure
* [OpenShift Walkthrough](docs/openshift.md) - shows how to create an OpenShift enabled Docker cluster on Azure
* [Swarm Walkthrough](docs/swarm.md) - shows how to create a Swarm enabled Docker cluster on Azure
* [Swarm Mode Walkthrough](docs/swarmmode.md) - shows how to create a Swarm Mode cluster on Azure
* [Custom VNET](examples/vnet) - shows how to use a custom VNET
* [Attached Disks](examples/disks-storageaccount) - shows how to attach up to 4 disks per node
* [Managed Disks](examples/disks-managed) - shows how to use managed disks
* [Large Clusters](examples/largeclusters) - shows how to create cluster sizes of up to 1200 nodes

## Contributing

Follow the [developers guide](docs/developers.md) to set up your environment.

To build acs-engine, run `make build`. If you are developing with a working [Docker environment](https://docs.docker.com/engine), you can also run `make dev` first to start a Docker container and run `make build` inside the container.

Please follow these instructions before submitting a PR:

1. Execute `make test` to run unit tests.

2. Manually test deployments if you are making modifications to the templates.
   For example, if you have to change the expected resulting templates then you
   should deploy the relevant example cluster definitions to ensure that you are not introducing any regressions.

3. Make sure that your changes are properly documented and include relevant unit tests.

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
