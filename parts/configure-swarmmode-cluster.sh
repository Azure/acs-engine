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

DOCKER_COMPOSE_VERSION="1.12.0"
#############
# Parameters
#############

MASTERCOUNT=${1}
MASTERPREFIX=${2}
MASTERFIRSTADDR=${3}
AZUREUSER=${4}
POSTINSTALLSCRIPTURI=${5}
BASESUBNET=${6}
VMNAME=`hostname`
VMNUMBER=`echo $VMNAME | sed 's/.*[^0-9]\([0-9]\+\)*$/\1/'`
VMPREFIX=`echo $VMNAME | sed 's/\(.*[^0-9]\)*[0-9]\+$/\1/'`

echo "Master Count: $MASTERCOUNT"
echo "Master Prefix: $MASTERPREFIX"
echo "Master First Addr: $MASTERFIRSTADDR"
echo "vmname: $VMNAME"
echo "VMNUMBER: $VMNUMBER, VMPREFIX: $VMPREFIX"
echo "BASESUBNET: $BASESUBNET"
echo "AZUREUSER: $AZUREUSER"

###################
# Common Functions
###################

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
    exit 2
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

echo "$HOSTADDR $VMNAME" | sudo tee -a /etc/hosts

################
# Install Docker
################

echo "Installing and configuring Docker"

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

echo "Updating Docker daemon options"

updateDockerDaemonOptions()
{
    sudo mkdir -p /etc/systemd/system/docker.service.d
    # Start Docker and listen on :2375 (no auth, but in vnet) and
    # also have it bind to the unix socket at /var/run/docker.sock
    sudo bash -c 'echo "[Service]
    ExecStart=
    ExecStart=/usr/bin/docker daemon -H tcp://0.0.0.0:2375 -H unix:///var/run/docker.sock
  " > /etc/systemd/system/docker.service.d/override.conf'
}
time updateDockerDaemonOptions

echo "Installing Docker Compose"
installDockerCompose()
{
  # sudo -i

  for i in {1..10}; do
    wget --tries 4 --retry-connrefused --waitretry=15 -qO- https://github.com/docker/compose/releases/download/$DOCKER_COMPOSE_VERSION/docker-compose-`uname -s`-`uname -m` > /usr/local/bin/docker-compose
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
            exit 2
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
        exit 2
    fi
    docker swarm join --token $swarmmodetoken $MASTER0IPADDR:2377
fi

if [ $POSTINSTALLSCRIPTURI != "disabled" ]
then
  echo "downloading, and kicking off post install script"
  /bin/bash -c "wget --tries 20 --retry-connrefused --waitretry=15 -qO- $POSTINSTALLSCRIPTURI | nohup /bin/bash >> /var/log/azure/cluster-bootstrap-postinstall.log 2>&1 &"
fi

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
