#!/bin/bash

###########################################################
# Configure Mesos One Box
#
# This installs the following components
# - zookeepr
# - mesos master
# - marathon
# - mesos agent
###########################################################
set -x

echo "starting mesos cluster configuration"
date
ps ax

SWARM_VERSION="swarm:1.1.0"
#############
# Parameters
#############
MASTERCOUNT=${1}
MASTERPREFIX=${2}
MASTERFIRSTADDR=${3}
SWARMENABLED=`echo ${4} | awk '{print tolower($0)}'`
MARATHONENABLED=`echo ${5} | awk '{print tolower($0)}'`
CHRONOSENABLED=`echo ${6} | awk '{print tolower($0)}'`
ACCOUNTNAME=${7}
set +x
ACCOUNTKEY=${8}
set -x
AZUREUSER=${9}
POSTINSTALLSCRIPTURI=${10}
BASESUBNET=${11}
HOMEDIR="/home/$AZUREUSER"
VMNAME=`hostname`
VMNUMBER=`echo $VMNAME | sed 's/.*[^0-9]\([0-9]\+\)*$/\1/'`
VMPREFIX=`echo $VMNAME | sed 's/\(.*[^0-9]\)*[0-9]\+$/\1/'`

echo "Master Count: $MASTERCOUNT"
echo "Master Prefix: $MASTERPREFIX"
echo "Master First Addr: $MASTERFIRSTADDR"
echo "vmname: $VMNAME"
echo "VMNUMBER: $VMNUMBER, VMPREFIX: $VMPREFIX"
echo "SWARMENABLED: $SWARMENABLED, MARATHONENABLED: $MARATHONENABLED, CHRONOSENABLED: $CHRONOSENABLED"
echo "ACCOUNTNAME: $ACCOUNTNAME"
echo "BASESUBNET: $BASESUBNET"

###################
# Common Functions
###################
ensureAzureNetwork()
{
  # ensure the host name is resolvable
  hostResolveHealthy=1
  for i in {1..120}; do
    host $VMNAME
    if [ $? -eq 0 ]
    then
      # hostname has been found continue
      hostResolveHealthy=0
      echo "the host name resolves"
      break
    fi
    sleep 1
  done
  if [ $hostResolveHealthy -ne 0 ]
  then
    echo "host name does not resolve, aborting install"
    exit 1
  fi

  # ensure the network works
  networkHealthy=1
  for i in {1..12}; do
    wget -O/dev/null http://bing.com
    if [ $? -eq 0 ]
    then
      # hostname has been found continue
      networkHealthy=0
      echo "the network is healthy"
      break
    fi
    sleep 10
  done
  if [ $networkHealthy -ne 0 ]
  then
    echo "the network is not healthy, cannot download from bing, aborting install"
    ifconfig
    ip a
    exit 2
  fi
  # ensure the hostname -i works
  networkHealthy=1
  for i in {1..120}; do
    hostname -i
    if [ $? -eq 0 ]
    then
      # hostname has been found continue
      networkHealthy=0
      echo "the network is healthy"
      break
    fi
    sleep 1
  done
  if [ $networkHealthy -ne 0 ]
  then
    echo "the network is not healthy, cannot resolve ip address, aborting install"
    ifconfig
    ip a
    exit 2
  fi
  # ensure hostname -f works
  networkHealthy=1
  for i in {1..120}; do
    hostname -f
    if [ $? -eq 0 ]
    then
      # hostname has been found continue
      networkHealthy=0
      echo "the network is healthy"
      break
    fi
    sleep 1
  done
  if [ $networkHealthy -ne 0 ]
  then
    echo "the network is not healthy, cannot resolve hostname, aborting install"
    ifconfig
    ip a
    exit 2
  fi
}

ensureAzureNetwork
HOSTADDR=`hostname -i`

ismaster ()
{
  if [ "$MASTERPREFIX" == "$VMPREFIX" ]
  then
    return 0
  else
    return 1
  fi
}
if ismaster ; then
  echo "this node is a master"
fi

isagent()
{
  if ismaster ; then
    return 1
  else
    return 0
  fi
}
if isagent ; then
  echo "this node is an agent"
fi

isomsrequired()
{
  if [ $ACCOUNTNAME != "none" ]
  then
    return 0
  else
    return 1
  fi
}
if isomsrequired ; then
  echo "this node requires oms"
fi

zkhosts()
{
  zkhosts=""
  for i in `seq 0 $((MASTERCOUNT-1))` ;
  do
    if [ "$i" -gt "0" ]
    then
      zkhosts="${zkhosts},"
    fi

    MASTEROCTET=`expr $MASTERFIRSTADDR + $i`
    IPADDR="${BASESUBNET}${MASTEROCTET}"
    zkhosts="${zkhosts}${IPADDR}:2181"
    # due to mesos team experience ip addresses are chosen over dns names
    #zkhosts="${zkhosts}${MASTERPREFIX}${i}:2181"
  done
  echo $zkhosts
}

