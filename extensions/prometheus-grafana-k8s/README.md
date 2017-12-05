# prometheus-grafana Extension


This is the prometheus-grafana extension.  Add this extension to the api model you pass as input into acs-engine as shown below to automatically enable prometheus and grafana in your new Kubernetes cluster.

```
{
  "apiVersion": "vlabs",
  "properties": {
    "orchestratorProfile": {
      "orchestratorType": "Kubernetes"
    },
    "masterProfile": {
      "count": 1,
      "dnsPrefix": "",
      "vmSize": "Standard_DS2_v2",
      "extensions": [
        { 
          "name": "prometheus-grafana-k8s"
        }
      ]
    },
    "agentPoolProfiles": [
      {
        "name": "agentpool1",
        "count": 3,
        "vmSize": "Standard_DS2_v2",
        "availabilityProfile": "AvailabilitySet"
      }
    ],
    "linuxProfile": {
      "adminUsername": "azureuser",
      "ssh": {
        "publicKeys": [
          {
            "keyData": ""
          }
        ]
      }
    },
    "extensionProfiles": [
      { 
        "name": "prometheus-grafana-k8s", 
        "version": "v1",
        "rootURL": "https://raw.githubusercontent.com/Azure/acs-engine/master/"
      }
    ],
    "servicePrincipalProfile": {
      "clientId": "",
      "secret": ""
    }
  }
}
```


The following script will be executed on the master:

```
$ prometheus-grafana-k8s.sh
```

You can validate that the extension is running as expected with the following commands:

```
$ kubectl get pods --show-all

$ NAMESPACE=default
$ K8S_SECRET_NAME=dashboard-grafana

# Get user name and password for the grafana dashboard

$ GF_USER_NAME=$(kubectl get secret $K8S_SECRET_NAME -o jsonpath="{.data.grafana-admin-user}" | base64 --decode)
$ echo $GF_USER_NAME
$ GF_PASSWORD=$(kubectl get secret $K8S_SECRET_NAME -o jsonpath="{.data.grafana-admin-password}" | base64 --decode)
$ echo $GF_PASSWORD

# Forwarding Grafana port to localhost in a background job
$ GF_POD_NAME=$(kubectl get po -n $NAMESPACE -l "component=grafana" -o jsonpath="{.items[0].metadata.name}")
$ kubectl port-forward $GF_POD_NAME 3000:3000 &

```

# Configuration
|Name|Required|Acceptable Value|
|---|---|---|
|name|yes|prometheus-grafana-k8s|
|version|yes|v1|
|extensionParameters|no||
|rootURL|optional||

_Note_: specify a string for `extensionParameters` for a non-default namespace in the Kubernetes cluster

# Example
``` javascript
{ "name": "prometheus-grafana-k8s", "version": "v1", "extensionParameters": "monitoring" }
```

# Supported Orchestrators
Kubernetes
