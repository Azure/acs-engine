# Microsoft Azure Container Service Engine - Build Windows Kubernetes Binaries

## Building Windows Kubernetes Binaries and deploy to an Azure storage account

The following instructions show how to deploy the Windows Kubernetes Binaries and deploy them to an Azure Storage Account.

1. Deploy a linux VM

2. Install docker
 1. `sudo -s`
 2. `wget --tries 4 --retry-connrefused --waitretry=15 -qO- https://get.docker.com | sh`
 3. `sudo usermod -aG docker azureuser - this assumes user of azure user`
 4. `Restart ssh connection to get env settings`

3. Update linux
 1. `sudo apt-get update`
 2. `sudo apt-get install linux-image-extra-$(uname -r) linux-image-extra-virtual apt-transport-https ca-certificates make gcc gcc-aarch64-linux-gnu zip`

4. Install Go - this must be installed in local user account since kubernetes cross compile install binaries.
 1. `wget -qO- https://storage.googleapis.com/golang/go1.7.3.linux-amd64.tar.gz | tar zx -C $HOME`
 2. `echo "export GOPATH=$HOME/gopath" >> $HOME/.profile`
 3. `echo 'export PATH=$PATH:$HOME/go/bin:$GOPATH/bin' >> $HOME/.profile`
 4. `echo 'export GOROOT=$HOME/go' >> $HOME/.profile`
 5. Restart ssh connection to get new env settings

5. Install Go Dependencies
 1. `go get -u github.com/jteeuwen/go-bindata/go-bindata`

6. Build windows kubelet.exe
 1. `cd $HOME`
 2. `git clone https://github.com/kubernetes/kubernetes $GOPATH/src/k8s.io/kubernetes`
 3. `cd $GOPATH/src/k8s.io/kubernetes`
 4. `make WHAT=cmd/kubelet`
 5. `make WHAT=cmd/kubelet KUBE_BUILD_PLATFORMS=windows/amd64`

7. Build windows kube-proxy.exe
 1. `cd $GOPATH/src/k8s.io/kubernetes`
 2. `make WHAT=cmd/kube-proxy`
 3. `make WHAT=cmd/kube-proxy KUBE_BUILD_PLATFORMS=windows/amd64`

8. Build the zip file
 ```
cd $HOME
rm -rf uploadk
mkdir uploadk
cd uploadk
mkdir k
cp $HOME/kubelet/kubernetes/_output/local/bin/windows/amd64/kubelet.exe k
cp $HOME/kube-proxy/kubernetes/_output/local/bin/windows/amd64/kube-proxy.exe k
wget https://storage.googleapis.com/kubernetes-release/release/v1.4.6/bin/windows/amd64/kubectl.exe -P k
chmod 775 k/kubectl.exe
cat <<EOT >> k/Dockerfile
FROM microsoft/windowsservercore

ADD pause.ps1 /pause/pause.ps1
	
CMD powershell /pause/pause.ps1
EOT
cat <<EOT >> k/pause.ps1
while(\$true)
{
    Start-Sleep -Seconds 60
}
EOT
wget https://nssm.cc/release/nssm-2.24.zip
unzip nssm-2.24.zip
cp nssm-2.24/win64/nssm.exe k
rm -rf nssm-2.24*
chmod 775 k/nssm.exe
zip -r k.zip *
  ```

9. Upload to storage.  This assumes you have az cli installed and you have a storage account
 1. `export AZURE_STORAGE_CONNECTION_STRING="DefaultEndpointsProtocol=https;AccountName=enter_your_account;AccountKey=enter_your_key"`
 2. `azure storage blob upload k.zip k8scontainer k.zip -q`