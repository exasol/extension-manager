package extensionApi

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type ExaAllScriptsTableSuite struct {
	IntegrationTestSuite
}

func TestExaAllScriptsTableSuite(t *testing.T) {
	suite.Run(t, new(ExaAllScriptsTableSuite))
}

func (suite *ExaAllScriptsTableSuite) TestReadScript() {
	connection := suite.Exasol.CreateConnection()
	defer func() { suite.NoError(connection.Close()) }()
	suite.createTestScript()
	defer suite.deleteTestScript()
	result, err := ReadExaAllScriptTable(connection)
	suite.NoError(err)
	suite.Assert().Equal(ExaAllScriptTable{Rows: []ExaAllScriptRow{{Name: "TEST.MY_SCRIPT", Text: "CREATE LUA SET SCRIPT \"MY_SCRIPT\" (\"a\" DOUBLE) RETURNS DOUBLE AS\nfunction run(ctx)\n  return 1\nend\n"}}}, *result)
}

func (suite *ExaAllScriptsTableSuite) createTestScript() {
	suite.ExecSQL("CREATE SCHEMA TEST")
	suite.ExecSQL(`CREATE LUA SET SCRIPT test.my_script (a DOUBLE)
    RETURNS DOUBLE AS
function run(ctx)
  return 1
end
/`)
}

func (suite *ExaAllScriptsTableSuite) deleteTestScript() {
	suite.ExecSQL("DROP SCRIPT test.my_script")
	suite.ExecSQL("DROP SCHEMA TEST")
}
