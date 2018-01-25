package armhelpers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Azure/acs-engine/pkg/api"
	"github.com/Azure/acs-engine/pkg/apierror"
	"github.com/Azure/azure-sdk-for-go/arm/resources/resources"
	"github.com/sirupsen/logrus"
)

func parseDeploymentOperation(logger *logrus.Entry, operation resources.DeploymentOperation) (*apierror.Error, error) {
	if operation.Properties == nil || operation.Properties.StatusMessage == nil {
		return nil, fmt.Errorf("DeploymentOperation.Properties is not set")
	}
	b, err := json.MarshalIndent(operation.Properties.StatusMessage, "", "  ")
	if err != nil {
		logger.Errorf("Error occurred marshalling JSON: '%v'", err)
		return nil, err
	}
	return toError(logger, b)
}

func toError(logger *logrus.Entry, b []byte) (*apierror.Error, error) {
	errresp := &apierror.ErrorResponse{}

	if err := json.Unmarshal(b, errresp); err != nil {
		logger.Errorf("Error occurred unmarshalling JSON: '%v' JSON: '%s'", err, string(b))
		return nil, err
	}

	armError := &errresp.Body
	// If error code is ResourceDeploymentFailure then RP error is defined in the child object field: "details
	switch armError.Code {
	case apierror.ResourceDeploymentFailure,
		apierror.InvalidTemplateDeployment,
		apierror.DeploymentFailed:
		// StatusMessage.error.details array supports multiple errors but in this particular case
		// DeploymentOperationProperties contains error from one specific resource type so the
		// chances of multiple deployment errors being returned for a single resource type is slim
		// (but possible) based on current error/QoS analysis. In those cases where multiple errors
		// are returned ACS will pick the first error code for determining whether this is an internal
		// or a client error. This can be reevaluated later based on practical experience.
		// However, note that customer will be returned the entire contents of "StatusMessage" object
		// (like before) so they have access to all the errors returned by ARM.
		logger.Infof("Found %s error code - error response = '%v'", armError.Code, armError)
		if len(armError.Details) > 0 {
			armError = &armError.Details[0]
		}
	}
	armError.Category = getErrorCategory(armError.Code)
	return armError, nil
}

func getErrorCategory(code apierror.ErrorCode) apierror.ErrorCategory {
	switch code {
	case apierror.InvalidParameter,
		apierror.BadRequest,
		apierror.OperationNotAllowed,
		apierror.PropertyChangeNotAllowed,
		apierror.UnregisterWithResourcesNotAllowed,
		apierror.InvalidParameterConflictingProperties,
		apierror.SubscriptionNotRegistered,
		apierror.ConflictingUserInput,
		apierror.QuotaExceeded,
		apierror.Unauthorized,
		apierror.ResourcesOverConstrained:
		return apierror.ClientError
	default:
		return apierror.InternalError
	}
}

