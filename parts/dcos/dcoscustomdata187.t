bootcmd:
- bash -c "if [ ! -f /var/lib/sdb-gpt ];then echo DCOS-5890;parted -s /dev/sdb mklabel
  gpt;touch /var/lib/sdb-gpt;fi"
disk_setup:
  ephemeral0:
    layout:
    - 50
    - 50
    overwrite: true
    table_type: gpt
fs_setup:
- device: ephemeral0.1
  filesystem: ext4
  overwrite: true
- device: ephemeral0.2
  filesystem: ext4
  overwrite: true
mounts:
- - ephemeral0.1
  - /var/lib/mesos
- - ephemeral0.2
  - /var/lib/docker
runcmd: PREPROVISION_EXTENSION
- /usr/lib/apt/apt.systemd.daily
- echo 2dd1ce17-079e-403c-b352-a1921ee207ee > /sys/bus/vmbus/drivers/hv_util/unbind # mitigation for bug https://bugs.launchpad.net/ubuntu/+source/linux/+bug/1676635
- sed -i "13i\echo 2dd1ce17-079e-403c-b352-a1921ee207ee > /sys/bus/vmbus/drivers/hv_util/unbind\n" /etc/rc.local # mitigation for bug https://bugs.launchpad.net/ubuntu/+source/linux/+bug/1676635
- - ln
  - -s
  - /bin/rm
  - /usr/bin/rm
- - ln
  - -s
  - /bin/mkdir
  - /usr/bin/mkdir
- - ln
  - -s
  - /bin/tar
  - /usr/bin/tar
- - ln
  - -s
  - /bin/ln
  - /usr/bin/ln
- - ln
  - -s
  - /bin/cp
  - /usr/bin/cp
- - ln
  - -s
  - /bin/systemctl
  - /usr/bin/systemctl
- - ln
  - -s
  - /bin/mount
  - /usr/bin/mount
- - ln
  - -s
  - /bin/bash
  - /usr/bin/bash
- - ln
  - -s
  - /usr/sbin/useradd
  - /usr/bin/useradd
- - systemctl
  - disable
  - --now
  - resolvconf.service
- - systemctl
  - mask
  - --now
  - lxc-net.service
- /opt/azure/containers/provision.sh
- - cp
  - -p
  - /etc/resolv.conf
  - /tmp/resolv.conf
- - rm
  - -f
  - /etc/resolv.conf
- - cp
  - -p
  - /tmp/resolv.conf
  - /etc/resolv.conf
- - systemctl
  - start
  - dcos-docker-install.service
- - systemctl
  - start
  - dcos-config-writer.service
- - systemctl
  - restart
  - systemd-journald.service
- - systemctl
  - restart
  - docker.service
- - systemctl
  - start
  - dcos-link-env.service
- - systemctl
  - enable
  - dcos-setup.service
- - systemctl
  - --no-block
  - start
  - dcos-setup.service
write_files:
- content: 'https://dcosio.azureedge.net/dcos/stable

    '
  owner: root
  path: /etc/mesosphere/setup-flags/repository-url
  permissions: '0644'
- content: '["dcos-config--setup_DCOSGUID", "dcos-metadata--setup_DCOSGUID"]

    '
  owner: root
  path: /etc/mesosphere/setup-flags/cluster-packages.json
  permissions: '0644'
- content: 'DCOS_ENVIRONMENT={{{targetEnvironment}}}

    '
  owner: root
  path: /etc/mesosphere/setup-flags/dcos-deploy-environment
  permissions: '0644'
- content: '[Journal]

    MaxLevelConsole=warning

    RateLimitInterval=1s

    RateLimitBurst=20000

    '
  owner: root
  path: /etc/systemd/journald.conf.d/dcos.conf
  permissions: '0644'
- content: "rexray:\n  loglevel: info\n  modules:\n    default-admin:\n      host:\
    \ tcp://127.0.0.1:61003\n    default-docker:\n      disabled: true\n"
  path: /etc/rexray/config.yml
  permissions: '0644'
- content: '[Unit]

    After=network-online.target

    Wants=network-online.target

    [Service]

    Type=oneshot

    Environment=DEBIAN_FRONTEND=noninteractive

    StandardOutput=journal+console

    StandardError=journal+console

    ExecStart=/usr/bin/bash -c "try=1;until dpkg -D3 -i /var/lib/mesos/dl/d.deb || ((try>9));do
    echo retry $((try++));sleep $((try*try));done;systemctl --now start docker;systemctl
    restart docker.socket"

    '
  path: /etc/systemd/system/dcos-docker-install.service
  permissions: '0644'
- content: '[Service]

    Restart=always

    StartLimitInterval=0

    RestartSec=15

    ExecStartPre=-/sbin/ip link del docker0

    ExecStart=

    ExecStart=/usr/bin/docker daemon -H fd:// --storage-driver=overlay

    '
  path: /etc/systemd/system/docker.service.d/execstart.conf
  permissions: '0644'
