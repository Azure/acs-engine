#!/bin/bash

source /opt/azure/containers/provision_source.sh
source /opt/azure/dcos/environment

TMPDIR="/tmp/dcos"
mkdir -p $TMPDIR

# default dc/os component download address (Azure CDN)
LIBIPSET_DOWNLOAD_URL=https://dcos-mirror.azureedge.net/pkg/libipset3_6.29-1_amd64.deb
IPSET_DOWNLOAD_URL=https://dcos-mirror.azureedge.net/pkg/ipset_6.29-1_amd64.deb
UNZIP_DOWNLOAD_URL=https://dcos-mirror.azureedge.net/pkg/unzip_6.0-20ubuntu1_amd64.deb
LIBLTDL_DOWNLOAD_URL=https://dcos-mirror.azureedge.net/pkg/libltdl7_2.4.6-0.1_amd64.deb
DOCKER_CE_DOWNLOAD_URL=https://dcos-mirror.azureedge.net/pkg/docker-ce_17.09.0~ce-0~ubuntu_amd64.deb
SELINUX_DOWNLOAD_URL=https://dcos-mirror.azureedge.net/pkg/selinux-utils_2.4-3build2_amd64.deb

case $DCOS_ENVIRONMENT in
    # because of Chinese GreatWall Firewall, the default packages on Azure CDN is blocked. So the following Chinese local mirror url should be used instead.
    AzureChinaCloud)
        LIBIPSET_DOWNLOAD_URL=http://acsengine.blob.core.chinacloudapi.cn/dcos/libipset3_6.29-1_amd64.deb
        IPSET_DOWNLOAD_URL=http://acsengine.blob.core.chinacloudapi.cn/dcos/ipset_6.29-1_amd64.deb
        UNZIP_DOWNLOAD_URL=http://acsengine.blob.core.chinacloudapi.cn/dcos/unzip_6.0-20ubuntu1_amd64.deb
        LIBLTDL_DOWNLOAD_URL=http://acsengine.blob.core.chinacloudapi.cn/dcos/libltdl7_2.4.6-0.1_amd64.deb
        DOCKER_CE_DOWNLOAD_URL=http://mirror.kaiyuanshe.cn/docker-ce/linux/ubuntu/dists/xenial/pool/stable/amd64/docker-ce_17.09.0~ce-0~ubuntu_amd64.deb
        SELINUX_DOWNLOAD_URL=http://mirror.kaiyuanshe.cn/docker-ce/linux/ubuntu/dists/xenial/pool/stable/amd64/selinux-utils_2.4-3build2_amd64.deb
    ;;
esac

for url in $LIBIPSET_DOWNLOAD_URL $IPSET_DOWNLOAD_URL $UNZIP_DOWNLOAD_URL $LIBLTDL_DOWNLOAD_URL $DOCKER_CE_DOWNLOAD_URL $SELINUX_DOWNLOAD_URL; do
  retry_get_install_deb 10 10 120 $url
  if [ $? -ne 0  ]; then
    exit 1
  fi
done

retrycmd_if_failure 10 10 120 curl -fsSL -o $TMPDIR/dcos_install.sh http://BOOTSTRAP_IP:8086/dcos_install.sh
if [ $? -ne 0  ]; then
  exit 1
fi
