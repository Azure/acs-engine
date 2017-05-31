package kubernetes

import (
	"encoding/json"
	"log"
	"os/exec"
	"regexp"
	"time"
)

// PodList is a container that holds all pods returned from doing a kubectl get pods
type PodList struct {
	Pods []Pod `json:"items"`
}

// Pod is used to parse data from kubectl get pods
type Pod struct {
	PodMetadata PodMetadata `json:"metadata"`
	PodStatus   PodStatus   `json:"status"`
}

// PodMetadata holds information like name, createdat, labels, and namespace
type PodMetadata struct {
	CreatedAt time.Time         `json:"creationTimestamp"`
	Labels    map[string]string `json:"labels"`
	Name      string            `json:"name"`
	Namespace string            `json:"namespace"`
}

// PodStatus holds information like hostIP and phase
type PodStatus struct {
	HostIP    string    `json:"hostIP"`
	Phase     string    `json:"phase"`
	PodIP     string    `json:"podIP"`
	StartTime time.Time `json:"startTime"`
}

// AreAllPodsRunning will return true if all pods are in a Running State
func AreAllPodsRunning(podName, namespace string) (bool, error) {
	status := false
	out, err := exec.Command("kubectl", "get", "pods", "-n", namespace, "-o", "json").Output()
	if err != nil {
		log.Printf("Error trying to run 'kubectl version':%s\n", err)
		return status, err
	}
	pl := PodList{}
	err = json.Unmarshal(out, &pl)
	if err != nil {
		log.Printf("Error unmarshalling nodes json:%s\n", err)
		return status, err
	}

	for _, pod := range pl.Pods {
		matched, err := regexp.MatchString(podName+"-.*", pod.PodMetadata.Name)
		if err != nil {
			log.Printf("Error trying to match pod name:%s\n", err)
			return status, err
		}
		if matched {
			if pod.PodStatus.Phase == "Running" {
				status = true
			} else {
				status = false
			}
		}
	}

	return status, nil
}
