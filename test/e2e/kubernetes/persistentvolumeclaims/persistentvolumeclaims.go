package persistentvolumeclaims

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"time"
)

// PersistentVolumeClaims is used to parse data from kubectl get pvc
type PersistentVolumeClaims struct {
	Metadata Metadata `json:"metadata"`
	Spec     Spec     `json:"spec"`
	Status   Status   `json:"status"`
}

// Metadata holds information like name, create time, and namespace
type Metadata struct {
	CreatedAt time.Time `json:"creationTimestamp"`
	Name      string    `json:"name"`
	NameSpace string    `json:"namespace"`
}

// Spec holds information like storageClassName, volumeName
type Spec struct {
	StorageClassName string `json:"storageClassName"`
	VolumeName       string `json:"volumeName"`
}

// Status holds information like phase
type Status struct {
	Phase string `json:"phase"`
}

// CreatePersistentVolumeClaimsFromFile will create a StorageClass from file with a name
func CreatePersistentVolumeClaimsFromFile(filename, name, namespace string) (*PersistentVolumeClaims, error) {
	out, err := exec.Command("kubectl", "apply", "-f", filename).CombinedOutput()
	if err != nil {
		log.Printf("Error trying to create PersistentVolumeClaims %s in namespace %s:%s\n", name, namespace, string(out))
		return nil, err
	}
	pvc, err := Get(name, namespace)
	if err != nil {
		log.Printf("Error while trying to fetch PersistentVolumeClaims %s in namespace %s:%s\n", name, namespace, err)
		return nil, err
	}
	return pvc, nil
}

// Get will return a PersistentVolumeClaims with a given name and namespace
func Get(pvcName, namespace string) (*PersistentVolumeClaims, error) {
	out, err := exec.Command("kubectl", "get", "pvc", pvcName, "-n", namespace, "-o", "json").CombinedOutput()
	if err != nil {
		return nil, err
	}
	pvc := PersistentVolumeClaims{}
	err = json.Unmarshal(out, &pvc)
	if err != nil {
		log.Printf("Error unmarshalling PersistentVolumeClaims json:%s\n", err)
		return nil, err
	}
	return &pvc, nil
}

// WaitOnReady will block until PersistentVolumeClaims is available
func (pvc *PersistentVolumeClaims) WaitOnReady(namespace string, sleep, duration time.Duration) (bool, error) {
	readyCh := make(chan bool, 1)
	errCh := make(chan error)
	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()
	go func() {
		for {
			select {
			case <-ctx.Done():
				errCh <- fmt.Errorf("Timeout exceeded (%s) while waiting for PersistentVolumeClaims (%s) to become ready", duration.String(), pvc.Metadata.Name)
			default:
				query, _ := Get(pvc.Metadata.Name, namespace)
				if query != nil && query.Status.Phase == "Bound" {
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
