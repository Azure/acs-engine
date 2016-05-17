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
DOCKER_COMPOSE_VERSION="1.6.2"
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

consulstr()
{
  consulargs=""
  for i in `seq 0 $((MASTERCOUNT-1))` ;
  do
    MASTEROCTET=`expr $MASTERFIRSTADDR + $i`
    IPADDR="${BASESUBNET}${MASTEROCTET}"

    if [ "$VMNUMBER" -eq "0" ]
    then
      consulargs="${consulargs}-bootstrap-expect $MASTERCOUNT "
    fi
    if [ "$VMNUMBER" -eq "$i" ]
    then
      consulargs="${consulargs}-advertise $IPADDR "
    else
      consulargs="${consulargs}-retry-join $IPADDR "
    fi
  done
  echo $consulargs
}

consulargs=$(consulstr)
MASTER0IPADDR="${BASESUBNET}${MASTERFIRSTADDR}"

######################
# resolve self in DNS
######################

echo "$HOSTADDR $VMNAME" | sudo tee -a /etc/hosts

################
# Install Docker
################

echo "Installing and configuring docker"

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
  echo 'DOCKER_OPTS="-H unix:///var/run/docker.sock -H 0.0.0.0:2375 --cluster-store=consul://'$MASTER0IPADDR:8500 --cluster-advertise=$HOSTADDR:2375'"' | sudo tee -a /etc/default/docker
fi

echo "Installing docker compose"
installDockerCompose()
{
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
  mkdir -p /data/consul
  echo "consul:
  image: \"progrium/consul\"
  command: -server -node $VMNAME $consulargs
  ports:
    - \"8500:8500\"
    - \"8300:8300\"
    - \"8301:8301\"
    - \"8301:8301/udp\"
    - \"8302:8302\"
    - \"8302:8302/udp\"
    - \"8400:8400\"
  volumes:
    - \"/data/consul:/data\"
  restart: \"always\"
swarm:
  image: \"$SWARM_VERSION\"
  command: manage --replication --advertise $HOSTADDR:2375 --discovery-opt kv.path=docker/nodes consul://$MASTER0IPADDR:8500
  ports:
    - \"2375:2375\"
  links:
    - \"consul\"
  volumes:
    - \"/etc/docker:/etc/docker\"
  restart: \"always\"
" > /opt/azure/containers/docker-compose.yml

  pushd /opt/azure/containers/
  docker-compose up -d
  popd
  echo "completed starting docker swarm on the master"
fi

if [ $POSTINSTALLSCRIPTURI != "disabled" ]
then
  echo "downloading, and kicking off post install script"
  /bin/bash -c "wget --tries 20 --retry-connrefused --waitretry=15 -qO- $POSTINSTALLSCRIPTURI | nohup /bin/bash >> /var/log/azure/cluster-bootstrap-postinstall.log 2>&1 &"
fi

echo "processes at end of script"
ps ax
date
echo "completed mesos cluster configuration"

echo "restart system to install any remaining software"
if isagent ; then
  shutdown -r now
else
  # wait 1 minute to restart master
  /bin/bash -c "shutdown -r 1 &"
fi
