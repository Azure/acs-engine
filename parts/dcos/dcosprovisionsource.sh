#!/bin/sh

retrycmd_if_failure() {
    retries=$1; wait_sleep=$2; timeout=$3; shift && shift && shift
    for i in $(seq 1 $retries); do
        timeout $timeout ${@}
        [ $? -eq 0  ] && break || \
        if [ $i -eq $retries ]; then
            echo "Error: Failed to execute \"$@\" after $i attempts"
            return 1
        else
            sleep $wait_sleep
        fi
    done
    echo Executed \"$@\" $i times;
}
retry_get_install_deb() {
  retries=$1; wait_sleep=$2; timeout=$3; url=$4;
  deb=$(mktemp)
  trap "rm -f $deb" RETURN
  retrycmd_if_failure $retries $wait_sleep $timeout curl -fsSL $url -o $deb
  if [ $? -ne 0  ]; then
    echo "Error: Failed to download $url"
    return 1
  fi
  retrycmd_if_failure $retries $wait_sleep $timeout dpkg -i $deb
  if [ $? -ne 0  ]; then
    echo "Error: Failed to install $url"
    return 1
  fi
}
