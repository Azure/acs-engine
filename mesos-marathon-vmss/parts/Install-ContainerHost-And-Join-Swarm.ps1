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
        Updates Windows Docker Binary from DockerPath
        Rewrites runDockerDaemon.cmd startup script to join a Docker Swarm
        Start Docker

    .DESCRIPTION
        Installs the prerequisites for creating Windows containers
        Opens TCP ports (80,443,2375,8080) in Windows Firewall.
        Updates Windows Docker Binary from DockerPath
        Rewrites runDockerDaemon.cmd startup script to join a Docker Swarm
        Start Docker

    .PARAMETER DockerPath
        Path to Docker.exe, can be local or URI

    .PARAMETER DockerDPath
        Path to DockerD.exe, can be local or URI

    .PARAMETER ExternalNetAdapter
        Specify a specific network adapter to bind to a DHCP network

    .PARAMETER Force 
        If a restart is required, forces an immediate restart.
        
    .PARAMETER HyperV 
        If passed, prepare the machine for Hyper-V containers

    .PARAMETER NoRestart
        If a restart is required the script will terminate and will not reboot the machine

    .PARAMETER SkipImageImport
        Skips import of the base WindowsServerCore image.

    .PARAMETER TransparentNetwork
        If passed, use DHCP configuration.  Otherwise, will use default docker network (NAT). (alias -UseDHCP)

    .PARAMETER WimPath
        Path to .wim file that contains the base package image
        
    .PARAMETER SwarmMasterIP
        IP Address of Docker Swarm Master

    .EXAMPLE
        .\Install-ContainerHost.ps1

#>
#Requires -Version 5.0

[CmdletBinding(DefaultParameterSetName="Standard")]
param(
    [string]
    [ValidateNotNullOrEmpty()]
    $DockerPath = "https://aka.ms/tp5/docker",

    [string]
    [ValidateNotNullOrEmpty()]
    $DockerDPath = "https://aka.ms/tp5/dockerd",

    [string]
    $ExternalNetAdapter,

    [switch]
    $Force,

    [switch]
    $HyperV,

    [switch]
    $NoRestart,

    [Parameter(DontShow)]
    [switch]
    $PSDirect,

    [switch]
    $SkipImageImport,

    [Parameter(ParameterSetName="Staging", Mandatory)]
    [switch]
    $Staging,

    [switch]
    [alias("UseDHCP")]
    $TransparentNetwork,

    [string]
    [ValidateNotNullOrEmpty()]
    $WimPath,
    
    [string]
    [ValidateNotNullOrEmpty()]
    $SwarmMasterIP = "172.16.0.5"
)

$global:RebootRequired = $false

$global:ErrorFile = "$pwd\Install-ContainerHost.err"

$global:BootstrapTask = "ContainerBootstrap"

$global:HyperVImage = "NanoServer"

filter Timestamp {"$(Get-Date -Format o): $_"}

function Write-Log($message)
{
    Write-Output $message | timestamp
}

function
Restart-And-Run()
{
    Test-Admin

    Write-Log "Restart is required; restarting now..." | timestamp

    $argList = $script:MyInvocation.Line.replace($script:MyInvocation.InvocationName, "")

    $scriptPath = $script:MyInvocation.MyCommand.Path

    $argList = $argList -replace "\.\\", "$pwd\"

    $logPath = $scriptPath.Replace(".ps1", ".log")

    $argList = $argList -replace "\.\\", "$scriptPath\"

    # We wrap the powershell call to a call to CMD like custom script extension so that
    # appending the output to the same file isn't written as Unicode. Existing file is not Unicode.
    Write-Log "Creating scheduled task action cmd /c powershell.exe -NoExit ($scriptPath $argList >> $logPath 2>&1)..."
    $action = New-ScheduledTaskAction -Execute "cmd" -Argument " /c powershell.exe -ExecutionPolicy Unrestricted $scriptPath $argList >> $logPath 2>&1" 

    Write-Log "Creating scheduled task trigger..."
    $trigger = New-ScheduledTaskTrigger -AtStartup

    # Custom script extension runs code as SYSTEM
    Write-Log "Registering script to re-run at next startup..."
    Register-ScheduledTask -TaskName $global:BootstrapTask -Action $action -Trigger $trigger -RunLevel Highest -User System | Out-Null

    try
    {
        if ($Force)
        {
            Restart-Computer -Force
        }
        else
        {
            Restart-Computer
        }
    }
    catch
    {
        Write-Error $_

        Write-Log "Please restart your computer manually to continue script execution."
    }

    exit
}


