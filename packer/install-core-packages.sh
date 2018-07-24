#!/bin/bash -eux
export DEBIAN_FRONTEND=noninteractive
apt_flags=(-o "Dpkg::Options::=--force-confnew" -qy)

apt-get update -q
apt-get upgrade "${apt_flags[@]}"

curl -fsSL https://packages.microsoft.com/config/ubuntu/16.04/packages-microsoft-prod.deb > /tmp/packages-microsoft-prod.deb
dpkg -i /tmp/packages-microsoft-prod.deb

# See https://github.com/kubernetes/kubernetes/blob/master/build/debian-hyperkube-base/Dockerfile#L25-L44
apt-get install "${apt_flags[@]}" apt-transport-https ca-certificates iptables iproute2 ebtables socat util-linux mount ethtool init-system-helpers nfs-common ceph-common conntrack glusterfs-client ipset jq cgroup-lite git pigz xz-utils fuse
