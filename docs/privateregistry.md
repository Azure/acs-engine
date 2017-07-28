# Private Registry for DCOS clusters

DCOS needs credentials to deploy containers from a private registry. This could be an Azure Container Registry or another private registry.

The [credentials are placed on the agent nodes]https://mesosphere.github.io/marathon/docs/native-docker-private-registry.html() for DCOS to present them to the registry when pulling a container.

When deploying a service you simply reference the credentials file created by the DCOS container deployment:

```
{
  "id": "/first",
  "fetch": [
    {
      "uri": "file:///etc/docker.tar.gz"
    }
  ],
  "container": {
    "docker": {
      "image": "xtophreg-microsoft.azurecr.io/first",
      "forcePullImage": false,
      "privileged": false,
      "network": "HOST"
    }
  }
} 
```

You define the credentials in the ACS api model:
```
{
  "apiVersion": "vlabs",
  "properties": {
    "orchestratorProfile": {
      "orchestratorType": "DCOS184",
      "Registry" : "yourregistry-microsoft.azurecr.io",
      "RegistryUser" : "registryadminuser",
      "RegistryPassword" : "registrypassword"
    }
    /// ...
  }
}
```

and then you deploy the cluster via ```acs-engine.exe```

