bootcmd:
- bash -c "if [ ! -f /var/lib/sdb-gpt ];then echo DCOS-5890;parted -s /dev/sdb mklabel
  gpt;touch /var/lib/sdb-gpt;fi"
- bash -c "if [ ! -f /var/lib/sdc-gpt ];then echo DCOS-5890;parted -s /dev/sdc mklabel
  gpt&&touch /var/lib/sdc-gpt;fi"
- bash -c "if [ ! -f /var/lib/sdd-gpt ];then echo DCOS-5890;parted -s /dev/sdd mklabel
  gpt&&touch /var/lib/sdd-gpt;fi"
- bash -c "if [ ! -f /var/lib/sde-gpt ];then echo DCOS-5890;parted -s /dev/sde mklabel
  gpt&&touch /var/lib/sde-gpt;fi"
- bash -c "if [ ! -f /var/lib/sdf-gpt ];then echo DCOS-5890;parted -s /dev/sdf mklabel
  gpt&&touch /var/lib/sdf-gpt;fi"
- bash -c "mkdir -p /dcos/volume{0,1,2,3}"
disk_setup:
  ephemeral0:
    layout:
    - 45
    - 45
    - 10
    overwrite: true
    table_type: gpt
  /dev/sdc:
    layout: true
    overwrite: true
    table_type: gpt
  /dev/sdd:
    layout: true
    overwrite: true
    table_type: gpt
  /dev/sde:
    layout: true
    overwrite: true
    table_type: gpt
  /dev/sdf:
    layout: true
    overwrite: true
    table_type: gpt
fs_setup:
- device: ephemeral0.1
  filesystem: ext4
  overwrite: true
- device: ephemeral0.2
  filesystem: ext4
  overwrite: true
- device: ephemeral0.3
  filesystem: ext4
  overwrite: true
- device: /dev/sdc1
  filesystem: ext4
  overwrite: true
- device: /dev/sdd1
  filesystem: ext4
  overwrite: true
- device: /dev/sde1
  filesystem: ext4
  overwrite: true
- device: /dev/sdf1
  filesystem: ext4
  overwrite: true
mounts:
- - ephemeral0.1
  - /var/lib/mesos
- - ephemeral0.2
  - /var/lib/docker
- - ephemeral0.3
  - /var/tmp
- - /dev/sdc1
  - /dcos/volume0
- - /dev/sdd1
  - /dcos/volume1
- - /dev/sde1
  - /dcos/volume2
- - /dev/sdf1
  - /dcos/volume3
runcmd: PREPROVISION_EXTENSION
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
- - systemctl
  - disable
  - --now
  - unscd.service
- - systemctl
  - stop
  - --now
  - unscd.service
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
- content: '{{{dcosRepositoryURL}}}
'
  owner: root
  path: /etc/mesosphere/setup-flags/repository-url
  permissions: '0644'
- content: |-
- content: '{{{dcosClusterPackageListID}}}'
  owner: root
  path: /etc/mesosphere/setup-flags/cluster-package-list
  permissions: '0644'
- content: |
    [Journal]
    MaxLevelConsole=warning
    RateLimitInterval=1s
    RateLimitBurst=20000
  owner: root
  path: /etc/systemd/journald.conf.d/dcos.conf
  permissions: '0644'
- content: |
    rexray:
      loglevel: info
      modules:
        default-docker:
          disabled: true
      service: vfs
  path: /etc/rexray/config.yml
  permissions: '0644'
- content: |
    [Unit]
    After=network-online.target
    Wants=network-online.target
    [Service]
    Type=oneshot
    Environment=DEBIAN_FRONTEND=noninteractive
    StandardOutput=journal+console
    StandardError=journal+console
    ExecStartPre=/usr/bin/curl -fLsSv --retry 20 -Y 100000 -y 60 -o /var/tmp/d.deb https://download.docker.com/linux/ubuntu/dists/xenial/pool/stable/amd64/docker-ce_17.09.0~ce-0~ubuntu_amd64.deb
    ExecStart=/usr/bin/bash -c "try=1;until dpkg -D3 -i /var/tmp/d.deb || ((try>9));do echo retry $((try++));sleep $((try*try));done;systemctl --now start docker;systemctl restart docker.socket"
  path: /etc/systemd/system/dcos-docker-install.service
  permissions: '0644'
