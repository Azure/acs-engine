#!/bin/bash

set -x
OS=$(cat /etc/*-release | grep ^ID= | tr -d 'ID="' | awk '{print toupper($0)}')
UBUNTU_OS_NAME="UBUNTU"
RHEL_OS_NAME="RHEL"
COREOS_OS_NAME="COREOS"
KUBECTL=/usr/local/bin/kubectl
DOCKER=/usr/bin/docker

set +x
ETCD_PEER_CERT=$(echo ${ETCD_PEER_CERTIFICATES} | cut -d'[' -f 2 | cut -d']' -f 1 | cut -d',' -f $((${MASTER_INDEX}+1)))
ETCD_PEER_KEY=$(echo ${ETCD_PEER_PRIVATE_KEYS} | cut -d'[' -f 2 | cut -d']' -f 1 | cut -d',' -f $((${MASTER_INDEX}+1)))
set -x

if [[ $OS == $COREOS_OS_NAME ]]; then
    KUBECTL=/opt/kubectl
fi

ensureRunCommandCompleted()
{
    for i in {1..900}; do
        if [ -e /opt/azure/containers/runcmd.complete ]; then
            break
        fi
        sleep 1
    done
}

ensureDockerInstallCompleted()
{
    for i in {1..900}; do
        if [ -e /opt/azure/containers/dockerinstall.complete ]; then
            break
        fi
        sleep 1
    done
}

echo `date`,`hostname`, startscript>>/opt/m

if [ -f /var/run/reboot-required ]; then
    REBOOTREQUIRED=true
else
    REBOOTREQUIRED=false
fi

if [[ ! -z "${MASTER_NODE}" ]]; then
    useradd -U "etcd"
    usermod -p "$(head -c 32 /dev/urandom | base64)" "etcd"
    passwd -u "etcd"
    id "etcd"
    
    echo `date`,`hostname`, beginGettingEtcdCerts>>/opt/m
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

    echo `date`,`hostname`, endGettingEtcdCerts>>/opt/m
    mkdir -p /opt/azure/containers && touch /opt/azure/containers/certs.ready
fi

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
    "subnetName": "${SUBNET}",
    "securityGroupName": "${NETWORK_SECURITY_GROUP}",
    "vnetName": "${VIRTUAL_NETWORK}",
    "vnetResourceGroup": "${VIRTUAL_NETWORK_RESOURCE_GROUP}",
    "routeTableName": "${ROUTE_TABLE}",
    "primaryAvailabilitySetName": "${PRIMARY_AVAILABILITY_SET}",
    "cloudProviderBackoff": ${CLOUDPROVIDER_BACKOFF},
    "cloudProviderBackoffRetries": ${CLOUDPROVIDER_BACKOFF_RETRIES},
    "cloudProviderBackoffExponent": ${CLOUDPROVIDER_BACKOFF_EXPONENT},
    "cloudProviderBackoffDuration": ${CLOUDPROVIDER_BACKOFF_DURATION},
    "cloudProviderBackoffJitter": ${CLOUDPROVIDER_BACKOFF_JITTER},
    "cloudProviderRatelimit": ${CLOUDPROVIDER_RATELIMIT},
    "cloudProviderRateLimitQPS": ${CLOUDPROVIDER_RATELIMIT_QPS},
    "cloudProviderRateLimitBucket": ${CLOUDPROVIDER_RATELIMIT_BUCKET},
    "useManagedIdentityExtension": ${USE_MANAGED_IDENTITY_EXTENSION},
    "useInstanceMetadata": ${USE_INSTANCE_METADATA}
}
EOF

set -x
function ensureFilepath() {
    if $REBOOTREQUIRED; then
        return
    fi
    found=1
    for i in {1..600}; do
        if [ -e $1 ]
        then
            found=0
            break
        fi
        sleep 1
    done
    if [ $found -ne 0 ]
    then
        exit 1
    fi
}

function retrycmd_if_failure() { retries=$1; wait=$2; shift && shift; for i in $(seq 1 $retries); do ${@}; [ $? -eq 0  ] && break || sleep $wait; done; echo Executed \"$@\" $i times; }

function downloadUrl () {
	for i in 1 2 3 4 5; do curl --max-time 60 -fsSL ${1}; [ $? -eq 0 ] && break || sleep 10; done
    echo Executed curl for \"${1}\" $i times
}

function setMaxPods () {
    sed -i "s/^KUBELET_MAX_PODS=.*/KUBELET_MAX_PODS=${1}/" /etc/default/kubelet
}

