package context

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/exasol/extension-manager/pkg/backend"
	"github.com/exasol/extension-manager/pkg/extensionAPI/exaMetadata"
	"github.com/exasol/extension-manager/pkg/extensionController/transaction"
	"github.com/stretchr/testify/suite"
)

type ContextSuite struct {
	suite.Suite
	db                 *sql.DB
	dbMock             sqlmock.Sqlmock
	bucketFSMock       *BucketFsContextMock
	metadataReaderMock *exaMetadata.ExaMetaDataReaderMock
}

func TestContextSuite(t *testing.T) {
	suite.Run(t, new(ContextSuite))
}

const EXTENSION_SCHEMA = "EXT_SCHEMA"
const BUCKETFS_BASE_PATH = "bucketfs-base-path"

func (suite *ContextSuite) SetupTest() {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
	suite.Require().NoError(err)
	suite.db = db
	suite.dbMock = mock
	suite.dbMock.MatchExpectationsInOrder(true)
	suite.bucketFSMock = CreateBucketFsContextMock()
	suite.metadataReaderMock = exaMetadata.CreateExaMetaDataReaderMock(EXTENSION_SCHEMA)
}

func (suite *ContextSuite) AfterTest(suiteName, testName string) {
	suite.NoError(suite.dbMock.ExpectationsWereMet())
	suite.bucketFSMock.AssertExpectations(suite.T())
	suite.metadataReaderMock.AssertExpectations(suite.T())
}

func (suite *ContextSuite) TestCreate() {
	ctx := suite.createContext()
	suite.NotNil(ctx)
	suite.Equal("EXT_SCHEMA", ctx.ExtensionSchemaName)
	suite.NotNil(ctx.BucketFs)
	suite.NotNil(ctx.SqlClient)
}

/* [utest -> dsn~extension-context-sql-client~1]. */
func (suite *ContextSuite) TestSqlClientQuerySuccess() {
	ctx := suite.createContext()
	suite.dbMock.ExpectQuery("select 1").WillReturnRows(sqlmock.NewRowsWithColumnDefinition(
		sqlmock.NewColumn("col1").OfType("type1", "sample"),
		sqlmock.NewColumn("col2").OfType("type2", "sample"),
	).
		AddRow(1, "a").
		AddRow(2, "b")).
		RowsWillBeClosed()
	result := ctx.SqlClient.Query("select 1")
	suite.Equal(backend.QueryResult{
		Columns: []backend.Column{{Name: "col1", TypeName: "type1"}, {Name: "col2", TypeName: "type2"}},
		Rows:    []backend.Row{{int64(1), "a"}, {int64(2), "b"}}}, result)
}

func (suite *ContextSuite) TestSqlClientQueryFailure() {
	ctx := suite.createContext()
	suite.dbMock.ExpectQuery("invalid").WillReturnError(errors.New("mock error"))
	suite.PanicsWithError("error executing query 'invalid': mock error", func() {
		ctx.SqlClient.Query("invalid")
	})
}

/* [utest -> dsn~extension-context-sql-client~1]. */
func (suite *ContextSuite) TestSqlClientExecuteSuccess() {
	ctx := suite.createContext()
	suite.dbMock.ExpectExec("create script").WillReturnResult(sqlmock.NewResult(1, 1))
	suite.NotPanics(func() {
		ctx.SqlClient.Execute("create script")
	})
}

func (suite *ContextSuite) TestSqlClientExecuteFailure() {
	ctx := suite.createContext()
	suite.dbMock.ExpectExec("invalid").WillReturnError(errors.New("mock error"))
	suite.PanicsWithError("error executing statement 'invalid': mock error", func() {
		ctx.SqlClient.Execute("invalid")
	})
}

/* [utest -> dsn~extension-context-bucketfs~1]. */
func (suite *ContextSuite) TestBucketFsResolvePath() {
	ctx := suite.createContextWithClients()
	suite.bucketFSMock.SimulateResolvePath("file.txt", "/absolute/path/file.txt")
	suite.Equal("/absolute/path/file.txt", ctx.BucketFs.ResolvePath("file.txt"))
}

func (suite *ContextSuite) TestBucketFsResolvePathError() {
	ctx := suite.createContextWithClients()
	suite.bucketFSMock.SimulateResolvePathPanics("file.txt", "mock error")
	suite.PanicsWithValue("mock error", func() {
		ctx.BucketFs.ResolvePath("file.txt")
	})
}

/* [utest -> dsn~extension-context-metadata~1]. */
func (suite *ContextSuite) TestMetadataGetScriptByName() {
	ctx := suite.createContextWithClients()
	suite.metadataReaderMock.SimulateGetScriptByNameScriptText("script", "scriptText")
	suite.Equal(&exaMetadata.ExaScriptRow{Schema: "?", Name: "script", Type: "", InputType: "", ResultType: "", Text: "scriptText", Comment: ""}, ctx.Metadata.GetScriptByName("script"))
}

func (suite *ContextSuite) TestMetadataGetScriptByNameNoScriptFound() {
	ctx := suite.createContextWithClients()
	suite.metadataReaderMock.SimulateGetScriptByName("script", nil)
	suite.Nil(ctx.Metadata.GetScriptByName("script"))
}

func (suite *ContextSuite) TestMetadataGetScriptByNameFails() {
	ctx := suite.createContextWithClients()
	suite.metadataReaderMock.SimulateGetScriptByNameFails("script", errors.New("mock error"))
	suite.PanicsWithError(`failed to find script "EXT_SCHEMA"."script". Caused by: mock error`, func() {
		ctx.Metadata.GetScriptByName("script")
	})
}

func (suite *ContextSuite) createContext() *ExtensionContext {
	suite.dbMock.ExpectBegin()
	txCtx, err := transaction.BeginTransaction(context.Background(), suite.db, BUCKETFS_BASE_PATH)
	suite.Require().NoError(err)
	return CreateContext(txCtx, "EXT_SCHEMA")
}

func (suite *ContextSuite) createContextWithClients() *ExtensionContext {
	suite.dbMock.ExpectBegin()
	txCtx, err := transaction.BeginTransaction(context.Background(), suite.db, BUCKETFS_BASE_PATH)
	suite.Require().NoError(err)
	return CreateContextWithClient("EXT_SCHEMA", txCtx, nil, suite.bucketFSMock, suite.metadataReaderMock)
}
