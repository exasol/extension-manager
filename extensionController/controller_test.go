package extensionController

import (
	"context"
	"database/sql"
	"fmt"
	"path"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/exasol/extension-manager/apiErrors"
	"github.com/exasol/extension-manager/extensionAPI"
	"github.com/exasol/extension-manager/integrationTesting"

	"github.com/stretchr/testify/suite"
)

type ControllerUTestSuite struct {
	suite.Suite
	tempExtensionRepo string
	controller        TransactionController
	db                *sql.DB
	dbMock            sqlmock.Sqlmock
	bucketFsMock      bucketFsMock
	metaDataMock      exaMetaDataReaderMock
}

func TestControllerUTestSuite(t *testing.T) {
	suite.Run(t, new(ControllerUTestSuite))
}

func (suite *ControllerUTestSuite) SetupTest() {
	tempExtensionRepo := suite.T().TempDir()
	suite.tempExtensionRepo = tempExtensionRepo
	suite.bucketFsMock = createBucketFsMock()
	suite.metaDataMock = createExaMetaDataReaderMock(EXTENSION_SCHEMA)
	ctrl := &controllerImpl{extensionFolder: suite.tempExtensionRepo, schema: EXTENSION_SCHEMA, metaDataReader: &suite.metaDataMock}
	suite.controller = &transactionControllerImpl{controller: ctrl, bucketFs: &suite.bucketFsMock}

	db, mock, err := sqlmock.New()
	if err != nil {
		suite.Failf("an error '%v' was not expected when opening a stub database connection", err.Error())
	}
	suite.db = db
	suite.dbMock = mock
	suite.dbMock.MatchExpectationsInOrder(true)
	suite.dbMock.ExpectBegin()
	suite.dbMock.ExpectExec(`CREATE SCHEMA IF NOT EXISTS "test"`).WillReturnResult(sqlmock.NewResult(0, 0))
}

func (suite *ControllerUTestSuite) TeardownTest() {
	suite.NoError(suite.dbMock.ExpectationsWereMet())
}

func (suite *ControllerUTestSuite) TestGetAllExtensions() {
	suite.writeDefaultExtension()
	suite.bucketFsMock.simulateFiles([]BfsFile{{Name: "my-extension.1.2.3.jar", Size: 3}})
	extensions, err := suite.controller.GetAllExtensions(mockContext(), suite.db)
	suite.NoError(err)
	suite.Equal([]*Extension{{Name: "MyDemoExtension", Id: "testing-extension.js", Description: "An extension for testing.",
		InstallableVersions: []string{"0.1.0"}}}, extensions)
}

func (suite *ControllerUTestSuite) TestGetAllExtensionsWithMissingJar() {
	suite.writeDefaultExtension()
	suite.bucketFsMock.simulateFiles([]BfsFile{})
	extensions, err := suite.controller.GetAllExtensions(mockContext(), suite.db)
	suite.NoError(err)
	suite.Empty(extensions)
}

func (suite *ControllerUTestSuite) TestGetAllInstallations_failsWithGenericError() {
	integrationTesting.CreateTestExtensionBuilder(suite.T()).
		WithFindInstallationsFunc("throw Error(`mock error from js`)").
		Build().
		WriteToFile(path.Join(suite.tempExtensionRepo, DEFAULT_EXTENSION_ID))
	suite.metaDataMock.simulateExaMetaData(extensionAPI.ExaMetadata{})
	extensions, err := suite.controller.GetInstalledExtensions(mockContext(), suite.db)
	suite.EqualError(err, `failed to find installations for extension "MyDemoExtension": mock error from js`)
	suite.Nil(extensions)
}

func (suite *ControllerUTestSuite) TestGetAllInstallations_failsWithApiError() {
	integrationTesting.CreateTestExtensionBuilder(suite.T()).
		WithFindInstallationsFunc("throw new ApiError(400, `mock error from js`)").
		Build().
		WriteToFile(path.Join(suite.tempExtensionRepo, DEFAULT_EXTENSION_ID))
	suite.metaDataMock.simulateExaMetaData(extensionAPI.ExaMetadata{})
	extensions, err := suite.controller.GetInstalledExtensions(mockContext(), suite.db)
	suite.assertApiError(err, 400, "failed to find installations for extension \"MyDemoExtension\": mock error from js")
	suite.Nil(extensions)
}

func (suite *ControllerUTestSuite) TestGetAllInstallations() {
	suite.writeDefaultExtension()
	suite.metaDataMock.simulateExaAllScripts([]extensionAPI.ExaAllScriptRow{{Schema: "schema", Name: "script"}})
	installations, err := suite.controller.GetInstalledExtensions(mockContext(), suite.db)
	suite.NoError(err)
	suite.Assert().Equal([]*extensionAPI.JsExtInstallation{{Name: "schema.script", Version: "0.1.0",
		InstanceParameters: []interface{}{map[string]interface{}{"id": "p1", "name": "param1", "type": "string"}}}}, installations)
}

