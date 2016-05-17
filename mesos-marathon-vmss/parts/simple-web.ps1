<#code used from https://gist.github.com/wagnerandrade/5424431#>
$ip = (Get-NetIPAddress | where {$_.IPAddress -Like '*.*.*.*'})[0].IPAddress
$url = "http://"+$ip+":80/"
$listener = New-Object System.Net.HttpListener
$listener.Prefixes.Add($url)
$listener.Start()
$callerCounts = @{}

Write-Host('Listening at {0}...' -f $url)

while ($listener.IsListening) {
    $context = $listener.GetContext()
    $requestUrl = $context.Request.Url
    $clientIP = $context.Request.RemoteEndPoint.Address
    $response = $context.Response

    Write-Host ''
    Write-Host('> {0}' -f $requestUrl)

    $count = 1
    $k=$callerCounts.Get_Item($clientIP)
    if ($k -ne $null) { $count += $k }
    $callerCounts.Set_Item($clientIP, $count)
    $header="<html><body><H1>Windows Container Web Server</H1>"
    $callerCountsString=""
    $callerCounts.Keys | % { $callerCountsString+='<p>IP {0} callerCount {1} ' -f $_,$callerCounts.Item($_) }
    $footer="</body></html>"
    $content='{0}{1}{2}' -f $header,$callerCountsString,$footer
    Write-Output $content
    $buffer = [System.Text.Encoding]::UTF8.GetBytes($content)
    $response.ContentLength64 = $buffer.Length
    $response.OutputStream.Write($buffer, 0, $buffer.Length)
    $response.Close()

    $responseStatus = $response.StatusCode
    Write-Host('< {0}' -f $responseStatus)
  }
