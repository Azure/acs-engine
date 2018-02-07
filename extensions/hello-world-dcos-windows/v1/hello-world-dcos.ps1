# Script file to run hello-world in dcos

#!/usr/bin/pwsh


Write-Host "$(date) - Starting Script"

# Deploy container
Write-Host "$(date)  - Deploying hello-world"

$uri =  "http://"+($env:DCOS_AGENT_IP)+":5051/metrics/snapshot"
while($true) {
    $obj = ((Invoke-Webrequest -Method GET -URI $uri ).Content | ConvertFrom-JSON )
    Write-Host "$(date) - system/cpus_total = " ($obj.'system/cpus_total') ", mem free bytes = " ($obj.'system/mem_free_bytes') ", mem total bytes = " ($obj.'system/mem_total_bytes') 
    sleep 5
}

Write-Host "$(date) - view resources in mesos UI to validate"
Write-Host "$(date) - Script complete"

