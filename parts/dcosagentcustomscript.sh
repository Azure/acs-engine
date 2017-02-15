#!/bin/bash

set -e

update_attribute()
{
    if [ ! -a /etc/mesosphere/roles/master ];
    then
        echo adding FD / UD attributes
        meta=$( curl --silent --fail  http://169.254.169.254/metadata/v1/InstanceInfo )
        # parsing with bash to avoid jq install. This might break when the metadata service changes
        ud=$(echo $meta | sed "s/^.*\"UD\": *\"\([0-9]*\)\".*$/\1/")
        fd=$(echo $meta | sed "s/^.*\"FD\": *\"\([0-9]*\)\".*$/\1/")

        if [ -a /var/lib/dcos/mesos-slave-common ];
        then
            echo Adding $ud $fd
            sed -e "s/;*$/;UD:"$ud";FD:"$fd"/" -i /var/lib/dcos/mesos-slave-common
        else
            echo new attributes
            echo "MESOS_ATTRIBUTES=UD:$ud;FD:$fd" >> /var/lib/dcos/mesos-slave-common
        fi
    fi
}

create_raid()
{
    raidname="md127"
    raiddev="/dev/"$raidname
    mpoint="/dcos/volume0"
    disks=( $( ls /dev/sd[cdef] ) )
    if [ ${#disks[@]} -eq 0 ]
    then
        echo no dataDisks found. Exiting
        exit 0
    fi

    echo Found ${#disks[@]} Disks
    echo found ${disks[@]}

    if [ ! -a $raiddev ]; then
        echo setting up raid
        # install mdadm
        until apt-get -y update && apt-get -y install mdadm
        do
            echo "Trying again"
            sleep 2
        done

        # partition
        for d in "${disks[@]}"
        do
            parted --script  $d mklabel gpt
            parted --script --  $d mkpart primary 0% 100%
            #parted --script --  $d set 1 raid on
        done

        # raid
        partitions=( $( ls /dev/sd[cdef]1 ) )
        mdadm --create $raidname --level=0 --raid-devices=${#partitions[@]} ${partitions[@]}
        mkfs.ext4 $raiddev

        echo mounting on $mpoint
        mkdir -p $mpoint

        echo UUID=$(sudo /sbin/blkid | grep $raiddev | cut -d\" -f 2) $mpoint ext4 defaults 0 2 >> /etc/fstab
        mount -a
        chmod 777 $mpoint
    fi
}

update_attribute
create_raid