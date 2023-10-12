package backend

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/suite"
)

type ExasolSqlClientUTestSuite struct {
	suite.Suite
	db     *sql.DB
	dbMock sqlmock.Sqlmock
}

func TestExasolSqlClientUTestSuite(t *testing.T) {
	suite.Run(t, new(ExasolSqlClientUTestSuite))
}

func (suite *ExasolSqlClientUTestSuite) SetupTest() {
	db, mock, err := sqlmock.New()
	if err != nil {
		suite.Failf("an error '%v' was not expected when opening a stub database connection", err.Error())
	}
	suite.db = db
	suite.dbMock = mock
	suite.dbMock.MatchExpectationsInOrder(true)
}

func (suite *ExasolSqlClientUTestSuite) createClient() SimpleSQLClient {
	return NewSqlClient(context.Background(), suite.createMockTransaction())
}

func (suite *ExasolSqlClientUTestSuite) createMockTransaction() *sql.Tx {
	suite.dbMock.ExpectBegin()
	tx, err := suite.db.Begin()
	suite.NoError(err)
	return tx
}

/* [utest -> dsn~extension-context-sql-client~1]. */
func (suite *ExasolSqlClientUTestSuite) TestExecuteSucceeds() {
	client := suite.createClient()
	suite.dbMock.ExpectExec("select 1").WillReturnResult(sqlmock.NewResult(1, 1))
	result, err := client.Execute("select 1")
	suite.NoError(err)
	suite.NotNil(result)
}

func (suite *ExasolSqlClientUTestSuite) TestExecuteFails() {
	client := suite.createClient()
	suite.dbMock.ExpectExec("invalid").WillReturnError(fmt.Errorf("expected"))
	result, err := client.Execute("invalid")
	suite.EqualError(err, "error executing statement \"invalid\": expected")
	suite.Nil(result)
}

var forbiddenCommandTests = []struct {
	statement        string
	forbiddenCommand string
}{
	{"select 1", ""},
	{"com mit", ""},
	{"roll back", ""},
	{"do rollback", ""},
	{"do commit", ""},
	{"commit", "commit"},
	{"rollback", "rollback"},
	{"COMMIT", "commit"},
	{"ROLLBACK", "rollback"},
	{" commit; ", "commit"},
	{"\t\r\n ; commit \t\r\n ; ", "commit"},
	{"\t\r\n ; COMMIT \t\r\n ; ", "commit"}}

func (suite *ExasolSqlClientUTestSuite) TestExecuteValidation() {
	for _, test := range forbiddenCommandTests {
		suite.Run(fmt.Sprintf("running statement %q contains forbidden command %q", test.statement, test.forbiddenCommand), func() {
			client := suite.createClient()
			if test.forbiddenCommand != "" {
				expectedError := fmt.Sprintf("statement %q contains forbidden command %q. Transaction handling is done by extension manager", test.statement, test.forbiddenCommand)
				result, err := client.Execute(test.statement)
				suite.EqualError(err, expectedError)
				suite.Nil(result)
			} else {
				suite.dbMock.ExpectExec(test.statement).WillReturnResult(sqlmock.NewResult(1, 0))
				result, err := client.Execute(test.statement)
				suite.NoError(err)
				suite.NotNil(result)
			}
		})
	}
}

func (suite *ExasolSqlClientUTestSuite) TestQueryFails() {
	client := suite.createClient()
	suite.dbMock.ExpectQuery("invalid").WillReturnError(fmt.Errorf("expected")).RowsWillBeClosed()
	result, err := client.Query("invalid")
	suite.EqualError(err, "error executing statement \"invalid\": expected")
	suite.Nil(result)
}

/* [utest -> dsn~extension-context-sql-client~1]. */
func (suite *ExasolSqlClientUTestSuite) TestQuerySucceeds() {
	client := suite.createClient()
	suite.dbMock.ExpectQuery("query").WillReturnRows(sqlmock.NewRows([]string{"col1", "col2"}).AddRow(1, "a").AddRow(2, "b")).RowsWillBeClosed()
	result, err := client.Query("query")
	suite.NoError(err)
	suite.NotNil(result)
	suite.Equal(&QueryResult{
		Columns: []Column{{Name: "col1", TypeName: "type"}, {Name: "col2", TypeName: "type"}},
		Rows:    []Row{{int64(1), "a"}, {int64(2), "b"}}}, result)
}

func (suite *ExasolSqlClientUTestSuite) TestQueryRowFails() {
	client := suite.createClient()
	suite.dbMock.ExpectQuery("query").WillReturnRows(sqlmock.NewRows([]string{"col1", "col2"}).AddRow(2, "b").RowError(0, fmt.Errorf("mock"))).RowsWillBeClosed()
	result, err := client.Query("query")
	suite.EqualError(err, "error while iterating result: mock")
	suite.Nil(result)
}

func (suite *ExasolSqlClientUTestSuite) TestQueryCloseFails() {
	client := suite.createClient()
	suite.dbMock.ExpectQuery("query").WillReturnRows(sqlmock.NewRows([]string{"col1", "col2"}).CloseError(fmt.Errorf("mock error"))).RowsWillBeClosed()
	result, err := client.Query("query")
	suite.EqualError(err, "error while iterating result: mock error")
	suite.Nil(result)
}

func (suite *ExasolSqlClientUTestSuite) TestQueryValidation() {
	for _, test := range forbiddenCommandTests {
		suite.Run(fmt.Sprintf("running statement %q contains forbidden command %q", test.statement, test.forbiddenCommand), func() {
			client := suite.createClient()
			if test.forbiddenCommand != "" {
				expectedError := fmt.Sprintf("statement %q contains forbidden command %q. Transaction handling is done by extension manager", test.statement, test.forbiddenCommand)
				result, err := client.Query(test.statement)
				suite.EqualError(err, expectedError)
				suite.Nil(result)
			} else {
				suite.dbMock.ExpectQuery(test.statement).WillReturnRows(sqlmock.NewRows([]string{"col1"})).RowsWillBeClosed()
				result, err := client.Query(test.statement)
				suite.NoError(err)
				suite.NotNil(result)
			}
		})
	}
}
