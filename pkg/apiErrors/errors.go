package apiErrors

import (
	"errors"
	"fmt"
	"net/http"
)

func NewInternalServerError(originalError error) *APIError {
	return &APIError{
		Status:        http.StatusInternalServerError,
		Message:       "Internal server error",
		RequestID:     "",
		OriginalError: originalError,
	}
}

func NewNotFoundErrorF(format string, a ...interface{}) error {
	return NewAPIErrorF(http.StatusNotFound, format, a...)
}

func NewBadRequestErrorF(format string, a ...interface{}) error {
	return NewAPIErrorF(http.StatusBadRequest, format, a...)
}

func NewUnauthorizedErrorF(format string, a ...interface{}) error {
	return NewAPIErrorF(http.StatusUnauthorized, format, a...)
}

func NewAPIErrorF(status int, format string, a ...interface{}) error {
	return &APIError{
		Status:        status,
		Message:       fmt.Sprintf(format, a...),
		RequestID:     "",
		OriginalError: nil,
	}
}

func NewAPIError(status int, message string) error {
	return &APIError{
		Status:        status,
		Message:       message,
		RequestID:     "",
		OriginalError: nil,
	}
}

func NewAPIErrorWithCause(message string, cause error) error {
	var apiErr *APIError
	if errors.As(cause, &apiErr) {
		return &APIError{
			Status:        apiErr.Status,
			Message:       fmt.Sprintf("%s: %s", message, apiErr.Message),
			RequestID:     "",
			OriginalError: cause,
		}
	}
	return fmt.Errorf("%s: %w", message, cause)
}

// UnwrapAPIError returns an API Error if one exists in the given error's chain or nil if no API Error exists in the chain.
func UnwrapAPIError(err error) *APIError {
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr
	}
	return NewInternalServerError(err)
}

type APIError struct {
	Status        int    `json:"code"`                // HTTP status code
	Message       string `json:"message"`             // human-readable message
	RequestID     string `json:"requestID,omitempty"` // ID to identify the request that caused this error
	OriginalError error  `json:"-"`
}

func (a *APIError) Error() string {
	return a.Message
}

func AsAPIError(err error) (*APIError, bool) {
	//nolint:errorlint // Type assertion via cast is OK here as the error is not wrapped
	if apiError, ok := err.(*APIError); ok {
		return apiError, true
	} else {
		return nil, false
	}
}
