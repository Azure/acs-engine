package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
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

	// deternime timeout
	timeoutMin, err := strconv.Atoi(os.Getenv("STAGE_TIMEOUT_MIN"))
	if err != nil {
		fmt.Printf("Error [Atoi STAGE_TIMEOUT_MIN]: %v\n", err)
		return err
	}
	timeout := time.Duration(time.Minute * time.Duration(timeoutMin))

	// login to Azure
	if err := runStep("set_azure_account", m.rootDir, "main", os.Environ(), fmt.Sprintf("%s/main.log", logDir), timeout); err != nil {
		return err
	}

	// return values for tests
	retvals := make([]byte, n)

	m.wg.Add(n)
	for i, d := range m.config.Deployments {
		go func(i int, d Deployment) {
			defer m.wg.Done()

			name := strings.TrimSuffix(d.ClusterDefinition, filepath.Ext(d.ClusterDefinition))
			instanceName := fmt.Sprintf("test-acs-%s-%s-%s-%d", strings.Replace(name, "/", "-", -1), d.Location, os.Getenv("BUILD_NUMBER"), i)
			logFile := fmt.Sprintf("%s/%s.log", logDir, instanceName)

			// determine orchestrator
			env := os.Environ()
			env = append(env, fmt.Sprintf("CLUSTER_DEFINITION=%s", d.ClusterDefinition))
			cmd := exec.Command("test/step.sh", "get_orchestrator_type")
			cmd.Env = env
			out, err := cmd.Output()
			if err != nil {
				wrileLog(logFile, []byte(err.Error()))
				fmt.Printf("Error [getOrchestrator %s] : %v\n", d.ClusterDefinition, err)
				retvals[i] = 1
				return
			}
			orchestrator := strings.TrimSpace(string(out))

			// update environment
			env = append(env, fmt.Sprintf("LOCATION=%s", d.Location))
			env = append(env, fmt.Sprintf("ORCHESTRATOR=%s", orchestrator))
			env = append(env, fmt.Sprintf("INSTANCE_NAME=%s", instanceName))
			env = append(env, fmt.Sprintf("DEPLOYMENT_NAME=%s", instanceName))
			env = append(env, fmt.Sprintf("RESOURCE_GROUP=%s", instanceName))

			steps := []string{"generate_template", "deploy_template"}

			// determine validation script
			validate := fmt.Sprintf("test/cluster-tests/%s/test.sh", orchestrator)
			if _, err = os.Stat(fmt.Sprintf("%s/%s", m.rootDir, validate)); err == nil {
				env = append(env, fmt.Sprintf("VALIDATE=%s", validate))
				steps = append(steps, "validate")
			}

			for _, step := range steps {
				if err = runStep(step, m.rootDir, instanceName, env, logFile, timeout); err != nil {
					retvals[i] = 1
					break
				}
				if step == "generate_template" {
					// set up extra environment variables available after template generation
					cmd := exec.Command("test/step.sh", "get_orchestrator_version")
					cmd.Env = env
					out, err := cmd.Output()
					if err != nil {
						retvals[i] = 1
						break
					}
					env = append(env, fmt.Sprintf("EXPECTED_ORCHESTRATOR_VERSION=%s", strings.TrimSpace(string(out))))

					if orchestrator == "kubernetes" {
						cmd = exec.Command("test/step.sh", "get_node_count")
						cmd.Env = env
						out, err = cmd.Output()
						if err != nil {
							retvals[i] = 1
							break
						}
						env = append(env, fmt.Sprintf("EXPECTED_NODE_COUNT=%s", strings.TrimSpace(string(out))))
					}
				}
			}
			// clean up
			runStep("cleanup", m.rootDir, instanceName, env, logFile, timeout)
		}(i, d)
	}
	m.wg.Wait()
	for _, retval := range retvals {
		if retval != 0 {
			return errors.New("Test failed")
		}
	}
	return nil
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

func runStep(step, dir, instanceName string, env []string, logFile string, timeout time.Duration) error {
	cmd := exec.Command("/bin/bash", "-c", fmt.Sprintf("%s %s", script, step))
	cmd.Dir = dir
	cmd.Env = env

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	if err := cmd.Start(); err != nil {
		fmt.Printf("Error [%s %s] : %v\n", step, instanceName, err)
		return err
	}
	timer := time.AfterFunc(timeout, func() {
		cmd.Process.Kill()
	})
	err := cmd.Wait()
	timer.Stop()

	wrileLog(logFile, out.Bytes())
	if err != nil {
		fmt.Printf("Error [%s %s] : %v\n", step, instanceName, err)
		return err
	}
	fmt.Printf("SUCCESS [%s %s]\n", step, instanceName)
	return nil
}

func wrileLog(fname string, data []byte) {
	f, err := os.OpenFile(fname, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Error [OpenFile %s] : %v\n", fname, err)
		return
	}
	defer f.Close()

	if _, err = f.Write(data); err != nil {
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
