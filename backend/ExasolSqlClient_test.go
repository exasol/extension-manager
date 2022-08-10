package backend

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/suite"
)

type ExasolSqlClientTestSuite struct {
	suite.Suite
	db     *sql.DB
	dbMock sqlmock.Sqlmock
}

func TestExasolSqlClient(t *testing.T) {
	suite.Run(t, new(ExasolSqlClientTestSuite))
}

func (suite *ExasolSqlClientTestSuite) SetupTest() {
	db, mock, err := sqlmock.New()
	if err != nil {
		suite.Failf("an error '%v' was not expected when opening a stub database connection", err.Error())
	}
	suite.db = db
	suite.dbMock = mock
}

func (suite *ExasolSqlClientTestSuite) TestRun_succeeds() {
	client := NewSqlClient(suite.db)
	suite.dbMock.ExpectExec("select 1").WillReturnResult(sqlmock.NewResult(1, 1))
	client.RunQuery("select 1")
}

func (suite *ExasolSqlClientTestSuite) TestRun_fails() {
	client := NewSqlClient(suite.db)
	suite.dbMock.ExpectExec("invalid").WillReturnError(fmt.Errorf("expected"))
	suite.Panics(func() { client.RunQuery("invalid") })
}
