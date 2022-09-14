package extensionController

import (
	"fmt"
	"path"
	"testing"

	"github.com/exasol/extension-manager/apiErrors"
	"github.com/exasol/extension-manager/extensionAPI"
	"github.com/exasol/extension-manager/integrationTesting"

	"github.com/stretchr/testify/suite"
)

const (
	EXTENSION_SCHEMA = "test"
	EXTENSION_ID     = "testing-extension.js"
)

type ControllerITestSuite struct {
	suite.Suite
	exasol            *integrationTesting.DbTestSetup
	tempExtensionRepo string
}

func TestControllerITestSuite(t *testing.T) {
	suite.Run(t, new(ControllerITestSuite))
}

func (suite *ControllerITestSuite) SetupSuite() {
	suite.exasol = integrationTesting.StartDbSetup(&suite.Suite)
}

func (suite *ControllerITestSuite) TearDownSuite() {
	suite.exasol.StopDb()
}

func (suite *ControllerITestSuite) SetupTest() {
	suite.exasol.CreateConnection()
	suite.T().Cleanup(func() {
		suite.exasol.CloseConnection()
	})
	tempExtensionRepo := suite.T().TempDir()
	suite.tempExtensionRepo = tempExtensionRepo
}

func (suite *ControllerITestSuite) TestGetAllExtensions() {
	suite.writeDefaultExtension()
	suite.uploadBucketFsFile("123", "my-extension.1.2.3.jar") // create file with 3B size
	controller := Create(suite.tempExtensionRepo, EXTENSION_SCHEMA)
	extensions, err := controller.GetAllExtensions(mockContext(), suite.exasol.GetConnection())
	suite.NoError(err)
	suite.Equal(1, len(extensions))
	suite.Equal("MyDemoExtension", extensions[0].Name, "name")
	suite.Equal(EXTENSION_ID, extensions[0].Id, "id")
}

func (suite *ControllerITestSuite) writeDefaultExtension() {
	integrationTesting.CreateTestExtensionBuilder(suite.T()).
		WithBucketFsUpload(integrationTesting.BucketFsUploadParams{Name: "extension jar", BucketFsFilename: "my-extension.1.2.3.jar", FileSize: 3}).
		WithFindInstallationsFunc(`
		return metadata.allScripts.rows.map(row => {
			return {name: row.schema + "." + row.name, version: "0.1.0", instanceParameters: []}
		});`).
		Build().
		WriteToFile(path.Join(suite.tempExtensionRepo, EXTENSION_ID))
}

func (suite *ControllerITestSuite) TestGetAllExtensionsWithMissingJar() {
	integrationTesting.CreateTestExtensionBuilder(suite.T()).
		WithBucketFsUpload(integrationTesting.BucketFsUploadParams{Name: "extension jar", BucketFsFilename: "missing-jar.jar", FileSize: 3}).
		Build().
		WriteToFile(path.Join(suite.tempExtensionRepo, EXTENSION_ID))
	controller := Create(suite.tempExtensionRepo, EXTENSION_SCHEMA)
	extensions, err := controller.GetAllExtensions(mockContext(), suite.exasol.GetConnection())
	suite.NoError(err)
	suite.Empty(extensions)
}

func (suite *ControllerITestSuite) GetInstalledExtensions_failsWithGenericError() {
	const jarName = "my-failing-extension-1.2.3.jar"
	integrationTesting.CreateTestExtensionBuilder(suite.T()).
		WithBucketFsUpload(integrationTesting.BucketFsUploadParams{Name: "extension jar", BucketFsFilename: jarName, FileSize: 3}).
		WithFindInstallationsFunc("throw Error(`mock error from js`)").
		Build().
		WriteToFile(path.Join(suite.tempExtensionRepo, EXTENSION_ID))
	suite.uploadBucketFsFile("123", jarName) // create file with 3B size
	controller := Create(suite.tempExtensionRepo, EXTENSION_SCHEMA)
	extensions, err := controller.GetInstalledExtensions(mockContext(), suite.exasol.GetConnection())
	suite.ErrorContains(err, `failed to find installations: failed to find installations for extension "testing-extension.js": Error: mock error from js at`)
	suite.Nil(extensions)
}

