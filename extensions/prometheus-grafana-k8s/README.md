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


The following script will be executed on the agent nodes:

```
$ prometheus-grafana-k8s.sh
```

You can validate that the extension is running as expected with the following commands:

```
$ kubectl get pods --all-namespaces

NAMESPACE     NAME                                                        READY     STATUS    RESTARTS   AGE
default       cadvisor-pshlh                                              1/1       Running   0          32m
default       cadvisor-z84d6                                              1/1       Running   0          32m
default       dashboard-grafana-5d864656b8-m47q4                          1/1       Running   0          31m
default       monitoring-prometheus-kube-state-metrics-6bc48cd465-mfrzm   1/1       Running   0          32m
default       monitoring-prometheus-node-exporter-74s6d                   1/1       Running   0          32m
default       monitoring-prometheus-node-exporter-c44ff                   1/1       Running   0          32m
default       monitoring-prometheus-pushgateway-d4f679b7-zstwf            1/1       Running   0          32m
default       monitoring-prometheus-server-75bb797794-rmjft               2/2       Running   0          32m
kube-system   calico-node-2b6tz                                           2/2       Running   0          34m
kube-system   calico-node-cf29f                                           2/2       Running   0          34m
kube-system   calico-node-x86bl                                           2/2       Running   0          34m
kube-system   heapster-568476f785-hp8c6                                   2/2       Running   0          32m
kube-system   kube-addon-manager-k8s-master-35213955-0                    1/1       Running   0          33m
kube-system   kube-apiserver-k8s-master-35213955-0                        1/1       Running   2          33m
kube-system   kube-controller-manager-k8s-master-35213955-0               1/1       Running   0          33m
kube-system   kube-dns-v20-59b4f7dc55-hjc7q                               3/3       Running   0          34m
kube-system   kube-dns-v20-59b4f7dc55-mqs7d                               3/3       Running   0          34m
kube-system   kube-proxy-l8crw                                            1/1       Running   0          34m
kube-system   kube-proxy-p966n                                            1/1       Running   0          34m
kube-system   kube-proxy-zvjgk                                            1/1       Running   0          34m
kube-system   kube-scheduler-k8s-master-35213955-0                        1/1       Running   0          33m
kube-system   kubernetes-dashboard-64dcf5784f-sr464                       1/1       Running   0          34m
kube-system   metrics-server-7fcdc5dbb9-6scj6                             1/1       Running   1          34m
kube-system   tiller-deploy-d85ccb55c-rz897                               1/1       Running   0          34m

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
|extensionParameters|no|see below|
|rootURL|optional||

_Note_: the format for `extensionParameters` is the following: `"<namespace>;<prometheus_values_config_url>;<cadvisor_daemonset_config_url>"`. Each of these placeholders are optional (as is the entire `extensionParameters` itself)

# Example
``` javascript
{ "name": "prometheus-grafana-k8s", "version": "v1", "extensionParameters": "monitoring;;" }
```

# Supported Orchestrators
Kubernetes