- content: |
    [Service]
    Restart=always
    StartLimitInterval=0
    RestartSec=15
    ExecStartPre=-/sbin/ip link del docker0
    ExecStart=
    ExecStart=/usr/bin/dockerd --storage-driver=overlay
  path: /etc/systemd/system/docker.service.d/execstart.conf
  permissions: '0644'
- content: |
    [Unit]
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
  path: /etc/systemd/system/docker.socket
  permissions: '0644'
- content: |
      [Unit]
      Requires=dcos-setup.service
      After=dcos-setup.service
      [Service]
      Type=oneshot
      EnvironmentFile=/etc/environment
      EnvironmentFile=/opt/mesosphere/environment
      ExecStart=/usr/bin/bash -c "echo $(detect_ip) $(hostname) > /etc/hosts"
  path: /etc/systemd/system/dcos-config-writer.service
  permissions: '0644'
- content: |
    "bound_values":
      "adminrouter_auth_enabled": |-
        {{{oauthEnabled}}}
      "cluster_name": |-
        {{{masterPublicIPAddressName}}}
      "exhibitor_azure_account_key": |-
        ', listKeys(resourceId('Microsoft.Storage/storageAccounts', variables('masterStorageAccountExhibitorName')), '2015-06-15').key1, '
      "exhibitor_azure_account_name": |-
        {{{masterStorageAccountExhibitorName}}}
      "exhibitor_azure_prefix": |-
        {{{masterPublicIPAddressName}}}
      "master_list": |-
        ["', DCOSCUSTOMDATAPUBLICIPSTR'"]
      "oauth_enabled": |-
        {{{oauthEnabled}}}
    "late_bound_package_id": |-
      dcos-provider-{{{dcosProviderPackageID}}}-azure--setup
  owner: root
  path: /etc/mesosphere/setup-flags/late-config.yaml
  permissions: '0644'
- content: |
    [Unit]
    Before=dcos.target
    [Service]
    Type=oneshot
    StandardOutput=journal+console
    StandardError=journal+console
    ExecStartPre=/usr/bin/mkdir -p /etc/profile.d
    ExecStart=/usr/bin/ln -sf /opt/mesosphere/bin/add_dcos_path.sh /etc/profile.d/dcos.sh
  path: /etc/systemd/system/dcos-link-env.service
  permissions: '0644'
- content: |
    [Unit]
    Description=Pkgpanda: Download DC/OS to this host.
    After=network-online.target
    Wants=network-online.target
    ConditionPathExists=!/opt/mesosphere/
    [Service]
    Type=oneshot
    StandardOutput=journal+console
    StandardError=journal+console
    ExecStartPre=/usr/bin/curl --keepalive-time 2 -fLsSv --retry 20 -Y 100000 -y 60 -o //var/tmp/bootstrap.tar.xz {{{dcosBootstrapURL}}}
    ExecStartPre=/usr/bin/mkdir -p /opt/mesosphere
    ExecStart=/usr/bin/tar -axf //var/tmp/bootstrap.tar.xz -C /opt/mesosphere
    ExecStartPost=-/usr/bin/rm -f //var/tmp/bootstrap.tar.xz
  path: /etc/systemd/system/dcos-download.service
  permissions: '0644'
- content: |
    [Unit]
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
  path: /etc/systemd/system/dcos-setup.service
  permissions: '0644'
- content: ''
  path: /etc/mesosphere/roles/azure
- path: /var/lib/dcos/mesos-slave-common
  content: 'ATTRIBUTES_STR'
  permissions: "0644"
  owner: "root"
- content: 'PROVISION_STR'
  path: /opt/azure/containers/provision.sh
  permissions: "0744"
  owner: "root"
