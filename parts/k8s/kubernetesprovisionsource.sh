#!/bin/sh

retrycmd_if_failure() { retries=$1; wait=$2; shift && shift; for i in $(seq 1 $retries); do ${@}; [ $? -eq 0  ] && break || sleep $wait; done; echo Executed \"$@\" $i times; }
retrycmd_if_failure_no_stats() { retries=$1; wait=$2; shift && shift; for i in $(seq 1 $retries); do ${@}; [ $? -eq 0  ] && break || sleep $wait; done; }