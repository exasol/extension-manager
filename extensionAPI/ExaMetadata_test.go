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
	fixture := integrationTesting.CreateLuaScriptFixture(suite.Connection)
	defer fixture.Close()
	result, err := ReadMetadataTables(suite.Connection, fixture.GetSchemaName())
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

func (suite *ExaAllScriptsTableSuite) TestReadMetadataOfJavaAdapterScript() {
	fixture := integrationTesting.CreateJavaAdapterScriptFixture(suite.Connection)
	defer fixture.Close()
	result, err := ReadMetadataTables(suite.Connection, fixture.GetSchemaName())
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

func (suite *ExaAllScriptsTableSuite) TestReadMetadataOfJavaSetScript() {
	fixture := integrationTesting.CreateJavaSetScriptFixture(suite.Connection)
	defer fixture.Close()
	result, err := ReadMetadataTables(suite.Connection, fixture.GetSchemaName())
	suite.NoError(err)
	suite.Assert().Equal(
		ExaAllScriptTable{Rows: []ExaAllScriptRow{{
			Schema:     "TEST",
			Name:       "IMPORT_FROM_S3_DOCUMENT_FILES",
			Type:       "UDF",
			InputType:  "SET",
			ResultType: "EMITS",
			Text:       "CREATE JAVA SET SCRIPT \"IMPORT_FROM_S3_DOCUMENT_FILES\" (\"DATA_LOADER\" VARCHAR(2000000) UTF8, \"SCHEMA_MAPPING_REQUEST\" VARCHAR(2000000) UTF8, \"CONNECTION_NAME\" VARCHAR(500) UTF8) EMITS (...) AS\n%scriptclass com.exasol.adapter.document.UdfEntryPoint;\n%jar /buckets/bfsdefault/default/vs.jar;",
			Comment:    ""}}}, result.AllScripts)
}
