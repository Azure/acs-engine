# microsoft-oms-agent-k8s Extension

The microsoft-oms-agent-k8s extension installs Microsoft OMS agent container on a kubernetes cluster.

# Configuration
|Name|Required|Acceptable Value|
|---|---|---|
|name|yes|microsoft-oms-agent-k8s|
|version|yes|v1|
|extensionParameters|yes|your WSID and KEY in json format|
|rootURL|no||

# Example
``` javascript
	{ 
        "name": "hello-world-k8s", 
        "version": "v1" 
		"extensionsParameters": {
			"WSID": "0a26d157-1234-4ed5-a1b1-ae4577e5ab88",
			"KEY": "Av/H0yNkanOicUUfnnf2vh1cCfEzFCupzdi6zBGn4Bo6g2HfRNcPK0OfeA2cAcA+Kfj+3hRu2xNJuG8MNWp5Pg=="
		}
    }
```

# extensionParameters
The Microsoft OMS Agent k8s extension requires your OMS WSID and Key to be placed in extensionParameters.  You can find this in your OMS Settings in the OMS Portal.  

# Supported Orchestrators
Kubernetes