#!/bin/bash

###########################################################
# Configure Swarm Mode One Box
#
# This installs the following components
# - Docker
# - Docker Compose
# - Swarm Mode masters
# - Swarm Mode agents
###########################################################

set -x

echo "starting Swarm Mode cluster configuration"
date
ps ax

#############
# Parameters
#############

DOCKER_CE_VERSION=${1}
DOCKER_COMPOSE_VERSION=${2}
MASTERCOUNT=${3}
MASTERPREFIX=${4}
MASTERFIRSTADDR=${5}
AZUREUSER=${6}
POSTINSTALLSCRIPTURI=${7}
BASESUBNET=${8}
DOCKERENGINEDOWNLOADREPO=${9}
DOCKERCOMPOSEDOWNLOADURL=${10}
VMNAME=`hostname`
VMNUMBER=`echo $VMNAME | sed 's/.*[^0-9]\([0-9]\+\)*$/\1/'`
VMPREFIX=`echo $VMNAME | sed 's/\(.*[^0-9]\)*[0-9]\+$/\1/'`
OS="$(. /etc/os-release; echo $ID)"

echo "Master Count: $MASTERCOUNT"
echo "Master Prefix: $MASTERPREFIX"
echo "Master First Addr: $MASTERFIRSTADDR"
echo "vmname: $VMNAME"
echo "VMNUMBER: $VMNUMBER, VMPREFIX: $VMPREFIX"
echo "BASESUBNET: $BASESUBNET"
echo "AZUREUSER: $AZUREUSER"
echo "OS ID: $OS"

###################
# Common Functions
###################

isUbuntu()
{
  if [ "$OS" == "ubuntu" ]
  then
    return 0
  else
    return 1
  fi
}

isRHEL()
{
  if [ "$OS" == "rhel" ]
  then
    return 0
  else
    return 1
  fi
}

ensureAzureNetwork()
{
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
    echo "the network is not healthy, aborting install"
    ifconfig
    ip a
    exit 1
  fi
  # ensure the host ip can resolve
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
  # attempt to fix hostname, in case dns is not resolving Azure IPs (but can resolve public ips)
  if [ $networkHealthy -ne 0 ]
  then
    HOSTNAME=`hostname`
    HOSTADDR=`ip address show dev eth0 | grep -Eo 'inet (addr:)?([0-9]*\.){3}[0-9]*' | grep -Eo '([0-9]*\.){3}[0-9]*'`
    echo $HOSTADDR $HOSTNAME >> /etc/hosts
    hostname -i
    if [ $? -eq 0 ]
    then
      # hostname has been found continue
      networkHealthy=0
      echo "the network is healthy by updating /etc/hosts"
    fi
  fi
  if [ $networkHealthy -ne 0 ]
  then
    echo "the network is not healthy, cannot resolve ip address, aborting install"
    ifconfig
    ip a
    exit 2
  fi
}
ensureAzureNetwork
HOSTADDR=`hostname -i`

# apply all Canonical security updates during provisioning
/usr/lib/apt/apt.systemd.daily

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

MASTER0IPADDR="${BASESUBNET}${MASTERFIRSTADDR}"

######################
# resolve self in DNS
######################

if [ -z "$(grep "$HOSTADDR $VMNAME" /etc/hosts)" ]; then
    echo "$HOSTADDR $VMNAME" | sudo tee -a /etc/hosts
fi

################
# Install Docker
################

echo "Installing and configuring Docker"

installDockerUbuntu()
{
  for i in {1..10}; do
    apt-get install -y apt-transport-https ca-certificates curl software-properties-common
    curl --max-time 60 -fsSL https://download.docker.com/linux/ubuntu/gpg | apt-key add - 
    add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable"
    apt-get update
    apt-get install -y docker-ce=${DOCKER_CE_VERSION}
    if [ $? -eq 0 ]
    then
      systemctl restart docker
      # hostname has been found continue
      echo "Docker installed successfully"
      break
    fi
    sleep 10
  done
}

installDockerRHEL()
{
  for i in {1..10}; do
    yum-config-manager --add-repo https://download.docker.com/linux/centos/docker-ce.repo
    yum makecache fast
    yum -y install docker-ce
    if [ $? -eq 0 ]
    then
      systemctl enable docker
      systemctl start docker
      echo "Docker installed successfully"
      break
    fi
    sleep 10
  done
}

installDocker()
{
  if isUbuntu ; then
    installDockerUbuntu
  elif isRHEL ; then
    installDockerRHEL
  else
    echo "OS not supported, aborting install"
    exit 5
  fi
}

time installDocker

