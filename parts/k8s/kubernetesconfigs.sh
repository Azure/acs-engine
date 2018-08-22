#!/bin/bash

function systemctlEnableAndStart() {
    systemctl_restart 100 5 30 $1
    RESTART_STATUS=$?
    systemctl status $1 --no-pager -l > /var/log/azure/$1-status.log
    if [ $RESTART_STATUS -ne 0 ]; then
        echo "$1 could not be started"
        exit $ERR_SYSTEMCTL_START_FAIL
    fi
    retrycmd_if_failure 10 5 3 systemctl enable $1
    if [ $? -ne 0 ]; then
        echo "$1 could not be enabled by systemctl"
        exit $ERR_SYSTEMCTL_ENABLE_FAIL
    fi
}

function configureEtcd() {
    useradd -U "etcd"
    usermod -p "$(head -c 32 /dev/urandom | base64)" "etcd"
    passwd -u "etcd"
    id "etcd"

    APISERVER_PRIVATE_KEY_PATH="/etc/kubernetes/certs/apiserver.key"
    touch "${APISERVER_PRIVATE_KEY_PATH}"
    chmod 0600 "${APISERVER_PRIVATE_KEY_PATH}"
    chown root:root "${APISERVER_PRIVATE_KEY_PATH}"

    CA_PRIVATE_KEY_PATH="/etc/kubernetes/certs/ca.key"
    touch "${CA_PRIVATE_KEY_PATH}"
    chmod 0600 "${CA_PRIVATE_KEY_PATH}"
    chown root:root "${CA_PRIVATE_KEY_PATH}"

    ETCD_SERVER_PRIVATE_KEY_PATH="/etc/kubernetes/certs/etcdserver.key"
    touch "${ETCD_SERVER_PRIVATE_KEY_PATH}"
    chmod 0600 "${ETCD_SERVER_PRIVATE_KEY_PATH}"
    chown etcd:etcd "${ETCD_SERVER_PRIVATE_KEY_PATH}"

    ETCD_CLIENT_PRIVATE_KEY_PATH="/etc/kubernetes/certs/etcdclient.key"
    touch "${ETCD_CLIENT_PRIVATE_KEY_PATH}"
    chmod 0600 "${ETCD_CLIENT_PRIVATE_KEY_PATH}"
    chown root:root "${ETCD_CLIENT_PRIVATE_KEY_PATH}"

    ETCD_PEER_PRIVATE_KEY_PATH="/etc/kubernetes/certs/etcdpeer${MASTER_INDEX}.key"
    touch "${ETCD_PEER_PRIVATE_KEY_PATH}"
    chmod 0600 "${ETCD_PEER_PRIVATE_KEY_PATH}"
    chown etcd:etcd "${ETCD_PEER_PRIVATE_KEY_PATH}"

    ETCD_SERVER_CERTIFICATE_PATH="/etc/kubernetes/certs/etcdserver.crt"
    touch "${ETCD_SERVER_CERTIFICATE_PATH}"
    chmod 0644 "${ETCD_SERVER_CERTIFICATE_PATH}"
    chown root:root "${ETCD_SERVER_CERTIFICATE_PATH}"

    ETCD_CLIENT_CERTIFICATE_PATH="/etc/kubernetes/certs/etcdclient.crt"
    touch "${ETCD_CLIENT_CERTIFICATE_PATH}"
    chmod 0644 "${ETCD_CLIENT_CERTIFICATE_PATH}"
    chown root:root "${ETCD_CLIENT_CERTIFICATE_PATH}"

    ETCD_PEER_CERTIFICATE_PATH="/etc/kubernetes/certs/etcdpeer${MASTER_INDEX}.crt"
    touch "${ETCD_PEER_CERTIFICATE_PATH}"
    chmod 0644 "${ETCD_PEER_CERTIFICATE_PATH}"
    chown root:root "${ETCD_PEER_CERTIFICATE_PATH}"

    set +x
    echo "${APISERVER_PRIVATE_KEY}" | base64 --decode > "${APISERVER_PRIVATE_KEY_PATH}"
    echo "${CA_PRIVATE_KEY}" | base64 --decode > "${CA_PRIVATE_KEY_PATH}"
    echo "${ETCD_SERVER_PRIVATE_KEY}" | base64 --decode > "${ETCD_SERVER_PRIVATE_KEY_PATH}"
    echo "${ETCD_CLIENT_PRIVATE_KEY}" | base64 --decode > "${ETCD_CLIENT_PRIVATE_KEY_PATH}"
    echo "${ETCD_PEER_KEY}" | base64 --decode > "${ETCD_PEER_PRIVATE_KEY_PATH}"
    echo "${ETCD_SERVER_CERTIFICATE}" | base64 --decode > "${ETCD_SERVER_CERTIFICATE_PATH}"
    echo "${ETCD_CLIENT_CERTIFICATE}" | base64 --decode > "${ETCD_CLIENT_CERTIFICATE_PATH}"
    echo "${ETCD_PEER_CERT}" | base64 --decode > "${ETCD_PEER_CERTIFICATE_PATH}"
    set -x

    /opt/azure/containers/setup-etcd.sh > /opt/azure/containers/setup-etcd.log 2>&1
    RET=$?
    if [ $RET -ne 0 ]; then
        exit $RET
    fi

    /opt/azure/containers/mountetcd.sh || exit $ERR_ETCD_VOL_MOUNT_FAIL
    systemctlEnableAndStart etcd
    for i in $(seq 1 600); do
        MEMBER="$(sudo etcdctl member list | grep -E ${MASTER_VM_NAME} | cut -d':' -f 1)"
        if [ "$MEMBER" != "" ]; then
            break
        else
            sleep 1
        fi
    done
    retrycmd_if_failure 10 1 5 sudo etcdctl member update $MEMBER ${ETCD_PEER_URL} || exit $ERR_ETCD_CONFIG_FAIL
}

