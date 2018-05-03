# Build Docker image

**Bash**
```bash
$ VERSION=0.16.0
$ docker build --no-cache --build-arg BUILD_DATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"` --build-arg ACSENGINE_VERSION="$VERSION" -t microsoft/acs-engine:$VERSION --file ./Dockerfile.linux .
```
**PowerShell**
```powershell
PS> $VERSION="0.16.0"
PS> docker build --no-cache --build-arg BUILD_DATE=$(Get-Date((Get-Date).ToUniversalTime()) -UFormat "%Y-%m-%dT%H:%M:%SZ") --build-arg ACSENGINE_VERSION="$VERSION" -t microsoft/acs-engine:$VERSION --file .\Dockerfile.linux .
```

# Inspect Docker image metadata

**Bash**
```bash
$ docker image inspect microsoft/acs-engine:0.16.0 --format "{{json .Config.Labels}}" | jq
{
  "maintainer": "Microsoft",
  "org.label-schema.build-date": "2017-10-25T04:35:06Z",
  "org.label-schema.description": "The Azure Container Service Engine (acs-engine) generates ARM (Azure Resource Manager) templates for Docker enabled clusters on Microsoft Azure with your choice of DCOS, Kubernetes, or Swarm orchestrators.",
  "org.label-schema.docker.cmd": "docker run -v ${PWD}:/acs-engine/workspace -it --rm microsoft/acs-engine:0.16.0",
  "org.label-schema.license": "MIT",
  "org.label-schema.name": "Azure Container Service Engine (acs-engine)",
  "org.label-schema.schema-version": "1.0",
  "org.label-schema.url": "https://github.com/Azure/acs-engine",
  "org.label-schema.usage": "https://github.com/Azure/acs-engine/blob/master/docs/acsengine.md",
  "org.label-schema.vcs-url": "https://github.com/Azure/acs-engine.git",
  "org.label-schema.vendor": "Microsoft",
  "org.label-schema.version": "0.16.0"
}
```

**PowerShell**
```powershell
PS> docker image inspect microsoft/acs-engine:0.16.0 --format "{{json .Config.Labels}}" | ConvertFrom-Json | ConvertTo-Json
{
    "maintainer":  "Microsoft",
    "org.label-schema.build-date":  "2017-10-25T04:35:06Z",
    "org.label-schema.description":  "The Azure Container Service Engine (acs-engine) generates ARM (Azure Resource Manager) templates for Docker enabled clusters on Microsoft Azure with your choice of DCOS, Kubernetes, or Swarm orchestrators.",
    "org.label-schema.docker.cmd":  "docker run -v ${PWD}:/acs-engine/workspace -it --rm microsoft/acs-engine:0.16.0",
    "org.label-schema.license":  "MIT",
    "org.label-schema.name":  "Azure Container Service Engine (acs-engine)",
    "org.label-schema.schema-version":  "1.0",
    "org.label-schema.url":  "https://github.com/Azure/acs-engine",
    "org.label-schema.usage":  "https://github.com/Azure/acs-engine/blob/master/docs/acsengine.md",
    "org.label-schema.vcs-url":  "https://github.com/Azure/acs-engine.git",
    "org.label-schema.vendor":  "Microsoft",
    "org.label-schema.version":  "0.16.0"
}
```

# Run Docker image

```
$ docker run -v ${PWD}:/acs-engine/workspace -it --rm microsoft/acs-engine:0.16.0
```
