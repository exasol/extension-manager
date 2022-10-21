package extensionAPI

import (
	"fmt"
	"testing"

	"github.com/exasol/extension-manager/integrationTesting"
	"github.com/stretchr/testify/suite"
)

type ExaMetadataSuite struct {
	suite.Suite
	exasol *integrationTesting.DbTestSetup
}

func TestExaAllScriptsTableSuite(t *testing.T) {
	suite.Run(t, new(ExaMetadataSuite))
}

func (suite *ExaMetadataSuite) SetupSuite() {
	suite.exasol = integrationTesting.StartDbSetup(&suite.Suite)
}

func (suite *ExaMetadataSuite) TearDownSuite() {
	suite.exasol.StopDb()
}

func (suite *ExaMetadataSuite) BeforeTest(suiteName, testName string) {
	suite.exasol.CreateConnection()
	suite.T().Cleanup(func() {
		suite.exasol.CloseConnection()
	})
}

func (suite *ExaMetadataSuite) TestReadMetadataWithAllColumnsDefined() {
	fixture := integrationTesting.CreateLuaScriptFixture(suite.exasol.GetConnection())
	fixture.Cleanup(suite.T())
	result := suite.readMetaDataTables(fixture.GetSchemaName())
	suite.Equal(
		ExaScriptTable{Rows: []ExaScriptRow{{
			Schema:     "TEST",
			Name:       "MY_SCRIPT",
			Type:       "UDF",
			InputType:  "SET",
			ResultType: "RETURNS",
			Text:       "CREATE LUA SET SCRIPT \"MY_SCRIPT\" (\"a\" DOUBLE) RETURNS DOUBLE AS\nfunction run(ctx) return 1 end",
			Comment:    "my comment"}}}, result.AllScripts)
}

func (suite *ExaMetadataSuite) TestReadMetadataOfJavaAdapterScript() {
	fixture := integrationTesting.CreateJavaAdapterScriptFixture(suite.exasol.GetConnection())
	fixture.Cleanup(suite.T())
	result := suite.readMetaDataTables(fixture.GetSchemaName())
	suite.Equal(
		ExaScriptTable{Rows: []ExaScriptRow{{
			Schema:     "TEST",
			Name:       "VS_ADAPTER",
			Type:       "ADAPTER",
			InputType:  "",
			ResultType: "",
			Text:       "CREATE JAVA  ADAPTER SCRIPT \"VS_ADAPTER\" AS\n%scriptclass com.exasol.adapter.RequestDispatcher;\n%jar /buckets/bfsdefault/default/vs.jar;",
			Comment:    ""}}}, result.AllScripts)
}

func (suite *ExaMetadataSuite) TestReadMetadataOfJavaSetScript() {
	fixture := integrationTesting.CreateJavaSetScriptFixture(suite.exasol.GetConnection())
	fixture.Cleanup(suite.T())
	result := suite.readMetaDataTables(fixture.GetSchemaName())
	suite.Equal(
		ExaScriptTable{Rows: []ExaScriptRow{{
			Schema:     "TEST",
			Name:       "IMPORT_FROM_S3_DOCUMENT_FILES",
			Type:       "UDF",
			InputType:  "SET",
			ResultType: "EMITS",
			Text:       "CREATE JAVA SET SCRIPT \"IMPORT_FROM_S3_DOCUMENT_FILES\" (\"DATA_LOADER\" VARCHAR(2000000) UTF8, \"SCHEMA_MAPPING_REQUEST\" VARCHAR(2000000) UTF8, \"CONNECTION_NAME\" VARCHAR(500) UTF8) EMITS (...) AS\n%scriptclass com.exasol.adapter.document.UdfEntryPoint;\n%jar /buckets/bfsdefault/default/vs.jar;",
			Comment:    ""}}}, result.AllScripts)
}

func (suite *ExaMetadataSuite) TestReadMetadataScripts_NoResult() {
	result := suite.readMetaDataTables("dummy")
	suite.Equal(ExaScriptTable{Rows: []ExaScriptRow{}}, result.AllScripts)
}

func (suite *ExaMetadataSuite) TestReadMetadataVirtualSchemas_Empty() {
	result := suite.readMetaDataTables("dummy")
	suite.Equal(ExaVirtualSchemasTable{Rows: []ExaVirtualSchemaRow{}}, result.AllVirtualSchemas)
}

func (suite *ExaMetadataSuite) TestExtractSchemaAndName() {
	var tests = []struct {
		input          string
		expectedSchema string
		expectedName   string
		expectedError  bool
	}{
		{"", "", "", true},
		{"invalid", "", "", true},
		{"invalid_separator", "", "", true},
		{".name", "", "name", false},
		{"schema.", "schema", "", false},
		{"schema.name", "schema", "name", false},
		{"SCHEMA.NAME", "SCHEMA", "NAME", false},
	}
	for _, t := range tests {
		suite.Run(t.input, func() {
			schema, name, err := extractSchemaAndName(t.input)
			if t.expectedError {
				suite.EqualError(err, fmt.Sprintf("invalid format for adapter script: %q", t.input))
				suite.Equal("", schema)
				suite.Equal("", name)
			} else {
				suite.NoError(err)
				suite.Equal(t.expectedSchema, schema)
				suite.Equal(t.expectedName, name)
			}
		})
	}
}

func (suite *ExaMetadataSuite) readMetaDataTables(schemaName string) *ExaMetadata {
	tx, err := suite.exasol.GetConnection().Begin()
	suite.NoError(err)
	metaData, err := CreateExaMetaDataReader().ReadMetadataTables(tx, schemaName)
	suite.NoError(err)
	return metaData
}
