package extensionController

import (
	"database/sql"

	"github.com/exasol/extension-manager/pkg/extensionAPI"
	"github.com/exasol/extension-manager/pkg/extensionController/bfs"
	"github.com/exasol/extension-manager/pkg/extensionController/transactionContext"
	"github.com/exasol/extension-manager/pkg/parameterValidator"
	"github.com/stretchr/testify/mock"
)

// Exa metadata reader

type exaMetaDataReaderMock struct {
	mock.Mock
	extensionSchema string
}

func createExaMetaDataReaderMock(extensionSchema string) exaMetaDataReaderMock {
	var _ extensionAPI.ExaMetadataReader = &exaMetaDataReaderMock{}
	return exaMetaDataReaderMock{extensionSchema: extensionSchema}
}

func (m *exaMetaDataReaderMock) simulateExaAllScripts(scripts []extensionAPI.ExaScriptRow) {
	m.simulateExaMetaData(extensionAPI.ExaMetadata{AllScripts: extensionAPI.ExaScriptTable{Rows: scripts}})
}

func (m *exaMetaDataReaderMock) simulateExaMetaData(metaData extensionAPI.ExaMetadata) {
	m.On("ReadMetadataTables", mock.Anything, m.extensionSchema).Return(&metaData, nil)
}

func (mock *exaMetaDataReaderMock) ReadMetadataTables(tx *sql.Tx, schemaName string) (*extensionAPI.ExaMetadata, error) {
	args := mock.Called(tx, schemaName)
	if buckets, ok := args.Get(0).(*extensionAPI.ExaMetadata); ok {
		return buckets, args.Error(1)
	}
	return nil, args.Error(1)
}

// controller

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

func (mock *mockControllerImpl) GetParameterDefinitions(txCtx *transactionContext.TransactionContext, extensionId string, extensionVersion string) ([]parameterValidator.ParameterDefinition, error) {
	args := mock.Called(extensionId, extensionVersion)
	if result, ok := args.Get(0).([]parameterValidator.ParameterDefinition); ok {
		return result, args.Error(1)
	}
	return nil, args.Error(1)
}

func (mock *mockControllerImpl) GetAllInstallations(txCtx *transactionContext.TransactionContext) ([]*extensionAPI.JsExtInstallation, error) {
	args := mock.Called(txCtx)
	if result, ok := args.Get(0).([]*extensionAPI.JsExtInstallation); ok {
		return result, args.Error(1)
	}
	return nil, args.Error(1)
}

func (mock *mockControllerImpl) InstallExtension(txCtx *transactionContext.TransactionContext, extensionId string, extensionVersion string) error {
	args := mock.Called(txCtx, extensionId, extensionVersion)
	return args.Error(0)
}

func (mock *mockControllerImpl) UninstallExtension(txCtx *transactionContext.TransactionContext, extensionId string, extensionVersion string) error {
	args := mock.Called(txCtx, extensionId, extensionVersion)
	return args.Error(0)
}

func (mock *mockControllerImpl) CreateInstance(txCtx *transactionContext.TransactionContext, extensionId string, extensionVersion string, parameterValues []ParameterValue) (*extensionAPI.JsExtInstance, error) {
	args := mock.Called(txCtx, extensionId, extensionVersion, parameterValues)
	if result, ok := args.Get(0).(*extensionAPI.JsExtInstance); ok {
		return result, args.Error(1)
	}
	return nil, args.Error(1)
}

func (mock *mockControllerImpl) FindInstances(txCtx *transactionContext.TransactionContext, extensionId string, extensionVersion string) ([]*extensionAPI.JsExtInstance, error) {
	args := mock.Called(txCtx, extensionId, extensionVersion)
	if result, ok := args.Get(0).([]*extensionAPI.JsExtInstance); ok {
		return result, args.Error(1)
	}
	return nil, args.Error(1)
}

func (mock *mockControllerImpl) DeleteInstance(txCtx *transactionContext.TransactionContext, extensionId, extensionVersion, instanceId string) error {
	args := mock.Called(txCtx, extensionId, extensionVersion, instanceId)
	return args.Error(0)
}
