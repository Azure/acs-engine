package node

import (
	"errors"
	"fmt"
	"log"
	"os/exec"
	"regexp"
	"strings"
)

// Version returns the version of an OpenShift cluster.
func Version() (string, error) {
	cmd := exec.Command("oc", "version")
	fmt.Printf("\n$ %s\n", strings.Join(cmd.Args, " "))
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Error trying to run 'oc version':%s", string(out))
		return "", err
	}
	exp := regexp.MustCompile(`(openshift\s)+(v\d+.\d+.\d+)+`)
	for _, line := range strings.Split(string(out), "\n") {
		if strings.HasPrefix(line, "openshift") {
			s := exp.FindStringSubmatch(line)
			return s[2], nil
		}
	}
	return "", errors.New("cannot find openshift version")
}
