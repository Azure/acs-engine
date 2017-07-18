package main

import (
	"bufio"
	"bytes"
	"errors"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/Azure/acs-engine/test/acs-engine-test/config"
	"github.com/Azure/acs-engine/test/acs-engine-test/metrics"
	"github.com/Azure/acs-engine/test/acs-engine-test/report"
)

const (
	script = "test/step.sh"

	testReport     = "TestReport.json"
	combinedReport = "CombinedReport.json"

	metricsEndpoint = ":8125"
	metricsNS       = "ACSEngine"
	metricName      = "Error"
)

const usage = `Usage:
  acs-engine-test -c <configuration.json> -d <acs-engine root directory>

  Options:
    -c <configuration.json> : JSON file containing a list of deployment configurations.
		Refer to acs-engine/test/acs-engine-test/acs-engine-test.json for examples
	-d <acs-engine root directory>
`

var logDir string
var orchestratorRe *regexp.Regexp

func init() {
	orchestratorRe = regexp.MustCompile(`"orchestratorType": "(\S+)"`)
}

type ErrorStat struct {
	errorInfo    *report.ErrorInfo
	testCategory string
	count        int64
}

type TestManager struct {
	config    *config.TestConfig
	reportMgr *report.ReportMgr
	lock      sync.Mutex
	wg        sync.WaitGroup
	rootDir   string
}

func (m *TestManager) Run() error {
	n := len(m.config.Deployments)
	if n == 0 {
		return nil
	}

	// determine timeout
	timeoutMin, err := strconv.Atoi(os.Getenv("STAGE_TIMEOUT_MIN"))
	if err != nil {
		return fmt.Errorf("Error [Atoi STAGE_TIMEOUT_MIN]: %v", err)
	}
	timeout := time.Duration(time.Minute * time.Duration(timeoutMin))

	// determine number of retries
	retries, err := strconv.Atoi(os.Getenv("NUM_OF_RETRIES"))
	if err != nil {
		fmt.Println("Warning: NUM_OF_RETRIES is not set or invalid. Assuming 1")
		retries = 1
	}
	// login to Azure
	if _, err := m.runStep("init", "set_azure_account", os.Environ(), timeout); err != nil {
		return err
	}

	// return values for tests
	success := make([]bool, n)
	rand.Seed(time.Now().UnixNano())

	m.wg.Add(n)
	for index, dep := range m.config.Deployments {
		go func(index int, dep config.Deployment) {
			defer m.wg.Done()
			resMap := make(map[string]*ErrorStat)
			for attempt := 0; attempt < retries; attempt++ {
				errorInfo := m.testRun(dep, index, attempt, timeout)
				// do not retry if successful
				if errorInfo == nil {
					success[index] = true
					break
				}
				if errorStat, ok := resMap[errorInfo.ErrName]; !ok {
					resMap[errorInfo.ErrName] = &ErrorStat{errorInfo: errorInfo, testCategory: dep.TestCategory, count: 1}
				} else {
					errorStat.count++
				}
			}
			sendMetrics(resMap)
		}(index, dep)
	}
	m.wg.Wait()
	//create reports
	if err = m.reportMgr.CreateTestReport(fmt.Sprintf("%s/%s", logDir, testReport)); err != nil {
		fmt.Printf("Failed to create %s: %v\n", testReport, err)
	}
	if err = m.reportMgr.CreateCombinedReport(fmt.Sprintf("%s/%s", logDir, combinedReport), testReport); err != nil {
		fmt.Printf("Failed to create %s: %v\n", combinedReport, err)
	}
	// fail the test on error
	for _, ok := range success {
		if !ok {
			return errors.New("Test failed")
		}
	}
	return nil
}

