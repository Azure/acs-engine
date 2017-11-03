# Make sure to git clone Kubernetes repo with symlink
# git clone -c core.symlinks=true https://github.com/Azure/kubernetes ${GOPATH}/src/k8s.io/kubernetes
$k8spath = Join-Path -Path $Env:GOPATH -ChildPath "src\k8s.io\kubernetes"
if (!(Test-Path -Path $k8spath))
{
	Write-Host "Kubernetes path $k8spath does not exist!"
	exit
}

$makefile = Join-Path -Path $k8spath -ChildPath "Makefile"
if (!((Get-Item $makefile).Attributes.ToString() -match "ReparsePoint"))
{
	Write-Host "Kubernetes Makefile $makefile is not a symlink!"
	Write-Host "Please use '-c core.symlinks=true' option in git clone."
	Write-Host "E.g., git clone -c core.symlinks=true https://github.com/Azure/kubernetes `${GOPATH}/src/k8s.io/kubernetes"
	exit
}

Get-Content Dockerfile.k8swin | docker build --pull -t k8swin -
docker run --security-opt seccomp:unconfined -it `
	-v ${k8spath}:/gopath/src/k8s.io/kubernetes `
	-w /gopath/src/k8s.io/kubernetes `
		k8swin /bin/bash
