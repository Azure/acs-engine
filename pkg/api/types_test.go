package api

import (
	"log"
	"testing"
)

const exampleCustomHyperkubeImage = `example.azurecr.io/example/hyperkube-amd64:custom`

const exampleAPIModel = `{
		"apiVersion": "vlabs",
	"properties": {
		"orchestratorProfile": {
			"orchestratorType": "Kubernetes",
			"kubernetesConfig": {
				"customHyperkubeImage": "` + exampleCustomHyperkubeImage + `"
			}
		},
		"masterProfile": { "count": 1, "dnsPrefix": "", "vmSize": "Standard_D2_v2" },
		"agentPoolProfiles": [ { "name": "linuxpool1", "count": 2, "vmSize": "Standard_D2_v2", "availabilityProfile": "AvailabilitySet" } ],
		"windowsProfile": { "adminUsername": "azureuser", "adminPassword": "replacepassword1234$" },
		"linuxProfile": { "adminUsername": "azureuser", "ssh": { "publicKeys": [ { "keyData": "" } ] }
		},
		"servicePrincipalProfile": { "clientId": "", "secret": "" }
	}
}
`

func TestIsDCOS(t *testing.T) {
	dCOSProfile := &OrchestratorProfile{
		OrchestratorType: "DCOS",
	}
	if !dCOSProfile.IsDCOS() {
		t.Fatalf("unable to detect DCOS orchestrator profile from OrchestratorType=%s", dCOSProfile.OrchestratorType)
	}
	kubernetesProfile := &OrchestratorProfile{
		OrchestratorType: "Kubernetes",
	}
	if kubernetesProfile.IsDCOS() {
		t.Fatalf("unexpectedly detected DCOS orchestrator profile from OrchestratorType=%s", kubernetesProfile.OrchestratorType)
	}
}

func TestCustomHyperkubeImageField(t *testing.T) {
	log.Println(exampleAPIModel)
	apiloader := &Apiloader{
		Translator: nil,
	}
	apimodel, _, err := apiloader.DeserializeContainerService([]byte(exampleAPIModel), false, false, nil)
	if err != nil {
		t.Fatalf("unexpectedly error deserializing the example apimodel: %s", err)
	}

	actualCustomHyperkubeImage := apimodel.Properties.OrchestratorProfile.KubernetesConfig.CustomHyperkubeImage
	if actualCustomHyperkubeImage != exampleCustomHyperkubeImage {
		t.Fatalf("kubernetesConfig->customHyperkubeImage field value was unexpected: got(%s), expected(%s)", actualCustomHyperkubeImage, exampleCustomHyperkubeImage)
	}
}
