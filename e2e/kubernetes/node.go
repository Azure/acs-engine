package kubernetes

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
	CreatedAt  time.Time         `json:"creationTimestamp"`
	Labels     map[string]string `json:"labels"`
	Name       string            `json:"name"`
	NodeStatus NodeStatus        `json:"status"`
}

// NodeStatus parses information from the status key
type NodeStatus struct {
	NodeInfo      NodeInfo  `json:"nodeInfo"`
	NodeAddresses []Address `json:"addresses"`
}

// Address contains an address and a type
type Address struct {
	Address string `json:"address"`
	Type    string `json:"type"`
}

// NodeInfo contains information like what version the kubelet is running
type NodeInfo struct {
	ContainerRuntimeVersion string `json:"containerRuntimeVersion"`
	KubeProxyVersion        string `json:"kubeProxyVersion"`
	KubeletProxyVersion     string `json:"kubeletVersion"`
	OperatingSystem         string `json:"operatingSystem"`
}

// NodeList is used to parse out Nodes from a list
type NodeList struct {
	Nodes []Node `json:"items"`
}

// GetNodes returns the current nodes for a given kubeconfig
func GetNodes() (*NodeList, error) {
	out, err := exec.Command("kubectl", "get", "nodes", "-o", "json").Output()
	if err != nil {
		log.Printf("Error trying to run 'kubectl get nodes':%s\n", err)
		return nil, err
	}
	nl := NodeList{}
	err = json.Unmarshal(out, &nl)
	if err != nil {
		log.Printf("Error unmarshalling nodes json:%s\n", err)
	}
	return &nl, nil
}

// GetVersion get the version of the server
func GetVersion() (string, error) {
	out, err := exec.Command("kubectl", "version", "--short").Output()
	if err != nil {
		log.Printf("Error trying to run 'kubectl version':%s\n", err)
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
func (ns *NodeStatus) GetAddressByType(t string) *Address {
	for _, a := range ns.NodeAddresses {
		if a.Type == t {
			return &a
		}
	}
	return nil
}
