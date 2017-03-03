#!/usr/bin/env groovy

node {
  withCredentials([[$class: 'UsernamePasswordMultiBinding', credentialsId: 'AZURE_CLI_SPN_ACS_TEST',
                  passwordVariable: 'SPN_PASSWORD', usernameVariable: 'SPN_USER']]) {
    timestamps {
      wrap([$class: 'AnsiColorBuildWrapper', 'colorMapName': 'XTerm']) {

        def tests = [:]
        def orch = ["kubernetes", "dcos", "swarm", "swarmmode"]
        def jobname = "acs-engine-deploy"
        for(int i = 0; i < orch.size(); i++) {
          def orchestrator = orch[i]
          echo orchestrator
          tests[orchestrator] = {
            stage(orchestrator) {
              build job: jobname,
              parameters:
               [[$class: 'StringParameterValue', name: 'ORCHESTRATOR', value: orchestrator],
                [$class: 'StringParameterValue', name: 'LOCATION',     value: 'eastus']]
            }
          }
        }
        parallel tests
      }
    }
  }
}
