package pod

import (
	"context"
	"encoding/json"
	"fmt"
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

// GetAll will return all pods in a given namespace
func GetAll(namespace string) (*List, error) {
	out, err := exec.Command("kubectl", "get", "pods", "-n", namespace, "-o", "json").CombinedOutput()
	if err != nil {
		log.Printf("Error trying to run 'kubectl get pods':%s\n", string(out))
		return nil, err
	}
	pl := List{}
	err = json.Unmarshal(out, &pl)
	if err != nil {
		log.Printf("Error unmarshalling nodes json:%s\n", err)
		return nil, err
	}
	return &pl, nil
}

// Get will return a pod with a given name and namespace
func Get(podName, namespace string) (*Pod, error) {
	out, err := exec.Command("kubectl", "get", "pods", podName, "-n", namespace, "-o", "json").CombinedOutput()
	if err != nil {
		log.Printf("Error trying to run 'kubectl get pods':%s\n", string(out))
		return nil, err
	}
	p := Pod{}
	err = json.Unmarshal(out, &p)
	if err != nil {
		log.Printf("Error unmarshalling nodes json:%s\n", err)
		return nil, err
	}
	return &p, nil
}

// AreAllPodsRunning will return true if all pods are in a Running State
func AreAllPodsRunning(podPrefix, namespace string) (bool, error) {
	pl, err := GetAll(namespace)
	if err != nil {
		log.Printf("Error while trying to check if all pods are in running state:%s", err)
		return false, err
	}

	var status []bool
	for _, pod := range pl.Pods {
		matched, err := regexp.MatchString(podPrefix+"-.*", pod.Metadata.Name)
		if err != nil {
			log.Printf("Error trying to match pod name:%s\n", err)
			return false, err
		}
		if matched {
			if pod.Status.Phase != "Running" {
				status = append(status, false)
			} else {
				status = append(status, true)
			}
		}
	}

	if len(status) == 0 {
		return false, nil
	}

	for _, s := range status {
		if s == false {
			return false, nil
		}
	}

	return true, nil
}

// WaitOnReady will block until all nodes are in ready state
func WaitOnReady(podPrefix, namespace string, sleep, duration time.Duration) (bool, error) {
	readyCh := make(chan bool, 1)
	errCh := make(chan error)
	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()
	go func() {
		for {
			select {
			case <-ctx.Done():
				errCh <- fmt.Errorf("Timeout exceeded (%s) while waiting for Pods (%s) to become ready in namespace (%s)", duration.String(), podPrefix, namespace)
			default:
				ready, err := AreAllPodsRunning(podPrefix, namespace)
				if err != nil {
					log.Printf("Error while waiting on pods to become ready:%s", err)
				}
				if ready == true {
					readyCh <- true
				} else {
					time.Sleep(sleep)
				}
			}
		}
	}()
	for {
		select {
		case err := <-errCh:
			return false, err
		case ready := <-readyCh:
			return ready, nil
		}
	}
}
