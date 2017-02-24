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
        env.ORCHESTRATOR = "${ORCHESTRATOR}"
        String locations_str = "${LOCATIONS}"
        String sendTo = "${SEND_TO}".trim()

        if(locations_str.equals("all")) {
          locations_str = "\
australiaeast australiasoutheast \
brazilsouth \
canadacentral canadaeast \
centralindia southindia \
centralus eastus2 eastus northcentralus southcentralus westcentralus westus2 westus \
eastasia southeastasia \
japaneast japanwest \
northeurope westeurope \
uksouth ukwest"
        }
        def locations = locations_str.tokenize('[ \t\n]+')

        dir(clone_dir) {
          def img = null
          stage('Init') {
            deleteDir()
            checkout scm
            img = docker.build('acs-engine-test', '--pull .')
          }

          img.inside("-u root:root") {
            def junit_dir = "_junit"
            stage('Setup') {
              // Set up Azure
              sh("az login --service-principal -u ${SPN_USER} -p ${SPN_PASSWORD} --tenant ${TENANT_ID}")
              sh("az account set --subscription ${SUBSCRIPTION_ID}")
              // Create report directory
              sh("mkdir ${junit_dir}")
              // Build and test acs-engine
              sh('make ci')
              // Create template
              sh("printf 'acs-test%x' \$(date '+%s') > INSTANCE_NAME")
              env.INSTANCE_NAME = readFile('INSTANCE_NAME').trim()
              env.CLUSTER_DEFINITION="examples/${ORCHESTRATOR}.json"
              env.CLUSTER_SERVICE_PRINCIPAL_CLIENT_ID="${SERVICE_PRINCIPAL_CLIENT_ID}"
              env.CLUSTER_SERVICE_PRINCIPAL_CLIENT_SECRET="${SERVICE_PRINCIPAL_CLIENT_SECRET}"
              sh('./test/step.sh generate_template')
            }

            for (i = 0; i <locations.size(); i++) {
              env.LOCATION = locations[i]
              env.RESOURCE_GROUP = "test-acs-${ORCHESTRATOR}-${env.LOCATION}-${env.BUILD_NUMBER}"
              env.DEPLOYMENT_NAME = "${env.RESOURCE_GROUP}"
              def ok = true
              // Deploy
              try {
                stage("${env.LOCATION} deploy") {
                  def test = "deploy-${env.LOCATION}"
                  sh("mkdir -p ${junit_dir}/${test}")
                  sh("cp ./test/shunit/deploy_template.sh ${junit_dir}/${test}/t.sh")
                  sh("cd ${junit_dir}; shunit.sh -t ${test} > ${test}/junit.xml")
                  sh("grep 'failures=\"0\"' ${junit_dir}/${test}/junit.xml")
                }
              }
              catch(exc) {
                echo "Exception ${exc}"
                ok = false
              }
              // Verify deployment
              try {
                stage("${env.LOCATION} validate") {
                  if(ok) {
                    def test = "validate-${env.LOCATION}"
                    sh("mkdir -p ${junit_dir}/${test}")
                    sh("cp ./test/shunit/validate_deployment.sh ${junit_dir}/${test}/t.sh")
                    sh("cd ${junit_dir}; shunit.sh -t ${test} > ${test}/junit.xml")
                    sh("grep 'failures=\"0\"' ${junit_dir}/${test}/junit.xml")
                  }
                  else {
                    echo "Skipped verification for ${env.RESOURCE_GROUP}"
                  }
                }
              }
              catch(exc) {
                echo "Exception ${exc}"
              }
              // Clean up
              try {
                sh('./test/step.sh cleanup')
              }
              catch(exc) {
                echo "Exception ${exc}"
              }
            } // for (i = 0; i <locations...
            // Generate reports
            try {
              junit("${junit_dir}/**/junit.xml")
              if(currentBuild.result == "UNSTABLE") {
                currentBuild.result = "FAILURE"
                if(sendTo != "") {
                  emailext(
                    to: "${sendTo}",
                    subject: "[ACS Engine Jenkins Failure] ${env.JOB_NAME} #${env.BUILD_NUMBER}",
                    body: "${env.BUILD_URL}testReport")
                }
              }
            }
            catch(exc) {
              echo "Exception ${exc}"
            }
            // Final clean up
            sh("rm -rf ${clone_dir}/_output")
            sh("rm -rf ${clone_dir}/.azure")
            sh("rm -rf ${junit_dir}")
          }
        }
      }
    }
  }
}
