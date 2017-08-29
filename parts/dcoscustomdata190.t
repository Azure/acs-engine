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
- /opt/azure/containers/add_admin_to_docker_group.sh
write_files:
- content: 'https://dcosio.azureedge.net/dcos/stable

'
  owner: root
  path: /etc/mesosphere/setup-flags/repository-url
  permissions: '0644'
- content: '["3dt--7847ebb24bf6756c3103902971b34c3f09c3afbd", "adminrouter--0493a6fdaed08e1971871818e194aa4607df4f09",
    "avro-cpp--760c214063f6b038b522eaf4b768b905fed56ebc", "boost-libs--2015ccb58fb756f61c02ee6aa05cc1e27459a9ec",
    "bootstrap--59a905ecee27e71168ed44cefda4481fb76b816d", "boto--6344d31eef082c7bd13259b17034ea7b5c34aedf",
    "check-time--be7d0ba757ec87f9965378fee7c76a6ee5ae996d", "cni--e48337da39a8cd379414acfe0da52a9226a10d24",
    "cosmos--20decef90f0623ed253a12ec4cf5c148b18d8249", "curl--fc3486c43f98e63f9b12675f1356e8fe842f26b0",
    "dcos-config--setup_DCOSGUID", "dcos-history--77b0e97d7b25c8bedf8f7da0689cac65b83e3813",
    "dcos-image--bda6a02bcb2eb21c4218453a870cc584f921a800", "dcos-image-deps--83584fd868e5b470f7cf754424a9a75b328e9b68",
    "dcos-integration-test--c28bcb2347799dca43083f55e4c7b28503176f9c", "dcos-log--4d630df863228f38c6333e44670b4c4b20a74832",
    "dcos-metadata--setup_DCOSGUID", "dcos-metrics--23ee2f89c58b1258bc959f1d0dd7debcbb3d79d2",
    "dcos-oauth--0079529da183c0f23a06d2b069721b6fa6cc7b52", "dcos-signal--1bcd3b612cbdc379380dcba17cdf9a3b6652d9dc",
    "dcos-ui--d4afd695796404a5b35950c3daddcae322481ac4", "dnspython--0f833eb9a8abeba3179b43f3a200a8cd42d3795a",
    "docker-gc--59a98ed6446a084bf74e4ff4b8e3479f59ea8528", "dvdcli--5374dd4ffb519f1dcefdec89b2247e3404f2e2e3",
    "erlang--a9ee2530357a3301e53056b36a93420847b339a3", "exhibitor--72d9d8f947e5411eda524d40dde1a58edeb158ed",
    "flask--26d1bcdb2d1c3dcf1d2c03bc0d4f29c86d321b21", "java--cd5e921ce66b0d3303883c06d73a657314044304",
    "libevent--208be855d2be29c9271a7bd6c04723ff79946e02", "libffi--83ce3bd7eda2ef089e57efd2bc16c144d5a1f094",
    "libsodium--9ff915db08c6bba7d6738af5084e782b13c84bf8", "logrotate--7f7bc4416d3ad101d0c5218872858483b516be07",
    "marathon--bfb24f7f90cb3cd52a1cb22a07caafa5013bba21", "mesos--aaedd03eee0d57f5c0d49c74ff1e5721862cad98",
    "mesos-dns--0401501b2b5152d01bfa84ff6d007fdafe414b16", "mesos-modules--311849eaae42696b8a7eefe86b9ab3ebd9bd48f5",
    "metronome--467e4c64f804dbd4cd8572516e111a3f9298c10d", "navstar--1128db0234105a64fb4be52f4453cd6aa895ff30",
    "ncurses--d889894b71aa1a5b311bafef0e85479025b4dacb", "octarine--e86d3312691b12523280d56f6260216729aaa0ad",
    "openssl--b01a32a42e3ccba52b417276e9509a441e1d4a82", "pkgpanda-api--541feb8a8be58bdde8fecf1d2e5bfa0515f5a7d0",
    "pkgpanda-role--f8a749a4a821476ad2ef7e9dd9d12b6a8c4643a4", "pytest--78aee3e58a049cdab0d266af74f77d658b360b4f",
    "python--b7a144a49577a223d37d447c568f51330ee95390", "python-azure-mgmt-resource--03c05550f43b0e7a4455c33fe43b0deb755d87f0",
    "python-cryptography--4184767c68e48801dd394072cb370c610a05029d", "python-dateutil--fdc6ff929f65dd0918cf75a9ad56704683d31781",
    "python-docopt--beba78faa13e5bf4c52393b4b82d81f3c391aa65", "python-gunicorn--a537f95661fb2689c52fe12510eb0d01cb83af60",
    "python-isodate--40d378c688e6badfd16676dd8b51b742bfebc8d5", "python-jinja2--7450f5ae5a822f63f7a58c717207be0456df51ed",
    "python-kazoo--cb7ce13a1068cd82dd84ea0de32b529a760a4bdd", "python-markupsafe--dd46d2a3c58611656a235f96d4adc51b2a7a590e",
    "python-passlib--802ec3605c0b82428fedba60983b1bafaa036bb8", "python-pyyaml--81dd44cc4a24db7cefa7016c6586a131acf279c3",
    "python-requests--1b2cadbd3811cc0c2ee235ce927e13ea1d6af41d", "python-retrying--eb7b8bac133f50492b1e1349cbe77c3e38bd02c3",
    "python-tox--07244f8a939a10353634c952c6d88ec4a3c05736", "rexray--869621bb411c9f2a793ea42cdfeed489e1972aaa",
    "six--f06424b68523c4dfa2a7c3e7475d479f3d361e42", "spartan--9cc57a3d55452b905d90e3201f56913140914ecc",
    "strace--7d01796d64994451c1b2b82d161a335cbe90569b", "teamcity-messages--e623a4d86eb3a8d199cefcc240dd4c5460cb2962",
    "toybox--f235594ab8ea9a2864ee72abe86723d76f92e848"]

'
  owner: root
  path: /etc/mesosphere/setup-flags/cluster-packages.json
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
- content: |
    #!/bin/bash
    adduser {{{adminUsername}}} docker
  path: "/opt/azure/containers/add_admin_to_docker_group.sh"
  permissions: "0744"
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
