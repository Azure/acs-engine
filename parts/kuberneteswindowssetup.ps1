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

    [parameter(Mandatory=$true)]
    [ValidateNotNullOrEmpty()]
    $CACertificate,

    [parameter(Mandatory=$true)]
    [ValidateNotNullOrEmpty()]
    $MasterFQDN,

    [parameter(Mandatory=$true)]
    [ValidateNotNullOrEmpty()]
    $Location,

    [parameter(Mandatory=$true)]
    [ValidateNotNullOrEmpty()]
    $AgentCertificate,

    [parameter(Mandatory=$true)]
    [ValidateNotNullOrEmpty()]
    $AgentKey,

    [parameter(Mandatory=$true)]
    [ValidateNotNullOrEmpty()]
    $AzureHostname,

    [parameter(Mandatory=$true)]
    [ValidateNotNullOrEmpty()]
    $SecondaryNICIP,

    [parameter(Mandatory=$true)]
    [ValidateNotNullOrEmpty()]
    $PODCIDRSubnet
)

$global:DockerServiceName = "Docker"
$global:RRASServiceName = "RemoteAccess"
$global:KubeDir = "c:\k"
$global:KubeBinariesSASURL = "https://acsengine.blob.core.windows.net/windows/k.zip?st=2016-11-04T15%3A36%3A00Z&se=2020-10-05T15%3A36%3A00Z&sp=rl&sv=2015-12-11&sr=b&sig=I46lfcEGqKm5tlj1tisjeb4vlijB%2FD1qqBRDsHXA658%3D"

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
    server: https://$MasterFQDN.$Location.cloudapp.azure.com
  name: "$MasterFQDN"
contexts:
- context:
    cluster: "$MasterFQDN"
    user: "$MasterFQDN-admin"
  name: "$MasterFQDN"
current-context: "$MasterFQDN"
kind: Config
users:
- name: "$MasterFQDN-admin"
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
    $kubeletStartFile = $global:KubeDir + "\kubeletstart.ps1"

    $kubeConfig = @"
docker network rm podnetwork
docker network create -d transparent --gateway $SecondaryNICIP --subnet $PODCIDRSubnet podnetwork
SET CONTAINER_NETWORK=podnetwork
.\kubelet.exe --hostname-override=$AzureHostname --pod-infra-container-image="kubletwin/pause" --resolv-conf="" --api-servers=https://${MasterIP}:443 --kubeconfig=c:\k\config
"@

    $kubeConfig | Out-File -encoding ASCII -filepath "$kubeletStartFile"
}

try
{
    Write-Log "Provisioning $global:DockerServiceName... with IP $MasterIP"

    Write-Log "Enable RRAS"
    New-Item -ItemType directory -Path $global:KubeDir

    #Write-Log "Enable and start RRAS"
    #Set-Service -Name $global:RRASServiceName -StartupType automatic -Status Running

    Write-Log "download kubelet binaries and unzip"
    Get-KubeBinaries

    Write-Log "Write kube config"
    Write-KubeConfig

    Write-Log "Create the Pause Container kubletwin/pause"
    New-InfraContainer

    Write-Log "write kubelet startfile"
    Write-KubeletStartFile
    
    Write-Log "Setup Complete"
}
catch
{
    Write-Error $_
}