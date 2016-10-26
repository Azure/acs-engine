#!/bin/bash

set -e

update_attribute()
{
    if [ ! -a /etc/mesosphere/roles/master ];
    then
        echo adding FD / UD attributes
        meta=$( curl http://169.254.169.254/metadata/v1/InstanceInfo )
        # parsing with bash to avoid jq install. This might break when the metadata service changes
        ud="UD"$( echo $meta | cut -d\" -f 8)
        fd="FD"$( echo $meta | cut -d\" -f 12)

        if [ -a /var/lib/dcos/mesos-slave-common ];
        then
            echo Adding $ud $fd
            sed -e 's/;*$/;UD:'$ud';FD:'$fd'/' -i /var/lib/dcos/mesos-slave-common
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

    if [ ! -a /dev/sd[cdef] ]
    then
        echo no dataDisks found. Exiting
        exit 1
    fi

    disks=( $( ls /dev/sd[cdef] ) )
    echo Found ${#disks[@]} Disks
    echo found ${disks[@]}

    if [ ! -e $( $raiddev ) ]; then
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
            sudo parted $d --script mklabel gpt
            sudo parted $d --script -- mkpart primary 0 -1
            sudo parted $d --script -- set 1 raid on
        done

        # raid
        mdadm --create $raidname --level 0 --raid-devices ${#disks[@]} ${disks[@]}
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
