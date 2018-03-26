bootcmd:
- bash -c "if [ ! -f /var/lib/sdb-gpt ];then echo DCOS-5890;parted -s /dev/sdb mklabel
  gpt;touch /var/lib/sdb-gpt;fi"
disk_setup:
  ephemeral0:
    layout:
    - 45
    - 45
    - 10
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
mounts:
- - ephemeral0.1
  - /var/lib/mesos
- - ephemeral0.2
  - /var/lib/docker
- - ephemeral0.3
  - /var/tmp
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
- - systemctl
  - disable
  - --now
  - unscd.service
- - systemctl
  - stop
  - --now
  - unscd.service
- sed -i "s/^Port 22$/Port 22\nPort 2222/1" /etc/ssh/sshd_config
- service ssh restart 
- /opt/azure/containers/setup_ephemeral_disk.sh
- - tar
  - czf 
  - /etc/docker.tar.gz
  - -C
  - /tmp/xtoph
  - .docker
- - rm 
  - -rf 
  - /tmp/xtoph
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
- content: '{{{dcosRepositoryURL}}}'
  owner: root
  path: /etc/mesosphere/setup-flags/repository-url
  permissions: '0644'
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
        default-admin:
          host: tcp://127.0.0.1:61003
        default-docker:
          disabled: true
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
    ExecStartPre=/usr/bin/curl -fLsSv --retry 20 -Y 100000 -y 60 -o /var/tmp/d.deb https://az837203.vo.msecnd.net/dcos-deps/docker-engine_1.13.1-0-ubuntu-xenial_amd64.deb
    ExecStart=/usr/bin/bash -c "try=1;until dpkg -D3 -i /var/tmp/d.deb || ((try>9));do echo retry $((try++));sleep $((try*try));done;systemctl --now start docker;systemctl restart docker.socket"
  path: /etc/systemd/system/dcos-docker-install.service
  permissions: '0644'
- content: |
    [Service]
    Restart=always
    StartLimitInterval=0
    RestartSec=15
    LimitNOFILE=16384
    ExecStartPre=-/sbin/ip link del docker0
    ExecStart=
    ExecStart=/usr/bin/docker daemon -H fd:// --storage-driver=overlay
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
  content: |
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
      dcos-provider-DCOSGUID-azure--setup
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
- content: 'PROVISION_STR'
  path: "/opt/azure/containers/provision.sh"
  permissions: "0744"
  owner: "root"
- path: /var/lib/dcos/mesos-slave-common
  content: 'ATTRIBUTES_STR'
  permissions: "0644"
  owner: "root"
- content: '{ "auths": { "{{{registry}}}": { "auth" : "{{{registryKey}}}" } } }'
  path: "/tmp/xtoph/.docker/config.json"
  owner: "root"
- content: |
    #!/bin/bash
    # Check the partitions on /dev/sdb created by cloudinit and force a detach and
    # reformat of the parition.  After which, all will be remounted.
    EPHEMERAL_DISK="/dev/sdb"
    PARTITIONS=`fdisk -l $EPHEMERAL_DISK | grep "^$EPHEMERAL_DISK" | cut -d" " -f1 | sed "s~$EPHEMERAL_DISK~~"`
    if [ -n "$PARTITIONS" ]; then
        for f in $PARTITIONS; do
            df -k | grep "/dev/sdb$f"
            if [ $? -eq 0 ]; then
                umount -f /dev/sdb$f
            fi
            mkfs.ext4 /dev/sdb$f
        done
        mount -a
    fi
    # If there is a /var/tmp partition on the ephemeral disk, create a symlink such
    # that the /var/log/mesos and /var/log/journal placed on the ephemeral disk.
    VAR_TMP_PARTITION=`df -P /var/tmp | tail -1 | cut -d" " -f 1`
    echo $VAR_TMP_PARTITION | grep "^$EPHEMERAL_DISK"
    if [ $? -eq 0 ]; then
        # Handle the /var/log/mesos directory
        mkdir -p /var/tmp/log/mesos
        if [ -d "/var/log/mesos" ]; then
            cp -rp /var/log/mesos/* /var/tmp/log/mesos/
            rm -rf /var/log/mesos
        fi
        ln -s /var/tmp/log/mesos /var/log/mesos
        # Handle the /var/log/journal direcotry
        mkdir -p /var/tmp/log/journal
        if [ -d "/var/log/journal" ]; then
            cp -rp /var/log/journal/* /var/tmp/log/journal/
            rm -rf /var/log/journal
        fi
        ln -s /var/tmp/log/journal /var/log/journal
    fi
  path: "/opt/azure/containers/setup_ephemeral_disk.sh"
  permissions: "0744"
  owner: "root"
