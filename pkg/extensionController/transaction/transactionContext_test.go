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
	suite.dbMock.ExpectBegin().WillReturnError(fmt.Errorf("mock error"))
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

func (suite *TransactionContextSuite) TestRollback() {
	suite.dbMock.ExpectBegin()
	txCtx, _ := suite.beginTransaction()
	suite.dbMock.ExpectRollback()
	txCtx.Rollback()
}

func (suite *TransactionContextSuite) TestCommit() {
	suite.dbMock.ExpectBegin()
	txCtx, _ := suite.beginTransaction()
	suite.dbMock.ExpectCommit()
	suite.NoError(txCtx.Commit())
}

func (suite *TransactionContextSuite) beginTransaction() (*TransactionContext, error) {
	return BeginTransaction(context.Background(), suite.db, BUCKETFS_BASE_PATH)
}
