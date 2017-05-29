package kubernetes

import (
	"encoding/json"
	"log"
	"os/exec"
	"time"
)

// ServiceList holds the items from a kubectl get svc
type ServiceList struct {
	Services []Service `json:"items"`
}

// Service represents a kubernetes service
type Service struct {
	ServiceMetadata ServiceMetadata `json:"metadata"`
	ServiceSpec     ServiceSpec     `json:"spec"`
}

// ServiceMetadata holds information like name, namespace, and labels
type ServiceMetadata struct {
	CreatedAt time.Time         `json:"creationTimestamp"`
	Labels    map[string]string `json:"labels"`
	Name      string            `json:"name"`
	Namespace string            `json:"namespace"`
}

// ServiceSpec holds information like clusterIP and port
type ServiceSpec struct {
	ClusterIP string `json:"clusterIP"`
	Ports     []Port `json:"ports"`
	Type      string `json:"type"`
}

// Port represents a service port definition
type Port struct {
	NodePort   int    `json:"nodePort"`
	Port       int    `json:"port"`
	Protocol   string `json:"protocol"`
	TargetPort int    `json:"targetPort"`
}

// GetServices returns the current service definition for a given namespace
func GetServices(namespace string) (*ServiceList, error) {
	out, err := exec.Command("kubectl", "get", "svc", "-o", "json", "-n", namespace).Output()
	if err != nil {
		log.Printf("Error trying to run 'kubectl get svc':%s\n", err)
		return nil, err
	}
	sl := ServiceList{}
	err = json.Unmarshal(out, &sl)
	if err != nil {
		log.Printf("Error unmarshalling service json:%s\n", err)
		return nil, err
	}
	return &sl, nil
}

// GetNodePortForPort will return the node port for a given pod
func (sl *ServiceList) GetNodePortForPort(service string, port int) int {
	for _, s := range sl.Services {
		if s.ServiceMetadata.Name == service {
			for _, p := range s.ServiceSpec.Ports {
				if p.Port == port {
					return p.NodePort
				}
			}
		}
	}
	return 0
}
