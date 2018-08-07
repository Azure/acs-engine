# WinRM Extension

This extension enabled WinRM with HTTPS and a self-signed certificate on each Windows node. This makes it possible to troubleshoot Windows nodes from the master over ssh, without needing Remote Desktop.

If you want to pull all the logs off Windows nodes quickly, deploy with this extension, then SSH to the master node and run [logslurp](https://github.com/PatrickLang/logslurp) to gather them all at once.

# Configuration

|Name               |Required|Acceptable Value     |
|-------------------|--------|---------------------|
|name               |yes     |winrm                |
|version            |yes     |v1                   |
|rootURL            |optional|                     |

# Example

``` javascript
    ...
    "agentPoolProfiles": [
      {
        "name": "windowspool1",
        "extensions": [
          {
            "name": "winrm"
          }
        ]
      }
    ],
    ...
    "extensionProfiles": [
      {
        "name": "winrm",
        "version": "v1"
      }
    ]
    ...
```


# Supported Orchestrators

Kubernetes

# Troubleshoot

Extension execution output is logged to files found under the following directory on the target virtual machine.

```sh
C:\WindowsAzure\Logs\Plugins\Microsoft.Compute.CustomScriptExtension
```

The specified files are downloaded into the following directory on the target virtual machine.

```sh
C:\Packages\Plugins\Microsoft.Compute.CustomScriptExtension\1.*\Downloads\<n>
```
