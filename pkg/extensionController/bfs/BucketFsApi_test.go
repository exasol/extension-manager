package bfs

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/suite"
)

const BUCKETFS_BASE_PATH = "/basePath/"
const FILE_NAME = "file.txt"

type BucketFsAPIUTestSuite struct {
	suite.Suite
	db     *sql.DB
	dbMock sqlmock.Sqlmock
}

func TestBucketFsApiUTestSuite(t *testing.T) {
	suite.Run(t, new(BucketFsAPIUTestSuite))
}

func (suite *BucketFsAPIUTestSuite) SetupTest() {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
	suite.NoError(err)
	suite.db = db
	suite.dbMock = mock
	suite.dbMock.MatchExpectationsInOrder(true)
}

func (suite *BucketFsAPIUTestSuite) AfterTest(suiteName, testName string) {
	suite.NoError(suite.dbMock.ExpectationsWereMet())
}

// CreateBucketFsAPI
func (suite *BucketFsAPIUTestSuite) TestCreateBucketFsAPI() {
	suite.dbMock.ExpectBegin()
	suite.dbMock.ExpectExec("CREATE SCHEMA INTERNAL_\\d+").WillReturnResult(sqlmock.NewResult(0, 1))
	suite.dbMock.ExpectExec("(?m)CREATE OR REPLACE PYTHON3 SCALAR SCRIPT.*").WillReturnResult(sqlmock.NewResult(0, 1))
	client, err := CreateBucketFsAPI(BUCKETFS_BASE_PATH, context.Background(), suite.db)
	suite.NoError(err)
	suite.NotNil(client)
}

func (suite *BucketFsAPIUTestSuite) TestCreateBucketFsAPIFailsCreatingTransaction() {
	suite.dbMock.ExpectBegin().WillReturnError(fmt.Errorf("mock error"))
	client, err := CreateBucketFsAPI(BUCKETFS_BASE_PATH, context.Background(), suite.db)
	suite.EqualError(err, "failed to create a transaction. Cause: mock error")
	suite.Nil(client)
}

func (suite *BucketFsAPIUTestSuite) TestCreateBucketFsAPIFailsCreatingSchema() {
	suite.dbMock.ExpectBegin()
	suite.dbMock.ExpectExec("CREATE SCHEMA INTERNAL_\\d+").WillReturnError(fmt.Errorf("mock error"))
	suite.dbMock.ExpectRollback()
	client, err := CreateBucketFsAPI(BUCKETFS_BASE_PATH, context.Background(), suite.db)
	suite.EqualError(err, "failed to create a schema for BucketFS list script. Cause: mock error")
	suite.Nil(client)
}

func (suite *BucketFsAPIUTestSuite) TestCreateBucketFsAPIFailsCreatingUDFScript() {
	suite.dbMock.ExpectBegin()
	suite.dbMock.ExpectExec("CREATE SCHEMA INTERNAL_\\d+").WillReturnResult(sqlmock.NewResult(0, 1))
	suite.dbMock.ExpectExec("(?m)CREATE OR REPLACE PYTHON3 SCALAR SCRIPT.*").WillReturnError(fmt.Errorf("mock error"))
	suite.dbMock.ExpectRollback()
	client, err := CreateBucketFsAPI(BUCKETFS_BASE_PATH, context.Background(), suite.db)
	suite.EqualError(err, "failed to create UDF script for listing bucket. Cause: mock error")
	suite.Nil(client)
}

// ListFiles

/* [utest -> dsn~configure-bucketfs-path~1]. */
func (suite *BucketFsAPIUTestSuite) TestListFiles() {
	client := suite.createBucketFsClientHandleError()
	suite.dbMock.ExpectPrepare(`SELECT "INTERNAL_.* ORDER BY FULL_PATH`).
		WillBeClosed().
		ExpectQuery().WithArgs(BUCKETFS_BASE_PATH).WillReturnRows(sqlmock.NewRows([]string{"FILE_NAME", "FULL_PATH", "SIZE"}).
		AddRow("file1.txt", "/base/file1.txt", 10).
		AddRow("file2.txt", "/base2/file2.txt", 20)).
		RowsWillBeClosed()
	result, err := client.ListFiles()
	suite.NoError(err)
	suite.Equal([]BfsFile{{Name: "file1.txt", Path: "/base/file1.txt", Size: 10}, {Name: "file2.txt", Path: "/base2/file2.txt", Size: 20}}, result)
}

func (suite *BucketFsAPIUTestSuite) TestListFilesPrepareQueryFails() {
	client := suite.createBucketFsClientHandleError()
	suite.dbMock.ExpectPrepare(`SELECT "INTERNAL_.* ORDER BY FULL_PATH`).WillReturnError(fmt.Errorf("mock error"))
	result, err := client.ListFiles()
	suite.EqualError(err, "failed to create prepared statement for running list files UDF. Cause: mock error")
	suite.Empty(result)
}

func (suite *BucketFsAPIUTestSuite) TestListFilesExecuteQueryFails() {
	client := suite.createBucketFsClientHandleError()
	suite.dbMock.ExpectPrepare(`SELECT "INTERNAL_.* ORDER BY FULL_PATH`).
		WillBeClosed().
		ExpectQuery().WillReturnError(fmt.Errorf("mock error"))
	result, err := client.ListFiles()
	suite.EqualError(err, "failed to list files in BucketFS using UDF. Cause: mock error")
	suite.Empty(result)
}

