//------------------------------------------------------------
// Copyright (c) Microsoft Corporation.  All rights reserved.
//------------------------------------------------------------

package apierror

// New creates an ErrorResponse
func New(errorCategory ErrorCategory, errorCode ErrorCode, message string) *ErrorResponse {
	return &ErrorResponse{
		Body: Error{
			Code:     errorCode,
			Message:  message,
			Category: errorCategory,
		},
	}
}