function
Install-Feature
{
    [CmdletBinding()]
    param(
        [ValidateNotNullOrEmpty()]
        [string]
        $FeatureName
    )

    Write-Log "Querying status of Windows feature: $FeatureName..."
    if (Get-Command Get-WindowsFeature -ErrorAction SilentlyContinue)
    {
        if ((Get-WindowsFeature $FeatureName).Installed)
        {
            Write-Log "Feature $FeatureName is already enabled."
        }
        else
        {
            Test-Admin

            Write-Log "Enabling feature $FeatureName..."
        }

        $time = Measure-Command {
            $featureInstall = Add-WindowsFeature $FeatureName

            if ($featureInstall.RestartNeeded -eq "Yes")
            {
                $global:RebootRequired = $true;
            }
        }
        
        Write-Log "Add-WindowsFeature $FeatureName took $time"
    }
    else
    {
        if ((Get-WindowsOptionalFeature -Online -FeatureName $FeatureName).State -eq "Disabled")
        {
            if (Test-Nano)
            {
                throw "This NanoServer deployment does not include $FeatureName.  Please add the appropriate package"
            }

            Test-Admin

            Write-Log "Enabling feature $FeatureName..."
            $feature = Enable-WindowsOptionalFeature -Online -FeatureName $FeatureName -All -NoRestart

            if ($feature.RestartNeeded -eq "True")
            {
                $global:RebootRequired = $true;
            }
        }
        else
        {
            Write-Log "Feature $FeatureName is already enabled."

            if (Test-Nano)
            {
                #
                # Get-WindowsEdition is not present on Nano.  On Nano, we assume reboot is not needed
                #
            }
            elseif ((Get-WindowsEdition -Online).RestartNeeded)
            {
                $global:RebootRequired = $true;
            }
        }
    }
}


function
New-ContainerTransparentNetwork
{
    if ($ExternalNetAdapter)
    {
        $netAdapter = (Get-NetAdapter |? {$_.Name -eq "$ExternalNetAdapter"})[0]
    }
    else
    {
        $netAdapter = (Get-NetAdapter |? {($_.Status -eq 'Up') -and ($_.ConnectorPresent)})[0]
    }

    Write-Log "Creating container network (Transparent)..."
    New-ContainerNetwork -Name "Transparent" -Mode Transparent -NetworkAdapterName $netAdapter.Name | Out-Null
}


