## Deployment

Here are the steps to deploy a simple Kubernetes cluster:

1. [install acs-engine](../acsengine.md#downloading-and-building-acs-engine)
2. [generate your ssh key](../ssh.md#ssh-key-generation)
3. [generate your service principal](../serviceprincipal.md)
4. edit the [Kubernetes example](../examples/kubernetes.json) and fill in the blank strings
5. [generate the template](../acsengine.md#generating-a-template)
6. [deploy the output azuredeploy.json and azuredeploy.parameters.json](../README.md#deployment-usage)
  * To enable the optional network policy enforcement using calico, you have to
    set the parameter during this step according to this [guide](../kuberntes.md#optional-enable-network-policy-enforcement-using-calico)
7. Temporary workaround when deploying a cluster in a custom VNET with
   Kubernetes 1.6.0:
    1. After a cluster has been created in step 6 get id of the route table resource from Microsoft.Network provider in your resource group. 
       The route table resource id is of the format:
       `/subscriptions/SUBSCRIPTIONID/resourceGroups/RESOURCEGROUPNAME/providers/Microsoft.Network/routeTables/ROUTETABLENAME`
    2. Update properties of all subnets in the newly created VNET that are used by Kubernetes cluster to refer to the route table resource by appending the following to subnet properties:
        ```shell
        "routeTable": {
                "id": "/subscriptions/<SubscriptionId>/resourceGroups/<ResourceGroupName>/providers/Microsoft.Network/routeTables/<RouteTableResourceName>"
              }
        ```

        E.g.:
        ```shell
        "subnets": [
            {
              "name": "subnetname",
              "id": "/subscriptions/<SubscriptionId>/resourceGroups/<ResourceGroupName>/providers/Microsoft.Network/virtualNetworks/<VirtualNetworkName>/subnets/<SubnetName>",
              "properties": {
                "provisioningState": "Succeeded",
                "addressPrefix": "10.240.0.0/16",
                "routeTable": {
                  "id": "/subscriptions/<SubscriptionId>/resourceGroups/<ResourceGroupName>/providers/Microsoft.Network/routeTables/<RouteTableResourceName>"
                }
              ....
              }
              ....
            }
        ]
        ```