func (suite *BucketFsAPIUTestSuite) TestListFilesWrongResultColumnCount() {
	client := suite.createBucketFsClientHandleError()
	suite.dbMock.ExpectPrepare(`SELECT "INTERNAL_.* ORDER BY FULL_PATH`).
		WillBeClosed().
		ExpectQuery().WillReturnRows(sqlmock.NewRows([]string{"FILE_NAME", "FULL_PATH"}).
		AddRow("file1.txt", "/base/file1.txt").
		AddRow("file2.txt", "/base2/file2.txt")).
		RowsWillBeClosed()
	result, err := client.ListFiles()
	suite.EqualError(err, "failed reading result of BucketFS list UDF. Cause: sql: expected 2 destination arguments in Scan, not 3")
	suite.Empty(result)
}

// FindAbsolutePath

/* [utest -> dsn~configure-bucketfs-path~1] */
/* [utest -> dsn~resolving-files-in-bucketfs~1]. */
func (suite *BucketFsAPIUTestSuite) TestFindAbsolutePath() {
	client := suite.createBucketFsClientHandleError()
	suite.dbMock.ExpectPrepare(`SELECT FULL_PATH FROM.*`).
		WillBeClosed().
		ExpectQuery().WithArgs(BUCKETFS_BASE_PATH, FILE_NAME).WillReturnRows(sqlmock.NewRows([]string{"FULL_PATH"}).AddRow("/abs/path/file.txt")).
		RowsWillBeClosed()
	result, err := client.FindAbsolutePath(FILE_NAME)
	suite.NoError(err)
	suite.Equal("/abs/path/file.txt", result)
}

func (suite *BucketFsAPIUTestSuite) TestFindAbsolutePathPrepareQueryFails() {
	client := suite.createBucketFsClientHandleError()
	suite.dbMock.ExpectPrepare(`SELECT FULL_PATH FROM.*`).WillReturnError(fmt.Errorf("mock error"))
	result, err := client.FindAbsolutePath(FILE_NAME)
	suite.EqualError(err, "failed to create prepared statement for running list files UDF. Cause: mock error")
	suite.Equal("", result)
}

func (suite *BucketFsAPIUTestSuite) TestFindAbsolutePathExecuteQueryFails() {
	client := suite.createBucketFsClientHandleError()
	suite.dbMock.ExpectPrepare(`SELECT FULL_PATH FROM.*`).
		WillBeClosed().
		ExpectQuery().WithArgs(BUCKETFS_BASE_PATH, FILE_NAME).WillReturnError(fmt.Errorf("mock error"))
	result, err := client.FindAbsolutePath(FILE_NAME)
	suite.EqualError(err, "failed to find absolute path in BucketFS using UDF. Cause: mock error")
	suite.Equal("", result)
}

func (suite *BucketFsAPIUTestSuite) TestFindAbsolutePathNoResult() {
	client := suite.createBucketFsClientHandleError()
	suite.dbMock.ExpectPrepare(`SELECT FULL_PATH FROM.*`).
		WillBeClosed().
		ExpectQuery().WithArgs(BUCKETFS_BASE_PATH, FILE_NAME).WillReturnRows(sqlmock.NewRows([]string{"FULL_PATH"})).RowsWillBeClosed()
	result, err := client.FindAbsolutePath(FILE_NAME)
	suite.EqualError(err, `file "file.txt" not found in BucketFS`)
	suite.Equal("", result)
}

func (suite *BucketFsAPIUTestSuite) TestFindAbsolutePathRowError() {
	client := suite.createBucketFsClientHandleError()
	suite.dbMock.ExpectPrepare(`SELECT FULL_PATH FROM.*`).
		WillBeClosed().
		ExpectQuery().WithArgs(BUCKETFS_BASE_PATH, FILE_NAME).WillReturnRows(sqlmock.NewRows([]string{"FULL_PATH"}).AddRow("/abs/path/file.txt").RowError(0, fmt.Errorf("mock error"))).
		RowsWillBeClosed()
	result, err := client.FindAbsolutePath(FILE_NAME)
	suite.EqualError(err, "failed iterating absolute path results. Cause: mock error")
	suite.Equal("", result)
}

func (suite *BucketFsAPIUTestSuite) TestFindAbsolutePathWrongResultColumnCount() {
	client := suite.createBucketFsClientHandleError()
	suite.dbMock.ExpectPrepare(`SELECT FULL_PATH FROM.*`).
		WillBeClosed().
		ExpectQuery().WithArgs(BUCKETFS_BASE_PATH, FILE_NAME).WillReturnRows(sqlmock.NewRows([]string{"FULL_PATH", "Additional Column"}).AddRow("/abs/path/file.txt", "a")).
		RowsWillBeClosed()
	result, err := client.FindAbsolutePath(FILE_NAME)
	suite.EqualError(err, `failed reading absolute path. Cause: sql: expected 2 destination arguments in Scan, not 1`)
	suite.Equal("", result)
}

func (suite *BucketFsAPIUTestSuite) createBucketFsClientHandleError() BucketFsAPI {
	bfsClient, err := suite.createBucketFsClient()
	if err != nil {
		suite.FailNow("Creating BFS API failed: " + err.Error())
	}
	return bfsClient
}

func (suite *BucketFsAPIUTestSuite) createBucketFsClient() (BucketFsAPI, error) {
	suite.dbMock.ExpectBegin()
	suite.dbMock.ExpectExec("CREATE SCHEMA INTERNAL_\\d+").WillReturnResult(sqlmock.NewResult(0, 1))
	suite.dbMock.ExpectExec("(?m)CREATE OR REPLACE PYTHON3 SCALAR SCRIPT.*").WillReturnResult(sqlmock.NewResult(0, 1))
	return CreateBucketFsAPI(BUCKETFS_BASE_PATH, context.Background(), suite.db)
}
