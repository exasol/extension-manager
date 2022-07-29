package extensionController

import (
	"fmt"
	"os"
	"path"
	"testing"

	"github.com/exasol/extension-manager/integrationTesting"

	"github.com/stretchr/testify/suite"
)

const (
	EXTENSION_SCHEMA     = "test"
	DEFAULT_EXTENSION_ID = "testing-extension.js"
)

type ExtensionControllerSuite struct {
	integrationTesting.IntegrationTestSuite
	tempExtensionRepo string
}

func TestExtensionControllerSuite(t *testing.T) {
	suite.Run(t, new(ExtensionControllerSuite))
}

func (suite *ExtensionControllerSuite) SetupSuite() {
	suite.IntegrationTestSuite.SetupSuite()
}

func (suite *ExtensionControllerSuite) TearDownSuite() {
	suite.IntegrationTestSuite.TearDownSuite()
}

func (suite *ExtensionControllerSuite) SetupTest() {
	tempExtensionRepo, err := os.MkdirTemp(os.TempDir(), "ExtensionControllerSuite")
	if err != nil {
		panic(err)
	}
	suite.tempExtensionRepo = tempExtensionRepo
}

func (suite *ExtensionControllerSuite) AfterTest(suiteName, testName string) {
	err := os.RemoveAll(suite.tempExtensionRepo)
	if err != nil {
		panic(err)
	}
}

func (suite *ExtensionControllerSuite) TestGetAllExtensions() {
	suite.writeDefaultExtension()
	suite.NoError(suite.Exasol.UploadStringContent("123", "my-extension.1.2.3.jar")) // create file with 3B size
	defer func() { suite.NoError(suite.Exasol.DeleteFile("my-extension.1.2.3.jar")) }()
	controller := Create(suite.tempExtensionRepo, EXTENSION_SCHEMA)
	extensions, err := controller.GetAllExtensions(suite.Connection)
	suite.NoError(err)
	suite.Assert().Equal(1, len(extensions))
	suite.Assert().Equal("MyDemoExtension", extensions[0].Name, "name")
	suite.Assert().Equal(DEFAULT_EXTENSION_ID, extensions[0].Id, "id")
}

func (suite *ExtensionControllerSuite) writeDefaultExtension() {
	integrationTesting.CreateTestExtensionBuilder().
		WithBucketFsUpload(integrationTesting.BucketFsUploadParams{Name: "extension jar", BucketFsFilename: "my-extension.1.2.3.jar", FileSize: 3}).
		WithFindInstallationsFunc(`
		return metadata.allScripts.rows.map(row => {
			return {name: row.schema + "." + row.name, version: "0.1.0", instanceParameters: []}
		});`).
		Build().
		WriteToFile(path.Join(suite.tempExtensionRepo, DEFAULT_EXTENSION_ID))
}

func (suite *ExtensionControllerSuite) TestGetAllExtensionsWithMissingJar() {
	integrationTesting.CreateTestExtensionBuilder().
		WithBucketFsUpload(integrationTesting.BucketFsUploadParams{Name: "extension jar", BucketFsFilename: "missing-jar.jar", FileSize: 3}).
		Build().
		WriteToFile(path.Join(suite.tempExtensionRepo, DEFAULT_EXTENSION_ID))
	controller := Create(suite.tempExtensionRepo, EXTENSION_SCHEMA)
	dbConnectionWithNoAutocommit, err := suite.Exasol.CreateConnectionWithConfig(false)
	suite.NoError(err)
	defer func() { suite.NoError(dbConnectionWithNoAutocommit.Close()) }()
	extensions, err := controller.GetAllExtensions(dbConnectionWithNoAutocommit)
	suite.NoError(err)
	suite.Assert().Empty(extensions)
}

func (suite *ExtensionControllerSuite) TestGetAllExtensionsThrowingJSError() {
	const jarName = "my-failing-extension-1.2.3.jar"
	integrationTesting.CreateTestExtensionBuilder().
		WithBucketFsUpload(integrationTesting.BucketFsUploadParams{Name: "extension jar", BucketFsFilename: jarName, FileSize: 3}).
		WithFindInstallationsFunc("throw Error(`mock error from js`)").
		Build().
		WriteToFile(path.Join(suite.tempExtensionRepo, DEFAULT_EXTENSION_ID))
	suite.NoError(suite.Exasol.UploadStringContent("123", jarName)) // create file with 3B size
	defer func() { suite.NoError(suite.Exasol.DeleteFile(jarName)) }()
	controller := Create(suite.tempExtensionRepo, EXTENSION_SCHEMA)
	extensions, err := controller.GetAllInstallations(suite.Connection)
	suite.ErrorContains(err, `failed to find installations: failed to find installations for extension "testing-extension.js": Error: mock error from js at`)
	suite.Nil(extensions)
}

func (suite *ExtensionControllerSuite) TestGetAllInstallations() {
	suite.writeDefaultExtension()
	fixture := integrationTesting.CreateLuaScriptFixture(suite.Connection)
	controller := Create(suite.tempExtensionRepo, fixture.GetSchemaName())
	defer fixture.Close()
	installations, err := controller.GetAllInstallations(suite.Connection)
	suite.NoError(err)
	suite.Assert().Equal(1, len(installations))
	suite.Assert().Equal(fixture.GetSchemaName()+".MY_SCRIPT", installations[0].Name)
}

func (suite *ExtensionControllerSuite) TestInstallFailsForUnknownExtensionId() {
	controller := Create(suite.tempExtensionRepo, EXTENSION_SCHEMA)
	err := controller.InstallExtension(suite.Connection, "unknown-extension-id", "ver")
	suite.ErrorContains(err, "failed to load extension with id \"unknown-extension-id\": failed to load extension from file")
}