function ensureRPC() {
    systemctlEnableAndStart rpcbind
    systemctlEnableAndStart rpc-statd
}

function runAptDaily() {
    /usr/lib/apt/apt.systemd.daily
}

function generateAggregatedAPICerts() {
    wait_for_file 1 1 /etc/kubernetes/generate-proxy-certs.sh && /etc/kubernetes/generate-proxy-certs.sh
}

function configureK8s() {
    KUBELET_PRIVATE_KEY_PATH="/etc/kubernetes/certs/client.key"
    touch "${KUBELET_PRIVATE_KEY_PATH}"
    chmod 0600 "${KUBELET_PRIVATE_KEY_PATH}"
    chown root:root "${KUBELET_PRIVATE_KEY_PATH}"

    APISERVER_PUBLIC_KEY_PATH="/etc/kubernetes/certs/apiserver.crt"
    touch "${APISERVER_PUBLIC_KEY_PATH}"
    chmod 0644 "${APISERVER_PUBLIC_KEY_PATH}"
    chown root:root "${APISERVER_PUBLIC_KEY_PATH}"

    AZURE_JSON_PATH="/etc/kubernetes/azure.json"
    touch "${AZURE_JSON_PATH}"
    chmod 0600 "${AZURE_JSON_PATH}"
    chown root:root "${AZURE_JSON_PATH}"

    set +x
    echo "${KUBELET_PRIVATE_KEY}" | base64 --decode > "${KUBELET_PRIVATE_KEY_PATH}"
    echo "${APISERVER_PUBLIC_KEY}" | base64 --decode > "${APISERVER_PUBLIC_KEY_PATH}"
    cat << EOF > "${AZURE_JSON_PATH}"
{
    "cloud":"${TARGET_ENVIRONMENT}",
    "tenantId": "${TENANT_ID}",
    "subscriptionId": "${SUBSCRIPTION_ID}",
    "aadClientId": "${SERVICE_PRINCIPAL_CLIENT_ID}",
    "aadClientSecret": "${SERVICE_PRINCIPAL_CLIENT_SECRET}",
    "resourceGroup": "${RESOURCE_GROUP}",
    "location": "${LOCATION}",
    "vmType": "${VM_TYPE}",
    "subnetName": "${SUBNET}",
    "securityGroupName": "${NETWORK_SECURITY_GROUP}",
    "vnetName": "${VIRTUAL_NETWORK}",
    "vnetResourceGroup": "${VIRTUAL_NETWORK_RESOURCE_GROUP}",
    "routeTableName": "${ROUTE_TABLE}",
    "primaryAvailabilitySetName": "${PRIMARY_AVAILABILITY_SET}",
    "primaryScaleSetName": "${PRIMARY_SCALE_SET}",
    "cloudProviderBackoff": ${CLOUDPROVIDER_BACKOFF},
    "cloudProviderBackoffRetries": ${CLOUDPROVIDER_BACKOFF_RETRIES},
    "cloudProviderBackoffExponent": ${CLOUDPROVIDER_BACKOFF_EXPONENT},
    "cloudProviderBackoffDuration": ${CLOUDPROVIDER_BACKOFF_DURATION},
    "cloudProviderBackoffJitter": ${CLOUDPROVIDER_BACKOFF_JITTER},
    "cloudProviderRatelimit": ${CLOUDPROVIDER_RATELIMIT},
    "cloudProviderRateLimitQPS": ${CLOUDPROVIDER_RATELIMIT_QPS},
    "cloudProviderRateLimitBucket": ${CLOUDPROVIDER_RATELIMIT_BUCKET},
    "useManagedIdentityExtension": ${USE_MANAGED_IDENTITY_EXTENSION},
    "useInstanceMetadata": ${USE_INSTANCE_METADATA},
    "loadBalancerSku": "${LOAD_BALANCER_SKU}",
    "excludeMasterFromStandardLB": ${EXCLUDE_MASTER_FROM_STANDARD_LB},
    "providerVaultName": "${KMS_PROVIDER_VAULT_NAME}",
    "providerKeyName": "k8s",
    "providerKeyVersion": ""
}
EOF
    set -x
    generateAggregatedAPICerts
}

