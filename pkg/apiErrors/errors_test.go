package apiErrors_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/exasol/extension-manager/pkg/apiErrors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewInternalServerError(t *testing.T) {
	orgErr := errors.New("mock")
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

func TestNewAPIErrorWithCauseApiErrorCause(t *testing.T) {
	cause := apiErrors.NewAPIError(123, "cause")
	err := apiErrors.NewAPIErrorWithCause("msg", cause)
	assertApiError(t, err, "msg: cause", 123, cause)
}

func TestNewAPIErrorWithCauseNonApiErrorCause(t *testing.T) {
	cause := errors.New("cause")
	err := apiErrors.NewAPIErrorWithCause("msg", cause)
	assert.EqualError(t, err, "msg: cause")
}

func TestUnwrapAPIErrorNormalError(t *testing.T) {
	orgErr := errors.New("mock")
	err := apiErrors.UnwrapAPIError(orgErr)
	assertApiError(t, err, "Internal server error", 500, orgErr)
}

func TestUnwrapAPIErrorAPIError(t *testing.T) {
	orgErr := apiErrors.NewNotFoundErrorF("extension not found")
	err := apiErrors.UnwrapAPIError(orgErr)
	assertApiError(t, err, "extension not found", 404, nil)
}

func TestUnwrapAPIErrorWrappedAPIError(t *testing.T) {
	orgErr := apiErrors.NewNotFoundErrorF("extension not found")
	orgErr = fmt.Errorf("wrapper %w", orgErr)
	err := apiErrors.UnwrapAPIError(orgErr)
	assertApiError(t, err, "extension not found", 404, nil)
}

func TestUnwrapAPIErrorMultipleWrapperAPIError(t *testing.T) {
	orgErr := apiErrors.NewNotFoundErrorF("extension not found")
	orgErr = fmt.Errorf("wrapper1 %w", orgErr)
	orgErr = fmt.Errorf("wrapper2 %w", orgErr)
	orgErr = fmt.Errorf("wrapper3 %w", orgErr)
	err := apiErrors.UnwrapAPIError(orgErr)
	assertApiError(t, err, "extension not found", 404, nil)
}

func assertApiError(t *testing.T, err error, expectedMsg string, expectedStatus int, expectedOrgError error) {
	t.Helper()
	if apiErr, ok := apiErrors.AsAPIError(err); ok {
		assert.Equal(t, expectedMsg, apiErr.Message)
		assert.Equal(t, expectedStatus, apiErr.Status)
		if expectedOrgError == nil {
			require.NoError(t, apiErr.OriginalError, "original error")
		} else {
			assert.Same(t, expectedOrgError, apiErr.OriginalError, "original error")
		}
	} else {
		t.Errorf("Expected an ApiError but got %T", err)
	}
}
