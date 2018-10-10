$global:KubeDir = "c:\k"
$global:VolumePluginDir = [Io.path]::Combine("$global:KubeDir", "volumeplugins")
$global:KubeletStartFile = [io.path]::Combine($global:KubeDir, "kubeletstart.ps1")
$global:KubeProxyStartFile = [io.path]::Combine($global:KubeDir, "kubeproxystart.ps1")

function
Write-AzureConfig
{
    Param(
        
        [string][Parameter(Mandatory=$true)]
        AADClientId,
        [string][Parameter(Mandatory=$true)]
        AADClientSecret,
        [string][Parameter(Mandatory=$true)]
        TenantId,
        [string][Parameter(Mandatory=$true)]
        SubscriptionId,
        [string][Parameter(Mandatory=$true)]
        ResourceGroup,
        [string][Parameter(Mandatory=$true)]
        Location,
        [string][Parameter(Mandatory=$true)]
        VmType,
        [string][Parameter(Mandatory=$true)]
        SubnetName,
        [string][Parameter(Mandatory=$true)]
        SecurityGroupName,
        [string][Parameter(Mandatory=$true)]
        VNetName,
        [string][Parameter(Mandatory=$true)]
        RouteTableName,
        [string][Parameter(Mandatory=$true)]
        PrimaryAvailabilitySetName,
        [string][Parameter(Mandatory=$true)]
        PrimaryScaleSetName,
        [string][Parameter(Mandatory=$true)]
        UseManagedIdentityExtension,
        [string][Parameter(Mandatory=$true)]
        UserAssignedClientID,
        [string][Parameter(Mandatory=$true)]
        UseInstanceMetadata,
        [string][Parameter(Mandatory=$true)]
        LoadBalancerSku,
        [string][Parameter(Mandatory=$true)]
        ExcludeMasterFromStandardLB
    )
    $azureConfigFile = [io.path]::Combine($global:KubeDir, "azure.json")

    $azureConfig = @"
{
    "tenantId": "$TenantId",
    "subscriptionId": "$global:SubscriptionId",
    "aadClientId": "$AADClientId",
    "aadClientSecret": "$AADClientSecret",
    "resourceGroup": "$global:ResourceGroup",
    "location": "$Location",
    "vmType": "$global:VmType",
    "subnetName": "$global:SubnetName",
    "securityGroupName": "$global:SecurityGroupName",
    "vnetName": "$global:VNetName",
    "routeTableName": "$global:RouteTableName",
    "primaryAvailabilitySetName": "$global:PrimaryAvailabilitySetName",
    "primaryScaleSetName": "$global:PrimaryScaleSetName",
    "useManagedIdentityExtension": $global:UseManagedIdentityExtension,
    "userAssignedIdentityID": $global:UserAssignedClientID,
    "useInstanceMetadata": $global:UseInstanceMetadata,
    "loadBalancerSku": "$global:LoadBalancerSku",
    "excludeMasterFromStandardLB": $global:ExcludeMasterFromStandardLB
}
"@

    $azureConfig | Out-File -encoding ASCII -filepath "$azureConfigFile"
}


function
Write-CACert
{
    Param(
        [string][Parameter(Mandatory=$true)]
        CACertificate
    )
    $caFile = [io.path]::Combine($global:KubeDir, "ca.crt")
    [System.Text.Encoding]::ASCII.GetString([System.Convert]::FromBase64String($CACertificate)) | Out-File -Encoding ascii $caFile
}

function
Write-KubeConfig
{
    Param(
        [string][Parameter(Mandatory=$true)]
        CACertificate,
        [string][Parameter(Mandatory=$true)]
        MasterFQDNPrefix,
        [string][Parameter(Mandatory=$true)]
        MasterIP,
        [string][Parameter(Mandatory=$true)]
        AgentKey,
        [string][Parameter(Mandatory=$true)]
        AgentCertificate
    )
    $kubeConfigFile = [io.path]::Combine($global:KubeDir, "config")

    $kubeConfig = @"
---
apiVersion: v1
clusters:
- cluster:
    certificate-authority-data: "$CACertificate"
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
    client-certificate-data: "$AgentCertificate"
    client-key-data: "$AgentKey"
"@

    $kubeConfig | Out-File -encoding ASCII -filepath "$kubeConfigFile"
}

function
New-InfraContainer
{
    cd $global:KubeDir
    $computerInfo = Get-ComputerInfo
    $windowsBase = if ($computerInfo.WindowsVersion -eq "1709") {
        "microsoft/nanoserver:1709"
    } elseif ($computerInfo.WindowsVersion -eq "1803") {
        "microsoft/nanoserver:1803"
    } elseif ($computerInfo.WindowsVersion -eq "1809") {
        # TODO: unsure if 2019 will report 1809 or not
        "microsoft/nanoserver:1809"
    } else {
        "mcr.microsoft.com/nanoserver-insider"
    }

    "FROM $($windowsBase)" | Out-File -encoding ascii -FilePath Dockerfile
    "CMD cmd /c ping -t localhost" | Out-File -encoding ascii -FilePath Dockerfile -Append
    docker build -t kubletwin/pause .
}


function
Get-KubeBinaries
{
    Param(
        # TODO: Deprecate this and replace with methods that get individual components instead of zip containing everything
        [string]
        KubeBinariesSASURL
    )
    
    $zipfile = "c:\k.zip"
    for ($i=0; $i -le 10; $i++)
    {
        DownloadFileOverHttp -Source $KubeBinariesSASURL -DestinationPath $zipfile
        if ($?) {
            break
        } else {
            Write-Log $Error[0].Exception.Message
        }
    }
    Expand-Archive -path $zipfile -DestinationPath C:\
}
