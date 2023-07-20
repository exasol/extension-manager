package context

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/exasol/extension-manager/pkg/backend"
	"github.com/exasol/extension-manager/pkg/extensionController/bfs"
	"github.com/exasol/extension-manager/pkg/extensionController/transaction"
	"github.com/stretchr/testify/suite"
)

type ContextSuite struct {
	suite.Suite
	db           *sql.DB
	dbMock       sqlmock.Sqlmock
	bucketFSMock *bfs.BucketFsMock
}

func TestContextSuite(t *testing.T) {
	suite.Run(t, new(ContextSuite))
}

func (suite *ContextSuite) SetupTest() {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
	suite.NoError(err)
	suite.db = db
	suite.dbMock = mock
	suite.dbMock.MatchExpectationsInOrder(true)
	suite.bucketFSMock = &bfs.BucketFsMock{}
}

func (suite *ContextSuite) AfterTest(suiteName, testName string) {
	suite.NoError(suite.dbMock.ExpectationsWereMet())
}

func (suite *ContextSuite) TestCreate() {
	ctx := suite.createContext()
	suite.NotNil(ctx)
	suite.Equal("EXT_SCHEMA", ctx.ExtensionSchemaName)
	suite.NotNil(ctx.BucketFs)
	suite.NotNil(ctx.SqlClient)
}

func (suite *ContextSuite) TestSqlClientQuerySuccess() {
	ctx := suite.createContext()
	suite.dbMock.ExpectQuery("select 1").WillReturnRows(sqlmock.NewRows([]string{"col1", "col2"}).AddRow(1, "a").AddRow(2, "b")).RowsWillBeClosed()
	result := ctx.SqlClient.Query("select 1")
	suite.Equal(backend.QueryResult{
		Columns: []backend.Column{{Name: "col1"}, {Name: "col2"}},
		Rows:    []backend.Row{{int64(1), "a"}, {int64(2), "b"}}}, result)
}

func (suite *ContextSuite) TestSqlClientQueryFailure() {
	ctx := suite.createContext()
	suite.dbMock.ExpectQuery("invalid").WillReturnError(fmt.Errorf("mock error"))
	suite.PanicsWithError("error executing statement \"invalid\": mock error", func() {
		ctx.SqlClient.Query("invalid")
	})
}

func (suite *ContextSuite) TestSqlClientExecuteSuccess() {
	ctx := suite.createContext()
	suite.dbMock.ExpectExec("create script").WillReturnResult(sqlmock.NewResult(1, 1))
	suite.NotPanics(func() {
		ctx.SqlClient.Execute("create script")
	})
}

func (suite *ContextSuite) TestSqlClientExecuteFailure() {
	ctx := suite.createContext()
	suite.dbMock.ExpectExec("invalid").WillReturnError(fmt.Errorf("mock error"))
	suite.PanicsWithError("error executing statement \"invalid\": mock error", func() {
		ctx.SqlClient.Execute("invalid")
	})
}

func (suite *ContextSuite) TestBucketFsResolvePath() {
	ctx := suite.createContextWithClients()
	suite.bucketFSMock.SimulateAbsolutePath("file.txt", "/absolute/path/file.txt")
	suite.Equal("/absolute/path/file.txt", ctx.BucketFs.ResolvePath("file.txt"))
}

func (suite *ContextSuite) TestBucketFsResolvePathError() {
	ctx := suite.createContextWithClients()
	suite.bucketFSMock.SimulateAbsolutePathError("file.txt", fmt.Errorf("mock error"))
	suite.PanicsWithError("failed to find absolute path for file \"file.txt\": mock error", func() {
		ctx.BucketFs.ResolvePath("file.txt")
	})
}

func (suite *ContextSuite) createContext() *ExtensionContext {
	suite.dbMock.ExpectBegin()
	txCtx, err := transaction.BeginTransaction(context.Background(), suite.db)
	suite.NoError(err)
	return CreateContext(txCtx, "EXT_SCHEMA", "/bucketfs/base/path/")
}

func (suite *ContextSuite) createContextWithClients() *ExtensionContext {
	suite.dbMock.ExpectBegin()
	txCtx, err := transaction.BeginTransaction(context.Background(), suite.db)
	suite.NoError(err)
	return CreateContextWithClient("EXT_SCHEMA", txCtx, nil, suite.bucketFSMock)
}
