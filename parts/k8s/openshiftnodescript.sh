#!/bin/bash -x

# TODO: /etc/dnsmasq.d/origin-upstream-dns.conf is currently hardcoded; it
# probably shouldn't be
SERVICE_TYPE=origin
if [ -f "/etc/sysconfig/atomic-openshift-node" ]; then
    SERVICE_TYPE=atomic-openshift
fi

# TODO: remove this once we generate the registry certificate
cat >>/etc/sysconfig/docker <<'EOF'
INSECURE_REGISTRY='--insecure-registry 172.30.0.0/16'
EOF

systemctl restart docker.service

{{if .IsInfra}}
echo "BOOTSTRAP_CONFIG_NAME=node-config-infra" >>/etc/sysconfig/${SERVICE_TYPE}-node
{{else}}
echo "BOOTSTRAP_CONFIG_NAME=node-config-compute" >>/etc/sysconfig/${SERVICE_TYPE}-node
{{end}}

rm -rf /etc/etcd/* /etc/origin/master/* /etc/origin/node/*

( cd / && base64 -d <<< {{ .ConfigBundle }} | tar -xz)

cp /etc/origin/node/ca.crt /etc/pki/ca-trust/source/anchors/openshift-ca.crt
update-ca-trust

# TODO: when enabling secure registry, may need:
# ln -s /etc/origin/node/node-client-ca.crt /etc/docker/certs.d/docker-registry.default.svc:5000

# note: ${SERVICE_TYPE}-node crash loops until master is up
systemctl enable ${SERVICE_TYPE}-node.service
systemctl start ${SERVICE_TYPE}-node.service

exit 0
