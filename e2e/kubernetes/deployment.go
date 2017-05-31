package kubernetes

import (
	"encoding/json"
	"log"
	"os/exec"
	"time"
)

// DeploymentList holds a list of deployments returned from kubectl get deploy
type DeploymentList struct {
	Deployments []Deployment `json:"items"`
}

// Deployment repesentes a kubernetes deployment
type Deployment struct {
	DeploymentMetadata DeploymentMetadata `json:"metadata"`
}

// DeploymentMetadata holds information like labels, name, and namespace
type DeploymentMetadata struct {
	CreatedAt time.Time         `json:"creationTimestamp"`
	Labels    map[string]string `json:"labels"`
	Name      string            `json:"name"`
	Namespace string            `json:"namespace"`
}

// DeploymentSpec holds information the deployment strategy and number of replicas
type DeploymentSpec struct {
	Replicas           int                `json:"replicas"`
	DeploymentTemplate DeploymentTemplate `json:"template"`
}

// DeploymentTemplate is used for fetching the deployment spec -> containers
type DeploymentTemplate struct {
	DeploymentTemplateSpec DeploymentTemplateSpec `json:"spec"`
}

// DeploymentTemplateSpec holds the list of containers for a deployment, the dns policy, and restart policy
type DeploymentTemplateSpec struct {
	Containers    []Container `json:"containers"`
	DNSPolicy     string      `json:"dnsPolicy"`
	RestartPolicy string      `json:"restartPolicy"`
}

// Container holds information like image, pull policy, name, etc...
type Container struct {
	Image      string `json:"image"`
	PullPolicy string `json:"imagePullPolicy"`
	Name       string `json:"name"`
}

// Deploy will create a deployment for a given image with a name in a namespace
func Deploy(image, name, namespace string) (*Deployment, error) {
	_, err := exec.Command("kubectl", "run", "--image", "library/nginx:latest", "nginx", "--namespace", "deis").Output()
	if err != nil {
		log.Printf("Error trying to deploy %s [%s] in namespace %s:%s\n", name, image, namespace, err)
		return nil, err
	}
	d, err := GetDeployment(name, namespace)
	if err != nil {
		log.Printf("Error while trying to fetch Deployment %s in namespace %s:%s\n", name, namespace, err)
		return nil, err
	}
	return d, nil
}

// GetDeployment returns a deployment from a name and namespace
func GetDeployment(name, namespace string) (*Deployment, error) {
	out, err := exec.Command("kubectl", "get", "deploy", "-n", "namespace", "name").Output()
	if err != nil {
		log.Printf("Error while trying to fetch deployment %s in namespace %s:%s\n", name, namespace, err)
		return nil, err
	}
	d := Deployment{}
	err = json.Unmarshal(out, &d)
	if err != nil {
		log.Printf("Error while trying to unmarshal deployment json:%s\n", err)
		return nil, err
	}
	return &d, nil
}
