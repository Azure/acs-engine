#!/usr/bin/env bash
NGINX_INSTANCES=20
PODS_PER_NODE=50

while true; do
    echo Cleaning up Helm Releases...
    for RELEASE in `helm list | egrep -v 'UPDATED|static-release' | awk '{print $1}'`;
    do
        helm delete $RELEASE
    done
    echo Running activity scripts on a $1 node cluster
    echo Installing 20 instances of nginx-ingress
    echo And maybe installing a few instances of wordpress
    i=1
    while [[ $i -le $NGINX_INSTANCES ]]
    do
        rand=$RANDOM
        echo "###"
        echo Installing nginx
        echo "###"
        helm install stable/nginx-ingress
        let "rand %= 10"
        if (( $rand > 5 )); then
            echo "###"
            echo Installing wordpress
            echo "###"
            helm install stable/wordpress
        fi
        echo Taking a little break...
        let "rand %= 300"; sleep $rand
        ((i = i + 1))
    done
    PERCENT_COVERAGE=`bc <<<"scale=2; ${NGINX_INSTANCES} / ${1}"`
    REPLICAS=`bc <<<"scale=2; ${PODS_PER_NODE} / ${PERCENT_COVERAGE}"`
    REPLICAS_INT=${REPLICAS%.*}
    echo Scaling up nginx backend deployments to $REPLICAS_INT
    for deployment in `kubectl get deployments | egrep 'nginx-ingress-default-backend' | awk '{print $1}'`; do
        kubectl scale deployment $deployment --replicas=$REPLICAS_INT
    done
    SECONDS=0
    while true; do
        sleep 30
        kubectl get svc | egrep -q 'pending'
        RESPONSE=$?
        if [ $RESPONSE -eq 1 ]; then
            break
        else
            echo -ne "Waiting for LoadBalancer Ingress...\r"
        fi
    done
    MINUTES=$(($SECONDS % 60))
    echo -e "\e[32mTook $MINUTES minutes to get all public IP address assignments"
done