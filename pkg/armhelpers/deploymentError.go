package armhelpers

import (
	"encoding/json"
	"strings"

	"github.com/Azure/acs-engine/pkg/api"
	"github.com/Azure/acs-engine/pkg/apierror"
	"github.com/Azure/azure-sdk-for-go/arm/resources/resources"
	"github.com/sirupsen/logrus"
)

func toErrorResponse(logger *logrus.Entry, operation resources.DeploymentOperation) (*apierror.ErrorResponse, error) {
	errresp := &apierror.ErrorResponse{}
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
	errresp.Body.Category = getErrorCategory(errresp.Body.Code)
	return errresp, nil
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

// GetDeploymentError returns deployment error
func GetDeploymentError(res *resources.DeploymentExtended, az ACSEngineClient, logger *logrus.Entry, resourceGroupName, deploymentName string) (*apierror.Error, error) {
	logger.Infof("Getting detailed deployment errors for %s", deploymentName)

	if res == nil || res.Properties == nil || res.Properties.ProvisioningState == nil {
		return nil, nil
	}
	properties := res.Properties

	switch *properties.ProvisioningState {
	case string(api.Canceled):
		logger.Warning("template deployment has been canceled")
		return &apierror.Error{
			Code:     apierror.ProvisioningFailed,
			Message:  "template deployment has been canceled",
			Category: apierror.ClientError}, nil

	case string(api.Failed):
		var top int32 = 1
		results := make([]resources.DeploymentOperationsListResult, top)
		res, err := az.ListDeploymentOperations(resourceGroupName, deploymentName, &top)
		if err != nil {
			logger.Errorf("unable to list deployment operations %s. error: %v", deploymentName, err)
			return nil, err
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
		return analyzeDeploymentResultAndSaveError(resourceGroupName, deploymentName, results, logger)

	default:
		return nil, nil
	}
}

func analyzeDeploymentResultAndSaveError(resourceGroupName, deploymentName string,
	operationLists []resources.DeploymentOperationsListResult, logger *logrus.Entry) (*apierror.Error, error) {
	var errresp *apierror.ErrorResponse
	var err error
	errs := []string{}
	isInternalErr := false
	failedCnt := 0
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

			failedCnt++
			errresp, err = toErrorResponse(logger, operation)
			if err != nil {
				logger.Errorf("unable to convert deployment operation to error response in deployment %s from ARM. error: %v", deploymentName, err)
				return nil, err
			}
			if errresp.Body.Category == apierror.InternalError {
				isInternalErr = true
			}
			errs = append(errs, errresp.Error())
		}
	}
	provisionErr := &apierror.Error{}
	if failedCnt > 0 {
		if isInternalErr {
			provisionErr.Category = apierror.InternalError
		} else {
			provisionErr.Category = apierror.ClientError
		}
		if failedCnt == 1 {
			provisionErr = &errresp.Body
		} else {
			provisionErr.Code = apierror.ProvisioningFailed
			provisionErr.Message = strings.Join(errs, "\n")
		}
		return provisionErr, nil
	}

	return nil, nil
}
