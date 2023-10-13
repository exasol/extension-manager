package transaction

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/suite"
)

const BUCKETFS_BASE_PATH = "bucketfs-base-path"

var mockError = fmt.Errorf("mock error")

type TransactionContextSuite struct {
	suite.Suite
	db     *sql.DB
	dbMock sqlmock.Sqlmock
}

func TestTransactionContextSuite(t *testing.T) {
	suite.Run(t, new(TransactionContextSuite))
}

func (suite *TransactionContextSuite) SetupTest() {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
	suite.NoError(err)
	suite.db = db
	suite.dbMock = mock
	suite.dbMock.MatchExpectationsInOrder(true)
}

func (suite *TransactionContextSuite) AfterTest(suiteName, testName string) {
	suite.NoError(suite.dbMock.ExpectationsWereMet())
}

func (suite *TransactionContextSuite) TestBeginTransaction() {
	suite.dbMock.ExpectBegin()
	txCtx, err := suite.beginTransaction()
	suite.NoError(err)
	suite.NotNil(txCtx)
}

func (suite *TransactionContextSuite) TestBeginTransactionFails() {
	suite.dbMock.ExpectBegin().WillReturnError(mockError)
	txCtx, err := suite.beginTransaction()
	suite.EqualError(err, "failed to begin transaction: mock error")
	suite.Nil(txCtx)
}

func (suite *TransactionContextSuite) TestBeginTransactionAuthenticationFails() {
	suite.dbMock.ExpectBegin().WillReturnError(fmt.Errorf("mock error: 'Connection exception - authentication failed'"))
	txCtx, err := suite.beginTransaction()
	suite.EqualError(err, "invalid database credentials")
	suite.Nil(txCtx)
}

func (suite *TransactionContextSuite) TestGetContext() {
	suite.dbMock.ExpectBegin()
	txCtx, _ := suite.beginTransaction()
	suite.NotNil(txCtx.GetContext())
}

func (suite *TransactionContextSuite) TestGetDBConnection() {
	suite.dbMock.ExpectBegin()
	txCtx, _ := suite.beginTransaction()
	suite.NotNil(txCtx.GetDBConnection())
}

func (suite *TransactionContextSuite) TestGetTransaction() {
	suite.dbMock.ExpectBegin()
	txCtx, _ := suite.beginTransaction()
	suite.NotNil(txCtx.GetTransaction())
}

// GetBucketFsClient()

func (suite *TransactionContextSuite) TestGetBucketFsClient() {
	suite.dbMock.ExpectBegin()
	txCtx, _ := suite.beginTransaction()
	suite.dbMock.ExpectBegin()
	suite.dbMock.ExpectExec("CREATE SCHEMA INTERNAL_\\d+").WillReturnResult(sqlmock.NewResult(0, 1))
	suite.dbMock.ExpectExec("(?m)CREATE OR REPLACE PYTHON3 SCALAR SCRIPT.*").WillReturnResult(sqlmock.NewResult(0, 1))
	bfsClient, err := txCtx.GetBucketFsClient()
	suite.NoError(err)
	suite.NotNil(bfsClient)
}

func (suite *TransactionContextSuite) TestGetBucketFsClientTwiceReturnsSameObject() {
	suite.dbMock.ExpectBegin()
	txCtx, _ := suite.beginTransaction()
	suite.dbMock.ExpectBegin()
	suite.dbMock.ExpectExec("CREATE SCHEMA INTERNAL_\\d+").WillReturnResult(sqlmock.NewResult(0, 1))
	suite.dbMock.ExpectExec("(?m)CREATE OR REPLACE PYTHON3 SCALAR SCRIPT.*").WillReturnResult(sqlmock.NewResult(0, 1))
	bfsClient1, err := txCtx.GetBucketFsClient()
	suite.NoError(err)
	suite.NotNil(bfsClient1)

	bfsClient2, err := txCtx.GetBucketFsClient()
	suite.NoError(err)
	suite.NotNil(bfsClient2)

	suite.Same(bfsClient1, bfsClient2)
}

func (suite *TransactionContextSuite) TestGetBucketFsClientFailsCreatingSchema() {
	suite.dbMock.ExpectBegin()
	txCtx, _ := suite.beginTransaction()
	suite.dbMock.ExpectBegin()
	suite.dbMock.ExpectExec("CREATE SCHEMA INTERNAL_\\d+").WillReturnError(mockError)
	bfsClient, err := txCtx.GetBucketFsClient()
	suite.EqualError(err, "failed to create a schema for BucketFS list script. Cause: mock error")
	suite.Nil(bfsClient)
}