func (suite *ControllerITestSuite) GetInstalledExtensions_failsWithApiError() {
	const jarName = "my-failing-extension-1.2.3.jar"
	integrationTesting.CreateTestExtensionBuilder(suite.T()).
		WithBucketFsUpload(integrationTesting.BucketFsUploadParams{Name: "extension jar", BucketFsFilename: jarName, FileSize: 3}).
		WithFindInstallationsFunc("throw new ApiError(400, `mock error from js`)").
		Build().
		WriteToFile(path.Join(suite.tempExtensionRepo, EXTENSION_ID))
	suite.uploadBucketFsFile("123", jarName)
	controller := Create(suite.tempExtensionRepo, EXTENSION_SCHEMA)
	extensions, err := controller.GetInstalledExtensions(mockContext(), suite.exasol.GetConnection())
	if apiError, ok := err.(*apiErrors.APIError); ok {
		suite.Equal("mock error from js", apiError.Message)
		suite.Equal(400, apiError.Status)
	} else {
		suite.Fail("wrong error type", "Expected APIError but got %T: %v", err, err)
	}
	suite.Nil(extensions)
}

func (suite *ControllerITestSuite) TestGetAllInstallations() {
	suite.writeDefaultExtension()
	fixture := integrationTesting.CreateLuaScriptFixture(suite.exasol.GetConnection())
	controller := Create(suite.tempExtensionRepo, fixture.GetSchemaName())
	fixture.Cleanup(suite.T())
	installations, err := controller.GetInstalledExtensions(mockContext(), suite.exasol.GetConnection())
	suite.NoError(err)
	suite.Equal(1, len(installations))
	suite.Equal(fixture.GetSchemaName()+".MY_SCRIPT", installations[0].Name)
}

func (suite *ControllerITestSuite) TestInstallFailsForUnknownExtensionId() {
	controller := Create(suite.tempExtensionRepo, EXTENSION_SCHEMA)
	err := controller.InstallExtension(mockContext(), suite.exasol.GetConnection(), "unknown-extension-id", "ver")
	suite.ErrorContains(err, "failed to load extension with id \"unknown-extension-id\": failed to load extension from file")
}

func (suite *ControllerITestSuite) TestInstallSucceeds() {
	suite.writeDefaultExtension()
	controller := Create(suite.tempExtensionRepo, EXTENSION_SCHEMA)
	err := controller.InstallExtension(mockContext(), suite.exasol.GetConnection(), EXTENSION_ID, "ver")
	suite.NoError(err)
}

func (suite *ControllerITestSuite) TestUninstallFailsForUnknownExtensionId() {
	controller := Create(suite.tempExtensionRepo, EXTENSION_SCHEMA)
	err := controller.UninstallExtension(mockContext(), suite.exasol.GetConnection(), "unknown-extension-id", "ver")
	suite.ErrorContains(err, "failed to load extension with id \"unknown-extension-id\": failed to load extension from file")
}

func (suite *ControllerITestSuite) TestUninstallSucceeds() {
	suite.writeDefaultExtension()
	controller := Create(suite.tempExtensionRepo, EXTENSION_SCHEMA)
	err := controller.UninstallExtension(mockContext(), suite.exasol.GetConnection(), EXTENSION_ID, "ver")
	suite.NoError(err)
}

func (suite *ControllerITestSuite) TestEnsureSchemaExistsCreatesSchemaIfItDoesNotExist() {
	suite.writeDefaultExtension()
	const schemaName = "my_testing_schema"
	suite.dropSchema(schemaName)
	defer suite.dropSchema(schemaName)
	controller := Create(suite.tempExtensionRepo, schemaName)
	suite.NotContains(suite.getAllSchemaNames(), schemaName)
	err := controller.InstallExtension(mockContext(), suite.exasol.GetConnection(), EXTENSION_ID, "ver")
	suite.NoError(err)
	suite.Contains(suite.getAllSchemaNames(), schemaName)
}

