package bfs

import (
	"fmt"

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
	m.On("ListFiles").Return(files, nil)
}

func (m *BucketFsMock) SimulateFilesError(err error) {
	m.On("ListFiles").Return(nil, err)
}

func (m *BucketFsMock) SimulateAbsolutePath(fileName, absolutePath string) {
	m.On("FindAbsolutePath", fileName).Return(absolutePath, nil)
}

func (m *BucketFsMock) SimulateAbsolutePathError(fileName string, err error) {
	m.On("FindAbsolutePath", fileName).Return("", err)
}

func (m *BucketFsMock) SimulateCloseSuccess() {
	fmt.Println("### Expecting close success")
	m.On("Close").Return(nil)
}

func (m *BucketFsMock) SimulateCloseFails(err error) {
	fmt.Println("### Expecting close failure")
	m.On("Close").Return(err)
}

func (mock *BucketFsMock) ListFiles() ([]BfsFile, error) {
	args := mock.Called()
	if buckets, ok := args.Get(0).([]BfsFile); ok {
		return buckets, args.Error(1)
	}
	return nil, args.Error(1)
}

func (mock *BucketFsMock) FindAbsolutePath(fileName string) (absolutePath string, retErr error) {
	args := mock.Called(fileName)
	return args.String(0), args.Error(1)
}

func (mock *BucketFsMock) Close() error {
	fmt.Println("BFS Mock is closed now!")
	args := mock.Called()
	return args.Error(0)
}
