//nolint:unused // Mock functions are actually used in tests
package extensionController

import (
	"github.com/exasol/extension-manager/pkg/extensionAPI"
	"github.com/exasol/extension-manager/pkg/extensionController/bfs"
	"github.com/exasol/extension-manager/pkg/extensionController/transaction"
	"github.com/exasol/extension-manager/pkg/parameterValidator"
	"github.com/stretchr/testify/mock"
)

type mockControllerImpl struct {
	mock.Mock
}

func (mock *mockControllerImpl) GetAllExtensions(bfsFiles []bfs.BfsFile) ([]*Extension, error) {
	args := mock.Called(bfsFiles)
	if ext, ok := args.Get(0).([]*Extension); ok {
		return ext, args.Error(1)
	}
	return nil, args.Error(1)
}

func (mock *mockControllerImpl) GetParameterDefinitions(txCtx *transaction.TransactionContext, extensionId string, extensionVersion string) ([]parameterValidator.ParameterDefinition, error) {
	args := mock.Called(extensionId, extensionVersion)
	if result, ok := args.Get(0).([]parameterValidator.ParameterDefinition); ok {
		return result, args.Error(1)
	}
	return nil, args.Error(1)
}

func (mock *mockControllerImpl) GetAllInstallations(txCtx *transaction.TransactionContext) ([]*extensionAPI.JsExtInstallation, error) {
	args := mock.Called(txCtx)
	if result, ok := args.Get(0).([]*extensionAPI.JsExtInstallation); ok {
		return result, args.Error(1)
	}
	return nil, args.Error(1)
}

func (mock *mockControllerImpl) InstallExtension(txCtx *transaction.TransactionContext, extensionId string, extensionVersion string) error {
	args := mock.Called(txCtx, extensionId, extensionVersion)
	return args.Error(0)
}

func (mock *mockControllerImpl) UninstallExtension(txCtx *transaction.TransactionContext, extensionId string, extensionVersion string) error {
	args := mock.Called(txCtx, extensionId, extensionVersion)
	return args.Error(0)
}

func (mock *mockControllerImpl) UpgradeExtension(txCtx *transaction.TransactionContext, extensionId string) (*extensionAPI.JsUpgradeResult, error) {
	args := mock.Called(txCtx, extensionId)
	if result, ok := args.Get(0).(*extensionAPI.JsUpgradeResult); ok {
		return result, args.Error(1)
	}
	return nil, args.Error(1)
}

func (mock *mockControllerImpl) CreateInstance(txCtx *transaction.TransactionContext, extensionId string, extensionVersion string, parameterValues []ParameterValue) (*extensionAPI.JsExtInstance, error) {
	args := mock.Called(txCtx, extensionId, extensionVersion, parameterValues)
	if result, ok := args.Get(0).(*extensionAPI.JsExtInstance); ok {
		return result, args.Error(1)
	}
	return nil, args.Error(1)
}

func (mock *mockControllerImpl) FindInstances(txCtx *transaction.TransactionContext, extensionId string, extensionVersion string) ([]*extensionAPI.JsExtInstance, error) {
	args := mock.Called(txCtx, extensionId, extensionVersion)
	if result, ok := args.Get(0).([]*extensionAPI.JsExtInstance); ok {
		return result, args.Error(1)
	}
	return nil, args.Error(1)
}

func (mock *mockControllerImpl) DeleteInstance(txCtx *transaction.TransactionContext, extensionId, extensionVersion, instanceId string) error {
	args := mock.Called(txCtx, extensionId, extensionVersion, instanceId)
	return args.Error(0)
}
