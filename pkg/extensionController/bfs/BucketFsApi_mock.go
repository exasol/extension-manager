package bfs

import (
	"context"
	"database/sql"

	"github.com/stretchr/testify/mock"
)

type BucketFsMock struct {
	mock.Mock
}

func CreateBucketFsMock() *BucketFsMock {
	//nolint:exhaustruct // Empty struct is OK for Mock
	return &BucketFsMock{}
}

func (m *BucketFsMock) SimulateFiles(files []BfsFile) {
	m.On("ListFiles", mock.Anything, mock.Anything).Return(files, nil)
}

func (m *BucketFsMock) SimulateFilesError(err error) {
	m.On("ListFiles", mock.Anything, mock.Anything).Return(nil, err)
}

func (m *BucketFsMock) SimulateAbsolutePath(fileName, absolutePath string) {
	m.On("FindAbsolutePath", mock.Anything, mock.Anything, fileName).Return(absolutePath, nil)
}

func (m *BucketFsMock) SimulateAbsolutePathError(fileName string, err error) {
	m.On("FindAbsolutePath", mock.Anything, mock.Anything, fileName).Return("", err)
}

func (mock *BucketFsMock) ListFiles(ctx context.Context, db *sql.DB) ([]BfsFile, error) {
	args := mock.Called(ctx, db)
	if buckets, ok := args.Get(0).([]BfsFile); ok {
		return buckets, args.Error(1)
	}
	return nil, args.Error(1)
}

func (mock *BucketFsMock) FindAbsolutePath(ctx context.Context, db *sql.DB, fileName string) (absolutePath string, retErr error) {
	args := mock.Called(ctx, db, fileName)
	return args.String(0), args.Error(1)
}
