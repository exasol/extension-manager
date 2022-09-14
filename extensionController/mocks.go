package extensionController

import (
	"context"
	"database/sql"

	"github.com/exasol/extension-manager/extensionAPI"
	"github.com/stretchr/testify/mock"
)

// BucketFs

type bucketFsMock struct {
	mock.Mock
}

func createBucketFsMock() bucketFsMock {
	var _ BucketFsAPI = &bucketFsMock{}
	return bucketFsMock{}
}

func (m *bucketFsMock) simulateFiles(files []BfsFile) {
	m.On("ListFiles", mock.Anything, mock.Anything, "default").Return(files, nil)
}

func (m *bucketFsMock) ListBuckets(ctx context.Context, db *sql.DB) ([]string, error) {
	args := m.Called(ctx, db)
	if buckets, ok := args.Get(0).([]string); ok {
		return buckets, args.Error(1)
	}
	return nil, args.Error(1)
}

func (mock *bucketFsMock) ListFiles(ctx context.Context, db *sql.DB, bucket string) ([]BfsFile, error) {
	args := mock.Called(ctx, db, bucket)
	if buckets, ok := args.Get(0).([]BfsFile); ok {
		return buckets, args.Error(1)
	}
	return nil, args.Error(1)
}

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

func (mock *mockControllerImpl) GetAllExtensions(bfsFiles []BfsFile) ([]*Extension, error) {
	args := mock.Called(bfsFiles)
	if ext, ok := args.Get(0).([]*Extension); ok {
		return ext, args.Error(1)
	}
	return nil, args.Error(1)
}

func (mock *mockControllerImpl) GetAllInstallations(tx *sql.Tx) ([]*extensionAPI.JsExtInstallation, error) {
	args := mock.Called(tx)
	if result, ok := args.Get(0).([]*extensionAPI.JsExtInstallation); ok {
		return result, args.Error(1)
	}
	return nil, args.Error(1)
}

func (mock *mockControllerImpl) InstallExtension(tx *sql.Tx, extensionId string, extensionVersion string) error {
	args := mock.Called(tx, extensionId, extensionVersion)
	return args.Error(0)
}

func (mock *mockControllerImpl) UninstallExtension(tx *sql.Tx, extensionId string, extensionVersion string) error {
	args := mock.Called(tx, extensionId, extensionVersion)
	return args.Error(0)
}

func (mock *mockControllerImpl) CreateInstance(tx *sql.Tx, extensionId string, extensionVersion string, parameterValues []ParameterValue) (*extensionAPI.JsExtInstance, error) {
	args := mock.Called(tx, extensionId, extensionVersion, parameterValues)
	if result, ok := args.Get(0).(*extensionAPI.JsExtInstance); ok {
		return result, args.Error(1)
	}
	return nil, args.Error(1)
}

func (mock *mockControllerImpl) FindInstances(tx *sql.Tx, extensionId string, extensionVersion string) ([]*extensionAPI.JsExtInstance, error) {
	args := mock.Called(tx, extensionId, extensionVersion)
	if result, ok := args.Get(0).([]*extensionAPI.JsExtInstance); ok {
		return result, args.Error(1)
	}
	return nil, args.Error(1)
}

func (mock *mockControllerImpl) DeleteInstance(tx *sql.Tx, extensionId string, instanceId string) error {
	args := mock.Called(tx, extensionId, instanceId)
	return args.Error(0)
}
