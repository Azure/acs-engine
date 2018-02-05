package armhelpers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Azure/acs-engine/pkg/api"
	"github.com/Azure/azure-sdk-for-go/arm/resources/resources"
	"github.com/sirupsen/logrus"
)

// ArmError is the OData v4 format, used by the RPC and
// will go into the v2.2 Azure REST API guidelines
type ArmError struct {
	Code    string     `json:"code"`
	Message string     `json:"message"`
	Target  string     `json:"target,omitempty"`
	Details []ArmError `json:"details,omitempty"`
}

// ArmErrorResponse  defines Resource Provider API 2.0 Error Response Content structure
type ArmErrorResponse struct {
	Body ArmError `json:"error"`
}

// Error implements error interface to return error in json
func (e *ArmError) Error() string {
	output, err := json.MarshalIndent(e, " ", " ")
	if err != nil {
		return err.Error()
	}
	return string(output)
}

// Error implements error interface to return error in json
func (e *ArmErrorResponse) Error() string {
	return e.Body.Error()
}

func parseDeploymentOperation(logger *logrus.Entry, operation resources.DeploymentOperation) (*ArmError, error) {
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

func toError(logger *logrus.Entry, b []byte) (*ArmError, error) {
	errresp := &ArmErrorResponse{}

	if err := json.Unmarshal(b, errresp); err != nil {
		logger.Errorf("Error occurred unmarshalling JSON: '%v' JSON: '%s'", err, string(b))
		return nil, err
	}
	return &errresp.Body, nil
}

// DeployTemplateSync deploys the template and returns ArmError
func DeployTemplateSync(az ACSEngineClient, logger *logrus.Entry, resourceGroupName, deploymentName string, template map[string]interface{}, parameters map[string]interface{}) error {
	depExt, depErr := az.DeployTemplate(resourceGroupName, deploymentName, template, parameters, nil)
	if depErr == nil {
		return nil
	}

	logger.Infof("Getting detailed deployment errors for %s", deploymentName)

	if depExt == nil {
		logger.Warn("DeploymentExtended is nil")
		return &ArmError{
			Code:    "DeploymentFailed",
			Message: depErr.Error()}
	}

	// try to extract error from ARM Response
	var armErr *ArmError
	if depExt.Response.Response != nil && depExt.Body != nil {
		buf := new(bytes.Buffer)
		buf.ReadFrom(depExt.Body)
		logger.Infof("StatusCode: %d, Error: %s", depExt.Response.StatusCode, buf.String())
		if resp, err := toError(logger, buf.Bytes()); err == nil {
			armErr = resp
		} else {
			logger.Errorf("unable to unmarshal response into apierror: %v", err)
		}
	} else {
		logger.Errorf("Got error from Azure SDK without response from ARM")
		// This is the failed sdk validation before calling ARM path
		return &ArmError{
			Code:    "DeploymentFailed",
			Message: depErr.Error()}
	}

	// Check that ARM returned an error
	if armErr == nil || len(armErr.Message) == 0 || len(armErr.Code) == 0 {
		logger.Warn("Not an ARM Error")
		return &ArmError{
			Code:    "DeploymentFailed",
			Message: depErr.Error()}
	}

	if depExt.Properties == nil || depExt.Properties.ProvisioningState == nil {
		logger.Warn("No resources.DeploymentExtended.Properties")
		return armErr
	}
	properties := depExt.Properties

	switch *properties.ProvisioningState {
	case string(api.Canceled):
		logger.Warning("template deployment has been canceled")
		return &ArmError{
			Code:    "DeploymentFailed",
			Message: "template deployment has been canceled"}

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
	operationLists []resources.DeploymentOperationsListResult, logger *logrus.Entry) (*ArmError, error) {
	var apierr *ArmError
	var err error
	errs := []string{}

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
			errs = append(errs, apierr.Error())
		}
	}

	if len(errs) > 0 {
		provisionErr := &ArmError{}
		if len(errs) == 1 {
			provisionErr = apierr
		} else {
			provisionErr.Code = "ProvisioningFailed"
			provisionErr.Message = strings.Join(errs, "\n")
		}
		return provisionErr, nil
	}
	return nil, nil
}
