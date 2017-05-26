# hello-world-k8s Extension

Sample hello-world extension.  Calls the following on the master:

```
 kubectl run hello-world --image=hello-world
```

You can validate that the extension was run by running:
```
kubectl get pods 
kubectl logs <name from kubectl get pods>
```

# Configuration
|Name|Required|Acceptable Value|
|---|---|---|
|name|yes|hello-world-k8s|
|version|yes|v1|
|extensionParameters|no||
|rootURL|optional||

# Example
``` javascript
{ "name": "hello-world-k8s", "version": "v1" }
```

# Supported Orchestrators
Kubernetes