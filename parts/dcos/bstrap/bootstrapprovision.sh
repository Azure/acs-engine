#!/bin/bash

source /opt/azure/dcos/environment

retrycmd_if_failure() { retries=$1; wait=$2; timeout=$3; shift && shift && shift; for i in $(seq 1 $retries); do timeout $timeout ${@}; [ $? -eq 0  ] && break || sleep $wait; done; echo Executed \"$@\" $i times; }

TMPDIR="/tmp/dcos"
mkdir -p $TMPDIR

curl -fLsSv --retry 20 -Y 100000 -y 60 -o $TMPDIR/key https://download.docker.com/linux/ubuntu/gpg &
wait

apt-key add $TMPDIR/key
apt-key fingerprint 0EBFCD88
add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable"
retrycmd_if_failure 10 10 120 apt-get update
retrycmd_if_failure 10 10 120 apt-get install docker-ce=17.06.2~ce-0~ubuntu -y
