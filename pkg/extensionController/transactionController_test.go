package extensionController

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/exasol/extension-manager/pkg/extensionAPI"
	"github.com/exasol/extension-manager/pkg/extensionController/bfs"
	"github.com/exasol/extension-manager/pkg/extensionController/transaction"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

const mockErrorMsg = "mock error"

var mockError = fmt.Errorf(mockErrorMsg)

const beginMockTransactionFailedErrorMsg = "failed to start transaction: failed to start mock transaction: " + mockErrorMsg

type extCtrlUnitTestSuite struct {
	suite.Suite
	ctrl                   TransactionController
	db                     *sql.DB
	dbMock                 sqlmock.Sqlmock
	mockCtrl               mockControllerImpl
	bucketFsMock           *bfs.BucketFsMock
	transactionStarterMock *transaction.TransactionStarterMock
}

func TestExtensionControllerUnitTestSuite(t *testing.T) {
	suite.Run(t, new(extCtrlUnitTestSuite))
}

func (suite *extCtrlUnitTestSuite) SetupTest() {
	suite.setupDbMock()
	suite.mockCtrl = createMockControllerImpl()
	suite.bucketFsMock = bfs.CreateBucketFsMock()
	suite.transactionStarterMock = transaction.CreateTransactionStarterMock(suite.db, suite.bucketFsMock)
	suite.ctrl = &transactionControllerImpl{
		controller:         &suite.mockCtrl,
		transactionStarter: suite.transactionStarterMock.GetTransactionStarter(),
		config: ExtensionManagerConfig{
			ExtensionRegistryURL: "registry-url",
			BucketFSBasePath:     "bfs-base-path",
			ExtensionSchema:      "ext-schema",
		},
	}
}

func (suite *extCtrlUnitTestSuite) setupDbMock() {
	db, dbMock, err := sqlmock.New()
	if err != nil {
		suite.Failf("error '%v' was not expected when opening a stub database connection", err.Error())
	}
	suite.db = db
	suite.dbMock = dbMock
	suite.dbMock.MatchExpectationsInOrder(true)
}

func (suite *extCtrlUnitTestSuite) AfterTest(suiteName, testName string) {
	if err := suite.dbMock.ExpectationsWereMet(); err != nil {
		suite.Failf("unfulfilled expectations", err.Error())
	}
	suite.bucketFsMock.AssertExpectations(suite.T())
	suite.mockCtrl.AssertExpectations(suite.T())
}

// CreateWithValidatedConfig

func (suite *extCtrlUnitTestSuite) TestCreateWithValidatedConfigSuccess() {
	ctrl, err := CreateWithValidatedConfig(ExtensionManagerConfig{ExtensionRegistryURL: "url", BucketFSBasePath: "bfspath", ExtensionSchema: "schema"})
	suite.NoError(err)
	suite.NotNil(ctrl)
}

