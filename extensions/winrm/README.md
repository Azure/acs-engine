# WinRM Extension

This extension enabled WinRM with HTTPS and a self-signed certificate on each Windows node. This makes it possible to troubleshoot Windows nodes from the master over ssh, without needing Remote Desktop.


## Use Cases

### Remote Shells

PowerShell can be used to connect to the Windows node from the same Azure vNet. If you want to connect from the Linux master, then it's easiest to run PowerShell as a container:

```bash
docker run -it mcr.microsoft.com/powershell pwsh
```

Then, you can run a few PowerShell commands to connect to the Windows node:

```powershell
# This will prompt for the Windows node's username & password
$credential = Get-Credential

# Now connect
Enter-PsSession <hostname> -Credential $credential -Authentication Basic -UseSSL
```

### Log Gathering

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
