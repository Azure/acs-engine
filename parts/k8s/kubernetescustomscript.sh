#!/bin/bash
set -x
echo `date`,`hostname`, startcustomscript>>/opt/m
source /opt/azure/containers/provision_source.sh
source /opt/azure/containers/provision_installs.sh
source /opt/azure/containers/provision_configs.sh

CUSTOM_SEARCH_DOMAIN_SCRIPT=/opt/azure/containers/setup-custom-search-domains.sh

set +x
ETCD_PEER_CERT=$(echo ${ETCD_PEER_CERTIFICATES} | cut -d'[' -f 2 | cut -d']' -f 1 | cut -d',' -f $((${NODE_INDEX}+1)))
ETCD_PEER_KEY=$(echo ${ETCD_PEER_PRIVATE_KEYS} | cut -d'[' -f 2 | cut -d']' -f 1 | cut -d',' -f $((${NODE_INDEX}+1)))
set -x

if [[ $OS == $COREOS_OS_NAME ]]; then
    echo "Changing default kubectl bin location"
    KUBECTL=/opt/kubectl
fi

if [ -f /var/run/reboot-required ]; then
    REBOOTREQUIRED=true
else
    REBOOTREQUIRED=false
fi

if [ -f /var/log/azure/golden-image-install.complete ]; then
    echo "detected golden image pre-install"
    FULL_INSTALL_REQUIRED=false
else
    FULL_INSTALL_REQUIRED=true
fi

function testOutboundConnection() {
    retrycmd_if_failure 40 1 3 nc -vz www.google.com 443 || retrycmd_if_failure 40 1 3 nc -vz www.1688.com 443 || exit $ERR_OUTBOUND_CONN_FAIL
}

function holdWALinuxAgent() {
    if [[ $OS == $UBUNTU_OS_NAME ]]; then
        # make sure walinuxagent doesn't get updated in the middle of running this script
        retrycmd_if_failure 20 5 30 apt-mark hold walinuxagent || exit $ERR_HOLD_WALINUXAGENT
    fi
}

testOutboundConnection

if [[ ! -z "${MASTER_NODE}" ]]; then
    installEtcd
fi

if $FULL_INSTALL_REQUIRED; then
    holdWALinuxAgent
    installDeps
else 
    echo "Golden image; skipping dependencies installation"
fi

installContainerRuntime
installNetworkPlugin
installContainerd
extractHyperkube
ensureRPC
createKubeManifestDir

if [[ ! -z "${MASTER_NODE}" ]]; then
    configureEtcd
else
    removeEtcd
fi

if [ -f $CUSTOM_SEARCH_DOMAIN_SCRIPT ]; then
    $CUSTOM_SEARCH_DOMAIN_SCRIPT > /opt/azure/containers/setup-custom-search-domain.log 2>&1 || exit $ERR_CUSTOM_SEARCH_DOMAINS_FAIL
fi

if [[ "$CONTAINER_RUNTIME" == "docker" ]]; then
    ensureDocker
elif [[ "$CONTAINER_RUNTIME" == "clear-containers" ]]; then
	# Ensure we can nest virtualization
	if grep -q vmx /proc/cpuinfo; then
        ensureCCProxy
	fi
elif [[ "$CONTAINER_RUNTIME" == "kata-containers" ]]; then
    # Ensure we can nest virtualization
    if grep -q vmx /proc/cpuinfo; then
        installKataContainersRuntime
    fi
fi

configureK8s
configureCNI

if [[ ! -z "${MASTER_NODE}" ]]; then
    configAddons
fi

ensureContainerd

if [[ ! -z "${MASTER_NODE}" && "${KMS_PROVIDER_VAULT_NAME}" != "" ]]; then
    ensureKMS
fi

ensureKubelet
ensureJournal


if [[ ! -z "${MASTER_NODE}" ]]; then
    writeKubeConfig
    ensureEtcd
    ensureK8sControlPlane
    # workaround for 1.12 bug https://github.com/Azure/acs-engine/issues/3681 will remove once upstream is fixed
    if [[ "${KUBERNETES_VERSION}" = 1.12.* ]]; then
        ensureKubelet 
    fi
fi

if [[ "${GPU_NODE}" = true ]]; then
    if $FULL_INSTALL_REQUIRED; then
        installGPUDrivers
    fi
    ensureGPUDrivers
fi

if $FULL_INSTALL_REQUIRED; then
    if [[ $OS == $UBUNTU_OS_NAME ]]; then
        # mitigation for bug https://bugs.launchpad.net/ubuntu/+source/linux/+bug/1676635
        echo 2dd1ce17-079e-403c-b352-a1921ee207ee > /sys/bus/vmbus/drivers/hv_util/unbind
        sed -i "13i\echo 2dd1ce17-079e-403c-b352-a1921ee207ee > /sys/bus/vmbus/drivers/hv_util/unbind\n" /etc/rc.local

        retrycmd_if_failure 20 5 30 apt-mark unhold walinuxagent || exit $ERR_RELEASE_HOLD_WALINUXAGENT
    fi
fi

echo "Custom script finished successfully"

echo `date`,`hostname`, endcustomscript>>/opt/m
mkdir -p /opt/azure/containers && touch /opt/azure/containers/provision.complete
ps auxfww > /opt/azure/provision-ps.log &

if $REBOOTREQUIRED; then
  # wait 1 minute to restart node, so that the custom script extension can complete
  echo 'reboot required, rebooting node in 1 minute'
  /bin/bash -c "shutdown -r 1 &"
else
  runAptDaily &
fi