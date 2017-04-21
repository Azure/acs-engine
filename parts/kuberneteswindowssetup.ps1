<#
    .SYNOPSIS
        Provisions VM as a Kubernetes agent.

    .DESCRIPTION
        Provisions VM as a Kubernetes agent.
#>
[CmdletBinding(DefaultParameterSetName="Standard")]
param(
    [string]
    [ValidateNotNullOrEmpty()]
    $MasterIP,

    [parameter()]
    [ValidateNotNullOrEmpty()]
    $KubeDnsServiceIp,

    [parameter(Mandatory=$true)]
    [ValidateNotNullOrEmpty()]
    $MasterFQDNPrefix,

    [parameter(Mandatory=$true)]
    [ValidateNotNullOrEmpty()]
    $Location,

    [parameter(Mandatory=$true)]
    [ValidateNotNullOrEmpty()]
    $AgentKey,

    [parameter(Mandatory=$true)]
    [ValidateNotNullOrEmpty()]
    $AzureHostname,

    [parameter(Mandatory=$true)]
    [ValidateNotNullOrEmpty()]
    $AADClientId,

    [parameter(Mandatory=$true)]
    [ValidateNotNullOrEmpty()]
    $AADClientSecret
)

$global:CACertificate = "{{{caCertificate}}}"
$global:AgentCertificate = "{{{clientCertificate}}}"
$global:DockerServiceName = "Docker"
$global:RRASServiceName = "RemoteAccess"
$global:KubeDir = "c:\k"
$global:KubeBinariesSASURL = "{{{kubeBinariesSASURL}}}"
$global:KubeBinariesVersion = "{{{kubeBinariesVersion}}}"
$global:KubeletStartFile = $global:KubeDir + "\kubeletstart.ps1"
$global:KubeProxyStartFile = $global:KubeDir + "\kubeproxystart.ps1"
$global:NatNetworkName="nat"
$global:TransparentNetworkName="transparentNet"

$global:TenantId = "{{{tenantID}}}"
$global:SubscriptionId = "{{{subscriptionId}}}"
$global:ResourceGroup = "{{{resourceGroup}}}"
$global:SubnetName = "{{{subnetName}}}"
$global:SecurityGroupName = "{{{nsgName}}}"
$global:VNetName = "{{{virtualNetworkName}}}"
$global:RouteTableName = "{{{routeTableName}}}"
$global:PrimaryAvailabilitySetName = "{{{primaryAvailablitySetName}}}"

filter Timestamp {"$(Get-Date -Format o): $_"}

function
Write-Log($message)
{
    $msg = $message | Timestamp
    Write-Output $msg
}

function
Expand-ZIPFile($file, $destination)
{
    $shell = new-object -com shell.application
    $zip = $shell.NameSpace($file)
    foreach($item in $zip.items())
    {
        $shell.Namespace($destination).copyhere($item)
    }
}

function
Get-KubeBinaries()
{
    $zipfile = "c:\k.zip"
    Invoke-WebRequest -Uri $global:KubeBinariesSASURL -OutFile $zipfile
    Expand-ZIPFile -File $zipfile -Destination C:\
}

function
Write-AzureConfig()
{
    $azureConfigFile = $global:KubeDir + "\azure.json"

    $azureConfig = @"
{
    "tenantId": "$global:TenantId",
    "subscriptionId": "$global:SubscriptionId",
    "aadClientId": "$AADClientId",
    "aadClientSecret": "$AADClientSecret",
    "resourceGroup": "$global:ResourceGroup",
    "location": "$Location",
    "subnetName": "$global:SubnetName",
    "securityGroupName": "$global:SecurityGroupName",
    "vnetName": "$global:VNetName",
    "routeTableName": "$global:RouteTableName",
    "primaryAvailabilitySetName": "$global:PrimaryAvailabilitySetName"
}
"@

    $azureConfig | Out-File -encoding ASCII -filepath "$azureConfigFile"    
}