func (suite *ExtensionControllerSuite) TestInstallSucceeds() {
	suite.writeDefaultExtension()
	controller := Create(suite.tempExtensionRepo, EXTENSION_SCHEMA)
	err := controller.InstallExtension(suite.Connection, DEFAULT_EXTENSION_ID, "ver")
	suite.NoError(err)
}

func (suite *ExtensionControllerSuite) TestEnsureSchemaExistsCreatesSchemaIfItDoesNotExist() {
	suite.writeDefaultExtension()
	const schemaName = "my_testing_schema"
	defer suite.dropSchema(schemaName)
	controller := Create(suite.tempExtensionRepo, schemaName)
	suite.Assert().NotContains(suite.getAllSchemaNames(), schemaName)
	err := controller.InstallExtension(suite.Connection, DEFAULT_EXTENSION_ID, "ver")
	suite.NoError(err)
	suite.Assert().Contains(suite.getAllSchemaNames(), schemaName)
}

func (suite *ExtensionControllerSuite) TestEnsureSchemaDoesNotFailIfSchemaAlreadyExists() {
	suite.writeDefaultExtension()
	const schemaName = "my_testing_schema"
	defer suite.dropSchema(schemaName)
	controller := Create(suite.tempExtensionRepo, schemaName)
	suite.createSchema(schemaName)
	err := controller.InstallExtension(suite.Connection, DEFAULT_EXTENSION_ID, "ver")
	suite.NoError(err)
	suite.Assert().Contains(suite.getAllSchemaNames(), schemaName)
}

func (suite *ExtensionControllerSuite) TestAddInstance_wrongVersion() {
	integrationTesting.CreateTestExtensionBuilder().
		WithFindInstallationsFunc(integrationTesting.MockFindInstallationsFunction("test", "0.1.0", `[]`)).
		WithAddInstanceFunc("context.sqlClient.runQuery('select 1'); return {name: `ext_${version}_${params.values[0].name}_${params.values[0].value}`};").
		Build().
		WriteToFile(path.Join(suite.tempExtensionRepo, DEFAULT_EXTENSION_ID))
	controller := Create(suite.tempExtensionRepo, EXTENSION_SCHEMA)
	instanceName, err := controller.CreateInstance(suite.Connection, DEFAULT_EXTENSION_ID, "wrongVersion", []ParameterValue{})
	suite.EqualError(err, `failed to find installations: version "wrongVersion" not found for extension "testing-extension.js", available versions: ["0.1.0"]`)
	suite.Equal("", instanceName)
}

func (suite *ExtensionControllerSuite) TestAddInstance_invalidParameters() {
	integrationTesting.CreateTestExtensionBuilder().
		WithFindInstallationsFunc(integrationTesting.MockFindInstallationsFunction("test", "0.1.0", `[{
		id: "p1",
		name: "My param",
		type: "string",
		required: true
	}]`)).WithAddInstanceFunc("context.sqlClient.runQuery('select 1'); return {name: `ext_${version}_${params.values[0].name}_${params.values[0].value}`};").
		Build().
		WriteToFile(path.Join(suite.tempExtensionRepo, DEFAULT_EXTENSION_ID))
	controller := Create(suite.tempExtensionRepo, EXTENSION_SCHEMA)
	instanceName, err := controller.CreateInstance(suite.Connection, DEFAULT_EXTENSION_ID, "0.1.0", []ParameterValue{{Name: "p2", Value: "val"}})
	suite.EqualError(err, `invalid parameters: Failed to validate parameter "My param": This is a required field.`)
	suite.Equal("", instanceName)
}

func (suite *ExtensionControllerSuite) TestAddInstance_validParameters() {
	integrationTesting.CreateTestExtensionBuilder().
		WithFindInstallationsFunc(integrationTesting.MockFindInstallationsFunction("test", "0.1.0", `[{
		id: "p1",
		name: "My param",
		type: "string"
	}]`)).WithAddInstanceFunc("context.sqlClient.runQuery('select 1'); return {name: `ext_${version}_${params.values[0].name}_${params.values[0].value}`};").
		Build().
		WriteToFile(path.Join(suite.tempExtensionRepo, DEFAULT_EXTENSION_ID))
	controller := Create(suite.tempExtensionRepo, EXTENSION_SCHEMA)
	instanceName, err := controller.CreateInstance(suite.Connection, DEFAULT_EXTENSION_ID, "0.1.0", []ParameterValue{{Name: "p1", Value: "val"}})
	suite.NoError(err)
	suite.Equal("ext_0.1.0_p1_val", instanceName)
}

func (suite *ExtensionControllerSuite) createSchema(schemaName string) {
	_, err := suite.Connection.Exec(fmt.Sprintf(`CREATE SCHEMA "%s"`, schemaName))
	suite.NoError(err)
}

func (suite *ExtensionControllerSuite) dropSchema(schemaName string) {
	_, err := suite.Connection.Exec(fmt.Sprintf(`DROP SCHEMA IF EXISTS "%s"`, schemaName))
	suite.NoError(err)
}

func (suite *ExtensionControllerSuite) getAllSchemaNames() []string {
	rows, err := suite.Connection.Query("SELECT SCHEMA_NAME FROM EXA_ALL_SCHEMAS ORDER BY SCHEMA_NAME")
	suite.NoError(err)
	defer rows.Close()
	var schemaNames []string
	for rows.Next() {
		var schemaName string
		suite.NoError(rows.Scan(&schemaName))
		schemaNames = append(schemaNames, schemaName)
	}
	return schemaNames
}
