//------------------------------------------------------------
// Copyright (c) Microsoft Corporation.  All rights reserved.
//------------------------------------------------------------

package apierror

import (
	"testing"

	. "github.com/onsi/gomega"
)

func TestNewAPIError(t *testing.T) {
	RegisterTestingT(t)

	apiError := New(
		ClientError,
		InvalidParameter,
		"error test")

	Expect(apiError.Body.Code).Should(Equal(ErrorCode("InvalidParameter")))
}

func TestAcsNewAPIError(t *testing.T) {
	RegisterTestingT(t)

	apiError := New(
		ClientError,
		ScaleDownInternalError,
		"error test")

	Expect(apiError.Body.Code).Should(Equal(ErrorCode("ScaleDownInternalError")))
}
