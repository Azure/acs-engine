#!/bin/bash -eux

CRI_CONTAINERD_VERSION="1.1.0"

wget -q --show-progress --https-only --timestamping \
"https://storage.googleapis.com/cri-containerd-release/cri-containerd-${CRI_CONTAINERD_VERSION}.linux-amd64.tar.gz"
sudo tar -xvf cri-containerd-${CRI_CONTAINERD_VERSION}.linux-amd64.tar.gz -C /

sed -i '/\[Service\]/a ExecStartPost=\/sbin\/iptables -P FORWARD ACCEPT' /etc/systemd/system/containerd.service
