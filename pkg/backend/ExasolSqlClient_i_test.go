package backend_test

import (
	"context"
	"testing"

	"github.com/exasol/extension-manager/pkg/backend"
	"github.com/exasol/extension-manager/pkg/integrationTesting"
	"github.com/stretchr/testify/suite"
)

type ExasolSqlClientITestSuite struct {
	suite.Suite
	client backend.SimpleSQLClient
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
	suite.client = backend.NewSqlClient(context.Background(), tx)
}

func (suite *ExasolSqlClientITestSuite) TestExecuteSucceeds() {
	result, err := suite.client.Execute("select 1")
	suite.NoError(err)
	suite.NotNil(result)
}

func (suite *ExasolSqlClientITestSuite) TestExecuteWithArgumentSucceeds() {
	result, err := suite.client.Execute("select 1 from dual where 1 = ?", 1)
	suite.NoError(err)
	suite.NotNil(result)
}

func (suite *ExasolSqlClientITestSuite) TestExecuteFails() {
	result, err := suite.client.Execute("invalid")
	suite.ErrorContains(err, "error executing statement \"invalid\": E-EGOD-11: execution failed")
	suite.Nil(result)
}

func (suite *ExasolSqlClientITestSuite) TestQueryFails() {
	result, err := suite.client.Query("invalid")
	suite.ErrorContains(err, "error executing statement \"invalid\": E-EGOD-11: execution failed")
	suite.Nil(result)
}

func (suite *ExasolSqlClientITestSuite) TestQuerySucceeds() {
	result, err := suite.client.Query("select 1 as col")
	suite.NoError(err)
	suite.Equal(&backend.QueryResult{Columns: []backend.Column{{Name: "COL", TypeName: "DECIMAL"}}, Rows: []backend.Row{{(float64(1))}}}, result)
}

func (suite *ExasolSqlClientITestSuite) TestQueryWithArgument() {
	result, err := suite.client.Query("select 1 as col from dual where 1 = ?", 1)
	suite.NoError(err)
	suite.Equal(&backend.QueryResult{Columns: []backend.Column{{Name: "COL", TypeName: "DECIMAL"}}, Rows: []backend.Row{{(float64(1))}}}, result)
}

func (suite *ExasolSqlClientITestSuite) TestQueryNoRow() {
	result, err := suite.client.Query("select 1 as col from dual where 1=2")
	suite.NoError(err)
	suite.Equal(&backend.QueryResult{Columns: []backend.Column{{Name: "COL", TypeName: "DECIMAL"}}, Rows: []backend.Row{}}, result)
}

func (suite *ExasolSqlClientITestSuite) TestQueryMultipleRows() {
	result, err := suite.client.Query("select t.* from (values (1, 'a'), (2, 'b'), (3, 'c')) as t(num, txt)")
	suite.NoError(err)
	suite.Equal(&backend.QueryResult{
		Columns: []backend.Column{{Name: "NUM", TypeName: "DECIMAL"}, {Name: "TXT", TypeName: "VARCHAR"}},
		Rows:    []backend.Row{{float64(1), "a"}, {float64(2), "b"}, {float64(3), "c"}},
	}, result)
}
