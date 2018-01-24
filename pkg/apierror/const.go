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
	InvalidParameter                      ErrorCode = "InvalidParameter"
	BadRequest                            ErrorCode = "BadRequest"
	NotFound                              ErrorCode = "NotFound"
	Conflict                              ErrorCode = "Conflict"
	PreconditionFailed                    ErrorCode = "PreconditionFailed"
	OperationNotAllowed                   ErrorCode = "OperationNotAllowed"
	OperationPreempted                    ErrorCode = "OperationPreempted"
	PropertyChangeNotAllowed              ErrorCode = "PropertyChangeNotAllowed"
	InternalOperationError                ErrorCode = "InternalOperationError"
	InvalidSubscriptionStateTransition    ErrorCode = "InvalidSubscriptionStateTransition"
	UnregisterWithResourcesNotAllowed     ErrorCode = "UnregisterWithResourcesNotAllowed"
	InvalidParameterConflictingProperties ErrorCode = "InvalidParameterConflictingProperties"
	SubscriptionNotRegistered             ErrorCode = "SubscriptionNotRegistered"
	ConflictingUserInput                  ErrorCode = "ConflictingUserInput"
	ProvisioningInternalError             ErrorCode = "ProvisioningInternalError"
	ProvisioningFailed                    ErrorCode = "ProvisioningFailed"
	NetworkingInternalOperationError      ErrorCode = "NetworkingInternalOperationError"
	QuotaExceeded                         ErrorCode = "QuotaExceeded"
	Unauthorized                          ErrorCode = "Unauthorized"
	ResourcesOverConstrained              ErrorCode = "ResourcesOverConstrained"
	ControlPlaneProvisioningInternalError ErrorCode = "ControlPlaneProvisioningInternalError"
	ControlPlaneProvisioningSyncError     ErrorCode = "ControlPlaneProvisioningSyncError"
	ControlPlaneInternalError             ErrorCode = "ControlPlaneInternalError"
	ControlPlaneCloudProviderNotSet       ErrorCode = "ControlPlaneCloudProviderNotSet"

	// From Microsoft.WindowsAzure.ContainerService.API.AcsErrorCode
	ScaleDownInternalError ErrorCode = "ScaleDownInternalError"

	// New
	PreconditionCheckTimeOut     ErrorCode = "PreconditionCheckTimeOut"
	UpgradeFailed                ErrorCode = "UpgradeFailed"
	ScaleError                   ErrorCode = "ScaleError"
	CreateRoleAssignmentError    ErrorCode = "CreateRoleAssignmentError"
	ServicePrincipalNotFound     ErrorCode = "ServicePrincipalNotFound"
	ClusterResourceGroupNotFound ErrorCode = "ClusterResourceGroupNotFound"

	// Error codes returned by HCP
	UnderlayNotFound         ErrorCode = "UnderlayNotFound"
	UnderlaysOverConstrained ErrorCode = "UnderlaysOverConstrained"
	UnexpectedUnderlayCount  ErrorCode = "UnexpectedUnderlayCount"
)