function
Install-ContainerHost
{
    "If this file exists when Install-ContainerHost.ps1 exits, the script failed!" | Out-File -FilePath $global:ErrorFile

    if (Test-Client)
    {
        if (-not $HyperV)
        {
            Write-Log "Enabling Hyper-V containers by default for Client SKU"
            $HyperV = $true
        }    
    }
    #
    # Validate required Windows features
    #
    Install-Feature -FeatureName Containers

    if ($HyperV)
    {
        Install-Feature -FeatureName Hyper-V

        #
        # TODO: remove if/else when IUM and DirectMap can coexist
        #
        #if ((Get-ItemProperty -Path "HKLM:\SYSTEM\CurrentControlSet\Control\DeviceGuard" -Name HyperVVirtualizationBasedSecurityOptOut -ErrorAction SilentlyContinue).HyperVVirtualizationBasedSecurityOptOut -eq 1)
        #{
        #    Write-Log "IUM is already disabled (DirectMap will be operational)."
        #}
        #else
        #{
        #    Write-Log "Disabling IUM to enable DirectMap"
        #    if (-not (Get-Item -Path "HKLM:\SYSTEM\CurrentControlSet\Control\DeviceGuard" -ErrorAction SilentlyContinue))
        #    {
        #        New-Item -Path "HKLM:\SYSTEM\CurrentControlSet\Control\DeviceGuard"
        #    }
        #
        #    Set-ItemProperty -Path "HKLM:\SYSTEM\CurrentControlSet\Control\DeviceGuard" -Name HyperVVirtualizationBasedSecurityOptOut -Value 1
        #    $global:RebootRequired = $true
        #}
    }

    if ($global:RebootRequired)
    {
        if ($NoRestart)
        {
            Write-Warning "A reboot is required; stopping script execution"
            exit
        }

        Restart-And-Run
    }

    #
    # Unregister the bootstrap task, if it was previously created
    #
    if ((Get-ScheduledTask -TaskName $global:BootstrapTask -ErrorAction SilentlyContinue) -ne $null)
    {
        schtasks /DELETE /TN $global:BootstrapTask /F
        #Unregister-ScheduledTask -TaskName $global:BootstrapTask -Confirm $true
    }
    

    #
    # Configure networking
    #
    if ($($PSCmdlet.ParameterSetName) -ne "Staging")
    {
        Write-Log "Configuring ICMP firewall rules for containers..."
        netsh advfirewall firewall add rule name="ICMP for containers" dir=in protocol=icmpv4 action=allow | Out-Null
        netsh advfirewall firewall add rule name="ICMP for containers" dir=out protocol=icmpv4 action=allow | Out-Null
        
        if ($TransparentNetwork)
        {
            Write-Log "Waiting for Hyper-V Management..."
            $networks = $null

            try
            {
                $networks = Get-ContainerNetwork -ErrorAction SilentlyContinue
            }
            catch
            {
                #
                # If we can't query network, we are in bootstrap mode.  Assume no networks
                #
            }

            if ($networks.Count -eq 0)
            {
                Write-Log "Enabling container networking..."
                New-ContainerTransparentNetwork
            }
            else
            {
                Write-Log "Networking is already configured.  Confirming configuration..."
                
                $transparentNetwork = $networks |? { $_.Mode -eq "Transparent" }

                if ($transparentNetwork -eq $null)
                {
                    Write-Log "We didn't find a configured external network; configuring now..."
                    New-ContainerTransparentNetwork
                }
                else
                {
                    if ($ExternalNetAdapter)
                    {
                        $netAdapters = (Get-NetAdapter |? {$_.Name -eq "$ExternalNetAdapter"})

                        if ($netAdapters.Count -eq 0)
                        {
                            throw "No adapters found that match the name $ExternalNetAdapter"
                        }

                        $netAdapter = $netAdapters[0]
                        $transparentNetwork = $networks |? { $_.NetworkAdapterName -eq $netAdapter.InterfaceDescription }

                        if ($transparentNetwork -eq $null)
                        {
                            throw "One or more external networks are configured, but not on the requested adapter ($ExternalNetAdapter)"
                        }

                        Write-Log "Configured transparent network found: $($transparentNetwork.Name)"
                    }
                    else
                    {
                        Write-Log "Configured transparent network found: $($transparentNetwork.Name)"
                    }
                }
            }
        }
    }

    #
    # Install, register, and start Docker
    #
    if (Test-Docker)
    {
        Write-Log "Docker is already installed."
    }
    else
    {
        $time = Measure-Command {
            Install-Docker -DockerPath $DockerPath -DockerDPath $DockerDPath
        }
        
        Write-Log "Install-Docker -DockerPath $DockerPath -DockerDPath $DockerDPath took $time"
    }

    $newBaseImages = @()

    if (-not $SkipImageImport)
    {        
        if ($WimPath -eq "")
        {
            $imageName = "WindowsServerCore"

            if ($HyperV -or (Test-Nano))
            {
                $imageName = "NanoServer"
            }

            #
            # Install the base package
            #
            if (Test-InstalledContainerImage $imageName)
            {
                Write-Log "Image $imageName is already installed on this machine."
            }
            else
            {
                Test-ContainerImageProvider

                $hostBuildInfo = (gp "HKLM:\SOFTWARE\Microsoft\Windows NT\CurrentVersion").BuildLabEx.Split(".")
                $version = $hostBuildInfo[0]

                $InstallParams = @{
                    ErrorAction = "Stop"
                    Name = $imageName
                }

                if ($version -eq "14300")
                {
                    $InstallParams.Add("MinimumVersion", "10.0.14300.1000")
                    $InstallParams.Add("MaximumVersion", "10.0.14300.1010")
                    $versionString = "-MinimumVersion 10.0.14300.1000 -MaximumVersion 10.0.14300.1010"
                }
                else
                {
                    if (Test-Client)
                    {
                        $versionString = " [latest version]"
                    }
                    else
                    {
                        $qfe = $hostBuildInfo[1]

                        $InstallParams.Add("RequiredVersion", "10.0.$version.$qfe")
                        $versionString = "-RequiredVersion 10.0.$version.$qfe"
                    }                    
                }

                $time = Measure-Command {
                    Write-Log "Getting Container OS image ($imageName $versionString) from OneGet (this may take a few minutes)..."
                    #
                    # TODO: expect the follow to have default ErrorAction of stop
                    #
                    Install-ContainerImage @InstallParams
            
                    Write-Log "Container base image install complete."
                }
        
                Write-Log "Install-ContainerImage ($imageName $versionString) took $time"
                
                $newBaseImages += $imageName
            }
        }
        else
        {
            Write-Log "Installing Container OS image from $WimPath (this may take a few minutes)..."

            if (Test-Path $WimPath)
            {
                #
                # .wim is present and local
                #
            }
            elseif (($WimPath -as [System.URI]).AbsoluteURI -ne $null)
            {
                #
                # .wim is on a URI and must be downloaded
                #
                $localWimPath = "$pwd\ContainerBaseImage.wim"

                Copy-File -SourcePath $WimPath -DestinationPath $localWimPath

                $WimPath = $localWimPath
            }
            else
            {
                throw "Cannot copy from invalid WimPath $WimPath"
            }

            $imageName = (get-windowsimage -imagepath $WimPath -LogPath ($env:temp+"dism_$(random)_GetImageInfo.log") -Index 1).imagename
                        
            if ($PSDirect -and (Test-Nano))
            {
                #
                # This is a gross hack for TP5 to avoid a CoreCLR issue
                #
                $modulePath = "$($env:Temp)\Containers2.psm1"

                $cmdletContent = gc $env:windir\System32\WindowsPowerShell\v1.0\Modules\Containers\1.0.0.0\Containers.psm1

                $cmdletContent = $cmdletContent.replace('Set-Acl $fileToReAcl -AclObject $acl', '[System.IO.FileSystemAclExtensions]::SetAccessControl($fileToReAcl, $acl)')
                $cmdletContent = $cmdletContent.replace('function Install-ContainerOSImage','function Install-ContainerOSImage2')

                $cmdletContent | sc $modulePath

                Import-Module $modulePath -DisableNameChecking
                Install-ContainerOSImage2 -WimPath $WimPath -Force
                Remove-Item $modulePath
            }
            else
            {
                Install-ContainerOsImage -WimPath $WimPath -Force
            }

            $newBaseImages += $imageName
        }

        #
        # Optionally OneGet the Hyper-V container image if it isn't just installed
        #
        if ($HyperV -and (-not (Test-Nano)))
        {
            if ((Test-InstalledContainerImage $global:HyperVImage))
            {
                Write-Log "OS image ($global:HyperVImage) is already installed."
            }
            else
            {
                Test-ContainerImageProvider

                Write-Log "Getting Container OS image ($global:HyperVImage) from OneGet (this may take a few minutes)..."
                Install-ContainerImage $global:HyperVImage

                $newBaseImages += $global:HyperVImage
            }
        }
    }

    if ($newBaseImages.Count -gt 0)
    {
        foreach ($baseImage in $newBaseImages)
        {
            Write-DockerImageTag -BaseImageName $baseImage
        }

        "tag complete" | Out-File -FilePath "$dockerData\tag.txt" -Encoding ASCII

        #
        # if certs.d exists, restart docker in TLS mode
        #
        $dockerCerts = "$($env:ProgramData)\docker\certs.d"

        if (Test-Path $dockerCerts)
        {
            if ((Get-ChildItem $dockerCerts).Count -gt 0)
            {
                Stop-Docker
                Start-Docker
            }
        }
    }

    Remove-Item $global:ErrorFile

    Write-Log "Script complete!"
}$global:AdminPriviledges = $false
$global:DockerServiceName = "Docker"

