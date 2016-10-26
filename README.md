# Azure Container Service ARM Template Generator

## Overview

`acs-engine` is a tool that produces ARM (Azure Resource Manager) templates for deploying various container orchestrators into Azure.
The input to the tool is a cluster definition. The cluster definition is very similar to (in many cases the same as) the ARM template
syntax used to deploy an Azure Container Service cluster.

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

## Usage (Deployment)

Generated templates can be deployed using
[the Azure CLI](https://github.com/Azure/azure-cli) or
[Powershell](https://github.com/Azure/azure-powershell).

### Deploying with Azure CLI

```bash
$ az login

$ az account set --name "<SUBSCRIPTION NAME OR ID>"

$ az group create \
    --name="<RESOURCE_GROUP_NAME>" \
    --location="<LOCATION>"

$ az group deployment create \
    --name="<DEPLOYMENT NAME>" \
    --resource-group="<RESOURCE_GROUP_NAME>" \
    --template-file="./_output/<INSTANCE>/azuredeploy.json"
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
    -TemplateFile _output\<INSTANCE>\azuredeploy.json
```

