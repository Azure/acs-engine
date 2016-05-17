# DCOS Templates for a 3 master clusters with attached disks

This template deploys a DCOS cluster on Azure an attached disk.

# Approach to enable larger persistent data disks on machines

Instructions below contributed by @Radek44

Disclaimer: this approach is not a replacement for the current DCOS given as it has a few tradeoffs however a potential alternative approach for customers that need large amounts of persistent storage within the cluster for running.

Step 1 : Deploy the 3 cluster DCOS template in this project.

Step 2: In the template it is possible to configure additional Data Disks . Maximum size of a single data disk is 1023Gb, depending on VM class up to 32 disks can be attached to a single VM.
"dataDisks" : [
     {
          "name": ....
           "diskSizeGb": "1023" ...
     }
]

Step 3: Once the cluster is created the disks should be configured before anything else is deployed on the agents.
In order to prepare the disks and register them permanently here is a quick bash script
```#!/bin/bash
# Create a folder for the dcos first extended volume
# Volumes have to be made discoverable by mesos by being created as /dcos/volume{n}
# Create a partition table and one partition spawning the entire disk
# Format the partition as ext4
# Mount the partition
# Register the partition in fstab
# Reboot the machine after 1 minute
#
sudo mkdir -p /dcos/volume0&&sudo parted -s /dev/sdc mklabel gpt mkpart primary ext4 0% 100%&&sudo mkfs -t ext4 /dev/sdc1&&sudo mount /dev/sdc1 /dcos/volume0&&sudo sh -c "echo '/dev/sdc1\t/dcos/volume0\text4\tdefaults\t0\t2' >> /etc/fstab"&&sudo shutdown --reboot 1
Notice that the Azure documentation suggests that volumes are registered using their UUID. However doing that I have had failures with agents restarting
In cases anyone still wants to try this is the alternate approach for the registration step becomes:
sudo sh -c "echo 'UUID=$(sudo blkid | grep '/dev/sdc1' | sed -n 's/.*UUID=\"\([^\"]*\)\".*/\1/p')\t/dcos/volume0\text4\tdefaults\t0\t2' >> /etc/fstab"
```
Step 4: Once the machine rebooted, edit the DCOS configuration in order to start leveraging the new volume.
```#!/bin/bash
#
# This script will take a disk previsouly mounted on /dcos/volume0 and ensure that it is visible and used by mesos
# It removes current mesos resource configuration
# It then restarts relevant dcos services to have them pick up the config
# Then quick reboot and we are golden
sudo rm /var/lib/dcos/mesos-resources&&sudo rm -f /var/lib/mesos/slave/meta/slaves/latest&&sudo systemctl restart dcos-vol-discovery-priv-agent&&sleep 5&&sudo systemctl restart dcos-mesos-slave&&sudo shutdown --reboot 1
```
Step 5: Once the machine rebooted Mesos should now be reporting ~1Gb extra for each node on which this was executed.
A quick helper script to execute the above in parallel on all target agents (requires a private_agents file listing the target IPs).
```#!/bin/bash
#
# Assumes you have a file called private_agents that contains the list of hosts that you want impacted by the script (typically your private agent pool IPs so 10.32.0.x range
# Install parallel-ssh
sudo apt-get install pssh
#
# First part of the parallel executions creates and mounts the disks and restarts the agent
parallel-ssh -O StrictHostKeyChecking=no -l loopadmin -h private_agents -P -I < ~/datadisk.sh
#
# Get a coffee - wait a bit for agent to restart
sleep 2m
#
# Second part of the parallel execution configures mesos to see the new disks and restarts
parallel-ssh -O StrictHostKeyChecking=no -l loopadmin -h private_agents -P -I < ~/mesosdisk.sh
#
echo "When You Play The Game Of Thrones, You Win Or You Die"
```

For further reading on attached disks browse to https://azure.microsoft.com/en-us/documentation/articles/virtual-machines-linux-classic-attach-disk/.