function setNetworkPlugin () {
    sed -i "s/^KUBELET_NETWORK_PLUGIN=.*/KUBELET_NETWORK_PLUGIN=${1}/" /etc/default/kubelet
}

function setKubeletOpts () {
	sed -i "s#^KUBELET_OPTS=.*#KUBELET_OPTS=${1}#" /etc/default/kubelet
}

function setDockerOpts () {
    sed -i "s#^DOCKER_OPTS=.*#DOCKER_OPTS=${1}#" /etc/default/kubelet
}

function configAzureNetworkPolicy() {
    CNI_CONFIG_DIR=/etc/cni/net.d
    mkdir -p $CNI_CONFIG_DIR
    chown -R root:root $CNI_CONFIG_DIR
    chmod 755 $CNI_CONFIG_DIR
    CNI_BIN_DIR=/opt/cni/bin
    mkdir -p $CNI_BIN_DIR
    downloadUrl ${VNET_CNI_PLUGINS_URL} | tar -xz -C $CNI_BIN_DIR
    downloadUrl ${CNI_PLUGINS_URL} | tar -xz -C $CNI_BIN_DIR ./loopback ./portmap
    chown -R root:root $CNI_BIN_DIR
    chmod -R 755 $CNI_BIN_DIR

    mv $CNI_BIN_DIR/10-azure.conflist $CNI_CONFIG_DIR/
    chmod 600 $CNI_CONFIG_DIR/10-azure.conflist

    /sbin/ebtables -t nat --list
    setNetworkPlugin cni
    setDockerOpts " --volume=/etc/cni/:/etc/cni:ro --volume=/opt/cni/:/opt/cni:ro"
}

function configCalicoNetworkPolicy() {
    setNetworkPlugin cni
    setDockerOpts " --volume=/etc/cni/:/etc/cni:ro --volume=/opt/cni/:/opt/cni:ro"
}

function configNetworkPolicy() {
    if [[ "${NETWORK_POLICY}" = "azure" ]]; then
        configAzureNetworkPolicy
    elif [[ "${NETWORK_POLICY}" = "calico" ]]; then
        configCalicoNetworkPolicy
    else
        setNetworkPlugin kubenet
        setDockerOpts ""
    fi
}

function installClearContainersRuntime() {
	curl -sSL --retry 5 --retry-delay 10 --retry-max-time 30 "https://download.opensuse.org/repositories/home:clearcontainers:clear-containers-3/xUbuntu_16.04/Release.key" | apt-key add -
	echo 'deb http://download.opensuse.org/repositories/home:/clearcontainers:/clear-containers-3/xUbuntu_16.04/ /' > /etc/apt/sources.list.d/cc-runtime.list

	apt-get update
	apt-get install --no-install-recommends -y \
		cc-runtime

	apt-get install --no-install-recommends -y \
		lvm2 \
		thin-provisioning-tools

	retrycmd_if_failure 5 5 timeout 30s systemctl daemon-reload

	retrycmd_if_failure 5 5 timeout 10s systemctl enable cc-proxy
	retrycmd_if_failure 5 5 timeout 30s systemctl start cc-proxy

	configAzureNetworkPolicy
	setKubeletOpts " --container-runtime=remote --container-runtime-endpoint=/var/run/crio.sock"
	setDockerOpts " --volume=/etc/cni/:/etc/cni:ro --volume=/opt/cni/:/opt/cni:ro"
}

