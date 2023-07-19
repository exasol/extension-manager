package bfs

import (
	"context"
	"database/sql"

	"github.com/stretchr/testify/mock"
)

type BucketFsMock struct {
	mock.Mock
}

func (m *BucketFsMock) SimulateFiles(files []BfsFile) {
	m.On("ListFiles", mock.Anything, mock.Anything).Return(files, nil)
}

func (mock *BucketFsMock) ListFiles(ctx context.Context, db *sql.DB) ([]BfsFile, error) {
	args := mock.Called(ctx, db)
	if buckets, ok := args.Get(0).([]BfsFile); ok {
		return buckets, args.Error(1)
	}
	return nil, args.Error(1)
}
