# prometheus-grafana Extension


Sample prometheus-grafana extension.  Runs the following on the master:

```
 
```

You can validate that the extension was run by running:
```
kubectl get pods --show-all

```

# Configuration
|Name|Required|Acceptable Value|
|---|---|---|
|name|yes|prometheus-grafana-k8s|
|version|yes|v1|
|extensionParameters|no||
|rootURL|optional||

# Example
``` javascript
{ "name": "prometheus-grafana-k8s", "version": "v1" }
```

# Supported Orchestrators
Kubernetes