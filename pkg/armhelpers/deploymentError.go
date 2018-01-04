package armhelpers

import (
	"encoding/json"
	"fmt"

	"github.com/Azure/acs-engine/pkg/api"
	"github.com/Azure/azure-sdk-for-go/arm/resources/resources"
	"github.com/sirupsen/logrus"
)

// Error is the OData v4 format, used by the RPC and
// will go into the v2.2 Azure REST API guidelines
type Error struct {
	Code    string  `json:"code"`
	Message string  `json:"message"`
	Target  string  `json:"target,omitempty"`
	Details []Error `json:"details,omitempty"`
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

	// In some cases ARM returns  "ResourceDeploymentFailure" as the top error code and error
	// code returned by the resource provider is in the child array "details". The following code
	// covers both of those scenarios. Also note, that details array is recursive i.e. the objects
	// in details array can contain another details array with no max depth of this tree enforced
	// by ARM. The following code only evaluates the first level property of "StatusMessage": "error"
	// and its child array: "details" for error codes.

	// If error code is ResourceDeploymentFailure then RP error is defined in the child object field: "details
	if errresp.Body.Code == "ResourceDeploymentFailure" {
		// StatusMessage.error.details array supports multiple errors but in this particular case
		// DeploymentOperationProperties contains error from one specific resource type so the
		// chances of multiple deployment errors being returned for a single resource type is slim
		// (but possible) based on current error/QoS analysis. In those cases where multiple errors
		// are returned ACS will pick the first error code for determining whether this is an internal
		// or a client error. This can be reevaluated later based on practical experience.
		// However, note that customer will be returned the entire contents of "StatusMessage" object
		// (like before) so they have access to all the errors returned by ARM.
		logger.Infof("Found ResourceDeploymentFailure error code - error response = '%v'", *errresp)
		details := errresp.Body.Details
		if len(details) > 0 {
			errresp.Body.Code = details[0].Code
			errresp.Body.Message = details[0].Message
			errresp.Body.Target = details[0].Target
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
		if operation.Properties == nil || *operation.Properties.ProvisioningState != string(api.Failed) {
			continue
		}

		errresp, err := toArmError(logger, operation)
		if err != nil {
			logger.Errorf("unable to convert deployment operation to error response in deployment %s from ARM. error: %v", deploymentName, err)
			return ret, err
		}

		if len(errresp.Body.Code) > 0 {
			logger.Warnf("got failed deployment operation in deployment %s. error: %v", deploymentName, errresp.Error())
		}
		ret = append(ret, errresp)
	}
	return ret, nil
}

// GetDeploymentError returns deployment error
func GetDeploymentError(az ACSEngineClient, logger *logrus.Entry, resourceGroupName, deploymentName string) (*ErrorResponse, error) {
	errList := []*ErrorResponse{}
	logger.Infof("Getting detailed deployment errors for %s", deploymentName)

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
	errList = append(errList, eList...)
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
		errList = append(errList, eList...)
	}
	switch len(errList) {
	case 0:
		return nil, fmt.Errorf("unable to extract deployment error for %s", deploymentName)
	case 1:
		return errList[0], nil
	default:
		// combine all errors
		for _, e := range errList[1:] {
			errList[0].Body.Details = append(errList[0].Body.Details, e.Body.Details...)
		}
		return errList[0], nil
	}
}
