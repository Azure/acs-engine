package main

import (
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
)

const script = "test/step.sh"

const usage = `Usage:
  acs-engine-test -c <configuration.json> -d <acs-engine root directory>

  Options:
    -c <configuration.json> : JSON file containing a list of deployment configurations.
		Refer to acs-engine/test/acs-engine-test/acs-engine-test.json for examples
	-d <acs-engine root directory>
`

var logDir string
var orchestrator_re *regexp.Regexp

func init() {
	orchestrator_re = regexp.MustCompile(`"orchestratorType": "(\S+)"`)
}

type TestManager struct {
	config  *TestConfig
	lock    sync.Mutex
	wg      sync.WaitGroup
	rootDir string
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
		go func(index int, dep Deployment) {
			defer m.wg.Done()
			for attempt := 0; attempt < retries; attempt++ {
				success[index] = m.testRun(dep, index, attempt, timeout)
				// do not retry if successful
				if success[index] {
					break
				}
			}
		}(index, dep)
	}
	m.wg.Wait()
	for _, ok := range success {
		if !ok {
			return errors.New("Test failed")
		}
	}
	return nil
}

func (m *TestManager) testRun(d Deployment, index, attempt int, timeout time.Duration) bool {
	name := strings.TrimSuffix(d.ClusterDefinition, filepath.Ext(d.ClusterDefinition))
	instanceName := fmt.Sprintf("acse-%d-%s-%s-%d-%d", rand.Intn(0x0ffffff), d.Location, os.Getenv("BUILD_NUMBER"), index, attempt)
	resourceGroup := fmt.Sprintf("x-%s-%s-%s-%d-%d", strings.Replace(name, "/", "-", -1), d.Location, os.Getenv("BUILD_NUMBER"), index, attempt)
	logFile := fmt.Sprintf("%s/%s.log", logDir, resourceGroup)
	validateLogFile := fmt.Sprintf("%s/validate-%s.log", logDir, resourceGroup)
	success := true

	// determine orchestrator
	env := os.Environ()
	env = append(env, fmt.Sprintf("CLUSTER_DEFINITION=examples/%s", d.ClusterDefinition))
	cmd := exec.Command("test/step.sh", "get_orchestrator_type")
	cmd.Env = env
	out, err := cmd.Output()
	if err != nil {
		wrileLog(logFile, "Error [getOrchestrator %s] : %v", d.ClusterDefinition, err)
		return false
	}
	orchestrator := strings.TrimSpace(string(out))

	// update environment
	env = append(env, fmt.Sprintf("LOCATION=%s", d.Location))
	env = append(env, fmt.Sprintf("ORCHESTRATOR=%s", orchestrator))
	env = append(env, fmt.Sprintf("INSTANCE_NAME=%s", instanceName))
	env = append(env, fmt.Sprintf("DEPLOYMENT_NAME=%s", instanceName))
	env = append(env, fmt.Sprintf("RESOURCE_GROUP=%s", resourceGroup))

	steps := []string{"generate_template", "deploy_template"}

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
			wrileLog(logFile, "Error [%s:%s] %v\nOutput: %s", step, resourceGroup, err, txt)
			success = false
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
				success = false
				break
			}
			env = append(env, fmt.Sprintf("EXPECTED_ORCHESTRATOR_VERSION=%s", strings.TrimSpace(string(out))))

			if orchestrator == "kubernetes" {
				cmd = exec.Command("test/step.sh", "get_node_count")
				cmd.Env = env
				out, err = cmd.Output()
				if err != nil {
					wrileLog(logFile, "Error [%s:%s] %v", "get_node_count", resourceGroup, err)
					success = false
					break
				}
				env = append(env, fmt.Sprintf("EXPECTED_NODE_COUNT=%s", strings.TrimSpace(string(out))))
			}
		}
	}
	// clean up
	if txt, err := m.runStep(resourceGroup, "cleanup", env, timeout); err != nil {
		wrileLog(logFile, "Error: %v\nOutput: %s", err, txt)
	}
	if success {
		// do not keep logs for successful test
		for _, fname := range []string{logFile, validateLogFile} {
			if _, err := os.Stat(fname); !os.IsNotExist(err) {
				if err = os.Remove(fname); err != nil {
					fmt.Printf("Failed to remove %s : %v\n", fname, err)
				}
			}
		}
	}
	return success
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
	// work around az-cli parallelization issue https://github.com/Azure/azure-cli/issues/3255
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

func main_internal() error {
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
	testManager.config, err = getTestConfig(configFile)
	if err != nil {
		return err
	}
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
	if err := main_internal(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	os.Exit(0)
}
