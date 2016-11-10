<#
    .SYNOPSIS
        Provisions VM as a Kubernetes agent.

    .DESCRIPTION
        Provisions VM as a Kubernetes agent.

    .PARAMETER MasterIP
        IP Address of Docker Swarm Master

    .EXAMPLE
        .\kuberneteswindowssetup.ps1 -MasterIP 192.168.255.5

#>
[CmdletBinding(DefaultParameterSetName="Standard")]
param(
    [string]
    [ValidateNotNullOrEmpty()]
    $MasterIP = "10.240.255.5",

    [parameter()]
    [ValidateNotNullOrEmpty()]
    $KubeDnsServiceIp = "10.0.0.10",

    [parameter()]
    [ValidateNotNullOrEmpty()]
    $CACertificate = "<<<caCertificate>>>",

    [parameter(Mandatory=$true)]
    [ValidateNotNullOrEmpty()]
    $MasterFQDNPrefix,

    [parameter(Mandatory=$true)]
    [ValidateNotNullOrEmpty()]
    $Location,

    [parameter()]
    [ValidateNotNullOrEmpty()]
    $AgentCertificate = "<<<clientCertificate>>>",

    [parameter(Mandatory=$true)]
    [ValidateNotNullOrEmpty()]
    $AgentKey,

    [parameter(Mandatory=$true)]
    [ValidateNotNullOrEmpty()]
    $AzureHostname
)

$global:DockerServiceName = "Docker"
$global:RRASServiceName = "RemoteAccess"
$global:KubeDir = "c:\k"
$global:KubeBinariesSASURL = "https://acsengine.blob.core.windows.net/windows/k.zip?st=2016-11-08T02%3A27%3A00Z&se=2020-11-09T02%3A27%3A00Z&sp=rl&sv=2015-12-11&sr=b&sig=8vloXDY9rG1XCGf4gqreBbI%2Bj1IWkxWiDZQSZM5S9mY%3D"
$global:KubeletStartFile = $global:KubeDir + "\kubeletstart.ps1"
$global:KubeProxyStartFile = $global:KubeDir + "\kubeproxystart.ps1"

filter Timestamp {"$(Get-Date -Format o): $_"}

function Write-Log($message)
{
    $msg = $message | Timestamp
    Write-Output $msg
}

function Expand-ZIPFile($file, $destination)
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
Write-KubeConfig()
{
    $kubeConfigFile = $global:KubeDir + "\config"

    $kubeConfig = @"
---
apiVersion: v1
clusters:
- cluster:
    certificate-authority-data: "$CACertificate"
    server: https://$MasterFQDNPrefix.$Location.cloudapp.azure.com
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
    client-certificate-data: "$AgentCertificate"
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
Write-KubeletStartFile
{
    $argList = @("--hostname-override=$AzureHostname","--pod-infra-container-image=kubletwin/pause","--resolv-conf=""""","--api-servers=https://${MasterIP}:443","--kubeconfig=c:\k\config")
    $process = Start-Process -FilePath c:\k\kubelet.exe -PassThru -ArgumentList $argList

    $podCidrDiscovered=$false
    $podCIDR=""
    # run kubelet until podCidr is discovered
    while (-not $podCidrDiscovered)
    {
        $podCIDR=c:\k\kubectl.exe --kubeconfig=c:\k\config get nodes/$AzureHostname -o custom-columns=podCidr:.spec.podCIDR --no-headers

        if ($podCIDR.length -gt 0)
        {
            Write-Log "'$podCIDR' found"
            $podCidrDiscovered=$true
        }
        else
        {
            Write-Log "sleeping for 10 seconds..."
            Start-Sleep -sec 10    
        }
    }

    # stop the kubelet process now that we have our CIDR
    $process | Stop-Process

    $kubeConfig = @"
`$netResult=docker network ls | findstr podnetwork
if (`$netResult.length -eq 0)
{
    New-ContainerNetwork -Name podnetwork -Mode L2Bridge -SubnetPrefix $podCIDR -GatewayAddress 10.240.0.1
    Restart-Service $global:DockerServiceName
}
`$env:CONTAINER_NETWORK="podnetwork"
c:\k\kubelet.exe --hostname-override=$AzureHostname --pod-infra-container-image=kubletwin/pause --resolv-conf="" --allow-privileged=true --enable-debugging-handlers --api-servers=https://${MasterIP}:443 --cluster-dns=$KubeDnsServiceIp --cluster-domain=cluster.local  --kubeconfig=c:\k\config --hairpin-mode=promiscuous-bridge --v=2
"@
    $kubeConfig | Out-File -encoding ASCII -filepath $global:KubeletStartFile

    $kubeProxyStartStr = @"
`$nodeIP=""
`$aliasName="vEthernet (HNS Internal NIC)"
while (`$true)
{
    try
    {
        `$nodeNic=Get-NetIPaddress -InterfaceAlias `$aliasName -AddressFamily IPv4
        #bind to the docker IP address
        `$nodeIP=`$nodeNic.IPAddress | Where-Object {`$_.StartsWith("172.")} | Select-Object -First 1
        break
    }
    catch
    {
        Write-Output "sleeping for 10s since `$aliasName is not defined"
        Start-Sleep -sec 10
    }
}

`$env:INTERFACE_TO_ADD_SERVICE_IP=`$aliasName
c:\k\kube-proxy.exe --v=3 --proxy-mode=userspace --hostname-override=$AzureHostname --master=${MasterIP}:8080 --bind-address=`$nodeIP --kubeconfig=c:\k\config
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
}

try
{
    if ($true) {
        Write-Log "Provisioning $global:DockerServiceName... with IP $MasterIP"

        Write-Log "download kubelet binaries and unzip"
        Get-KubeBinaries

        Write-Log "Write kube config"
        Write-KubeConfig

        Write-Log "Create the Pause Container kubletwin/pause"
        New-InfraContainer

        Write-Log "write kubelet startfile"
        Write-KubeletStartFile 

        Write-Log "install the NSSM service"
        New-NSSMService

        Write-Log "Turn off Firewall to enable pods to talk to service endpoints. (Kubelet should eventually do this)"
        netsh advfirewall set allprofiles state off

        Write-Log "Install hyperv to expose vfpext"
        dism /Online /Enable-Feature /FeatureName:Microsoft-Hyper-V /All /NoRestart
        
        Write-Log "Setup Complete"
        Restart-Computer -Force
    }
    else 
    {
        # keep for debugging purposes
        Write-Log "kuberneteswindowssetup.ps1 -MasterIP $MasterIP -KubeDnsServiceIp $KubeDnsServiceIp -MasterFQDNPrefix $MasterFQDNPrefix -Location $Location -CACertificate $CACertificate -AgentCertificate $AgentCertificate -AgentKey $AgentKey -AzureHostname $AzureHostname"
    }
}
catch
{
    Write-Error $_
}