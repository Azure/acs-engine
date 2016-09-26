package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"

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

func usage(errs ...error) {
	for _, err := range errs {
		fmt.Printf("error: %s\n\n", err.Error())
	}
	fmt.Printf("usage: %s ClusterDefinitionFile\n", os.Args[0])
	fmt.Println("       read the ClusterDefinitionFile and output an arm template")
	fmt.Println()
	fmt.Println("options:")
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
		usage(err)
		os.Exit(1)
	}

	if acsCluster, err = loadAcsCluster(jsonFile); err != nil {
		usage(fmt.Errorf("error while loading %s: %s", jsonFile, err.Error()))
		os.Exit(1)
	}

	if template, err = clustertemplate.GenerateTemplate(acsCluster, *templateDirectory); err != nil {
		usage(fmt.Errorf("error generating template %s: %s", jsonFile, err.Error()))
		os.Exit(1)
	}

	fmt.Print(template)
}
