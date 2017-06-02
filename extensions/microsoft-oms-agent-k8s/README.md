# microsoft-oms-agent-k8s Extension

The `microsoft-oms-agent-k8s` extension installs Microsoft OMS agent container on all the nodes (master and agent) in a kubernetes cluster.

# Configuration

|Name|Required|Acceptable Value|
|---|---|---|
|name|yes|microsoft-oms-agent-k8s|
|version|yes|v1|
|extensionParameters|yes|The base64 encoded json representation of your OMS WSID and KEY values.|
|rootURL|no||

# Extension Parameters

The Microsoft OMS Agent k8s extension requires your OMS WSID and Key to be placed in extensionParameters.  You can find this in your OMS Settings in the OMS Portal.  

The parameters for this extension must be provided in the following json format. 

``` javascript
{ 
  "WSID": "c714f34a-74cd-4bea-b1cb-b1af58a2ec1a", 
  "KEY": "dGhlIG9tcyBrZXkgdmFsdWUgZm9yIHdzaWQgYzcxNGYzNGEtNzRjZC00YmVhLWIxY2ItYjFhZjU4YTJlYzFhCg==" 
}
```
The json must then be base64 encoded before being passed into the `extensionParameters` value.

Here is an example in bash.
``` bash
$ printf '{ "WSID": "c714f34a-74cd-4bea-b1cb-b1af58a2ec1a", "KEY": "dGhlIG9tcyBrZXkgdmFsdWUgZm9yIHdzaWQgYzcxNGYzNGEtNzRjZC00YmVhLWIxY2ItYjFhZjU4YTJlYzFhCg==" }' | base64 -w0
eyAiV1NJRCI6ICJjNzE0ZjM0YS03NGNkLTRiZWEtYjFjYi1iMWFmNThhMmVjMWEiLCAiS0VZIjogImRHaGxJRzl0Y3lCclpYa2dkbUZzZFdVZ1ptOXlJSGR6YVdRZ1l6Y3hOR1l6TkdFdE56UmpaQzAwWW1WaExXSXhZMkl0WWpGaFpqVTRZVEpsWXpGaENnPT0iIH0=
```

Here is an example in PowerShell.
``` powershell
PS> $json = '{ "WSID": "c714f34a-74cd-4bea-b1cb-b1af58a2ec1a", "KEY": "dGhlIG9tcyBrZXkgdmFsdWUgZm9yIHdzaWQgYzcxNGYzNGEtNzRjZC00YmVhLWIxY2ItYjFhZjU4YTJlYzFhCg==" }'
PS> [Convert]::ToBase64String([System.Text.Encoding]::UTF8.GetBytes($json))
eyAiV1NJRCI6ICJjNzE0ZjM0YS03NGNkLTRiZWEtYjFjYi1iMWFmNThhMmVjMWEiLCAiS0VZIjogImRHaGxJRzl0Y3lCclpYa2dkbUZzZFdVZ1ptOXlJSGR6YVdRZ1l6Y3hOR1l6TkdFdE56UmpaQzAwWW1WaExXSXhZMkl0WWpGaFpqVTRZVEpsWXpGaENnPT0iIH0=
```

# Example



``` javascript
{ 
  "name": "microsoft-oms-agent-k8s", 
  "version": "v1" 
  "extensionsParameters": "eyAiV1NJRCI6ICJjNzE0ZjM0YS03NGNkLTRiZWEtYjFjYi1iMWFmNThhMmVjMWEiLCAiS0VZIjogImRHaGxJRzl0Y3lCclpYa2dkbUZzZFdVZ1ptOXlJSGR6YVdRZ1l6Y3hOR1l6TkdFdE56UmpaQzAwWW1WaExXSXhZMkl0WWpGaFpqVTRZVEpsWXpGaENnPT0iIH0="
}
```



# Supported Orchestrators
Kubernetes