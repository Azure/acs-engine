# Microsoft Azure Container Service Engine

The Azure Container Service Engine (`acs-engine`) generates ARM (Azure Resource Manager) templates for Docker enabled clusters on Microsoft Azure with your choice of DCOS, [Kubernetes](kubernetes/deploy.md), or Swarm orchestrators. The input to acs-engine is a cluster definition file which describes the desired cluster, including orchestrator, features, and agents. The structure of the input files is very similar to the public API for Azure Container Service.

<a href="#install-acs-engine"></a>

## Install

Binary downloads for the latest version of acs-engine for are available [here](https://github.com/Azure/acs-engine/releases/latest). Download `acs-engine` for your operating system. Extract the binary and copy it to your `$PATH`.

If would prefer to build `acs-engine` from source or are you are interested in contributing to `acs-engine` see [building from source](#build-acs-engine-from-source) below.

## Usage

`acs-engine` reads a JSON [cluster definition](./clusterdefinition.md) and generates a number of files that may be submitted to Azure Resource Manager (ARM). The generated files include:

1. **apimodel.json**: is an expanded version of the cluster definition provided to the generate command. All default or computed values will be expanded during the generate phase.
2. **azuredeploy.json**: represents a complete description of all Azure resources required to fulfill the cluster definition from `apimodel.json`.
3. **azuredeploy.parameters.json**: the parameters file holds a series of custom variables which are used in various locations throughout `azuredeploy.json`.
4. **certificate and access config files**: orchestrators like Kubernetes require certificates and additional configuration files (e.g. Kubernetes apiserver certificates and kubeconfig).

### Generate Templates

ACS Engine consumes a cluster definition which outlines the desired shape, size, and configuration of Kubernetes. There are a number of features that can be enabled through the cluster definition.

See [ACS Engine The Long Way](kubernetes/deploy.md#acs-engine-the-long-way) for an example on generating templates by hand.

<a href="#deployment-usage"></a>

### Deploy Templates

Generated templates can be deployed using [the Azure CLI 2.0](https://github.com/Azure/azure-cli) or [Powershell](https://github.com/Azure/azure-powershell).

#### Deploying with Azure CLI 2.0

Azure CLI 2.0 is the latest CLI maintained and supported by Microsoft. For installation instructions see [the Azure CLI GitHub repository](https://github.com/Azure/azure-cli#installation) for the latest release.

```bash
$ az login

$ az account set --subscription "<SUBSCRIPTION NAME OR ID>"

$ az group create \
    --name "<RESOURCE_GROUP_NAME>" \
    --location "<LOCATION>"

$ az group deployment create \
    --name "<DEPLOYMENT NAME>" \
    --resource-group "<RESOURCE_GROUP_NAME>" \
    --template-file "./_output/<INSTANCE>/azuredeploy.json" \
    --parameters "./_output/<INSTANCE>/azuredeploy.parameters.json"
```

#### Deploying with Powershell

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

<a href="#build-from-source"></a>

## Build ACS Engine from Source

### Docker Development Environment

The easiest way to start hacking on `acs-engine` is to use a Docker-based environment. If you already have Docker installed then you can get started with a few commands.

* Windows (PowerShell): `.\scripts\devenv.ps1`
* Linux/OSX (bash): `./scripts/devenv.sh`

This script mounts the `acs-engine` source directory as a volume into the Docker container, which means you can edit your source code in your favorite editor on your machine, while still being able to compile and test inside of the Docker container. This environment mirrors the environment used in the acs-engine continuous integration (CI) system.

When the script `devenv.ps1` or `devenv.sh` completes, you will be left at a command prompt.

Run the following commands to pull the latest dependencies and build the `acs-engine` tool.

```
# install and download build dependencies
make bootstrap
# build the `acs-engine` binary
make build
```

The build process leaves the compiled `acs-engine` binary in the `bin` directory. Make sure everything completed successfully by running `bin/acs-engine` without any arguments:

```
# ./bin/acs-engine
ACS-Engine deploys and manages Kubernetes, Swarm Mode, and DC/OS clusters in Azure

Usage:
  acs-engine [command]

Available Commands:
  deploy      deploy an Azure Resource Manager template
  generate    Generate an Azure Resource Manager template
  help        Help about any command
  version     Print the version of ACS-Engine
  scale       Scale a existing cluster deployed by acs-engine

Flags:
      --debug   enable verbose debug logs
  -h, --help    help for acs-engine

Use "acs-engine [command] --help" for more information about a command.
```

[Here's a quick demo video showing the dev/build/test cycle with this setup.](https://www.youtube.com/watch?v=lc6UZmqxQMs)

## Building on Windows, OSX, and Linux

Building ACS Engine from source has a few requirements for each of the platforms. Download and install the pre-reqs for your platform, Windows, Linux, or Mac:

### Prerequisite
1. Go version 1.8 [installation instructions](https://golang.org/doc/install)
2. Git Version Control [installation instructions](https://git-scm.com/download/)

### Windows

Setup steps:

1. Setup your go workspace. This guide assumes you are using `c:\gopath` as your Go workspace:
  1. Type Windows key-R to open the run prompt
  2. Type `rundll32 sysdm.cpl,EditEnvironmentVariables` to open the system variables
  3. Add `c:\go\bin` and `c:\gopath\bin` to your PATH variables
  4. Click "new" and add new environment variable named `GOPATH` and set the value to `c:\gopath`

Build acs-engine:
  1. Type Windows key-R to open the run prompt
  2. Type `cmd` to open a command prompt
  3. Type `mkdir %GOPATH%` to create your gopath
  4. Type `cd %GOPATH%`
  5. Type `go get -d github.com/Azure/acs-engine` to download acs-engine from GitHub
  6. Type `go get all` to get the supporting components
  7. Type `go get -u github.com/jteeuwen/go-bindata/...`
  8. Type `cd %GOPATH%\src\github.com\Azure\acs-engine\pkg\acsengine`
  9. Type `go generate`
  10. Type `cd %GOPATH%\src\github.com\Azure\acs-engine\pkg\i18n`
  11. Type `go generate`
  12. Type `cd %GOPATH%\src\github.com\Azure\acs-engine`
  13. Type `go build` to build the project
  14. Type `go install` to install the project
  15. Run `acs-engine.exe` to see the command line parameters

### OS X and Linux

Setup steps:

  1. Open a command prompt to setup your gopath:
  2. `mkdir $HOME/go`
  3. edit `$HOME/.bash_profile` and add the following lines to setup your go path
  ```
  export GOPATH=$HOME/go
  export PATH=$PATH:/usr/local/go/bin:$GOPATH/bin
  ```
  4. `source $HOME/.bash_profile`

Build acs-engine:
  1. Type `go get github.com/Azure/acs-engine` to get the acs-engine Github project
  2. Type `cd $GOPATH/src/github.com/Azure/acs-engine` to change to the source directory
  3. Type `make bootstrap` to install supporting components
  4. Type `make build` to build the project
  5. Type `./bin/acs-engine` to see the command line parameters