function
Copy-File
{
    [CmdletBinding()]
    param(
        [string]
        $SourcePath,
        
        [string]
        $DestinationPath
    )
    
    if ($SourcePath -eq $DestinationPath)
    {
        return
    }
          
    if (Test-Path $SourcePath)
    {
        Copy-Item -Path $SourcePath -Destination $DestinationPath
    }
    elseif (($SourcePath -as [System.URI]).AbsoluteURI -ne $null)
    {
        if (Test-Nano)
        {
            $handler = New-Object System.Net.Http.HttpClientHandler
            $client = New-Object System.Net.Http.HttpClient($handler)
            $client.Timeout = New-Object System.TimeSpan(0, 30, 0)
            $cancelTokenSource = [System.Threading.CancellationTokenSource]::new() 
            $responseMsg = $client.GetAsync([System.Uri]::new($SourcePath), $cancelTokenSource.Token)
            $responseMsg.Wait()

            if (!$responseMsg.IsCanceled)
            {
                $response = $responseMsg.Result
                if ($response.IsSuccessStatusCode)
                {
                    $downloadedFileStream = [System.IO.FileStream]::new($DestinationPath, [System.IO.FileMode]::Create, [System.IO.FileAccess]::Write)
                    $copyStreamOp = $response.Content.CopyToAsync($downloadedFileStream)
                    $copyStreamOp.Wait()
                    $downloadedFileStream.Close()
                    if ($copyStreamOp.Exception -ne $null)
                    {
                        throw $copyStreamOp.Exception
                    }      
                }
            }  
        }
        elseif ($PSVersionTable.PSVersion.Major -ge 5)
        {
            #
            # We disable progress display because it kills performance for large downloads (at least on 64-bit PowerShell)
            #
            $ProgressPreference = 'SilentlyContinue'
            wget -Uri $SourcePath -OutFile $DestinationPath -UseBasicParsing
            $ProgressPreference = 'Continue'
        }
        else
        {
            $webClient = New-Object System.Net.WebClient
            $webClient.DownloadFile($SourcePath, $DestinationPath)
        } 
    }
    else
    {
        throw "Cannot copy from $SourcePath"
    }
}


