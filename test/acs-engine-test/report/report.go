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

var errorRegexpMap map[string]string

func init() {
	errorRegexpMap = map[string]string{
		"azcli run":  "_init__.py",
		"azcli load": "Error loading command module",

		"VMStartTimedOut":                                "VMStartTimedOut",
		"OSProvisioningTimedOut":                         "OSProvisioningTimedOut",
		"VMExtensionProvisioningError":                   "VMExtensionProvisioningError",
		"VMExtensionProvisioningTimeout":                 "VMExtensionProvisioningTimeout",
		"InternalExecutionError":                         "InternalExecutionError",
		"SkuNotAvailable":                                "SkuNotAvailable",
		"MaxStorageAccountsCountPerSubscriptionExceeded": "MaxStorageAccountsCountPerSubscriptionExceeded",
		"ImageManagementOperationError":                  "ImageManagementOperationError",
		"DiskProcessingError":                            "DiskProcessingError",
		"DiskServiceInternalError":                       "DiskServiceInternalError",
		"AllocationFailed":                               "AllocationFailed",
		"NetworkingInternalOperationError":               "NetworkingInternalOperationError",
		"PlatformFaultDomainCount":                       "platformFaultDomainCount",

		"K8S curl error":                  "curl_error",
		"K8S no external IP":              "gave up waiting for loadbalancer to get an ingress ip",
		"K8S nodes not ready":             "gave up waiting for apiserver",
		"K8S service unreachable":         "gave up waiting for service to be externally reachable",
		"K8S nginx unreachable":           "failed to get expected response from nginx through the loadbalancer",
		"DCOS nodes not ready":            "gave up waiting for DCOS nodes",
		"DCOS marathon validation failed": "dcos/test.sh] marathon validation failed",
		"DCOS marathon not added":         "dcos/test.sh] gave up waiting for marathon to be added",
		"DCOS marathon-lb not installed":  "Failed to install marathon-lb",
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

func (h *ReportMgr) Process(txt, location string) {
	for key, regex := range errorRegexpMap {
		if match, _ := regexp.MatchString(regex, txt); match {
			h.addFailure(key, map[string]int{location: 1})
			return
		}
	}
	h.addFailure("Unspecified error", map[string]int{location: 1})
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
