#!/bin/bash

set -x

echo "starting swarm cluster configuration"
date
ps ax

#############
# Parameters
#############

SWARM_VERSION=${1}
DOCKER_COMPOSE_VERSION=${2}
MASTERCOUNT=${3}
MASTERPREFIX=${4}
MASTERFIRSTADDR=${5}
AZUREUSER=${6}
POSTINSTALLSCRIPTURI=${7}
BASESUBNET=${8}
DOCKERENGINEDOWNLOADREPO=${9}
DOCKERCOMPOSEDOWNLOADURL=${10}
DOCKER_CE_VERSION=17.03.*
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

# simple general command retry function
retrycmd_if_failure() { for i in 1 2 3 4 5; do $@; [ $? -eq 0  ] && break || sleep 5; done ; }

installDocker()
{
  for i in {1..10}; do
    apt-get install -y apt-transport-https ca-certificates curl software-properties-common
    curl --max-time 60 -fsSL https://download.docker.com/linux/ubuntu/gpg | apt-key add - 
    add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable"
    apt-get update
    apt-get install -y docker-ce=${DOCKER_CE_VERSION}
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

if ismaster ; then
  echo "Having ssh listen to port 2222 as well as 22"
  sudo sed  -i "s/^Port 22$/Port 22\nPort 2222/1" /etc/ssh/sshd_config
fi

if [ $POSTINSTALLSCRIPTURI != "disabled" ]
then
  echo "downloading, and kicking off post install script"
  /bin/bash -c "wget --tries 20 --retry-connrefused --waitretry=15 -qO- $POSTINSTALLSCRIPTURI | nohup /bin/bash >> /var/log/azure/cluster-bootstrap-postinstall.log 2>&1 &"
fi

echo "processes at end of script"
ps ax
date
echo "completed Swarm cluster configuration"

echo "restart system to install any remaining software"
if isagent ; then
  shutdown -r now
else
  # wait 1 minute to restart master
  /bin/bash -c "shutdown -r 1 &"
fi
