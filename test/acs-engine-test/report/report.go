package report

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"sync"
)

type TestFailure struct {
	Error string `json:"error"`
	Count int    `json:"count"`
}

type TestReport struct {
	Build    string        `json:"build"`
	Failures []TestFailure `json:"failures"`
}

type ReportManager struct {
	lock     sync.Mutex
	build    string
	failures map[string]*TestFailure
}

var errorRegexpMap map[string]string

func init() {
	errorRegexpMap = map[string]string{
		"Azure CLI error ": "_init__.py",

		"Deployment error: VMExtensionProvisioningTimeout":                 "VMExtensionProvisioningTimeout",
		"Deployment error: InternalExecutionError":                         "InternalExecutionError",
		"Deployment error: OSProvisioningTimedOut":                         "OSProvisioningTimedOut",
		"Deployment error: VMStartTimedOut":                                "VMStartTimedOut",
		"Deployment error: DiskServiceInternalError":                       "DiskServiceInternalError",
		"Deployment error: VMExtensionProvisioningError":                   "VMExtensionProvisioningError",
		"Deployment error: MaxStorageAccountsCountPerSubscriptionExceeded": "MaxStorageAccountsCountPerSubscriptionExceeded",
		"Deployment error: ImageManagementOperationError":                  "ImageManagementOperationError",
		"Deployment error: DiskProcessingError":                            "DiskProcessingError",
		"Deployment error: AllocationFailed":                               "AllocationFailed",
		"Deployment error: NetworkingInternalOperationError":               "NetworkingInternalOperationError",

		"K8S validattion: curl error":                 "curl_error",
		"K8S validation: external IP":                 "gave up waiting for loadbalancer to get an ingress ip",
		"K8S validation: nodes not ready":             "gave up waiting for apiserver",
		"K8S validation: service unreachable":         "gave up waiting for service to be externally reachable",
		"K8S validation: nginx unreachable":           "failed to get expected response from nginx through the loadbalancer",
		"DCOS validation: nodes not ready":            "gave up waiting for DCOS nodes",
		"DCOS validation: marathon validation failed": "dcos/test.sh] marathon validation failed",
		"DCOS validation: marathon not added":         "dcos/test.sh] gave up waiting for marathon to be added",
		"DCOS validation: marathon-lb not installed":  "Failed to install marathon-lb",
	}
}

func New(build string) *ReportManager {
	h := &ReportManager{}
	h.build = build
	h.failures = map[string]*TestFailure{}
	return h
}

func (h *ReportManager) Process(txt string) {
	fmt.Println("ReportManager.Process")
	for key, regex := range errorRegexpMap {
		if match, _ := regexp.MatchString(regex, txt); match {
			h.addFailure(key)
			return
		}
	}
	h.addFailure("Unspecified error")
}

func (h *ReportManager) addFailure(key string) {
	fmt.Println("ReportManager.addFailure")
	h.lock.Lock()
	defer h.lock.Unlock()

	if testFailure, ok := h.failures[key]; !ok {
		h.failures[key] = &TestFailure{Error: key, Count: 1}
	} else {
		testFailure.Count++
	}
}

func (h *ReportManager) CreateReport(filepath string) error {
	fmt.Println("ReportManager.CreateReport")
	testReport := &TestReport{}
	testReport.Build = h.build
	testReport.Failures = make([]TestFailure, len(h.failures))
	i := 0
	for _, f := range h.failures {
		testReport.Failures[i] = *f
		i++
	}
	data, err := json.Marshal(testReport)
	if err != nil {
		return err
	}
	file, err := os.OpenFile(filepath, os.O_CREATE|os.O_WRONLY, os.FileMode(0644))
	if err != nil {
		return err
	}
	defer file.Close()
	return file.Write(data)
}
