package main

import "testing"

func TestConfigParse(t *testing.T) {

	test_cfg := `
{"deployments":
  [
    {
      "cluster_definition":"examples/kubernetes.json",
      "location":"westus",
      "skip_validation":true
    },
    {
      "cluster_definition":"examples/dcos.json",
      "location":"eastus",
      "skip_validation":false
    },
    {
      "cluster_definition":"examples/swarm.json",
      "location":"southcentralus"
    },
    {
      "cluster_definition":"examples/swarmmode.json",
      "location":"westus2"
    }
  ]
}
`

	testConfig := TestConfig{}
	if err := testConfig.Read([]byte(test_cfg)); err != nil {
		t.Fatal(err)
	}
	if err := testConfig.Validate(); err != nil {
		t.Fatal(err)
	}
	if len(testConfig.Deployments) != 4 {
		t.Fatalf("Wrong number of deployments: %d instead of 4", len(testConfig.Deployments))
	}
}
