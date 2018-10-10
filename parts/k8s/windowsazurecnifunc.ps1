

# TODO: remove - dead code?
function
Set-VnetPluginMode($mode)
{
    Param(
        [string][Parameter(Mandatory=$true)]
        AzureCNIConfDir
    )
    # Sets Azure VNET CNI plugin operational mode.
    $fileName  = [Io.path]::Combine("$AzureCNIConfDir", "10-azure.conflist")
    (Get-Content $fileName) | %{$_ -replace "`"mode`":.*", "`"mode`": `"$mode`","} | Out-File -encoding ASCII -filepath $fileName
}


function
Install-VnetPlugins
{
    Param(
        [string][Parameter(Mandatory=$true)]
        AzureCNIConfDir,
        [string][Parameter(Mandatory=$true)]
        AzureCNIBinDir
    )
    # Create CNI directories.
    mkdir $AzureCNIBinDir
    mkdir $AzureCNIConfDir

    # Download Azure VNET CNI plugins.
    # Mirror from https://github.com/Azure/azure-container-networking/releases
    $zipfile =  [Io.path]::Combine("$AzureCNIDir", "azure-vnet.zip")
    Invoke-WebRequest -Uri $VNetCNIPluginsURL -OutFile $zipfile
    Expand-Archive -path $zipfile -DestinationPath $AzureCNIBinDir
    del $zipfile

    # Windows does not need a separate CNI loopback plugin because the Windows
    # kernel automatically creates a loopback interface for each network namespace.
    # Copy CNI network config file and set bridge mode.
    move $AzureCNIBinDir/*.conflist $AzureCNIConfDir
}

# TODO: remove - dead code?
function
Set-AzureNetworkPlugin()
{
    # Azure VNET network policy requires tunnel (hairpin) mode because policy is enforced in the host.
    Set-VnetPluginMode "tunnel"
}

function
Set-AzureCNIConfig
{
    Param(
        [string][Parameter(Mandatory=$true)]
        AzureCNIConfDir,
        [string][Parameter(Mandatory=$true)]
        KubeDnsSearchPath,
        [string][Parameter(Mandatory=$true)]
        KubeClusterCIDR,
        [string][Parameter(Mandatory=$true)]
        MasterSubnet,
        [string][Parameter(Mandatory=$true)]
        KubeServiceCIDR
    )
    # Fill in DNS information for kubernetes.
    $fileName  = [Io.path]::Combine("$AzureCNIConfDir", "10-azure.conflist")
    $configJson = Get-Content $fileName | ConvertFrom-Json
    $configJson.plugins.dns.Nameservers[0] = $KubeDnsServiceIp
    $configJson.plugins.dns.Search[0] = $KubeDnsSearchPath
    $configJson.plugins.AdditionalArgs[0].Value.ExceptionList[0] = $KubeClusterCIDR
    $configJson.plugins.AdditionalArgs[0].Value.ExceptionList[1] = $MasterSubnet
    $configJson.plugins.AdditionalArgs[1].Value.DestinationPrefix  = $KubeServiceCIDR

    $configJson | ConvertTo-Json -depth 20 | Out-File -encoding ASCII -filepath $fileName
}
