package api

import (
	"testing"

	"github.com/Azure/acs-engine/pkg/api/vlabs"
)

func TestConvertKubernetesDebugToVlabs(t *testing.T) {
	a := &KubernetesConfig{
		Debug: map[string]string{
			"waitForNodes": "true",
		},
	}
	v := &vlabs.KubernetesConfig{}
	convertKubernetesDebugToVlabs(a, v)
	for key, val := range v.Debug {
		if a.Debug[key] != val {
			t.Fatalf("got unexpected kubernetes debug config value for %s: %s, expected %s",
				key, a.Debug[key], val)
		}
	}
}
