#!/bin/bash

if [ ! -e /opt/azure/containers/runcmd.complete ]; then
  exit 0
fi

READY=$(kubectl --kubeconfig /var/lib/kubelet/kubeconfig get nodes/$(hostname) -o jsonpath='{.status.conditions[?(@.type=="Ready")].status}')
if [ "$ready" != "Unknown" ]; then
  exit 0
fi

# Setting "cooling" interval between subsequent restarts to 30 minutes (1800 sec.)
COOLING=1800
FNAME=/var/run/kubelet-last-restart
if [ ! -e $FNAME ]; then
  touch $FNAME
  SINCE_LAST=${COOLING}
else
  SINCE_LAST=$(($(date +%s) - $(stat -c %X $FNAME)))
fi

if [ ${SINCE_LAST} -ge ${COOLING} ]; then
  systemctl restart kubelet
  touch $FNAME
fi