function 
Expand-ArchiveNano
{
    [CmdletBinding()]
    param 
    (
        [string] $Path,
        [string] $DestinationPath
    )

    [System.IO.Compression.ZipFile]::ExtractToDirectory($Path, $DestinationPath)
}


function 
Expand-ArchivePrivate
{
    [CmdletBinding()]
    param 
    (
        [Parameter(Mandatory=$true)]
        [string] 
        $Path,

        [Parameter(Mandatory=$true)]        
        [string] 
        $DestinationPath
    )
        
    $shell = New-Object -com Shell.Application
    $zipFile = $shell.NameSpace($Path)
    
    $shell.NameSpace($DestinationPath).CopyHere($zipFile.items())
    
}


function
Test-InstalledContainerImage
{
    [CmdletBinding()]
    param(
        [Parameter(Mandatory=$true)]
        [string]
        [ValidateNotNullOrEmpty()]
        $BaseImageName
    )

    $path = Join-Path (Join-Path $env:ProgramData "Microsoft\Windows\Images") "*$BaseImageName*"
    
    return Test-Path $path
}


function
Get-Nsmm
{
    [CmdletBinding()]
    param(
        [Parameter(Mandatory=$true)]
        [string]
        [ValidateNotNullOrEmpty()]
        $Destination,

        [string]
        [ValidateNotNullOrEmpty()]
        $WorkingDir = "$env:temp"
    )
    
    Write-Log "This script uses a third party tool: NSSM. For more information, see https://nssm.cc/usage"       
    Write-Log "Downloading NSSM..."

    $nssmUri = "https://nssm.cc/release/nssm-2.24.zip"            
    $nssmZip = "$($env:temp)\$(Split-Path $nssmUri -Leaf)"
            
    Write-Verbose "Creating working directory..."
    $tempDirectory = New-Item -ItemType Directory -Force -Path "$($env:temp)\nssm"
    
    Copy-File -SourcePath $nssmUri -DestinationPath $nssmZip
            
    Write-Log "Extracting NSSM from archive..."
    if (Test-Nano)
    {
        Expand-ArchiveNano -Path $nssmZip -DestinationPath $tempDirectory.FullName
    }
    elseif ($PSVersionTable.PSVersion.Major -ge 5)
    {
        Expand-Archive -Path $nssmZip -DestinationPath $tempDirectory.FullName
    }
    else
    {
        Expand-ArchivePrivate -Path $nssmZip -DestinationPath $tempDirectory.FullName
    }
    Remove-Item $nssmZip

    Write-Verbose "Copying NSSM to $Destination..."
    Copy-Item -Path "$($tempDirectory.FullName)\nssm-2.24\win64\nssm.exe" -Destination "$Destination"

    Write-Verbose "Removing temporary directory..."
    $tempDirectory | Remove-Item -Recurse
}


