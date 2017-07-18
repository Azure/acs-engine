package report

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"sync"
	"time"
)

type ErrorInfo struct {
	TestName string
	ErrName  string
	ErrClass string
	Location string
}

type ErrorStat struct {
	Count     int            `json:"count"`
	Locations map[string]int `json:"locations"`
}

type ReportMgr struct {
	lock        sync.Mutex
	JobName     string    `json:"job"`
	BuildNum    int       `json:"build"`
	Deployments int       `json:"deployments"`
	Errors      int       `json:"errors"`
	StartTime   time.Time `json:"startTime"`
	Duration    string    `json:"duration"`
	// Failure map: key=error, value=locations
	Failures map[string]*ErrorStat `json:"failures"`
}

type logError struct {
	name  string
	class string
	regex string
}

const (
	errClassDeployment = "Deployment"
	errClassValidation = "Validation"
	errClassAzcli      = "AzCLI"
	errClassNone       = "None"

	errSuccess = "Success"
	errUnknown = "UnspecifiedError"
)

var logErrors []logError

func init() {
	logErrors = []logError{
		{name: "AzCliRunError", class: errClassAzcli, regex: "_init__.py"},
		{name: "AzCliLoadError", class: errClassAzcli, regex: "Error loading command module"},

		{name: "VMStartTimedOut", class: errClassDeployment, regex: "VMStartTimedOut"},
		{name: "OSProvisioningTimedOut", class: errClassDeployment, regex: "OSProvisioningTimedOut"},
		{name: "VMExtensionProvisioningError", class: errClassDeployment, regex: "VMExtensionProvisioningError"},
		{name: "VMExtensionProvisioningTimeout", class: errClassDeployment, regex: "VMExtensionProvisioningTimeout"},
		{name: "InternalExecutionError", class: errClassDeployment, regex: "InternalExecutionError"},
		{name: "SkuNotAvailable", class: errClassDeployment, regex: "SkuNotAvailable"},
		{name: "MaxStorageAccountsCountPerSubscriptionExceeded", class: errClassDeployment, regex: "MaxStorageAccountsCountPerSubscriptionExceeded"},
		{name: "ImageManagementOperationError", class: errClassDeployment, regex: "ImageManagementOperationError"},
		{name: "DiskProcessingError", class: errClassDeployment, regex: "DiskProcessingError"},
		{name: "DiskServiceInternalError", class: errClassDeployment, regex: "DiskServiceInternalError"},
		{name: "AllocationFailed", class: errClassDeployment, regex: "AllocationFailed"},
		{name: "NetworkingInternalOperationError", class: errClassDeployment, regex: "NetworkingInternalOperationError"},
		{name: "PlatformFaultDomainCount", class: errClassDeployment, regex: "platformFaultDomainCount"},

		{name: "K8sNodeNotReady", class: errClassValidation, regex: "K8S: gave up waiting for apiserver"},
		{name: "K8sUnexpectedVersion", class: errClassValidation, regex: "K8S: unexpected kubernetes version"},
		{name: "K8sContainerNotCreated", class: errClassValidation, regex: "K8S: gave up waiting for containers"},
		{name: "K8sPodNotRunning", class: errClassValidation, regex: "K8S: gave up waiting for running pods"},
		{name: "K8sKubeDnsNotRunning", class: errClassValidation, regex: "K8S: gave up waiting for kube-dns"},
		{name: "K8sDashboardNotRunning", class: errClassValidation, regex: "K8S: gave up waiting for kubernetes-dashboard"},
		{name: "K8sKubeProxyNotRunning", class: errClassValidation, regex: "K8S: gave up waiting for kube-proxy"},
		{name: "K8sProxyNotWorking", class: errClassValidation, regex: "K8S: gave up verifying proxy"},
		{name: "K8sDeploymentNotReady", class: errClassValidation, regex: "K8S: gave up waiting for deployment"},
		{name: "K8sNoExternalIP", class: errClassValidation, regex: "K8S: gave up waiting for loadbalancer to get an ingress ip"},
		{name: "K8sNginxUnreachable", class: errClassValidation, regex: "K8S: failed to get expected response from nginx through the loadbalancer"},

		{name: "DcosNodeNotReady", class: errClassValidation, regex: "gave up waiting for DCOS nodes"},
		{name: "DcosMarathonValidationFailed", class: errClassValidation, regex: "dcos/test.sh] marathon validation failed"},
		{name: "DcosMarathonNotAdded", class: errClassValidation, regex: "dcos/test.sh] gave up waiting for marathon to be added"},
		{name: "DcosMarathonLbNotInstalled", class: errClassValidation, regex: "Failed to install marathon-lb"},

		{name: "DockerCeNetworkNotReady", class: errClassValidation, regex: "DockerCE: gave up waiting for network to be created"},
		{name: "DockerCeServiceNotReady", class: errClassValidation, regex: "DockerCE: gave up waiting for service to be created"},
		{name: "DockerCeServiceUnreachable", class: errClassValidation, regex: "DockerCE: gave up waiting for service to be externally reachable"},
	}
}

