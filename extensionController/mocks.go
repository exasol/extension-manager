package extensionController

import (
	"context"
	"database/sql"

	"github.com/exasol/extension-manager/extensionAPI"
	"github.com/stretchr/testify/mock"
)

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
	} else {
		return args.Get(0).([]string), args.Error(1)
	}
}

func (mock *bucketFsMock) ListFiles(ctx context.Context, db *sql.DB, bucket string) ([]BfsFile, error) {
	args := mock.Called(ctx, db, bucket)
	if buckets, ok := args.Get(0).([]BfsFile); ok {
		return buckets, args.Error(1)
	} else {
		return args.Get(0).([]BfsFile), args.Error(1)
	}
}

type exaMetaDataReaderMock struct {
	mock.Mock
	extensionSchema string
}

func createExaMetaDataReaderMock(extensionSchema string) exaMetaDataReaderMock {
	var _ extensionAPI.ExaMetadataReader = &exaMetaDataReaderMock{}
	return exaMetaDataReaderMock{extensionSchema: extensionSchema}
}

func (m *exaMetaDataReaderMock) simulateExaAllScripts(scripts []extensionAPI.ExaAllScriptRow) {
	m.simulateExaMetaData(extensionAPI.ExaMetadata{AllScripts: extensionAPI.ExaAllScriptTable{Rows: scripts}})
}
func (m *exaMetaDataReaderMock) simulateExaMetaData(metaData extensionAPI.ExaMetadata) {
	m.On("ReadMetadataTables", mock.Anything, m.extensionSchema).Return(&metaData, nil)
}

func (mock *exaMetaDataReaderMock) ReadMetadataTables(tx *sql.Tx, schemaName string) (*extensionAPI.ExaMetadata, error) {
	args := mock.Called(tx, schemaName)
	if buckets, ok := args.Get(0).(*extensionAPI.ExaMetadata); ok {
		return buckets, args.Error(1)
	} else {
		return args.Get(0).(*extensionAPI.ExaMetadata), args.Error(1)
	}
}
