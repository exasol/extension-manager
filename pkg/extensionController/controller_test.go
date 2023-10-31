package extensionController

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/exasol/extension-manager/pkg/apiErrors"
	"github.com/exasol/extension-manager/pkg/extensionAPI"
	"github.com/exasol/extension-manager/pkg/extensionAPI/exaMetadata"
	"github.com/exasol/extension-manager/pkg/extensionController/bfs"
	"github.com/exasol/extension-manager/pkg/extensionController/registry"
	"github.com/exasol/extension-manager/pkg/extensionController/transaction"
	"github.com/exasol/extension-manager/pkg/integrationTesting"
	"github.com/exasol/extension-manager/pkg/parameterValidator"

	"github.com/stretchr/testify/suite"
)

const beginTransactionFailedErrorMsg = "failed to start transaction: " + mockErrorMsg

type ControllerUTestSuite struct {
	suite.Suite
	tempExtensionRepo      string
	controller             *transactionControllerImpl
	db                     *sql.DB
	dbMock                 sqlmock.Sqlmock
	bucketFsMock           *bfs.BucketFsMock
	metaDataMock           *exaMetadata.ExaMetaDataReaderMock
	transactionStarterMock *transaction.TransactionStarterMock
}

func TestControllerUTestSuite(t *testing.T) {
	suite.Run(t, new(ControllerUTestSuite))
}

