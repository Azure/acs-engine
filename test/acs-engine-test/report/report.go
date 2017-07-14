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

	"github.com/Azure/acs-engine/test/acs-engine-test/metrics"
)

type ErrorStat struct {
	Count     int            `json:"count"`
	Locations map[string]int `json:"locations"`
}

type ReportMgr struct {
	lock       sync.Mutex
	metricsNS  string
	metricName string

	JobName     string    `json:"job"`
	BuildNum    int       `json:"build"`
	Deployments int       `json:"deployments"`
	Errors      int       `json:"errors"`
	StartTime   time.Time `json:"startTime"`
	Duration    string    `json:"duration"`
	// Failure map: key=error, value=locations
	Failures map[string]*ErrorStat `json:"failures"`
}

type errorInstance struct {
	errName  string
	errClass string
	errRegex string
}

const (
	errClassDeployment = "deployment"
	errClassValidation = "validation"
	errClassAzcli      = "azcli"
	errClassNone       = "none"

	errUnknown = "Unspecified error"
)

var errors []errorInstance

func init() {
	errors = []errorInstance{
		{errName: "azcli run", errClass: errClassAzcli, errRegex: "_init__.py"},
		{errName: "azcli load", errClass: errClassAzcli, errRegex: "Error loading command module"},

		{errName: "VMStartTimedOut", errClass: errClassDeployment, errRegex: "VMStartTimedOut"},
		{errName: "OSProvisioningTimedOut", errClass: errClassDeployment, errRegex: "OSProvisioningTimedOut"},
		{errName: "VMExtensionProvisioningError", errClass: errClassDeployment, errRegex: "VMExtensionProvisioningError"},
		{errName: "VMExtensionProvisioningTimeout", errClass: errClassDeployment, errRegex: "VMExtensionProvisioningTimeout"},
		{errName: "InternalExecutionError", errClass: errClassDeployment, errRegex: "InternalExecutionError"},
		{errName: "SkuNotAvailable", errClass: errClassDeployment, errRegex: "SkuNotAvailable"},
		{errName: "MaxStorageAccountsCountPerSubscriptionExceeded", errClass: errClassDeployment, errRegex: "MaxStorageAccountsCountPerSubscriptionExceeded"},
		{errName: "ImageManagementOperationError", errClass: errClassDeployment, errRegex: "ImageManagementOperationError"},
		{errName: "DiskProcessingError", errClass: errClassDeployment, errRegex: "DiskProcessingError"},
		{errName: "DiskServiceInternalError", errClass: errClassDeployment, errRegex: "DiskServiceInternalError"},
		{errName: "AllocationFailed", errClass: errClassDeployment, errRegex: "AllocationFailed"},
		{errName: "NetworkingInternalOperationError", errClass: errClassDeployment, errRegex: "NetworkingInternalOperationError"},
		{errName: "PlatformFaultDomainCount", errClass: errClassDeployment, errRegex: "platformFaultDomainCount"},

		{errName: "K8S nodes not ready", errClass: errClassValidation, errRegex: "K8S: gave up waiting for apiserver"},
		{errName: "K8S unexpected version", errClass: errClassValidation, errRegex: "K8S: unexpected kubernetes version"},
		{errName: "K8S containers not created", errClass: errClassValidation, errRegex: "K8S: gave up waiting for containers"},
		{errName: "K8S pods not running", errClass: errClassValidation, errRegex: "K8S: gave up waiting for running pods"},
		{errName: "K8S kube-dns not running", errClass: errClassValidation, errRegex: "K8S: gave up waiting for kube-dns"},
		{errName: "K8S dashboard not running", errClass: errClassValidation, errRegex: "K8S: gave up waiting for kubernetes-dashboard"},
		{errName: "K8S kube-proxy not running", errClass: errClassValidation, errRegex: "K8S: gave up waiting for kube-proxy"},
		{errName: "K8S proxy not working", errClass: errClassValidation, errRegex: "K8S: gave up verifying proxy"},
		{errName: "K8S deployment not ready", errClass: errClassValidation, errRegex: "K8S: gave up waiting for deployment"},
		{errName: "K8S no external IP", errClass: errClassValidation, errRegex: "K8S: gave up waiting for loadbalancer to get an ingress ip"},
		{errName: "K8S nginx unreachable", errClass: errClassValidation, errRegex: "K8S: failed to get expected response from nginx through the loadbalancer"},

		{errName: "DCOS nodes not ready", errClass: errClassValidation, errRegex: "gave up waiting for DCOS nodes"},
		{errName: "DCOS marathon validation failed", errClass: errClassValidation, errRegex: "dcos/test.sh] marathon validation failed"},
		{errName: "DCOS marathon not added", errClass: errClassValidation, errRegex: "dcos/test.sh] gave up waiting for marathon to be added"},
		{errName: "DCOS marathon-lb not installed", errClass: errClassValidation, errRegex: "Failed to install marathon-lb"},

		{errName: "DockerCE failed to create network", errClass: errClassValidation, errRegex: "DockerCE: gave up waiting for network to be created"},
		{errName: "DockerCE failed to create service", errClass: errClassValidation, errRegex: "DockerCE: gave up waiting for service to be created"},
		{errName: "DockerCE service unreachable", errClass: errClassValidation, errRegex: "DockerCE: gave up waiting for service to be externally reachable"},
	}
}

func New(metricsNS, metricName, jobName string, buildNum int, nDeploys int) *ReportMgr {
	h := &ReportMgr{}
	h.metricsNS = metricsNS
	h.metricName = metricName
	h.JobName = jobName
	h.BuildNum = buildNum
	h.Deployments = nDeploys
	h.Errors = 0
	h.StartTime = time.Now().UTC()
	h.Failures = make(map[string]*ErrorStat)
	return h
}

func (h *ReportMgr) Copy() *ReportMgr {
	n := New(h.metricsNS, h.metricName, h.JobName, h.BuildNum, h.Deployments)
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

func (h *ReportMgr) Process(txt, testName, location string) {
	for _, errInst := range errors {
		if match, _ := regexp.MatchString(errInst.errRegex, txt); match {
			h.addFailure(errInst.errName, map[string]int{location: 1})
			h.sendMetric(testName, location, errInst.errName, errInst.errClass)
			return
		}
	}
	h.addFailure(errUnknown, map[string]int{location: 1})
	h.sendMetric(testName, location, errUnknown, errClassNone)
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

func (h *ReportMgr) sendMetric(testName, location, errName, errClass string) {
	// add metrics
	dims := map[string]string{
		"test":     testName,
		"location": location,
		"errName":  errName,
		"errClass": errClass,
	}
	err := metrics.AddMetric(h.metricsNS, h.metricName, dims)
	if err != nil {
		fmt.Printf("Failed to send metric: %v", err)
	}
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
