package acsengine

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"fmt"
	"math/big"
	"net"
	"os"
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

func CreatePki(extraFQDNs []string, extraIPs []net.IP, clusterDomain string, caPair *PkiKeyCertPair) (*PkiKeyCertPair, *PkiKeyCertPair, *PkiKeyCertPair, error) {
	start := time.Now()
	defer func(s time.Time) {
		fmt.Fprintf(os.Stderr, "cert creation took %s\n", time.Since(s))
	}(start)
	extraFQDNs = append(extraFQDNs, fmt.Sprintf("kubernetes"))
	extraFQDNs = append(extraFQDNs, fmt.Sprintf("kubernetes.default"))
	extraFQDNs = append(extraFQDNs, fmt.Sprintf("kubernetes.default.svc"))
	extraFQDNs = append(extraFQDNs, fmt.Sprintf("kubernetes.default.svc.%s", clusterDomain))
	extraFQDNs = append(extraFQDNs, fmt.Sprintf("kubernetes.kube-system"))
	extraFQDNs = append(extraFQDNs, fmt.Sprintf("kubernetes.kube-system.svc"))
	extraFQDNs = append(extraFQDNs, fmt.Sprintf("kubernetes.kube-system.svc.%s", clusterDomain))

	var (
		caCertificate         *x509.Certificate
		caPrivateKey          *rsa.PrivateKey
		apiServerCertificate  *x509.Certificate
		apiServerPrivateKey   *rsa.PrivateKey
		clientCertificate     *x509.Certificate
		clientPrivateKey      *rsa.PrivateKey
		kubeConfigCertificate *x509.Certificate
		kubeConfigPrivateKey  *rsa.PrivateKey
	)
	errors := make(chan error)

	var err error
	caCertificate, err = pemToCertificate(caPair.CertificatePem)
	if err != nil {
		return nil, nil, nil, err
	}
	caPrivateKey, err = pemToKey(caPair.PrivateKeyPem)
	if err != nil {
		return nil, nil, nil, err
	}

	go func() {
		var err error
		apiServerCertificate, apiServerPrivateKey, err = createCertificate("apiserver", caCertificate, caPrivateKey, true, extraFQDNs, extraIPs)
		errors <- err
	}()

	go func() {
		var err error
		clientCertificate, clientPrivateKey, err = createCertificate("client", caCertificate, caPrivateKey, false, nil, nil)
		errors <- err
	}()

	go func() {
		var err error
		kubeConfigCertificate, kubeConfigPrivateKey, err = createCertificate("client", caCertificate, caPrivateKey, false, nil, nil)
		errors <- err
	}()

	e1 := <-errors
	e2 := <-errors
	e3 := <-errors
	if e1 != nil {
		return nil, nil, nil, e1
	}
	if e2 != nil {
		return nil, nil, nil, e2
	}
	if e3 != nil {
		return nil, nil, nil, e2
	}

	return &PkiKeyCertPair{CertificatePem: string(certificateToPem(apiServerCertificate.Raw)), PrivateKeyPem: string(privateKeyToPem(apiServerPrivateKey))},
		&PkiKeyCertPair{CertificatePem: string(certificateToPem(clientCertificate.Raw)), PrivateKeyPem: string(privateKeyToPem(clientPrivateKey))},
		&PkiKeyCertPair{CertificatePem: string(certificateToPem(kubeConfigCertificate.Raw)), PrivateKeyPem: string(privateKeyToPem(kubeConfigPrivateKey))},
		nil
}

func createCertificate(commonName string, caCertificate *x509.Certificate, caPrivateKey *rsa.PrivateKey, isServer bool, extraFQDNs []string, extraIPs []net.IP) (*x509.Certificate, *rsa.PrivateKey, error) {
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

func pemToCertificate(raw string) (*x509.Certificate, error) {
	cpb, _ := pem.Decode([]byte(raw))
	if cpb == nil {
		return nil, errors.New("The raw pem is not a valid PEM formatted block.")
	}
	return x509.ParseCertificate(cpb.Bytes)
}

func pemToKey(raw string) (*rsa.PrivateKey, error) {
	kpb, _ := pem.Decode([]byte(raw))
	if kpb == nil {
		return nil, errors.New("The raw pem is not a valid PEM formatted block.")
	}
	return x509.ParsePKCS1PrivateKey(kpb.Bytes)
}
