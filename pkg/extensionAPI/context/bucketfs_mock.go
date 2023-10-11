package context

import (
	"github.com/stretchr/testify/mock"
)

type BucketFsContextMock struct {
	mock.Mock
}

func CreateBucketFsContextMock() *BucketFsContextMock {
	//nolint:exhaustruct // Empty struct is OK for Mock
	return &BucketFsContextMock{}
}

func (mock *BucketFsContextMock) SimulateResolvePath(fileName string, result string) {
	mock.On("ResolvePath", fileName).Return(result)
}

func (mock *BucketFsContextMock) SimulateResolvePathPanics(fileName string, panicMessage string) {
	mock.On("ResolvePath", fileName).Panic(panicMessage)
}

func (mock *BucketFsContextMock) ResolvePath(fileName string) string {
	mockArgs := mock.Called(fileName)
	return mockArgs.String(0)
}