func (suite *ControllerUTestSuite) TestInstallFailsForUnknownExtensionId() {
	err := suite.controller.InstallExtension(mockContext(), suite.db, "unknown-extension-id", "ver")
	suite.ErrorContains(err, "failed to load extension with id \"unknown-extension-id\": failed to load extension from file")
}

func (suite *ControllerUTestSuite) TestInstallSucceeds() {
	integrationTesting.CreateTestExtensionBuilder(suite.T()).
		WithInstallFunc("context.sqlClient.runQuery('install extension')").
		Build().
		WriteToFile(path.Join(suite.tempExtensionRepo, DEFAULT_EXTENSION_ID))
	suite.dbMock.ExpectExec("install extension").WillReturnResult(sqlmock.NewResult(0, 0))
	suite.dbMock.ExpectCommit()
	err := suite.controller.InstallExtension(mockContext(), suite.db, DEFAULT_EXTENSION_ID, "ver")
	suite.NoError(err)
}

func (suite *ControllerUTestSuite) TestInstall_QueryFails() {
	integrationTesting.CreateTestExtensionBuilder(suite.T()).
		WithInstallFunc("context.sqlClient.runQuery('install extension')").
		Build().
		WriteToFile(path.Join(suite.tempExtensionRepo, DEFAULT_EXTENSION_ID))
	suite.dbMock.ExpectExec("install extension").WillReturnError(fmt.Errorf("mock"))
	suite.dbMock.ExpectRollback()
	err := suite.controller.InstallExtension(mockContext(), suite.db, DEFAULT_EXTENSION_ID, "ver")
	suite.EqualError(err, "failed to install extension \"testing-extension.js\": error executing statement \"install extension\": mock")
}

func (suite *ControllerUTestSuite) TestInstall_FailsWithGenericError() {
	integrationTesting.CreateTestExtensionBuilder(suite.T()).
		WithInstallFunc("throw Error('mock error from js')").
		Build().
		WriteToFile(path.Join(suite.tempExtensionRepo, DEFAULT_EXTENSION_ID))
	suite.dbMock.ExpectCommit()
	err := suite.controller.InstallExtension(mockContext(), suite.db, DEFAULT_EXTENSION_ID, "ver")
	suite.EqualError(err, "mock error from js")
	suite.dbMock.ExpectRollback()
}

func (suite *ControllerUTestSuite) TestInstall_FailsWithApiError() {
	integrationTesting.CreateTestExtensionBuilder(suite.T()).
		WithInstallFunc("throw new ApiError(400, 'mock error from js')").
		Build().
		WriteToFile(path.Join(suite.tempExtensionRepo, DEFAULT_EXTENSION_ID))
	suite.dbMock.ExpectCommit()
	err := suite.controller.InstallExtension(mockContext(), suite.db, DEFAULT_EXTENSION_ID, "ver")
	suite.assertApiError(err, 400, "mock error from js")
	suite.dbMock.ExpectRollback()
}

func (suite *ControllerUTestSuite) TestAddInstance_wrongVersion() {
	integrationTesting.CreateTestExtensionBuilder(suite.T()).
		WithFindInstallationsFunc(integrationTesting.MockFindInstallationsFunction("test", "0.1.0", `[]`)).
		Build().
		WriteToFile(path.Join(suite.tempExtensionRepo, DEFAULT_EXTENSION_ID))
	suite.metaDataMock.simulateExaAllScripts([]extensionAPI.ExaAllScriptRow{{Schema: "schema", Name: "script"}})
	instanceName, err := suite.controller.CreateInstance(mockContext(), suite.db, DEFAULT_EXTENSION_ID, "wrongVersion", []ParameterValue{})
	suite.EqualError(err, `failed to find installations: version "wrongVersion" not found for extension "testing-extension.js", available versions: ["0.1.0"]`)
	suite.Equal("", instanceName)
	suite.dbMock.ExpectRollback()
}

