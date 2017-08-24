#!/bin/bash
set -eo pipefail

ACS_ENGINE_HOME=${GOPATH}/src/github.com/Azure/acs-engine

usage() {
	echo "$0 [-v version] [-p acs_patch_version]"
	echo " -v <version>: version"
	echo " -p <patched version>: acs_patch_version"
	exit 0
}

while getopts ":v:p:" opt; do
  case ${opt} in
    v)
      version=${OPTARG}
      ;;
    p)
      acs_patch_version=${OPTARG}
      ;;
    *)
			usage
      ;;
  esac
done

if [ -z "${version}" ] || [ -z "${acs_patch_version}" ]; then
    usage
fi

KUBERNETES_RELEASE=$(echo $version | cut -d'.' -f1,2)
KUBERNETES_RELEASE_BRANCH_NAME=release-${KUBERNETES_RELEASE}
ACS_VERSION=${version}-${acs_patch_version}
ACS_BRANCH_NAME=acs-v${ACS_VERSION}
DIST_DIR=${ACS_ENGINE_HOME}/_dist/k8s-windows-v${ACS_VERSION}

fetch_k8s() {
	git clone https://github.com/Azure/kubernetes ${GOPATH}/src/k8s.io/kubernetes || true
	cd ${GOPATH}/src/k8s.io/kubernetes
	git remote add upstream https://github.com/kubernetes/kubernetes || true
	git fetch upstream
}

set_git_config() {
	git config user.name "Sean Knox"
	git config user.email "sean.knox@microsoft.com"
}

create_version_branch() {
	git checkout -b ${ACS_BRANCH_NAME} ${KUBERNETES_RELEASE_BRANCH_NAME} || true
}

create_dist_dir() {
	mkdir -p ${DIST_DIR}
}

build_kubelet() {
	echo "building kubelet.exe..."
	build/run.sh make WHAT=cmd/kubelet KUBE_BUILD_PLATFORMS=windows/amd64
	cp ${GOPATH}/src/k8s.io/kubernetes/_output/dockerized/bin/windows/amd64/kubelet.exe ${DIST_DIR}
}

build_kubeproxy() {
	echo "building kube-proxy.exe..."
	build/run.sh make WHAT=cmd/kube-proxy KUBE_BUILD_PLATFORMS=windows/amd64
	cp ${GOPATH}/src/k8s.io/kubernetes/_output/dockerized/bin/windows/amd64/kube-proxy.exe ${DIST_DIR}
}

download_kubectl() {
	kubectl="https://storage.googleapis.com/kubernetes-release/release/v${version}/bin/windows/amd64/kubectl.exe"
	echo "dowloading ${kubectl} ..."
	wget ${kubectl} -P k
	curl ${kubectl} -o ${DIST_DIR}/kubectl.exe
	chmod 775 ${DIST_DIR}/kubectl.exe
}

download_nccm() {
	NSSM_VERSION=2.24
	NSSM_URL=https://nssm.cc/release/nssm-${NSSM_VERSION}.zip
	echo "downloading nssm ..."
	curl ${NSSM_URL} -o /tmp/nssm-${NSSM_VERSION}.zip
	unzip -q -d /tmp /tmp/nssm-${NSSM_VERSION}.zip
	cp /tmp/nssm-${NSSM_VERSION}/win64/nssm.exe ${DIST_DIR}
	chmod 775 ${DIST_DIR}/nssm.exe
	rm -rf /tmp/nssm-${NSSM_VERSION}*
}

download_winnat() {
	echo "copying winnat.sys ..."
	cp $HOME/winnat/winnat.sys ${DIST_DIR}
}

copy_dockerfile_and_pause_ps1() {
  cp ${ACS_ENGINE_HOME}/windows/* ${DIST_DIR}
}

create_zip() {
	zip -r ${DIST_DIR}/../v${version}intwinnat.zip ${DIST_DIR}/*
}

create_dist_dir
fetch_k8s
set_git_config
create_version_branch

# Due to what appears to be a bug in the Kubernetes Windows build system, one
# has to first build a linux binary to generate _output/bin/deepcopy-gen.
# Building to Windows w/o doing this will generate an empty deepcopy-gen.
build/run.sh make WHAT=cmd/kubelet KUBE_BUILD_PLATFORMS=linux/amd64

build_kubelet
build_kubeproxy
download_kubectl
download_nccm
copy_dockerfile_and_pause_ps1
create_zip
