package backend

import (
	"context"
	"testing"

	"github.com/exasol/extension-manager/integrationTesting"
	"github.com/stretchr/testify/suite"
)

type ExasolSqlClientITestSuite struct {
	suite.Suite
	client *ExasolSqlClient
	exasol *integrationTesting.DbTestSetup
}

func TestExasolSqlClientITestSuite(t *testing.T) {
	suite.Run(t, new(ExasolSqlClientITestSuite))
}

func (suite *ExasolSqlClientITestSuite) SetupSuite() {
	suite.exasol = integrationTesting.StartDbSetup(&suite.Suite)
}

func (suite *ExasolSqlClientITestSuite) TearDownSuite() {
	suite.exasol.StopDb()
}

func (suite *ExasolSqlClientITestSuite) SetupTest() {
	suite.exasol.CreateConnection()
	suite.T().Cleanup(func() {
		suite.exasol.CloseConnection()
	})
	tx, err := suite.exasol.GetConnection().Begin()
	if err != nil {
		suite.T().Fatalf("failed to begin DB connection: %v", err)
	}
	suite.T().Cleanup(func() {
		err := tx.Rollback()
		suite.NoError(err, "failed to rollback transaction")
	})
	suite.client = NewSqlClient(context.Background(), tx)
}

func (suite *ExasolSqlClientITestSuite) TestExecute_Succeeds() {
	suite.NotPanics(func() { suite.client.Execute("select 1") })
}

func (suite *ExasolSqlClientITestSuite) TestExecute_WithArgumentSucceeds() {
	suite.NotPanics(func() { suite.client.Execute("select 1 from dual where 1 = ?", 1) })
}

func (suite *ExasolSqlClientITestSuite) TestExecute_WithArgsSucceeds() {
	suite.NotPanics(func() { suite.client.Execute("select 1 from dual where 1 = ?", 1) })
}

func (suite *ExasolSqlClientITestSuite) TestExecute_fails() {
	suite.Panics(func() { suite.client.Execute("invalid") })
}

func (suite *ExasolSqlClientITestSuite) TestQuery_fails() {
	suite.Panics(func() { suite.client.Query("invalid") })
}

func (suite *ExasolSqlClientITestSuite) TestQuery_Succeeds() {
	result := suite.client.Query("select 1 as col")
	suite.Equal(QueryResult{Columns: []Column{{Name: "COL", TypeName: "DECIMAL"}}, Rows: []Row{{(float64(1))}}}, result)
}

func (suite *ExasolSqlClientITestSuite) TestQuery_WithArgument() {
	result := suite.client.Query("select 1 as col from dual where 1 = ?", 1)
	suite.Equal(QueryResult{Columns: []Column{{Name: "COL", TypeName: "DECIMAL"}}, Rows: []Row{{(float64(1))}}}, result)
}

func (suite *ExasolSqlClientITestSuite) TestQuery_NoRow() {
	result := suite.client.Query("select 1 as col from dual where 1=2")
	suite.Equal(QueryResult{Columns: []Column{{Name: "COL", TypeName: "DECIMAL"}}, Rows: []Row{}}, result)
}

func (suite *ExasolSqlClientITestSuite) TestQuery_MultipleRows() {
	result := suite.client.Query("select t.* from (values (1, 'a'), (2, 'b'), (3, 'c')) as t(num, txt)")
	suite.Equal(QueryResult{
		Columns: []Column{{Name: "NUM", TypeName: "DECIMAL"}, {Name: "TXT", TypeName: "VARCHAR"}},
		Rows:    []Row{{float64(1), "a"}, {float64(2), "b"}, {float64(3), "c"}},
	}, result)
}