func (suite *ControllerITestSuite) TestEnsureSchemaDoesNotFailIfSchemaAlreadyExists() {
	suite.writeDefaultExtension()
	const schemaName = "my_testing_schema"
	defer suite.dropSchema(schemaName)
	controller := Create(suite.tempExtensionRepo, schemaName)
	suite.createSchema(schemaName)
	err := controller.InstallExtension(mockContext(), suite.exasol.GetConnection(), EXTENSION_ID, "ver")
	suite.NoError(err)
	suite.Contains(suite.getAllSchemaNames(), schemaName)
}

func (suite *ControllerITestSuite) TestAddInstance_wrongVersion() {
	integrationTesting.CreateTestExtensionBuilder(suite.T()).
		WithFindInstallationsFunc(integrationTesting.MockFindInstallationsFunction("test", "0.1.0", `[]`)).
		WithAddInstanceFunc("context.sqlClient.execute('select 1'); return {id: 'instId', name: `ext_${version}_${params.values[0].name}_${params.values[0].value}`};").
		Build().
		WriteToFile(path.Join(suite.tempExtensionRepo, EXTENSION_ID))
	controller := Create(suite.tempExtensionRepo, EXTENSION_SCHEMA)
	instance, err := controller.CreateInstance(mockContext(), suite.exasol.GetConnection(), EXTENSION_ID, "wrongVersion", []ParameterValue{})
	suite.EqualError(err, `failed to find installations: version "wrongVersion" not found for extension "testing-extension.js", available versions: ["0.1.0"]`)
	suite.Nil(instance)
}

func (suite *ControllerITestSuite) TestAddInstance_invalidParameters() {
	integrationTesting.CreateTestExtensionBuilder(suite.T()).
		WithFindInstallationsFunc(integrationTesting.MockFindInstallationsFunction("test", "0.1.0", `[{
		id: "p1",
		name: "My param",
		type: "string",
		required: true
	}]`)).WithAddInstanceFunc("context.sqlClient.execute('select 1'); return {id: 'instId', name: `ext_${version}_${params.values[0].name}_${params.values[0].value}`};").
		Build().
		WriteToFile(path.Join(suite.tempExtensionRepo, EXTENSION_ID))
	controller := Create(suite.tempExtensionRepo, EXTENSION_SCHEMA)
	instance, err := controller.CreateInstance(mockContext(), suite.exasol.GetConnection(), EXTENSION_ID, "0.1.0", []ParameterValue{})
	suite.EqualError(err, `invalid parameters: Failed to validate parameter 'My param': This is a required parameter.`)
	suite.Nil(instance)
}

func (suite *ControllerITestSuite) TestAddInstance_validParameters() {
	integrationTesting.CreateTestExtensionBuilder(suite.T()).
		WithFindInstallationsFunc(integrationTesting.MockFindInstallationsFunction("test", "0.1.0", `[{
		id: "p1",
		name: "My param",
		type: "string"
	}]`)).WithAddInstanceFunc("context.sqlClient.execute('select 1'); return {id: 'instId', name: `ext_${version}_${params.values[0].name}_${params.values[0].value}`};").
		Build().
		WriteToFile(path.Join(suite.tempExtensionRepo, EXTENSION_ID))
	controller := Create(suite.tempExtensionRepo, EXTENSION_SCHEMA)
	instance, err := controller.CreateInstance(mockContext(), suite.exasol.GetConnection(), EXTENSION_ID, "0.1.0", []ParameterValue{{Name: "p1", Value: "val"}})
	suite.NoError(err)
	suite.Equal(&extensionAPI.JsExtInstance{Id: "instId", Name: "ext_0.1.0_p1_val"}, instance)
}

func (suite *ControllerITestSuite) TestFindInstances() {
	integrationTesting.CreateTestExtensionBuilder(suite.T()).
		WithFindInstancesFunc("context.sqlClient.execute('select 1'); return [{id: 'instId', name: 'instName_ver'+version}]").
		Build().WriteToFile(path.Join(suite.tempExtensionRepo, EXTENSION_ID))
	controller := Create(suite.tempExtensionRepo, EXTENSION_SCHEMA)
	instances, err := controller.FindInstances(mockContext(), suite.exasol.GetConnection(), EXTENSION_ID, "0.1.0")
	suite.NoError(err)
	suite.Equal([]*extensionAPI.JsExtInstance{{Id: "instId", Name: "instName_ver0.1.0"}}, instances)
}