function configureCNI() {
    # Turn on br_netfilter (needed for the iptables rules to work on bridges)
    # and permanently enable it
    retrycmd_if_failure 30 6 10 modprobe br_netfilter || exit $ERR_MODPROBE_FAIL
    # /etc/modules-load.d is the location used by systemd to load modules
    echo -n "br_netfilter" > /etc/modules-load.d/br_netfilter.conf
    if [[ "${NETWORK_PLUGIN}" = "azure" ]]; then
        mv $CNI_BIN_DIR/10-azure.conflist $CNI_CONFIG_DIR/
        chmod 600 $CNI_CONFIG_DIR/10-azure.conflist
        /sbin/ebtables -t nat --list
    fi
}

function setKubeletOpts () {
    sed -i "s#^KUBELET_OPTS=.*#KUBELET_OPTS=${1}#" /etc/default/kubelet
}

function ensureCCProxy() {
    cat $CC_SERVICE_IN_TMP | sed 's#@libexecdir@#/usr/libexec#' > /etc/systemd/system/cc-proxy.service
    cat $CC_SOCKET_IN_TMP sed 's#@localstatedir@#/var#' > /etc/systemd/system/cc-proxy.socket
    # Enable and start Clear Containers proxy service
	echo "Enabling and starting Clear Containers proxy service..."
	systemctlEnableAndStart cc-proxy
}