// Rollback()

func (suite *TransactionContextSuite) TestRollback() {
	suite.dbMock.ExpectBegin()
	txCtx, _ := suite.beginTransaction()
	suite.dbMock.ExpectRollback()
	txCtx.Rollback()
}

func (suite *TransactionContextSuite) TestRollbackClosesBfsClient() {
	suite.dbMock.ExpectBegin()
	txCtx, _ := suite.beginTransaction()

	suite.dbMock.ExpectBegin()
	suite.dbMock.ExpectExec("CREATE SCHEMA INTERNAL_\\d+").WillReturnResult(sqlmock.NewResult(0, 1))
	suite.dbMock.ExpectExec("(?m)CREATE OR REPLACE PYTHON3 SCALAR SCRIPT.*").WillReturnResult(sqlmock.NewResult(0, 1))
	_, err := txCtx.GetBucketFsClient()
	suite.NoError(err)

	suite.dbMock.ExpectRollback() // Rollback from BFS client
	suite.dbMock.ExpectRollback() // Rollback transaction
	txCtx.Rollback()
}

func (suite *TransactionContextSuite) TestRollbackClosingBfsClientFails() {
	suite.dbMock.ExpectBegin()
	txCtx, _ := suite.beginTransaction()

	suite.dbMock.ExpectBegin()
	suite.dbMock.ExpectExec("CREATE SCHEMA INTERNAL_\\d+").WillReturnResult(sqlmock.NewResult(0, 1))
	suite.dbMock.ExpectExec("(?m)CREATE OR REPLACE PYTHON3 SCALAR SCRIPT.*").WillReturnResult(sqlmock.NewResult(0, 1))
	_, err := txCtx.GetBucketFsClient()
	suite.NoError(err)

	suite.dbMock.ExpectRollback().WillReturnError(mockError) // Rollback from BFS client
	suite.dbMock.ExpectRollback()
	txCtx.Rollback()
}

// Commit()

func (suite *TransactionContextSuite) TestCommit() {
	suite.dbMock.ExpectBegin()
	txCtx, _ := suite.beginTransaction()
	suite.dbMock.ExpectCommit()
	suite.NoError(txCtx.Commit())
}

func (suite *TransactionContextSuite) TestCommitClosesBfsClient() {
	suite.dbMock.ExpectBegin()
	txCtx, _ := suite.beginTransaction()

	suite.dbMock.ExpectBegin()
	suite.dbMock.ExpectExec("CREATE SCHEMA INTERNAL_\\d+").WillReturnResult(sqlmock.NewResult(0, 1))
	suite.dbMock.ExpectExec("(?m)CREATE OR REPLACE PYTHON3 SCALAR SCRIPT.*").WillReturnResult(sqlmock.NewResult(0, 1))
	_, err := txCtx.GetBucketFsClient()
	suite.NoError(err)

	suite.dbMock.ExpectRollback() // Rollback from BFS client
	suite.dbMock.ExpectCommit()   // Commit transaction
	suite.NoError(txCtx.Commit())
}

func (suite *TransactionContextSuite) TestCommitClosingBfsClientFails() {
	suite.dbMock.ExpectBegin()
	txCtx, _ := suite.beginTransaction()

	suite.dbMock.ExpectBegin()
	suite.dbMock.ExpectExec("CREATE SCHEMA INTERNAL_\\d+").WillReturnResult(sqlmock.NewResult(0, 1))
	suite.dbMock.ExpectExec("(?m)CREATE OR REPLACE PYTHON3 SCALAR SCRIPT.*").WillReturnResult(sqlmock.NewResult(0, 1))
	_, err := txCtx.GetBucketFsClient()
	suite.NoError(err)

	suite.dbMock.ExpectRollback().WillReturnError(mockError) // Rollback from BFS client
	suite.EqualError(txCtx.Commit(), "failed to close BucketFS client: failed to rollback transaction to cleanup resources. Cause: mock error")
}

func (suite *TransactionContextSuite) beginTransaction() (*TransactionContext, error) {
	return BeginTransaction(context.Background(), suite.db, BUCKETFS_BASE_PATH)
}
