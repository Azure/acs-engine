package dcos_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestDcos(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Dcos Suite")
}
