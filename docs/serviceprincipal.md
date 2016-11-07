# Microsoft Azure Container Service Engine

## Service Principals

### Overview

Service Accounts in Azure are tied to Active Directory Service Principals. You can read more about
Service Principals and AD Applications: ["Application and service principal objects in Azure Active Directory"](https://azure.microsoft.com/en-us/documentation/articles/active-directory-application-objects/).

Kubernetes uses a Service Principal to talk to Azure APIs to dynamically manage
resources such as
[User Defined Routes](https://azure.microsoft.com/en-us/documentation/articles/virtual-networks-udr-overview/)
and [L4 Load Balancers](https://azure.microsoft.com/en-us/documentation/articles/load-balancer-overview/).

### Creating a Service Principal


There are several ways to create a Service Principal in Azure Active Directory:

* **With the [Azure CLI](https://github.com/Azure/azure-cli)**
   
   ```shell
   az login
   az account set --name="${SUBSCRIPTION_ID}"
   az ad sp create-for-rbac --role="Contributor" --scopes="/subscriptions/${SUBSCRIPTION_ID}"
   ```
   
   This will output your `client_id` and `client_secret` (`password`).


* **With the legacy [Azure XPlat CLI](https://github.com/Azure/azure-xplat-cli)**

   Instructions: ["Use Azure CLI to create a service principal to access resources"](https://azure.microsoft.com/en-us/documentation/articles/resource-group-authenticate-service-principal-cli/)

* **With [PowerShell](https://azure.microsoft.com/en-us/documentation/articles/resource-group-authenticate-service-principal-cli/)**

   Instructions: ["Use Azure PowerShell to create a service principal to access resources"](https://azure.microsoft.com/en-us/documentation/articles/resource-group-authenticate-service-principal-cli/)

* **With the [Legacy Portal](https://azure.microsoft.com/en-us/documentation/articles/resource-group-create-service-principal-portal/)**

   Instructions: ["Use portal to create Active Directory application and service principal that can access resources"](https://azure.microsoft.com/en-us/documentation/articles/resource-group-create-service-principal-portal/)
