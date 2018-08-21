$cert = New-SelfSignedCertificate -DnsName (hostname) -CertStoreLocation Cert:\LocalMachine\My
winrm create winrm/config/Listener?Address=*+Transport=HTTPS "@{Hostname=`"$(hostname)`"; CertificateThumbprint=`"$($cert.Thumbprint)`"}"
winrm set winrm/config/service/auth "@{Basic=`"true`"}"
New-NetFirewallRule -DisplayName "Windows Remote Management (HTTPS-In)" -Name "WINRM-HTTP-In-TCP-Any" -Profile Any -LocalPort 5986 -Protocol TCP