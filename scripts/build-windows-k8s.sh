#!/bin/bash
set -eo pipefail

ACS_ENGINE_HOME=${GOPATH}/src/github.com/Azure/acs-engine

version=""
acs_patch_version=""
winnat="false"

while getopts "hv:p" opt; do
  case $opt in
    h)
      echo "$0 [-v version] [-p acs_patch_version] [-w]"
      echo " -v <version>: version"
      echo " -p <patched version>: acs_patch_version"
      exit 0
      ;;
    v)
      version="${OPTARG}"
      ;;
    p)
      acs_patch_version="${OPTARG}"
      ;;
    \?)
      echo "$0 [-v version] [-p acs_patch_version] [-w]"
      exit 1
      ;;
  esac
done

if [[ -z $version ]]; then
	echo "Unknown or no Kubernetes version provided"
	exit
fi

if [[ -z $acs_patch_version ]]; then
	DIST_DIR=${ACS_ENGINE_HOME}/_dist/k8s-windows-v${version}
else
	DIST_DIR=${ACS_ENGINE_HOME}/_dist/k8s-windows-v${version}-${acs_patch_version}
fi

fetch_k8s() {
	git clone https://github.com/kubernetes/kubernetes ${GOPATH}/src/k8s.io/kubernetes || true
	cd ${GOPATH}/src/k8s.io/kubernetes
	git remote add acs https://github.com/JiangtianLi/kubernetes || true
	git fetch acs
}

set_git_config() {
	git config user.name "Sean Knox"
	git config user.email "sean.knox@microsoft.com"
}

create_version_branch() {
	git checkout -b acs-windows-v${version} v${version}

}

create_dist_dir() {
	mkdir -p ${DIST_DIR}
}

# These are Microsoft patches only present in our fork, they are not yet upstream.
if [[ "$version" =~ "1.6" ]]; then
	STARTING_PATCH_SHA=d74e09bb4e4e7026f45becbed8310665ddcb8514
	ENDING_PATCH_SHA=2adf0591bf70cf8affd61be5f2aa2495676172dd
	echo "Patching for v1.6"
elif [[ "$version" =~ "1.7" ]]; then
	STARTING_PATCH_SHA=41e11915a5289fb7074fdcb2a96cac2e13543845
	ENDING_PATCH_SHA=b0c9ea2463aba41c30a671760c875bf4aaea9845
	echo "Patching for v1.7"
fi

apply_acs_patches() {
	git cherry-pick ${STARTING_PATCH_SHA}..${ENDING_PATCH_SHA}
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
apply_acs_patches

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