function setupContainerd() {
    echo "Configuring cri-containerd..."
    mkdir -p "/etc/containerd"
    CRI_CONTAINERD_CONFIG="/etc/containerd/config.toml"
    echo "subreaper = false" > "$CRI_CONTAINERD_CONFIG"
    echo "oom_score = 0" >> "$CRI_CONTAINERD_CONFIG"
    echo "[plugins.cri]" >> "$CRI_CONTAINERD_CONFIG" 
    echo "sandbox_image = \"$POD_INFRA_CONTAINER_SPEC\"" >> "$CRI_CONTAINERD_CONFIG"
    echo "[plugins.cri.containerd.untrusted_workload_runtime]" >> "$CRI_CONTAINERD_CONFIG"
    echo "runtime_type = 'io.containerd.runtime.v1.linux'" >> "$CRI_CONTAINERD_CONFIG"
    if [[ "$CONTAINER_RUNTIME" == "clear-containers" ]]; then
        echo "runtime_engine = '/usr/bin/cc-runtime'" >> "$CRI_CONTAINERD_CONFIG"
    elif [[ "$CONTAINER_RUNTIME" == "kata-containers" ]]; then
        echo "runtime_engine = '/usr/bin/kata-runtime'" >> "$CRI_CONTAINERD_CONFIG"
    else
        echo "runtime_engine = '/usr/local/sbin/runc'" >> "$CRI_CONTAINERD_CONFIG"
    fi
    echo "[plugins.cri.containerd.default_runtime]" >> "$CRI_CONTAINERD_CONFIG"
    echo "runtime_type = 'io.containerd.runtime.v1.linux'" >> "$CRI_CONTAINERD_CONFIG"
    echo "runtime_engine = '/usr/local/sbin/runc'" >> "$CRI_CONTAINERD_CONFIG"
    setKubeletOpts " --container-runtime=remote --runtime-request-timeout=15m --container-runtime-endpoint=unix:///run/containerd/containerd.sock"
}

function ensureContainerd() {
    if [[ "$CONTAINER_RUNTIME" == "clear-containers" ]] || [[ "$CONTAINER_RUNTIME" == "kata-containers" ]] || [[ "$CONTAINER_RUNTIME" == "containerd" ]]; then
        setupContainerd
        # Enable and start cri-containerd service
        # Make sure this is done after networking plugins are installed
        echo "Enabling and starting cri-containerd service..."
        systemctlEnableAndStart containerd
    fi
}

function ensureDocker() {
    echo "ExecStartPost=/sbin/iptables -P FORWARD ACCEPT" >> /etc/systemd/system/docker.service.d/exec_start.conf
    usermod -aG docker ${ADMINUSER}
    wait_for_file 600 1 $DOCKER || exit $ERR_FILE_WATCH_TIMEOUT
    systemctlEnableAndStart docker
    retrycmd_if_failure 6 1 10 docker pull busybox # pre-pull busybox, but don't exit if fail 
    systemctlEnableAndStart docker-health-probe
}
function ensureKMS() {
    systemctlEnableAndStart kms
}

function ensureKubelet() {
    systemctlEnableAndStart kubelet
}

function ensureJournal(){
    echo "Storage=persistent" >> /etc/systemd/journald.conf
    echo "SystemMaxUse=1G" >> /etc/systemd/journald.conf
    echo "RuntimeMaxUse=1G" >> /etc/systemd/journald.conf
    echo "ForwardToSyslog=no" >> /etc/systemd/journald.conf
    systemctlEnableAndStart systemd-journald
}

function ensurePodSecurityPolicy() {
    POD_SECURITY_POLICY_FILE="/etc/kubernetes/manifests/pod-security-policy.yaml"
    if [ -f $POD_SECURITY_POLICY_FILE ]; then
        $KUBECTL create -f $POD_SECURITY_POLICY_FILE
    fi
}

function ensureK8sControlPlane() {
    if $REBOOTREQUIRED; then
        return
    fi
    wait_for_file 600 1 $KUBECTL || exit $ERR_FILE_WATCH_TIMEOUT
    retrycmd_if_failure 900 1 20 $KUBECTL 2>/dev/null cluster-info || exit $ERR_K8S_RUNNING_TIMEOUT
    ensurePodSecurityPolicy
}

