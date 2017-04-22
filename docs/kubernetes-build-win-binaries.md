# Microsoft Azure Container Service Engine - Build Windows Kubernetes Binaries

## Building Windows Kubernetes Binaries and deploy to an Azure storage account

The following instructions show how to deploy the Windows Kubernetes Binaries and deploy them to an Azure Storage Account.

1. Deploy a linux VM

2. Install docker
   ```
   sudo wget --tries 4 --retry-connrefused --waitretry=15 -qO- https://get.docker.com | sh
   sudo usermod -aG docker $USER
   newgrp docker
   ```

3. Update linux
   ```
   sudo apt-get update
   sudo apt-get install -y linux-image-extra-$(uname -r) linux-image-extra-virtual apt-transport-https ca-certificates make gcc gcc-aarch64-linux-gnu zip
   ```

4. Install Go - this must be installed in local user account since kubernetes cross compile install binaries.
   ```
   wget -qO- https://storage.googleapis.com/golang/go1.7.3.linux-amd64.tar.gz | tar zx -C $HOME
   echo "export GOPATH=$HOME/gopath" >> $HOME/.profile
   echo 'export PATH=$PATH:$HOME/go/bin:$GOPATH/bin' >> $HOME/.profile
   echo 'export GOROOT=$HOME/go' >> $HOME/.profile
   source $HOME/.profile
   ```

5. Install Go Dependencies
   ```
   go get -u github.com/jteeuwen/go-bindata/go-bindata
   ```

6. Build windows kubelet.exe
   ```
   cd $HOME
   git clone https://github.com/kubernetes/kubernetes $GOPATH/src/k8s.io/kubernetes
   cd $GOPATH/src/k8s.io/kubernetes
   make WHAT=cmd/kubelet
   make WHAT=cmd/kubelet KUBE_BUILD_PLATFORMS=windows/amd64
   ```
7. Build windows kube-proxy.exe
   ```
   cd $GOPATH/src/k8s.io/kubernetes
   make WHAT=cmd/kube-proxy
   make WHAT=cmd/kube-proxy KUBE_BUILD_PLATFORMS=windows/amd64
   ```

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
   ```
   export AZURE_STORAGE_CONNECTION_STRING="DefaultEndpointsProtocol=https;AccountName=enter_your_account;AccountKey=enter_your_key"
   azure storage blob upload k.zip k8scontainer k.zip -q
   ```