function installGo() {
	export GO_SRC=/usr/local/go
	export GOPATH="${HOME}/.go"

	if [[ -d "$GO_SRC" ]]; then
		rm -rf "$GO_SRC"
	fi

	if [[ -d "$GOPATH" ]]; then
		rm -rf "$GOPATH"
	fi

	GO_VERSION=$(curl --retry 5 --retry-delay 10 --retry-max-time 30 -sSL "https://golang.org/VERSION?m=text")

	(
	curl --retry 5 --retry-delay 10 --retry-max-time 30 -sSL "https://storage.googleapis.com/golang/${GO_VERSION}.linux-amd64.tar.gz" | sudo tar -v -C /usr/local -xz
	)

	export PATH="${GO_SRC}/bin:${PATH}:${GOPATH}/bin"
}

function buildRunc() {
	mkdir -p "${GOPATH}/src/github.com/opencontainers"
	(
	cd "${GOPATH}/src/github.com/opencontainers"
	git clone "https://github.com/opencontainers/runc.git"
	cd runc
	git reset --hard v1.0.0-rc4
	make BUILDTAGS="seccomp apparmor"
	make install
	)
}

function buildCRIO() {
	add-apt-repository -y ppa:projectatomic/ppa
	add-apt-repository -y ppa:alexlarsson/flatpak
	apt-get update

	apt-get install --no-install-recommends -y \
		btrfs-tools \
		gcc \
		git \
		libapparmor-dev \
		libassuan-dev \
		libc6-dev \
		libdevmapper-dev \
		libglib2.0-dev \
		libgpg-error-dev \
		libgpgme11-dev \
		libostree-dev \
		libseccomp-dev \
		libselinux1-dev \
		make \
		pkg-config \
		skopeo-containers

	installGo;

	go get github.com/cpuguy83/go-md2man

	(
	go get -u github.com/docker/docker/daemon/logger/templates
	cd "${GOPATH}/src/github.com/docker/docker"
	mkdir -p utils
	cp -r daemon/logger/templates utils/
	)

	buildRunc;

	mkdir -p "${GOPATH}/src/github.com/kubernetes-incubator"
	(
	cd "${GOPATH}/src/github.com/kubernetes-incubator"
	git clone "https://github.com/kubernetes-incubator/cri-o.git"
	cd cri-o
	git reset --hard v1.0.0
	make BUILDTAGS="seccomp apparmor"
	make install
	make install.config
	make install.systemd
	)

	rm -vrf "$tmpd"

	rm -vrf "$GO_SRC" "$GOPATH"

	setupCRIO;
}

function setupCRIO() {
	SYSTEMD_CRI_O_SERVICE_FILE="/usr/local/lib/systemd/system/crio.service"
	sed -i 's#ExecStart=/usr/local/bin/crio#ExecStart=/usr/local/bin/crio -log-level debug#' "$SYSTEMD_CRI_O_SERVICE_FILE"

	CRI_O_CONFIG="/etc/crio/crio.conf"
	sed -i 's#storage_driver = ""#storage_driver = "devicemapper"#' "$CRI_O_CONFIG"
	sed -i 's#storage_option = \[#storage_option = \["dm.directlvm_device=/dev/sdc", "dm.thinp_percent=95", "dm.thinp_metapercent=1", "dm.thinp_autoextend_threshold=80", "dm.thinp_autoextend_percent=20", "dm.directlvm_device_force=true"#' "$CRI_O_CONFIG"
	sed -i 's#runtime = "/usr/bin/runc"#runtime = "/usr/local/sbin/runc"#' "$CRI_O_CONFIG"
	sed -i 's#runtime_untrusted_workload = ""#runtime_untrusted_workload = "/usr/bin/cc-runtime"#' "$CRI_O_CONFIG"
	sed -i 's#default_workload_trust = "trusted"#default_workload_trust = "untrusted"#' "$CRI_O_CONFIG"

	retrycmd_if_failure 5 5 timeout 30s systemctl daemon-reload
}

