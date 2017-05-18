# helm-k8s Extension

The helm-k8s extension installs the helm-tiller service on a kubernetes cluster.

# Configuration
|Name|Required|Acceptable Value|
|---|---|---|
|name|yes|helm-k8s|
|version|yes|v1|
|rootURL|no||

# Example
``` javascript
    { 
        "name": "helm-k8s", 
        "version": "v1" 
    }
```

# Supported Orchestrators
Kubernetes