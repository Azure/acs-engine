# Microsoft Azure Container Service Engine - Kubernetes Windows Walkthrough

## Supported Windows versions
Prior to acs-engine v0.9.2, Kubernetes Windows cluster uses Windows Server 2016. There are a few restrictions in Windows Networking for Kubernetes as documented in https://blogs.technet.microsoft.com/networking/2017/04/04/windows-networking-for-kubernetes/. Besides, Windows POD deployment performanace is limited due to the bottleneck of container image size and configuration at container start time. 

With the release of new Windows Server version 1709, acs-engine v0.9.2 and beyond has leveraged the new Windows version to deploy Kubernetes Windows cluster with signifcant improvement in Windows container and networking performance, as well as new features in storage. Specifically,
1. Windows is now on par with Linux in terms of networking. New features including hostport have been implemented in kube-proxy and Windows platform and CNI to enhance networking performance. Please refer to http://blog.kubernetes.io/2017/09/windows-networking-at-parity-with-linux.html for details.
2. Azure Files and Disks are now supported to mount on Kubernetes Windows cluster with the new SMB feature in Windows.
3. Multiple containers in POD are now supported on Kubernetes Windows cluster.

Note, with the rollout of new Windows Server versions in acs-engine, the workload deployed on Windows cluster should match as documented here: https://docs.microsoft.com/en-us/virtualization/windowscontainers/deploy-containers/version-compatibility . Both the containers being deployed,
as well as the `kubletwin/pause` container must match the version of the Windows host. Otherwise, [pods may get stuck at ContainerCreating](https://docs.microsoft.com/en-us/virtualization/windowscontainers/kubernetes/common-problems#my-kubernetes-pods-are-stuck-at-containercreating) state.

## Deployment

Here are the steps to deploy a simple Kubernetes cluster with Windows:

1. [install acs-engine](../acsengine.md#downloading-and-building-acs-engine)
2. [generate your ssh key](../ssh.md#ssh-key-generation)
3. [generate your service principal](../serviceprincipal.md)
4. edit the [Kubernetes windows example](../../examples/windows/kubernetes.json) and fill in the blank strings
5. [generate the template](../acsengine.md#generating-a-template)
6. [deploy the output azuredeploy.json and azuredeploy.parameters.json](../acsengine.md#deployment-usage)

### Common customizations

As part of step 4, edit the [Kubernetes windows example](../../examples/windows/kubernetes.json), you can also make some changes to how Windows is deployed.

#### Changing the OS disk size

The Windows Server deployments default to 30GB for the OS drive (C:), which may not be enough. You can change this size by adding `osDiskSizeGB` under the `agentPoolProfiles`, such as:

```
"agentPoolProfiles": [
      {
        "name": "windowspool2",
        "count": 2,
        "vmSize": "Standard_D2_v3",
        "availabilityProfile": "AvailabilitySet",
        "osType": "Windows",
        "osDiskSizeGB": 127
     }
```

#### Choosing the Windows Server version

If you want to deploy a specific Windows Server version, you can find available versions with `az vm image list --publisher MicrosoftWindowsServer --all -o table`

```
$ az vm image list --publisher MicrosoftWindowsServer --all -o table                                                                                        

Offer                    Publisher                      Sku                                             Urn                                                                                                            Version
-----------------------  -----------------------------  ----------------------------------------------  -------------------------------------------------------------------------------------------------------------  -----------------
...
WindowsServerSemiAnnual  MicrosoftWindowsServer         Datacenter-Core-1709-with-Containers-smalldisk  MicrosoftWindowsServer:WindowsServerSemiAnnual:Datacenter-Core-1709-with-Containers-smalldisk:1709.0.20180412  1709.0.20180412
WindowsServerSemiAnnual  MicrosoftWindowsServer         Datacenter-Core-1803-with-Containers-smalldisk  MicrosoftWindowsServer:WindowsServerSemiAnnual:Datacenter-Core-1803-with-Containers-smalldisk:1803.0.20180504  1803.0.20180504
```

You can use the Offer, Publisher and Sku to pick a specific version by adding `windowsOffer`, `windowsPublisher`, `windowsSku` and (optionally) `widndowsVersion` to the `windowsProfile` section. In this example, the latest Windows Server version 1803 image would be deployed.

```
"windowsProfile": {
            "adminUsername": "azureuser",
            "adminPassword": "...",
            "windowsPublisher": "MicrosoftWindowsServer",
            "windowsOffer": "WindowsServerSemiAnnual",
            "windowsSku": "Datacenter-Core-1803-with-Containers-smalldisk"
     },
```

## Walkthrough

Once your Kubernetes cluster has been created you will have a resource group containing:

1. 1 master accessible by SSH on port 22 or kubectl on port 443

2. a set of windows and linux nodes.  The windows nodes can be accessed through an RDP SSH tunnel via the master node.  To do this, follow these [instructions](../ssh.md#ssh-to-the-machine), replacing port 80 with 3389.  Since your windows machine is already using port 3389, it is recommended to use 3390 to Windows Node 0, 10.240.0.4, 3391 to Windows Node 1, 10.240.0.5, and so on as shown in the following image:

![Image of Windows RDP tunnels](../images/rdptunnels.png)

The following image shows the architecture of a container service cluster with 1 master, and 2 agents:

![Image of Kubernetes cluster on azure with Windows](../images/kubernetes-windows.png)

In the image above, you can see the following parts:

1. **Master Components** - The master runs the Kubernetes scheduler, api server, and controller manager.  Port 443 is exposed for remote management with the kubectl cli.
2. **Linux Nodes** - the Kubernetes nodes run in an availability set.  Azure load balancers are dynamically added to the cluster depending on exposed services.
3. **Windows Nodes** - the Kubernetes windows nodes run in an availability set.
3. **Common Components** - All VMs run a kubelet, Docker, and a Proxy.
4. **Networking** - All VMs are assigned an ip address in the 10.240.0.0/16 network.  Each VM is assigned a /24 subnet for their pod CIDR enabling IP per pod.  The proxy running on each VM implements the service network 10.0.0.0/16.

All VMs are in the same private VNET and are fully accessible to each other.

## Create your First Kubernetes Service

After completing this walkthrough you will know how to:
 * access Kubernetes cluster via SSH,
 * deploy a simple Windows Docker application and expose to the world,
 * and deploy a hybrid Windows / Linux Docker application.
 
1. After successfully deploying the template write down the master FQDN (Fully Qualified Domain Name).
   1. If using Powershell or CLI, the output parameter is in the OutputsString section named 'masterFQDN'
   2. If using Portal, to get the output you need to:
     1. navigate to "resource group"
     2. click on the resource group you just created
     3. then click on "Succeeded" under *last deployment*
     4. then click on the "Microsoft.Template"
     5. now you can copy the output FQDNs and sample SSH commands

   ![Image of docker scaling](../images/portal-kubernetes-outputs.png)

2. SSH to the master FQDN obtained in step 1.

3. Explore your nodes and running pods:
  1. to see a list of your nodes type `kubectl get nodes`.  If you want full detail of the nodes, add `-o yaml` to become `kubectl get nodes -o yaml`.
  2. to see a list of running pods type `kubectl get pods --all-namespaces`.  By default DNS, heapster, and the dashboard pods will be assigned to the Linux nodes.

4. Start your first Docker image by editing a file named `simpleweb.yaml` filling in the contents below, and then apply by typing `kubectl apply -f simpleweb.yaml`.  This will start a windows simple web application and expose to the world.

  ```yaml
  apiVersion: v1
  kind: Service
  metadata:
    name: win-webserver
    labels:
      app: win-webserver
  spec:
    ports:
      # the port that this service should serve on
    - port: 80
      targetPort: 80
    selector:
      app: win-webserver
    type: LoadBalancer
  ---
  apiVersion: extensions/v1beta1
  kind: Deployment
  metadata:
    labels:
      app: win-webserver
    name: win-webserver
  spec:
    replicas: 1
    template:
      metadata:
        labels:
          app: win-webserver
        name: win-webserver
      spec:
        containers:
        - name: windowswebserver
          image: microsoft/windowsservercore:1803
          command:
          - powershell.exe
          - -command
          - "<#code used from https://gist.github.com/wagnerandrade/5424431#> ; $$listener = New-Object System.Net.HttpListener ; $$listener.Prefixes.Add('http://*:80/') ; $$listener.Start() ; $$callerCounts = @{} ; Write-Host('Listening at http://*:80/') ; while ($$listener.IsListening) { ;$$context = $$listener.GetContext() ;$$requestUrl = $$context.Request.Url ;$$clientIP = $$context.Request.RemoteEndPoint.Address ;$$response = $$context.Response ;Write-Host '' ;Write-Host('> {0}' -f $$requestUrl) ;  ;$$count = 1 ;$$k=$$callerCounts.Get_Item($$clientIP) ;if ($$k -ne $$null) { $$count += $$k } ;$$callerCounts.Set_Item($$clientIP, $$count) ;$$header='<html><body><H1>Windows Container Web Server</H1>' ;$$callerCountsString='' ;$$callerCounts.Keys | % { $$callerCountsString+='<p>IP {0} callerCount {1} ' -f $$_,$$callerCounts.Item($$_) } ;$$footer='</body></html>' ;$$content='{0}{1}{2}' -f $$header,$$callerCountsString,$$footer ;Write-Output $$content ;$$buffer = [System.Text.Encoding]::UTF8.GetBytes($$content) ;$$response.ContentLength64 = $$buffer.Length ;$$response.OutputStream.Write($$buffer, 0, $$buffer.Length) ;$$response.Close() ;$$responseStatus = $$response.StatusCode ;Write-Host('< {0}' -f $$responseStatus)  } ; "
        nodeSelector:
          beta.kubernetes.io/os: windows
  ```

5. Type `kubectl get pods -w` to watch the deployment of the service that takes about 30 seconds.  Once running, type `kubectl get svc` and curl the 10.x address to see the output, eg. `curl 10.244.1.4`

6. Type `kubectl get svc -w` to watch the addition of the external IP address that will take about 2-5 minutes.  Once there, you can take the external IP and view in your web browser.

## Example using Azure Files and Azure Disks
### Create Azure File workload
This example is modified after https://github.com/andyzhangx/Demo/tree/master/windows/azurefile/rs3

#### 1. Create an azure file storage class
```kubectl apply -f https://raw.githubusercontent.com/JiangtianLi/Examples/master/windows/azurefile/storageclass-azurefile.yaml```

#### make sure storageclass is created successfully
```
kubectl get storageclass/azurefile -o wide
```

#### 2. Create a pvc for azure file
```kubectl apply -f https://raw.githubusercontent.com/JiangtianLi/Examples/master/windows/azurefile/pvc-azurefile.yaml```

#### make sure pvc is created successfully
```
kubectl get pvc/pvc-azurefile -o wide
```

#### 3. Create a pod with azure file pvc
```kubectl apply -f https://raw.githubusercontent.com/JiangtianLi/Examples/master/windows/azurefile/iis-azurefile.yaml```

#### watch the status of pod until its `STATUS` is `Running`
```
watch kubectl get po/iis-azurefile -o wide
```

#### 4. Enter the pod container to validate
```
kubectl exec -it iis-azurefile -- cmd
```

```
C:\>dir c:\mnt\azure
 Volume in drive C has no label.
 Volume Serial Number is F878-8D74

 Directory of c:\mnt\azure

11/16/2017  09:45 PM    <DIR>          .
11/16/2017  09:45 PM    <DIR>          ..
               0 File(s)              0 bytes
               2 Dir(s)   5,368,709,120 bytes free

```

### Create Azure Disk workload
This example is modified after https://github.com/andyzhangx/Demo/tree/master/windows/azuredisk/rs3

#### 1. Create an azure disk storage class

##### option#1: k8s agent pool is based on blob disk VM
```kubectl apply -f https://raw.githubusercontent.com/JiangtianLi/Examples/master/windows/azuredisk/storageclass-azuredisk.yaml```

##### option#2: k8s agent pool is based on managed disk VM
```kubectl apply -f https://raw.githubusercontent.com/JiangtianLi/Examples/master/windows/azuredisk/storageclass-azuredisk-managed.yaml```

#### make sure storageclass is created successfully
```
kubectl get storageclass/azuredisk -o wide
```

#### 2. Create a pvc for azure disk
```kubectl apply -f https://raw.githubusercontent.com/JiangtianLi/Examples/master/windows/azuredisk/pvc-azuredisk.yaml```

#### make sure pvc is created successfully
```
kubectl get pvc/pvc-azuredisk -o wide
```

#### 3. Create a pod with azure disk pvc
```kubectl apply -f https://raw.githubusercontent.com/JiangtianLi/Examples/master/windows/azuredisk/iis-azuredisk.yaml```

#### watch the status of pod until its `STATUS` is `Running`
```
watch kubectl get po/iis-azuredisk -o wide
```

#### 4. Enter the pod container to validate
```
kubectl exec -it iis-azuredisk -- cmd
```


## Example using multiple containers in a POD
```yaml
apiVersion: extensions/v1beta1
kind: Deployment
metadata:
  labels:
    app: two-containers
  name: two-containers
spec:
  replicas: 1
  template:
    metadata:
      labels:
        app: two-containers
      name: two-containers
    spec:
      volumes:
      - name: shared-data
        emptyDir: {}

      containers:

        - name: iis-container
          image: microsoft/iis:windowsservercore-1803
          volumeMounts:
          - name: shared-data
            mountPath: /wwwcache
          command: 
          - powershell.exe
          - -command 
          - "while ($true) { Start-Sleep -Seconds 10; Copy-Item -Path C:\\wwwcache\\iisstart.htm -Destination C:\\inetpub\\wwwroot\\iisstart.htm; }"            

        - name: servercore-container
          image: microsoft/windowsservercore:1803
          volumeMounts:
          - name: shared-data
            mountPath: /poddata
          command: 
          - powershell.exe
          - -command 
          - "$i=0; while ($true) { Start-Sleep -Seconds 10; $msg = 'Hello from the servercore container, count is {0}' -f $i; Set-Content -Path C:\\poddata\\iisstart.htm -Value $msg; $i++; }"

      nodeSelector:
        beta.kubernetes.io/os: windows
```

## Real-world Workload
TODO


## Windows-specific Troubleshooting

Windows support is still in active development with many changes each week. Read on for more info on known per-version issues and troubleshooting if you run into problems.

### Checking versions

Please be sure to include this info with any Windows bug reports.

Kubernetes
`kubectl version`
-	“Server Version”
`kubectl describe node <windows node>`
-	“kernel version”
-	Also note the IP Address for the next step, but you don't need to share it

Windows config
Connect to the Windows node with remote desktop. This is easiest forwarding a port through SSH from your Kubernetes management endpoint.

1.	`ssh -L 5500:<internal ip>:3389 user@masterFQDN`
2.	Once connected, run `mstsc.exe /v:localhost:5500` to connect. Log in with the username & password you set for the Windows agents.

The Azure CNI plugin version and configuration is stored in `C:\k\azurecni\netconf\10-azure.conflist`. Get
-	mode
-	dns.Nameservers
-	dns.Search

Get the Azure CNI build by running `C:\k\azurecni\bin\azure-vnet.exe --help`. It will dump some errors, but the version such as ` v1.0.4-1-gf0f090e` will be listed.

```
...
2018/05/23 01:28:57 "Start Flag false CniSucceeded false Name CNI Version v1.0.4-1-gf0f090e ErrorMessage required env variables missing vnet []
...
```

### Known Issues per Version


ACS-Engine | Windows Server |	Kubernetes | Azure CNI | Notes
-----------|----------------|------------|-----------|----------
V0.16.2	| Windows Server version 1709 (10.0.16299.____)	| V1.9.7 | ? | DNS resolution is not configured
V0.17.0 | Windows Server version 1709	| V1.10.2 | v1.0.4 | Acs-engine version 0.17 defaults to Windows Server version 1803. You can override it to use 1709 instead [here](#choosing-the-windows-server-version). Manual workarounds needed on Windows for DNS Server list, DNS search suffix
V0.17.0 | Windows Server version 1803 (10.0.17134.1) | V1.10.2 | v1.0.4 | Manual workarounds needed on Windows for DNS Server list, DNS search suffix, and dropped packets
v0.17.1 | Windows Server version 1709 | v1.10.3 | v1.0.4-1-gf0f090e | Manual workarounds needed on Windows for DNS Server list and DNS search suffix. This ACS-Engine version defaults to Windows Server version 1803, but you can override it to use 1709 instead [here](#choosing-the-windows-server-version)
v0.18.3 | Windows Server version 1803 | v1.10.3 | v1.0.6 | Manual workaround needed for DNS search suffix

### Known problems

#### Packets from Windows pods are dropped

Affects: Windows Server version 1803 (10.0.17134.1)

Issues: https://github.com/Azure/acs-engine/issues/3037 

There is a problem with the “L2Tunnel” networking mode not forwarding packets correctly specific to Windows Server version 1803. Windows Server version 1709 is not affected.

Workarounds:
**Fixes are still in development.** A Windows hotfix is needed, and willbe deployed by ACS-Engine once it's ready. The hotfix will be removed later when it's in a future cumulative rollup.


#### Pods cannot resolve public DNS names

Affects: Some builds of Azure CNI

Issues: https://github.com/Azure/azure-container-networking/issues/147

Run `ipconfig /all` in a pod, and check that the first DNS server listed is within your cluster IP range (10.x.x.x). If it's not listed, or not the first in the list, then an azure-cni update is needed.

Workaround:

1.	Get the kube-dns service IP with `kubectl get svc -n kube-system kube-dns`
2.  Cordon & drain the node
3.	Modify `C:\k\azurecni\netconf\10-azure.conflist` and make it the first entry under Nameservers
4.  Uncordon the node

Example:
```
{
    "cniVersion":  "0.3.0",
    "name":  "azure",
    "plugins":  [
                    {
                        "type":  "azure-vnet",
                        "mode":  "tunnel",
                        "bridge":  "azure0",
                        "ipam":  {
                                     "type":  "azure-vnet-ipam"
                                 },
                        "dns":  {
                                    "Nameservers":  [
                                                        "10.0.0.10",
                                                        "168.63.129.16"
                                                    ],
                                    "Search":  [
                                                   "default.svc.cluster.local"
                                               ]
                                },
…
```

#### Pods cannot resolve cluster DNS names

Affects: Azure CNI plugin <= 0.3.0

Issues: https://github.com/Azure/azure-container-networking/issues/146

If you can't resolve internal service names within the same namespace, run `ipconfig /all` in a pod, and check that the DNS Suffix Search List matches the form `<namespace>.svc.cluster.local`. An Azure CNI update is needed to set the right DNS suffix.

Workaround:
1.	Use the FQDN in DNS lookups such as `kubernetes.kube-system.svc.cluster.local`
2.	Instead of DNS, use environment variables `* _SERVICE_HOST` and `*_SERVICE_PORT` to find service IPs and ports in the same namespace


#### Pods cannot ping default route or internet IPs

Affects: All acs-engine deployed clusters

ICMP traffic is not routed between private Azure vNETs or to the internet.

Workaround: test network connections with another protocol (TCP/UDP). For example `Invoke-WebRequest -UseBasicParsing https://www.azure.com` or `curl https://www.azure.com`.



## Cluster Troubleshooting

If your cluster is not reachable, you can run the following command to check for common failures.

### Misconfigured Service Principal

If your Service Principal is misconfigured, none of the Kubernetes components will come up in a healthy manner.
You can check to see if this the problem:

```shell
ssh -i ~/.ssh/id_rsa USER@MASTERFQDN sudo journalctl -u kubelet | grep --text autorest
```

If you see output that looks like the following, then you have **not** configured the Service Principal correctly.
You may need to check to ensure the credentials were provided accurately, and that the configured Service Principal has
read and **write** permissions to the target Subscription.

`Nov 10 16:35:22 k8s-master-43D6F832-0 docker[3177]: E1110 16:35:22.840688    3201 kubelet_node_status.go:69] Unable to construct api.Node object for kubelet: failed to get external ID from cloud provider: autorest#WithErrorUnlessStatusCode: POST https://login.microsoftonline.com/72f988bf-86f1-41af-91ab-2d7cd011db47/oauth2/token?api-version=1.0 failed with 400 Bad Request: StatusCode=400`

## Learning More

Here are recommended links to learn more about Kubernetes:

1. [Kubernetes Bootcamp](https://kubernetesbootcamp.github.io/kubernetes-bootcamp/index.html) - shows you how to deploy, scale, update, and debug containerized applications.
2. [Kubernetes Userguide](http://kubernetes.io/docs/user-guide/) - provides information on running programs in an existing Kubernetes cluster.
3. [Kubernetes Examples](https://github.com/kubernetes/kubernetes/tree/master/examples) - provides a number of examples on how to run real applications with Kubernetes.
