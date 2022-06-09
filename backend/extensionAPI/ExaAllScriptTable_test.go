package extensionAPI

import (
	"backend/integrationTesting"
	"github.com/stretchr/testify/suite"
	"testing"
)

type ExaAllScriptsTableSuite struct {
	integrationTesting.IntegrationTestSuite
}

func TestExaAllScriptsTableSuite(t *testing.T) {
	suite.Run(t, new(ExaAllScriptsTableSuite))
}

func (suite *ExaAllScriptsTableSuite) TestReadScript() {
	luaScriptFixture := integrationTesting.CreateLuaScriptFixture(suite.Connection)
	defer luaScriptFixture.Close()
	result, err := ReadExaAllScriptTable(suite.Connection)
	suite.NoError(err)
	suite.Assert().Equal(ExaAllScriptTable{Rows: []ExaAllScriptRow{{Name: "TEST.MY_SCRIPT", Text: "CREATE LUA SET SCRIPT \"MY_SCRIPT\" (\"a\" DOUBLE) RETURNS DOUBLE AS\nfunction run(ctx)\n  return 1\nend\n"}}}, *result)
}