function ensureCRIO() {
	if [[ "$CONTAINER_RUNTIME" == "clear-containers" ]]; then
		if grep -q vmx /proc/cpuinfo; then
			retrycmd_if_failure 5 5 timeout 10s systemctl enable crio crio-shutdown
			retrycmd_if_failure 5 5 timeout 30s systemctl start crio
		fi
	fi
}

function systemctlEnableAndCheck() {
    systemctl enable $1
    systemctl is-enabled $1
    enabled=$?
    for i in {1..900}; do
        if [ $enabled -ne 0 ]; then
            systemctl enable $1
            systemctl is-enabled $1
            enabled=$?
        else
            break
        fi
        sleep 1
    done
    if [ $enabled -ne 0 ]
    then
        exit 5
    fi
}

function ensureDocker() {
    systemctlEnableAndCheck docker
    if ! $REBOOTREQUIRED; then
        retrycmd_if_failure 5 5 timeout 60s systemctl restart docker
        dockerStarted=1
        for i in {1..900}; do
            if ! /usr/bin/docker info; then
                retrycmd_if_failure 5 5 timeout 60s systemctl restart docker
            else
                dockerStarted=0
                break
            fi
            sleep 1
        done
        if [ $dockerStarted -ne 0 ]
        then
            exit 2
        fi
    fi
}

function ensureKubelet() {
    systemctlEnableAndCheck kubelet
    if ! $REBOOTREQUIRED; then
        retrycmd_if_failure 20 10 timeout 60s systemctl restart kubelet
    fi
}

function extractKubectl(){
    systemctlEnableAndCheck kubectl-extract
    if ! $REBOOTREQUIRED; then
        systemctl restart kubectl-extract
    fi
}

function ensureJournal(){
    retrycmd_if_failure 5 5 timeout 30s systemctl daemon-reload
    systemctlEnableAndCheck systemd-journald.service
    echo "Storage=persistent" >> /etc/systemd/journald.conf
    echo "SystemMaxUse=1G" >> /etc/systemd/journald.conf
    echo "RuntimeMaxUse=1G" >> /etc/systemd/journald.conf
    echo "ForwardToSyslog=no" >> /etc/systemd/journald.conf
    if ! $REBOOTREQUIRED; then
        systemctl restart systemd-journald.service
    fi
}

function ensureK8s() {
    if $REBOOTREQUIRED; then
        return
    fi
    k8sHealthy=1
    nodesActive=1
    nodesReady=1
    for i in {1..600}; do
        if [ -e $KUBECTL ]
        then
            break
        fi
        sleep 1
    done
    for i in {1..600}; do
        $KUBECTL 2>/dev/null cluster-info
            if [ "$?" = "0" ]
            then
                k8sHealthy=0
                break
            fi
        sleep 1
    done
    if [ $k8sHealthy -ne 0 ]
    then
        exit 3
    fi
    for i in {1..1800}; do
        nodes=$(${KUBECTL} get nodes 2>/dev/null | grep 'Ready' | wc -l)
            if [ $nodes -eq $TOTAL_NODES ]
            then
                nodesActive=0
                break
            fi
        sleep 1
    done
    if [ $nodesActive -ne 0 ]
    then
        exit 3
    fi
    for i in {1..600}; do
        notReady=$(${KUBECTL} get nodes 2>/dev/null | grep 'NotReady' | wc -l)
            if [ $notReady -eq 0 ]
            then
                nodesReady=0
                break
            fi
        sleep 1
    done
    if [ $nodesReady -ne 0 ]
    then
        exit 3
    fi
}

