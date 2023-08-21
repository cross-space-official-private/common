package failure

import (
	"fmt"

	"github.com/cross-space-official-private/common/businesserror"
	"github.com/joomcode/errorx"
)

const (
	ResourceNotFoundError = "resource_not_found"
	InvalidArgumentError  = "invalid_arguments"
	NotAuthenticatedError = "authentication_failed"
	NoAuthorizationError  = "authorization_failed"
	InternalServiceError  = "internal_service_error"
	TooManyRequestsError  = "too_many_requests"
)

type Failure struct {
	errorCode  string
	message    string
	stacktrace string
}

func GeneratePlainFailure(errorCode string, message string, stacktrace string) *Failure {
	return &Failure{
		errorCode:  errorCode,
		message:    message,
		stacktrace: stacktrace,
	}
}

func GenerateFailureFromError(errorCode string, err businesserror.XSpaceBusinessError) *Failure {
	return &Failure{
		errorCode:  errorCode,
		message:    err.Error(),
		stacktrace: err.Stacktrace(),
	}
}

func GenerateFailureFromErrorx(errorCode string, err *errorx.Error) *Failure {
	return &Failure{
		errorCode:  errorCode,
		message:    err.Message(),
		stacktrace: fmt.Sprintf("%+v", err),
	}
}

func (failure *Failure) Error() string {
	return failure.errorCode
}

func (failure *Failure) Stacktrace() string {
	return failure.stacktrace
}

func (failure *Failure) Message() string {
	return failure.message
}
