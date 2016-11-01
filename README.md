# Microsoft Azure Container Service Engine - Builds Docker Enabled Clusters

## Overview

The Azure Container Service Engine (`acs-engine`) generates ARM (Azure Resource Manager) templates for Docker enabled clusters on Microsoft Azure with your choice of DCOS, Kubernetes, or Swarm orchestrators. The input to the tool is a cluster definition. The cluster definition is very similar to (in many cases the same as) the ARM template syntax used to deploy a Microsoft Azure Container Service cluster.

The cluster definition file enables the following customizations to your Docker enabled cluster:
* choice of DCOS, Kubernetes, or Swarm orchestrators
* multiple agent pools where each agent pool can specify:
 * standard or premium VM Sizes,
 * node count, 
 * Virtual Machine ScaleSets or Availability Sets,
 * Storage Account Disks or Managed Disks (under private preview),
 * and Linux or Microsoft Windows.
* Docker cluster sizes of 1200
* Custom VNET

## User guides

* [ACS Engine](docs/acsengine.md) - shows you how to build and use the ACS engine to generate custom Docker enabled container clusters
* [Cluster Definition](docs/clusterdefinition.md) - describes the components of the cluster definition file
* [DCOS Walkthrough](docs/dcos.md) - shows how to create a DCOS enabled Docker cluster on Azure
* [Kubernetes Walkthrough](docs/kubernetes.md) - shows how to create a Kubernetes enabled Docker cluster on Azure
* [Swarm Walkthrough](docs/swarm.md) - shows how to create a Swarm enabled Docker cluster on Azure
* [Custom VNET](examples/vnet) - shows how to use a custom VNET 
* [Attached Disks](examples/disks-storageaccount) - shows how to attach up to 4 disks per node
* [Managed Disks](examples/disks-managed) (under private preview) - shows how to use managed disks 
* [Large Clusters](examples/largeclusters) - shows how to create cluster sizes of up to 1200 nodes
* [Windows Clusters](examples/windows) - shows how to create windows or mixed Microsoft Windows and Linux Docker clusters on Microsoft Azure

## Development (Docker)

The easiest way to get started developing on `acs-engine` is to use Docker.
If you already have Docker or "Docker for {Windows,Mac}" then you can get started
without needing to install anything extra.

* Windows (PowerShell): `.\scripts\devenv.ps1`
* Linux (bash): `.\scripts\devenv.sh`

This setup mounts the `acs-engine` source directory as a volume into the Docker container.
This means that you can edit your source code normally in your favorite editor on your
machine, while still being able to compile and test inside of the Docker container (the
same environment used in our Continuous Integration system).

[Here's a quick demo video showing the dev/build/test cycle with this setup.](https://www.youtube.com/watch?v=lc6UZmqxQMs)

## Development (Native)

### Requirements
- PowerShell (Windows)
- `bash` + `make` (Linux)
- `git`
- `go` (with a properly configured GOPATH)

### Building (Linux)

```shell
make build
```

### Building (Windows, PowerShell)

```shell
cd ${env:GOPATH}/github.com/Azure/acs-engine
go get .
go build .
```


## Contributing

Please follow these instructions before submitting a PR:

1. Execute `make ci` to run the checkin validation tests.

2. Manually test deployments if you are making modifications to the templates.
   For example, if you have to change the expected resulting templates then you
   should deploy the relevant example cluster definitions to ensure you're not
   introducing any sort of regression.

## Usage (Template Generation)

Usage is best demonstrated with an example:

```shell
$ vim examples/kubernetes.classic.json

# insert your preferred, unique DNS prefix
# insert your SSH public key

$ ./acs-engine examples/kubernetes.classic.json
```

This produces a new directory inside `_output/` that contains an ARM template
for deploying Kubernetes into Azure. (In the case of Kubernetes, some additional
needed assets are generated and placed in the output directory.)

## Deployment Usage

Generated templates can be deployed using
[the Azure CLI 1.0](https://github.com/Azure/azure-xplat-cli),
[the Azure CLI 2.0](https://github.com/Azure/azure-cli) or
[Powershell](https://github.com/Azure/azure-powershell).

### Deploying with Azure CLI 1.0

```bash
$ azure login

$ azure account set --name "<SUBSCRIPTION NAME OR ID>"

$ azure group create \
    --name="<RESOURCE_GROUP_NAME>" \
    --location="<LOCATION>"

$ azure group deployment create \
    --name="<DEPLOYMENT NAME>" \
    --resource-group="<RESOURCE_GROUP_NAME>" \
    --template-file="./_output/<INSTANCE>/azuredeploy.json" \
    --parameters-file="./_output/<INSTANCE>azuredeploy.parameters.json"
```

### Deploying with Azure CLI 2.0

```bash
$ az login

$ az account set --name "<SUBSCRIPTION NAME OR ID>"

$ az group create \
    --name="<RESOURCE_GROUP_NAME>" \
    --location="<LOCATION>"

$ az resource group deployment create \
    --name="<DEPLOYMENT NAME>" \
    --resource-group="<RESOURCE_GROUP_NAME>" \
    --template-file-path="./_output/<INSTANCE>/azuredeploy.json" \
    --parameters-file-path="./_output/<INSTANCE>azuredeploy.parameters.json"
```

### Deploying with Powershell

```powershell
Add-AzureRmAccount

Select-AzureRmSubscription -SubscriptionID <SUBSCRIPTION_ID>

New-AzureRmResourceGroup `
    -Name <RESOURCE_GROUP_NAME> `
    -Location <LOCATION>

New-AzureRmResourceGroupDeployment `
    -Name <DEPLOYMENT_NAME> `
    -ResourceGroupName <RESOURCE_GROUP_NAME> `
    -TemplateFile _output\<INSTANCE>\azuredeploy.json `
    -TemplateParameterFile _output\<INSTANCE>\azuredeploy.parameters.json
```