func (m *TestManager) testRun(d config.Deployment, index, attempt int, timeout time.Duration) *report.ErrorInfo {
	rgPrefix := os.Getenv("RESOURCE_GROUP_PREFIX")
	if rgPrefix == "" {
		rgPrefix = "x"
		fmt.Printf("RESOURCE_GROUP_PREFIX is not set. Using default '%s'\n", rgPrefix)
	}
	testName := strings.TrimSuffix(d.ClusterDefinition, filepath.Ext(d.ClusterDefinition))
	instanceName := fmt.Sprintf("acse-%d-%s-%s-%d-%d", rand.Intn(0x0ffffff), d.Location, os.Getenv("BUILD_NUMBER"), index, attempt)
	resourceGroup := fmt.Sprintf("%s-%s-%s-%s-%d-%d", rgPrefix, strings.Replace(testName, "/", "-", -1), d.Location, os.Getenv("BUILD_NUMBER"), index, attempt)
	logFile := fmt.Sprintf("%s/%s.log", logDir, resourceGroup)
	validateLogFile := fmt.Sprintf("%s/validate-%s.log", logDir, resourceGroup)

	// determine orchestrator
	env := os.Environ()
	env = append(env, fmt.Sprintf("CLUSTER_DEFINITION=examples/%s", d.ClusterDefinition))
	cmd := exec.Command("test/step.sh", "get_orchestrator_type")
	cmd.Env = env
	out, err := cmd.Output()
	if err != nil {
		wrileLog(logFile, "Error [getOrchestrator %s] : %v", d.ClusterDefinition, err)
		return report.NewErrorInfo(testName, "OrchestratorTypeParsingError", "PreRun", d.Location)
	}
	orchestrator := strings.TrimSpace(string(out))

	// update environment
	env = append(env, fmt.Sprintf("LOCATION=%s", d.Location))
	env = append(env, fmt.Sprintf("ORCHESTRATOR=%s", orchestrator))
	env = append(env, fmt.Sprintf("INSTANCE_NAME=%s", instanceName))
	env = append(env, fmt.Sprintf("DEPLOYMENT_NAME=%s", instanceName))
	env = append(env, fmt.Sprintf("RESOURCE_GROUP=%s", resourceGroup))

	// add scenario-specific environment variables
	envFile := fmt.Sprintf("examples/%s.env", d.ClusterDefinition)
	if _, err = os.Stat(envFile); err == nil {
		envHandle, err := os.Open(envFile)
		if err != nil {
			wrileLog(logFile, "Error [open %s] : %v", envFile, err)
			return report.NewErrorInfo(testName, "FileAccessError", "PreRun", d.Location)
		}
		defer envHandle.Close()

		fileScanner := bufio.NewScanner(envHandle)
		for fileScanner.Scan() {
			str := strings.TrimSpace(fileScanner.Text())
			if match, _ := regexp.MatchString(`^\S+=\S+$`, str); match {
				env = append(env, str)
			}
		}
	}

	var errorInfo *report.ErrorInfo
	steps := []string{"create_resource_group", "predeploy", "generate_template", "deploy_template", "postdeploy"}

	// determine validation script
	if !d.SkipValidation {
		validate := fmt.Sprintf("test/cluster-tests/%s/test.sh", orchestrator)
		if _, err = os.Stat(fmt.Sprintf("%s/%s", m.rootDir, validate)); err == nil {
			env = append(env, fmt.Sprintf("VALIDATE=%s", validate))
			steps = append(steps, "validate")
		}
	}
	for _, step := range steps {
		txt, err := m.runStep(resourceGroup, step, env, timeout)
		if err != nil {
			errorInfo = m.reportMgr.Process(txt, testName, d.Location)
			wrileLog(logFile, "Error [%s:%s] %v\nOutput: %s", step, resourceGroup, err, txt)
			// check AUTOCLEAN flag: if set to 'n', don't remove deployment
			if os.Getenv("AUTOCLEAN") == "n" {
				env = append(env, "CLEANUP=n")
			}
			break
		}
		wrileLog(logFile, txt)
		if step == "generate_template" {
			// set up extra environment variables available after template generation
			validateLogFile = fmt.Sprintf("%s/validate-%s.log", logDir, resourceGroup)
			env = append(env, fmt.Sprintf("LOGFILE=%s", validateLogFile))

			cmd := exec.Command("test/step.sh", "get_orchestrator_version")
			cmd.Env = env
			out, err := cmd.Output()
			if err != nil {
				wrileLog(logFile, "Error [%s:%s] %v", "get_orchestrator_version", resourceGroup, err)
				errorInfo = report.NewErrorInfo(testName, "OrchestratorVersionParsingError", "PreRun", d.Location)
				break
			}
			env = append(env, fmt.Sprintf("EXPECTED_ORCHESTRATOR_VERSION=%s", strings.TrimSpace(string(out))))

			cmd = exec.Command("test/step.sh", "get_node_count")
			cmd.Env = env
			out, err = cmd.Output()
			if err != nil {
				wrileLog(logFile, "Error [%s:%s] %v", "get_node_count", resourceGroup, err)
				errorInfo = report.NewErrorInfo(testName, "NodeCountParsingError", "PreRun", d.Location)
				break
			}
			nodesCount := strings.Split(strings.TrimSpace(string(out)), ":")
			if len(nodesCount) != 2 {
				wrileLog(logFile, "get_node_count: unexpected output '%s'", string(out))
				errorInfo = report.NewErrorInfo(testName, "NodeCountParsingError", "PreRun", d.Location)
				break
			}
			env = append(env, fmt.Sprintf("EXPECTED_NODE_COUNT=%s", nodesCount[0]))
			env = append(env, fmt.Sprintf("EXPECTED_LINUX_NODE_COUNT=%s", nodesCount[1]))
		}
	}
	// clean up
	if txt, err := m.runStep(resourceGroup, "cleanup", env, timeout); err != nil {
		wrileLog(logFile, "Error: %v\nOutput: %s", err, txt)
	}
	if errorInfo == nil {
		// do not keep logs for successful test
		for _, fname := range []string{logFile, validateLogFile} {
			if _, err := os.Stat(fname); !os.IsNotExist(err) {
				if err = os.Remove(fname); err != nil {
					fmt.Printf("Failed to remove %s : %v\n", fname, err)
				}
			}
		}
	}
	return errorInfo
}

