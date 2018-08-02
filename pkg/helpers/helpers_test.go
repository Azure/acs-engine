package helpers

import (
	"bytes"
	"crypto/x509"
	"encoding/pem"
	"math/rand"
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

	expectedPublicKeyString := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQCyx5MHXjJvJAx5DJ9FZNIDa/QTWorSF+Ra21Tz49DQWfdSESnCGFFVBh/MQUFGv5kCenbmqEjsWF177kFOdv1vOTz4sKRlHg7u3I9uCyyZQrWx4X4RdNk7eX+isQVjFXYw2W1rRDUrnK/82qVTv1f0gu1DV4Z7GoIa2jfJ0zBUY3IW0VN9jYaPVuwv4t5y2GwSZF+HBRuOfLfiUgt4+qVFOz4KwRaEBsVfWxlidlT3K3/+ztWpFOmaKIOjQreEWV10ZSo3f9g6j/HdMPtwYvRCtYStbFCRmcbPr9nuR84SAX/4f95KvBAKLnXwb5Bt71D2vAlZSW1Ylv2VbcaZ73+43EpyphYCSg3kOCdwsqE/EU+Swued82SguLALD3mNKbxHGJppFjz3GMyPpJuSH5EE1OANyPxABCwCYycKiNWbOPi3l6o4tMrASYRXi8l3l9JCvioUJ3bXXH6cDpcP4P6QgsuxhwVkUiECU+dbjJXK4gAUVuWKkMOdY7ITh82oU3wOWXbk8K3bdIUp2ylcHeAd2pekGMuaEKGbrXGRiBitCEjl67Bj5opQflgSmI63g8Sa3mKOPGRYMI5MXHMVj4Rns5JFHoENuImrlvrbLv3izAwO61vgN7iK26BwzO7jz92fNOHGviejNWYJyi4vZlq07153NZXP8D2xYTebh9hwHQ==\n"

	expectedPrivateKeyString := `-----BEGIN RSA PRIVATE KEY-----
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
-----END RSA PRIVATE KEY-----
`

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

	if string(pemBuffer.Bytes()) != expectedPrivateKeyString {
		t.Fatalf("Private Key did not match expected format/value")
	}

	if publicKey != expectedPublicKeyString {
		t.Fatalf("Public Key did not match expected format/value")
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
