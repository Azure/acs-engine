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

// GetAllByPrefix will return all pods in a given namespace that match a prefix
func GetAllByPrefix(prefix, namespace string) ([]Pod, error) {
	pl, err := GetAll(namespace)
	if err != nil {
		return nil, err
	}
	pods := []Pod{}
	for _, p := range pl.Pods {
		matched, err := regexp.MatchString(prefix+"-.*", p.Metadata.Name)
		if err != nil {
			log.Printf("Error trying to match pod name:%s\n", err)
			return nil, err
		}
		if matched {
			pods = append(pods, p)
		}
	}
	return pods, nil
}

// AreAllPodsRunning will return true if all pods are in a Running State
func AreAllPodsRunning(podPrefix, namespace string) (bool, error) {
	pl, err := GetAll(namespace)
	if err != nil {
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
				ready, _ := AreAllPodsRunning(podPrefix, namespace)
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

// Exec will execute the given command in the pod
func (p *Pod) Exec(cmd ...string) ([]byte, error) {
	execCmd := []string{"exec", p.Metadata.Name, "-n", p.Metadata.Namespace}
	for _, s := range cmd {
		execCmd = append(execCmd, s)
	}
	out, err := exec.Command("kubectl", execCmd...).CombinedOutput()
	if err != nil {
		log.Printf("Error trying to run 'kubectl exec':%s\n", string(out))
		log.Printf("Command:kubectl exec %s -n %s %s \n", p.Metadata.Name, p.Metadata.Namespace, cmd)
		return nil, err
	}
	return out, nil
}

// CheckLinuxOutboundConnection will keep retrying the check if an error is received until the timeout occurs or it passes. This helps us when DNS may not be available for some time after a pod starts.
func (p *Pod) CheckLinuxOutboundConnection(sleep, duration time.Duration) (bool, error) {
	exp, err := regexp.Compile("200 OK")
	if err != nil {
		log.Printf("Error while trying to create regex for linux outbound check:%s\n", err)
		return false, err
	}
	readyCh := make(chan bool, 1)
	errCh := make(chan error)
	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()
	go func() {
		for {
			select {
			case <-ctx.Done():
				errCh <- fmt.Errorf("Timeout exceeded (%s) while waiting for Pod (%s) to check outbound internet connection", duration.String(), p.Metadata.Name)
			default:
				_, err := p.Exec("--", "/usr/bin/apt", "update")
				if err != nil {
					break
				}
				_, err = p.Exec("--", "/usr/bin/apt", "install", "-y", "curl")
				if err != nil {
					break
				}
				out, err := p.Exec("--", "curl", "-I", "http://www.bing.com")
				if err == nil {
					matched := exp.MatchString(string(out))
					if matched {
						readyCh <- true
					} else {
						readyCh <- false
					}
				} else {
					log.Printf("Error:%s\n", err)
					log.Printf("Out:%s\n", out)
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

// CheckWindowsOutboundConnection will keep retrying the check if an error is received until the timeout occurs or it passes. This helps us when DNS may not be available for some time after a pod starts.
func (p *Pod) CheckWindowsOutboundConnection(sleep, duration time.Duration) (bool, error) {
	exp, err := regexp.Compile("(StatusCode\\s*:\\s*200)")
	if err != nil {
		log.Printf("Error while trying to create regex for windows outbound check:%s\n", err)
		return false, err
	}
	readyCh := make(chan bool, 1)
	errCh := make(chan error)
	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()
	go func() {
		for {
			select {
			case <-ctx.Done():
				errCh <- fmt.Errorf("Timeout exceeded (%s) while waiting for Pod (%s) to check outbound internet connection", duration.String(), p.Metadata.Name)
			default:
				out, err := p.Exec("--", "powershell", "iwr", "-UseBasicParsing", "-TimeoutSec", "60", "www.bing.com")
				if err == nil {
					matched := exp.MatchString(string(out))
					if matched {
						readyCh <- true
					} else {
						readyCh <- false
					}
				} else {
					log.Printf("Error:%s\n", err)
					log.Printf("Out:%s\n", out)
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
