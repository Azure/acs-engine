# vamp-dcos Extension

Installs VAMP 0.9.4 as per the [VAMP Standard Install](http://vamp.io/documentation/installation/v0.9.4/dcos/#standard-install).
Should only be ran on a single master after the waitforall extension has succeeded

# Configuration
|Name|Required|Acceptable Value|
|---|---|---|
|name|yes|vamp-dcos|
|version|yes|v1|
|extensionParameters|no||
|rootURL|optional||
|singleOrAll|Required|single|

# Example
```javascript
    "masterProfile": {
        ...
      "extensions": [
        { 
          "name": "vamp-dcos", 
          "singleOrAll": "single"
        }
        ...
     ]
    },
    ...
    "extensionsProfile": [
      { 
        "name": "vamp-dcos", 
        "version": "v1"
      }
    ]
    
```

# Supported Orchestrators
DCOS