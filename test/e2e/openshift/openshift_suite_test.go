package openshift_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestOpenShift(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "OpenShift Suite")
}
