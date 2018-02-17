#!/bin/bash

###########################################################
# START SECRET DATA - ECHO DISABLED
###########################################################

# Following parameters now read from environment variable
# Fields for `azure.json`
# TENANT_ID SUBSCRIPTION_ID RESOURCE_GROUP LOCATION SUBNET
# NETWORK_SECURITY_GROUP VIRTUAL_NETWORK VIRTUAL_NETWORK_RESOURCE_GROUP ROUTE_TABLE PRIMARY_AVAILABILITY_SET
# SERVICE_PRINCIPAL_CLIENT_ID SERVICE_PRINCIPAL_CLIENT_SECRET KUBELET_PRIVATE_KEY TARGET_ENVIRONMENT NETWORK_POLICY
# FQDNSuffix VNET_CNI_PLUGINS_URL CNI_PLUGINS_URL MAX_PODS

# Default values for backoff configuration
# CLOUDPROVIDER_BACKOFF CLOUDPROVIDER_BACKOFF_RETRIES CLOUDPROVIDER_BACKOFF_EXPONENT CLOUDPROVIDER_BACKOFF_DURATION CLOUDPROVIDER_BACKOFF_JITTER
# Default values for rate limit configuration
# CLOUDPROVIDER_RATELIMIT CLOUDPROVIDER_RATELIMIT_QPS CLOUDPROVIDER_RATELIMIT_BUCKET

# USE_MANAGED_IDENTITY_EXTENSION USE_INSTANCE_METADATA

# Master only secrets
# APISERVER_PRIVATE_KEY CA_CERTIFICATE CA_PRIVATE_KEY MASTER_FQDN KUBECONFIG_CERTIFICATE
# KUBECONFIG_KEY ETCD_SERVER_CERTIFICATE ETCD_SERVER_PRIVATE_KEY ETCD_CLIENT_CERTIFICATE ETCD_CLIENT_PRIVATE_KEY
# ETCD_PEER_CERTIFICATES ETCD_PEER_PRIVATE_KEYS ADMINUSER MASTER_INDEX

set -x
# Capture Interesting Network Stuffs during provision
packetCaptureProvision() {
    tcpdump -G 600 -W 1 -n -vv -w /var/log/azure/dnsdump.pcap -Z root -i eth0 udp port 53 > /dev/null 2>&1 &
}

packetCaptureProvision

