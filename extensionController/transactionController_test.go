package extensionController

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/exasol/extension-manager/extensionAPI"

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
		suite.Failf("there were unfulfilled expectations: %v", err.Error())
	}
}

func (suite *extCtrlUnitTestSuite) TestGetAllExtensions_Success() {
	suite.mockBfs.On("ListFiles", mock.Anything, mock.Anything, "default").Return([]BfsFile{}, nil)
	suite.mockCtrl.On("GetAllExtensions", mock.Anything).Return([]*Extension{}, nil)
	extensions, err := suite.ctrl.GetAllExtensions(mockContext(), suite.db)
	suite.NoError(err)
	suite.Len(extensions, 0)
}

func (suite *extCtrlUnitTestSuite) TestGetAllExtensions_BucketFsListFails() {
	suite.mockBfs.On("ListFiles", mock.Anything, mock.Anything, "default").Return(nil, fmt.Errorf("mock"))
	extensions, err := suite.ctrl.GetAllExtensions(mockContext(), suite.db)
	suite.EqualError(err, "failed to search for required files in BucketFS. Cause: mock")
	suite.Nil(extensions)
}

func (suite *extCtrlUnitTestSuite) TestGetAllExtensions_GetFails() {
	suite.mockBfs.On("ListFiles", mock.Anything, mock.Anything, "default").Return([]BfsFile{}, nil)
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

func (suite *extCtrlUnitTestSuite) TestGetAllInstallations_Success() {
	suite.dbMock.ExpectBegin()
	mockResult := []*extensionAPI.JsExtInstallation{{Name: "ext"}}
	suite.mockCtrl.On("GetAllInstallations", mock.Anything).Return(mockResult, nil)
	suite.dbMock.ExpectRollback()
	installations, err := suite.ctrl.GetInstalledExtensions(mockContext(), suite.db)
	suite.NoError(err)
	suite.Equal(mockResult, installations)
}

func (suite *extCtrlUnitTestSuite) TestGetAllInstallations_Failure() {
	suite.dbMock.ExpectBegin()
	suite.mockCtrl.On("GetAllInstallations", mock.Anything).Return(nil, fmt.Errorf("mock"))
	suite.dbMock.ExpectRollback()
	installations, err := suite.ctrl.GetInstalledExtensions(mockContext(), suite.db)
	suite.EqualError(err, "mock")
	suite.Nil(installations)
}

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
