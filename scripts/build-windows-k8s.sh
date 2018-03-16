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
	if [ ! "${version}" \< "1.7.10" ]; then
		echo "version 1.7.10 and after..."
		# In 1.7.10, the following commit is not needed and has conflict with 137f4cb16e
		# due to the out-of-order back porting into Azure 1.7. So removing it.
		# cee32e92f7 fix#50150: azure disk mount failure on coreos
		git revert --no-edit cee32e92f7 || true

		if [ ! "${version}" \< "1.7.12" ]; then
			echo "version 1.7.12 and after..."
			# In 1.7.12, the following commits are cherry-picked in upstream and has conflict
			# with 137f4cb16e. So removing them.
			git revert --no-edit 593653c384 || true #only for linux
			git revert --no-edit 7305738dd1 || true #add tests only
			git revert --no-edit e01bafcf80 || true #only for linux
			git revert --no-edit afd79db7a6 || true #only for linux
			git revert --no-edit 3a4abca2f7 || true #covered by commit 3aa179744f 
			git revert --no-edit 6a2e2f47d3 || true #covered by commit 3aa179744f

			if [ ! "${version}" \< "1.7.13" ]; then
				echo "version 1.7.13 and after..."
				# In 1.7.13, the following commit is cherry-picked in upstream and has conflict
				# with 137f4cb16e. So removing it.
				git revert --no-edit 3aa179744f || true #only for linux

				if [ ! "${version}" \< "1.7.14" ]; then
					echo "version 1.7.14 and after..."
					# From 1.7.14, 975d0a4bb9 conflict with 137f4cb16e. We use a36d59ddda in Azure repo instead of 975d0a4bb9
					git revert --no-edit 975d0a4bb9 || true
					# From 1.7.14, 51584188ee and 3a0db21dcb conflict with af3a93b07e. We use f5c45d3def in Azure repo instead of them
					git revert --no-edit 51584188ee || true
					git revert --no-edit 3a0db21dcb || true
					# From 1.7.14, 273411cc90 conflict with e9591ef03e. We use d18812a049 in Azure repo instead of 273411cc90
					git revert --no-edit 273411cc90 || true
					# From 1.7.14, a97f60fbf3 conflict with a36d59ddda. We use 69c56c6037 in Azure repo instead of a97f60fbf3
					git revert --no-edit a97f60fbf3 || true
				fi
			fi
		fi
	fi

	# 32ceaa7918 fix #60625: add remount logic for azure file plugin on Windows
	# ...
	# b8fe713754 Use adapter vEthernet (HNSTransparent) on Windows host network to find node IP
	# From 1.7.13, acbdec96da is not needed since 060111c603 supercedes it
	git cherry-pick --allow-empty --keep-redundant-commits b8fe713754^..e912889a7f
	if [ "${version}" \< "1.7.13" ]; then
		git cherry-pick --allow-empty --keep-redundant-commits acbdec96da
	fi
	git cherry-pick --allow-empty --keep-redundant-commits 76d7c23f62^..32ceaa7918

	if [ ! "${version}" \< "1.7.14" ]; then
		# From 1.7.14, 975d0a4bb9 conflict with 137f4cb16e. We use a36d59ddda in Azure repo instead of 975d0a4bb9
		git cherry-pick --allow-empty --keep-redundant-commits a36d59ddda
		# From 1.7.14, 51584188ee and 3a0db21dcb conflict with af3a93b07e. We use f5c45d3def in Azure repo instead of them
		git cherry-pick -X theirs --allow-empty --keep-redundant-commits f5c45d3def # git complains about a conflict and just take theirs
		# From 1.7.14, 273411cc90 conflict with e9591ef03e. We use d18812a049 in Azure repo instead of 273411cc90
		git cherry-pick --allow-empty --keep-redundant-commits d18812a049
		# From 1.7.14, a97f60fbf3 conflict with a36d59ddda. We use 69c56c6037 in Azure repo instead of a97f60fbf3
		git cherry-pick --allow-empty --keep-redundant-commits 69c56c6037

		git cherry-pick --allow-empty --keep-redundant-commits 3e930be6bc
	fi
}

