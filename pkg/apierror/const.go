//------------------------------------------------------------
// Copyright (c) Microsoft Corporation.  All rights reserved.
//------------------------------------------------------------

package apierror

// ErrorCategory indicates the kind of error
type ErrorCategory string

const (
	// ClientError is expected error
	ClientError ErrorCategory = "ClientError"

	// InternalError is system or internal error
	InternalError ErrorCategory = "InternalError"
)

// Common Azure Resource Provider API error code
type ErrorCode string

const (
	// From Microsoft.Azure.ResourceProvider.API.ErrorCode

	// InvalidParameter error
	InvalidParameter ErrorCode = "InvalidParameter"
	// BadRequest error
	BadRequest ErrorCode = "BadRequest"
	// NotFound error
	NotFound ErrorCode = "NotFound"
	// Conflict error
	Conflict ErrorCode = "Conflict"
	// PreconditionFailed error
	PreconditionFailed ErrorCode = "PreconditionFailed"
	// OperationNotAllowed error
	OperationNotAllowed ErrorCode = "OperationNotAllowed"
	// OperationPreempted error
	OperationPreempted ErrorCode = "OperationPreempted"
	// PropertyChangeNotAllowed error
	PropertyChangeNotAllowed ErrorCode = "PropertyChangeNotAllowed"
	// InternalOperationError error
	InternalOperationError ErrorCode = "InternalOperationError"
	// InvalidSubscriptionStateTransition error
	InvalidSubscriptionStateTransition ErrorCode = "InvalidSubscriptionStateTransition"
	// UnregisterWithResourcesNotAllowed error
	UnregisterWithResourcesNotAllowed ErrorCode = "UnregisterWithResourcesNotAllowed"
	// InvalidParameterConflictingProperties error
	InvalidParameterConflictingProperties ErrorCode = "InvalidParameterConflictingProperties"
	// SubscriptionNotRegistered error
	SubscriptionNotRegistered ErrorCode = "SubscriptionNotRegistered"
	// ConflictingUserInput error
	ConflictingUserInput ErrorCode = "ConflictingUserInput"
	// ProvisioningInternalError error
	ProvisioningInternalError ErrorCode = "ProvisioningInternalError"
	// ProvisioningFailed error
	ProvisioningFailed ErrorCode = "ProvisioningFailed"
	// NetworkingInternalOperationError error
	NetworkingInternalOperationError ErrorCode = "NetworkingInternalOperationError"
	// QuotaExceeded error
	QuotaExceeded ErrorCode = "QuotaExceeded"
	// Unauthorized error
	Unauthorized ErrorCode = "Unauthorized"
	// ResourcesOverConstrained error
	ResourcesOverConstrained ErrorCode = "ResourcesOverConstrained"

	// ResourceDeploymentFailure error
	ResourceDeploymentFailure ErrorCode = "ResourceDeploymentFailure"
	// InvalidTemplateDeployment error
	InvalidTemplateDeployment ErrorCode = "InvalidTemplateDeployment"
	// DeploymentFailed error
	DeploymentFailed ErrorCode = "DeploymentFailed"
)