// DeployTemplateSync deploys the template and returns apierror
func DeployTemplateSync(az ACSEngineClient, logger *logrus.Entry, resourceGroupName, deploymentName string, template map[string]interface{}, parameters map[string]interface{}) *apierror.Error {
	depExt, depErr := az.DeployTemplate(resourceGroupName, deploymentName, template, parameters, nil)
	if depErr == nil {
		return nil
	}

	logger.Infof("Getting detailed deployment errors for %s", deploymentName)

	if depExt == nil {
		logger.Warn("DeploymentExtended is nil")
		return &apierror.Error{
			Code:     apierror.InternalOperationError,
			Message:  depErr.Error(),
			Category: apierror.InternalError}
	}

	// try to extract error from ARM Response
	var armErr *apierror.Error
	if depExt.Response.Response != nil && depExt.Body != nil {
		buf := new(bytes.Buffer)
		buf.ReadFrom(depExt.Body)
		logger.Infof("StatusCode: %d, Error: %s", depExt.Response.StatusCode, buf.String())
		if resp, err := toError(logger, buf.Bytes()); err == nil {
			switch {
			case depExt.Response.StatusCode < 500 && depExt.Response.StatusCode >= 400:
				resp.Category = apierror.ClientError
			case depExt.Response.StatusCode >= 500:
				resp.Category = apierror.InternalError
			}
			armErr = resp
		} else {
			logger.Errorf("unable to unmarshal response into apierror: %v", err)
		}
	} else {
		logger.Errorf("Got error from Azure SDK without response from ARM")
		// This is the failed sdk validation before calling ARM path
		return &apierror.Error{
			Code:     apierror.InternalOperationError,
			Message:  depErr.Error(),
			Category: apierror.InternalError}
	}

	// Check that ARM returned ErrorResponse
	if armErr == nil || len(armErr.Message) == 0 || len(armErr.Code) == 0 {
		logger.Warn("Not an ARM Response")
		return &apierror.Error{
			Code:     apierror.InternalOperationError,
			Message:  depErr.Error(),
			Category: apierror.InternalError}
	}

	if depExt.Properties == nil || depExt.Properties.ProvisioningState == nil {
		logger.Warn("No resources.DeploymentExtended.Properties")
		return armErr
	}
	properties := depExt.Properties

	switch *properties.ProvisioningState {
	case string(api.Canceled):
		logger.Warning("template deployment has been canceled")
		return &apierror.Error{
			Code:     apierror.ProvisioningFailed,
			Message:  "template deployment has been canceled",
			Category: apierror.ClientError}

	case string(api.Failed):
		var top int32 = 1
		results := make([]resources.DeploymentOperationsListResult, top)
		res, err := az.ListDeploymentOperations(resourceGroupName, deploymentName, &top)
		if err != nil {
			logger.Errorf("unable to list deployment operations %s. error: %v", deploymentName, err)
			return armErr
		}
		results[0] = res

		for res.NextLink != nil {
			res, err = az.ListDeploymentOperationsNextResults(res)
			if err != nil {
				logger.Warningf("unable to list next deployment operations %s. error: %v", deploymentName, err)
				break
			}

			results = append(results, res)
		}
		apierr, err := analyzeDeploymentResultAndSaveError(resourceGroupName, deploymentName, results, logger)
		if err != nil || apierr == nil {
			return armErr
		}
		return apierr

	default:
		logger.Warningf("Unexpected ProvisioningState %s", *properties.ProvisioningState)
		return armErr
	}
}

func analyzeDeploymentResultAndSaveError(resourceGroupName, deploymentName string,
	operationLists []resources.DeploymentOperationsListResult, logger *logrus.Entry) (*apierror.Error, error) {
	var apierr *apierror.Error
	var err error
	errs := []string{}
	isInternalErr := false

	for _, operationsList := range operationLists {
		if operationsList.Value == nil {
			continue
		}

		for _, operation := range *operationsList.Value {
			if operation.Properties == nil || *operation.Properties.ProvisioningState != string(api.Failed) {
				continue
			}

			// log the full deployment operation error response
			if operation.ID != nil && operation.OperationID != nil {
				b, _ := json.Marshal(operation.Properties)
				logger.Infof("deployment operation ID %s, operationID %s, prooperties: %s", *operation.ID, *operation.OperationID, b)
			} else {
				logger.Error("either deployment ID or operationID is nil")
			}

			apierr, err = parseDeploymentOperation(logger, operation)
			if err != nil {
				logger.Errorf("unable to convert deployment operation to error response in deployment %s from ARM. error: %v", deploymentName, err)
				return nil, err
			}
			if apierr.Category == apierror.InternalError {
				isInternalErr = true
			}
			errs = append(errs, apierr.Error())
		}
	}
	provisionErr := &apierror.Error{}
	if len(errs) > 0 {
		if isInternalErr {
			provisionErr.Category = apierror.InternalError
		} else {
			provisionErr.Category = apierror.ClientError
		}
		if len(errs) == 1 {
			provisionErr = apierr
		} else {
			provisionErr.Code = apierror.ProvisioningFailed
			provisionErr.Message = strings.Join(errs, "\n")
		}
		return provisionErr, nil
	}
	return nil, nil
}