zkconfig()
{
  postfix="$1"
  zkhosts=$(zkhosts)
  zkconfigstr="zk://${zkhosts}/${postfix}"
  echo $zkconfigstr
}

######################
# resolve self in DNS
######################
echo "$HOSTADDR $VMNAME" | sudo tee -a /etc/hosts

################
# Install Docker
################
echo "Installing and configuring docker and swarm"
installDocker()
{
  for i in {1..10}; do
    wget --tries 4 --retry-connrefused --waitretry=15 -qO- https://get.docker.com | sh
    if [ $? -eq 0 ]
    then
      # hostname has been found continue
      echo "Docker installed successfully"
      break
    fi
    sleep 10
  done
}

time installDocker

sudo usermod -aG docker $AZUREUSER
if isagent ; then
  # Start Docker and listen on :2375 (no auth, but in vnet)
  echo 'DOCKER_OPTS="-H unix:///var/run/docker.sock -H 0.0.0.0:2375"' | sudo tee -a /etc/default/docker
fi

if isomsrequired ; then
  # the following insecure registry is for OMS
  echo 'DOCKER_OPTS="$DOCKER_OPTS --insecure-registry 137.135.93.9"' | sudo tee -a /etc/default/docker
fi

sudo service docker restart

ensureDocker()
{
  # ensure that docker is healthy
  dockerHealthy=1
  for i in {1..3}; do
    sudo docker info
    if [ $? -eq 0 ]
    then
      # hostname has been found continue
      dockerHealthy=0
      echo "Docker is healthy"
      sudo docker ps -a
      break
    fi
    sleep 10
  done
  if [ $dockerHealthy -ne 0 ]
  then
    echo "Docker is not healthy"
  fi
}

ensureDocker

############
# setup OMS
############
if isomsrequired ; then
  set +x
  EPSTRING="DefaultEndpointsProtocol=https;AccountName=${ACCOUNTNAME};AccountKey=${ACCOUNTKEY}"
  docker run --restart=always -d 137.135.93.9/msdockeragentv3 http://${VMNAME}:2375 "${EPSTRING}"
  set -x
fi

##################
# Install Mesos
##################
sudo apt-key adv --keyserver keyserver.ubuntu.com --recv E56151BF
DISTRO=$(lsb_release -is | tr '[:upper:]' '[:lower:]')
CODENAME=$(lsb_release -cs)
echo "deb http://repos.mesosphere.io/${DISTRO} ${CODENAME} main" | sudo tee /etc/apt/sources.list.d/mesosphere.list
time sudo add-apt-repository -y ppa:openjdk-r/ppa
time sudo apt-get -y update
time sudo DEBIAN_FRONTEND=noninteractive apt-get -y --force-yes install openjdk-8-jre-headless
if ismaster ; then
  time sudo DEBIAN_FRONTEND=noninteractive apt-get -y --force-yes install mesosphere
else
  time sudo DEBIAN_FRONTEND=noninteractive apt-get -y --force-yes install mesos
fi

#########################
# Configure ZooKeeper
#########################
zkmesosconfig=$(zkconfig "mesos")
echo $zkmesosconfig | sudo tee /etc/mesos/zk

if ismaster ; then
  echo $VMNUMBER | sudo tee /etc/zookeeper/conf/myid
  for i in `seq 0 $((MASTERCOUNT-1))` ;
  do
    MASTEROCTET=`expr $MASTERFIRSTADDR + $i`
    IPADDR="${BASESUBNET}${MASTEROCTET}"
    echo "server.${i}=${IPADDR}:2888:3888" | sudo tee -a /etc/zookeeper/conf/zoo.cfg
    # due to mesos team experience ip addresses are chosen over dns names
    #echo "server.${i}=${MASTERPREFIX}${i}:2888:3888" | sudo tee -a /etc/zookeeper/conf/zoo.cfg
  done
fi

#########################################
# Configure Mesos Master and Frameworks
#########################################
if ismaster ; then
  quorum=`expr $MASTERCOUNT / 2 + 1`
  echo $quorum | sudo tee /etc/mesos-master/quorum
  hostname -i | sudo tee /etc/mesos-master/ip
  hostname | sudo tee /etc/mesos-master/hostname
  echo 'Mesos Cluster on Microsoft Azure' | sudo tee /etc/mesos-master/cluster
fi

