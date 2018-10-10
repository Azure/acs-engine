# This is a temporary file to test dot-sourcing functions stored in separate scripts in a zip file

filter Timestamp {"$(Get-Date -Format o): $_"}

function
Write-Log($message)
{
    $msg = $message | Timestamp
    Write-Output $msg
}

function DownloadFileOverHttp($Url, $DestinationPath)
{
    $secureProtocols = @()
    $insecureProtocols = @([System.Net.SecurityProtocolType]::SystemDefault, [System.Net.SecurityProtocolType]::Ssl3)

    foreach ($protocol in [System.Enum]::GetValues([System.Net.SecurityProtocolType]))
    {
        if ($insecureProtocols -notcontains $protocol)
        {
            $secureProtocols += $protocol
        }
    }
    [System.Net.ServicePointManager]::SecurityProtocol = $secureProtocols

    Invoke-WebRequest $Url -UseBasicParsing -OutFile $DestinationPath -Verbose
    Write-Log "Downloaded file to $DestinationPath"
}