func (suite *extCtrlUnitTestSuite) TestCreateWithValidatedConfigFailure() {
	var tests = []struct {
		name          string
		config        ExtensionManagerConfig
		expectedError string
	}{
		{name: "missing registry url", config: ExtensionManagerConfig{ExtensionRegistryURL: "", BucketFSBasePath: "bfspath", ExtensionSchema: "schema"}, expectedError: "invalid configuration: missing ExtensionRegistryURL"},
		{name: "empty registry url", config: ExtensionManagerConfig{ExtensionRegistryURL: "", BucketFSBasePath: "bfspath", ExtensionSchema: "schema"}, expectedError: "invalid configuration: missing ExtensionRegistryURL"},
		{name: "missing bucketfs base path", config: ExtensionManagerConfig{ExtensionRegistryURL: "url", BucketFSBasePath: "", ExtensionSchema: "schema"}, expectedError: "invalid configuration: missing BucketFSBasePath"},
		{name: "empty bucketfs base path", config: ExtensionManagerConfig{ExtensionRegistryURL: "url", BucketFSBasePath: "", ExtensionSchema: "schema"}, expectedError: "invalid configuration: missing BucketFSBasePath"},
		{name: "missing schema", config: ExtensionManagerConfig{ExtensionRegistryURL: "url", BucketFSBasePath: "bfspath", ExtensionSchema: ""}, expectedError: "invalid configuration: missing ExtensionSchema"},
		{name: "empty schema", config: ExtensionManagerConfig{ExtensionRegistryURL: "url", BucketFSBasePath: "bfspath", ExtensionSchema: ""}, expectedError: "invalid configuration: missing ExtensionSchema"},
		{name: "all missing", config: ExtensionManagerConfig{ExtensionRegistryURL: "", BucketFSBasePath: "", ExtensionSchema: ""}, expectedError: "invalid configuration: missing BucketFSBasePath"},
	}
	for _, test := range tests {
		suite.Run(test.name, func() {
			ctrl, err := CreateWithValidatedConfig(test.config)
			suite.EqualError(err, test.expectedError)
			suite.Nil(ctrl)
		})
	}
}

// GetAllExtensions

func (suite *extCtrlUnitTestSuite) TestGetAllExtensionsSuccess() {
	suite.dbMock.ExpectBegin()
	suite.bucketFsMock.SimulateFiles([]bfs.BfsFile{})
	suite.bucketFsMock.SimulateCloseSuccess()
	suite.mockCtrl.On("GetAllExtensions", mock.Anything).Return([]*Extension{}, nil)
	extensions, err := suite.ctrl.GetAllExtensions(mockContext(), suite.db)
	suite.NoError(err)
	suite.Len(extensions, 0)
}

func (suite *extCtrlUnitTestSuite) TestGetAllExtensionsBucketFsListFails() {
	suite.dbMock.ExpectBegin()
	suite.bucketFsMock.SimulateFilesError(mockError)
	suite.bucketFsMock.SimulateCloseSuccess()
	extensions, err := suite.ctrl.GetAllExtensions(mockContext(), suite.db)
	suite.EqualError(err, "failed to search for required files in BucketFS. Cause: mock error")
	suite.Nil(extensions)
}

func (suite *extCtrlUnitTestSuite) TestGetAllExtensionsGetFails() {
	suite.dbMock.ExpectBegin()
	suite.bucketFsMock.SimulateFiles([]bfs.BfsFile{})
	suite.bucketFsMock.SimulateCloseSuccess()
	suite.mockCtrl.On("GetAllExtensions", mock.Anything).Return(nil, mockError)
	extensions, err := suite.ctrl.GetAllExtensions(mockContext(), suite.db)
	suite.EqualError(err, mockErrorMsg)
	suite.Nil(extensions)
}

func (suite *extCtrlUnitTestSuite) TestGetAllInstallationsBeginTransactionFailure() {
	suite.dbMock.ExpectBegin().WillReturnError(mockError)
	installations, err := suite.ctrl.GetInstalledExtensions(mockContext(), suite.db)
	suite.EqualError(err, beginMockTransactionFailedErrorMsg)
	suite.Nil(installations)
}

// GetInstalledExtensions

func (suite *extCtrlUnitTestSuite) TestGetInstalledExtensionsSuccess() {
	suite.dbMock.ExpectBegin()
	mockResult := []*extensionAPI.JsExtInstallation{{ID: "mock-ID", Name: "ext", Version: "mock-version"}}
	suite.mockCtrl.On("GetAllInstallations", mock.Anything).Return(mockResult, nil)
	suite.dbMock.ExpectRollback()
	installations, err := suite.ctrl.GetInstalledExtensions(mockContext(), suite.db)
	suite.NoError(err)
	suite.Equal(mockResult, installations)
}