function
Write-KubeConfig()
{
    $kubeConfigFile = $global:KubeDir + "\config"

    $kubeConfig = @"
---
apiVersion: v1
clusters:
- cluster:
    certificate-authority-data: "$global:CACertificate"
    server: https://${MasterIP}:443
  name: "$MasterFQDNPrefix"
contexts:
- context:
    cluster: "$MasterFQDNPrefix"
    user: "$MasterFQDNPrefix-admin"
  name: "$MasterFQDNPrefix"
current-context: "$MasterFQDNPrefix"
kind: Config
users:
- name: "$MasterFQDNPrefix-admin"
  user:
    client-certificate-data: "$global:AgentCertificate"
    client-key-data: "$AgentKey"
"@

    $kubeConfig | Out-File -encoding ASCII -filepath "$kubeConfigFile"    
}

function
New-InfraContainer()
{
    cd $global:KubeDir
    docker build -t kubletwin/pause . 
}

function
Write-KubernetesStartFiles($podCIDR)
{
    $kubeStartStr = @"
`$global:TransparentNetworkName="$global:TransparentNetworkName"
`$global:AzureHostname="$AzureHostname"
`$global:MasterIP="$MasterIP"
`$global:NatNetworkName="$global:NatNetworkName"
`$global:KubeDnsServiceIp="$KubeDnsServiceIp"
`$global:KubeBinariesVersion="$global:KubeBinariesVersion"

function
Get-PodGateway(`$podCIDR)
{
    return `$podCIDR.substring(0,`$podCIDR.lastIndexOf(".")) + ".1"
}

function
Set-DockerNetwork(`$podCIDR)
{
    # Turn off Firewall to enable pods to talk to service endpoints. (Kubelet should eventually do this)
    netsh advfirewall set allprofiles state off

    `$dockerTransparentNet=docker network ls --quiet --filter "NAME=`$global:TransparentNetworkName"
    if (`$dockerTransparentNet.length -eq 0)
    {
        `$podGW=Get-PodGateway(`$podCIDR)

        # create new transparent network
        docker network create --driver=transparent --subnet=`$podCIDR --gateway=`$podGW `$global:TransparentNetworkName

        # create host vnic for gateway ip to forward the traffic and kubeproxy to listen over VIP
        Add-VMNetworkAdapter -ManagementOS -Name forwarder -SwitchName "Layered Ethernet 3"

        # Assign gateway IP to new adapter and enable forwarding on host adapters:
        netsh interface ipv4 add address "vEthernet (forwarder)" `$podGW 255.255.255.0
        netsh interface ipv4 set interface "vEthernet (forwarder)" for=en
        netsh interface ipv4 set interface "vEthernet (HNSTransparent)" for=en
    }
}

function
Get-PodCIDR()
{
    `$podCIDR=c:\k\kubectl.exe --kubeconfig=c:\k\config get nodes/`$(`$global:AzureHostname.ToLower()) -o custom-columns=podCidr:.spec.podCIDR --no-headers
    return `$podCIDR
}

function
Test-PodCIDR(`$podCIDR)
{
    return `$podCIDR.length -gt 0
}

try
{
    `$podCIDR=Get-PodCIDR
    `$podCidrDiscovered=Test-PodCIDR(`$podCIDR)

    # if the podCIDR has not yet been assigned to this node, start the kubelet process to get the podCIDR, and then promptly kill it.
    if (-not `$podCidrDiscovered)
    {
        `$argList = @("--hostname-override=`$global:AzureHostname","--pod-infra-container-image=kubletwin/pause","--resolv-conf=""""","--api-servers=https://`${global:MasterIP}:443","--kubeconfig=c:\k\config","--enable-cri=false")

        if (`$global:KubeBinariesVersion -eq "1.5.3")
        {
            `$argList = @("--hostname-override=`$global:AzureHostname","--pod-infra-container-image=kubletwin/pause","--resolv-conf=""""","--api-servers=https://`${global:MasterIP}:443","--kubeconfig=c:\k\config")
        }

        `$process = Start-Process -FilePath c:\k\kubelet.exe -PassThru -ArgumentList `$argList

        # run kubelet until podCidr is discovered
        Write-Host "waiting to discover pod CIDR"
        while (-not `$podCidrDiscovered)
        {
            Write-Host "Sleeping for 10s, and then waiting to discover pod CIDR"
            Start-Sleep -sec 10
            
            `$podCIDR=Get-PodCIDR
            `$podCidrDiscovered=Test-PodCIDR(`$podCIDR)
        }
    
        # stop the kubelet process now that we have our CIDR, discard the process output
        `$process | Stop-Process | Out-Null
    }
    
    Set-DockerNetwork(`$podCIDR)

    # startup the service
    `$podGW=Get-PodGateway(`$podCIDR)
    `$env:CONTAINER_NETWORK="`$global:TransparentNetworkName"
    `$env:NAT_NETWORK="`$global:NatNetworkName"
    `$env:POD_GW="`$podGW"
    `$env:VIP_CIDR="10.0.0.0/8"
    
    if (`$global:KubeBinariesVersion -eq "1.5.3")
    {
        c:\k\kubelet.exe --hostname-override=`$global:AzureHostname --pod-infra-container-image=kubletwin/pause --resolv-conf="" --allow-privileged=true --enable-debugging-handlers --api-servers=https://`${global:MasterIP}:443 --cluster-dns=`$global:KubeDnsServiceIp --cluster-domain=cluster.local  --kubeconfig=c:\k\config --hairpin-mode=promiscuous-bridge --v=2 --azure-container-registry-config=c:\k\azure.json
    }
    else
    {
        c:\k\kubelet.exe --hostname-override=`$global:AzureHostname --pod-infra-container-image=kubletwin/pause --resolv-conf="" --allow-privileged=true --enable-debugging-handlers --api-servers=https://`${global:MasterIP}:443 --cluster-dns=`$global:KubeDnsServiceIp --cluster-domain=cluster.local  --kubeconfig=c:\k\config --hairpin-mode=promiscuous-bridge --v=2 --azure-container-registry-config=c:\k\azure.json --enable-cri=false --image-pull-progress-deadline=20m
    }
}
catch
{
    Write-Error `$_
}
"@
    $kubeStartStr | Out-File -encoding ASCII -filepath $global:KubeletStartFile

    $kubeProxyStartStr = @"
`$env:INTERFACE_TO_ADD_SERVICE_IP="vEthernet (forwarder)"
c:\k\kube-proxy.exe --v=3 --proxy-mode=userspace --hostname-override=$AzureHostname --kubeconfig=c:\k\config
"@

    $kubeProxyStartStr | Out-File -encoding ASCII -filepath $global:KubeProxyStartFile
}