if ismaster  && [ "$MARATHONENABLED" == "true" ] ; then
  # setup marathon
  sudo mkdir -p /etc/marathon/conf
  sudo cp /etc/mesos-master/hostname /etc/marathon/conf
  sudo cp /etc/mesos/zk /etc/marathon/conf/master
  zkmarathonconfig=$(zkconfig "marathon")
  echo $zkmarathonconfig | sudo tee /etc/marathon/conf/zk
  # enable marathon to failover tasks to other nodes immediately
  echo 0 | sudo tee /etc/marathon/conf/failover_timeout
  #echo false | sudo tee /etc/marathon/conf/checkpoint
fi

#########################################
# Configure Mesos Master and Frameworks
#########################################
if ismaster ; then
  # Download and install mesos-dns
  sudo mkdir -p /usr/local/mesos-dns
  sudo wget --tries 4 --retry-connrefused --waitretry=15 https://github.com/mesosphere/mesos-dns/releases/download/v0.5.1/mesos-dns-v0.5.1-linux-amd64 -O mesos-dns-linux
  sudo chmod +x mesos-dns-linux
  sudo mv mesos-dns-linux /usr/local/mesos-dns/mesos-dns
  RESOLVER=`cat /etc/resolv.conf | grep nameserver | tail -n 1 | awk '{print $2}'`

  COUNT=$((MASTERCOUNT-1))
  #generate a list of master's for input to zk config
  MASTERS=""
  for i  in `seq 0 $COUNT` ;
   do
     MASTEROCTET=`expr $MASTERFIRSTADDR + $i`
     IPADDR="$BASESUBNET$MASTEROCTET:5050"
     MASTERS="$MASTERS\"${IPADDR}\""
     if [ "$i" -lt "$COUNT" ]; then
        MASTERS="$MASTERS,"
     fi
done
echo "
{
  \"zk\": \"${zkmesosconfig}\",
  \"masters\": ["${MASTERS}"],
  \"refreshSeconds\": 1,
  \"ttl\": 1,
  \"domain\": \"mesos\",
  \"port\": 53,
  \"timeout\": 1,
  \"listener\": \"0.0.0.0\",
  \"email\": \"root.mesos-dns.mesos\",
  \"resolvers\": [\"$RESOLVER\"]
}
" > mesos-dns.json
  sudo mv mesos-dns.json /usr/local/mesos-dns/mesos-dns.json

  echo "
description \"mesos dns\"

# Start just after the System-V jobs (rc) to ensure networking and zookeeper
# are started. This is as simple as possible to ensure compatibility with
# Ubuntu, Debian, CentOS, and RHEL distros. See:
# http://upstart.ubuntu.com/cookbook/#standard-idioms
start on stopped rc RUNLEVEL=[2345]
respawn

exec /usr/local/mesos-dns/mesos-dns -config /usr/local/mesos-dns/mesos-dns.json" > mesos-dns.conf
  sudo mv mesos-dns.conf /etc/init
fi

#########################
# Configure Mesos Agent
#########################
if isagent ; then
  # Add docker containerizer
  echo "docker,mesos" | sudo tee /etc/mesos-slave/containerizers
  # Add timeout for agent to download docker
  echo '5mins' > /etc/mesos-slave/executor_registration_timeout
  # Add resources configuration
  if ismaster ; then
    echo "ports:[1-21,23-79,81-4399,4401-5049,5052-8079,8081-32000]" | sudo tee /etc/mesos-slave/resources
  else
    echo "ports:[1-21,23-5050,5052-32000]" | sudo tee /etc/mesos-slave/resources
  fi
  hostname -i | sudo tee /etc/mesos-slave/ip
  hostname | sudo tee /etc/mesos-slave/hostname
fi
# Add mesos-dns IP addresses to the head file, so they are at the top of the file
for i in `seq 0 $((MASTERCOUNT-1))` ;
do
    MASTEROCTET=`expr $MASTERFIRSTADDR + $i`
    IPADDR="${BASESUBNET}${MASTEROCTET}"
    echo nameserver $IPADDR | sudo tee -a /etc/resolvconf/resolv.conf.d/head
done
cat /etc/resolvconf/resolv.conf.d/head

