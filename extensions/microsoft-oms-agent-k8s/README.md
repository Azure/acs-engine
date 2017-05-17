# sysdig-cloud-k8s Extension

The sysdig-cloud-k8s extension installs sysdig cloud on a kubernetes cluster.

# Configuration
|Name|Required|Acceptable Value|
|---|---|---|
|name|yes|sysdig-cloud-k8s|
|version|yes|v1|
|extensionParameters|yes|your sysdig key|
|rootURL|no||

# Example
``` javascript
{ "name": "sysdig-cloud-k8s", "version": "v1", "extensionParameters": "c18fa5f1-cda4-32ed-c725-s3ac0ae62110" }
```

# extensionParameters
The sysdig cloud k8s extension requires your sysdig cloud access key to be placed in extensionParameters.  You can find this in your user profile.  

# Supported Orchestrators
Kubernetes