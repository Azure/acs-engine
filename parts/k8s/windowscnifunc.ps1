function Get-HnsPsm1
{
    Param(
        [string]
        HnsUrl = "https://github.com/Microsoft/SDN/raw/master/Kubernetes/windows/hns.psm1",
        [string]
        [Parameter(Mandatory=$true)]
        HNSModule
    )
    DownloadFileOverHttp $HnsUrl "$HNSModule"
}

function Update-WinCNI
{
    Param(
        [string]
        WinCniUrl = "https://github.com/Microsoft/SDN/raw/master/Kubernetes/windows/cni/wincni.exe",
        [string]
        [Parameter(Mandatory=$true)]
        CNIPath
    )
    $wincni = "wincni.exe"
    $wincniFile = [Io.path]::Combine($CNIPath, $wincni)
    DownloadFileOverHttp $WinCniUrl $wincniFile
}

# TODO: Move the code that creates the wincni configuration file out of windowskubeletfunc.ps1 and put it here