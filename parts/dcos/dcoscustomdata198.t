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
- content: '{{{dcosRepositoryURL}}}

'
  owner: root
  path: /etc/mesosphere/setup-flags/repository-url
  permissions: '0644'
- content: '["3dt--4eb6a10d16421bc87cb6e93ac97746f36aded925", "adminrouter--31f3f6390c8ef79a2774f42390d6340a24d67f08",
    "avro-cpp--6194e9a67928c357c1c1b2bb409536ceef888e04", "boost-libs--2015ccb58fb756f61c02ee6aa05cc1e27459a9ec",
    "bootstrap--d50592de9bf45937df7bcc7008e84a8739239c99", "boto--471853efd730e52e4ed7bfb890587432a576982a",
    "check-time--be7d0ba757ec87f9965378fee7c76a6ee5ae996d", "cni--e48337da39a8cd379414acfe0da52a9226a10d24",
    "cosmos--74e0339c91c278622d9f45b5fb0771872f443140", "curl--e7fd5880e4f94db05692d7e43279d8fe6348cb21",
    "dcos-config--setup_{{{dcosProviderPackageID}}}", "dcos-history--787ce2fd81cb7469590c12951033f0482e879d2a",
    "dcos-image--078703170a2f218447abea4b1be00b7431b340f1", "dcos-image-deps--5512ff49cdbba7f404759a5751a4ab1eae44c677",
    "dcos-integration-test--bad12974ed31ace44432ad9a451c5b5dc3e20e81", "dcos-log--4d630df863228f38c6333e44670b4c4b20a74832",
    "dcos-metadata--setup_{{{dcosProviderPackageID}}}", "dcos-metrics--e65d65e1b65335efdaa6bf7609a671f4288e7af9",
    "dcos-oauth--23d8ca77549c1ac6087c11c9f7e8f8a4fddfc948", "dcos-signal--5633dc8da7e864cb34e3d29ed13e6756c7a6df94",
    "dcos-ui--6f4af319cf4dd9bb8366de22ec37775beaa96747", "dnspython--1118f0ffaa60e6a779d4614f0ed692d215005f0e",
    "docker-gc--9737ec72de5d1edc71175028762f06fe22c8a48c", "dvdcli--5374dd4ffb519f1dcefdec89b2247e3404f2e2e3",
    "erlang--984871e11f69e37aeb76a471d4a4b90e93fdf355", "exhibitor--300da0c612afcf27541dbc681da5de3a6408de7e",
    "flask--2936647fa917d16ee289d34e61fd1afcc49157b5", "java--091eb5a0f3dcbd7762a43e84c3e2d6aac8891111",
    "libevent--468f4ae789f659e452e8356a9d2309c4f41135a8", "libffi--83ce3bd7eda2ef089e57efd2bc16c144d5a1f094",
    "libsodium--9ff915db08c6bba7d6738af5084e782b13c84bf8", "logrotate--7f7bc4416d3ad101d0c5218872858483b516be07",
    "marathon--99d0cbc65da6be31872878174f3a28fa63d0fa34", "mesos--0c992033b8d43e00dc69f0c548c826d573c82642",
    "mesos-dns--ca591a18f9b010999106285fedddd010606c0d06", "mesos-modules--4c176c23a4fd3670d059fec55e2d4c8c7dbf1f6c",
    "metronome--138ec50cd4da05bce74b6cd2c84ae873c2bd67ab", "navstar--fdf7e79fdf210548d183badfde00d60c1a540257",
    "ncurses--d889894b71aa1a5b311bafef0e85479025b4dacb", "octarine--4e37c062d2f145f9c2ce01d30dadf72c2aac5c4a",
    "openssl--ef04a6f76f6e5e194c783bc129fdabad16816aff", "pkgpanda-api--220e45fbd93403f8b4fd7f9c8c3d5178aff6e34b",
    "pkgpanda-role--f8a749a4a821476ad2ef7e9dd9d12b6a8c4643a4", "pytest--63ab7e9520e4da70202b81076880fcdf2c1236cf",
    "python--3c96ab7f21312f4d7d54a9b901cfe6382aa66b8a", "python-azure-mgmt-resource--2313114eec2adcb37ef61082cd2cfdceabf5c21e",
    "python-cryptography--39ee7d59411569700f3343e64c32e9711a83decc", "python-dateutil--d098c1933ca6d754a90734afd366d556cc3107a8",
    "python-docopt--85e7726dbb777584a9f5d4dd7bd58ed8ca5466d8", "python-gunicorn--bd425f55abd9236b5ead7e68a3c40c39b8d75bb7",
    "python-isodate--9a15007db453e141892966ebf50a9175ee0ba08b", "python-jinja2--9fbc35d1405f06f1959c54629ab7d443cef79076",
    "python-kazoo--050358610274815ebacabcdfca874729e53f4e0b", "python-markupsafe--09c65e6cdedd4783137a203cbc1b5a64ef3124eb",
    "python-passlib--27056b95ad1a067b7992402e679c6260e673a554", "python-pyyaml--5be319fd73348558d69a03fb6dcb134e9b7f4c48",
    "python-requests--63e1c3f4f03efc4607a4c20c5492026a9af7a9c7", "python-retrying--692b1a298d22436e25b2d14fc4f980be444adbe7",
    "python-tox--7962137d89dae9eb45dd80b0ea59731fa3f5bbc9", "rexray--f07795e2c10f9a1a27de9d8e67ab171029db2e1d",
    "six--9229b1a9d7d57bc086fa50f73fc9a753d9a4605d", "spartan--3dc1785bf698e65ceb2fecf26b2a439de219269f",
    "strace--7d01796d64994451c1b2b82d161a335cbe90569b", "teamcity-messages--d13bc3f52ed0e30de3a71d86ff8718984b60b65f",
    "toybox--c0e85790eb8aaeefe5037b053c2fcd140ab800a4"]

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