function ensureEtcd() {
    etcdIsRunning=1
    for i in {1..600}; do
        curl --cacert /etc/kubernetes/certs/ca.crt --cert /etc/kubernetes/certs/etcdclient.crt --key /etc/kubernetes/certs/etcdclient.key --max-time 60 https://127.0.0.1:2379/v2/machines;
        if [ $? -eq 0 ]
        then
            etcdIsRunning=0
            break
        fi
        sleep 1
    done
    if [ $etcdIsRunning -ne 0 ]
    then
        exit 3
    fi
}

function ensureEtcdDataDir() {
    mount | grep /dev/sdc1 | grep /var/lib/etcddisk
    if [ "$?" = "0" ]
    then
        return
    else
        s = 5
        for i in {1..60}; do
            sudo mount -a && mount | grep /dev/sdc1 | grep /var/lib/etcddisk;
            if [ "$?" = "0" ]
            then
                (( t = ${i} * ${s} ))
                return
            fi
            sleep $s
        done
    fi

   exit 4
}

function ensurePodSecurityPolicy(){
    if $REBOOTREQUIRED; then
        return
    fi
    POD_SECURITY_POLICY_FILE="/etc/kubernetes/manifests/pod-security-policy.yaml"
    if [ -f $POD_SECURITY_POLICY_FILE ]; then
        $KUBECTL create -f $POD_SECURITY_POLICY_FILE
    fi
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
    set -x
}

if [[ "$CONTAINER_RUNTIME" == "clear-containers" ]]; then
	ensureRunCommandCompleted
	echo `date`,`hostname`, RunCmdCompleted>>/opt/m
fi

if [[ $OS == $UBUNTU_OS_NAME ]]; then
	apt-mark hold walinuxagent
fi

echo `date`,`hostname`, EnsureDockerStart>>/opt/m
ensureDockerInstallCompleted
ensureDocker
echo `date`,`hostname`, configNetworkPolicyStart>>/opt/m
configNetworkPolicy
if [[ "$CONTAINER_RUNTIME" == "clear-containers" ]]; then
	if grep -q vmx /proc/cpuinfo; then
		echo `date`,`hostname`, installClearContainersRuntimeStart>>/opt/m
		installClearContainersRuntime
		echo `date`,`hostname`, buildCRIOStart>>/opt/m
		buildCRIO
	fi
fi
echo `date`,`hostname`, setMaxPodsStart>>/opt/m
setMaxPods ${MAX_PODS}
echo `date`,`hostname`, ensureCRIOStart>>/opt/m
ensureCRIO
echo `date`,`hostname`, ensureKubeletStart>>/opt/m
ensureKubelet
echo `date`,`hostname`, extractKubctlStart>>/opt/m
extractKubectl
echo `date`,`hostname`, ensureJournalStart>>/opt/m
ensureJournal
echo `date`,`hostname`, ensureJournalDone>>/opt/m

ensureRunCommandCompleted
echo `date`,`hostname`, RunCmdCompleted>>/opt/m

if [[ ! -z "${MASTER_NODE}" ]]; then
    writeKubeConfig
    ensureFilepath $KUBECTL
    ensureFilepath $DOCKER
    ensureEtcdDataDir
    ensureEtcd
    ensureK8s
    ensurePodSecurityPolicy
fi

if [[ $OS == $UBUNTU_OS_NAME ]]; then
    # mitigation for bug https://bugs.launchpad.net/ubuntu/+source/linux/+bug/1676635
    echo 2dd1ce17-079e-403c-b352-a1921ee207ee > /sys/bus/vmbus/drivers/hv_util/unbind
    sed -i "13i\echo 2dd1ce17-079e-403c-b352-a1921ee207ee > /sys/bus/vmbus/drivers/hv_util/unbind\n" /etc/rc.local
    apt-mark unhold walinuxagent
fi

if $REBOOTREQUIRED; then
  /bin/bash -c "shutdown -r 1 &"
fi

echo `date`,`hostname`, endscript>>/opt/m

mkdir -p /opt/azure/containers && touch /opt/azure/containers/provision.complete
ps auxfww > /opt/azure/provision-ps.log &
