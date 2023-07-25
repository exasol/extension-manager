package restAPI

import (
	"context"
	"database/sql"

	"github.com/exasol/extension-manager/pkg/extensionAPI"
	"github.com/exasol/extension-manager/pkg/extensionController"
	"github.com/exasol/extension-manager/pkg/parameterValidator"
	"github.com/stretchr/testify/mock"
)

type mockExtensionController struct {
	mock.Mock
}

func (m *mockExtensionController) InstallExtension(ctx context.Context, db *sql.DB, extensionId string, extensionVersion string) error {
	args := m.Called(ctx, db, extensionId, extensionVersion)
	return args.Error(0)
}

func (m *mockExtensionController) UninstallExtension(ctx context.Context, db *sql.DB, extensionId string, extensionVersion string) error {
	args := m.Called(ctx, db, extensionId, extensionVersion)
	return args.Error(0)
}

func (m *mockExtensionController) UpgradeExtension(ctx context.Context, db *sql.DB, extensionId string) (*extensionAPI.JsUpgradeResult, error) {
	args := m.Called(ctx, db, extensionId)
	if result, ok := args.Get(0).(*extensionAPI.JsUpgradeResult); ok {
		return result, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockExtensionController) GetInstalledExtensions(ctx context.Context, db *sql.DB) ([]*extensionAPI.JsExtInstallation, error) {
	args := m.Called(ctx, db)
	if installations, ok := args.Get(0).([]*extensionAPI.JsExtInstallation); ok {
		return installations, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockExtensionController) GetParameterDefinitions(ctx context.Context, db *sql.DB, extensionId string, extensionVersion string) ([]parameterValidator.ParameterDefinition, error) {
	args := m.Called(ctx, db, extensionId, extensionVersion)
	if paramDefinitions, ok := args.Get(0).([]parameterValidator.ParameterDefinition); ok {
		return paramDefinitions, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockExtensionController) GetAllExtensions(ctx context.Context, db *sql.DB) ([]*extensionController.Extension, error) {
	args := m.Called(ctx, db)
	if extensions, ok := args.Get(0).([]*extensionController.Extension); ok {
		return extensions, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockExtensionController) CreateInstance(ctx context.Context, db *sql.DB, extensionId string, extensionVersion string, parameterValues []extensionController.ParameterValue) (*extensionAPI.JsExtInstance, error) {
	args := m.Called(ctx, db, extensionId, extensionVersion, parameterValues)
	if instance, ok := args.Get(0).(*extensionAPI.JsExtInstance); ok {
		return instance, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockExtensionController) FindInstances(ctx context.Context, db *sql.DB, extensionId string, extensionVersion string) ([]*extensionAPI.JsExtInstance, error) {
	args := m.Called(ctx, db, extensionId, extensionVersion)
	if instances, ok := args.Get(0).([]*extensionAPI.JsExtInstance); ok {
		return instances, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockExtensionController) DeleteInstance(ctx context.Context, db *sql.DB, extensionId, extensionVersion, instanceId string) error {
	args := m.Called(ctx, db, extensionId, extensionVersion, instanceId)
	return args.Error(0)
}
