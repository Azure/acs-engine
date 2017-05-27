# hello-world-k8s Extension

Sample hello-world extension.  Calls the following on the master:

```
 kubectl run hello-world --quiet --image=busybox --restart=OnFailure -- echo "Hello Kubernetes!"
```

You can validate that the extension was run by running:
```
kubectl get pods --show-all
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