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
        Connects Docker to a Swarm Mode master.

    .DESCRIPTION
        Installs the prerequisites for creating Windows containers
        Opens TCP ports (80,443,2375,8080) in Windows Firewall.
        Connects Docker to a Swarm Mode master.

    .PARAMETER SwarmMasterIP
        IP Address of Docker Swarm Mode Master

    .EXAMPLE
        .\Join-SwarmMode-cluster.ps1 -SwarmMasterIP 192.168.255.5

#>
#Requires -Version 5.0

[CmdletBinding(DefaultParameterSetName="Standard")]
param(
    [string]
    [ValidateNotNullOrEmpty()]
    $SwarmMasterIP = "172.16.0.5"
)

$global:DockerServiceName = "Docker"
$global:DockerBinariesURL = "https://acsengine.blob.core.windows.net/swarmm/docker.zip"
$global:DockerExePath = "C:\Program Files\Docker"
$global:IsNewDockerVersion = $false

filter Timestamp {"$(Get-Date -Format o): $_"}

function Write-Log($message)
{
    $msg = $message | Timestamp
    Write-Output $msg
}

function Start-Docker()
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

function Stop-Docker()
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

function Expand-ZIPFile($file, $destination)
{
    $shell = new-object -com shell.application
    $zip = $shell.NameSpace($file)
    foreach($item in $zip.items())
    {
        $shell.Namespace($destination).copyhere($item, 0x14)
    }
}

function Install-DockerBinaries()
{
    if( $global:IsNewDockerVersion)
    {
        Write-Log "Skipping installation of new Docker binaries because latest is already installed."
        return
    }

    $currentRetry = 0;
    $success = $false;

    $zipfile = "c:\swarmm.zip"

    do {
        try
        {
            Write-Log "Downloading and installing Docker binaries...."
            Invoke-WebRequest -Uri $global:DockerBinariesURL -OutFile $zipfile
            $success = $true;
            Write-Log "Successfully downloaded Docker binaries. Number of retries: $currentRetry";
        }
        catch [System.Exception]
        {
            $message = 'Exception occurred while trying to download binaries:' + $_.Exception.ToString();
            Write-Log $message;
            if ($currentRetry -gt 5) {
                $message = "Could not download Docker binaries, aborting install. Error: " + $_.Exception.ToString();
                throw $message;
            } else {
                Write-Log "Sleeping before retry number: $currentRetry to download binaries.";
                Start-Sleep -sec 5;
            }
            $currentRetry = $currentRetry + 1;
        }
    } while (!$success);
      
    Write-Log "Expanding zip file at destination: $global:DockerExePath"
    Expand-ZIPFile -File $zipfile -Destination $global:DockerExePath

    Write-Log "Deleting zip file at: $zipfile"
    Remove-Item $zipfile
}

function Update-DockerServiceRecoveryPolicy()
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
    $tcpports = @(80,443,2375,8080,2377,7946,4789)
    foreach ($tcpport in $tcpports)
    {
        $netsh = "netsh advfirewall firewall add rule name='Open Port $tcpport' dir=in action=allow protocol=TCP localport=$tcpport"
        Write-Log "enabling port with command $netsh"
        Invoke-Expression -Command:$netsh
    }

    $udpports = @(7946,4789)
    foreach ($udpport in $udpports)
    {
        $netsh = "netsh advfirewall firewall add rule name='Open Port $udpport' dir=in action=allow protocol=UDP localport=$udpport"
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

    Write-Log "Configure Docker Engine to accept incoming connections on port 2375"
    $OutFile = @"
{
    "hosts": ["tcp://0.0.0.0:2375", "npipe://"]
}
"@

    $OutFile | Out-File -encoding ASCII -filepath "$targetDir\daemon.json"
}

function Join-Swarm()
{
    $currentRetry = 0;
    $success = $false;
    $getTokenCommand = "docker -H $($SwarmMasterIP):2375 swarm join-token -q worker"
    $swarmmodetoken;

    do {
        try
        {
            Write-Log "Executing [$getTokenCommand] command...."
            <#& $swarmmodetoken#>
            $swarmmodetoken = Invoke-Expression -Command:$getTokenCommand
            $success = $true;
            Write-Log "Successfully executed [$getTokenCommand] command. Number of entries: $currentRetry. Token: [$swarmmodetoken]";
        }
        catch [System.Exception]
        {
            $message = 'Exception occurred while trying to execute command [$swarmmodetoken]:' + $_.Exception.ToString();
            Write-Log $message;
            if ($currentRetry -gt 120) {
                $message = "Agent couldn't join Swarm, aborting install. Error: " + $_.Exception.ToString();
                throw $message;
            } else {
                Write-Log "Sleeping before $currentRetry retry of [$getTokenCommand] command";
                Start-Sleep -sec 5;
            }
            $currentRetry = $currentRetry + 1;
        }
    } while (!$success);

    $joinSwarmCommand = "docker swarm join --token $($swarmmodetoken) $($SwarmMasterIP):2377"
    Write-Log "Joining Swarm. Command [$joinSwarmCommand]...."
    Invoke-Expression -Command:$joinSwarmCommand
}

function Confirm-DockerVersion()
{
   $dockerServerVersionCmd = "docker version --format '{{.Server.Version}}'"
   Write-Log "Running command: $dockerServerVersionCmd"
   $dockerServerVersion = Invoke-Expression -Command:$dockerServerVersionCmd

   $dockerClientVersionCmd = "docker version --format '{{.Client.Version}}'"
   Write-Log "Running command: $dockerClientVersionCmd"
   $dockerClientVersion = Invoke-Expression -Command:$dockerClientVersionCmd

   Write-Log "Docker Server version: $dockerServerVersion, Docker Client verison: $dockerClientVersion"
   
   $serverVersionData = $dockerServerVersion.Split(".")
   $isNewServerVersion = $false;
   if(($serverVersionData[0] -ge 1) -and ($serverVersionData[1] -ge 13)){
       $isNewServerVersion = $true;
       Write-Log "Setting isNewServerVersion to $isNewServerVersion"
   }

   $clientVersionData = $dockerClientVersion.Split(".")
   $isNewClientVersion = $false;
   if(($clientVersionData[0] -ge 1) -and ($clientVersionData[1] -ge 13)){
       $isNewClientVersion = $true;
       Write-Log "Setting  isNewClientVersion to $isNewClientVersion"   
   }

   if($isNewServerVersion -and $isNewClientVersion)
   {
       $global:IsNewDockerVersion = $true;
       Write-Log "Setting IsNewDockerVersion to $global:IsNewDockerVersion"
   }
}

try
{
    Write-Log "Provisioning $global:DockerServiceName... with Swarm IP $SwarmMasterIP"

    Write-Log "Checking Docker version"
    Confirm-DockerVersion

    Write-Log "Stop Docker"
    Stop-Docker

    Write-Log "Installing Docker binaries"
    Install-DockerBinaries

    Write-Log "Opening firewall ports"
    Open-FirewallPorts

    Write-Log "Write Docker Configuration"
    Write-DockerDaemonJson

    Write-Log "Update Docker restart policy"
    Update-DockerServiceRecoveryPolicy
    
    Write-Log "Start Docker"
    Start-Docker
    
    Write-Log "Join existing Swarm"
    Join-Swarm

    Write-Log "Setup Complete"
}
catch
{
    Write-Error $_
}