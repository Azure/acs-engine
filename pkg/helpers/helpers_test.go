package helpers

import "testing"

func TestPointerToBool(t *testing.T) {
	boolVar := true
	ret := PointerToBool(boolVar)
	if *ret != boolVar {
		t.Fatalf("expected PointerToBool(true) to return *true, instead returned %#v", ret)
	}
}
