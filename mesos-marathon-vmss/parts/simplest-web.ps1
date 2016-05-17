<#code used from https://gist.github.com/wagnerandrade/5424431#>
$ip = (Get-NetIPAddress | where {$_.IPAddress -Like '*.*.*.*'})[0].IPAddress
$url = 'http://{0}:80/' -f $ip
$listener = New-Object System.Net.HttpListener
$listener.Prefixes.Add($url)
$listener.Start()

Write-Host('Listening at {0}...' -f $url)

while ($listener.IsListening) {
    $context = $listener.GetContext()
    $requestUrl = $context.Request.Url
    $clientIP = $context.Request.RemoteEndPoint.Address
    $response = $context.Response
    Write-Host ''
    Write-Host('> {0}' -f $requestUrl)
    $content='<html><body><H1>helloworld</H1></body></html>'
    Write-Output $content
    $buffer = [System.Text.Encoding]::UTF8.GetBytes($content)
    $response.ContentLength64 = $buffer.Length
    $response.OutputStream.Write($buffer, 0, $buffer.Length)
    $response.Close()
    $responseStatus = $response.StatusCode
    Write-Host('< {0}' -f $responseStatus)
  }
