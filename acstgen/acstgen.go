package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"./api/vlabs"
	"./clustertemplate"
)

// loadAcsCluster loads an ACS Cluster API Model from a JSON file
func loadAcsCluster(jsonFile string) (*vlabs.AcsCluster, error) {
	contents, e := ioutil.ReadFile(jsonFile)
	if e != nil {
		return nil, fmt.Errorf("error reading file %s: %s", jsonFile, e.Error())
	}

	acsCluster := &vlabs.AcsCluster{}
	if e := json.Unmarshal(contents, &acsCluster); e != nil {
		return nil, fmt.Errorf("error unmarshalling file %s: %s", jsonFile, e.Error())
	}
	acsCluster.SetDefaults()
	if e := acsCluster.Validate(); e != nil {
		return nil, fmt.Errorf("error validating acs cluster from file %s: %s", jsonFile, e.Error())
	}

	return acsCluster, nil
}

func translateJSON(content string, translateParams [][]string, reverseTranslate bool) string {
	for _, tuple := range translateParams {
		if len(tuple) != 2 {
			panic("string tuples must be of size 2")
		}
		a := tuple[0]
		b := tuple[1]
		if reverseTranslate {
			content = strings.Replace(content, b, a, -1)
		} else {
			content = strings.Replace(content, a, b, -1)
		}
	}
	return content
}

func prettyPrintJSON(content string) (string, error) {
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(content), &data); err != nil {
		return "", err
	}
	prettyprint, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return "", err
	}
	return string(prettyprint), nil
}

func prettyPrintArmTemplate(template string) (string, error) {
	translateParams := [][]string{
		{"parameters", "dparameters"},
		{"variables", "eparameters"},
		{"resources", "fresources"},
		{"outputs", "zoutputs"},
	}

	template = translateJSON(template, translateParams, false)
	var err error
	if template, err = prettyPrintJSON(template); err != nil {
		return "", err
	}
	template = translateJSON(template, translateParams, true)

	return template, nil
}

func usage(errs ...error) {
	for _, err := range errs {
		fmt.Fprintf(os.Stderr, "error: %s\n\n", err.Error())
	}
	fmt.Fprintf(os.Stderr, "usage: %s ClusterDefinitionFile\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "       read the ClusterDefinitionFile and output an arm template")
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "options:")
	flag.PrintDefaults()
}

var templateDirectory = flag.String("templateDirectory", "./parts", "directory containing base template files")

func main() {
	var acsCluster *vlabs.AcsCluster
	var template string
	var err error

	flag.Parse()

	if argCount := len(flag.Args()); argCount == 0 {
		usage()
		os.Exit(1)
	}

	jsonFile := flag.Arg(0)
	if _, err = os.Stat(jsonFile); os.IsNotExist(err) {
		usage(fmt.Errorf("file %s does not exist", jsonFile))
		os.Exit(1)
	}

	if _, err = os.Stat(*templateDirectory); os.IsNotExist(err) {
		usage(fmt.Errorf("base templates directory %s does not exist", jsonFile))
		os.Exit(1)
	}

	if err = clustertemplate.VerifyFiles(*templateDirectory); err != nil {
		fmt.Fprintf(os.Stderr, "verification failed: %s\n", err.Error())
		os.Exit(1)
	}

	if acsCluster, err = loadAcsCluster(jsonFile); err != nil {
		fmt.Fprintf(os.Stderr, "error while loading %s: %s", jsonFile, err.Error())
		os.Exit(1)
	}

	if template, err = clustertemplate.GenerateTemplate(acsCluster, *templateDirectory); err != nil {
		fmt.Fprintf(os.Stderr, "error generating template %s: %s", jsonFile, err.Error())
		os.Exit(1)
	}

	if template, err = prettyPrintArmTemplate(template); err != nil {
		fmt.Fprintf(os.Stderr, "error pretty printing template %s", err.Error())
		os.Exit(1)
	}

	fmt.Print(template)
}