k8s_18_cherry_pick() {
	# 4fd355d04a fix #60625: add remount logic for azure file plugin on Windows
	# ...
	# 4647f2f616 merge #52401: add windows implementation of GetMountRefs

	# d75ef50170 merge#54334: fix azure disk mount failure on coreos and some other distros

	# Only cherry pick before 1.8.6 due to merge conflict with 1.8.6 and up: b8594873f4 merge #50673: azure disk fix of SovereignCloud support

	# cb29df51c0 Fixing 'targetport' to service 'port' mapping
	# ...
	# b42981f90b merge#53629: fix azure file mount limit issue on windows due to using drive letter

	# We have to skip 8b3345114e due to merge conflict with 1.8.4 and up, 4647f2f616 has the change.

	# 8d477271f7 update bazel
	# ...
	# 69644018c8 Use adapter vEthernet (HNSTransparent) on Windows host network to find node IP

	if [ ! "${version}" \< "1.8.9" ]; then
		# From 1.8.9, 63b4f60e43 conflict with b42981f90b. We use 6a8305e419 in Azure repo instead of 63b4f60e43
		git revert --no-edit 63b4f60e43 || true
		# From 1.8.9, 40d5e0a34f conflict with 6a8305e419. We use b90d61a48c in Azure repo instead of 40d5e0a34f
		git revert --no-edit 40d5e0a34f || true
	fi

	git cherry-pick --allow-empty --keep-redundant-commits 69644018c8^..8d477271f7
	git cherry-pick --allow-empty --keep-redundant-commits b42981f90b^..cb29df51c0
	if [ "${version}" \< "1.8.6" ]; then
		git cherry-pick --allow-empty --keep-redundant-commits b8594873f4
	fi
	git cherry-pick --allow-empty --keep-redundant-commits d75ef50170
	git cherry-pick --allow-empty --keep-redundant-commits 4647f2f616^..4fd355d04a

	if [ ! "${version}" \< "1.8.9" ]; then
		# From 1.8.9, 63b4f60e43 conflict with b42981f90b. We use 6a8305e419 in Azure repo instead of 63b4f60e43
		git cherry-pick --allow-empty --keep-redundant-commits 6a8305e419
		# From 1.8.9, 40d5e0a34f conflict with 6a8305e419. We use b90d61a48c in Azure repo instead of 40d5e0a34f
		git cherry-pick --allow-empty --keep-redundant-commits b90d61a48c
	fi
}

apply_acs_cherry_picks() {
	if [ "${KUBERNETES_RELEASE}" == "1.6" ]; then
		k8s_16_cherry_pick
	elif [ "${KUBERNETES_RELEASE}" == "1.7" ]; then
		k8s_17_cherry_pick
	elif [ "${KUBERNETES_RELEASE}" == "1.8" ]; then
		k8s_18_cherry_pick
	elif [ "${KUBERNETES_RELEASE}" == "1.9" ]; then
		echo "No need to cherry-pick for 1.9!"
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

download_wincni() {
	mkdir -p ${DIST_DIR}/cni/config
	az storage blob download -f ${DIST_DIR}/cni/wincni.exe -c ${AZURE_STORAGE_CONTAINER_NAME} -n wincni.exe
}

copy_dockerfile_and_hns_psm1() {
  cp ${ACS_ENGINE_HOME}/windows/* ${DIST_DIR}
}

create_zip() {
	cd ${DIST_DIR}/..
	zip -r ../v${ACS_VERSION}int.zip k/*
	cd -
}

upload_zip_to_blob_storage() {
	az storage blob upload -f ${DIST_DIR}/../../v${ACS_VERSION}int.zip -c ${AZURE_STORAGE_CONTAINER_NAME} -n v${ACS_VERSION}int.zip
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
download_wincni
copy_dockerfile_and_hns_psm1
create_zip
upload_zip_to_blob_storage
push_acs_branch
