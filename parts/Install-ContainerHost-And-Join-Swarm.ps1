############################################################
# Script adapted from
# https://raw.githubusercontent.com/Microsoft/Virtualization-Documentation/master/windows-server-container-tools/Install-ContainerHost/Install-ContainerHost.ps1

<#
    .NOTES
        Copyright (c) Microsoft Corporation.  All rights reserved.

        Use of this sample source code is subject to the terms of the Microsoft
        license agreement under which you licensed this sample source code. If
        you did not accept the terms of the license agreement, you are not
        authorized to use this sample source code. For the terms of the license,
        please see the license agreement between you and Microsoft or, if applicable,
        see the LICENSE.RTF on your install media or the root of your tools installation.
        THE SAMPLE SOURCE CODE IS PROVIDED "AS IS", WITH NO WARRANTIES.

    .SYNOPSIS
        Installs the prerequisites for creating Windows containers
        Opens TCP ports (80,443,2375,8080) in Windows Firewall.
        Connects Docker to a swarm master.

    .DESCRIPTION
        Installs the prerequisites for creating Windows containers
        Opens TCP ports (80,443,2375,8080) in Windows Firewall.
        Connects Docker to a swarm master.

    .PARAMETER SwarmMasterIP
        IP Address of Docker Swarm Master

    .EXAMPLE
        .\Install-ContainerHost.ps1 -SwarmMasterIP 192.168.255.5

#>
#Requires -Version 5.0

[CmdletBinding(DefaultParameterSetName="Standard")]
param(
    [string]
    [ValidateNotNullOrEmpty()]
    $SwarmMasterIP = "172.16.0.5"
)

$global:DockerServiceName = "Docker"
$global:HNSServiceName = "hns"

filter Timestamp {"$(Get-Date -Format o): $_"}

function Write-Log($message)
{
    $msg = $message | Timestamp
    Write-Output $msg
}

function
Start-Docker()
{
    Write-Log "Starting $global:DockerServiceName..."
    $startTime = Get-Date
        
    while (-not $dockerReady)
    {
        try
        {
            Start-Service -Name $global:DockerServiceName -ea Stop

            $dockerReady = $true            
        }
        catch
        {
            $timeElapsed = $(Get-Date) - $startTime
            if ($($timeElapsed).TotalMinutes -ge 5)
            {
                Write-Log "Docker Daemon did not start successfully within 5 minutes."
                break
            }

            $errorStr = $_.Exception.Message
            Write-Log "Starting Service failed: $errorStr" 
            Write-Log "sleeping for 10 seconds..."
            Start-Sleep -sec 10
        }
    }
}


function
Stop-Docker()
{
    Write-Log "Stopping $global:DockerServiceName..."
    try
    {
        Stop-Service -Name $global:DockerServiceName -ea Stop   
    }
    catch
    {
        Write-Log "Failed to stop Docker"
    }
}

function
Update-DockerServiceRecoveryPolicy()
{
    $dockerReady = $false
    $startTime = Get-Date
    
    # wait until the service exists
    while (-not $dockerReady)
    {
        if (Get-Service $global:DockerServiceName -ErrorAction SilentlyContinue)
        {
            $dockerReady = $true
        }
        else 
        {
            $timeElapsed = $(Get-Date) - $startTime
            if ($($timeElapsed).TotalMinutes -ge 5)
            {
                Write-Log "Unable to find service $global:DockerServiceName within 5 minutes."
                break
            }
            Write-Log "failed to find $global:DockerServiceName, sleeping for 5 seconds"
            Start-Sleep -sec 5
        }
    }
    
    Write-Log "Updating docker restart policy, to ensure it restarts on error"
    $services = Get-WMIObject win32_service | Where-Object {$_.name -imatch $global:DockerServiceName}
    foreach ($service in $services)
    {
        sc.exe failure $service.name reset= 86400 actions= restart/5000
    }
}

# Open Windows Firewall Ports Needed
function Open-FirewallPorts()
{
    $ports = @(80,443,2375,8080)
    foreach ($port in $ports)
    {
        $netsh = "netsh advfirewall firewall add rule name='Open Port $port' dir=in action=allow protocol=TCP localport=$port"
        Write-Log "enabling port with command $netsh"
        Invoke-Expression -Command:$netsh
    }
}

# Update Docker Config to have cluster-store=consul:// address configured for Swarm cluster.
function Write-DockerDaemonJson()
{
    $dataDir = $env:ProgramData

    # create the target directory
    $targetDir = $dataDir + '\docker\config'
    if(!(Test-Path -Path $targetDir )){
        New-Item -ItemType directory -Path $targetDir
    }

    Write-Log "Delete key file, so that this node is unique to swarm"
    $keyFileName = "$targetDir\key.json"
    Write-Log "Removing $($keyFileName)"
    if (Test-Path $keyFileName) {
      Remove-Item $keyFileName
    }

    $ipAddress = Get-IPAddress

    Write-Log "Advertise $($ipAddress) to consul://$($SwarmMasterIP):8500"
    $OutFile = @"
{
    "hosts": ["tcp://0.0.0.0:2375", "npipe://"],
    "cluster-store": "consul://$($SwarmMasterIP):8500",
    "cluster-advertise": "$($ipAddress):2375"
}
"@

    $OutFile | Out-File -encoding ASCII -filepath "$targetDir\daemon.json"
}

# Get Node IPV4 Address
function Get-IPAddress()
{
    return (Get-NetIPAddress | where {$_.IPAddress -Like '10.*' -and $_.AddressFamily -eq 'IPV4'})[0].IPAddress
}

try
{
    Write-Log "Provisioning $global:DockerServiceName... with Swarm IP $SwarmMasterIP"

    Write-Log "Stop Docker"
    Stop-Docker

    Write-Log "Opening firewall ports"
    Open-FirewallPorts

    Write-Log "Write Docker Configuration"
    Write-DockerDaemonJson

    Write-Log "Update Docker restart policy"
    Update-DockerServiceRecoveryPolicy
    
    Write-Log "Start Docker"
    Start-Docker

    #remove-ItemProperty -Path "HKLM:\SYSTEM\CurrentControlSet\Control\Wininit"  Headless
    #Write-Log "shutdown /r /f /t 60"
    #shutdown /r /f /t 60

    Write-Log "Setup Complete"
}
catch
{
    Write-Error $_
}


