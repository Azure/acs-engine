<#
    .SYNOPSIS
        Provisions VM as a DCOS agent.

    .DESCRIPTION
        Provisions VM as a DCOS agent.

     Invoke by:
       
#>

[CmdletBinding(DefaultParameterSetName="Standard")]
param(
    [string]
    [ValidateNotNullOrEmpty()]
    $masterCount,

    [string]
    [ValidateNotNullOrEmpty()]
    $firstMasterIP,
    
    [string]
    [ValidateNotNullOrEmpty()]
    $bootstrapUri,

    [parameter()]
    [ValidateNotNullOrEmpty()]
    $isAgent,

    [parameter()]
    [ValidateNotNullOrEmpty()]
    $subnet,

    [parameter()]
    [AllowNull()]
    $isPublic = $false,

    [string]
    [AllowNull()]
    $customAttrs = ""
)




$global:BootstrapInstallDir = "C:\AzureData"

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
        $shell.Namespace($destination).copyhere($item, 0x14)
    }
}


function 
Remove-Directory($dirname)
{

    try {
        #Get-ChildItem $dirname -Recurse | Remove-Item  -force -confirm:$false
        # This doesn't work because of long file names
        # But this does:
        Invoke-Expression ("cmd /C rmdir /s /q "+$dirname)
    }
    catch {
        # If this fails we don't want it to stop

    }
}


function 
Check-Subnet ([string]$cidr, [string]$ip)
{
    try {

        $network, [int]$subnetlen = $cidr.Split('/')
    
        if ($subnetlen -eq 0)
        {
            $subnetlen = 8 # Default in case we get an IP addr, not CIDR
        }
        $a = ([IPAddress] $network)
        [uint32] $unetwork = [uint32]$a.Address
    
        $mask = -bnot ((-bnot [uint32]0) -shl (32 - $subnetlen))
    
        $a = [IPAddress]$ip
        [uint32] $uip = [uint32]$a.Address
    
        return ($unetwork -eq ($mask -band $uip))
    }
    catch {
        return $false
    }
}

#
# Gets the bootstrap script from the blob store and places it in c:\AzureData

function
Get-BootstrapScript($download_uri, $download_dir)
{
    # Get Mesos Binaries
    $scriptfile = "DCOSWindowsAgentSetup.ps1"

    Write-Log " get script "+ ($download_uri+"/"+$scriptfile) + "and put it "+ ($download_dir+"\"+$scriptfile)

    Invoke-WebRequest -Uri ($download_uri+"/"+$scriptfile) -OutFile ($download_dir+"\"+$scriptfile)
   
    $scriptfile = "packages.ps1"
    Write-Log " get package file "+ ($download_uri+"/"+$scriptfile) + "and put it "+ ($download_dir+"\"+$scriptfile)
    Invoke-WebRequest -Uri ($download_uri+"/"+$scriptfile) -OutFile ($download_dir+"\"+$scriptfile)

}


try
{
    # Set to false for debugging.  This will output the start script to
    # c:\AzureData\dcosProvisionScript.log, and then you can RDP 
    # to the windows machine, and run the script manually to watch
    # the output.
    Write-Log "Get the install script"

    Write-Log ("Parameters = isAgent = ["+ $isAgent + "] mastercount = ["+$MasterCount + "] First master ip= [" + $firstMasterIp+ "] boostrap URI = ["+ $bootstrapUri+"] Subnet = ["+ $subnet +"]" + " -customAttrs " + $customAttrs )

    # Get the boostrap script

    Get-BootstrapScript $bootstrapUri $global:BootstrapInstallDir

    # Convert Master count and first IP to a JSON array of IPAddresses
    $ip = ([IPAddress]$firstMasterIp).getAddressBytes()
    [Array]::Reverse($ip)
    $ip = ([IPAddress]($ip -join '.')).Address

    $MasterIP = @([IPAddress]$null)
    
    for ($i = 0; $i -lt $MasterCount; $i++ ) 
    {
       $new_ip = ([IPAddress]$ip).getAddressBytes()
       [Array]::Reverse($new_ip)
       $new_ip = [IPAddress]($new_ip -join '.')
       $MasterIP += $new_ip
      
       $ip++
     
    }
    $master_str  = $MasterIP.IPAddressToString

    # Add the port numbers
    if ($master_str.count -eq 1) {
        $master_str += ":2181"
    }
    else {
        for ($i = 0; $i -lt $master_str.count; $i++) 
        {
            $master_str[$i] += ":2181"
        }
    }
    $master_json = ConvertTo-Json $master_str
    $master_json = $master_json -replace [Environment]::NewLine,""

    $private_ip = ( Get-NetIPAddress | where { $_.AddressFamily -eq "IPv4" } | where { Check-Subnet $subnet $_.IPAddress } )  # We know the subnet we are on. Makes it easier and more robust
    [Environment]::SetEnvironmentVariable("DCOS_AGENT_IP", $private_ip.IPAddress, "Machine")

    if ($isAgent)
    {
        $run_cmd = $global:BootstrapInstallDir+"\DCOSWindowsAgentSetup.ps1 -MasterIP '$master_json' -AgentPrivateIP "+($private_ip.IPAddress) +" -BootstrapUrl '$bootstrapUri' " 
        if ($isPublic) 
        {
            $run_cmd += " -isPublic:`$true "
        }
        if ($customAttrs) 
        {
            $run_cmd += " -customAttrs '$customAttrs'"
        }
        $run_cmd += ">"+$global:BootstrapInstallDir+"\DCOSWindowsAgentSetup.log 2>&1"
        Write-Log "run setup script $run_cmd"
        Invoke-Expression $run_cmd
    }
    else # We must be deploying a master
    {
        $run_cmd = $global:BootstrapInstallDir+"\DCOSWindowsMasterSetup.ps1 -MasterIP '$master_json' -MasterPrivateIP $privateIP.IPAddress -BootstrapUrl '$bootstrapUri'"
        Write-Log "run setup script $run_cmd"
        Invoke-Expression $run_cmd
    }
}
catch
{
    Write-Error $_
}
