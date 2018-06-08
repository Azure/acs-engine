package helpers

import "testing"

func TestCloneAndConform(t *testing.T) {
	type a struct {
		S string `json:"s" conform:"redact"`
	}

	type b struct {
		A *a `json:"a"`
	}

	src := b{
		A: &a{
			S: "secret",
		},
	}
	var dst b
	err := CloneAndConform(&src, &dst)
	if err != nil {
		t.Fatalf("CloneAndConform should be error free")
	}

	if dst.A == nil {
		t.Fatalf("dst.A should not be nil")
	}

	if dst.A.S == src.A.S {
		t.Fatalf("dst.A.S should not be equal to src.A.S")
	}

	if dst.A.S != "REDACTED" {
		t.Fatalf("dst.A.S should be equal to REDACTED")
	}
}

func TestPointerToBool(t *testing.T) {
	boolVar := true
	ret := PointerToBool(boolVar)
	if *ret != boolVar {
		t.Fatalf("expected PointerToBool(true) to return *true, instead returned %#v", ret)
	}
}

func TestIsRegionNormalized(t *testing.T) {
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