func (suite *ControllerUTestSuite) TestAddInstance_invalidParameters() {
	integrationTesting.CreateTestExtensionBuilder(suite.T()).
		WithFindInstallationsFunc(integrationTesting.MockFindInstallationsFunction("test", "0.1.0", `[{
		id: "p1",
		name: "My param",
		type: "string",
		required: true
	}]`)).WithAddInstanceFunc("return {name: `ext_${version}_${params.values[0].name}_${params.values[0].value}`};").
		Build().
		WriteToFile(path.Join(suite.tempExtensionRepo, DEFAULT_EXTENSION_ID))
	suite.metaDataMock.simulateExaAllScripts([]extensionAPI.ExaAllScriptRow{})
	instanceName, err := suite.controller.CreateInstance(mockContext(), suite.db, DEFAULT_EXTENSION_ID, "0.1.0", []ParameterValue{})
	suite.EqualError(err, `invalid parameters: Failed to validate parameter 'My param': This is a required parameter.`)
	suite.Equal("", instanceName)
	suite.dbMock.ExpectRollback()
}

func (suite *ControllerUTestSuite) TestAddInstance_failsWithGenericError() {
	integrationTesting.CreateTestExtensionBuilder(suite.T()).
		WithFindInstallationsFunc(integrationTesting.MockFindInstallationsFunction("test", "0.1.0", `[]`)).
		WithAddInstanceFunc("throw Error('mock error from js')").
		Build().
		WriteToFile(path.Join(suite.tempExtensionRepo, DEFAULT_EXTENSION_ID))
	suite.metaDataMock.simulateExaAllScripts([]extensionAPI.ExaAllScriptRow{})
	instanceName, err := suite.controller.CreateInstance(mockContext(), suite.db, DEFAULT_EXTENSION_ID, "0.1.0", []ParameterValue{})
	suite.EqualError(err, `mock error from js`)
	suite.Equal("", instanceName)
	suite.dbMock.ExpectRollback()
}

func (suite *ControllerUTestSuite) TestAddInstance_failsWithApiError() {
	integrationTesting.CreateTestExtensionBuilder(suite.T()).
		WithFindInstallationsFunc(integrationTesting.MockFindInstallationsFunction("test", "0.1.0", `[]`)).
		WithAddInstanceFunc("throw new ApiError(400, 'mock error from js')").
		Build().
		WriteToFile(path.Join(suite.tempExtensionRepo, DEFAULT_EXTENSION_ID))
	suite.metaDataMock.simulateExaAllScripts([]extensionAPI.ExaAllScriptRow{})
	instanceName, err := suite.controller.CreateInstance(mockContext(), suite.db, DEFAULT_EXTENSION_ID, "0.1.0", []ParameterValue{})
	suite.assertApiError(err, 400, `mock error from js`)
	suite.Equal("", instanceName)
	suite.dbMock.ExpectRollback()
}

func (suite *ControllerUTestSuite) TestAddInstance_validParameters() {
	integrationTesting.CreateTestExtensionBuilder(suite.T()).
		WithFindInstallationsFunc(integrationTesting.MockFindInstallationsFunction("test", "0.1.0", `[{
		id: "p1",
		name: "My param",
		type: "string"
	}]`)).WithAddInstanceFunc("return {name: `ext_${version}_${params.values[0].name}_${params.values[0].value}`};").
		Build().
		WriteToFile(path.Join(suite.tempExtensionRepo, DEFAULT_EXTENSION_ID))
	suite.metaDataMock.simulateExaAllScripts([]extensionAPI.ExaAllScriptRow{})
	suite.dbMock.ExpectCommit()
	instanceName, err := suite.controller.CreateInstance(mockContext(), suite.db, DEFAULT_EXTENSION_ID, "0.1.0", []ParameterValue{{Name: "p1", Value: "val"}})
	suite.NoError(err)
	suite.Equal("ext_0.1.0_p1_val", instanceName)
}

func (suite *ControllerUTestSuite) writeDefaultExtension() {
	integrationTesting.CreateTestExtensionBuilder(suite.T()).
		WithBucketFsUpload(integrationTesting.BucketFsUploadParams{Name: "extension jar", BucketFsFilename: "my-extension.1.2.3.jar", FileSize: 3}).
		WithFindInstallationsFunc(`
		return metadata.allScripts.rows.map(row => {
			return {name: row.schema + "." + row.name, version: "0.1.0", instanceParameters: [{id:"p1", name:"param1", type:"string"}]}
		});`).
		WithInstallFunc("context.sqlClient.runQuery('install extension')").
		Build().
		WriteToFile(path.Join(suite.tempExtensionRepo, DEFAULT_EXTENSION_ID))
}

func mockContext() context.Context {
	return context.Background()
}

func (suite *ControllerUTestSuite) assertApiError(err error, expectedStatus int, expectedMessage string) {
	if apiError, ok := err.(*apiErrors.APIError); ok {
		suite.Equal(expectedMessage, apiError.Message)
		suite.Equal(400, apiError.Status)
	} else {
		suite.Fail("wrong error type", "Expected APIError but got %T: %v", err, err)
	}
}
