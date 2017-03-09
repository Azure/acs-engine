#!/bin/bash

set -e

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
        until apt-get -y update && apt-get -y install mdadm
        do
            echo "Trying again"
            sleep 2
        done
        for d in "${disks[@]}"
        do
            parted --script  $d mklabel gpt
            parted --script --  $d mkpart primary 0% 100%
        done
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

create_raid