function
New-NSSMService
{
    # setup kubelet
    c:\k\nssm install Kubelet C:\Windows\System32\WindowsPowerShell\v1.0\powershell.exe
    c:\k\nssm set Kubelet AppDirectory $global:KubeDir
    c:\k\nssm set Kubelet AppParameters $global:KubeletStartFile
    c:\k\nssm set Kubelet DisplayName Kubelet
    c:\k\nssm set Kubelet Description Kubelet
    c:\k\nssm set Kubelet Start SERVICE_AUTO_START
    c:\k\nssm set Kubelet ObjectName LocalSystem
    c:\k\nssm set Kubelet Type SERVICE_WIN32_OWN_PROCESS
    c:\k\nssm set Kubelet AppThrottle 1500
    c:\k\nssm set Kubelet AppStdout C:\k\kubelet.log
    c:\k\nssm set Kubelet AppStderr C:\k\kubelet.err.log
    c:\k\nssm set Kubelet AppStdoutCreationDisposition 4
    c:\k\nssm set Kubelet AppStderrCreationDisposition 4
    c:\k\nssm set Kubelet AppRotateFiles 1
    c:\k\nssm set Kubelet AppRotateOnline 1
    c:\k\nssm set Kubelet AppRotateSeconds 86400
    c:\k\nssm set Kubelet AppRotateBytes 1048576
    net start Kubelet
    
    # setup kubeproxy
    c:\k\nssm install Kubeproxy C:\Windows\System32\WindowsPowerShell\v1.0\powershell.exe
    c:\k\nssm set Kubeproxy AppDirectory $global:KubeDir
    c:\k\nssm set Kubeproxy AppParameters $global:KubeProxyStartFile
    c:\k\nssm set Kubeproxy DisplayName Kubeproxy
    c:\k\nssm set Kubeproxy DependOnService Kubelet
    c:\k\nssm set Kubeproxy Description Kubeproxy
    c:\k\nssm set Kubeproxy Start SERVICE_AUTO_START
    c:\k\nssm set Kubeproxy ObjectName LocalSystem
    c:\k\nssm set Kubeproxy Type SERVICE_WIN32_OWN_PROCESS
    c:\k\nssm set Kubeproxy AppThrottle 1500
    c:\k\nssm set Kubeproxy AppStdout C:\k\kubeproxy.log
    c:\k\nssm set Kubeproxy AppStderr C:\k\kubeproxy.err.log
    c:\k\nssm set Kubeproxy AppRotateFiles 1
    c:\k\nssm set Kubeproxy AppRotateOnline 1
    c:\k\nssm set Kubeproxy AppRotateSeconds 86400
    c:\k\nssm set Kubeproxy AppRotateBytes 1048576
    net start Kubeproxy
}

