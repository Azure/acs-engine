package node

import (
	"encoding/json"
	"log"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

const (
	//ServerVersion is used to parse out the version of the API running
	ServerVersion = `(Server Version:\s)+(v\d.\d.\d)+`
)

// Node represents the kubernetes Node Resource
type Node struct {
	CreatedAt time.Time         `json:"creationTimestamp"`
	Labels    map[string]string `json:"labels"`
	Name      string            `json:"name"`
	Status    Status            `json:"status"`
}

// Status parses information from the status key
type Status struct {
	Info          Info      `json:"Info"`
	NodeAddresses []Address `json:"addresses"`
}

// Address contains an address and a type
type Address struct {
	Address string `json:"address"`
	Type    string `json:"type"`
}

// Info contains information like what version the kubelet is running
type Info struct {
	ContainerRuntimeVersion string `json:"containerRuntimeVersion"`
	KubeProxyVersion        string `json:"kubeProxyVersion"`
	KubeletProxyVersion     string `json:"kubeletVersion"`
	OperatingSystem         string `json:"operatingSystem"`
}

// List is used to parse out Nodes from a list
type List struct {
	Nodes []Node `json:"items"`
}

// Get returns the current nodes for a given kubeconfig
func Get() (*List, error) {
	out, err := exec.Command("kubectl", "get", "nodes", "-o", "json").CombinedOutput()
	if err != nil {
		log.Printf("Error trying to run 'kubectl get nodes':%s\n", string(out))
		return nil, err
	}
	nl := List{}
	err = json.Unmarshal(out, &nl)
	if err != nil {
		log.Printf("Error unmarshalling nodes json:%s\n", err)
	}
	return &nl, nil
}

// Version get the version of the server
func Version() (string, error) {
	out, err := exec.Command("kubectl", "version", "--short").CombinedOutput()
	if err != nil {
		log.Printf("Error trying to run 'kubectl version':%s\n", string(out))
		return "", err
	}
	split := strings.Split(string(out), "\n")
	exp, err := regexp.Compile(ServerVersion)
	if err != nil {
		log.Printf("Error while compiling regexp:%s\n", ServerVersion)
	}
	s := exp.FindStringSubmatch(split[1])
	return s[2], nil
}

// GetAddressByType will return the Address object for a given Kubernetes node
func (ns *Status) GetAddressByType(t string) *Address {
	for _, a := range ns.NodeAddresses {
		if a.Type == t {
			return &a
		}
	}
	return nil
}
