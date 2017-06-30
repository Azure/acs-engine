package report

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"sync"
	"time"
)

type TestFailure struct {
	Error string `json:"error"`
	Count int    `json:"count"`
}

type TestReport struct {
	BuildNum    int           `json:"build"`
	Deployments int           `json:"deployments"`
	Errors      int           `json:"errors"`
	StartTime   time.Time     `json:"startTime"`
	Duration    string        `json:"duration"`
	Failures    []TestFailure `json:"failures"`
}

type ReportManager struct {
	lock      sync.Mutex
	buildNum  int
	nDeploys  int
	nErrors   int
	timestamp time.Time
	failures  map[string]*TestFailure
}

var errorRegexpMap map[string]string

func init() {
	errorRegexpMap = map[string]string{
		"Azure CLI error ": "_init__.py",

		"Deployment error: VMStartTimedOut":                                "VMStartTimedOut",
		"Deployment error: OSProvisioningTimedOut":                         "OSProvisioningTimedOut",
		"Deployment error: VMExtensionProvisioningError":                   "VMExtensionProvisioningError",
		"Deployment error: VMExtensionProvisioningTimeout":                 "VMExtensionProvisioningTimeout",
		"Deployment error: InternalExecutionError":                         "InternalExecutionError",
		"Deployment error: SkuNotAvailable":                                "SkuNotAvailable",
		"Deployment error: MaxStorageAccountsCountPerSubscriptionExceeded": "MaxStorageAccountsCountPerSubscriptionExceeded",
		"Deployment error: ImageManagementOperationError":                  "ImageManagementOperationError",
		"Deployment error: DiskProcessingError":                            "DiskProcessingError",
		"Deployment error: DiskServiceInternalError":                       "DiskServiceInternalError",
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

func New(buildNum int, nDeploys int) *ReportManager {
	h := &ReportManager{}
	h.buildNum = buildNum
	h.nDeploys = nDeploys
	h.nErrors = 0
	h.timestamp = time.Now().UTC()
	h.failures = map[string]*TestFailure{}
	return h
}

func (h *ReportManager) Process(txt string) {
	for key, regex := range errorRegexpMap {
		if match, _ := regexp.MatchString(regex, txt); match {
			h.addFailure(key, 1)
			return
		}
	}
	h.addFailure("Unspecified error", 1)
}

func (h *ReportManager) addFailure(key string, n int) {
	h.lock.Lock()
	defer h.lock.Unlock()

	h.nErrors += n

	if testFailure, ok := h.failures[key]; !ok {
		h.failures[key] = &TestFailure{Error: key, Count: n}
	} else {
		testFailure.Count += n
	}
}

func (h *ReportManager) CreateTestReport(filepath string) error {
	testReport := &TestReport{}
	testReport.BuildNum = h.buildNum
	testReport.Deployments = h.nDeploys
	testReport.Errors = h.nErrors
	testReport.StartTime = h.timestamp
	testReport.Duration = time.Now().UTC().Sub(h.timestamp).String()
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
	_, err = file.Write(data)
	if err != nil {
		return err
	}
	return nil
}

func (h *ReportManager) CreateCombinedReport(filepath, testReportFname string) error {
	basedir := fmt.Sprintf("/var/lib/jenkins/jobs/%s/builds", os.Getenv("JOB_BASE_NAME"))
	if _, err := os.Stat(basedir); err != nil {
		return err
	}
	now := time.Now().UTC()
	combinedReport := New(h.buildNum, 0)
	for n := h.buildNum; n > 0; n-- {
		data, err := ioutil.ReadFile(fmt.Sprintf("%s/%d/%s", basedir, n, testReportFname))
		if err != nil {
			break
		}
		testReport := &TestReport{}
		if err := json.Unmarshal(data, &testReport); err != nil {
			break
		}
		// get combined report for past 24 hours
		if now.Sub(testReport.StartTime) > time.Duration(time.Hour*24) {
			break
		}
		combinedReport.timestamp = testReport.StartTime
		combinedReport.nDeploys += testReport.Deployments

		for _, f := range testReport.Failures {
			combinedReport.addFailure(f.Error, f.Count)
		}
	}
	return combinedReport.CreateTestReport(filepath)
}
