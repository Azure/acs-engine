// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

package api

import (
	"reflect"
	"testing"
)

// testContainerService and testOpenShiftCluster are defined in
// converterfromosaapi_test.go.

func TestConvertContainerServiceToVLabsOpenShiftCluster(t *testing.T) {
	oc := ConvertContainerServiceToVLabsOpenShiftCluster(testContainerService)
	if !reflect.DeepEqual(oc, testOpenShiftCluster) {
		t.Errorf("ConvertContainerServiceToVLabsOpenShiftCluster returned unexpected result\n%#v\n", oc)
	}
}
