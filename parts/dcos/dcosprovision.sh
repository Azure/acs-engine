#!/bin/bash

MESOSDIR=/var/lib/mesos/dl
mkdir $MESOSDIR

# load the env vars
. /etc/mesosphere/setup-flags/dcos-deploy-environment

# default dc/os component download address (Azure CDN)
DOCKER_ENGINE_DOWNLOAD_URL=https://mesosphere.blob.core.windows.net/dcos-deps/docker-engine_1.13.1-0-ubuntu-xenial_amd64.deb
LIBIPSET_DOWNLOAD_URL=https://az837203.vo.msecnd.net/dcos-deps/libipset3_6.29-1_amd64.deb
IPSET_DOWNLOAD_URL=https://az837203.vo.msecnd.net/dcos-deps/ipset_6.29-1_amd64.deb
UNZIP_DOWNLOAD_URL=https://az837203.vo.msecnd.net/dcos-deps/unzip_6.0-20ubuntu1_amd64.deb
LIBLTDL_DOWNLOAD_URL=https://az837203.vo.msecnd.net/dcos-deps/libltdl7_2.4.6-0.1_amd64.deb

case $DCOS_ENVIRONMENT in
    # because of Chinese GreatWall Firewall, the default packages on Azure CDN is blocked. So the following Chinese local mirror url should be used instead.
    AzureChinaCloud)
        DOCKER_ENGINE_DOWNLOAD_URL=http://acsengine.blob.core.chinacloudapi.cn/dcos/docker-engine_1.11.2-0~xenial_amd64.deb
        LIBIPSET_DOWNLOAD_URL=http://acsengine.blob.core.chinacloudapi.cn/dcos/libipset3_6.29-1_amd64.deb
        IPSET_DOWNLOAD_URL=http://acsengine.blob.core.chinacloudapi.cn/dcos/ipset_6.29-1_amd64.deb
        UNZIP_DOWNLOAD_URL=http://acsengine.blob.core.chinacloudapi.cn/dcos/unzip_6.0-20ubuntu1_amd64.deb
        LIBLTDL_DOWNLOAD_URL=http://acsengine.blob.core.chinacloudapi.cn/dcos/libltdl7_2.4.6-0.1_amd64.deb
    ;;
esac

curl -fLsSv --retry 20 -Y 100000 -y 60 -o $MESOSDIR/d.deb $DOCKER_ENGINE_DOWNLOAD_URL &
curl -fLsSv --retry 20 -Y 100000 -y 60 -o $MESOSDIR/1.deb $LIBIPSET_DOWNLOAD_URL &
curl -fLsSv --retry 20 -Y 100000 -y 60 -o $MESOSDIR/2.deb $IPSET_DOWNLOAD_URL &
curl -fLsSv --retry 20 -Y 100000 -y 60 -o $MESOSDIR/3.deb $UNZIP_DOWNLOAD_URL &
curl -fLsSv --retry 20 -Y 100000 -y 60 -o $MESOSDIR/4.deb $LIBLTDL_DOWNLOAD_URL &
wait

for i in {1..300}; do
    dpkg -i $MESOSDIR/{1,2,3,4}.deb
    if [ "$?" = "0" ]
    then
        echo "succeeded"
        break
    fi
    sleep 1
done

ROLESFILECONTENTS

# add Azure update domain and fault domain attributes
ud=$( curl -H Metadata:true "http://169.254.169.254/metadata/instance/compute/platformUpdateDomain?api-version=2017-04-02&format=text" )
fd=$( curl -H Metadata:true "http://169.254.169.254/metadata/instance/compute/platformFaultDomain?api-version=2017-04-02&format=text" )
echo ";azure.faultdomain:$fd;azure.updatedomain:$ud" >> /var/lib/dcos/mesos-slave-common