#!/bin/bash

source /opt/azure/containers/provision_source.sh
source /opt/azure/dcos/environment

# default dc/os component download address (Azure CDN)
LIBLTDL_DOWNLOAD_URL=https://dcos-mirror.azureedge.net/pkg/libltdl7_2.4.6-0.1_amd64.deb
DOCKER_CE_DOWNLOAD_URL=https://dcos-mirror.azureedge.net/pkg/docker-ce_17.09.0~ce-0~ubuntu_amd64.deb

case $DCOS_ENVIRONMENT in
    # because of Chinese GreatWall Firewall, the default packages on Azure CDN is blocked. So the following Chinese local mirror url should be used instead.
    AzureChinaCloud)
        LIBLTDL_DOWNLOAD_URL=http://acsengine.blob.core.chinacloudapi.cn/dcos/libltdl7_2.4.6-0.1_amd64.deb
        DOCKER_CE_DOWNLOAD_URL=http://mirror.kaiyuanshe.cn/docker-ce/linux/ubuntu/dists/xenial/pool/stable/amd64/docker-ce_17.09.0~ce-0~ubuntu_amd64.deb
    ;;
esac

for url in $LIBLTDL_DOWNLOAD_URL $DOCKER_CE_DOWNLOAD_URL; do
  retry_get_install_deb 10 10 120 $url
  if [ $? -ne 0  ]; then
    exit 1
  fi
done
