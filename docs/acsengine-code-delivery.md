# ACS Engine code delivery guide

[![ACS Engine](https://azurecomcdn.azureedge.net/mediahandler/acomblog/media/Default/blog/a8f28783-3ddc-4081-a57d-6d97147467bf.png)](https://github.com/azure/acs-engine)

ACS Engine is an open source project to generate ARM (Azure Resource Manager) templates DC/OS, Kubernetes, Swarm Mode clusters on Microsoft Azure.
This documents provides guidelines to the acs-engine testing and continuous integration process.

## Development pipeline

ACS Engine employs CI system that incorporates Jenkins server, configured to interact with ACS Engine GitHub project.
A recommended way to contribute to ACS Engine is to fork github.com/Azure/acs-engine project.
and create a separated branch (a feature branch) for the feature you are working on.

The following steps constitute ACS Engine delivery pipeline

 1. Complete the current iteration of the code change, and check it into the feature branch
 2. Invoke unit test. Return to step (1) if failed.
```sh
    $ make ci
```
 3. Create a template. Return to step (1) if failed.
```sh
    $ acs-engine generate --api-model kubernetes.json
```
 4. Deploy the template in Azure. Return to step (1) if failed.
```sh
    $ az group create --name=<RESOURCE_GROUP> --location=<LOCATION>
    $ az group deployment create \
      --name <DEPLOYMENT_NAME> \
      --resource-group <RESOURCE_GROUP> \
     --template-file azuredeploy.json \
     --parameters @azuredeploy.parameters.json
```
  5. Create a pull request (PR) from github.com/Azure/acs-engine portal.
  6. The PR triggers a Jenkins job that 
  + applies the changes to the HEAD of the master branch
  + generates multiple ARM templates for different deployment scenarios
  + simultaneously provisions the clusters based on generated templates in Azure
  This test might take 20-40 minutes.
  If the test fails, review the logs. If the failure was caused by your code change, return to step (1).
  Sometimes the test might fail because of intermittent Azure issues, such as resource unavailability of provisioning timeout. In this case manually trigger Jenkins PR job from your GitHub PR page.
  7. The PR is reviewed by the members of ACS Engine team. Should the changes have been requested, return to step (1).
  8. Once the PR is approved, and Jenkins PR job has passed, the PR could be merged into the master branch
  9. Once merged, another Jenkins job is triggered, to verify integrity of the master branch. This job is similar to the PR job.
