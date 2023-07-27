package backend

import (
	"database/sql"

	"github.com/stretchr/testify/mock"
)

type SimpleSqlClientMock struct {
	mock.Mock
}

func (mock *SimpleSqlClientMock) SimulateExecuteSuccess(query string, args ...any) {
	var mockResult sql.Result = &mockSqlResult{}
	if len(args) > 0 {
		mock.On("Execute", query, args).Return(mockResult, nil)
	} else {
		mock.On("Execute", query, []interface{}{}).Return(mockResult, nil)
	}
}

func (mock *SimpleSqlClientMock) SimulateExecuteError(err error, query string, args ...any) {
	mock.On("Execute", query, args).Return(nil, err)
}

func (mock *SimpleSqlClientMock) Execute(query string, args ...any) (sql.Result, error) {
	mockArgs := mock.Called(query, args)
	return mockArgs.Get(0).(sql.Result), mockArgs.Error(1)
}

func (mock *SimpleSqlClientMock) SimulateQuerySuccess(result *QueryResult, query string, args ...any) {
	mock.On("Query", query, args).Return(result, nil)
}

func (mock *SimpleSqlClientMock) SimulateQueryError(err error, query string, args ...any) {
	mock.On("Query", query, args).Return(nil, err)
}

func (mock *SimpleSqlClientMock) Query(query string, args ...any) (*QueryResult, error) {
	mockArgs := mock.Called(query, args)
	return mockArgs.Get(0).(*QueryResult), mockArgs.Error(1)
}

type mockSqlResult struct{}

func (m *mockSqlResult) LastInsertId() (int64, error) {
	return 0, nil
}
func (m *mockSqlResult) RowsAffected() (int64, error) {
	return 0, nil
}
