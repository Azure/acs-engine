# Microsoft Azure Container Service Engine

The Azure Container Service Engine (`acs-engine`) generates ARM (Azure Resource Manager) templates for Docker enabled clusters on Microsoft Azure with your choice of DCOS, Kubernetes, or Swarm orchestrators. The input to acs-engine is a cluster definition file which describes the desired cluster, including orchestrator, features, and agents. The structure of the input files is very similar to the public API for Azure Container Service.

## Install

Binary downloads for the latest verison of acs-engine are available here:

* [OSX](https://github.com/Azure/acs-engine/releases/download/v0.4.0/acs-engine-v0.4.0-darwin-amd64.tar.gz)
* [Linux 64bit](https://github.com/Azure/acs-engine/releases/download/v0.4.0/acs-engine-v0.4.0-linux-amd64.tar.gz)
* [Windows 64bit](https://github.com/Azure/acs-engine/releases/download/v0.4.0/acs-engine-v0.4.0-windows-amd64.zip)

Download `acs-engine` for your operating system. Extract the binary and copy it to your `$PATH`.

If would prefer to build `acs-engine` from source or are you are interested in contributing to `acs-engine` see [building from source](#build-from-source) below.

## Usage

### Overview

`acs-engine` reads a JSON [cluster definition](../clusterdefinition.md) and generates a number of files that may be submitted to Azure Resource Manager (ARM). The possible outputs include:

1. **apimodel.json**: is an expanded version of the cluster definition provided to the generate command. All default or computed values will be expanded during the generate phase.
2. **azuredeploy.json**: represents a complete description of all Azure resources required to fulfill the cluster definition from `apimodel.json`.
3. **azuredeploy.parameters.json**: the parameters file holds a series of custom variables which are used in various locations throughout `azuredeploy.json`.
4. **certificate and access config files**: orchestrators like Kubernetes require certificates and additional configuration files (e.g. Kubernetes apiserver certificates and kubeconfig).

If you want to customize cluster configuaration after the `generate` step, make sure to modify `apimodel.json`. This ensures that any computed settings and certificates are correctly preserved. For example, if you want to add a second agent pool, you would edit `apimodel.json` and then run `acs-engine` against that file to generate the new ARM templates. This ensures that during the deploy steps resources remain untouched and only the new agent pools are created.

### Generating a template

Here is an example of how to generate a new deployment. This example assumes you are using [examples/kubernetes.json](../examples/kubernetes.json).

1. Before starting ensure you have generated a valid [SSH Public/Private key pair](ssh.md#ssh-key-generation).
2. edit [examples/kubernetes.json](../examples/kubernetes.json) and fill in the blanks.
3. run `./acs-engine generate examples/kubernetes.json` to generate the templates in the _output/Kubernetes-UNIQUEID directory.  The UNIQUEID is a hash of your master's FQDN prefix.
4. now you can use the `azuredeploy.json` and `azuredeploy.parameters.json` for deployment as described in [deployment usage](../README.md#deployment-usage).

### Deploying a template

For deployment see [deployment usage](../README.md#deployment-usage).

<a href="#build-from-source"></a>
## Build ACS Engine from Source

### Windows

Requirements:
- Git for Windows. Download and install [here](https://git-scm.com/download/win)
- Go for Windows. Download and install [here](https://golang.org/dl/), accept all defaults.
- Powershell 

Build Steps: 
 
1. Setup your go workspace.  This example assumes you are using `c:\gopath` as your workspace:
  1. Windows key-R to open the run prompt
  2. `rundll32 sysdm.cpl,EditEnvironmentVariables` to open the system variables
  3. add `c:\go\bin` to your PATH variables
  4. click "new" and add new environment variable GOPATH and set to `c:\gopath`
  
Build acs-engine:
  1. Windows key-R to open the run prompt
  2. `cmd` to open command prompt
  3. mkdir %GOPATH%
  4. cd %GOPATH%
  5. type `go get github.com/Azure/acs-engine` to get the acs-engine Github project
  6. type `go get all` to get the supporting components
  7. `cd %GOPATH%\src\github.com\Azure\acs-engine`
  8. `go build` to build the project
3. `acs-engine` to see the command line parameters

### OS X

Requirements:
- Go for OS X. Download and install [here](https://golang.org/dl/)

Build Steps: 

  1. Open a command prompt to setup your gopath:
  2. `mkdir $HOME/gopath`
  3. edit `$HOME/.bash_profile` and add the following lines to setup your go path
  ```
  export PATH=$PATH:/usr/local/go/bin
  export GOPATH=$HOME/gopath
  ```
  4. `source $HOME/.bash_profile`
Build acs-engine:
  1. type `go get github.com/Azure/acs-engine` to get the acs-engine Github project
  2. type `go get all` to get the supporting components
  3. `cd $GOPATH/src/github.com/Azure/acs-engine`
  4. `go build` to build the project
  5. `./acs-engine` to see the command line parameters

### Linux

Requirements:
- Go for Linux
  - Download the appropriate archive for your system [here](https://golang.org/dl/)
  - sudo tar -C /usr/local -xzf go$VERSION.$OS-$ARCH.tar.gz (replace with your downloaded archive)
- `git`

Build Steps: 

  1. Setup Go path:
  2. `mkdir $HOME/gopath`
  3. edit `$HOME/.profile` and add the following lines to setup your go path
  ```
  export PATH=$PATH:/usr/local/go/bin
  export GOPATH=$HOME/gopath
  ```
  4. `source $HOME/.profile`
 
Build acs-engine:
  1. type `go get github.com/Azure/acs-engine` to get the acs-engine Github project
  2. type `go get all` to get the supporting components
  3. `cd $GOPATH/src/github.com/Azure/acs-engine`
  4. `go build` to build the project
  5. `./acs-engine` to see the command line parameters

## Docker Development Environment

The easiest way to start hacking on `acs-engine` is to use a Docker-based environment. If you already have Docker installed then you can get started with a few commands.

* Windows (PowerShell): `.\scripts\devenv.ps1`
* Linux/OSX (bash): `./scripts/devenv.sh`

This script mounts the `acs-engine` source directory as a volume into the Docker container, which means you can edit your source code in your favorite editor on your machine, while still being able to compile and test inside of the Docker container. This environment mirrors the environment used in the acs-engine continuous integration (CI) system.

When the script `devenv.ps1` or `devenv.sh` completes, you will be left at a command prompt.

Run the following commands to pull the latest dependencies and build the `acs-engine` tool.

```
# install and download build dependencies
make prereqs
# build the `acs-engine` binary
make build
```

Th build process leaves the compiled `acs-engine` binary in the current directly. Make sure everything completed successfully bu running `acs-engine` without any arguments:

```
# ./acs-engine
ACS-Engine deploys and manages Kubernetes, Swarm Mode, and DC/OS clusters in Azure

Usage:
  acs-engine [command]

Available Commands:
  deploy      deploy an Azure Resource Manager template
  generate    Generate an Azure Resource Manager template
  help        Help about any command
  version     Print the version of ACS-Engine

Flags:
      --debug   enable verbose debug logs
  -h, --help    help for acs-engine

Use "acs-engine [command] --help" for more information about a command.
```

[Here's a quick demo video showing the dev/build/test cycle with this setup.](https://www.youtube.com/watch?v=lc6UZmqxQMs)
