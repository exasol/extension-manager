package extensionController

import (
	"fmt"
	"path"
	"testing"

	"github.com/exasol/extension-manager/pkg/apiErrors"
	"github.com/exasol/extension-manager/pkg/extensionAPI"
	"github.com/exasol/extension-manager/pkg/integrationTesting"

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

/* [itest -> dsn~extension-definition~1] */
func (suite *ControllerITestSuite) TestGetAllExtensions() {
	suite.writeDefaultExtension()
	suite.uploadBucketFsFile("123", "my-extension.1.2.3.jar") // create file with 3B size
	extensions, err := suite.createController().GetAllExtensions(mockContext(), suite.exasol.GetConnection())
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
			return {name: row.schema + "." + row.name, version: "0.1.0"}
		});`).
		Build().
		WriteToFile(path.Join(suite.tempExtensionRepo, EXTENSION_ID))
}

func (suite *ControllerITestSuite) TestGetAllExtensionsWithMissingJar() {
	integrationTesting.CreateTestExtensionBuilder(suite.T()).
		WithBucketFsUpload(integrationTesting.BucketFsUploadParams{Name: "extension jar", BucketFsFilename: "missing-jar.jar", FileSize: 3}).
		Build().
		WriteToFile(path.Join(suite.tempExtensionRepo, EXTENSION_ID))
	extensions, err := suite.createController().GetAllExtensions(mockContext(), suite.exasol.GetConnection())
	suite.NoError(err)
	suite.Empty(extensions)
}

func (suite *ControllerITestSuite) GetInstalledExtensionsFailsWithGenericError() {
	const jarName = "my-failing-extension-1.2.3.jar"
	integrationTesting.CreateTestExtensionBuilder(suite.T()).
		WithBucketFsUpload(integrationTesting.BucketFsUploadParams{Name: "extension jar", BucketFsFilename: jarName, FileSize: 3}).
		WithFindInstallationsFunc("throw Error(`mock error from js`)").
		Build().
		WriteToFile(path.Join(suite.tempExtensionRepo, EXTENSION_ID))
	suite.uploadBucketFsFile("123", jarName) // create file with 3B size
	extensions, err := suite.createController().GetInstalledExtensions(mockContext(), suite.exasol.GetConnection())
	suite.ErrorContains(err, `failed to find installations: failed to find installations for extension "testing-extension.js": Error: mock error from js at`)
	suite.Nil(extensions)
}

func (suite *ControllerITestSuite) GetInstalledExtensionsFailsWithApiError() {
	const jarName = "my-failing-extension-1.2.3.jar"
	integrationTesting.CreateTestExtensionBuilder(suite.T()).
		WithBucketFsUpload(integrationTesting.BucketFsUploadParams{Name: "extension jar", BucketFsFilename: jarName, FileSize: 3}).
		WithFindInstallationsFunc("throw new ApiError(400, `mock error from js`)").
		Build().
		WriteToFile(path.Join(suite.tempExtensionRepo, EXTENSION_ID))
	suite.uploadBucketFsFile("123", jarName)
	extensions, err := suite.createController().GetInstalledExtensions(mockContext(), suite.exasol.GetConnection())
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
	fixture.Cleanup(suite.T())
	installations, err := suite.createControllerWithSchema(fixture.GetSchemaName()).
		GetInstalledExtensions(mockContext(), suite.exasol.GetConnection())
	suite.NoError(err)
	suite.Equal(1, len(installations))
	suite.Equal(fixture.GetSchemaName()+".MY_SCRIPT", installations[0].Name)
}

// Install

func (suite *ControllerITestSuite) TestInstallFailsForUnknownExtensionId() {
	err := suite.createController().InstallExtension(mockContext(), suite.exasol.GetConnection(), "unknown-extension-id", "ver")
	suite.ErrorContains(err, "failed to load extension \"unknown-extension-id\"")
}

func (suite *ControllerITestSuite) TestInstallSucceeds() {
	suite.writeDefaultExtension()
	err := suite.createController().InstallExtension(mockContext(), suite.exasol.GetConnection(), EXTENSION_ID, "ver")
	suite.NoError(err)
}

// Uninstall

func (suite *ControllerITestSuite) TestUninstallFailsForUnknownExtensionId() {
	err := suite.createController().UninstallExtension(mockContext(), suite.exasol.GetConnection(), "unknown-extension-id", "ver")
	suite.ErrorContains(err, "failed to load extension \"unknown-extension-id\"")
}

func (suite *ControllerITestSuite) TestUninstallSucceeds() {
	suite.writeDefaultExtension()
	err := suite.createController().UninstallExtension(mockContext(), suite.exasol.GetConnection(), EXTENSION_ID, "ver")
	suite.NoError(err)
}

// Upgrade

func (suite *ControllerITestSuite) TestUpgradeFailsForUnknownExtensionId() {
	result, err := suite.createController().UpgradeExtension(mockContext(), suite.exasol.GetConnection(), "unknown-extension-id")
	suite.ErrorContains(err, "failed to load extension \"unknown-extension-id\"")
	suite.Nil(result)
}

/* [itest -> dsn~upgrade-extension~1] */
func (suite *ControllerITestSuite) TestUpgradeSucceeds() {
	suite.writeDefaultExtension()
	result, err := suite.createController().UpgradeExtension(mockContext(), suite.exasol.GetConnection(), EXTENSION_ID)
	suite.NoError(err)
	suite.Equal(&extensionAPI.JsUpgradeResult{PreviousVersion: "0.1.0", NewVersion: "0.2.0"}, result)
}

/* [itest -> const~use-reserved-schema~1] */
func (suite *ControllerITestSuite) TestEnsureSchemaExistsCreatesSchemaIfItDoesNotExist() {
	suite.writeDefaultExtension()
	const schemaName = "my_testing_schema"
	suite.dropSchema(schemaName)
	defer suite.dropSchema(schemaName)
	suite.NotContains(suite.getAllSchemaNames(), schemaName)
	err := suite.createControllerWithSchema(schemaName).InstallExtension(mockContext(), suite.exasol.GetConnection(), EXTENSION_ID, "ver")
	suite.NoError(err)
	suite.Contains(suite.getAllSchemaNames(), schemaName)
}

func (suite *ControllerITestSuite) TestEnsureSchemaDoesNotFailIfSchemaAlreadyExists() {
	suite.writeDefaultExtension()
	const schemaName = "my_testing_schema"
	defer suite.dropSchema(schemaName)
	suite.createSchema(schemaName)
	err := suite.createControllerWithSchema(schemaName).InstallExtension(mockContext(), suite.exasol.GetConnection(), EXTENSION_ID, "ver")
	suite.NoError(err)
	suite.Contains(suite.getAllSchemaNames(), schemaName)
}

/* [itest -> dsn~validate-parameters~1] */
/* [itest -> dsn~parameter-definitions~1] */
/* [itest -> dsn~extension-context-sql-client~1] */
func (suite *ControllerITestSuite) TestAddInstanceInvalidParameters() {
	integrationTesting.CreateTestExtensionBuilder(suite.T()).
		WithFindInstallationsFunc(integrationTesting.MockFindInstallationsFunction("test", "0.1.0")).
		WithAddInstanceFunc("context.sqlClient.execute('select 1'); return {id: 'instId', name: `ext_${version}_${params.values[0].name}_${params.values[0].value}`};").
		WithGetInstanceParameterDefinitionFunc(`return [{id: "param1", name: "My param", type: "string", required: true}]`).
		Build().
		WriteToFile(path.Join(suite.tempExtensionRepo, EXTENSION_ID))
	instance, err := suite.createController().CreateInstance(mockContext(), suite.exasol.GetConnection(), EXTENSION_ID, "0.1.0", []ParameterValue{})
	suite.EqualError(err, `invalid parameters: Failed to validate parameter 'My param': This is a required parameter.`)
	suite.Nil(instance)
}

func (suite *ControllerITestSuite) TestAddInstanceValidParameters() {
	integrationTesting.CreateTestExtensionBuilder(suite.T()).
		WithFindInstallationsFunc(integrationTesting.MockFindInstallationsFunction("test", "0.1.0")).
		WithAddInstanceFunc("context.sqlClient.execute('select 1'); return {id: 'instId', name: `ext_${version}_${params.values[0].name}_${params.values[0].value}`};").
		Build().
		WriteToFile(path.Join(suite.tempExtensionRepo, EXTENSION_ID))
	instance, err := suite.createController().CreateInstance(mockContext(), suite.exasol.GetConnection(), EXTENSION_ID, "0.1.0", []ParameterValue{{Name: "p1", Value: "val"}})
	suite.NoError(err)
	suite.Equal(&extensionAPI.JsExtInstance{Id: "instId", Name: "ext_0.1.0_p1_val"}, instance)
}

func (suite *ControllerITestSuite) TestFindInstances() {
	integrationTesting.CreateTestExtensionBuilder(suite.T()).
		WithFindInstancesFunc("context.sqlClient.execute('select 1'); return [{id: 'instId', name: 'instName_ver'+version}]").
		Build().WriteToFile(path.Join(suite.tempExtensionRepo, EXTENSION_ID))
	instances, err := suite.createController().FindInstances(mockContext(), suite.exasol.GetConnection(), EXTENSION_ID, "0.1.0")
	suite.NoError(err)
	suite.Equal([]*extensionAPI.JsExtInstance{{Id: "instId", Name: "instName_ver0.1.0"}}, instances)
}

/* [itest -> dsn~extension-context-sql-client~1] */
func (suite *ControllerITestSuite) TestFindInstancesUseSqlQueryResult() {
	integrationTesting.CreateTestExtensionBuilder(suite.T()).
		WithFindInstancesFunc(`const result=context.sqlClient.query("select 1 as c1, 'a' as c2 from dual where 1=?", 1);
			const col1 = result.columns[0];
			const col2 = result.columns[1];
			const row1 = result.rows[0];` +
			"return [{id: 'instId', name: `${col1.name}: ${col1.typeName}/${typeof(row1[0])} = ${row1[0]}, ${col2.name}: ${col2.typeName}/${typeof(row1[1])} = ${row1[1]}`}]").
		Build().WriteToFile(path.Join(suite.tempExtensionRepo, EXTENSION_ID))
	instances, err := suite.createController().FindInstances(mockContext(), suite.exasol.GetConnection(), EXTENSION_ID, "0.1.0")
	suite.NoError(err)
	suite.Equal([]*extensionAPI.JsExtInstance{{Id: "instId", Name: "C1: DECIMAL/number = 1, C2: CHAR/string = a"}}, instances)
}

func (suite *ControllerITestSuite) TestDeleteInstancesFailsWithInvalidQuery() {
	integrationTesting.CreateTestExtensionBuilder(suite.T()).
		WithDeleteInstanceFunc("context.sqlClient.execute('drop instance')").
		Build().WriteToFile(path.Join(suite.tempExtensionRepo, EXTENSION_ID))
	err := suite.createController().DeleteInstance(mockContext(), suite.exasol.GetConnection(), EXTENSION_ID, "extVersion", "instId")
	suite.ErrorContains(err, `failed to delete instance "instId" for extension "testing-extension.js": error executing statement "drop instance"`)
}

func (suite *ControllerITestSuite) TestDeleteInstancesSucceeds() {
	integrationTesting.CreateTestExtensionBuilder(suite.T()).
		WithDeleteInstanceFunc("context.sqlClient.execute('select 1')").
		Build().WriteToFile(path.Join(suite.tempExtensionRepo, EXTENSION_ID))
	err := suite.createController().DeleteInstance(mockContext(), suite.exasol.GetConnection(), EXTENSION_ID, "extVersion", "instId")
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

func (suite *ControllerITestSuite) createControllerWithSchema(schema string) TransactionController {
	return Create(suite.tempExtensionRepo, schema)
}

func (suite *ControllerITestSuite) createController() TransactionController {
	return suite.createControllerWithSchema(EXTENSION_SCHEMA)
}
