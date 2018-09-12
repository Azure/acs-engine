package helpers

import (
	"bytes"
	"crypto/x509"
	"encoding/pem"
	"math/rand"
	"strings"
	"testing"

	"github.com/Azure/acs-engine/pkg/i18n"
	"github.com/pkg/errors"
)

type ContainerService struct {
	ID       string `json:"id"`
	Location string `json:"location"`
	Name     string `json:"name"`
}

func TestJSONMarshal(t *testing.T) {
	input := &ContainerService{}
	result, _ := JSONMarshal(input, false)
	expected := "{\"id\":\"\",\"location\":\"\",\"name\":\"\"}\n"
	if string(result) != expected {
		t.Fatalf("JSONMarshal returned unexpected result: expected %s but got %s", expected, string(result))
	}
	result, _ = JSONMarshalIndent(input, "", "", false)
	expected = "{\n\"id\": \"\",\n\"location\": \"\",\n\"name\": \"\"\n}\n"
	if string(result) != expected {
		t.Fatalf("JSONMarshal returned unexpected result: expected \n%sbut got \n%s", expected, result)
	}
}

func TestNormalizeAzureRegion(t *testing.T) {
	cases := []struct {
		input          string
		expectedResult string
	}{
		{
			input:          "westus",
			expectedResult: "westus",
		},
		{
			input:          "West US",
			expectedResult: "westus",
		},
		{
			input:          "Eastern Africa",
			expectedResult: "easternafrica",
		},
		{
			input:          "",
			expectedResult: "",
		},
	}

	for _, c := range cases {
		result := NormalizeAzureRegion(c.input)
		if c.expectedResult != result {
			t.Fatalf("NormalizeAzureRegion returned unexpected result: expected %s but got %s", c.expectedResult, result)
		}
	}
}

func TestPointerToBool(t *testing.T) {
	boolVar := true
	ret := PointerToBool(boolVar)
	if *ret != boolVar {
		t.Fatalf("expected PointerToBool(true) to return *true, instead returned %#v", ret)
	}

	if IsTrueBoolPointer(ret) != boolVar {
		t.Fatalf("expected IsTrueBoolPointer(*true) to return true, instead returned %#v", IsTrueBoolPointer(ret))
	}

	boolVar = false
	ret = PointerToBool(boolVar)
	if *ret != boolVar {
		t.Fatalf("expected PointerToBool(false) to return *false, instead returned %#v", ret)
	}

	if IsTrueBoolPointer(ret) != boolVar {
		t.Fatalf("expected IsTrueBoolPointer(*false) to return false, instead returned %#v", IsTrueBoolPointer(ret))
	}
}

func TestPointerToInt(t *testing.T) {
	int1 := 1
	int2 := 2
	ret1 := PointerToInt(int1)
	if *ret1 != int1 {
		t.Fatalf("expected PointerToInt(1) to return *1, instead returned %#v", ret1)
	}
	ret2 := PointerToInt(int2)
	if *ret2 != int2 {
		t.Fatalf("expected PointerToInt(2) to return *2, instead returned %#v", ret2)
	}

	if *ret2 <= *ret1 {
		t.Fatalf("Pointers to ints messed up their values and made 2 <= 1")
	}
}

func TestCreateSSH(t *testing.T) {
	rg := rand.New(rand.NewSource(42))

	translator := &i18n.Translator{
		Locale: nil,
	}

	privateKey, publicKey, err := CreateSSH(rg, translator)
	if err != nil {
		t.Fatalf("failed to generate SSH: %s", err)
	}
	pemBlock := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}
	pemBuffer := bytes.Buffer{}
	pem.Encode(&pemBuffer, pemBlock)

	if !strings.HasPrefix(pemBuffer.String(), "-----BEGIN RSA PRIVATE KEY-----") {
		t.Fatalf("Private Key did not start with expected header")
	}

	if privateKey.N.BitLen() != SSHKeySize {
		t.Fatalf("Private Key was of length %d but %d was expected", privateKey.N.BitLen(), SSHKeySize)
	}

	if err := privateKey.Validate(); err != nil {
		t.Fatalf("Private Key failed validation: %v", err)
	}

	if !strings.HasPrefix(publicKey, "ssh-rsa ") {
		t.Fatalf("Public Key did not start with expected header")
	}
}

func TestAcceleratedNetworkingSupported(t *testing.T) {
	cases := []struct {
		input          string
		expectedResult bool
	}{
		{
			input:          "Standard_A1",
			expectedResult: false,
		},
		{
			input:          "Standard_G4",
			expectedResult: false,
		},
		{
			input:          "Standard_B3",
			expectedResult: false,
		},
		{
			input:          "Standard_D1_v2",
			expectedResult: false,
		},
		{
			input:          "Standard_L3",
			expectedResult: false,
		},
		{
			input:          "Standard_NC6",
			expectedResult: false,
		},
		{
			input:          "Standard_G4",
			expectedResult: false,
		},
		{
			input:          "Standard_D2_v2",
			expectedResult: true,
		},
		{
			input:          "Standard_DS2_v2",
			expectedResult: true,
		},
		{
			input:          "",
			expectedResult: false,
		},
	}

	for _, c := range cases {
		result := AcceleratedNetworkingSupported(c.input)
		if c.expectedResult != result {
			t.Fatalf("AcceleratedNetworkingSupported returned unexpected result for %s: expected %t but got %t", c.input, c.expectedResult, result)
		}
	}
}

func TestEqualError(t *testing.T) {
	testcases := []struct {
		errA     error
		errB     error
		expected bool
	}{
		{
			errA:     nil,
			errB:     nil,
			expected: true,
		},
		{
			errA:     errors.New("sample error"),
			errB:     nil,
			expected: false,
		},
		{
			errA:     nil,
			errB:     errors.New("sample error"),
			expected: false,
		},
		{
			errA:     errors.New("sample error"),
			errB:     errors.New("sample error"),
			expected: true,
		},
		{
			errA:     errors.New("sample error 1"),
			errB:     errors.New("sample error 2"),
			expected: false,
		},
	}

	for _, test := range testcases {
		if EqualError(test.errA, test.errB) != test.expected {
			t.Errorf("expected EqualError to return %t for errors %s and %s", test.expected, test.errA, test.errB)
		}
	}
}
