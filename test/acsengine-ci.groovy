#!/usr/bin/env groovy

node {
  withCredentials([[$class: 'UsernamePasswordMultiBinding', credentialsId: 'AZURE_CLI_SPN_ACS_TEST',
                  passwordVariable: 'SPN_PASSWORD', usernameVariable: 'SPN_USER']]) {
    timestamps {
      wrap([$class: 'AnsiColorBuildWrapper', 'colorMapName': 'XTerm']) {
        env.GOPATH="${WORKSPACE}"
        env.PATH="${env.PATH}:${env.GOPATH}/bin"
        def clone_dir = "${env.GOPATH}/src/github.com/Azure/acs-engine"
        env.HOME=clone_dir
        def success = true

        dir(clone_dir) {
          def img = null
          try {
            stage('Init') {
              deleteDir()
              checkout scm
              img = docker.build('acs-engine-test', '--pull .')
            }
          }
          catch(exc) {
            echo "Exception ${exc}"
            success = false
          }
          img.inside("-u root:root") {
            String errorMsg = ""
            try {
              stage('Test') {
                if(success) {
                  // Create template, deploy and test
                  env.SERVICE_PRINCIPAL_CLIENT_ID="${SPN_USER}"
                  env.SERVICE_PRINCIPAL_CLIENT_SECRET="${SPN_PASSWORD}"
                  env.TENANT_ID="${TENANT_ID}"
                  env.SUBSCRIPTION_ID="${SUBSCRIPTION_ID}"
                  
                  sh("printf 'acs-ci-test%x' \$(date '+%s') > INSTANCE_NAME")
                  env.INSTANCE_NAME = readFile('INSTANCE_NAME').trim()
                  env.INSTANCE_NAME_PREFIX = "acs-ci"
                  env.ORCHESTRATOR = "${ORCHESTRATOR}"
                  env.CLUSTER_DEFINITION="examples/${ORCHESTRATOR}.json"
                  env.CLUSTER_SERVICE_PRINCIPAL_CLIENT_ID="${SERVICE_PRINCIPAL_CLIENT_ID}"
                  env.CLUSTER_SERVICE_PRINCIPAL_CLIENT_SECRET="${SERVICE_PRINCIPAL_CLIENT_SECRET}"

                  env.LOCATION = "${LOCATION}"
                  env.RESOURCE_GROUP = "test-acs-${ORCHESTRATOR}-${env.LOCATION}-${env.BUILD_NUMBER}"
                  env.DEPLOYMENT_NAME = "${env.RESOURCE_GROUP}"
                  env.CLEANUP = "y"

                  sh('./test/deploy.sh')
                }
              }
            }
            catch(exc) {
              echo "Exception ${exc}"
              success = false
              errorMsg = "Please run \"make ci\" for verification"
            }
            // Final clean up
            sh("rm -rf ${clone_dir}/_output")
            sh("rm -rf ${clone_dir}/.azure")
            sh("rm -rf ${clone_dir}/.kube")
            if(!success) {
              currentBuild.result = "FAILURE"
              String to = "${SEND_TO}".trim()
              if(errorMsg != "") {
                if(to != "") {
                  to += ";"
                }
                to += emailextrecipients([[$class: 'CulpritsRecipientProvider']])
              }
              if(to != "") {
                def url = "${env.BUILD_URL}\n\n"
                for(String addr : to.tokenize('[ \t\n;,]+')) {
                  if(!addr.endsWith("@microsoft.com")) {
                    url = ""
                  }
                }
                gitCommit = sh(returnStdout: true, script: 'git rev-parse HEAD').trim()
                emailext(
                  to: to,
                  subject: "[ACS Engine is BROKEN] ${env.JOB_NAME} #${env.BUILD_NUMBER}",
                  body: "Commit: ${gitCommit}\n\n${url}${errorMsg}"
                )
              }
            }
          }
        }
      }
    }
  }
}
