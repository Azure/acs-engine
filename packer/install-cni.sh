#!/bin/bash -eux

VNET_CNI_PLUGINS_URL="https://acs-mirror.azureedge.net/cni/azure-vnet-cni-linux-amd64-latest.tgz"
CNI_PLUGINS_URL="https://acs-mirror.azureedge.net/cni/cni-plugins-amd64-latest.tgz"

CNI_CONFIG_DIR="/etc/cni/net.d"
CNI_BIN_DIR="/opt/cni/bin"
AZURE_CNI_TGZ_TMP="/tmp/azure_cni.tgz"
CONTAINERNETWORKING_CNI_TGZ_TMP="/tmp/containernetworking_cni.tgz"

# Create CNI conf and bin directories
mkdir -p \
  $CNI_CONFIG_DIR \
  $CNI_BIN_DIR

# install Azure CNI
curl -fsSL $VNET_CNI_PLUGINS_URL -o $AZURE_CNI_TGZ_TMP
tar -xzf $AZURE_CNI_TGZ_TMP -C $CNI_BIN_DIR

cat << EOF > $CNI_CONFIG_DIR/10-azure.conflist
{
   "cniVersion":"0.3.0",
   "name":"azure",
   "plugins":[
      {
         "type":"azure-vnet",
         "mode":"bridge",
         "bridge":"azure0",
         "ipam":{
            "type":"azure-vnet-ipam"
         }
      },
      {
         "type":"portmap",
         "capabilities":{
            "portMappings":true
         },
         "snat":true
      }
   ]
}
EOF

# install CNI plugins
curl -fsSL ${CNI_PLUGINS_URL} -o ${CONTAINERNETWORKING_CNI_TGZ_TMP}
tar -xzf ${CONTAINERNETWORKING_CNI_TGZ_TMP} -C $CNI_BIN_DIR
chown -R root:root $CNI_BIN_DIR
chmod -R 755 $CNI_BIN_DIR

chmod 600 $CNI_CONFIG_DIR/10-azure.conflist
/sbin/ebtables -t nat --list