##############################################
# configure init rules restart all processes
##############################################
echo "stop mesos and framework processes, they will restart after reboot"
if ismaster ; then
  echo manual | sudo tee /etc/init/mesos-slave.override
  sudo service mesos-slave stop

  # stop all running services
  sudo service marathon stop
  sudo service chronos stop
  sudo service mesos-dns stop
  sudo service mesos-master stop
  sudo service zookeeper stop

  # the following will clear out any corrupt zookeeper state, and zookeeper will
  # reconstruct this on the reboot at the end of provisioning
  sudo mkdir /var/lib/zookeeperbackup
  sudo mv /var/lib/zookeeper/* /var/lib/zookeeperbackup
  sudo cp /var/lib/zookeeperbackup/myid /var/lib/zookeeper/
else
  echo manual | sudo tee /etc/init/zookeeper.override
  sudo service zookeeper stop
  echo manual | sudo tee /etc/init/mesos-master.override
  sudo service mesos-master stop
fi

if ismaster && [ "$SWARMENABLED" == "true" ] && [ $VMNUMBER -eq "0" ]; then
  echo "starting docker swarm version $SWARM_VERSION"
  echo "sleep 10 seconds to give master time to come up"
  sleep 10
  echo sudo docker run -d --net=host -e SWARM_MESOS_USER=root \
      --restart=always \
      $SWARM_VERSION manage \
      -c mesos-experimental \
      --cluster-opt mesos.address=$HOSTADDR \
      --cluster-opt mesos.port=3375 $zkmesosconfig
  sudo docker run -d --net=host -e SWARM_MESOS_USER=root \
      --restart=always \
      $SWARM_VERSION manage \
      -c mesos-experimental \
      --cluster-opt mesos.address=$HOSTADDR \
      --cluster-opt mesos.port=3375 $zkmesosconfig
  sudo docker ps
  echo "completed starting docker swarm"
fi

###################
# Install Admin Router
###################
installMesosAdminRouter()
{
  sudo DEBIAN_FRONTEND=noninteractive apt-get -y --force-yes install nginx-extras lua-cjson
  # the admin router comes from https://github.com/mesosphere/adminrouter-public
  ADMIN_ROUTER_GITHUB_URL=https://raw.githubusercontent.com/anhowe/adminrouter-public/master
  NGINX_CONF_PATH=/usr/share/nginx/conf
  sudo mkdir -p $NGINX_CONF_PATH
  wget --tries 4 --retry-connrefused --waitretry=15 -qO$NGINX_CONF_PATH/common.lua $ADMIN_ROUTER_GITHUB_URL/common.lua
  wget --tries 4 --retry-connrefused --waitretry=15 -qO$NGINX_CONF_PATH/metadata.lua $ADMIN_ROUTER_GITHUB_URL/metadata.lua
  wget --tries 4 --retry-connrefused --waitretry=15 -qO$NGINX_CONF_PATH/service.lua $ADMIN_ROUTER_GITHUB_URL/service.lua
  wget --tries 4 --retry-connrefused --waitretry=15 -qO$NGINX_CONF_PATH/slave.lua $ADMIN_ROUTER_GITHUB_URL/slave.lua
  wget --tries 4 --retry-connrefused --waitretry=15 -qO$NGINX_CONF_PATH/slavehostname.lua $ADMIN_ROUTER_GITHUB_URL/slavehostname.lua
  wget --tries 4 --retry-connrefused --waitretry=15 -qO$NGINX_CONF_PATH/url.lua $ADMIN_ROUTER_GITHUB_URL/url.lua

  sudo mv /etc/nginx/nginx.conf /etc/nginx/nginx.conf.orig
  sudo cp /opt/azure/containers/nginx.conf /etc/nginx/nginx.conf
}

# only install the mesos dcos cli on the master
if ismaster ; then
  time installMesosAdminRouter
fi

###################
# Install Mesos DCOS CLI
###################
installMesosDCOSCLI()
{
  sudo DEBIAN_FRONTEND=noninteractive apt-get -y --force-yes install -y python-pip
  sudo pip install virtualenv
  sudo -i -u $AZUREUSER mkdir $HOMEDIR/dcos
  for i in {1..10}; do
    wget --tries 4 --retry-connrefused --waitretry=15 -qO- https://raw.githubusercontent.com/mesosphere/dcos-cli/master/bin/install/install-optout-dcos-cli.sh | sudo -i -u $AZUREUSER /bin/bash -s $HOMEDIR/dcos/. http://leader.mesos --add-path yes
    if [ $? -eq 0 ]
    then
      echo "Mesos DCOS-CLI installed successfully"
      break
    fi
    sleep 10
  done
}

# only install the mesos dcos cli on the master
if ismaster ; then
  time installMesosDCOSCLI
fi

###################
# Post Install
###################
if [ $POSTINSTALLSCRIPTURI != "disabled" ]
then
  echo "downloading, and kicking off post install script"
  /bin/bash -c "wget --tries 20 --retry-connrefused --waitretry=15 -qO- $POSTINSTALLSCRIPTURI | nohup /bin/bash >> /var/log/azure/cluster-bootstrap-postinstall.log 2>&1 &"
fi

ps ax
echo "Finished installing and configuring docker and swarm"
date
echo "completed mesos cluster configuration"

echo "restart system to install any remaining software"
if isagent ; then
  shutdown -r now
else
  # wait 30s for guest agent to communicate back success and then reboot
  /usr/bin/nohup /bin/bash -c "sleep 30s; shutdown -r now" &
fi
