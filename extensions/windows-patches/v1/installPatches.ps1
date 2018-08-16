
# Return codes:
#  0 - success
#  1 - install failure
#  2 - download failure

param(
    [string[]] $URIs # TODO: this may need some more parsing
)

function DownloadFile([string] $URI)
{
    try {
        $fileName = Split-Path $URI -Leaf
        $fullName = [io.path]::Combine($env:TEMP, $filename)
        Invoke-WebRequest -UseBasicParsing $URI -OutFile $fullName
        return $fullName
    } catch {
        Write-Error $_.Exception.Message
        exit 2
    }
}


$URIs | ForEach-Object {
    $ext = [io.path]::GetExtension($_)
    switch ($ext) {
        ".exe" { 
            $localPath = DownloadFile($_)
            Write-Host "Starting $localPath"
            $proc = Start-Process -Passthru -FilePath "$localPath" /q /norestart
            Wait-Process -InputObject $proc
            if ($proc.ExitCode -eq 0)
            {
                Write-Host "Finished running $localPath"
            } else {
                Write-Error "Error running $localPath, exitcode $($proc.ExitCode)"
            }
        }
        ".msu" {
            Write-Error "MSU Not Implemented Yet"
            exit 1
        }
        Default {
            Write-Error "Cannot install $_ - unknown file type"
            exit 1
        }
    }
}

exit 0