func (suite *extCtrlUnitTestSuite) TestGetInstalledExtensionsFailure() {
	suite.dbMock.ExpectBegin()
	suite.mockCtrl.On("GetAllInstallations", mock.Anything).Return(nil, mockError)
	suite.dbMock.ExpectRollback()
	installations, err := suite.ctrl.GetInstalledExtensions(mockContext(), suite.db)
	suite.EqualError(err, mockErrorMsg)
	suite.Nil(installations)
}

// InstallExtension

func (suite *extCtrlUnitTestSuite) TestInstallExtensionBeginTransactionFailure() {
	suite.dbMock.ExpectBegin().WillReturnError(mockError)
	err := suite.ctrl.InstallExtension(mockContext(), suite.db, "extId", "extVer")
	suite.EqualError(err, beginMockTransactionFailedErrorMsg)
}

func (suite *extCtrlUnitTestSuite) TestInstallExtensionSuccess() {
	suite.dbMock.ExpectBegin()
	suite.mockCtrl.On("InstallExtension", mock.Anything, "extId", "extVer").Return(nil)
	suite.dbMock.ExpectCommit()
	err := suite.ctrl.InstallExtension(mockContext(), suite.db, "extId", "extVer")
	suite.NoError(err)
}

func (suite *extCtrlUnitTestSuite) TestInstallExtensionFailureRollback() {
	suite.dbMock.ExpectBegin()
	suite.mockCtrl.On("InstallExtension", mock.Anything, "extId", "extVer").Return(mockError)
	suite.dbMock.ExpectRollback()
	err := suite.ctrl.InstallExtension(mockContext(), suite.db, "extId", "extVer")
	suite.EqualError(err, mockErrorMsg)
}

func (suite *extCtrlUnitTestSuite) TestInstallExtensionCommitFailure() {
	suite.dbMock.ExpectBegin()
	suite.mockCtrl.On("InstallExtension", mock.Anything, "extId", "extVer").Return(nil)
	suite.dbMock.ExpectCommit().WillReturnError(mockError)
	err := suite.ctrl.InstallExtension(mockContext(), suite.db, "extId", "extVer")
	suite.EqualError(err, mockErrorMsg)
}

// UninstallExtension

func (suite *extCtrlUnitTestSuite) TestUninstallExtensionBeginTransactionFailure() {
	suite.dbMock.ExpectBegin().WillReturnError(mockError)
	err := suite.ctrl.UninstallExtension(mockContext(), suite.db, "extId", "extVer")
	suite.EqualError(err, beginMockTransactionFailedErrorMsg)
}

func (suite *extCtrlUnitTestSuite) TestUninstallExtensionSuccess() {
	suite.dbMock.ExpectBegin()
	suite.mockCtrl.On("UninstallExtension", mock.Anything, "extId", "extVer").Return(nil)
	suite.dbMock.ExpectCommit()
	err := suite.ctrl.UninstallExtension(mockContext(), suite.db, "extId", "extVer")
	suite.NoError(err)
}

func (suite *extCtrlUnitTestSuite) TestUninstallExtensionFailureRollback() {
	suite.dbMock.ExpectBegin()
	suite.mockCtrl.On("UninstallExtension", mock.Anything, "extId", "extVer").Return(mockError)
	suite.dbMock.ExpectRollback()
	err := suite.ctrl.UninstallExtension(mockContext(), suite.db, "extId", "extVer")
	suite.EqualError(err, mockErrorMsg)
}

func (suite *extCtrlUnitTestSuite) TestUninstallExtensionCommitFailure() {
	suite.dbMock.ExpectBegin()
	suite.mockCtrl.On("UninstallExtension", mock.Anything, "extId", "extVer").Return(nil)
	suite.dbMock.ExpectCommit().WillReturnError(mockError)
	err := suite.ctrl.UninstallExtension(mockContext(), suite.db, "extId", "extVer")
	suite.EqualError(err, mockErrorMsg)
}

// Upgrade

