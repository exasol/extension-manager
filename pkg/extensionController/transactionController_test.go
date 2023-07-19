package extensionController

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/exasol/extension-manager/pkg/extensionAPI"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type extCtrlUnitTestSuite struct {
	suite.Suite
	ctrl     TransactionController
	db       *sql.DB
	dbMock   sqlmock.Sqlmock
	mockCtrl mockControllerImpl
	mockBfs  bucketFsMock
}

func TestExtensionControllerUnitTestSuite(t *testing.T) {
	suite.Run(t, new(extCtrlUnitTestSuite))
}

func (suite *extCtrlUnitTestSuite) SetupTest() {
	suite.mockCtrl = mockControllerImpl{}
	suite.mockBfs = bucketFsMock{}
	suite.ctrl = &transactionControllerImpl{controller: &suite.mockCtrl, bucketFs: &suite.mockBfs}
	db, dbMock, err := sqlmock.New()
	if err != nil {
		suite.Failf("error '%v' was not expected when opening a stub database connection", err.Error())
	}
	dbMock.MatchExpectationsInOrder(true)
	suite.db = db
	suite.dbMock = dbMock
}

func (suite *extCtrlUnitTestSuite) AfterTest(suiteName, testName string) {
	if err := suite.dbMock.ExpectationsWereMet(); err != nil {
		suite.Failf("unfulfilled expectations", err.Error())
	}
}

// GetAllExtensions

func (suite *extCtrlUnitTestSuite) TestGetAllExtensions_Success() {
	suite.mockBfs.On("ListFiles", mock.Anything, mock.Anything).Return([]BfsFile{}, nil)
	suite.mockCtrl.On("GetAllExtensions", mock.Anything).Return([]*Extension{}, nil)
	extensions, err := suite.ctrl.GetAllExtensions(mockContext(), suite.db)
	suite.NoError(err)
	suite.Len(extensions, 0)
}

func (suite *extCtrlUnitTestSuite) TestGetAllExtensions_BucketFsListFails() {
	suite.mockBfs.On("ListFiles", mock.Anything, mock.Anything).Return(nil, fmt.Errorf("mock"))
	extensions, err := suite.ctrl.GetAllExtensions(mockContext(), suite.db)
	suite.EqualError(err, "failed to search for required files in BucketFS. Cause: mock")
	suite.Nil(extensions)
}

func (suite *extCtrlUnitTestSuite) TestGetAllExtensions_GetFails() {
	suite.mockBfs.On("ListFiles", mock.Anything, mock.Anything).Return([]BfsFile{}, nil)
	suite.mockCtrl.On("GetAllExtensions", mock.Anything).Return(nil, fmt.Errorf("mock"))
	extensions, err := suite.ctrl.GetAllExtensions(mockContext(), suite.db)
	suite.EqualError(err, "mock")
	suite.Nil(extensions)
}

func (suite *extCtrlUnitTestSuite) TestGetAllInstallations_BeginTransactionFailure() {
	suite.dbMock.ExpectBegin().WillReturnError(fmt.Errorf("mock"))
	installations, err := suite.ctrl.GetInstalledExtensions(mockContext(), suite.db)
	suite.EqualError(err, "failed to begin transaction: mock")
	suite.Nil(installations)
}

// GetInstalledExtensions

func (suite *extCtrlUnitTestSuite) TestGetInstalledExtensions_Success() {
	suite.dbMock.ExpectBegin()
	mockResult := []*extensionAPI.JsExtInstallation{{Name: "ext"}}
	suite.mockCtrl.On("GetAllInstallations", mock.Anything).Return(mockResult, nil)
	suite.dbMock.ExpectRollback()
	installations, err := suite.ctrl.GetInstalledExtensions(mockContext(), suite.db)
	suite.NoError(err)
	suite.Equal(mockResult, installations)
}

func (suite *extCtrlUnitTestSuite) TestGetInstalledExtensions_Failure() {
	suite.dbMock.ExpectBegin()
	suite.mockCtrl.On("GetAllInstallations", mock.Anything).Return(nil, fmt.Errorf("mock"))
	suite.dbMock.ExpectRollback()
	installations, err := suite.ctrl.GetInstalledExtensions(mockContext(), suite.db)
	suite.EqualError(err, "mock")
	suite.Nil(installations)
}

// InstallExtension

func (suite *extCtrlUnitTestSuite) TestInstallExtension_BeginTransactionFailure() {
	suite.dbMock.ExpectBegin().WillReturnError(fmt.Errorf("mock"))
	err := suite.ctrl.InstallExtension(mockContext(), suite.db, "extId", "extVer")
	suite.EqualError(err, "failed to begin transaction: mock")
}

