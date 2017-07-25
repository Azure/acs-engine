<#
    .SYNOPSIS
        Provisions VM as a DCOS agent.

    .DESCRIPTION
        Provisions VM as a DCOS agent.
#>
[CmdletBinding(DefaultParameterSetName="Standard")]
param(
    [string]
    [ValidateNotNullOrEmpty()]
    $MasterIP,

    [parameter()]
    [ValidateNotNullOrEmpty()]
    $DnsServiceIp,

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

$global:CACertificate = "{{WrapAsVariable "caCertificate"}}"
$global:AgentCertificate = "{{WrapAsVariable "clientCertificate"}}"
$global:DockerServiceName = "Docker"
$global:RRASServiceName = "RemoteAccess"
$global:DcosDir = "c:\dcos"
$global:DcosBinariesSASURL = "{{WrapAsVariable "dcosBinariesSASURL"}}"
$global:DcosBinariesVersion = "{{WrapAsVariable "dcosBinariesVersion"}}"
$global:DCOSStartFile = $global:DcosDir + "\dcosletstart.ps1"
$global:DcosProxyStartFile = $global:DcosDir + "\dcosproxystart.ps1"
$global:NatNetworkName="nat"
$global:TransparentNetworkName="transparentNet"

$global:TenantId = "{{WrapAsVariable "tenantID"}}"
$global:SubscriptionId = "{{WrapAsVariable "subscriptionId"}}"
$global:ResourceGroup = "{{WrapAsVariable "resourceGroup"}}"
$global:SubnetName = "{{WrapAsVariable "subnetName"}}"
$global:SecurityGroupName = "{{WrapAsVariable "nsgName"}}"
$global:VNetName = "{{WrapAsVariable "virtualNetworkName"}}"
$global:RouteTableName = "{{WrapAsVariable "routeTableName"}}"
$global:PrimaryAvailabilitySetName = "{{WrapAsVariable "primaryAvailablitySetName"}}"
$global:NeedPatchWinNAT = $false

$global:UseManagedIdentityExtension = "{{WrapAsVariable "useManagedIdentityExtension"}}"
$global:UseInstanceMetadata = "{{WrapAsVariable "useInstanceMetadata"}}"

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
Get-DcosBinaries()
{
    $zipfile = "c:\k.zip"
    Invoke-WebRequest -Uri $global:DcosBinariesSASURL -OutFile $zipfile
    Expand-ZIPFile -File $zipfile -Destination C:\
}

function
Patch-WinNATBinary()
{
    $winnatcurr = $global:DcosDir + "\winnat.sys"
    if (Test-Path $winnatcurr)
    {
        $global:NeedPatchWinNAT = $true
        $winnatsys = "$env:SystemRoot\System32\drivers\winnat.sys"
        Stop-Service winnat
        takeown /f $winnatsys
        icacls $winnatsys /grant "Administrators:(F)"    
        Copy-Item $winnatcurr $winnatsys
        bcdedit /set TESTSIGNING on
    }
}

function
Write-AzureConfig()
{
    $azureConfigFile = $global:DcosDir + "\azure.json"

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
    "primaryAvailabilitySetName": "$global:PrimaryAvailabilitySetName",
    "useManagedIdentityExtension": $global:UseManagedIdentityExtension,
    "useInstanceMetadata": $global:UseInstanceMetadata
}
"@

    $azureConfig | Out-File -encoding ASCII -filepath "$azureConfigFile"    
}

function
Write-DcosConfig()
{
    $dcosConfigFile = $global:DcosDir + "\config"

    $dcosConfig = @"
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

    $dcosConfig | Out-File -encoding ASCII -filepath "$dcosConfigFile"    
}

function
New-InfraContainer()
{
    cd $global:DcosDir
    docker build -t dcoswin/pause . 
}

function
Write-DCOSStartFiles($podCIDR)
{
    $DCOSArgList = @("--hostname-override=`$global:AzureHostname","--pod-infra-container-image=dcoswin/pause","--resolv-conf=""""""""","--api-servers=https://`${global:MasterIP}:443","--dcosconfig=c:\dcos\config")
    $DCOSCommandLine = @"
c:\dcos\dcoslet.exe --hostname-override=`$global:AzureHostname --pod-infra-container-image=dcoswin/pause --resolv-conf="" --allow-privileged=true --enable-debugging-handlers --api-servers=https://`${global:MasterIP}:443 --cluster-dns=`$global:DcosDnsServiceIp --cluster-domain=cluster.local  --dcosconfig=c:\dcos\config --hairpin-mode=promiscuous-bridge --v=2 --azure-container-registry-config=c:\dcos\azure.json --runtime-request-timeout=10m
"@

    if ($global:DcosBinariesVersion -ge "1.6.0")
    {
        # stop using container runtime interface from 1.6.0+ (officially deprecated from 1.7.0)
        if ($global:DcosBinariesVersion -lt "1.7.0")
        {
            $DCOSArgList += "--enable-cri=false"
            $DCOSCommandLine += " --enable-cri=false"
        }
        # more time is needed to pull windows server images (flag supported from 1.6.0)
        $DCOSCommandLine += " --image-pull-progress-deadline=20m --cgroups-per-qos=false --enforce-node-allocatable=`"`""
    }
    $DCOSArgListStr = "`"" + ($DCOSArgList -join "`",`"") + "`""

    $DCOSArgListStr = "@`($DCOSArgListStr`)"

    $dcosStartStr = @"
