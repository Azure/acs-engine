package helpers

import "testing"

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