sudo usermod -aG docker $AZUREUSER

echo "Updating Docker daemon options"

updateDockerDaemonOptions()
{
    sudo mkdir -p /etc/systemd/system/docker.service.d
    # Start Docker and listen on :2375 (no auth, but in vnet) and
    # also have it bind to the unix socket at /var/run/docker.sock
    sudo bash -c 'echo "[Service]
    ExecStart=
    ExecStart=/usr/bin/dockerd -H tcp://0.0.0.0:2375 -H unix:///var/run/docker.sock
  " > /etc/systemd/system/docker.service.d/override.conf'
}
time updateDockerDaemonOptions

echo "Installing Docker Compose"
installDockerCompose()
{
  # sudo -i

  for i in {1..10}; do
    wget --tries 4 --retry-connrefused --waitretry=15 -qO- $DOCKERCOMPOSEDOWNLOADURL/$DOCKER_COMPOSE_VERSION/docker-compose-`uname -s`-`uname -m` > /usr/local/bin/docker-compose
    if [ $? -eq 0 ]
    then
      # hostname has been found continue
      echo "docker-compose installed successfully"
      break
    fi
    sleep 10
  done
}
time installDockerCompose
chmod +x /usr/local/bin/docker-compose

if ismaster && isRHEL ; then
  echo "Opening Docker ports"
  firewall-cmd --add-port=2375/tcp --permanent
  firewall-cmd --add-port=2377/tcp --permanent
  firewall-cmd --reload
fi

echo "Restarting Docker"
sudo systemctl daemon-reload
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

##############################################
# configure init rules restart all processes
##############################################

if ismaster ; then
    if [ "$HOSTADDR" = "$MASTER0IPADDR" ]; then
          echo "Creating a new Swarm on first master"
          docker swarm init --advertise-addr $(hostname -i):2377 --listen-addr $(hostname -i):2377
    else
        echo "Secondary master attempting to join an existing Swarm"
        swarmmodetoken=""
        swarmmodetokenAcquired=1
        for i in {1..120}; do
            swarmmodetoken=$(docker -H $MASTER0IPADDR:2375 swarm join-token -q manager)
            if [ $? -eq 0 ]; then
                swarmmodetokenAcquired=0
                break
            fi
            sleep 5
        done
        if [ $swarmmodetokenAcquired -ne 0 ]
        then
            echo "Secondary master couldn't connect to Swarm, aborting install"
            exit 3
        fi
        docker swarm join --token $swarmmodetoken $MASTER0IPADDR:2377
    fi
fi

if ismaster ; then
  echo "Having ssh listen to port 2222 as well as 22"
  sudo sed  -i "s/^Port 22$/Port 22\nPort 2222/1" /etc/ssh/sshd_config
fi

if ismaster ; then
  echo "Setting availability of master node: '$VMNAME' to pause"
  docker node update --availability pause $VMNAME
fi

if isagent ; then
    echo "Agent attempting to join an existing Swarm"
    swarmmodetoken=""
    swarmmodetokenAcquired=1
    for i in {1..120}; do
        swarmmodetoken=$(docker -H $MASTER0IPADDR:2375 swarm join-token -q worker)
        if [ $? -eq 0 ]; then
            swarmmodetokenAcquired=0
            break
        fi
        sleep 5
    done
    if [ $swarmmodetokenAcquired -ne 0 ]
    then
        echo "Agent couldn't join Swarm, aborting install"
        exit 4
    fi
    docker swarm join --token $swarmmodetoken $MASTER0IPADDR:2377
fi

if [ $POSTINSTALLSCRIPTURI != "disabled" ]
then
  echo "downloading, and kicking off post install script"
  /bin/bash -c "wget --tries 20 --retry-connrefused --waitretry=15 -qO- $POSTINSTALLSCRIPTURI | nohup /bin/bash >> /var/log/azure/cluster-bootstrap-postinstall.log 2>&1 &"
fi

# mitigation for bug https://bugs.launchpad.net/ubuntu/+source/linux/+bug/1676635
echo 2dd1ce17-079e-403c-b352-a1921ee207ee > /sys/bus/vmbus/drivers/hv_util/unbind
sed -i "13i\echo 2dd1ce17-079e-403c-b352-a1921ee207ee > /sys/bus/vmbus/drivers/hv_util/unbind\n" /etc/rc.local

echo "processes at end of script"
ps ax
date
echo "completed Swarm Mode cluster configuration"

echo "restart system to install any remaining software"
if isagent ; then
  shutdown -r now
else
  # wait 1 minute to restart master
  /bin/bash -c "shutdown -r 1 &"
fi