function 
Test-Admin()
{
    # Get the ID and security principal of the current user account
    $myWindowsID=[System.Security.Principal.WindowsIdentity]::GetCurrent()
    $myWindowsPrincipal=new-object System.Security.Principal.WindowsPrincipal($myWindowsID)
  
    # Get the security principal for the Administrator role
    $adminRole=[System.Security.Principal.WindowsBuiltInRole]::Administrator
  
    # Check to see if we are currently running "as Administrator"
    if ($myWindowsPrincipal.IsInRole($adminRole))
    {
        $global:AdminPriviledges = $true
        return
    }
    else
    {
        #
        # We are not running "as Administrator"
        # Exit from the current, unelevated, process
        #
        throw "You must run this script as administrator"   
    }
}


function 
Test-ContainerImageProvider()
{
    if (-not (Get-Command Install-ContainerImage -ea SilentlyContinue))
    {   
        Wait-Network

        Write-Log "Installing ContainerImage provider..."
        Install-PackageProvider ContainerImage -Force | Out-Null
    }

    if (-not (Get-Command Install-ContainerImage -ea SilentlyContinue))
    {
        throw "Could not install ContainerImage provider"
    }
}


function 
Test-Client()
{
    return (-not ((Get-Command Get-WindowsFeature -ErrorAction SilentlyContinue) -or (Test-Nano)))
}


function 
Test-Nano()
{
    $EditionId = (Get-ItemProperty -Path 'HKLM:\SOFTWARE\Microsoft\Windows NT\CurrentVersion' -Name 'EditionID').EditionId

    return (($EditionId -eq "ServerStandardNano") -or 
            ($EditionId -eq "ServerDataCenterNano") -or 
            ($EditionId -eq "NanoServer") -or 
            ($EditionId -eq "ServerTuva"))
}


function 
Wait-Network()
{
    $connectedAdapter = Get-NetAdapter |? ConnectorPresent

    if ($connectedAdapter -eq $null)
    {
        throw "No connected network"
    }
       
    $startTime = Get-Date
    $timeElapsed = $(Get-Date) - $startTime

    while ($($timeElapsed).TotalMinutes -lt 5)
    {
        $readyNetAdapter = $connectedAdapter |? Status -eq 'Up'

        if ($readyNetAdapter -ne $null)
        {
            return;
        }

        Write-Log "Waiting for network connectivity..."
        Start-Sleep -sec 5

        $timeElapsed = $(Get-Date) - $startTime
    }

    throw "Network not connected after 5 minutes"
}


function
Get-DockerImages
{
    return docker images
}

function
Find-DockerImages
{
    [CmdletBinding()]
    param(
        [Parameter(Mandatory=$true)]
        [string]
        [ValidateNotNullOrEmpty()]
        $BaseImageName
    )

    return docker images | Where { $_ -match $BaseImageName.tolower() }
}