func (suite *extCtrlUnitTestSuite) TestUpgradeBeginTransactionFailure() {
	suite.dbMock.ExpectBegin().WillReturnError(mockError)
	result, err := suite.ctrl.UpgradeExtension(mockContext(), suite.db, "extId")
	suite.EqualError(err, beginMockTransactionFailedErrorMsg)
	suite.Nil(result)
}

/* [utest -> dsn~upgrade-extension~1]. */
func (suite *extCtrlUnitTestSuite) TestUpgradeSuccess() {
	suite.dbMock.ExpectBegin()
	suite.mockCtrl.On("UpgradeExtension", mock.Anything, "extId").Return(&extensionAPI.JsUpgradeResult{PreviousVersion: "old", NewVersion: "new"}, nil)
	suite.dbMock.ExpectCommit()
	result, err := suite.ctrl.UpgradeExtension(mockContext(), suite.db, "extId")
	suite.NoError(err)
	suite.Equal(&extensionAPI.JsUpgradeResult{PreviousVersion: "old", NewVersion: "new"}, result)
}

func (suite *extCtrlUnitTestSuite) TestUpgradeFailure() {
	suite.dbMock.ExpectBegin()
	suite.mockCtrl.On("UpgradeExtension", mock.Anything, "extId").Return(nil, mockError)
	suite.dbMock.ExpectRollback()
	result, err := suite.ctrl.UpgradeExtension(mockContext(), suite.db, "extId")
	suite.EqualError(err, mockErrorMsg)
	suite.Nil(result)
}

func (suite *extCtrlUnitTestSuite) TestUpgradeCommitFailure() {
	suite.dbMock.ExpectBegin()
	suite.mockCtrl.On("UpgradeExtension", mock.Anything, "extId").Return(&extensionAPI.JsUpgradeResult{PreviousVersion: "old", NewVersion: "new"}, nil)
	suite.dbMock.ExpectCommit().WillReturnError(mockError)
	result, err := suite.ctrl.UpgradeExtension(mockContext(), suite.db, "extId")
	suite.EqualError(err, mockErrorMsg)
	suite.Nil(result)
}

// CreateInstance

func (suite *extCtrlUnitTestSuite) TestCreateInstanceBeginTransactionFailure() {
	suite.dbMock.ExpectBegin().WillReturnError(mockError)
	instance, err := suite.ctrl.CreateInstance(mockContext(), suite.db, "extId", "extVer", []ParameterValue{})
	suite.EqualError(err, beginMockTransactionFailedErrorMsg)
	suite.Nil(instance)
}

func (suite *extCtrlUnitTestSuite) TestCreateInstanceSuccess() {
	suite.dbMock.ExpectBegin()
	suite.mockCtrl.On("CreateInstance", mock.Anything, "extId", "extVer", mock.Anything).Return(&extensionAPI.JsExtInstance{Id: "instId", Name: "newInst"}, nil)
	suite.dbMock.ExpectCommit()
	instance, err := suite.ctrl.CreateInstance(mockContext(), suite.db, "extId", "extVer", []ParameterValue{})
	suite.NoError(err)
	suite.Equal(&extensionAPI.JsExtInstance{Id: "instId", Name: "newInst"}, instance)
}

func (suite *extCtrlUnitTestSuite) TestCreateInstanceFailure() {
	suite.dbMock.ExpectBegin()
	suite.mockCtrl.On("CreateInstance", mock.Anything, "extId", "extVer", mock.Anything).Return(nil, mockError)
	suite.dbMock.ExpectRollback()
	instance, err := suite.ctrl.CreateInstance(mockContext(), suite.db, "extId", "extVer", []ParameterValue{})
	suite.EqualError(err, mockErrorMsg)
	suite.Nil(instance)
}

