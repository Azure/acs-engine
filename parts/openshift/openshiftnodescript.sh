#!/bin/bash -ex

# TODO: /etc/dnsmasq.d/origin-upstream-dns.conf is currently hardcoded; it
# probably shouldn't be
SERVICE_TYPE=origin
if [ -f "/etc/sysconfig/atomic-openshift-node" ]; then
    SERVICE_TYPE=atomic-openshift
fi

{{if eq .Role "infra"}}
echo "BOOTSTRAP_CONFIG_NAME=node-config-infra" >>/etc/sysconfig/${SERVICE_TYPE}-node
{{else}}
echo "BOOTSTRAP_CONFIG_NAME=node-config-compute" >>/etc/sysconfig/${SERVICE_TYPE}-node
{{end}}

rm -rf /etc/etcd/* /etc/origin/master/* /etc/origin/node/*

( cd / && base64 -d <<< {{ .ConfigBundle }} | tar -xz)

cp /etc/origin/node/ca.crt /etc/pki/ca-trust/source/anchors/openshift-ca.crt
update-ca-trust

# note: ${SERVICE_TYPE}-node crash loops until master is up
systemctl enable ${SERVICE_TYPE}-node.service
systemctl start ${SERVICE_TYPE}-node.service &
