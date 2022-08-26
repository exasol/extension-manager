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

func (mock *bucketFsMock) ListBuckets(ctx context.Context, db *sql.DB) ([]string, error) {
	args := mock.Called(ctx, db)
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
}

func createExaMetaDataReaderMock() exaMetaDataReaderMock {
	var _ extensionAPI.ExaMetadataReader = &exaMetaDataReaderMock{}
	return exaMetaDataReaderMock{}
}

func (mock *exaMetaDataReaderMock) ReadMetadataTables(tx *sql.Tx, schemaName string) (*extensionAPI.ExaMetadata, error) {
	args := mock.Called(tx, schemaName)
	if buckets, ok := args.Get(0).(*extensionAPI.ExaMetadata); ok {
		return buckets, args.Error(1)
	} else {
		return args.Get(0).(*extensionAPI.ExaMetadata), args.Error(1)
	}
}
