# marathonBaseUrl="http://localhost:8080/marathon" # handy for local testing with SSH
marathonBaseUrl="http://marathon.mesos:8080" 
# dcosBaseUrl="http://localhost/" # handy for local testing with SSH
dcosBaseUrl="http://leader.mesos/"

# definitions for elasticsearch and vamp applications (from http://vamp.io/documentation/installation/v0.9.5/dcos/)
elasticsearch='{"id":"elasticsearch","instances":1,"cpus":0.2,"mem":1024.0,"container":{"type":"DOCKER","docker":{"image":"magneticio/elastic:2.2","network":"HOST","forcePullImage":true}},"healthChecks":[{"protocol":"TCP","gracePeriodSeconds":30,"intervalSeconds":10,"timeoutSeconds":5,"port":9200,"maxConsecutiveFailures":0}]}'
vamp='{"id":"vamp/vamp","instances":1,"cpus":0.5,"mem":1024,"container":{"type":"DOCKER","docker":{"image":"magneticio/vamp:0.9.5-dcos","network":"BRIDGE","portMappings":[{"containerPort":8080,"hostPort":0,"name":"vip0","labels":{"VIP_0":"10.20.0.100:8080"}}],"forcePullImage":true}},"labels":{"DCOS_SERVICE_NAME":"vamp","DCOS_SERVICE_SCHEME":"http","DCOS_SERVICE_PORT_INDEX":"0"},"env":{"VAMP_WAIT_FOR":"http://elasticsearch.marathon.mesos:9200/.kibana","VAMP_WORKFLOW_DRIVER_VAMP_URL":"http://10.20.0.100:8080","VAMP_ELASTICSEARCH_URL":"http://elasticsearch.marathon.mesos:9200","VAMP_DB_TYPE":"mysql","VAMP_DB_URL":"jdbc:mysql://mysql.marathon.l4lb.thisdcos.directory:3306/vamp-${namespace}?useSSL=false","VAMP_DB_CREATE_URL":"jdbc:mysql://mysql.marathon.l4lb.thisdcos.directory:3306?useSSL=false","VAMP_DB_USER":"root","VAMP_DB_PASSWORD":"root"},"healthChecks":[{"protocol":"TCP","gracePeriodSeconds":30,"intervalSeconds":10,"timeoutSeconds":5,"portIndex":0,"maxConsecutiveFailures":0}]}'

function wait_for_running_app(){
    ## args
    ## $1 - app name
    appname=$1

    ## wait for deployment
    waitCount=0
    appRunning=0
    while [ $waitCount -lt 12 ]
    do
        echo "Testing for $appname running..."
        appResponse=$(curl "$marathonBaseUrl/v2/apps/$appname" 2>/dev/null) # url for $appname app
        if [[ $appResponse =~ \"state\":\"TASK_RUNNING\" ]]; then # test for TASK_RUNNING
            echo "$appname running!"
            appRunning=1
            return 0
        fi
        
        sleep 10s
        waitCount=$[$waitCount+1]
    done
    return 1
}

## Install mysql from universe
echo "Starting mysql package deployment"
curl "$dcosBaseUrl/package/install" -H 'Accept: application/vnd.dcos.package.install-response+json;charset=utf-8;version=v1' -H 'Content-Type: application/vnd.dcos.package.install-request+json;charset=utf-8;version=v1' -X POST -d '{"packageName": "mysql"}'
echo ""

wait_for_running_app "mysql"
mysqlRunning=$?

if [[ $mysqlRunning -ne 0 ]]; then
    echo "mysql not running - quitting"
    exit 100
fi

## Start elasticsearch deployment
echo "Starting elasticsearch deployment"
curl -X POST "$marathonBaseUrl/v2/apps" -H 'Content-Type: application/json' -d $elasticsearch
echo ""

wait_for_running_app "elasticsearch"
elasticRunning=$?

if [[ $elasticRunning -ne 0 ]]; then
    echo "elasticsearch not running - quitting"
    exit 100
fi


## Start vamp deployment
echo "Starting vamp deployment"
curl -X POST "$marathonBaseUrl/v2/apps" -H 'Content-Type: application/json' -d $vamp
echo ""

wait_for_running_app "vamp/vamp"
running=$?
if [[ $running -ne 0 ]]; then
    echo "vamp/vamp not running - quitting"
    exit 100
fi

wait_for_running_app "vamp/vamp-gateway-agent"
running=$?
if [[ $running -ne 0 ]]; then
    echo "vamp/vamp-gateway-agent not running - quitting"
    exit 100
fi

wait_for_running_app "vamp/workflow-health"
running=$?
if [[ $running -ne 0 ]]; then
    echo "vamp/workflow-health not running - quitting"
    exit 100
fi

wait_for_running_app "vamp/workflow-kibana"
running=$?
if [[ $running -ne 0 ]]; then
    echo "vamp/workflow-kibana not running - quitting"
    exit 100
fi

wait_for_running_app "vamp/workflow-metrics"
running=$?
if [[ $running -ne 0 ]]; then
    echo "vamp/workflow-metrics not running - quitting"
    exit 100
fi

wait_for_running_app "vamp/workflow-vga"
running=$?
if [[ $running -ne 0 ]]; then
    echo "vamp/workflow-vga not running - quitting"
    exit 100
fi