- content: '[Unit]

    PartOf=docker.service

    [Socket]

    ListenStream=/var/run/docker.sock

    SocketMode=0660

    SocketUser=root

    SocketGroup=docker

    ListenStream=2375

    BindIPv6Only=both

    [Install]

    WantedBy=sockets.target

    '
  path: /etc/systemd/system/docker.socket
  permissions: '0644'
- content: '[Unit]

    Requires=dcos-setup.service

    After=dcos-setup.service

    [Service]

    Type=oneshot

    EnvironmentFile=/etc/environment

    EnvironmentFile=/opt/mesosphere/environment

    ExecStart=/usr/bin/bash -c "echo $(detect_ip) $(hostname) > /etc/hosts"

    '
  path: /etc/systemd/system/dcos-config-writer.service
  permissions: '0644'
- content: 'MESOS_CLUSTER={{{masterPublicIPAddressName}}}

    '
  path: /etc/mesosphere/setup-packages/dcos-provider-azure--setup/etc/mesos-master-provider
- content: 'ADMINROUTER_ACTIVATE_AUTH_MODULE={{{oauthEnabled}}}

    '
  path: /etc/mesosphere/setup-packages/dcos-provider-azure--setup/etc/adminrouter.env
- content: '["'', DCOSCUSTOMDATAPUBLICIPSTR''"]

    '
  path: /etc/mesosphere/setup-packages/dcos-provider-azure--setup/etc/master_list
- content: 'EXHIBITOR_BACKEND=AZURE

    AZURE_CONTAINER=dcos-exhibitor

    AZURE_PREFIX={{{masterPublicIPAddressName}}}

    '
  path: /etc/mesosphere/setup-packages/dcos-provider-azure--setup/etc/exhibitor
- content: 'com.netflix.exhibitor.azure.account-name={{{masterStorageAccountExhibitorName}}}

    com.netflix.exhibitor.azure.account-key='', listKeys(resourceId(''Microsoft.Storage/storageAccounts'',
    variables(''masterStorageAccountExhibitorName'')), ''2015-06-15'').key1,''

    '
  path: /etc/mesosphere/setup-packages/dcos-provider-azure--setup/etc/exhibitor.properties
- content: '{"uiConfiguration":{"plugins":{"banner":{"enabled":false,"backgroundColor":"#1E232F","foregroundColor":"#FFFFFF","headerTitle":null,"headerContent":null,"footerContent":null,"imagePath":null,"dismissible":null},"branding":{"enabled":false},"external-links":
    {"enabled": false},


    "authentication":{"enabled":false},


    "oauth":{"enabled":{{{oauthEnabled}}},"authHost":"https://dcos.auth0.com"},



    "tracking":{"enabled":false}}}}

    '
  path: /etc/mesosphere/setup-packages/dcos-provider-azure--setup/etc/ui-config.json
- content: '{}'
  path: /etc/mesosphere/setup-packages/dcos-provider-azure--setup/pkginfo.json
- content: '[Unit]

    Before=dcos.target

    [Service]

    Type=oneshot

    StandardOutput=journal+console

    StandardError=journal+console

    ExecStartPre=/usr/bin/mkdir -p /etc/profile.d

    ExecStart=/usr/bin/ln -sf /opt/mesosphere/environment.export /etc/profile.d/dcos.sh

    '
  path: /etc/systemd/system/dcos-link-env.service
  permissions: '0644'
- content: '[Unit]

    Description=Pkgpanda: Download DC/OS to this host.

    After=network-online.target

    Wants=network-online.target

    ConditionPathExists=!/opt/mesosphere/

    [Service]

    Type=oneshot

    StandardOutput=journal+console

    StandardError=journal+console

    ExecStartPre=/usr/bin/curl --keepalive-time 2 -fLsSv --retry 20 -Y 100000 -y 60 -o /var/lib/mesos/dl/bootstrap.tar.xz {{{dcosBootstrapURL}}}

    ExecStartPre=/usr/bin/mkdir -p /opt/mesosphere

    ExecStart=/usr/bin/tar -axf /var/lib/mesos/dl/bootstrap.tar.xz -C /opt/mesosphere

    ExecStartPost=-/usr/bin/rm -f /var/lib/mesos/dl/bootstrap.tar.xz

    '
  path: /etc/systemd/system/dcos-download.service
  permissions: '0644'
- content: '[Unit]

    Description=Pkgpanda: Specialize DC/OS for this host.

    Requires=dcos-download.service

    After=dcos-download.service

    [Service]

    Type=oneshot

    StandardOutput=journal+console

    StandardError=journal+console

    EnvironmentFile=/opt/mesosphere/environment

    ExecStart=/opt/mesosphere/bin/pkgpanda setup --no-block-systemd

    [Install]

    WantedBy=multi-user.target

    '
  path: /etc/systemd/system/dcos-setup.service
  permissions: '0644'
- content: ''
  path: /etc/mesosphere/roles/azure
- content: 'PROVISION_STR'
  path: "/opt/azure/containers/provision.sh"
  permissions: "0744"
  owner: "root"
- path: /var/lib/dcos/mesos-slave-common
  content: 'ATTRIBUTES_STR'
  permissions: "0644"
  owner: "root"