function ensureEtcd() {
    retrycmd_if_failure 120 5 10 curl --cacert /etc/kubernetes/certs/ca.crt --cert /etc/kubernetes/certs/etcdclient.crt --key /etc/kubernetes/certs/etcdclient.key ${ETCD_CLIENT_URL}/v2/machines || exit $ERR_ETCD_RUNNING_TIMEOUT
}

function writeKubeConfig() {
    KUBECONFIGDIR=/home/$ADMINUSER/.kube
    KUBECONFIGFILE=$KUBECONFIGDIR/config
    mkdir -p $KUBECONFIGDIR
    touch $KUBECONFIGFILE
    chown $ADMINUSER:$ADMINUSER $KUBECONFIGDIR
    chown $ADMINUSER:$ADMINUSER $KUBECONFIGFILE
    chmod 700 $KUBECONFIGDIR
    chmod 600 $KUBECONFIGFILE

    # disable logging after secret output
    set +x
    echo "
---
apiVersion: v1
clusters:
- cluster:
    certificate-authority-data: \"$CA_CERTIFICATE\"
    server: $KUBECONFIG_SERVER
  name: \"$MASTER_FQDN\"
contexts:
- context:
    cluster: \"$MASTER_FQDN\"
    user: \"$MASTER_FQDN-admin\"
  name: \"$MASTER_FQDN\"
current-context: \"$MASTER_FQDN\"
kind: Config
users:
- name: \"$MASTER_FQDN-admin\"
  user:
    client-certificate-data: \"$KUBECONFIG_CERTIFICATE\"
    client-key-data: \"$KUBECONFIG_KEY\"
" > $KUBECONFIGFILE
    # renable logging after secrets
    set -x
}

function configClusterAutoscalerAddon() {
    if [[ "${USE_MANAGED_IDENTITY_EXTENSION}" == true ]]; then
        CLUSTER_AUTOSCALER_MSI_VOLUME_MOUNT="- mountPath: /var/lib/waagent/\n\          name: waagent\n\          readOnly: true"
        CLUSTER_AUTOSCALER_MSI_VOLUME="- hostPath:\n\          path: /var/lib/waagent/\n\        name: waagent"
        CLUSTER_AUTOSCALER_MSI_HOST_NETWORK="hostNetwork: true"

        sed -i "s|<kubernetesClusterAutoscalerVolumeMounts>|${CLUSTER_AUTOSCALER_MSI_VOLUME_MOUNT}|g" "/etc/kubernetes/addons/cluster-autoscaler-deployment.yaml"
        sed -i "s|<kubernetesClusterAutoscalerVolumes>|${CLUSTER_AUTOSCALER_MSI_VOLUME}|g" "/etc/kubernetes/addons/cluster-autoscaler-deployment.yaml"
        sed -i "s|<kubernetesClusterAutoscalerHostNetwork>|$(echo "${CLUSTER_AUTOSCALER_MSI_HOST_NETWORK}")|g" "/etc/kubernetes/addons/cluster-autoscaler-deployment.yaml"
    elif [[ "${USE_MANAGED_IDENTITY_EXTENSION}" == false ]]; then
        sed -i "s|<kubernetesClusterAutoscalerVolumeMounts>|""|g" "/etc/kubernetes/addons/cluster-autoscaler-deployment.yaml"
        sed -i "s|<kubernetesClusterAutoscalerVolumes>|""|g" "/etc/kubernetes/addons/cluster-autoscaler-deployment.yaml"
        sed -i "s|<kubernetesClusterAutoscalerHostNetwork>|""|g" "/etc/kubernetes/addons/cluster-autoscaler-deployment.yaml"
    fi

    sed -i "s|<kubernetesClusterAutoscalerClientId>|$(echo $SERVICE_PRINCIPAL_CLIENT_ID | base64)|g" "/etc/kubernetes/addons/cluster-autoscaler-deployment.yaml"
    sed -i "s|<kubernetesClusterAutoscalerClientSecret>|$(echo $SERVICE_PRINCIPAL_CLIENT_SECRET | base64)|g" "/etc/kubernetes/addons/cluster-autoscaler-deployment.yaml"
    sed -i "s|<kubernetesClusterAutoscalerSubscriptionId>|$(echo $SUBSCRIPTION_ID | base64)|g" "/etc/kubernetes/addons/cluster-autoscaler-deployment.yaml"
    sed -i "s|<kubernetesClusterAutoscalerTenantId>|$(echo $TENANT_ID | base64)|g" "/etc/kubernetes/addons/cluster-autoscaler-deployment.yaml"
    sed -i "s|<kubernetesClusterAutoscalerResourceGroup>|$(echo $RESOURCE_GROUP | base64)|g" "/etc/kubernetes/addons/cluster-autoscaler-deployment.yaml"
    sed -i "s|<kubernetesClusterAutoscalerVmType>|$(echo $VM_TYPE | base64)|g" "/etc/kubernetes/addons/cluster-autoscaler-deployment.yaml"
    sed -i "s|<kubernetesClusterAutoscalerVMSSName>|$(echo $PRIMARY_SCALE_SET)|g" "/etc/kubernetes/addons/cluster-autoscaler-deployment.yaml"
}

