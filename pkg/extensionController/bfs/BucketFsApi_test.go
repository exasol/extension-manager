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
func (suite *BucketFsAPIUTestSuite) TestListFiles() {
	suite.dbMock.ExpectBegin()
	suite.dbMock.ExpectExec("CREATE SCHEMA INTERNAL_\\d+").WillReturnResult(sqlmock.NewResult(0, 1))
	suite.dbMock.ExpectExec("(?m)CREATE OR REPLACE PYTHON3 SCALAR SCRIPT.*").WillReturnResult(sqlmock.NewResult(0, 1))
	suite.dbMock.ExpectPrepare(`SELECT "INTERNAL_.* ORDER BY FULL_PATH`).
		WillBeClosed().
		ExpectQuery().WillReturnRows(sqlmock.NewRows([]string{"FILE_NAME", "FULL_PATH", "SIZE"}).
		AddRow("file1.txt", "/base/file1.txt", 10).
		AddRow("file2.txt", "/base2/file2.txt", 20)).
		RowsWillBeClosed()
	suite.dbMock.ExpectRollback()
	result, err := suite.listFiles()
	suite.NoError(err)
	suite.Equal([]BfsFile{{Name: "file1.txt", Path: "/base/file1.txt", Size: 10}, {Name: "file2.txt", Path: "/base2/file2.txt", Size: 20}}, result)
}

func (suite *BucketFsAPIUTestSuite) TestListFilesBeginTransactionFails() {
	suite.dbMock.ExpectBegin().WillReturnError(fmt.Errorf("mock error"))
	suite.dbMock.ExpectRollback()
	result, err := suite.listFiles()
	suite.EqualError(err, "failed to create a transaction. Cause: mock error")
	suite.Empty(result)
}

func (suite *BucketFsAPIUTestSuite) TestListFilesCreateSchemaFails() {
	suite.dbMock.ExpectBegin()
	suite.dbMock.ExpectExec("CREATE SCHEMA INTERNAL_\\d+").WillReturnError(fmt.Errorf("mock error"))
	suite.dbMock.ExpectRollback()
	result, err := suite.listFiles()
	suite.EqualError(err, "failed to create a schema for BucketFS list script. Cause: mock error")
	suite.Empty(result)
}

// TODO

func (suite *BucketFsAPIUTestSuite) listFiles() ([]BfsFile, error) {
	return CreateBucketFsAPI(BUCKETFS_BASE_PATH).ListFiles(context.Background(), suite.db)
}
