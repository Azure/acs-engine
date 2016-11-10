# Microsoft Azure Container Service Engine

The Azure Container Service Engine (`acs-engine`) generates ARM (Azure Resource Manager) templates for Docker enabled clusters on Microsoft Azure with your choice of DCOS, Kubernetes, or Swarm orchestrators. The input to the tool is a cluster definition. The cluster definition is very similar to (in many cases the same as) the ARM template syntax used to deploy a Microsoft Azure Container Service cluster.

# Downloading and Building ACS Engine

As mentioned in the README, the simpliest way to run ACS Engine is via a Docker container (insrtuctions in the README). 

That said, ACS Engine can also be built and run natively on Windows OS X, and Linux. Instructions below: 

## Windows

Requirements:
- Powershell 
- Git for Windows. Download and install [here](https://git-scm.com/download/win)
- Go for Windows. Download and install [here](https://golang.org/dl/), accept all defaults.

Build Steps: 
 
1. Setup your go workspace.  This example assumes you are using `c:\gopath` as your workspace:
  1. Windows key-R to open the run prompt
  2. `rundll32 sysdm.cpl,EditEnvironmentVariables` to open the system variables
  3. add `c:\go\bin` to your PATH variables
  4. click "new" and add new environment variable GOPATH and set to `c:\gopath`
2. Build acs-engine:
  1. Windows key-R to open the run prompt
  2. `cmd` to open command prompt
  3. mkdir %GOPATH%
  4. cd %GOPATH%
  5. type `go get github.com/Azure/acs-engine` to get the acs-engine Github project
  6. type `go get all` to get the supporting components
  7. `cd %GOPATH%\src\github.com\Azure\acs-engine`
  8. `go build` to build the project
3. `acs-engine` to see the command line parameters

## OS X

Requirements:
- Go for OS X. Download and install [Go for OS X](https://golang.org/dl/)

Build Steps: 

1. Open a command prompt to setup your gopath:
  1. `mkdir $HOME/go`
  2. edit `$HOME/.sh_profile` and add the following line to setup your go path
  ```
  export PATH=$PATH:/usr/local/go/bin
  export GOPATH=$HOME/go
  ```
  3. `source $HOME/.sh_profile`
2. Build acs-engine:
  1. type `go get github.com/Azure/acs-engine` to get the acs-engine Github project
  2. type `go get all` to get the supporting components
  3. `cd $GOPATH/src/github.com/Azure/acs-engine`
  4. `go build` to build the project
3. `./acs-engine` to see the command line parameters

## Linux

Requirements:
- Go for Linux. Download and install [Go for Linux](https://golang.org/dl/)
- `git`
- https://github.com/Azure/acs-engine#development-docker

Build Steps: 

1. Setup Go path:
  1. `mkdir $HOME/go`
  2. edit `$HOME/.profile` and add the following line to setup your go path
  ```
  export PATH=$PATH:/usr/local/go/bin
  export GOPATH=$HOME/go
  ```
  3. `source $HOME/.profile`
2. Build acs-engine:
  1. type `go get github.com/Azure/acs-engine` to get the acs-engine Github project
  2. type `go get all` to get the supporting components
  3. `cd $GOPATH/src/github.com/Azure/acs-engine`
  4. `go build` to build the project
3. `./acs-engine` to see the command line parameters


# Template Generation

The `acs-engine` takes a json [cluster definition file](clusterdefinition.md) as a parameter and generates 3 or more of the following files:

1. **apimodel.json** - this is the cluster configuration file used for generation
2. **azuredeploy.json** - this is the main ARM (Azure Resource Model) template used to deploy a full Docker enabled cluster
3. **azuredeploy.parameters.json** - this is the parameters file used along with azurdeploy.json during deployment and contains configurable parameters
4. **certificate and access config files** - some orchestrators like Kubernetes require certificate generation, and these generated files and access files like the kube config files are stored along side the model and ARM template files.

As a rule of thumb you should always work with the `apimodel.json` when modifying an existing running deployment.  This ensures that all the same settings and certificates are correctly preserved.  For example, if you want to add a second agent pool, you would edit `apimodel.json` and then run `acs-engine` against that file to generate the new ARM templates. Then during deployment all existing deployments remain untouched, and only the new agent pools resources are created.

# Generating a template

Here is an example of how to generate a new deployment.  This example assumes you are using [examples/dcos.json](../examples/dcos.json).

1. Before starting ensure you have generated a valid [SSH Public/Private key pair](ssh.md#ssh-key-generation).
2. edit [examples/dcos.json](../examples/dcos.json) and fill in the blanks.
3. run `acs-engine examples/dcos.json` to generate the templates in the _output/DCOS184-UNIQUEID directory.  The UNIQUEID is a hash of your master's FQDN prefix.
4. now you can use the `azuredeploy.json` and `azuredeploy.parameters.json` for deployment as described in [deployment usage](../README.md#deployment-usage).

# Deploying templates

For deployment see [deployment usage](../README.md#deployment-usage).