func (suite *ControllerUTestSuite) BeforeTest(suiteName, testName string) {
	suite.tempExtensionRepo = suite.T().TempDir()
	suite.createController()
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

func (suite *ControllerUTestSuite) createController() {
	suite.initDbMock()
	suite.bucketFsMock = bfs.CreateBucketFsMock()
	suite.metaDataMock = exaMetadata.CreateExaMetaDataReaderMock(EXTENSION_SCHEMA)

	config := ExtensionManagerConfig{
		ExtensionSchema:      EXTENSION_SCHEMA,
		ExtensionRegistryURL: "registryUrl",
		BucketFSBasePath:     "bfsBasePath",
	}
	ctrl := &controllerImpl{
		registry:       registry.NewRegistry(suite.tempExtensionRepo),
		config:         config,
		metaDataReader: suite.metaDataMock,
	}

	suite.transactionStarterMock = transaction.CreateTransactionStarterMock(suite.db, suite.bucketFsMock)
	suite.controller = &transactionControllerImpl{
		config:             config,
		controller:         ctrl,
		transactionStarter: suite.transactionStarterMock.GetTransactionStarter(),
	}
}

func (suite *ControllerUTestSuite) simulateTransactionBeginFails(err error) {
	suite.transactionStarterMock.SimulateTransactionFailed(err)
	suite.controller.transactionStarter = suite.transactionStarterMock.GetTransactionStarter()
}

func (suite *ControllerUTestSuite) AfterTest(suiteName, testName string) {
	suite.NoError(suite.dbMock.ExpectationsWereMet())
	suite.bucketFsMock.AssertExpectations(suite.T())
	suite.metaDataMock.AssertExpectations(suite.T())
}

// GetAllExtensions

/* [utest -> dsn~list-extensions~1]. */
func (suite *ControllerUTestSuite) TestGetAllExtensionsAndBucketFSContainsJar() {
	suite.registerDefaultExtensionDefinition()
	suite.dbMock.ExpectBegin()
	suite.bucketFsMock.SimulateFiles([]bfs.BfsFile{{Name: "my-extension.1.2.3.jar", Size: 3, Path: "path"}})
	suite.bucketFsMock.SimulateCloseSuccess()
	suite.dbMock.ExpectRollback()
	extensions, err := suite.controller.GetAllExtensions(mockContext(), suite.db)
	suite.NoError(err)
	suite.Equal([]*Extension{{Name: "MyDemoExtension", Id: "testing-extension.js", Category: "Demo category", Description: "An extension for testing.",
		InstallableVersions: []extensionAPI.JsExtensionVersion{{Name: "0.1.0", Latest: true, Deprecated: false}}}}, extensions)
}
func (suite *ControllerUTestSuite) TestGetAllExtensionsFailsStartingTransaction() {
	suite.simulateTransactionBeginFails(mockError)
	extensions, err := suite.controller.GetAllExtensions(mockContext(), suite.db)
	suite.EqualError(err, beginTransactionFailedErrorMsg)
	suite.Nil(extensions)
}

/* [utest -> dsn~list-extensions~1]. */
func (suite *ControllerUTestSuite) TestGetAllExtensionsButNoJarInBucketFS() {
	suite.registerDefaultExtensionDefinition()
	suite.dbMock.ExpectBegin()
	suite.bucketFsMock.SimulateFiles([]bfs.BfsFile{})
	suite.bucketFsMock.SimulateCloseSuccess()
	suite.dbMock.ExpectRollback()
	extensions, err := suite.controller.GetAllExtensions(mockContext(), suite.db)
	suite.NoError(err)
	suite.Empty(extensions)
}

func (suite *ControllerUTestSuite) TestGetAllExtensionsFailsForInvalidExtension() {
	suite.writeFile("broken-extension.js", "invalid javascript")
	suite.dbMock.ExpectBegin()
	suite.bucketFsMock.SimulateFiles([]bfs.BfsFile{})
	suite.bucketFsMock.SimulateCloseSuccess()
	suite.dbMock.ExpectRollback()
	extensions, err := suite.controller.GetAllExtensions(mockContext(), suite.db)
	suite.ErrorContains(err, `failed to load extension "broken-extension.js": failed to run extension "broken-extension.js" with content "invalid javascript": SyntaxError`)
	suite.Empty(extensions)
}

func (suite *ControllerUTestSuite) writeFile(fileName, content string) {
	filePath := path.Join(suite.tempExtensionRepo, fileName)
	err := os.WriteFile(filePath, []byte(content), 0600)
	if err != nil {
		suite.T().Errorf("failed to write to %q: %v", filePath, err)
	}
}

type errorTest struct {
	testName        string
	throwCommand    string
	expectedStatus  int
	expectedMessage string
}

var errorTests = []errorTest{
	{testName: "generic", throwCommand: "throw Error(`mock error from js`);", expectedStatus: -1, expectedMessage: ""},
	{testName: "internal server error", throwCommand: "throw new InternalServerError(`mock error from js`);", expectedStatus: -1, expectedMessage: ""},
	{testName: "bad request", throwCommand: "throw new BadRequestError(`mock error from js`);", expectedStatus: 400, expectedMessage: ""},
	{testName: "null pointer", throwCommand: `({}).a.b; throw Error("mock");`, expectedStatus: -1, expectedMessage: "TypeError: Cannot read property 'b' of undefined"},
}

// GetInstalledExtensions

func (suite *ControllerUTestSuite) TestGetAllInstallations() {
	suite.registerDefaultExtensionDefinition()
	//nolint:exhaustruct // Non-exhaustive struct is fine here
	suite.metaDataMock.SimulateExaAllScripts([]exaMetadata.ExaScriptRow{{Schema: "schema", Name: "script"}})
	suite.dbMock.ExpectBegin()
	suite.dbMock.ExpectRollback()
	installations, err := suite.controller.GetInstalledExtensions(mockContext(), suite.db)
	suite.NoError(err)
	suite.Equal([]*extensionAPI.JsExtInstallation{{ID: "testing-extension.js", Name: "schema.script", Version: "0.1.0"}}, installations)
}

func (suite *ControllerUTestSuite) TestGetAllInstallationsFailsStartingTransaction() {
	suite.simulateTransactionBeginFails(mockError)
	installations, err := suite.controller.GetInstalledExtensions(mockContext(), suite.db)
	suite.EqualError(err, beginTransactionFailedErrorMsg)
	suite.Nil(installations)
}

func (suite *ControllerUTestSuite) TestGetAllInstallationsReturnsEmptyList() {
	integrationTesting.CreateTestExtensionBuilder(suite.T()).
		WithFindInstallationsFunc("return []").
		Build().
		WriteToFile(path.Join(suite.tempExtensionRepo, EXTENSION_ID))

	//nolint:exhaustruct // Non-exhaustive struct is fine here
	suite.metaDataMock.SimulateExaAllScripts([]exaMetadata.ExaScriptRow{{Schema: "schema", Name: "script"}})
	suite.dbMock.ExpectBegin()
	suite.dbMock.ExpectRollback()
	installations, err := suite.controller.GetInstalledExtensions(mockContext(), suite.db)
	suite.NoError(err)
	suite.Empty(installations)
}

func (suite *ControllerUTestSuite) TestGetAllInstallationsFails() {
	for _, t := range errorTests {
		suite.Run(t.testName, func() {
			integrationTesting.CreateTestExtensionBuilder(suite.T()).
				WithFindInstallationsFunc(t.throwCommand).
				Build().
				WriteToFile(path.Join(suite.tempExtensionRepo, EXTENSION_ID))
			suite.createController()
			//nolint:exhaustruct // Empty metadata is OK for this test
			suite.metaDataMock.SimulateExaMetaData(exaMetadata.ExaMetadata{})
			suite.dbMock.ExpectBegin()
			suite.dbMock.ExpectRollback()
			extensions, err := suite.controller.GetInstalledExtensions(mockContext(), suite.db)
			suite.assertError(t, err)
			suite.Nil(extensions)
			suite.metaDataMock.AssertExpectations(suite.T())
		})
	}
}

// FindInstances

func (suite *ControllerUTestSuite) TestFindInstancesReturnsEmptyList() {
	integrationTesting.CreateTestExtensionBuilder(suite.T()).
		WithFindInstancesFunc(`return []`).
		Build().
		WriteToFile(path.Join(suite.tempExtensionRepo, EXTENSION_ID))
	suite.dbMock.ExpectBegin()
	suite.dbMock.ExpectRollback()
	extensions, err := suite.controller.FindInstances(mockContext(), suite.db, EXTENSION_ID, "ver")
	suite.NoError(err)
	suite.Empty(extensions)
}

func (suite *ControllerUTestSuite) TestFindInstancesReturnsEntries() {
	integrationTesting.CreateTestExtensionBuilder(suite.T()).
		WithFindInstancesFunc(`return [{"id":"ext1","name":"ext-name1"},{"id":"ext2","name":"ext-name2"}]`).
		Build().
		WriteToFile(path.Join(suite.tempExtensionRepo, EXTENSION_ID))
	suite.dbMock.ExpectBegin()
	suite.dbMock.ExpectRollback()
	extensions, err := suite.controller.FindInstances(mockContext(), suite.db, EXTENSION_ID, "ver")
	suite.NoError(err)
	suite.Equal([]*extensionAPI.JsExtInstance{{Id: "ext1", Name: "ext-name1"}, {Id: "ext2", Name: "ext-name2"}}, extensions)
}

func (suite *ControllerUTestSuite) TestFindInstancesFailsStartingTransaction() {
	suite.simulateTransactionBeginFails(mockError)
	extensions, err := suite.controller.FindInstances(mockContext(), suite.db, EXTENSION_ID, "ver")
	suite.EqualError(err, beginTransactionFailedErrorMsg)
	suite.Nil(extensions)
}

func (suite *ControllerUTestSuite) TestFindInstancesFails() {
	for _, t := range errorTests {
		suite.Run(t.testName, func() {
			integrationTesting.CreateTestExtensionBuilder(suite.T()).
				WithFindInstancesFunc(t.throwCommand).
				Build().
				WriteToFile(path.Join(suite.tempExtensionRepo, EXTENSION_ID))
			suite.createController()
			suite.dbMock.ExpectBegin()
			suite.dbMock.ExpectRollback()
			extensions, err := suite.controller.FindInstances(mockContext(), suite.db, EXTENSION_ID, "ver")
			suite.assertError(t, err)
			suite.Nil(extensions)
		})
	}
}

// GetParameterDefinitions

func (suite *ControllerUTestSuite) TestGetParameterDefinitionsFails() {
	for _, t := range errorTests {
		suite.Run(t.testName, func() {
			integrationTesting.CreateTestExtensionBuilder(suite.T()).
				WithGetInstanceParameterDefinitionFunc(t.throwCommand).
				Build().
				WriteToFile(path.Join(suite.tempExtensionRepo, EXTENSION_ID))
			suite.createController()
			suite.dbMock.ExpectBegin()
			suite.dbMock.ExpectRollback()
			extensions, err := suite.controller.GetParameterDefinitions(mockContext(), suite.db, EXTENSION_ID, "ver")
			suite.assertError(t, err)
			suite.Nil(extensions)
		})
	}
}

/* [utest -> dsn~parameter-versioning~1]. */
func (suite *ControllerUTestSuite) TestGetParameterDefinitionsSucceeds() {
	integrationTesting.CreateTestExtensionBuilder(suite.T()).
		WithGetInstanceParameterDefinitionFunc(`context.sqlClient.query('get param definitions'); return [{id: "param1", name: "My param:"+version, type: "string"}]`).
		Build().
		WriteToFile(path.Join(suite.tempExtensionRepo, EXTENSION_ID))
	suite.dbMock.ExpectBegin()
	suite.dbMock.ExpectQuery("param definitions").WillReturnRows(sqlmock.NewRows([]string{"col1"}))
	suite.dbMock.ExpectRollback()
	definitions, err := suite.controller.GetParameterDefinitions(mockContext(), suite.db, EXTENSION_ID, "ext-version")
	suite.NoError(err)
	suite.Equal([]parameterValidator.ParameterDefinition{{Id: "param1", Name: "My param:ext-version",
		RawDefinition: map[string]interface{}{"id": "param1", "name": "My param:ext-version", "type": "string"}}}, definitions)
}

func (suite *ControllerUTestSuite) TestGetParameterDefinitionsFailsStartingTransaction() {
	suite.simulateTransactionBeginFails(mockError)
	definitions, err := suite.controller.GetParameterDefinitions(mockContext(), suite.db, EXTENSION_ID, "ext-version")
	suite.EqualError(err, beginTransactionFailedErrorMsg)
	suite.Nil(definitions)
}

func (suite *ControllerUTestSuite) assertError(t errorTest, actualError error) {
	suite.T().Helper()
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

// InstallExtension

func (suite *ControllerUTestSuite) TestInstallFailsForUnknownExtensionId() {
	suite.dbMock.ExpectBegin()
	suite.dbMock.ExpectRollback()
	err := suite.controller.InstallExtension(mockContext(), suite.db, "unknown-extension-id", "ver")
	suite.ErrorContains(err, `failed to load extension "unknown-extension-id"`)
	suite.ErrorContains(err, `unknown-extension-id" not found`)
}

func (suite *ControllerUTestSuite) TestInstallFailsStartingTransaction() {
	suite.simulateTransactionBeginFails(mockError)
	err := suite.controller.InstallExtension(mockContext(), suite.db, EXTENSION_ID, "ver")
	suite.EqualError(err, beginTransactionFailedErrorMsg)
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

func (suite *ControllerUTestSuite) TestInstallQueryFails() {
	integrationTesting.CreateTestExtensionBuilder(suite.T()).
		WithInstallFunc("context.sqlClient.execute('install extension')").
		Build().
		WriteToFile(path.Join(suite.tempExtensionRepo, EXTENSION_ID))
	suite.dbMock.ExpectBegin()
	suite.dbMock.ExpectExec(`CREATE SCHEMA IF NOT EXISTS "test"`).WillReturnResult(sqlmock.NewResult(0, 0))
	suite.dbMock.ExpectExec("install extension").WillReturnError(fmt.Errorf("mock"))
	suite.dbMock.ExpectRollback()
	err := suite.controller.InstallExtension(mockContext(), suite.db, EXTENSION_ID, "ver")
	suite.EqualError(err, "failed to install extension \"testing-extension.js\": error executing statement 'install extension': mock")
}

func (suite *ControllerUTestSuite) TestInstallFails() {
	for _, t := range errorTests {
		suite.Run(t.testName, func() {
			integrationTesting.CreateTestExtensionBuilder(suite.T()).
				WithInstallFunc(t.throwCommand).
				Build().
				WriteToFile(path.Join(suite.tempExtensionRepo, EXTENSION_ID))
			suite.createController()
			suite.dbMock.ExpectBegin()
			suite.dbMock.ExpectExec(`CREATE SCHEMA IF NOT EXISTS "test"`).WillReturnResult(sqlmock.NewResult(0, 0))
			suite.dbMock.ExpectRollback()
			err := suite.controller.InstallExtension(mockContext(), suite.db, EXTENSION_ID, "ver")
			suite.assertError(t, err)
		})
	}
}

// UninstallExtension

func (suite *ControllerUTestSuite) TestUninstallFailsForUnknownExtensionId() {
	suite.dbMock.ExpectBegin()
	suite.dbMock.ExpectRollback()
	err := suite.controller.UninstallExtension(mockContext(), suite.db, "unknown-extension-id", "ver")
	suite.ErrorContains(err, `failed to load extension "unknown-extension-id"`)
	suite.ErrorContains(err, `unknown-extension-id" not found`)
}

func (suite *ControllerUTestSuite) TestUninstallFailsStartingTransaction() {
	suite.simulateTransactionBeginFails(mockError)
	err := suite.controller.UninstallExtension(mockContext(), suite.db, "unknown-extension-id", "ver")
	suite.EqualError(err, beginTransactionFailedErrorMsg)
}

func (suite *ControllerUTestSuite) TestUninstallSucceeds() {
	integrationTesting.CreateTestExtensionBuilder(suite.T()).
		WithUninstallFunc("context.sqlClient.execute(`uninstall extension version ${version}`)").
		Build().
		WriteToFile(path.Join(suite.tempExtensionRepo, EXTENSION_ID))
	suite.dbMock.ExpectBegin()
	suite.dbMock.ExpectExec("uninstall extension version ver").WillReturnResult(sqlmock.NewResult(0, 0))
	suite.dbMock.ExpectCommit()
	err := suite.controller.UninstallExtension(mockContext(), suite.db, EXTENSION_ID, "ver")
	suite.NoError(err)
}

func (suite *ControllerUTestSuite) TestUninstallQueryFails() {
	integrationTesting.CreateTestExtensionBuilder(suite.T()).
		WithUninstallFunc("context.sqlClient.execute(`uninstall extension version ${version}`)").
		Build().
		WriteToFile(path.Join(suite.tempExtensionRepo, EXTENSION_ID))
	suite.dbMock.ExpectBegin()
	suite.dbMock.ExpectExec("uninstall extension version ver").WillReturnError(fmt.Errorf("mock"))
	suite.dbMock.ExpectRollback()
	err := suite.controller.UninstallExtension(mockContext(), suite.db, EXTENSION_ID, "ver")
	suite.EqualError(err, "failed to uninstall extension \"testing-extension.js\": error executing statement 'uninstall extension version ver': mock")
}

func (suite *ControllerUTestSuite) TestUninstallFails() {
	for _, t := range errorTests {
		suite.Run(t.testName, func() {
			integrationTesting.CreateTestExtensionBuilder(suite.T()).
				WithUninstallFunc(t.throwCommand).
				Build().
				WriteToFile(path.Join(suite.tempExtensionRepo, EXTENSION_ID))
			suite.createController()
			suite.dbMock.ExpectBegin()
			suite.dbMock.ExpectRollback()
			err := suite.controller.UninstallExtension(mockContext(), suite.db, EXTENSION_ID, "ver")
			suite.assertError(t, err)
		})
	}
}

func (suite *ControllerUTestSuite) TestUninstallFailsWhenInstanceExists() {
	integrationTesting.CreateTestExtensionBuilder(suite.T()).
		WithUninstallFunc("context.sqlClient.execute(`uninstall extension version ${version}`)").
		WithFindInstancesFunc(`return [{"id":"ext1","name":"ext-name1"}, {"id":"ext2","name":"ext-name2"}]`).
		Build().
		WriteToFile(path.Join(suite.tempExtensionRepo, EXTENSION_ID))
	suite.dbMock.ExpectBegin()
	suite.dbMock.ExpectRollback()
	err := suite.controller.UninstallExtension(mockContext(), suite.db, EXTENSION_ID, "ver")
	suite.EqualError(err, "cannot uninstall extension because 2 instance(s) still exist: ext-name1, ext-name2")
	if apiErr, ok := apiErrors.AsAPIError(err); ok {
		suite.Equal(400, apiErr.Status)
	} else {
		suite.Fail(fmt.Sprintf("expected an APIError but got %T: %v", err, err))
	}
}

func (suite *ControllerUTestSuite) TestUninstallFailsListingInstances() {
	integrationTesting.CreateTestExtensionBuilder(suite.T()).
		WithUninstallFunc("context.sqlClient.execute(`uninstall extension version ${version}`)").
		WithFindInstancesFunc(`throw new Error('mock js error')`).
		Build().
		WriteToFile(path.Join(suite.tempExtensionRepo, EXTENSION_ID))
	suite.dbMock.ExpectBegin()
	suite.dbMock.ExpectRollback()
	err := suite.controller.UninstallExtension(mockContext(), suite.db, EXTENSION_ID, "ver")
	suite.ErrorContains(err, `failed to check existing instances: failed to list instances for extension "testing-extension.js" in version "ver": Error: mock js error`)
}

// Upgrade

func (suite *ControllerUTestSuite) TestUpgradeFailsForUnknownExtensionId() {
	suite.dbMock.ExpectBegin()
	suite.dbMock.ExpectRollback()
	result, err := suite.controller.UpgradeExtension(mockContext(), suite.db, "unknown-extension-id")
	suite.ErrorContains(err, `failed to load extension "unknown-extension-id"`)
	suite.ErrorContains(err, `unknown-extension-id" not found`)
	suite.Nil(result)
}

func (suite *ControllerUTestSuite) TestUpgradeFailsStartingTransaction() {
	suite.simulateTransactionBeginFails(mockError)
	result, err := suite.controller.UpgradeExtension(mockContext(), suite.db, "unknown-extension-id")
	suite.EqualError(err, beginTransactionFailedErrorMsg)
	suite.Nil(result)
}

/* [utest -> dsn~upgrade-extension~1]. */
func (suite *ControllerUTestSuite) TestUpgradeSucceeds() {
	integrationTesting.CreateTestExtensionBuilder(suite.T()).
		WithUpgradeFunc("context.sqlClient.execute(`upgrade extension`); return { previousVersion: 'old', newVersion: 'new' };").
		Build().
		WriteToFile(path.Join(suite.tempExtensionRepo, EXTENSION_ID))
	suite.dbMock.ExpectBegin()
	suite.dbMock.ExpectExec("upgrade extension").WillReturnResult(sqlmock.NewResult(0, 0))
	suite.dbMock.ExpectCommit()
	result, err := suite.controller.UpgradeExtension(mockContext(), suite.db, EXTENSION_ID)
	suite.NoError(err)
	suite.Equal(&extensionAPI.JsUpgradeResult{PreviousVersion: "old", NewVersion: "new"}, result)
}

func (suite *ControllerUTestSuite) TestUpgradeQueryFails() {
	integrationTesting.CreateTestExtensionBuilder(suite.T()).
		WithUpgradeFunc("context.sqlClient.execute(`upgrade extension`); return { previousVersion: 'old', newVersion: 'new' };").
		Build().
		WriteToFile(path.Join(suite.tempExtensionRepo, EXTENSION_ID))
	suite.dbMock.ExpectBegin()
	suite.dbMock.ExpectExec("upgrade extension").WillReturnError(fmt.Errorf("mock"))
	suite.dbMock.ExpectRollback()
	result, err := suite.controller.UpgradeExtension(mockContext(), suite.db, EXTENSION_ID)
	suite.EqualError(err, "failed to upgrade extension \"testing-extension.js\": error executing statement 'upgrade extension': mock")
	suite.Nil(result)
}

func (suite *ControllerUTestSuite) TestUpgradeFails() {
	for _, t := range errorTests {
		suite.Run(t.testName, func() {
			integrationTesting.CreateTestExtensionBuilder(suite.T()).
				WithUpgradeFunc(t.throwCommand).
				Build().
				WriteToFile(path.Join(suite.tempExtensionRepo, EXTENSION_ID))
			suite.createController()
			suite.dbMock.ExpectBegin()
			suite.dbMock.ExpectRollback()
			result, err := suite.controller.UpgradeExtension(mockContext(), suite.db, EXTENSION_ID)
			suite.assertError(t, err)
			suite.Nil(result)
		})
	}
}

// CreateInstance

/* [utest -> dsn~parameter-types~1]. */
func (suite *ControllerUTestSuite) TestCreateInstanceInvalidParameters() {
	integrationTesting.CreateTestExtensionBuilder(suite.T()).
		WithFindInstallationsFunc(integrationTesting.MockFindInstallationsFunction("test", "0.1.0")).
		WithAddInstanceFunc("throw new Error('This should not be called.')").
		WithGetInstanceParameterDefinitionFunc(`return [{id: "param1", name: "My param", type: "string", required: true}]`).
		Build().
		WriteToFile(path.Join(suite.tempExtensionRepo, EXTENSION_ID))
	suite.dbMock.ExpectBegin()
	suite.dbMock.ExpectExec(`CREATE SCHEMA IF NOT EXISTS "test"`).WillReturnResult(sqlmock.NewResult(0, 0))
	suite.dbMock.ExpectRollback()
	instance, err := suite.controller.CreateInstance(mockContext(), suite.db, EXTENSION_ID, "0.1.0", []ParameterValue{})
	suite.EqualError(err, `invalid parameters: Failed to validate parameter 'My param' (param1): This is a required parameter.`)
	suite.Nil(instance)
}

func (suite *ControllerUTestSuite) TestCreateInstanceFailsStartingTransaction() {
	suite.simulateTransactionBeginFails(mockError)
	instance, err := suite.controller.CreateInstance(mockContext(), suite.db, EXTENSION_ID, "0.1.0", []ParameterValue{})
	suite.EqualError(err, beginTransactionFailedErrorMsg)
	suite.Nil(instance)
}

func (suite *ControllerUTestSuite) TestCreateInstanceFails() {
	for _, t := range errorTests {
		suite.Run(t.testName, func() {
			integrationTesting.CreateTestExtensionBuilder(suite.T()).
				WithFindInstallationsFunc(integrationTesting.MockFindInstallationsFunction("test", "0.1.0")).
				WithAddInstanceFunc(t.throwCommand).
				Build().
				WriteToFile(path.Join(suite.tempExtensionRepo, EXTENSION_ID))
			suite.dbMock.ExpectBegin()
			suite.dbMock.ExpectExec(`CREATE SCHEMA IF NOT EXISTS "test"`).WillReturnResult(sqlmock.NewResult(0, 0))
			suite.dbMock.ExpectRollback()
			instance, err := suite.controller.CreateInstance(mockContext(), suite.db, EXTENSION_ID, "0.1.0", []ParameterValue{})
			suite.assertError(t, err)
			suite.Nil(instance)
		})
	}
}

func (suite *ControllerUTestSuite) TestCreateInstanceValidParameters() {
	integrationTesting.CreateTestExtensionBuilder(suite.T()).
		WithFindInstallationsFunc(integrationTesting.MockFindInstallationsFunction("test", "0.1.0")).
		WithAddInstanceFunc("return {id: 'instId', name: `ext_${version}_${params.values[0].name}_${params.values[0].value}`};").
		Build().
		WriteToFile(path.Join(suite.tempExtensionRepo, EXTENSION_ID))
	suite.dbMock.ExpectBegin()
	suite.dbMock.ExpectExec(`CREATE SCHEMA IF NOT EXISTS "test"`).WillReturnResult(sqlmock.NewResult(0, 0))
	suite.dbMock.ExpectCommit()
	instance, err := suite.controller.CreateInstance(mockContext(), suite.db, EXTENSION_ID, "0.1.0", []ParameterValue{{Name: "p1", Value: "val"}})
	suite.NoError(err)
	suite.Equal(&extensionAPI.JsExtInstance{Id: "instId", Name: "ext_0.1.0_p1_val"}, instance)
}

// DeleteInstance

func (suite *ControllerUTestSuite) TestDeleteInstancesFails() {
	for _, t := range errorTests {
		suite.Run(t.testName, func() {
			integrationTesting.CreateTestExtensionBuilder(suite.T()).
				WithDeleteInstanceFunc(t.throwCommand).
				Build().
				WriteToFile(path.Join(suite.tempExtensionRepo, EXTENSION_ID))
			suite.createController()
			suite.dbMock.ExpectBegin()
			suite.dbMock.ExpectRollback()
			err := suite.controller.DeleteInstance(mockContext(), suite.db, EXTENSION_ID, "extVersion", "instId")
			suite.assertError(t, err)
		})
	}
}

func (suite *ControllerUTestSuite) TestDeleteInstanceSucceeds() {
	integrationTesting.CreateTestExtensionBuilder(suite.T()).
		WithDeleteInstanceFunc("context.sqlClient.execute(`delete instance ${instanceId}`)").
		Build().
		WriteToFile(path.Join(suite.tempExtensionRepo, EXTENSION_ID))
	suite.dbMock.ExpectBegin()
	suite.dbMock.ExpectExec("delete instance instId").WillReturnResult(sqlmock.NewResult(0, 0))
	suite.dbMock.ExpectCommit()
	err := suite.controller.DeleteInstance(mockContext(), suite.db, EXTENSION_ID, "extVersion", "instId")
	suite.NoError(err)
}

func (suite *ControllerUTestSuite) TestDeleteInstanceFailsStartingTransaction() {
	suite.simulateTransactionBeginFails(mockError)
	err := suite.controller.DeleteInstance(mockContext(), suite.db, EXTENSION_ID, "extVersion", "instId")
	suite.EqualError(err, beginTransactionFailedErrorMsg)
}

func (suite *ControllerUTestSuite) registerDefaultExtensionDefinition() {
	integrationTesting.CreateTestExtensionBuilder(suite.T()).
		WithBucketFsUpload(integrationTesting.BucketFsUploadParams{Name: "extension jar", BucketFsFilename: "my-extension.1.2.3.jar", FileSize: 3, DownloadUrl: "", LicenseUrl: "", LicenseAgreementRequired: false}).
		WithFindInstallationsFunc(`
		return metadata.allScripts.rows.map(row => {
			return {name: row.schema + "." + row.name, version: "0.1.0"}
		});`).
		WithInstallFunc("context.sqlClient.execute('install extension')").
		Build().
		WriteToFile(path.Join(suite.tempExtensionRepo, EXTENSION_ID))
}

func mockContext() context.Context {
	return context.Background()
}

func (suite *ControllerUTestSuite) assertApiError(err error, expectedStatus int, expectedMessage string) {
	suite.T().Helper()
	if apiError, ok := apiErrors.AsAPIError(err); ok {
		suite.ErrorContains(apiError, expectedMessage)
		suite.Contains(apiError.Message, expectedMessage)
		suite.Equal(expectedStatus, apiError.Status)
	} else {
		suite.Fail("wrong error type", "Expected APIError but got %T: %v", err, err)
	}
}

func (suite *ControllerUTestSuite) assertNonApiError(err error, expectedMessage string) {
	suite.T().Helper()
	if _, ok := apiErrors.AsAPIError(err); ok {
		suite.Fail("wrong error type", "Expected non-APIError but got %T: %v", err, err)
	} else {
		suite.ErrorContains(err, expectedMessage)
	}
}
