package namespace

import (
	"encoding/json"
	"log"
	"os/exec"
	"time"
)

// Namespace holds namespace metadata
type Namespace struct {
	Metadata Metadata `json:"metadata"`
}

// Metadata holds information like name and created timestamp
type Metadata struct {
	CreatedAt time.Time `json:"creationTimestamp"`
	Name      string    `json:"name"`
}

// Create a namespace with the given name
func Create(name string) (*Namespace, error) {
	out, err := exec.Command("kubectl", "create", "namespace", name).CombinedOutput()
	if err != nil {
		log.Printf("Error trying to create namespace (%s):%s\n", name, string(out))
		return nil, err
	}
	return Get(name)
}

// Get returns a namespace for with a given name
func Get(name string) (*Namespace, error) {
	out, err := exec.Command("kubectl", "get", "namespace", name, "-o", "json").CombinedOutput()
	if err != nil {
		log.Printf("Error trying to get namespace (%s):%s\n", name, string(out))
		return nil, err
	}
	n := Namespace{}
	err = json.Unmarshal(out, &n)
	if err != nil {
		log.Printf("Error unmarshalling namespace json:%s\n", err)
	}
	return &n, nil
}

// Delete a namespace
func (n *Namespace) Delete() error {
	out, err := exec.Command("kubectl", "delete", "namespace", n.Metadata.Name).CombinedOutput()
	if err != nil {
		log.Printf("Error while trying to delete namespace (%s):%s\n", n.Metadata.Name, out)
		return err
	}
	return nil
}
