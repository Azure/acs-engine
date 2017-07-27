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

// ErrorInfo represents the CI error
type ErrorInfo struct {
	TestName string
	ErrName  string
	ErrClass string
	Location string
}

// ErrorStat represents the aggregate error count and region
type ErrorStat struct {
	Count     int            `json:"count"`
	Locations map[string]int `json:"locations"`
}

// Manager represents the details about a build and errors in that build
type Manager struct {
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
	// ErrClassDeployment represents an error during deployment
	ErrClassDeployment = "Deployment"
	// ErrClassValidation represents an error during validation (tests)
	ErrClassValidation = "Validation"
	// ErrClassAzcli represents an error with Azure CLI
	ErrClassAzcli = "AzCLI"
	// ErrClassNone represents absence of error
	ErrClassNone = "None"

	// ErrSuccess represents a success, for some reason
	ErrSuccess = "Success"
	// ErrUnknown represents an unknown error
	ErrUnknown = "UnspecifiedError"
)

var logErrors []logError

func init() {
	logErrors = []logError{
		{name: "AzCliRunError", class: ErrClassAzcli, regex: "_init__.py"},
		{name: "AzCliLoadError", class: ErrClassAzcli, regex: "Error loading command module"},

		{name: "VMStartTimedOut", class: ErrClassDeployment, regex: "VMStartTimedOut"},
		{name: "OSProvisioningTimedOut", class: ErrClassDeployment, regex: "OSProvisioningTimedOut"},
		{name: "VMExtensionProvisioningError", class: ErrClassDeployment, regex: "VMExtensionProvisioningError"},
		{name: "VMExtensionProvisioningTimeout", class: ErrClassDeployment, regex: "VMExtensionProvisioningTimeout"},
		{name: "InternalExecutionError", class: ErrClassDeployment, regex: "InternalExecutionError"},
		{name: "SkuNotAvailable", class: ErrClassDeployment, regex: "SkuNotAvailable"},
		{name: "MaxStorageAccountsCountPerSubscriptionExceeded", class: ErrClassDeployment, regex: "MaxStorageAccountsCountPerSubscriptionExceeded"},
		{name: "ImageManagementOperationError", class: ErrClassDeployment, regex: "ImageManagementOperationError"},
		{name: "DiskProcessingError", class: ErrClassDeployment, regex: "DiskProcessingError"},
		{name: "DiskServiceInternalError", class: ErrClassDeployment, regex: "DiskServiceInternalError"},
		{name: "AllocationFailed", class: ErrClassDeployment, regex: "AllocationFailed"},
		{name: "NetworkingInternalOperationError", class: ErrClassDeployment, regex: "NetworkingInternalOperationError"},
		{name: "PlatformFaultDomainCount", class: ErrClassDeployment, regex: "platformFaultDomainCount"},

		{name: "K8sNodeNotReady", class: ErrClassValidation, regex: "K8S: gave up waiting for apiserver"},
		{name: "K8sUnexpectedVersion", class: ErrClassValidation, regex: "K8S: unexpected kubernetes version"},
		{name: "K8sContainerNotCreated", class: ErrClassValidation, regex: "K8S: gave up waiting for containers"},
		{name: "K8sPodNotRunning", class: ErrClassValidation, regex: "K8S: gave up waiting for running pods"},
		{name: "K8sKubeDnsNotRunning", class: ErrClassValidation, regex: "K8S: gave up waiting for kube-dns"},
		{name: "K8sDashboardNotRunning", class: ErrClassValidation, regex: "K8S: gave up waiting for kubernetes-dashboard"},
		{name: "K8sKubeProxyNotRunning", class: ErrClassValidation, regex: "K8S: gave up waiting for kube-proxy"},
		{name: "K8sProxyNotWorking", class: ErrClassValidation, regex: "K8S: gave up verifying proxy"},
		{name: "K8sLinuxDeploymentNotReady", class: ErrClassValidation, regex: "K8S-Linux: gave up waiting for deployment"},
		{name: "K8sWindowsDeploymentNotReady", class: ErrClassValidation, regex: "K8S-Windows: gave up waiting for deployment"},
		{name: "K8sLinuxNoExternalIP", class: ErrClassValidation, regex: "K8S-Linux: gave up waiting for loadbalancer to get an ingress ip"},
		{name: "K8sWindowsNoExternalIP", class: ErrClassValidation, regex: "K8S-Windows: gave up waiting for loadbalancer to get an ingress ip"},
		{name: "K8sLinuxNginxUnreachable", class: ErrClassValidation, regex: "K8S-Linux: failed to get expected response from nginx through the loadbalancer"},
		{name: "K8sWindowsSimpleWebUnreachable", class: ErrClassValidation, regex: "K8S-Windows: failed to get expected response from simpleweb through the loadbalancer"},
		{name: "K8sWindowsNoSimpleWebPodName", class: ErrClassValidation, regex: "K8S-Windows: failed to get expected pod name for simpleweb"},
		{name: "K8sWindowsNoSimpleWebOutboundInternet", class: ErrClassValidation, regex: "K8S-Windows: failed to get outbound internet connection inside simpleweb container"},

		{name: "DcosNodeNotReady", class: ErrClassValidation, regex: "gave up waiting for DCOS nodes"},
		{name: "DcosMarathonValidationFailed", class: ErrClassValidation, regex: "dcos/test.sh] marathon validation failed"},
		{name: "DcosMarathonNotAdded", class: ErrClassValidation, regex: "dcos/test.sh] gave up waiting for marathon to be added"},
		{name: "DcosMarathonLbNotInstalled", class: ErrClassValidation, regex: "Failed to install marathon-lb"},

		{name: "DockerCeNetworkNotReady", class: ErrClassValidation, regex: "DockerCE: gave up waiting for network to be created"},
		{name: "DockerCeServiceNotReady", class: ErrClassValidation, regex: "DockerCE: gave up waiting for service to be created"},
		{name: "DockerCeServiceUnreachable", class: ErrClassValidation, regex: "DockerCE: gave up waiting for service to be externally reachable"},
	}
}

// New creates a new error report
func New(jobName string, buildNum int, nDeploys int) *Manager {
	h := &Manager{}
	h.JobName = jobName
	h.BuildNum = buildNum
	h.Deployments = nDeploys
	h.Errors = 0
	h.StartTime = time.Now().UTC()
	h.Failures = make(map[string]*ErrorStat)
	return h
}

// TBD
func (h *Manager) Copy() *Manager {
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

// TBD
func (h *Manager) Process(txt, testName, location string) *ErrorInfo {
	for _, logErr := range logErrors {
		if match, _ := regexp.MatchString(logErr.regex, txt); match {
			h.addFailure(logErr.name, map[string]int{location: 1})
			return NewErrorInfo(testName, logErr.name, logErr.class, location)
		}
	}
	h.addFailure(ErrUnknown, map[string]int{location: 1})
	return NewErrorInfo(testName, ErrUnknown, ErrClassNone, location)
}

func (h *Manager) addFailure(key string, locations map[string]int) {
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

// TBD
func (h *Manager) CreateTestReport(filepath string) error {
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

// TBD
func (h *Manager) CreateCombinedReport(filepath, testReportFname string) error {
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
		testReport := &Manager{}
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

// TBD
func NewErrorInfo(testName, errName, errClass, location string) *ErrorInfo {
	return &ErrorInfo{TestName: testName, ErrName: errName, ErrClass: errClass, Location: location}
}
