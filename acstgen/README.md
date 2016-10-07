# ACSTGEN - Template generator

Template generator builds a custom template based on user requirements.  Examples exist under clusterdefinitions folder.

# Building the Application

You must first install Go (or use a Docker container, see below). Once you have an environment with Go installed you can uild the applicaton with ```go
build```.

# Generating a Template

Once build run the generator with the command ```./acstgen
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