func (suite *extCtrlUnitTestSuite) TestInstallExtension_Success() {
	suite.dbMock.ExpectBegin()
	suite.mockCtrl.On("InstallExtension", mock.Anything, "extId", "extVer").Return(nil)
	suite.dbMock.ExpectCommit()
	err := suite.ctrl.InstallExtension(mockContext(), suite.db, "extId", "extVer")
	suite.NoError(err)
}

func (suite *extCtrlUnitTestSuite) TestInstallExtension_FailureRollback() {
	suite.dbMock.ExpectBegin()
	suite.mockCtrl.On("InstallExtension", mock.Anything, "extId", "extVer").Return(fmt.Errorf("mock"))
	suite.dbMock.ExpectRollback()
	err := suite.ctrl.InstallExtension(mockContext(), suite.db, "extId", "extVer")
	suite.EqualError(err, "mock")
}

func (suite *extCtrlUnitTestSuite) TestInstallExtension_CommitFailure() {
	suite.dbMock.ExpectBegin()
	suite.mockCtrl.On("InstallExtension", mock.Anything, "extId", "extVer").Return(nil)
	suite.dbMock.ExpectCommit().WillReturnError(fmt.Errorf("commit failed"))
	err := suite.ctrl.InstallExtension(mockContext(), suite.db, "extId", "extVer")
	suite.EqualError(err, "commit failed")
}

// UninstallExtension

func (suite *extCtrlUnitTestSuite) TestUninstallExtension_BeginTransactionFailure() {
	suite.dbMock.ExpectBegin().WillReturnError(fmt.Errorf("mock"))
	err := suite.ctrl.UninstallExtension(mockContext(), suite.db, "extId", "extVer")
	suite.EqualError(err, "failed to begin transaction: mock")
}

func (suite *extCtrlUnitTestSuite) TestUninstallExtension_Success() {
	suite.dbMock.ExpectBegin()
	suite.mockCtrl.On("UninstallExtension", mock.Anything, "extId", "extVer").Return(nil)
	suite.dbMock.ExpectCommit()
	err := suite.ctrl.UninstallExtension(mockContext(), suite.db, "extId", "extVer")
	suite.NoError(err)
}

func (suite *extCtrlUnitTestSuite) TestUninstallExtension_FailureRollback() {
	suite.dbMock.ExpectBegin()
	suite.mockCtrl.On("UninstallExtension", mock.Anything, "extId", "extVer").Return(fmt.Errorf("mock"))
	suite.dbMock.ExpectRollback()
	err := suite.ctrl.UninstallExtension(mockContext(), suite.db, "extId", "extVer")
	suite.EqualError(err, "mock")
}

func (suite *extCtrlUnitTestSuite) TestUninstallExtension_CommitFailure() {
	suite.dbMock.ExpectBegin()
	suite.mockCtrl.On("UninstallExtension", mock.Anything, "extId", "extVer").Return(nil)
	suite.dbMock.ExpectCommit().WillReturnError(fmt.Errorf("commit failed"))
	err := suite.ctrl.UninstallExtension(mockContext(), suite.db, "extId", "extVer")
	suite.EqualError(err, "commit failed")
}

// CreateInstance

func (suite *extCtrlUnitTestSuite) TestCreateInstance_BeginTransactionFailure() {
	suite.dbMock.ExpectBegin().WillReturnError(fmt.Errorf("mock"))
	instance, err := suite.ctrl.CreateInstance(mockContext(), suite.db, "extId", "extVer", []ParameterValue{})
	suite.EqualError(err, "failed to begin transaction: mock")
	suite.Nil(instance)
}

func (suite *extCtrlUnitTestSuite) TestCreateInstance_Success() {
	suite.dbMock.ExpectBegin()
	suite.mockCtrl.On("CreateInstance", mock.Anything, "extId", "extVer", mock.Anything).Return(&extensionAPI.JsExtInstance{Id: "instId", Name: "newInst"}, nil)
	suite.dbMock.ExpectCommit()
	instance, err := suite.ctrl.CreateInstance(mockContext(), suite.db, "extId", "extVer", []ParameterValue{})
	suite.NoError(err)
	suite.Equal(&extensionAPI.JsExtInstance{Id: "instId", Name: "newInst"}, instance)
}

