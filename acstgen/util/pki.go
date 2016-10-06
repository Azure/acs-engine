// util package comes from https://github.com/colemickens/azkube/blob/master/util/pki.go
package util

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"net"
	"time"
)

const (
	ValidityDuration = time.Hour * 24 * 365 * 2
	PkiKeySize       = 4096
)

type PkiKeyCertPair struct {
	CertificatePem string
	PrivateKeyPem  string
}

func CreatePki(masterFQDN string, extraFQDNs []string, extraIPs []net.IP, clusterDomain string) (*PkiKeyCertPair, *PkiKeyCertPair, *PkiKeyCertPair, error) {
	extraFQDNs = append(extraFQDNs, fmt.Sprintf("kubernetes"))
	extraFQDNs = append(extraFQDNs, fmt.Sprintf("kubernetes.default"))
	extraFQDNs = append(extraFQDNs, fmt.Sprintf("kubernetes.default.svc"))
	extraFQDNs = append(extraFQDNs, fmt.Sprintf("kubernetes.default.svc.%s", clusterDomain))
	extraFQDNs = append(extraFQDNs, fmt.Sprintf("kubernetes.kube-system"))
	extraFQDNs = append(extraFQDNs, fmt.Sprintf("kubernetes.kube-system.svc"))
	extraFQDNs = append(extraFQDNs, fmt.Sprintf("kubernetes.kube-system.svc.%s", clusterDomain))

	caCertificate, caPrivateKey, err := createCertificate("ca", nil, nil, false, "", nil, nil)
	if err != nil {
		return nil, nil, nil, err
	}
	apiServerCertificate, apiServerPrivateKey, err := createCertificate("apiserver", caCertificate, caPrivateKey, true, masterFQDN, extraFQDNs, extraIPs)
	if err != nil {
		return nil, nil, nil, err
	}
	clientCertificate, clientPrivateKey, err := createCertificate("client", caCertificate, caPrivateKey, false, "", nil, nil)
	if err != nil {
		return nil, nil, nil, err
	}

	return &PkiKeyCertPair{CertificatePem: string(certificateToPem(caCertificate.Raw)), PrivateKeyPem: string(privateKeyToPem(caPrivateKey))},
		&PkiKeyCertPair{CertificatePem: string(certificateToPem(apiServerCertificate.Raw)), PrivateKeyPem: string(privateKeyToPem(apiServerPrivateKey))},
		&PkiKeyCertPair{CertificatePem: string(certificateToPem(clientCertificate.Raw)), PrivateKeyPem: string(privateKeyToPem(clientPrivateKey))}, nil
}

func createCertificate(commonName string, caCertificate *x509.Certificate, caPrivateKey *rsa.PrivateKey, isServer bool, FQDN string, extraFQDNs []string, extraIPs []net.IP) (*x509.Certificate, *rsa.PrivateKey, error) {
	var err error

	isCA := (caCertificate == nil)

	now := time.Now()

	template := x509.Certificate{
		Subject:   pkix.Name{CommonName: commonName},
		NotBefore: now,
		NotAfter:  now.Add(ValidityDuration),

		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		BasicConstraintsValid: true,
	}

	if isCA {
		template.KeyUsage |= x509.KeyUsageCertSign
		template.IsCA = isCA
	} else if isServer {
		extraFQDNs = append(extraFQDNs, FQDN)
		extraIPs = append(extraIPs, net.ParseIP("10.3.0.1"))

		template.DNSNames = extraFQDNs
		template.IPAddresses = extraIPs
		template.ExtKeyUsage = append(template.ExtKeyUsage, x509.ExtKeyUsageServerAuth)
	} else {
		template.ExtKeyUsage = append(template.ExtKeyUsage, x509.ExtKeyUsageClientAuth)
	}

	snMax := new(big.Int).Lsh(big.NewInt(1), 128)
	template.SerialNumber, err = rand.Int(rand.Reader, snMax)
	if err != nil {
		return nil, nil, err
	}

	privateKey, err := rsa.GenerateKey(rand.Reader, PkiKeySize)

	var privateKeyToUse *rsa.PrivateKey
	var certificateToUse *x509.Certificate
	if !isCA {
		privateKeyToUse = caPrivateKey
		certificateToUse = caCertificate
	} else {
		privateKeyToUse = privateKey
		certificateToUse = &template
	}

	certDerBytes, err := x509.CreateCertificate(rand.Reader, &template, certificateToUse, &privateKey.PublicKey, privateKeyToUse)
	if err != nil {
		return nil, nil, err
	}

	certificate, err := x509.ParseCertificate(certDerBytes)
	if err != nil {
		return nil, nil, err
	}

	return certificate, privateKey, nil
}

func certificateToPem(derBytes []byte) []byte {
	pemBlock := &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: derBytes,
	}
	pemBuffer := bytes.Buffer{}
	pem.Encode(&pemBuffer, pemBlock)

	return pemBuffer.Bytes()
}

func privateKeyToPem(privateKey *rsa.PrivateKey) []byte {
	pemBlock := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}
	pemBuffer := bytes.Buffer{}
	pem.Encode(&pemBuffer, pemBlock)

	return pemBuffer.Bytes()
}
