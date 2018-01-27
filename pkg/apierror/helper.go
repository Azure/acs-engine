package apierror

import (
	"encoding/json"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/arm/resources/resources"
)

// ExtractCodeFromARMHttpResponse returns the ARM error's Code field
// If not found return defaultCode
func ExtractCodeFromARMHttpResponse(resp *http.Response, defaultCode ErrorCode) ErrorCode {
	if resp == nil {
		return defaultCode
	}
	decoder := json.NewDecoder(resp.Body)
	errorJSON := ErrorResponse{}
	if err := decoder.Decode(&errorJSON); err != nil {
		return defaultCode
	}

	if errorJSON.Body.Code == "" {
		return defaultCode
	}
	return ErrorCode(errorJSON.Body.Code)
}

//ConvertToAPIError turns a ManagementErrorWithDetails into a apierror.Error
func ConvertToAPIError(mError *resources.ManagementErrorWithDetails) *Error {
	retVal := &Error{}
	if mError.Code != nil {
		retVal.Code = ErrorCode(*mError.Code)
	}
	if mError.Message != nil {
		retVal.Message = *mError.Message
	}
	if mError.Target != nil {
		retVal.Target = *mError.Target
	}
	if mError.Details != nil {
		retVal.Details = []Error{}
		for _, me := range *mError.Details {
			retVal.Details = append(retVal.Details, *ConvertToAPIError(&me))
		}
	}
	return retVal
}
