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

func (suite *ControllerUTestSuite) BeforeTest(suiteName, testName string) {
	tempExtensionRepo := suite.T().TempDir()
	suite.tempExtensionRepo = tempExtensionRepo
	suite.createController()
	suite.initDbMock()
}

func (suite *ControllerUTestSuite) createController() {
	suite.bucketFsMock = createBucketFsMock()
	suite.metaDataMock = createExaMetaDataReaderMock(EXTENSION_SCHEMA)
	ctrl := &controllerImpl{extensionFolder: suite.tempExtensionRepo, schema: EXTENSION_SCHEMA, metaDataReader: &suite.metaDataMock}
	suite.controller = &transactionControllerImpl{controller: ctrl, bucketFs: &suite.bucketFsMock}
}

func (suite *ControllerUTestSuite) initDbMock() {
	db, mock, err := sqlmock.New()
	if err != nil {
		suite.Failf("an error '%v' was not expected when opening a stub database connection", err.Error())
	}
	suite.db = db
	suite.dbMock = mock
	suite.dbMock.MatchExpectationsInOrder(true)
}

func (suite *ControllerUTestSuite) AfterTest(suiteName, testName string) {
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

type errorTest struct {
	testName        string
	throwCommand    string
	expectedStatus  int
	expectedMessage string
}

var errorTests = []errorTest{
	{testName: "generic", throwCommand: "throw Error(`mock error from js`)", expectedStatus: -1},
	{testName: "internal server error", throwCommand: "throw new InternalServerError(`mock error from js`)", expectedStatus: -1},
	{testName: "bad request", throwCommand: "throw new BadRequestError(`mock error from js`)", expectedStatus: 400},
	{testName: "null pointer", throwCommand: `(<any>{}).a.b; throw Error("mock")`, expectedStatus: -1, expectedMessage: "TypeError: Cannot read property 'b' of undefined"},
}

func (suite *ControllerUTestSuite) TestGetAllInstallationsFails() {
	for _, t := range errorTests {
		suite.Run(t.testName, func() {
			integrationTesting.CreateTestExtensionBuilder(suite.T()).
				WithFindInstallationsFunc(t.throwCommand).
				Build().
				WriteToFile(path.Join(suite.tempExtensionRepo, EXTENSION_ID))
			suite.metaDataMock.simulateExaMetaData(extensionAPI.ExaMetadata{})
			suite.initDbMock()
			suite.dbMock.ExpectBegin()
			suite.dbMock.ExpectRollback()
			extensions, err := suite.controller.GetInstalledExtensions(mockContext(), suite.db)
			suite.assertError(t, err)
			suite.Nil(extensions)
		})
	}
}

func (suite *ControllerUTestSuite) TestFindInstancesFails() {
	for _, t := range errorTests {
		suite.Run(t.testName, func() {
			integrationTesting.CreateTestExtensionBuilder(suite.T()).
				WithFindInstancesFunc(t.throwCommand).
				Build().
				WriteToFile(path.Join(suite.tempExtensionRepo, EXTENSION_ID))
			suite.metaDataMock.simulateExaMetaData(extensionAPI.ExaMetadata{})
			suite.initDbMock()
			suite.dbMock.ExpectBegin()
			suite.dbMock.ExpectRollback()
			extensions, err := suite.controller.FindInstances(mockContext(), suite.db, EXTENSION_ID, "ver")
			suite.assertError(t, err)
			suite.Nil(extensions)
		})
	}
}

func (suite *ControllerUTestSuite) TestDeleteInstancesFails() {
	for _, t := range errorTests {
		suite.Run(t.testName, func() {
			integrationTesting.CreateTestExtensionBuilder(suite.T()).
				WithDeleteInstanceFunc(t.throwCommand).
				Build().
				WriteToFile(path.Join(suite.tempExtensionRepo, EXTENSION_ID))
			suite.initDbMock()
			suite.dbMock.ExpectBegin()
			suite.dbMock.ExpectRollback()
			err := suite.controller.DeleteInstance(mockContext(), suite.db, EXTENSION_ID, "instId")
			suite.assertError(t, err)
		})
	}
}

func (suite *ControllerUTestSuite) assertError(t errorTest, actualError error) {
	expectedErrorMessage := "mock error from js"
	if t.expectedMessage != "" {
		expectedErrorMessage = t.expectedMessage
	}
	if t.expectedStatus > 0 {
		suite.assertApiError(actualError, t.expectedStatus, expectedErrorMessage)
	} else {
		suite.assertNonApiError(actualError, expectedErrorMessage)
	}
}

func (suite *ControllerUTestSuite) TestGetAllInstallations() {
	suite.writeDefaultExtension()
	suite.metaDataMock.simulateExaAllScripts([]extensionAPI.ExaScriptRow{{Schema: "schema", Name: "script"}})
	suite.dbMock.ExpectBegin()
	suite.dbMock.ExpectRollback()
	installations, err := suite.controller.GetInstalledExtensions(mockContext(), suite.db)
	suite.NoError(err)
	suite.Equal([]*extensionAPI.JsExtInstallation{{Name: "schema.script", Version: "0.1.0",
		InstanceParameters: []interface{}{map[string]interface{}{"id": "p1", "name": "param1", "type": "string"}}}}, installations)
}

func (suite *ControllerUTestSuite) TestInstallFailsForUnknownExtensionId() {
	suite.dbMock.ExpectBegin()
	suite.dbMock.ExpectRollback()
	err := suite.controller.InstallExtension(mockContext(), suite.db, "unknown-extension-id", "ver")
	suite.ErrorContains(err, "failed to load extension with id \"unknown-extension-id\": failed to load extension from file")
}

func (suite *ControllerUTestSuite) TestInstallSucceeds() {
	integrationTesting.CreateTestExtensionBuilder(suite.T()).
		WithInstallFunc("context.sqlClient.execute('install extension')").
		Build().
		WriteToFile(path.Join(suite.tempExtensionRepo, EXTENSION_ID))
	suite.dbMock.ExpectBegin()
	suite.dbMock.ExpectExec(`CREATE SCHEMA IF NOT EXISTS "test"`).WillReturnResult(sqlmock.NewResult(0, 0))
	suite.dbMock.ExpectExec("install extension").WillReturnResult(sqlmock.NewResult(0, 0))
	suite.dbMock.ExpectCommit()
	err := suite.controller.InstallExtension(mockContext(), suite.db, EXTENSION_ID, "ver")
	suite.NoError(err)
}

func (suite *ControllerUTestSuite) TestInstall_QueryFails() {
	integrationTesting.CreateTestExtensionBuilder(suite.T()).
		WithInstallFunc("context.sqlClient.execute('install extension')").
		Build().
		WriteToFile(path.Join(suite.tempExtensionRepo, EXTENSION_ID))
	suite.dbMock.ExpectBegin()
	suite.dbMock.ExpectExec(`CREATE SCHEMA IF NOT EXISTS "test"`).WillReturnResult(sqlmock.NewResult(0, 0))
	suite.dbMock.ExpectExec("install extension").WillReturnError(fmt.Errorf("mock"))
	suite.dbMock.ExpectRollback()
	err := suite.controller.InstallExtension(mockContext(), suite.db, EXTENSION_ID, "ver")
	suite.EqualError(err, "failed to install extension \"testing-extension.js\": error executing statement \"install extension\": mock")
}

func (suite *ControllerUTestSuite) TestInstallFails() {
	for _, t := range errorTests {
		suite.Run(t.testName, func() {
			integrationTesting.CreateTestExtensionBuilder(suite.T()).
				WithInstallFunc(t.throwCommand).
				Build().
				WriteToFile(path.Join(suite.tempExtensionRepo, EXTENSION_ID))
			suite.initDbMock()
			suite.dbMock.ExpectBegin()
			suite.dbMock.ExpectExec(`CREATE SCHEMA IF NOT EXISTS "test"`).WillReturnResult(sqlmock.NewResult(0, 0))
			suite.dbMock.ExpectRollback()
			err := suite.controller.InstallExtension(mockContext(), suite.db, EXTENSION_ID, "ver")
			suite.assertError(t, err)
		})
	}
}

func (suite *ControllerUTestSuite) TestAddInstance_wrongVersion() {
	integrationTesting.CreateTestExtensionBuilder(suite.T()).
		WithFindInstallationsFunc(integrationTesting.MockFindInstallationsFunction("test", "0.1.0", `[]`)).
		Build().
		WriteToFile(path.Join(suite.tempExtensionRepo, EXTENSION_ID))
	suite.metaDataMock.simulateExaAllScripts([]extensionAPI.ExaScriptRow{{Schema: "schema", Name: "script"}})
	suite.dbMock.ExpectBegin()
	suite.dbMock.ExpectExec(`CREATE SCHEMA IF NOT EXISTS "test"`).WillReturnResult(sqlmock.NewResult(0, 0))
	suite.dbMock.ExpectRollback()
	instance, err := suite.controller.CreateInstance(mockContext(), suite.db, EXTENSION_ID, "wrongVersion", []ParameterValue{})
	suite.EqualError(err, `failed to find installations: version "wrongVersion" not found for extension "testing-extension.js", available versions: ["0.1.0"]`)
	suite.Nil(instance)
}

func (suite *ControllerUTestSuite) TestAddInstance_invalidParameters() {
	integrationTesting.CreateTestExtensionBuilder(suite.T()).
		WithFindInstallationsFunc(integrationTesting.MockFindInstallationsFunction("test", "0.1.0", `[{
		id: "p1",
		name: "My param",
		type: "string",
		required: true
	}]`)).WithAddInstanceFunc("return {id: 'instId', name: `ext_${version}_${params.values[0].name}_${params.values[0].value}`};").
		Build().
		WriteToFile(path.Join(suite.tempExtensionRepo, EXTENSION_ID))
	suite.metaDataMock.simulateExaAllScripts([]extensionAPI.ExaScriptRow{})
	suite.dbMock.ExpectBegin()
	suite.dbMock.ExpectExec(`CREATE SCHEMA IF NOT EXISTS "test"`).WillReturnResult(sqlmock.NewResult(0, 0))
	suite.dbMock.ExpectRollback()
	instance, err := suite.controller.CreateInstance(mockContext(), suite.db, EXTENSION_ID, "0.1.0", []ParameterValue{})
	suite.EqualError(err, `invalid parameters: Failed to validate parameter 'My param': This is a required parameter.`)
	suite.Nil(instance)
}

func (suite *ControllerUTestSuite) TestAddInstanceFails() {
	for _, t := range errorTests {
		suite.Run(t.testName, func() {
			integrationTesting.CreateTestExtensionBuilder(suite.T()).
				WithFindInstallationsFunc(integrationTesting.MockFindInstallationsFunction("test", "0.1.0", `[]`)).
				WithAddInstanceFunc(t.throwCommand).
				Build().
				WriteToFile(path.Join(suite.tempExtensionRepo, EXTENSION_ID))
			suite.metaDataMock.simulateExaAllScripts([]extensionAPI.ExaScriptRow{})
			suite.dbMock.ExpectBegin()
			suite.dbMock.ExpectExec(`CREATE SCHEMA IF NOT EXISTS "test"`).WillReturnResult(sqlmock.NewResult(0, 0))
			suite.dbMock.ExpectRollback()
			instance, err := suite.controller.CreateInstance(mockContext(), suite.db, EXTENSION_ID, "0.1.0", []ParameterValue{})
			suite.assertError(t, err)
			suite.Nil(instance)
		})
	}
}

func (suite *ControllerUTestSuite) TestAddInstance_validParameters() {
	integrationTesting.CreateTestExtensionBuilder(suite.T()).
		WithFindInstallationsFunc(integrationTesting.MockFindInstallationsFunction("test", "0.1.0", `[{
		id: "p1",
		name: "My param",
		type: "string"
	}]`)).WithAddInstanceFunc("return {id: 'instId', name: `ext_${version}_${params.values[0].name}_${params.values[0].value}`};").
		Build().
		WriteToFile(path.Join(suite.tempExtensionRepo, EXTENSION_ID))
	suite.metaDataMock.simulateExaAllScripts([]extensionAPI.ExaScriptRow{})
	suite.dbMock.ExpectBegin()
	suite.dbMock.ExpectExec(`CREATE SCHEMA IF NOT EXISTS "test"`).WillReturnResult(sqlmock.NewResult(0, 0))
	suite.dbMock.ExpectCommit()
	instance, err := suite.controller.CreateInstance(mockContext(), suite.db, EXTENSION_ID, "0.1.0", []ParameterValue{{Name: "p1", Value: "val"}})
	suite.NoError(err)
	suite.Equal(&extensionAPI.JsExtInstance{Id: "instId", Name: "ext_0.1.0_p1_val"}, instance)
}

func (suite *ControllerUTestSuite) writeDefaultExtension() {
	integrationTesting.CreateTestExtensionBuilder(suite.T()).
		WithBucketFsUpload(integrationTesting.BucketFsUploadParams{Name: "extension jar", BucketFsFilename: "my-extension.1.2.3.jar", FileSize: 3}).
		WithFindInstallationsFunc(`
		return metadata.allScripts.rows.map(row => {
			return {name: row.schema + "." + row.name, version: "0.1.0", instanceParameters: [{id:"p1", name:"param1", type:"string"}]}
		});`).
		WithInstallFunc("context.sqlClient.execute('install extension')").
		Build().
		WriteToFile(path.Join(suite.tempExtensionRepo, EXTENSION_ID))
}

func mockContext() context.Context {
	return context.Background()
}

func (suite *ControllerUTestSuite) assertApiError(err error, expectedStatus int, expectedMessage string) {
	if apiError, ok := err.(*apiErrors.APIError); ok {
		suite.ErrorContains(apiError, expectedMessage)
		suite.Contains(apiError.Message, expectedMessage)
		suite.Equal(expectedStatus, apiError.Status)
	} else {
		suite.Fail("wrong error type", "Expected APIError but got %T: %v", err, err)
	}
}

func (suite *ControllerUTestSuite) assertNonApiError(err error, expectedMessage string) {
	if _, ok := err.(*apiErrors.APIError); ok {
		suite.Fail("wrong error type", "Expected non-APIError but got %T: %v", err, err)
	} else {
		suite.ErrorContains(err, expectedMessage)
	}
}