func New(jobName string, buildNum int, nDeploys int) *ReportMgr {
	h := &ReportMgr{}
	h.JobName = jobName
	h.BuildNum = buildNum
	h.Deployments = nDeploys
	h.Errors = 0
	h.StartTime = time.Now().UTC()
	h.Failures = make(map[string]*ErrorStat)
	return h
}

func (h *ReportMgr) Copy() *ReportMgr {
	n := New(h.JobName, h.BuildNum, h.Deployments)
	n.Errors = h.Errors
	n.StartTime = h.StartTime
	for e, f := range h.Failures {
		locs := make(map[string]int)
		for l, c := range f.Locations {
			locs[l] = c
		}
		n.Failures[e] = &ErrorStat{Count: f.Count, Locations: locs}
	}
	return n
}

func (h *ReportMgr) Process(txt, testName, location string) *ErrorInfo {
	for _, logErr := range logErrors {
		if match, _ := regexp.MatchString(logErr.regex, txt); match {
			h.addFailure(logErr.name, map[string]int{location: 1})
			return NewErrorInfo(testName, logErr.name, logErr.class, location)
		}
	}
	h.addFailure(errUnknown, map[string]int{location: 1})
	return NewErrorInfo(testName, errUnknown, errClassNone, location)
}

func (h *ReportMgr) addFailure(key string, locations map[string]int) {
	h.lock.Lock()
	defer h.lock.Unlock()

	cnt := 0
	if failure, ok := h.Failures[key]; !ok {
		locs := make(map[string]int)
		for l, c := range locations {
			locs[l] = c
			cnt += c
		}
		h.Failures[key] = &ErrorStat{Count: cnt, Locations: locs}
	} else {
		for l, c := range locations {
			cnt += c
			if _, ok := failure.Locations[l]; !ok {
				failure.Locations[l] = c
			} else {
				failure.Locations[l] += c
			}
		}
		failure.Count += cnt
	}
	h.Errors += cnt
}

func (h *ReportMgr) CreateTestReport(filepath string) error {
	h.Duration = time.Now().UTC().Sub(h.StartTime).String()
	data, err := json.MarshalIndent(h, "", "  ")
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

func (h *ReportMgr) CreateCombinedReport(filepath, testReportFname string) error {
	// "COMBINED_PAST_REPORTS" is the number of recent reports in the combined report
	reports, err := strconv.Atoi(os.Getenv("COMBINED_PAST_REPORTS"))
	if err != nil || reports <= 0 {
		fmt.Println("Warning: COMBINED_PAST_REPORTS is not set or invalid. Ignoring")
		return nil
	}
	combinedReport := h.Copy()
	for i := 1; i <= reports; i++ {
		data, err := ioutil.ReadFile(fmt.Sprintf("%s/%d/%s/%s",
			os.Getenv("JOB_BUILD_ROOTDIR"), h.BuildNum-i, os.Getenv("JOB_BUILD_SUBDIR"), testReportFname))
		if err != nil {
			break
		}
		testReport := &ReportMgr{}
		if err := json.Unmarshal(data, &testReport); err != nil {
			break
		}
		combinedReport.StartTime = testReport.StartTime
		combinedReport.Deployments += testReport.Deployments

		for e, f := range testReport.Failures {
			combinedReport.addFailure(e, f.Locations)
		}
	}
	return combinedReport.CreateTestReport(filepath)
}

func NewErrorInfo(testName, errName, errClass, location string) *ErrorInfo {
	return &ErrorInfo{TestName: testName, ErrName: errName, ErrClass: errClass, Location: location}
}
