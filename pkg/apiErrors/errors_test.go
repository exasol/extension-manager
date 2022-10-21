package apiErrors_test

import (
	"fmt"
	"testing"

	"github.com/exasol/extension-manager/pkg/apiErrors"
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

func TestUnwrapAPIError_normalError(t *testing.T) {
	orgErr := fmt.Errorf("mock")
	err := apiErrors.UnwrapAPIError(orgErr)
	assertApiError(t, err, "Internal server error", 500, orgErr)
}

func TestUnwrapAPIError_APIError(t *testing.T) {
	orgErr := apiErrors.NewNotFoundErrorF("extension not found")
	err := apiErrors.UnwrapAPIError(orgErr)
	assertApiError(t, err, "extension not found", 404, nil)
}

func TestUnwrapAPIError_wrappedAPIError(t *testing.T) {
	orgErr := apiErrors.NewNotFoundErrorF("extension not found")
	orgErr = fmt.Errorf("wrapper %w", orgErr)
	err := apiErrors.UnwrapAPIError(orgErr)
	assertApiError(t, err, "extension not found", 404, nil)
}

func TestUnwrapAPIError_multipleWrapperAPIError(t *testing.T) {
	orgErr := apiErrors.NewNotFoundErrorF("extension not found")
	orgErr = fmt.Errorf("wrapper1 %w", orgErr)
	orgErr = fmt.Errorf("wrapper2 %w", orgErr)
	orgErr = fmt.Errorf("wrapper3 %w", orgErr)
	err := apiErrors.UnwrapAPIError(orgErr)
	assertApiError(t, err, "extension not found", 404, nil)
}

func assertApiError(t *testing.T, err error, expectedMsg string, expectedStatus int, expectedOrgError error) {
	t.Helper()
	if apiErr, ok := err.(*apiErrors.APIError); ok {
		assert.Equal(t, expectedMsg, apiErr.Message)
		assert.Equal(t, expectedStatus, apiErr.Status)
		if expectedOrgError == nil {
			assert.Nil(t, apiErr.OriginalError, "original error")
		} else {
			assert.Same(t, expectedOrgError, apiErr.OriginalError, "original error")
		}
	} else {
		t.Errorf("Expected an ApiError but got %T", err)
	}
}
