package extensionController

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/exasol/extension-manager/extensionAPI"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type mockControllerImpl struct {
	mock.Mock
}

func (mock *mockControllerImpl) GetAllExtensions(bfsFiles []BfsFile) ([]*Extension, error) {
	args := mock.Called(bfsFiles)
	if ext, ok := args.Get(0).([]*Extension); ok {
		return ext, args.Error(1)
	} else {
		return nil, args.Error(1)
	}
}
func (mock *mockControllerImpl) GetAllInstallations(tx *sql.Tx) ([]*extensionAPI.JsExtInstallation, error) {
	args := mock.Called(tx)
	if result, ok := args.Get(0).([]*extensionAPI.JsExtInstallation); ok {
		return result, args.Error(1)
	} else {
		return nil, args.Error(1)
	}
}
func (mock *mockControllerImpl) InstallExtension(tx *sql.Tx, extensionId string, extensionVersion string) error {
	args := mock.Called(tx, extensionId, extensionVersion)
	return args.Error(0)
}
func (mock *mockControllerImpl) CreateInstance(tx *sql.Tx, extensionId string, extensionVersion string, parameterValues []ParameterValue) (string, error) {
	args := mock.Called(tx, extensionId, extensionVersion, parameterValues)
	return args.String(0), args.Error(1)
}

type mockBucketFs struct {
	mock.Mock
}

func (mock *mockBucketFs) ListBuckets(ctx context.Context, db *sql.DB) ([]string, error) {
	args := mock.Called(ctx, db)
	if buckets, ok := args.Get(0).([]string); ok {
		return buckets, args.Error(1)
	} else {
		return args.Get(0).([]string), args.Error(1)
	}
}
func (mock *mockBucketFs) ListFiles(ctx context.Context, db *sql.DB, bucket string) ([]BfsFile, error) {
	args := mock.Called(ctx, db, bucket)
	if files, ok := args.Get(0).([]BfsFile); ok {
		return files, args.Error(1)
	} else {
		return nil, args.Error(1)
	}
}

type extCtrlUnitTestSuite struct {
	suite.Suite
	ctrl     TransactionController
	db       *sql.DB
	dbMock   sqlmock.Sqlmock
	mockCtrl mockControllerImpl
	mockBfs  mockBucketFs
}

func TestExtensionControllerUnitTestSuite(t *testing.T) {
	suite.Run(t, new(extCtrlUnitTestSuite))
}

func (suite *extCtrlUnitTestSuite) SetupTest() {
	suite.mockCtrl = mockControllerImpl{}
	suite.mockBfs = mockBucketFs{}
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

func (suite *extCtrlUnitTestSuite) TestGetAllInstallations_Success() {
	suite.dbMock.ExpectBegin()
	mockResult := []*extensionAPI.JsExtInstallation{{Name: "ext"}}
	suite.mockCtrl.On("GetAllInstallations", mock.Anything).Return(mockResult, nil)
	suite.dbMock.ExpectRollback()
	installations, err := suite.ctrl.GetAllInstallations(mockContext(), suite.db)
	suite.NoError(err)
	suite.Equal(mockResult, installations)
}

func (suite *extCtrlUnitTestSuite) TestGetAllInstallations_Failure() {
	suite.dbMock.ExpectBegin()
	suite.mockCtrl.On("GetAllInstallations", mock.Anything).Return(nil, fmt.Errorf("mock"))
	suite.dbMock.ExpectRollback()
	installations, err := suite.ctrl.GetAllInstallations(mockContext(), suite.db)
	suite.EqualError(err, "mock")
	suite.Nil(installations)
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
