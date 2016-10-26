#!/bin/bash

MESOSDIR=/var/lib/mesos/dl
mkdir $MESOSDIR

# load the env vars
. /etc/mesosphere/setup-flags/bootstrap-id

curl -fLsSv --retry 20 -Y 100000 -y 60 -o $MESOSDIR/bootstrap.tar.xz https://dcosio.azureedge.net/dcos/testing/bootstrap/${BOOTSTRAP_ID}.bootstrap.tar.xz &
curl -fLsSv --retry 20 -Y 100000 -y 60 -o $MESOSDIR/bootstrap.tar.xz https://az837203.vo.msecnd.net/dcos/testing/bootstrap/${BOOTSTRAP_ID}.bootstrap.tar.xz &
curl -fLsSv --retry 20 -Y 100000 -y 60 -o $MESOSDIR/d.deb https://az837203.vo.msecnd.net/dcos-deps/docker-engine_1.11.2-0~xenial_amd64.deb &
curl -fLsSv --retry 20 -Y 100000 -y 60 -o $MESOSDIR/1.deb https://az837203.vo.msecnd.net/dcos-deps/libipset3_6.29-1_amd64.deb &
curl -fLsSv --retry 20 -Y 100000 -y 60 -o $MESOSDIR/2.deb https://az837203.vo.msecnd.net/dcos-deps/ipset_6.29-1_amd64.deb &
curl -fLsSv --retry 20 -Y 100000 -y 60 -o $MESOSDIR/3.deb https://az837203.vo.msecnd.net/dcos-deps/unzip_6.0-20ubuntu1_amd64.deb &
curl -fLsSv --retry 20 -Y 100000 -y 60 -o $MESOSDIR/4.deb https://az837203.vo.msecnd.net/dcos-deps/libltdl7_2.4.6-0.1_amd64.deb &
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