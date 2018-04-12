#!/bin/sh

retrycmd_if_failure() { retries=$1; wait=$2; timeout=$3; shift && shift && shift; for i in $(seq 1 $retries); do timeout $timeout ${@}; [ $? -eq 0  ] && break || sleep $wait; done; echo Executed \"$@\" $i times; }
retrycmd_if_failure_no_stats() { retries=$1; wait=$2; timeout=$3; shift && shift && shift; for i in $(seq 1 $retries); do timeout $timeout ${@}; [ $? -eq 0  ] && break || sleep $wait; done; }
retrycmd_get_tarball() { retries=$1; wait=$2; tarball=$3; url=$4; for i in $(seq 1 $retries); do tar -tzf $tarball; [ $? -eq 0  ] && break || retrycmd_if_failure_no_stats $retries 1 10 curl -fsSL $url -o $tarball; sleep $wait; done; }
wait_for_file() { retries=$1; wait=$2; filepath=$3; for i in $(seq 1 $retries); do if [ -f $filepath ]; then break; fi; sleep $wait; done }
apt_get_update() { for i in $(seq 1 100); do apt-get update 2>&1 | grep -x "[WE]:.*"; [ $? -ne 0  ] && break || sleep 1; done; echo Executed apt-get update $i times; }
systemctl_restart() { retries=$1; wait=$2; timeout=$3 svcname=$4; for i in $(seq 1 $retries); do timeout $timeout systemctl daemon-reload && systemctl restart $svcname && systemctl daemon-reload; [ $? -eq 0  ] && break || sleep $wait; done; }