func isValidEnv() bool {
	valid := true
	envVars := []string{
		"SERVICE_PRINCIPAL_CLIENT_ID",
		"SERVICE_PRINCIPAL_CLIENT_SECRET",
		"TENANT_ID",
		"SUBSCRIPTION_ID",
		"CLUSTER_SERVICE_PRINCIPAL_CLIENT_ID",
		"CLUSTER_SERVICE_PRINCIPAL_CLIENT_SECRET",
		"STAGE_TIMEOUT_MIN"}

	for _, envVar := range envVars {
		if os.Getenv(envVar) == "" {
			fmt.Printf("Must specify environment variable %s\n", envVar)
			valid = false
		}
	}
	return valid
}

func (m *TestManager) runStep(name, step string, env []string, timeout time.Duration) (string, error) {
	// prevent ARM throttling
	m.lock.Lock()
	go func() {
		time.Sleep(2 * time.Second)
		m.lock.Unlock()
	}()

	cmd := exec.Command("/bin/bash", "-c", fmt.Sprintf("%s %s", script, step))
	cmd.Dir = m.rootDir
	cmd.Env = env

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	if err := cmd.Start(); err != nil {
		return "", err
	}
	timer := time.AfterFunc(timeout, func() {
		cmd.Process.Kill()
	})
	err := cmd.Wait()
	timer.Stop()

	now := time.Now().Format("15:04:05")
	if err != nil {
		fmt.Printf("ERROR [%s] [%s %s]\n", now, step, name)
		return out.String(), err
	}
	fmt.Printf("SUCCESS [%s] [%s %s]\n", now, step, name)
	return out.String(), nil
}

func wrileLog(fname string, format string, args ...interface{}) {
	str := fmt.Sprintf(format, args...)

	f, err := os.OpenFile(fname, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Error [OpenFile %s] : %v\n", fname, err)
		return
	}
	defer f.Close()

	if _, err = f.Write([]byte(str)); err != nil {
		fmt.Printf("Error [Write %s] : %v\n", fname, err)
	}
}

func sendMetrics(resMap map[string]*ErrorStat) {
	for _, errorStat := range resMap {
		var severity string
		if errorStat.count > 1 {
			severity = "Critical"
		} else {
			severity = "Intermittent"
		}
		category := errorStat.testCategory
		if len(category) == 0 {
			category = "generic"
		}
		// add metrics
		dims := map[string]string{
			"TestName":     errorStat.errorInfo.TestName,
			"TestCategory": category,
			"Location":     errorStat.errorInfo.Location,
			"Error":        errorStat.errorInfo.ErrName,
			"Class":        errorStat.errorInfo.ErrClass,
			"Severity":     severity,
		}
		err := metrics.AddMetric(metricsEndpoint, metricsNS, metricName, errorStat.count, dims)
		if err != nil {
			fmt.Printf("Failed to send metric: %v\n", err)
		}
	}
}

func mainInternal() error {
	var configFile string
	var rootDir string
	var err error
	flag.StringVar(&configFile, "c", "", "deployment configurations")
	flag.StringVar(&rootDir, "d", "", "acs-engine root directory")
	flag.Usage = func() {
		fmt.Println(usage)
	}
	flag.Parse()

	testManager := TestManager{}

	// validate environment
	if !isValidEnv() {
		return fmt.Errorf("environment is not set")
	}
	// get test configuration
	if configFile == "" {
		return fmt.Errorf("test configuration is not provided")
	}
	testManager.config, err = config.GetTestConfig(configFile)
	if err != nil {
		return err
	}
	// get Jenkins build number
	buildNum, err := strconv.Atoi(os.Getenv("BUILD_NUMBER"))
	if err != nil {
		fmt.Println("Warning: BUILD_NUMBER is not set or invalid. Assuming 0")
		buildNum = 0
	}
	// initialize report manager
	testManager.reportMgr = report.New(os.Getenv("JOB_BASE_NAME"), buildNum, len(testManager.config.Deployments))
	// check root directory
	if rootDir == "" {
		return fmt.Errorf("acs-engine root directory is not provided")
	}
	testManager.rootDir = rootDir
	if _, err = os.Stat(fmt.Sprintf("%s/%s", rootDir, script)); err != nil {
		return err
	}
	// make logs directory
	logDir = fmt.Sprintf("%s/_logs", rootDir)
	os.RemoveAll(logDir)
	if err = os.Mkdir(logDir, os.FileMode(0755)); err != nil {
		return err
	}
	// run tests
	return testManager.Run()
}

func main() {
	if err := mainInternal(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	os.Exit(0)
}