# Find distro name via ID value in releases files and upcase
OS=$(cat /etc/*-release | grep ^ID= | tr -d 'ID="' | awk '{print toupper($0)}')
UBUNTU_OS_NAME="UBUNTU"
RHEL_OS_NAME="RHEL"
COREOS_OS_NAME="COREOS"

# Set default filepaths
KUBECTL=/usr/local/bin/kubectl
DOCKER=/usr/bin/docker

set +x
ETCD_PEER_CERT=$(echo ${ETCD_PEER_CERTIFICATES} | cut -d'[' -f 2 | cut -d']' -f 1 | cut -d',' -f $((${MASTER_INDEX}+1)))
ETCD_PEER_KEY=$(echo ${ETCD_PEER_PRIVATE_KEYS} | cut -d'[' -f 2 | cut -d']' -f 1 | cut -d',' -f $((${MASTER_INDEX}+1)))
set -x

# CoreOS: /usr is read-only; therefore kubectl is installed at /opt/kubectl
#   Details on install at kubernetetsmastercustomdataforcoreos.yml
if [[ $OS == $COREOS_OS_NAME ]]; then
    echo "Changing default kubectl bin location"
    KUBECTL=/opt/kubectl
fi

# cloudinit runcmd and the extension will run in parallel, this is to ensure
# runcmd finishes
ensureRunCommandCompleted()
{
    echo "waiting for runcmd to finish"
    for i in {1..900}; do
        if [ -e /opt/azure/containers/runcmd.complete ]; then
            echo "runcmd finished, took $i seconds"
            break
        fi
        sleep 1
    done
}

# cloudinit runcmd and the extension will run in parallel, this is to ensure
# runcmd finishes
ensureDockerInstallCompleted()
{
    echo "waiting for docker install to finish"
    for i in {1..900}; do
        if [ -e /opt/azure/containers/dockerinstall.complete ]; then
            echo "docker install finished, took $i seconds"
            break
        fi
        sleep 1
    done
}

echo `date`,`hostname`, startscript>>/opt/m

# A delay to start the kubernetes processes is necessary
# if a reboot is required.  Otherwise, the agents will encounter issue:
# https://github.com/kubernetes/kubernetes/issues/41185
if [ -f /var/run/reboot-required ]; then
    REBOOTREQUIRED=true
else
    REBOOTREQUIRED=false
fi

if [[ ! -z "${MASTER_NODE}" ]]; then
    echo "executing master node provision operations"
    
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
else
    echo "skipping master node provision operations, this is an agent node"
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

###########################################################
# END OF SECRET DATA
###########################################################

set -x

# wait for presence of a file
function ensureFilepath() {
    if $REBOOTREQUIRED; then
        return
    fi
    found=1
    for i in {1..600}; do
        if [ -e $1 ]
        then
            found=0
            echo "$1 is present, took $i seconds to verify"
            break
        fi
        sleep 1
    done
    if [ $found -ne 0 ]
    then
        echo "$1 is not present after $i seconds of trying to verify"
        exit 1
    fi
}

function downloadUrl () {
	# Wrapper around curl to download blobs more reliably.
	# Workaround the --retry issues with a for loop and set a max timeout.
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

    # Download Azure VNET CNI plugins.
    CNI_BIN_DIR=/opt/cni/bin
    mkdir -p $CNI_BIN_DIR

    # Mirror from https://github.com/Azure/azure-container-networking/releases/tag/$AZURE_PLUGIN_VER/azure-vnet-cni-linux-amd64-$AZURE_PLUGIN_VER.tgz
    downloadUrl ${VNET_CNI_PLUGINS_URL} | tar -xz -C $CNI_BIN_DIR
    # Mirror from https://github.com/containernetworking/cni/releases/download/$CNI_RELEASE_VER/cni-amd64-$CNI_RELEASE_VERSION.tgz
    downloadUrl ${CNI_PLUGINS_URL} | tar -xz -C $CNI_BIN_DIR ./loopback ./portmap
    chown -R root:root $CNI_BIN_DIR
    chmod -R 755 $CNI_BIN_DIR

    # Copy config file
    mv $CNI_BIN_DIR/10-azure.conflist $CNI_CONFIG_DIR/
    chmod 600 $CNI_CONFIG_DIR/10-azure.conflist

    # Dump ebtables rules.
    /sbin/ebtables -t nat --list

    # Enable CNI.
    setNetworkPlugin cni
    setDockerOpts " --volume=/etc/cni/:/etc/cni:ro --volume=/opt/cni/:/opt/cni:ro"
}

# Configures Kubelet to use CNI and mount the appropriate hostpaths
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
        # No policy, defaults to kubenet.
        setNetworkPlugin kubenet
        setDockerOpts ""
    fi
}

# Install the Clear Containers runtime
function installClearContainersRuntime() {
	# Add Clear Containers repository key
	echo "Adding Clear Containers repository key..."
	curl -sSL --retry 5 --retry-delay 10 --retry-max-time 30 "https://download.opensuse.org/repositories/home:clearcontainers:clear-containers-3/xUbuntu_16.04/Release.key" | apt-key add -

	# Add Clear Container repository
	echo "Adding Clear Containers repository..."
	echo 'deb http://download.opensuse.org/repositories/home:/clearcontainers:/clear-containers-3/xUbuntu_16.04/ /' > /etc/apt/sources.list.d/cc-runtime.list

	# Install Clear Containers runtime
	echo "Installing Clear Containers runtime..."
	apt-get update
	apt-get install --no-install-recommends -y \
		cc-runtime

	# Install thin tools for devicemapper configuration
	echo "Installing thin tools to provision devicemapper..."
	apt-get install --no-install-recommends -y \
		lvm2 \
		thin-provisioning-tools

	# Load systemd changes
	echo "Loading changes to systemd service files..."
	systemctl daemon-reload

	# Enable and start Clear Containers proxy service
	echo "Enabling and starting Clear Containers proxy service..."
	systemctl enable cc-proxy
	systemctl start cc-proxy

	# CRIO has only been tested with the azure plugin
	configAzureNetworkPolicy
	setKubeletOpts " --container-runtime=remote --container-runtime-endpoint=/var/run/crio.sock"
	setDockerOpts " --volume=/etc/cni/:/etc/cni:ro --volume=/opt/cni/:/opt/cni:ro"
}

# Install Go from source
function installGo() {
	export GO_SRC=/usr/local/go
	export GOPATH="${HOME}/.go"

	# Remove any old version of Go
	if [[ -d "$GO_SRC" ]]; then
		rm -rf "$GO_SRC"
	fi

	# Remove any old GOPATH
	if [[ -d "$GOPATH" ]]; then
		rm -rf "$GOPATH"
	fi

	# Get the latest Go version
	GO_VERSION=$(curl --retry 5 --retry-delay 10 --retry-max-time 30 -sSL "https://golang.org/VERSION?m=text")

	echo "Installing Go version $GO_VERSION..."

	# subshell
	(
	curl --retry 5 --retry-delay 10 --retry-max-time 30 -sSL "https://storage.googleapis.com/golang/${GO_VERSION}.linux-amd64.tar.gz" | sudo tar -v -C /usr/local -xz
	)

	# Set GOPATH and update PATH
	echo "Setting GOPATH and updating PATH"
	export PATH="${GO_SRC}/bin:${PATH}:${GOPATH}/bin"
}

# Build and install runc
function buildRunc() {
	# Clone the runc source
	echo "Cloning the runc source..."
	mkdir -p "${GOPATH}/src/github.com/opencontainers"
	(
	cd "${GOPATH}/src/github.com/opencontainers"
	git clone "https://github.com/opencontainers/runc.git"
	cd runc
	git reset --hard v1.0.0-rc4
	make BUILDTAGS="seccomp apparmor"
	make install
	)

	echo "Successfully built and installed runc..."
}

# Build and install CRI-O
function buildCRIO() {
	# Add CRI-O repositories
	echo "Adding repositories required for cri-o..."
	add-apt-repository -y ppa:projectatomic/ppa
	add-apt-repository -y ppa:alexlarsson/flatpak
	apt-get update

	# Install CRI-O dependencies
	echo "Installing dependencies for CRI-O..."
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

	# Install md2man
	go get github.com/cpuguy83/go-md2man

	# Fix for templates dependency
	(
	go get -u github.com/docker/docker/daemon/logger/templates
	cd "${GOPATH}/src/github.com/docker/docker"
	mkdir -p utils
	cp -r daemon/logger/templates utils/
	)

	buildRunc;

	# Clone the CRI-O source
	echo "Cloning the CRI-O source..."
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

	echo "Successfully built and installed CRI-O..."

	# Cleanup the temporary directory
	rm -vrf "$tmpd"

	# Cleanup the Go install
	rm -vrf "$GO_SRC" "$GOPATH"

	setupCRIO;
}

# Setup CRI-O
function setupCRIO() {
	# Configure CRI-O
	echo "Configuring CRI-O..."

	# Configure crio systemd service file
	SYSTEMD_CRI_O_SERVICE_FILE="/usr/local/lib/systemd/system/crio.service"
	sed -i 's#ExecStart=/usr/local/bin/crio#ExecStart=/usr/local/bin/crio -log-level debug#' "$SYSTEMD_CRI_O_SERVICE_FILE"

	# Configure /etc/crio/crio.conf
	CRI_O_CONFIG="/etc/crio/crio.conf"
	sed -i 's#storage_driver = ""#storage_driver = "devicemapper"#' "$CRI_O_CONFIG"
	sed -i 's#storage_option = \[#storage_option = \["dm.directlvm_device=/dev/sdc", "dm.thinp_percent=95", "dm.thinp_metapercent=1", "dm.thinp_autoextend_threshold=80", "dm.thinp_autoextend_percent=20", "dm.directlvm_device_force=true"#' "$CRI_O_CONFIG"
	sed -i 's#runtime = "/usr/bin/runc"#runtime = "/usr/local/sbin/runc"#' "$CRI_O_CONFIG"
	sed -i 's#runtime_untrusted_workload = ""#runtime_untrusted_workload = "/usr/bin/cc-runtime"#' "$CRI_O_CONFIG"
	sed -i 's#default_workload_trust = "trusted"#default_workload_trust = "untrusted"#' "$CRI_O_CONFIG"

	# Load systemd changes
	echo "Loading changes to systemd service files..."
	systemctl daemon-reload
}

function ensureCRIO() {
	if [[ "$CONTAINER_RUNTIME" == "clear-containers" ]]; then
		# Make sure we can nest virtualization
		if grep -q vmx /proc/cpuinfo; then
			# Enable and start cri-o service
			# Make sure this is done after networking plugins are installed
			echo "Enabling and starting cri-o service..."
			systemctl enable crio crio-shutdown
			systemctl start crio
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
            echo "$1 took $i seconds to be enabled by systemctl"
            break
        fi
        sleep 1
    done
    if [ $enabled -ne 0 ]
    then
        echo "$1 could not be enabled by systemctl"
        exit 5
    fi
}

function ensureDocker() {
    systemctlEnableAndCheck docker
    # only start if a reboot is not required
    if ! $REBOOTREQUIRED; then
        systemctl restart docker
        dockerStarted=1
        for i in {1..900}; do
            if ! /usr/bin/docker info; then
                echo "status $?"
                /bin/systemctl restart docker
            else
                echo "docker started, took $i seconds"
                dockerStarted=0
                break
            fi
            sleep 1
        done
        if [ $dockerStarted -ne 0 ]
        then
            echo "docker did not start"
            exit 2
        fi
    fi
}

function ensureKubelet() {
    systemctlEnableAndCheck kubelet
    # only start if a reboot is not required
    if ! $REBOOTREQUIRED; then
        systemctl restart kubelet
    fi
}

function extractKubectl(){
    systemctlEnableAndCheck kubectl-extract
    # only start if a reboot is not required
    if ! $REBOOTREQUIRED; then
        systemctl restart kubectl-extract
    fi
}

function ensureJournal(){
    systemctl daemon-reload
    systemctlEnableAndCheck systemd-journald.service
    echo "Storage=persistent" >> /etc/systemd/journald.conf
    echo "SystemMaxUse=1G" >> /etc/systemd/journald.conf
    echo "RuntimeMaxUse=1G" >> /etc/systemd/journald.conf
    echo "ForwardToSyslog=no" >> /etc/systemd/journald.conf
    # only start if a reboot is not required
    if ! $REBOOTREQUIRED; then
        systemctl restart systemd-journald.service
    fi
}

function ensureApiserver() {
    if $REBOOTREQUIRED; then
        return
    fi
    kubernetesStarted=1
    for i in {1..600}; do
        if [ -e $KUBECTL ]
        then
            $KUBECTL cluster-info
            if [ "$?" = "0" ]
            then
                echo "kubernetes started, took $i seconds"
                kubernetesStarted=0
                break
            fi
        else
            /usr/bin/docker ps | grep apiserver
            if [ "$?" = "0" ]
            then
                echo "kubernetes started, took $i seconds"
                kubernetesStarted=0
                break
            fi
        fi
        sleep 1
    done
    if [ $kubernetesStarted -ne 0 ]
    then
        echo "kubernetes did not start"
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
            echo "Etcd setup successfully, took $i seconds"
            break
        fi
        sleep 1
    done
    if [ $etcdIsRunning -ne 0 ]
    then
        echo "Etcd not accessible after $i seconds"
        exit 3
    fi
}

function ensureEtcdDataDir() {
    mount | grep /dev/sdc1 | grep /var/lib/etcddisk
    if [ "$?" = "0" ]
    then
        echo "Etcd is running with data dir at: /var/lib/etcddisk"
        return
    else
        echo "/var/lib/etcddisk was not found at /dev/sdc1. Trying to mount all devices."
        s = 5
        for i in {1..60}; do
            sudo mount -a && mount | grep /dev/sdc1 | grep /var/lib/etcddisk;
            if [ "$?" = "0" ]
            then
                (( t = ${i} * ${s} ))
                echo "/var/lib/etcddisk mounted at: /dev/sdc1, took $t seconds"
                return
            fi
            sleep $s
        done
    fi

   echo "Etcd data dir was not found at: /var/lib/etcddisk"
   exit 4
}

function ensurePodSecurityPolicy(){
    if $REBOOTREQUIRED; then
        return
    fi
    POD_SECURITY_POLICY_FILE="/etc/kubernetes/manifests/pod-security-policy.yaml"
    if [ -f $POD_SECURITY_POLICY_FILE ]; then
        kubectl create -f $POD_SECURITY_POLICY_FILE
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

    # disable logging after secret output
    set +x
    echo "
---
apiVersion: v1
clusters:
- cluster:
    certificate-authority-data: \"$CA_CERTIFICATE\"
    server: https://$MASTER_FQDN.$LOCATION.$FQDNSuffix
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

if [[ "$CONTAINER_RUNTIME" == "clear-containers" ]]; then
	# If the container runtime is "clear-containers" we need to ensure the
	# run command is completed _before_ we start installing all the dependencies
	# for clear-containers to make sure there is not a dpkg lock.
	ensureRunCommandCompleted
	echo `date`,`hostname`, RunCmdCompleted>>/opt/m
fi

if [[ $OS == $UBUNTU_OS_NAME ]]; then
	# make sure walinuxagent doesn't get updated in the middle of running this script
	apt-mark hold walinuxagent
fi

# master and node
echo `date`,`hostname`, EnsureDockerStart>>/opt/m
ensureDockerInstallCompleted
ensureDocker
echo `date`,`hostname`, configNetworkPolicyStart>>/opt/m
configNetworkPolicy
if [[ "$CONTAINER_RUNTIME" == "clear-containers" ]]; then
	# Ensure we can nest virtualization
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

# On all other runtimes, but "clear-containers" we can ensure the run command
# completed here to allow for parallelizing the custom script
ensureRunCommandCompleted
echo `date`,`hostname`, RunCmdCompleted>>/opt/m

# master only
if [[ ! -z "${MASTER_NODE}" ]]; then
    writeKubeConfig
    ensureFilepath $KUBECTL
    ensureFilepath $DOCKER
    ensureEtcdDataDir
    ensureEtcd
    ensureApiserver
    ensurePodSecurityPolicy
fi

if [[ $OS == $UBUNTU_OS_NAME ]]; then
    # mitigation for bug https://bugs.launchpad.net/ubuntu/+source/linux/+bug/1676635
    echo 2dd1ce17-079e-403c-b352-a1921ee207ee > /sys/bus/vmbus/drivers/hv_util/unbind
    sed -i "13i\echo 2dd1ce17-079e-403c-b352-a1921ee207ee > /sys/bus/vmbus/drivers/hv_util/unbind\n" /etc/rc.local

    apt-mark unhold walinuxagent
fi

echo "Install complete successfully"

if $REBOOTREQUIRED; then
  # wait 1 minute to restart node, so that the custom script extension can complete
  echo 'reboot required, rebooting node in 1 minute'
  /bin/bash -c "shutdown -r 1 &"
fi

echo `date`,`hostname`, endscript>>/opt/m

mkdir -p /opt/azure/containers && touch /opt/azure/containers/provision.complete
ps auxfww > /opt/azure/provision-ps.log &
