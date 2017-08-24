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

if [ -z "${AZURE_STORAGE_CONNECTION_STRING}" ] || [ -z "${AZURE_STORAGE_CONTAINER_NAME}" ]; then
    echo '$AZURE_STORAGE_CONNECTION_STRING and $AZURE_STORAGE_CONTAINER_NAME need to be set for upload to Azure Blob Storage.'
		exit 1
fi

KUBERNETES_RELEASE=$(echo $version | cut -d'.' -f1,2)
KUBERNETES_RELEASE_BRANCH_NAME=release-${KUBERNETES_RELEASE}
ACS_VERSION=${version}-${acs_patch_version}
ACS_BRANCH_NAME=acs-v${ACS_VERSION}
DIST_DIR=${ACS_ENGINE_HOME}/_dist/k8s-windows-v${ACS_VERSION}/k

fetch_k8s() {
	git clone https://github.com/Azure/kubernetes ${GOPATH}/src/k8s.io/kubernetes || true
	cd ${GOPATH}/src/k8s.io/kubernetes
	git remote add upstream https://github.com/kubernetes/kubernetes || true
	git fetch upstream
}

set_git_config() {
	git config user.name "ACS CI"
	git config user.email "containers@microsoft.com"
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

download_nssm() {
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
	az storage blob download -f ${DIST_DIR}/winnat.sys -c ${AZURE_STORAGE_CONTAINER_NAME} -n winnat.sys
}

copy_dockerfile_and_pause_ps1() {
  cp ${ACS_ENGINE_HOME}/windows/* ${DIST_DIR}
}

create_zip() {
	cd ${DIST_DIR}/..
	zip -r ../v${version}intwinnat.zip k/*
	cd -
}

upload_zip_to_blob_storage() {
	az storage blob upload -f ${DIST_DIR}/../../v${version}intwinnat.zip -c ${AZURE_STORAGE_CONTAINER_NAME} -n v${version}intwinnat.zip
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
download_nssm
download_winnat
copy_dockerfile_and_pause_ps1
create_zip
upload_zip_to_blob_storage
