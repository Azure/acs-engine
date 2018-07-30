$cert = New-SelfSignedCertificate -DnsName (hostname) -CertStoreLocation Cert:\LocalMachine\My
winrm create winrm/config/Listener?Address=*+Transport=HTTPS "@{Hostname=`"$(hostname)`"; CertificateThumbprint=`"$($cert.Thumbprint)`"}"
winrm set winrm/config/service/auth "@{Basic=`"true`"}"