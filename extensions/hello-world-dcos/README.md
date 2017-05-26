# hello-world-dcos Extension

Sample hello-world extension.  Calls the following on the master:

```
 curl -X post http://localhost:8080/v2/apps -d "{ \"id\": \"hello-marathon\", \"cmd\": \"while [ true ] ; do echo 'Hello World' ; sleep 5 ; done\", \"cpus\": 0.1, \"mem\": 10.0, \"instances\": 1 }" -H "Content-type:application/json"
```

You can validate that the extension was run by running (make sure you have tunneled into the master):
```
dcos auth login
dcos task log hello-marathon 
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
    "masterProfile": {
      ...
      "extensions": [
        { 
          "name": "hello-world-dcos"
        }
     ]
    },
    ...
    "extensionsProfile": [
      { 
        "name": "hello-world-dcos", 
        "version": "v1", 
        "rootURL": "https://bagbyimages.blob.core.windows.net:443/" 
      }
    ]
    

```

# Supported Orchestrators
"DCOS", "DCOS173", "DCOS184", "DCOS188"