`$global:TransparentNetworkName="$global:TransparentNetworkName"
`$global:AzureHostname="$AzureHostname"
`$global:MasterIP="$MasterIP"
`$global:NatNetworkName="$global:NatNetworkName"
`$global:DcosDnsServiceIp="$DcosDnsServiceIp"
`$global:DcosBinariesVersion="$global:DcosBinariesVersion"

function
Get-PodGateway(`$podCIDR)
{
    return `$podCIDR.substring(0,`$podCIDR.lastIndexOf(".")) + ".1"
}

function
Set-DockerNetwork(`$podCIDR)
{
    # Turn off Firewall to enable pods to talk to service endpoints. (DCOS should eventually do this)
    netsh advfirewall set allprofiles state off

    `$dockerTransparentNet=docker network ls --quiet --filter "NAME=`$global:TransparentNetworkName"
    if (`$dockerTransparentNet.length -eq 0)
    {
        `$podGW=Get-PodGateway(`$podCIDR)

        # create new transparent network
        docker network create --driver=transparent --subnet=`$podCIDR --gateway=`$podGW `$global:TransparentNetworkName

        
        `$vmswitch = get-vmSwitch  | ? SwitchType -EQ External
        # create host vnic for gateway ip to forward the traffic and dcosproxy to listen over VIP
        Add-VMNetworkAdapter -ManagementOS -Name forwarder -SwitchName `$vmswitch.Name

        # Assign gateway IP to new adapter and enable forwarding on host adapters:
        netsh interface ipv4 add address "vEthernet (forwarder)" `$podGW 255.255.255.0
        netsh interface ipv4 set interface "vEthernet (forwarder)" for=en
        netsh interface ipv4 set interface "vEthernet (HNSTransparent)" for=en
    }
}

function
Get-PodCIDR()
{
    `$podCIDR=c:\dcos\dcosctl.exe --dcosconfig=c:\dcos\config get nodes/`$(`$global:AzureHostname.ToLower()) -o custom-columns=podCidr:.spec.podCIDR --no-headers
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

    # if the podCIDR has not yet been assigned to this node, start the dcoslet process to get the podCIDR, and then promptly kill it.
    if (-not `$podCidrDiscovered)
    {
        `$argList = $DCOSArgListStr

        `$process = Start-Process -FilePath c:\dcos\dcoslet.exe -PassThru -ArgumentList `$argList

        # run dcoslet until podCidr is discovered
        Write-Host "waiting to discover pod CIDR"
        while (-not `$podCidrDiscovered)
        {
            Write-Host "Sleeping for 10s, and then waiting to discover pod CIDR"
            Start-Sleep -sec 10
            
            `$podCIDR=Get-PodCIDR
            `$podCidrDiscovered=Test-PodCIDR(`$podCIDR)
        }
    
        # stop the dcoslet process now that we have our CIDR, discard the process output
        `$process | Stop-Process | Out-Null
    }
    
    Set-DockerNetwork(`$podCIDR)

    # startup the service
    `$podGW=Get-PodGateway(`$podCIDR)
    `$env:CONTAINER_NETWORK="`$global:TransparentNetworkName"
    `$env:NAT_NETWORK="`$global:NatNetworkName"
    `$env:POD_GW="`$podGW"
    `$env:VIP_CIDR="10.0.0.0/8"

    $DCOSCommandLine
}
catch
{
    Write-Error `$_
}
"@
    $dcosStartStr | Out-File -encoding ASCII -filepath $global:DCOSStartFile

    $dcosProxyStartStr = @"