func (suite *extCtrlUnitTestSuite) TestCreateInstance_Failure() {
	suite.dbMock.ExpectBegin()
	suite.mockCtrl.On("CreateInstance", mock.Anything, "extId", "extVer", mock.Anything).Return(nil, fmt.Errorf("mock"))
	suite.dbMock.ExpectRollback()
	instance, err := suite.ctrl.CreateInstance(mockContext(), suite.db, "extId", "extVer", []ParameterValue{})
	suite.EqualError(err, "mock")
	suite.Nil(instance)
}

func (suite *extCtrlUnitTestSuite) TestCreateInstance_CommitFailure() {
	suite.dbMock.ExpectBegin()
	suite.mockCtrl.On("CreateInstance", mock.Anything, "extId", "extVer", mock.Anything).Return(&extensionAPI.JsExtInstance{Id: "instId", Name: "newInst"}, nil)
	suite.dbMock.ExpectCommit().WillReturnError(fmt.Errorf("mock"))
	instance, err := suite.ctrl.CreateInstance(mockContext(), suite.db, "extId", "extVer", []ParameterValue{})
	suite.EqualError(err, "mock")
	suite.Nil(instance)
}

// FindInstances

func (suite *extCtrlUnitTestSuite) TestFindInstances_BeginTransactionFailure() {
	suite.dbMock.ExpectBegin().WillReturnError(fmt.Errorf("mock"))
	instances, err := suite.ctrl.FindInstances(mockContext(), suite.db, "extId", "extVer")
	suite.EqualError(err, "failed to begin transaction: mock")
	suite.Nil(instances)
}

func (suite *extCtrlUnitTestSuite) TestFindInstances_Success() {
	suite.dbMock.ExpectBegin()
	suite.mockCtrl.On("FindInstances", mock.Anything, "extId", "extVer").Return([]*extensionAPI.JsExtInstance{{Id: "instId", Name: "newInst"}}, nil)
	suite.dbMock.ExpectRollback()
	instances, err := suite.ctrl.FindInstances(mockContext(), suite.db, "extId", "extVer")
	suite.NoError(err)
	suite.Equal([]*extensionAPI.JsExtInstance{{Id: "instId", Name: "newInst"}}, instances)
}

func (suite *extCtrlUnitTestSuite) TestFindInstances_Failure() {
	suite.dbMock.ExpectBegin()
	suite.mockCtrl.On("FindInstances", mock.Anything, "extId", "extVer").Return(nil, fmt.Errorf("mock"))
	suite.dbMock.ExpectRollback()
	instances, err := suite.ctrl.FindInstances(mockContext(), suite.db, "extId", "extVer")
	suite.EqualError(err, "mock")
	suite.Nil(instances)
}

// DeleteInstance

func (suite *extCtrlUnitTestSuite) TestDeleteInstance_BeginTransactionFailure() {
	suite.dbMock.ExpectBegin().WillReturnError(fmt.Errorf("mock"))
	err := suite.ctrl.DeleteInstance(mockContext(), suite.db, "extId", "extVers", "instId")
	suite.EqualError(err, "failed to begin transaction: mock")
}

func (suite *extCtrlUnitTestSuite) TestDeleteInstance_Success() {
	suite.dbMock.ExpectBegin()
	suite.mockCtrl.On("DeleteInstance", mock.Anything, "extId", "extVers", "instId", mock.Anything).Return(nil)
	suite.dbMock.ExpectCommit()
	err := suite.ctrl.DeleteInstance(mockContext(), suite.db, "extId", "extVers", "instId")
	suite.NoError(err)
}

func (suite *extCtrlUnitTestSuite) TestDeleteInstance_Failure() {
	suite.dbMock.ExpectBegin()
	suite.mockCtrl.On("DeleteInstance", mock.Anything, "extId", "extVers", "instId", mock.Anything).Return(fmt.Errorf("mock"))
	suite.dbMock.ExpectRollback()
	err := suite.ctrl.DeleteInstance(mockContext(), suite.db, "extId", "extVers", "instId")
	suite.EqualError(err, "mock")
}

func (suite *extCtrlUnitTestSuite) TestDeleteInstance_CommitFailure() {
	suite.dbMock.ExpectBegin()
	suite.mockCtrl.On("DeleteInstance", mock.Anything, "extId", "extVers", "instId", mock.Anything).Return(nil)
	suite.dbMock.ExpectCommit().WillReturnError(fmt.Errorf("mock"))
	err := suite.ctrl.DeleteInstance(mockContext(), suite.db, "extId", "extVers", "instId")
	suite.EqualError(err, "mock")
}
