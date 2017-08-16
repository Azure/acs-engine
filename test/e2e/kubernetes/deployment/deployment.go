package deployment

import (
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"time"
)

// List holds a list of deployments returned from kubectl get deploy
type List struct {
	Deployments []Deployment `json:"items"`
}

// Deployment repesentes a kubernetes deployment
type Deployment struct {
	Metadata Metadata `json:"metadata"`
}

// Metadata holds information like labels, name, and namespace
type Metadata struct {
	CreatedAt time.Time         `json:"creationTimestamp"`
	Labels    map[string]string `json:"labels"`
	Name      string            `json:"name"`
	Namespace string            `json:"namespace"`
}

// Spec holds information the deployment strategy and number of replicas
type Spec struct {
	Replicas int      `json:"replicas"`
	Template Template `json:"template"`
}

// Template is used for fetching the deployment spec -> containers
type Template struct {
	TemplateSpec TemplateSpec `json:"spec"`
}

// TemplateSpec holds the list of containers for a deployment, the dns policy, and restart policy
type TemplateSpec struct {
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

// Create will create a deployment for a given image with a name in a namespace
func Create(image, name, namespace string) (*Deployment, error) {
	out, err := exec.Command("kubectl", "run", "-n", namespace, "--image", "library/nginx:latest", name).CombinedOutput()
	if err != nil {
		log.Printf("Error trying to deploy %s [%s] in namespace %s:%s\n", name, image, namespace, string(out))
		return nil, err
	}
	d, err := Get(name, namespace)
	if err != nil {
		log.Printf("Error while trying to fetch Deployment %s in namespace %s:%s\n", name, namespace, err)
		return nil, err
	}
	return d, nil
}

// Get returns a deployment from a name and namespace
func Get(name, namespace string) (*Deployment, error) {
	out, err := exec.Command("kubectl", "get", "deploy", "-o", "json", "-n", namespace, name).CombinedOutput()
	if err != nil {
		log.Printf("Error while trying to fetch deployment %s in namespace %s:%s\n", name, namespace, string(out))
		return nil, err
	}
	d := Deployment{}
	err = json.Unmarshal(out, &d)
	if err != nil {
		log.Printf("Error while trying to unmarshal deployment json:%s\n%s\n", err, string(out))
		return nil, err
	}
	return &d, nil
}

// Delete will delete a deployment in a given namespace
func (d *Deployment) Delete() error {
	out, err := exec.Command("kubectl", "delete", "deploy", "-n", d.Metadata.Namespace, d.Metadata.Name).CombinedOutput()
	if err != nil {
		log.Printf("Error while trying to delete deployment %s in namespace %s:%s\n", d.Metadata.Namespace, d.Metadata.Name, string(out))
		return err
	}
	return nil
}

// Expose will create a load balancer and expose the deployment on a given port
func (d *Deployment) Expose(port int) error {
	ref := fmt.Sprintf("deployments/%s", d.Metadata.Name)
	out, err := exec.Command("kubectl", "expose", ref, "--type", "LoadBalancer", "-n", d.Metadata.Namespace, "--port", strconv.Itoa(port)).CombinedOutput()
	if err != nil {
		log.Printf("Error while trying to expose deployment %s in namespace %s on port %v:%s\n", d.Metadata.Name, d.Metadata.Namespace, port, string(out))
		return err
	}
	return nil
}
