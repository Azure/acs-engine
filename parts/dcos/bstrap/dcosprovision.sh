#!/bin/bash

source /opt/azure/containers/provision_source.sh
source /opt/azure/dcos/environment

TMPDIR="/tmp/dcos"
mkdir -p $TMPDIR

# default dc/os component download address (Azure CDN)
packages=(
  https://dcos-mirror.azureedge.net/pkg/libipset3_6.29-1_amd64.deb
  https://dcos-mirror.azureedge.net/pkg/ipset_6.29-1_amd64.deb
  https://dcos-mirror.azureedge.net/pkg/unzip_6.0-20ubuntu1_amd64.deb
  https://dcos-mirror.azureedge.net/pkg/libltdl7_2.4.6-0.1_amd64.deb
  https://dcos-mirror.azureedge.net/pkg/docker-ce_17.09.0~ce-0~ubuntu_amd64.deb
  https://dcos-mirror.azureedge.net/pkg/selinux-utils_2.4-3build2_amd64.deb
)

# sha1sum checksums for @packages
sha1sums=(
  f88d09688291917c8bb65682fea9f5d571ec8d6a
  807dc11f5bfa39bb4b0dc9024fc51bb309905a21
  57ae2bb6ded1fdf91b6d518294134df1ff13fcca
  9a0f9f2769d3dc834737aa7df50aaaea369af98d
  94f6e89be6d45d9988269a237eb27c7d6a844d7f
  77bdb5847060845c0a158f567b1ddd7fa34b7236
)

case $DCOS_ENVIRONMENT in
  # because of Chinese GreatWall Firewall, the default packages on Azure CDN is blocked. So the following Chinese local mirror url should be used instead.
  AzureChinaCloud)
    packages=(
      http://acsengine.blob.core.chinacloudapi.cn/dcos/libipset3_6.29-1_amd64.deb
      http://acsengine.blob.core.chinacloudapi.cn/dcos/ipset_6.29-1_amd64.deb
      http://acsengine.blob.core.chinacloudapi.cn/dcos/unzip_6.0-20ubuntu1_amd64.deb
      http://acsengine.blob.core.chinacloudapi.cn/dcos/libltdl7_2.4.6-0.1_amd64.deb
      http://mirror.kaiyuanshe.cn/docker-ce/linux/ubuntu/dists/xenial/pool/stable/amd64/docker-ce_17.09.0~ce-0~ubuntu_amd64.deb
      http://acsengine.blob.core.chinacloudapi.cn/dcos/selinux-utils_2.4-3build2_amd64.deb
    )
  ;;
esac

len=$((${#packages[@]}-1))
for i in $(seq 0 $len); do
  retry_get_install_deb 10 10 120 ${packages[$i]} ${sha1sums[$i]}
    if [ $? -ne 0  ]; then
    exit 1
  fi
done

retrycmd_if_failure 10 10 120 curl -fsSL -o $TMPDIR/dcos_install.sh http://BOOTSTRAP_IP:8086/dcos_install.sh
if [ $? -ne 0  ]; then
  exit 1
fi