configACIConnectorAddon() {
    ACI_CONNECTOR_CREDENTIALS=$(printf "{\"clientId\": \"$(echo $SERVICE_PRINCIPAL_CLIENT_ID)\", \"clientSecret\": \"$(echo $SERVICE_PRINCIPAL_CLIENT_SECRET)\", \"tenantId\": \"$(echo $TENANT_ID)\", \"subscriptionId\": \"$(echo $SUBSCRIPTION_ID)\", \"activeDirectoryEndpointUrl\": \"https://login.microsoftonline.com\",\"resourceManagerEndpointUrl\": \"https://management.azure.com/\", \"activeDirectoryGraphResourceId\": \"https://graph.windows.net/\", \"sqlManagementEndpointUrl\": \"https://management.core.windows.net:8443/\", \"galleryEndpointUrl\": \"https://gallery.azure.com/\", \"managementEndpointUrl\": \"https://management.core.windows.net/\"}" | base64 -w 0)

    openssl req -newkey rsa:4096 -new -nodes -x509 -days 3650 -keyout /etc/kubernetes/certs/aci-connector-key.pem -out /etc/kubernetes/certs/aci-connector-cert.pem -subj "/C=US/ST=CA/L=virtualkubelet/O=virtualkubelet/OU=virtualkubelet/CN=virtualkubelet"
    ACI_CONNECTOR_KEY=$(base64 /etc/kubernetes/certs/aci-connector-key.pem -w0)
    ACI_CONNECTOR_CERT=$(base64 /etc/kubernetes/certs/aci-connector-cert.pem -w0)

    sed -i "s|<kubernetesACIConnectorCredentials>|$ACI_CONNECTOR_CREDENTIALS|g" "/etc/kubernetes/addons/aci-connector-deployment.yaml"
    sed -i "s|<kubernetesACIConnectorResourceGroup>|$(echo $RESOURCE_GROUP)|g" "/etc/kubernetes/addons/aci-connector-deployment.yaml"
    sed -i "s|<kubernetesACIConnectorCert>|$(echo $ACI_CONNECTOR_CERT)|g" "/etc/kubernetes/addons/aci-connector-deployment.yaml"
    sed -i "s|<kubernetesACIConnectorKey>|$(echo $ACI_CONNECTOR_KEY)|g" "/etc/kubernetes/addons/aci-connector-deployment.yaml"
}

configAddons() {
    if [[ "${CLUSTER_AUTOSCALER_ADDON}" = True ]]; then
        configClusterAutoscalerAddon
    fi

    if [[ "${ACI_CONNECTOR_ADDON}" = True ]]; then
        configACIConnectorAddon
    fi
}