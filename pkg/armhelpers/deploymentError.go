package armhelpers

import (
	"encoding/json"
	"strings"

	"github.com/Azure/acs-engine/pkg/api"
	"github.com/Azure/azure-sdk-for-go/arm/resources/resources"
	"github.com/sirupsen/logrus"
)

//TODO move pkg/core/apierror/ to acs-engine

// ErrorCategory indicates the kind of error
type ErrorCategory string

const (
	// ClientError is expected error
	ClientError ErrorCategory = "ClientError"

	// InternalError is system or internal error
	InternalError ErrorCategory = "InternalError"
)

// Error is the OData v4 format, used by the RPC and
// will go into the v2.2 Azure REST API guidelines
type Error struct {
	Code    string  `json:"code"`
	Message string  `json:"message"`
	Target  string  `json:"target,omitempty"`
	Details []Error `json:"details,omitempty"`

	Category ErrorCategory `json:"-"`
}

// ErrorResponse defines Resource Provider API 2.0 Error Response Content structure
type ErrorResponse struct {
	Body Error `json:"error"`
}

// DeploymentError defines deployment error along with deployment operation errors
type DeploymentError struct {
	RootError       error
	OperationErrors []*ErrorResponse
}

// Error implements error interface to return error in json
func (e *DeploymentError) Error() string {
	if len(e.OperationErrors) == 0 {
		return e.RootError.Error()
	}
	errStrList := make([]string, len(e.OperationErrors)+1)
	errStrList[0] = e.RootError.Error()
	for i, errResp := range e.OperationErrors {
		errStrList[i+1] = errResp.Error()
	}
	return strings.Join(errStrList, " | ")
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

func toArmError(logger *logrus.Entry, operation resources.DeploymentOperation) (*ErrorResponse, error) {
	errresp := &ErrorResponse{}
	if operation.Properties != nil && operation.Properties.StatusMessage != nil {
		b, err := json.MarshalIndent(operation.Properties.StatusMessage, "", "  ")
		if err != nil {
			logger.Errorf("Error occurred marshalling JSON: '%v'", err)
			return nil, err
		}
		if err := json.Unmarshal(b, errresp); err != nil {
			logger.Errorf("Error occurred unmarshalling JSON: '%v' JSON: '%s'", err, string(b))
			return nil, err
		}
	}
	return errresp, nil
}

func toArmErrors(logger *logrus.Entry, deploymentName string, operationsList resources.DeploymentOperationsListResult) ([]*ErrorResponse, error) {
	ret := []*ErrorResponse{}

	if operationsList.Value == nil {
		return ret, nil
	}

	for _, operation := range *operationsList.Value {
		if operation.Properties == nil || operation.Properties.ProvisioningState == nil || *operation.Properties.ProvisioningState != string(api.Failed) {
			continue
		}

		errresp, err := toArmError(logger, operation)
		if err != nil {
			logger.Warnf("unable to convert deployment operation to error response in deployment %s from ARM. error: %v", deploymentName, err)
			continue
		}

		if len(errresp.Body.Code) > 0 {
			logger.Warnf("got failed deployment operation in deployment %s. error: %v", deploymentName, errresp.Error())
		}
		ret = append(ret, errresp)
	}
	return ret, nil
}

//TODO errorCode is ErrorCode
func newErrorResponse(errorCategory ErrorCategory, errorCode string, message string) *ErrorResponse {
	return &ErrorResponse{
		Body: Error{
			Code:     errorCode,
			Message:  message,
			Category: errorCategory,
		},
	}
}

// GetDeploymentError returns deployment error
func GetDeploymentError(res *resources.DeploymentExtended, rootError error, az ACSEngineClient, logger *logrus.Entry, resourceGroupName, deploymentName string) (*DeploymentError, error) {
	if rootError == nil {
		return nil, nil
	}
	logger.Infof("Getting detailed deployment errors for %s", deploymentName)
	deploymentError := &DeploymentError{RootError: rootError}

	if res != nil && res.Response.Response != nil && res.Body != nil {
		armErr := &ErrorResponse{}
		if d := json.NewDecoder(res.Body); d != nil {
			if err := d.Decode(armErr); err == nil {
				logger.Errorf("StatusCode: %d, ErrorCode: %s, ErrorMessage: %s", res.Response.StatusCode, armErr.Body.Code, armErr.Body.Message)
				deploymentError.OperationErrors = append(deploymentError.OperationErrors, armErr)
				switch {
				case res.Response.StatusCode < 500 && res.Response.StatusCode >= 400:
					armErr.Body.Category = ClientError
					return deploymentError, nil
				case res.Response.StatusCode >= 500:
					armErr.Body.Category = InternalError
					return deploymentError, nil
				}
			} else {
				logger.Errorf("unable to unmarshal response into apierror: %v", err)
			}
		}
	} else {
		logger.Errorf("Got error from Azure SDK without response from ARM, error: %v", rootError)
		// This is the failed sdk validation before calling ARM path
		deploymentError.OperationErrors = append(deploymentError.OperationErrors, newErrorResponse(InternalError, "InternalOperationError", rootError.Error()))
		return deploymentError, nil
	}

	var top int32 = 1
	operationList, err := az.ListDeploymentOperations(resourceGroupName, deploymentName, &top)
	if err != nil {
		logger.Warnf("unable to list deployment operations: %v", err)
		return nil, err
	}
	eList, err := toArmErrors(logger, deploymentName, operationList)
	if err != nil {
		return nil, err
	}
	deploymentError.OperationErrors = append(deploymentError.OperationErrors, eList...)
	for operationList.NextLink != nil {
		operationList, err = az.ListDeploymentOperationsNextResults(operationList)
		if err != nil {
			logger.Warnf("unable to list next deployment operations: %v", err)
			break
		}
		eList, err := toArmErrors(logger, deploymentName, operationList)
		if err != nil {
			return nil, err
		}
		deploymentError.OperationErrors = append(deploymentError.OperationErrors, eList...)
	}
	return deploymentError, nil
}