`$env:INTERFACE_TO_ADD_SERVICE_IP="vEthernet (forwarder)"
c:\dcos\dcos-proxy.exe --v=3 --proxy-mode=userspace --hostname-override=$AzureHostname --dcosconfig=c:\dcos\config
"@

    $dcosProxyStartStr | Out-File -encoding ASCII -filepath $global:DcosProxyStartFile
}

function
New-NSSMService
{
    # setup dcoslet
    c:\dcos\nssm install DCOS C:\Windows\System32\WindowsPowerShell\v1.0\powershell.exe
    c:\dcos\nssm set DCOS AppDirectory $global:DcosDir
    c:\dcos\nssm set DCOS AppParameters $global:DCOSStartFile
    c:\dcos\nssm set DCOS DisplayName DCOS
    c:\dcos\nssm set DCOS Description DCOS
    c:\dcos\nssm set DCOS Start SERVICE_AUTO_START
    c:\dcos\nssm set DCOS ObjectName LocalSystem
    c:\dcos\nssm set DCOS Type SERVICE_WIN32_OWN_PROCESS
    c:\dcos\nssm set DCOS AppThrottle 1500
    c:\dcos\nssm set DCOS AppStdout C:\dcos\dcoslet.log
    c:\dcos\nssm set DCOS AppStderr C:\dcos\dcoslet.err.log
    c:\dcos\nssm set DCOS AppStdoutCreationDisposition 4
    c:\dcos\nssm set DCOS AppStderrCreationDisposition 4
    c:\dcos\nssm set DCOS AppRotateFiles 1
    c:\dcos\nssm set DCOS AppRotateOnline 1
    c:\dcos\nssm set DCOS AppRotateSeconds 86400
    c:\dcos\nssm set DCOS AppRotateBytes 1048576
    if ($global:NeedPatchWinNAT -eq $false)
    {
        net start DCOS
    }

    # setup dcosproxy
    c:\dcos\nssm install Dcosproxy C:\Windows\System32\WindowsPowerShell\v1.0\powershell.exe
    c:\dcos\nssm set Dcosproxy AppDirectory $global:DcosDir
    c:\dcos\nssm set Dcosproxy AppParameters $global:DcosProxyStartFile
    c:\dcos\nssm set Dcosproxy DisplayName Dcosproxy
    c:\dcos\nssm set Dcosproxy DependOnService DCOS
    c:\dcos\nssm set Dcosproxy Description Dcosproxy
    c:\dcos\nssm set Dcosproxy Start SERVICE_AUTO_START
    c:\dcos\nssm set Dcosproxy ObjectName LocalSystem
    c:\dcos\nssm set Dcosproxy Type SERVICE_WIN32_OWN_PROCESS
    c:\dcos\nssm set Dcosproxy AppThrottle 1500
    c:\dcos\nssm set Dcosproxy AppStdout C:\dcos\dcosproxy.log
    c:\dcos\nssm set Dcosproxy AppStderr C:\dcos\dcosproxy.err.log
    c:\dcos\nssm set Dcosproxy AppRotateFiles 1
    c:\dcos\nssm set Dcosproxy AppRotateOnline 1
    c:\dcos\nssm set Dcosproxy AppRotateSeconds 86400
    c:\dcos\nssm set Dcosproxy AppRotateBytes 1048576
    if ($global:NeedPatchWinNAT -eq $false)
    {
        net start Dcosproxy
    }
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
    # c:\AzureData\dcosProvisionScript.log, and then you can RDP 
    # to the windows machine, and run the script manually to watch
    # the output.
    if ($true) {
        Write-Log "Provisioning $global:DockerServiceName... with IP $MasterIP"

        Write-Log "download dcoslet binaries and unzip"
        Get-DcosBinaries

        Write-Log "Write azure config"
        Write-AzureConfig

        Write-Log "Write dcos config"
        Write-DcosConfig

        Write-Log "Create the Pause Container dcoswin/pause"
        New-InfraContainer

        Write-Log "write dcoslet startfile with pod CIDR of $podCIDR"
        Write-DCOSStartFiles $podCIDR

        Write-Log "install the NSSM service"
        New-NSSMService

        Write-Log "Set Internet Explorer"
        Set-Explorer

        Write-Log "Patch winnat binary"
        Patch-WinNATBinary

        Write-Log "Setup Complete"
        if ($global:NeedPatchWinNAT -eq $true)
        {
            Write-Log "Reboot for patching winnat to be effective and start dcoslet/dcosproxy service"
            Restart-Computer
        }
    }
    else 
    {
        # keep for debugging purposes
        Write-Log ".\dcosProvisioncript.ps1 -MasterIP $MasterIP -DcosDnsServiceIp $DcosDnsServiceIp -MasterFQDNPrefix $MasterFQDNPrefix -Location $Location -AgentKey $AgentKey -AzureHostname $AzureHostname -AADClientId $AADClientId -AADClientSecret $AADClientSecret"
    }
}
catch
{
    Write-Error $_
}
