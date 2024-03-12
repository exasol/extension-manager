package exaMetadata_test

import (
	"testing"

	"github.com/exasol/extension-manager/pkg/extensionAPI/exaMetadata"
	"github.com/exasol/extension-manager/pkg/integrationTesting"
	"github.com/stretchr/testify/suite"
)

type ExaMetadataITestSuite struct {
	suite.Suite
	exasol *integrationTesting.DbTestSetup
}

func TestExaMetadataITestSuite(t *testing.T) {
	suite.Run(t, new(ExaMetadataITestSuite))
}

func (suite *ExaMetadataITestSuite) SetupSuite() {
	suite.exasol = integrationTesting.StartDbSetup(&suite.Suite)
}

func (suite *ExaMetadataITestSuite) TearDownSuite() {
	suite.exasol.StopDb()
}

func (suite *ExaMetadataITestSuite) BeforeTest(suiteName, testName string) {
	suite.exasol.CreateConnection()
	suite.T().Cleanup(func() {
		suite.exasol.CloseConnection()
	})
}

/* [utest -> dsn~extension-components~1]. */
func (suite *ExaMetadataITestSuite) TestReadMetadataWithAllColumnsDefined() {
	fixture := integrationTesting.CreateLuaScriptFixture(suite.exasol.GetConnection())
	fixture.Cleanup(suite.T())
	result := suite.readMetaDataTables(fixture.GetSchemaName())
	suite.Equal(
		exaMetadata.ExaScriptTable{Rows: []exaMetadata.ExaScriptRow{{
			Schema:     "TEST",
			Name:       "MY_SCRIPT",
			Type:       "UDF",
			InputType:  "SET",
			ResultType: "RETURNS",
			Text:       "CREATE LUA SET SCRIPT \"MY_SCRIPT\" (\"a\" DOUBLE) RETURNS DOUBLE AS\nfunction run(ctx) return 1 end",
			Comment:    "my comment"}}}, result.AllScripts)
}

func (suite *ExaMetadataITestSuite) TestReadMetadataOfJavaAdapterScript() {
	fixture := integrationTesting.CreateJavaAdapterScriptFixture(suite.exasol.GetConnection())
	fixture.Cleanup(suite.T())
	result := suite.readMetaDataTables(fixture.GetSchemaName())
	suite.Equal(
		exaMetadata.ExaScriptTable{Rows: []exaMetadata.ExaScriptRow{{
			Schema:     "TEST",
			Name:       "VS_ADAPTER",
			Type:       "ADAPTER",
			InputType:  "",
			ResultType: "",
			Text:       "CREATE JAVA  ADAPTER SCRIPT \"VS_ADAPTER\" AS\n%scriptclass com.exasol.adapter.RequestDispatcher;\n%jar /buckets/bfsdefault/default/vs.jar;",
			Comment:    ""}}}, result.AllScripts)
}

func (suite *ExaMetadataITestSuite) TestReadMetadataOfJavaSetScript() {
	fixture := integrationTesting.CreateJavaSetScriptFixture(suite.exasol.GetConnection())
	fixture.Cleanup(suite.T())
	result := suite.readMetaDataTables(fixture.GetSchemaName())
	suite.Equal(
		exaMetadata.ExaScriptTable{Rows: []exaMetadata.ExaScriptRow{{
			Schema:     "TEST",
			Name:       "IMPORT_FROM_S3_DOCUMENT_FILES",
			Type:       "UDF",
			InputType:  "SET",
			ResultType: "EMITS",
			Text:       "CREATE JAVA SET SCRIPT \"IMPORT_FROM_S3_DOCUMENT_FILES\" (\"DATA_LOADER\" VARCHAR(2000000) UTF8, \"SCHEMA_MAPPING_REQUEST\" VARCHAR(2000000) UTF8, \"CONNECTION_NAME\" VARCHAR(500) UTF8) EMITS (...) AS\n%scriptclass com.exasol.adapter.document.UdfEntryPoint;\n%jar /buckets/bfsdefault/default/vs.jar;",
			Comment:    ""}}}, result.AllScripts)
}

func (suite *ExaMetadataITestSuite) TestReadMetadataScriptsNoResult() {
	result := suite.readMetaDataTables("dummy")
	suite.Equal(exaMetadata.ExaScriptTable{Rows: []exaMetadata.ExaScriptRow{}}, result.AllScripts)
}

func (suite *ExaMetadataITestSuite) TestReadMetadataVirtualSchemasEmpty() {
	result := suite.readMetaDataTables("dummy")
	suite.Equal(exaMetadata.ExaVirtualSchemasTable{Rows: []exaMetadata.ExaVirtualSchemaRow{}}, result.AllVirtualSchemas)
}

/* [itest -> dsn~extension-context-metadata~1]. */
func (suite *ExaMetadataITestSuite) TestGetScriptByName() {
	fixture := integrationTesting.CreateJavaAdapterScriptFixture(suite.exasol.GetConnection())
	fixture.Cleanup(suite.T())
	result, err := suite.getScriptByName(fixture.GetSchemaName(), "VS_ADAPTER")
	suite.Require().NoError(err)
	suite.Equal(
		&exaMetadata.ExaScriptRow{
			Schema:     "TEST",
			Name:       "VS_ADAPTER",
			Type:       "ADAPTER",
			InputType:  "",
			ResultType: "",
			Text:       "CREATE JAVA  ADAPTER SCRIPT \"VS_ADAPTER\" AS\n%scriptclass com.exasol.adapter.RequestDispatcher;\n%jar /buckets/bfsdefault/default/vs.jar;",
			Comment:    ""}, result)
}

func (suite *ExaMetadataITestSuite) TestGetScriptByNameNoResult() {
	result, err := suite.getScriptByName("schema", "script")
	suite.Require().NoError(err)
	suite.Nil(result)
}

func (suite *ExaMetadataITestSuite) readMetaDataTables(schemaName string) *exaMetadata.ExaMetadata {
	tx, err := suite.exasol.GetConnection().Begin()
	suite.Require().NoError(err)
	metaData, err := exaMetadata.CreateExaMetaDataReader().ReadMetadataTables(tx, schemaName)
	suite.Require().NoError(err)
	return metaData
}

func (suite *ExaMetadataITestSuite) getScriptByName(schemaName, scriptName string) (*exaMetadata.ExaScriptRow, error) {
	tx, err := suite.exasol.GetConnection().Begin()
	suite.Require().NoError(err)
	return exaMetadata.CreateExaMetaDataReader().GetScriptByName(tx, schemaName, scriptName)
}
