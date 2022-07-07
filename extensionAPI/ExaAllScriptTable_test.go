package extensionAPI

import (
	"testing"

	"github.com/exasol/extension-manager/integrationTesting"
	"github.com/stretchr/testify/suite"
)

type ExaAllScriptsTableSuite struct {
	integrationTesting.IntegrationTestSuite
}

func TestExaAllScriptsTableSuite(t *testing.T) {
	suite.Run(t, new(ExaAllScriptsTableSuite))
}

func (suite *ExaAllScriptsTableSuite) TestReadMetadataWithAllColumnsDefined() {
	luaScriptFixture := integrationTesting.CreateLuaScriptFixture(suite.Connection)
	defer luaScriptFixture.Close()
	result, err := ReadMetadataTables(suite.Connection)
	suite.NoError(err)
	suite.Assert().Equal(
		ExaAllScriptTable{Rows: []ExaAllScriptRow{{
			Schema:     "TEST",
			Name:       "MY_SCRIPT",
			Type:       "UDF",
			InputType:  "SET",
			ResultType: "RETURNS",
			Text:       "CREATE LUA SET SCRIPT \"MY_SCRIPT\" (\"a\" DOUBLE) RETURNS DOUBLE AS\nfunction run(ctx) return 1 end",
			Comment:    "my comment"}}}, result.AllScripts)
}

func (suite *ExaAllScriptsTableSuite) TestReadMetadataWithMissingValues() {
	luaScriptFixture := integrationTesting.CreateJavaAdapterScriptFixture(suite.Connection)
	defer luaScriptFixture.Close()
	result, err := ReadMetadataTables(suite.Connection)
	suite.NoError(err)
	suite.Assert().Equal(
		ExaAllScriptTable{Rows: []ExaAllScriptRow{{
			Schema:     "TEST",
			Name:       "VS_ADAPTER",
			Type:       "ADAPTER",
			InputType:  "",
			ResultType: "",
			Text:       "CREATE JAVA  ADAPTER SCRIPT \"VS_ADAPTER\" AS\n%scriptclass com.exasol.adapter.RequestDispatcher;\n%jar /buckets/bfsdefault/default/vs.jar;",
			Comment:    ""}}}, result.AllScripts)
}
