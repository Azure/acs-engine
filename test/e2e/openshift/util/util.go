package util

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
)

func printCmd(cmd *exec.Cmd) {
	fmt.Printf("\n$ %s\n", strings.Join(cmd.Args, " "))
}

// CreateFromTemplate processes and creates the provided templateName/templateNamespace template
// in the provided namespace.
func CreateFromTemplate(templateName, templateNamespace, namespace string) error {
	processCmd := exec.Command("oc", "process", templateName, "-n", templateNamespace)
	printCmd(processCmd)
	out, err := processCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("cannot process template %s: %v\noutput: %s", templateName, err, string(out))
	}
	if err := ioutil.WriteFile(templateName, out, 0644); err != nil {
		return fmt.Errorf("cannot create tempfile for processed template %s: %v", templateName, err)
	}
	defer os.Remove(templateName)
	createCmd := exec.Command("oc", "create", "-n", namespace, "-f", templateName)
	printCmd(createCmd)
	out, err = createCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("cannot create processed template %s: %v\noutput: %s", templateName, err, string(out))
	}
	return nil
}

// WaitForDeploymentConfig waits until the provided deploymentconfig namespace/name
// gets deployed.
func WaitForDeploymentConfig(name, namespace string) error {
	cmd := exec.Command("oc", "rollout", "status", fmt.Sprintf("dc/%s", name), "-n", namespace)
	printCmd(cmd)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Error trying to see the rollout status of dc/%s: %s", name, string(out))
		return err
	}
	return nil
}

// GetHost expects the name and namespace of a route in order to
// return its host.
func GetHost(name, namespace string) (string, error) {
	cmd := exec.Command("oc", "get", fmt.Sprintf("route/%s", name), "-n", namespace, "-o", "jsonpath={.spec.host}")
	printCmd(cmd)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Error trying to get the hostname of route/%s: %s", name, string(out))
		return "", err
	}
	return string(out), nil
}

// TestHost tries to access host and retries maxRetries times with a retryDelay
// that is doubled on every retry.
func TestHost(host string, maxRetries int, retryDelay time.Duration) error {
	backoff := retryDelay
	url := fmt.Sprintf("http://%s", host)

	resp, err := http.Get(url)
	if err == nil && resp.StatusCode == http.StatusOK {
		return nil
	}
	if err == nil {
		log.Printf("got status %q while trying to access %s", resp.Status, host)
		resp.Body.Close()
	} else {
		log.Printf("error while trying to access %s: %v", host, err)
	}
	for retries := 1; retries <= maxRetries; retries++ {
		log.Printf("Retry #%d to access %s", retries, host)
		resp, err = http.Get(url)
		if err != nil {
			log.Printf("error while trying to access %s: %v", host, err)
			continue
		}
		resp.Body.Close()
		if resp.StatusCode == http.StatusOK {
			return nil
		}
		log.Printf("got status %q while trying to access %s", resp.Status, host)
		time.Sleep(backoff)
		backoff *= 2
	}
	if err != nil {
		return err
	}
	return fmt.Errorf("unexpected response status: %v", resp.Status)
}
