package apiErrors_test

import (
	"fmt"
	"testing"

	"github.com/exasol/extension-manager/apiErrors"
	"github.com/stretchr/testify/assert"
)

func TestNewInternalServerError(t *testing.T) {
	orgErr := fmt.Errorf("mock")
	err := apiErrors.NewInternalServerError(orgErr)
	assertApiError(t, err, "Internal server error", 500, orgErr)
}

func TestNewBadRequestErrorF(t *testing.T) {
	err := apiErrors.NewBadRequestErrorF("err %d", 42)
	assertApiError(t, err, "err 42", 400, nil)
}

func TestNewUnauthorizedErrorF(t *testing.T) {
	err := apiErrors.NewUnauthorizedErrorF("err %d", 42)
	assertApiError(t, err, "err 42", 401, nil)
}

func TestNewAPIErrorF(t *testing.T) {
	err := apiErrors.NewAPIErrorF(123, "err %d", 42)
	assertApiError(t, err, "err 42", 123, nil)
}

func TestNewAPIError(t *testing.T) {
	err := apiErrors.NewAPIError(123, "err")
	assertApiError(t, err, "err", 123, nil)
}

func TestNewAPIErrorWithCause_ApiErrorCause(t *testing.T) {
	cause := apiErrors.NewAPIError(123, "cause")
	err := apiErrors.NewAPIErrorWithCause("msg", cause)
	assertApiError(t, err, "msg: cause", 123, cause)
}

func TestNewAPIErrorWithCause_nonApiErrorCause(t *testing.T) {
	cause := fmt.Errorf("cause")
	err := apiErrors.NewAPIErrorWithCause("msg", cause)
	assert.EqualError(t, err, "msg: cause")
}

func assertApiError(t *testing.T, err error, expectedMsg string, expectedStatus int, expectedOrgError error) {
	t.Helper()
	if apiErr, ok := err.(*apiErrors.APIError); ok {
		assert.Equal(t, expectedMsg, apiErr.Message)
		assert.Equal(t, expectedStatus, apiErr.Status)
		if expectedOrgError == nil {
			assert.Nil(t, apiErr.OriginalError)
		} else {
			assert.Same(t, expectedOrgError, apiErr.OriginalError)
		}
	} else {
		t.Errorf("Expected an ApiError but got %T", err)
	}
}
