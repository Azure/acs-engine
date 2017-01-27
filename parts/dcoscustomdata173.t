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
runcmd:
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
- - systemctl
  - stop
  - resolvconf.service
- - systemctl
  - disable
  - resolvconf.service
- - systemctl
  - stop
  - lxc-net.service
- - systemctl
  - disable
  - lxc-net.service
- - systemctl
  - mask
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
  - dcos-link-env.service
- - systemctl
  - start
  - dcos-download.service
- - systemctl
  - start
  - dcos-setup.service
- - systemctl
  - start
  - dcos-config-writer.service
write_files:
- content: 'https://az837203.vo.msecnd.net/dcos/testing

    '
  owner: root
  path: /etc/mesosphere/setup-flags/repository-url
  permissions: '0644'
- content: 'BOOTSTRAP_ID=df308b6fc3bd91e1277baa5a3db928ae70964722

    '
  owner: root
  path: /etc/mesosphere/setup-flags/bootstrap-id
  permissions: '0644'
- content: '["dcos-config--setup_DCOSGUID", "dcos-metadata--setup_DCOSGUID"]

    '
  owner: root
  path: /etc/mesosphere/setup-flags/cluster-packages.json
  permissions: '0644'
- content: '[Journal]

    MaxLevelConsole=warning

    '
  owner: root
  path: /etc/systemd/journald.conf.d/dcos.conf
  permissions: '0644'
- content: '{{{nameSuffix}}}

    '
  path: /etc/mesosphere/cluster-id
  permissions: '0644'
- content: "\nrexray:\n  loglevel: info\n  modules:\n    default-docker:\n      disabled:\
    \ true\n"
  path: /etc/rexray/config.yml
  permissions: '0644'
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

    After=network-online.target

    Wants=network-online.target

    ConditionPathExists=!/opt/mesosphere/

    [Service]

    EnvironmentFile=/etc/mesosphere/setup-flags/bootstrap-id

    Type=oneshot

    StandardOutput=journal+console

    StandardError=journal+console

    ExecStartPre=/usr/bin/curl -fLsSv --retry 20 -Y 100000 -y 60
    -o /var/lib/mesos/dl/bootstrap.tar.xz https://az837203.vo.msecnd.net/dcos/testing/bootstrap/${BOOTSTRAP_ID}.bootstrap.tar.xz

    ExecStartPre=/usr/bin/mkdir -p /opt/mesosphere

    ExecStart=/usr/bin/tar -axf /var/lib/mesos/dl/bootstrap.tar.xz -C /opt/mesosphere

    ExecStartPost=-/usr/bin/rm -f /var/lib/mesos/dl/bootstrap.tar.xz

    '
  path: /etc/systemd/system/dcos-download.service
  permissions: '0644'
- content: '[Unit]

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

    ExecStart=/usr/bin/bash -c "echo -e \"127.0.0.1 localhost\n$(detect_ip) $(hostname)\"
    > /etc/hosts"

    '
  path: /etc/systemd/system/dcos-config-writer.service
  permissions: '0644'
- content: 'MESOS_CLUSTER={{{masterPublicIPAddressName}}}

    '
  path: /etc/mesosphere/setup-packages/dcos-provider-azure--setup/etc/mesos-master-provider
- content: '

    ADMINROUTER_ACTIVATE_AUTH_MODULE={{{oauthEnabled}}}

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
- content: 'com.netflix.exhibitor.azure.account-name={{{masterStorageAccountName}}}

    com.netflix.exhibitor.azure.account-key='', listKeys(resourceId(''Microsoft.Storage/storageAccounts'',
    variables(''masterStorageAccountName'')), ''2015-06-15'').key1,''

    '
  path: /etc/mesosphere/setup-packages/dcos-provider-azure--setup/etc/exhibitor.properties
- content: '{"uiConfiguration":{"plugins":{"banner":{"enabled":false,"backgroundColor":"#1E232F","foregroundColor":"#FFFFFF","headerTitle":null,"headerContent":null,"footerContent":null,"imagePath":null,"dismissible":null},"branding":{"enabled":false},"external-links":
    {"enabled": false},


    "authentication":{"enabled":false},


    "oauth":{"enabled":{{{oauthEnabled}}},"authHost":"https://dcos.auth0.com"},



    "networking":{"enabled":false},"organization":{"enabled":false},"tracking":{"enabled":false}}}}

    '
  path: /etc/mesosphere/setup-packages/dcos-provider-azure--setup/etc/ui-config.json
- content: '{}'
  path: /etc/mesosphere/setup-packages/dcos-provider-azure--setup/pkginfo.json
- content: ''
  path: /etc/mesosphere/roles/azure
- content: 'PROVISION_STR'
  path: "/opt/azure/containers/provision.sh"
  permissions: "0744"
  owner: "root"