function 
Install-Docker()
{
    [CmdletBinding()]
    param(
        [string]
        [ValidateNotNullOrEmpty()]
        $DockerPath = "https://aka.ms/tp5/docker",

        [string]
        [ValidateNotNullOrEmpty()]
        $DockerDPath = "https://aka.ms/tp5/dockerd"
    )

    Test-Admin

    Write-Log "Installing Docker..."
    Copy-File -SourcePath $DockerPath -DestinationPath $env:windir\System32\docker.exe

    try
    {
        Write-Log "Installing Docker daemon..."
        Copy-File -SourcePath $DockerDPath -DestinationPath $env:windir\System32\dockerd.exe
    }
    catch 
    {
        Write-Warning "DockerD not yet present."
    }

    $dockerData = "$($env:ProgramData)\docker"
    $dockerLog = "$dockerData\daemon.log"

    if (-not (Test-Path $dockerData))
    {
        Write-Log "Creating Docker program data..."
        New-Item -ItemType Directory -Force -Path $dockerData | Out-Null
    }

    $dockerDaemonScript = "$dockerData\runDockerDaemon.cmd"

    New-DockerDaemonRunText | Out-File -FilePath $dockerDaemonScript -Encoding ASCII

    if (Test-Nano)
    {
        Write-Log "Creating scheduled task action..."
        $action = New-ScheduledTaskAction -Execute "cmd.exe" -Argument "/c $dockerDaemonScript > $dockerLog 2>&1"

        Write-Log "Creating scheduled task trigger..."
        $trigger = New-ScheduledTaskTrigger -AtStartup

        Write-Log "Creating scheduled task settings..."
        $settings = New-ScheduledTaskSettingsSet -Priority 5

        Write-Log "Registering Docker daemon to launch at startup..."
        Register-ScheduledTask -TaskName $global:DockerServiceName -Action $action -Trigger $trigger -Settings $settings -User SYSTEM -RunLevel Highest | Out-Null

        Write-Log "Launching daemon..."
        Start-ScheduledTask -TaskName $global:DockerServiceName
    }
    else
    {
        if (Test-Path "$($env:SystemRoot)\System32\nssm.exe")
        {
            Write-Log "NSSM is already installed"
        }
        else
        {
            Get-Nsmm -Destination "$($env:SystemRoot)\System32" -WorkingDir "$env:temp"
        }

        Write-Log "Configuring NSSM for $global:DockerServiceName service..."
        Start-Process -Wait "nssm" -ArgumentList "install $global:DockerServiceName $($env:SystemRoot)\System32\cmd.exe /s /c $dockerDaemonScript < nul"
        Start-Process -Wait "nssm" -ArgumentList "set $global:DockerServiceName DisplayName Docker Daemon"
        Start-Process -Wait "nssm" -ArgumentList "set $global:DockerServiceName Description The Docker Daemon provides management capabilities of containers for docker clients"
        # Pipe output to daemon.log
        Start-Process -Wait "nssm" -ArgumentList "set $global:DockerServiceName AppStderr $dockerLog"
        Start-Process -Wait "nssm" -ArgumentList "set $global:DockerServiceName AppStdout $dockerLog"
        # Allow 30 seconds for graceful shutdown before process is terminated
        Start-Process -Wait "nssm" -ArgumentList "set $global:DockerServiceName AppStopMethodConsole 30000"

        Start-Service -Name $global:DockerServiceName
    }

    #
    # Waiting for docker to come to steady state
    #
    Wait-Docker

    Write-Log "The following images are present on this machine:"
    foreach ($image in (Get-DockerImages))
    {
        Write-Log "    $image"
    }
    Write-Log ""
}


function
New-DockerDaemonRunText
{
    return @"

@echo off
set certs=%ProgramData%\docker\certs.d

if exist %ProgramData%\docker (goto :run)
mkdir %ProgramData%\docker

:run
if exist %certs%\server-cert.pem (if exist %ProgramData%\docker\tag.txt (goto :secure))

if not exist %systemroot%\system32\dockerd.exe (goto :legacy)

dockerd -H npipe:// 
goto :eof

:legacy
docker daemon -H npipe:// 
goto :eof

:secure
if not exist %systemroot%\system32\dockerd.exe (goto :legacysecure)
dockerd -H npipe:// -H 0.0.0.0:2376 --tlsverify --tlscacert=%certs%\ca.pem --tlscert=%certs%\server-cert.pem --tlskey=%certs%\server-key.pem
goto :eof

:legacysecure
docker daemon -H npipe:// -H 0.0.0.0:2376 --tlsverify --tlscacert=%certs%\ca.pem --tlscert=%certs%\server-cert.pem --tlskey=%certs%\server-key.pem

"@

}


function 
Start-Docker()
{
    Write-Log "Starting $global:DockerServiceName..."
    if (Test-Nano)
    {
        Start-ScheduledTask -TaskName $global:DockerServiceName
    }
    else
    {
        Start-Service -Name $global:DockerServiceName
    }
}


function 
Stop-Docker()
{
    Write-Log "Stopping $global:DockerServiceName..."
    if (Test-Nano)
    {
        Stop-ScheduledTask -TaskName $global:DockerServiceName

        #
        # ISSUE: can we do this more gently?
        #
        Get-Process $global:DockerServiceName | Stop-Process -Force
    }
    else
    {
        Stop-Service -Name $global:DockerServiceName
    }
}


