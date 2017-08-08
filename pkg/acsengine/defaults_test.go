// +build ignore

package acsengine

import (
	"testing"

	"github.com/Azure/acs-engine/pkg/api"
	. "github.com/onsi/gomega"
)

func TestCertAlreadyPresent(t *testing.T) {
	RegisterTestingT(t)
	var cert *api.CertificateProfile

	Expect(certAlreadyPresent(nil)).To(BeFalse())

	cert = &api.CertificateProfile{}
	Expect(certAlreadyPresent(cert)).To(BeFalse())

	cert = &api.CertificateProfile{
		APIServerCertificate: "a",
	}
	Expect(certAlreadyPresent(cert)).To(BeTrue())

	cert = &api.CertificateProfile{
		APIServerPrivateKey: "b",
	}
	Expect(certAlreadyPresent(cert)).To(BeTrue())

	cert = &api.CertificateProfile{
		ClientCertificate: "c",
	}
	Expect(certAlreadyPresent(cert)).To(BeTrue())

	cert = &api.CertificateProfile{
		ClientPrivateKey: "d",
	}
	Expect(certAlreadyPresent(cert)).To(BeTrue())
}