func (suite *extCtrlUnitTestSuite) TestCreateInstanceCommitFailure() {
	suite.dbMock.ExpectBegin()
	suite.mockCtrl.On("CreateInstance", mock.Anything, "extId", "extVer", mock.Anything).Return(&extensionAPI.JsExtInstance{Id: "instId", Name: "newInst"}, nil)
	suite.dbMock.ExpectCommit().WillReturnError(mockError)
	instance, err := suite.ctrl.CreateInstance(mockContext(), suite.db, "extId", "extVer", []ParameterValue{})
	suite.EqualError(err, mockErrorMsg)
	suite.Nil(instance)
}

// FindInstances

func (suite *extCtrlUnitTestSuite) TestFindInstancesBeginTransactionFailure() {
	suite.dbMock.ExpectBegin().WillReturnError(mockError)
	instances, err := suite.ctrl.FindInstances(mockContext(), suite.db, "extId", "extVer")
	suite.EqualError(err, beginMockTransactionFailedErrorMsg)
	suite.Nil(instances)
}

func (suite *extCtrlUnitTestSuite) TestFindInstancesSuccess() {
	suite.dbMock.ExpectBegin()
	suite.mockCtrl.On("FindInstances", mock.Anything, "extId", "extVer").Return([]*extensionAPI.JsExtInstance{{Id: "instId", Name: "newInst"}}, nil)
	suite.dbMock.ExpectRollback()
	instances, err := suite.ctrl.FindInstances(mockContext(), suite.db, "extId", "extVer")
	suite.NoError(err)
	suite.Equal([]*extensionAPI.JsExtInstance{{Id: "instId", Name: "newInst"}}, instances)
}

func (suite *extCtrlUnitTestSuite) TestFindInstancesFailure() {
	suite.dbMock.ExpectBegin()
	suite.mockCtrl.On("FindInstances", mock.Anything, "extId", "extVer").Return(nil, mockError)
	suite.dbMock.ExpectRollback()
	instances, err := suite.ctrl.FindInstances(mockContext(), suite.db, "extId", "extVer")
	suite.EqualError(err, mockErrorMsg)
	suite.Nil(instances)
}

// DeleteInstance

func (suite *extCtrlUnitTestSuite) TestDeleteInstanceBeginTransactionFailure() {
	suite.dbMock.ExpectBegin().WillReturnError(mockError)
	err := suite.ctrl.DeleteInstance(mockContext(), suite.db, "extId", "extVers", "instId")
	suite.EqualError(err, beginMockTransactionFailedErrorMsg)
}

func (suite *extCtrlUnitTestSuite) TestDeleteInstanceSuccess() {
	suite.dbMock.ExpectBegin()
	suite.mockCtrl.On("DeleteInstance", mock.Anything, "extId", "extVers", "instId", mock.Anything).Return(nil)
	suite.dbMock.ExpectCommit()
	err := suite.ctrl.DeleteInstance(mockContext(), suite.db, "extId", "extVers", "instId")
	suite.NoError(err)
}

func (suite *extCtrlUnitTestSuite) TestDeleteInstanceFailure() {
	suite.dbMock.ExpectBegin()
	suite.mockCtrl.On("DeleteInstance", mock.Anything, "extId", "extVers", "instId", mock.Anything).Return(mockError)
	suite.dbMock.ExpectRollback()
	err := suite.ctrl.DeleteInstance(mockContext(), suite.db, "extId", "extVers", "instId")
	suite.EqualError(err, mockErrorMsg)
}

func (suite *extCtrlUnitTestSuite) TestDeleteInstanceCommitFailure() {
	suite.dbMock.ExpectBegin()
	suite.mockCtrl.On("DeleteInstance", mock.Anything, "extId", "extVers", "instId", mock.Anything).Return(nil)
	suite.dbMock.ExpectCommit().WillReturnError(mockError)
	err := suite.ctrl.DeleteInstance(mockContext(), suite.db, "extId", "extVers", "instId")
	suite.EqualError(err, mockErrorMsg)
}
