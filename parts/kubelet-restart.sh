#!/bin/bash

if [ ! -e /opt/azure/containers/runcmd.complete ]; then
  exit 0
fi

ready=$(kubectl --kubeconfig /var/lib/kubelet/kubeconfig get nodes/$(hostname) -o jsonpath='{.status.conditions[?(@.type=="Ready")].status}')
if [ "$ready" = "Unknown" ]; then
  systemctl restart kubelet
fi
