package apiErrors

import (
	"fmt"
	"net/http"
)

func NewInternalServerError(originalError error) error {
	return &APIError{
		Status:        http.StatusInternalServerError,
		Message:       "Internal server error",
		OriginalError: originalError,
	}
}

func NewBadRequestErrorF(format string, a ...interface{}) error {
	return NewAPIErrorF(http.StatusBadRequest, format, a...)
}

func NewUnauthorizedErrorF(format string, a ...interface{}) error {
	return NewAPIErrorF(http.StatusUnauthorized, format, a...)
}

func NewAPIErrorF(status int, format string, a ...interface{}) error {
	return &APIError{
		Status:  status,
		Message: fmt.Sprintf(format, a...),
	}
}

func NewAPIError(status int, message string) error {
	return &APIError{
		Status:  status,
		Message: message,
	}
}

func NewAPIErrorWithCause(message string, cause error) error {
	if apiErr, ok := cause.(*APIError); ok {
		return &APIError{
			Status:        apiErr.Status,
			Message:       fmt.Sprintf("%s: %s", message, apiErr.Message),
			OriginalError: cause,
		}
	}
	return fmt.Errorf("%s: %w", message, cause)
}

type APIError struct {
	Status        int    `json:"code"`                // HTTP status code
	Message       string `json:"message"`             // human-readable message
	RequestID     string `json:"requestID,omitempty"` // ID to identify the request that caused this error
	Timestamp     string `json:"timestamp,omitempty" jsonschema:"format=date-time" `
	APIID         string `json:"apiID,omitempty"` // Corresponding API action
	OriginalError error  `json:"-"`
}

func (a *APIError) Error() string {
	return a.Message
}
