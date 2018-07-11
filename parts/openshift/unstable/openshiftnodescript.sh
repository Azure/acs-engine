#!/bin/bash -ex

# TODO: /etc/dnsmasq.d/origin-upstream-dns.conf is currently hardcoded; it
# probably shouldn't be
SERVICE_TYPE=origin
if [ -f "/etc/sysconfig/atomic-openshift-node" ]; then
    SERVICE_TYPE=atomic-openshift
fi

if grep -q ^ResourceDisk.Filesystem=xfs /etc/waagent.conf; then
	# Bad image: docker and waagent are racing.  Try to fix up.  Leave this code
	# until the bad images have gone away.
	set +e

	# stop docker if it hasn't failed already
	systemctl stop docker.service

	# wait until waagent has run mkfs and mounted /var/lib/docker
	while ! mountpoint -q /var/lib/docker; do
		sleep 1
	done

	# now roll us back. /var/lib/docker/* may be mounted if docker lost the
	# race.
	umount /var/lib/docker
	umount /var/lib/docker/*

	# disable waagent from racing again if we reboot.
	sed -i -e '/^ResourceDisk.Format=/ s/=.*/=n/' /etc/waagent.conf
	set -e
fi

systemctl stop docker.service
# Also a bad image: the umount should also go away.
umount /var/lib/docker || true
mkfs.xfs -f /dev/sdb1
echo '/dev/sdb1  /var/lib/docker  xfs  grpquota  0 0' >>/etc/fstab
mount /var/lib/docker
restorecon -R /var/lib/docker
systemctl start docker.service

{{if eq .Role "infra"}}
echo "BOOTSTRAP_CONFIG_NAME=node-config-infra" >>/etc/sysconfig/${SERVICE_TYPE}-node
{{else}}
echo "BOOTSTRAP_CONFIG_NAME=node-config-compute" >>/etc/sysconfig/${SERVICE_TYPE}-node
{{end}}

sed -i -e "s#DEBUG_LOGLEVEL=2#DEBUG_LOGLEVEL=4#" /etc/sysconfig/${SERVICE_TYPE}-node

rm -rf /etc/etcd/* /etc/origin/master/*

( cd / && base64 -d <<< {{ .ConfigBundle | shellQuote }} | tar -xz)

cp /etc/origin/node/ca.crt /etc/pki/ca-trust/source/anchors/openshift-ca.crt
update-ca-trust

# note: ${SERVICE_TYPE}-node crash loops until master is up
systemctl enable ${SERVICE_TYPE}-node.service
systemctl start ${SERVICE_TYPE}-node.service &

while [[ $(KUBECONFIG=/etc/origin/node/node.kubeconfig oc get node $(hostname) -o template \
    --template '{{`{{range .status.conditions}}{{if eq .type "Ready"}}{{.status}}{{end}}{{end}}`}}') != True ]]; do
    sleep 1
done
