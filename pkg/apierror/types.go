//------------------------------------------------------------
// Copyright (c) Microsoft Corporation.  All rights reserved.
//------------------------------------------------------------

package apierror

import "encoding/json"

// Error is the OData v4 format, used by the RPC and
// will go into the v2.2 Azure REST API guidelines
type Error struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
	Target  string    `json:"target,omitempty"`
	Details []Error   `json:"details,omitempty"`

	Category       ErrorCategory `json:"-"`
	ExceptionType  string        `json:"-"`
	InternalDetail string        `json:"-"`
}

// ErrorResponse  defines Resource Provider API 2.0 Error Response Content structure
type ErrorResponse struct {
	Body Error `json:"error"`
}

// Error implements error interface to return error in json
func (e *ErrorResponse) Error() string {
	return e.Body.Error()
}

// Error implements error interface to return error in json
func (e *Error) Error() string {
	output, err := json.MarshalIndent(e, " ", " ")
	if err != nil {
		return err.Error()
	}
	return string(output)
}
