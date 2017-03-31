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
- sed -i "s/^Port 22$/Port 22\nPort 2222/1" /etc/ssh/sshd_config
- service ssh restart 
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
- content: '["3dt--6a71ec3c3407eb25c6bf2330326dc49b3de3c2eb", "adminrouter--ffc5b908bbba1c7e87ce09c84f0835e3f960fc8e", "avro-cpp--7b355a85f39ca6dbe2468ec50b71f3787c6c7c3d", "boost-libs--2015ccb58fb756f61c02ee6aa05cc1e27459a9ec", "bootstrap--6e05035d265bd327d2ec114101fd292dc0aaf3a3", "boto--6344d31eef082c7bd13259b17034ea7b5c34aedf", "check-time--be7d0ba757ec87f9965378fee7c76a6ee5ae996d", "cni--e48337da39a8cd379414acfe0da52a9226a10d24", "cosmos--93d021389b92d4c08c7e2236da510da69b1c632f", "curl--1148f64e03819f381cda4dc2e8a6199fb3c53a7e", "dcos-config--setup_DCOSGUID", "dcos-history--f8e3cc66dc1b9e01800e721ee980c09f3a8dfe46", "dcos-image--e637ab1daad8d81eea7f9be042394a94c42a39d6", "dcos-image-deps--3ed9dee844359c415123cb6fb6b306f215faab2a", "dcos-integration-test--0fb256ff2c38ff751eaf2ce4748273a8338b4441", "dcos-log--b542bb89a5af9642e04df35869beee4ce253e535", "dcos-metadata--setup_DCOSGUID", "dcos-metrics--41f4d0b1b84b8e8fe2876baeb3bd07ce873a54e0", "dcos-oauth--0079529da183c0f23a06d2b069721b6fa6cc7b52", "dcos-signal--1bcd3b612cbdc379380dcba17cdf9a3b6652d9dc", "dcos-ui--da8a5003a3c5ec478f89b18a5a216a0ea7bb1d62", "dnspython--0f833eb9a8abeba3179b43f3a200a8cd42d3795a", "docker-gc--59a98ed6446a084bf74e4ff4b8e3479f59ea8528", "dvdcli--5374dd4ffb519f1dcefdec89b2247e3404f2e2e3", "erlang--c88d0e71b0bd2900612498095d3ac320ae9ff80d", "exhibitor--72d9d8f947e5411eda524d40dde1a58edeb158ed", "flask--26d1bcdb2d1c3dcf1d2c03bc0d4f29c86d321b21", "java--cd5e921ce66b0d3303883c06d73a657314044304", "libevent--208be855d2be29c9271a7bd6c04723ff79946e02", "libsodium--9ff915db08c6bba7d6738af5084e782b13c84bf8", "logrotate--faf6c640a994ac549afe734e05d322ab9052448b", "marathon--fa629c85fc11eceffce921aeaf43d1eac2ee4a7d", "mesos--3ee073c6f436f77d94bcd0af0648d6f26e2ec197", "mesos-dns--f374ceda1dfade3eacdbdfed0d57bcf88c905242", "mesos-modules--7ef1d3c2691c64e84f1b60da4f014aea926daef7", "metronome--4328a268b5139ab5bc2e942b28d748d6815763b5", "navstar--b1ed66efe8fe7bd7e0138a66a51558c8cc486060", "ncurses--d889894b71aa1a5b311bafef0e85479025b4dacb", "octarine--521813a6f6459dc1e0e32e161999b95ed9eacbac", "openssl--b01a32a42e3ccba52b417276e9509a441e1d4a82", "pkgpanda-api--20de028f4e65672f301a187e46f12330d9f836cc", "pkgpanda-role--f8a749a4a821476ad2ef7e9dd9d12b6a8c4643a4", "pytest--78aee3e58a049cdab0d266af74f77d658b360b4f", "python--b7a144a49577a223d37d447c568f51330ee95390", "python-azure-mgmt-resource--9e68c5bacce73c50d9b313d660f402dffca9d39e", "python-dateutil--fdc6ff929f65dd0918cf75a9ad56704683d31781", "python-docopt--beba78faa13e5bf4c52393b4b82d81f3c391aa65", "python-gunicorn--a537f95661fb2689c52fe12510eb0d01cb83af60", "python-isodate--40d378c688e6badfd16676dd8b51b742bfebc8d5", "python-jinja2--7450f5ae5a822f63f7a58c717207be0456df51ed", "python-kazoo--cb7ce13a1068cd82dd84ea0de32b529a760a4bdd", "python-markupsafe--dd46d2a3c58611656a235f96d4adc51b2a7a590e", "python-passlib--802ec3605c0b82428fedba60983b1bafaa036bb8", "python-pyyaml--81dd44cc4a24db7cefa7016c6586a131acf279c3", "python-requests--1b2cadbd3811cc0c2ee235ce927e13ea1d6af41d", "python-retrying--eb7b8bac133f50492b1e1349cbe77c3e38bd02c3", "python-tox--07244f8a939a10353634c952c6d88ec4a3c05736", "rexray--869621bb411c9f2a793ea42cdfeed489e1972aaa", "six--f06424b68523c4dfa2a7c3e7475d479f3d361e42", "spartan--58a5611725de935357a0d96b2caef838ebc99b79", "strace--7d01796d64994451c1b2b82d161a335cbe90569b", "teamcity-messages--e623a4d86eb3a8d199cefcc240dd4c5460cb2962", "toybox--f235594ab8ea9a2864ee72abe86723d76f92e848"]

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
- content: "rexray:\n  loglevel: info\n  modules:\n    default-admin:\n      host:\
    \ tcp://127.0.0.1:61003\n    default-docker:\n      disabled: true\n"
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
    ExecStartPre=/usr/bin/curl -fLsSv --retry 20 -Y 100000 -y 60 -o /var/tmp/d.deb https://az837203.vo.msecnd.net/dcos-deps/docker-engine_1.11.2-0~xenial_amd64.deb
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