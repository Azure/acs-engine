#!/bin/bash
set -eo pipefail

ACS_ENGINE_HOME=${GOPATH}/src/github.com/Azure/acs-engine

usage() {
	echo "$0 [-v version] [-p acs_patch_version]"
	echo " -v <version>: version"
	echo " -p <patched version>: acs_patch_version"
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
			exit
      ;;
  esac
done

if [ -z "${version}" ] || [ -z "${acs_patch_version}" ]; then
    usage
		exit 1
fi

if [ -z "${AZURE_STORAGE_CONNECTION_STRING}" ] || [ -z "${AZURE_STORAGE_CONTAINER_NAME}" ]; then
    echo '$AZURE_STORAGE_CONNECTION_STRING and $AZURE_STORAGE_CONTAINER_NAME need to be set for upload to Azure Blob Storage.'
		exit 1
fi

KUBERNETES_RELEASE=$(echo $version | cut -d'.' -f1,2)
KUBERNETES_TAG_BRANCH=v${version}
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
	git checkout -b ${ACS_BRANCH_NAME} ${KUBERNETES_TAG_BRANCH} || true
}

k8s_16_cherry_pick() {
	# 232fa6e5bc (HEAD -> release-1.6, origin/release-1.6) Fix the delay caused by network setup in POD Infra container
	# 02b1c2b9e2 Use dns policy to determine setting DNS servers on the correct NIC in Windows container
	# 4c2a2d79aa Fix getting host DNS for Windows node when dnsPolicy is Default
	# caa314ccdc Update docker version parsing to allow nonsemantic versions as they have changed how they do their versions
	# f18be40948 Fix the issue in unqualified name where DNS client such as ping or iwr validate name in response and original question. Switch to use miekg's DNS library
	# c862b583c9 Remove DNS server from NAT network adapter inside container
	# f6c27f9375 Merged libCNI-on-Windows changes from CNI release 0.5.0, PRs 359 and 361
	# 4f196c6cac Fix the issue that ping uses the incorrect NIC to resolve name sometimes
	# 2c9fd27449 Workaround for Outbound Internet traffic in Azure Kubernetes
	# 5fa0725025 Use adapter vEthernet (HNSTransparent) on Windows host network to find node IP

	git cherry-pick 5fa0725025^..232fa6e5bc
}

k8s_17_cherry_pick() {
        # d9f95d1152 Add start time for root container spec
        # ee4e187451 Fix windows docker stats cpu units issue
	# f33e310ede Implement metrics for Windows Containers
	# 08ab2f2636 Use dns policy to determine setting DNS servers on the correct NIC in Windows container
	# 1d4a5780cb Fix network config due to the split of start POD sandbox and start container from 1.7.0
	# 4d00a52142 Workaround for Outbound Internet traffic in Azure Kubernetes (*) Connect a Nat Network to the container (Second adapter) (*) Modify the route so that internet traffic goes via Nat network, and POD traffic goes over the CONTAINER_NETWORK (*) Modify getContainerIP to return the IP corresponding to POD network, and ignore Nat Network (*) DNS Fix for ACS Kubernetes in Windows
	# b8fe713754 Use adapter vEthernet (HNSTransparent) on Windows host network to find node IP

	git cherry-pick b8fe713754^..d9f95d1152
}

k8s_18_cherry_pick() {
        # d07a7d498f Resolve merge conflict
        # f96286ae9e Add start time for root container spec
        # 6d9542b1d6 Fix windows docker stats cpu units issue
        # 6800851a7e Implement metrics for Windows Containers
        # 9802a3178c Use dns policy to determine setting DNS servers on the correct NIC in Windows container
        # cba79368f9 Fix network config due to the split of start POD sandbox and start container from 1.7.0
        # 0e9f3ac18c Workaround for Outbound Internet traffic in Azure Kubernetes (*) Connect a Nat Network to the container (Second adapter) (*) Modify the route so that internet traffic goes via Nat network, and POD traffic goes over the CONTAINER_NETWORK (*) Modify getContainerIP to return the IP corresponding to POD network, and ignore Nat Network (*) DNS Fix for ACS Kubernetes in Windows
        # 69644018c8 Use adapter vEthernet (HNSTransparent) on Windows host network to find node IP

        git cherry-pick 69644018c8^..d07a7d498f
}

apply_acs_cherry_picks() {
	if [ "${KUBERNETES_RELEASE}" == "1.6" ]; then
		k8s_16_cherry_pick
	elif [ "${KUBERNETES_RELEASE}" == "1.7" ]; then
		k8s_17_cherry_pick
        elif [ "${KUBERNETES_RELEASE}" == "1.8" ]; then
		k8s_18_cherry_pick
	else
		echo "Unable to apply cherry picks for ${KUBERNETES_RELEASE}."
		exit 1
	fi
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
	zip -r ../v${ACS_VERSION}intwinnat.zip k/*
	cd -
}

upload_zip_to_blob_storage() {
	az storage blob upload -f ${DIST_DIR}/../../v${ACS_VERSION}intwinnat.zip -c ${AZURE_STORAGE_CONTAINER_NAME} -n v${ACS_VERSION}intwinnat.zip
}

push_acs_branch() {
  cd ${GOPATH}/src/k8s.io/kubernetes
  git push origin ${ACS_BRANCH_NAME}
}

create_dist_dir
fetch_k8s
set_git_config
create_version_branch
apply_acs_cherry_picks

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
push_acs_branch