function
Set-Explorer
{
    # setup explorer so that it is usable
    New-Item -Path HKLM:"\\SOFTWARE\\Policies\\Microsoft\\Internet Explorer"
    New-Item -Path HKLM:"\\SOFTWARE\\Policies\\Microsoft\\Internet Explorer\\BrowserEmulation"
    New-ItemProperty -Path HKLM:"\\SOFTWARE\\Policies\\Microsoft\\Internet Explorer\\BrowserEmulation" -Name IntranetCompatibilityMode -Value 0 -Type DWord
    New-Item -Path HKLM:"\\SOFTWARE\\Policies\\Microsoft\\Internet Explorer\\Main"
    New-ItemProperty -Path HKLM:"\\SOFTWARE\\Policies\\Microsoft\\Internet Explorer\\Main" -Name "Start Page" -Type String -Value http://bing.com
}

try
{
    # Set to false for debugging.  This will output the start script to
    # c:\AzureData\CustomDataSetupScript.log, and then you can RDP 
    # to the windows machine, and run the script manually to watch
    # the output.
    if ($true) {
        Write-Log "Provisioning $global:DockerServiceName... with IP $MasterIP"

        Write-Log "download kubelet binaries and unzip"
        Get-KubeBinaries

        Write-Log "Write azure config"
        Write-AzureConfig

        Write-Log "Write kube config"
        Write-KubeConfig

        Write-Log "Create the Pause Container kubletwin/pause"
        New-InfraContainer

        Write-Log "write kubelet startfile with pod CIDR of $podCIDR"
        Write-KubernetesStartFiles $podCIDR

        Write-Log "install the NSSM service"
        New-NSSMService

        Write-Log "Set Internet Explorer"
        Set-Explorer

        Write-Log "Setup Complete"
    }
    else 
    {
        # keep for debugging purposes
        Write-Log ".\CustomDataSetupScript.ps1 -MasterIP $MasterIP -KubeDnsServiceIp $KubeDnsServiceIp -MasterFQDNPrefix $MasterFQDNPrefix -Location $Location -AgentKey $AgentKey -AzureHostname $AzureHostname -AADClientId $AADClientId -AADClientSecret $AADClientSecret"
    }
}
catch
{
    Write-Error $_
}