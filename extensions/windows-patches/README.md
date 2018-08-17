# Windows Patching Extension

This extension will install Windows Server patches, including prerelease hotfixes. It's useful for the following cases:

1. Microsoft support has provided a pre-release hotfix for your testing
2. Installing additional Windows update packages (MSU) that are not yet included in the default Windows Server with Containers VM on the Azure Marketplace.


# Configuration

|Name               |Required| Acceptable Value     |
|-------------------|--------|----------------------|
|name               |yes     | windows-patches      |
|version            |yes     | v1                   |
|rootURL            |optional| `https://raw.githubusercontent.com/Azure/acs-engine/master/` or any repo with the same extensions/... directory structure |
|extensionParameters|yes     | comma-delimited list of URIs enclosed with ' such as `'https://privateupdates.domain.ext/Windows10.0-KB999999-x64-InstallForTestingPurposesOnly.exe', 'https://privateupdates.domain.ext/Windows10.0-KB123456-x64-InstallForTestingPurposesOnly.exe'` |

# Example

```json
    ...
    "agentPoolProfiles": [
      {
        "name": "windowspool1",
        "extensions": [
          {
            "name": "windows-patches"
          }
        ]
      }
    ],
    ...
    "extensionProfiles": [
      {
        "name": "windows-patches",
        "version": "v1",
        "rootURL": "https://raw.githubusercontent.com/Azure/acs-engine/master/",
        "extensionParameters": "'https://privateupdates.domain.ext/Windows10.0-KB999999-x64-InstallForTestingPurposesOnly.exe', 'https://privateupdates.domain.ext/Windows10.0-KB123456-x64-InstallForTestingPurposesOnly.exe'"
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
