# Microsoft Azure Container Service Engine

The Azure Container Service Engine (`acs-engine`) generates ARM (Azure Resource Manager) templates for Docker enabled clusters on Microsoft Azure with your choice of DCOS, Kubernetes, or Swarm orchestrators. The input to the tool is a cluster definition. The cluster definition is very similar to (in many cases the same as) the ARM template syntax used to deploy a Microsoft Azure Container Service cluster.

# Development in Docker

The easiest way to get started developing on `acs-engine` is to use Docker. If you already have Docker or "Docker for {Windows,Mac}" then you can get started without needing to install anything extra.

* Windows (PowerShell): `.\scripts\devenv.ps1`
* Linux/OSX (bash): `./scripts/devenv.sh`

This setup mounts the `acs-engine` source directory as a volume into the Docker container.
This means that you can edit your source code normally in your favorite editor on your
machine, while still being able to compile and test inside of the Docker container (the
same environment used in our Continuous Integration system).

When the execution of `devenv.{ps1,sh}` completes, you should find the console logged into the container. As a final step, in order to get the `acs-engine` tool ready, you should build the sources with:

```
make build
```

When the build process completes, verify that `acs-engine` is available, invoking the command without parameters. 
You should see something like this:

```
# acs-engine
Usage of acs-engine:
  -artifacts string
    	directory where artifacts will be written
  -caCertificatePath string
    	the path to the CA Certificate file
  -caKeyPath string
    	the path to the CA key file
  -classicMode
    	enable classic parameters and outputs
  -noPrettyPrint
    	do not pretty print output
  -parametersOnly
    	only output the parameters
```

[Here's a quick demo video showing the dev/build/test cycle with this setup.](https://www.youtube.com/watch?v=lc6UZmqxQMs)

# Downloading and Building ACS Engine Locally 

ACS Engine can also be built and run natively on Windows, OS X, and Linux. Instructions below: 

## Windows

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

## OS X

Requirements:
- Go for OS X. Download and install [here](https://golang.org/dl/)

Build Steps: 

  1. Open a command prompt to setup your gopath:
  2. `mkdir $HOME/gopath`
  3. edit `$HOME/.bash_profile` and add the following lines to setup your go path
  ```
  export PATH=$PATH:/usr/local/go/bin
  export GOPATH=$HOME/gopath
  source $HOME/.sh_profile
  ```
Build acs-engine:
  1. type `go get github.com/Azure/acs-engine` to get the acs-engine Github project
  2. type `go get all` to get the supporting components
  3. `cd $GOPATH/src/github.com/Azure/acs-engine`
  4. `go build` to build the project
  5. `./acs-engine` to see the command line parameters

## Linux

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
  4. source $HOME/.profile
 
Build acs-engine:
  1. type `go get github.com/Azure/acs-engine` to get the acs-engine Github project
  2. type `go get all` to get the supporting components
  3. `cd $GOPATH/src/github.com/Azure/acs-engine`
  4. `go build` to build the project
  5. `./acs-engine` to see the command line parameters


# Template Generation

The `acs-engine` takes a json [cluster definition file](clusterdefinition.md) as a parameter and generates 3 or more of the following files:

1. **apimodel.json** - this is the cluster configuration file used for generation
2. **azuredeploy.json** - this is the main ARM (Azure Resource Model) template used to deploy a full Docker enabled cluster
3. **azuredeploy.parameters.json** - this is the parameters file used along with azurdeploy.json during deployment and contains configurable parameters
4. **certificate and access config files** - some orchestrators like Kubernetes require certificate generation, and these generated files and access files like the kube config files are stored along side the model and ARM template files.

As a rule of thumb you should always work with the `apimodel.json` when modifying an existing running deployment.  This ensures that all the same settings and certificates are correctly preserved.  For example, if you want to add a second agent pool, you would edit `apimodel.json` and then run `acs-engine` against that file to generate the new ARM templates. Then during deployment all existing deployments remain untouched, and only the new agent pools resources are created.

# Generating a template

Here is an example of how to generate a new deployment.  This example assumes you are using [examples/kubernetes.json](../examples/kubernetes.json).

1. Before starting ensure you have generated a valid [SSH Public/Private key pair](ssh.md#ssh-key-generation).
2. edit [examples/kubernetes.json](../examples/kubernetes.json) and fill in the blanks.
3. run `acs-engine examples/kubernetes.json` to generate the templates in the _output/Kubernetes-UNIQUEID directory.  The UNIQUEID is a hash of your master's FQDN prefix.
4. now you can use the `azuredeploy.json` and `azuredeploy.parameters.json` for deployment as described in [deployment usage](../README.md#deployment-usage).

# Deploying templates

For deployment see [deployment usage](../README.md#deployment-usage).
