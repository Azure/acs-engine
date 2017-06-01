package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
)

type Deployment struct {
	ClusterDefinition string `json:"cluster_definition"`
	Location          string `json:"location"`
	SkipValidation    bool   `json:"skip_validation,omitempty"`
}

type TestConfig struct {
	Deployments []Deployment `json:"deployments"`
}

func (c *TestConfig) Read(data []byte) error {
	return json.Unmarshal(data, c)
}

func (c *TestConfig) Validate() error {
	for _, d := range c.Deployments {
		if d.ClusterDefinition == "" {
			return errors.New("Cluster definition is not set")
		}
		if d.Location == "" {
			return errors.New("Location is not set")
		}
	}
	return nil
}

func getTestConfig(fname string) (*TestConfig, error) {
	data, err := ioutil.ReadFile(fname)
	if err != nil {
		return nil, err
	}
	config := &TestConfig{}
	if err = config.Read(data); err != nil {
		return nil, err
	}
	if err = config.Validate(); err != nil {
		return nil, err
	}
	return config, nil
}
