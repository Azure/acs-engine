package acsengine

import (
	"crypto/rsa"
	"crypto/x509"
	"testing"
)

func TestCreateCertificateWithOrganisation(t *testing.T) {
	var err error
	var caPair *PkiKeyCertPair

	var (
		caCertificate   *x509.Certificate
		caPrivateKey    *rsa.PrivateKey
		testCertificate *x509.Certificate
	)

	caCertificate, caPrivateKey, err = createCertificate("ca", nil, nil, false, false, nil, nil, nil)
	if err != nil {
		t.Fatalf("failed to generate certificate: %s", err)
	}
	caPair = &PkiKeyCertPair{CertificatePem: string(certificateToPem(caCertificate.Raw)), PrivateKeyPem: string(privateKeyToPem(caPrivateKey))}

	caCertificate, err = pemToCertificate(caPair.CertificatePem)
	if err != nil {
		t.Fatalf("failed to generate certificate: %s", err)
	}
	caPrivateKey, err = pemToKey(caPair.PrivateKeyPem)
	if err != nil {
		t.Fatalf("failed to generate certificate: %s", err)
	}

	organization := make([]string, 1)
	organization[0] = "system:masters"
	testCertificate, _, err = createCertificate("client", caCertificate, caPrivateKey, false, false, nil, nil, organization)
	if err != nil {
		t.Fatalf("failed to generate certificate: %s", err)
	}

	certificationOrganization := testCertificate.Subject.Organization

	if certificationOrganization[0] != organization[0] || len(certificationOrganization) != len(organization) {
		t.Fatalf("certificate organisation did not match")
	}
}

func TestCreateCertificateWithoutOrganisation(t *testing.T) {
	var err error
	var caPair *PkiKeyCertPair

	var (
		caCertificate   *x509.Certificate
		caPrivateKey    *rsa.PrivateKey
		testCertificate *x509.Certificate
	)

	caCertificate, caPrivateKey, err = createCertificate("ca", nil, nil, false, false, nil, nil, nil)
	if err != nil {
		t.Fatalf("failed to generate certificate: %s", err)
	}
	caPair = &PkiKeyCertPair{CertificatePem: string(certificateToPem(caCertificate.Raw)), PrivateKeyPem: string(privateKeyToPem(caPrivateKey))}

	caCertificate, err = pemToCertificate(caPair.CertificatePem)
	if err != nil {
		t.Fatalf("failed to generate certificate: %s", err)
	}
	caPrivateKey, err = pemToKey(caPair.PrivateKeyPem)
	if err != nil {
		t.Fatalf("failed to generate certificate: %s", err)
	}

	testCertificate, _, err = createCertificate("client", caCertificate, caPrivateKey, false, false, nil, nil, nil)
	if err != nil {
		t.Fatalf("failed to generate certificate: %s", err)
	}

	certificationOrganization := testCertificate.Subject.Organization

	if len(certificationOrganization) != 0 {
		t.Fatalf("certificate organisation should be empty but has length %d", len(certificationOrganization))
	}
}
