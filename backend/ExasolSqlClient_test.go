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
}

func (suite *ExasolSqlClientUTestSuite) createClient() *ExasolSqlClient {
	return NewSqlClient(context.Background(), suite.createMockTransaction())
}

func (suite *ExasolSqlClientUTestSuite) createMockTransaction() *sql.Tx {
	suite.dbMock.ExpectBegin()
	tx, err := suite.db.Begin()
	suite.NoError(err)
	return tx
}

func (suite *ExasolSqlClientUTestSuite) TestExecute_succeeds() {
	client := suite.createClient()
	suite.dbMock.ExpectExec("select 1").WillReturnResult(sqlmock.NewResult(1, 1))
	client.Execute("select 1")
}

func (suite *ExasolSqlClientUTestSuite) TestExecute_fails() {
	client := suite.createClient()
	suite.dbMock.ExpectExec("invalid").WillReturnError(fmt.Errorf("expected"))
	suite.PanicsWithError("error executing statement \"invalid\": expected", func() { client.Execute("invalid") })
}

var forbiddenCommandTests = []struct {
	statement        string
	forbiddenCommand string
}{{"select 1", ""}, {"com mit", ""}, {"roll back", ""},
	{"commit", "commit"}, {"rollback", "rollback"}, {"COMMIT", "commit"}, {"ROLLBACK", "rollback"},
	{" commit; ", "commit"}, {"\t\r\n ; commit \t\r\n ; ", "commit"}, {"\t\r\n ; COMMIT \t\r\n ; ", "commit"}}

func (suite *ExasolSqlClientUTestSuite) TestExecute_validation() {
	for _, test := range forbiddenCommandTests {
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

func (suite *ExasolSqlClientUTestSuite) TestQuery_fails() {
	client := suite.createClient()
	suite.dbMock.ExpectQuery("invalid").WillReturnError(fmt.Errorf("expected"))
	suite.PanicsWithError("error executing statement \"invalid\": expected", func() { client.Query("invalid") })
}

func (suite *ExasolSqlClientUTestSuite) TestQuery_succeeds() {
	client := suite.createClient()
	suite.dbMock.ExpectQuery("query").WillReturnRows(sqlmock.NewRows([]string{"col1", "col2"}).AddRow(1, "a").AddRow(2, "b"))
	result := client.Query("query")
	suite.NotNil(result)
	suite.Equal(QueryResult{
		Columns: []Column{{Name: "col1"}, {Name: "col2"}},
		Rows:    []Row{{int64(1), "a"}, {int64(2), "b"}}}, result)
}

func (suite *ExasolSqlClientUTestSuite) TestQuery_rowFails() {
	client := suite.createClient()
	suite.dbMock.ExpectQuery("query").WillReturnRows(sqlmock.NewRows([]string{"col1", "col2"}).AddRow(2, "b").RowError(0, fmt.Errorf("mock")))
	suite.PanicsWithError("error while iterating result: mock", func() { client.Query("query") })
}

func (suite *ExasolSqlClientUTestSuite) TestQuery_validation() {
	for _, test := range forbiddenCommandTests {
		suite.Run(fmt.Sprintf("running statement %q contains forbidden command %q", test.statement, test.forbiddenCommand), func() {
			client := suite.createClient()
			if test.forbiddenCommand != "" {
				expectedError := fmt.Sprintf("statement %q contains forbidden command %q. Transaction handling is done by extension manager", test.statement, test.forbiddenCommand)
				suite.PanicsWithError(expectedError, func() { client.Query(test.statement) })
			} else {
				suite.dbMock.ExpectQuery(test.statement).WillReturnRows(sqlmock.NewRows([]string{"col1"}))
				result := client.Query(test.statement)
				suite.NotNil(result)
			}
		})
	}
}