function 
Test-Docker()
{
    $service = $null

    if (Test-Nano)
    {
        $service = Get-ScheduledTask -TaskName $global:DockerServiceName -ErrorAction SilentlyContinue
    }
    else
    {
        $service = Get-Service -Name $global:DockerServiceName -ErrorAction SilentlyContinue
    }

    return ($service -ne $null)
}


function 
Wait-Docker()
{
    Write-Log "Waiting for Docker daemon..."
    $dockerReady = $false
    $startTime = Get-Date

    while (-not $dockerReady)
    {
        try
        {
            docker version | Out-Null

            if (-not $?)
            {
                throw "Docker daemon is not running yet"
            }

            $dockerReady = $true
        }
        catch 
        {
            $timeElapsed = $(Get-Date) - $startTime

            if ($($timeElapsed).TotalMinutes -ge 1)
            {
                throw "Docker Daemon did not start successfully within 1 minute."
            } 

            # Swallow error and try again
            Start-Sleep -sec 1
        }
    }
    Write-Log "Successfully connected to Docker Daemon."
}


function 
Write-DockerImageTag()
{
    [CmdletBinding()]
    param(
        [Parameter(Mandatory=$true)]
        [string]
        $BaseImageName
    )

    $dockerOutput = Find-DockerImages $BaseImageName

    if ($dockerOutput.Count -gt 1)
    {
        Write-Log "Base image is already tagged:"
    }
    else
    {
        if ($dockerOutput.Count -lt 1)
        {
            #
            # Docker restart required if the image was installed after Docker was 
            # last started
            #
            Stop-Docker
            Start-Docker

            $dockerOutput = Find-DockerImages $BaseImageName

            if ($dockerOutput.Count -lt 1)
            {
                throw "Could not find Docker image to match '$BaseImageName'"
            }
        }

        if ($dockerOutput.Count -gt 1)
        {
            Write-Log "Base image is already tagged:"
        }
        else
        {
            #
            # Register the base image with Docker
            #
            $imageId = ($dockerOutput -split "\s+")[2]

            Write-Log "Tagging new base image ($imageId)..."
            
            docker tag $imageId "$($BaseImageName.tolower()):latest"
            Write-Log "Base image is now tagged:"

            $dockerOutput = Find-DockerImages $BaseImageName
        }
    }
    
    Write-Log $dockerOutput
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

# Get Node IPV4 Address
function Get-IPAddress()
{
    return (Get-NetIPAddress | where {$_.IPAddress -Like '10.*' -and $_.AddressFamily -eq 'IPV4'})[0].IPAddress
}

# Update Docker Config to have cluster-store=consul:// address configured for Swarm cluster.
function Write-DockerStartupScriptWithSwarmClusterInfo()
{
    $dataDir = $env:ProgramData

    # create the target directory
    $targetDir = $dataDir + '\docker'
    if(!(Test-Path -Path $targetDir )){
        New-Item -ItemType directory -Path $targetDir
    }

    $ipAddress = Get-IPAddress
    $OutFile = @"
@echo off
set certs=%ProgramData%\docker\certs.d

if exist %ProgramData%\docker (goto :run)
mkdir %ProgramData%\docker

:run
if not exist %systemroot%\system32\dockerd.exe (goto :legacy)

dockerd -H npipe:// -H 0.0.0.0:2375 --cluster-store=consul://$($SwarmMasterIP):8500 --cluster-advertise=$($ipAddress):2375
goto :eof

:legacy
docker daemon -H npipe:// -H 0.0.0.0:2375 --cluster-store=consul://$($SwarmMasterIP):8500 --cluster-advertise=$($ipAddress):2375
goto :eof
"@

    $OutFile | Out-File -encoding ASCII -filepath "$targetDir\runDockerDaemon.cmd"
}

try
{
    Write-Log "Install Windows Container feature and Docker"
    Install-ContainerHost

    Write-Log "Stop Docker"
    Stop-Docker

    Write-Log "Opening firewall ports"
    Open-FirewallPorts

    Write-Log "Setup Docker Startup Script With Swarm Cluster Info"
    Write-DockerStartupScriptWithSwarmClusterInfo

    Write-Log "Start Docker"
    Start-Docker

    Write-Log "Setup Complete"
}
catch 
{
    Write-Error $_
}