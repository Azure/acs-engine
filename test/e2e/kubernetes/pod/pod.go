package pod

import (
	"encoding/json"
	"log"
	"os/exec"
	"regexp"
	"time"
)

// List is a container that holds all pods returned from doing a kubectl get pods
type List struct {
	Pods []Pod `json:"items"`
}

// Pod is used to parse data from kubectl get pods
type Pod struct {
	Metadata Metadata `json:"metadata"`
	Status   Status   `json:"status"`
}

// Metadata holds information like name, createdat, labels, and namespace
type Metadata struct {
	CreatedAt time.Time         `json:"creationTimestamp"`
	Labels    map[string]string `json:"labels"`
	Name      string            `json:"name"`
	Namespace string            `json:"namespace"`
}

// Status holds information like hostIP and phase
type Status struct {
	HostIP    string    `json:"hostIP"`
	Phase     string    `json:"phase"`
	PodIP     string    `json:"podIP"`
	StartTime time.Time `json:"startTime"`
}

// AreAllPodsRunning will return true if all pods are in a Running State
func AreAllPodsRunning(podName, namespace string) (bool, error) {
	status := false
	out, err := exec.Command("kubectl", "get", "pods", "-n", namespace, "-o", "json").CombinedOutput()
	if err != nil {
		log.Printf("Error trying to run 'kubectl version':%s\n", string(out))
		return status, err
	}
	pl := List{}
	err = json.Unmarshal(out, &pl)
	if err != nil {
		log.Printf("Error unmarshalling nodes json:%s\n", err)
		return status, err
	}

	for _, pod := range pl.Pods {
		matched, err := regexp.MatchString(podName+"-.*", pod.Metadata.Name)
		if err != nil {
			log.Printf("Error trying to match pod name:%s\n", err)
			return status, err
		}
		if matched {
			if pod.Status.Phase == "Running" {
				status = true
			} else {
				status = false
			}
		}
	}

	return status, nil
}
