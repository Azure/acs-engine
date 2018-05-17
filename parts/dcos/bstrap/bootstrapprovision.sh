#!/bin/bash

source /opt/azure/dcos/environment

retrycmd_if_failure() { retries=$1; wait=$2; timeout=$3; shift && shift && shift; for i in $(seq 1 $retries); do timeout $timeout ${@}; [ $? -eq 0  ] && break || sleep $wait; done; echo Executed \"$@\" $i times; }

TMPDIR="/tmp/dcos"
mkdir -p $TMPDIR

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

curl -fLsSv --retry 20 -Y 100000 -y 60 -o $TMPDIR/1.deb $LIBLTDL_DOWNLOAD_URL &
curl -fLsSv --retry 20 -Y 100000 -y 60 -o $TMPDIR/2.deb $DOCKER_CE_DOWNLOAD_URL &
wait

retrycmd_if_failure 10 10 120 dpkg -i $TMPDIR/{1,2}.deb
