# Chocolately Extension

This extension installs packages passed as parameters via the Chocolatey package manager.

Information about the Chocolatey package can be found at https://chocolatey.org.

# Configuration

|Name               |Required|Acceptable Value     |
|-------------------|--------|---------------------|
|name               |yes     |choco                |
|version            |yes     |v1                   |
|extensionParameters|yes     |microsoft-build-tools|
|rootURL            |optional|                     |

# Example

``` javascript
    ...
    "agentPoolProfiles": [
      {
        "name": "windowspool1",
        "extensions": [
          {
            "name": "choco"
          }
        ]
      }
    ],
    ...
    "extensionProfiles": [
      {
        "name": "choco",
        "version": "v1",
        "extensionParameters": "microsoft-build-tools"
      }
    ]
    ...
```

> Note: For multiple chocolatey packages you may provide a comma or semicolon separated list.

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
