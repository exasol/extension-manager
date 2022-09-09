package backend

import (
	"context"
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

func (suite *ExasolSqlClientTestSuite) createClient() *ExasolSqlClient {
	return NewSqlClient(context.Background(), suite.createMockTransaction())
}

func (suite *ExasolSqlClientTestSuite) TestRun_succeeds() {
	client := suite.createClient()
	suite.dbMock.ExpectExec("select 1").WillReturnResult(sqlmock.NewResult(1, 1))
	client.Execute("select 1")
}

func (suite *ExasolSqlClientTestSuite) TestRun_fails() {
	client := suite.createClient()
	suite.dbMock.ExpectExec("invalid").WillReturnError(fmt.Errorf("expected"))
	suite.PanicsWithError("error executing statement \"invalid\": expected", func() { client.Execute("invalid") })
}

func (suite *ExasolSqlClientTestSuite) TestRun_validation() {
	var tests = []struct {
		statement        string
		forbiddenCommand string
	}{{"select 1", ""}, {"com mit", ""}, {"roll back", ""},
		{"commit", "commit"}, {"rollback", "rollback"}, {"COMMIT", "commit"}, {"ROLLBACK", "rollback"},
		{" commit; ", "commit"}, {"\t\r\n ; commit \t\r\n ; ", "commit"}, {"\t\r\n ; COMMIT \t\r\n ; ", "commit"}}
	for _, test := range tests {
		suite.Run(fmt.Sprintf("running statement %q contains forbidden command %q", test.statement, test.forbiddenCommand), func() {
			client := suite.createClient()
			if test.forbiddenCommand != "" {
				expectedError := fmt.Sprintf("statement %q contains forbidden command %q. Transaction handling is done by extension manager", test.statement, test.forbiddenCommand)
				suite.PanicsWithError(expectedError, func() { client.Execute(test.statement) })
			} else {
				suite.dbMock.ExpectExec(test.statement).WillReturnResult(sqlmock.NewResult(1, 0))
				client.Execute(test.statement)
			}
		})
	}
}

func (suite *ExasolSqlClientTestSuite) createMockTransaction() *sql.Tx {
	suite.dbMock.ExpectBegin()
	tx, err := suite.db.Begin()
	suite.NoError(err)
	return tx
}
