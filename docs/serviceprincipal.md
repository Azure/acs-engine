# Microsoft Azure Container Service Engine - Creating a Service Principal

# Overview

Orchestrators such as Kubernetes require a service principal to dynamically adjust Azure resources.  This guide shows you how to create a service principal.

# Service Principal Creation

Here is how to create a service principal for your cluster:

1. Install the [Microsoft Azure CLI 2.0](https://github.com/azure/azure-cli) for your dev environment.
2. open up a command prompt
3. `az login` to login
4. If you have more than one subscription ID, execute `az account set -n SUBSCRIPTIONID` to select the correct subscription ID
5.  `az ad sp create-for-rbac --role contributor --scopes /subscriptions/SUBSCRIPTIONID` replacing SUBSCRIPTIONID with your subscription ID to create your client id and secret.  Copy and paste the result in a secure location.  `client_id` maps directly to servicePrincipalProfile.servicePrincipalClientID, and `client_secret` maps directly to servicePrincipalClientSecret.  Note that, you can further scope to resource group if you would like to target a specific resource group, eg scope of `/subscriptions/SUBSCRIPTIONID/resourcegroups/mygroup`
