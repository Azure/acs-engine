#!/bin/bash -eux

HYPERKUBE_VERSION="v1.10.5"
HYPERKUBE_URL="k8s.gcr.io/hyperkube-amd64:${HYPERKUBE_VERSION}"

TMP_DIR=$(mktemp -d)
curl -sSL -o /usr/local/bin/img "https://acs-mirror.azureedge.net/img/img-linux-amd64-v0.4.6"
chmod +x /usr/local/bin/img
img pull $HYPERKUBE_URL
img unpack $HYPERKUBE_URL

path=$(find /home/packer/rootfs -name "hyperkube")
cp "$path" "/usr/local/bin/kubelet"
cp "$path" "/usr/local/bin/kubectl"

chmod a+x /usr/local/bin/kubelet /usr/local/bin/kubectl
rm -rf /home/packer/rootfs

echo "Install complete successfully" > /var/log/azure/golden-image-install.complete
