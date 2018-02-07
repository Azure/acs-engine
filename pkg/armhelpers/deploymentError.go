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

// DeploymentError contains the root deployment error along with deployment operation errors
type DeploymentError struct {
	TopError          error
	StatusCode        int
	Response          []byte
	ProvisioningState string
	Operations        [][]byte
}

// Error implements error interface
func (e *DeploymentError) Error() string {
	var str string
	if e.TopError != nil {
		str = e.TopError.Error()
	}
	var ops []string
	for _, b := range e.Operations {
		ops = append(ops, string(b))
	}
	return fmt.Sprintf("TopError[%s] StatusCode[%d] Response[%s] ProvisioningState[%s] Operations[%s]",
		str, e.StatusCode, e.Response, e.ProvisioningState, strings.Join(ops, " | "))
}

// DeployTemplateSync deploys the template and returns ArmError
func DeployTemplateSync(az ACSEngineClient, logger *logrus.Entry, resourceGroupName, deploymentName string, template map[string]interface{}, parameters map[string]interface{}) error {
	deploymentExtended, err := az.DeployTemplate(resourceGroupName, deploymentName, template, parameters, nil)
	if err == nil {
		return nil
	}

	logger.Infof("Getting detailed deployment errors for %s", deploymentName)
	deploymentErr := &DeploymentError{}
	deploymentErr.TopError = err

	if deploymentExtended == nil {
		logger.Warn("DeploymentExtended is nil")
		return deploymentErr
	}

	// try to extract error from ARM Response
	if deploymentExtended.Response.Response != nil && deploymentExtended.Body != nil {
		buf := new(bytes.Buffer)
		buf.ReadFrom(deploymentExtended.Body)
		logger.Infof("StatusCode: %d, Error: %s", deploymentExtended.Response.StatusCode, buf.String())
		deploymentErr.Response = buf.Bytes()
		deploymentErr.StatusCode = deploymentExtended.Response.StatusCode
	} else {
		logger.Errorf("Got error from Azure SDK without response from ARM")
		// This is the failed sdk validation before calling ARM path
		return deploymentErr
	}

	if deploymentExtended.Properties == nil || deploymentExtended.Properties.ProvisioningState == nil {
		logger.Warn("No resources.DeploymentExtended.Properties")
		return deploymentErr
	}
	properties := deploymentExtended.Properties
	deploymentErr.ProvisioningState = *properties.ProvisioningState

	var top int32 = 1
	res, err := az.ListDeploymentOperations(resourceGroupName, deploymentName, &top)
	if err != nil {
		logger.Errorf("unable to list deployment operations %s. error: %v", deploymentName, err)
		return deploymentErr
	}
	getOperationError(res, deploymentErr, logger)

	for res.NextLink != nil {
		res, err = az.ListDeploymentOperationsNextResults(res)
		if err != nil {
			logger.Warningf("unable to list next deployment operations %s. error: %v", deploymentName, err)
			break
		}
		getOperationError(res, deploymentErr, logger)
	}
	return deploymentErr
}

func getOperationError(operationsList resources.DeploymentOperationsListResult, deploymentErr *DeploymentError, logger *logrus.Entry) {
	if operationsList.Value == nil {
		return
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

		if operation.Properties == nil || operation.Properties.StatusMessage == nil {
			logger.Error("DeploymentOperation.Properties is not set")
			continue
		}
		b, err := json.MarshalIndent(operation.Properties.StatusMessage, "", "  ")
		if err != nil {
			logger.Errorf("Error occurred marshalling JSON: '%v'", err)
			continue
		}
		deploymentErr.Operations = append(deploymentErr.Operations, b)
	}
}
