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
        String sendTo = "${SEND_TO}".trim()

        dir(clone_dir) {
          def img = null
          stage('Init') {
            deleteDir()
            checkout scm
            img = docker.build('acs-engine-test', '--pull .')
          }

          img.inside("-u root:root") {
            def success = true
            def junit_dir = "_junit"
            def prefix = ""
            stage('Setup') {
              // Set up Azure
              sh("az login --service-principal -u ${SPN_USER} -p ${SPN_PASSWORD} --tenant ${TENANT_ID}")
              sh("az account set --subscription ${SUBSCRIPTION_ID}")

              env.SERVICE_PRINCIPAL_CLIENT_ID="${SPN_USER}"
              env.SERVICE_PRINCIPAL_CLIENT_SECRET="${SPN_PASSWORD}"
              env.TENANT_ID="${TENANT_ID}"
              env.SUBSCRIPTION_ID="${SUBSCRIPTION_ID}"
              env.LOCATION = "${LOCATION}"
              env.CLUSTER_SERVICE_PRINCIPAL_CLIENT_ID="${SERVICE_PRINCIPAL_CLIENT_ID}"
              env.CLUSTER_SERVICE_PRINCIPAL_CLIENT_SECRET="${SERVICE_PRINCIPAL_CLIENT_SECRET}"
              env.CLEANUP = "y"

              sh("printf 'acs-features-test%x' \$(date '+%s') > INSTANCE_NAME_PREFIX")
              prefix = readFile('INSTANCE_NAME_PREFIX').trim()
              // Create report directory
              sh("mkdir ${junit_dir}")
              // Build and test acs-engine
              sh('make ci')
            }
            def subdirs = "${SCENARIOS}".tokenize('[ \t\n]+')
            for (i = 0; i < subdirs.size(); i++) {
              def subdir = subdirs[i]
              def names = sh(returnStdout: true, script: "cd examples; ls ${subdir}/*.json").split("\\r?\\n")
              for(j = 0; j< names.size(); j++) {
                def name = names[j].trim()
                env.CLUSTER_DEFINITION = pwd()+"/examples/${name}"
                env.INSTANCE_NAME = "${prefix}-${i}-${j}"
                env.RESOURCE_GROUP = "test-acs-${env.LOCATION}-${env.BUILD_NUMBER}-${i}-${j}"
                env.DEPLOYMENT_NAME = "${env.RESOURCE_GROUP}"
                env.ORCHESTRATOR = sh(returnStdout: true, script: "awk '/\\\"orchestratorType\\\"\\s*:/ { print \$2 }' ${env.CLUSTER_DEFINITION} | awk -F\\\" '{print \$2}'").toLowerCase().trim()
                echo env.ORCHESTRATOR
                // Generate and deploy template, validate deployments
                try {
                  stage(name) {
                    def scripts = ["generate_template.sh", "deploy_template.sh"]
                    if(env.ORCHESTRATOR == "dcos" || env.ORCHESTRATOR == "swarmmode" || env.ORCHESTRATOR == "kubernetes") {
                      scripts += "validate_deployment.sh"
                    }
                    for(k = 0; k < scripts.size(); k++) {
                      def script = scripts[k]
                      def test = "${name}.${script}"
                      sh("mkdir -p ${junit_dir}/${test}")
                      sh("cp ./test/shunit/${script} ${junit_dir}/${test}/t.sh")
                      sh("cd ${junit_dir}; shunit.sh -t ${test} > ${test}/junit.xml")
                      sh("grep 'failures=\"0\"' ${junit_dir}/${test}/junit.xml")
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
              } // for (j = 0; j <files...
            } // for (i = 0; i <subdirs...
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
            sh("rm -rf ${clone_dir}/.kube")
            sh("rm -rf ${junit_dir}")
          }
        }
      }
    }
  }
}