func (suite *ControllerITestSuite) TestFindInstances_useSqlQueryResult() {
	integrationTesting.CreateTestExtensionBuilder(suite.T()).
		WithFindInstancesFunc(`const result=context.sqlClient.query("select 1 as c1, 'a' as c2 from dual where 1=?", 1);
			const col1 = result.columns[0];
			const col2 = result.columns[1];
			const row1 = result.rows[0];` +
			"return [{id: 'instId', name: `${col1.name}: ${col1.typeName}/${typeof(row1[0])} = ${row1[0]}, ${col2.name}: ${col2.typeName}/${typeof(row1[1])} = ${row1[1]}`}]").
		Build().WriteToFile(path.Join(suite.tempExtensionRepo, EXTENSION_ID))
	controller := Create(suite.tempExtensionRepo, EXTENSION_SCHEMA)
	instances, err := controller.FindInstances(mockContext(), suite.exasol.GetConnection(), EXTENSION_ID, "0.1.0")
	suite.NoError(err)
	suite.Equal([]*extensionAPI.JsExtInstance{{Id: "instId", Name: "C1: DECIMAL/number = 1, C2: CHAR/string = a"}}, instances)
}

func (suite *ControllerITestSuite) TestDeleteInstances_failsWithInvalidQuery() {
	integrationTesting.CreateTestExtensionBuilder(suite.T()).
		WithDeleteInstanceFunc("context.sqlClient.execute('drop instance')").
		Build().WriteToFile(path.Join(suite.tempExtensionRepo, EXTENSION_ID))
	controller := Create(suite.tempExtensionRepo, EXTENSION_SCHEMA)
	err := controller.DeleteInstance(mockContext(), suite.exasol.GetConnection(), EXTENSION_ID, "instId")
	suite.ErrorContains(err, `failed to delete instance "instId" for extension "testing-extension.js": error executing statement "drop instance"`)
}

func (suite *ControllerITestSuite) TestDeleteInstances_succeeds() {
	integrationTesting.CreateTestExtensionBuilder(suite.T()).
		WithDeleteInstanceFunc("context.sqlClient.execute('select 1')").
		Build().WriteToFile(path.Join(suite.tempExtensionRepo, EXTENSION_ID))
	controller := Create(suite.tempExtensionRepo, EXTENSION_SCHEMA)
	err := controller.DeleteInstance(mockContext(), suite.exasol.GetConnection(), EXTENSION_ID, "instId")
	suite.NoError(err)
}

func (suite *ControllerITestSuite) createSchema(schemaName string) {
	_, err := suite.exasol.GetConnection().Exec(fmt.Sprintf(`CREATE SCHEMA "%s"`, schemaName))
	if err != nil {
		suite.FailNowf("failed to create schema %s: %v", schemaName, err.Error())
	}
}

func (suite *ControllerITestSuite) dropSchema(schemaName string) {
	_, err := suite.exasol.GetConnection().Exec(fmt.Sprintf(`DROP SCHEMA IF EXISTS "%s" CASCADE`, schemaName))
	if err != nil {
		suite.FailNowf("failed to drop schema %s: %v", schemaName, err.Error())
	}
}

func (suite *ControllerITestSuite) getAllSchemaNames() []string {
	rows, err := suite.exasol.GetConnection().Query("SELECT SCHEMA_NAME FROM EXA_ALL_SCHEMAS ORDER BY SCHEMA_NAME")
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

func (suite *ControllerITestSuite) uploadBucketFsFile(content, fileName string) {
	err := suite.exasol.Exasol.UploadStringContent(content, fileName)
	if err != nil {
		suite.FailNowf("upload failed", "failed to upload string content: %v", err)
	}
	suite.T().Cleanup(func() {
		suite.NoError(suite.exasol.Exasol.DeleteFile(fileName))
	})
}
