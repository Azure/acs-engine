# ACSEngine - Template generator

Template generator builds a custom template based on user requirements.  Examples exist under clusterdefinitions folder.

# Building the Application

To build the application:
1. set $GOPATH (example c:\gopath)
2. go get github.com/Azure/acs-engine
3. cd $GOPATH/github.com/Azure/acs-engine/src
4. do "go get -u github.com/jteeuwen/go-bindata/..." if you don't have it yet
5. go generate to generate source files needed
6. go build to produce the acsengine binary

# Generating a Template

Once build run the generator with the command ```./acsengine
CLUSTER_DEFINITION_FILE``` where ```CLUSTER_DEFINITION_FILE``` is the
path to the cluster definition you want to use. The application
outputs therequired ARM template. This can be piped to a file for
later use (see below).

There are some example definition files in the folder `clusterdefinitions`.

# Using a Template

Generated templates can be deployed using the Azure CLI or Powershell. 

## Deploying with Powershell

To deploy an ACS instance using Powershell and your generated template
run the following two commands:

``` Powershell
New-AzureRmResourceGroup -Name <RESOURCE_GROUP_NAME> -Location <LOCATION> -Force
New-AzureRmResourceGroupDeployment -Name <DEPLOYMENT_NAME> -ResourceGroupName <RESOURCE_GROUP_NAME> RGName  -TemplateFile <TEMPLATE_FILE>
```

## Deploying with Azure CLI

To deploy an ACS instance using the Azure CLI and your generated
template run the following two commands:

``` bash
azure group create <RESOURCE_GROUP_NAME> <LOCATION>
azure group deployment create <RESOURCE_GROUP_NAME> <DEPLOYMENT NAME> -f <TEMPLATE_FILE>

```

# Using a Docker container

The ```scripts``` directory contains helper scripts that will assist
with the using a Docker container to work with the appication.

## dev.sh

Run a Go container with the application source code mounted into the
container. You can edit the code in your favorite editor on the client
while building and running the container.

## generate.sh

Generate a template from a given configuration and store it in the
```generated``` folder. For example, to generate a Swarm template use:

``` bash
./scripts/generate.sh swarm.json
```
